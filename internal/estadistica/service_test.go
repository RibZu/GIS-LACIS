package estadistica

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type MockStorage struct {
	nombres     []string
	productos   []ProductoResumen
	integrantes []IntegranteResumen
	tesis       []TesisResumen
	proyectos   []ProyectoFilaEstadistica

	errNombres     error
	errProductos   error
	errIntegrantes error
	errTesis       error
	errProyectos   error
}

func (m *MockStorage) NombresParticipantesExternos() ([]string, error) {
	return m.nombres, m.errNombres
}
func (m *MockStorage) ResumenProductos() ([]ProductoResumen, error) {
	return m.productos, m.errProductos
}
func (m *MockStorage) ResumenIntegrantes() ([]IntegranteResumen, error) {
	return m.integrantes, m.errIntegrantes
}
func (m *MockStorage) ResumenTesis() ([]TesisResumen, error) {
	return m.tesis, m.errTesis
}
func (m *MockStorage) ResumenProyectos() ([]ProyectoFilaEstadistica, error) {
	return m.proyectos, m.errProyectos
}

func nuevoServicio(m *MockStorage) *Service {
	return NewService(m, zap.NewNop())
}

func anio(a int) *int { return &a }

func TestPorcentajes(t *testing.T) {
	casos := []struct {
		parte, total, p1, p2 int
	}{
		{0, 0, 0, 0},
		{1, 3, 33, 67},
		{2, 3, 67, 33},
		{3, 3, 100, 0},
		{0, 4, 0, 100},
		{8, 25, 32, 68},
	}
	for _, c := range casos {
		p1, p2 := porcentajes(c.parte, c.total)
		assert.Equal(t, c.p1, p1, "parte %d de %d", c.parte, c.total)
		assert.Equal(t, c.p2, p2, "parte %d de %d", c.parte, c.total)
	}
}

func TestSeriePorAnio(t *testing.T) {
	assert.Nil(t, seriePorAnio(nil))

	assert.Equal(t, []CantidadPorAnio{
		{Anio: 2019, Cantidad: 2, PorcentajeDelMaximo: 100},
		{Anio: 2020, Cantidad: 0, PorcentajeDelMaximo: 0},
		{Anio: 2021, Cantidad: 0, PorcentajeDelMaximo: 0},
		{Anio: 2022, Cantidad: 1, PorcentajeDelMaximo: 50},
	}, seriePorAnio([]int{2022, 2019, 2019}))

	assert.Equal(t, []CantidadPorAnio{{Anio: 2024, Cantidad: 3, PorcentajeDelMaximo: 100}},
		seriePorAnio([]int{2024, 2024, 2024}))
}

func TestEstadisticasIntegrantes_Vacio(t *testing.T) {
	assert.Equal(t, EstadisticasIntegrantes{}, estadisticasIntegrantes(nil))
}

func TestEstadisticasIntegrantes_CuentaComoLaLista(t *testing.T) {
	e := estadisticasIntegrantes([]IntegranteResumen{
		{Activo: true, PerteneceLacis: true},
		{Activo: true, PerteneceGrupoSoftware: true},
		{Activo: true, PerteneceLacis: true, PerteneceGrupoSoftware: true},
		{Activo: false, PerteneceLacis: true},
		{Activo: false},
	})
	assert.Equal(t, EstadisticasIntegrantes{Registrados: 5, Activos: 3, Inactivos: 2, Lacis: 3, Software: 2}, e)
	assert.Equal(t, e.Registrados, e.Activos+e.Inactivos)
}

func TestEsPosgrado_MismaReglaQueLaLista(t *testing.T) {
	casos := map[string]bool{
		"Doctorado":       true,
		"Maestría":        true,
		"Especialización": true,
		"MAESTRÍA EN X":   true,
		"doctorado":       true,
		"Posgrado":        true,
		"Grado":           false,
		"Licenciatura":    false,
		"":                false,
	}
	for nivel, esperado := range casos {
		assert.Equal(t, esperado, esPosgrado(nivel), "nivel %q", nivel)
	}
}

func TestEstadisticasTesis_Vacio(t *testing.T) {
	e := estadisticasTesis(nil)
	assert.Equal(t, EstadisticasTesis{}, e)
	assert.False(t, e.HayTesis)
	assert.Nil(t, e.PorAnio)
}

