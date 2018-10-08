package main

import (
    "fmt"
    "log"
    "os"
    "strings"

    "github.com/streadway/amqp"

)

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
        panic(fmt.Sprintf("%s: %s", msg, err))
    }
}

func bodyFrom(args []string) string {
        var s string
        if (len(args) < 2) || os.Args[1] == "" {
                s = "hello"
        } else {
                s = strings.Join(args[1:], " ")
        }
        return s
}

func main() {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close() 

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()

    err = ch.ExchangeDeclare(
        "logs",
        "fanout", 
        true,
        false,
        false,
        false,
        nil,
    )
    failOnError(err, "failed to declare a queue")

    body := bodyFrom(os.Args)

    err = ch.Publish(
        "logs", 
        "",
        false,
        false,
        amqp.Publishing {
            DeliveryMode: amqp.Persistent,
            ContentType: "text/plain",
            Body: []byte(body),
        })
    log.Printf(" [x] %s", body)
    failOnError(err, "failed to publicsh a message")
}
