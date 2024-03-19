package test

import (
	"fmt"
	"sync"
	"time"

	logs "github.com/Mayvid0/proxy_server/internal/AccessLog"
)

func TestWriting() {

	// Number of goroutines to simulate concurrent log writes
	numGoroutines := 20

	// Wait group to synchronize goroutines
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Start goroutines to simulate concurrent log writes
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			logs.WriteLogToFile(fmt.Sprintf("Log message from Goroutine %d at %s\n", id, time.Now()))
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

}
