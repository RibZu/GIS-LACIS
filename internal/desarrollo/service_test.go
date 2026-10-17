package desarrollo

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"PaginaSEG/internal/integrante"

	"go.uber.org/zap"
)

func ptr[T any](v T) *T {
	return &v
}

type mockStorage struct {
	llamadas map[string]int

	ultimoCrear    *Desarrollo
	crearIDs       []int
	crearExternos  []string
	ultimoActID    int
	ultimoActCampo UpdateFields
	actIDs         []int
	actExternos    []string

	errCrear       error
	errActualizar  error
	errDelete      error
	errRead        error
	errIntegrantes error
	errExternos    error

	integrantesPorDesarrollo map[int][]integrante.Integrante
	externosPorDesarrollo    map[int][]string
}

func nuevoMock() *mockStorage {
	return &mockStorage{llamadas: map[string]int{}}
}

func (m *mockStorage) Read(id int) (*Desarrollo, error) {
	m.llamadas["Read"]++
	if m.errRead != nil {
		return nil, m.errRead
	}
	return &Desarrollo{ID: id, Titulo: "x", Anio: 2024}, nil
}

func (m *mockStorage) Delete(id int) error {
	m.llamadas["Delete"]++
	return m.errDelete
}

func (m *mockStorage) GetAll() ([]Desarrollo, error) {
	m.llamadas["GetAll"]++
	return nil, nil
}

func (m *mockStorage) CrearCompleto(d *Desarrollo, ids []int, externos []string) error {
	m.llamadas["CrearCompleto"]++
	copia := *d
	m.ultimoCrear = &copia
	m.crearIDs = ids
	m.crearExternos = externos
	if m.errCrear != nil {
		return m.errCrear
	}
	d.ID = 42
	return nil
}

func (m *mockStorage) ActualizarCompleto(id int, fields UpdateFields, ids []int, externos []string) error {
	m.llamadas["ActualizarCompleto"]++
	m.ultimoActID = id
	m.ultimoActCampo = fields
	m.actIDs = ids
	m.actExternos = externos
	return m.errActualizar
}

func (m *mockStorage) GetIntegrantes(id int) ([]integrante.Integrante, error) {
	m.llamadas["GetIntegrantes"]++
	return nil, nil
}

func (m *mockStorage) GetParticipantesExternos(id int) ([]string, error) {
	m.llamadas["GetParticipantesExternos"]++
	return nil, nil
}

func (m *mockStorage) IntegrantesPorDesarrollo() (map[int][]integrante.Integrante, error) {
	m.llamadas["IntegrantesPorDesarrollo"]++
	return m.integrantesPorDesarrollo, m.errIntegrantes
}

func (m *mockStorage) ParticipantesExternosPorDesarrollo() (map[int][]string, error) {
	m.llamadas["ParticipantesExternosPorDesarrollo"]++
	return m.externosPorDesarrollo, m.errExternos
}

func servicio(m *mockStorage) *Service {
	return NewService(m, zap.NewNop())
}

func productoValido() *Desarrollo {
	return &Desarrollo{Titulo: "Sistema de prueba", Anio: 2024, URL: "ejemplo.com.ar", Contacto: "a@b.com", Descripcion: "desc"}
}

func TestCrearCompleto_AnioInvalidoNoLlegaAlStorage(t *testing.T) {
	anioActual := time.Now().Year()
	for _, anio := range []int{0, -5, 1989, anioActual + 2, 99999} {
		m := nuevoMock()
		d := productoValido()
		d.Anio = anio
		if err := servicio(m).CrearCompleto(d, nil, nil); !errors.Is(err, ErrAnioInvalido) {
			t.Errorf("año %d: se esperaba ErrAnioInvalido, llegó %v", anio, err)
		}
		if m.llamadas["CrearCompleto"] != 0 {
			t.Errorf("año %d: no se debería llamar al storage", anio)
		}
	}
}

func TestCrearCompleto_LosExtremosDelAnioSonValidos(t *testing.T) {
	for _, anio := range []int{1990, time.Now().Year() + 1} {
		m := nuevoMock()
		d := productoValido()
		d.Anio = anio
		if err := servicio(m).CrearCompleto(d, nil, nil); err != nil {
			t.Errorf("año %d debería ser válido, llegó %v", anio, err)
		}
	}
}

