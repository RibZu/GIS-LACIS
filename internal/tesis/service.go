package tesis

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

var (
	ErrTituloRequerido  = errors.New("el título de la tesis es obligatorio")
	ErrNivelRequerido   = errors.New("debes seleccionar el nivel académico de la tesis")
	ErrCarreraRequerida = errors.New("debes indicar la carrera o posgrado de origen")
	ErrAnioInvalido     = errors.New("el año debe ser un número válido de 4 dígitos")
	ErrAutorRequerido   = errors.New("debes indicar o seleccionar un autor para la tesis")
	ErrIDInvalido       = errors.New("el ID debe ser mayor a 0")
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

func (s *Service) Create(t *Tesis) error {
	t.Titulo = strings.TrimSpace(t.Titulo)
	if t.Titulo == "" {
		return ErrTituloRequerido
	}

	t.Nivel = strings.TrimSpace(t.Nivel)
	if t.Nivel == "" {
		return ErrNivelRequerido
	}

	t.CarreraOrigen = strings.TrimSpace(t.CarreraOrigen)
	if t.CarreraOrigen == "" {
		return ErrCarreraRequerida
	}

	// Validar año si viene provisto
	if t.Anio != nil {
		actual := time.Now().Year() + 5
		if *t.Anio < 1970 || *t.Anio > actual {
			return ErrAnioInvalido
		}
	}

	// Validar que exista al menos un autor (por ID o texto histórico)
	tieneAutorID := t.AutorID != nil && *t.AutorID > 0
	tieneAutorHist := strings.TrimSpace(t.AutorHistorico) != ""
	if !tieneAutorID && !tieneAutorHist {
		return ErrAutorRequerido
	}

	err := s.storage.Create(t)
	if err != nil {
		s.logger.Error("Error al crear tesis en el servicio", zap.Error(err))
		return fmt.Errorf("servicio create tesis: %w", err)
	}

	s.logger.Info("Tesis creada exitosamente", zap.Int("id", t.ID), zap.String("titulo", t.Titulo))
	return nil
}

func (s *Service) Read(id int) (*Tesis, error) {
	if id <= 0 {
		return nil, ErrIDInvalido
	}
	t, err := s.storage.Read(id)
	if err != nil {
		s.logger.Error("Error al leer tesis por ID en el servicio", zap.Int("id", id), zap.Error(err))
		return nil, err
	}
	return t, nil
}

func (s *Service) GetAll() ([]Tesis, error) {
	lista, err := s.storage.GetAll()
	if err != nil {
		s.logger.Error("Error al listar tesis en el servicio", zap.Error(err))
		return nil, err
	}
	return lista, nil
}

func (s *Service) GetByCarrera(carrera string) ([]Tesis, error) {
	if carrera == "" {
		return s.GetAll()
	}
	return s.storage.GetByCarrera(carrera)
}

func (s *Service) GetByNivel(nivel string) ([]Tesis, error) {
	if nivel == "" {
		return s.GetAll()
	}
	return s.storage.GetByNivel(nivel)
}

func (s *Service) Update(id int, fields UpdateFields) error {
	if id <= 0 {
		return ErrIDInvalido
	}
	if fields.Titulo != nil {
		*fields.Titulo = strings.TrimSpace(*fields.Titulo)
		if *fields.Titulo == "" {
			return ErrTituloRequerido
		}
	}
	if fields.Nivel != nil {
		*fields.Nivel = strings.TrimSpace(*fields.Nivel)
		if *fields.Nivel == "" {
			return ErrNivelRequerido
		}
	}
	if fields.CarreraOrigen != nil {
		*fields.CarreraOrigen = strings.TrimSpace(*fields.CarreraOrigen)
		if *fields.CarreraOrigen == "" {
			return ErrCarreraRequerida
		}
	}
	if fields.Anio != nil && *fields.Anio > 0 {
		actual := time.Now().Year() + 5
		if *fields.Anio < 1970 || *fields.Anio > actual {
			return ErrAnioInvalido
		}
	}

	err := s.storage.Update(id, fields)
	if err != nil {
		s.logger.Error("Error al actualizar tesis", zap.Int("id", id), zap.Error(err))
		return err
	}
	s.logger.Info("Tesis actualizada exitosamente", zap.Int("id", id))
	return nil
}

func (s *Service) Delete(id int) error {
	if id <= 0 {
		return ErrIDInvalido
	}
	err := s.storage.Delete(id)
	if err != nil {
		s.logger.Error("Error al eliminar tesis en el servicio", zap.Int("id", id), zap.Error(err))
		return err
	}
	s.logger.Info("Tesis eliminada exitosamente", zap.Int("id", id))
	return nil
}
