package usuario

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

type Storage interface {
	Create(usuario *UsuarioGestor) error
	Read(id int) (*UsuarioGestor, error)
	ReadByUsername(username string) (*UsuarioGestor, error)
}

type PostgressStorage struct {
	db *sql.DB
}

func NewPostgressStorage(db *sql.DB) *PostgressStorage /*<-- indica que tipo de dato entrega la funcion al terminar */ {
	return &PostgressStorage{
		db: db, //Aaca se esta asignando
	}
}
func (c *PostgressStorage) Create(u *UsuarioGestor) error { /*En este caso Postgreetorage se utiliza utilizar siempre la mimsa conexxion a la base de datos*/
	query := `INSERT INTO usuario_gestor(integrante_id, username, password_hash, email, rol) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := c.db.QueryRow(query, u.IntegrantesId, u.Username, u.PasswordHash, u.Email, u.Rol).Scan(&u.ID)

	//En este caso el uso de los asterisoc viene dado ya que con c, se evita copiar toda la coenxion a al base de datos
	//Cada vez que se llama un metodo y con u se
	if err != nil {
		return fmt.Errorf("Error al insertar usuario en PostgreSQL: %w", err)
	}
	return nil
}
func (c *PostgressStorage) ReadByUsername(username string) (*UsuarioGestor, error) {
	query := `SELECT id, integrante_id, username, password_hash, email, ultimo_acceso, rol FROM usuario_gestor WHERE username = $1`

	row := c.db.QueryRow(query, username)
	var u UsuarioGestor

	err := row.Scan(
		&u.ID,
		&u.IntegrantesId,
		&u.Username,
		&u.PasswordHash,
		&u.Email,
		&u.UltimoAcceso,
		&u.Rol,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil

}

func (c *PostgressStorage) Read(id int) (*UsuarioGestor, error) {
	query := `SELECT id,integrante_id,username,password_hash,emial,ultimo_acceso,rol FROM usuario_gestor where id = $1`
	row := c.db.QueryRow(query, id)

	var u UsuarioGestor
	err := row.Scan(
		&u.ID,
		&u.IntegrantesId,
		&u.Username,
		&u.PasswordHash,
		&u.Email,
		&u.UltimoAcceso,
		&u.Rol,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}
