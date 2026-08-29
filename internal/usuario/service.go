package usuario

import (
	"errors"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameRequerido = errors.New("Username requerido")
	ErrPasswordRequerido = errors.New("Password requerido")
	ErrEmailRequerido    = errors.New("Email requerido")
	ErrRolRequerido      = errors.New("Rol requerido")
	ErrIDInvalido        = errors.New("ID inválido")
	ErrModulosRequerido  = errors.New("Debes seleccionar al menos un modulo")
	ErrUsernameInvalido  = errors.New("El username no puede contener espacios")
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
	if strings.ContainsAny(*u.Username, " \t\n\r") {
		return ErrUsernameInvalido
	}
	if u.PasswordHash == nil || *u.PasswordHash == "" {
		return ErrPasswordRequerido
	}
	if u.Email == nil || *u.Email == "" {
		return ErrEmailRequerido

	}
	if u.Modulos == nil || *u.Modulos == "" {
		return ErrModulosRequerido
	}
	rol := calcularRol(*u.Modulos)
	u.Rol = &rol
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(*u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Error al realizar hash de password: %w", err)
	}
	hash := string(hashBytes)
	u.PasswordHash = &hash

	err = s.storage.Create(u)
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

func (s *Service) Read(id int) (*UsuarioGestor, error) {
	if id <= 0 {
		return nil, ErrIDInvalido
	}
	usuario, err := s.storage.Read(id)
	if err != nil {
		s.logger.Error("Errore durante read", zap.Int("id", id), zap.Error(err))
		return nil, err
	}
	return usuario, nil
}

func (s *Service) Delete(id int) error {
	if id <= 0 {
		return ErrIDInvalido
	}
	err := s.storage.Delete(id)
	if err != nil {
		s.logger.Error("Errore durante delete", zap.Int("id", id), zap.Error(err))
		return err
	}

	s.logger.Info("Usuario eliminado", zap.Int("id", id))
	return nil

}

func (s *Service) GetAll() ([]UsuarioGestor, error) {
	usuarios, err := s.storage.GetAll()
	if err != nil {
		s.logger.Error("Error para obtener usuarios", zap.Error(err))
		return nil, err
	}
	return usuarios, nil
}

func (s *Service) Update(id int, fields *UpdateFieldGestor) error {
	if fields.Username != nil {
		if *fields.Username == "" {
			return ErrUsernameRequerido
		}
		if strings.ContainsAny(*fields.Username, " \t\n\r") {
			return ErrUsernameInvalido
		}
	}
	if fields.Email != nil && *fields.Email == "" {
		return ErrEmailRequerido
	}
	if fields.Rol != nil && *fields.Rol == "" {
		return ErrRolRequerido
	}
	if fields.Modulos != nil {
		if *fields.Modulos == "" {
			return ErrModulosRequerido
		}
		rol := calcularRol(*fields.Modulos)
		fields.Rol = &rol
	}
	if fields.PasswordHash != nil {
		if *fields.PasswordHash == "" {
			return ErrPasswordRequerido
		}
		hashBytes, err := bcrypt.GenerateFromPassword(
			[]byte(*fields.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("error al hashear password: %w", err)
		}
		hashed := string(hashBytes)
		fields.PasswordHash = &hashed
	}
	err := s.storage.Update(id, fields)
	if err != nil {
		s.logger.Error("Error al actualizar usuario", zap.Int("id", id), zap.Error(err))
		return err
	}
	s.logger.Info("Usuario actualizado", zap.Int("id", id))

	return nil
}

func calcularRol(modulos string) string {
	var validos []string
	for _, m := range strings.Split(modulos, ",") {
		m = strings.TrimSpace(m)
		if m != "" {
			validos = append(validos, m)
		}
	}
	if len(validos) == 1 {
		return "GESTOR " + strings.ToUpper(validos[0])
	}
	return "GESTOR"
}
