package models

import "time"

// SwaggerApplication is a Swagger-friendly version of Application
// @Description Application for a financial assistance scheme
type SwaggerApplication struct {
	ID              string     `json:"id" example:"01913b7a-4493-74b2-93f8-e684c4ca935c"`
	ApplicantID     string     `json:"applicant_id" example:"01913b7a-4493-74b2-93f8-e684c4ca935c"`
	SchemeID        string     `json:"scheme_id" example:"01913b89-9a43-7163-8757-01cc254783f3"`
	Status          string     `json:"status" example:"pending" enums:"pending,approved,rejected"`
	ApplicationDate time.Time  `json:"application_date"`
	DecisionDate    *time.Time `json:"decision_date,omitempty"`
	Notes           string     `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at,omitempty"`
	UpdatedAt       time.Time  `json:"updated_at,omitempty"`
	Applicant       *Applicant `json:"applicant,omitempty"`
	Scheme          *Scheme    `json:"scheme,omitempty"`
}

// SwaggerApplicationResponse is a Swagger-friendly version of ApplicationResponse
// @Description Response containing an application with applicant and scheme details
type SwaggerApplicationResponse struct {
	SwaggerApplication
	Applicant ApplicantResponse `json:"applicant"`
	Scheme    SchemeResponse    `json:"scheme"`
}
