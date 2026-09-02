package services

import (
	"context"

	"github.com/frp274/ingsoft3-tp01/backend/internal/models"
	"github.com/frp274/ingsoft3-tp01/backend/internal/repository"
)

// PlanService define la lógica de negocio asociada a los planes.
type PlanService interface {
	GetPlanesPublicos(ctx context.Context) ([]models.Plan, error)
	GetPlanByID(ctx context.Context, id int) (*models.Plan, error)
	CrearPlan(ctx context.Context, plan *models.Plan) error
}

type planService struct {
	planRepo repository.PlanRepository
}

// NewPlanService crea una nueva instancia del servicio de planes.
func NewPlanService(planRepo repository.PlanRepository) PlanService {
	return &planService{
		planRepo: planRepo,
	}
}

// GetPlanesPublicos obtiene únicamente los planes con visibilidad pública.
func (s *planService) GetPlanesPublicos(ctx context.Context) ([]models.Plan, error) {
	return s.planRepo.ListPublicos(ctx)
}

// GetPlanByID obtiene un plan específico por su ID.
func (s *planService) GetPlanByID(ctx context.Context, id int) (*models.Plan, error) {
	return s.planRepo.GetByID(ctx, id)
}

// CrearPlan registra un nuevo plan.
func (s *planService) CrearPlan(ctx context.Context, plan *models.Plan) error {
	return s.planRepo.Create(ctx, plan)
}