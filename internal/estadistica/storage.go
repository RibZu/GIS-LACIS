package estadistica

import (
	"database/sql"
	"fmt"
)

type Storage interface {
	NombresParticipantesExternos() ([]string, error)
	ResumenProductos() ([]ProductoResumen, error)
	ResumenIntegrantes() ([]IntegranteResumen, error)
	ResumenTesis() ([]TesisResumen, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) NombresParticipantesExternos() ([]string, error) {
	rows, err := s.db.Query(`SELECT nombre FROM participante_externo`)
	if err != nil {
		return nil, fmt.Errorf("error al consultar participantes externos: %w", err)
	}
	defer rows.Close()

	var nombres []string
	for rows.Next() {
		var nombre string
		if err := rows.Scan(&nombre); err != nil {
			return nil, fmt.Errorf("error al escanear participante externo: %w", err)
		}
		nombres = append(nombres, nombre)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return nombres, nil
}

func (s *PostgresStorage) ResumenProductos() ([]ProductoResumen, error) {
	query := `SELECT d.anio, COALESCE(d.url, '')
	          FROM desarrollo d
	          ORDER BY d.anio DESC, d.id DESC`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar el resumen de productos: %w", err)
	}
	defer rows.Close()

	var productos []ProductoResumen
	for rows.Next() {
		var p ProductoResumen
		if err := rows.Scan(&p.Anio, &p.URL); err != nil {
			return nil, fmt.Errorf("error al escanear el resumen de un producto: %w", err)
		}
		productos = append(productos, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return productos, nil
}

func (s *PostgresStorage) ResumenIntegrantes() ([]IntegranteResumen, error) {
	query := `SELECT COALESCE(activo, true), COALESCE(pertenece_lacis, false), COALESCE(pertenece_grupo_software, false)
	          FROM integrante`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar el resumen de integrantes: %w", err)
	}
	defer rows.Close()

	var integrantes []IntegranteResumen
	for rows.Next() {
		var i IntegranteResumen
		if err := rows.Scan(&i.Activo, &i.PerteneceLacis, &i.PerteneceGrupoSoftware); err != nil {
			return nil, fmt.Errorf("error al escanear el resumen de un integrante: %w", err)
		}
		integrantes = append(integrantes, i)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return integrantes, nil
}

func (s *PostgresStorage) ResumenTesis() ([]TesisResumen, error) {
	query := `SELECT anio, COALESCE(nivel, ''), COALESCE(TRIM(archivo_pdf), '') <> ''
	          FROM tesis`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar el resumen de tesis: %w", err)
	}
	defer rows.Close()

	var tesis []TesisResumen
	for rows.Next() {
		var t TesisResumen
		var anio sql.NullInt64
		if err := rows.Scan(&anio, &t.Nivel, &t.TienePDF); err != nil {
			return nil, fmt.Errorf("error al escanear el resumen de una tesis: %w", err)
		}
		if anio.Valid {
			a := int(anio.Int64)
			t.Anio = &a
		}
		tesis = append(tesis, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tesis, nil
}
