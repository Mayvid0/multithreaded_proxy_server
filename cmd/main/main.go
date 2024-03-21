package main

import (
	"log"
	"net/http"

	"github.com/Mayvid0/multithreaded_proxy_server/internal/proxy"
)

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", proxy.ForwardProxy)
	mux.HandleFunc("/reverse", proxy.ForwardProxy)
	return mux
}

func main() {
	// Create a new HTTP server with the handleRequest function as the handler
	router := routes()
	server := http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	// test code to check concurrent write to log file
	// test.TestWriting()

	// Start the server and log any errors
	log.Println("Starting proxy server on :8000")
	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Error starting proxy server: ", err)
	}

}
