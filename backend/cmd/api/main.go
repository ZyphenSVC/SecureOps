package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ZyphenSVC/SecureOps/backend/internal/admin"
	"github.com/ZyphenSVC/SecureOps/backend/internal/audit"
	"github.com/ZyphenSVC/SecureOps/backend/internal/auth"
	"github.com/ZyphenSVC/SecureOps/backend/internal/config"
	"github.com/ZyphenSVC/SecureOps/backend/internal/db"
	"github.com/ZyphenSVC/SecureOps/backend/internal/rbac"
	"github.com/ZyphenSVC/SecureOps/backend/internal/users"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer pool.Close()

	userRepo := users.NewRepository(pool)
	rbacRepo := rbac.NewRepository(pool)
	auditRepo := audit.NewRepository(pool)
	authService := auth.NewService(userRepo, rbacRepo, auditRepo, cfg.JWTSecret, cfg.BcryptCost)
	authHandler := auth.NewHandler(authService)
	adminHandler := admin.NewHandler(userRepo, auditRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/admin/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		authService.RequirePermission("users:read", adminHandler.ListUsers)(w, r)
	})

	mux.HandleFunc("/admin/audit-logs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		authService.RequirePermission("audit_logs:read", adminHandler.ListAuditLogs)(w, r)
	})

	mux.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		authHandler.Register(w, r)
	})

	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		authHandler.Login(w, r)
	})

	mux.HandleFunc("/auth/me", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		authService.Authenticate(authHandler.Me)(w, r)
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
			"env":    cfg.AppEnv,
			"build":  "auth-routes-v1",
		})
	})

	mux.HandleFunc("GET /db-health", func(w http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{
				"status": "database unavailable",
			})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"status": "database ok",
		})
	})

	addr := ":" + cfg.Port
	log.Printf("SecureOps API listening on %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("write json error: %v", err)
	}
}
