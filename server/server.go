package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
	address    string
}

func NewServer(address string) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    address,
			Handler: nil,
		},
		address: address,
	}
}

func (s *Server) Handlefunc(path string, handler http.HandlerFunc) {
	http.HandleFunc(path, handler)
	s.httpServer.Handler = http.DefaultServeMux
}

func (s *Server) Start() error {
	log.Printf("Server listening on %s\n", s.address)

	// Setup channel to handle stopping the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v\n", err)
		}
	}()

	// Wait for signal...
	<-quit
	log.Printf("Shutting down server...")

	// Provide 5 sec timout to let things stop (not really needed since requests are minimal)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown gracefully - make sure all requests complete
	if err := s.httpServer.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown error: %v\n", err)
		return err
	}

	log.Println("Server stopped gracefully")
	return nil
}
