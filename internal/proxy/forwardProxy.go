package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	logs "github.com/Mayvid0/multithreaded_proxy_server/internal/AccessLog"
	lru "github.com/Mayvid0/multithreaded_proxy_server/internal/lruCache"
)

var (
	customTransport = http.DefaultTransport
	lruCache        = lru.NewLRUCache(1000)
	mutex           sync.Mutex
)

func ForwardProxy(w http.ResponseWriter, r *http.Request) {
	// Create a context with a timeout of, for example, 10 seconds
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel() // Ensure context cancellation to release resources

	// Create a new HTTP request with the same method, URL, and body as the original request
	targetURL := r.URL

	mutex.Lock()
	if cachedResponse := lruCache.Get(targetURL.String()); cachedResponse != "" {
		searchStatus := "Success"
		if cachedResponse == "" {
			searchStatus = "Error"
		}

		// Create a reader from the cached response string
		stringResponseFromCache := strings.NewReader(cachedResponse + " (Brought from cache)\n")

		// Copy data from the reader to the response writer
		_, err := io.Copy(w, stringResponseFromCache)
		if err != nil {
			log.Fatal(err)
			return
		}

		accessLog := fmt.Sprintf("%s - accessed by ip: %s, date: %s, status: %s\n", targetURL, r.RemoteAddr, time.Now().Format("2006-01-02 15:04:05"), searchStatus)
		logs.WriteLogToFile(accessLog)
		mutex.Unlock()
	} else {
		mutex.Unlock()
		proxyReq, err := http.NewRequestWithContext(ctx, r.Method, targetURL.String(), r.Body)

		if err != nil {
			http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
			return
		}

		// Copy the headers from the original request to the proxy request
		for name, values := range r.Header {
			for _, value := range values {
				proxyReq.Header.Add(name, value)
			}
		}

		// Send the proxy request using the custom transport
		resp, err := customTransport.RoundTrip(proxyReq)

		if err != nil {
			http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
			return
		} else {
			fmt.Printf("sent proxy request to %s\n", r.URL)
		}
		defer resp.Body.Close()

		// Copy the headers from the proxy response to the original response
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		// Set the status code of the original response to the status code of the proxy response
		w.WriteHeader(resp.StatusCode)
		searchStatus := "Success"
		if resp.StatusCode >= 400 {
			searchStatus = "Error"
		}
		accessLog := fmt.Sprintf("%s - accessed by ip: %s, date: %s, status: %s\n", targetURL, r.RemoteAddr, time.Now().Format("2006-01-02 15:04:05"), searchStatus)
		logs.WriteLogToFile(accessLog)

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return
		}

		// Convert the byte slice to a string
		responseToBeStoredInCache := string(bodyBytes)

		mutex.Lock()
		lruCache.Put(targetURL.String(), responseToBeStoredInCache)
		mutex.Unlock()
		// Create a reader from the byte slice
		reader := strings.NewReader(responseToBeStoredInCache)

		// Copy the body of the proxy response to the original response
		_, err = io.Copy(w, reader)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}
