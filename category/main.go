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
	db := redis.NewClient(&redis.Options{
		Addr:     REDIS_HOST + ":" + REDIS_PORT,
		Password: "",
		DB:       0,
	})
	r.GET("/pets/categories/reptile", handlers.ReptileHandler(db))
	r.GET("/pets/categories/dog", handlers.DogHandler(db))
	r.GET("/pets/categories/cat", handlers.CatHandler(db))
	go func() {
		conn, err := amqp.Dial(RMQ_HOST)
		utils.FailOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		utils.FailOnError(err, "Failed to open a channel")
		defer ch.Close()
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

		err = ch.Qos(1, 0, false)
		utils.FailOnError(err, "failed to configure Qos")

		msgs, err := ch.Consume(
			"pets", // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		utils.FailOnError(err, "Failed to register a consumer")

		forever := make(chan bool)

		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			logrus.Info(fmt.Sprintf("received message %s", string(d.Body)))
			ctx := context.Background()
			pet := models.Pet{}
			err := json.Unmarshal(d.Body, &pet)
			if err != nil {
				logrus.Info(err)
			}
			err = db.LPush(ctx, pet.Category, 0, pet.Name).Err()
			if err != nil {
				log.Fatal(err)
			}

		}
		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		<-forever

	}()
	r.Run(":3000")
}
