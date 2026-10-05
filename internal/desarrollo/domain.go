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

// TieneRepositorio indica si el producto tiene un repositorio publicado, que es el campo
// "Enlace / URL". "Con licencia" es lo contrario: un producto sin enlace (o con solo espacios).
func TieneRepositorio(url string) bool {
	return strings.TrimSpace(url) != ""
}

// TieneRepositorio es la misma regla sobre el producto, para usarla desde las plantillas
// ({{ if $d.TieneRepositorio }}). "Con licencia" es !TieneRepositorio.
func (d Desarrollo) TieneRepositorio() bool {
	return TieneRepositorio(d.URL)
}

// QuitarPrefijoEst devuelve el nombre sin la palabra "Est"/"Est." del principio (sin distinguir
// mayúsculas), repetida o no. Si no había prefijo devuelve el nombre exactamente igual; si lo
// había, devuelve las palabras restantes separadas por un espacio (puede ser ""). Se usa solo
// para limpiar datos ya cargados, no al guardar desde el panel, donde cada nombre se respeta.
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
