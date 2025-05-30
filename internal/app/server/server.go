package server

import (
	"context"
	"github.com/go-chi/cors"
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

	h := handlers.NewZkKYCHandler(cfg)
	mh := handlers.NewMiddlewareHandler(cfg)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://localhost:3000",
			"http://frontend:3000",
			"https://dici24.ru/",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(mh.GzipHandle)
	r.Use(mh.UnpackHandle)

	r.Route("/api", func(r chi.Router) {

		r.Post("/user", h.APICreateUser)
		r.Get("/user/{eth}", h.APIGetUser)

		r.Post("/regulator/login", h.LoginHandler)

		r.Group(func(r chi.Router) {

			r.Use(mh.JwtAuthHandler)
			r.Get("/regulator/user/{eth}", h.APIGetUserForRegulator)
			
		})
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
