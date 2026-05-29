package admin

import (
	"net/http"
	"time"

	"github.com/ZyphenSVC/SecureOps/backend/internal/audit"
	"github.com/ZyphenSVC/SecureOps/backend/internal/auth"
	"github.com/ZyphenSVC/SecureOps/backend/internal/httpx"
	"github.com/ZyphenSVC/SecureOps/backend/internal/users"
)

type Handler struct {
	users *users.Repository
	audit *audit.Repository
}

func NewHandler(users *users.Repository, auditRepo *audit.Repository) *Handler {
	return &Handler{
		users: users,
		audit: auditRepo,
	}
}

func (h *Handler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := h.audit.List(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "could not list audit logs")
		return
	}
	
	actorUserID, _ := auth.UserIDFromContext(r.Context())
	_ = h.audit.Record(r.Context(), &actorUserID, "admin.audit_logs.listed", "audit_log", nil, "{}")

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"audit_logs": logs,
	})
}

type userResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	allUsers, err := h.users.List(r.Context())
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "could not list users")
		return
	}

	response := make([]userResponse, 0, len(allUsers))
	actorUserID, _ := auth.UserIDFromContext(r.Context())
	_ = h.audit.Record(r.Context(), &actorUserID, "admin.users.listed", "user", nil, "{}")

	for _, user := range allUsers {
		response = append(response, userResponse{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.FullName,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"users": response,
	})
}
