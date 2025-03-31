package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SchemeRepository handles database operations for schemes
type SchemeRepository struct {
	DB *sql.DB
}

// NewSchemeRepository creates a new repository with the given database connection
func NewSchemeRepository(db *sql.DB) *SchemeRepository {
	return &SchemeRepository{DB: db}
}

// GetAll retrieves all schemes from the database
func (r *SchemeRepository) GetAll() ([]Scheme, error) {
	query := `SELECT id, name, description, criteria, created_at, updated_at
			  FROM schemes
			  ORDER BY name ASC`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying schemes: %v", err)
	}
	defer rows.Close()

	var schemes []Scheme
	for rows.Next() {
		var s Scheme
		var criteriaJSON []byte

		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &criteriaJSON,
			&s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning scheme row: %v", err)
		}

		// Parse criteria JSON
		if err := json.Unmarshal(criteriaJSON, &s.Criteria); err != nil {
			return nil, fmt.Errorf("error unmarshaling criteria: %v", err)
		}

		// Get benefits for each scheme
		benefits, err := r.GetBenefits(s.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting benefits: %v", err)
		}
		s.Benefits = benefits

		schemes = append(schemes, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating scheme rows: %v", err)
	}

	return schemes, nil
}

// GetByID retrieves a scheme by ID
func (r *SchemeRepository) GetByID(id string) (*Scheme, error) {
	query := `SELECT id, name, description, criteria, created_at, updated_at
			  FROM schemes
			  WHERE id = ?`

	var s Scheme
	var criteriaJSON []byte

	err := r.DB.QueryRow(query, id).Scan(&s.ID, &s.Name, &s.Description, &criteriaJSON,
		&s.CreatedAt, &s.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No scheme found
		}
		return nil, fmt.Errorf("error querying scheme: %v", err)
	}

	// Parse criteria JSON
	if err := json.Unmarshal(criteriaJSON, &s.Criteria); err != nil {
		return nil, fmt.Errorf("error unmarshaling criteria: %v", err)
	}

	// Get benefits
	benefits, err := r.GetBenefits(s.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting benefits: %v", err)
	}
	s.Benefits = benefits

	return &s, nil
}

// Create inserts a new scheme into the database
func (r *SchemeRepository) Create(s *Scheme) error {
	// Generate UUID if not provided
	if s.ID == "" {
		s.ID = uuid.New().String()
	}

	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now

	// Convert criteria to JSON
	criteriaJSON, err := json.Marshal(s.Criteria)
	if err != nil {
		return fmt.Errorf("error marshaling criteria: %v", err)
	}

	query := `INSERT INTO schemes (id, name, description, criteria, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?)`

	_, err = r.DB.Exec(query, s.ID, s.Name, s.Description, criteriaJSON, s.CreatedAt, s.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error creating scheme: %v", err)
	}

	// Create benefits
	for i := range s.Benefits {
		s.Benefits[i].SchemeID = s.ID
		if err := r.CreateBenefit(&s.Benefits[i]); err != nil {
			return fmt.Errorf("error creating benefit: %v", err)
		}
	}

	return nil
}

// Update updates an existing scheme
func (r *SchemeRepository) Update(s *Scheme) error {
	s.UpdatedAt = time.Now()

	// Convert criteria to JSON
	criteriaJSON, err := json.Marshal(s.Criteria)
	if err != nil {
		return fmt.Errorf("error marshaling criteria: %v", err)
	}

	query := `UPDATE schemes
			  SET name = ?, description = ?, criteria = ?, updated_at = ?
			  WHERE id = ?`

	_, err = r.DB.Exec(query, s.Name, s.Description, criteriaJSON, s.UpdatedAt, s.ID)
	if err != nil {
		return fmt.Errorf("error updating scheme: %v", err)
	}

	return nil
}

