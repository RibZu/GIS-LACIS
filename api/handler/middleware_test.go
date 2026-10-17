package handler

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"fmt"

	"PaginaSEG/internal/usuario"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ptrStr(s string) *string {
	return &s
}

type mockUsuarios struct {
	usuarios       map[int]*usuario.UsuarioGestor
	lecturas       int
	borrados       []int
	errLectura     error
	errPorUsername error
}

func nuevoMockUsuarios(usuarios ...*usuario.UsuarioGestor) *mockUsuarios {
	m := &mockUsuarios{usuarios: map[int]*usuario.UsuarioGestor{}}
	for _, u := range usuarios {
		m.usuarios[u.ID] = u
	}
	return m
}

func (m *mockUsuarios) Create(u *usuario.UsuarioGestor) error {
	u.ID = len(m.usuarios) + 100
	m.usuarios[u.ID] = u
	return nil
}

func (m *mockUsuarios) Read(id int) (*usuario.UsuarioGestor, error) {
	m.lecturas++
	if m.errLectura != nil {
		return nil, m.errLectura
	}
	u, ok := m.usuarios[id]
	if !ok {
		return nil, usuario.ErrNotFound
	}
	return u, nil
}

func (m *mockUsuarios) ReadByUsername(username string) (*usuario.UsuarioGestor, error) {
	if m.errPorUsername != nil {
		return nil, m.errPorUsername
	}
	for _, u := range m.usuarios {
		if u.Username != nil && *u.Username == username {
			return u, nil
		}
	}
	return nil, usuario.ErrNotFound
}

func (m *mockUsuarios) GetAll() ([]usuario.UsuarioGestor, error) {
	var todos []usuario.UsuarioGestor
	for _, u := range m.usuarios {
		todos = append(todos, *u)
	}
	return todos, nil
}

func (m *mockUsuarios) Update(id int, fields *usuario.UpdateFieldGestor) error {
	return nil
}

func (m *mockUsuarios) Delete(id int) error {
	m.borrados = append(m.borrados, id)
	delete(m.usuarios, id)
	return nil
}

func (m *mockUsuarios) ContarAdmins() (int, error) {
	total := 0
	for _, u := range m.usuarios {
		if u.Rol != nil && *u.Rol == "ADMIN" {
			total++
		}
	}
	return total, nil
}

func adminDePrueba(id int) *usuario.UsuarioGestor {
	return &usuario.UsuarioGestor{ID: id, Username: ptrStr("admin"), Rol: ptrStr("ADMIN")}
}

func gestorDePrueba(id int, modulos string) *usuario.UsuarioGestor {
	return &usuario.UsuarioGestor{ID: id, Username: ptrStr("gestor"), Rol: ptrStr("GESTOR"), Modulos: ptrStr(modulos)}
}

func routerDePrueba(m *mockUsuarios) *gin.Engine {
	svc := usuario.NewService(m, zap.NewNop())
	uh := NewUsuarioHandler(svc, zap.NewNop())
	ok := func(c *gin.Context) { c.String(http.StatusOK, "ok") }

	r := gin.New()
	r.SetHTMLTemplate(template.Must(template.New("ErrorPagina.html").Parse(`{{ .Titulo }}|{{ .Mensaje }}`)))
	admin := r.Group("/admin", RequireLogin(svc))
	admin.GET("/dashboard", ok)

	soloAdmin := admin.Group("")
	soloAdmin.Use(RequireAdmin(svc))
	soloAdmin.GET("/usuarios", ok)
	soloAdmin.POST("/borrar-usuario", uh.Borrar)
	soloAdmin.GET("/borrar-usuario", uh.BorrarGet)

	conTesis := admin.Group("")
	conTesis.Use(RequireModule(svc, "tesis"))
	conTesis.GET("/tesis", ok)

	api := r.Group("/api")
	api.Use(RequireModuleAPI(svc, "tesis"))
	api.GET("/tesis", ok)
	return r
}

func cookieDeUsuario(id int) *http.Cookie {
	c, w := contextoConCookie("", nil)
	IniciarSesion(c, id)
	return cookieDeSesion(w)
}

