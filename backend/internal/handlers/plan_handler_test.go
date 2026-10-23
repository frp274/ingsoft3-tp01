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

// MockPlanService implementa services.PlanService usando testify/mock.
type MockPlanService struct {
	mock.Mock
}

func (m *MockPlanService) GetPlanesPublicos(ctx context.Context) ([]models.Plan, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Plan), args.Error(1)
}

func (m *MockPlanService) GetPlanByID(ctx context.Context, id int) (*models.Plan, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Plan), args.Error(1)
}

func (m *MockPlanService) CrearPlan(ctx context.Context, plan *models.Plan) error {
	args := m.Called(ctx, plan)
	return args.Error(0)
}

// TestPlanHandler_CrearPlan_Exitoso verifica el retorno 201 Created cuando el payload es válido.
func TestPlanHandler_CrearPlan_Exitoso(t *testing.T) {
	// Arrange
	mockService := new(MockPlanService)
	handler := handlers.NewPlanHandler(mockService)

	mockService.On("CrearPlan", mock.Anything, mock.MatchedBy(func(p *models.Plan) bool {
		return p.Nombre == "Paseo en Bici" && p.Capacidad == 10
	})).Return(nil)

	payload := `{
		"nombre": "Paseo en Bici",
		"fechaInicio": "2026-11-10T10:00:00Z",
		"fechaFin": "2026-11-10T13:00:00Z",
		"visibilidad": "publico",
		"idEtiqueta": 1,
		"ubicacion": "Parque Sarmiento",
		"capacidad": 10,
		"descripcion": "Llevar agua y casco"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/planes", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// Act
	handler.HandlePlanes(rec, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "Paseo en Bici")
	mockService.AssertExpectations(t)
}

// TestPlanHandler_CrearPlan_ErrorServicio valida 500 cuando el servicio falla.
func TestPlanHandler_CrearPlan_ErrorServicio(t *testing.T) {
	// Arrange
	mockService := new(MockPlanService)
	handler := handlers.NewPlanHandler(mockService)

	mockService.On("CrearPlan", mock.Anything, mock.Anything).Return(errors.New("fallo en base de datos"))

	payload := `{
		"nombre": "Paseo en Bici",
		"fechaInicio": "2026-11-10T10:00:00Z",
		"fechaFin": "2026-11-10T13:00:00Z",
		"idEtiqueta": 1,
		"ubicacion": "Parque",
		"capacidad": 10
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/planes", bytes.NewBufferString(payload))
	rec := httptest.NewRecorder()

	// Act
	handler.HandlePlanes(rec, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error al guardar el plan")
}

// TestPlanHandler_CrearPlan_ValidacionesPayload_TableDriven valida errores 400 Bad Request ante payloads inválidos.
func TestPlanHandler_CrearPlan_ValidacionesPayload_TableDriven(t *testing.T) {
	casos := []struct {
		nombrePrueba  string
		payload       string
		errorEsperado string
	}{
		{
			nombrePrueba: "nombre vacio",
			payload: `{
				"nombre": "   ",
				"fechaInicio": "2026-11-10T10:00:00Z",
				"fechaFin": "2026-11-10T13:00:00Z",
				"visibilidad": "publico",
				"idEtiqueta": 1,
				"ubicacion": "Parque Sarmiento",
				"capacidad": 10
			}`,
			errorEsperado: "El campo 'nombre' es obligatorio",
		},
		{
			nombrePrueba: "ubicacion vacia",
			payload: `{
				"nombre": "Caminata",
				"fechaInicio": "2026-11-10T10:00:00Z",
				"fechaFin": "2026-11-10T13:00:00Z",
				"visibilidad": "publico",
				"idEtiqueta": 1,
				"ubicacion": "   ",
				"capacidad": 10
			}`,
			errorEsperado: "El campo 'ubicacion' es obligatorio",
		},
		{
			nombrePrueba: "capacidad menor o igual a cero",
			payload: `{
				"nombre": "Caminata",
				"fechaInicio": "2026-11-10T10:00:00Z",
				"fechaFin": "2026-11-10T13:00:00Z",
				"visibilidad": "publico",
				"idEtiqueta": 1,
				"ubicacion": "Parque",
				"capacidad": 0
			}`,
			errorEsperado: "El campo 'capacidad' debe ser un número entero mayor a 0",
		},
		{
			nombrePrueba: "etiqueta invalida",
			payload: `{
				"nombre": "Caminata",
				"fechaInicio": "2026-11-10T10:00:00Z",
				"fechaFin": "2026-11-10T13:00:00Z",
				"visibilidad": "publico",
				"idEtiqueta": 0,
				"ubicacion": "Parque",
				"capacidad": 10
			}`,
			errorEsperado: "El campo 'idEtiqueta' es obligatorio",
		},
		{
			nombrePrueba: "visibilidad invalida",
			payload: `{
				"nombre": "Caminata",
				"fechaInicio": "2026-11-10T10:00:00Z",
				"fechaFin": "2026-11-10T13:00:00Z",
				"visibilidad": "invalida",
				"idEtiqueta": 1,
				"ubicacion": "Parque",
				"capacidad": 10
			}`,
			errorEsperado: "El campo 'visibilidad' debe ser 'publico' o 'privado'",
		},
		{
			nombrePrueba: "fechaInicio vacia",
			payload: `{
				"nombre": "Caminata",
				"fechaInicio": "",
				"fechaFin": "2026-11-10T13:00:00Z",
				"idEtiqueta": 1,
				"ubicacion": "Parque",
				"capacidad": 10
			}`,
			errorEsperado: "El campo 'fechaInicio' es obligatorio",
		},
		{
			nombrePrueba: "fechaFin vacia",
			payload: `{
				"nombre": "Caminata",
				"fechaInicio": "2026-11-10T10:00:00Z",
				"fechaFin": "",
				"idEtiqueta": 1,
				"ubicacion": "Parque",
				"capacidad": 10
			}`,
			errorEsperado: "El campo 'fechaFin' es obligatorio",
		},
		{
			nombrePrueba: "fechaInicio formato invalido",
			payload: `{
				"nombre": "Caminata",
				"fechaInicio": "fecha-rara",
				"fechaFin": "2026-11-10T13:00:00Z",
				"idEtiqueta": 1,
				"ubicacion": "Parque",
				"capacidad": 10
			}`,
			errorEsperado: "Formato inválido para 'fechaInicio'",
		},
		{
			nombrePrueba: "fechaFin formato invalido",
			payload: `{
				"nombre": "Caminata",
				"fechaInicio": "2026-11-10T10:00:00Z",
				"fechaFin": "fecha-rara",
				"idEtiqueta": 1,
				"ubicacion": "Parque",
				"capacidad": 10
			}`,
			errorEsperado: "Formato inválido para 'fechaFin'",
		},
		{
			nombrePrueba: "fecha fin anterior a fecha inicio",
			payload: `{
				"nombre": "Caminata",
				"fechaInicio": "2026-11-10T18:00:00Z",
				"fechaFin": "2026-11-10T10:00:00Z",
				"visibilidad": "publico",
				"idEtiqueta": 1,
				"ubicacion": "Parque",
				"capacidad": 10
			}`,
			errorEsperado: "La 'fechaFin' debe ser posterior a la 'fechaInicio'",
		},
		{
			nombrePrueba:  "json malformado",
			payload:       `{"nombre": "Caminata", invalid_json}`,
			errorEsperado: "Payload JSON inválido",
		},
	}

	for _, tc := range casos {
		t.Run(tc.nombrePrueba, func(t *testing.T) {
			// Arrange
			mockService := new(MockPlanService)
			handler := handlers.NewPlanHandler(mockService)

			req := httptest.NewRequest(http.MethodPost, "/api/planes", bytes.NewBufferString(tc.payload))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Act
			handler.HandlePlanes(rec, req)

			// Assert
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.errorEsperado)
			mockService.AssertNotCalled(t, "CrearPlan", mock.Anything, mock.Anything)
		})
	}
}

// TestPlanHandler_GetPlanesPublicos_Exitoso verifica el retorno 200 OK y el listado de planes en JSON.
func TestPlanHandler_GetPlanesPublicos_Exitoso(t *testing.T) {
	// Arrange
	mockService := new(MockPlanService)
	handler := handlers.NewPlanHandler(mockService)

	planes := []models.Plan{
		{ID: 1, Nombre: "Trekking", Visibilidad: "publico"},
		{ID: 2, Nombre: "Kayak", Visibilidad: "publico"},
	}

	mockService.On("GetPlanesPublicos", mock.Anything).Return(planes, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/planes", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.HandlePlanes(rec, req)

	// Assert
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Trekking")
	assert.Contains(t, rec.Body.String(), "Kayak")
	mockService.AssertExpectations(t)
}

// TestPlanHandler_GetPlanesPublicos_ErrorServicio valida HTTP 500 cuando el servicio falla.
func TestPlanHandler_GetPlanesPublicos_ErrorServicio(t *testing.T) {
	// Arrange
	mockService := new(MockPlanService)
	handler := handlers.NewPlanHandler(mockService)

	mockService.On("GetPlanesPublicos", mock.Anything).Return(nil, errors.New("fallo en servicio"))

	req := httptest.NewRequest(http.MethodGet, "/api/planes", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.HandlePlanes(rec, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error al consultar planes")
}

// TestPlanHandler_MetodoNoPermitido_CasoBorde prueba el caso de un método HTTP no soportado.
func TestPlanHandler_MetodoNoPermitido_CasoBorde(t *testing.T) {
	// Arrange
	mockService := new(MockPlanService)
	handler := handlers.NewPlanHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/api/planes", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.HandlePlanes(rec, req)

	// Assert
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	assert.Contains(t, rec.Body.String(), "Método no permitido")
}
