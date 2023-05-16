package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
	"vaccination-server/db"
	"vaccination-server/models"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	HASH_COST = 10
)

type SignUpResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type SignUpRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
}
type SignUpLoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type LoginResponse struct {
	Token string `json:"token"`
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var request = SignUpRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), HASH_COST)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var user = models.User{
		Email:    request.Email,
		Name:     request.Name,
		Password: string(hashedPassword),
	}
	createdUser := db.DB.Create(&user)
	errorCreate := createdUser.Error

	if errorCreate != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorCreate.Error()))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SignUpResponse{
		ID:    user.Id,
		Email: user.Email,
		Name:  user.Name,
	})

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request = SignUpLoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var user models.User
	db.DB.Where("email = ?", request.Email).First(&user)

	if user.Id == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid Credentials"))
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}
	claims := models.AppClaims{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Token: tokenString,
	})

}
