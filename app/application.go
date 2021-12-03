package app

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/DevNiklesh/go-cron-scheduler-rabbit-mq/config"
	crontab "github.com/DevNiklesh/go-cron-scheduler-rabbit-mq/cron"
	"github.com/DevNiklesh/go-cron-scheduler-rabbit-mq/models"
	"github.com/DevNiklesh/go-cron-scheduler-rabbit-mq/rabbit"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

var (
	router = gin.Default()
	conf   = config.New()
)

func StartApplication() {
	log.Println("Starting the Cron-RabbitMQ application")

	// RabbitMQ connection
	connMQ, err := rabbit.GetConn(conf.RabbitUrl)
	if err != nil {
		log.Fatalf("rabbit connection: %v", err)
	}
	defer connMQ.Close()

	// Declare Exchange topic
	err = connMQ.DeclareTopicExchange(conf.Exchange)
	if err != nil {
		log.Fatalf("declare exchange: %v", err)
	}

	// Starting Consumer
	connMQ.StartConsumer(conf.Exchange, conf.QueueWorker, conf.KeyWorker, func(d amqp.Delivery) bool {
		var msg models.Message
		err := json.Unmarshal(d.Body, &msg)
		if err != nil {
			log.Fatalf("unmarshal message: %v", err)
			return false
		}

		fmt.Println("Message from RabbitMQ: ", msg)
		return true
	})

	// crontab is a abstraction layer for Cron libraries.
	c := crontab.New()

	// Job 1 - Adding Cron Job that publishes data to rabbitMQ worker queue
	c.AddJob("*/1 * * * *", func() {
		inputMsg := models.Message{Text: "Every 1 minute", Source: "Worker 1", Time: time.Now().Unix()}
		message, err := json.Marshal(inputMsg)
		if err != nil {
			log.Fatalf("marshal message: %v", err)
		}

		// Publishing message to RabbitMQ
		err = connMQ.Publish(conf.Exchange, conf.KeyWorker, message)
		if err != nil {
			log.Fatalf("publish message: %v", err)
		}
	})

	// Job 2
	c.AddJob("*/2 * * * *", func() {
		inputMsg := models.Message{Text: "Every 2 minute", Source: "Worker 2", Time: time.Now().Unix()}
		message, err := json.Marshal(inputMsg)
		if err != nil {
			log.Fatalf("marshal message: %v", err)
		}

		err = connMQ.Publish(conf.Exchange, conf.KeyWorker, message)
		if err != nil {
			log.Fatalf("publish message: %v", err)
		}
	})

	c.StartCron()

	router.Run(":8000")
}
