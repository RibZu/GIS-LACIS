package integrante

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

type jsonIntegranteRaw struct {
	ClaseRol        string `json:"clase_rol"`
	Rol             string `json:"rol"`
	Nombre          string `json:"nombre"`
	Imagen          string `json:"imagen"`
	CV              string `json:"cv"`
	Mail            string `json:"mail"`
	Linkedin        string `json:"linkedin"`
	Especializacion string `json:"especializacion"`
	Descripcion     string `json:"descripcion"`
	Grupo           string `json:"grupo"`
}

type memberAgg struct {
	RawName         string
	PerteneceGIS    bool
	PerteneceLacis  bool
	RolSoftwareID   *int
	RolLacisID      *int
	Especializacion string
	Descripcion     string
	Mail            string
	Linkedin        string
	Imagen          string
	CV              string
}

func mapRoleID(item jsonIntegranteRaw, defaultRole *int) int {
	if defaultRole != nil {
		return *defaultRole
	}
	rol := strings.ToUpper(item.Rol)
	clase := strings.ToLower(item.ClaseRol)
	if strings.Contains(rol, "DIRECTOR") || strings.Contains(clase, "director") {
		return 1
	}
	if strings.Contains(rol, "INVESTIGADOR") || strings.Contains(clase, "investigador") {
		return 2
	}
	if strings.Contains(rol, "ASESOR") || strings.Contains(clase, "asesor") {
		return 3
	}
	if strings.Contains(rol, "ESTUDIANTE") || strings.Contains(rol, "BECARIO") ||
		strings.Contains(clase, "estudiante") || strings.Contains(clase, "becario") {
		return 4
	}
	return 2
}

