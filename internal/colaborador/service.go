package colaborador

import (
	"errors"

	"go.uber.org/zap"
)

var (
	ErrNotFound = errors.New("colaborador no encontrado")
)

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

// GetAllAdmin devuelve todos los colaboradores sin filtrar.
func (s *Service) GetAllAdmin() ([]Colaborador, error) {
	return s.storage.GetAll()
}

// GetAll devuelve solo los colaboradores activos.
func (s *Service) GetAll() ([]Colaborador, error) {
	all, err := s.storage.GetAll()
	if err != nil {
		return nil, err
	}

	activos := make([]Colaborador, 0)
	for _, c := range all {
		if c.Activo {
			activos = append(activos, c)
		}
	}
	return activos, nil
}

func (s *Service) Create(c *Colaborador) error {
	if c.LogoURL == "" {
		return errors.New("debe ingresar el logo de la empresa")
	}
	if err := s.storage.Create(c); err != nil {
		s.logger.Error("Error al crear colaborador", zap.Error(err))
		return err
	}
	return nil
}

func (s *Service) Read(id int) (*Colaborador, error) {
	colab, err := s.storage.Read(id)
	if err != nil {
		if err.Error() == "colaborador no encontrado" {
			return nil, ErrNotFound
		}
		s.logger.Error("Error al leer colaborador", zap.Int("id", id), zap.Error(err))
		return nil, err
	}
	return colab, nil
}

func (s *Service) Update(id int, fields UpdateFields) error {
	// Verificar existencia
	_, err := s.Read(id)
	if err != nil {
		return err
	}

	if err := s.storage.Update(id, fields); err != nil {
		s.logger.Error("Error al actualizar colaborador", zap.Int("id", id), zap.Error(err))
		return err
	}
	return nil
}

func (s *Service) Delete(id int) error {
	_, err := s.Read(id)
	if err != nil {
		return err
	}

	if err := s.storage.Delete(id); err != nil {
		s.logger.Error("Error al dar de baja colaborador", zap.Int("id", id), zap.Error(err))
		return err
	}
	return nil
}

func (s *Service) Restaurar(id int) error {
	// Reutiliza Update
	activo := true
	return s.Update(id, UpdateFields{Activo: &activo})
}