func TestCrearCompleto_Titulo(t *testing.T) {
	casos := map[string]error{
		"":                            ErrTituloRequerido,
		"    ":                        ErrTituloRequerido,
		"\t\n":                        ErrTituloRequerido,
		strings.Repeat("a", 256):      ErrTituloLargo,
		strings.Repeat("ñ", 256):      ErrTituloLargo,
		strings.Repeat("a", 255):      nil,
		strings.Repeat("ñ", 255):      nil,
		"Sistema con acentos y ñandú": nil,
	}
	for titulo, esperado := range casos {
		m := nuevoMock()
		d := productoValido()
		d.Titulo = titulo
		err := servicio(m).CrearCompleto(d, nil, nil)
		if esperado == nil && err != nil {
			t.Errorf("título de %d caracteres: no debería fallar, llegó %v", len([]rune(titulo)), err)
		}
		if esperado != nil && !errors.Is(err, esperado) {
			t.Errorf("título de %d caracteres: se esperaba %v, llegó %v", len([]rune(titulo)), esperado, err)
		}
		if esperado != nil && m.llamadas["CrearCompleto"] != 0 {
			t.Errorf("título inválido: no se debería llamar al storage")
		}
	}
}

func TestCrearCompleto_NormalizaLaURL(t *testing.T) {
	casos := map[string]string{
		"ejemplo.com.ar":         "https://ejemplo.com.ar",
		"  ejemplo.com.ar  ":     "https://ejemplo.com.ar",
		"http://ejemplo.com":     "http://ejemplo.com",
		"https://github.com/x/y": "https://github.com/x/y",
		"":                       "",
		"    ":                   "",
	}
	for entrada, esperada := range casos {
		m := nuevoMock()
		d := productoValido()
		d.URL = entrada
		if err := servicio(m).CrearCompleto(d, nil, nil); err != nil {
			t.Fatalf("URL %q: llegó %v", entrada, err)
		}
		if m.ultimoCrear.URL != esperada {
			t.Errorf("URL %q: se esperaba %q, llegó %q", entrada, esperada, m.ultimoCrear.URL)
		}
	}
}

func TestCrearCompleto_LimitesDeLongitud(t *testing.T) {
	m := nuevoMock()
	d := productoValido()
	d.URL = strings.Repeat("a", 495)
	if err := servicio(m).CrearCompleto(d, nil, nil); !errors.Is(err, ErrURLLarga) {
		t.Errorf("URL larga (con el https:// agregado): se esperaba ErrURLLarga, llegó %v", err)
	}

	d = productoValido()
	d.Contacto = strings.Repeat("c", 256)
	if err := servicio(m).CrearCompleto(d, nil, nil); !errors.Is(err, ErrContactoLargo) {
		t.Errorf("contacto largo: se esperaba ErrContactoLargo, llegó %v", err)
	}
	if m.llamadas["CrearCompleto"] != 0 {
		t.Error("ninguna validación fallida debería llegar al storage")
	}
}

func TestCrearCompleto_NormalizaExternos(t *testing.T) {
	m := nuevoMock()
	externos := []string{"  Ana   Pérez ", "ana pérez", "", "   ", "Luis", "LUIS", "Marta"}
	if err := servicio(m).CrearCompleto(productoValido(), nil, externos); err != nil {
		t.Fatalf("llegó %v", err)
	}
	esperado := []string{"Ana Pérez", "Luis", "Marta"}
	if !reflect.DeepEqual(m.crearExternos, esperado) {
		t.Fatalf("se esperaba %v, llegó %v", esperado, m.crearExternos)
	}
}

func TestCrearCompleto_ExternoDemasiadoLargo(t *testing.T) {
	m := nuevoMock()
	err := servicio(m).CrearCompleto(productoValido(), nil, []string{"Ok", strings.Repeat("x", 151)})
	if !errors.Is(err, ErrParticipanteLargo) {
		t.Fatalf("se esperaba ErrParticipanteLargo, llegó %v", err)
	}
	if m.llamadas["CrearCompleto"] != 0 {
		t.Fatal("no se debería llamar al storage")
	}
	if err := servicio(m).CrearCompleto(productoValido(), nil, []string{strings.Repeat("x", 150)}); err != nil {
		t.Fatalf("150 caracteres debería ser válido, llegó %v", err)
	}
}

func TestCrearCompleto_SinParticipantesEsValido(t *testing.T) {
	m := nuevoMock()
	if err := servicio(m).CrearCompleto(productoValido(), nil, nil); err != nil {
		t.Fatalf("un producto sin participantes debería guardarse, llegó %v", err)
	}
	if len(m.crearIDs) != 0 || len(m.crearExternos) != 0 {
		t.Fatalf("no debería haber participantes, llegó ids=%v externos=%v", m.crearIDs, m.crearExternos)
	}
}

func TestCrearCompleto_NormalizaIDsDeIntegrantes(t *testing.T) {
	m := nuevoMock()
	if err := servicio(m).CrearCompleto(productoValido(), []int{3, 3, 0, -1, 5, 3}, nil); err != nil {
		t.Fatalf("llegó %v", err)
	}
	if !reflect.DeepEqual(m.crearIDs, []int{3, 5}) {
		t.Fatalf("se esperaba [3 5], llegó %v", m.crearIDs)
	}
}

