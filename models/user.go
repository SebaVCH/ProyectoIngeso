package models

import (
	"gorm.io/gorm"
)

// Estructura que representa un usuario en la base de datos
type Usuario struct {
	gorm.Model
	NombreYapellido   string `json:"nameLastName"`
	NombreUsuario     string `gorm:"uniqueIndex" json:"username"`
	CorreoElectronico string `gorm:"uniqueIndex" json:"mail"`
	Contrasena        string `json:"password"`
}
