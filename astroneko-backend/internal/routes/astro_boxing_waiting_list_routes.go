package routes

import (
	"astroneko-backend/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupAstroBoxingWaitingListRoutes(api fiber.Router, astroBoxingWaitingListHandler *handlers.AstroBoxingWaitingListHTTPHandler) {
	astroBoxingWaitingList := api.Group("/astro-boxing-waiting-list")

	astroBoxingWaitingList.Post("/join", astroBoxingWaitingListHandler.JoinAstroBoxingWaitingList)
	astroBoxingWaitingList.Get("/check", astroBoxingWaitingListHandler.IsInAstroBoxingWaitingListByEmail)
}
