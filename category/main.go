package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/s1ntaxe770r/petstore/handlers"
	"github.com/s1ntaxe770r/petstore/models"
	"github.com/s1ntaxe770r/petstore/utils"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	REDIS_HOST = os.Getenv("REDIS_HOST")
	REDIS_PORT = os.Getenv("REDIS_PORT")
	RMQ_HOST   = os.Getenv("RMQ_HOST")
)

func main() {
	r := gin.New()

	// create redis client
	db := redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST + ":" + REDIS_PORT,
		Password: "",
		DB:       0,
	})
	r.GET("/pets/categories/reptile", handlers.ReptileHandler(db))
	r.GET("/pets/categories/dog", handlers.DogHandler(db))
	r.GET("/pets/categories/cat", handlers.CatHandler(db))

	conn, err := amqp.Dial(RMQ_HOST)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// create a new channel
	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()
	err = ch.ExchangeDeclare(
		"pets",   // name
		"fanout", // type
		false,    // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	utils.FailOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"pets", // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	// bind to pet-category queue
	err = ch.QueueBind(
		q.Name,          // queue name
		"pets-category", // routing key
		"pets",          // exchange
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to bind a queue")
	_ = ch.Qos(10, 0, false)
	logrus.Info("[*] started consumer")
	// q, err := ch.QueueDeclare(
	// 	"pets", // name
	// 	true,   // durable
	// 	false,  // delete when unused
	// 	false,  // exclusive
	// 	false,  // no-wait
	// 	nil,    // arguments
	// )
	// utils.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		"pets",     // queue
		"category", // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-waitx
		nil,        // args
	)
	utils.FailOnError(err, "Failed to register a consumer")

	// start the http server in a goroutine
	go func() {
		r.Run(":3000")
	}()
	forever := make(chan bool)

	go func() {
		logrus.Info(" [*] Waiting for messages. To exit press CTRL+C")
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			logrus.Info(fmt.Sprintf("received message %s", string(d.Body)))
			ctx := context.Background()
			var pet models.Pet
			err := json.Unmarshal(d.Body, &pet)
			if err != nil {
				logrus.Info(err)
			}

			// insert pet category in the format ["reptile","sam","alonso"]
			err = db.LPush(ctx, pet.Category, pet.Name).Err()
			if err != nil {
				logrus.Fatal(err)
			}
		}

	}()
	<-forever
}
