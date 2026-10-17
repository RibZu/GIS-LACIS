package handler

import (
	"errors"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"PaginaSEG/internal/usuario"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"golang.org/x/crypto/bcrypt"
)

const claveDePrueba = "super-secreta-123"

func usuarioConClave(t *testing.T, id int, nombre, clave string) *usuario.UsuarioGestor {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(clave), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	h := string(hash)
	return &usuario.UsuarioGestor{ID: id, Username: &nombre, PasswordHash: &h, Rol: ptrStr("GESTOR"), Modulos: ptrStr("tesis")}
}

func routerDeLogin(m *mockUsuarios, logger *zap.Logger) *gin.Engine {
	svc := usuario.NewService(m, zap.NewNop())
	ah := NewAuthHandler(svc, nil, logger)
	r := gin.New()
	r.SetHTMLTemplate(template.Must(template.New("Login.html").Parse(`[{{ .Error }}][{{ .Aviso }}]`)))
	r.GET("/login", ah.ShowLogin)
	r.POST("/login", ah.ProcessLogin)
	r.GET("/logout", ah.Logout)
	return r
}

func enviarLogin(r *gin.Engine, usuarioForm, claveForm string) *httptest.ResponseRecorder {
	cuerpo := "username=" + usuarioForm + "&password=" + claveForm
	return pedir(r, http.MethodPost, "/login", nil, cuerpo)
}

func TestLogin_UsuarioInexistenteYClaveIncorrectaSonIndistinguibles(t *testing.T) {
	usarClave(t, "clave-de-test")
	m := nuevoMockUsuarios(usuarioConClave(t, 5, "ana", claveDePrueba))
	r := routerDeLogin(m, zap.NewNop())

	inexistente := enviarLogin(r, "nadie", "cualquiera")
	incorrecta := enviarLogin(r, "ana", "mala-clave")
	vacio := enviarLogin(r, "", "")

	for nombre, w := range map[string]*httptest.ResponseRecorder{"inexistente": inexistente, "incorrecta": incorrecta, "vacío": vacio} {
		if w.Code != http.StatusUnauthorized {
			t.Errorf("%s: se esperaba 401, llegó %d", nombre, w.Code)
		}
		if !strings.Contains(w.Body.String(), mensajeLoginFallido) {
			t.Errorf("%s: falta el mensaje uniforme, el cuerpo es %q", nombre, w.Body.String())
		}
		if w.Header().Get("Set-Cookie") != "" {
			t.Errorf("%s: no debería crearse ninguna sesión", nombre)
		}
	}
	if inexistente.Body.String() != incorrecta.Body.String() {
		t.Errorf("las respuestas deberían ser idénticas:\n%q\n%q", inexistente.Body.String(), incorrecta.Body.String())
	}
	if inexistente.Body.String() != vacio.Body.String() {
		t.Errorf("las respuestas deberían ser idénticas:\n%q\n%q", inexistente.Body.String(), vacio.Body.String())
	}
}

func TestLogin_UsuarioSinHashDaLaMismaRespuesta(t *testing.T) {
	usarClave(t, "clave-de-test")
	sinHash := &usuario.UsuarioGestor{ID: 6, Username: ptrStr("fantasma"), Rol: ptrStr("GESTOR")}
	m := nuevoMockUsuarios(usuarioConClave(t, 5, "ana", claveDePrueba), sinHash)
	r := routerDeLogin(m, zap.NewNop())

	a := enviarLogin(r, "fantasma", "x")
	b := enviarLogin(r, "nadie", "x")
	if a.Code != http.StatusUnauthorized || a.Body.String() != b.Body.String() {
		t.Fatalf("un usuario sin hash debería responder igual que uno inexistente: %d %q vs %q", a.Code, a.Body.String(), b.Body.String())
	}
}

func TestLogin_CorrectoCreaUnaSesionFirmada(t *testing.T) {
	usarClave(t, "clave-de-test")
	m := nuevoMockUsuarios(usuarioConClave(t, 5, "ana", claveDePrueba))
	w := enviarLogin(routerDeLogin(m, zap.NewNop()), "ana", claveDePrueba)
	esperarRedireccion(t, w, http.StatusFound, "/admin/dashboard")

	ck := cookieDeSesion(w)
	if ck == nil {
		t.Fatal("el login correcto debería crear la cookie de sesión")
	}
	if ck.Value == "5" {
		t.Fatal("la cookie no puede ser el ID sin firma")
	}
	c, _ := contextoConCookie(ck.Value, nil)
	if id, ok := CurrentUserID(c); !ok || id != 5 {
		t.Fatalf("la cookie debería ser una sesión válida del usuario 5, dio id=%d ok=%v", id, ok)
	}
}

