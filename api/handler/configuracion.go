package handler

import (
	"net/http"

	"PaginaSEG/internal/configuracion"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ConfiguracionHandler struct {
	service *configuracion.Service
	logger  *zap.Logger
}

func NewConfiguracionHandler(s *configuracion.Service, l *zap.Logger) *ConfiguracionHandler {
	return &ConfiguracionHandler{
		service: s,
		logger:  l,
	}
}

// VerConfiguracion renderiza la vista de administración de configuración general
func (h *ConfiguracionHandler) VerConfiguracion(c *gin.Context) {
	cfg, err := h.service.ObtenerConfiguracion()
	if err != nil {
		h.logger.Error("Error al obtener configuración para vista admin", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error al cargar la configuración")
		return
	}

	_, loggedIn := CurrentUserID(c)
	c.HTML(http.StatusOK, "Configuracion.html", gin.H{
		"LoggedIn":      loggedIn,
		"EsAdmin":       true,
		"Configuracion": cfg,
	})
}

// GuardarConfiguracion recibe los datos del formulario o JSON y persiste los cambios
func (h *ConfiguracionHandler) GuardarConfiguracion(c *gin.Context) {
	var dto configuracion.UpdateConfiguracionDTO

	if c.ContentType() == "application/json" {
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
			return
		}
	} else {
		// Parsea formulario POST tradicional
		doc := c.PostForm("email_doctorado")
		mCal := c.PostForm("email_maestria_calidad")
		mSoft := c.PostForm("email_maestria_software")
		esp := c.PostForm("email_especializacion")
		ingInf := c.PostForm("email_ing_informatica")
		tecWeb := c.PostForm("email_tecnicatura_web")
		tel := c.PostForm("telefono_footer")
		mailFoot := c.PostForm("email_footer")
		dir := c.PostForm("direccion")
		ubic := c.PostForm("ubicacion_footer")
		ofic := c.PostForm("oficinas_footer")
		mapa := c.PostForm("mapa_url")

		dto = configuracion.UpdateConfiguracionDTO{
			EmailDoctorado:        &doc,
			EmailMaestriaCalidad:  &mCal,
			EmailMaestriaSoftware: &mSoft,
			EmailEspecializacion:  &esp,
			EmailIngInformatica:   &ingInf,
			EmailTecnicaturaWeb:   &tecWeb,
			TelefonoFooter:        &tel,
			EmailFooter:           &mailFoot,
			Direccion:             &dir,
			UbicacionFooter:       &ubic,
			OficinasFooter:        &ofic,
			MapaURL:               &mapa,
		}
	}

	cfgActualizada, err := h.service.GuardarConfiguracion(dto)
	if err != nil {
		h.logger.Warn("Error de validación o persistencia en configuración", zap.Error(err))
		if c.ContentType() == "application/json" || c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusBadRequest, "Configuracion.html", gin.H{
			"LoggedIn":      true,
			"EsAdmin":       true,
			"Configuracion": dto,
			"Error":         err.Error(),
		})
		return
	}

	if c.ContentType() == "application/json" || c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
		c.JSON(http.StatusOK, gin.H{
			"mensaje":       "Configuración guardada exitosamente",
			"configuracion": cfgActualizada,
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/configuracion?status=guardado")
}

// API_GetConfiguracion endpoint público para consultar configuración
func (h *ConfiguracionHandler) API_GetConfiguracion(c *gin.Context) {
	cfg, err := h.service.ObtenerConfiguracion()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al consultar la configuración"})
		return
	}
	c.JSON(http.StatusOK, cfg)
}
