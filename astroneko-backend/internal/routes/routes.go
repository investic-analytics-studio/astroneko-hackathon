package routes

import (
	"log"

	"astroneko-backend/internal/adapters"
	"astroneko-backend/internal/handlers"
	"astroneko-backend/internal/repositories"
	"astroneko-backend/internal/services"
	"astroneko-backend/pkg/databases/gorm"
	"astroneko-backend/pkg/firebase"
	"astroneko-backend/pkg/logger"
	"astroneko-backend/pkg/middleware"
	"astroneko-backend/pkg/validator"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// SetupAllRoutes initializes and sets up all application routes
func SetupAllRoutes(app *fiber.App, db *gorm.DB, zapLogger *zap.Logger) {
	// Initialize handlers
	healthHandler := handlers.NewHealthHTTPHandler()

	// Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// API v1
	api := app.Group("/v1/api")

	// Example endpoint
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to Astroneko Backend API",
			"version": "1.0.0",
		})
	})

	// Setup routes based on database availability
	if db != nil {
		setupApplicationRoutes(app, api, db, healthHandler, zapLogger)
	} else {
		// Basic health check without auth if no database
		app.Get("/health-check", healthHandler.HealthCheck)
	}
}

// setupApplicationRoutes sets up all application-specific routes with database dependencies
func setupApplicationRoutes(app *fiber.App, api fiber.Router, db *gorm.DB, healthHandler *handlers.HealthHTTPHandler, zapLogger *zap.Logger) {
	// Initialize logger
	appLogger := logger.NewDualLogger(zapLogger)

	// Initialize global Firebase client
	if err := firebase.InitFirebaseClient(appLogger); err != nil {
		log.Printf("Warning: Firebase initialization failed: %v", err)
	}
	firebaseClient := firebase.FirebaseClient
	if firebaseClient == nil {
		log.Printf("Warning: Firebase client is nil")
	}

	// Manual migrations are handled separately
	// Run migrations using: go run cmd/migrate/main.go

	// Initialize database adapter
	dbAdapter := adapters.NewGormAdapter(db.Postgres)

	// Initialize dependencies
	userRepo := repositories.NewUserRepository(dbAdapter)
	referralCodeRepo := repositories.NewReferralCodeRepository(dbAdapter)
	firebaseAdapter := firebase.NewFirebaseClientAdapter(firebaseClient)
	userService := services.NewUserService(userRepo, firebaseAdapter, "", appLogger, referralCodeRepo)
	referralCodeService := services.NewReferralCodeService(referralCodeRepo, userRepo, appLogger)
	userValidator := validator.New()
	userHandler := handlers.NewUserHTTPHandler(userService, referralCodeService, userValidator)

	// Waiting list dependencies
	waitingListRepo := repositories.NewWaitingListRepository(dbAdapter)
	waitingListService := services.NewWaitingListService(waitingListRepo, appLogger)
	waitingListValidator := validator.New()
	waitingListHandler := handlers.NewWaitingListHTTPHandler(waitingListService, waitingListValidator)

	// Agent dependencies
	agentRepo := repositories.NewAgentRepository()
	agentService := services.NewAgentService(agentRepo, appLogger)
	agentValidator := validator.New()
	agentHandler := handlers.NewAgentHTTPHandler(agentService, agentValidator)

	// Referral code dependencies
	referralCodeValidator := validator.New()
	referralCodeHandler := handlers.NewReferralCodeHTTPHandler(referralCodeService, referralCodeValidator)

	// CRM user dependencies
	crmUserRepo := repositories.NewCRMUserRepository(dbAdapter)
	crmUserService := services.NewCRMUserService(crmUserRepo, appLogger, "your-jwt-secret-key") // Use environment variable in production
	crmUserValidator := validator.New()
	crmUserHandler := handlers.NewCRMUserHTTPHandler(crmUserService, crmUserValidator)

	// User limit dependencies
	userLimitRepo := repositories.NewUserLimitRepository(dbAdapter)
	userLimitService := services.NewUserLimitService(userLimitRepo, userRepo, appLogger)
	userLimitValidator := validator.New()
	userLimitHandler := handlers.NewUserLimitHTTPHandler(userLimitService, userLimitValidator)

	// Astro boxing waiting list dependencies
	astroBoxingWaitingListRepo := repositories.NewAstroBoxingWaitingListRepository(dbAdapter)
	astroBoxingWaitingListService := services.NewAstroBoxingWaitingListService(astroBoxingWaitingListRepo, appLogger)
	astroBoxingWaitingListValidator := validator.New()
	astroBoxingWaitingListHandler := handlers.NewAstroBoxingWaitingListHTTPHandler(astroBoxingWaitingListService, astroBoxingWaitingListValidator)

	// Guest usage tracking dependencies
	guestUsageRepo := repositories.NewGuestUsageRepository(dbAdapter, appLogger)

	// History dependencies
	historyRepo := repositories.NewHistoryRepository(dbAdapter)
	historyService := services.NewHistoryService(historyRepo, appLogger)
	historyHandler := handlers.NewHistoryHTTPHandler(historyService)

	// Initialize middleware
	authMiddleware := middleware.NewFirebaseAuthMiddleware(firebaseClient, userService, appLogger)
	crmAuthMiddleware := middleware.NewCRMAuthMiddleware(crmUserService, appLogger)
	guestRateLimitMiddleware := middleware.NewGuestRateLimitMiddleware(guestUsageRepo, appLogger)

	// Setup all route modules
	SetupHealthRoutes(app, api, healthHandler, authMiddleware)
	SetupUserRoutes(api, userHandler, authMiddleware)
	SetupAuthRoutes(api, userHandler, authMiddleware)
	SetupWaitingListRoutes(api, waitingListHandler)
	SetupAgentRoutes(api, agentHandler, authMiddleware, guestRateLimitMiddleware)
	SetupCRMRoutes(api, crmUserHandler, userHandler, crmAuthMiddleware)
	SetupReferralCodeRoutes(api, referralCodeHandler, crmAuthMiddleware)
	SetupUserLimitRoutes(api, userLimitHandler, crmAuthMiddleware, authMiddleware)
	SetupAstroBoxingWaitingListRoutes(api, astroBoxingWaitingListHandler)
	SetupHistoryRoutes(api, historyHandler, authMiddleware)
}
