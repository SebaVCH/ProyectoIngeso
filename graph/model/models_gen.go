// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Carrito struct {
	CartID   string `json:"cartID"`
	UserID   string `json:"userID"`
	CourseID string `json:"courseID"`
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
