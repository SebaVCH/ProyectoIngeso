package models

type Usuario struct {
	UserID       string `gorm:"primaryKey;type:text" json:"userID"`
	NameLastName string `json:"nameLastName"`
	Username     string `gorm:"uniqueIndex" json:"username"`
	Email        string `gorm:"uniqueIndex" json:"email"`
	Password     string `json:"password"`
	Role         string `json:"role"`
}
