package main

import (
	"ProyectoIngeso/graph"
	"ProyectoIngeso/graph/model"
	"ProyectoIngeso/models"
	"ProyectoIngeso/mq"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors" // Importar el middleware CORS
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

	// Migrar el esquema de Curso y Usuario
	err = bd.AutoMigrate(&models.Carrito{})
	if err != nil {
		log.Fatal("Error al migrar la base de datos", err)
	}

	// Migración automática del modelo Usuario
	err = bd.AutoMigrate(
		&models.Usuario{},
		&models.Reseña{},
		&model.Carrito{},
		&models.Carrito{},
		&models.Pago{},
		&models.Notificación{},
	)
	if err != nil {
		return
	}
}

func main() {
	// Iniciar consumidor de RabbitMQ
	go func() {
		err := utils.StartUserConsumer() // Asegúrate de que la función StartUserConsumer sea pública
		if err != nil {
			log.Fatalf("Error al iniciar el consumidor de RabbitMQ: %s", err)
		}
	}()
	// Resolver
	resolver := graph.Resolver{DB: bd}

	// Servidor GraphQL
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver}))

	// Middleware CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Cambia esto si tu frontend está en otro dominio o puerto
		AllowCredentials: true,
	}).Handler(srv)

	http.Handle("/graphql", corsHandler)
	http.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))

	log.Println("Iniciando servidor en :8080...")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("No se pudo iniciar el servidor: %s\n", err)
	}
}
