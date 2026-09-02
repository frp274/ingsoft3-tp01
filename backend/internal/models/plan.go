package models

import "time"

// Plan representa un evento o juntada organizada por un usuario.
type Plan struct {
	ID            int       `json:"id" db:"id"`
	Nombre        string    `json:"nombre" db:"nombre"`
	IDOrganizador int       `json:"idOrganizador" db:"idOrganizador"`
	FechaInicio   time.Time `json:"fechaInicio" db:"fechaInicio"`
	FechaFin      time.Time `json:"fechaFin" db:"fechaFin"`
	Visibilidad   string    `json:"visibilidad" db:"visibilidad"` // "publico" o "privado"
	Descripcion   *string   `json:"descripcion" db:"descripcion"`
	IDEtiqueta    int       `json:"idEtiqueta" db:"idEtiqueta"`
	Ubicacion     string    `json:"ubicacion" db:"ubicacion"`
	Capacidad     int       `json:"capacidad" db:"capacidad"`
}