package main

import (
    "fmt"
    "log"
    "reflect"

    "github.com/streadway/amqp"

)

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
        panic(fmt.Sprintf("%s: %s", msg, err))
    }
}

func main() {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close() 

    notify := conn.NotifyClose(make(chan *amqp.Error))
    go func () {
    //for{  //receive loop
        select {  //check connection
            case err = <-notify:
            fmt.Println("Broken connection")
            //work with error
            break //reconnect
        }
    //}
    }()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()
    
    fmt.Println(reflect.TypeOf(ch))
    fmt.Println(reflect.TypeOf(err))

    q, err := ch.QueueDeclare(
        "hello", 
        false,
        false,
        false,
        false,
        nil,
    )
    failOnError(err, "failed to declare a queue")

    fmt.Println(reflect.TypeOf(q))

    body := "hello"

    err = ch.Publish(
        "", 
        q.Name,
        false,
        false,
        amqp.Publishing {
            ContentType: "text/plain",
            Body: []byte(body),
        })
    log.Printf(" [x] %s", body)
    failOnError(err, "failed to publicsh a message")
}
