package main

import (
	"ProyectoIngeso/graph" // Cambiar de graphql a graph
	"ProyectoIngeso/models"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
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
	bd.AutoMigrate(
		&models.Usuario{},
		&models.Reseña{},
		&models.Curso{},
		&models.Carrito{},
		&models.Pago{},
		&models.Notificación{},
	)
}

func main() {
	// Resolver
	resolver := graph.Resolver{DB: bd}

	// Servidor GraphQL
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver}))

	http.Handle("/graphql", srv)
	http.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))

	log.Println("Iniciando servidor en :8080...")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("No se pudo iniciar el servidor: %s\n", err)
	}
}
