package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"github.com/snisid/vehicle-alert-svc/internal/alerter"
	"github.com/snisid/vehicle-alert-svc/internal/consumer"
)

type AlertEvent struct {
	AlertID       string `json:"alert_id"`
	PlateNumber   string `json:"plate_number"`
	CrimeCategory string `json:"crime_category"`
	AlertLevel    string `json:"alert_level"`
	ReportingUnit string `json:"reporting_unit"`
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	redisAddr := getEnv("SIVC_REDIS_ADDR", "redis-master:6379")
	kafkaBrokers := getEnv("SIVC_KAFKA_BROKERS", "kafka:9092")
	radioEndpoint := getEnv("ALERT_RADIO_ENDPOINT", "http://radio-gateway:8080/alert")
	smsEndpoint := getEnv("ALERT_SMS_ENDPOINT", "http://sms-gateway:8080/send")

	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	defer rdb.Close()

	radioAlerter := alerter.NewRadioAlerter(radioEndpoint)
	smsAlerter := alerter.NewSMSAlerter(smsEndpoint)
	pushAlerter := alerter.NewPushAlerter(rdb)

	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBrokers},
		Topic:   "sivc.alerts.created",
		GroupID: "vehicle-alert-svc",
	})
	defer kafkaReader.Close()

	alertConsumer := consumer.NewKafkaConsumer(kafkaReader, radioAlerter, smsAlerter, pushAlerter, logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		logger.Info("Vehicle alert service started")
		if err := alertConsumer.Start(ctx); err != nil {
			logger.Error("Consumer error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down vehicle alert service...")
	cancel()
	time.Sleep(2 * time.Second)
	fmt.Println("Vehicle alert service stopped")
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
