package passing_webtest

import (
	"log"
	"sync"
	"time"
)

func Runner(qtyOfThreads int, startURL, finalURL string, limiter <-chan time.Time) {
	wg := sync.WaitGroup{}
	wg.Add(qtyOfThreads)

	successRate := 0

	for i := 0; i < qtyOfThreads; i++ {
		go func(n int) {
			worker, err := NewWorker(limiter, startURL, finalURL)
			if err != nil {
				log.Printf("Process #%d: %s", n, err)
				return
			}
			err = worker.PassingTest()
			if err == nil {
				log.Printf("Process #%d: Test successfully passed", n)
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
