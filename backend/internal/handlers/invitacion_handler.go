package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/frp274/ingsoft3-tp01/backend/internal/models"
	"github.com/frp274/ingsoft3-tp01/backend/internal/services"
)

// InvitacionHandler maneja las peticiones HTTP para Invitaciones y Mis Planes.
type InvitacionHandler struct {
	invitacionService services.InvitacionService
}

// NewInvitacionHandler inicializa el controlador.
func NewInvitacionHandler(invitacionService services.InvitacionService) *InvitacionHandler {
	return &InvitacionHandler{
		invitacionService: invitacionService,
	}
}

type SolicitarInvitacionDTO struct {
	IDPlan int `json:"idPlan"`
}

// HandleInvitaciones maneja POST /api/invitaciones y GET /api/invitaciones
func (h *InvitacionHandler) HandleInvitaciones(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.SolicitarInvitacion(w, r)
	case http.MethodGet:
		h.GetMisInvitaciones(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método no permitido"})
	}
}

// SolicitarInvitacion maneja POST /api/invitaciones (US-5).
func (h *InvitacionHandler) SolicitarInvitacion(w http.ResponseWriter, r *http.Request) {
	var dto SolicitarInvitacionDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Payload JSON inválido: " + err.Error()})
		return
	}

	if dto.IDPlan <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "El campo 'idPlan' es obligatorio y debe ser mayor a 0"})
		return
	}

	// Usuario activo temporal (id = 1) hasta implementar JWT
	idUsuarioActivo := 1

	invitacion, err := h.invitacionService.SolicitarInvitacion(r.Context(), idUsuarioActivo, dto.IDPlan)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al procesar solicitud: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invitacion)
}

// GetMisInvitaciones devuelve las invitaciones del usuario activo.
func (h *InvitacionHandler) GetMisInvitaciones(w http.ResponseWriter, r *http.Request) {
	idUsuarioActivo := 1
	invitaciones, err := h.invitacionService.GetMisInvitaciones(r.Context(), idUsuarioActivo)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al obtener invitaciones: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(invitaciones)
}

// HandleMisPlanes maneja GET /api/mis-planes (US-5).
func (h *InvitacionHandler) HandleMisPlanes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método no permitido"})
		return
	}

	idUsuarioActivo := 1
	misPlanes, err := h.invitacionService.GetMisPlanes(r.Context(), idUsuarioActivo)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al consultar mis planes: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if misPlanes == nil {
		misPlanes = []models.PlanConInvitacion{}
	}
	json.NewEncoder(w).Encode(misPlanes)
}