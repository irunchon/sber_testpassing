package passing_webtest

import (
	"log"
	"sync"
	"time"
)

func Runner(qtyOfThreads int, startURL, finalURL string) {
	wg := sync.WaitGroup{}
	wg.Add(qtyOfThreads)
	limiter := time.Tick(333 * time.Millisecond)
	worker := NewWorker(limiter)
	successRate := 0

	for i := 0; i < qtyOfThreads; i++ {
		go func(n int) {
			result := worker.PassingTest(startURL, finalURL)
			if result == nil {
				log.Printf("Process #%d: Test successfully passed", n)
				successRate++
			} else {
				log.Printf("Process #%d: %s", n, result)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Printf("\nSuccessfully passed %d tests of %d\n", successRate, qtyOfThreads)
}
