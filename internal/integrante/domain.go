package integrante

import (
	"PaginaSEG/internal/rol"
)

type Integrante struct {
	ID              int     `json:"id"`
	Nombre          string  `json:"nombre"`
	Apellido        string  `json:"apellido"`
	Imagen          string  `json:"imagen"`
	CV              string  `json:"cv"`       // Array de strings en Postgres
	Contacto        string  `json:"contacto"` // Mail o datos de contacto
	Especializacion string  `json:"especializacion"`
	Descripcion     string  `json:"descripcion"`
	RolID           int     `json:"rol_id"`        // Clave foránea (fk_integrantes_rol)
	Rol             rol.Rol `json:"rol,omitempty"` // Relación: Objeto Rol completo al hacer JOIN
}

type UpdateFields struct {
	Nombre          *string `json:"nombre"`
	Apellido        *string `json:"apellido"`
	Imagen          *string `json:"imagen"`
	CV              *string `json:"cv"`       // Array de strings en Postgres
	Contacto        *string `json:"contacto"` // Mail o datos de contacto
	Especializacion *string `json:"especializacion"`
	Descripcion     *string `json:"descripcion"`
	RolID           *int    `json:"rol_id"`
}
