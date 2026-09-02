package services

import (
	"context"
	"fmt"

	"github.com/frp274/ingsoft3-tp01/backend/internal/models"
	"github.com/frp274/ingsoft3-tp01/backend/internal/repository"
)

// InvitacionService define la lógica de negocio para solicitudes e inscripciones.
type InvitacionService interface {
	SolicitarInvitacion(ctx context.Context, idUsuario int, idPlan int) (*models.Invitacion, error)
	GetMisPlanes(ctx context.Context, idUsuario int) ([]models.PlanConInvitacion, error)
	GetMisInvitaciones(ctx context.Context, idUsuario int) ([]models.Invitacion, error)
}

type invitacionService struct {
	invitacionRepo repository.InvitacionRepository
	planRepo       repository.PlanRepository
}

// NewInvitacionService crea una nueva instancia del servicio de invitaciones.
func NewInvitacionService(invitacionRepo repository.InvitacionRepository, planRepo repository.PlanRepository) InvitacionService {
	return &invitacionService{
		invitacionRepo: invitacionRepo,
		planRepo:       planRepo,
	}
}

// SolicitarInvitacion registra una solicitud con estado 'pendiente'.
func (s *invitacionService) SolicitarInvitacion(ctx context.Context, idUsuario int, idPlan int) (*models.Invitacion, error) {
	// Verificar que el plan exista
	plan, err := s.planRepo.GetByID(ctx, idPlan)
	if err != nil {
		return nil, fmt.Errorf("error al verificar existencia del plan: %w", err)
	}
	if plan == nil {
		return nil, fmt.Errorf("el plan especificado (id %d) no existe", idPlan)
	}

	// Verificar si ya existe una invitación previa
	existente, err := s.invitacionRepo.GetByUsuarioYPlan(ctx, idUsuario, idPlan)
	if err != nil {
		return nil, fmt.Errorf("error al verificar invitaciones previas: %w", err)
	}
	if existente != nil {
		return existente, nil // Ya está solicitada o inscripta
	}

	nuevaInv := &models.Invitacion{
		IDPlan:    idPlan,
		IDUsuario: idUsuario,
		Estado:    "pendiente",
	}

	if err := s.invitacionRepo.Create(ctx, nuevaInv); err != nil {
		return nil, fmt.Errorf("error al crear invitación: %w", err)
	}

	return nuevaInv, nil
}

// GetMisPlanes obtiene los planes a los que el usuario se postuló mediante JOIN.
func (s *invitacionService) GetMisPlanes(ctx context.Context, idUsuario int) ([]models.PlanConInvitacion, error) {
	return s.invitacionRepo.ListPlanesConInvitacionByUsuario(ctx, idUsuario)
}

// GetMisInvitaciones obtiene la lista de invitaciones del usuario.
func (s *invitacionService) GetMisInvitaciones(ctx context.Context, idUsuario int) ([]models.Invitacion, error) {
	return s.invitacionRepo.ListByUsuario(ctx, idUsuario)
}