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

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var queryUser models.User
	if err := json.NewDecoder(r.Body).Decode(&queryUser); err != nil {
		h.logger.Printf("Error decoding the user from request: %v", err)
		http.Error(w, "Invalid payload", http.StatusBadRequest)

		return
	}

	var resultUser models.User

	err := h.collection.FindOne(r.Context(), bson.D{{Key: "email", Value: queryUser.Email}}).Decode(&resultUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			h.logger.Printf("No users found with the email %s", queryUser.Email)
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)

			return
		}

		h.logger.Printf("Error finding user in database: %v", err)
		http.Error(w, "No user found", http.StatusBadRequest)

		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(resultUser.Password), []byte(queryUser.Password)); err != nil {
		h.logger.Printf("Mismatch between passwords: %v", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)

		return
	}

	token, err := h.jwtManager.GenerateJWT(resultUser.Email)
	if err != nil {
		h.logger.Printf("Error generating JWT: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)

		return
	}

	response := map[string]string{
		"message": "User logged in successfully",
		"token":   token,
	}

	json.NewEncoder(w).Encode(response)
}
