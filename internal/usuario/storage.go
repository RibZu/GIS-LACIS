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
	GetAll() ([]UsuarioGestor, error)
	Update(id int, usuario *UpdateFieldGestor) error
	Delete(id int) error
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
	query := `INSERT INTO usuario_gestor(integrante_id, username, password_hash, email, rol, modulos) VALUES ($1, $2, $3, $4, $5,$6) RETURNING id`
	err := c.db.QueryRow(query, u.IntegrantesId, u.Username, u.PasswordHash, u.Email, u.Rol, u.Modulos).Scan(&u.ID)

	//En este caso el uso de los asterisoc viene dado ya que con c, se evita copiar toda la coenxion a al base de datos
	//Cada vez que se llama un metodo y con u se
	if err != nil {
		return fmt.Errorf("Error al insertar usuario en PostgreSQL: %w", err)
	}
	return nil
}
func (c *PostgressStorage) ReadByUsername(username string) (*UsuarioGestor, error) {
	query := `SELECT id, integrante_id, username, password_hash, email, ultimo_acceso, rol, modulos FROM usuario_gestor WHERE username = $1`

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
		&u.Modulos,
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
	query := `SELECT id,integrante_id,username,password_hash,email,ultimo_acceso,rol,modulos FROM usuario_gestor where id = $1`
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
		&u.Modulos,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (c *PostgressStorage) GetAll() ([]UsuarioGestor, error) {
	query := `SELECT id, username, password_hash, email, ultimo_acceso, rol, modulos FROM usuario_gestor`
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var usuarios []UsuarioGestor
	for rows.Next() {
		var u UsuarioGestor
		err := rows.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Email, &u.UltimoAcceso, &u.Rol, &u.Modulos) /*Vos al scan le podes pasar tantos parametros como direcciones recibe la consulta*/
		if err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return usuarios, nil
}

func (c *PostgressStorage) Delete(id int) error {
	query := `DELETE FROM usuario_gestor WHERE id = $1`
	res, err := c.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error al eliminar usuario de PostgreSQL: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	return nil
}
func (c *PostgressStorage) Update(id int, fields *UpdateFieldGestor) error {
	query := `UPDATE usuario_gestor SET `
	var args []interface{} /*Esto sirve para guardar datos de diferentes tipos, ya que go cuando uno define un slice en go solo puede guardar
	valores de un tipo pero si lo definimos como interface/any le decimos que puede recibir varios tipos de valores*/
	argID := 1
	if fields.Username != nil {
		query += fmt.Sprintf("username=$%d, ", argID)
		args = append(args, *fields.Username)
		argID++
	}
	if fields.PasswordHash != nil {
		query += fmt.Sprintf("password_hash=$%d, ", argID)
		args = append(args, *fields.PasswordHash)
		argID++
	}
	if fields.Email != nil {
		query += fmt.Sprintf("email=$%d, ", argID)
		args = append(args, *fields.Email)
		argID++
	}
	if fields.Rol != nil {
		query += fmt.Sprintf("rol=$%d, ", argID)
		args = append(args, *fields.Rol)
		argID++
	}
	if fields.Modulos != nil {
		query += fmt.Sprintf("modulos=$%d, ", argID)
		args = append(args, *fields.Modulos)
		argID++
	}
	if len(args) == 0 {
		return nil
	}

	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d", argID)
	args = append(args, id)
	res, err := c.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error al actualizar usuario en la base de datos: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
