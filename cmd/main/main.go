package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var customTransport = http.DefaultTransport

func init() {
	// Here, you can customize the transport, e.g., set timeouts or enable/disable keep-alive
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Create a new HTTP request with the same method, URL, and body as the original request
	targetURL := r.URL
	// stringUrl := targetURL.String()
	// if strings.HasPrefix(stringUrl, "http://localhost:8080/") {
	// 	http.Error(w, "Can't access this URL", http.StatusBadRequest)
	// 	return
	// }

	proxyReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)

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
	writeLogTofile(accessLog)
	// Copy the body of the proxy response to the original response
	io.Copy(w, resp.Body)
}

func writeLogTofile(accessLog string) {

	fmt.Print("writing to log file\n")
	f, err := os.OpenFile("logs/access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := f.Write([]byte(accessLog)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	// Create a new HTTP server with the handleRequest function as the handler
	server := http.Server{
		Addr:    ":8000",
		Handler: http.HandlerFunc(handleRequest),
	}

	// Start the server and log any errors
	log.Println("Starting proxy server on :8000")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting proxy server: ", err)
	}

}
