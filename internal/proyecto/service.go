package proyecto

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
)

var (
	ErrTituloRequerido = errors.New("el título del proyecto es obligatorio")
	ErrAniosInvalidos  = errors.New("los años deben ser mayores a 0 y el año de inicio no puede ser mayor al de fin")
	ErrIDInvalido      = errors.New("el ID debe ser mayor a 0")
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

func (s *Service) Create(p *Proyecto) error {
	if p.Titulo == "" {
		return ErrTituloRequerido
	}
	if p.AnioInicio <= 0 || p.AnioFin <= 0 || p.AnioInicio > p.AnioFin {
		return ErrAniosInvalidos
	}

	err := s.storage.Create(p)
	if err != nil {
		s.logger.Error("Error al crear proyecto en el servicio", zap.Error(err))
		return fmt.Errorf("servicio create: %w", err)
	}
	s.logger.Info("Proyecto creado exitosamente", zap.Int("id", p.ID))
	return nil
}

func (s *Service) Read(id int) (*Proyecto, error) {
	if id <= 0 {
		return nil, ErrIDInvalido
	}
	p, err := s.storage.Read(id)
	if err != nil {
		s.logger.Error("Error al leer proyecto por ID en el servicio", zap.Int("id", id), zap.Error(err))
		return nil, err
	}
	return p, nil
}

func (s *Service) Update(id int, fields UpdateFields) error {
	if id <= 0 {
		return ErrIDInvalido
	}

	if fields.Titulo != nil && *fields.Titulo == "" {
		return ErrTituloRequerido
	}

	if fields.AnioInicio != nil && fields.AnioFin != nil {
		if *fields.AnioInicio <= 0 || *fields.AnioFin <= 0 || *fields.AnioInicio > *fields.AnioFin {
			return ErrAniosInvalidos
		}
	} else if fields.AnioInicio != nil && *fields.AnioInicio <= 0 {
		return ErrAniosInvalidos
	} else if fields.AnioFin != nil && *fields.AnioFin <= 0 {
		return ErrAniosInvalidos
	}

	err := s.storage.Update(id, fields)
	if err != nil {
		s.logger.Error("Error al actualizar proyecto en el servicio", zap.Int("id", id), zap.Error(err))
		return err
	}
	s.logger.Info("Proyecto actualizado exitosamente", zap.Int("id", id))
	return nil
}

func (s *Service) Delete(id int) error {
	if id <= 0 {
		return ErrIDInvalido
	}
	err := s.storage.Delete(id)
	if err != nil {
		s.logger.Error("Error al eliminar proyecto en el servicio", zap.Int("id", id), zap.Error(err))
		return err
	}
	s.logger.Info("Proyecto eliminado exitosamente", zap.Int("id", id))
	return nil
}

func (s *Service) GetAll() ([]Proyecto, error) {
	proyectos, err := s.storage.GetAll()
	if err != nil {
		s.logger.Error("Error al obtener lista de proyectos en el servicio", zap.Error(err))
		return nil, err
	}
	if proyectos == nil {
		proyectos = make([]Proyecto, 0)
	}
	return proyectos, nil
}

func (s *Service) GetAllAdmin() ([]Proyecto, error) {
	proyectos, err := s.storage.GetAllAdmin()
	if err != nil {
		s.logger.Error("Error al obtener proyectos para admin", zap.Error(err))
		return nil, err
	}
	if proyectos == nil {
		proyectos = make([]Proyecto, 0)
	}
	return proyectos, nil
}

func (s *Service) Restaurar(id int) error {
	if id <= 0 {
		return ErrIDInvalido
	}
	err := s.storage.Restaurar(id)
	if err != nil {
		s.logger.Error("Error al restaurar proyecto en el servicio", zap.Int("id", id), zap.Error(err))
		return err
	}
	s.logger.Info("Proyecto restaurado exitosamente", zap.Int("id", id))
	return nil
}
