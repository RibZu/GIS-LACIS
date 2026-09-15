package desarrollo

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"PaginaSEG/internal/integrante"

	"go.uber.org/zap"
)

var (
	ErrTituloRequerido = errors.New("el título es obligatorio")
	ErrAnioInvalido    = errors.New("el año debe ser válido (entre 1990 y el año actual + 1)")
	ErrIDInvalido      = errors.New("el ID debe ser mayor a 0")
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

// normalizarURL le agrega "https://" a los enlaces que el usuario cargó sin
// protocolo (ej: "hola.com.ar"), para que el <a href="..."> de la lista
// funcione como link externo real y no como ruta relativa del propio sitio.
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

func (s *Service) Create(d *Desarrollo) error {
	if d.Titulo == "" {
		return ErrTituloRequerido
	}
	if !anioValido(d.Anio) {
		return ErrAnioInvalido
	}
	d.URL = normalizarURL(d.URL)

	err := s.storage.Create(d)
	if err != nil {
		s.logger.Error("Error al crear desarrollo en el servicio", zap.Error(err))
		return fmt.Errorf("servicio create: %w", err)
	}
	s.logger.Info("Desarrollo creado exitosamente", zap.Int("id", d.ID))
	return nil
}

func (s *Service) Update(id int, fields UpdateFields) error {
	if id <= 0 {
		return ErrIDInvalido
	}
	if fields.Titulo != nil && *fields.Titulo == "" {
		return ErrTituloRequerido
	}
	if fields.Anio != nil && !anioValido(*fields.Anio) {
		return ErrAnioInvalido
	}
	if fields.URL != nil {
		normalizada := normalizarURL(*fields.URL)
		fields.URL = &normalizada
	}

	err := s.storage.Update(id, fields)
	if err != nil {
		s.logger.Error("Error al actualizar desarrollo", zap.Int("id", id), zap.Error(err))
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
		s.logger.Error("Error al leer desarrollo por ID", zap.Int("id", id), zap.Error(err))
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
		s.logger.Error("Error al eliminar desarrollo", zap.Int("id", id), zap.Error(err))
		return err
	}
	s.logger.Info("Desarrollo eliminado exitosamente", zap.Int("id", id))
	return nil
}

// -------- Vínculo con integrantes registrados (desarrollo_integrantes) --------

func (s *Service) VincularIntegrantes(desarrolloID int, integranteIDs []int) error {
	if desarrolloID <= 0 {
		return ErrIDInvalido
	}
	return s.storage.SetIntegrantes(desarrolloID, integranteIDs)
}

func (s *Service) ObtenerIntegrantes(desarrolloID int) ([]integrante.Integrante, error) {
	if desarrolloID <= 0 {
		return nil, ErrIDInvalido
	}
	return s.storage.GetIntegrantes(desarrolloID)
}

// -------- Participantes externos (participante_externo) --------

func (s *Service) VincularParticipantesExternos(desarrolloID int, nombres []string) error {
	if desarrolloID <= 0 {
		return ErrIDInvalido
	}
	return s.storage.SetParticipantesExternos(desarrolloID, nombres)
}

func (s *Service) ObtenerParticipantesExternos(desarrolloID int) ([]string, error) {
	if desarrolloID <= 0 {
		return nil, ErrIDInvalido
	}
	return s.storage.GetParticipantesExternos(desarrolloID)
}