// Delete removes a scheme
func (r *SchemeRepository) Delete(id string) error {
	query := `DELETE FROM schemes WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting scheme: %v", err)
	}
	return nil
}

// GetBenefits retrieves all benefits for a scheme
func (r *SchemeRepository) GetBenefits(schemeID string) ([]Benefit, error) {
	query := `SELECT id, scheme_id, name, description, amount, created_at, updated_at
			  FROM benefits
			  WHERE scheme_id = ?
			  ORDER BY name ASC`

	rows, err := r.DB.Query(query, schemeID)
	if err != nil {
		return nil, fmt.Errorf("error querying benefits: %v", err)
	}
	defer rows.Close()

	var benefits []Benefit
	for rows.Next() {
		var b Benefit
		var description sql.NullString
		var amount sql.NullFloat64

		if err := rows.Scan(&b.ID, &b.SchemeID, &b.Name, &description, &amount,
			&b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning benefit row: %v", err)
		}

		if description.Valid {
			b.Description = description.String
		}
		if amount.Valid {
			b.Amount = amount.Float64
		}

		benefits = append(benefits, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating benefit rows: %v", err)
	}

	return benefits, nil
}

// CreateBenefit inserts a new benefit
func (r *SchemeRepository) CreateBenefit(b *Benefit) error {
	// Generate UUID if not provided
	if b.ID == "" {
		b.ID = uuid.New().String()
	}

	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now

	query := `INSERT INTO benefits (id, scheme_id, name, description, amount, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.DB.Exec(query, b.ID, b.SchemeID, b.Name, b.Description, b.Amount, b.CreatedAt, b.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error creating benefit: %v", err)
	}

	return nil
}

// DeleteBenefit removes a benefit
func (r *SchemeRepository) DeleteBenefit(id string) error {
	query := `DELETE FROM benefits WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting benefit: %v", err)
	}
	return nil
}

// GetEligibleSchemes finds all schemes for which an applicant is eligible
func (r *SchemeRepository) GetEligibleSchemes(applicantID string, applicantRepo *ApplicantRepository) ([]Scheme, error) {
	// Get applicant with household
	applicant, err := applicantRepo.GetByID(applicantID)
	if err != nil {
		return nil, fmt.Errorf("error getting applicant: %v", err)
	}
	if applicant == nil {
		return nil, fmt.Errorf("applicant not found: %s", applicantID)
	}

	// Get all schemes
	schemes, err := r.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error getting schemes: %v", err)
	}

	var eligibleSchemes []Scheme
	for _, scheme := range schemes {
		if isEligible(applicant, &scheme) {
			eligibleSchemes = append(eligibleSchemes, scheme)
		}
	}

	return eligibleSchemes, nil
}

// isEligible checks if an applicant is eligible for a scheme based on criteria
func isEligible(applicant *Applicant, scheme *Scheme) bool {
	criteria := scheme.Criteria

	// Check employment status
	if criteria.EmploymentStatus != "" &&
		strings.ToLower(criteria.EmploymentStatus) != strings.ToLower(applicant.EmploymentStatus) {
		return false
	}

	// Check marital status
	if criteria.MaritalStatus != "" &&
		strings.ToLower(criteria.MaritalStatus) != strings.ToLower(applicant.MaritalStatus) {
		return false
	}

	// Check children criteria
	if criteria.HasChildren.SchoolLevel != "" {
		hasEligibleChild := false
		for _, member := range applicant.Household {
			// Check if the member is a child
			if strings.Contains(strings.ToLower(member.Relation), "son") ||
				strings.Contains(strings.ToLower(member.Relation), "daughter") {
				// Check age for primary school (roughly 6-12 years)
				age := time.Now().Year() - member.DateOfBirth.Year()
				if age >= 6 && age <= 12 && criteria.HasChildren.SchoolLevel == "primary" {
					hasEligibleChild = true
					break
				}
			}
		}
		if !hasEligibleChild {
			return false
		}
	}

	return true
}
