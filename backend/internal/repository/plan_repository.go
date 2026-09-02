package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/frp274/ingsoft3-tp01/backend/internal/models"
)

// PlanRepository define los métodos de acceso a datos para Plan.
type PlanRepository interface {
	GetByID(ctx context.Context, id int) (*models.Plan, error)
	ListPublicos(ctx context.Context) ([]models.Plan, error)
	ListByOrganizador(ctx context.Context, idOrganizador int) ([]models.Plan, error)
	Create(ctx context.Context, p *models.Plan) error
}

type mysqlPlanRepository struct {
	db *sql.DB
}

// NewPlanRepository crea una nueva instancia inyectando el pool de conexiones.
func NewPlanRepository(db *sql.DB) PlanRepository {
	return &mysqlPlanRepository{db: db}
}

func (r *mysqlPlanRepository) GetByID(ctx context.Context, id int) (*models.Plan, error) {
	query := `SELECT id, nombre, idOrganizador, fechaInicio, fechaFin, visibilidad, descripcion, idEtiqueta, ubicacion, capacidad 
		FROM planes WHERE id = ?`
	var p models.Plan
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Nombre, &p.IDOrganizador, &p.FechaInicio, &p.FechaFin, &p.Visibilidad, &p.Descripcion, &p.IDEtiqueta, &p.Ubicacion, &p.Capacidad,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener plan por id: %w", err)
	}
	return &p, nil
}

func (r *mysqlPlanRepository) ListPublicos(ctx context.Context) ([]models.Plan, error) {
	query := `SELECT id, nombre, idOrganizador, fechaInicio, fechaFin, visibilidad, descripcion, idEtiqueta, ubicacion, capacidad 
		FROM planes WHERE visibilidad = 'publico' ORDER BY fechaInicio ASC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al listar planes públicos: %w", err)
	}
	defer rows.Close()

	var planes []models.Plan
	for rows.Next() {
		var p models.Plan
		if err := rows.Scan(
			&p.ID, &p.Nombre, &p.IDOrganizador, &p.FechaInicio, &p.FechaFin, &p.Visibilidad, &p.Descripcion, &p.IDEtiqueta, &p.Ubicacion, &p.Capacidad,
		); err != nil {
			return nil, err
		}
		planes = append(planes, p)
	}
	return planes, nil
}

func (r *mysqlPlanRepository) ListByOrganizador(ctx context.Context, idOrganizador int) ([]models.Plan, error) {
	query := `SELECT id, nombre, idOrganizador, fechaInicio, fechaFin, visibilidad, descripcion, idEtiqueta, ubicacion, capacidad 
		FROM planes WHERE idOrganizador = ? ORDER BY fechaInicio ASC`
	rows, err := r.db.QueryContext(ctx, query, idOrganizador)
	if err != nil {
		return nil, fmt.Errorf("error al listar planes por organizador: %w", err)
	}
	defer rows.Close()

	var planes []models.Plan
	for rows.Next() {
		var p models.Plan
		if err := rows.Scan(
			&p.ID, &p.Nombre, &p.IDOrganizador, &p.FechaInicio, &p.FechaFin, &p.Visibilidad, &p.Descripcion, &p.IDEtiqueta, &p.Ubicacion, &p.Capacidad,
		); err != nil {
			return nil, err
		}
		planes = append(planes, p)
	}
	return planes, nil
}

func (r *mysqlPlanRepository) Create(ctx context.Context, p *models.Plan) error {
	query := `INSERT INTO planes (nombre, idOrganizador, fechaInicio, fechaFin, visibilidad, descripcion, idEtiqueta, ubicacion, capacidad)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, p.Nombre, p.IDOrganizador, p.FechaInicio, p.FechaFin, p.Visibilidad, p.Descripcion, p.IDEtiqueta, p.Ubicacion, p.Capacidad)
	if err != nil {
		return fmt.Errorf("error al insertar plan: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)
	return nil
}