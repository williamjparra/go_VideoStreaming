package handler

import (
	"encoding/json"
	"net/http"

	"github.com/williamjparra/go_VideoStreaming/internal/auth"
	"github.com/williamjparra/go_VideoStreaming/internal/repository"
)

type AuthHandler struct {
	UserRepo *repository.UserRepository
}

func NewAuthHandler(repo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{UserRepo: repo}
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string           `json:"token,omitempty"`
	User  *repository.User `json:"user,omitempty"`
	Error string           `json:"error,omitempty"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(AuthResponse{Error: "Método no permitido"})
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AuthResponse{Error: "Datos de registro inválidos"})
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AuthResponse{Error: "Error al procesar la contraseña"})
		return
	}

	user, err := h.UserRepo.CreateUser(req.Username, hash)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(AuthResponse{Error: "El usuario ya existe"})
		return
	}

	token, _ := auth.GenerateJWT(user.ID, user.Username)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(AuthResponse{Error: "Método no permitido"})
		return
	}

	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(AuthResponse{Error: "Cuerpo de solicitud inválido"})
		return
	}

	user, err := h.UserRepo.GetUserByUsername(req.Username)
	if err != nil || !auth.CheckPasswordHash(req.Password, user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(AuthResponse{Error: "Credenciales incorrectas"})
		return
	}

	token, err := auth.GenerateJWT(user.ID, user.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(AuthResponse{Error: "Error al generar el token"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AuthResponse{
		Token: token,
		User:  user,
	})
}
