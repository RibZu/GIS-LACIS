package integrante

import (
	"errors"
	"fmt"
	"regexp"

	"go.uber.org/zap"
)

var regexSoloLetras = regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑ\s]+$`)

// Errores de validación de negocio
var (
	ErrNombreRequerido          = errors.New("el nombre es obligatorio")
	ErrApellidoRequerido        = errors.New("el apellido es obligatorio")
	ErrApellidoInvalido         = errors.New("el apellido solo debe contener letras")
	ErrNombreInvalido           = errors.New("el nombre solo debe contener letras")
	ErrRolRequerido             = errors.New("el rol es obligatorio")
	ErrEspecializacionRequerida = errors.New("la especialización es obligatoria")
	ErrContactoRequerido        = errors.New("el contacto es obligatorio")
	ErrDescripcionRequerida     = errors.New("la descripción es obligatoria")
	ErrIDInvalido               = errors.New("el ID debe ser mayor a 0")
)

// Service contiene la lógica de negocio e interactúa con Storage y Logger
type Service struct {
	storage Storage
	logger  *zap.Logger
}

// NewService es el constructor del servicio
func NewService(s Storage, l *zap.Logger) *Service {
	return &Service{
		storage: s,
		logger:  l,
	}
}

func esSoloLetras(texto string) bool {
	return regexSoloLetras.MatchString(texto)
}

func (s *Service) Create(integrante *Integrante) error {
	if integrante.Nombre == "" {
		return ErrNombreRequerido
	}
	if integrante.Apellido == "" {
		return ErrApellidoRequerido
	}
	if integrante.RolID == 0 {
		return ErrRolRequerido
	}
	if integrante.Especializacion == "" {
		return ErrEspecializacionRequerida
	}
	if integrante.Contacto == "" {
		return ErrContactoRequerido
	}
	if integrante.Descripcion == "" {
		return ErrDescripcionRequerida
	}
	if !esSoloLetras(integrante.Nombre) {
		return ErrNombreInvalido
	}

	if !esSoloLetras(integrante.Apellido) {
		return ErrApellidoInvalido
	}

	err := s.storage.Create(integrante)
	if err != nil {
		s.logger.Error("Error al crear integrante en el servicio", zap.Error(err))
		return fmt.Errorf("servicio create: %w", err)
	}
	s.logger.Info("Integrante creado exitosamente", zap.Int("id", integrante.ID))
	return nil
}

func (s *Service) Update(id int, fields UpdateFields) error {
	if id <= 0 {
		return ErrIDInvalido
	}
	if fields.Nombre != nil {
		if *fields.Nombre == "" {
			return ErrNombreRequerido
		}
		if !esSoloLetras(*fields.Nombre) {
			return ErrNombreInvalido
		}
	}
	if fields.Apellido != nil {
		if *fields.Apellido == "" {
			return ErrApellidoRequerido
		}
		if !esSoloLetras(*fields.Apellido) {
			return ErrApellidoInvalido
		}
	}
	if fields.RolID != nil && *fields.RolID <= 0 {
		return ErrRolRequerido
	}
	if fields.Especializacion != nil && *fields.Especializacion == "" {
		return ErrEspecializacionRequerida
	}
	if fields.Contacto != nil && *fields.Contacto == "" {
		return ErrContactoRequerido
	}
	if fields.Descripcion != nil && *fields.Descripcion == "" {
		return ErrDescripcionRequerida
	}
	err := s.storage.Update(id, fields)
	if err != nil {
		s.logger.Error("Error al actualizar integrante", zap.Int("id", id), zap.Error(err))
		return err
	}
	s.logger.Info("Integrante actualizado exitosamente", zap.Int("id", id))
	return nil
}

func (s *Service) Read(id int) (*Integrante, error) {
	if id <= 0 {
		return nil, ErrIDInvalido
	}
	integrante, err := s.storage.Read(id)

	if err != nil {
		s.logger.Error("Error al leer integrante por ID en el servicio", zap.Int("id", id), zap.Error(err))
		return nil, err
	}

	return integrante, nil

}

func (s *Service) GetAll() ([]Integrante, error) {

	integrantes, err := s.storage.GetAll()
	if err != nil {
		s.logger.Error("Error al obtener la lista de integrantes en el servicio", zap.Error(err))
		return nil, err
	}

	return integrantes, nil

}

func (s *Service) Delete(id int) error {
	if id <= 0 {
		return ErrIDInvalido
	}
	err := s.storage.Delete(id)
	if err != nil {
		s.logger.Error("Error al eliminar integrante en el servicio", zap.Int("id", id), zap.Error(err))
		return err
	}
	s.logger.Info("Integrante eliminado exitosamente", zap.Int("id", id))
	return nil
}
