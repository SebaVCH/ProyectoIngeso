package utils

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/streadway/amqp"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "ProyectoIngeso/models"
    "ProyectoIngeso/utils"
)

// Definir la estructura que corresponde al mensaje recibido desde RabbitMQ
type RabbitMQMessage struct {
    Pattern string `json:"pattern"`
    Data    string `json:"data"`
    ID      string `json:"id"`
}

// Esta función será llamada desde main.go para iniciar el consumer
func StartUserConsumer() error {
    // Conectar a RabbitMQ
    conn, ch, err := utils.ConnectRabbitMQ()
    if err != nil {
        return fmt.Errorf("error connecting to RabbitMQ: %w", err)
    }
    defer conn.Close()
    defer ch.Close()

    // Conectar a la base de datos
    db, err := gorm.Open(sqlite.Open("base_datos.db"), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("error connecting to database: %w", err)
    }

    // Declarar la cola para escuchar las solicitudes
    queueName := "users_queue"
    q, err := ch.QueueDeclare(
        queueName, // name
        false,     // durable
        false,     // delete when unused
        false,     // exclusive
        false,     // no-wait
        nil,       // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare a queue: %w", err)
    }

    msgs, err := ch.Consume(
        q.Name, // queue
        "",     // consumer
        true,   // auto-ack
        false,  // exclusive
        false,  // no-local
        false,  // no-wait
        nil,    // args
    )
    if err != nil {
        return fmt.Errorf("failed to register a consumer: %w", err)
    }

    go func() {
        for d := range msgs {
            // Imprimir el cuerpo del mensaje recibido para depuración
            fmt.Printf("Mensaje recibido: %s\n", string(d.Body))

            // Deserializar el mensaje recibido en la estructura RabbitMQMessage
            var msg RabbitMQMessage
            err := json.Unmarshal(d.Body, &msg)
            if err != nil {
                log.Printf("Error unmarshalling message: %s", err)
                continue
            }

            // Extraer el userID del campo Data
            userID := msg.Data
            fmt.Printf("UserID deserializado: %s\n", userID)

            // Obtener la información del usuario desde la base de datos
            var usuario models.Usuario
            result := db.Where("user_id = ?", userID).First(&usuario)
            if result.Error != nil {
                log.Printf("Error finding user with ID %s: %s", userID, result.Error)
                continue
            }

            // Responder con el email del usuario
            response := struct {
                Email string `json:"email"`
            }{
                Email: usuario.Email,
            }

            responseBody, err := json.Marshal(response)
            if err != nil {
                log.Printf("Error marshalling response: %s", err)
                continue
            }

            // Publicar la respuesta en la cola de respuesta
            err = ch.Publish(
                "",        // exchange
                d.ReplyTo, // routing key
                false,     // mandatory
                false,     // immediate
                amqp.Publishing{
                    ContentType:   "application/json",
                    CorrelationId: d.CorrelationId,
                    Body:          responseBody,
                })
            if err != nil {
                log.Printf("Failed to publish a response: %s", err)
            }
        }
    }()

    log.Printf("Waiting for messages. To exit press CTRL+C")
    <-make(chan bool)

    return nil
}
