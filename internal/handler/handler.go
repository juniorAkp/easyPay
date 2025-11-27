package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juniorAkp/easyPay/internal/config"
	"github.com/juniorAkp/easyPay/internal/queue"
	"github.com/juniorAkp/easyPay/internal/repository"
	"github.com/juniorAkp/easyPay/internal/services/momo"
	"github.com/juniorAkp/easyPay/pkg/types"
)

type Handler struct {
	router      *gin.Engine
	repo        *repository.Repository
	config      *config.Config
	momoService *momo.MoMoService
	rabbitMq    *queue.Rabbitmq
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
func NewHandler(r *gin.Engine, repo *repository.Repository, cfg *config.Config, momo *momo.MoMoService, rb *queue.Rabbitmq) *Handler {
	return &Handler{
		router:      r,
		repo:        repo,
		config:      cfg,
		momoService: momo,
		rabbitMq:    rb,
	}
}

func (h *Handler) RegisterRoutes() {
	h.router.GET("/webhook", h.VerifyWebhook)
	h.router.GET("/callback", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"complete": "true"})
	})
	h.router.POST("/webhook", h.HandleWebhook)
	h.router.GET("/health", h.HealthCheck)
}

func (h *Handler) VerifyWebhook(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == h.config.VerifyToken {
		c.String(http.StatusOK, challenge)
		return
	}

	c.JSON(http.StatusForbidden, gin.H{"error": "verification failed"})
}

// HandleWebhook processes incoming WhatsApp messages
func (h *Handler) HandleWebhook(c *gin.Context) {
	var payload types.WhatsAppWebhook
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if len(payload.Entry) > 0 && len(payload.Entry[0].Changes) > 0 && len(payload.Entry[0].Changes[0].Value.Messages) > 0 {
		fmt.Println(payload)
		msg := types.IncomingMessage{
			Phone:     payload.Entry[0].Changes[0].Value.Messages[0].From,
			Name:      payload.Entry[0].Changes[0].Value.Contacts[0].Profile.Name,
			Message:   payload.Entry[0].Changes[0].Value.Messages[0].Text.Body,
			Timestamp: payload.Entry[0].Changes[0].Value.Messages[0].Timestamp,
		}
		err := h.rabbitMq.Publish("incoming_message", msg)
		failOnError(err, "Failed to publish a message")
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

//func (h *Handler) processMessage(userMessage, userPhone, userName string) {
//	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
//	defer cancel()
//
//	fmt.Printf("Processing message from %s: %s\n", userPhone, userMessage)
//
//	user, err := h.repo.GetUserByPhone(ctx, userPhone)
//	if err != nil {
//		fmt.Printf("Error fetching user: %v\n", err)
//	}
//
//	if user == nil {
//		newUser := &model.Details{
//			Username:  userName,
//			Phone:     userPhone,
//			CreatedAt: time.Now(),
//			UpdatedAt: time.Now(),
//		}
//
//		if err := h.repo.CreateUser(ctx, newUser); err != nil {
//			fmt.Printf("Failed to create user: %v\n", err)
//			return
//		}
//
//		fmt.Printf("New user created: %s\n", userPhone)
//		whatsapp.SendTemplateMessage(userPhone, string(whatsapp.EasyPayIntro))
//		return
//	}
//
//	result, err := ai.ExtractMessageDetails(ctx, userMessage)
//	if err != nil {
//		fmt.Printf("AI extraction error: %v\n", err)
//		whatsapp.SendMessage(userPhone, "Sorry, I didn’t understand that. Please send the amount and account again.")
//		return
//	}
//
//	if result.Amount <= 0 {
//		whatsapp.SendMessage(userPhone, "Sorry, I didn’t understand that. Please send the amount and account again.")
//		return
//	}
//
//	whatsapp.SendMessage(userPhone, fmt.Sprintf("⌛Processing your payment of %.2f GHS...", result.Amount))
//
//	refId, err := h.momoService.RequestPay(
//		strconv.FormatFloat(result.Amount, 'f', -1, 64),
//		userPhone,
//		result.Message,
//	)
//	if err != nil {
//		fmt.Printf("MTN MoMo error: %v\n", err)
//		whatsapp.SendMessage(userPhone, "Payment failed. Please try again.")
//		return
//	}
//
//	//@TODO //change to webhooks
//	time.Sleep(5 * time.Second)
//	status, _ := h.momoService.CheckTransactionStatus(refId)
//
//	switch status {
//	case "SUCCESSFUL":
//		whatsapp.SendMessage(userPhone, "✅ Payment successful!")
//	case "PENDING":
//		whatsapp.SendMessage(userPhone, "⌛ Payment pending. Please approve it on your phone.")
//	default:
//		whatsapp.SendMessage(userPhone, "❌ Payment failed or was cancelled.")
//	}
//}

// HealthCheck endpoint
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "easyPay",
	})
}
