package main

import (
	"net/http"

	"github.com/AlifRJ/Go_Auth/app/controller"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/v1", func(r chi.Router) {

		r.Get("/user", controller.GetAllUsers)
		r.Get("/user/{id}", controller.GetUser)
		r.Post("/user", controller.StoreUser)
		r.Post("/user/{id}", controller.UpdateUser)
		r.Delete("/user/{id}", controller.DeleteUser)
	})


	http.ListenAndServe(":8080", r)
}