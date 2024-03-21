package test

import (
	"os/exec"
	"sync"
	"testing"
)

func TestWriting(t *testing.T) {
	// Number of goroutines to simulate concurrent log writes
	numGoroutines := 20

	// Wait group to synchronize goroutines
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Start goroutines to execute the curl command concurrently

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			// make sure to replace the port number and the server you want to request to
			cmd := exec.Command("curl", "-x", "http://localhost:8000", "http://localhost:8080/hello")
			err := cmd.Run()
			if err != nil {
				t.Errorf("Error executing curl command in Goroutine %d: %v\n", id, err)
				return
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}
