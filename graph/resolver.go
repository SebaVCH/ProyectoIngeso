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
	"net/http"
	"time"

	"github.com/google/uuid"
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

func (r *Resolver) ActualizarUsernameConEmail(ctx context.Context, email string, newUsername string) (*models.Usuario, error) {
	var usuario models.Usuario

	// Buscar el usuario por su email
	if err := r.DB.Where("email = ?", email).First(&usuario).Error; err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	// Verificar si el nuevo username ya está en uso
	var existingUser models.Usuario
	if err := r.DB.Where("username = ?", newUsername).First(&existingUser).Error; err == nil {
		return nil, errors.New("el nombre de usuario ya está en uso")
	}

	// Actualizar el username
	usuario.Username = newUsername
	if err := r.DB.Save(&usuario).Error; err != nil {
		return nil, errors.New("no se pudo actualizar el nombre de usuario")
	}

	return &usuario, nil
}

func (r *Resolver) ActualizarNombreCompleto(ctx context.Context, email string, newNameLastName string) (*models.Usuario, error) {
	var usuario models.Usuario

	// Buscar el usuario por su email
	if err := r.DB.Where("email = ?", email).First(&usuario).Error; err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	// Actualizar el nombre completo
	usuario.NameLastName = newNameLastName
	if err := r.DB.Save(&usuario).Error; err != nil {
		return nil, errors.New("no se pudo actualizar el nombre completo")
	}

	return &usuario, nil
}

func (r *Resolver) ActualizarEmail(ctx context.Context, email string, newEmail string) (*models.Usuario, error) {
	var usuario models.Usuario

	// Buscar el usuario por su email actual
	if err := r.DB.Where("email = ?", email).First(&usuario).Error; err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	// Verificar si el nuevo email ya está en uso
	var existingUser models.Usuario
	if err := r.DB.Where("email = ?", newEmail).First(&existingUser).Error; err == nil {
		return nil, errors.New("el email ya está en uso")
	}

	// Actualizar el email
	usuario.Email = newEmail
	if err := r.DB.Save(&usuario).Error; err != nil {
		return nil, errors.New("no se pudo actualizar el email")
	}

	return &usuario, nil
}
func (r *Resolver) ActualizarContrasena(ctx context.Context, email string, oldPassword string, newPassword string) (string, error) {
	var usuario models.Usuario

	// Buscar el usuario por el email
	if err := r.DB.Where("email = ?", email).First(&usuario).Error; err != nil {
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

// AddToCartByEmail agrega un curso al carrito del usuario utilizando el correo electrónico.
func (r *Resolver) AddToCartbyEmail(ctx context.Context, email string, courseID string) (*model.Carrito, error) {
	// Verificar si el usuario existe y obtener el userID mediante el email.
	userID, err := r.checkUserExistsByEmail(email)
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

	// Obtener los cursos del usuario mediante la consulta GetCoursesByEmail
	userCourses, err := r.GetCoursesByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error al obtener los cursos del usuario: %v", err)
	}

	// Verificar si el curso ya está en la lista de cursos del usuario
	for _, usuarioCurso := range userCourses {
		if usuarioCurso.CourseID == courseID {
			return nil, fmt.Errorf("el usuario ya tiene este curso en su lista de cursos")
		}
	}

	// Verificar si el curso ya está en el carrito del usuario
	existingCartItem := &model.Carrito{}
	err = r.DB.Where("user_id = ? AND course_id = ?", userID, courseID).First(&existingCartItem).Error
	if err == nil {
		// Si se encuentra un item en el carrito, no lo agregamos de nuevo.
		return nil, fmt.Errorf("el curso ya está en tu carrito")
	}

	// Crear un nuevo elemento en el carrito.
	cartItem := &model.Carrito{
		CartID:   uuid.New().String(),
		UserID:   userID,
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

// ViewCartByEmail permite ver el carrito del usuario utilizando el email.
func (r *Resolver) ViewCartByEmail(ctx context.Context, email string) ([]*model.Carrito, error) {
	// Buscar el usuario por su email y obtener el userID.
	var usuario model.Usuario
	if err := r.DB.Where("email = ?", email).First(&usuario).Error; err != nil {
		return nil, fmt.Errorf("usuario no encontrado")
	}

	// Buscar todos los elementos del carrito asociados al userID.
	var carrito []*model.Carrito
	if err := r.DB.Where("user_id = ?", usuario.UserID).Find(&carrito).Error; err != nil {
		return nil, fmt.Errorf("error al obtener el carrito: %v", err)
	}

	return carrito, nil
}

// AddCourseToUser agrega un curso a la lista de cursos de un usuario por su ID.
func (r *Resolver) AddCourseToUser(ctx context.Context, email string, courseID string) (string, error) {
	// Verificar si el curso existe usando la función `checkCourseExists`.
	exists, err := r.checkCourseExists(courseID)
	if err != nil {
		return "", fmt.Errorf("error al verificar la existencia del curso: %v", err)
	}
	if !exists {
		return "", fmt.Errorf("el curso con ID %s no existe", courseID)
	}

	// Verificar si el usuario existe usando el email.
	var usuario model.Usuario
	if err := r.DB.Where("email = ?", email).First(&usuario).Error; err != nil {
		return "", fmt.Errorf("usuario no encontrado: %v", err)
	}

	// Verificar si la relación usuario-curso ya existe.
	var relacion model.UsuarioCurso
	if err := r.DB.Where("email = ? AND course_id = ?", email, courseID).First(&relacion).Error; err == nil {
		return "", fmt.Errorf("el usuario ya tiene este curso agregado")
	}

	// Crear la nueva relación usuario-curso.
	nuevaRelacion := model.UsuarioCurso{
		ID:       generateUniqueID(),
		Email:    email, // Puedes optar por guardar el `username` del usuario si es necesario.
		CourseID: courseID,
	}
	if err := r.DB.Create(&nuevaRelacion).Error; err != nil {
		return "", fmt.Errorf("error al agregar el curso al usuario: %v", err)
	}

	return "Curso agregado exitosamente al usuario", nil
}

// GetCoursesByEmail obtiene los cursos asociados a un usuario dado su email.
func (r *Resolver) GetCoursesByEmail(ctx context.Context, email string) ([]model.UsuarioCurso, error) {
	// Verificar si existen cursos asociados al email proporcionado.
	var cursos []model.UsuarioCurso
	if err := r.DB.Where("email = ?", email).Find(&cursos).Error; err != nil {
		return nil, fmt.Errorf("error al obtener los cursos para el email %s: %v", email, err)
	}

	return cursos, nil
}

func (r *Resolver) ObtenerUsernamePorEmail(ctx context.Context, email string) (string, error) {
    var usuario models.Usuario

    // Buscar el usuario por el email
    if err := r.DB.Where("email = ?", email).First(&usuario).Error; err != nil {
        return "", errors.New("usuario no encontrado")
    }

    // Devolver el nombre de usuario
    return usuario.Username, nil
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
func (r *Resolver) checkUserExistsByEmail(email string) (string, error) {
	var usuario model.Usuario
	if err := r.DB.Where("email = ?", email).First(&usuario).Error; err != nil {
		return "", err
	}
	return usuario.UserID, nil
}

// checkCourseExists verifica si un curso existe en el servicio de cursos.
func (r *Resolver) checkCourseExists(courseID string) (bool, error) {
	url := "http://proyectoingesocursos:8081/graphql" // Asegúrate de que esta URL sea correcta.
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
