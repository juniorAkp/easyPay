package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/juniorAkp/easyPay/internal/config"
	"github.com/juniorAkp/easyPay/internal/consumer"
	"github.com/juniorAkp/easyPay/internal/database"
	"github.com/juniorAkp/easyPay/internal/handler"
	"github.com/juniorAkp/easyPay/internal/queue"
	"github.com/juniorAkp/easyPay/internal/repository"
	"github.com/juniorAkp/easyPay/internal/services/momo"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg := config.Load()

	conn, err := database.NewPostgresPool()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	repo := repository.NewRepository(conn)

	r := gin.Default()
	rabbit := queue.NewRabbitmq(cfg.RabbitmqUrl)

	//consumer
	consumer := consumer.NewConsumer(rabbit, repo)
	consumer.StartConsumers()

	h := handler.NewHandler(r, repo, cfg, momo.NewMoMoService(), rabbit)

	h.RegisterRoutes()

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
