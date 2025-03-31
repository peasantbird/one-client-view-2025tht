package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"one-client-view-2025tht/app/models"
)

// SchemeHandler handles HTTP requests related to schemes
type SchemeHandler struct {
	SchemeRepo    *models.SchemeRepository
	ApplicantRepo *models.ApplicantRepository
}

// NewSchemeHandler creates a new handler with the given repositories
func NewSchemeHandler(schemeRepo *models.SchemeRepository, applicantRepo *models.ApplicantRepository) *SchemeHandler {
	return &SchemeHandler{
		SchemeRepo:    schemeRepo,
		ApplicantRepo: applicantRepo,
	}
}

// GetSchemes handles GET /api/schemes
// @Summary Get all schemes
// @Description Retrieve a list of all financial assistance schemes
// @Tags schemes
// @Accept json
// @Produce json
// @Success 200 {array} models.SchemeResponse
// @Failure 500 {object} string "Internal server error"
// @Router /api/schemes [get]
func (h *SchemeHandler) GetSchemes(w http.ResponseWriter, r *http.Request) {
	schemes, err := h.SchemeRepo.GetAll()
	if err != nil {
		http.Error(w, "Failed to get schemes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response objects
	var response []models.SchemeResponse
	for _, s := range schemes {
		response = append(response, models.SchemeResponse{
			Scheme:   s,
			Benefits: s.Benefits,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetScheme handles GET /api/schemes/{id}
// @Summary Get scheme by ID
// @Description Retrieve a specific scheme by its ID
// @Tags schemes
// @Accept json
// @Produce json
// @Param id path string true "Scheme ID"
// @Success 200 {object} models.SchemeResponse
// @Failure 404 {object} string "Scheme not found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/schemes/{id} [get]
func (h *SchemeHandler) GetScheme(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	scheme, err := h.SchemeRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to get scheme: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if scheme == nil {
		http.Error(w, "Scheme not found", http.StatusNotFound)
		return
	}

	response := models.SchemeResponse{
		Scheme:   *scheme,
		Benefits: scheme.Benefits,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetEligibleSchemes handles GET /api/schemes/eligible?applicant={id}
// @Summary Get eligible schemes for an applicant
// @Description Retrieve all schemes that an applicant is eligible for
// @Tags schemes
// @Accept json
// @Produce json
// @Param applicant query string true "Applicant ID"
// @Success 200 {object} models.EligibleSchemesResponse
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Applicant not found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/schemes/eligible [get]
func (h *SchemeHandler) GetEligibleSchemes(w http.ResponseWriter, r *http.Request) {
	applicantID := r.URL.Query().Get("applicant")
	if applicantID == "" {
		http.Error(w, "Applicant ID is required", http.StatusBadRequest)
		return
	}

	// Check if applicant exists
	applicant, err := h.ApplicantRepo.GetByID(applicantID)
	if err != nil {
		http.Error(w, "Failed to get applicant: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if applicant == nil {
		http.Error(w, "Applicant not found", http.StatusNotFound)
		return
	}

	// Get eligible schemes
	schemes, err := h.SchemeRepo.GetEligibleSchemes(applicantID, h.ApplicantRepo)
	if err != nil {
		http.Error(w, "Failed to get eligible schemes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response objects
	var schemeResponses []models.SchemeResponse
	for _, s := range schemes {
		schemeResponses = append(schemeResponses, models.SchemeResponse{
			Scheme:   s,
			Benefits: s.Benefits,
		})
	}

	response := models.EligibleSchemesResponse{
		ApplicantID: applicantID,
		Schemes:     schemeResponses,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateScheme handles POST /api/schemes
// @Summary Create a new scheme
// @Description Add a new financial assistance scheme
// @Tags schemes
// @Accept json
// @Produce json
// @Param scheme body models.Scheme true "Scheme information"
// @Success 201 {object} models.SchemeResponse
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal server error"
// @Router /api/schemes [post]
func (h *SchemeHandler) CreateScheme(w http.ResponseWriter, r *http.Request) {
	var scheme models.Scheme
	err := json.NewDecoder(r.Body).Decode(&scheme)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Basic validation
	if scheme.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	if scheme.Description == "" {
		http.Error(w, "Description is required", http.StatusBadRequest)
		return
	}

	err = h.SchemeRepo.Create(&scheme)
	if err != nil {
		http.Error(w, "Failed to create scheme: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.SchemeResponse{
		Scheme:   scheme,
		Benefits: scheme.Benefits,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateScheme handles PUT /api/schemes/{id}
// @Summary Update scheme
// @Description Update an existing scheme's information
// @Tags schemes
// @Accept json
// @Produce json
// @Param id path string true "Scheme ID"
// @Param scheme body models.Scheme true "Updated scheme information"
// @Success 200 {object} models.SchemeResponse
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Scheme not found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/schemes/{id} [put]
func (h *SchemeHandler) UpdateScheme(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if scheme exists
	existing, err := h.SchemeRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to get scheme: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "Scheme not found", http.StatusNotFound)
		return
	}

	var scheme models.Scheme
	err = json.NewDecoder(r.Body).Decode(&scheme)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure ID matches path parameter
	scheme.ID = id

	// Basic validation
	if scheme.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	if scheme.Description == "" {
		http.Error(w, "Description is required", http.StatusBadRequest)
		return
	}

	// Preserve benefits
	scheme.Benefits = existing.Benefits

	err = h.SchemeRepo.Update(&scheme)
	if err != nil {
		http.Error(w, "Failed to update scheme: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.SchemeResponse{
		Scheme:   scheme,
		Benefits: scheme.Benefits,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteScheme handles DELETE /api/schemes/{id}
// @Summary Delete scheme
// @Description Remove a scheme from the system
// @Tags schemes
// @Accept json
// @Produce json
// @Param id path string true "Scheme ID"
// @Success 204 "No content"
// @Failure 404 {object} string "Scheme not found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/schemes/{id} [delete]
func (h *SchemeHandler) DeleteScheme(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if scheme exists
	existing, err := h.SchemeRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to get scheme: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "Scheme not found", http.StatusNotFound)
		return
	}

	err = h.SchemeRepo.Delete(id)
	if err != nil {
		http.Error(w, "Failed to delete scheme: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
