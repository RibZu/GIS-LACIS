package estadistica

type ProductoResumen struct {
	Anio int
	URL  string
}

type IntegranteResumen struct {
	Activo                 bool
	PerteneceLacis         bool
	PerteneceGrupoSoftware bool
}

type TesisResumen struct {
	Anio     *int
	Nivel    string
	TienePDF bool
}

type CantidadPorAnio struct {
	Anio                int
	Cantidad            int
	PorcentajeDelMaximo int
}

type EstadisticasIntegrantes struct {
	Registrados int
	Activos     int
	Inactivos   int
	Lacis       int
	Software    int
}

type EstadisticasTesis struct {
	Registradas      int
	Posgrado         int
	GradoOtros       int
	PorAnio          []CantidadPorAnio
	HayTesis         bool
	ConPDF           int
	SinPDF           int
	PorcentajeConPDF int
	PorcentajeSinPDF int
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
	Integrantes              EstadisticasIntegrantes
	Tesis                    EstadisticasTesis
}
