package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.53

import (
	"ProyectoIngeso/graph/model"
	"ProyectoIngeso/models"
	"ProyectoIngeso/utils"
	"context"
	"errors"
	"fmt"
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
		UserID:       usuario.UserID,
		NameLastName: usuario.NameLastName,
		Username:     usuario.Username,
		Email:        usuario.Email,
		Password:     usuario.Password,
		Role:         usuario.Role,
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

// UpdateUsername maneja la mutación para actualizar el nombre de usuario.
func (r *mutationResolver) ActualizarUsername(ctx context.Context, username string, newUsername string) (*model.Usuario, error) {
	var usuario models.Usuario

	// Buscar el usuario por el nombre de usuario actual
	if err := r.DB.Where("username = ?", username).First(&usuario).Error; err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	// Actualizar el nombre de usuario
	usuario.Username = newUsername

	// Guardar los cambios
	if err := r.DB.Save(&usuario).Error; err != nil {
		return nil, errors.New("no se pudo actualizar el nombre de usuario")
	}

	// Retornar el usuario actualizado
	return &model.Usuario{
		UserID:       usuario.UserID,
		NameLastName: usuario.NameLastName,
		Username:     usuario.Username,
		Email:        usuario.Email,
		Password:     usuario.Password,
		Role:         usuario.Role,
	}, nil
}

// UpdatePassword maneja la mutación para actualizar la contraseña.
func (r *mutationResolver) ActualizarPassword(ctx context.Context, username string, oldPassword string, newPassword string) (*string, error) {
	var usuario models.Usuario

	// Buscar el usuario por el nombre de usuario
	if err := r.DB.Where("username = ?", username).First(&usuario).Error; err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	// Verificar que la contraseña actual sea correcta
	if !utils.VerificarHashContrasena(oldPassword, usuario.Password) {
		return nil, errors.New("la contraseña actual es incorrecta")
	}

	// Cifrar la nueva contraseña
	newHashedPassword, err := utils.HashContrasena(newPassword)
	if err != nil {
		return nil, errors.New("error al cifrar la nueva contraseña")
	}

	// Actualizar la contraseña
	usuario.Password = newHashedPassword

	// Guardar los cambios
	if err := r.DB.Save(&usuario).Error; err != nil {
		return nil, errors.New("no se pudo actualizar la contraseña")
	}

	successMsg := "Contraseña actualizada exitosamente"
	return &successMsg, nil
}

// AddToCart is the resolver for the addToCart field.
func (r *mutationResolver) AddToCart(ctx context.Context, username string, courseID string, quantity int) (*model.Carrito, error) {
	return r.Resolver.AddToCart(ctx, username, courseID, quantity)
}

// RemoveFromCart is the resolver for the removeFromCart field.
func (r *mutationResolver) RemoveFromCart(ctx context.Context, username string, courseID string) (*bool, error) {
	success, err := r.Resolver.RemoveFromCart(ctx, username, courseID)
	if err != nil {
		return nil, err
	}
	return success, nil
}

// GetUsuario maneja la consulta para obtener un usuario por su ID.
func (r *queryResolver) GetUsuario(ctx context.Context, id string) (*model.Usuario, error) {
	var usuario models.Usuario
	if err := r.DB.First(&usuario, "user_id = ?", id).Error; err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	// Convertir el modelo de base de datos a modelo GraphQL
	return &model.Usuario{
		UserID:       usuario.UserID,
		NameLastName: usuario.NameLastName,
		Username:     usuario.Username,
		Email:        usuario.Email,
		Password:     usuario.Password,
		Role:         usuario.Role,
	}, nil
}

// UserByUsername is the resolver for the userByUsername field.
func (r *queryResolver) UserByUsername(ctx context.Context, username string) (*model.Usuario, error) {
	var usuario model.Usuario

	// Buscar el usuario por nombre de usuario en la base de datos
	if err := r.DB.Where("username = ?", username).First(&usuario).Error; err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %v", err)
	}

	return &usuario, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
