package configuracion

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"go.uber.org/zap"
)

var (
	ErrEmailInvalido = errors.New("formato de correo electrónico inválido")
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

func (s *Service) ObtenerConfiguracion() (*ConfiguracionSitio, error) {
	cfg, err := s.storage.Get()
	if err != nil {
		s.logger.Error("Error al obtener configuración en service", zap.Error(err))
		return nil, err
	}
	// Defaults si los campos están vacíos
	if cfg.EmailDoctorado == "" {
		cfg.EmailDoctorado = "agaris@gmail.com"
	}
	if cfg.EmailMaestriaCalidad == "" {
		cfg.EmailMaestriaCalidad = "magister@unsl.edu.ar"
	}
	if cfg.EmailMaestriaSoftware == "" {
		cfg.EmailMaestriaSoftware = "magister@unsl.edu.ar"
	}
	if cfg.EmailEspecializacion == "" {
		cfg.EmailEspecializacion = "mgperalta03@gmail.com"
	}
	if cfg.EmailIngInformatica == "" {
		cfg.EmailIngInformatica = "ingenieriainformaticaunsl@gmail.com"
	}
	// TecnicaturaWeb no tiene correo por defecto (permanece vacío a menos que se configure)
	return cfg, nil
}

func (s *Service) GuardarConfiguracion(dto UpdateConfiguracionDTO) (*ConfiguracionSitio, error) {
	current, err := s.storage.Get()
	if err != nil {
		current = &ConfiguracionSitio{}
	}

	validarMail := func(email, nombreCampo string) error {
		trimmed := strings.TrimSpace(email)
		if trimmed == "" {
			return nil
		}
		if _, err := mail.ParseAddress(trimmed); err != nil {
			return fmt.Errorf("%s: %s", ErrEmailInvalido.Error(), nombreCampo)
		}
		return nil
	}

	if dto.EmailDoctorado != nil {
		if err := validarMail(*dto.EmailDoctorado, "Email Doctorado"); err != nil {
			return nil, err
		}
		current.EmailDoctorado = strings.TrimSpace(*dto.EmailDoctorado)
	}

	if dto.EmailMaestriaCalidad != nil {
		if err := validarMail(*dto.EmailMaestriaCalidad, "Email Maestría en Calidad"); err != nil {
			return nil, err
		}
		current.EmailMaestriaCalidad = strings.TrimSpace(*dto.EmailMaestriaCalidad)
	}

	if dto.EmailMaestriaSoftware != nil {
		if err := validarMail(*dto.EmailMaestriaSoftware, "Email Maestría en Software"); err != nil {
			return nil, err
		}
		current.EmailMaestriaSoftware = strings.TrimSpace(*dto.EmailMaestriaSoftware)
	}

	if dto.EmailEspecializacion != nil {
		if err := validarMail(*dto.EmailEspecializacion, "Email Especialización"); err != nil {
			return nil, err
		}
		current.EmailEspecializacion = strings.TrimSpace(*dto.EmailEspecializacion)
	}

	if dto.EmailIngInformatica != nil {
		if err := validarMail(*dto.EmailIngInformatica, "Email Ingeniería en Informática"); err != nil {
			return nil, err
		}
		current.EmailIngInformatica = strings.TrimSpace(*dto.EmailIngInformatica)
	}

	if dto.EmailTecnicaturaWeb != nil {
		if err := validarMail(*dto.EmailTecnicaturaWeb, "Email Tecnicatura Universitaria en Web"); err != nil {
			return nil, err
		}
		current.EmailTecnicaturaWeb = strings.TrimSpace(*dto.EmailTecnicaturaWeb)
	}

	if dto.EmailFooter != nil {
		if err := validarMail(*dto.EmailFooter, "Email de Contacto General"); err != nil {
			return nil, err
		}
		current.EmailFooter = strings.TrimSpace(*dto.EmailFooter)
	}

	if dto.TelefonoFooter != nil {
		current.TelefonoFooter = strings.TrimSpace(*dto.TelefonoFooter)
	}
	if dto.Direccion != nil {
		current.Direccion = strings.TrimSpace(*dto.Direccion)
	}
	if dto.UbicacionFooter != nil {
		current.UbicacionFooter = strings.TrimSpace(*dto.UbicacionFooter)
	}
	if dto.OficinasFooter != nil {
		current.OficinasFooter = strings.TrimSpace(*dto.OficinasFooter)
	}
	if dto.MapaURL != nil {
		current.MapaURL = strings.TrimSpace(*dto.MapaURL)
	}

	if err := s.storage.Update(current); err != nil {
		s.logger.Error("Error al actualizar configuración en storage", zap.Error(err))
		return nil, fmt.Errorf("error al guardar configuración: %w", err)
	}

	s.logger.Info("Configuración del sitio actualizada exitosamente")
	return current, nil
}
