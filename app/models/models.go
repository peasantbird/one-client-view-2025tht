package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Applicant represents an individual applying for financial assistance
type Applicant struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	EmploymentStatus string            `json:"employment_status"`
	Sex              string            `json:"sex"`
	DateOfBirth      time.Time         `json:"date_of_birth"`
	MaritalStatus    string            `json:"marital_status"`
	CreatedAt        time.Time         `json:"created_at,omitempty"`
	UpdatedAt        time.Time         `json:"updated_at,omitempty"`
	Household        []HouseholdMember `json:"household,omitempty"`
}

// HouseholdMember represents a family member living with the applicant
type HouseholdMember struct {
	ID               string    `json:"id"`
	ApplicantID      string    `json:"applicant_id"`
	Name             string    `json:"name"`
	EmploymentStatus string    `json:"employment_status"`
	Sex              string    `json:"sex"`
	DateOfBirth      time.Time `json:"date_of_birth"`
	Relation         string    `json:"relation"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}

// Criteria represents the eligibility criteria for schemes
type Criteria struct {
	EmploymentStatus string        `json:"employment_status,omitempty"`
	MaritalStatus    string        `json:"marital_status,omitempty"`
	HasChildren      ChildCriteria `json:"has_children,omitempty"`
}

// ChildCriteria represents specific criteria related to children
type ChildCriteria struct {
	SchoolLevel string `json:"school_level,omitempty"`
}

// Scheme represents a financial assistance scheme
type Scheme struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Criteria    Criteria  `json:"criteria"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	Benefits    []Benefit `json:"benefits,omitempty"`
}

// Benefit represents benefits provided by a scheme
type Benefit struct {
	ID          string    `json:"id"`
	SchemeID    string    `json:"scheme_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Amount      float64   `json:"amount,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// Application represents an application for a financial assistance scheme
type Application struct {
	ID              string       `json:"id"`
	ApplicantID     string       `json:"applicant_id"`
	SchemeID        string       `json:"scheme_id"`
	Status          string       `json:"status"`
	ApplicationDate time.Time    `json:"application_date"`
	DecisionDate    sql.NullTime `json:"decision_date,omitempty"`
	Notes           string       `json:"notes,omitempty"`
	CreatedAt       time.Time    `json:"created_at,omitempty"`
	UpdatedAt       time.Time    `json:"updated_at,omitempty"`
	Applicant       *Applicant   `json:"applicant,omitempty"`
	Scheme          *Scheme      `json:"scheme,omitempty"`
}

// UnmarshalJSON custom unmarshaler for Scheme to handle the JSON criteria field
func (s *Scheme) UnmarshalJSON(data []byte) error {
	type Alias Scheme
	aux := &struct {
		Criteria json.RawMessage `json:"criteria"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return json.Unmarshal(aux.Criteria, &s.Criteria)
}

// MarshalJSON custom marshaler for Scheme to handle the JSON criteria field
func (s Scheme) MarshalJSON() ([]byte, error) {
	type Alias Scheme
	criteriaJSON, err := json.Marshal(s.Criteria)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&struct {
		Criteria json.RawMessage `json:"criteria"`
		*Alias
	}{
		Criteria: criteriaJSON,
		Alias:    (*Alias)(&s),
	})
}

// ApplicantResponse is used for API responses that include household members
type ApplicantResponse struct {
	Applicant
	Household []HouseholdMember `json:"household"`
}

// SchemeResponse is used for API responses that include benefits
type SchemeResponse struct {
	Scheme
	Benefits []Benefit `json:"benefits"`
}

// ApplicationRequest is used for creating a new application
type ApplicationRequest struct {
	ApplicantID string `json:"applicant_id"`
	SchemeID    string `json:"scheme_id"`
	Notes       string `json:"notes,omitempty"`
}

// ApplicationResponse is used for API responses
type ApplicationResponse struct {
	Application
	Applicant ApplicantResponse `json:"applicant"`
	Scheme    SchemeResponse    `json:"scheme"`
}

// EligibleSchemesResponse is used for returning eligible schemes for an applicant
type EligibleSchemesResponse struct {
	ApplicantID string           `json:"applicant_id"`
	Schemes     []SchemeResponse `json:"schemes"`
}
