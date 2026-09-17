package proyecto

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound = errors.New("proyecto no encontrado")
	ErrEmptyID  = errors.New("ID de proyecto inválido")
)

type Storage interface {
	Create(p *Proyecto) error
	Read(id int) (*Proyecto, error)
	Delete(id int) error
	Update(id int, fields UpdateFields) error
	GetAll() ([]Proyecto, error)
	GetAllAdmin() ([]Proyecto, error)
	Restaurar(id int) error

	GetEquipo(proyectoID int) ([]MiembroEquipo, error)
	AgregarMiembro(proyectoID, integranteID int, rolEnProyecto string) error
	QuitarMiembro(proyectoID, integranteID int) error
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) Create(p *Proyecto) error {
	query := `INSERT INTO proyecto (titulo, descripcion, enlace, equipo_historico, anio_inicio, anio_fin)
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := s.db.QueryRow(query, p.Titulo, p.Descripcion, p.Enlace, p.EquipoHistorico, p.AnioInicio, p.AnioFin).Scan(&p.ID)
	if err != nil {
		return fmt.Errorf("error al insertar proyecto en PostgreSQL: %w", err)
	}
	return nil
}

func (s *PostgresStorage) Read(id int) (*Proyecto, error) {
	if id <= 0 {
		return nil, ErrEmptyID
	}
	query := `SELECT id, titulo, descripcion, enlace, equipo_historico, anio_inicio, anio_fin 
	          FROM proyecto WHERE id = $1`
	var p Proyecto
	var descripcion, enlace, equipoHistorico sql.NullString

	err := s.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Titulo,
		&descripcion,
		&enlace,
		&equipoHistorico,
		&p.AnioInicio,
		&p.AnioFin,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error al leer proyecto con ID %d: %w", id, err)
	}

	if descripcion.Valid {
		p.Descripcion = descripcion.String
	}
	if enlace.Valid {
		p.Enlace = enlace.String
	}
	if equipoHistorico.Valid {
		p.EquipoHistorico = equipoHistorico.String
	}

	return &p, nil
}

func (s *PostgresStorage) Delete(id int) error {
	if id <= 0 {
		return ErrEmptyID
	}
	query := `UPDATE proyecto SET activo = false WHERE id = $1`
	res, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al dar de baja lógica el proyecto con ID %d: %w", id, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostgresStorage) Update(id int, fields UpdateFields) error {
	if id <= 0 {
		return ErrEmptyID
	}

	var setClauses []string
	var args []interface{}
	argID := 1

	if fields.Titulo != nil {
		setClauses = append(setClauses, fmt.Sprintf("titulo = $%d", argID))
		args = append(args, *fields.Titulo)
		argID++
	}
	if fields.Descripcion != nil {
		setClauses = append(setClauses, fmt.Sprintf("descripcion = $%d", argID))
		args = append(args, *fields.Descripcion)
		argID++
	}
	if fields.Enlace != nil {
		setClauses = append(setClauses, fmt.Sprintf("enlace = $%d", argID))
		args = append(args, *fields.Enlace)
		argID++
	}
	if fields.EquipoHistorico != nil {
		setClauses = append(setClauses, fmt.Sprintf("equipo_historico = $%d", argID))
		args = append(args, *fields.EquipoHistorico)
		argID++
	}
	if fields.AnioInicio != nil {
		setClauses = append(setClauses, fmt.Sprintf("anio_inicio = $%d", argID))
		args = append(args, *fields.AnioInicio)
		argID++
	}
	if fields.AnioFin != nil {
		setClauses = append(setClauses, fmt.Sprintf("anio_fin = $%d", argID))
		args = append(args, *fields.AnioFin)
		argID++
	}

	if len(setClauses) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE proyecto SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argID)
	args = append(args, id)

	res, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error al actualizar proyecto con ID %d: %w", id, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresStorage) GetAll() ([]Proyecto, error) {
	query := `SELECT id, titulo, descripcion, enlace, equipo_historico, anio_inicio, anio_fin, activo 
	          FROM proyecto WHERE activo = true ORDER BY anio_fin DESC, anio_inicio DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar proyectos activos: %w", err)
	}
	defer rows.Close()

	proyectos := make([]Proyecto, 0)
	for rows.Next() {
		var p Proyecto
		var descripcion, enlace, equipoHistorico sql.NullString

		if err := rows.Scan(
			&p.ID, &p.Titulo, &descripcion, &enlace, &equipoHistorico, &p.AnioInicio, &p.AnioFin, &p.Activo,
		); err != nil {
			return nil, fmt.Errorf("error al escanear fila: %w", err)
		}

		if descripcion.Valid {
			p.Descripcion = descripcion.String
		}
		if enlace.Valid {
			p.Enlace = enlace.String
		}
		if equipoHistorico.Valid {
			p.EquipoHistorico = equipoHistorico.String
		}

		proyectos = append(proyectos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error en iteración: %w", err)
	}
	return proyectos, nil
}
func (s *PostgresStorage) GetAllAdmin() ([]Proyecto, error) {
	query := `SELECT id, titulo, descripcion, enlace, equipo_historico, anio_inicio, anio_fin, activo 
	          FROM proyecto ORDER BY anio_fin DESC, anio_inicio DESC`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar todos los proyectos: %w", err)
	}
	defer rows.Close()

	proyectos := make([]Proyecto, 0)
	for rows.Next() {
		var p Proyecto
		var descripcion, enlace, equipoHistorico sql.NullString

		if err := rows.Scan(
			&p.ID, &p.Titulo, &descripcion, &enlace, &equipoHistorico, &p.AnioInicio, &p.AnioFin, &p.Activo,
		); err != nil {
			return nil, fmt.Errorf("error al escanear fila: %w", err)
		}

		if descripcion.Valid {
			p.Descripcion = descripcion.String
		}
		if enlace.Valid {
			p.Enlace = enlace.String
		}
		if equipoHistorico.Valid {
			p.EquipoHistorico = equipoHistorico.String
		}

		proyectos = append(proyectos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error en iteración: %w", err)
	}
	return proyectos, nil
}

