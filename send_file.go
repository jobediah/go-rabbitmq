package main

import (
    "fmt"
    "log"
    "os"
    "io/ioutil"

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

    dat, err := ioutil.ReadFile("./example.txt")
    failOnError(err, "Could not read file")
    //fmt.Print(string(dat))
    log.Print(string(dat))

    f, err := os.Open("./example.txt")
    failOnError(err, "Could not open file")

    fileinfo, err := f.Stat()
    failOnError(err, "Could not get file stats")

    filesize := fileinfo.Size()
    buffer := make([]byte, filesize)

    file, err := f.Read(buffer)
    failOnError(err, "Could not read file")

    log.Printf("%d bytes: %s\n", file, string(buffer)) 
   

    f.Close()

    err = ch.Publish(
        "", 
        q.Name,
        false,
        false,
        amqp.Publishing {
            ContentType: "text/plain",
            Body: []byte(buffer),
        })
    failOnError(err, "failed to publicsh a message")
}
