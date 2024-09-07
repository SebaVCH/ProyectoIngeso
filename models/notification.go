package models

type Notificación struct {
	NotificationID string `gorm:"primaryKey;type:text" json:"notificationID"`
	UserID         string `gorm:"not null;type:text" json:"userID"`
	Message        string `json:"message"`
	Status         string `json:"status"`
	CreatedAt      string `json:"createdAt"` // Puedes usar un tipo de fecha si lo prefieres

	User Usuario `gorm:"foreignKey:UserID"`
}
