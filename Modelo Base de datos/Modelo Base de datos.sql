-- =============================================================================
-- Script DDL para PostgreSQL
-- Basado en el Diagrama UML del Sistema
-- =============================================================================

-- Desactivar salidas intermedias y configurar manejo de transacciones
BEGIN;

-- -----------------------------------------------------------------------------
-- 1. PAQUETE: NÚCLEO Y ROLES
-- -----------------------------------------------------------------------------

CREATE TABLE rol (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL
);

CREATE TABLE integrante (
    id SERIAL PRIMARY KEY,
    rol_id INT REFERENCES rol(id) ON DELETE SET NULL,
    activo BOOLEAN NOT NULL DEFAULT TRUE, -- Baja lógica
    nombre VARCHAR(100) NOT NULL,
    apellido VARCHAR(100) NOT NULL,
    titulo_especializacion VARCHAR(255),
    descripcion TEXT,
    contacto_mail VARCHAR(255),
    contacto_linkedin VARCHAR(255),
    imagen_url VARCHAR(500),
    cv_url VARCHAR(500),
    pertenece_lacis BOOLEAN NOT NULL DEFAULT FALSE,
    pertenece_grupo_software BOOLEAN NOT NULL DEFAULT FALSE
);

-- -----------------------------------------------------------------------------
-- 2. PAQUETE: GESTIÓN E INFRAESTRUCTURA
-- -----------------------------------------------------------------------------

CREATE TABLE usuario_gestor (
    id SERIAL PRIMARY KEY,
    integrante_id INT UNIQUE REFERENCES integrante(id) ON DELETE SET NULL, -- Relación 0..1 opcional
    username VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    ultimo_acceso TIMESTAMP WITH TIME ZONE,
    rol VARCHAR(50) NOT NULL DEFAULT 'GESTOR', -- Calculado: ADMIN, nombre de módulo, o GESTOR si tiene 2+ módulos
    modulos VARCHAR(255) NOT NULL DEFAULT '' -- Lista de módulos separados por coma, ej: "integrantes,proyectos"
);

CREATE TABLE configuracion_sitio (
    id SERIAL PRIMARY KEY,
    usuario_gestor_id INT REFERENCES usuario_gestor(id) ON DELETE SET NULL, -- Modificado por
    telefono_footer VARCHAR(50),
    email_footer VARCHAR(255),
    direccion VARCHAR(255)
);

CREATE TABLE colaboradores (
    id BIGSERIAL PRIMARY KEY,
    usuario_gestor_id INT REFERENCES usuario_gestor(id) ON DELETE SET NULL,
    descripcion TEXT,
    logo_url VARCHAR(500)
);

-- -----------------------------------------------------------------------------
-- 3. PAQUETE: PROYECTOS Y DESARROLLOS
-- -----------------------------------------------------------------------------

CREATE TABLE proyecto (
    id SERIAL PRIMARY KEY,
    titulo VARCHAR(255) NOT NULL,
    descripcion TEXT,
    enlace VARCHAR(500),
    equipo_historico VARCHAR(500)
);

CREATE TABLE desarrollo (
    id SERIAL PRIMARY KEY,
    titulo VARCHAR(255) NOT NULL,
    url VARCHAR(500),
    contacto VARCHAR(255),
    descripcion TEXT
);

CREATE TABLE reconocimientos (
    id BIGSERIAL PRIMARY KEY,
    titulo VARCHAR(255) NOT NULL,
    descripcion TEXT
);

CREATE TABLE participante_externo (
    id SERIAL PRIMARY KEY,
    desarrollo_id INT NOT NULL REFERENCES desarrollo(id) ON DELETE CASCADE,
    nombre VARCHAR(150) NOT NULL
);

CREATE TABLE proyecto_integrantes (
    id SERIAL PRIMARY KEY,
    integrante_id INT NOT NULL REFERENCES integrante(id) ON DELETE CASCADE,
    proyecto_id INT NOT NULL REFERENCES proyecto(id) ON DELETE CASCADE,
    CONSTRAINT uq_proyecto_integrante UNIQUE (proyecto_id, integrante_id)
);

CREATE TABLE desarrollo_integrantes (
    id SERIAL PRIMARY KEY,
    desarrollo_id INT NOT NULL REFERENCES desarrollo(id) ON DELETE CASCADE,
    integrante_id INT NOT NULL REFERENCES integrante(id) ON DELETE CASCADE,
    CONSTRAINT uq_desarrollo_integrante UNIQUE (desarrollo_id, integrante_id)
);

CREATE TABLE proyecto_reconocimientos (
    id SERIAL PRIMARY KEY,
    proyecto_id INT NOT NULL REFERENCES proyecto(id) ON DELETE CASCADE,
    reconocimiento_id BIGINT NOT NULL REFERENCES reconocimientos(id) ON DELETE CASCADE,
    CONSTRAINT uq_proyecto_reconocimiento UNIQUE (proyecto_id, reconocimiento_id)
);

CREATE TABLE desarrollo_reconocimientos (
    id SERIAL PRIMARY KEY,
    desarrollo_id INT NOT NULL REFERENCES desarrollo(id) ON DELETE CASCADE,
    reconocimiento_id BIGINT NOT NULL REFERENCES reconocimientos(id) ON DELETE CASCADE,
    CONSTRAINT uq_desarrollo_reconocimiento UNIQUE (desarrollo_id, reconocimiento_id)
);

