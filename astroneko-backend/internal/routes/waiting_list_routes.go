package routes

import (
	"astroneko-backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupWaitingListRoutes(api fiber.Router, waitingListHandler *handlers.WaitingListHTTPHandler) {
	// Waiting list routes
	waitingList := api.Group("/waiting-list")
	waitingList.Post("/join", waitingListHandler.JoinWaitingList)
	waitingList.Post("/check", waitingListHandler.IsInWaitingListByEmail)
}
