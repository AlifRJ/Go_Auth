package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/AlifRJ/Go_Auth/app/repository"
	"github.com/AlifRJ/Go_Auth/app/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

type UserHandler struct {
	service *service.UserService
	accessTokenAuth *jwtauth.JWTAuth
}

func NewUserHandler(s *service.UserService, accessTokenAuth * jwtauth.JWTAuth) *UserHandler {
	return &UserHandler{service: s, accessTokenAuth: accessTokenAuth}
}

// Routes Definition
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/v1", func(r chi.Router) {
		r.Group(func(r chi.Router){
			r.Use(jwtauth.Verifier(h.accessTokenAuth))
			r.Use(jwtauth.Authenticator(h.accessTokenAuth))
			r.Get("/users", h.GetAll)
			r.Get("/users/{id}", h.GetByID)
			r.Post("/users", h.Create)
			r.Put("/users/{id}", h.UpdateUser)
			r.Delete("/users/{id}", h.DeleteUser)
			r.Delete("/users-delete/{id}", h.PermanentlyDeleteUser)
		})
	})
}

// Get All Users
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	// Get Users
	users, err := h.service.GetAll(r.Context(), limit, offset)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
		return
	}

	respondJSON(w, http.StatusOK, users)
}

// Get User by ID
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID parameter"})
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// Register new User
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name  string `json:"name"`
		Username  string `json:"username"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	user, err := h.service.RegisterUser(r.Context(), body.Name, body.Username, body.Email, body.Password)
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

// Update User
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Request Params
	id, err := parseUintParam(r, "id")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID parameter"})
		return
	}

	var body struct {
		Name  string `json:"name"`
		Username  string `json:"username"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	user, err := h.service.UpdateUser(r.Context(), id, body.Name, body.Username, body.Email, body.Password)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, user)
}

// Soft Delete User
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request){
	id, err := parseUintParam(r, "id")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID parameter"})
		return
	}

	if err := h.service.DeleteUser(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Permanent Delete User
func (h *UserHandler) PermanentlyDeleteUser(w http.ResponseWriter, r *http.Request){
	id, err := parseUintParam(r, "id")
	if err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid ID parameter"})
		return
	}

	if err := h.service.DeleteUserPermanently(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			respondJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
			return
		}
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Response & Param Helpers
func respondJSON(w http.ResponseWriter, code int, payload any) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func parseUintParam(r *http.Request, key string) (uint, error) {
	valStr := chi.URLParam(r, key)
	val, err := strconv.ParseUint(valStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}