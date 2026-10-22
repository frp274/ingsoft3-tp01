package services

import (
	"context"
	"errors"
	"strings"

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

// NewPlanService crea una nueva instancia del servicio de planes con inyección de dependencias.
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
	if id <= 0 {
		return nil, errors.New("el id del plan debe ser mayor a 0")
	}
	return s.planRepo.GetByID(ctx, id)
}

// CrearPlan valida las reglas de negocio y persiste un nuevo plan.
func (s *planService) CrearPlan(ctx context.Context, plan *models.Plan) error {
	if plan == nil {
		return errors.New("el plan no puede ser nulo")
	}

	plan.Nombre = strings.TrimSpace(plan.Nombre)
	if plan.Nombre == "" {
		return errors.New("el nombre del plan es obligatorio")
	}

	plan.Ubicacion = strings.TrimSpace(plan.Ubicacion)
	if plan.Ubicacion == "" {
		return errors.New("la ubicación del plan es obligatoria")
	}

	if plan.Capacidad <= 0 {
		return errors.New("la capacidad debe ser mayor a 0")
	}

	if plan.IDEtiqueta <= 0 {
		return errors.New("la etiqueta es obligatoria")
	}

	if plan.Visibilidad == "" {
		plan.Visibilidad = "publico"
	}
	if plan.Visibilidad != "publico" && plan.Visibilidad != "privado" {
		return errors.New("la visibilidad debe ser 'publico' o 'privado'")
	}

	if plan.FechaInicio.IsZero() {
		return errors.New("la fecha de inicio es obligatoria")
	}
	if plan.FechaFin.IsZero() {
		return errors.New("la fecha de fin es obligatoria")
	}
	if !plan.FechaFin.After(plan.FechaInicio) {
		return errors.New("la fecha de fin debe ser posterior a la fecha de inicio")
	}

	return s.planRepo.Create(ctx, plan)
}
