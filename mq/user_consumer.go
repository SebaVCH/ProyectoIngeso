package utils

import (
	"encoding/json"
	"fmt"
	"log"

	"ProyectoIngeso/models"
	"ProyectoIngeso/utils"

	"github.com/streadway/amqp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Definir la estructura que corresponde al mensaje recibido desde RabbitMQ
type RabbitMQMessage struct {
	Pattern string `json:"pattern"` // Ejemplo: "add_course_to_user"
	Data    string `json:"data"`    // Email, userID, etc.
	Extra   string `json:"extra"`   // Para datos adicionales como courseID
}

// Iniciar el consumidor de RabbitMQ desde main.go
func StartUserConsumer() error {
	// Conectar a RabbitMQ
	conn, ch, err := utils.ConnectRabbitMQ()
	if err != nil {
		return fmt.Errorf("error connecting to RabbitMQ: %w", err)
	}
	defer conn.Close()
	defer ch.Close()

	// Conectar a la base de datos SQLite
	db, err := gorm.Open(sqlite.Open("base_datos.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}

	// Declarar la cola para escuchar solicitudes
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

	// Consumir mensajes de la cola
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

	// Manejar los mensajes recibidos
	go func() {
		for d := range msgs {
			fmt.Printf("Mensaje recibido: %s\n", string(d.Body))

			var msg RabbitMQMessage
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				log.Printf("Error unmarshalling message: %s", err)
				continue
			}

			var responseBody []byte

			switch msg.Pattern {
			case "get_user_id":
				var usuario models.Usuario
				result := db.Where("email = ?", msg.Data).First(&usuario)
				if result.Error != nil {
					log.Printf("No se encontró un usuario con el correo %s: %s", msg.Data, result.Error)
					continue
				}
				response := struct {
					UserID string `json:"userID"`
				}{UserID: usuario.UserID}
				responseBody, err = json.Marshal(response)
				if err != nil {
					log.Printf("Error marshalling response: %s", err)
					continue
				}

			case "get_user_name":
				var usuario models.Usuario
				result := db.Where("user_id = ?", msg.Data).First(&usuario)
				if result.Error != nil {
					log.Printf("No se encontró un usuario con ID %s: %s", msg.Data, result.Error)
					continue
				}
				response := struct {
					Name string `json:"name"`
				}{Name: usuario.NameLastName}
				responseBody, err = json.Marshal(response)
				if err != nil {
					log.Printf("Error marshalling response: %s", err)
					continue
				}

			case "get_cart_courses":
				var carritos []models.Carrito
				result := db.Where("user_id = ?", msg.Data).Find(&carritos)
				if result.Error != nil {
					log.Printf("Error al obtener el carrito para el usuario %s: %s", msg.Data, result.Error)
					continue
				}
				responseBody, err = json.Marshal(carritos)
				if err != nil {
					log.Printf("Error al serializar la respuesta: %s", err)
					continue
				}

			case "clear_user_cart":
				// Vaciar el carrito del usuario
				userID := msg.Data
				result := db.Where("user_id = ?", userID).Delete(&models.Carrito{})
				if result.Error != nil {
					log.Printf("Error al vaciar el carrito para el usuario %s: %s", userID, result.Error)
					continue
				}
				log.Printf("Carrito vaciado para el usuario %s", userID)
				response := struct {
					Message string `json:"message"`
				}{Message: "Carrito vaciado exitosamente"}
				responseBody, err = json.Marshal(response)
				if err != nil {
					log.Printf("Error al serializar la respuesta: %s", err)
					continue
				}

			case "add_course_to_user":
				// Manejar el caso de agregar curso
				email := msg.Data
				courseID := msg.Extra

				if email == "" || courseID == "" {
					log.Printf("Datos incompletos: email=%s, courseID=%s", email, courseID)
					continue
				}

				// Verificar si el usuario existe
				var usuario models.Usuario
				result := db.Where("email = ?", email).First(&usuario)
				if result.Error != nil {
					log.Printf("No se encontró un usuario con el email %s: %s", email, result.Error)
					continue
				}

				// Verificar si la relación ya existe
				var usuarioCurso models.UsuarioCurso
				check := db.Where("email = ? AND course_id = ?", email, courseID).First(&usuarioCurso)
				if check.Error == nil {
					log.Printf("El curso %s ya está asociado al usuario %s", courseID, email)
					response := struct {
						Message string `json:"message"`
					}{Message: "El curso ya está asociado al usuario"}
					responseBody, _ = json.Marshal(response)
				} else {
					// Crear la nueva relación
					nuevoUsuarioCurso := models.UsuarioCurso{
						ID:       utils.GenerateID(), // Genera un ID único
						Email:    email,
						CourseID: courseID,
					}
					if err := db.Create(&nuevoUsuarioCurso).Error; err != nil {
						log.Printf("Error al asociar el curso %s con el usuario %s: %s", courseID, email, err)
						continue
					}

					log.Printf("Curso %s asociado exitosamente al usuario %s", courseID, email)
					response := struct {
						Message string `json:"message"`
					}{Message: "Curso agregado exitosamente"}
					responseBody, err = json.Marshal(response)
					if err != nil {
						log.Printf("Error al serializar la respuesta: %s", err)
						continue
					}
				}

			default:
				log.Printf("Patrón no soportado: %s", msg.Pattern)
				continue
			}

			// Publicar la respuesta al cliente
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

	log.Printf("Esperando mensajes. Presiona CTRL+C para salir.")
	<-make(chan bool)

	return nil
}
