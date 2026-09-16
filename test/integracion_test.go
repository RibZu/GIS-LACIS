package test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"PaginaSEG/api"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	_ = os.Chdir("..")
	os.Exit(m.Run())
}

func TestIntegration_InsertarIntegrante(t *testing.T) {
	dsn := "postgres://postgres:isma_mesa22@localhost:5433/lacis?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err)
	defer db.Close()

	r := gin.Default()
	api.InitRoutes(r)

	formData := url.Values{}
	formData.Set("nombre", "Test")
	formData.Set("apellido", "Integracion")
	formData.Set("contacto", "test@unsl.edu.ar")
	formData.Set("especializacion", "Testing")
	formData.Set("descripcion", "Prueba E2E")
	formData.Set("pertenece_lacis", "true")
	formData.Set("rol_lacis_id", "1")

	req, _ := http.NewRequest("POST", "/admin/insertar-integrante", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session", Value: "1"})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "/admin/integrantes", w.Header().Get("Location"))
}

func TestIntegration_ObtenerIntegrantesJSON(t *testing.T) {
	dsn := "postgres://postgres:isma_mesa22@localhost:5433/lacis?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err)
	defer db.Close()

	r := gin.Default()
	api.InitRoutes(r)

	req, _ := http.NewRequest("GET", "/api/v1/integrantes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_InsertarTesis(t *testing.T) {
	dsn := "postgres://postgres:isma_mesa22@localhost:5433/lacis?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err)
	defer db.Close()

	r := gin.Default()
	api.InitRoutes(r)

	formData := url.Values{}
	formData.Set("titulo", "Tesis Test Integración")
	formData.Set("nivel", "Doctorado")
	formData.Set("carrera_origen", "Doctorado en Ingeniería en Informática")
	formData.Set("anio", "2025")
	formData.Set("autor_tipo", "externo")
	formData.Set("autor_historico", "Lic. Autor de Prueba")
	formData.Set("director_tipo", "externo")
	formData.Set("director_historico", "Dr. Director de Prueba")
	formData.Set("palabras_clave", "Testing, Calidad, Go")
	formData.Set("resumen", "Resumen de prueba de integración")

	req, _ := http.NewRequest("POST", "/admin/insertar-tesis", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "session", Value: "1"})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusSeeOther, w.Code)
	assert.Equal(t, "/admin/tesis", w.Header().Get("Location"))
}

func TestIntegration_LoginRedirectsIfAuthenticated(t *testing.T) {
	dsn := "postgres://postgres:isma_mesa22@localhost:5433/lacis?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err)
	defer db.Close()

	r := gin.Default()
	api.InitRoutes(r)

	req, _ := http.NewRequest("GET", "/login", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "1"})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/admin/dashboard", w.Header().Get("Location"))
}

