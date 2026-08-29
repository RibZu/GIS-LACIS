package usuario

type UsuarioGestor struct {
	ID int `json:"id"`
	//Se pone json id como una manera de reformular la la variblae para pasarla en el formato
	//que debe estar en el json
	IntegrantesId *int `json:"integrantes_id"` //Con el operador * en este caso lo utilizamos para que go sepa que esa varaible puede estar vacia

	Username *string `json:"username"`

	PasswordHash *string `json:"password_hash"`

	Email *string `json:"email"`

	UltimoAcceso *string `json:"ultimo_acceso"`
	Rol          *string `json:"rol"`
}

type UpdateFieldGestor struct {
	Username     *string `json:"username"`
	PasswordHash *string `json:"password_hash"`
	Email        *string `json:"email"`
	Rol          *string `json:"rol"`
	//En este caso se estan usando punteros en los String ya que estos permiten que el campo sea nulo en go
	//Es decir que se puede dejar Vacio
}
