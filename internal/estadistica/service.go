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

func (s *Service) Obtener() (*Estadisticas, error) {
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
		IntegrantesActivosLacis: contarActivosLacis(integrantes),
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

func contarActivosLacis(filas []IntegranteResumen) int {
	total := 0
	for _, f := range filas {
		if f.Activo && f.PerteneceLacis {
			total++
		}
	}
	return total
}

func (s *Service) fallo(que string, err error) error {
	s.logger.Error("No se pudieron obtener las estadísticas", zap.String("lectura", que), zap.Error(err))
	return fmt.Errorf("estadisticas: %w", err)
}

func porcentajes(parte, total int) (int, int) {
	if total == 0 {
		return 0, 0
	}
	p := (200*parte + total) / (2 * total)
	return p, 100 - p
}

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

func esPosgrado(nivel string) bool {
	n := strings.ToLower(nivel)
	return strings.Contains(n, "doctor") || strings.Contains(n, "maestr") ||
		strings.Contains(n, "especializ") || strings.Contains(n, "posgrado")
}

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

