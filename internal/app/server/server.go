package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zkKYC-backend/internal/app/config"
	"zkKYC-backend/internal/app/handlers"

	"github.com/go-chi/chi/v5"
)

// Server stores *http.Server and *handlers.ZkKYCHandler
type Server struct {
	*http.Server
	sh *handlers.ZkKYCHandler
}

// Create new server
func NewServer(cfg config.Config) Server {

	h := handlers.NewShortenerHandler(cfg)

	r := chi.NewRouter()
	r.Use(handlers.GzipHandle)
	r.Use(handlers.UnpackHandle)

	r.Route("/", func(r chi.Router) {
		r.Post("/api/user", h.APICreateUser)
		r.Get("/api/user/{eth}", h.APIGetExitingUser)
	})

	s := Server{Server: &http.Server{}, sh: h}
	s.Addr = cfg.ServerAddress
	s.Handler = r

	return s
}

// Entrypoint for server
func (s *Server) Start(cfg config.Config) {

	defer s.sh.DB.Close()

	idleConnsClosed := make(chan struct{})

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {

		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		if err := s.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()

	if err := s.ListenAndServe(); err != nil {
		log.Println(err)
	}

	<-idleConnsClosed
	log.Println("Server Shutdown gracefully")

}
