package handlers

import (
	"ProyectoIngeso/models"
	"ProyectoIngeso/utils"
	"encoding/json"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

var bd *gorm.DB

func init() {
	var err error
	bd, err = gorm.Open(sqlite.Open("base_datos.db"), &gorm.Config{})
	if err != nil {
		panic("No se pudo conectar a la base de datos")
	}

	// Migración automática del modelo Usuario
	bd.AutoMigrate(&models.Usuario{})
}

// RegistrarUsuario maneja el registro de nuevos usuarios.
func RegistrarUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var usuario models.Usuario
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(usuario.Contrasena) < 8 {
		http.Error(w, "La contraseña debe tener al menos 8 caracteres", http.StatusBadRequest)
		return
	}

	// Cifrar la contraseña
	hash, err := utils.HashContrasena(usuario.Contrasena)
	if err != nil {
		http.Error(w, "Error al cifrar la contraseña", http.StatusInternalServerError)
		return
	}
	usuario.Contrasena = hash

	// Guardar el usuario en la base de datos
	if err := bd.Create(&usuario).Error; err != nil {
		http.Error(w, "No se pudo registrar el usuario", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Usuario registrado exitosamente")
}

// IniciarSesion maneja el proceso de login de los usuarios.
func IniciarSesion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var datosEntrada struct {
		Identificador string `json:"identificador"`
		Contrasena    string `json:"Contrasena"`
	}
	if err := json.NewDecoder(r.Body).Decode(&datosEntrada); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var usuario models.Usuario
	if err := bd.Where("correo_electronico = ? OR nombre_usuario = ?", datosEntrada.Identificador, datosEntrada.Identificador).First(&usuario).Error; err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusUnauthorized)
		return
	}

	if !utils.VerificarHashContrasena(datosEntrada.Contrasena, usuario.Contrasena) {
		http.Error(w, "Contraseña inválida", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Inicio de sesión exitoso")
}