func TestCrearCompleto_LlamaUnaSolaVezAlStorage(t *testing.T) {
	m := nuevoMock()
	d := productoValido()
	if err := servicio(m).CrearCompleto(d, []int{1, 2}, []string{"Ana"}); err != nil {
		t.Fatalf("llegó %v", err)
	}
	if m.llamadas["CrearCompleto"] != 1 {
		t.Fatalf("se esperaba una sola llamada, hubo %d", m.llamadas["CrearCompleto"])
	}
	if d.ID != 42 {
		t.Fatalf("el ID asignado por el storage debería quedar en el producto, es %d", d.ID)
	}
	if len(m.llamadas) != 1 {
		t.Fatalf("solo debería usarse CrearCompleto, se usaron %v", m.llamadas)
	}
}

func TestCrearCompleto_ErroresDelStorageSePropagan(t *testing.T) {
	falla := errors.New("base caída")
	m := nuevoMock()
	m.errCrear = falla
	if err := servicio(m).CrearCompleto(productoValido(), nil, nil); !errors.Is(err, falla) {
		t.Fatalf("se esperaba el error del storage, llegó %v", err)
	}

	m = nuevoMock()
	m.errCrear = ErrIntegranteInexistente
	if err := servicio(m).CrearCompleto(productoValido(), []int{999}, nil); !errors.Is(err, ErrIntegranteInexistente) {
		t.Fatalf("se esperaba ErrIntegranteInexistente, llegó %v", err)
	}
}

func TestActualizarCompleto_IDInvalido(t *testing.T) {
	m := nuevoMock()
	for _, id := range []int{0, -1} {
		if err := servicio(m).ActualizarCompleto(id, UpdateFields{}, nil, nil); !errors.Is(err, ErrIDInvalido) {
			t.Errorf("id %d: se esperaba ErrIDInvalido, llegó %v", id, err)
		}
	}
	if m.llamadas["ActualizarCompleto"] != 0 {
		t.Error("no se debería llamar al storage")
	}
}

func TestActualizarCompleto_Validaciones(t *testing.T) {
	m := nuevoMock()
	s := servicio(m)
	if err := s.ActualizarCompleto(1, UpdateFields{Titulo: ptr("  ")}, nil, nil); !errors.Is(err, ErrTituloRequerido) {
		t.Errorf("título en blanco: llegó %v", err)
	}
	if err := s.ActualizarCompleto(1, UpdateFields{Anio: ptr(1800)}, nil, nil); !errors.Is(err, ErrAnioInvalido) {
		t.Errorf("año inválido: llegó %v", err)
	}
	if err := s.ActualizarCompleto(1, UpdateFields{Contacto: ptr(strings.Repeat("c", 300))}, nil, nil); !errors.Is(err, ErrContactoLargo) {
		t.Errorf("contacto largo: llegó %v", err)
	}
	if err := s.ActualizarCompleto(1, UpdateFields{URL: ptr(strings.Repeat("u", 600))}, nil, nil); !errors.Is(err, ErrURLLarga) {
		t.Errorf("URL larga: llegó %v", err)
	}
	if err := s.ActualizarCompleto(1, UpdateFields{}, nil, []string{strings.Repeat("e", 200)}); !errors.Is(err, ErrParticipanteLargo) {
		t.Errorf("externo largo: llegó %v", err)
	}
	if m.llamadas["ActualizarCompleto"] != 0 {
		t.Error("ninguna validación fallida debería llegar al storage")
	}
}

func TestActualizarCompleto_EnviaTodoJuntoYNormalizado(t *testing.T) {
	m := nuevoMock()
	campos := UpdateFields{Titulo: ptr("Nuevo"), Anio: ptr(2025), URL: ptr("repo.com"), Contacto: ptr("x@y.com"), Descripcion: ptr("d")}
	err := servicio(m).ActualizarCompleto(7, campos, []int{2, 2, 4}, []string{" Ana ", "ana"})
	if err != nil {
		t.Fatalf("llegó %v", err)
	}
	if m.llamadas["ActualizarCompleto"] != 1 {
		t.Fatalf("se esperaba una sola llamada, hubo %d", m.llamadas["ActualizarCompleto"])
	}
	if m.ultimoActID != 7 {
		t.Errorf("id enviado: %d", m.ultimoActID)
	}
	if *m.ultimoActCampo.URL != "https://repo.com" {
		t.Errorf("la URL debería normalizarse, llegó %q", *m.ultimoActCampo.URL)
	}
	if !reflect.DeepEqual(m.actIDs, []int{2, 4}) {
		t.Errorf("ids: %v", m.actIDs)
	}
	if !reflect.DeepEqual(m.actExternos, []string{"Ana"}) {
		t.Errorf("externos: %v", m.actExternos)
	}
}

