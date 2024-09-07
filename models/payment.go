package models

type Pago struct {
	PaymentID     string  `gorm:"primaryKey;type:text" json:"paymentID"`
	UserID        string  `gorm:"not null;type:text" json:"userID"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	PaymentMethod string  `json:"paymentMethod"`
	PaymentDate   string  `json:"paymentDate"` // Puedes usar un tipo de fecha si lo prefieres

	User Usuario `gorm:"foreignKey:UserID"`
}
