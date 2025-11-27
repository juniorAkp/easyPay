package queue

import (
	"encoding/json"
	"log"

	"github.com/juniorAkp/easyPay/pkg/types"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbitmq struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func NewRabbitmq(url string) *Rabbitmq {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")

	// Channel config
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	return &Rabbitmq{
		conn:    conn,
		channel: ch,
	}
}

func (r *Rabbitmq) DeclareQueue(queueName string) error {
	q, err := r.channel.QueueDeclare(
		queueName, // name
		true,      // durable (survives broker restart)
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}
	r.queue = q
	log.Printf("Queue declared: %s", queueName)
	return nil
}

func (r *Rabbitmq) Publish(queueName string, message any) error {
	// Ensure queue exists before publishing
	_, err := r.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return r.channel.Publish("", queueName, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		DeliveryMode: amqp.Persistent, // Make message persistent
	})
}

func (r *Rabbitmq) Consume(queueName string, handler func(msg types.IncomingMessage) error) error {
	// Ensure queue exists before consuming
	_, err := r.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	messages, err := r.channel.Consume(
		queueName,
		"",    // consumer tag
		false, // auto-ack (we'll manually ack)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	log.Printf("Started consuming from queue: %s", queueName)

	go func() {
		for d := range messages {
			log.Printf("Received message: %s", d.Body)

			var message types.IncomingMessage
			if err := json.Unmarshal(d.Body, &message); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				d.Nack(false, false) // Don't requeue
				continue
			}

			if err := handler(message); err != nil {
				log.Printf("Failed to handle message: %v", err)
				d.Nack(false, true) // Requeue for retry
			} else {
				log.Printf("Message processed successfully")
				d.Ack(false)
			}
		}
	}()

	return nil
}
func (r *Rabbitmq) Close() {
	if r.channel != nil {
		err := r.channel.Close()
		failOnError(err, "Failed to close channel")
	}
	if r.conn != nil {
		err := r.conn.Close()
		failOnError(err, "Failed to close connection")
	}
}
