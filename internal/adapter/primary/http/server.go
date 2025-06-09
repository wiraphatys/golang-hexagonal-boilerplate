package http

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"basedir/internal/adapter/primary/http/order"
	"basedir/internal/adapter/primary/http/product"
	core "basedir/internal/core/infra/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
)

type HttpServer struct {
	cfg           core.AppConfigProvider
	log           *zap.Logger
	fiberApp      *fiber.App
	apiBaseRouter fiber.Router
	basePath      string
}

type RouteGroup struct {
	order   *order.OrderHttpHandler
	product *product.ProductHttpHandler
}

func NewRouteGroup(order *order.OrderHttpHandler, product *product.ProductHttpHandler) *RouteGroup {
	return &RouteGroup{
		order:   order,
		product: product,
	}
}

func NewHttpServer(cfg core.AppConfigProvider, zapLogger *zap.Logger, baseApiPrefix string) *HttpServer {
	validateArguments(cfg, zapLogger, &baseApiPrefix)

	app := fiber.New(fiber.Config{
		AppName: cfg.GetServerName(),
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} - ${method} ${path}\n",
		TimeFormat: "2006/01/02 15:04:05",
		TimeZone:   "Asia/Bangkok",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Origin,X-PINGOTHER,Accept,Authorization,Content-Type,X-CSRF-Token",
		ExposeHeaders:    "Link",
		AllowCredentials: false,
		MaxAge:           300,
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 60 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"message": "Too many requests, please try again later.",
			})
		},
	}))

	apiGroup := app.Group(baseApiPrefix)
	zapLogger.Info("Fiber HTTP server core initialized with middleware.", zap.String("baseApiPrefix", baseApiPrefix))

	return &HttpServer{
		cfg:           cfg,
		log:           zapLogger,
		fiberApp:      app,
		apiBaseRouter: apiGroup,
		basePath:      baseApiPrefix,
	}
}

func (s *HttpServer) Start() {
	serverAddr := fmt.Sprintf("%s:%s", s.cfg.GetServerHost(), s.cfg.GetServerPort())
	s.log.Info("Attempting to start HTTP server...", zap.String("address", serverAddr))

	go func() {
		if err := s.fiberApp.Listen(serverAddr); err != nil && err.Error() != "http: Server closed" {
			s.log.Error("Failed to start HTTP server listener", zap.Error(err))
		}
	}()

	s.log.Info("HTTP server listener started.", zap.String("address", serverAddr))
	s.gracefulShutdown()
}

func (s *HttpServer) SetupRoute(routeGroup *RouteGroup) {
	if routeGroup.order == nil || routeGroup.product == nil {
		s.log.Error("Failed to set up route")
	}

	order := routeGroup.order
	product := routeGroup.product

	s.registerAPIGroup("/orders", func(orderRouter fiber.Router) {
		orderRouter.Get("/:id", order.GetOrderByID)
	})

	s.registerAPIGroup("/products", func(productRouter fiber.Router) {
		productRouter.Get("/:id", product.GetProductByID)
	})
}

func (s *HttpServer) gracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	s.log.Info("Received shutdown signal", zap.String("signal", sig.String()))
	s.log.Info("Gracefully shutting down HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := s.fiberApp.ShutdownWithContext(shutdownCtx); err != nil {
		s.log.Error("Error during server shutdown", zap.Error(err))
	} else {
		s.log.Info("HTTP server shut down gracefully.")
	}

	s.log.Info("Cleanup finished. Exiting.")
}

func (s *HttpServer) registerAPIGroup(subPrefix string, groupRegistrar func(router fiber.Router)) {
	if !strings.HasPrefix(subPrefix, "/") && subPrefix != "" {
		subPrefix = "/" + subPrefix
	}
	fullPrefix := s.basePath + subPrefix
	if s.basePath == "/" && subPrefix == "" {
		fullPrefix = "/"
	} else if s.basePath == "/" && strings.HasPrefix(subPrefix, "/") {
		fullPrefix = subPrefix
	} else if s.basePath != "" && subPrefix == "" {
		fullPrefix = s.basePath
	}

	group := s.apiBaseRouter.Group(subPrefix)
	groupRegistrar(group)
	s.log.Info("Registered API group", zap.String("fullPrefix", fullPrefix))
}

func (s *HttpServer) addHttpRoute(method, relativePath string, handler fiber.Handler) fiber.Router {
	if !strings.HasPrefix(relativePath, "/") && relativePath != "" {
		relativePath = "/" + relativePath
	}
	fullPath := s.basePath + relativePath
	if s.basePath == "/" && relativePath == "" {
		fullPath = "/"
	} else if s.basePath == "/" && strings.HasPrefix(relativePath, "/") {
		fullPath = relativePath
	} else if s.basePath != "" && relativePath == "" {
		fullPath = s.basePath
	}

	s.log.Info("Adding HTTP route", zap.String("method", method), zap.String("fullPath", fullPath))
	return s.apiBaseRouter.Add(method, relativePath, handler)
}

func validateArguments(cfg core.AppConfigProvider, zapLogger *zap.Logger, baseApiPrefix *string) {
	if cfg == nil {
		zapLogger.Fatal("Server configuration is missing for HTTP server initialization")
	}
	if zapLogger == nil {
		zapLogger.Fatal("Logger instance is missing for HTTP server initialization")
	}
	if baseApiPrefix != nil {
		if *baseApiPrefix == "" {
			zapLogger.Warn("baseApiPrefix is empty, API routes will be registered at the root.")
		} else if !strings.HasPrefix(*baseApiPrefix, "/") {
			*baseApiPrefix = "/" + *baseApiPrefix
			zapLogger.Warn("baseApiPrefix did not start with '/', prepended it.", zap.String("newPrefix", *baseApiPrefix))
		}
	}
}
