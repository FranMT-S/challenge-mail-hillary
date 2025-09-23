package server

import (
	"api/config"
	"api/middleware"
	"api/routes"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"gorm.io/gorm"
)

type Server struct {
	Router *chi.Mux
	DB     *gorm.DB
	Port   int
}

func NewServer(db *gorm.DB, port int) *Server {
	return &Server{
		Router: chi.NewRouter(),
		DB:     db,
		Port:   port,
	}
}

func (s *Server) Start() error {
	s.setupRoutes()

	log.Println("Server started on port:", s.Port)
	return http.ListenAndServe(":"+strconv.Itoa(s.Port), s.Router)
}

func (s *Server) setupRoutes() *Server {

	s.configCORS()
	s.setupMiddleware(s.Router)
	s.setupMailRoutes()
	return s
}

func (s *Server) setupMiddleware(router *chi.Mux) *Server {
	router.Use(chiMiddleware.Recoverer)
	router.Use(chiMiddleware.Logger)
	router.Use(middleware.LogError)
	return s
}

func (s *Server) setupMailRoutes() *Server {
	routes.SetupMailRoutes(s.Router, s.DB)
	return s
}

func (s *Server) configCORS() *Server {

	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{config.GetConfig().ClientHost},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	return s
}
