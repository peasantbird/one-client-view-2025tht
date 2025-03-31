package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	_ "one-client-view-2025tht/docs" // This will be auto-generated

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	"one-client-view-2025tht/app/database"
	"one-client-view-2025tht/app/handlers"
	"one-client-view-2025tht/app/models"
)

// @host localhost:8080
// @BasePath /api
// @schemes http

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found. Using environment variables.")
	}

	// Configure database
	dbConfig := &database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 3306),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "one_client_view_2025tht"),
	}

	// Initialize database connection
	err = database.Initialize(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Create repositories
	db := database.GetDB()
	applicantRepo := models.NewApplicantRepository(db)
	schemeRepo := models.NewSchemeRepository(db)
	applicationRepo := models.NewApplicationRepository(db, applicantRepo, schemeRepo)

	// Create handlers
	applicantHandler := handlers.NewApplicantHandler(applicantRepo)
	schemeHandler := handlers.NewSchemeHandler(schemeRepo, applicantRepo)
	applicationHandler := handlers.NewApplicationHandler(applicationRepo, applicantRepo, schemeRepo)

	// Create router
	router := mux.NewRouter()

	// API routes
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Applicant routes
	apiRouter.HandleFunc("/applicants", applicantHandler.GetApplicants).Methods("GET")
	apiRouter.HandleFunc("/applicants", applicantHandler.CreateApplicant).Methods("POST")
	apiRouter.HandleFunc("/applicants/{id}", applicantHandler.GetApplicant).Methods("GET")
	apiRouter.HandleFunc("/applicants/{id}", applicantHandler.UpdateApplicant).Methods("PUT")
	apiRouter.HandleFunc("/applicants/{id}", applicantHandler.DeleteApplicant).Methods("DELETE")

	// Scheme routes
	apiRouter.HandleFunc("/schemes", schemeHandler.GetSchemes).Methods("GET")
	apiRouter.HandleFunc("/schemes", schemeHandler.CreateScheme).Methods("POST")
	apiRouter.HandleFunc("/schemes/eligible", schemeHandler.GetEligibleSchemes).Methods("GET")
	apiRouter.HandleFunc("/schemes/{id}", schemeHandler.GetScheme).Methods("GET")
	apiRouter.HandleFunc("/schemes/{id}", schemeHandler.UpdateScheme).Methods("PUT")
	apiRouter.HandleFunc("/schemes/{id}", schemeHandler.DeleteScheme).Methods("DELETE")

	// Application routes
	apiRouter.HandleFunc("/applications", applicationHandler.GetApplications).Methods("GET")
	apiRouter.HandleFunc("/applications", applicationHandler.CreateApplication).Methods("POST")
	apiRouter.HandleFunc("/applications/{id}", applicationHandler.GetApplication).Methods("GET")
	apiRouter.HandleFunc("/applications/{id}", applicationHandler.UpdateApplication).Methods("PUT")
	apiRouter.HandleFunc("/applications/{id}", applicationHandler.DeleteApplication).Methods("DELETE")

	// Swagger documentation
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// Configure CORS middleware
	router.Use(corsMiddleware)

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// CORS middleware to allow cross-origin requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to get environment variable as an integer
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
