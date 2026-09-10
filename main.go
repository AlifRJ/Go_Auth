package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AlifRJ/Go_Auth/app/handler"
	"github.com/AlifRJ/Go_Auth/app/migration"
	"github.com/AlifRJ/Go_Auth/app/repository"
	"github.com/AlifRJ/Go_Auth/app/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	// Load env
	_ = godotenv.Load()

	ctx := context.Background()

	// Open database connection pool
	dsn:="postgres://"+os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@localhost:"+os.Getenv("DB_PORT")+"/"+os.Getenv("DB_DATABASE")+"?sslmode=disable"
	
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Error parsing connection string: %v\n", err)
	}
	defer db.Close()

	// Ping to verify connection
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

	// Dependency Injection
	userRepo := repository.NewPostgresUserRepository(db)
	userService := service.NewUserService(userRepo)
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("SECRET_KEY")), nil)
	userHandler := handler.NewUserHandler(userService, tokenAuth)
	
	// Start Server
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
  	}))

	userHandler.RegisterRoutes(r)

	fmt.Println("Running on port: 8080")
	http.ListenAndServe(":8080", r)
}