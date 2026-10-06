package desarrollo

import "strings"

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

func TieneRepositorio(url string) bool {
	return strings.TrimSpace(url) != ""
}

func (d Desarrollo) TieneRepositorio() bool {
	return TieneRepositorio(d.URL)
}

func QuitarPrefijoEst(nombre string) string {
	palabras := strings.Fields(nombre)
	descartadas := 0
	for descartadas < len(palabras) {
		primera := strings.ToLower(palabras[descartadas])
		if primera != "est" && primera != "est." {
			break
		}
		descartadas++
	}
	if descartadas == 0 {
		return nombre
	}
	return strings.Join(palabras[descartadas:], " ")
}
