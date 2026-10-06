package usuario

import (
	"errors"
	"testing"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func ptr[T any](v T) *T {
	return &v
}

type mockStorage struct {
	usuarios     map[int]*UsuarioGestor
	llamadas     []string
	ultimoCreate *UsuarioGestor
	ultimoUpdate *UpdateFieldGestor
	errCreate    error
	errUpdate    error
	errDelete    error
	errContar    error
}

func nuevoMock(usuarios ...*UsuarioGestor) *mockStorage {
	m := &mockStorage{usuarios: map[int]*UsuarioGestor{}}
	for _, u := range usuarios {
		m.usuarios[u.ID] = u
	}
	return m
}

func (m *mockStorage) llamo(nombre string) int {
	n := 0
	for _, l := range m.llamadas {
		if l == nombre {
			n++
		}
	}
	return n
}

func (m *mockStorage) Create(u *UsuarioGestor) error {
	m.llamadas = append(m.llamadas, "Create")
	m.ultimoCreate = u
	if m.errCreate != nil {
		return m.errCreate
	}
	u.ID = len(m.usuarios) + 100
	m.usuarios[u.ID] = u
	return nil
}

func (m *mockStorage) Read(id int) (*UsuarioGestor, error) {
	m.llamadas = append(m.llamadas, "Read")
	u, ok := m.usuarios[id]
	if !ok {
		return nil, ErrNotFound
	}
	return u, nil
}

func (m *mockStorage) ReadByUsername(username string) (*UsuarioGestor, error) {
	for _, u := range m.usuarios {
		if u.Username != nil && *u.Username == username {
			return u, nil
		}
	}
	return nil, ErrNotFound
}

func (m *mockStorage) GetAll() ([]UsuarioGestor, error) {
	var todos []UsuarioGestor
	for _, u := range m.usuarios {
		todos = append(todos, *u)
	}
	return todos, nil
}

func (m *mockStorage) Update(id int, fields *UpdateFieldGestor) error {
	m.llamadas = append(m.llamadas, "Update")
	m.ultimoUpdate = fields
	return m.errUpdate
}

func (m *mockStorage) Delete(id int) error {
	m.llamadas = append(m.llamadas, "Delete")
	if m.errDelete != nil {
		return m.errDelete
	}
	delete(m.usuarios, id)
	return nil
}

func (m *mockStorage) ContarAdmins() (int, error) {
	m.llamadas = append(m.llamadas, "ContarAdmins")
	if m.errContar != nil {
		return 0, m.errContar
	}
	total := 0
	for _, u := range m.usuarios {
		if u.Rol != nil && *u.Rol == "ADMIN" {
			total++
		}
	}
	return total, nil
}

func admin(id int) *UsuarioGestor {
	return &UsuarioGestor{ID: id, Username: ptr("admin"), Rol: ptr("ADMIN")}
}

func gestor(id int, modulos string) *UsuarioGestor {
	return &UsuarioGestor{ID: id, Username: ptr("gestor"), Rol: ptr(calcularRol(modulos)), Modulos: ptr(modulos)}
}

func servicio(m *mockStorage) *Service {
	return NewService(m, zap.NewNop())
}

func TestDelete_AutoEliminacionSeBloquea(t *testing.T) {
	m := nuevoMock(admin(1), admin(2))
	err := servicio(m).Delete(1, 1)
	if !errors.Is(err, ErrAutoEliminacion) {
		t.Fatalf("se esperaba ErrAutoEliminacion, llegó %v", err)
	}
	if m.llamo("Delete") != 0 {
		t.Fatal("no se debería llamar a Delete del storage")
	}
}

func TestDelete_NoSePuedeBorrarAlUnicoAdmin(t *testing.T) {
	m := nuevoMock(admin(1), gestor(2, "estadisticas"))
	err := servicio(m).Delete(2, 1)
	if !errors.Is(err, ErrUltimoAdmin) {
		t.Fatalf("se esperaba ErrUltimoAdmin, llegó %v", err)
	}
	if m.llamo("Delete") != 0 {
		t.Fatal("no se debería llamar a Delete del storage")
	}
	if _, existe := m.usuarios[1]; !existe {
		t.Fatal("el único administrador no puede desaparecer")
	}
}

func TestDelete_SePuedeBorrarUnAdminSiHayOtro(t *testing.T) {
	m := nuevoMock(admin(1), admin(2))
	if err := servicio(m).Delete(1, 2); err != nil {
		t.Fatalf("con dos administradores se puede borrar uno, llegó %v", err)
	}
	if _, existe := m.usuarios[2]; existe {
		t.Fatal("el usuario 2 debería haberse borrado")
	}
	if _, existe := m.usuarios[1]; !existe {
		t.Fatal("el usuario 1 debería seguir existiendo")
	}
}

func TestDelete_UnGestorSeBorraAunqueHayaUnSoloAdmin(t *testing.T) {
	m := nuevoMock(admin(1), gestor(2, "tesis"))
	if err := servicio(m).Delete(1, 2); err != nil {
		t.Fatalf("un gestor se debería poder borrar, llegó %v", err)
	}
	if m.llamo("ContarAdmins") != 0 {
		t.Fatal("para un gestor no hace falta contar administradores")
	}
}

func TestDelete_Inexistente(t *testing.T) {
	m := nuevoMock(admin(1))
	if err := servicio(m).Delete(1, 99); !errors.Is(err, ErrNotFound) {
		t.Fatalf("se esperaba ErrNotFound, llegó %v", err)
	}
	if m.llamo("Delete") != 0 {
		t.Fatal("no se debería llamar a Delete del storage")
	}
}

func TestDelete_IDInvalido(t *testing.T) {
	m := nuevoMock(admin(1))
	for _, id := range []int{0, -3} {
		if err := servicio(m).Delete(1, id); !errors.Is(err, ErrIDInvalido) {
			t.Errorf("id %d: se esperaba ErrIDInvalido, llegó %v", id, err)
		}
	}
}

func TestDelete_ErrorDelStorageSePropaga(t *testing.T) {
	falla := errors.New("base caída")
	m := nuevoMock(admin(1), gestor(2, "tesis"))
	m.errDelete = falla
	if err := servicio(m).Delete(1, 2); !errors.Is(err, falla) {
		t.Fatalf("se esperaba el error del storage, llegó %v", err)
	}
	if _, existe := m.usuarios[2]; !existe {
		t.Fatal("si el borrado falla, el usuario debe seguir existiendo")
	}
}

func TestDelete_ErrorAlContarAdminsSePropaga(t *testing.T) {
	falla := errors.New("base caída")
	m := nuevoMock(admin(1), admin(2))
	m.errContar = falla
	if err := servicio(m).Delete(1, 2); !errors.Is(err, falla) {
		t.Fatalf("se esperaba el error de ContarAdmins, llegó %v", err)
	}
	if m.llamo("Delete") != 0 {
		t.Fatal("si no se pudo verificar, no se debe borrar")
	}
}

func TestUpdate_AdminNoPuedeQuitarseSuPropioRol(t *testing.T) {
	m := nuevoMock(admin(1), admin(2))
	err := servicio(m).Update(1, 1, &UpdateFieldGestor{Modulos: ptr("integrantes")})
	if !errors.Is(err, ErrAutoDegradacion) {
		t.Fatalf("se esperaba ErrAutoDegradacion, llegó %v", err)
	}
	if m.llamo("Update") != 0 {
		t.Fatal("no se debería llamar a Update del storage")
	}
}

func TestUpdate_NoSePuedeDegradarAlUltimoAdmin(t *testing.T) {
	m := nuevoMock(admin(1), gestor(2, "tesis"))
	err := servicio(m).Update(2, 1, &UpdateFieldGestor{Modulos: ptr("integrantes")})
	if !errors.Is(err, ErrUltimoAdmin) {
		t.Fatalf("se esperaba ErrUltimoAdmin, llegó %v", err)
	}
	if m.llamo("Update") != 0 {
		t.Fatal("no se debería llamar a Update del storage")
	}
}

func TestUpdate_SePuedeDegradarUnAdminSiHayOtro(t *testing.T) {
	m := nuevoMock(admin(1), admin(2))
	err := servicio(m).Update(1, 2, &UpdateFieldGestor{Modulos: ptr("integrantes")})
	if err != nil {
		t.Fatalf("con dos administradores se puede degradar uno, llegó %v", err)
	}
	if m.ultimoUpdate == nil || m.ultimoUpdate.Rol == nil || *m.ultimoUpdate.Rol != "GESTOR INTEGRANTES" {
		t.Fatalf("el rol debería quedar GESTOR INTEGRANTES, quedó %v", m.ultimoUpdate)
	}
}

func TestUpdate_EditarUnAdminSinModulosNoTocaElRol(t *testing.T) {
	m := nuevoMock(admin(1))
	err := servicio(m).Update(1, 1, &UpdateFieldGestor{Username: ptr("jefe"), Email: ptr("jefe@lacis.com")})
	if err != nil {
		t.Fatalf("editar los datos de un administrador debería funcionar, llegó %v", err)
	}
	if m.ultimoUpdate.Rol != nil || m.ultimoUpdate.Modulos != nil {
		t.Fatalf("el rol y los módulos no se deben enviar al storage, llegó rol=%v modulos=%v", m.ultimoUpdate.Rol, m.ultimoUpdate.Modulos)
	}
	if m.llamo("ContarAdmins") != 0 {
		t.Fatal("sin cambio de rol no hace falta contar administradores")
	}
}

func TestUpdate_GestorRecalculaElRolConSusModulos(t *testing.T) {
	casos := []struct {
		modulos string
		rol     string
	}{
		{"tesis", "GESTOR TESIS"},
		{"tesis,integrantes", "GESTOR"},
		{"estadisticas", "GESTOR ESTADISTICAS"},
	}
	for _, c := range casos {
		m := nuevoMock(admin(1), gestor(2, "proyectos"))
		if err := servicio(m).Update(1, 2, &UpdateFieldGestor{Modulos: ptr(c.modulos)}); err != nil {
			t.Fatalf("%s: llegó %v", c.modulos, err)
		}
		if got := *m.ultimoUpdate.Rol; got != c.rol {
			t.Errorf("%s: el rol debería ser %q, es %q", c.modulos, c.rol, got)
		}
	}
}

func TestUpdate_ModulosVaciosSeRechazan(t *testing.T) {
	m := nuevoMock(admin(1), gestor(2, "tesis"))
	if err := servicio(m).Update(1, 2, &UpdateFieldGestor{Modulos: ptr("")}); !errors.Is(err, ErrModulosRequerido) {
		t.Fatalf("se esperaba ErrModulosRequerido, llegó %v", err)
	}
}

func TestUpdate_ValidacionesDeCampos(t *testing.T) {
	m := nuevoMock(admin(1), gestor(2, "tesis"))
	s := servicio(m)
	if err := s.Update(1, 2, &UpdateFieldGestor{Username: ptr("")}); !errors.Is(err, ErrUsernameRequerido) {
		t.Errorf("usuario vacío: llegó %v", err)
	}
	if err := s.Update(1, 2, &UpdateFieldGestor{Username: ptr("con espacio")}); !errors.Is(err, ErrUsernameInvalido) {
		t.Errorf("usuario con espacio: llegó %v", err)
	}
	if err := s.Update(1, 2, &UpdateFieldGestor{Email: ptr("")}); !errors.Is(err, ErrEmailRequerido) {
		t.Errorf("email vacío: llegó %v", err)
	}
	if err := s.Update(1, 2, &UpdateFieldGestor{PasswordHash: ptr("")}); !errors.Is(err, ErrPasswordRequerido) {
		t.Errorf("contraseña vacía: llegó %v", err)
	}
	if err := s.Update(1, 0, &UpdateFieldGestor{}); !errors.Is(err, ErrIDInvalido) {
		t.Errorf("id 0: llegó %v", err)
	}
	if m.llamo("Update") != 0 {
		t.Error("ninguna validación fallida debería llegar al storage")
	}
}

func TestUpdate_UsuarioInexistenteConModulos(t *testing.T) {
	m := nuevoMock(admin(1))
	if err := servicio(m).Update(1, 99, &UpdateFieldGestor{Modulos: ptr("tesis")}); !errors.Is(err, ErrNotFound) {
		t.Fatalf("se esperaba ErrNotFound, llegó %v", err)
	}
}

func TestUpdate_LaContrasenaSeGuardaHasheada(t *testing.T) {
	m := nuevoMock(admin(1), gestor(2, "tesis"))
	if err := servicio(m).Update(1, 2, &UpdateFieldGestor{PasswordHash: ptr("clave-nueva")}); err != nil {
		t.Fatalf("llegó %v", err)
	}
	guardado := *m.ultimoUpdate.PasswordHash
	if guardado == "clave-nueva" {
		t.Fatal("la contraseña no se puede guardar en texto plano")
	}
	if bcrypt.CompareHashAndPassword([]byte(guardado), []byte("clave-nueva")) != nil {
		t.Fatal("el hash guardado no corresponde a la contraseña")
	}
}

func TestUpdate_DuplicadosSePropagan(t *testing.T) {
	casos := []error{ErrUsernameDuplicado, ErrEmailDuplicado}
	for _, dup := range casos {
		m := nuevoMock(admin(1), gestor(2, "tesis"))
		m.errUpdate = dup
		err := servicio(m).Update(1, 2, &UpdateFieldGestor{Username: ptr("otro")})
		if !errors.Is(err, dup) {
			t.Errorf("se esperaba %v, llegó %v", dup, err)
		}
	}
}

func TestCreate_RolDerivadoDeLosModulos(t *testing.T) {
	casos := []struct {
		modulos string
		rol     string
	}{
		{"tesis", "GESTOR TESIS"},
		{"tesis,integrantes", "GESTOR"},
	}
	for _, c := range casos {
		m := nuevoMock()
		u := &UsuarioGestor{Username: ptr("nuevo"), PasswordHash: ptr("clave"), Email: ptr("n@x.com"), Modulos: ptr(c.modulos)}
		if err := servicio(m).Create(u); err != nil {
			t.Fatalf("%s: llegó %v", c.modulos, err)
		}
		if *m.ultimoCreate.Rol != c.rol {
			t.Errorf("%s: el rol debería ser %q, es %q", c.modulos, c.rol, *m.ultimoCreate.Rol)
		}
		if *m.ultimoCreate.PasswordHash == "clave" {
			t.Error("la contraseña no se puede guardar en texto plano")
		}
	}
}

func TestCreate_NuncaCreaUnAdministrador(t *testing.T) {
	m := nuevoMock()
	u := &UsuarioGestor{Username: ptr("nuevo"), PasswordHash: ptr("clave"), Email: ptr("n@x.com"), Modulos: ptr("ADMIN")}
	if err := servicio(m).Create(u); err != nil {
		t.Fatalf("llegó %v", err)
	}
	if *m.ultimoCreate.Rol == "ADMIN" {
		t.Fatal("ADMIN no se puede obtener asignando un módulo")
	}
}

func TestCreate_Validaciones(t *testing.T) {
	m := nuevoMock()
	s := servicio(m)
	base := func() *UsuarioGestor {
		return &UsuarioGestor{Username: ptr("nuevo"), PasswordHash: ptr("clave"), Email: ptr("n@x.com"), Modulos: ptr("tesis")}
	}
	u := base()
	u.Username = ptr("")
	if err := s.Create(u); !errors.Is(err, ErrUsernameRequerido) {
		t.Errorf("usuario vacío: llegó %v", err)
	}
	u = base()
	u.Username = ptr("con espacio")
	if err := s.Create(u); !errors.Is(err, ErrUsernameInvalido) {
		t.Errorf("usuario con espacio: llegó %v", err)
	}
	u = base()
	u.PasswordHash = ptr("")
	if err := s.Create(u); !errors.Is(err, ErrPasswordRequerido) {
		t.Errorf("contraseña vacía: llegó %v", err)
	}
	u = base()
	u.Email = ptr("")
	if err := s.Create(u); !errors.Is(err, ErrEmailRequerido) {
		t.Errorf("email vacío: llegó %v", err)
	}
	u = base()
	u.Modulos = ptr("")
	if err := s.Create(u); !errors.Is(err, ErrModulosRequerido) {
		t.Errorf("sin módulos: llegó %v", err)
	}
	if m.llamo("Create") != 0 {
		t.Error("ninguna validación fallida debería llegar al storage")
	}
}

func TestCreate_DuplicadosSePropagan(t *testing.T) {
	for _, dup := range []error{ErrUsernameDuplicado, ErrEmailDuplicado} {
		m := nuevoMock()
		m.errCreate = dup
		u := &UsuarioGestor{Username: ptr("nuevo"), PasswordHash: ptr("clave"), Email: ptr("n@x.com"), Modulos: ptr("tesis")}
		if err := servicio(m).Create(u); !errors.Is(err, dup) {
			t.Errorf("se esperaba %v, llegó %v", dup, err)
		}
	}
}

func TestRead_NoEncontradoNoEsUnErrorDelSistema(t *testing.T) {
	m := nuevoMock()
	if _, err := servicio(m).Read(5); !errors.Is(err, ErrNotFound) {
		t.Fatalf("se esperaba ErrNotFound, llegó %v", err)
	}
	if _, err := servicio(m).ReadByUsername("nadie"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("se esperaba ErrNotFound, llegó %v", err)
	}
}
