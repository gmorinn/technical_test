package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gmorinn/technical_test/internal/api"
	"github.com/gmorinn/technical_test/internal/config"
	"github.com/gmorinn/technical_test/internal/db"
	"github.com/gmorinn/technical_test/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/gmorinn/technical_test/docs"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 60 * time.Second

	shutdownTimeout = 15 * time.Second
)

func mwServerHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Server", "GM API")
		c.Next()
	}
}

//	@title			FizzBuzz API
//	@version		1.0
//	@description	A FizzBuzz REST API that returns the classic sequence for caller-supplied divisors and words, and reports the most requested combination.
//	@BasePath		/api/v1

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	database, err := db.NewDB(cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Database, cfg.TZ, cfg.Database.Port)
	if err != nil {
		return err
	}
	defer database.Close()

	if err := performMigrations(cfg); err != nil {
		return err
	}

	r := gin.New()
	r.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/health"))
	r.Use(gin.Recovery())
	r.Use(mwServerHeader())
	r.Use(middleware.GinContextToContextMiddleware())
	r.Use(cors.New(cors.Config{
		AllowOrigins: cfg.Cors,
		AllowMethods: []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type", "Origin", "Accept", "X-Requested-With"},
		// AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	h := api.New(cfg, database)
	h.RegisterRoutes(r)

	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.Port),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	serveErr := make(chan error, 1)
	go func() {
		log.Printf("🚀 Running REST API server at http://localhost:%d/api/v1", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
			return
		}
		serveErr <- nil
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
		stop()
	}

	log.Println("⏳ Shutting down, draining in-flight requests…")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	log.Println("✅ Shutdown complete")
	return nil
}
