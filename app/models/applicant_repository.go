package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ApplicantRepository handles database operations for applicants
type ApplicantRepository struct {
	DB *sql.DB
}

// NewApplicantRepository creates a new repository with the given database connection
func NewApplicantRepository(db *sql.DB) *ApplicantRepository {
	return &ApplicantRepository{DB: db}
}

// GetAll retrieves all applicants from the database
func (r *ApplicantRepository) GetAll() ([]Applicant, error) {
	query := `SELECT id, name, employment_status, sex, date_of_birth, marital_status, created_at, updated_at
			  FROM applicants
			  ORDER BY name ASC`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying applicants: %v", err)
	}
	defer rows.Close()

	var applicants []Applicant
	for rows.Next() {
		var a Applicant
		if err := rows.Scan(&a.ID, &a.Name, &a.EmploymentStatus, &a.Sex, &a.DateOfBirth,
			&a.MaritalStatus, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning applicant row: %v", err)
		}

		// Get household members for each applicant
		members, err := r.GetHouseholdMembers(a.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting household members: %v", err)
		}
		a.Household = members

		applicants = append(applicants, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating applicant rows: %v", err)
	}

	return applicants, nil
}

// GetByID retrieves an applicant by ID
func (r *ApplicantRepository) GetByID(id string) (*Applicant, error) {
	query := `SELECT id, name, employment_status, sex, date_of_birth, marital_status, created_at, updated_at
			  FROM applicants
			  WHERE id = ?`

	var a Applicant
	err := r.DB.QueryRow(query, id).Scan(&a.ID, &a.Name, &a.EmploymentStatus, &a.Sex,
		&a.DateOfBirth, &a.MaritalStatus, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No applicant found
		}
		return nil, fmt.Errorf("error querying applicant: %v", err)
	}

	// Get household members
	members, err := r.GetHouseholdMembers(a.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting household members: %v", err)
	}
	a.Household = members

	return &a, nil
}

// Create inserts a new applicant into the database
func (r *ApplicantRepository) Create(a *Applicant) error {
	// Generate UUID if not provided
	if a.ID == "" {
		a.ID = uuid.New().String()
	}

	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

	query := `INSERT INTO applicants (id, name, employment_status, sex, date_of_birth, marital_status, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.DB.Exec(query, a.ID, a.Name, a.EmploymentStatus, a.Sex,
		a.DateOfBirth, a.MaritalStatus, a.CreatedAt, a.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating applicant: %v", err)
	}

	// Create household members
	for i := range a.Household {
		a.Household[i].ApplicantID = a.ID
		if err := r.CreateHouseholdMember(&a.Household[i]); err != nil {
			return fmt.Errorf("error creating household member: %v", err)
		}
	}

	return nil
}

// Update updates an existing applicant
func (r *ApplicantRepository) Update(a *Applicant) error {
	a.UpdatedAt = time.Now()

	query := `UPDATE applicants
			  SET name = ?, employment_status = ?, sex = ?,
				  date_of_birth = ?, marital_status = ?, updated_at = ?
			  WHERE id = ?`

	_, err := r.DB.Exec(query, a.Name, a.EmploymentStatus, a.Sex,
		a.DateOfBirth, a.MaritalStatus, a.UpdatedAt, a.ID)

	if err != nil {
		return fmt.Errorf("error updating applicant: %v", err)
	}

	return nil
}

// Delete removes an applicant
func (r *ApplicantRepository) Delete(id string) error {
	query := `DELETE FROM applicants WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting applicant: %v", err)
	}
	return nil
}

// GetHouseholdMembers retrieves all household members for an applicant
func (r *ApplicantRepository) GetHouseholdMembers(applicantID string) ([]HouseholdMember, error) {
	query := `SELECT id, applicant_id, name, employment_status, sex, date_of_birth, relation, created_at, updated_at
			  FROM household_members
			  WHERE applicant_id = ?
			  ORDER BY name ASC`

	rows, err := r.DB.Query(query, applicantID)
	if err != nil {
		return nil, fmt.Errorf("error querying household members: %v", err)
	}
	defer rows.Close()

	var members []HouseholdMember
	for rows.Next() {
		var m HouseholdMember
		if err := rows.Scan(&m.ID, &m.ApplicantID, &m.Name, &m.EmploymentStatus, &m.Sex,
			&m.DateOfBirth, &m.Relation, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning household member row: %v", err)
		}
		members = append(members, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating household member rows: %v", err)
	}

	return members, nil
}

// CreateHouseholdMember inserts a new household member
func (r *ApplicantRepository) CreateHouseholdMember(m *HouseholdMember) error {
	// Generate UUID if not provided
	if m.ID == "" {
		m.ID = uuid.New().String()
	}

	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now

	query := `INSERT INTO household_members (id, applicant_id, name, employment_status, sex, date_of_birth, relation, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.DB.Exec(query, m.ID, m.ApplicantID, m.Name, m.EmploymentStatus, m.Sex,
		m.DateOfBirth, m.Relation, m.CreatedAt, m.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error creating household member: %v", err)
	}

	return nil
}

// DeleteHouseholdMember removes a household member
func (r *ApplicantRepository) DeleteHouseholdMember(id string) error {
	query := `DELETE FROM household_members WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting household member: %v", err)
	}
	return nil
}
