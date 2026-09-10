package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AlifRJ/Go_Auth/app/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

type UserHandler struct {
	service *service.UserService
	tokenAuth *jwtauth.JWTAuth
}

func NewUserHandler(s *service.UserService, tokenAuth * jwtauth.JWTAuth) *UserHandler {
	return &UserHandler{service: s, tokenAuth: tokenAuth}
}

// Routes
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/v1", func(r chi.Router) {

		// Public Routes
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("PONG!"))
		})
		r.Post("/login", h.Login) 

		// Protected Routes
		r.Group(func(r chi.Router){
			r.Use(jwtauth.Verifier(h.tokenAuth))
			r.Use(jwtauth.Authenticator(h.tokenAuth))

			r.Get("/users", h.GetAll)
			r.Get("/users/{id}", h.GetByID)
			r.Post("/users", h.Create)
			r.Put("/users/{id}", h.UpdateUser)
			r.Delete("/users/{id}", h.DeleteUser)
			r.Delete("/users-delete/{id}", h.PermanentlyDeleteUser)
		})
	})
}

// Login Handler
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	// Perform credential check
	var body struct {
		Identity  string `json:"identity"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if body.Identity == "" || body.Password == "" {
        http.Error(w, "Identity and password are required", http.StatusBadRequest)
        return
    }

	// Get the user
	user, err := h.service.Login(r.Context(), body.Identity, body.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create the claims map
	claims := map[string]any{
		"id": user.ID,
		"name": user.Name,
		"username": user.Username,
		"email": user.Email,
	}

	// Set expiration using jwtauth helper functions (e.g., expires in 1 hour)
	jwtauth.SetExpiryIn(claims, 1*time.Hour)
	jwtauth.SetIssuedNow(claims)

	// Generate the signed JWT token string
	_, tokenString, err := h.tokenAuth.Encode(claims)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Return token string to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "token": tokenString,
    })
}

// Get All Users
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Request Query
	query := r.URL.Query()
	
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	// Get Users
	users, err := h.service.GetAll(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	// Return users
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// Get User by ID
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Request Params
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Get User
	user, err := h.service.GetByID(r.Context(), uint(id))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Return user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Register new User
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name  string `json:"name"`
		Username  string `json:"username"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	// Validate Request
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create User
	user, err := h.service.RegisterUser(r.Context(), body.Name, body.Username, body.Email, body.Password)
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Update User
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Request Params
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var body struct {
		Name  string `json:"name"`
		Username  string `json:"username"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	// Validate Request
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update User
	user, err := h.service.UpdateUser(r.Context(), uint(id), body.Name, body.Username, body.Email, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Soft Delete User
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request){
	// Request Params
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Soft Delete User
	if err := h.service.DeleteUser(r.Context(), uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Permanent Delete User
func (h *UserHandler) PermanentlyDeleteUser(w http.ResponseWriter, r *http.Request){
	// Request Params
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Hard Delete User
	if err := h.service.DeleteUserPermanently(r.Context(), uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}