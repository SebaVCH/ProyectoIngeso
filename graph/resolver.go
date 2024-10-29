package graph

import (
	"ProyectoIngeso/graph/model"
	"ProyectoIngeso/models"
	"ProyectoIngeso/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"time"
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
	// 1. Cifrar la contraseña
	hash, err := utils.HashContrasena(input.Contrasena)
	if err != nil {
		return nil, errors.New("error al cifrar la contraseña")
	}

	// 2. Crear el usuario
	usuario := models.Usuario{
		UserID:       generateUniqueID(), // Implementa esta función para generar un ID único
		NameLastName: input.NombreYapellido,
		Username:     input.NombreUsuario,
		Email:        input.CorreoElectronico,
		Password:     hash,
		Role:         "user", // Rol por defecto
	}

	// Guardar el usuario en la base de datos
	if err := r.DB.Create(&usuario).Error; err != nil {
		return nil, errors.New("error al crear el usuario")
	}

	/*// 3. Crear el carrito asociado al usuario recién creado
	carrito := models.Carrito{
		CartID:   generateUniqueID(), // Genera un ID único para el carrito
		UserID:   usuario.UserID,     // Asocia el carrito con el usuario
		CourseID: "",                 // Inicialmente sin curso
	}

	// Guardar el carrito en la base de datos
	if err := r.DB.Create(&carrito).Error; err != nil {
		return nil, errors.New("error al crear el carrito del usuario")
	}*/

	// 4. Retornar el usuario creado
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

// DeleteUserByUsername - elimina un usuario por su nombre de usuario
func (r *Resolver) DeleteUserByUsername(ctx context.Context, username string) (string, error) {
	var usuario models.Usuario

	// Buscar el usuario por el nombre de usuario
	if err := r.DB.Where("username = ?", username).First(&usuario).Error; err != nil {
		return "", errors.New("usuario no encontrado")
	}

	// Eliminar el usuario de la base de datos
	if err := r.DB.Delete(&usuario).Error; err != nil {
		return "", errors.New("no se pudo eliminar el usuario")
	}

	return "Usuario eliminado exitosamente", nil
}

// AddToCart agrega un curso al carrito del usuario.
func (r *Resolver) AddToCart(ctx context.Context, username string, courseID string) (*model.Carrito, error) {
	// Verificar si el usuario existe y obtener el userID.
	userID, err := r.checkUserExists(username)
	if err != nil {
		return nil, fmt.Errorf("error al verificar el usuario: %v", err)
	}

	// Verificar si el curso existe en el servicio de cursos.
	courseExists, err := r.checkCourseExists(courseID)
	if err != nil {
		return nil, fmt.Errorf("error al verificar el curso: %v", err)
	}
	if !courseExists {
		return nil, fmt.Errorf("curso con ID %s no encontrado", courseID)
	}

	// Crear un nuevo elemento en el carrito.
	cartItem := &model.Carrito{
		CartID:   uuid.New().String(),
		UserID:   userID, // Asegúrate de que esta variable no esté vacía.
		CourseID: courseID,
	}

	// Verificar si el campo UserID no está vacío.
	if cartItem.UserID == "" {
		return nil, fmt.Errorf("userID no puede estar vacío")
	}

	// Guardar el elemento en la base de datos.
	if err := r.DB.Create(cartItem).Error; err != nil {
		return nil, err
	}

	return cartItem, nil
}

// DeleteCartByID elimina un carrito por su ID
func (r *Resolver) DeleteCartByID(ctx context.Context, cartID string) (string, error) {
	// Buscar el carrito por su ID
	var carrito model.Carrito
	if err := r.DB.Where("cart_id = ?", cartID).First(&carrito).Error; err != nil {
		return "", errors.New("carrito no encontrado")
	}

	// Eliminar el carrito de la base de datos
	if err := r.DB.Delete(&carrito).Error; err != nil {
		return "", errors.New("no se pudo eliminar el carrito")
	}

	return "Carrito eliminado exitosamente", nil
}

// DeleteCartByCourseID elimina el carrito de un usuario por courseID.
func (r *Resolver) DeleteCartByCourseID(ctx context.Context, courseID string) (string, error) {
	// Verificar si el curso existe en el servicio de cursos.
	courseExists, err := r.checkCourseExists(courseID)
	if err != nil {
		return "", fmt.Errorf("error al verificar el curso: %v", err)
	}
	if !courseExists {
		return "", fmt.Errorf("curso con ID %s no encontrado", courseID)
	}

	// Eliminar todos los registros de carrito con el courseID especificado.
	if err := r.DB.Where("course_id = ?", courseID).Delete(&model.Carrito{}).Error; err != nil {
		return "", errors.New("no se pudo eliminar los carritos con el curso especificado")
	}

	return "Carritos eliminados exitosamente", nil
}

// RemoveFromCart elimina un curso del carrito del usuario.
func (r *Resolver) RemoveFromCart(ctx context.Context, username string, courseID string) (*bool, error) {
	// Verificar si el usuario existe y obtener su userID.
	userID, err := r.checkUserExists(username)
	if err != nil {
		return nil, fmt.Errorf("error al verificar el usuario: %v", err)
	}

	// Verificar si el curso existe en el servicio de cursos.
	courseExists, err := r.checkCourseExists(courseID)
	if err != nil {
		return nil, fmt.Errorf("error al verificar el curso: %v", err)
	}
	if !courseExists {
		return nil, fmt.Errorf("curso con ID %s no encontrado", courseID)
	}

	// Eliminar el curso del carrito del usuario usando el userID.
	if err := r.DB.Where("userID = ? AND courseID = ?", userID, courseID).Delete(&model.Carrito{}).Error; err != nil {
		return nil, err
	}

	success := true
	return &success, nil
}

// ViewCartByUserID permite ver el carrito del usuario utilizando el userID.
func (r *Resolver) ViewCartByUserID(ctx context.Context, userID string) ([]*model.Carrito, error) {
	var carrito []*model.Carrito

	// Buscar todos los elementos del carrito asociados al userID.
	if err := r.DB.Where("user_id = ?", userID).Find(&carrito).Error; err != nil {
		return nil, fmt.Errorf("error al obtener el carrito: %v", err)
	}

	return carrito, nil
}

// ViewCartByUsername permite ver el carrito del usuario utilizando el nombre de usuario.
func (r *Resolver) ViewCartByUsername(ctx context.Context, username string) ([]*model.Carrito, error) {
	// Verificar si el usuario existe y obtener el userID.
	userID, err := r.checkUserExists(username)
	if err != nil {
		return nil, fmt.Errorf("error al verificar el usuario: %v", err)
	}

	// Buscar todos los elementos del carrito asociados al userID.
	var carrito []*model.Carrito
	if err := r.DB.Where("user_id = ?", userID).Find(&carrito).Error; err != nil {
		return nil, fmt.Errorf("error al obtener el carrito: %v", err)
	}

	return carrito, nil
}

// GetAllUsers devuelve todos los usuarios.
func (r *Resolver) GetAllUsers(ctx context.Context) ([]*model.Usuario, error) {
	var users []*model.Usuario

	// Consultar todos los usuarios en la base de datos.
	if err := r.DB.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("error al obtener los usuarios: %v", err)
	}

	return users, nil
}

