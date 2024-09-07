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
