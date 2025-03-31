# one-client-view-2025tht

A backend application for managing financial assistance schemes and applications. This system allows administrators to manage financial schemes, applicants, and applications.

## Features

- Manage financial assistance schemes
- Add and update applicant records
- View eligible schemes for each applicant
- Submit and manage applications for financial assistance

## Technology Stack

- Backend: Go (Golang)
- Database: MySQL
- API: RESTful API with JSON

## Prerequisites

- Go 1.20 or higher
- MySQL 8.0 or higher

## Setup Instructions

### 1. Clone the repository

```bash
git clone https://github.com/peasantbird/one-client-view-2025tht
cd one-client-view-2025tht
```

### 2. Set up the database

Create a MySQL database and run the schema file to set up tables:

```bash
mysql -u root -p -e "CREATE DATABASE one_client_view_2025tht;"
mysql -u root -p one_client_view_2025tht < app/database/schema.sql
```

### 3. Configure environment variables

Create a `.env` file in the root directory:

```
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=one_client_view_2025tht
PORT=8080
```

### 4. Install dependencies

```bash
go mod download
```

### 5. Build and run the application

```bash
go build -o one-client-view-2025tht ./app
./one-client-view-2025tht
```

The server will start running at `http://localhost:8080` by default.

## API Endpoints

### Applicants

- `GET /api/applicants` - Get all applicants
- `POST /api/applicants` - Create a new applicant
- `GET /api/applicants/{id}` - Get applicant by ID
- `PUT /api/applicants/{id}` - Update applicant
- `DELETE /api/applicants/{id}` - Delete applicant

### Schemes

- `GET /api/schemes` - Get all schemes
- `POST /api/schemes` - Create a new scheme
- `GET /api/schemes/{id}` - Get scheme by ID
- `PUT /api/schemes/{id}` - Update scheme
- `DELETE /api/schemes/{id}` - Delete scheme
- `GET /api/schemes/eligible?applicant={id}` - Get eligible schemes for an applicant

### Applications

- `GET /api/applications` - Get all applications
- `POST /api/applications` - Create a new application
- `GET /api/applications/{id}` - Get application by ID
- `PUT /api/applications/{id}` - Update application
- `DELETE /api/applications/{id}` - Delete application

## Data Models

### Applicant

```json
{
  "id": "uuid",
  "name": "string",
  "employment_status": "employed|unemployed",
  "sex": "male|female|other",
  "date_of_birth": "date",
  "marital_status": "single|married|widowed|divorced",
  "household": [
    {
      "id": "uuid",
      "name": "string",
      "employment_status": "employed|unemployed",
      "sex": "male|female|other",
      "date_of_birth": "date",
      "relation": "string"
    }
  ]
}
```

### Scheme

```json
{
  "id": "uuid",
  "name": "string",
  "description": "string",
  "criteria": {
    "employment_status": "string",
    "marital_status": "string",
    "has_children": {
      "school_level": "string"
    }
  },
  "benefits": [
    {
      "id": "uuid",
      "name": "string",
      "description": "string",
      "amount": "number"
    }
  ]
}
```

### Application

```json
{
  "id": "uuid",
  "applicant_id": "uuid",
  "scheme_id": "uuid",
  "status": "pending|approved|rejected",
  "application_date": "datetime",
  "decision_date": "datetime",
  "notes": "string"
}
```
