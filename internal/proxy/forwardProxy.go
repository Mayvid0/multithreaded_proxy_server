package proxy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	logs "github.com/Mayvid0/proxy_server/internal/cacheAndLog"
)

var customTransport = http.DefaultTransport

func ForwardProxy(w http.ResponseWriter, r *http.Request) {
	// Create a context with a timeout of, for example, 10 seconds
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel() // Ensure context cancellation to release resources

	// Create a new HTTP request with the same method, URL, and body as the original request
	targetURL := r.URL
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
	// Copy the body of the proxy response to the original response
	io.Copy(w, resp.Body)
}
