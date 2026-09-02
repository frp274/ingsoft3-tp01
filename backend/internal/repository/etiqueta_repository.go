package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/frp274/ingsoft3-tp01/backend/internal/models"
)

// EtiquetaRepository define los métodos de acceso a datos para Etiqueta.
type EtiquetaRepository interface {
	GetByID(ctx context.Context, id int) (*models.Etiqueta, error)
	List(ctx context.Context) ([]models.Etiqueta, error)
	Create(ctx context.Context, e *models.Etiqueta) error
}

type mysqlEtiquetaRepository struct {
	db *sql.DB
}

// NewEtiquetaRepository crea una nueva instancia inyectando el pool de conexiones.
func NewEtiquetaRepository(db *sql.DB) EtiquetaRepository {
	return &mysqlEtiquetaRepository{db: db}
}

func (r *mysqlEtiquetaRepository) GetByID(ctx context.Context, id int) (*models.Etiqueta, error) {
	query := `SELECT id, nombre FROM etiquetas WHERE id = ?`
	var e models.Etiqueta
	err := r.db.QueryRowContext(ctx, query, id).Scan(&e.ID, &e.Nombre)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener etiqueta por id: %w", err)
	}
	return &e, nil
}

func (r *mysqlEtiquetaRepository) List(ctx context.Context) ([]models.Etiqueta, error) {
	query := `SELECT id, nombre FROM etiquetas ORDER BY id ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al listar etiquetas: %w", err)
	}
	defer rows.Close()

	var etiquetas []models.Etiqueta
	for rows.Next() {
		var e models.Etiqueta
		if err := rows.Scan(&e.ID, &e.Nombre); err != nil {
			return nil, err
		}
		etiquetas = append(etiquetas, e)
	}
	return etiquetas, nil
}

func (r *mysqlEtiquetaRepository) Create(ctx context.Context, e *models.Etiqueta) error {
	query := `INSERT INTO etiquetas (nombre) VALUES (?)`
	res, err := r.db.ExecContext(ctx, query, e.Nombre)
	if err != nil {
		return fmt.Errorf("error al insertar etiqueta: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	e.ID = int(id)
	return nil
}