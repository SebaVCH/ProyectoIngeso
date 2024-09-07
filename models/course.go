package models

type Curso struct {
	CourseID     string  `gorm:"primaryKey;type:text" json:"courseID"`
	InstructorID string  `gorm:"not null;type:text" json:"instructorID"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Category     string  `json:"category"`

	Instructor Usuario `gorm:"foreignKey:InstructorID"`
}
