package logs

import (
	"log"
	"os"
	"sync"
)

var logMutex sync.Mutex
var logChannel chan string
var logFile *os.File
var isInitialized bool

func init() {
	if isInitialized {
		return // Prevent reinitialization
	}
	logChannel = make(chan string, 1000) // Buffered channel with capacity 1000
	var err error
	logFile, err = os.OpenFile("logs/access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	go writeToLogFile()
	isInitialized = true
}

func WriteLogToFile(accessLog string) {
	if !isInitialized {
		return // Log system not initialized
	}
	logChannel <- accessLog
}

func writeToLogFile() {
	for accessLog := range logChannel {
		logMutex.Lock()
		if _, err := logFile.WriteString(accessLog); err != nil {
			log.Println("Error writing to log file:", err)
		}
		logMutex.Unlock()
	}
}

func CloseLogFile() {
	if logFile != nil {
		logFile.Close()
	}
	isInitialized = false
}
