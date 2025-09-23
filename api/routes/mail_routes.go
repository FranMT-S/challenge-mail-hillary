package routes

import (
	"api/controllers"
	"api/middleware"
	"api/services"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// SetupMailRoutes configures the mail routes
func SetupMailRoutes(router *chi.Mux, db *gorm.DB) {

	mailService := services.NewEmailService(db)
	mailController := controllers.NewMailController(mailService)

	// Setup mail routes
	router.Route("/api", func(r chi.Router) {
		r.Route("/mails", func(r chi.Router) {
			r.Use(middleware.Pagination)
			r.Post("/search", mailController.SearchMails)
		})
	})
}
