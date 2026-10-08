package configuracion

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("configuración no encontrada")
)

type Storage interface {
	Get() (*ConfiguracionSitio, error)
	Update(cfg *ConfiguracionSitio) error
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(db *sql.DB) (*PostgresStorage, error) {
	s := &PostgresStorage{db: db}
	if err := s.ensureSchema(); err != nil {
		return nil, fmt.Errorf("error al asegurar esquema de configuracion_sitio: %w", err)
	}
	return s, nil
}

func (s *PostgresStorage) ensureSchema() error {
	// Verificar si la tabla existe, crearla si no
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS configuracion_sitio (
		id SERIAL PRIMARY KEY,
		usuario_gestor_id INT REFERENCES usuario_gestor(id) ON DELETE SET NULL,
		telefono_footer VARCHAR(50),
		email_footer VARCHAR(255),
		direccion VARCHAR(255),
		ubicacion_footer VARCHAR(255),
		oficinas_footer VARCHAR(255),
		mapa_url TEXT,
		email_doctorado VARCHAR(255),
		email_maestria_calidad VARCHAR(255),
		email_maestria_software VARCHAR(255),
		email_especializacion VARCHAR(255),
		email_ing_informatica VARCHAR(255),
		email_tecnicatura_web VARCHAR(255)
	);`
	if _, err := s.db.Exec(createTableQuery); err != nil {
		return err
	}

	// Agregar columnas que pudieran faltar si la tabla existía previamente con menos campos
	alterQueries := []string{
		`ALTER TABLE configuracion_sitio ADD COLUMN IF NOT EXISTS ubicacion_footer VARCHAR(255);`,
		`ALTER TABLE configuracion_sitio ADD COLUMN IF NOT EXISTS oficinas_footer VARCHAR(255);`,
		`ALTER TABLE configuracion_sitio ADD COLUMN IF NOT EXISTS mapa_url TEXT;`,
		`ALTER TABLE configuracion_sitio ADD COLUMN IF NOT EXISTS email_doctorado VARCHAR(255);`,
		`ALTER TABLE configuracion_sitio ADD COLUMN IF NOT EXISTS email_maestria_calidad VARCHAR(255);`,
		`ALTER TABLE configuracion_sitio ADD COLUMN IF NOT EXISTS email_maestria_software VARCHAR(255);`,
		`ALTER TABLE configuracion_sitio ADD COLUMN IF NOT EXISTS email_especializacion VARCHAR(255);`,
		`ALTER TABLE configuracion_sitio ADD COLUMN IF NOT EXISTS email_ing_informatica VARCHAR(255);`,
		`ALTER TABLE configuracion_sitio ADD COLUMN IF NOT EXISTS email_tecnicatura_web VARCHAR(255);`,
	}

	for _, q := range alterQueries {
		if _, err := s.db.Exec(q); err != nil {
			return err
		}
	}

	// Verificar si existe al menos una fila; de no existir, insertar la configuración inicial por defecto
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM configuracion_sitio`).Scan(&count); err != nil {
		return err
	}

	if count == 0 {
		insertInitial := `
		INSERT INTO configuracion_sitio (
			telefono_footer, email_footer, direccion, ubicacion_footer, oficinas_footer,
			email_doctorado, email_maestria_calidad, email_maestria_software,
			email_especializacion, email_ing_informatica, email_tecnicatura_web
		) VALUES (
			'4520300 Int. 2101', 'contacto@unsl.edu.ar', 'Ejército de los Andes 950', 'Bloque II - 1 piso', 'N° 1-2-3-26-27-28-29',
			'agaris@gmail.com', 'magister@unsl.edu.ar', 'magister@unsl.edu.ar',
			'mgperalta03@gmail.com', 'ingenieriainformaticaunsl@gmail.com', ''
		);`
		if _, err := s.db.Exec(insertInitial); err != nil {
			return err
		}
	}

	return nil
}

func (s *PostgresStorage) Get() (*ConfiguracionSitio, error) {
	query := `
	SELECT id, usuario_gestor_id,
	       COALESCE(telefono_footer, ''),
	       COALESCE(email_footer, ''),
	       COALESCE(direccion, ''),
	       COALESCE(ubicacion_footer, ''),
	       COALESCE(oficinas_footer, ''),
	       COALESCE(mapa_url, ''),
	       COALESCE(email_doctorado, ''),
	       COALESCE(email_maestria_calidad, ''),
	       COALESCE(email_maestria_software, ''),
	       COALESCE(email_especializacion, ''),
	       COALESCE(email_ing_informatica, ''),
	       COALESCE(email_tecnicatura_web, '')
	FROM configuracion_sitio
	ORDER BY id ASC
	LIMIT 1;`

	var cfg ConfiguracionSitio
	err := s.db.QueryRow(query).Scan(
		&cfg.ID,
		&cfg.UsuarioGestorID,
		&cfg.TelefonoFooter,
		&cfg.EmailFooter,
		&cfg.Direccion,
		&cfg.UbicacionFooter,
		&cfg.OficinasFooter,
		&cfg.MapaURL,
		&cfg.EmailDoctorado,
		&cfg.EmailMaestriaCalidad,
		&cfg.EmailMaestriaSoftware,
		&cfg.EmailEspecializacion,
		&cfg.EmailIngInformatica,
		&cfg.EmailTecnicaturaWeb,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("error al obtener configuración: %w", err)
	}

	return &cfg, nil
}

func (s *PostgresStorage) Update(cfg *ConfiguracionSitio) error {
	query := `
	UPDATE configuracion_sitio
	SET telefono_footer = $1,
	    email_footer = $2,
	    direccion = $3,
	    ubicacion_footer = $4,
	    oficinas_footer = $5,
	    mapa_url = $6,
	    email_doctorado = $7,
	    email_maestria_calidad = $8,
	    email_maestria_software = $9,
	    email_especializacion = $10,
	    email_ing_informatica = $11,
	    email_tecnicatura_web = $12
	WHERE id = (SELECT id FROM configuracion_sitio ORDER BY id ASC LIMIT 1);`

	res, err := s.db.Exec(query,
		cfg.TelefonoFooter,
		cfg.EmailFooter,
		cfg.Direccion,
		cfg.UbicacionFooter,
		cfg.OficinasFooter,
		cfg.MapaURL,
		cfg.EmailDoctorado,
		cfg.EmailMaestriaCalidad,
		cfg.EmailMaestriaSoftware,
		cfg.EmailEspecializacion,
		cfg.EmailIngInformatica,
		cfg.EmailTecnicaturaWeb,
	)
	if err != nil {
		return fmt.Errorf("error al actualizar configuración: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		// En caso de que no hubiera fila, la insertamos
		insertQuery := `
		INSERT INTO configuracion_sitio (
			telefono_footer, email_footer, direccion, ubicacion_footer, oficinas_footer, mapa_url,
			email_doctorado, email_maestria_calidad, email_maestria_software,
			email_especializacion, email_ing_informatica, email_tecnicatura_web
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`
		_, err = s.db.Exec(insertQuery,
			cfg.TelefonoFooter,
			cfg.EmailFooter,
			cfg.Direccion,
			cfg.UbicacionFooter,
			cfg.OficinasFooter,
			cfg.MapaURL,
			cfg.EmailDoctorado,
			cfg.EmailMaestriaCalidad,
			cfg.EmailMaestriaSoftware,
			cfg.EmailEspecializacion,
			cfg.EmailIngInformatica,
			cfg.EmailTecnicaturaWeb,
		)
		if err != nil {
			return fmt.Errorf("error al insertar fila inicial de configuración: %w", err)
		}
	}

	return nil
}
