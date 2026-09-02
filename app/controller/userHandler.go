package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AlifRJ/Go_Auth/app/model"
	"github.com/AlifRJ/Go_Auth/app/payload"
	"github.com/AlifRJ/Go_Auth/app/utils"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

var users []model.User

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Get Users
	data := users

	// Return user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	// Get Params
	id := chi.URLParam(r, "id")

	// Get Users
	userID,_ := strconv.ParseUint(id, 10, 0)
	data, err := utils.GetUserByID(users, uint(userID))
	if err != nil {
		log.Printf("[SERVER ERROR] Failed to fetch user data: %v", err)
		http.Error(w, "User not Found!", http.StatusNotFound)
		return
	}

	// Return user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func StoreUser(w http.ResponseWriter, r *http.Request) {
	// Get request body
	var body model.User
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

	// Hash Request Password
	newPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[SERVER ERROR] Failed to fetch user data: %v", err)
		http.Error(w, "cannot encrypt password", http.StatusInternalServerError)
		return
	}

	// Auto Increment ID
	fmt.Println(len(users))
	newID := uint(len(users)+1)

	// Create User
	newUser := model.User{
		ID: newID,
		Name: body.Name,
		Username: body.Username,
		Email: body.Email,
		Password: string(newPassword),
		Created_at: time.Now(),
		Updated_at: nil,
		Deleted_at: nil,
	}

	// Store User
	users = append(users, newUser)

	// Return Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode("created")
}


func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get params
	id := chi.URLParam(r, "id")
	
	// Convert params to uid
	userID,err := strconv.ParseUint(id, 10, 64)
	if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

	// Get request body
	var body payload.UpdateUserPayload
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

	// Hash request password
	var newPassword string
	if body.Password != nil{
		password, err := bcrypt.GenerateFromPassword([]byte(*body.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("[SERVER ERROR] Failed to fetch user data: %v", err)
			http.Error(w, "Cannot encrypt password", http.StatusInternalServerError)
			return
		}
		newPassword = string(password)
	}

	// Update user
	if err := utils.UpdateUserByID(users, uint(userID), body, newPassword); err != nil{
		http.Error(w, "User not Found!", http.StatusNotFound)
		return
	}

	// Return Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode("updated")
}

func DeleteUser(w http.ResponseWriter, r *http.Request){
	id := chi.URLParam(r, "id")

	// Convert params to uid
	userID,err := strconv.ParseUint(id, 10, 64)
	if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

	// Delete user
    if err := utils.DeleteUserByID(users, uint(userID)); err != nil{
		http.Error(w, "User not Found!", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode("deleted")
}

func PermanentlyDeleteUser(w http.ResponseWriter, r *http.Request){
	id := chi.URLParam(r, "id")

	// Convert params to uid
	userID,err := strconv.ParseUint(id, 10, 64)
	if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

	index := -1
	// Find the index of the target struct
    for i, user := range users {
        if user.ID == uint(userID) {
            index = i
            break
        }
    }
	// Delete if found
	if index != -1 {
        users = append(users[:index], users[index+1:]...)
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode("deleted")
}