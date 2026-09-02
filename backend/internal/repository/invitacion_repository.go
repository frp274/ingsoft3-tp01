package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/frp274/ingsoft3-tp01/backend/internal/models"
)

// InvitacionRepository define los métodos de acceso a datos para Invitacion.
type InvitacionRepository interface {
	GetByID(ctx context.Context, idInvitacion int) (*models.Invitacion, error)
	ListByUsuario(ctx context.Context, idUsuario int) ([]models.Invitacion, error)
	ListByPlan(ctx context.Context, idPlan int) ([]models.Invitacion, error)
	Create(ctx context.Context, inv *models.Invitacion) error
	UpdateEstado(ctx context.Context, idInvitacion int, estado string) error
}

type mysqlInvitacionRepository struct {
	db *sql.DB
}

// NewInvitacionRepository crea una nueva instancia inyectando el pool de conexiones.
func NewInvitacionRepository(db *sql.DB) InvitacionRepository {
	return &mysqlInvitacionRepository{db: db}
}

func (r *mysqlInvitacionRepository) GetByID(ctx context.Context, idInvitacion int) (*models.Invitacion, error) {
	query := `SELECT idInvitacion, idPlan, idUsuario, estado FROM invitaciones WHERE idInvitacion = ?`
	var inv models.Invitacion
	err := r.db.QueryRowContext(ctx, query, idInvitacion).Scan(&inv.IDInvitacion, &inv.IDPlan, &inv.IDUsuario, &inv.Estado)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener invitación por id: %w", err)
	}
	return &inv, nil
}

func (r *mysqlInvitacionRepository) ListByUsuario(ctx context.Context, idUsuario int) ([]models.Invitacion, error) {
	query := `SELECT idInvitacion, idPlan, idUsuario, estado FROM invitaciones WHERE idUsuario = ?`
	rows, err := r.db.QueryContext(ctx, query, idUsuario)
	if err != nil {
		return nil, fmt.Errorf("error al listar invitaciones por usuario: %w", err)
	}
	defer rows.Close()

	var invitaciones []models.Invitacion
	for rows.Next() {
		var inv models.Invitacion
		if err := rows.Scan(&inv.IDInvitacion, &inv.IDPlan, &inv.IDUsuario, &inv.Estado); err != nil {
			return nil, err
		}
		invitaciones = append(invitaciones, inv)
	}
	return invitaciones, nil
}

func (r *mysqlInvitacionRepository) ListByPlan(ctx context.Context, idPlan int) ([]models.Invitacion, error) {
	query := `SELECT idInvitacion, idPlan, idUsuario, estado FROM invitaciones WHERE idPlan = ?`
	rows, err := r.db.QueryContext(ctx, query, idPlan)
	if err != nil {
		return nil, fmt.Errorf("error al listar invitaciones por plan: %w", err)
	}
	defer rows.Close()

	var invitaciones []models.Invitacion
	for rows.Next() {
		var inv models.Invitacion
		if err := rows.Scan(&inv.IDInvitacion, &inv.IDPlan, &inv.IDUsuario, &inv.Estado); err != nil {
			return nil, err
		}
		invitaciones = append(invitaciones, inv)
	}
	return invitaciones, nil
}

func (r *mysqlInvitacionRepository) Create(ctx context.Context, inv *models.Invitacion) error {
	query := `INSERT INTO invitaciones (idPlan, idUsuario, estado) VALUES (?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, inv.IDPlan, inv.IDUsuario, inv.Estado)
	if err != nil {
		return fmt.Errorf("error al insertar invitación: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	inv.IDInvitacion = int(id)
	return nil
}

func (r *mysqlInvitacionRepository) UpdateEstado(ctx context.Context, idInvitacion int, estado string) error {
	query := `UPDATE invitaciones SET estado = ? WHERE idInvitacion = ?`
	_, err := r.db.ExecContext(ctx, query, estado, idInvitacion)
	if err != nil {
		return fmt.Errorf("error al actualizar estado de la invitación: %w", err)
	}
	return nil
}