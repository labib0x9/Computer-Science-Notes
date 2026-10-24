package main

import (
	"context"
	"log/slog"
	"os"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Starting application", "version", "1.0.0", "port", 8080)
	slog.Warn("Missing configuration", "config", "db_url")
	slog.Error("Failed to connect", "error", "timeout")

	// JSON
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger = slog.New(jsonHandler)
	slog.SetDefault(logger)

	//
	ctx := context.WithValue(context.Background(), "trace_id", "abc123")
	slog.InfoContext(ctx, "processing request", "user_id", 42)

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo, // Only log Info and higher
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger = slog.New(handler)
	slog.SetDefault(logger)

}

// When Should You Use slog?

//     Building production-grade Go services

//     Working with cloud infrastructure (e.g., GCP, AWS, k8s)

//     Need structured logs for log aggregation systems (e.g., Loki, Elastic, Datadog)

//     Want fine-grained control over log levels and context propagation
