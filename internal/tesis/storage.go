package tesis

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound = errors.New("tesis no encontrada")
	ErrEmptyID  = errors.New("ID de tesis vacío o inválido")
)

type Storage interface {
	Create(t *Tesis) error
	Read(id int) (*Tesis, error)
	Delete(id int) error
	Update(id int, fields UpdateFields) error
	GetAll() ([]Tesis, error)
	GetByCarrera(carrera string) ([]Tesis, error)
	GetByNivel(nivel string) ([]Tesis, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}

func (s *PostgresStorage) Create(t *Tesis) error {
	query := `
		INSERT INTO tesis (
			titulo, 
			anio, 
			nivel, 
			carrera_origen, 
			palabras_clave, 
			resumen, 
			archivo_pdf, 
			autor_id, 
			autor_historico, 
			director_id, 
			director_historico, 
			coodirector_id, 
			coodirector_historico, 
			proyecto_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		t.Titulo,
		t.Anio,
		t.Nivel,
		t.CarreraOrigen,
		t.PalabrasClave,
		t.Resumen,
		t.ArchivoPDF,
		t.AutorID,
		t.AutorHistorico,
		t.DirectorID,
		t.DirectorHistorico,
		t.CoodirectorID,
		t.CoodirectorHistorico,
		t.ProyectoID,
	).Scan(&t.ID)

	if err != nil {
		return fmt.Errorf("error al insertar tesis en PostgreSQL: %w", err)
	}

	// Insertar integrantes secundarios en integrantes_tesis si existen
	if len(t.IntegrantesIDs) > 0 {
		for _, intID := range t.IntegrantesIDs {
			if intID > 0 {
				_, _ = s.db.Exec(`INSERT INTO integrantes_tesis (tesis_id, integrante_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, t.ID, intID)
			}
		}
	}

	return nil
}

func (s *PostgresStorage) Read(id int) (*Tesis, error) {
	query := `
		SELECT 
			t.id, 
			t.titulo, 
			t.anio, 
			COALESCE(t.nivel, ''), 
			COALESCE(t.carrera_origen, ''), 
			COALESCE(t.palabras_clave, ''), 
			COALESCE(t.resumen, ''), 
			COALESCE(t.archivo_pdf, ''), 
			t.autor_id, 
			COALESCE(t.autor_historico, ''), 
			t.director_id, 
			COALESCE(t.director_historico, ''), 
			t.coodirector_id, 
			COALESCE(t.coodirector_historico, ''), 
			t.proyecto_id,
			COALESCE(ia.nombre || ' ' || ia.apellido, ''),
			COALESCE(id_dir.nombre || ' ' || id_dir.apellido, ''),
			COALESCE(ic.nombre || ' ' || ic.apellido, '')
		FROM tesis t
		LEFT JOIN integrante ia ON t.autor_id = ia.id
		LEFT JOIN integrante id_dir ON t.director_id = id_dir.id
		LEFT JOIN integrante ic ON t.coodirector_id = ic.id
		WHERE t.id = $1`

	var t Tesis
	var anioVal sql.NullInt64
	var autorID, dirID, coodirID, proyID sql.NullInt64
	var autorNom, dirNom, coodirNom string

	row := s.db.QueryRow(query, id)
	err := row.Scan(
		&t.ID,
		&t.Titulo,
		&anioVal,
		&t.Nivel,
		&t.CarreraOrigen,
		&t.PalabrasClave,
		&t.Resumen,
		&t.ArchivoPDF,
		&autorID,
		&t.AutorHistorico,
		&dirID,
		&t.DirectorHistorico,
		&coodirID,
		&t.CoodirectorHistorico,
		&proyID,
		&autorNom,
		&dirNom,
		&coodirNom,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error al leer tesis en PostgreSQL: %w", err)
	}

	if anioVal.Valid {
		val := int(anioVal.Int64)
		t.Anio = &val
	}
	if autorID.Valid {
		val := int(autorID.Int64)
		t.AutorID = &val
		t.AutorNombre = autorNom
	}
	if dirID.Valid {
		val := int(dirID.Int64)
		t.DirectorID = &val
		t.DirectorNombre = dirNom
	}
	if coodirID.Valid {
		val := int(coodirID.Int64)
		t.CoodirectorID = &val
		t.CoodirectorNombre = coodirNom
	}
	if proyID.Valid {
		val := int(proyID.Int64)
		t.ProyectoID = &val
	}

	// Cargar integrantes_tesis vinculados
	intRows, err := s.db.Query(`SELECT integrante_id FROM integrantes_tesis WHERE tesis_id = $1`, id)
	if err == nil {
		defer intRows.Close()
		for intRows.Next() {
			var iID int
			if errScan := intRows.Scan(&iID); errScan == nil {
				t.IntegrantesIDs = append(t.IntegrantesIDs, iID)
			}
		}
	}

	return &t, nil
}

