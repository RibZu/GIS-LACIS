package estadistica

import (
	"fmt"
	"strings"

	"PaginaSEG/internal/desarrollo"

	"go.uber.org/zap"
)

type Service struct {
	storage Storage
	logger  *zap.Logger
}

func NewService(s Storage, l *zap.Logger) *Service {
	return &Service{storage: s, logger: l}
}

// Obtener calcula las estadísticas en el momento. Si cualquier lectura falla devuelve un error y
// ningún número, para no mostrar resultados parciales.
func (s *Service) Obtener() (*Estadisticas, error) {
	integrantes, err := s.storage.ContarIntegrantesActivosLacis()
	if err != nil {
		return nil, s.fallo("integrantes activos", err)
	}
	nombres, err := s.storage.NombresParticipantesExternos()
	if err != nil {
		return nil, s.fallo("participantes externos", err)
	}
	productos, err := s.storage.ResumenProductos()
	if err != nil {
		return nil, s.fallo("productos", err)
	}

	est := &Estadisticas{
		IntegrantesActivosLacis: integrantes,
		ParticipantesExternos:   contarPersonasDistintas(nombres),
		Productos:               len(productos),
	}
	for _, p := range productos {
		if desarrollo.TieneRepositorio(p.URL) {
			est.ProductosConRepositorio++
		}
	}
	est.ProductosConLicencia = est.Productos - est.ProductosConRepositorio
	est.HayProductos = est.Productos > 0
	if est.HayProductos {
		// Redondeo al entero más cercano; el segundo porcentaje es el resto, así suman 100.
		est.PorcentajeConRepositorio = (200*est.ProductosConRepositorio + est.Productos) / (2 * est.Productos)
		est.PorcentajeConLicencia = 100 - est.PorcentajeConRepositorio
	}
	est.PorAnio = productosPorAnio(productos)
	est.Incompletos = productosIncompletos(productos)
	return est, nil
}

func (s *Service) fallo(que string, err error) error {
	s.logger.Error("No se pudieron obtener las estadísticas", zap.String("lectura", que), zap.Error(err))
	return fmt.Errorf("estadisticas: %w", err)
}

// contarPersonasDistintas cuenta los nombres distintos ignorando mayúsculas y espacios de más, así
// la misma persona cargada en varios productos cuenta una sola vez. Los nombres vacíos no cuentan.
func contarPersonasDistintas(nombres []string) int {
	vistos := make(map[string]struct{}, len(nombres))
	for _, n := range nombres {
		clave := strings.ToLower(strings.Join(strings.Fields(n), " "))
		if clave != "" {
			vistos[clave] = struct{}{}
		}
	}
	return len(vistos)
}

// productosPorAnio arma una fila por cada año entre el primero y el último con productos, incluso
// los años sin ninguno, para que la serie no salte años.
func productosPorAnio(productos []ProductoResumen) []CantidadPorAnio {
	if len(productos) == 0 {
		return nil
	}
	cantidades := make(map[int]int)
	minAnio, maxAnio := productos[0].Anio, productos[0].Anio
	for _, p := range productos {
		cantidades[p.Anio]++
		minAnio = min(minAnio, p.Anio)
		maxAnio = max(maxAnio, p.Anio)
	}
	maxCantidad := 0
	for _, c := range cantidades {
		maxCantidad = max(maxCantidad, c)
	}

	serie := make([]CantidadPorAnio, 0, maxAnio-minAnio+1)
	for anio := minAnio; anio <= maxAnio; anio++ {
		serie = append(serie, CantidadPorAnio{
			Anio:                anio,
			Cantidad:            cantidades[anio],
			PorcentajeDelMaximo: cantidades[anio] * 100 / maxCantidad,
		})
	}
	return serie
}

// productosIncompletos devuelve, en el orden recibido, los productos a los que les falta la
// descripción, el contacto o los participantes. El repositorio no cuenta.
func productosIncompletos(productos []ProductoResumen) []ProductoIncompleto {
	var incompletos []ProductoIncompleto
	for _, p := range productos {
		var faltantes []string
		if strings.TrimSpace(p.Descripcion) == "" {
			faltantes = append(faltantes, "descripción")
		}
		if strings.TrimSpace(p.Contacto) == "" {
			faltantes = append(faltantes, "contacto")
		}
		if !p.TieneParticipantes {
			faltantes = append(faltantes, "participantes")
		}
		if len(faltantes) > 0 {
			incompletos = append(incompletos, ProductoIncompleto{ID: p.ID, Titulo: p.Titulo, Anio: p.Anio, Faltantes: faltantes})
		}
	}
	return incompletos
}
