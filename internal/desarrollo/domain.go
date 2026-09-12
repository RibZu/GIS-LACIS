package desarrollo

type Desarrollo struct {
	ID          int    `json:"id"`
	Titulo      string `json:"titulo"`
	Anio        int    `json:"anio"`
	URL         string `json:"url"`         // url en Postgres
	Contacto    string `json:"contacto"`    // contacto en Postgres
	Descripcion string `json:"descripcion"` // descripcion en Postgres
}

type UpdateFields struct {
	Titulo      *string `json:"titulo"`
	Anio        *int    `json:"anio"`
	URL         *string `json:"url"`
	Contacto    *string `json:"contacto"`
	Descripcion *string `json:"descripcion"`
}
