package mossservice

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/isayme/go-amqp-reconnect/rabbitmq"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

// MossTask is a template for a moss task
type MossTask struct {
	ProblemID   string   `json:"problem_id"`
	ClassID     string   `json:"class_id"`
	Language    string   `json:"language"`
	Submissions []string `json:"submission_ids"`
}

var supportedLanguage = map[string]bool{
	"java":    true,
	"python3": true,
	"clang":   true,
	"cpp":     true,
}
var conn *rabbitmq.Connection
var channel *rabbitmq.Channel
var queue amqp.Queue

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Setup connect rabbitmq
func Setup() {
	var err error

	if gin.Mode() == "test" {
		return
	}
	if gin.Mode() == "debug" {
		err = godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	}

	conn, err = rabbitmq.Dial(os.Getenv("RABBITMQ_HOST"))
	failOnError(err, "Failed to connect to RabbitMQ")
	channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
	queue, err = channel.QueueDeclare(
		"program_moss", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")
}

// ErrUnsupportedLanguage is returned when the language is not supported
var ErrUnsupportedLanguage = errors.New("unsupported language")

// Validate validates the style task
func (j *MossTask) Validate() error {
	if _, ok := supportedLanguage[j.Language]; !ok {
		return ErrUnsupportedLanguage
	}
	return nil
}

// Run  Run a new submission
func (j *MossTask) Run() (err error) {
	var data []byte
	if data, err = json.Marshal(j); err != nil {
		return
	}

	if gin.Mode() == "test" {
		return
	}

	err = channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/json",
			Body:         data,
		},
	)

	return
}
