package routes

import (
	"astroneko-backend/internal/handlers"
	"astroneko-backend/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(api fiber.Router, userHandler *handlers.UserHTTPHandler, authMiddleware *middleware.AuthMiddleware) {
	// User routes
	users := api.Group("/users")
	users.Post("/", userHandler.CreateUser)
	users.Get("/:id", userHandler.GetUserByID)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}
