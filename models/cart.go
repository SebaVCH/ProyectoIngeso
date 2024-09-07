package models

type Carrito struct {
	CartID   string `gorm:"primaryKey;type:text" json:"cartID"`
	UserID   string `gorm:"not null;type:text" json:"userID"`
	CourseID string `gorm:"type:text" json:"courseID"`
	Quantity int    `json:"quantity"`

	User  Usuario `gorm:"foreignKey:UserID"`
	Curso Curso   `gorm:"foreignKey:CourseID"`
}
