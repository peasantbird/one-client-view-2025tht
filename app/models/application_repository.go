package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ApplicationRepository handles database operations for applications
type ApplicationRepository struct {
	DB            *sql.DB
	ApplicantRepo *ApplicantRepository
	SchemeRepo    *SchemeRepository
}

// NewApplicationRepository creates a new repository with the given database connection
func NewApplicationRepository(db *sql.DB, applicantRepo *ApplicantRepository, schemeRepo *SchemeRepository) *ApplicationRepository {
	return &ApplicationRepository{
		DB:            db,
		ApplicantRepo: applicantRepo,
		SchemeRepo:    schemeRepo,
	}
}

// GetAll retrieves all applications from the database
func (r *ApplicationRepository) GetAll() ([]Application, error) {
	query := `SELECT id, applicant_id, scheme_id, status, application_date, decision_date, notes, created_at, updated_at
			  FROM applications
			  ORDER BY application_date DESC`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying applications: %v", err)
	}
	defer rows.Close()

	var applications []Application
	for rows.Next() {
		var a Application
		var decisionDate sql.NullTime
		var notes sql.NullString

		if err := rows.Scan(&a.ID, &a.ApplicantID, &a.SchemeID, &a.Status,
			&a.ApplicationDate, &decisionDate, &notes, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning application row: %v", err)
		}

		if decisionDate.Valid {
			a.DecisionDate = decisionDate
		}
		if notes.Valid {
			a.Notes = notes.String
		}

		// Get applicant and scheme details
		applicant, err := r.ApplicantRepo.GetByID(a.ApplicantID)
		if err != nil {
			return nil, fmt.Errorf("error getting applicant: %v", err)
		}
		a.Applicant = applicant

		scheme, err := r.SchemeRepo.GetByID(a.SchemeID)
		if err != nil {
			return nil, fmt.Errorf("error getting scheme: %v", err)
		}
		a.Scheme = scheme

		applications = append(applications, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating application rows: %v", err)
	}

	return applications, nil
}

// GetByID retrieves an application by ID
func (r *ApplicationRepository) GetByID(id string) (*Application, error) {
	query := `SELECT id, applicant_id, scheme_id, status, application_date, decision_date, notes, created_at, updated_at
			  FROM applications
			  WHERE id = ?`

	var a Application
	var decisionDate sql.NullTime
	var notes sql.NullString

	err := r.DB.QueryRow(query, id).Scan(&a.ID, &a.ApplicantID, &a.SchemeID, &a.Status,
		&a.ApplicationDate, &decisionDate, &notes, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No application found
		}
		return nil, fmt.Errorf("error querying application: %v", err)
	}

	if decisionDate.Valid {
		a.DecisionDate = decisionDate
	}
	if notes.Valid {
		a.Notes = notes.String
	}

	// Get applicant and scheme details
	applicant, err := r.ApplicantRepo.GetByID(a.ApplicantID)
	if err != nil {
		return nil, fmt.Errorf("error getting applicant: %v", err)
	}
	a.Applicant = applicant

	scheme, err := r.SchemeRepo.GetByID(a.SchemeID)
	if err != nil {
		return nil, fmt.Errorf("error getting scheme: %v", err)
	}
	a.Scheme = scheme

	return &a, nil
}

// GetByApplicantID retrieves all applications for an applicant
func (r *ApplicationRepository) GetByApplicantID(applicantID string) ([]Application, error) {
	query := `SELECT id, applicant_id, scheme_id, status, application_date, decision_date, notes, created_at, updated_at
			  FROM applications
			  WHERE applicant_id = ?
			  ORDER BY application_date DESC`

	rows, err := r.DB.Query(query, applicantID)
	if err != nil {
		return nil, fmt.Errorf("error querying applications: %v", err)
	}
	defer rows.Close()

	var applications []Application
	for rows.Next() {
		var a Application
		var decisionDate sql.NullTime
		var notes sql.NullString

		if err := rows.Scan(&a.ID, &a.ApplicantID, &a.SchemeID, &a.Status,
			&a.ApplicationDate, &decisionDate, &notes, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning application row: %v", err)
		}

		if decisionDate.Valid {
			a.DecisionDate = decisionDate
		}
		if notes.Valid {
			a.Notes = notes.String
		}

		// Get scheme details
		scheme, err := r.SchemeRepo.GetByID(a.SchemeID)
		if err != nil {
			return nil, fmt.Errorf("error getting scheme: %v", err)
		}
		a.Scheme = scheme

		applications = append(applications, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating application rows: %v", err)
	}

	return applications, nil
}

// Create inserts a new application into the database
func (r *ApplicationRepository) Create(a *Application) error {
	// Validate applicant and scheme exist
	applicant, err := r.ApplicantRepo.GetByID(a.ApplicantID)
	if err != nil {
		return fmt.Errorf("error validating applicant: %v", err)
	}
	if applicant == nil {
		return fmt.Errorf("applicant not found: %s", a.ApplicantID)
	}

	scheme, err := r.SchemeRepo.GetByID(a.SchemeID)
	if err != nil {
		return fmt.Errorf("error validating scheme: %v", err)
	}
	if scheme == nil {
		return fmt.Errorf("scheme not found: %s", a.SchemeID)
	}

	// Check if applicant is eligible for the scheme
	if !isEligible(applicant, scheme) {
		return fmt.Errorf("applicant is not eligible for this scheme")
	}

	// Generate UUID if not provided
	if a.ID == "" {
		a.ID = uuid.New().String()
	}

	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
	a.ApplicationDate = now

	// Set default status if not provided
	if a.Status == "" {
		a.Status = "pending"
	}

	query := `INSERT INTO applications (id, applicant_id, scheme_id, status, application_date, notes, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = r.DB.Exec(query, a.ID, a.ApplicantID, a.SchemeID, a.Status,
		a.ApplicationDate, a.Notes, a.CreatedAt, a.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating application: %v", err)
	}

	return nil
}

// Update updates an existing application
func (r *ApplicationRepository) Update(a *Application) error {
	a.UpdatedAt = time.Now()

	var decisionDate interface{}
	if a.DecisionDate.Valid {
		decisionDate = a.DecisionDate.Time
	} else {
		decisionDate = nil
	}

	query := `UPDATE applications
			  SET status = ?, decision_date = ?, notes = ?, updated_at = ?
			  WHERE id = ?`

	_, err := r.DB.Exec(query, a.Status, decisionDate, a.Notes, a.UpdatedAt, a.ID)
	if err != nil {
		return fmt.Errorf("error updating application: %v", err)
	}

	return nil
}

// UpdateStatus updates the status of an application
func (r *ApplicationRepository) UpdateStatus(id, status string) error {
	now := time.Now()
	var decisionDate interface{}

	// If status is approved or rejected, set decision date
	if status == "approved" || status == "rejected" {
		decisionDate = now
	} else {
		decisionDate = nil
	}

	query := `UPDATE applications
			  SET status = ?, decision_date = ?, updated_at = ?
			  WHERE id = ?`

	_, err := r.DB.Exec(query, status, decisionDate, now, id)
	if err != nil {
		return fmt.Errorf("error updating application status: %v", err)
	}

	return nil
}

// Delete removes an application
func (r *ApplicationRepository) Delete(id string) error {
	query := `DELETE FROM applications WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting application: %v", err)
	}
	return nil
}
