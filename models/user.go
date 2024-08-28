package models

import (
	"gorm.io/gorm"
)

// Estructura que representa un usuario en la base de datos
type Usuario struct {
	gorm.Model
	CorreoElectronico string `gorm:"uniqueIndex"`
	ContrasenaHash    string
}
