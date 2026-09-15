package tesis

import (
	"PaginaSEG/internal/integrante"
)

// Tesis representa una tesis o trabajo final de grado/posgrado
type Tesis struct {
	ID                   int                    `json:"id"`
	Titulo               string                 `json:"titulo"`
	Anio                 *int                   `json:"anio"`
	Nivel                string                 `json:"nivel"`                 // Ej: "Doctorado", "Maestría", "Especialización", "Grado"
	CarreraOrigen        string                 `json:"carrera_origen"`        // Ej: "Doctorado en Ingeniería Informática", etc.
	PalabrasClave        string                 `json:"palabras_clave"`
	Resumen              string                 `json:"resumen"`
	ArchivoPDF           string                 `json:"archivo_pdf"`
	AutorID              *int                   `json:"autor_id,omitempty"`
	AutorHistorico       string                 `json:"autor_historico"`
	DirectorID           *int                   `json:"director_id,omitempty"`
	DirectorHistorico    string                 `json:"director_historico"`
	CoodirectorID        *int                   `json:"coodirector_id,omitempty"`
	CoodirectorHistorico string                 `json:"coodirector_historico"`
	ProyectoID           *int                   `json:"proyecto_id,omitempty"`

	// Objetos o datos enriquecidos para vista/JSON
	Autor                *integrante.Integrante `json:"autor,omitempty"`
	Director             *integrante.Integrante `json:"director,omitempty"`
	Coodirector          *integrante.Integrante `json:"coodirector,omitempty"`
	AutorNombre          string                 `json:"autor_nombre"`
	DirectorNombre       string                 `json:"director_nombre"`
	CoodirectorNombre    string                 `json:"coodirector_nombre"`
	IntegrantesIDs       []int                  `json:"integrantes_ids,omitempty"`
}

// GetAutorDisplay devuelve el nombre visible del autor (registrado o histórico)
func (t Tesis) GetAutorDisplay() string {
	if t.AutorNombre != "" {
		return t.AutorNombre
	}
	if t.AutorHistorico != "" {
		return t.AutorHistorico
	}
	if t.Autor != nil {
		return t.Autor.Nombre + " " + t.Autor.Apellido
	}
	return "No especificado"
}

// GetDirectorDisplay devuelve el nombre visible del director
func (t Tesis) GetDirectorDisplay() string {
	if t.DirectorNombre != "" {
		return t.DirectorNombre
	}
	if t.DirectorHistorico != "" {
		return t.DirectorHistorico
	}
	if t.Director != nil {
		return t.Director.Nombre + " " + t.Director.Apellido
	}
	return "No especificado"
}

// GetCoodirectorDisplay devuelve el nombre visible del codirector
func (t Tesis) GetCoodirectorDisplay() string {
	if t.CoodirectorNombre != "" {
		return t.CoodirectorNombre
	}
	if t.CoodirectorHistorico != "" {
		return t.CoodirectorHistorico
	}
	if t.Coodirector != nil {
		return t.Coodirector.Nombre + " " + t.Coodirector.Apellido
	}
	return ""
}

// GetAnioDisplay devuelve el año como string o "-" si no está especificado
func (t Tesis) GetAnioDisplay() int {
	if t.Anio != nil {
		return *t.Anio
	}
	return 0
}

// UpdateFields contiene los campos para actualizar una tesis
type UpdateFields struct {
	Titulo               *string `json:"titulo"`
	Anio                 *int    `json:"anio"`
	Nivel                *string `json:"nivel"`
	CarreraOrigen        *string `json:"carrera_origen"`
	PalabrasClave        *string `json:"palabras_clave"`
	Resumen              *string `json:"resumen"`
	ArchivoPDF           *string `json:"archivo_pdf"`
	AutorID              *int    `json:"autor_id"`
	AutorHistorico       *string `json:"autor_historico"`
	DirectorID           *int    `json:"director_id"`
	DirectorHistorico    *string `json:"director_historico"`
	CoodirectorID        *int    `json:"coodirector_id"`
	CoodirectorHistorico *string `json:"coodirector_historico"`
	ProyectoID           *int    `json:"proyecto_id"`
	IntegrantesIDs       []int   `json:"integrantes_ids"`
}
