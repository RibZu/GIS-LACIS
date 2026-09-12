package proyecto

type Proyecto struct {
	ID              int    `json:"id"`
	Titulo          string `json:"titulo"`
	Descripcion     string `json:"descripcion"`
	Enlace          string `json:"enlace"`
	EquipoHistorico string `json:"equipo_historico"`
	AnioInicio      int    `json:"anio_inicio"`
	AnioFin         int    `json:"anio_fin"`
	Activo          bool   `json:"activo"`
}

type UpdateFields struct {
	Titulo          *string `json:"titulo"`
	Descripcion     *string `json:"descripcion"`
	Enlace          *string `json:"enlace"`
	EquipoHistorico *string `json:"equipo_historico"`
	AnioInicio      *int    `json:"anio_inicio"`
	AnioFin         *int    `json:"anio_fin"`
}

type MiembroEquipo struct {
	IntegranteID  int    `json:"integrante_id"`
	Nombre        string `json:"nombre"`
	Apellido      string `json:"apellido"`
	RolEnProyecto string `json:"rol_en_proyecto"`
}
