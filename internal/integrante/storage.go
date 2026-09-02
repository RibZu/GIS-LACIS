package integrante

import (
	"PaginaSEG/internal/rol"
	"database/sql"
	"errors"
	"fmt"
)

// ErrNotFound devuelve el error cuando no se encuentra el integrante con el ID proporcionado
var ErrNotFound = errors.New("integrante no encontrado")

// ErrEmptyID retorna error cuando el ID es inválido
var ErrEmptyID = errors.New("ID de integrante vacío o inválido")

type Storage interface {
	Create(integrante *Integrante) error
	Read(id int) (*Integrante, error)
	Delete(id int) error
	Update(id int, fields UpdateFields) error
	GetAll() ([]Integrante, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}

func (c *PostgresStorage) Create(integrante *Integrante) error {
	// Mantener rol_id con el primer rol disponible para retrocompatibilidad
	if integrante.RolID <= 0 {
		if integrante.RolLacisID != nil && *integrante.RolLacisID > 0 {
			integrante.RolID = *integrante.RolLacisID
		} else if integrante.RolSoftwareID != nil && *integrante.RolSoftwareID > 0 {
			integrante.RolID = *integrante.RolSoftwareID
		}
	}

	query := `INSERT INTO integrante (
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
			  )
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	          RETURNING id`

	err := c.db.QueryRow(
		query,
		integrante.Nombre,
		integrante.Apellido,
		integrante.Especializacion,
		integrante.Descripcion,
		integrante.Contacto,
		integrante.ContactoLinkedin,
		integrante.Imagen,
		integrante.CV,
		integrante.PerteneceLacis,
		integrante.PerteneceGrupoSoftware,
		integrante.Activo,
		integrante.RolID,
		integrante.RolLacisID,
		integrante.RolSoftwareID,
	).Scan(&integrante.ID)

	if err != nil {
		return fmt.Errorf("error al insertar integrante en PostgreSQL: %w", err)
	}

	return nil
}

func (s *PostgresStorage) Read(id int) (*Integrante, error) {
	query := `
		SELECT i.id, 
		       i.nombre, 
		       i.apellido, 
		       COALESCE(i.titulo_especializacion, ''), 
		       COALESCE(i.descripcion, ''), 
		       COALESCE(i.contacto_mail, ''), 
		       COALESCE(i.contacto_linkedin, ''), 
		       COALESCE(i.imagen_url, ''), 
		       COALESCE(i.cv_url, ''), 
		       COALESCE(i.pertenece_lacis, false), 
		       COALESCE(i.pertenece_grupo_software, false), 
		       COALESCE(i.activo, true), 
		       COALESCE(i.rol_id, 0), 
		       COALESCE(r.nombre, ''),
		       i.rol_lacis_id,
		       COALESCE(r_lacis.nombre, ''),
		       i.rol_software_id,
		       COALESCE(r_soft.nombre, '')
		FROM integrante i
		LEFT JOIN rol r ON i.rol_id = r.id
		LEFT JOIN rol r_lacis ON i.rol_lacis_id = r_lacis.id
		LEFT JOIN rol r_soft ON i.rol_software_id = r_soft.id
		WHERE i.id = $1`
	row := s.db.QueryRow(query, id)
	var u Integrante
	var rolLacisID, rolSoftwareID sql.NullInt64
	var rolLacisNombre, rolSoftwareNombre string

	err := row.Scan(
		&u.ID,
		&u.Nombre,
		&u.Apellido,
		&u.Especializacion,
		&u.Descripcion,
		&u.Contacto,
		&u.ContactoLinkedin,
		&u.Imagen,
		&u.CV,
		&u.PerteneceLacis,
		&u.PerteneceGrupoSoftware,
		&u.Activo,
		&u.RolID,
		&u.Rol.Nombre,
		&rolLacisID,
		&rolLacisNombre,
		&rolSoftwareID,
		&rolSoftwareNombre,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error al leer integrante en PostgreSQL: %w", err)
	}
	u.Rol.ID = u.RolID
	if rolLacisID.Valid {
		idVal := int(rolLacisID.Int64)
		u.RolLacisID = &idVal
		u.RolLacis = &rol.Rol{
			ID:     idVal,
			Nombre: rolLacisNombre,
		}
	}
	if rolSoftwareID.Valid {
		idVal := int(rolSoftwareID.Int64)
		u.RolSoftwareID = &idVal
		u.RolSoftware = &rol.Rol{
			ID:     idVal,
			Nombre: rolSoftwareNombre,
		}
	}
	return &u, nil
}

func (s *PostgresStorage) GetAll() ([]Integrante, error) {
	query := `
		SELECT i.id, 
		       i.nombre, 
		       i.apellido, 
		       COALESCE(i.titulo_especializacion, ''), 
		       COALESCE(i.descripcion, ''), 
		       COALESCE(i.contacto_mail, ''), 
		       COALESCE(i.contacto_linkedin, ''), 
		       COALESCE(i.imagen_url, ''), 
		       COALESCE(i.cv_url, ''), 
		       COALESCE(i.pertenece_lacis, false), 
		       COALESCE(i.pertenece_grupo_software, false), 
		       COALESCE(i.activo, true), 
		       COALESCE(i.rol_id, 0), 
		       COALESCE(r.nombre, ''),
		       i.rol_lacis_id,
		       COALESCE(r_lacis.nombre, ''),
		       i.rol_software_id,
		       COALESCE(r_soft.nombre, '')
		FROM integrante i
		LEFT JOIN rol r ON i.rol_id = r.id
		LEFT JOIN rol r_lacis ON i.rol_lacis_id = r_lacis.id
		LEFT JOIN rol r_soft ON i.rol_software_id = r_soft.id
		ORDER BY i.id ASC`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar integrantes en PostgreSQL: %w", err)
	}
	defer rows.Close()
	var integrantes []Integrante
	for rows.Next() {
		var u Integrante
		var rolLacisID, rolSoftwareID sql.NullInt64
		var rolLacisNombre, rolSoftwareNombre string
		err := rows.Scan(
			&u.ID,
			&u.Nombre,
			&u.Apellido,
			&u.Especializacion,
			&u.Descripcion,
			&u.Contacto,
			&u.ContactoLinkedin,
			&u.Imagen,
			&u.CV,
			&u.PerteneceLacis,
			&u.PerteneceGrupoSoftware,
			&u.Activo,
			&u.RolID,
			&u.Rol.Nombre,
			&rolLacisID,
			&rolLacisNombre,
			&rolSoftwareID,
			&rolSoftwareNombre,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear integrante: %w", err)
		}
		u.Rol.ID = u.RolID
		if rolLacisID.Valid {
			idVal := int(rolLacisID.Int64)
			u.RolLacisID = &idVal
			u.RolLacis = &rol.Rol{
				ID:     idVal,
				Nombre: rolLacisNombre,
			}
		}
		if rolSoftwareID.Valid {
			idVal := int(rolSoftwareID.Int64)
			u.RolSoftwareID = &idVal
			u.RolSoftware = &rol.Rol{
				ID:     idVal,
				Nombre: rolSoftwareNombre,
			}
		}
		integrantes = append(integrantes, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return integrantes, nil
}

func (s *PostgresStorage) Delete(id int) error {
	query := `DELETE FROM integrante WHERE id = $1`
	res, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar integrante en PostgreSQL: %w", err)
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
	var query string = "UPDATE integrante SET "
	var args []interface{}
	var argID int = 1

	if fields.Nombre != nil {
		query += fmt.Sprintf("nombre = $%d, ", argID)
		args = append(args, *fields.Nombre)
		argID++
	}
	if fields.Apellido != nil {
		query += fmt.Sprintf("apellido = $%d, ", argID)
		args = append(args, *fields.Apellido)
		argID++
	}
	if fields.Especializacion != nil {
		query += fmt.Sprintf("titulo_especializacion = $%d, ", argID)
		args = append(args, *fields.Especializacion)
		argID++
	}
	if fields.Descripcion != nil {
		query += fmt.Sprintf("descripcion = $%d, ", argID)
		args = append(args, *fields.Descripcion)
		argID++
	}
	if fields.Contacto != nil {
		query += fmt.Sprintf("contacto_mail = $%d, ", argID)
		args = append(args, *fields.Contacto)
		argID++
	}
	if fields.ContactoLinkedin != nil {
		query += fmt.Sprintf("contacto_linkedin = $%d, ", argID)
		args = append(args, *fields.ContactoLinkedin)
		argID++
	}
	if fields.Imagen != nil {
		query += fmt.Sprintf("imagen_url = $%d, ", argID)
		args = append(args, *fields.Imagen)
		argID++
	}
	if fields.CV != nil {
		query += fmt.Sprintf("cv_url = $%d, ", argID)
		args = append(args, *fields.CV)
		argID++
	}
	if fields.PerteneceLacis != nil {
		query += fmt.Sprintf("pertenece_lacis = $%d, ", argID)
		args = append(args, *fields.PerteneceLacis)
		argID++
	}
	if fields.PerteneceGrupoSoftware != nil {
		query += fmt.Sprintf("pertenece_grupo_software = $%d, ", argID)
		args = append(args, *fields.PerteneceGrupoSoftware)
		argID++
	}
	if fields.Activo != nil {
		query += fmt.Sprintf("activo = $%d, ", argID)
		args = append(args, *fields.Activo)
		argID++
	}
	if fields.RolID != nil {
		query += fmt.Sprintf("rol_id = $%d, ", argID)
		args = append(args, *fields.RolID)
		argID++
	}
	if fields.RolLacisID != nil {
		if *fields.RolLacisID > 0 {
			query += fmt.Sprintf("rol_lacis_id = $%d, ", argID)
			args = append(args, *fields.RolLacisID)
			argID++
		} else {
			query += "rol_lacis_id = NULL, "
		}
	}
	if fields.RolSoftwareID != nil {
		if *fields.RolSoftwareID > 0 {
			query += fmt.Sprintf("rol_software_id = $%d, ", argID)
			args = append(args, *fields.RolSoftwareID)
			argID++
		} else {
			query += "rol_software_id = NULL, "
		}
	}

	if len(args) == 0 && !stringsContainsNull(query) {
		return nil
	}

	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", argID)
	args = append(args, id)
	res, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error al actualizar integrante en PostgreSQL: %w", err)
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

func stringsContainsNull(q string) bool {
	return len(q) > len("UPDATE integrante SET ")
}
