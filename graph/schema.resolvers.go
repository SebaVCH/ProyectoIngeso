package graph

import (
	"ProyectoIngeso/graph/model" // Importa el paquete generado por gqlgen
	"ProyectoIngeso/models"      // Importa el paquete models
	"ProyectoIngeso/utils"       // Importa el paquete utils
	"context"
	"errors"
	"github.com/google/uuid"
)

// RegisterUsuario maneja la mutación para registrar un usuario.
func (r *mutationResolver) RegisterUsuario(ctx context.Context, nameLastName string, username string, email string, password string) (*model.Usuario, error) {
	hash, err := utils.HashContrasena(password)
	if err != nil {
		return nil, errors.New("error al cifrar la contraseña")
	}

	// Crear el modelo de usuario en la base de datos
	usuario := &models.Usuario{
		UserID:       generateUniqueID(), // Implementa esta función para generar un ID único
		NameLastName: nameLastName,
		Username:     username,
		Email:        email,
		Password:     hash,
		Role:         "user", // Asigna un rol por defecto si es necesario
	}

	if err := r.DB.Create(usuario).Error; err != nil {
		return nil, errors.New("no se pudo registrar el usuario")
	}

	// Convertir el modelo de base de datos a modelo GraphQL
	return &model.Usuario{
		ID:           usuario.UserID,
		NameLastName: &usuario.NameLastName,
		Username:     &usuario.Username,
		Email:        &usuario.Email,
		Role:         &usuario.Role,
	}, nil
}

// LoginUsuario maneja la mutación para iniciar sesión.
func (r *mutationResolver) LoginUsuario(ctx context.Context, identificador string, password string) (*string, error) {
	var usuario models.Usuario
	if err := r.DB.Where("email = ? OR username = ?", identificador, identificador).First(&usuario).Error; err != nil {
		msg := "usuario no encontrado"
		return &msg, nil
	}

	if !utils.VerificarHashContrasena(password, usuario.Password) {
		msg := "contraseña inválida"
		return &msg, nil
	}

	successMsg := "Inicio de sesión exitoso"
	return &successMsg, nil
}

// GetUsuario maneja la consulta para obtener un usuario por su ID.
func (r *queryResolver) GetUsuario(ctx context.Context, id string) (*model.Usuario, error) {
	var usuario models.Usuario
	if err := r.DB.First(&usuario, "user_id = ?", id).Error; err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	// Convertir el modelo de base de datos a modelo GraphQL
	return &model.Usuario{
		ID:           usuario.UserID,
		NameLastName: &usuario.NameLastName,
		Username:     &usuario.Username,
		Email:        &usuario.Email,
		Role:         &usuario.Role,
	}, nil
}

// Mutation devuelve la implementación de MutationResolver.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query devuelve la implementación de QueryResolver.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// Genera un ID único (puedes ajustar esto según tus necesidades)
func generateUniqueID() string {
	return uuid.NewString()
}
