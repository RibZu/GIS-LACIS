package desarrollo

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("desarrollo no encontrado")
var ErrEmptyID = errors.New("ID de desarrollo vacío o inválido")

type Storage interface {
	Create(d *Desarrollo) error
	Read(id int) (*Desarrollo, error)
	Delete(id int) error
	Update(id int, fields UpdateFields) error
	GetAll() ([]Desarrollo, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) Create(d *Desarrollo) error {
	query := `INSERT INTO desarrollo (titulo, anio, url, contacto, descripcion)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING id`

	err := s.db.QueryRow(query, d.Titulo, d.Anio, d.URL, d.Contacto, d.Descripcion).Scan(&d.ID)
	if err != nil {
		return fmt.Errorf("error al insertar desarrollo en PostgreSQL: %w", err)
	}
	return nil
}

func (s *PostgresStorage) Read(id int) (*Desarrollo, error) {
	query := `SELECT id, titulo, anio, COALESCE(url, ''), COALESCE(contacto, ''), COALESCE(descripcion, '')
	          FROM desarrollo WHERE id = $1`
	row := s.db.QueryRow(query, id)

	var d Desarrollo
	err := row.Scan(&d.ID, &d.Titulo, &d.Anio, &d.URL, &d.Contacto, &d.Descripcion)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error al leer desarrollo en PostgreSQL: %w", err)
	}
	return &d, nil
}

func (s *PostgresStorage) GetAll() ([]Desarrollo, error) {
	query := `SELECT id, titulo, anio, COALESCE(url, ''), COALESCE(contacto, ''), COALESCE(descripcion, '')
	          FROM desarrollo ORDER BY anio DESC, id DESC`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar desarrollos en PostgreSQL: %w", err)
	}
	defer rows.Close()

	var lista []Desarrollo
	for rows.Next() {
		var d Desarrollo
		if err := rows.Scan(&d.ID, &d.Titulo, &d.Anio, &d.URL, &d.Contacto, &d.Descripcion); err != nil {
			return nil, fmt.Errorf("error al escanear desarrollo: %w", err)
		}
		lista = append(lista, d)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return lista, nil
}

func (s *PostgresStorage) Delete(id int) error {
	query := `DELETE FROM desarrollo WHERE id = $1`
	res, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar desarrollo en PostgreSQL: %w", err)
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
	var query string = "UPDATE desarrollo SET "
	var args []interface{}
	var argID int = 1

	if fields.Titulo != nil {
		query += fmt.Sprintf("titulo = $%d, ", argID)
		args = append(args, *fields.Titulo)
		argID++
	}
	if fields.Anio != nil {
		query += fmt.Sprintf("anio = $%d, ", argID)
		args = append(args, *fields.Anio)
		argID++
	}
	if fields.URL != nil {
		query += fmt.Sprintf("url = $%d, ", argID)
		args = append(args, *fields.URL)
		argID++
	}
	if fields.Contacto != nil {
		query += fmt.Sprintf("contacto = $%d, ", argID)
		args = append(args, *fields.Contacto)
		argID++
	}
	if fields.Descripcion != nil {
		query += fmt.Sprintf("descripcion = $%d, ", argID)
		args = append(args, *fields.Descripcion)
		argID++
	}

	if len(args) == 0 {
		return nil
	}

	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", argID)
	args = append(args, id)
	res, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error al actualizar desarrollo en PostgreSQL: %w", err)
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
