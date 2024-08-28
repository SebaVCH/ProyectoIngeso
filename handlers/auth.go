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

	// Cifrar la contraseña
	hash, err := utils.HashContrasena(usuario.ContrasenaHash)
	if err != nil {
		http.Error(w, "Error al cifrar la contraseña", http.StatusInternalServerError)
		return
	}
	usuario.ContrasenaHash = hash

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
		CorreoElectronico string `json:"correo_electronico"`
		Contrasena        string `json:"contrasena"`
	}
	if err := json.NewDecoder(r.Body).Decode(&datosEntrada); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var usuario models.Usuario
	if err := bd.Where("correo_electronico = ?", datosEntrada.CorreoElectronico).First(&usuario).Error; err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusUnauthorized)
		return
	}

	if !utils.VerificarHashContrasena(datosEntrada.Contrasena, usuario.ContrasenaHash) {
		http.Error(w, "Contraseña inválida", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Inicio de sesión exitoso")
}
