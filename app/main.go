package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AlifRJ/Go_Auth/app/controller"
	"github.com/AlifRJ/Go_Auth/app/migration"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	// ******************** //
	// 		Load ENV		//
	// ********************	//

	_ = godotenv.Load()

	// ******************** //
	// 		Configure DB	//
	// ********************	//

	// Open database connection pool
	dsn:="postgres://"+os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@localhost:"+os.Getenv("DB_PORT")+"/"+os.Getenv("DB_DATABASE")+"?sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Error parsing connection string: %v\n", err)
	}
	defer db.Close()

	// Ping to verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Cannot connect to database: %v\n", err)
	}
	fmt.Println("Successfully connected to PostgreSQL!")

	// Migrate users table
	if err := migration.UserMigration(db); err != nil{
		log.Fatalf("Failed to create table: %v\n", err)
	} else{
		// seeder.UserSeeder(db)
	}
	
	// ******************** //
	// 	Listen & Serve HTTP	//
	// ********************	//

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/v1", func(r chi.Router) {

		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("PONG!"))
		})
		r.Get("/user", controller.GetAllUsers)
		r.Get("/user/{id}", controller.GetUser)
		r.Post("/user", controller.StoreUser)
		r.Post("/user/{id}", controller.UpdateUser)
		r.Delete("/user/{id}", controller.DeleteUser)
	})


	fmt.Println("Running on port: 8080")
	http.ListenAndServe(":8080", r)
}