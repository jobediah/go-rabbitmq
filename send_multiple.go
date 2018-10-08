package main

import (
    "fmt"
    "log"
    "time"
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
    flag := false
    go func () {
    for{  //receive loop
        select {  //check connection
            case err = <-notify:
            flag = true
            fmt.Println("Broken connection")
            //work with error
            break //reconnect
        }
    time.Sleep(3*time.Second)
    }
    }()

    ch, err := conn.Channel()
    failOnError(err, "Failed to open a channel")
    defer ch.Close()
    
    q, err := ch.QueueDeclare(
        "hello", 
        false,
        false,
        false,
        false,
        nil,
    )
    failOnError(err, "failed to declare a queue")

    body := "hello"
    
    for {
    if flag {
    conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
    //failOnError(err, "Failed to connect to RabbitMQ")
    defer conn.Close() 
    if err == nil { continue }
    ch, err = conn.Channel()
    //failOnError(err, "Failed to open a channel")
    if err == nil { continue }
    flag = false
    defer ch.Close()

    q, err = ch.QueueDeclare(
        "hello", 
        false,
        false,
        false,
        false,
        nil,
    )
    if err == nil { continue }
    //failOnError(err, "failed to declare a queue")

    }
    _ = ch.Publish(
        "", 
        q.Name,
        false,
        false,
        amqp.Publishing {
            ContentType: "text/plain",
            Body: []byte(body),
        })
    log.Printf(" [x] %s", body)
    //failOnError(err, "failed to publicsh a message")
    time.Sleep(3*time.Second)
    }
}
