package services_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/frp274/ingsoft3-tp01/backend/internal/models"
	"github.com/frp274/ingsoft3-tp01/backend/internal/services"
)

// MockInvitacionRepository implementa repository.InvitacionRepository usando testify/mock.
type MockInvitacionRepository struct {
	mock.Mock
}

func (m *MockInvitacionRepository) GetByID(ctx context.Context, idInvitacion int) (*models.Invitacion, error) {
	args := m.Called(ctx, idInvitacion)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Invitacion), args.Error(1)
}

func (m *MockInvitacionRepository) GetByUsuarioYPlan(ctx context.Context, idUsuario int, idPlan int) (*models.Invitacion, error) {
	args := m.Called(ctx, idUsuario, idPlan)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Invitacion), args.Error(1)
}

func (m *MockInvitacionRepository) ListByUsuario(ctx context.Context, idUsuario int) ([]models.Invitacion, error) {
	args := m.Called(ctx, idUsuario)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Invitacion), args.Error(1)
}

func (m *MockInvitacionRepository) ListByPlan(ctx context.Context, idPlan int) ([]models.Invitacion, error) {
	args := m.Called(ctx, idPlan)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Invitacion), args.Error(1)
}

func (m *MockInvitacionRepository) ListPlanesConInvitacionByUsuario(ctx context.Context, idUsuario int) ([]models.PlanConInvitacion, error) {
	args := m.Called(ctx, idUsuario)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PlanConInvitacion), args.Error(1)
}

func (m *MockInvitacionRepository) Create(ctx context.Context, inv *models.Invitacion) error {
	args := m.Called(ctx, inv)
	return args.Error(0)
}

func (m *MockInvitacionRepository) UpdateEstado(ctx context.Context, idInvitacion int, estado string) error {
	args := m.Called(ctx, idInvitacion, estado)
	return args.Error(0)
}

// TestInvitacionService_SolicitarInvitacion_Exitoso verifica el flujo exitoso cuando no hay solicitud previa.
func TestInvitacionService_SolicitarInvitacion_Exitoso(t *testing.T) {
	// Arrange
	mockPlanRepo := new(MockPlanRepository)
	mockInvRepo := new(MockInvitacionRepository)
	service := services.NewInvitacionService(mockInvRepo, mockPlanRepo)

	idUsuario := 1
	idPlan := 10

	mockPlanRepo.On("GetByID", mock.Anything, idPlan).Return(&models.Plan{ID: idPlan, Nombre: "Torneo"}, nil)
	mockInvRepo.On("GetByUsuarioYPlan", mock.Anything, idUsuario, idPlan).Return(nil, nil)
	mockInvRepo.On("Create", mock.Anything, mock.MatchedBy(func(inv *models.Invitacion) bool {
		return inv.IDPlan == idPlan && inv.IDUsuario == idUsuario && inv.Estado == "pendiente"
	})).Return(nil)

	// Act
	inv, err := service.SolicitarInvitacion(context.Background(), idUsuario, idPlan)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, inv)
	assert.Equal(t, idPlan, inv.IDPlan)
	assert.Equal(t, idUsuario, inv.IDUsuario)
	assert.Equal(t, "pendiente", inv.Estado)
	mockPlanRepo.AssertExpectations(t)
	mockInvRepo.AssertExpectations(t)
}

// TestInvitacionService_SolicitarInvitacion_PlanNoExiste_RetornaError valida error cuando el plan no existe.
func TestInvitacionService_SolicitarInvitacion_PlanNoExiste_RetornaError(t *testing.T) {
	// Arrange
	mockPlanRepo := new(MockPlanRepository)
	mockInvRepo := new(MockInvitacionRepository)
	service := services.NewInvitacionService(mockInvRepo, mockPlanRepo)

	idUsuario := 1
	idPlanInexistente := 999

	mockPlanRepo.On("GetByID", mock.Anything, idPlanInexistente).Return(nil, nil)

	// Act
	inv, err := service.SolicitarInvitacion(context.Background(), idUsuario, idPlanInexistente)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, inv)
	assert.Contains(t, err.Error(), "no existe")
	mockInvRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	mockPlanRepo.AssertExpectations(t)
}

// TestInvitacionService_SolicitarInvitacion_YaExiste_Idempotente verifica que si ya existe no se vuelve a crear.
func TestInvitacionService_SolicitarInvitacion_YaExiste_Idempotente(t *testing.T) {
	// Arrange
	mockPlanRepo := new(MockPlanRepository)
	mockInvRepo := new(MockInvitacionRepository)
	service := services.NewInvitacionService(mockInvRepo, mockPlanRepo)

	idUsuario := 1
	idPlan := 5
	invitacionExistente := &models.Invitacion{
		IDInvitacion: 42,
		IDPlan:       idPlan,
		IDUsuario:    idUsuario,
		Estado:       "pendiente",
	}

	mockPlanRepo.On("GetByID", mock.Anything, idPlan).Return(&models.Plan{ID: idPlan}, nil)
	mockInvRepo.On("GetByUsuarioYPlan", mock.Anything, idUsuario, idPlan).Return(invitacionExistente, nil)

	// Act
	inv, err := service.SolicitarInvitacion(context.Background(), idUsuario, idPlan)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, inv)
	assert.Equal(t, 42, inv.IDInvitacion)
	mockInvRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	mockPlanRepo.AssertExpectations(t)
	mockInvRepo.AssertExpectations(t)
}

