package auth

import (
	"encoding/json"
	"net/http"

	"github.com/ZyphenSVC/SecureOps/backend/internal/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

type userResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	IsActive bool   `json:"is_active"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.FullName == "" {
		httpx.WriteError(w, http.StatusBadRequest, "email, password, and full_name are required")
		return
	}

	user, token, err := h.service.Register(r.Context(), req.Email, req.Password, req.FullName)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "could not register user")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, authResponse{
		Token: token,
		User: userResponse{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			IsActive: user.IsActive,
		},
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		httpx.WriteError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	user, token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, authResponse{
		Token: token,
		User: userResponse{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			IsActive: user.IsActive,
		},
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "unauthenticated")
		return
	}

	user, err := h.service.users.FindByID(r.Context(), userID)
	if err != nil {
		httpx.WriteError(w, http.StatusNotFound, "user not found")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, userResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		IsActive: user.IsActive,
	})
}
