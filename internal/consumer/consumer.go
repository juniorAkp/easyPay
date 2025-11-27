package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/juniorAkp/easyPay/internal/model"
	"github.com/juniorAkp/easyPay/internal/queue"
	"github.com/juniorAkp/easyPay/internal/repository"
	"github.com/juniorAkp/easyPay/internal/services/ai"
	"github.com/juniorAkp/easyPay/internal/services/momo"
	"github.com/juniorAkp/easyPay/internal/services/whatsapp"
	"github.com/juniorAkp/easyPay/pkg/types"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s:%s", msg, err)
	}
}

type Consumer struct {
	rabbitMq *queue.Rabbitmq
	repo     *repository.Repository
	momo     *momo.MoMoService
}

func NewConsumer(rabbit *queue.Rabbitmq, repo *repository.Repository) *Consumer {
	return &Consumer{
		rabbitMq: rabbit,
		repo:     repo,
	}
}

func (r *Consumer) StartConsumers() {
	log.Println("Starting all RabbitMQ consumers...")

	go r.consumeIncomingMessages()
	go r.consumePaymentRequests()
	go r.consumeNotifications()

	log.Println("All consumers started successfully")
}

func (r *Consumer) consumeIncomingMessages() {
	err := r.rabbitMq.Consume("incoming_message", func(msg types.IncomingMessage) error {
		ctx := context.Background()

		log.Printf("Processing message from: %s", msg.Phone)

		user, err := r.repo.GetUserByPhone(ctx, msg.Phone)

		if err != nil || user == nil {
			newUser := &model.Details{
				Username: msg.Name,
				Phone:    msg.Phone,
			}

			if err := r.repo.CreateUser(ctx, newUser); err != nil {
				log.Printf("Failed to create user: %v", err)
				return err
			}

			log.Printf("New user created: %s (%s)", msg.Name, msg.Phone)

			whatsapp.SendTemplateMessage(msg.Phone, string(whatsapp.EasyPayIntro))
			return nil
		}

		result, err := ai.ExtractMessageDetails(ctx, msg.Message)
		failOnError(err, "AI extraction error")

		fmt.Println(result)

		payment := &types.MessageDetails{
			Account: result.Account,
			Message: result.Message,
			Amount:  result.Amount,
		}

		r.rabbitMq.Publish("payment_request", payment)
		return nil
	})

	if err != nil {
		log.Fatalf("Failed to start incoming messages consumer: %v", err)
	}
}

func (r *Consumer) consumePaymentRequests() {
	err := r.rabbitMq.Consume("payment_request", func(msg types.IncomingMessage) error {

		return nil
	})
}

func (r *Consumer) consumeNotifications() {

}
