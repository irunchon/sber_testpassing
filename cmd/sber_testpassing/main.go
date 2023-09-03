package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/irunchon/sber_testpassing/internal/app/passing_webtest"
)

func main() {
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
