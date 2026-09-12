package tesis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type MockStorage struct{}

func (m *MockStorage) Create(t *Tesis) error {
	t.ID = 1
	return nil
}
func (m *MockStorage) Read(id int) (*Tesis, error) {
	anio := 2024
	return &Tesis{
		ID:            id,
		Titulo:        "Optimización de Procesos de Calidad de Software",
		Nivel:         "Doctorado",
		CarreraOrigen: "Doctorado en Ingeniería Informática",
		Anio:          &anio,
		AutorHistorico: "Ing. Juan Pérez",
	}, nil
}
func (m *MockStorage) Delete(id int) error                      { return nil }
func (m *MockStorage) Update(id int, fields UpdateFields) error { return nil }
func (m *MockStorage) GetAll() ([]Tesis, error) {
	anio := 2024
	return []Tesis{
		{
			ID:            1,
			Titulo:        "Optimización de Procesos de Calidad de Software",
			Nivel:         "Doctorado",
			CarreraOrigen: "Doctorado en Ingeniería Informática",
			Anio:          &anio,
			AutorHistorico: "Ing. Juan Pérez",
		},
	}, nil
}
func (m *MockStorage) GetByCarrera(carrera string) ([]Tesis, error) {
	return m.GetAll()
}
func (m *MockStorage) GetByNivel(nivel string) ([]Tesis, error) {
	return m.GetAll()
}

func TestCreate_SinTitulo_DebeDevolverError(t *testing.T) {
	mock := &MockStorage{}
	srv := NewService(mock, zap.NewNop())
	anio := 2024

	tesis := &Tesis{
		Titulo:         "",
		Nivel:          "Doctorado",
		CarreraOrigen:  "Doctorado en Ingeniería Informática",
		Anio:           &anio,
		AutorHistorico: "Juan Pérez",
	}

	err := srv.Create(tesis)
	assert.Error(t, err)
	assert.Equal(t, ErrTituloRequerido, err)
}

func TestCreate_SinNivel_DebeDevolverError(t *testing.T) {
	mock := &MockStorage{}
	srv := NewService(mock, zap.NewNop())

	tesis := &Tesis{
		Titulo:         "Título de Tesis",
		Nivel:          "",
		CarreraOrigen:  "Doctorado en Ingeniería Informática",
		AutorHistorico: "Juan Pérez",
	}

	err := srv.Create(tesis)
	assert.Error(t, err)
	assert.Equal(t, ErrNivelRequerido, err)
}

func TestCreate_SinCarrera_DebeDevolverError(t *testing.T) {
	mock := &MockStorage{}
	srv := NewService(mock, zap.NewNop())

	tesis := &Tesis{
		Titulo:         "Título de Tesis",
		Nivel:          "Doctorado",
		CarreraOrigen:  "",
		AutorHistorico: "Juan Pérez",
	}

	err := srv.Create(tesis)
	assert.Error(t, err)
	assert.Equal(t, ErrCarreraRequerida, err)
}

func TestCreate_SinAutor_DebeDevolverError(t *testing.T) {
	mock := &MockStorage{}
	srv := NewService(mock, zap.NewNop())

	tesis := &Tesis{
		Titulo:        "Título de Tesis",
		Nivel:         "Doctorado",
		CarreraOrigen: "Doctorado en Ingeniería Informática",
	}

	err := srv.Create(tesis)
	assert.Error(t, err)
	assert.Equal(t, ErrAutorRequerido, err)
}

func TestCreate_AnioInvalido_DebeDevolverError(t *testing.T) {
	mock := &MockStorage{}
	srv := NewService(mock, zap.NewNop())
	anioInvalido := 1800

	tesis := &Tesis{
		Titulo:         "Título de Tesis",
		Nivel:          "Doctorado",
		CarreraOrigen:  "Doctorado en Ingeniería Informática",
		Anio:           &anioInvalido,
		AutorHistorico: "Juan Pérez",
	}

	err := srv.Create(tesis)
	assert.Error(t, err)
	assert.Equal(t, ErrAnioInvalido, err)
}

func TestCreate_Exitoso(t *testing.T) {
	mock := &MockStorage{}
	srv := NewService(mock, zap.NewNop())
	anio := 2024
	autorID := 5

	tesis := &Tesis{
		Titulo:        "Investigación sobre Calidad de Software",
		Nivel:         "Maestría",
		CarreraOrigen: "Maestría en Calidad de Software",
		Anio:          &anio,
		AutorID:       &autorID,
		PalabrasClave: "Testing, Calidad, Métricas",
		Resumen:       "Resumen detallado de la tesis",
	}

	err := srv.Create(tesis)
	assert.NoError(t, err)
	assert.Equal(t, 1, tesis.ID)
}

func TestRead_IDInvalido_DebeDevolverError(t *testing.T) {
	mock := &MockStorage{}
	srv := NewService(mock, zap.NewNop())

	_, err := srv.Read(0)
	assert.Error(t, err)
	assert.Equal(t, ErrIDInvalido, err)
}

func TestRead_Exitoso(t *testing.T) {
	mock := &MockStorage{}
	srv := NewService(mock, zap.NewNop())

	tesis, err := srv.Read(1)
	assert.NoError(t, err)
	assert.NotNil(t, tesis)
	assert.Equal(t, 1, tesis.ID)
}

func TestGetAll_Exitoso(t *testing.T) {
	mock := &MockStorage{}
	srv := NewService(mock, zap.NewNop())

	lista, err := srv.GetAll()
	assert.NoError(t, err)
	assert.Len(t, lista, 1)
}

func TestUpdate_IDInvalido_DebeDevolverError(t *testing.T) {
	mock := &MockStorage{}
	srv := NewService(mock, zap.NewNop())

	err := srv.Update(0, UpdateFields{})
	assert.Error(t, err)
	assert.Equal(t, ErrIDInvalido, err)
}

func TestUpdate_TituloVacio_DebeDevolverError(t *testing.T) {
	mock := &MockStorage{}
	srv := NewService(mock, zap.NewNop())
	tituloVacio := "   "

	err := srv.Update(1, UpdateFields{Titulo: &tituloVacio})
	assert.Error(t, err)
	assert.Equal(t, ErrTituloRequerido, err)
}

func TestUpdate_Exitoso(t *testing.T) {
	mock := &MockStorage{}
	srv := NewService(mock, zap.NewNop())
	nuevoTitulo := "Nuevo Título de Tesis Actualizado"

	err := srv.Update(1, UpdateFields{Titulo: &nuevoTitulo})
	assert.NoError(t, err)
}

func TestDelete_IDInvalido_DebeDevolverError(t *testing.T) {
	mock := &MockStorage{}
	srv := NewService(mock, zap.NewNop())

	err := srv.Delete(-1)
	assert.Error(t, err)
	assert.Equal(t, ErrIDInvalido, err)
}

func TestDelete_Exitoso(t *testing.T) {
	mock := &MockStorage{}
	srv := NewService(mock, zap.NewNop())

	err := srv.Delete(1)
	assert.NoError(t, err)
}
