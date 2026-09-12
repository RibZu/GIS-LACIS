package integrante

import (
	"errors"
	"fmt"
	"regexp"

	"go.uber.org/zap"
)

var regexSoloLetras = regexp.MustCompile(`^[a-zA-ZáéíóúÁÉÍÓÚñÑ\s\.\-']+$`)

// Errores de validación de negocio
var (
	ErrNombreRequerido          = errors.New("el nombre es obligatorio")
	ErrApellidoRequerido        = errors.New("el apellido es obligatorio")
	ErrApellidoInvalido         = errors.New("el apellido solo debe contener letras y caracteres válidos")
	ErrNombreInvalido           = errors.New("el nombre solo debe contener letras y caracteres válidos")
	ErrRolRequerido             = errors.New("el rol es obligatorio")
	ErrRolLacisRequerido        = errors.New("debes seleccionar un rol para LaCIS")
	ErrRolSoftwareRequerido     = errors.New("debes seleccionar un rol para Grupo Software")
	ErrPertenenciaRequerida     = errors.New("el integrante debe pertenecer al menos a un grupo (LaCIS o Grupo Software)")
	ErrEspecializacionRequerida = errors.New("la especialización o título es obligatorio")
	ErrContactoRequerido        = errors.New("el contacto o email es obligatorio")
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
	if !esSoloLetras(integrante.Nombre) {
		return ErrNombreInvalido
	}
	if !esSoloLetras(integrante.Apellido) {
		return ErrApellidoInvalido
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

	// Si no tiene ningún grupo seleccionado
	if !integrante.PerteneceLacis && !integrante.PerteneceGrupoSoftware {
		return ErrPertenenciaRequerida
	}

	// Asignar rol general si viene RolID legacy
	if integrante.RolID > 0 {
		if integrante.PerteneceLacis && integrante.RolLacisID == nil {
			val := integrante.RolID
			integrante.RolLacisID = &val
		}
		if integrante.PerteneceGrupoSoftware && integrante.RolSoftwareID == nil {
			val := integrante.RolID
			integrante.RolSoftwareID = &val
		}
	}

	// Validar rol de LaCIS si pertenece
	if integrante.PerteneceLacis {
		if integrante.RolLacisID == nil || *integrante.RolLacisID <= 0 {
			return ErrRolLacisRequerido
		}
	} else {
		integrante.RolLacisID = nil
	}

	// Validar rol de Grupo Software si pertenece
	if integrante.PerteneceGrupoSoftware {
		if integrante.RolSoftwareID == nil || *integrante.RolSoftwareID <= 0 {
			return ErrRolSoftwareRequerido
		}
	} else {
		integrante.RolSoftwareID = nil
	}

	// Asegurar RolID para retrocompatibilidad
	if integrante.RolID <= 0 {
		if integrante.RolLacisID != nil {
			integrante.RolID = *integrante.RolLacisID
		} else if integrante.RolSoftwareID != nil {
			integrante.RolID = *integrante.RolSoftwareID
		}
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
	if fields.Especializacion != nil && *fields.Especializacion == "" {
		return ErrEspecializacionRequerida
	}
	if fields.Contacto != nil && *fields.Contacto == "" {
		return ErrContactoRequerido
	}
	if fields.Descripcion != nil && *fields.Descripcion == "" {
		return ErrDescripcionRequerida
	}

	// Validaciones de roles por grupo si se envían
	if fields.PerteneceLacis != nil && *fields.PerteneceLacis {
		if fields.RolLacisID != nil && *fields.RolLacisID <= 0 {
			return ErrRolLacisRequerido
		}
	}
	if fields.PerteneceGrupoSoftware != nil && *fields.PerteneceGrupoSoftware {
		if fields.RolSoftwareID != nil && *fields.RolSoftwareID <= 0 {
			return ErrRolSoftwareRequerido
		}
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
