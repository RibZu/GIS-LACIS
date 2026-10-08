package configuracion_test

import (
	"testing"

	"PaginaSEG/internal/configuracion"

	"go.uber.org/zap"
)

type mockStorage struct {
	cfg *configuracion.ConfiguracionSitio
}

func (m *mockStorage) Get() (*configuracion.ConfiguracionSitio, error) {
	if m.cfg == nil {
		m.cfg = &configuracion.ConfiguracionSitio{
			EmailDoctorado: "agaris@gmail.com",
		}
	}
	return m.cfg, nil
}

func (m *mockStorage) Update(cfg *configuracion.ConfiguracionSitio) error {
	m.cfg = cfg
	return nil
}

func TestGuardarConfiguracion_EmailsValidos(t *testing.T) {
	storage := &mockStorage{}
	logger := zap.NewNop()
	svc := configuracion.NewService(storage, logger)

	docMail := "nuevo_doctorado@unsl.edu.ar"
	dto := configuracion.UpdateConfiguracionDTO{
		EmailDoctorado: &docMail,
	}

	updated, err := svc.GuardarConfiguracion(dto)
	if err != nil {
		t.Fatalf("se esperaba éxito al guardar, se obtuvo error: %v", err)
	}

	if updated.EmailDoctorado != docMail {
		t.Errorf("esperado %s, obtenido %s", docMail, updated.EmailDoctorado)
	}
}

func TestGuardarConfiguracion_EmailInvalido(t *testing.T) {
	storage := &mockStorage{}
	logger := zap.NewNop()
	svc := configuracion.NewService(storage, logger)

	invalido := "correo-invalido-sin-arroba"
	dto := configuracion.UpdateConfiguracionDTO{
		EmailDoctorado: &invalido,
	}

	_, err := svc.GuardarConfiguracion(dto)
	if err == nil {
		t.Fatal("se esperaba error de validación de email, pero no ocurrió")
	}
}
