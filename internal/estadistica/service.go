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
	integrantesLacis, err := s.storage.ContarIntegrantesActivosLacis()
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
	integrantes, err := s.storage.ResumenIntegrantes()
	if err != nil {
		return nil, s.fallo("integrantes", err)
	}
	tesis, err := s.storage.ResumenTesis()
	if err != nil {
		return nil, s.fallo("tesis", err)
	}
	proyectos, err := s.storage.ResumenProyectos()
	if err != nil {
		return nil, s.fallo("proyectos", err)
	}

	est := &Estadisticas{
		IntegrantesActivosLacis: integrantesLacis,
		ParticipantesExternos:   contarPersonasDistintas(nombres),
		Productos:               len(productos),
	}
	anios := make([]int, 0, len(productos))
	for _, p := range productos {
		if desarrollo.TieneRepositorio(p.URL) {
			est.ProductosConRepositorio++
		}
		anios = append(anios, p.Anio)
	}
	est.ProductosConLicencia = est.Productos - est.ProductosConRepositorio
	est.HayProductos = est.Productos > 0
	est.PorcentajeConRepositorio, est.PorcentajeConLicencia = porcentajes(est.ProductosConRepositorio, est.Productos)
	est.PorAnio = seriePorAnio(anios)
	est.Integrantes = estadisticasIntegrantes(integrantes)
	est.Tesis = estadisticasTesis(tesis)
	est.Proyectos = estadisticasProyectos(proyectos)
	return est, nil
}

func (s *Service) fallo(que string, err error) error {
	s.logger.Error("No se pudieron obtener las estadísticas", zap.String("lectura", que), zap.Error(err))
	return fmt.Errorf("estadisticas: %w", err)
}

// porcentajes reparte 100 entre parte y el resto de total. Redondea parte al entero más cercano y
// el segundo porcentaje es lo que falta, así suman 100. Sin total devuelve 0 y 0.
func porcentajes(parte, total int) (int, int) {
	if total == 0 {
		return 0, 0
	}
	p := (200*parte + total) / (2 * total)
	return p, 100 - p
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

// seriePorAnio arma una fila por cada año entre el primero y el último de la lista, incluso los
// años sin ninguno, para que la serie no salte años. La usan productos y tesis.
func seriePorAnio(anios []int) []CantidadPorAnio {
	if len(anios) == 0 {
		return nil
	}
	cantidades := make(map[int]int)
	minAnio, maxAnio := anios[0], anios[0]
	for _, a := range anios {
		cantidades[a]++
		minAnio = min(minAnio, a)
		maxAnio = max(maxAnio, a)
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

// estadisticasIntegrantes cuenta igual que lo hacía antes lista.js en la lista de integrantes: todo
// el que no está activo es inactivo, y LaCIS y Grupo Software se cuentan por separado.
func estadisticasIntegrantes(filas []IntegranteResumen) EstadisticasIntegrantes {
	e := EstadisticasIntegrantes{Registrados: len(filas)}
	for _, f := range filas {
		if f.Activo {
			e.Activos++
		} else {
			e.Inactivos++
		}
		if f.PerteneceLacis {
			e.Lacis++
		}
		if f.PerteneceGrupoSoftware {
			e.Software++
		}
	}
	return e
}

// esPosgrado es la misma regla que usaba antes listaTesis.js en la lista de tesis. Se dejó tal cual
// a pedido del usuario: no cambiarla aunque dependa del texto del nivel.
func esPosgrado(nivel string) bool {
	n := strings.ToLower(nivel)
	return strings.Contains(n, "doctor") || strings.Contains(n, "maestr") ||
		strings.Contains(n, "especializ") || strings.Contains(n, "posgrado")
}

// estadisticasTesis arma las tarjetas de tesis y los gráficos por año y de PDF. Las tesis sin año
// cuentan en todo menos en la serie por año.
func estadisticasTesis(filas []TesisResumen) EstadisticasTesis {
	e := EstadisticasTesis{Registradas: len(filas), HayTesis: len(filas) > 0}
	var anios []int
	for _, f := range filas {
		if esPosgrado(f.Nivel) {
			e.Posgrado++
		}
		if f.TienePDF {
			e.ConPDF++
		}
		if f.Anio != nil {
			anios = append(anios, *f.Anio)
		}
	}
	e.GradoOtros = e.Registradas - e.Posgrado
	e.SinPDF = e.Registradas - e.ConPDF
	e.PorcentajeConPDF, e.PorcentajeSinPDF = porcentajes(e.ConPDF, e.Registradas)
	e.PorAnio = seriePorAnio(anios)
	return e
}

// estadisticasProyectos agrupa las filas por proyecto y calcula la cantidad de integrantes por rol.
// Los proyectos antiguos que no tienen integrantes cargados (Cantidad = 0) se conservan con TieneIntegrantes = false.
func estadisticasProyectos(filas []ProyectoFilaEstadistica) EstadisticasProyectos {
	res := EstadisticasProyectos{
		HayProyectos: len(filas) > 0,
	}
	if len(filas) == 0 {
		return res
	}

	proyectosMap := make(map[int]*ProyectoEstadistica)
	var proyectosOrdenados []*ProyectoEstadistica
	globalRolesMap := make(map[string]int)

	for _, f := range filas {
		p, existe := proyectosMap[f.ID]
		if !existe {
			p = &ProyectoEstadistica{
				ID:              f.ID,
				Titulo:          f.Titulo,
				AnioInicio:      f.AnioInicio,
				AnioFin:         f.AnioFin,
				Activo:          f.Activo,
				EquipoHistorico: f.EquipoHistorico,
			}
			proyectosMap[f.ID] = p
			proyectosOrdenados = append(proyectosOrdenados, p)
		}

		if f.Cantidad > 0 {
			rolNombre := strings.TrimSpace(f.Rol)
			if rolNombre == "" {
				rolNombre = "Sin rol asignado"
			}
			p.Roles = append(p.Roles, RolCantidad{
				Rol:      rolNombre,
				Cantidad: f.Cantidad,
			})
			p.TotalIntegrantes += f.Cantidad
			p.TieneIntegrantes = true

			globalRolesMap[rolNombre] += f.Cantidad
			res.TotalParticipaciones += f.Cantidad
		}
	}

	res.TotalProyectos = len(proyectosOrdenados)
	res.Lista = make([]ProyectoEstadistica, 0, len(proyectosOrdenados))

	for _, p := range proyectosOrdenados {
		if p.TieneIntegrantes {
			res.ProyectosConIntegrantes++
		} else {
			res.ProyectosSinDatos++
		}
		res.Lista = append(res.Lista, *p)
	}

	// Orden canónico para los roles globales más comunes del sistema
	rolesOrdenCanonico := []string{
		"Director / Co-Director",
		"Investigador",
		"Asesor Externo",
		"Estudiante / Becario",
		"Sin rol asignado",
	}

	rolesVistos := make(map[string]bool)
	for _, rNombre := range rolesOrdenCanonico {
		if cant, ok := globalRolesMap[rNombre]; ok && cant > 0 {
			res.RolesGlobal = append(res.RolesGlobal, RolCantidad{Rol: rNombre, Cantidad: cant})
			rolesVistos[rNombre] = true
		}
	}
	for rNombre, cant := range globalRolesMap {
		if !rolesVistos[rNombre] && cant > 0 {
			res.RolesGlobal = append(res.RolesGlobal, RolCantidad{Rol: rNombre, Cantidad: cant})
		}
	}

	return res
}

