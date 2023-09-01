package passing_webtest

import (
	"log"
	"sync"
)

func Runner(qtyOfThreads int, startURL, finalURL string) {
	wg := sync.WaitGroup{}
	wg.Add(qtyOfThreads)

	for i := 0; i < qtyOfThreads; i++ {
		go func(n int) {
			result := PassingTest(startURL, finalURL)
			log.Printf("Process #%d: ", n)
			if result == nil {
				log.Println("Test successfully passed")
			} else {
				log.Println(result)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