// checkUserExists verifica si un usuario existe en la base de datos y devuelve su userID.
func (r *Resolver) checkUserExists(username string) (string, error) {
	var usuario model.Usuario
	if err := r.DB.Where("username = ?", username).First(&usuario).Error; err != nil {
		return "", fmt.Errorf("usuario no encontrado")
	}
	return usuario.UserID, nil
}

// checkCourseExists verifica si un curso existe en el servicio de cursos.
func (r *Resolver) checkCourseExists(courseID string) (bool, error) {
	url := "http://localhost:8081/graphql" // Asegúrate de que esta URL sea correcta.
	query := fmt.Sprintf(`{"query": "query { cursoByID(courseID: \"%s\") { courseID } }"}`, courseID)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("error al verificar el curso, código de respuesta: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	// Verificar si la respuesta contiene el curso.
	if data, found := result["data"].(map[string]interface{}); found {
		if course, exists := data["cursoByID"].(map[string]interface{}); exists && course["courseID"] != nil {
			return true, nil
		}
	}

	return false, nil
}

func (r *queryResolver) UsuarioByUsername(ctx context.Context, username string) (*model.Usuario, error) {
	var usuario model.Usuario

	// Buscar el usuario por nombre de usuario en la base de datos
	if err := r.DB.Where("username = ?", username).First(&usuario).Error; err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %v", err)
	}

	return &usuario, nil
}

func generateUniqueID() string {
	return uuid.NewString()
}
