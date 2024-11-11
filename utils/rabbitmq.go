package utils

import (
    "log"
    "github.com/streadway/amqp"
)

func ConnectRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
    conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %s", err)
        return nil, nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %s", err)
        return conn, nil, err
    }

    return conn, ch, nil
}
