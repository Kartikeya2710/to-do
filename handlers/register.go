package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Kartikeya2710/to-do/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.logger.Printf("Error decoding the user from request: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)

		return
	}

	filter := bson.D{{Key: "email", Value: user.Email}}
	var resultUser models.User

	err := h.collection.FindOne(r.Context(), filter).Decode(&resultUser)
	// we found a user with the same email
	if err == nil {
		h.logger.Printf("User with email %s found in database", user.Email)

		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "EmailAlreadyRegistered",
			"message": "An account with this email already exists. Please log in or reset your password.",
		})

		return
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		h.logger.Printf("Error finding user with email: %s", user.Email)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Printf("Error hashing password: %v", err)
		http.Error(w, "Invalid password", http.StatusBadRequest)

		return
	}

	user.Password = string(hashedPassword)

	_, err = h.collection.InsertOne(r.Context(), user)
	if err != nil {
		h.logger.Printf("Error inserting user in database: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)

		return
	}

	token, err := h.jwtManager.GenerateJWT(user.Email)
	if err != nil {
		h.logger.Printf("Error generating JWT: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)

		return
	}

	response := map[string]string{
		"message": "User registered successfully",
		"token":   token,
	}
	json.NewEncoder(w).Encode(response)
}
