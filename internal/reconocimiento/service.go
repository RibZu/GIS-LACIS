package reconocimiento

import (
	"errors"

	"go.uber.org/zap"
)

var ErrNotFound = errors.New("reconocimiento no encontrado")

type Service struct {
	storage Storage
	logger  *zap.Logger
}

func NewService(s Storage, l *zap.Logger) *Service {
	return &Service{
		storage: s,
		logger:  l,
	}
}

func (s *Service) GetAll() ([]Reconocimiento, error) {
	reconocimientos, err := s.storage.GetAll()
	if err != nil {
		s.logger.Error("Error al obtener reconocimientos en el servicio", zap.Error(err))
		return nil, err
	}
	if reconocimientos == nil {
		reconocimientos = make([]Reconocimiento, 0)
	}
	return reconocimientos, nil
}

// GetAllAdmin retorna todos los reconocimientos (incluyendo inactivos) para el panel admin
func (s *Service) GetAllAdmin() ([]Reconocimiento, error) {
	return s.GetAll()
}

func (s *Service) Create(r *Reconocimiento) error {
	if r.Titulo == "" {
		return errors.New("el título es requerido")
	}
	if err := s.storage.Create(r); err != nil {
		s.logger.Error("Error al crear reconocimiento", zap.Error(err))
		return err
	}
	return nil
}

func (s *Service) Read(id int) (*Reconocimiento, error) {
	r, err := s.storage.Read(id)
	if err != nil {
		if err.Error() == "reconocimiento no encontrado" {
			return nil, ErrNotFound
		}
		s.logger.Error("Error al leer reconocimiento", zap.Int("id", id), zap.Error(err))
		return nil, err
	}
	return r, nil
}

func (s *Service) Update(id int, fields UpdateFields) error {
	if err := s.storage.Update(id, fields); err != nil {
		s.logger.Error("Error al actualizar reconocimiento", zap.Int("id", id), zap.Error(err))
		return err
	}
	return nil
}

// Delete realiza baja lógica (activo = false)
func (s *Service) Delete(id int) error {
	if err := s.storage.Delete(id); err != nil {
		if err.Error() == "reconocimiento no encontrado" {
			return ErrNotFound
		}
		s.logger.Error("Error al dar de baja reconocimiento", zap.Int("id", id), zap.Error(err))
		return err
	}
	return nil
}

// Restaurar activa un reconocimiento dado de baja lógica
func (s *Service) Restaurar(id int) error {
	activo := true
	return s.Update(id, UpdateFields{Activo: &activo})
}
