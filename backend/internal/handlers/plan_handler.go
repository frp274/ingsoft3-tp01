package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/frp274/ingsoft3-tp01/backend/internal/models"
	"github.com/frp274/ingsoft3-tp01/backend/internal/services"
)

// PlanHandler maneja las peticiones HTTP relacionadas con Planes.
type PlanHandler struct {
	planService services.PlanService
}

// NewPlanHandler inicializa el controlador de planes.
func NewPlanHandler(planService services.PlanService) *PlanHandler {
	return &PlanHandler{
		planService: planService,
	}
}

// CreatePlanDTO representa el payload JSON recibido para crear un plan.
type CreatePlanDTO struct {
	Nombre      string  `json:"nombre"`
	FechaInicio string  `json:"fechaInicio"`
	FechaFin    string  `json:"fechaFin"`
	Visibilidad string  `json:"visibilidad"`
	Descripcion *string `json:"descripcion"` // Opcional
	IDEtiqueta  int     `json:"idEtiqueta"`
	Ubicacion   string  `json:"ubicacion"`
	Capacidad   int     `json:"capacidad"`
}

// HandlePlanes enruta GET /api/planes y POST /api/planes
func (h *PlanHandler) HandlePlanes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetPlanesPublicos(w, r)
	case http.MethodPost:
		h.CrearPlan(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método no permitido"})
	}
}

// GetPlanesPublicos maneja el endpoint GET /api/planes.
func (h *PlanHandler) GetPlanesPublicos(w http.ResponseWriter, r *http.Request) {
	planes, err := h.planService.GetPlanesPublicos(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar planes: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(planes)
}

// CrearPlan maneja el endpoint POST /api/planes (US-4).
func (h *PlanHandler) CrearPlan(w http.ResponseWriter, r *http.Request) {
	var dto CreatePlanDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Payload JSON inválido: " + err.Error()})
		return
	}

	// 1. Validaciones de campos obligatorios
	dto.Nombre = strings.TrimSpace(dto.Nombre)
	if dto.Nombre == "" {
		h.sendError(w, http.StatusBadRequest, "El campo 'nombre' es obligatorio")
		return
	}

	dto.Ubicacion = strings.TrimSpace(dto.Ubicacion)
	if dto.Ubicacion == "" {
		h.sendError(w, http.StatusBadRequest, "El campo 'ubicacion' es obligatorio")
		return
	}

	if dto.Capacidad <= 0 {
		h.sendError(w, http.StatusBadRequest, "El campo 'capacidad' debe ser un número entero mayor a 0")
		return
	}

	if dto.IDEtiqueta <= 0 {
		h.sendError(w, http.StatusBadRequest, "El campo 'idEtiqueta' es obligatorio y debe ser válido")
		return
	}

	if dto.Visibilidad == "" {
		dto.Visibilidad = "publico"
	}
	if dto.Visibilidad != "publico" && dto.Visibilidad != "privado" {
		h.sendError(w, http.StatusBadRequest, "El campo 'visibilidad' debe ser 'publico' o 'privado'")
		return
	}

	// 2. Validación y parseo de fechas
	if dto.FechaInicio == "" {
		h.sendError(w, http.StatusBadRequest, "El campo 'fechaInicio' es obligatorio")
		return
	}
	if dto.FechaFin == "" {
		h.sendError(w, http.StatusBadRequest, "El campo 'fechaFin' es obligatorio")
		return
	}

	fechaInicio, err := parseDateTime(dto.FechaInicio)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Formato inválido para 'fechaInicio': "+err.Error())
		return
	}

	fechaFin, err := parseDateTime(dto.FechaFin)
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "Formato inválido para 'fechaFin': "+err.Error())
		return
	}

	if !fechaFin.After(fechaInicio) {
		h.sendError(w, http.StatusBadRequest, "La 'fechaFin' debe ser posterior a la 'fechaInicio'")
		return
	}

	// 3. Limpieza de campo opcional descripción
	var descripcion *string
	if dto.Descripcion != nil && strings.TrimSpace(*dto.Descripcion) != "" {
		descTrim := strings.TrimSpace(*dto.Descripcion)
		descripcion = &descTrim
	}

	// 4. Asignar temporalmente idOrganizador = 1 (usuario semilla) hasta incorporar JWT
	idOrganizadorTemporal := 1

	nuevoPlan := &models.Plan{
		Nombre:        dto.Nombre,
		IDOrganizador: idOrganizadorTemporal,
		FechaInicio:   fechaInicio,
		FechaFin:      fechaFin,
		Visibilidad:   dto.Visibilidad,
		Descripcion:   descripcion,
		IDEtiqueta:    dto.IDEtiqueta,
		Ubicacion:     dto.Ubicacion,
		Capacidad:     dto.Capacidad,
	}

	if err := h.planService.CrearPlan(r.Context(), nuevoPlan); err != nil {
		h.sendError(w, http.StatusInternalServerError, "Error al guardar el plan: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevoPlan)
}

func (h *PlanHandler) sendError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// parseDateTime intenta parsear fechas en formatos comunes (RFC3339, datetime-local de HTML, etc.)
func parseDateTime(val string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, val); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("use formato ISO 8601 o AAAA-MM-DDTHH:MM")
}