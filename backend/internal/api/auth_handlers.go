package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/auth"
	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/store"
)

type ctxKey int

const userCtxKey ctxKey = iota

const sessionCookie = "blo_session"

func (h *Handler) setSession(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(h.auth.TTL().Seconds()),
	})
}

func (h *Handler) clearSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

func toUserResponse(u *store.User) userResponse {
	return userResponse{ID: u.ID, Username: u.Username, Role: u.Role}
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	u, err := h.store.GetUserByUsername(r.Context(), req.Username)
	if err != nil || !auth.CheckPassword(u.PasswordHash, req.Password) {
		writeErr(w, http.StatusUnauthorized, "invalid username or password")
		return
	}
	token, _, err := h.auth.Issue(u.ID, u.Username, u.Role)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not issue session")
		return
	}
	h.setSession(w, token)
	writeJSON(w, http.StatusOK, toUserResponse(u))
}

func (h *Handler) logout(w http.ResponseWriter, _ *http.Request) {
	h.clearSession(w)
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	u := h.userFromRequest(r)
	if u == nil {
		writeErr(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	writeJSON(w, http.StatusOK, toUserResponse(u))
}

// userFromRequest resolves the session cookie to a user, or nil.
func (h *Handler) userFromRequest(r *http.Request) *store.User {
	c, err := r.Cookie(sessionCookie)
	if err != nil {
		return nil
	}
	claims, err := h.auth.Parse(c.Value)
	if err != nil {
		return nil
	}
	u, err := h.store.GetUserByID(r.Context(), claims.Subject)
	if err != nil {
		return nil
	}
	return u
}

// requireAuth rejects unauthenticated requests and stashes the user in context.
func (h *Handler) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := h.userFromRequest(r)
		if u == nil {
			writeErr(w, http.StatusUnauthorized, "authentication required")
			return
		}
		ctx := context.WithValue(r.Context(), userCtxKey, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// requireAdmin must run after requireAuth; it enforces the admin role.
func (h *Handler) requireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := currentUser(r.Context())
		if u == nil || u.Role != auth.RoleAdmin {
			writeErr(w, http.StatusForbidden, "admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func currentUser(ctx context.Context) *store.User {
	u, _ := ctx.Value(userCtxKey).(*store.User)
	return u
}

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// createUser is admin-only account provisioning.
func (h *Handler) createUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Username == "" || req.Password == "" {
		writeErr(w, http.StatusBadRequest, "username and password are required")
		return
	}
	role := req.Role
	if role != auth.RoleAdmin {
		role = auth.RoleUser
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not hash password")
		return
	}
	u, err := h.store.CreateUser(r.Context(), req.Username, hash, role)
	if err != nil {
		writeErr(w, http.StatusConflict, "could not create user (username may already exist)")
		return
	}
	writeJSON(w, http.StatusCreated, toUserResponse(u))
}