func TestEstadisticasTesis_CasoDeLaSpec(t *testing.T) {
	e := estadisticasTesis([]TesisResumen{
		{Anio: anio(2019), Nivel: "Doctorado", TienePDF: true},
		{Anio: anio(2019), Nivel: "Maestría", TienePDF: true},
		{Anio: anio(2022), Nivel: "Maestría"},
		{Anio: anio(2023), Nivel: "Grado", TienePDF: true},
		{Nivel: "Grado"},
		{Nivel: ""},
	})
	assert.Equal(t, 6, e.Registradas)
	assert.Equal(t, 3, e.Posgrado)
	assert.Equal(t, 3, e.GradoOtros)
	assert.True(t, e.HayTesis)
	assert.Equal(t, []CantidadPorAnio{
		{Anio: 2019, Cantidad: 2, PorcentajeDelMaximo: 100},
		{Anio: 2020, Cantidad: 0, PorcentajeDelMaximo: 0},
		{Anio: 2021, Cantidad: 0, PorcentajeDelMaximo: 0},
		{Anio: 2022, Cantidad: 1, PorcentajeDelMaximo: 50},
		{Anio: 2023, Cantidad: 1, PorcentajeDelMaximo: 50},
	}, e.PorAnio)
	assert.Equal(t, 3, e.ConPDF)
	assert.Equal(t, 3, e.SinPDF)
	assert.Equal(t, 50, e.PorcentajeConPDF)
	assert.Equal(t, 50, e.PorcentajeSinPDF)
}

func TestEstadisticasTesis_PDFEnLosExtremos(t *testing.T) {
	todas := estadisticasTesis([]TesisResumen{{TienePDF: true}, {TienePDF: true}})
	assert.Equal(t, 100, todas.PorcentajeConPDF)
	assert.Equal(t, 0, todas.PorcentajeSinPDF)

	ninguna := estadisticasTesis([]TesisResumen{{}, {}, {}})
	assert.Equal(t, 0, ninguna.PorcentajeConPDF)
	assert.Equal(t, 100, ninguna.PorcentajeSinPDF)
	assert.Equal(t, 3, ninguna.SinPDF)
}

func TestEstadisticasTesis_TodasSinAnio(t *testing.T) {
	e := estadisticasTesis([]TesisResumen{{Nivel: "Doctorado"}, {Nivel: "Grado", TienePDF: true}})
	assert.Nil(t, e.PorAnio)
	assert.Equal(t, 2, e.Registradas)
	assert.Equal(t, 1, e.Posgrado)
	assert.Equal(t, 1, e.GradoOtros)
	assert.Equal(t, 1, e.ConPDF)
}

func TestObtener_ArmaTodosLosModulos(t *testing.T) {
	m := &MockStorage{
		nombres: []string{"Ana Pérez", " ana  pérez ", "Luis Gómez", ""},
		productos: []ProductoResumen{
			{Anio: 2024, URL: "https://github.com/lacis/uno"},
			{Anio: 2022, URL: ""},
			{Anio: 2022, URL: ""},
		},
		integrantes: []IntegranteResumen{{Activo: true, PerteneceLacis: true}, {Activo: true, PerteneceLacis: true}, {Activo: false}},
		tesis:       []TesisResumen{{Anio: anio(2021), Nivel: "Maestría", TienePDF: true}},
	}
	est, err := nuevoServicio(m).Obtener()
	assert.NoError(t, err)

	assert.Equal(t, 2, est.IntegrantesActivosLacis)
	assert.Equal(t, 2, est.ParticipantesExternos)
	assert.Equal(t, 3, est.Productos)
	assert.Equal(t, 1, est.ProductosConRepositorio)
	assert.Equal(t, 2, est.ProductosConLicencia)
	assert.Equal(t, 33, est.PorcentajeConRepositorio)
	assert.Equal(t, 67, est.PorcentajeConLicencia)
	assert.True(t, est.HayProductos)
	assert.Equal(t, []CantidadPorAnio{
		{Anio: 2022, Cantidad: 2, PorcentajeDelMaximo: 100},
		{Anio: 2023, Cantidad: 0, PorcentajeDelMaximo: 0},
		{Anio: 2024, Cantidad: 1, PorcentajeDelMaximo: 50},
	}, est.PorAnio)

	assert.Equal(t, EstadisticasIntegrantes{Registrados: 3, Activos: 2, Inactivos: 1, Lacis: 2}, est.Integrantes)

	assert.Equal(t, 1, est.Tesis.Registradas)
	assert.Equal(t, 1, est.Tesis.Posgrado)
	assert.Equal(t, 100, est.Tesis.PorcentajeConPDF)
}

func TestObtener_BaseVacia(t *testing.T) {
	est, err := nuevoServicio(&MockStorage{}).Obtener()
	assert.NoError(t, err)
	assert.False(t, est.HayProductos)
	assert.Nil(t, est.PorAnio)
	assert.Equal(t, 0, est.PorcentajeConRepositorio)
	assert.Equal(t, EstadisticasIntegrantes{}, est.Integrantes)
	assert.Equal(t, EstadisticasTesis{}, est.Tesis)
}

