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
    Pattern string `json:"pattern"` // "get_user_id" o "get_user_name"
    Data    string `json:"data"`    // Email del usuario para "get_user_id"
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

            // Definir la variable para la respuesta
            var responseBody []byte

            switch msg.Pattern {
            case "get_user_id":
                // Buscar el usuario por email
                var usuario models.Usuario
                result := db.Where("email = ?", msg.Data).First(&usuario)
                if result.Error != nil {
                    log.Printf("No se encontró un usuario con el correo %s: %s", msg.Data, result.Error)
                    continue
                }

                // Responder con la userId del usuario
                response := struct {
                    UserID string `json:"userID"`
                }{
                    UserID: usuario.UserID,
                }

                responseBody, err = json.Marshal(response)
                if err != nil {
                    log.Printf("Error marshalling response: %s", err)
                    continue
                }

            case "get_user_name":
                // Buscar el usuario por su ID
                var usuario models.Usuario
                result := db.Where("user_id = ?", msg.Data).First(&usuario)
                if result.Error != nil {
                    log.Printf("No se encontró un usuario con ID %s: %s", msg.Data, result.Error)
                    continue
                }

                // Responder con el nombre del usuario
                response := struct {
                    Name string `json:"name"`
                }{
                    Name: usuario.NameLastName, // Usar la columna de nombre completo
                }

                responseBody, err = json.Marshal(response)
                if err != nil {
                    log.Printf("Error marshalling response: %s", err)
                    continue
                }

            default:
                log.Printf("Patrón no soportado: %s", msg.Pattern)
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
