package models

type Carrito struct {
	CartID   string `json:"cartID" gorm:"primaryKey"`
	UserID   string `json:"userID"`
	CourseID string `json:"courseID"`
}

// TableName especifica el nombre de la tabla en la base de datos.
func (Carrito) TableName() string {
	return "carritos" // Aquí defines el nombre de la tabla en la base de datos.
}