func (s *PostgresStorage) GetAll() ([]Tesis, error) {
	query := `
		SELECT 
			t.id, 
			t.titulo, 
			t.anio, 
			COALESCE(t.nivel, ''), 
			COALESCE(t.carrera_origen, ''), 
			COALESCE(t.palabras_clave, ''), 
			COALESCE(t.resumen, ''), 
			COALESCE(t.archivo_pdf, ''), 
			t.autor_id, 
			COALESCE(t.autor_historico, ''), 
			t.director_id, 
			COALESCE(t.director_historico, ''), 
			t.coodirector_id, 
			COALESCE(t.coodirector_historico, ''), 
			t.proyecto_id,
			COALESCE(ia.nombre || ' ' || ia.apellido, ''),
			COALESCE(id_dir.nombre || ' ' || id_dir.apellido, ''),
			COALESCE(ic.nombre || ' ' || ic.apellido, '')
		FROM tesis t
		LEFT JOIN integrante ia ON t.autor_id = ia.id
		LEFT JOIN integrante id_dir ON t.director_id = id_dir.id
		LEFT JOIN integrante ic ON t.coodirector_id = ic.id
		ORDER BY t.anio DESC NULLS LAST, t.id DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar tesis en PostgreSQL: %w", err)
	}
	defer rows.Close()

	var lista []Tesis
	for rows.Next() {
		var t Tesis
		var anioVal sql.NullInt64
		var autorID, dirID, coodirID, proyID sql.NullInt64
		var autorNom, dirNom, coodirNom string

		errScan := rows.Scan(
			&t.ID,
			&t.Titulo,
			&anioVal,
			&t.Nivel,
			&t.CarreraOrigen,
			&t.PalabrasClave,
			&t.Resumen,
			&t.ArchivoPDF,
			&autorID,
			&t.AutorHistorico,
			&dirID,
			&t.DirectorHistorico,
			&coodirID,
			&t.CoodirectorHistorico,
			&proyID,
			&autorNom,
			&dirNom,
			&coodirNom,
		)
		if errScan != nil {
			return nil, fmt.Errorf("error al escanear tesis: %w", errScan)
		}

		if anioVal.Valid {
			val := int(anioVal.Int64)
			t.Anio = &val
		}
		if autorID.Valid {
			val := int(autorID.Int64)
			t.AutorID = &val
			t.AutorNombre = autorNom
		}
		if dirID.Valid {
			val := int(dirID.Int64)
			t.DirectorID = &val
			t.DirectorNombre = dirNom
		}
		if coodirID.Valid {
			val := int(coodirID.Int64)
			t.CoodirectorID = &val
			t.CoodirectorNombre = coodirNom
		}
		if proyID.Valid {
			val := int(proyID.Int64)
			t.ProyectoID = &val
		}

		lista = append(lista, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return lista, nil
}

func (s *PostgresStorage) GetByCarrera(carrera string) ([]Tesis, error) {
	query := `
		SELECT 
			t.id, 
			t.titulo, 
			t.anio, 
			COALESCE(t.nivel, ''), 
			COALESCE(t.carrera_origen, ''), 
			COALESCE(t.palabras_clave, ''), 
			COALESCE(t.resumen, ''), 
			COALESCE(t.archivo_pdf, ''), 
			t.autor_id, 
			COALESCE(t.autor_historico, ''), 
			t.director_id, 
			COALESCE(t.director_historico, ''), 
			t.coodirector_id, 
			COALESCE(t.coodirector_historico, ''), 
			t.proyecto_id,
			COALESCE(ia.nombre || ' ' || ia.apellido, ''),
			COALESCE(id_dir.nombre || ' ' || id_dir.apellido, ''),
			COALESCE(ic.nombre || ' ' || ic.apellido, '')
		FROM tesis t
		LEFT JOIN integrante ia ON t.autor_id = ia.id
		LEFT JOIN integrante id_dir ON t.director_id = id_dir.id
		LEFT JOIN integrante ic ON t.coodirector_id = ic.id
		WHERE LOWER(t.carrera_origen) = LOWER($1)
		ORDER BY t.anio DESC NULLS LAST, t.id DESC`

	rows, err := s.db.Query(query, carrera)
	if err != nil {
		return nil, fmt.Errorf("error al consultar tesis por carrera en PostgreSQL: %w", err)
	}
	defer rows.Close()

	var lista []Tesis
	for rows.Next() {
		var t Tesis
		var anioVal sql.NullInt64
		var autorID, dirID, coodirID, proyID sql.NullInt64
		var autorNom, dirNom, coodirNom string

		errScan := rows.Scan(
			&t.ID,
			&t.Titulo,
			&anioVal,
			&t.Nivel,
			&t.CarreraOrigen,
			&t.PalabrasClave,
			&t.Resumen,
			&t.ArchivoPDF,
			&autorID,
			&t.AutorHistorico,
			&dirID,
			&t.DirectorHistorico,
			&coodirID,
			&t.CoodirectorHistorico,
			&proyID,
			&autorNom,
			&dirNom,
			&coodirNom,
		)
		if errScan != nil {
			return nil, fmt.Errorf("error al escanear tesis: %w", errScan)
		}

		if anioVal.Valid {
			val := int(anioVal.Int64)
			t.Anio = &val
		}
		if autorID.Valid {
			val := int(autorID.Int64)
			t.AutorID = &val
			t.AutorNombre = autorNom
		}
		if dirID.Valid {
			val := int(dirID.Int64)
			t.DirectorID = &val
			t.DirectorNombre = dirNom
		}
		if coodirID.Valid {
			val := int(coodirID.Int64)
			t.CoodirectorID = &val
			t.CoodirectorNombre = coodirNom
		}
		if proyID.Valid {
			val := int(proyID.Int64)
			t.ProyectoID = &val
		}

		lista = append(lista, t)
	}

	return lista, nil
}

func (s *PostgresStorage) GetByNivel(nivel string) ([]Tesis, error) {
	query := `
		SELECT 
			t.id, 
			t.titulo, 
			t.anio, 
			COALESCE(t.nivel, ''), 
			COALESCE(t.carrera_origen, ''), 
			COALESCE(t.palabras_clave, ''), 
			COALESCE(t.resumen, ''), 
			COALESCE(t.archivo_pdf, ''), 
			t.autor_id, 
			COALESCE(t.autor_historico, ''), 
			t.director_id, 
			COALESCE(t.director_historico, ''), 
			t.coodirector_id, 
			COALESCE(t.coodirector_historico, ''), 
			t.proyecto_id,
			COALESCE(ia.nombre || ' ' || ia.apellido, ''),
			COALESCE(id_dir.nombre || ' ' || id_dir.apellido, ''),
			COALESCE(ic.nombre || ' ' || ic.apellido, '')
		FROM tesis t
		LEFT JOIN integrante ia ON t.autor_id = ia.id
		LEFT JOIN integrante id_dir ON t.director_id = id_dir.id
		LEFT JOIN integrante ic ON t.coodirector_id = ic.id
		WHERE LOWER(t.nivel) = LOWER($1)
		ORDER BY t.anio DESC NULLS LAST, t.id DESC`

	rows, err := s.db.Query(query, nivel)
	if err != nil {
		return nil, fmt.Errorf("error al consultar tesis por nivel en PostgreSQL: %w", err)
	}
	defer rows.Close()

	var lista []Tesis
	for rows.Next() {
		var t Tesis
		var anioVal sql.NullInt64
		var autorID, dirID, coodirID, proyID sql.NullInt64
		var autorNom, dirNom, coodirNom string

		errScan := rows.Scan(
			&t.ID,
			&t.Titulo,
			&anioVal,
			&t.Nivel,
			&t.CarreraOrigen,
			&t.PalabrasClave,
			&t.Resumen,
			&t.ArchivoPDF,
			&autorID,
			&t.AutorHistorico,
			&dirID,
			&t.DirectorHistorico,
			&coodirID,
			&t.CoodirectorHistorico,
			&proyID,
			&autorNom,
			&dirNom,
			&coodirNom,
		)
		if errScan != nil {
			return nil, fmt.Errorf("error al escanear tesis: %w", errScan)
		}

		if anioVal.Valid {
			val := int(anioVal.Int64)
			t.Anio = &val
		}
		if autorID.Valid {
			val := int(autorID.Int64)
			t.AutorID = &val
			t.AutorNombre = autorNom
		}
		if dirID.Valid {
			val := int(dirID.Int64)
			t.DirectorID = &val
			t.DirectorNombre = dirNom
		}
		if coodirID.Valid {
			val := int(coodirID.Int64)
			t.CoodirectorID = &val
			t.CoodirectorNombre = coodirNom
		}
		if proyID.Valid {
			val := int(proyID.Int64)
			t.ProyectoID = &val
		}

		lista = append(lista, t)
	}

	return lista, nil
}

func (s *PostgresStorage) Delete(id int) error {
	query := `DELETE FROM tesis WHERE id = $1`
	res, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar tesis en PostgreSQL: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostgresStorage) Update(id int, fields UpdateFields) error {
	var query string = "UPDATE tesis SET "
	var args []interface{}
	var argID int = 1

	if fields.Titulo != nil {
		query += fmt.Sprintf("titulo = $%d, ", argID)
		args = append(args, *fields.Titulo)
		argID++
	}
	if fields.Anio != nil {
		if *fields.Anio > 0 {
			query += fmt.Sprintf("anio = $%d, ", argID)
			args = append(args, *fields.Anio)
			argID++
		} else {
			query += "anio = NULL, "
		}
	}
	if fields.Nivel != nil {
		query += fmt.Sprintf("nivel = $%d, ", argID)
		args = append(args, *fields.Nivel)
		argID++
	}
	if fields.CarreraOrigen != nil {
		query += fmt.Sprintf("carrera_origen = $%d, ", argID)
		args = append(args, *fields.CarreraOrigen)
		argID++
	}
	if fields.PalabrasClave != nil {
		query += fmt.Sprintf("palabras_clave = $%d, ", argID)
		args = append(args, *fields.PalabrasClave)
		argID++
	}
	if fields.Resumen != nil {
		query += fmt.Sprintf("resumen = $%d, ", argID)
		args = append(args, *fields.Resumen)
		argID++
	}
	if fields.ArchivoPDF != nil {
		query += fmt.Sprintf("archivo_pdf = $%d, ", argID)
		args = append(args, *fields.ArchivoPDF)
		argID++
	}
	if fields.AutorID != nil {
		if *fields.AutorID > 0 {
			query += fmt.Sprintf("autor_id = $%d, ", argID)
			args = append(args, *fields.AutorID)
			argID++
		} else {
			query += "autor_id = NULL, "
		}
	}
	if fields.AutorHistorico != nil {
		query += fmt.Sprintf("autor_historico = $%d, ", argID)
		args = append(args, *fields.AutorHistorico)
		argID++
	}
	if fields.DirectorID != nil {
		if *fields.DirectorID > 0 {
			query += fmt.Sprintf("director_id = $%d, ", argID)
			args = append(args, *fields.DirectorID)
			argID++
		} else {
			query += "director_id = NULL, "
		}
	}
	if fields.DirectorHistorico != nil {
		query += fmt.Sprintf("director_historico = $%d, ", argID)
		args = append(args, *fields.DirectorHistorico)
		argID++
	}
	if fields.CoodirectorID != nil {
		if *fields.CoodirectorID > 0 {
			query += fmt.Sprintf("coodirector_id = $%d, ", argID)
			args = append(args, *fields.CoodirectorID)
			argID++
		} else {
			query += "coodirector_id = NULL, "
		}
	}
	if fields.CoodirectorHistorico != nil {
		query += fmt.Sprintf("coodirector_historico = $%d, ", argID)
		args = append(args, *fields.CoodirectorHistorico)
		argID++
	}
	if fields.ProyectoID != nil {
		if *fields.ProyectoID > 0 {
			query += fmt.Sprintf("proyecto_id = $%d, ", argID)
			args = append(args, *fields.ProyectoID)
			argID++
		} else {
			query += "proyecto_id = NULL, "
		}
	}

	if len(args) == 0 && !strings.Contains(query, "NULL") {
		return nil
	}

	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", argID)
	args = append(args, id)

	res, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error al actualizar tesis en PostgreSQL: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	// Si se enviaron IntegrantesIDs, sincronizar tabla intermedia
	if fields.IntegrantesIDs != nil {
		_, _ = s.db.Exec(`DELETE FROM integrantes_tesis WHERE tesis_id = $1`, id)
		for _, intID := range fields.IntegrantesIDs {
			if intID > 0 {
				_, _ = s.db.Exec(`INSERT INTO integrantes_tesis (tesis_id, integrante_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, id, intID)
			}
		}
	}

	return nil
}
