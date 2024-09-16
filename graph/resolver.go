package graph

import (
	"ProyectoIngeso/models"
	"ProyectoIngeso/utils"
	"context"
	"errors"
	"gorm.io/gorm"
)

type Resolver struct {
	DB *gorm.DB
}

// RegistrarUsuario - maneja el registro de usuario
func (r *Resolver) RegistrarUsuario(ctx context.Context, input struct {
	NombreYapellido   string
	NombreUsuario     string
	CorreoElectronico string
	Contrasena        string
}) (*models.Usuario, error) {
	// Cifrar la contraseña
	hash, err := utils.HashContrasena(input.Contrasena)
	if err != nil {
		return nil, errors.New("error al cifrar la contraseña")
	}

	usuario := models.Usuario{
		UserID:       generateUniqueID(), // Implementa esta función para generar un ID único
		NameLastName: input.NombreYapellido,
		Username:     input.NombreUsuario,
		Email:        input.CorreoElectronico,
		Password:     hash,
		Role:         "user", // Asigna un rol por defecto si es necesario
	}

	if err := r.DB.Create(&usuario).Error; err != nil {
		return nil, errors.New("error al crear el usuario")
	}

	return &usuario, nil
}

// IniciarSesion - maneja el inicio de sesión del usuario
func (r *Resolver) IniciarSesion(ctx context.Context, input struct {
	Identificador string
	Contrasena    string
}) (string, error) {
	var usuario models.Usuario
	if err := r.DB.Where("email = ? OR username = ?", input.Identificador, input.Identificador).First(&usuario).Error; err != nil {
		return "", errors.New("usuario no encontrado")
	}

	if !utils.VerificarHashContrasena(input.Contrasena, usuario.Password) {
		return "", errors.New("contraseña inválida")
	}

	return "Inicio de sesión exitoso", nil
}

// UpdateUsername - maneja la actualización del nombre de usuario
func (r *Resolver) UpdateUsername(ctx context.Context, username string, newUsername string) (*models.Usuario, error) {
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
	return &usuario, nil
}

// UpdatePassword - maneja la actualización de la contraseña
func (r *Resolver) UpdatePassword(ctx context.Context, username string, oldPassword string, newPassword string) (string, error) {
	var usuario models.Usuario

	// Buscar el usuario por el nombre de usuario
	if err := r.DB.Where("username = ?", username).First(&usuario).Error; err != nil {
		return "", errors.New("usuario no encontrado")
	}

	// Verificar que la contraseña actual sea correcta
	if !utils.VerificarHashContrasena(oldPassword, usuario.Password) {
		return "", errors.New("la contraseña actual es incorrecta")
	}

	// Cifrar la nueva contraseña
	newHashedPassword, err := utils.HashContrasena(newPassword)
	if err != nil {
		return "", errors.New("error al cifrar la nueva contraseña")
	}

	// Actualizar la contraseña
	usuario.Password = newHashedPassword

	// Guardar los cambios
	if err := r.DB.Save(&usuario).Error; err != nil {
		return "", errors.New("no se pudo actualizar la contraseña")
	}

	return "Contraseña actualizada exitosamente", nil
}
