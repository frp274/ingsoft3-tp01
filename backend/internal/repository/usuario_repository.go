package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/frp274/ingsoft3-tp01/backend/internal/models"
)

// UsuarioRepository define los métodos de acceso a datos para Usuario.
type UsuarioRepository interface {
	GetByID(ctx context.Context, id int) (*models.Usuario, error)
	GetByEmail(ctx context.Context, email string) (*models.Usuario, error)
	Create(ctx context.Context, u *models.Usuario) error
	List(ctx context.Context) ([]models.Usuario, error)
}

type mysqlUsuarioRepository struct {
	db *sql.DB
}

// NewUsuarioRepository crea una nueva instancia inyectando el pool de conexiones.
func NewUsuarioRepository(db *sql.DB) UsuarioRepository {
	return &mysqlUsuarioRepository{db: db}
}

func (r *mysqlUsuarioRepository) GetByID(ctx context.Context, id int) (*models.Usuario, error) {
	query := `SELECT id, email, contrasena, nombre, apellido, ciudad, edad, etiquetaFavorita 
		FROM usuarios WHERE id = ?`
	var u models.Usuario
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID, &u.Email, &u.Contrasena, &u.Nombre, &u.Apellido, &u.Ciudad, &u.Edad, &u.EtiquetaFavorita,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario por id: %w", err)
	}
	return &u, nil
}

func (r *mysqlUsuarioRepository) GetByEmail(ctx context.Context, email string) (*models.Usuario, error) {
	query := `SELECT id, email, contrasena, nombre, apellido, ciudad, edad, etiquetaFavorita 
		FROM usuarios WHERE email = ?`
	var u models.Usuario
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID, &u.Email, &u.Contrasena, &u.Nombre, &u.Apellido, &u.Ciudad, &u.Edad, &u.EtiquetaFavorita,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuario por email: %w", err)
	}
	return &u, nil
}

func (r *mysqlUsuarioRepository) Create(ctx context.Context, u *models.Usuario) error {
	query := `INSERT INTO usuarios (email, contrasena, nombre, apellido, ciudad, edad, etiquetaFavorita)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, u.Email, u.Contrasena, u.Nombre, u.Apellido, u.Ciudad, u.Edad, u.EtiquetaFavorita)
	if err != nil {
		return fmt.Errorf("error al insertar usuario: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = int(id)
	return nil
}

func (r *mysqlUsuarioRepository) List(ctx context.Context) ([]models.Usuario, error) {
	query := `SELECT id, email, contrasena, nombre, apellido, ciudad, edad, etiquetaFavorita FROM usuarios`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al listar usuarios: %w", err)
	}
	defer rows.Close()

	var usuarios []models.Usuario
	for rows.Next() {
		var u models.Usuario
		if err := rows.Scan(&u.ID, &u.Email, &u.Contrasena, &u.Nombre, &u.Apellido, &u.Ciudad, &u.Edad, &u.EtiquetaFavorita); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}
	return usuarios, nil
}