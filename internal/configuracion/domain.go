package configuracion

type ConfiguracionSitio struct {
	ID                    int    `json:"id"`
	UsuarioGestorID       *int   `json:"usuario_gestor_id"`
	TelefonoFooter        string `json:"telefono_footer"`
	EmailFooter           string `json:"email_footer"`
	Direccion             string `json:"direccion"`
	UbicacionFooter       string `json:"ubicacion_footer"`
	OficinasFooter        string `json:"oficinas_footer"`
	MapaURL               string `json:"mapa_url"`
	EmailDoctorado        string `json:"email_doctorado"`
	EmailMaestriaCalidad  string `json:"email_maestria_calidad"`
	EmailMaestriaSoftware string `json:"email_maestria_software"`
	EmailEspecializacion  string `json:"email_especializacion"`
	EmailIngInformatica   string `json:"email_ing_informatica"`
	EmailTecnicaturaWeb   string `json:"email_tecnicatura_web"`
}

type UpdateConfiguracionDTO struct {
	TelefonoFooter        *string `json:"telefono_footer"`
	EmailFooter           *string `json:"email_footer"`
	Direccion             *string `json:"direccion"`
	UbicacionFooter       *string `json:"ubicacion_footer"`
	OficinasFooter        *string `json:"oficinas_footer"`
	MapaURL               *string `json:"mapa_url"`
	EmailDoctorado        *string `json:"email_doctorado"`
	EmailMaestriaCalidad  *string `json:"email_maestria_calidad"`
	EmailMaestriaSoftware *string `json:"email_maestria_software"`
	EmailEspecializacion  *string `json:"email_especializacion"`
	EmailIngInformatica   *string `json:"email_ing_informatica"`
	EmailTecnicaturaWeb   *string `json:"email_tecnicatura_web"`
}
