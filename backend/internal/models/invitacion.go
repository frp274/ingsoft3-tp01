package models

import "time"

// Invitacion representa la tabla pasarela M:N entre Usuario y Plan.
type Invitacion struct {
	IDInvitacion int    `json:"idInvitacion" db:"idInvitacion"`
	IDPlan       int    `json:"idPlan" db:"idPlan"`
	IDUsuario    int    `json:"idUsuario" db:"idUsuario"`
	Estado       string `json:"estado" db:"estado"` // "pendiente", "aceptada", "rechazada"
}

// PlanConInvitacion representa un Plan con la información de la invitación del usuario (JOIN).
type PlanConInvitacion struct {
	ID            int       `json:"id" db:"id"`
	Nombre        string    `json:"nombre" db:"nombre"`
	IDOrganizador int       `json:"idOrganizador" db:"idOrganizador"`
	FechaInicio   time.Time `json:"fechaInicio" db:"fechaInicio"`
	FechaFin      time.Time `json:"fechaFin" db:"fechaFin"`
	Visibilidad   string    `json:"visibilidad" db:"visibilidad"`
	Descripcion   *string   `json:"descripcion" db:"descripcion"`
	IDEtiqueta    int       `json:"idEtiqueta" db:"idEtiqueta"`
	Ubicacion     string    `json:"ubicacion" db:"ubicacion"`
	Capacidad     int       `json:"capacidad" db:"capacidad"`
	IDInvitacion  int       `json:"idInvitacion" db:"idInvitacion"`
	Estado        string    `json:"estado" db:"estado"`
}