// TestInvitacionService_SolicitarInvitacion_ParametrosInvalidos_TableDriven prueba casos bordes de IDs inválidos.
func TestInvitacionService_SolicitarInvitacion_ParametrosInvalidos_TableDriven(t *testing.T) {
	casos := []struct {
		nombre        string
		idUsuario     int
		idPlan        int
		errorEsperado string
	}{
		{
			nombre:        "usuario cero",
			idUsuario:     0,
			idPlan:        10,
			errorEsperado: "el id de usuario debe ser mayor a 0",
		},
		{
			nombre:        "usuario negativo",
			idUsuario:     -2,
			idPlan:        10,
			errorEsperado: "el id de usuario debe ser mayor a 0",
		},
		{
			nombre:        "plan cero",
			idUsuario:     1,
			idPlan:        0,
			errorEsperado: "el id de plan debe ser mayor a 0",
		},
		{
			nombre:        "plan negativo",
			idUsuario:     1,
			idPlan:        -5,
			errorEsperado: "el id de plan debe ser mayor a 0",
		},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			// Arrange
			mockPlanRepo := new(MockPlanRepository)
			mockInvRepo := new(MockInvitacionRepository)
			service := services.NewInvitacionService(mockInvRepo, mockPlanRepo)

			// Act
			inv, err := service.SolicitarInvitacion(context.Background(), tc.idUsuario, tc.idPlan)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, inv)
			assert.Contains(t, err.Error(), tc.errorEsperado)
			mockPlanRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
			mockInvRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
		})
	}
}

// TestInvitacionService_GetMisPlanes_Exitoso verifica el listado de planes con invitaciones del usuario.
func TestInvitacionService_GetMisPlanes_Exitoso(t *testing.T) {
	// Arrange
	mockPlanRepo := new(MockPlanRepository)
	mockInvRepo := new(MockInvitacionRepository)
	service := services.NewInvitacionService(mockInvRepo, mockPlanRepo)

	idUsuario := 1
	planesEsperados := []models.PlanConInvitacion{
		{ID: 10, Nombre: "Trekking", IDInvitacion: 1, Estado: "pendiente"},
		{ID: 20, Nombre: "Asado", IDInvitacion: 2, Estado: "aceptada"},
	}

	mockInvRepo.On("ListPlanesConInvitacionByUsuario", mock.Anything, idUsuario).Return(planesEsperados, nil)

	// Act
	resultado, err := service.GetMisPlanes(context.Background(), idUsuario)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, resultado, 2)
	assert.Equal(t, "Trekking", resultado[0].Nombre)
	assert.Equal(t, "aceptada", resultado[1].Estado)
	mockInvRepo.AssertExpectations(t)
}

// TestInvitacionService_GetMisPlanes_UsuarioInvalido_CasoBorde prueba el error al pasar un ID de usuario inválido.
func TestInvitacionService_GetMisPlanes_UsuarioInvalido_CasoBorde(t *testing.T) {
	// Arrange
	mockPlanRepo := new(MockPlanRepository)
	mockInvRepo := new(MockInvitacionRepository)
	service := services.NewInvitacionService(mockInvRepo, mockPlanRepo)

	// Act
	resultado, err := service.GetMisPlanes(context.Background(), -1)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultado)
	assert.Contains(t, err.Error(), "el id de usuario debe ser mayor a 0")
	mockInvRepo.AssertNotCalled(t, "ListPlanesConInvitacionByUsuario", mock.Anything, mock.Anything)
}

// TestInvitacionService_GetMisInvitaciones_Exitoso prueba la obtención de invitaciones del usuario.
func TestInvitacionService_GetMisInvitaciones_Exitoso(t *testing.T) {
	// Arrange
	mockPlanRepo := new(MockPlanRepository)
	mockInvRepo := new(MockInvitacionRepository)
	service := services.NewInvitacionService(mockInvRepo, mockPlanRepo)

	invsEsperadas := []models.Invitacion{
		{IDInvitacion: 1, IDPlan: 5, IDUsuario: 1, Estado: "pendiente"},
	}
	mockInvRepo.On("ListByUsuario", mock.Anything, 1).Return(invsEsperadas, nil)

	// Act
	invs, err := service.GetMisInvitaciones(context.Background(), 1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, invsEsperadas, invs)
	mockInvRepo.AssertExpectations(t)
}

// TestInvitacionService_GetMisInvitaciones_UsuarioInvalido prueba la validación con ID no positivo.
func TestInvitacionService_GetMisInvitaciones_UsuarioInvalido(t *testing.T) {
	// Arrange
	mockPlanRepo := new(MockPlanRepository)
	mockInvRepo := new(MockInvitacionRepository)
	service := services.NewInvitacionService(mockInvRepo, mockPlanRepo)

	// Act
	invs, err := service.GetMisInvitaciones(context.Background(), 0)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, invs)
	assert.Contains(t, err.Error(), "el id de usuario debe ser mayor a 0")
}
