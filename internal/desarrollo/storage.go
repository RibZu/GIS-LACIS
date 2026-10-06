package desarrollo

import (
	"database/sql"
	"errors"
	"fmt"

	"PaginaSEG/internal/integrante"

	"github.com/lib/pq"
)

var ErrNotFound = errors.New("desarrollo no encontrado")
var ErrEmptyID = errors.New("ID de desarrollo vacío o inválido")
var ErrIntegranteInexistente = errors.New("uno de los integrantes elegidos no existe")

type Storage interface {
	Read(id int) (*Desarrollo, error)
	Delete(id int) error
	GetAll() ([]Desarrollo, error)
	CrearCompleto(d *Desarrollo, integranteIDs []int, externos []string) error
	ActualizarCompleto(id int, fields UpdateFields, integranteIDs []int, externos []string) error
	GetIntegrantes(desarrolloID int) ([]integrante.Integrante, error)
	GetParticipantesExternos(desarrolloID int) ([]string, error)
	IntegrantesPorDesarrollo() (map[int][]integrante.Integrante, error)
	ParticipantesExternosPorDesarrollo() (map[int][]string, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
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

func construirUpdate(id int, fields UpdateFields) (string, []interface{}) {
	query := "UPDATE desarrollo SET "
	var args []interface{}
	argID := 1

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
		return "", nil
	}

	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", argID)
	args = append(args, id)
	return query, args
}

func vincularEnTx(tx *sql.Tx, desarrolloID int, integranteIDs []int, externos []string) error {
	if _, err := tx.Exec(`DELETE FROM desarrollo_integrantes WHERE desarrollo_id = $1`, desarrolloID); err != nil {
		return fmt.Errorf("error al limpiar vínculos previos: %w", err)
	}
	for _, integranteID := range integranteIDs {
		if _, err := tx.Exec(
			`INSERT INTO desarrollo_integrantes (desarrollo_id, integrante_id) VALUES ($1, $2)`,
			desarrolloID, integranteID,
		); err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23503" {
				return ErrIntegranteInexistente
			}
			return fmt.Errorf("error al vincular integrante %d: %w", integranteID, err)
		}
	}

	if _, err := tx.Exec(`DELETE FROM participante_externo WHERE desarrollo_id = $1`, desarrolloID); err != nil {
		return fmt.Errorf("error al limpiar participantes externos previos: %w", err)
	}
	for _, nombre := range externos {
		if _, err := tx.Exec(
			`INSERT INTO participante_externo (desarrollo_id, nombre) VALUES ($1, $2)`,
			desarrolloID, nombre,
		); err != nil {
			return fmt.Errorf("error al agregar participante externo %q: %w", nombre, err)
		}
	}
	return nil
}

func (s *PostgresStorage) CrearCompleto(d *Desarrollo, integranteIDs []int, externos []string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO desarrollo (titulo, anio, url, contacto, descripcion)
	          VALUES ($1, $2, $3, $4, $5)
	          RETURNING id`
	if err := tx.QueryRow(query, d.Titulo, d.Anio, d.URL, d.Contacto, d.Descripcion).Scan(&d.ID); err != nil {
		return fmt.Errorf("error al insertar desarrollo en PostgreSQL: %w", err)
	}

	if err := vincularEnTx(tx, d.ID, integranteIDs, externos); err != nil {
		d.ID = 0
		return err
	}
	if err := tx.Commit(); err != nil {
		d.ID = 0
		return fmt.Errorf("error al confirmar la transacción: %w", err)
	}
	return nil
}

func (s *PostgresStorage) ActualizarCompleto(id int, fields UpdateFields, integranteIDs []int, externos []string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("error al iniciar transacción: %w", err)
	}
	defer tx.Rollback()

	query, args := construirUpdate(id, fields)
	if query == "" {
		var existente int
		if err := tx.QueryRow(`SELECT id FROM desarrollo WHERE id = $1`, id).Scan(&existente); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrNotFound
			}
			return fmt.Errorf("error al verificar el desarrollo: %w", err)
		}
	} else {
		res, err := tx.Exec(query, args...)
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
	}

	if err := vincularEnTx(tx, id, integranteIDs, externos); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error al confirmar la transacción: %w", err)
	}
	return nil
}

func (s *PostgresStorage) GetIntegrantes(desarrolloID int) ([]integrante.Integrante, error) {
	query := `
		SELECT i.id, i.nombre, i.apellido
		FROM integrante i
		JOIN desarrollo_integrantes di ON di.integrante_id = i.id
		WHERE di.desarrollo_id = $1
		ORDER BY i.apellido, i.nombre`
	rows, err := s.db.Query(query, desarrolloID)
	if err != nil {
		return nil, fmt.Errorf("error al consultar integrantes del desarrollo: %w", err)
	}
	defer rows.Close()

	var lista []integrante.Integrante
	for rows.Next() {
		var it integrante.Integrante
		if err := rows.Scan(&it.ID, &it.Nombre, &it.Apellido); err != nil {
			return nil, fmt.Errorf("error al escanear integrante vinculado: %w", err)
		}
		lista = append(lista, it)
	}
	return lista, rows.Err()
}

func (s *PostgresStorage) GetParticipantesExternos(desarrolloID int) ([]string, error) {
	query := `SELECT nombre FROM participante_externo WHERE desarrollo_id = $1 ORDER BY nombre`
	rows, err := s.db.Query(query, desarrolloID)
	if err != nil {
		return nil, fmt.Errorf("error al consultar participantes externos: %w", err)
	}
	defer rows.Close()

	var lista []string
	for rows.Next() {
		var nombre string
		if err := rows.Scan(&nombre); err != nil {
			return nil, fmt.Errorf("error al escanear participante externo: %w", err)
		}
		lista = append(lista, nombre)
	}
	return lista, rows.Err()
}

func (s *PostgresStorage) IntegrantesPorDesarrollo() (map[int][]integrante.Integrante, error) {
	query := `
		SELECT di.desarrollo_id, i.id, i.nombre, i.apellido
		FROM desarrollo_integrantes di
		JOIN integrante i ON i.id = di.integrante_id
		ORDER BY di.desarrollo_id, i.apellido, i.nombre`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar los integrantes de todos los desarrollos: %w", err)
	}
	defer rows.Close()

	porDesarrollo := make(map[int][]integrante.Integrante)
	for rows.Next() {
		var desarrolloID int
		var it integrante.Integrante
		if err := rows.Scan(&desarrolloID, &it.ID, &it.Nombre, &it.Apellido); err != nil {
			return nil, fmt.Errorf("error al escanear integrante vinculado: %w", err)
		}
		porDesarrollo[desarrolloID] = append(porDesarrollo[desarrolloID], it)
	}
	return porDesarrollo, rows.Err()
}

func (s *PostgresStorage) ParticipantesExternosPorDesarrollo() (map[int][]string, error) {
	query := `SELECT desarrollo_id, nombre FROM participante_externo ORDER BY desarrollo_id, nombre`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar los participantes externos de todos los desarrollos: %w", err)
	}
	defer rows.Close()

	porDesarrollo := make(map[int][]string)
	for rows.Next() {
		var desarrolloID int
		var nombre string
		if err := rows.Scan(&desarrolloID, &nombre); err != nil {
			return nil, fmt.Errorf("error al escanear participante externo: %w", err)
		}
		porDesarrollo[desarrolloID] = append(porDesarrollo[desarrolloID], nombre)
	}
	return porDesarrollo, rows.Err()
}
