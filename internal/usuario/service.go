package usuario

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
)

var (
	ErrUsernameRequerido = errors.New("Username requerido")
	ErrPasswordRequerido = errors.New("Password requerido")
	ErrEmailRequerido    = errors.New("Email requerido")
	ErrRolRequerido      = errors.New("Rol requerido")
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

func (s *Service) Create(u *UsuarioGestor) error {
	if u.Username == nil || *u.Username == "" {
		return ErrUsernameRequerido
	}
	if u.PasswordHash == nil || *u.PasswordHash == "" {
		return ErrPasswordRequerido
	}
	if u.Email == nil || *u.Email == "" {
		return ErrEmailRequerido

	}
	if u.Rol == nil || *u.Rol == "" {
		return ErrRolRequerido
	}
	err := s.storage.Create(u)
	if err != nil {
		s.logger.Error("Errore durante crear", zap.Error(err))
		return fmt.Errorf("Errore durante crear: %w", err)
	}
	s.logger.Info("Usuario creado", zap.Int("id", u.ID))
	return nil
}

func (s *Service) ReadByUsername(username string) (*UsuarioGestor, error) {
	if username == "" {
		return nil, ErrUsernameRequerido
	}
	usuario, err := s.storage.ReadByUsername(username)
	if err != nil {
		s.logger.Error("Errore durante read", zap.String("username", username), zap.Error(err))
		return nil, err
	}
	return usuario, nil
}