func pedir(r *gin.Engine, metodo, ruta string, cookie *http.Cookie, cuerpo string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(metodo, ruta, strings.NewReader(cuerpo))
	if cuerpo != "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func esperarRedireccion(t *testing.T, w *httptest.ResponseRecorder, codigo int, destino string) {
	t.Helper()
	if w.Code != codigo {
		t.Fatalf("se esperaba status %d, llegó %d", codigo, w.Code)
	}
	if loc := w.Header().Get("Location"); loc != destino {
		t.Fatalf("se esperaba redirigir a %q, redirigió a %q", destino, loc)
	}
}

func TestAdmin_SinCookieVaAlLogin(t *testing.T) {
	usarClave(t, "clave-de-test")
	r := routerDePrueba(nuevoMockUsuarios(adminDePrueba(1)))
	esperarRedireccion(t, pedir(r, http.MethodGet, "/admin/dashboard", nil, ""), http.StatusFound, "/login")
}

func TestAdmin_CookiesFalsificadasVanAlLogin(t *testing.T) {
	usarClave(t, "clave-de-test")
	r := routerDePrueba(nuevoMockUsuarios(adminDePrueba(1)))
	valida := cookieDeUsuario(1)
	partes := strings.Split(valida.Value, ".")
	falsas := []string{
		"1",
		"azar",
		valida.Value + "x",
		fmt.Sprintf("2.%s.%s", partes[1], partes[2]),
		fmt.Sprintf("1.%d.%s", time.Now().Unix()+99999, partes[2]),
	}
	for _, valor := range falsas {
		cookie := &http.Cookie{Name: sessionCookieName, Value: valor}
		esperarRedireccion(t, pedir(r, http.MethodGet, "/admin/dashboard", cookie, ""), http.StatusFound, "/login")
	}
}

func TestAdmin_SesionVencidaAvisaQueExpiro(t *testing.T) {
	usarClave(t, "clave-de-test")
	r := routerDePrueba(nuevoMockUsuarios(adminDePrueba(1)))
	vence := time.Now().Unix() - 30
	cookie := &http.Cookie{Name: sessionCookieName, Value: fmt.Sprintf("1.%d.%s", vence, firmar(1, vence))}
	esperarRedireccion(t, pedir(r, http.MethodGet, "/admin/dashboard", cookie, ""), http.StatusFound, "/login?sesion=expirada")
	esperarRedireccion(t, pedir(r, http.MethodGet, "/admin/usuarios", cookie, ""), http.StatusFound, "/login?sesion=expirada")
	esperarRedireccion(t, pedir(r, http.MethodGet, "/admin/tesis", cookie, ""), http.StatusFound, "/login?sesion=expirada")
}

func TestAdmin_UsuarioBorradoMientrasEstabaConectadoVaAlLogin(t *testing.T) {
	usarClave(t, "clave-de-test")
	r := routerDePrueba(nuevoMockUsuarios(adminDePrueba(1)))
	esperarRedireccion(t, pedir(r, http.MethodGet, "/admin/dashboard", cookieDeUsuario(77), ""), http.StatusFound, "/login")
}

func TestAdmin_ConSesionValidaSeVeLaPaginaYNoSeGuardaEnCache(t *testing.T) {
	usarClave(t, "clave-de-test")
	r := routerDePrueba(nuevoMockUsuarios(adminDePrueba(1)))
	w := pedir(r, http.MethodGet, "/admin/dashboard", cookieDeUsuario(1), "")
	if w.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, llegó %d", w.Code)
	}
	if cc := w.Header().Get("Cache-Control"); cc != "no-store" {
		t.Fatalf("Cache-Control debería ser no-store, es %q", cc)
	}
}

func TestAdmin_UnaSolaLecturaDelUsuarioPorPeticion(t *testing.T) {
	usarClave(t, "clave-de-test")
	m := nuevoMockUsuarios(adminDePrueba(1))
	r := routerDePrueba(m)
	w := pedir(r, http.MethodGet, "/admin/usuarios", cookieDeUsuario(1), "")
	if w.Code != http.StatusOK {
		t.Fatalf("se esperaba 200, llegó %d", w.Code)
	}
	if m.lecturas != 1 {
		t.Fatalf("RequireLogin + RequireAdmin deberían leer el usuario una sola vez, leyeron %d", m.lecturas)
	}
}

