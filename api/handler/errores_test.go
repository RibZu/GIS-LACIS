package handler

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func routerDeErrores() *gin.Engine {
	r := gin.New()
	r.SetHTMLTemplate(template.Must(template.New("ErrorPagina.html").Parse(`{{ .Titulo }}|{{ .Mensaje }}`)))
	r.Use(gin.CustomRecovery(NuevoRecuperador(zap.NewNop())))
	r.NoRoute(PaginaNoEncontrada)
	r.GET("/boom", func(c *gin.Context) {
		panic("falla simulada")
	})
	return r
}

func TestRecuperador_UnPanicoDaUnaPaginaAmable(t *testing.T) {
	w := httptest.NewRecorder()
	routerDeErrores().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/boom", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("se esperaba 500, llegó %d", w.Code)
	}
	cuerpo := w.Body.String()
	if !strings.Contains(cuerpo, "Algo salió mal. Probá de nuevo en unos minutos.") {
		t.Errorf("falta el mensaje amable, el cuerpo es %q", cuerpo)
	}
	if strings.Contains(cuerpo, "falla simulada") {
		t.Error("el detalle técnico del pánico no puede llegar a la persona")
	}
}

func TestPaginaNoEncontrada_Da404ConMensaje(t *testing.T) {
	w := httptest.NewRecorder()
	routerDeErrores().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/no-existe", nil))

	if w.Code != http.StatusNotFound {
		t.Fatalf("se esperaba 404, llegó %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "No encontramos esa página.") {
		t.Errorf("falta el mensaje, el cuerpo es %q", w.Body.String())
	}
}

func contextoConQuery(query string) *gin.Context {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/?"+query, nil)
	return c
}

func TestAvisosDesdeQuery_SoloTextosDeLaListaBlanca(t *testing.T) {
	ok, aviso := avisosDesdeQuery(contextoConQuery("ok=producto-eliminado&error=ultimo-admin"))
	if ok != "Producto eliminado." {
		t.Errorf("aviso de éxito inesperado: %q", ok)
	}
	if aviso != "Tiene que quedar al menos un administrador." {
		t.Errorf("aviso de error inesperado: %q", aviso)
	}

	ok, aviso = avisosDesdeQuery(contextoConQuery("ok=%3Cscript%3Ealert(1)%3C%2Fscript%3E&error=texto-libre"))
	if ok != "" || aviso != "" {
		t.Errorf("un código desconocido debe ignorarse, llegó ok=%q error=%q", ok, aviso)
	}

	ok, aviso = avisosDesdeQuery(contextoConQuery(""))
	if ok != "" || aviso != "" {
		t.Errorf("sin parámetros no debe haber avisos, llegó ok=%q error=%q", ok, aviso)
	}
}

func TestAvisos_TodosLosCodigosTienenTexto(t *testing.T) {
	esperados := map[string]map[string]string{
		"ok": {
			"producto-creado":      "Producto guardado.",
			"producto-actualizado": "Cambios guardados.",
			"producto-eliminado":   "Producto eliminado.",
			"usuario-creado":       "Usuario creado.",
			"usuario-actualizado":  "Cambios guardados.",
			"usuario-eliminado":    "Usuario eliminado.",
		},
		"error": {
			"no-encontrado":   "No encontramos ese registro.",
			"error-eliminar":  "No se pudo eliminar. Probá de nuevo en unos minutos.",
			"autoeliminacion": "No podés eliminar tu propia cuenta.",
			"ultimo-admin":    "Tiene que quedar al menos un administrador.",
		},
	}
	for codigo, texto := range esperados["ok"] {
		if avisosOK[codigo] != texto {
			t.Errorf("aviso ok %q: se esperaba %q, es %q", codigo, texto, avisosOK[codigo])
		}
	}
	for codigo, texto := range esperados["error"] {
		if avisosError[codigo] != texto {
			t.Errorf("aviso de error %q: se esperaba %q, es %q", codigo, texto, avisosError[codigo])
		}
	}
	if len(avisosOK) != len(esperados["ok"]) || len(avisosError) != len(esperados["error"]) {
		t.Error("las listas blancas tienen códigos que la especificación no define")
	}
}
