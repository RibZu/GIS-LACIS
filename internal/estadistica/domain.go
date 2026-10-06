// Package estadistica calcula indicadores de solo lectura sobre los productos de software, los
// integrantes, los participantes externos, las tesis y los proyectos. Es un reporte que cruza varias tablas
// y no tiene altas, bajas ni ediciones, por eso no define UpdateFields.
package estadistica

// ProductoResumen tiene solo lo que usan las estadísticas de productos: el año y el enlace.
type ProductoResumen struct {
	Anio int
	URL  string
}

// IntegranteResumen es una fila de integrante con los tres datos sí/no que cuentan las tarjetas.
type IntegranteResumen struct {
	Activo                 bool
	PerteneceLacis         bool
	PerteneceGrupoSoftware bool
}

// TesisResumen es una fila de tesis. Anio es nil cuando la tesis no tiene año cargado.
type TesisResumen struct {
	Anio     *int
	Nivel    string
	TienePDF bool
}

// CantidadPorAnio es una fila de "Productos por año" o "Tesis por año". PorcentajeDelMaximo es el
// ancho de la barra: 100 para el año con más registros y proporcional para los demás.
type CantidadPorAnio struct {
	Anio                int
	Cantidad            int
	PorcentajeDelMaximo int
}

// EstadisticasIntegrantes son las tarjetas que antes calculaba lista.js en la lista de integrantes,
// con la misma cuenta.
type EstadisticasIntegrantes struct {
	Registrados int
	Activos     int
	Inactivos   int
	Lacis       int
	Software    int
}

// EstadisticasTesis junta las tarjetas que antes calculaba listaTesis.js en la lista de tesis
// (Registradas, Posgrado, GradoOtros) con los gráficos de tesis por año y con PDF vs. sin PDF.
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

// RolCantidad representa la cantidad de integrantes que cumplen un rol específico en un proyecto o a nivel global.
type RolCantidad struct {
	Rol      string
	Cantidad int
}

// ProyectoEstadistica representa un proyecto con su resumen de integrantes por rol.
type ProyectoEstadistica struct {
	ID               int
	Titulo           string
	AnioInicio       int
	AnioFin          int
	Activo           bool
	EquipoHistorico  string
	TotalIntegrantes int
	TieneIntegrantes bool
	Roles            []RolCantidad
}

// EstadisticasProyectos consolida los indicadores de proyectos y el detalle individual.
type EstadisticasProyectos struct {
	TotalProyectos          int
	ProyectosConIntegrantes int
	ProyectosSinDatos       int
	TotalParticipaciones    int
	HayProyectos            bool
	RolesGlobal             []RolCantidad
	Lista                   []ProyectoEstadistica
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
	Proyectos                EstadisticasProyectos
}
