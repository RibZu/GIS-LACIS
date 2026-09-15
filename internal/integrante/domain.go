package integrante

import (
	"PaginaSEG/internal/rol"
)

type Integrante struct {
	ID                     int      `json:"id"`
	Nombre                 string   `json:"nombre"`
	Apellido               string   `json:"apellido"`
	Especializacion        string   `json:"especializacion"`        // titulo_especializacion en Postgres
	Descripcion            string   `json:"descripcion"`            // descripcion en Postgres
	Contacto               string   `json:"contacto"`               // contacto_mail en Postgres
	ContactoLinkedin       string   `json:"contacto_linkedin"`      // contacto_linkedin en Postgres
	Imagen                 string   `json:"imagen"`                 // imagen_url en Postgres
	CV                     string   `json:"cv"`                     // cv_url en Postgres
	PerteneceLacis         bool     `json:"pertenece_lacis"`        // pertenece_lacis en Postgres
	PerteneceGrupoSoftware bool     `json:"pertenece_grupo_software"`// pertenece_grupo_software en Postgres
	Activo                 bool     `json:"activo"`                 // activo (baja lógica) en Postgres
	RolID                  int      `json:"rol_id"`                 // rol_id legado/general en Postgres
	RolLacisID             *int     `json:"rol_lacis_id,omitempty"` // rol_lacis_id en Postgres
	RolSoftwareID          *int     `json:"rol_software_id,omitempty"` // rol_software_id en Postgres
	Rol                    rol.Rol  `json:"rol,omitempty"`          // Relación: Objeto Rol general
	RolLacis               *rol.Rol `json:"rol_lacis,omitempty"`    // Relación: Objeto Rol en LaCIS
	RolSoftware            *rol.Rol `json:"rol_software,omitempty"` // Relación: Objeto Rol en Grupo Software
}

type UpdateFields struct {
	Nombre                 *string `json:"nombre"`
	Apellido               *string `json:"apellido"`
	Especializacion        *string `json:"especializacion"`
	Descripcion            *string `json:"descripcion"`
	Contacto               *string `json:"contacto"`
	ContactoLinkedin       *string `json:"contacto_linkedin"`
	Imagen                 *string `json:"imagen"`
	CV                     *string `json:"cv"`
	PerteneceLacis         *bool   `json:"pertenece_lacis"`
	PerteneceGrupoSoftware *bool   `json:"pertenece_grupo_software"`
	Activo                 *bool   `json:"activo"`
	RolID                  *int    `json:"rol_id"`
	RolLacisID             *int    `json:"rol_lacis_id"`
	RolSoftwareID          *int    `json:"rol_software_id"`
}

func (i Integrante) GetRolLacisID() int {
	if i.RolLacisID != nil {
		return *i.RolLacisID
	}
	return 0
}

func (i Integrante) GetRolSoftwareID() int {
	if i.RolSoftwareID != nil {
		return *i.RolSoftwareID
	}
	return 0
}
