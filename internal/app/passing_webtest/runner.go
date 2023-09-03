package passing_webtest

import (
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"
)

func Runner(qtyOfThreads int, startURL, finalURL string, limiter <-chan time.Time) {
	wg := sync.WaitGroup{}
	wg.Add(qtyOfThreads)

	successRate := 0

	slog.Info("Test runner is working")
	for i := 0; i < qtyOfThreads; i++ {
		go func(n int) {
			worker, err := NewWorker(limiter, startURL, finalURL)
			if err != nil {
				slog.Info(fmt.Sprintf("Process #%d: %s", n, err))
				return
			}
			slog.Info(fmt.Sprintf("Process #%d: starting the test", n))
			err = worker.PassingTest()
			if err == nil {
				slog.Info(fmt.Sprintf("Process #%d: Test successfully passed", n))
				successRate++
			} else {
				log.Printf("Process #%d: %s", n, err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Printf("Successfully passed %d tests of %d\n", successRate, qtyOfThreads)
}
