package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/auth"
	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/domain"
	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/orchestrator"
	"github.com/Sreenivas-Sadhu-Prabhakara/business-launch-orchestrator/backend/internal/store"
)

// loadOwned fetches a business and enforces access: admins see all, other users
// only their own. Returns false (and writes the response) if not allowed.
func (h *Handler) loadOwned(w http.ResponseWriter, r *http.Request, id string) (*domain.Business, bool) {
	b, err := h.store.GetBusiness(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		writeErr(w, http.StatusNotFound, "business not found")
		return nil, false
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}
	u := currentUser(r.Context())
	if u == nil || (u.Role != auth.RoleAdmin && b.OwnerID != u.ID) {
		// Hide existence from non-owners.
		writeErr(w, http.StatusNotFound, "business not found")
		return nil, false
	}
	return b, true
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// countryInfo is the public description of a supported jurisdiction.
type countryInfo struct {
	Code domain.Country       `json:"code"`
	Name string               `json:"name"`
	Plan []domain.PlannedStep `json:"plan"`
}

var countryNames = map[domain.Country]string{
	domain.CountryIndia:       "India",
	domain.CountryPhilippines: "Philippines",
	domain.CountryUS:          "United States",
}

func (h *Handler) listCountries(w http.ResponseWriter, _ *http.Request) {
	reg := h.svc.Registry()
	out := make([]countryInfo, 0, len(countryNames))
	for _, code := range []domain.Country{domain.CountryIndia, domain.CountryPhilippines, domain.CountryUS} {
		out = append(out, countryInfo{Code: code, Name: countryNames[code], Plan: reg.Plan(code)})
	}
	writeJSON(w, http.StatusOK, map[string]any{"countries": out})
}

func (h *Handler) countryPlan(w http.ResponseWriter, r *http.Request) {
	code := domain.Country(chi.URLParam(r, "code"))
	if !code.Valid() {
		writeErr(w, http.StatusNotFound, "unsupported country")
		return
	}
	writeJSON(w, http.StatusOK, countryInfo{
		Code: code, Name: countryNames[code], Plan: h.svc.Registry().Plan(code),
	})
}

// createBusinessRequest is the POST /businesses body.
type createBusinessRequest struct {
	Country         domain.Country `json:"country"`
	EntityType      string         `json:"entity_type"`
	LegalName       string         `json:"legal_name"`
	FounderName     string         `json:"founder_name"`
	FounderEmail    string         `json:"founder_email"`
	FounderPhone    string         `json:"founder_phone"`
	FounderIDNumber string         `json:"founder_id_number"`
	Address         domain.Address `json:"address"`
}

func (h *Handler) createBusiness(w http.ResponseWriter, r *http.Request) {
	var req createBusinessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.LegalName == "" || req.FounderName == "" {
		writeErr(w, http.StatusBadRequest, "legal_name and founder_name are required")
		return
	}
	if req.EntityType == "" {
		req.EntityType = defaultEntity(req.Country)
	}

	b := &domain.Business{
		Country:         req.Country,
		EntityType:      req.EntityType,
		LegalName:       req.LegalName,
		FounderName:     req.FounderName,
		FounderEmail:    req.FounderEmail,
		FounderPhone:    req.FounderPhone,
		FounderIDNumber: req.FounderIDNumber,
		Address:         req.Address,
	}

	if u := currentUser(r.Context()); u != nil {
		b.OwnerID = u.ID
	}

	if err := h.svc.CreateLaunch(r.Context(), b); err != nil {
		if errors.Is(err, orchestrator.ErrUnsupportedCountry) {
			writeErr(w, http.StatusBadRequest, "unsupported country (use IN, PH or US)")
			return
		}
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	h.respondWithDetail(w, r.Context(), http.StatusCreated, b.ID)
}

func (h *Handler) listBusinesses(w http.ResponseWriter, r *http.Request) {
	owner := "" // admin: all
	if u := currentUser(r.Context()); u != nil && u.Role != auth.RoleAdmin {
		owner = u.ID
	}
	items, err := h.store.ListBusinesses(r.Context(), owner, 50)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"businesses": items})
}

func (h *Handler) getBusiness(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, ok := h.loadOwned(w, r, id); !ok {
		return
	}
	h.respondWithDetail(w, r.Context(), http.StatusOK, id)
}

func (h *Handler) getSteps(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, ok := h.loadOwned(w, r, id); !ok {
		return
	}
	steps, err := h.store.GetSteps(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"steps": steps})
}

func (h *Handler) advance(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, ok := h.loadOwned(w, r, id); !ok {
		return
	}
	step, err := h.svc.AdvanceOne(r.Context(), id)
	switch {
	case errors.Is(err, orchestrator.ErrNoPendingSteps):
		writeJSON(w, http.StatusOK, map[string]any{"done": true, "message": "all steps completed"})
	case errors.Is(err, store.ErrNotFound):
		writeErr(w, http.StatusNotFound, "business not found")
	case err != nil:
		writeErr(w, http.StatusInternalServerError, err.Error())
	default:
		writeJSON(w, http.StatusOK, map[string]any{"step": step})
	}
}

func (h *Handler) runAll(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, ok := h.loadOwned(w, r, id); !ok {
		return
	}
	if _, err := h.svc.RunAll(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			writeErr(w, http.StatusNotFound, "business not found")
			return
		}
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondWithDetail(w, r.Context(), http.StatusOK, id)
}

// respondWithDetail returns {business, steps} for an id.
func (h *Handler) respondWithDetail(w http.ResponseWriter, ctx context.Context, status int, id string) {
	b, err := h.store.GetBusiness(ctx, id)
	if errors.Is(err, store.ErrNotFound) {
		writeErr(w, http.StatusNotFound, "business not found")
		return
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	steps, err := h.store.GetSteps(ctx, id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, status, map[string]any{"business": b, "steps": steps})
}

func defaultEntity(c domain.Country) string {
	switch c {
	case domain.CountryIndia:
		return "Private Limited Company"
	case domain.CountryUS:
		return "LLC"
	case domain.CountryPhilippines:
		return "Domestic Corporation"
	default:
		return "Company"
	}
}
