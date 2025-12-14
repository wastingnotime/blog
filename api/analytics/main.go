package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/wastingnotime/blog/internal/analytics/eventsink"
)

func main() {
	queueURL := os.Getenv("EVENT_QUEUE_URL")
	if queueURL == "" {
		log.Fatal("EVENT_QUEUE_URL is required")
	}

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		log.Fatal("AWS_REGION is required")
	}

	sharedSecret := os.Getenv("PLAUSIBLE_SHARED_SECRET")

	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	awsCfg.Region = region

	processor := eventsink.NewProcessor(sqs.NewFromConfig(awsCfg), queueURL)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover(), middleware.RequestID(), middleware.Logger())

	e.POST("/events/plausible", func(c echo.Context) error {
		if sharedSecret != "" {
			if token := c.Request().Header.Get("X-Plausible-Token"); token == "" || token != sharedSecret {
				return echo.NewHTTPError(http.StatusForbidden, "forbidden")
			}
		}

		var in eventsink.InboundEvent
		if err := c.Bind(&in); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
		}

		meta := eventsink.Meta{
			UserAgent: c.Request().UserAgent(),
			IP:        clientIP(c.Request()),
		}

		if err := processor.Handle(c.Request().Context(), in, meta); err != nil {
			if errors.Is(err, eventsink.ErrInvalidEvent) {
				return echo.NewHTTPError(http.StatusBadRequest, "invalid event")
			}

			log.Printf("failed to queue event: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to queue event")
		}

		return c.JSON(http.StatusAccepted, map[string]string{"status": "queued"})
	})

	// Basic health/readiness probes for Swarm.
	e.GET("/healthz", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	e.GET("/readyz", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	go func() {
		addr := ":" + port
		if err := e.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Printf("error during shutdown: %v", err)
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	if xr := r.Header.Get("X-Real-IP"); xr != "" {
		return strings.TrimSpace(xr)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}
