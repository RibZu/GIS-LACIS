package tesis

import (
	"strings"

	"PaginaSEG/internal/integrante"
)

type Tesis struct {
	ID                   int    `json:"id"`
	Titulo               string `json:"titulo"`
	Anio                 *int   `json:"anio"`
	Nivel                string `json:"nivel"`
	CarreraOrigen        string `json:"carrera_origen"`
	PalabrasClave        string `json:"palabras_clave"`
	Resumen              string `json:"resumen"`
	ArchivoPDF           string `json:"archivo_pdf"`
	AutorID              *int   `json:"autor_id,omitempty"`
	AutorHistorico       string `json:"autor_historico"`
	DirectorID           *int   `json:"director_id,omitempty"`
	DirectorHistorico    string `json:"director_historico"`
	CoodirectorID        *int   `json:"coodirector_id,omitempty"`
	CoodirectorHistorico string `json:"coodirector_historico"`
	ProyectoID           *int   `json:"proyecto_id,omitempty"`

	Autor              *integrante.Integrante `json:"autor,omitempty"`
	Director           *integrante.Integrante `json:"director,omitempty"`
	Coodirector        *integrante.Integrante `json:"coodirector,omitempty"`
	AutorNombre        string                 `json:"autor_nombre"`
	DirectorNombre     string                 `json:"director_nombre"`
	CoodirectorNombre  string                 `json:"coodirector_nombre"`
	IntegrantesIDs     []int                  `json:"integrantes_ids,omitempty"`
	AutoresSecundarios string                 `json:"autores_secundarios,omitempty"`
}

func (t Tesis) GetAutorDisplay() string {
	var partes []string
	if t.AutorNombre != "" {
		partes = append(partes, t.AutorNombre)
	} else if t.Autor != nil {
		partes = append(partes, t.Autor.Nombre+" "+t.Autor.Apellido)
	}
	if t.AutoresSecundarios != "" {
		partes = append(partes, t.AutoresSecundarios)
	}
	if t.AutorHistorico != "" {
		partes = append(partes, t.AutorHistorico)
	}
	if len(partes) > 0 {
		return strings.Join(partes, ", ")
	}
	return "No especificado"
}

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

func (t Tesis) GetAnioDisplay() int {
	if t.Anio != nil {
		return *t.Anio
	}
	return 0
}

func (t Tesis) GetAutorID() int {
	if t.AutorID != nil {
		return *t.AutorID
	}
	return 0
}

func (t Tesis) GetDirectorID() int {
	if t.DirectorID != nil {
		return *t.DirectorID
	}
	return 0
}

func (t Tesis) GetCoodirectorID() int {
	if t.CoodirectorID != nil {
		return *t.CoodirectorID
	}
	return 0
}

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
