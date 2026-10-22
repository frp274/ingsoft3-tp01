package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/frp274/ingsoft3-tp01/backend/internal/models"
	"github.com/frp274/ingsoft3-tp01/backend/internal/services"
)

// MockPlanRepository implementa repository.PlanRepository usando testify/mock.
type MockPlanRepository struct {
	mock.Mock
}

func (m *MockPlanRepository) GetByID(ctx context.Context, id int) (*models.Plan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Plan), args.Error(1)
}

func (m *MockPlanRepository) ListPublicos(ctx context.Context) ([]models.Plan, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Plan), args.Error(1)
}

func (m *MockPlanRepository) ListByOrganizador(ctx context.Context, idOrganizador int) ([]models.Plan, error) {
	args := m.Called(ctx, idOrganizador)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Plan), args.Error(1)
}

func (m *MockPlanRepository) Create(ctx context.Context, p *models.Plan) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

// TestPlanService_CrearPlan_Exitoso verifica la creación exitosa de un plan válido.
func TestPlanService_CrearPlan_Exitoso(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlanRepository)
	service := services.NewPlanService(mockRepo)

	now := time.Now()
	nuevoPlan := &models.Plan{
		Nombre:        "Trekking La Cumbrecita",
		IDOrganizador: 1,
		FechaInicio:   now.Add(24 * time.Hour),
		FechaFin:      now.Add(28 * time.Hour),
		Visibilidad:   "publico",
		IDEtiqueta:    1,
		Ubicacion:     "La Cumbrecita, Córdoba",
		Capacidad:     20,
	}

	mockRepo.On("Create", mock.Anything, nuevoPlan).Return(nil)

	// Act
	err := service.CrearPlan(context.Background(), nuevoPlan)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestPlanService_CrearPlan_Validaciones_TableDriven prueba múltiples reglas de validación usando Table-Driven Tests.
func TestPlanService_CrearPlan_Validaciones_TableDriven(t *testing.T) {
	now := time.Now()

	casos := []struct {
		nombrePrueba  string
		plan          *models.Plan
		errorEsperado string
	}{
		{
			nombrePrueba: "nombre vacio",
			plan: &models.Plan{
				Nombre:      "   ",
				Ubicacion:   "Córdoba",
				Capacidad:   10,
				IDEtiqueta:  1,
				FechaInicio: now,
				FechaFin:    now.Add(2 * time.Hour),
			},
			errorEsperado: "el nombre del plan es obligatorio",
		},
		{
			nombrePrueba: "ubicacion vacia",
			plan: &models.Plan{
				Nombre:      "Cine nocturno",
				Ubicacion:   "",
				Capacidad:   10,
				IDEtiqueta:  1,
				FechaInicio: now,
				FechaFin:    now.Add(2 * time.Hour),
			},
			errorEsperado: "la ubicación del plan es obligatoria",
		},
		{
			nombrePrueba: "capacidad cero o negativa",
			plan: &models.Plan{
				Nombre:      "Cine nocturno",
				Ubicacion:   "Patio Olmos",
				Capacidad:   0,
				IDEtiqueta:  1,
				FechaInicio: now,
				FechaFin:    now.Add(2 * time.Hour),
			},
			errorEsperado: "la capacidad debe ser mayor a 0",
		},
		{
			nombrePrueba: "etiqueta no positiva",
			plan: &models.Plan{
				Nombre:      "Cine nocturno",
				Ubicacion:   "Patio Olmos",
				Capacidad:   10,
				IDEtiqueta:  0,
				FechaInicio: now,
				FechaFin:    now.Add(2 * time.Hour),
			},
			errorEsperado: "la etiqueta es obligatoria",
		},
		{
			nombrePrueba: "visibilidad invalida",
			plan: &models.Plan{
				Nombre:      "Cine nocturno",
				Ubicacion:   "Patio Olmos",
				Capacidad:   10,
				IDEtiqueta:  1,
				Visibilidad: "oculto",
				FechaInicio: now,
				FechaFin:    now.Add(2 * time.Hour),
			},
			errorEsperado: "la visibilidad debe ser 'publico' o 'privado'",
		},
		{
			nombrePrueba: "fecha de fin anterior a inicio",
			plan: &models.Plan{
				Nombre:      "Cine nocturno",
				Ubicacion:   "Patio Olmos",
				Capacidad:   10,
				IDEtiqueta:  1,
				FechaInicio: now.Add(5 * time.Hour),
				FechaFin:    now.Add(2 * time.Hour),
			},
			errorEsperado: "la fecha de fin debe ser posterior a la fecha de inicio",
		},
		{
			nombrePrueba: "fechas de inicio y fin identicas",
			plan: &models.Plan{
				Nombre:      "Cine nocturno",
				Ubicacion:   "Patio Olmos",
				Capacidad:   10,
				IDEtiqueta:  1,
				FechaInicio: now,
				FechaFin:    now,
			},
			errorEsperado: "la fecha de fin debe ser posterior a la fecha de inicio",
		},
	}

	for _, tc := range casos {
		t.Run(tc.nombrePrueba, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockPlanRepository)
			service := services.NewPlanService(mockRepo)

			// Act
			err := service.CrearPlan(context.Background(), tc.plan)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.errorEsperado)
			mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
		})
	}
}

// TestPlanService_CrearPlan_PlanNulo_CasoBorde verifica el manejo de entrada nula (caso borde).
func TestPlanService_CrearPlan_PlanNulo_CasoBorde(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlanRepository)
	service := services.NewPlanService(mockRepo)

	// Act
	err := service.CrearPlan(context.Background(), nil)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "el plan no puede ser nulo", err.Error())
	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

// TestPlanService_GetPlanByID_Exitoso verifica la búsqueda de un plan por ID existente.
func TestPlanService_GetPlanByID_Exitoso(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlanRepository)
	service := services.NewPlanService(mockRepo)

	planEsperado := &models.Plan{ID: 1, Nombre: "Trekking Champaquí"}
	mockRepo.On("GetByID", mock.Anything, 1).Return(planEsperado, nil)

	// Act
	plan, err := service.GetPlanByID(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, planEsperado, plan)
	mockRepo.AssertExpectations(t)
}

// TestPlanService_GetPlanByID_IDInvalido_CasoBorde prueba el caso borde con IDs no positivos.
func TestPlanService_GetPlanByID_IDInvalido_CasoBorde(t *testing.T) {
	idsInvalidos := []int{0, -1, -99}

	for _, id := range idsInvalidos {
		// Arrange
		mockRepo := new(MockPlanRepository)
		service := services.NewPlanService(mockRepo)

		// Act
		plan, err := service.GetPlanByID(context.Background(), id)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, plan)
		assert.Contains(t, err.Error(), "el id del plan debe ser mayor a 0")
		mockRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	}
}

// TestPlanService_GetPlanesPublicos_Exitoso verifica la obtención correcta de planes públicos.
func TestPlanService_GetPlanesPublicos_Exitoso(t *testing.T) {
	// Arrange
	mockRepo := new(MockPlanRepository)
	service := services.NewPlanService(mockRepo)

	planesEsperados := []models.Plan{
		{ID: 1, Nombre: "Trekking", Visibilidad: "publico"},
		{ID: 2, Nombre: "Torneo de Padel", Visibilidad: "publico"},
	}

	mockRepo.On("ListPublicos", mock.Anything).Return(planesEsperados, nil)

	// Act
	resultado, err := service.GetPlanesPublicos(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resultado))
	assert.Equal(t, "Trekking", resultado[0].Nombre)
	mockRepo.AssertExpectations(t)
}
