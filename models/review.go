package models

type Reseña struct {
	ReviewID string `gorm:"primaryKey;type:text" json:"reviewID"`
	UserID   string `gorm:"not null;type:text" json:"userID"`
	CourseID string `gorm:"not null;type:text" json:"courseID"`
	Rating   int    `json:"rating"`
	Comments string `json:"comments"`

	User Usuario `gorm:"foreignKey:UserID"`
}
