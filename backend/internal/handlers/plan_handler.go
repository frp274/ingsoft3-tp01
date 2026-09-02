package handlers

import (
	"encoding/json"
	"net/http"

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

// GetPlanesPublicos maneja el endpoint GET /api/planes.
// Devuelve una lista de todos los planes con visibilidad 'publico'.
func (h *PlanHandler) GetPlanesPublicos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método no permitido"})
		return
	}

	planes, err := h.planService.GetPlanesPublicos(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar planes: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(planes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}