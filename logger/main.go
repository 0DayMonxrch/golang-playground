package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

type LogEntry struct {
	Level   string    `json:"level"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type Logger struct {
	file     *os.File
	logCh    chan LogEntry
	quitCh   chan struct{}
	wg       sync.WaitGroup
	dropLogs bool
}

// Constructor
func NewLogger(filename string, bufferSize int, dropLogs bool) (*Logger, error) {

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	logger := &Logger{
		file:     file,
		logCh:    make(chan LogEntry, bufferSize),
		quitCh:   make(chan struct{}),
		dropLogs: dropLogs,
	}

	logger.wg.Add(1)
	go logger.process()

	return logger, nil
}

func main() {

	logger, err := NewLogger("eventlog.json", 100, true)
	if err != nil {
		log.Fatal(err)
	}

	logger.Log("INFO", "Application Started")
	logger.Log("INFO", "Application event occurred")

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			memUsage := float64(m.Alloc) / float64(m.Sys) * 100
			logger.Log("METRIC",
				fmt.Sprintf("Current memory usage: %.2f%%", memUsage))
		}
	}()

	time.Sleep(10 * time.Second)

	logger.Log("INFO", "Application shutting down")
	logger.Close()
}

// Background worker
func (l *Logger) process() {
	defer l.wg.Done()

	encoder := json.NewEncoder(l.file)

	for {
		select {

		case entry := <-l.logCh:
			if err := encoder.Encode(entry); err != nil {
				fmt.Println("Failed to write log:", err)
			}

		case <-l.quitCh:
			// Drain remaining logs before exiting
			for len(l.logCh) > 0 {
				entry := <-l.logCh
				encoder.Encode(entry)
			}
			l.file.Close()
			return
		}
	}
}

// Public logging method
func (l *Logger) Log(level, message string) {

	entry := LogEntry{
		Level:   level,
		Message: sanitize(message),
		Time:    time.Now(),
	}

	if l.dropLogs {
		select {
		case l.logCh <- entry:
		default:
			// Drop log if buffer full
		}
	} else {
		l.logCh <- entry
	}
}

// Graceful shutdown
func (l *Logger) Close() {
	close(l.quitCh)
	l.wg.Wait()
}

// Basic log injection prevention
func sanitize(input string) string {
	return string([]byte(input))
}
