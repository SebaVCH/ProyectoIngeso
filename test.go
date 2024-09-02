package main

import (
	"ProyectoIngeso/handlers"
	"log"
	"net/http"
)

func main() {
	// Asocia las rutas /registro y /iniciar-sesion a las funciones correspondientes
	http.HandleFunc("/register", handlers.RegistrarUsuario)
	http.HandleFunc("/login", handlers.IniciarSesion)

	log.Println("Iniciando servidor en :8080...")

	var err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("No se pudo iniciar el servidor: %s\n", err)
	}

}