-- -----------------------------------------------------------------------------
-- 4. PAQUETE: PRODUCCIÓN ACADÉMICA
-- -----------------------------------------------------------------------------

CREATE TABLE tesis (
    id SERIAL PRIMARY KEY,
    autor_id INT REFERENCES integrante(id) ON DELETE SET NULL,
    autor_historico VARCHAR(255),
    director_id INT REFERENCES integrante(id) ON DELETE SET NULL,
    director_historico VARCHAR(255),
    coodirector_id INT REFERENCES integrante(id) ON DELETE SET NULL,
    coodirector_historico VARCHAR(255),
    proyecto_id INT REFERENCES proyecto(id) ON DELETE SET NULL,
    titulo VARCHAR(300) NOT NULL,
    anio INT,
    palabras_clave VARCHAR(500),
    resumen TEXT,
    archivo_pdf VARCHAR(500),
    nivel VARCHAR(100),
    carrera_origen VARCHAR(200)
);

CREATE TABLE integrantes_tesis (
    id SERIAL PRIMARY KEY,
    tesis_id INT NOT NULL REFERENCES tesis(id) ON DELETE CASCADE,
    integrante_id INT NOT NULL REFERENCES integrante(id) ON DELETE CASCADE,
    CONSTRAINT uq_tesis_integrante UNIQUE (tesis_id, integrante_id)
);

-- -----------------------------------------------------------------------------
-- 5. PAQUETE: INSTITUCIONAL (LACIS)
-- -----------------------------------------------------------------------------

CREATE TABLE lacis (
    id SERIAL PRIMARY KEY,
    proyecto_id INT REFERENCES proyecto(id) ON DELETE SET NULL,
    director_id INT REFERENCES integrante(id) ON DELETE SET NULL,
    comite_id INT REFERENCES integrante(id) ON DELETE SET NULL
);

CREATE TABLE integrantes_lacis (
    id SERIAL PRIMARY KEY,
    lacis_id INT NOT NULL REFERENCES lacis(id) ON DELETE CASCADE,
    integrante_id INT NOT NULL REFERENCES integrante(id) ON DELETE CASCADE,
    CONSTRAINT uq_lacis_integrante UNIQUE (lacis_id, integrante_id)
);

-- -----------------------------------------------------------------------------
-- ÍNDICES PARA MEJORA DE RENDIMIENTO EN CLAVES FORÁNEAS
-- -----------------------------------------------------------------------------

CREATE INDEX idx_integrante_rol ON integrante(rol_id);
CREATE INDEX idx_usuario_gestor_integrante ON usuario_gestor(integrante_id);
CREATE INDEX idx_configuracion_gestor ON configuracion_sitio(usuario_gestor_id);
CREATE INDEX idx_colaboradores_gestor ON colaboradores(usuario_gestor_id);
CREATE INDEX idx_tesis_autor ON tesis(autor_id);
CREATE INDEX idx_tesis_director ON tesis(director_id);
CREATE INDEX idx_tesis_coodirector ON tesis(coodirector_id);
CREATE INDEX idx_tesis_proyecto ON tesis(proyecto_id);
CREATE INDEX idx_integrantes_tesis_tesis ON integrantes_tesis(tesis_id);
CREATE INDEX idx_integrantes_tesis_integrante ON integrantes_tesis(integrante_id);
CREATE INDEX idx_lacis_proyecto ON lacis(proyecto_id);
CREATE INDEX idx_lacis_director ON lacis(director_id);
CREATE INDEX idx_lacis_comite ON lacis(comite_id);
CREATE INDEX idx_integrantes_lacis_lacis ON integrantes_lacis(lacis_id);
CREATE INDEX idx_integrantes_lacis_integrante ON integrantes_lacis(integrante_id);
CREATE INDEX idx_participante_externo_desarrollo ON participante_externo(desarrollo_id);
CREATE INDEX idx_proyecto_integrantes_proyecto ON proyecto_integrantes(proyecto_id);
CREATE INDEX idx_proyecto_integrantes_integrante ON proyecto_integrantes(integrante_id);
CREATE INDEX idx_desarrollo_integrantes_desarrollo ON desarrollo_integrantes(desarrollo_id);
CREATE INDEX idx_desarrollo_integrantes_integrante ON desarrollo_integrantes(integrante_id);
CREATE INDEX idx_proyecto_reconocimientos_proyecto ON proyecto_reconocimientos(proyecto_id);
CREATE INDEX idx_proyecto_reconocimientos_rec ON proyecto_reconocimientos(reconocimiento_id);
CREATE INDEX idx_desarrollo_reconocimientos_des ON desarrollo_reconocimientos(desarrollo_id);
CREATE INDEX idx_desarrollo_reconocimientos_rec ON desarrollo_reconocimientos(reconocimiento_id);

-- -----------------------------------------------------------------------------
-- DATOS POR DEFECTO
-- -----------------------------------------------------------------------------
-- Creamos el usuario administrador por defecto para poder iniciar sesión
INSERT INTO usuario_gestor (username, password_hash, email, rol)
VALUES ('admin', 'admin123', 'admin@lacis.com', 'ADMIN');

COMMIT;