func (s *PostgresStorage) Restaurar(id int) error {
	if id <= 0 {
		return ErrEmptyID
	}
	query := `UPDATE proyecto SET activo = true WHERE id = $1`
	res, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al restaurar el proyecto con ID %d: %w", id, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostgresStorage) GetEquipo(proyectoID int) ([]MiembroEquipo, error) {
	query := `SELECT pi.integrante_id, i.nombre, i.apellido, COALESCE(pi.rol_en_proyecto, '')
	          FROM proyecto_integrantes pi
	          JOIN integrante i ON i.id = pi.integrante_id
	          WHERE pi.proyecto_id = $1
	          ORDER BY pi.id ASC`
	rows, err := s.db.Query(query, proyectoID)
	if err != nil {
		return nil, fmt.Errorf("error al consultar equipo del proyecto: %w", err)
	}
	defer rows.Close()

	equipo := make([]MiembroEquipo, 0)
	for rows.Next() {
		var m MiembroEquipo
		if err := rows.Scan(&m.IntegranteID, &m.Nombre, &m.Apellido, &m.RolEnProyecto); err != nil {
			return nil, fmt.Errorf("error al escanear miembro de equipo: %w", err)
		}
		equipo = append(equipo, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error en iteración de equipo: %w", err)
	}
	return equipo, nil
}

func (s *PostgresStorage) AgregarMiembro(proyectoID, integranteID int, rolEnProyecto string) error {
	query := `INSERT INTO proyecto_integrantes (proyecto_id, integrante_id, rol_en_proyecto)
	          VALUES ($1, $2, $3)
	          ON CONFLICT (proyecto_id, integrante_id) DO UPDATE SET rol_en_proyecto = EXCLUDED.rol_en_proyecto`
	_, err := s.db.Exec(query, proyectoID, integranteID, rolEnProyecto)
	if err != nil {
		return fmt.Errorf("error al agregar miembro al equipo: %w", err)
	}
	return nil
}

func (s *PostgresStorage) QuitarMiembro(proyectoID, integranteID int) error {
	query := `DELETE FROM proyecto_integrantes WHERE proyecto_id = $1 AND integrante_id = $2`
	res, err := s.db.Exec(query, proyectoID, integranteID)
	if err != nil {
		return fmt.Errorf("error al quitar miembro del equipo: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
