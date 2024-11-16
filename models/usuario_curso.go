package models

type UsuarioCurso struct {
	ID       string `gorm:"primaryKey;column:id;type:text;default:(hex(randomblob(16)))" json:"id"`
	UserID   string `gorm:"column:user_id;type:text" json:"userID"`
	CourseID string `gorm:"column:course_id;type:text" json:"courseID"`
}

// TableName especifica el nombre de la tabla en la base de datos.
func (UsuarioCurso) TableName() string {
	return "usuario_cursos"
}