func TestRequireAdmin_UnGestorNoEntra(t *testing.T) {
	usarClave(t, "clave-de-test")
	r := routerDePrueba(nuevoMockUsuarios(adminDePrueba(1), gestorDePrueba(2, "tesis")))
	esperarRedireccion(t, pedir(r, http.MethodGet, "/admin/usuarios", cookieDeUsuario(2), ""), http.StatusFound, "/admin/dashboard")
	esperarRedireccion(t, pedir(r, http.MethodPost, "/admin/borrar-usuario", cookieDeUsuario(2), "id=1"), http.StatusFound, "/admin/dashboard")
}

func TestRequireModule(t *testing.T) {
	usarClave(t, "clave-de-test")
	m := nuevoMockUsuarios(adminDePrueba(1), gestorDePrueba(2, "tesis,estadisticas"), gestorDePrueba(3, "proyectos"))
	r := routerDePrueba(m)

	if w := pedir(r, http.MethodGet, "/admin/tesis", cookieDeUsuario(2), ""); w.Code != http.StatusOK {
		t.Errorf("un gestor con el módulo debería entrar, llegó %d", w.Code)
	}
	esperarRedireccion(t, pedir(r, http.MethodGet, "/admin/tesis", cookieDeUsuario(3), ""), http.StatusFound, "/admin/dashboard")
	if w := pedir(r, http.MethodGet, "/admin/tesis", cookieDeUsuario(1), ""); w.Code != http.StatusOK {
		t.Errorf("un administrador entra a cualquier módulo, llegó %d", w.Code)
	}
}

func TestRequireModuleAPI(t *testing.T) {
	usarClave(t, "clave-de-test")
	m := nuevoMockUsuarios(adminDePrueba(1), gestorDePrueba(2, "tesis"), gestorDePrueba(3, "proyectos"))
	r := routerDePrueba(m)

	if w := pedir(r, http.MethodGet, "/api/tesis", nil, ""); w.Code != http.StatusUnauthorized {
		t.Errorf("sin sesión se esperaba 401, llegó %d", w.Code)
	}
	forjada := &http.Cookie{Name: sessionCookieName, Value: "1"}
	if w := pedir(r, http.MethodGet, "/api/tesis", forjada, ""); w.Code != http.StatusUnauthorized {
		t.Errorf("con la cookie vieja sin firma se esperaba 401, llegó %d", w.Code)
	}
	if w := pedir(r, http.MethodGet, "/api/tesis", cookieDeUsuario(3), ""); w.Code != http.StatusForbidden {
		t.Errorf("sin el módulo se esperaba 403, llegó %d", w.Code)
	}
	if w := pedir(r, http.MethodGet, "/api/tesis", cookieDeUsuario(2), ""); w.Code != http.StatusOK {
		t.Errorf("con el módulo se esperaba 200, llegó %d", w.Code)
	}
}

func TestBorrarUsuario_PorGETNoBorraNada(t *testing.T) {
	usarClave(t, "clave-de-test")
	m := nuevoMockUsuarios(adminDePrueba(1), gestorDePrueba(2, "tesis"))
	r := routerDePrueba(m)
	esperarRedireccion(t, pedir(r, http.MethodGet, "/admin/borrar-usuario?id=2", cookieDeUsuario(1), ""), http.StatusSeeOther, "/admin/usuarios")
	if len(m.borrados) != 0 {
		t.Fatalf("un GET no puede borrar, se borró %v", m.borrados)
	}
}

func TestBorrarUsuario_PorPOST(t *testing.T) {
	usarClave(t, "clave-de-test")
	m := nuevoMockUsuarios(adminDePrueba(1), gestorDePrueba(2, "tesis"))
	r := routerDePrueba(m)
	esperarRedireccion(t, pedir(r, http.MethodPost, "/admin/borrar-usuario", cookieDeUsuario(1), "id=2"), http.StatusSeeOther, "/admin/usuarios?ok=usuario-eliminado")
	if len(m.borrados) != 1 || m.borrados[0] != 2 {
		t.Fatalf("se esperaba borrar el usuario 2, se borró %v", m.borrados)
	}
}

