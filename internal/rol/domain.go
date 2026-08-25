package rol

type Rol struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"` // Ej: "INVESTIGADOR", "DIRECTOR", "ESTUDIANTE"
}
