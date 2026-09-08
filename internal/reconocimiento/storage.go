package reconocimiento

import (
    "database/sql"
    "fmt"
)

// Storage defines the data access methods for reconocimiento.
type Storage interface {
    GetAll() ([]Reconocimiento, error)
    Create(r *Reconocimiento) error
    Read(id int) (*Reconocimiento, error)
    Update(id int, fields UpdateFields) error
    Delete(id int) error // logical delete sets activo = false
}

type PostgresStorage struct {
    db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
    return &PostgresStorage{db: db}
}

func (s *PostgresStorage) GetAll() ([]Reconocimiento, error) {
    query := `SELECT id, titulo, descripcion, activo FROM reconocimientos ORDER BY id ASC`
    rows, err := s.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("error al consultar reconocimientos: %w", err)
    }
    defer rows.Close()

    reconocimientos := make([]Reconocimiento, 0)
    for rows.Next() {
        var r Reconocimiento
        var descripcion sql.NullString
        var activo sql.NullBool
        if err := rows.Scan(&r.ID, &r.Titulo, &descripcion, &activo); err != nil {
            return nil, fmt.Errorf("error al escanear fila de reconocimiento: %w", err)
        }
        if descripcion.Valid { r.Descripcion = descripcion.String }
        if activo.Valid { r.Activo = activo.Bool }
        reconocimientos = append(reconocimientos, r)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error en iteracion de filas de reconocimientos: %w", err)
    }
    return reconocimientos, nil
}

func (s *PostgresStorage) Create(r *Reconocimiento) error {
    query := `INSERT INTO reconocimientos (titulo, descripcion, activo) VALUES ($1, $2, true) RETURNING id`
    err := s.db.QueryRow(query, r.Titulo, r.Descripcion).Scan(&r.ID)
    if err != nil {
        return fmt.Errorf("error al crear reconocimiento: %w", err)
    }
    r.Activo = true
    return nil
}

func (s *PostgresStorage) Read(id int) (*Reconocimiento, error) {
    query := `SELECT id, titulo, descripcion, activo FROM reconocimientos WHERE id = $1`
    var r Reconocimiento
    var descripcion sql.NullString
    var activo sql.NullBool
    err := s.db.QueryRow(query, id).Scan(&r.ID, &r.Titulo, &descripcion, &activo)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("reconocimiento no encontrado")
        }
        return nil, fmt.Errorf("error al leer reconocimiento: %w", err)
    }
    if descripcion.Valid { r.Descripcion = descripcion.String }
    if activo.Valid { r.Activo = activo.Bool }
    return &r, nil
}

func (s *PostgresStorage) Update(id int, fields UpdateFields) error {
    var query string = "UPDATE reconocimientos SET "
    var args []interface{}
    var argID int = 1
    if fields.Titulo != nil {
        query += fmt.Sprintf("titulo = $%d, ", argID)
        args = append(args, *fields.Titulo)
        argID++
    }
    if fields.Descripcion != nil {
        query += fmt.Sprintf("descripcion = $%d, ", argID)
        args = append(args, *fields.Descripcion)
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
    // Trim trailing comma and space
    query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", argID)
    args = append(args, id)
    _, err := s.db.Exec(query, args...)
    return err
}

func (s *PostgresStorage) Delete(id int) error {
    // Logical delete: set activo to false
    query := `UPDATE reconocimientos SET activo = false WHERE id = $1`
    res, err := s.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("error al eliminar (lógico) reconocimiento: %w", err)
    }
    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if rowsAffected == 0 {
        return fmt.Errorf("reconocimiento no encontrado")
    }
    return nil
}