func TestIntegration_LoginShowsFormWithoutSession(t *testing.T) {
	dsn := "postgres://postgres:isma_mesa22@localhost:5433/lacis?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err)
	defer db.Close()

	r := gin.Default()
	api.InitRoutes(r)

	req, _ := http.NewRequest("GET", "/login", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_ObtenerTesisJSON(t *testing.T) {
	dsn := "postgres://postgres:isma_mesa22@localhost:5433/lacis?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err)
	defer db.Close()

	r := gin.Default()
	api.InitRoutes(r)

	req, _ := http.NewRequest("GET", "/api/v1/tesis", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestIntegration_PosgradoUsaEncabezadoYPieCompartidos(t *testing.T) {
	dsn := "postgres://postgres:isma_mesa22@localhost:5433/lacis?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err)
	defer db.Close()

	r := gin.Default()
	api.InitRoutes(r)

	rutasPosgrado := []string{
		"/doctorado-ing-software",
		"/especializacion-ing-software",
		"/maestria-calidad-software",
		"/maestria-ing-software",
	}

	for _, ruta := range rutasPosgrado {
		req, _ := http.NewRequest("GET", ruta, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "ruta: %s", ruta)

		body := w.Body.String()
		assert.Contains(t, body, `href="/tesis"`, "ruta %s: falta la pestaña Tesis del encabezado compartido", ruta)
		assert.Contains(t, body, `href="/login"`, "ruta %s: falta el enlace de Iniciar Sesión del pie compartido", ruta)
		assert.Contains(t, body, "Iniciar Sesión", "ruta %s: falta el texto del botón Iniciar Sesión", ruta)
	}
}

func TestIntegration_AdminSinSesionRedirigeALogin(t *testing.T) {
	dsn := "postgres://postgres:isma_mesa22@localhost:5433/lacis?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err)
	defer db.Close()

	r := gin.Default()
	api.InitRoutes(r)

	rutasGET := []string{
		"/admin/dashboard",
		"/admin/integrantes",
		"/admin/desarrollos",
		"/admin/usuarios",
	}
	for _, ruta := range rutasGET {
		req, _ := http.NewRequest("GET", ruta, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusFound, w.Code, "ruta: %s", ruta)
		assert.Equal(t, "/login", w.Header().Get("Location"), "ruta: %s", ruta)
	}

	var countAntes int
	err = db.QueryRow("SELECT COUNT(*) FROM integrante").Scan(&countAntes)
	assert.NoError(t, err)

	formData := url.Values{}
	formData.Set("nombre", "NoDeberia")
	formData.Set("apellido", "Insertarse")
	formData.Set("contacto", "sin-sesion@unsl.edu.ar")
	formData.Set("especializacion", "Testing")
	formData.Set("descripcion", "No debería insertarse sin sesión")
	formData.Set("pertenece_lacis", "true")
	formData.Set("rol_lacis_id", "1")

	req, _ := http.NewRequest("POST", "/admin/insertar-integrante", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))

	var countDespues int
	err = db.QueryRow("SELECT COUNT(*) FROM integrante").Scan(&countDespues)
	assert.NoError(t, err)
	assert.Equal(t, countAntes, countDespues, "no debe haberse insertado ninguna fila sin sesión")

	req2, _ := http.NewRequest("GET", "/admin/dashboard", nil)
	req2.AddCookie(&http.Cookie{Name: "session", Value: "999999"})
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusFound, w2.Code)
	assert.Equal(t, "/login", w2.Header().Get("Location"))
}

func TestIntegration_AdminAPISinSesionDevuelve401(t *testing.T) {
	dsn := "postgres://postgres:isma_mesa22@localhost:5433/lacis?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err)
	defer db.Close()

	r := gin.Default()
	api.InitRoutes(r)

	req, _ := http.NewRequest("GET", "/api/v1/admin/proyectos-todos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "No autenticado")
}

func TestIntegration_PaginasPublicasSiguenAbiertasSinSesion(t *testing.T) {
	dsn := "postgres://postgres:isma_mesa22@localhost:5433/lacis?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err)
	defer db.Close()

	r := gin.Default()
	api.InitRoutes(r)

	rutasPublicas := []string{
		"/",
		"/lacis",
		"/tesis",
		"/doctorado-ing-software",
	}
	for _, ruta := range rutasPublicas {
		req, _ := http.NewRequest("GET", ruta, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "ruta pública sin sesión: %s", ruta)
	}
}

func TestIntegration_DesarrollosRedirigeALacisProductos(t *testing.T) {
	dsn := "postgres://postgres:isma_mesa22@localhost:5433/lacis?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err)
	defer db.Close()

	r := gin.Default()
	api.InitRoutes(r)

	req, _ := http.NewRequest("GET", "/desarrollos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusMovedPermanently, w.Code)
	assert.Equal(t, "/lacis#productos-software", w.Header().Get("Location"))
}

func TestIntegration_LacisMuestraProductosDeSoftware(t *testing.T) {
	dsn := "postgres://postgres:isma_mesa22@localhost:5433/lacis?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err)
	defer db.Close()

	r := gin.Default()
	api.InitRoutes(r)

	req, _ := http.NewRequest("GET", "/lacis", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, `id="productos-software"`)

	assert.Contains(t, body, "Productos de")
	assert.Contains(t, body, ">Software</span>")
}
