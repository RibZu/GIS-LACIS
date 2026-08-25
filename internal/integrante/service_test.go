package integrante

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// 1. MockStorage: Simulador de la base de datos PostgreSQL para pruebas unitarias
type MockStorage struct{}

func (m *MockStorage) Create(i *Integrante) error { return nil }
func (m *MockStorage) Read(id int) (*Integrante, error) {
	return &Integrante{ID: id, Nombre: "Corina", Apellido: "Abdelahad"}, nil
}
func (m *MockStorage) Delete(id int) error                      { return nil }
func (m *MockStorage) Update(id int, fields UpdateFields) error { return nil }
func (m *MockStorage) GetAll() ([]Integrante, error) {
	return []Integrante{{ID: 1, Nombre: "Corina", Apellido: "Abdelahad"}}, nil
}

func TestCreate_SinNombre_DebeDevolverError(t *testing.T) {

	mockStorage := &MockStorage{}
	logger := zap.NewNop()
	service := NewService(mockStorage, logger)

	nuevo := &Integrante{
		Nombre:          "",
		Apellido:        "Abdelahad",
		RolID:           1,
		Especializacion: "Modelos",
		Contacto:        "corina@email.unsl.edu.ar",
		Descripcion:     "Docente",
	}

	err := service.Create(nuevo)
	assert.Error(t, err)
	assert.Equal(t, ErrNombreRequerido, err)
}

func TestCreate_Exitoso(t *testing.T) {
	mockStorage := &MockStorage{}
	logger := zap.NewNop()
	service := NewService(mockStorage, logger)
	nuevo := &Integrante{
		Nombre:          "Corina",
		Apellido:        "Abdelahad",
		RolID:           1,
		Especializacion: "Modelos",
		Contacto:        "corina@email.unsl.edu.ar",
		Descripcion:     "Docente",
	}
	err := service.Create(nuevo)
	// Esperamos que NO devuelva ningún error (nil)
	assert.NoError(t, err)
}

func TestCreate_NombreConNumeros_DebeDevolverError(t *testing.T) {

	mockStorage := &MockStorage{}
	logger := zap.NewNop()
	service := NewService(mockStorage, logger)

	nuevo := &Integrante{
		Nombre:          "Corina123",
		Apellido:        "Abdelahad",
		RolID:           1,
		Especializacion: "Modelos",
		Contacto:        "docenne@gmail.com",
		Descripcion:     "Docente",
	}
	err := service.Create(nuevo)

	assert.Error(t, err)
	assert.Equal(t, ErrNombreInvalido, err)

}

func TestRead_IDInvalido_DebeDevolverError(t *testing.T) {

	mockStorage := &MockStorage{}
	service := NewService(mockStorage, zap.NewNop())

	_, err := service.Read(0)

	assert.Error(t, err)
	assert.Equal(t, ErrIDInvalido, err)

}

func TestCreate_SinRol_DebeDevolverError(t *testing.T) {

	mockStorage := &MockStorage{}
	service := NewService(mockStorage, zap.NewNop())

	nuevo := &Integrante{
		Nombre:          "Corina",
		Apellido:        "Abdelahad",
		RolID:           0,
		Especializacion: "Modelos",
		Contacto:        "ismae@gmail.com",
		Descripcion:     "Docente",
	}
	err := service.Create(nuevo)

	assert.Error(t, err)
	assert.Equal(t, ErrRolRequerido, err)

}

func TestCreate_SinEspecializacion_DebeDevolverError(t *testing.T) {
	mockStorage := &MockStorage{}
	service := NewService(mockStorage, zap.NewNop())
	nuevo := &Integrante{
		Nombre:          "Corina",
		Apellido:        "Abdelahad",
		RolID:           1,
		Especializacion: "",
		Contacto:        "corina@email.unsl.edu.ar",
		Descripcion:     "Docente",
	}
	err := service.Create(nuevo)
	assert.Error(t, err)
	assert.Equal(t, ErrEspecializacionRequerida, err)
}

func TestCreate_ApellidoConNumeros_DebeDevolverError(t *testing.T) {
	mockStorage := &MockStorage{}
	service := NewService(mockStorage, zap.NewNop())
	nuevo := &Integrante{
		Nombre:          "Corina",
		Apellido:        "Abdelahad123",
		RolID:           1,
		Especializacion: "Modelos",
		Contacto:        "corina@email.unsl.edu.ar",
		Descripcion:     "Docente",
	}
	err := service.Create(nuevo)
	assert.Error(t, err)
	assert.Equal(t, ErrApellidoInvalido, err)
}

func TestRead_Exitoso(t *testing.T) {
	mockStorage := &MockStorage{}
	service := NewService(mockStorage, zap.NewNop())

	IntegranteObj, err := service.Read(1) // valor 1 exitoso

	assert.NoError(t, err)
	assert.NotNil(t, IntegranteObj)

}

func TestGetAll_Exitoso(t *testing.T) {
	mockStorage := &MockStorage{}
	service := NewService(mockStorage, zap.NewNop())
	lista, err := service.GetAll()
	assert.NoError(t, err)
	assert.NotNil(t, lista)
	assert.GreaterOrEqual(t, len(lista), 0)
}

func TestUpdate_IDInvalido_DebeDevolverError(t *testing.T) {
	mockStorage := &MockStorage{}
	service := NewService(mockStorage, zap.NewNop())
	nuevoNombre := "Corina"
	fields := UpdateFields{Nombre: &nuevoNombre}
	err := service.Update(0, fields)
	assert.Error(t, err)
	assert.Equal(t, ErrIDInvalido, err)
}

func TestUpdate_NombreConNumeros_DebeDevolverError(t *testing.T) {
	mockStorage := &MockStorage{}
	service := NewService(mockStorage, zap.NewNop())
	nombreInvalido := "Corina123"
	fields := UpdateFields{Nombre: &nombreInvalido}
	err := service.Update(1, fields)
	assert.Error(t, err)
	assert.Equal(t, ErrNombreInvalido, err)
}

func TestUpdate_Exitoso(t *testing.T) {
	mockStorage := &MockStorage{}
	service := NewService(mockStorage, zap.NewNop())
	nuevoContacto := "nuevo_email@unsl.edu.ar"
	fields := UpdateFields{Contacto: &nuevoContacto}
	err := service.Update(1, fields)
	assert.NoError(t, err)
}

func TestDelete_IDInvalido_DebeDevolverError(t *testing.T) {
	mockStorage := &MockStorage{}
	service := NewService(mockStorage, zap.NewNop())
	err := service.Delete(-1)
	assert.Error(t, err)
	assert.Equal(t, ErrIDInvalido, err)
}

func TestDelete_Exitoso(t *testing.T) {
	mockStorage := &MockStorage{}
	service := NewService(mockStorage, zap.NewNop())
	err := service.Delete(1)
	assert.NoError(t, err)
}
