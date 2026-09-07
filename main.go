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
	userHandler := handler.NewUserHandler(userService)
	
	// Start Server
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	userHandler.RegisterRoutes(r)

	fmt.Println("Running on port: 8080")
	http.ListenAndServe(":8080", r)
}