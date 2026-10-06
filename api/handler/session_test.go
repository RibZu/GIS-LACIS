package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func usarClave(t *testing.T, clave string) {
	t.Helper()
	anterior := claveSesion
	claveSesion = []byte(clave)
	t.Cleanup(func() { claveSesion = anterior })
}

func contextoConCookie(valor string, cabeceras map[string]string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if valor != "" {
		req.AddCookie(&http.Cookie{Name: sessionCookieName, Value: valor})
	}
	for k, v := range cabeceras {
		req.Header.Set(k, v)
	}
	c.Request = req
	return c, w
}

func cookieDeSesion(w *httptest.ResponseRecorder) *http.Cookie {
	for _, ck := range w.Result().Cookies() {
		if ck.Name == sessionCookieName {
			return ck
		}
	}
	return nil
}

func valorDeSesionIniciada(t *testing.T, id int) string {
	t.Helper()
	c, w := contextoConCookie("", nil)
	IniciarSesion(c, id)
	ck := cookieDeSesion(w)
	if ck == nil {
		t.Fatal("IniciarSesion no escribió la cookie de sesión")
	}
	return ck.Value
}

func TestSesion_CreadaSeLeeComoValida(t *testing.T) {
	usarClave(t, "clave-de-test")
	c, _ := contextoConCookie(valorDeSesionIniciada(t, 7), nil)
	id, estado := leerSesion(c)
	if estado != sesionValida || id != 7 {
		t.Fatalf("se esperaba sesión válida del usuario 7, llegó id=%d estado=%d", id, estado)
	}
	if id, ok := CurrentUserID(c); !ok || id != 7 {
		t.Fatalf("CurrentUserID debería devolver 7 y true, devolvió %d y %v", id, ok)
	}
}

func TestSesion_SinCookieEsNinguna(t *testing.T) {
	usarClave(t, "clave-de-test")
	c, _ := contextoConCookie("", nil)
	if _, estado := leerSesion(c); estado != sesionNinguna {
		t.Fatalf("sin cookie se esperaba sesionNinguna, llegó %d", estado)
	}
	if _, ok := CurrentUserID(c); ok {
		t.Fatal("sin cookie CurrentUserID no puede dar true")
	}
}

func TestSesion_ValoresFalsificadosSonInvalidos(t *testing.T) {
	usarClave(t, "clave-de-test")
	valida := valorDeSesionIniciada(t, 7)
	partes := strings.Split(valida, ".")

	casos := map[string]string{
		"formato viejo, solo el id": "1",
		"valor al azar":             "cualquier-cosa",
		"id cambiado con la firma":  fmt.Sprintf("1.%s.%s", partes[1], partes[2]),
		"vencimiento cambiado":      fmt.Sprintf("%s.%d.%s", partes[0], time.Now().Unix()+999999, partes[2]),
		"firma vacía":               fmt.Sprintf("%s.%s.", partes[0], partes[1]),
		"firma alterada":            valida + "x",
		"id no numérico":            fmt.Sprintf("abc.%s.%s", partes[1], partes[2]),
		"id cero":                   fmt.Sprintf("0.%s.%s", partes[1], partes[2]),
		"vencimiento no numérico":   fmt.Sprintf("%s.xyz.%s", partes[0], partes[2]),
		"partes de más":             valida + ".extra",
	}
	for nombre, valor := range casos {
		c, _ := contextoConCookie(valor, nil)
		id, estado := leerSesion(c)
		if estado != sesionInvalida || id != 0 {
			t.Errorf("%s: se esperaba sesionInvalida con id 0, llegó id=%d estado=%d", nombre, id, estado)
		}
		if _, ok := CurrentUserID(c); ok {
			t.Errorf("%s: CurrentUserID no puede dar true", nombre)
		}
	}
}

func TestSesion_VencidaConFirmaCorrecta(t *testing.T) {
	usarClave(t, "clave-de-test")
	vence := time.Now().Unix() - 10
	valor := fmt.Sprintf("5.%d.%s", vence, firmar(5, vence))
	c, _ := contextoConCookie(valor, nil)
	id, estado := leerSesion(c)
	if estado != sesionVencida || id != 0 {
		t.Fatalf("se esperaba sesionVencida con id 0, llegó id=%d estado=%d", id, estado)
	}
	if _, ok := CurrentUserID(c); ok {
		t.Fatal("una sesión vencida no puede dar CurrentUserID true")
	}
}

