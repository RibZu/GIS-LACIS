package integrante

import (
	"PaginaSEG/internal/rol"
)

type Integrante struct {
	ID                     int      `json:"id"`
	Nombre                 string   `json:"nombre"`
	Apellido               string   `json:"apellido"`
	Especializacion        string   `json:"especializacion"`
	Descripcion            string   `json:"descripcion"`
	Contacto               string   `json:"contacto"`
	ContactoLinkedin       string   `json:"contacto_linkedin"`
	Imagen                 string   `json:"imagen"`
	CV                     string   `json:"cv"`
	PerteneceLacis         bool     `json:"pertenece_lacis"`
	PerteneceGrupoSoftware bool     `json:"pertenece_grupo_software"`
	Activo                 bool     `json:"activo"`
	RolID                  int      `json:"rol_id"`
	RolLacisID             *int     `json:"rol_lacis_id,omitempty"`
	RolSoftwareID          *int     `json:"rol_software_id,omitempty"`
	Rol                    rol.Rol  `json:"rol,omitempty"`
	RolLacis               *rol.Rol `json:"rol_lacis,omitempty"`
	RolSoftware            *rol.Rol `json:"rol_software,omitempty"`
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
