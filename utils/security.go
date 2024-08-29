package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashContrasena cifra la contraseña usando bcrypt.
func HashContrasena(contrasena string) (string, error) {

	//Costo = iteraciones para crear un hash
	bytes, err := bcrypt.GenerateFromPassword([]byte(contrasena), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerificarHashContrasena compara la contraseña en texto plano con la contraseña cifrada.
func VerificarHashContrasena(contrasena, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(contrasena))
	return err == nil
}
