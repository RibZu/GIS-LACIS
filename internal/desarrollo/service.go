package desarrollo

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"PaginaSEG/internal/integrante"

	"go.uber.org/zap"
)

var (
	ErrTituloRequerido   = errors.New("el título es obligatorio")
	ErrAnioInvalido      = errors.New("el año debe ser válido (entre 1990 y el año actual + 1)")
	ErrIDInvalido        = errors.New("el ID debe ser mayor a 0")
	ErrTituloLargo       = errors.New("el título no puede superar los 255 caracteres")
	ErrURLLarga          = errors.New("el enlace no puede superar los 500 caracteres")
	ErrContactoLargo     = errors.New("el contacto no puede superar los 255 caracteres")
	ErrParticipanteLargo = errors.New("cada participante externo puede tener hasta 150 caracteres")
)

const (
	maxTitulo       = 255
	maxURL          = 500
	maxContacto     = 255
	maxParticipante = 150
)

type Service struct {
	storage Storage
	logger  *zap.Logger
}

func NewService(s Storage, l *zap.Logger) *Service {
	return &Service{storage: s, logger: l}
}

func anioValido(anio int) bool {
	limite := time.Now().Year() + 1
	return anio >= 1990 && anio <= limite
}

func normalizarURL(url string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return ""
	}
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	return "https://" + url
}

func excede(texto string, maximo int) bool {
	return utf8.RuneCountInString(texto) > maximo
}

func validarTitulo(titulo string) error {
	if strings.TrimSpace(titulo) == "" {
		return ErrTituloRequerido
	}
	if excede(titulo, maxTitulo) {
		return ErrTituloLargo
	}
	return nil
}

func normalizarExternos(nombres []string) ([]string, error) {
	vistos := make(map[string]struct{}, len(nombres))
	limpios := make([]string, 0, len(nombres))
	for _, nombre := range nombres {
		nombre = strings.Join(strings.Fields(nombre), " ")
		if nombre == "" {
			continue
		}
		if excede(nombre, maxParticipante) {
			return nil, ErrParticipanteLargo
		}
		clave := strings.ToLower(nombre)
		if _, repetido := vistos[clave]; repetido {
			continue
		}
		vistos[clave] = struct{}{}
		limpios = append(limpios, nombre)
	}
	return limpios, nil
}

func normalizarIDs(ids []int) []int {
	vistos := make(map[int]struct{}, len(ids))
	limpios := make([]int, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, repetido := vistos[id]; repetido {
			continue
		}
		vistos[id] = struct{}{}
		limpios = append(limpios, id)
	}
	return limpios
}

func (s *Service) CrearCompleto(d *Desarrollo, integranteIDs []int, externos []string) error {
	if err := validarTitulo(d.Titulo); err != nil {
		return err
	}
	if !anioValido(d.Anio) {
		return ErrAnioInvalido
	}
	d.URL = normalizarURL(d.URL)
	if excede(d.URL, maxURL) {
		return ErrURLLarga
	}
	if excede(d.Contacto, maxContacto) {
		return ErrContactoLargo
	}
	limpios, err := normalizarExternos(externos)
	if err != nil {
		return err
	}

	if err := s.storage.CrearCompleto(d, normalizarIDs(integranteIDs), limpios); err != nil {
		if !errors.Is(err, ErrIntegranteInexistente) {
			s.logger.Error("Error al crear desarrollo", zap.Error(err))
		}
		return fmt.Errorf("crear desarrollo: %w", err)
	}
	s.logger.Info("Desarrollo creado exitosamente", zap.Int("id", d.ID))
	return nil
}

func (s *Service) ActualizarCompleto(id int, fields UpdateFields, integranteIDs []int, externos []string) error {
	if id <= 0 {
		return ErrIDInvalido
	}
	if fields.Titulo != nil {
		if err := validarTitulo(*fields.Titulo); err != nil {
			return err
		}
	}
	if fields.Anio != nil && !anioValido(*fields.Anio) {
		return ErrAnioInvalido
	}
	if fields.URL != nil {
		normalizada := normalizarURL(*fields.URL)
		if excede(normalizada, maxURL) {
			return ErrURLLarga
		}
		fields.URL = &normalizada
	}
	if fields.Contacto != nil && excede(*fields.Contacto, maxContacto) {
		return ErrContactoLargo
	}
	limpios, err := normalizarExternos(externos)
	if err != nil {
		return err
	}

	if err := s.storage.ActualizarCompleto(id, fields, normalizarIDs(integranteIDs), limpios); err != nil {
		if !errors.Is(err, ErrNotFound) && !errors.Is(err, ErrIntegranteInexistente) {
			s.logger.Error("Error al actualizar desarrollo", zap.Int("id", id), zap.Error(err))
		}
		return err
	}
	s.logger.Info("Desarrollo actualizado exitosamente", zap.Int("id", id))
	return nil
}

func (s *Service) Read(id int) (*Desarrollo, error) {
	if id <= 0 {
		return nil, ErrIDInvalido
	}
	d, err := s.storage.Read(id)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			s.logger.Error("Error al leer desarrollo por ID", zap.Int("id", id), zap.Error(err))
		}
		return nil, err
	}
	return d, nil
}

func (s *Service) GetAll() ([]Desarrollo, error) {
	lista, err := s.storage.GetAll()
	if err != nil {
		s.logger.Error("Error al obtener la lista de desarrollos", zap.Error(err))
		return nil, err
	}
	return lista, nil
}

func (s *Service) Delete(id int) error {
	if id <= 0 {
		return ErrIDInvalido
	}
	err := s.storage.Delete(id)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			s.logger.Error("Error al eliminar desarrollo", zap.Int("id", id), zap.Error(err))
		}
		return err
	}
	s.logger.Info("Desarrollo eliminado exitosamente", zap.Int("id", id))
	return nil
}

func (s *Service) ObtenerIntegrantes(desarrolloID int) ([]integrante.Integrante, error) {
	if desarrolloID <= 0 {
		return nil, ErrIDInvalido
	}
	return s.storage.GetIntegrantes(desarrolloID)
}

func (s *Service) ObtenerParticipantesExternos(desarrolloID int) ([]string, error) {
	if desarrolloID <= 0 {
		return nil, ErrIDInvalido
	}
	return s.storage.GetParticipantesExternos(desarrolloID)
}

func (s *Service) ParticipantesPorDesarrollo() (map[int][]integrante.Integrante, map[int][]string, error) {
	integrantes, err := s.storage.IntegrantesPorDesarrollo()
	if err != nil {
		s.logger.Error("Error al obtener los integrantes de los desarrollos", zap.Error(err))
		return nil, nil, err
	}
	externos, err := s.storage.ParticipantesExternosPorDesarrollo()
	if err != nil {
		s.logger.Error("Error al obtener los participantes externos de los desarrollos", zap.Error(err))
		return nil, nil, err
	}
	return integrantes, externos, nil
}
