package models

// Invitacion representa la tabla pasarela M:N entre Usuario y Plan.
type Invitacion struct {
	IDInvitacion int    `json:"idInvitacion" db:"idInvitacion"`
	IDPlan       int    `json:"idPlan" db:"idPlan"`
	IDUsuario    int    `json:"idUsuario" db:"idUsuario"`
	Estado       string `json:"estado" db:"estado"` // "pendiente", "aceptada", "rechazada"
}