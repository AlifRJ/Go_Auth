package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AlifRJ/Go_Auth/app/handler"
	"github.com/AlifRJ/Go_Auth/app/migration"
	"github.com/AlifRJ/Go_Auth/app/repository"
	"github.com/AlifRJ/Go_Auth/app/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Info: No .env file found, relying on system env")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Connection Pool Setup
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Error parsing database DSN: %v\n", err)
	}

	// Set pool limits
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	
	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Error creating database pool: %v\n", err)
	}
	defer db.Close()

	// Database Conn Check
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Cannot connect to database: %v\n", err)
	}
	fmt.Println("Successfully connected to PostgreSQL!")

	// Migrate users table
	if err := migration.UserMigration(ctx, db); err != nil{
		log.Fatalf("Failed to create table: %v\n", err)
	} else{
		// seeder.UserSeeder(db)
	}

	// Open Cache database connection
	rdb := redis.NewClient(&redis.Options{
		Addr:		os.Getenv("REDIS_HOST")+":"+os.Getenv("REDIS_PORT"),
		Password: 	os.Getenv("REDIS_PASSWORD"),
		DB:			0,
	})
	defer rdb.Close()

	// Redis Conn Check
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Cannot connect to Redis: %v\n", err)
	}
	fmt.Println("Successfully connected to Redis!")


	// JWT Setup
	accessSecret := os.Getenv("ACCESS_SECRET_KEY")
	refreshSecret := os.Getenv("REFRESH_SECRET_KEY")
	if accessSecret == "" || refreshSecret == "" {
		log.Fatal("JWT Secret Keys must be configured in environment variables")
	}
	accessTokenAuth := jwtauth.New("HS256", []byte(accessSecret), nil)
	refreshTokenAuth := jwtauth.New("HS256", []byte(refreshSecret), nil)
	
	// Dependency Injection
	userRepo := repository.NewPostgresUserRepository(db)
	cachedUserRepo := repository.NewCachedUserRepository(rdb, userRepo, 15*time.Minute)
	userService := service.NewUserService(cachedUserRepo)
	userHandler := handler.NewUserHandler(userService, accessTokenAuth)
	
	authRepo := repository.NewPostgresAuthRepository(db)
	cachedAuthRepo := repository.NewCachedAuthRepository(rdb, authRepo)
	authService := service.NewAuthService(cachedUserRepo, cachedAuthRepo)
	authHandler := handler.NewAuthHandler(authService, accessTokenAuth, refreshTokenAuth)
	
	// Chi Router & Middlewares
	r := chi.NewRouter()

	// Middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
  	}))

	r.Use(middleware.Logger)

	// Health Check Route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
    	w.Write([]byte("OK"))
	})

	// Routes Registry
	authHandler.RegisterRoutes(r)
	userHandler.RegisterRoutes(r)


	port := os.Getenv("APP_PORT")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		fmt.Printf("Server starting on port: %s\n", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server error: %v\n", err)
		}
	}()

	// Listen for shutdown signal (CTRL+C / SIGTERM)
	<-ctx.Done()
	fmt.Println("\nShutting down server gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	fmt.Println("Server stopped cleanly.")

	// fmt.Println("Running on port: 8080")
	// http.ListenAndServe(":8080", r)
}