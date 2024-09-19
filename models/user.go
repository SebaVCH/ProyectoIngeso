package models

type Usuario struct {
	UserID       string `gorm:"primaryKey;column:user_id;type:text" json:"userID"`
	NameLastName string `gorm:"column:name_last_name" json:"nameLastName"`
	Username     string `gorm:"uniqueIndex;column:username" json:"username"`
	Email        string `gorm:"uniqueIndex;column:email" json:"email"`
	Password     string `gorm:"column:password" json:"password"`
	Role         string `gorm:"column:role" json:"role"`
}

// Especificar el nombre de la tabla
func (Usuario) TableName() string {
	return "usuarios"
}
