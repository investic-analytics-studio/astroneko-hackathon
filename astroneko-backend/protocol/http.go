package protocol

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"astroneko-backend/configs"
	"astroneko-backend/docs"
	"astroneko-backend/internal/routes"
	"astroneko-backend/pkg/databases/gorm"
	"astroneko-backend/pkg/logger"
	"astroneko-backend/pkg/middleware"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

const signalBufferSize = 1

type Server struct {
	app         *fiber.App
	config      *configs.Config
	zapLogger   *zap.Logger
	csrfManager *middleware.CSRFManager
}

func NewServer() *Server {
	zapLogger, _ := logger.NewZapLogger()
	app := fiber.New(fiber.Config{
		BodyLimit:                 10 * 1024 * 1024, // 10MB
		ReadBufferSize:            16384,            // 16KB - increased for larger headers
		WriteBufferSize:           16384,            // 16KB - increased for larger headers
		ReadTimeout:               time.Second * 30,
		WriteTimeout:              time.Second * 30,
		IdleTimeout:               time.Second * 120,
		DisableKeepalive:          false,
		DisableDefaultDate:        false,
		DisableDefaultContentType: false,
		DisableHeaderNormalizing:  false,
		DisableStartupMessage:     false,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://staging.astro-boxing-next.pages.dev,https://astro-boxing-next.pages.dev,http://localhost:3000,http://localhost:5173,http://localhost:5174,http://localhost:5175,https://astroneko.com,https://astroneko.net,https://staging.luckycat-frontend.pages.dev,https://luckycat-frontend.pages.dev,https://fix-login.luckycat-frontend.pages.dev,https://dev.astroneko-crm-frontend.pages.dev,https://staging.astroneko-crm-frontend.pages.dev,https://astroneko-crm-frontend.pages.dev,https://astrofight.ai",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Csrf-Token,X-Requested-With",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	app.Use(logger.ZapLoggerMiddleware(zapLogger))
	app.Use(logger.ZapRecoveryMiddleware(zapLogger))

	// Add middleware to handle header size issues
	app.Use(func(c *fiber.Ctx) error {
		// Check for extremely large headers and reject early
		headers := c.GetReqHeaders()
		if len(headers) > 100 { // Limit number of headers
			return c.Status(fiber.StatusRequestHeaderFieldsTooLarge).JSON(fiber.Map{
				"error": "Too many request headers",
			})
		}

		// Check individual header sizes
		for key, value := range headers {
			if len(key) > 8192 || len(value) > 8192 { // 8KB limit per header
				return c.Status(fiber.StatusRequestHeaderFieldsTooLarge).JSON(fiber.Map{
					"error": "Request header field too large",
				})
			}
		}

		return c.Next()
	})

	// Security middleware
	securityConfig := middleware.SecurityConfig{
		CSRFSecret:    middleware.GenerateSecretKey(),
		CookieSecret:  middleware.GenerateSecretKey(),
		SessionSecret: middleware.GenerateSecretKey(),
		Domain:        "", // Will be set based on environment
	}

	// Rate limiting (applied first to prevent abuse)
	app.Use(middleware.SetupRateLimitMiddleware())

	// Helmet for security headers
	app.Use(middleware.SetupHelmetMiddleware())

	// Additional security headers
	app.Use(middleware.SetupSecureHeadersMiddleware())

	// XSS Protection
	app.Use(middleware.SetupXSSProtectionMiddleware())

	// Encrypt cookies
	app.Use(middleware.SetupEncryptCookieMiddleware(securityConfig))

	// Initialize CSRF manager
	csrfManager := middleware.NewCSRFManager(securityConfig)

	// Setup CSRF routes
	middleware.SetupCSRFRoutes(app, csrfManager)

	return &Server{
		app:         app,
		config:      nil,
		zapLogger:   zapLogger,
		csrfManager: csrfManager,
	}
}

func (s *Server) Initialize(configPath string) error {
	configs.InitViper(configPath)
	s.config = configs.GetViper()
	return nil
}

func (s *Server) setupDatabase() (*gorm.DB, error) {
	if s.config.App.Env == "local" {
		return gorm.ConnectToPostgreSQL(
			s.config.Postgres.Host,
			s.config.Postgres.Port,
			s.config.Postgres.Username,
			s.config.Postgres.Password,
			s.config.Postgres.DbName,
			s.config.Postgres.SSLMode,
		)
	}
	return gorm.ConnectToCloudSQL(
		s.config.Postgres.InstanceConnectionName,
		s.config.Postgres.Username,
		s.config.Postgres.Password,
		s.config.Postgres.DbName,
	)
}

func (s *Server) setupGracefulShutdown(db *gorm.DB) {
	c := make(chan os.Signal, signalBufferSize)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		for range c {
			log.Println("Graceful shutdown initiated...")
			gorm.DisconnectPostgres(db.Postgres)

			if err := s.app.Shutdown(); err != nil {
				log.Printf("Error during shutdown: %v\n", err)
			}
		}
	}()
}

func (s *Server) setupRoutes(db *gorm.DB) {
	// Setup all routes using the routes package
	routes.SetupAllRoutes(s.app, db, s.zapLogger)
}

func (s *Server) swagger() {
	docs.SwaggerInfo.Title = "Astroneko Backend API"
	docs.SwaggerInfo.Description = "REST API for Astroneko Backend"
	docs.SwaggerInfo.Version = "1.0"

	if s.config != nil {
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%s", s.config.App.Port)
	} else {
		docs.SwaggerInfo.Host = "localhost:8080"
	}

	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

func ServeHTTP() error {
	server := NewServer()

	// Initialize configuration - will load prod-config.yml
	if err := server.Initialize("./configs"); err != nil {
		log.Printf("Warning: Could not load config: %v, using default configuration", err)
	}

	server.swagger()

	// Try to setup database, but don't fail if it's not available
	dbConGorm, err := server.setupDatabase()
	if err != nil {
		log.Printf("Warning: Database not available: %v", err)
	} else {
		server.setupGracefulShutdown(dbConGorm)
	}

	// Setup routes (pass database connection and logger)
	server.setupRoutes(dbConGorm)

	port := "8080"
	if server.config != nil && server.config.App.Port != "" {
		port = server.config.App.Port
	}

	log.Printf("Starting Astroneko Backend API server...")
	log.Printf("Listening on port: %s", port)

	return server.app.Listen(":" + port)
}