func TestActualizarCompleto_ErroresSePropagan(t *testing.T) {
	for _, esperado := range []error{ErrNotFound, ErrIntegranteInexistente, errors.New("base caída")} {
		m := nuevoMock()
		m.errActualizar = esperado
		if err := servicio(m).ActualizarCompleto(1, UpdateFields{}, nil, nil); !errors.Is(err, esperado) {
			t.Errorf("se esperaba %v, llegó %v", esperado, err)
		}
	}
}

func TestDelete(t *testing.T) {
	m := nuevoMock()
	s := servicio(m)
	if err := s.Delete(0); !errors.Is(err, ErrIDInvalido) {
		t.Errorf("id 0: llegó %v", err)
	}
	if m.llamadas["Delete"] != 0 {
		t.Error("un id inválido no debería llegar al storage")
	}
	if err := s.Delete(3); err != nil {
		t.Errorf("borrado normal: llegó %v", err)
	}
	m.errDelete = ErrNotFound
	if err := s.Delete(3); !errors.Is(err, ErrNotFound) {
		t.Errorf("inexistente: llegó %v", err)
	}
	falla := errors.New("base caída")
	m.errDelete = falla
	if err := s.Delete(3); !errors.Is(err, falla) {
		t.Errorf("falla del storage: llegó %v", err)
	}
}

func TestRead(t *testing.T) {
	m := nuevoMock()
	s := servicio(m)
	if _, err := s.Read(-1); !errors.Is(err, ErrIDInvalido) {
		t.Errorf("id inválido: llegó %v", err)
	}
	if d, err := s.Read(5); err != nil || d.ID != 5 {
		t.Errorf("lectura normal: llegó %v %v", d, err)
	}
	m.errRead = ErrNotFound
	if _, err := s.Read(5); !errors.Is(err, ErrNotFound) {
		t.Errorf("inexistente: llegó %v", err)
	}
}

func TestParticipantesPorDesarrollo(t *testing.T) {
	m := nuevoMock()
	m.integrantesPorDesarrollo = map[int][]integrante.Integrante{1: {{ID: 9, Nombre: "Ana", Apellido: "Gómez"}}}
	m.externosPorDesarrollo = map[int][]string{1: {"Luis"}, 2: {"Marta"}}
	integrantes, externos, err := servicio(m).ParticipantesPorDesarrollo()
	if err != nil {
		t.Fatalf("llegó %v", err)
	}
	if len(integrantes[1]) != 1 || integrantes[1][0].ID != 9 {
		t.Errorf("integrantes: %v", integrantes)
	}
	if !reflect.DeepEqual(externos[2], []string{"Marta"}) {
		t.Errorf("externos: %v", externos)
	}
	if m.llamadas["IntegrantesPorDesarrollo"] != 1 || m.llamadas["ParticipantesExternosPorDesarrollo"] != 1 {
		t.Errorf("deberían ser dos lecturas en total, hubo %v", m.llamadas)
	}
}

func TestParticipantesPorDesarrollo_UnErrorNoDaResultadosParciales(t *testing.T) {
	falla := errors.New("base caída")

	m := nuevoMock()
	m.errIntegrantes = falla
	integrantes, externos, err := servicio(m).ParticipantesPorDesarrollo()
	if !errors.Is(err, falla) || integrantes != nil || externos != nil {
		t.Errorf("con falla en integrantes: llegó %v %v %v", integrantes, externos, err)
	}

	m = nuevoMock()
	m.errExternos = falla
	integrantes, externos, err = servicio(m).ParticipantesPorDesarrollo()
	if !errors.Is(err, falla) || integrantes != nil || externos != nil {
		t.Errorf("con falla en externos: llegó %v %v %v", integrantes, externos, err)
	}
}

func TestTieneRepositorio(t *testing.T) {
	casos := map[string]bool{"": false, "   ": false, "https://x.com": true, "x.com": true}
	for url, esperado := range casos {
		if TieneRepositorio(url) != esperado {
			t.Errorf("TieneRepositorio(%q) debería ser %v", url, esperado)
		}
		if (Desarrollo{URL: url}).TieneRepositorio() != esperado {
			t.Errorf("Desarrollo{URL: %q}.TieneRepositorio() debería ser %v", url, esperado)
		}
	}
}

func TestQuitarPrefijoEst(t *testing.T) {
	casos := map[string]string{
		"Est. Juan Pérez": "Juan Pérez",
		"est Juan":        "Juan",
		"Est. Est. Ana":   "Ana",
		"Estela Ruiz":     "Estela Ruiz",
		"Juan Pérez":      "Juan Pérez",
		"Est.":            "",
	}
	for entrada, esperado := range casos {
		if got := QuitarPrefijoEst(entrada); got != esperado {
			t.Errorf("QuitarPrefijoEst(%q) = %q, se esperaba %q", entrada, got, esperado)
		}
	}
}
