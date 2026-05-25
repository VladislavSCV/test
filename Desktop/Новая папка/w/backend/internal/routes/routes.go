package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"w/backend/internal/handlers"
	"w/backend/internal/middleware"
)

func Register(app *fiber.App, db *gorm.DB) {
	h := &handlers.Handlers{DB: db}
	api := app.Group("/api")
	api.Post("/register", h.Register)
	api.Post("/login", h.Login)
	auth := api.Group("", middleware.RequireAuth(db))
	auth.Get("/me", h.Me)
	auth.Post("/records", h.CreateRecord)
	auth.Get("/records/my", h.MyRecords)
	auth.Post("/reviews", h.CreateReview)
	admin := api.Group("/admin", middleware.RequireAuth(db), middleware.RequireAdmin())
	admin.Get("/records", h.AdminRecords)
	admin.Patch("/records/:id/status", h.UpdateStatus)
}
