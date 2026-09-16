package desarrollo

type Desarrollo struct {
	ID          int    `json:"id"`
	Titulo      string `json:"titulo"`
	Anio        int    `json:"anio"`
	URL         string `json:"url"`
	Contacto    string `json:"contacto"`
	Descripcion string `json:"descripcion"`
}

type UpdateFields struct {
	Titulo      *string `json:"titulo"`
	Anio        *int    `json:"anio"`
	URL         *string `json:"url"`
	Contacto    *string `json:"contacto"`
	Descripcion *string `json:"descripcion"`
}
