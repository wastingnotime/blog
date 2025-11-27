package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wastingnotime/blog/internal/analytics/consumer"
)

func main() {
	cfg, err := consumer.ConfigFromEnv()
	if err != nil {
		log.Fatalf("invalid config: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Printf("starting analytics consumer: region=%s queue=%s plausible=%s",
		cfg.AWSRegion, cfg.QueueURL, cfg.PlausibleURL)

	if err := consumer.Run(ctx, cfg); err != nil {
		log.Printf("consumer stopped with error: %v", err)
		// small grace period so logs flush
		time.Sleep(2 * time.Second)
	}
}
