package colaborador

import (
	"database/sql"
	"fmt"
)

// Storage define la interfaz de acceso a datos para colaboradores.
type Storage interface {
	GetAll() ([]Colaborador, error)
	Create(c *Colaborador) error
	Read(id int) (*Colaborador, error)
	Update(id int, fields UpdateFields) error
	Delete(id int) error
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) GetAll() ([]Colaborador, error) {
	query := `SELECT id, COALESCE(descripcion, ''), COALESCE(logo_url, ''), COALESCE(activo, true) 
	          FROM colaboradores ORDER BY id ASC`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar colaboradores: %w", err)
	}
	defer rows.Close()
	colaboradores := make([]Colaborador, 0)
	for rows.Next() {
		var c Colaborador
		if err := rows.Scan(&c.ID, &c.Descripcion, &c.LogoURL, &c.Activo); err != nil {
			return nil, fmt.Errorf("error al escanear fila de colaborador: %w", err)
		}
		colaboradores = append(colaboradores, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error en iteración de filas de colaboradores: %w", err)
	}
	return colaboradores, nil
}

func (s *PostgresStorage) Create(c *Colaborador) error {
	query := `INSERT INTO colaboradores (descripcion, logo_url, activo) 
	          VALUES ($1, $2, true) RETURNING id`
	err := s.db.QueryRow(query, c.Descripcion, c.LogoURL).Scan(&c.ID)
	if err != nil {
		return fmt.Errorf("error al crear colaborador: %w", err)
	}
	c.Activo = true
	return nil
}

func (s *PostgresStorage) Read(id int) (*Colaborador, error) {
	query := `SELECT id, COALESCE(descripcion, ''), COALESCE(logo_url, ''), COALESCE(activo, true) 
	          FROM colaboradores WHERE id = $1`
	var c Colaborador
	err := s.db.QueryRow(query, id).Scan(&c.ID, &c.Descripcion, &c.LogoURL, &c.Activo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("colaborador no encontrado")
		}
		return nil, fmt.Errorf("error al leer colaborador: %w", err)
	}
	return &c, nil
}

func (s *PostgresStorage) Update(id int, fields UpdateFields) error {
	var query string = "UPDATE colaboradores SET "
	var args []interface{}
	var argID int = 1
	if fields.Descripcion != nil {
		query += fmt.Sprintf("descripcion = $%d, ", argID)
		args = append(args, *fields.Descripcion)
		argID++
	}
	if fields.LogoURL != nil {
		query += fmt.Sprintf("logo_url = $%d, ", argID)
		args = append(args, *fields.LogoURL)
		argID++
	}
	if fields.Activo != nil {
		query += fmt.Sprintf("activo = $%d, ", argID)
		args = append(args, *fields.Activo)
		argID++
	}
	if len(args) == 0 {
		return nil
	}
	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", argID)
	args = append(args, id)
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *PostgresStorage) Delete(id int) error {
	query := `UPDATE colaboradores SET activo = false WHERE id = $1`
	res, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar (lógico) colaborador: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("colaborador no encontrado")
	}
	return nil
}
