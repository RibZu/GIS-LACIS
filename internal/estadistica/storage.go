package estadistica

import (
	"database/sql"
	"fmt"
)

type Storage interface {
	ContarIntegrantesActivosLacis() (int, error)
	NombresParticipantesExternos() ([]string, error)
	ResumenProductos() ([]ProductoResumen, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func (s *PostgresStorage) ContarIntegrantesActivosLacis() (int, error) {
	var total int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM integrante WHERE activo AND pertenece_lacis`).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("error al contar integrantes activos de LaCIS: %w", err)
	}
	return total, nil
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
	query := `SELECT d.id, d.titulo, d.anio, COALESCE(d.url, ''), COALESCE(d.contacto, ''), COALESCE(d.descripcion, ''),
	                 EXISTS (SELECT 1 FROM desarrollo_integrantes di WHERE di.desarrollo_id = d.id)
	                 OR EXISTS (SELECT 1 FROM participante_externo pe WHERE pe.desarrollo_id = d.id)
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
		if err := rows.Scan(&p.ID, &p.Titulo, &p.Anio, &p.URL, &p.Contacto, &p.Descripcion, &p.TieneParticipantes); err != nil {
			return nil, fmt.Errorf("error al escanear el resumen de un producto: %w", err)
		}
		productos = append(productos, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return productos, nil
}
