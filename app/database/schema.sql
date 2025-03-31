-- Schema

-- Applicants table
CREATE TABLE applicants (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    employment_status ENUM('employed', 'unemployed') NOT NULL,
    sex ENUM('male', 'female', 'other') NOT NULL,
    date_of_birth DATE NOT NULL,
    marital_status ENUM('single', 'married', 'widowed', 'divorced') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Household members table
CREATE TABLE household_members (
    id VARCHAR(36) PRIMARY KEY,
    applicant_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    employment_status ENUM('employed', 'unemployed') NOT NULL,
    sex ENUM('male', 'female', 'other') NOT NULL,
    date_of_birth DATE NOT NULL,
    relation VARCHAR(50) NOT NULL, -- e.g., 'son', 'daughter', 'spouse', etc.
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (applicant_id) REFERENCES applicants(id) ON DELETE CASCADE
);

-- Schemes table
CREATE TABLE schemes (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    criteria JSON NOT NULL, -- Store eligibility criteria as JSON
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Benefits table
CREATE TABLE benefits (
    id VARCHAR(36) PRIMARY KEY,
    scheme_id VARCHAR(36) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    amount DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (scheme_id) REFERENCES schemes(id) ON DELETE CASCADE
);

-- Applications table
CREATE TABLE applications (
    id VARCHAR(36) PRIMARY KEY,
    applicant_id VARCHAR(36) NOT NULL,
    scheme_id VARCHAR(36) NOT NULL,
    status ENUM('pending', 'approved', 'rejected') NOT NULL DEFAULT 'pending',
    application_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    decision_date TIMESTAMP NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (applicant_id) REFERENCES applicants(id),
    FOREIGN KEY (scheme_id) REFERENCES schemes(id)
);

-- Indexes for performance
CREATE INDEX idx_household_applicant ON household_members(applicant_id);
CREATE INDEX idx_benefits_scheme ON benefits(scheme_id);
CREATE INDEX idx_applications_applicant ON applications(applicant_id);
CREATE INDEX idx_applications_scheme ON applications(scheme_id);

-- Sample data for testing

-- Sample applicants
INSERT INTO applicants (id, name, employment_status, sex, date_of_birth, marital_status)
VALUES 
('01913b7a-4493-74b2-93f8-e684c4ca935c', 'James', 'unemployed', 'male', '1990-07-01', 'single'),
('01913b80-2c04-7f9d-86a4-497ef68cb3a0', 'Mary', 'unemployed', 'female', '1984-10-06', 'married');

-- Sample household members
INSERT INTO household_members (id, applicant_id, name, employment_status, sex, date_of_birth, relation)
VALUES
('01913b88-1d4d-7152-a7ce-75796a2e8ecf', '01913b80-2c04-7f9d-86a4-497ef68cb3a0', 'Gwen', 'unemployed', 'female', '2016-02-01', 'daughter'),
('01913b88-65c6-7255-820f-9c4dd1e5ce79', '01913b80-2c04-7f9d-86a4-497ef68cb3a0', 'Jayden', 'unemployed', 'male', '2018-03-15', 'son');

-- Sample schemes
INSERT INTO schemes (id, name, description, criteria)
VALUES
('01913b89-9a43-7163-8757-01cc254783f3', 'Retrenchment Assistance Scheme', 'Financial assistance for retrenched workers', '{"employment_status": "unemployed"}'),
('01913b89-befc-7ae3-bb37-3079aa7f1be0', 'Retrenchment Assistance Scheme (families)', 'Financial assistance for retrenched workers with primary school children', '{"employment_status": "unemployed", "has_children": {"school_level": "primary"}}');

-- Sample benefits
INSERT INTO benefits (id, scheme_id, name, description, amount)
VALUES
('01913b8b-9b12-7d2c-a1fa-ea613b802ebc', '01913b89-9a43-7163-8757-01cc254783f3', 'SkillsFuture Credits', 'Additional SkillsFuture credits for training', 500.00),
('01913b8c-5d33-7e9a-b2fa-fb723c904def', '01913b89-befc-7ae3-bb37-3079aa7f1be0', 'School Meal Vouchers', 'Daily school meal vouchers for primary school children', 200.00); 