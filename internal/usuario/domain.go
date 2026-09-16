package usuario

type UsuarioGestor struct {
	ID int `json:"id"`

	IntegrantesId *int `json:"integrantes_id"`

	Username *string `json:"username"`

	PasswordHash *string `json:"password_hash"`

	Email *string `json:"email"`

	UltimoAcceso *string `json:"ultimo_acceso"`
	Rol          *string `json:"rol"`
	Modulos      *string `json:"modulos"`
}

type UpdateFieldGestor struct {
	Username     *string `json:"username"`
	PasswordHash *string `json:"password_hash"`
	Email        *string `json:"email"`
	Rol          *string `json:"rol"`
	Modulos      *string `json:"modulos"`
}
