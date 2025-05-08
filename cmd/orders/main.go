package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	handler "github.com/lioarce01/real-time-order-processing-system/internal/adapters/interface/http"
	"github.com/lioarce01/real-time-order-processing-system/internal/adapters/messaging"
	pgrepo "github.com/lioarce01/real-time-order-processing-system/internal/adapters/repository/postgres"
	"github.com/lioarce01/real-time-order-processing-system/internal/core/service"
	"go.uber.org/zap"
)

func main() {
	// ─── Logger ───────────────────────────────────────────────────────────
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync()

	// ─── Configuration ───────────────────────────────────────────────────
	kafkaBrokers := []string{getEnv("KAFKA_BROKER", "kafka:9092")}
	kafkaGroupID := getEnv("KAFKA_GROUP_ID", "order-service-group")
	postgresDSN := getEnv("POSTGRES_DSN", "host=postgres user=postgres password=postgres dbname=ordersdb sslmode=disable")

	// Kafka topics to consume
	kafkaTopics := []string{"orders"}

	// ─── Repository (Postgres / GORM) ────────────────────────────────────
	repo, err := pgrepo.NewOrderRepository(postgresDSN, logger)
	if err != nil {
		logger.Fatal("failed to init Postgres repo", zap.Error(err))
	}

	// ─── 4. Publisher (Kafka) ────────────────────────────────────────────────
	pub, err := messaging.NewKafkaPublisher(kafkaBrokers, logger)
	if err != nil {
		logger.Fatal("failed to init Kafka publisher", zap.Error(err))
	}

	// ─── Core Service ─────────────────────────────────────────────────────
	orderSvc := service.NewOrderService(repo, pub, logger)

	// ─── Consumer Group (Kafka) ───────────────────────────────────────────
	// Build Sarama config
	saramaCfg := sarama.NewConfig()
	saramaCfg.Consumer.Return.Errors = true
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	// Create and run group
	go func() {
		if err := messaging.NewKafkaConsumerGroup(
			kafkaBrokers,
			kafkaGroupID,
			kafkaTopics,
			orderSvc,
			logger,
			saramaCfg,
		); err != nil {
			logger.Fatal("failed to start Kafka consumer group", zap.Error(err))
		}
	}()

	// ─── HTTP Handlers ────────────────────────────────────────────────────
	orderHandler := handler.NewOrderHandler(orderSvc, logger)

	http.HandleFunc("/orders", orderHandler.CreateOrder)
	http.HandleFunc("/orders/get", orderHandler.GetOrder)

	// Start HTTP Server
	go func() {
		port := "8080"
		fmt.Printf("Starting server on :%s...\n", port)
		if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
			logger.Fatal("failed to start HTTP server", zap.Error(err))
		}
	}()

	// ─── Shutdown ────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down service...")

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if pubImpl, ok := pub.(*messaging.KafkaPublisher); ok {
		if err := pubImpl.Close(); err != nil {
			logger.Warn("failed to close Kafka producer", zap.Error(err))
		}
	}

	logger.Info("service stopped")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
