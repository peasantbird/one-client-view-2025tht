package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"one-client-view-2025tht/app/models"
)

// ApplicationHandler handles HTTP requests related to applications
type ApplicationHandler struct {
	ApplicationRepo *models.ApplicationRepository
	ApplicantRepo   *models.ApplicantRepository
	SchemeRepo      *models.SchemeRepository
}

// NewApplicationHandler creates a new handler with the given repositories
func NewApplicationHandler(appRepo *models.ApplicationRepository, applicantRepo *models.ApplicantRepository, schemeRepo *models.SchemeRepository) *ApplicationHandler {
	return &ApplicationHandler{
		ApplicationRepo: appRepo,
		ApplicantRepo:   applicantRepo,
		SchemeRepo:      schemeRepo,
	}
}

// GetApplications handles GET /api/applications
// @Summary Get all applications
// @Description Retrieve a list of all financial assistance applications
// @Tags applications
// @Accept json
// @Produce json
// @Success 200 {array} models.SwaggerApplicationResponse
// @Failure 500 {object} string "Internal server error"
// @Router /api/applications [get]
func (h *ApplicationHandler) GetApplications(w http.ResponseWriter, r *http.Request) {
	applications, err := h.ApplicationRepo.GetAll()
	if err != nil {
		http.Error(w, "Failed to get applications: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert to response objects
	var response []models.ApplicationResponse
	for _, a := range applications {
		if a.Applicant == nil || a.Scheme == nil {
			continue // Skip invalid applications
		}

		response = append(response, models.ApplicationResponse{
			Application: a,
			Applicant: models.ApplicantResponse{
				Applicant: *a.Applicant,
				Household: a.Applicant.Household,
			},
			Scheme: models.SchemeResponse{
				Scheme:   *a.Scheme,
				Benefits: a.Scheme.Benefits,
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetApplication handles GET /api/applications/{id}
// @Summary Get application by ID
// @Description Retrieve a specific application by its ID
// @Tags applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 200 {object} models.SwaggerApplicationResponse
// @Failure 404 {object} string "Application not found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/applications/{id} [get]
func (h *ApplicationHandler) GetApplication(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	application, err := h.ApplicationRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to get application: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if application == nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}

	if application.Applicant == nil || application.Scheme == nil {
		http.Error(w, "Invalid application data", http.StatusInternalServerError)
		return
	}

	response := models.ApplicationResponse{
		Application: *application,
		Applicant: models.ApplicantResponse{
			Applicant: *application.Applicant,
			Household: application.Applicant.Household,
		},
		Scheme: models.SchemeResponse{
			Scheme:   *application.Scheme,
			Benefits: application.Scheme.Benefits,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateApplication handles POST /api/applications
// @Summary Create a new application
// @Description Submit a new application for a financial assistance scheme
// @Tags applications
// @Accept json
// @Produce json
// @Param application body models.ApplicationRequest true "Application information"
// @Success 201 {object} models.SwaggerApplicationResponse
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Applicant or scheme not found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/applications [post]
func (h *ApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	var request models.ApplicationRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Basic validation
	if request.ApplicantID == "" {
		http.Error(w, "Applicant ID is required", http.StatusBadRequest)
		return
	}
	if request.SchemeID == "" {
		http.Error(w, "Scheme ID is required", http.StatusBadRequest)
		return
	}

	// Check if applicant exists
	applicant, err := h.ApplicantRepo.GetByID(request.ApplicantID)
	if err != nil {
		http.Error(w, "Failed to get applicant: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if applicant == nil {
		http.Error(w, "Applicant not found", http.StatusNotFound)
		return
	}

	// Check if scheme exists
	scheme, err := h.SchemeRepo.GetByID(request.SchemeID)
	if err != nil {
		http.Error(w, "Failed to get scheme: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if scheme == nil {
		http.Error(w, "Scheme not found", http.StatusNotFound)
		return
	}

	// Create application
	application := &models.Application{
		ApplicantID: request.ApplicantID,
		SchemeID:    request.SchemeID,
		Notes:       request.Notes,
		Status:      "pending",
	}

	// Try to create the application
	err = h.ApplicationRepo.Create(application)
	if err != nil {
		http.Error(w, "Failed to create application: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the created application with all details
	createdApp, err := h.ApplicationRepo.GetByID(application.ID)
	if err != nil {
		http.Error(w, "Application created but failed to retrieve details: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.ApplicationResponse{
		Application: *createdApp,
		Applicant: models.ApplicantResponse{
			Applicant: *createdApp.Applicant,
			Household: createdApp.Applicant.Household,
		},
		Scheme: models.SchemeResponse{
			Scheme:   *createdApp.Scheme,
			Benefits: createdApp.Scheme.Benefits,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateApplication handles PUT /api/applications/{id}
// @Summary Update application
// @Description Update an existing application's status or notes
// @Tags applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Param application body object{status=string,notes=string} true "Updated application information"
// @Success 200 {object} models.SwaggerApplicationResponse
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Application not found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/applications/{id} [put]
func (h *ApplicationHandler) UpdateApplication(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if application exists
	existing, err := h.ApplicationRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to get application: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}

	var request struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Update only status and notes
	if request.Status != "" {
		existing.Status = request.Status
	}
	if request.Notes != "" {
		existing.Notes = request.Notes
	}

	err = h.ApplicationRepo.Update(existing)
	if err != nil {
		http.Error(w, "Failed to update application: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the updated application with all details
	updatedApp, err := h.ApplicationRepo.GetByID(existing.ID)
	if err != nil {
		http.Error(w, "Application updated but failed to retrieve details: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.ApplicationResponse{
		Application: *updatedApp,
		Applicant: models.ApplicantResponse{
			Applicant: *updatedApp.Applicant,
			Household: updatedApp.Applicant.Household,
		},
		Scheme: models.SchemeResponse{
			Scheme:   *updatedApp.Scheme,
			Benefits: updatedApp.Scheme.Benefits,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteApplication handles DELETE /api/applications/{id}
// @Summary Delete application
// @Description Remove an application from the system
// @Tags applications
// @Accept json
// @Produce json
// @Param id path string true "Application ID"
// @Success 204 "No content"
// @Failure 404 {object} string "Application not found"
// @Failure 500 {object} string "Internal server error"
// @Router /api/applications/{id} [delete]
func (h *ApplicationHandler) DeleteApplication(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if application exists
	existing, err := h.ApplicationRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Failed to get application: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}

	err = h.ApplicationRepo.Delete(id)
	if err != nil {
		http.Error(w, "Failed to delete application: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
