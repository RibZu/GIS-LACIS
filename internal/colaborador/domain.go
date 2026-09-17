package colaborador

type Colaborador struct {
	ID          int64  `json:"id"`
	Descripcion string `json:"descripcion"`
	LogoURL     string `json:"logo_url"`
	Activo      bool   `json:"activo"`
}

type UpdateFields struct {
	Descripcion *string `json:"descripcion"`
	LogoURL     *string `json:"logo_url"`
	Activo      *bool   `json:"activo"`
}
