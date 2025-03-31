package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"one-client-view-2025tht/app/models"
)

// ApplicantHandler handles HTTP requests related to applicants
type ApplicantHandler struct {
	ApplicantRepo *models.ApplicantRepository
}

// NewApplicantHandler creates a new handler with the given repository
func NewApplicantHandler(repo *models.ApplicantRepository) *ApplicantHandler {
	return &ApplicantHandler{ApplicantRepo: repo}
}

// GetApplicants handles GET /api/applicants
func (h *ApplicantHandler) GetApplicants(w http.ResponseWriter, r *http.Request) {
	applicants, err := h.ApplicantRepo.GetAll()
	if err != nil {
		http.Error(w, "Failed to get applicants: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response objects
	var response []models.ApplicantResponse
	for _, a := range applicants {
		response = append(response, models.ApplicantResponse{
			Applicant: a,
			Household: a.Household,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetApplicant handles GET /api/applicants/{id}
func (h *ApplicantHandler) GetApplicant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	applicant, err := h.ApplicantRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to get applicant: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if applicant == nil {
		http.Error(w, "Applicant not found", http.StatusNotFound)
		return
	}

	response := models.ApplicantResponse{
		Applicant: *applicant,
		Household: applicant.Household,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateApplicant handles POST /api/applicants
func (h *ApplicantHandler) CreateApplicant(w http.ResponseWriter, r *http.Request) {
	var applicant models.Applicant
	err := json.NewDecoder(r.Body).Decode(&applicant)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Basic validation
	if applicant.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Parse date strings if they came in a different format
	if applicant.DateOfBirth.IsZero() {
		dateStr := r.FormValue("date_of_birth")
		if dateStr != "" {
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				http.Error(w, "Invalid date format for date_of_birth: "+err.Error(), http.StatusBadRequest)
				return
			}
			applicant.DateOfBirth = date
		}
	}

	// Parse household member dates if needed
	for i := range applicant.Household {
		if applicant.Household[i].DateOfBirth.IsZero() {
			dateStr := r.FormValue("household[" + strconv.Itoa(i) + "].date_of_birth")
			if dateStr != "" {
				date, err := time.Parse("2006-01-02", dateStr)
				if err != nil {
					http.Error(w, "Invalid date format for household member date_of_birth: "+err.Error(), http.StatusBadRequest)
					return
				}
				applicant.Household[i].DateOfBirth = date
			}
		}
	}

	err = h.ApplicantRepo.Create(&applicant)
	if err != nil {
		http.Error(w, "Failed to create applicant: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.ApplicantResponse{
		Applicant: applicant,
		Household: applicant.Household,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateApplicant handles PUT /api/applicants/{id}
func (h *ApplicantHandler) UpdateApplicant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if applicant exists
	existing, err := h.ApplicantRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to get applicant: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "Applicant not found", http.StatusNotFound)
		return
	}

	var applicant models.Applicant
	err = json.NewDecoder(r.Body).Decode(&applicant)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure ID matches path parameter
	applicant.ID = id

	// Basic validation
	if applicant.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Parse date strings if they came in a different format
	if applicant.DateOfBirth.IsZero() {
		dateStr := r.FormValue("date_of_birth")
		if dateStr != "" {
			date, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				http.Error(w, "Invalid date format for date_of_birth: "+err.Error(), http.StatusBadRequest)
				return
			}
			applicant.DateOfBirth = date
		} else {
			applicant.DateOfBirth = existing.DateOfBirth
		}
	}

	err = h.ApplicantRepo.Update(&applicant)
	if err != nil {
		http.Error(w, "Failed to update applicant: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Note: this doesn't update household members - would need separate endpoints for that

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applicant)
}

// DeleteApplicant handles DELETE /api/applicants/{id}
func (h *ApplicantHandler) DeleteApplicant(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if applicant exists
	existing, err := h.ApplicantRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to get applicant: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "Applicant not found", http.StatusNotFound)
		return
	}

	err = h.ApplicantRepo.Delete(id)
	if err != nil {
		http.Error(w, "Failed to delete applicant: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
