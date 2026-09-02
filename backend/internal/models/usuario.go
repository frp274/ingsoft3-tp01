package models

// Usuario representa la entidad de un usuario en el sistema.
type Usuario struct {
	ID               int     `json:"id" db:"id"`
	Email            string  `json:"email" db:"email"`
	Contrasena       string  `json:"contrasena,omitempty" db:"contrasena"`
	Nombre           string  `json:"nombre" db:"nombre"`
	Apellido         string  `json:"apellido" db:"apellido"`
	Ciudad           string  `json:"ciudad" db:"ciudad"`
	Edad             int     `json:"edad" db:"edad"`
	EtiquetaFavorita *string `json:"etiquetaFavorita" db:"etiquetaFavorita"`
}