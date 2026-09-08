package reconocimiento

type Reconocimiento struct {
	ID          int64  `json:"id"`
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	Activo      bool   `json:"activo"`
}

type UpdateFields struct {
	Titulo      *string `json:"titulo"`
	Descripcion *string `json:"descripcion"`
	Activo      *bool   `json:"activo"`
}
