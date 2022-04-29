package main

import (
	"P2pSecuritySimulator/dataCollector"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				dataCollector.AppendRequireCertificateTime(int64(123456789))
				dataCollector.AppendAuthentificateTime(int64(123456789))
			}
		}()
	}

	wg.Wait()
}
