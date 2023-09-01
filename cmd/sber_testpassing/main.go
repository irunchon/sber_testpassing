package main

import (
	"log"
	"os"
	"strconv"

	"github.com/irunchon/sber_testpassing/internal/app/passing_webtest"
)

func main() {
	startURL := os.Getenv("START_PAGE")
	finalURL := os.Getenv("FINAL_PAGE")
	qtyOfThreads, err := strconv.Atoi(os.Getenv("QTY_OF_THREADS"))
	if err != nil {
		log.Fatalf("Failed to parse quantity of threads: %s\n", os.Getenv("QTY_OF_THREADS"))
	}
	passing_webtest.Runner(qtyOfThreads, startURL, finalURL)
}
