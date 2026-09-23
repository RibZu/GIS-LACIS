SET client_encoding = 'UTF8';
BEGIN;

CREATE TABLE rol (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL
);

CREATE TABLE integrante (
    id SERIAL PRIMARY KEY,
    rol_id INT REFERENCES rol(id) ON DELETE SET NULL,
    rol_lacis_id INT REFERENCES rol(id) ON DELETE SET NULL,
    rol_software_id INT REFERENCES rol(id) ON DELETE SET NULL,
    activo BOOLEAN NOT NULL DEFAULT TRUE,
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

CREATE TABLE usuario_gestor (
    id SERIAL PRIMARY KEY,
    integrante_id INT UNIQUE REFERENCES integrante(id) ON DELETE SET NULL,
    username VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    ultimo_acceso TIMESTAMP WITH TIME ZONE,
    rol VARCHAR(50) NOT NULL DEFAULT 'GESTOR',
    modulos VARCHAR(255) NOT NULL DEFAULT ''
);

CREATE TABLE configuracion_sitio (
    id SERIAL PRIMARY KEY,
    usuario_gestor_id INT REFERENCES usuario_gestor(id) ON DELETE SET NULL,
    telefono_footer VARCHAR(50),
    email_footer VARCHAR(255),
    direccion VARCHAR(255)
);

CREATE TABLE colaboradores (
    id BIGSERIAL PRIMARY KEY,
    usuario_gestor_id INT REFERENCES usuario_gestor(id) ON DELETE SET NULL,
    descripcion TEXT,
    logo_url VARCHAR(500),
    activo BOOLEAN DEFAULT TRUE
);

CREATE TABLE proyecto (
    id SERIAL PRIMARY KEY,
    titulo TEXT NOT NULL,
    descripcion TEXT,
    enlace TEXT,
    equipo_historico TEXT,
    anio_inicio INT,
    anio_fin INT,
    activo BOOLEAN DEFAULT TRUE
);

CREATE TABLE desarrollo (
    id SERIAL PRIMARY KEY,
    titulo VARCHAR(255) NOT NULL,
    anio INT NOT NULL,
    url VARCHAR(500),
    contacto VARCHAR(255),
    descripcion TEXT
);

CREATE TABLE reconocimientos (
    id BIGSERIAL PRIMARY KEY,
    titulo VARCHAR(255) NOT NULL,
    descripcion TEXT,
    activo BOOLEAN NOT NULL DEFAULT TRUE
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
    rol_en_proyecto VARCHAR(150),
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

CREATE INDEX idx_integrante_rol ON integrante(rol_id);
CREATE INDEX idx_integrante_rol_lacis ON integrante(rol_lacis_id);
CREATE INDEX idx_integrante_rol_software ON integrante(rol_software_id);
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

INSERT INTO rol (id, nombre) VALUES
    (1, 'Director / Co-Director'),
    (2, 'Investigador'),
    (3, 'Asesor Externo'),
    (4, 'Estudiante / Becario');
SELECT setval('rol_id_seq', (SELECT MAX(id) FROM rol));

INSERT INTO usuario_gestor (username, password_hash, email, rol)
VALUES ('admin', '$2a$10$0GSeTb6bSZBT1/kZbrf41.giDbz3GaU2ABeL20lew0NdX5dDw2cSK', 'admin@lacis.com', 'ADMIN');



INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (1, 1, NULL, 1, TRUE, 'Dr. Montejano', 'German', 'Ingenieria de Software', 'Doctor en Ingenieria de Software, Investigador categoría 1', 'german.a.montejano@gmail.com', '', '../../static/assets/img-GIS/IMGintegrantes/Foto German - German A Montejano.jpg', '../../static/assets/CVs-GIS/Bio Germán Montejano-CV resumido - German A Montejano.pdf', FALSE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (2, 1, NULL, 1, TRUE, 'Dr. Berón', 'Mario', 'Seguridad de Sistemas, Ingenieria Inversa', 'Doctor en Ciencias de la Computacion de la UNSL, Doutor en Ciencias da Computação , Universidade do Minho-BragaPortugal.', 'mberon@unsl.edu.ar', '', '../../static/assets/img-GIS/IMGintegrantes/mario beron.jpeg', '../../static/assets/CVs-GIS/MBeron.pdf', FALSE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (3, 1, 1, 1, TRUE, 'Dr. Salgado', 'Carlos Humberto', 'Metodologías de Desarrollo Seguro de Software', 'Calidad de Software y Testing. Docente e investigador en UNSL. Participación en equipos de I+D+I de la UNSL', 'csalgado@email.unsl.edu.ar', '', '../../static/assets/img-GIS/IMGintegrantes/salgado.jpeg', '../../static/assets/CVs-GIS/CV - HumbertoSalgado.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (4, 1, 1, 1, TRUE, 'Dr. Garis', 'Ana Gabriela', 'Protección contra Amenazas, Vulnerabilidades y Riesgos vinculados a los Artefactos de Software', 'Docente e investigadora en UNSL. Participación en equipos de I+D+I de la UNSL, del IIST (China), de la Universidad de Minho y la de Porto (Portugal).', '', 'https://www.linkedin.com/in/ana-garis-5738364', '../../static/assets/img-GIS/IMGintegrantes/foto_ana_garis 1 - Ana Garis.png', '../../static/assets/CVs-GIS/CV Ana Garis resumido short bio-2025 - Ana Garis.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (5, 2, 1, 2, TRUE, 'Mg. Baigorria', 'Lorena', 'Blockchain, Desarrollo Web, Tecnologías, Metodologías de Desarrollo, Modelos', 'Blockchain, Seguridad, Lenguajes de programación, Desarrollo Web', '', 'https://www.linkedin.com/in/lorena-baigorria/', '../../static/assets/img-GIS/IMGintegrantes/baigorria_lorena - Lorena Baigorria.jpeg', '../../static/assets/CVs-GIS/CV_Baigorria - Lorena Baigorria.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (6, 2, 1, 2, TRUE, 'Esp. Peralta', 'Mario Gabriel', 'Desarrollo Web, movil', 'Docente e investigador en UNSL. Participación en equipos de I+D+I de la UNSL', 'mperalta@unsl.edu.ar', '', '../../static/assets/img-GIS/IMGintegrantes/marioPeralta.jpeg', '../../static/assets/CVs-GIS/CV- MarioPeralta.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (7, 2, 1, 2, TRUE, 'Mg. Sanchez', 'Alberto', 'Calidad de procesos, productos y servicios', 'Magíster en Economía y Negocios, Licenciado en Ciencias de la Computación y Profesor en Computación. Consultor, auditor y docente universitario', '', 'https://www.linkedin.com/in/albertosanchez-gpfsoluciones', '../../static/assets/img-GIS/IMGintegrantes/FotoPerfil20231019_1 - Alberto Sanchez.jpg', '../../static/assets/CVs-GIS/CV_Resumen_Alberto_Sanchez - Alberto Sanchez.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (8, 2, 2, 2, TRUE, 'Dr. Abdelahad', 'Corina', 'Modelos, Transformación de Modelos', 'Profesor Adjunto e Investigadora en el Departamento de Informática de la FCFMyN (UNSL)', '', 'https://www.linkedin.com/in/corina-abdelahad-09a0482a1', '../../static/assets/img-GIS/IMGintegrantes/Foto Cori - Corina Abdelahad.jpg', '../../static/assets/CVs-GIS/CV Abdelahad Resumido - Corina Abdelahad.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (9, 2, 2, 2, TRUE, 'Esp. Albornoz', 'María Claudia', 'Interfaz Gráfica de Usuario', 'Soy M. Claudia Albornoz Docente e Investigadora de la UNSL.', 'albornoz@email.unsl.edu.ar', '', '../../static/assets/img-GIS/IMGintegrantes/Clau2025 - Maria Claudia Albornoz.jpg', '../../static/assets/CVs-GIS/CV-Resum-MCAloroz2025 - Maria Claudia Albornoz.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (10, 2, 2, 2, TRUE, 'Lic. Arroyuelo', 'Mónica del Valle', 'Programación - Análisis de Código - Seguridad', 'Docente e investigadora de la Universidad Nacional de San Luis. Sus áreas de interés incluyen la programación, el análisis de código, la seguridad de software y la protección frente a amenazas y vulnerabilidades en sistemas informáticos. Complementa su actividad académica con proyectos de innovación educativa y formación en tecnologías emergentes.', 'Mdarroyu@email.unsl.edu.ar', '', '../../static/assets/img-GIS/IMGintegrantes/Monica Arroyuelo.jpeg', '../../static/assets/CVs-GIS/CVArroyuelo - Monica Arroyuelo.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (11, 3, 3, 3, TRUE, 'Ing. Ayala', 'Camila', 'Frontend, Backend, Diseño, Análisis, Base de datos', 'Ingeniera en Informática - Desarrollo de soluciones escalables y sistemas multiplataforma con Metodologías Agiles', '', 'https://www.linkedin.com/in/camila-ayala-44288b3a1/', '../../static/assets/img-GIS/IMGintegrantes/foto_perfil - Camila Ayala.jpeg', '../../static/assets/CVs-GIS/Camila_Ayala_CV_Ingenieria_Informatica - Camila Ayala.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (12, 3, 3, 3, TRUE, 'Ing. Bascuñan', 'Agustin Ezequiel', 'Frontend, Backend, Base de Datos, Diseño, Análisis', 'Ingeniero en Informática por la Universidad Nacional de San Luis. Desarrollo de software (frontend y backend), gestión de bases de datos y trabajo colaborativo mediante metodologías ágiles (Scrum).', '', 'https://www.linkedin.com/in/agustin-bascunan/', '../../static/assets/img-GIS/IMGintegrantes/fotoCuadrada - Agustin Bascuñan.jpg', '../../static/assets/CVs-GIS/Agustin_Bascunan_CV_spanish_1PAGINA - Agustin Bascuñan.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (13, 2, 2, 2, TRUE, 'Lic. Bernardis', 'Edgardo', 'Ingeniería del Software, Seguridad.', 'Gestión de proyectos; actualmente investigo ofuscación de código en el área de Seguridad Informática.', '', 'https://www.linkedin.com/in/edgardo-bernardis-b58a8732a/', '../../static/assets/img-GIS/IMGintegrantes/eBernardisFoto - Edgardo Bernardis.jpg', '../../static/assets/CVs-GIS/ebernardis_CV RESUMIDO - Edgardo Bernardis.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (14, 2, 2, 2, TRUE, 'Mg. Bernardis', 'Hernán', 'Seguridad Informática, Ingeniería Reversa, Comprensión de Programas', 'Ciberseguridad y Seguridad Informática | Blue Team | Docente e Investigador', '', 'https://www.linkedin.com/in/hbernardis/', '../../static/assets/img-GIS/IMGintegrantes/HB - Hernán Bernardis.jpg', '../../static/assets/CVs-GIS/Hernán bernardis CV resumido - Hernán Bernardis.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (15, 3, 3, 3, TRUE, 'Mg. Doria', 'Maria Vanesa', 'Gestion de proyectos,gestión de talento humano, acceso abierto', 'Especialista en Ingeniería de Software (gestion de proyectos y gestión detalento humano), Docencia Universitaria y Acceso Abierto', 'vanesadoria@tecno.unca.edu.ar', '', '../../static/assets/img-GIS/IMGintegrantes/Foto de perfil Vanesa _fondo Blanco - Mg. Vanesa Doria.png', '../../static/assets/CVs-GIS//2025-María Vanesa Doria CV-Abreviado - Mg. Vanesa Doria.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (16, 3, 3, 3, TRUE, 'Mg. Flores', 'Carola Victoria', 'Ingeniería de requisitos, educación en ingeniería,acceso abierto', 'Especialista en Ingeniería de requisitos de software, Docencia Universitaria y Acceso abierto', 'carolaflores@tecno.unca.edu.ar', '', '../../static/assets/img-GIS/IMGintegrantes/2024-07-Foto de perfil Carola - Carola Flores.jpeg', '../../static/assets/CVs-GIS/2025-Carola Victoria Flores CV-Abreviado - Carola Flores.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (17, 4, NULL, 4, TRUE, 'Ing. Godoy', 'Inalen', 'Informática forense, análisis geoespacial, georreferenciación', 'Soy Técnico Universitario en Web y estudiante avanzada de Ingeniería en Informática. Actualmente me encuentro desarrollando mi Proyecto Integrador, enfocado en el análisis geoespacial aplicado a la informática forense.', '', 'www.linkedin.com/in/inalén-godoy-79238b13a', '../../static/assets/img-GIS/IMGintegrantes/Inalén Godoy.jpeg', '../../static/assets/CVs-GIS/curriculum_godoy - Inalén Godoy.pdf', FALSE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (18, 2, 2, 2, TRUE, 'Ing. Guiñazu', 'Joaquin', 'TOPSIS, MADM', 'Software Engineering Student | Software & Web Developer', '', 'https://www.linkedin.com/in/joaquin-gui%C3%B1azu-16b0332aa?utm_source=share_via&utm_content=profile&utm_medium=member_ios', '../../static/assets/img-GIS/IMGintegrantes/Joaquin Guiñazu.jpeg', '../../static/assets/CVs-GIS/CV - Guinazu Joaquin.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (19, 2, 2, 2, TRUE, 'Ing. Loyola', 'Maria Luciana', 'IA', 'Ingeniera en Informática con formación en desarrollo de software y bases de la computación.', '', 'https://www.linkedin.com/in/maria-luciana-loyola', '../../static/assets/img-GIS/IMGintegrantes/Maria Luciana Loyola.jpeg', '../../static/assets/CVs-GIS/CV_Py_Loyola - Maria Luciana Loyola.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (20, 2, 2, 2, TRUE, 'Dr. Miranda', 'Enrique', 'Informática Forense,Comprensión de Programas, Ingeniería Reversa', 'Doctor en Ingeniería Informática (beca CONICET), docente/investigador de grado yposgrado en UNSL en Ing. de Software e Informática Forense.', '', 'https://www.linkedin.com/in/eamiranda-ok', '../../static/assets/img-GIS/IMGintegrantes/foto perfil - Enrique Miranda.jpeg', '../../static/assets/CVs-GIS/CV-Miranda Enrique - FINAL - Enrique Miranda.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (21, 3, 3, 3, TRUE, 'Dr. Novillo Rangone', 'Gabriel', 'Inteligencia Artificial, Mineria de Datos, Ciencia de Datos', 'Secretario de Ciencia y Tecnica y Vinculacion Tecnologica de la UNVIME, Profesor Adj. Exc en UNVIME, Profesor de Posgrado en UNSL FCFQyM, y en UTN-FRC', '', 'https://www.linkedin.com/in/gabriel-novillo-rangone-331314237/', '../../static/assets/img-GIS/IMGintegrantes/novillo rangone.jpeg', '../../static/assets/CVs-GIS/CV Novillo Rangone.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (22, 3, 3, 3, TRUE, 'Dr.Sci Ochoa', 'Claudio', 'Backend architecture, scalability', 'Senior Principal Software Engineer', '', 'https://www.linkedin.com/in/claudioochoa/', '../../static/assets/img-GIS/IMGintegrantes/claudio - Claudio Ochoa.jpg', '../../static/assets/CVs-GIS/ClaudioOchoaResumeEn - Claudio Ochoa.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (23, 4, 4, 4, TRUE, 'Est. Olguin Quevedo', 'Gonzalo Gabriel', 'Backend, Frontend, IA', 'Soy un desarrollador con iniciativa, actitud positiva y habilidades interdisciplinarias, adaptable a diversos entornos. De rápido aprendizaje, mi enfoque carismático y actitud solidaria hacen que disfrute compartir conocimientos y experiencias con el equipo. Apasionado por enseñar lo que domino y aportar significativamente al éxito del proyecto asignado.', '', 'https://www.linkedin.com/in/gonzalo-olguin-quevedo-689b25212/', '../../static/assets/img-GIS/IMGintegrantes/Olguin.jpeg', '../../static/assets/CVs-GIS/Currículum Programador moderno turquesa - Gonzalo Olguin.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (24, 2, 2, 2, TRUE, 'Lic. Paez', 'Braian', 'Web, calidad', 'Soy Licenciado en Ciencia de la Computacion, me gusta poder encontrar soluciones a los problemas que se presenten mediante el desarrollo de software.', '', 'https://www.linkedin.com/in/braian-paez-9891b6306/', '../../static/assets/img-GIS/IMGintegrantes/Foto Braian Paez - Braian Paez.jpeg', '../../static/assets/CVs-GIS/CV_Braian_Paez - Braian Paez.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (25, 2, 2, 2, TRUE, 'Mg. Perez', 'Norma Beatriz', 'Ingeniería de la información - Gamificación - Seguridad - Base de Datos - Lenguajes', 'Mis investigaciones están focalizadas en el generar nuevas propuestas que permitan optimizar y obtener un software reutilizable y de calidad.', 'nbperez@email.unsl.edu.ar', '', '../../static/assets/img-GIS/IMGintegrantes/foto_Norma - Perez Norma Beatriz.png', '../../static/assets/CVs-GIS/CV_Resumido_NormaPerez2025 - Perez Norma Beatriz.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (26, 2, 2, 2, TRUE, 'Ing. Riera Bauer', 'Fabricio José', 'Web - IA', 'Estudiante avanzado, auxiliar de segunda categoría y apasionado nunca parar de aprender.', 'fjrbauer@email.unsl.edu.ar', '', '../../static/assets/img-GIS/IMGintegrantes/Foto perfil - Fabrizio José Riera Bauer.jpg', '../../static/assets/CVs-GIS/CV - Riera Bauer Fabrizio  - Fabrizio José Riera Bauer.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (27, 2, 2, 2, TRUE, 'Dr. Riesco', 'Daniel', 'Arquitectura de Software, Cloud Computing, Métodos Formales', 'Doctor en Informática por la Univ. de Vigo, con Maestría en Ing. del Conocimiento por la Univ. Politécnica de Madrid y Lic. en Cs de la Computación', 'driesco@unsl.edu.ar', '', '../../static/assets/img-GIS/IMGintegrantes/Foto Riesco - Daniel Riesco.jpg', '../../static/assets/CVs-GIS/CV_ORG_1_CV-ARGENTINO-COMPLETO_20163167973 2023 - Daniel Riesco.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (28, 2, 2, 2, TRUE, 'Ing. Salinas', 'Francisco Rodrigo', 'Calidad, Ingenieria de Software', 'Estudiante de quinto año de la carrera de Ingenieria en Informatica. Acreedor de la beca ARFITEC, actualmente desempeñando tareas de pasante en el laboratorio ETIS', '', 'https://www.linkedin.com/in/francisco-rodrigo-salinas-b753b6316/', '../../static/assets/img-GIS/IMGintegrantes/Francisco Rodrigo.jpeg', '../../static/assets/CVs-GIS/Curriculum vitae-3 - Francisco Rodrigo.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (29, 2, 2, 2, TRUE, 'Mg. Sanchez', 'Alejandro', 'Métodos formales', 'Docente, investigador y practicante de la ingeniería de software.  Fue miembro del Center for Electronic Governance, centro que dio lugar al actual Open Unit on Policy-Driven Electronic Governance (UNU-EGOV) en el área de interoperabilidad semántica para gobierno electrónico. Previo a este período, trabajó en la industria de software ocupando roles de programador, arquitecto de software, y líder de proyecto. Es autor de 35 publicaciones: 3 artículos de revista, 3 capítulos de libro, 29 papers en conferencias internacionales y nacionales.', '', 'https://www.linkedin.com/in/aljsanchez/', '../../static/assets/img-GIS/IMGintegrantes/Alejandro Sanchez.jpeg', '../../static/assets/CVs-GIS/CV - Alejandro Sanchez.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (30, 3, 3, 3, TRUE, 'Ing. Velez', 'Priscila Zoe', 'Frontend, backend, Base de datos, microservicios, Reingeniería, Diseño, Análisis', 'Enfoque en desarrollo Full Stack y bases de datos. Experiencia en aplicaciones web y mobile utilizando React, Next.js, NestJS y React Native, trabajando bajo metodologías ágiles (SCRUM) y abarcando todo el ciclo de desarrollo', '', 'https://www.linkedin.com/in/priscila-zoe-velez-941635276', '../../static/assets/img-GIS/IMGintegrantes/fotodecente - Zoé Velez.jpeg', '../../static/assets/CVs-GIS/Velez Priscila- CV (resumido) - Zoé Velez.pdf', TRUE, TRUE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (31, 2, 2, NULL, TRUE, 'Ing. Baglioni', 'Valentina', 'Calidad, IA', 'Ingeniera en Informática, docente universitaria e investigadora.', '', 'https://www.linkedin.com/in/valentina-baglioni-26b51b255/', '../../static/assets/img-GIS/IMGintegrantes/Valen Bagli.jpeg', '../../static/assets/CVs-GIS/CV_Valentina_Baglioni_Resumen_1_Pagina - Valen Bagli.pdf', TRUE, FALSE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (32, 4, 4, NULL, TRUE, 'Est. Farias', 'Ismael', 'Desarrollo Web, BackEnd - FrontEnd', 'Desarrollador web en formación, estudiante avanzado de Tecnicatura en Desarrollo Web, con experiencia en proyectos reales y académicos utilizando PHP, Java y MySQL. Orientado al desarrollo Back-End, bases de datos y metodologías ágiles', '', 'https://www.linkedin.com/in/ismael-alexander-farias-7a64b520b/', '../../static/assets/img-GIS/IMGintegrantes/Ismael Farias.jpg', '../../static/assets/CVs-GIS/IsmaelAlexander_Farias_CV.pdf', TRUE, FALSE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (33, 3, 3, NULL, TRUE, 'Lic. Figueroa', 'Agustín', 'IA', 'Soy Licenciado en Ciencias de la Computación, donde adquirí una solida base en este campo; formándome en Informática Forense y actualmente en IA', '', 'https://www.linkedin.com/in/agustin-figueroa-5303602a2', '../../static/assets/img-GIS/IMGintegrantes/Foto-Perfil - Agustín Figueroa.jpg', '../../static/assets/CVs-GIS/CV - Figueroa - Agustín Figueroa.pdf', TRUE, FALSE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (34, 3, 3, NULL, TRUE, 'Ing. Orellano Zalazar', 'Ignacio', 'Desarrollo Web', 'Ingeniero en informática con experiencia en el desarrollo web', '', 'https://www.linkedin.com/in/ignacio-orellano', '../../static/assets/img-GIS/IMGintegrantes/IMG_3075 - Ignacio Orellano.jpeg', '../../static/assets/CVs-GIS/Curriculum_Ignacio_Orellano.pdf - Ignacio Orellano.pdf', TRUE, FALSE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (35, 2, 2, NULL, TRUE, 'Dr. Vilallonga', 'Gabriel Domingo', 'Gestión de Procesos', 'Doctor en Ingeniería de Software. Profesor Titular e investigador en procesos, transformación digital y desarrollos con IAG.', 'gvilallo@tecno.unca.edu.ar', '', '../../static/assets/img-GIS/IMGintegrantes/Gabriel Vilallonga - Gabriel Domingo Vilallonga.jpeg', '../../static/assets/CVs-GIS/CV Resumido - Vilallonga, Gabriel D - Gabriel Domingo Vilallonga.pdf', TRUE, FALSE);
INSERT INTO integrante (id, rol_id, rol_lacis_id, rol_software_id, activo, nombre, apellido, titulo_especializacion, descripcion, contacto_mail, contacto_linkedin, imagen_url, cv_url, pertenece_lacis, pertenece_grupo_software) VALUES (36, 4, 4, NULL, TRUE, 'Est. Riberi Zunino', 'Simon Martin', 'Desarrollo Web (Backend y Frontend), Bases de Datos', 'Estudiante avanzado de la Tecnicatura Universitaria en Web (UNSL).', '', 'https://www.linkedin.com/in/simon-riberi-5a28bb238/', '', '../../static/assets/CVs-GIS/Simon_Riberi_CV.pdf', TRUE, FALSE);
SELECT setval('integrante_id_seq', (SELECT MAX(id) FROM integrante));

INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (1, 'Aplicación Web Autogestionable', 2017, NULL, 'ingenieriainformatica@gmail.com', 'Sistema web para PyMes, administracion de productos a traves de documentos word');
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (2, 'SL Parking: Gestión de estacionamiento on-line', 2017, NULL, 'ingenieriainformatica@gmail.com', 'Aplicacion para el pago de estacionamiento medido en ciudades inteligentes');
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (3, 'Draw Your Fish Bone', 2019, NULL, NULL, NULL);
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (4, 'Proyecto Paraísos', 2022, NULL, 'ingenieriainformatica@gmail.com', 'Sistema de relevamiento de expansion de hongos para arboles de la ciudad de San LUIS');
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (5, 'Sistema Administrativo y Financiero para Unidades Académicas (SAFUA)', 2022, NULL, 'ingenieriainformatica@gmail.com', 'Gestio asministrativa FCFMyN');
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (6, 'Diseño, construcción e implantación en entornos reales para uso industrial de una aplicación móvil para envío no convencional de paquetería bajo demanda', 2023, NULL, NULL, NULL);
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (7, 'Sistema de Gestión de Auditorías para el Tribunal de Cuentas de la Municipalidad de San Luis', 2023, NULL, 'ingenieriainformatica@gmail.com', 'Sistema de Gestión de Auditorías de gastos para el Tribunal de Cuentas de la Municipalidad de San Luis');
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (8, 'Sistema de Gestión de Espacios Académicos para la UNSL', 2023, 'http://sistemadeaulas.unsl.edu.ar/', 'ingenieriainformatica@gmail.com', 'Gestion y distribución de espacios aulicos de la UNSL');
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (9, 'Sistema de Pasantías y Prácticas Profesionales', 2023, NULL, NULL, NULL);
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (10, 'HEvaLog: Una Herramienta de Evaluación de Sistemas Complejos basada en Lógica, para la toma de Decisiones', 2024, NULL, NULL, NULL);
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (11, 'Sistema de Certificación Digital de Títulos Universitarios en Blockchain (UNSL Cert)', 2024, NULL, NULL, NULL);
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (12, 'Un Sistema Informático de Expendio de Órdenes para DOSPU', 2024, NULL, NULL, NULL);
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (13, 'Mi Gym: App de gestión de gimnasios y seguimiento integral del progreso y salud de los clientes', 2025, NULL, NULL, NULL);
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (14, 'PUPA: Asistente Conversacional basado en Inteligencia Artificial para la Prevención de Apuestas Online', 2025, 'https://prevencionenadicciones.unsl.edu.ar/', 'ingenieriainformatica@gmail.com', 'Asistente Virtual basado en IAG para prevenir ludopatias, apuestas online');
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (15, 'Parametrizando un Sistema Informático para Obra Social Universitaria con un Motor de Reglas', 2025, NULL, NULL, NULL);
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (16, 'TraiLead: Una herramienta para optimizar los procesos de inducción', 2025, NULL, NULL, NULL);
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (17, 'AVI: Asistente Virtual Conversacional basado en Inteligencia Artificial para Ingresantes de la FCFMyN', 2026, 'https://avi-unsl.vercel.app/', 'ingenieriainformatica@gmail.com', 'Asistente Vistual basado en IAG para guiar a los ingresantes de la FCFMYN de la UNSL en el ingreso a las carreras');
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (18, 'Diseño y desarrollo de un Asistente Virtual Web para la optimización de la trayectoria académica y el fomento del aprendizaje colaborativo en el Departamento de Informática de la UNSL', 2026, 'https://asistente-virtual-unsl.vercel.app', 'ingenieriainformatica@gmail.com', 'Diseño y desarrollo de un Asistente Virtual Web para la optimización de la trayectoria académica y el fomento del aprendizaje colaborativo en el Departamento de Informática de la UNSL');
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (19, 'Optimización de procesos mediante una arquitectura orientada a eventos', 2026, NULL, NULL, NULL);
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (20, 'SUAP: Sistema Universitario de Acompañamiento y Permanencia.', 2026, 'http://suap.dirinfo.unsl.edu.ar/', 'ingenieriainformatica@gmail.com', 'Gestion de calificaciones parciales y laboratorios para estudiantes universitario. Toma de asistencia con geoposicionamiento, Alerta y contactos de emergencia');
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (21, 'Sistema Informático para el Análisis de Sábanas de Telefonía', 2026, NULL, 'ingenieriainformatica@gmail.com', 'Sistema Informático para el Análisis de Sábanas de Telefonía para Poder Judicial de la Provincia de San Luis, VIsualizacion de comunicaciones en casos judiciales');
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (22, 'Sistema de Gestión Integral y Business Intelligence como Estrategia de Optimización del Capital de Trabajo', 2026, NULL, 'ingenieriainformatica@gmail.com', 'Desarrollo de un Sistema de Gestión Integral (ERP) y Business Intelligence para la pyme');
INSERT INTO desarrollo (id, titulo, anio, url, contacto, descripcion) VALUES (23, 'UniMaps: Una plataforma web interactiva para la localización efectiva de espacios de la UNSL', 2026, 'https://unimaps-pi.vercel.app/', 'ingenieriainformatica@gmail.com', 'UniMaps: Una plataforma web interactiva para la localización efectiva de espacios de la UNSL');
SELECT setval('desarrollo_id_seq', (SELECT MAX(id) FROM desarrollo));

INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (1, 1, 5);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (2, 1, 6);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (3, 2, 3);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (4, 2, 5);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (5, 3, 3);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (6, 3, 6);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (7, 4, 3);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (8, 4, 6);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (9, 5, 5);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (10, 6, 1);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (11, 6, 8);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (12, 7, 3);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (13, 7, 5);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (14, 7, 34);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (15, 8, 3);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (16, 8, 5);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (17, 8, 31);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (18, 9, 5);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (19, 9, 6);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (20, 10, 3);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (21, 10, 6);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (22, 11, 4);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (23, 12, 3);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (24, 12, 29);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (25, 13, 6);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (26, 13, 8);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (27, 14, 4);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (28, 14, 5);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (29, 14, 19);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (30, 14, 26);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (31, 15, 3);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (32, 15, 29);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (33, 16, 5);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (34, 16, 6);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (35, 17, 4);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (36, 17, 5);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (37, 18, 3);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (38, 18, 8);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (39, 19, 5);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (40, 19, 27);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (41, 20, 5);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (42, 20, 6);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (43, 20, 11);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (44, 20, 12);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (45, 20, 30);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (46, 21, 4);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (47, 21, 17);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (48, 21, 20);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (49, 22, 2);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (50, 22, 3);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (51, 23, 5);
INSERT INTO desarrollo_integrantes (id, desarrollo_id, integrante_id) VALUES (52, 23, 8);
SELECT setval('desarrollo_integrantes_id_seq', (SELECT MAX(id) FROM desarrollo_integrantes));

INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (1, 1, 'Marcos Picco');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (2, 2, 'Juan Pablo Imperiale');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (3, 3, 'Abatedaga Biole Nicolas');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (4, 4, 'Lautaro Emanuel');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (5, 5, 'Silvano Di Marco');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (6, 5, 'Dr. Cristian Tissera');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (7, 6, 'Santiago Jorge Bordarampe');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (8, 6, 'Sebastián Luis Riccardo');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (9, 7, 'Mathias Ezequiel Schuvab Arce');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (10, 7, 'Gabriel Petrino');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (11, 7, 'Tribunal de Cuentas Municipalidad de San Luis');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (12, 8, 'Rosa Lorenzo');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (13, 9, 'Juan Antonio Zini');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (14, 10, 'Luciano Martin Gurruchaga');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (15, 11, 'Vergés Federico Agustín');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (16, 11, 'Pablo Cristian Tissera');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (17, 12, 'Agustin Vela Luengo');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (18, 13, 'Maximiliano Hernán Pellegrino Dario');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (19, 14, 'Eliana Gonzalez');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (20, 15, 'Iván Brocas');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (21, 16, 'Agustín Ferrari');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (22, 17, 'Marcos Gelves');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (23, 17, 'Santiago Delias');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (24, 18, 'Juan Manuel Sanchez');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (25, 19, 'De Miguel Nicolás Agustín');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (26, 21, 'Departamento de Delitos Complejos del Poder Judicial de la Provincia de San LUIS');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (27, 22, 'Diaz Gustavo Emiliano');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (28, 22, 'Vilurón Gustavo Nicolás');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (29, 23, 'Lautaro Soria');
INSERT INTO participante_externo (id, desarrollo_id, nombre) VALUES (30, 23, 'Lucia Morandini');
SELECT setval('participante_externo_id_seq', (SELECT MAX(id) FROM participante_externo));

COMMIT;
