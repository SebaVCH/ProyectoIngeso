package main

import (
    "ProyectoIngeso/graph"
    "ProyectoIngeso/models"
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

    // Middleware CORS
    corsHandler := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000"}, // Cambia esto si tu frontend está en otro dominio o puerto
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