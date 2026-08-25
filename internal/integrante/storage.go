package integrante

import (
	"database/sql"
	"errors"
	"fmt"
)

// ErrNotFound devuldce el error Not found, cuando no se encuentra el USUARIO CON EL ID proporcionado
var ErrNotFound = errors.New("usuario no encontrado")

// ErrEmptyID retorna eso cuando el usuario no tiene el ID VACIO
var ErrEmptyID = errors.New("Id de usuario vacio")

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

	query := `INSERT INTO integrante (nombre, apellido, cv, imagen, contacto, especializacion, descripcion, rol)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	          RETURNING id`

	err := c.db.QueryRow(
		query,
		integrante.Nombre,
		integrante.Apellido,
		integrante.CV,
		integrante.Imagen,
		integrante.Contacto,
		integrante.Especializacion,
		integrante.Descripcion,
		integrante.RolID,
	).Scan(&integrante.ID)

	if err != nil {
		return fmt.Errorf("error al insertar integrante en PostgreSQL: %w", err)
	}

	return nil
}

func (s *PostgresStorage) Read(id int) (*Integrante, error) {
	query := `
		SELECT i.id, i.nombre, i.apellido, i.cv, i.imagen, i.contacto, i.especializacion, i.descripcion, i.rol, r.nombre
		FROM integrante i
		LEFT JOIN rol r ON i.rol = r.id
		WHERE i.id = $1`
	row := s.db.QueryRow(query, id)
	var u Integrante

	err := row.Scan(
		&u.ID,
		&u.Nombre,
		&u.Apellido,
		&u.CV,
		&u.Imagen,
		&u.Contacto,
		&u.Especializacion,
		&u.Descripcion,
		&u.RolID,
		&u.Rol.Nombre,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error al leer integrante en PostgreSQL: %w", err)
	}
	u.Rol.ID = u.RolID
	return &u, nil
}

func (s *PostgresStorage) GetAll() ([]Integrante, error) {
	query := `
		SELECT i.id, i.nombre, i.apellido, i.cv, i.imagen, i.contacto, i.especializacion, i.descripcion, i.rol, r.nombre
		FROM integrante i
		LEFT JOIN rol r ON i.rol = r.id
		ORDER BY i.apellido ASC`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar integrantes en PostgreSQL: %w", err)
	}
	defer rows.Close()
	var integrantes []Integrante
	for rows.Next() {
		var u Integrante
		err := rows.Scan(
			&u.ID,
			&u.Nombre,
			&u.Apellido,
			&u.CV,
			&u.Imagen,
			&u.Contacto,
			&u.Especializacion,
			&u.Descripcion,
			&u.RolID,
			&u.Rol.Nombre,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear integrante: %w", err)
		}
		u.Rol.ID = u.RolID
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
	// Verificamos si se eliminó alguna fila
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
	if fields.CV != nil {
		query += fmt.Sprintf("cv = $%d, ", argID)
		args = append(args, *fields.CV)
		argID++
	}

	if fields.Imagen != nil {
		query += fmt.Sprintf("imagen = $%d, ", argID)
		args = append(args, *fields.Imagen)
		argID++
	}

	if fields.Contacto != nil {
		query += fmt.Sprintf("contacto = $%d, ", argID)
		args = append(args, *fields.Contacto)
		argID++
	}
	if fields.Especializacion != nil {
		query += fmt.Sprintf("especializacion = $%d, ", argID)
		args = append(args, *fields.Especializacion)
		argID++
	}
	if fields.Descripcion != nil {
		query += fmt.Sprintf("descripcion = $%d, ", argID)
		args = append(args, *fields.Descripcion)
		argID++
	}
	if fields.RolID != nil {
		query += fmt.Sprintf("rol = $%d, ", argID)
		args = append(args, *fields.RolID)
		argID++
	}
	// Si no se envió ningún campo para actualizar, salimos sin error
	if len(args) == 0 {
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
