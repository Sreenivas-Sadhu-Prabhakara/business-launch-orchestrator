// Package api exposes the orchestrator over HTTP (chi router).
package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/auth"
	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/orchestrator"
	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/store"
)

// Handler bundles the dependencies the HTTP handlers need.
type Handler struct {
	svc        *orchestrator.Service
	store      *store.Store
	auth       *auth.Service
	corsOrigin string
}

// New builds the chi router with all routes mounted.
func New(svc *orchestrator.Service, st *store.Store, authSvc *auth.Service, corsOrigin string) http.Handler {
	h := &Handler{svc: svc, store: st, auth: authSvc, corsOrigin: corsOrigin}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(h.cors)

	r.Get("/healthz", h.health)

	r.Route("/api/v1", func(r chi.Router) {
		// Public auth endpoints.
		r.Post("/auth/login", h.login)
		r.Post("/auth/logout", h.logout)
		r.Get("/auth/me", h.me)

		// Everything else requires a valid session.
		r.Group(func(r chi.Router) {
			r.Use(h.requireAuth)

			r.With(h.requireAdmin).Get("/auth/users", h.listUsers)
			r.With(h.requireAdmin).Post("/auth/users", h.createUser)
			r.With(h.requireAdmin).Patch("/auth/users/{id}/role", h.updateUserRole)
			r.With(h.requireAdmin).Delete("/auth/users/{id}", h.deleteUser)

			r.Get("/countries", h.listCountries)
			r.Get("/countries/{code}/plan", h.countryPlan)

			r.Route("/businesses", func(r chi.Router) {
				r.Post("/", h.createBusiness)
				r.Get("/", h.listBusinesses)
				r.Get("/{id}", h.getBusiness)
				r.Get("/{id}/steps", h.getSteps)
				r.Post("/{id}/advance", h.advance)
				r.Post("/{id}/run", h.runAll)
			})
		})
	})

	return r
}

// cors is a permissive CORS middleware suitable for the bundled frontend.
func (h *Handler) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", h.corsOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