func TestSesion_FirmadaConOtraClaveEsInvalida(t *testing.T) {
	usarClave(t, "clave-uno")
	valor := valorDeSesionIniciada(t, 7)
	claveSesion = []byte("clave-dos")
	c, _ := contextoConCookie(valor, nil)
	if _, estado := leerSesion(c); estado != sesionInvalida {
		t.Fatalf("una cookie firmada con otra clave debería ser inválida, llegó estado=%d", estado)
	}
}

func TestConfigurarSesion_UsaElSecretoDado(t *testing.T) {
	anterior := claveSesion
	t.Cleanup(func() { claveSesion = anterior })

	ConfigurarSesion("secreto-de-railway", zap.NewNop())
	if string(claveSesion) != "secreto-de-railway" {
		t.Fatalf("la clave no se tomó del secreto, es %q", claveSesion)
	}

	antes := string(claveSesion)
	ConfigurarSesion("", zap.NewNop())
	if string(claveSesion) != antes {
		t.Fatal("sin secreto se debe conservar la clave que ya había")
	}
}

func TestClaveAleatoria_NoEsVaciaNiRepetida(t *testing.T) {
	a, b := claveAleatoria(), claveAleatoria()
	if len(a) != 32 || len(b) != 32 {
		t.Fatalf("la clave debería tener 32 bytes, tiene %d y %d", len(a), len(b))
	}
	if string(a) == string(b) {
		t.Fatal("dos claves aleatorias no deberían ser iguales")
	}
}

func TestIniciarSesion_AtributosDeLaCookie(t *testing.T) {
	usarClave(t, "clave-de-test")

	c, w := contextoConCookie("", nil)
	IniciarSesion(c, 3)
	ck := cookieDeSesion(w)
	if ck == nil {
		t.Fatal("no se escribió la cookie")
	}
	if !ck.HttpOnly {
		t.Error("la cookie debería ser HttpOnly")
	}
	if ck.SameSite != http.SameSiteLaxMode {
		t.Errorf("SameSite debería ser Lax, es %v", ck.SameSite)
	}
	if ck.Secure {
		t.Error("por HTTP la cookie no debería ser Secure")
	}
	if ck.MaxAge != duracionSesion {
		t.Errorf("MaxAge debería ser %d, es %d", duracionSesion, ck.MaxAge)
	}
	if ck.Path != "/" {
		t.Errorf("Path debería ser /, es %q", ck.Path)
	}

	c2, w2 := contextoConCookie("", map[string]string{"X-Forwarded-Proto": "https"})
	IniciarSesion(c2, 3)
	if ck2 := cookieDeSesion(w2); ck2 == nil || !ck2.Secure {
		t.Error("detrás de un proxy HTTPS la cookie debería ser Secure")
	}

	c3, w3 := contextoConCookie("", map[string]string{"X-Forwarded-Proto": "https, http"})
	IniciarSesion(c3, 3)
	if ck3 := cookieDeSesion(w3); ck3 == nil || !ck3.Secure {
		t.Error("con X-Forwarded-Proto en lista, el primer valor https debería marcar Secure")
	}

	c4, w4 := contextoConCookie("", map[string]string{"X-Forwarded-Proto": "http"})
	IniciarSesion(c4, 3)
	if ck4 := cookieDeSesion(w4); ck4 == nil || ck4.Secure {
		t.Error("con X-Forwarded-Proto http la cookie no debería ser Secure")
	}
}

func TestCerrarSesion_BorraLaCookie(t *testing.T) {
	c, w := contextoConCookie("", nil)
	CerrarSesion(c)
	ck := cookieDeSesion(w)
	if ck == nil {
		t.Fatal("CerrarSesion debería escribir una cookie para borrar la sesión")
	}
	if ck.MaxAge >= 0 {
		t.Errorf("MaxAge debería ser negativo, es %d", ck.MaxAge)
	}
	if ck.Value != "" {
		t.Errorf("el valor debería quedar vacío, es %q", ck.Value)
	}
	if !ck.HttpOnly {
		t.Error("la cookie de cierre debería mantener HttpOnly")
	}
}
