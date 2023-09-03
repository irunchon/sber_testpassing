package main

import (
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/irunchon/sber_testpassing/internal/app/passing_webtest"
)

func main() {
	loggerConfig()

	startURL := os.Getenv("START_PAGE")
	finalURL := os.Getenv("FINAL_PAGE")

	qtyOfThreads, err := strconv.Atoi(os.Getenv("QTY_OF_THREADS"))
	if err != nil {
		log.Fatalf("Failed to parse quantity of threads: %s\n", os.Getenv("QTY_OF_THREADS"))
	}

	maxRPS, err := strconv.Atoi(os.Getenv("MAX_SERVER_RPS"))
	if err != nil {
		log.Fatalf("Failed to parse RPS parameter: %s\n", os.Getenv("MAX_SERVER_RPS"))
	}
	limiter := time.Tick(getTimeLimit(maxRPS))

	passing_webtest.Runner(qtyOfThreads, startURL, finalURL, limiter)
}

func getTimeLimit(rps int) time.Duration {
	if rps <= 0 {
		return time.Second
	}
	return time.Duration(1000/rps) * time.Millisecond
}

func loggerConfig() {
	level := slog.LevelInfo
	err := level.UnmarshalText([]byte(os.Getenv("LOG_LEVEL")))
	if err != nil {
		slog.Info("Undefined log level")
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
