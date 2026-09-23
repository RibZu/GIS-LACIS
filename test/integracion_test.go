package test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

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
	// t.Cleanup (no defer): el DELETE registrado abajo corre antes que este Close (orden LIFO),
	// igual que en TestIntegration_LacisMuestraProductosDeSoftware.
	t.Cleanup(func() { db.Close() })
	t.Cleanup(func() {
		_, err := db.Exec(
			`DELETE FROM integrante WHERE nombre = $1 AND apellido = $2 AND contacto_mail = $3`,
			"Test", "Integracion", "test@unsl.edu.ar",
		)
		assert.NoError(t, err, "limpieza del integrante sembrado por el test")
	})

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
	assert.Contains(t, w.Header().Get("Location"), "/admin/integrantes")
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
	// Mismo orden LIFO que en TestIntegration_InsertarIntegrante: el DELETE corre antes del Close.
	// Los integrantes_tesis de la tesis se borran solos por ON DELETE CASCADE.
	t.Cleanup(func() { db.Close() })
	t.Cleanup(func() {
		_, err := db.Exec(`DELETE FROM tesis WHERE titulo = $1`, "Tesis Test Integración")
		assert.NoError(t, err, "limpieza de la tesis sembrada por el test")
	})

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
	assert.Contains(t, w.Header().Get("Location"), "/admin/tesis")
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
	// Mismo orden LIFO que en TestIntegration_InsertarIntegrante. Esta prueba no debería insertar
	// nada (es justo lo que verifica); el DELETE es una red de seguridad por si una regresión
	// dejara pasar el alta sin sesión.
	t.Cleanup(func() { db.Close() })
	t.Cleanup(func() {
		_, err := db.Exec(
			`DELETE FROM integrante WHERE nombre = $1 AND apellido = $2`,
			"NoDeberia", "Insertarse",
		)
		assert.NoError(t, err, "limpieza del integrante que no debería haberse insertado")
	})

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
	// t.Cleanup (no defer) para que el orden de LIFO sea el correcto frente al cleanup de
	// borrado registrado abajo: los t.Cleanup se ejecutan en orden inverso al de registro, así
	// que este Close (registrado primero) corre DESPUÉS del DELETE (registrado después) — con
	// `defer db.Close()` el cierre ocurriría antes de que t.Cleanup llegue a borrar nada.
	t.Cleanup(func() { db.Close() })

	r := gin.Default()
	api.InitRoutes(r)

	// Sembrado con limpieza garantizada al final (pase o falle el test): estos títulos son
	// exclusivos de este test, así que un DELETE por título no toca ningún desarrollo real. Sin
	// esto, cada corrida del test deja basura permanente en la base — y como "año actual + 1"
	// siempre gana el orden "más reciente primero", esa basura sepulta los productos reales en
	// la vista pública tras la primera corrida.
	tituloReciente := "ZZZ Test Recorte Reciente"
	tituloAntiguo := "AAA Test Recorte Antiguo"
	tituloConParticipantes := "Test Producto Con Participantes"
	t.Cleanup(func() {
		_, err := db.Exec(
			`DELETE FROM desarrollo WHERE titulo IN ($1, $2, $3)`,
			tituloReciente, tituloAntiguo, tituloConParticipantes,
		)
		assert.NoError(t, err, "limpieza de desarrollos sembrados por el test")
	})

	seedDesarrollo := func(titulo string, anio int, integranteIDs []string, externos []string) {
		formData := url.Values{}
		formData.Set("titulo", titulo)
		formData.Set("anio", strconv.Itoa(anio))
		for _, id := range integranteIDs {
			formData.Add("integrantes", id)
		}
		for _, ext := range externos {
			formData.Add("externos", ext)
		}
		req, _ := http.NewRequest("POST", "/admin/insertar-desarrollo", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{Name: "session", Value: "1"})
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusSeeOther, w.Code)
	}

	nombreExterno := "Participante Externo Test QA"
	seedDesarrollo(tituloReciente, time.Now().Year()+1, nil, nil)
	seedDesarrollo(tituloAntiguo, 1990, nil, nil)
	seedDesarrollo(tituloConParticipantes, 2024, []string{"5"}, []string{nombreExterno})

	var total int
	assert.NoError(t, db.QueryRow("SELECT COUNT(*) FROM desarrollo").Scan(&total))

	req, _ := http.NewRequest("GET", "/lacis", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, `id="productos-software"`)

	assert.Contains(t, body, "Productos de")
	assert.Contains(t, body, ">Software</span>")

	// US1: recorte a 6 — el HTML trae una tarjeta ("producto-item") por cada desarrollo, pero
	// las que exceden las 6 más recientes llevan además la clase "producto-oculto".
	assert.Equal(t, total, strings.Count(body, "producto-item"),
		"cada desarrollo debe tener su tarjeta en el HTML, aunque esté oculta")
	esperadasOcultas := total - 6
	if esperadasOcultas < 0 {
		esperadasOcultas = 0
	}
	assert.Equal(t, esperadasOcultas, strings.Count(body, "producto-oculto"))
	assert.Less(t, strings.Index(body, tituloReciente), strings.Index(body, tituloAntiguo),
		"el producto más reciente debe aparecer antes que el más antiguo")

	// US2: con más de 6 desarrollos (garantizado por el sembrado de arriba), debe aparecer el
	// control "ver más".
	assert.Contains(t, body, `id="btn-ver-mas-productos"`)

	// US3: el integrante vinculado y el participante externo sembrados arriba deben verse en la
	// tarjeta pública, igual que los vería un administrador en /admin/desarrollos.
	assert.Contains(t, body, "Mg. Baigorria")
	assert.Contains(t, body, nombreExterno)
}
