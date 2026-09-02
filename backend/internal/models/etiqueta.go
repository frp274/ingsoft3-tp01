package models

// Etiqueta representa una categoría o tag para clasificar planes.
type Etiqueta struct {
	ID     int    `json:"id" db:"id"`
	Nombre string `json:"nombre" db:"nombre"`
}