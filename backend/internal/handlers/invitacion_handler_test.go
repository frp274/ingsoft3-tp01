package handlers_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/frp274/ingsoft3-tp01/backend/internal/handlers"
	"github.com/frp274/ingsoft3-tp01/backend/internal/models"
)

// MockInvitacionService implementa services.InvitacionService usando testify/mock.
type MockInvitacionService struct {
	mock.Mock
}

func (m *MockInvitacionService) SolicitarInvitacion(ctx context.Context, idUsuario int, idPlan int) (*models.Invitacion, error) {
	args := m.Called(ctx, idUsuario, idPlan)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Invitacion), args.Error(1)
}

func (m *MockInvitacionService) GetMisPlanes(ctx context.Context, idUsuario int) ([]models.PlanConInvitacion, error) {
	args := m.Called(ctx, idUsuario)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PlanConInvitacion), args.Error(1)
}

func (m *MockInvitacionService) GetMisInvitaciones(ctx context.Context, idUsuario int) ([]models.Invitacion, error) {
	args := m.Called(ctx, idUsuario)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Invitacion), args.Error(1)
}

// TestInvitacionHandler_SolicitarInvitacion_Exitoso verifica el retorno 201 Created al solicitar invitación.
func TestInvitacionHandler_SolicitarInvitacion_Exitoso(t *testing.T) {
	// Arrange
	mockService := new(MockInvitacionService)
	handler := handlers.NewInvitacionHandler(mockService)

	mockService.On("SolicitarInvitacion", mock.Anything, 1, 5).Return(&models.Invitacion{
		IDInvitacion: 1,
		IDPlan:       5,
		IDUsuario:    1,
		Estado:       "pendiente",
	}, nil)

	payload := `{"idPlan": 5}`
	req := httptest.NewRequest(http.MethodPost, "/api/invitaciones", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Act
	handler.HandleInvitaciones(rec, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "pendiente")
	assert.Contains(t, rec.Body.String(), `"idPlan":5`)
	mockService.AssertExpectations(t)
}

// TestInvitacionHandler_SolicitarInvitacion_ErrorServicio verifica HTTP 500 cuando el servicio falla.
func TestInvitacionHandler_SolicitarInvitacion_ErrorServicio(t *testing.T) {
	// Arrange
	mockService := new(MockInvitacionService)
	handler := handlers.NewInvitacionHandler(mockService)

	mockService.On("SolicitarInvitacion", mock.Anything, 1, 5).Return(nil, errors.New("fallo en BD"))

	payload := `{"idPlan": 5}`
	req := httptest.NewRequest(http.MethodPost, "/api/invitaciones", bytes.NewBufferString(payload))
	rec := httptest.NewRecorder()

	// Act
	handler.HandleInvitaciones(rec, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error al procesar solicitud")
}

// TestInvitacionHandler_SolicitarInvitacion_CasosError_TableDriven valida errores 400 Bad Request en la solicitud.
func TestInvitacionHandler_SolicitarInvitacion_CasosError_TableDriven(t *testing.T) {
	casos := []struct {
		nombrePrueba  string
		payload       string
		errorEsperado string
	}{
		{
			nombrePrueba:  "idPlan cero",
			payload:       `{"idPlan": 0}`,
			errorEsperado: "El campo 'idPlan' es obligatorio y debe ser mayor a 0",
		},
		{
			nombrePrueba:  "idPlan negativo",
			payload:       `{"idPlan": -5}`,
			errorEsperado: "El campo 'idPlan' es obligatorio y debe ser mayor a 0",
		},
		{
			nombrePrueba:  "payload JSON roto",
			payload:       `{"idPlan": `,
			errorEsperado: "Payload JSON inválido",
		},
	}

	for _, tc := range casos {
		t.Run(tc.nombrePrueba, func(t *testing.T) {
			// Arrange
			mockService := new(MockInvitacionService)
			handler := handlers.NewInvitacionHandler(mockService)

			req := httptest.NewRequest(http.MethodPost, "/api/invitaciones", bytes.NewBufferString(tc.payload))
			rec := httptest.NewRecorder()

			// Act
			handler.HandleInvitaciones(rec, req)

			// Assert
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.errorEsperado)
			mockService.AssertNotCalled(t, "SolicitarInvitacion", mock.Anything, mock.Anything, mock.Anything)
		})
	}
}

// TestInvitacionHandler_HandleInvitaciones_MetodoNoPermitido verifica rechazo a métodos no soportados.
func TestInvitacionHandler_HandleInvitaciones_MetodoNoPermitido(t *testing.T) {
	// Arrange
	mockService := new(MockInvitacionService)
	handler := handlers.NewInvitacionHandler(mockService)

	req := httptest.NewRequest(http.MethodPut, "/api/invitaciones", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.HandleInvitaciones(rec, req)

	// Assert
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	assert.Contains(t, rec.Body.String(), "Método no permitido")
}

// TestInvitacionHandler_HandleMisPlanes_Exitoso verifica el listado de planes del usuario activo con HTTP 200 OK.
func TestInvitacionHandler_HandleMisPlanes_Exitoso(t *testing.T) {
	// Arrange
	mockService := new(MockInvitacionService)
	handler := handlers.NewInvitacionHandler(mockService)

	planesConInv := []models.PlanConInvitacion{
		{ID: 5, Nombre: "Trekking Champaquí", IDInvitacion: 1, Estado: "pendiente"},
		{ID: 8, Nombre: "Fiesta Primavera", IDInvitacion: 2, Estado: "aceptada"},
	}

	mockService.On("GetMisPlanes", mock.Anything, 1).Return(planesConInv, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/mis-planes", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.HandleMisPlanes(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Trekking Champaquí")
	assert.Contains(t, rec.Body.String(), "Fiesta Primavera")
	mockService.AssertExpectations(t)
}

// TestInvitacionHandler_HandleMisPlanes_ErrorServicio valida HTTP 500 cuando el servicio falla.
func TestInvitacionHandler_HandleMisPlanes_ErrorServicio(t *testing.T) {
	// Arrange
	mockService := new(MockInvitacionService)
	handler := handlers.NewInvitacionHandler(mockService)

	mockService.On("GetMisPlanes", mock.Anything, 1).Return(nil, errors.New("error en repositorio"))

	req := httptest.NewRequest(http.MethodGet, "/api/mis-planes", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.HandleMisPlanes(rec, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error al consultar mis planes")
}

// TestInvitacionHandler_HandleMisPlanes_MetodoNoPermitido prueba que sólo se acepte GET en /api/mis-planes.
func TestInvitacionHandler_HandleMisPlanes_MetodoNoPermitido(t *testing.T) {
	// Arrange
	mockService := new(MockInvitacionService)
	handler := handlers.NewInvitacionHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/mis-planes", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.HandleMisPlanes(rec, req)

	// Assert
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	assert.Contains(t, rec.Body.String(), "Método no permitido")
}

// TestInvitacionHandler_GetMisInvitaciones_Exitoso verifica el endpoint GET /api/invitaciones.
func TestInvitacionHandler_GetMisInvitaciones_Exitoso(t *testing.T) {
	// Arrange
	mockService := new(MockInvitacionService)
	handler := handlers.NewInvitacionHandler(mockService)

	invitaciones := []models.Invitacion{
		{IDInvitacion: 1, IDPlan: 10, IDUsuario: 1, Estado: "pendiente"},
	}

	mockService.On("GetMisInvitaciones", mock.Anything, 1).Return(invitaciones, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/invitaciones", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.HandleInvitaciones(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"idInvitacion":1`)
	mockService.AssertExpectations(t)
}

// TestInvitacionHandler_GetMisInvitaciones_ErrorServicio valida HTTP 500 al fallar el servicio.
func TestInvitacionHandler_GetMisInvitaciones_ErrorServicio(t *testing.T) {
	// Arrange
	mockService := new(MockInvitacionService)
	handler := handlers.NewInvitacionHandler(mockService)

	mockService.On("GetMisInvitaciones", mock.Anything, 1).Return(nil, errors.New("fallo de base de datos"))

	req := httptest.NewRequest(http.MethodGet, "/api/invitaciones", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.HandleInvitaciones(rec, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error al obtener invitaciones")
}