func TestBorrarUsuario_NoPuedeBorrarseASiMismo(t *testing.T) {
	usarClave(t, "clave-de-test")
	m := nuevoMockUsuarios(adminDePrueba(1), adminDePrueba(2))
	r := routerDePrueba(m)
	esperarRedireccion(t, pedir(r, http.MethodPost, "/admin/borrar-usuario", cookieDeUsuario(1), "id=1"), http.StatusSeeOther, "/admin/usuarios?error=autoeliminacion")
	if len(m.borrados) != 0 {
		t.Fatalf("no debería borrarse nada, se borró %v", m.borrados)
	}
}

func TestBorrarUsuario_UnAdminPuedeBorrarAOtroAdminSiHayDos(t *testing.T) {
	usarClave(t, "clave-de-test")
	m := nuevoMockUsuarios(adminDePrueba(1), adminDePrueba(2))
	r := routerDePrueba(m)
	esperarRedireccion(t, pedir(r, http.MethodPost, "/admin/borrar-usuario", cookieDeUsuario(1), "id=2"), http.StatusSeeOther, "/admin/usuarios?ok=usuario-eliminado")
}

func TestBorrarUsuario_IdsInvalidos(t *testing.T) {
	usarClave(t, "clave-de-test")
	m := nuevoMockUsuarios(adminDePrueba(1))
	r := routerDePrueba(m)
	for _, cuerpo := range []string{"", "id=abc", "id=0", "id=-4", "id=999"} {
		esperarRedireccion(t, pedir(r, http.MethodPost, "/admin/borrar-usuario", cookieDeUsuario(1), cuerpo), http.StatusSeeOther, "/admin/usuarios?error=no-encontrado")
	}
	if len(m.borrados) != 0 {
		t.Fatalf("no debería borrarse nada, se borró %v", m.borrados)
	}
}

func TestBorrarUsuario_SinSesionNoBorra(t *testing.T) {
	usarClave(t, "clave-de-test")
	m := nuevoMockUsuarios(adminDePrueba(1), gestorDePrueba(2, "tesis"))
	r := routerDePrueba(m)
	esperarRedireccion(t, pedir(r, http.MethodPost, "/admin/borrar-usuario", nil, "id=2"), http.StatusFound, "/login")
	if len(m.borrados) != 0 {
		t.Fatalf("sin sesión no puede borrarse nada, se borró %v", m.borrados)
	}
}

func TestAdmin_BaseCaidaMuestraAvisoYNoMandaAlLogin(t *testing.T) {
	usarClave(t, "clave-de-test")
	m := nuevoMockUsuarios(adminDePrueba(1))
	m.errLectura = errors.New("conexión rechazada")
	r := routerDePrueba(m)

	for _, ruta := range []string{"/admin/dashboard", "/admin/usuarios", "/admin/tesis"} {
		w := pedir(r, http.MethodGet, ruta, cookieDeUsuario(1), "")
		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("%s: se esperaba 503, llegó %d", ruta, w.Code)
		}
		if loc := w.Header().Get("Location"); loc != "" {
			t.Errorf("%s: no debería redirigir al login, redirigió a %q", ruta, loc)
		}
		if !strings.Contains(w.Body.String(), mensajeServicioNoDisponible) {
			t.Errorf("%s: falta el aviso, el cuerpo es %q", ruta, w.Body.String())
		}
		if strings.Contains(w.Body.String(), "conexión rechazada") {
			t.Errorf("%s: el detalle técnico no puede llegar a la persona", ruta)
		}
	}
}

func TestRequireModuleAPI_BaseCaidaDa503(t *testing.T) {
	usarClave(t, "clave-de-test")
	m := nuevoMockUsuarios(gestorDePrueba(2, "tesis"))
	m.errLectura = errors.New("conexión rechazada")
	w := pedir(routerDePrueba(m), http.MethodGet, "/api/tesis", cookieDeUsuario(2), "")
	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("se esperaba 503, llegó %d", w.Code)
	}
}