func TestObtener_SiFallaCualquierLecturaNoDevuelveNumeros(t *testing.T) {
	falla := errors.New("sin conexión")
	casos := map[string]*MockStorage{
		"participantes externos": {errNombres: falla},
		"productos":              {errProductos: falla},
		"integrantes":            {errIntegrantes: falla},
		"tesis":                  {errTesis: falla},
		"proyectos":              {errProyectos: falla},
	}
	for nombre, m := range casos {
		est, err := nuevoServicio(m).Obtener()
		assert.Error(t, err, nombre)
		assert.ErrorIs(t, err, falla, nombre)
		assert.Nil(t, est, nombre)
	}
}

func TestEstadisticasProductos_ConjuntoConocido(t *testing.T) {
	m := &MockStorage{
		productos: []ProductoResumen{
			{Anio: 2019, URL: "https://github.com/lacis/uno"},
			{Anio: 2019, URL: ""},
			{Anio: 2021, URL: "   "},
		},
		nombres: []string{"Ana Pérez", "ana  pérez", " LUIS ", "Luis", ""},
		integrantes: []IntegranteResumen{
			{Activo: true, PerteneceLacis: true},
			{Activo: true, PerteneceLacis: false},
			{Activo: false, PerteneceLacis: true},
		},
	}
	est, err := nuevoServicio(m).Obtener()
	assert.NoError(t, err)

	assert.Equal(t, 3, est.Productos)
	assert.Equal(t, 1, est.ProductosConRepositorio)
	assert.Equal(t, 2, est.ProductosConLicencia)
	assert.Equal(t, est.Productos, est.ProductosConRepositorio+est.ProductosConLicencia)
	assert.Equal(t, 33, est.PorcentajeConRepositorio)
	assert.Equal(t, 67, est.PorcentajeConLicencia)
	assert.Equal(t, 100, est.PorcentajeConRepositorio+est.PorcentajeConLicencia)
	assert.True(t, est.HayProductos)
	assert.Equal(t, []CantidadPorAnio{
		{Anio: 2019, Cantidad: 2, PorcentajeDelMaximo: 100},
		{Anio: 2020, Cantidad: 0, PorcentajeDelMaximo: 0},
		{Anio: 2021, Cantidad: 1, PorcentajeDelMaximo: 50},
	}, est.PorAnio)
	assert.Equal(t, 2, est.ParticipantesExternos)
	assert.Equal(t, 1, est.IntegrantesActivosLacis)
}

func TestEstadisticasProductos_UnSoloProducto(t *testing.T) {
	est, err := nuevoServicio(&MockStorage{productos: []ProductoResumen{{Anio: 2024, URL: "x.com"}}}).Obtener()
	assert.NoError(t, err)
	assert.Equal(t, 1, est.Productos)
	assert.Equal(t, 1, est.ProductosConRepositorio)
	assert.Equal(t, 0, est.ProductosConLicencia)
	assert.Equal(t, 100, est.PorcentajeConRepositorio)
	assert.Equal(t, 0, est.PorcentajeConLicencia)
	assert.Equal(t, []CantidadPorAnio{{Anio: 2024, Cantidad: 1, PorcentajeDelMaximo: 100}}, est.PorAnio)
}

func TestEstadisticasProductos_AniosDesordenadosYHuecos(t *testing.T) {
	est, err := nuevoServicio(&MockStorage{productos: []ProductoResumen{
		{Anio: 2026}, {Anio: 2023}, {Anio: 2026}, {Anio: 2026}, {Anio: 2024}, {Anio: 2023},
	}}).Obtener()
	assert.NoError(t, err)
	assert.Equal(t, []CantidadPorAnio{
		{Anio: 2023, Cantidad: 2, PorcentajeDelMaximo: 66},
		{Anio: 2024, Cantidad: 1, PorcentajeDelMaximo: 33},
		{Anio: 2025, Cantidad: 0, PorcentajeDelMaximo: 0},
		{Anio: 2026, Cantidad: 3, PorcentajeDelMaximo: 100},
	}, est.PorAnio)
}

func TestContarPersonasDistintas(t *testing.T) {
	assert.Equal(t, 0, contarPersonasDistintas(nil))
	assert.Equal(t, 0, contarPersonasDistintas([]string{"", "   ", "	"}))
	assert.Equal(t, 1, contarPersonasDistintas([]string{"Ana Pérez", " ana   pérez ", "ANA PÉREZ"}))
	assert.Equal(t, 3, contarPersonasDistintas([]string{"Ana Pérez", "Ana Perez", "Luis"}))
}

func TestContarActivosLacis(t *testing.T) {
	assert.Equal(t, 0, contarActivosLacis(nil))
	assert.Equal(t, 2, contarActivosLacis([]IntegranteResumen{
		{Activo: true, PerteneceLacis: true},
		{Activo: true, PerteneceLacis: true, PerteneceGrupoSoftware: true},
		{Activo: true},
		{Activo: false, PerteneceLacis: true},
	}))
}
