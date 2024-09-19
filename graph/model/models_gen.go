// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Curso struct {
	CourseID     string  `json:"courseID"`
	InstructorID string  `json:"instructorID"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Category     string  `json:"category"`
}

type Mutation struct {
}

type Query struct {
}

type Usuario struct {
	UserID       string `json:"userID"`
	NameLastName string `json:"nameLastName"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Role         string `json:"role"`
}