func MigrarIntegrantesDesdeJSON(db *sql.DB, baseDir string, logger *zap.Logger) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM integrante").Scan(&count)
	if err != nil {
		return fmt.Errorf("error al verificar cantidad de integrantes: %w", err)
	}

	if count > 0 {
		logger.Info("La tabla integrante ya contiene registros, no se requiere migración", zap.Int("total", count))
		return nil
	}

	logger.Info("Iniciando migración de integrantes desde JSON hacia PostgreSQL...")

	type fileConfig struct {
		Filename      string
		PerteneceGIS  bool
		PerteneceLaci bool
		RolGIS        *int
		RolLacis      *int
	}

	rolDirector := 1
	configs := []fileConfig{
		{Filename: "director.json", PerteneceGIS: true, PerteneceLaci: false, RolGIS: &rolDirector, RolLacis: nil},
		{Filename: "directores.json", PerteneceGIS: true, PerteneceLaci: false, RolGIS: &rolDirector, RolLacis: nil},
		{Filename: "directorLacis.json", PerteneceGIS: false, PerteneceLaci: true, RolGIS: nil, RolLacis: &rolDirector},
		{Filename: "directoresLacis.json", PerteneceGIS: false, PerteneceLaci: true, RolGIS: nil, RolLacis: &rolDirector},
		{Filename: "integrantes.json", PerteneceGIS: true, PerteneceLaci: false, RolGIS: nil, RolLacis: nil},
		{Filename: "integrantesLacis.json", PerteneceGIS: false, PerteneceLaci: true, RolGIS: nil, RolLacis: nil},
	}

	membersMap := make(map[string]*memberAgg)
	var orderedKeys []string

	for _, cfg := range configs {
		jsonPath := filepath.Join(baseDir, cfg.Filename)
		data, err := os.ReadFile(jsonPath)
		if err != nil {
			logger.Warn("No se pudo leer archivo JSON para migración", zap.String("path", jsonPath), zap.Error(err))
			continue
		}

		var items []jsonIntegranteRaw
		if err := json.Unmarshal(data, &items); err != nil {
			logger.Error("Error al deserializar JSON", zap.String("path", jsonPath), zap.Error(err))
			continue
		}

		for _, item := range items {
			rawName := strings.TrimSpace(item.Nombre)
			normName := strings.Join(strings.Fields(rawName), " ")
			if normName == "" {
				continue
			}

			m, exists := membersMap[normName]
			if !exists {
				m = &memberAgg{
					RawName: rawName,
				}
				membersMap[normName] = m
				orderedKeys = append(orderedKeys, normName)
			}

			if cfg.PerteneceGIS {
				m.PerteneceGIS = true
				rID := mapRoleID(item, cfg.RolGIS)
				if m.RolSoftwareID == nil || rID == 1 {
					m.RolSoftwareID = &rID
				}
			}

			if cfg.PerteneceLaci {
				m.PerteneceLacis = true
				rID := mapRoleID(item, cfg.RolLacis)
				if m.RolLacisID == nil || rID == 1 {
					m.RolLacisID = &rID
				}
			}

			if m.Especializacion == "" && item.Especializacion != "" {
				m.Especializacion = strings.TrimSpace(item.Especializacion)
			}
			if m.Descripcion == "" && item.Descripcion != "" {
				m.Descripcion = strings.TrimSpace(item.Descripcion)
			}
			if m.Mail == "" && item.Mail != "" {
				m.Mail = strings.TrimSpace(item.Mail)
			}
			if m.Linkedin == "" && item.Linkedin != "" {
				m.Linkedin = strings.TrimSpace(item.Linkedin)
			}
			if m.Imagen == "" && item.Imagen != "" {
				m.Imagen = strings.TrimSpace(item.Imagen)
			}
			if m.CV == "" && item.CV != "" {
				m.CV = strings.TrimSpace(item.CV)
			}
		}
	}

	prefixes := []string{"Dr.Sci", "Dr.", "Mg.", "Esp.", "Lic.", "Ing.", "Est."}
	compoundSurnames := map[string][2]string{
		"Novillo Rangone Gabriel":        {"Novillo Rangone", "Gabriel"},
		"Olguin Quevedo Gonzalo Gabriel": {"Olguin Quevedo", "Gonzalo Gabriel"},
		"Riera Bauer Fabricio José":      {"Riera Bauer", "Fabricio José"},
		"Orellano Zalazar Ignacio":       {"Orellano Zalazar", "Ignacio"},
		"Simon Martin Riberi Zunino":     {"Riberi Zunino", "Simon Martin"},
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	insertQuery := `INSERT INTO integrante (
		nombre, 
		apellido, 
		titulo_especializacion, 
		descripcion, 
		contacto_mail, 
		contacto_linkedin, 
		imagen_url, 
		cv_url, 
		pertenece_lacis, 
		pertenece_grupo_software, 
		activo, 
		rol_id,
		rol_lacis_id,
		rol_software_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	insertedCount := 0
	for _, normName := range orderedKeys {
		m := membersMap[normName]

		clean := normName
		foundPref := ""
		for _, p := range prefixes {
			if strings.HasPrefix(clean, p) {
				foundPref = p
				clean = strings.TrimSpace(strings.TrimPrefix(clean, p))
				break
			}
		}

		var surname, given string
		if val, ok := compoundSurnames[clean]; ok {
			surname = val[0]
			given = val[1]
		} else {
			parts := strings.Fields(clean)
			if len(parts) > 1 {
				surname = parts[0]
				given = strings.Join(parts[1:], " ")
			} else if len(parts) == 1 {
				surname = parts[0]
				given = ""
			}
		}

		nombreCol := strings.TrimSpace(foundPref + " " + surname)
		apellidoCol := strings.TrimSpace(given)
		if apellidoCol == "" {
			apellidoCol = surname
		}

		rolID := 2
		if m.RolSoftwareID != nil {
			rolID = *m.RolSoftwareID
		} else if m.RolLacisID != nil {
			rolID = *m.RolLacisID
		}

		_, err := tx.Exec(
			insertQuery,
			nombreCol,
			apellidoCol,
			m.Especializacion,
			m.Descripcion,
			m.Mail,
			m.Linkedin,
			m.Imagen,
			m.CV,
			m.PerteneceLacis,
			m.PerteneceGIS,
			true,
			rolID,
			m.RolLacisID,
			m.RolSoftwareID,
		)
		if err != nil {
			return fmt.Errorf("error al insertar integrante '%s': %w", normName, err)
		}
		insertedCount++
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar transacción de migración: %w", err)
	}

	logger.Info("Migración de integrantes finalizada exitosamente", zap.Int("total_insertados", insertedCount))
	return nil
}
