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

