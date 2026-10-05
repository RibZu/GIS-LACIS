// Package estadistica calcula indicadores de solo lectura sobre los productos de software, los
// integrantes y los participantes externos. Es un reporte que cruza varias tablas y no tiene
// altas, bajas ni ediciones, por eso no define UpdateFields.
package estadistica

type ProductoResumen struct {
	ID                 int
	Titulo             string
	Anio               int
	URL                string
	Contacto           string
	Descripcion        string
	TieneParticipantes bool
}

// CantidadPorAnio es una fila de "Productos por año". PorcentajeDelMaximo es el ancho de la barra:
// 100 para el año con más productos y proporcional para los demás.
type CantidadPorAnio struct {
	Anio                int
	Cantidad            int
	PorcentajeDelMaximo int
}

// ProductoIncompleto es un producto al que le falta información. Faltantes usa los nombres
// "descripción", "contacto" y "participantes", en ese orden. El repositorio no cuenta: un
// producto sin enlace es un producto con licencia, no un producto incompleto.
type ProductoIncompleto struct {
	ID        int
	Titulo    string
	Anio      int
	Faltantes []string
}

type Estadisticas struct {
	IntegrantesActivosLacis  int
	ParticipantesExternos    int
	Productos                int
	ProductosConRepositorio  int
	ProductosConLicencia     int
	PorcentajeConRepositorio int
	PorcentajeConLicencia    int
	HayProductos             bool
	PorAnio                  []CantidadPorAnio
	Incompletos              []ProductoIncompleto
}