func TestLogin_BaseCaidaNoSeHaceCompletarComoClaveIncorrecta(t *testing.T) {
	usarClave(t, "clave-de-test")
	m := nuevoMockUsuarios()
	m.errPorUsername = errors.New("conexión rechazada")
	w := enviarLogin(routerDeLogin(m, zap.NewNop()), "ana", claveDePrueba)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("se esperaba 503, llegó %d", w.Code)
	}
	cuerpo := w.Body.String()
	if !strings.Contains(cuerpo, mensajeLoginNoDisponible) {
		t.Errorf("falta el aviso de servicio no disponible, el cuerpo es %q", cuerpo)
	}
	if strings.Contains(cuerpo, mensajeLoginFallido) || strings.Contains(cuerpo, "conexión rechazada") {
		t.Errorf("no debería decir que la clave es incorrecta ni mostrar el detalle técnico: %q", cuerpo)
	}
	if w.Header().Get("Set-Cookie") != "" {
		t.Error("no debería crearse ninguna sesión")
	}
}

func TestLogin_LosFallosSeRegistranSinLaContrasena(t *testing.T) {
	usarClave(t, "clave-de-test")
	nucleo, registros := observer.New(zap.InfoLevel)
	m := nuevoMockUsuarios(usuarioConClave(t, 5, "ana", claveDePrueba))
	r := routerDeLogin(m, zap.New(nucleo))

	enviarLogin(r, "ana", "intento-con-esta-clave")
	enviarLogin(r, "nadie", "otra-clave-distinta")
	enviarLogin(r, "ana", claveDePrueba)

	var fallidos, correctos int
	for _, entrada := range registros.All() {
		volcado := entrada.Message
		for _, campo := range entrada.Context {
			volcado += " " + campo.String
		}
		for _, secreto := range []string{"intento-con-esta-clave", "otra-clave-distinta", claveDePrueba} {
			if strings.Contains(volcado, secreto) {
				t.Errorf("la contraseña %q no puede aparecer en el log: %q", secreto, volcado)
			}
		}
		switch entrada.Message {
		case "Intento de login fallido":
			fallidos++
			if entrada.Level != zap.WarnLevel {
				t.Errorf("los fallos deberían registrarse como Warn, fue %v", entrada.Level)
			}
		case "Inicio de sesión":
			correctos++
			if entrada.Level != zap.InfoLevel {
				t.Errorf("un login correcto debería registrarse como Info, fue %v", entrada.Level)
			}
		}
	}
	if fallidos != 2 || correctos != 1 {
		t.Errorf("se esperaban 2 fallos y 1 login correcto registrados, hubo %d y %d", fallidos, correctos)
	}
}

func TestShowLogin(t *testing.T) {
	usarClave(t, "clave-de-test")
	m := nuevoMockUsuarios(usuarioConClave(t, 5, "ana", claveDePrueba))
	r := routerDeLogin(m, zap.NewNop())

	w := pedir(r, http.MethodGet, "/login", nil, "")
	if w.Code != http.StatusOK || strings.Contains(w.Body.String(), "expiró") {
		t.Errorf("el login normal no debería avisar vencimiento: %d %q", w.Code, w.Body.String())
	}

	w = pedir(r, http.MethodGet, "/login?sesion=expirada", nil, "")
	if !strings.Contains(w.Body.String(), "Tu sesión expiró. Volvé a ingresar.") {
		t.Errorf("falta el aviso de sesión vencida: %q", w.Body.String())
	}

	w = pedir(r, http.MethodGet, "/login?sesion=<script>", nil, "")
	if strings.Contains(w.Body.String(), "<script>") {
		t.Errorf("un valor desconocido no debe reflejarse en la página: %q", w.Body.String())
	}

	esperarRedireccion(t, pedir(r, http.MethodGet, "/login", cookieDeUsuario(5), ""), http.StatusFound, "/admin/dashboard")

	forjada := &http.Cookie{Name: sessionCookieName, Value: "5"}
	if w := pedir(r, http.MethodGet, "/login", forjada, ""); w.Code != http.StatusOK {
		t.Errorf("con una cookie falsificada se debería mostrar el login, llegó %d", w.Code)
	}
}

func TestLogout_BorraLaSesion(t *testing.T) {
	usarClave(t, "clave-de-test")
	m := nuevoMockUsuarios(usuarioConClave(t, 5, "ana", claveDePrueba))
	w := pedir(routerDeLogin(m, zap.NewNop()), http.MethodGet, "/logout", cookieDeUsuario(5), "")
	esperarRedireccion(t, w, http.StatusFound, "/login")
	ck := cookieDeSesion(w)
	if ck == nil || ck.MaxAge >= 0 || ck.Value != "" {
		t.Fatalf("el logout debería borrar la cookie, llegó %+v", ck)
	}
}
