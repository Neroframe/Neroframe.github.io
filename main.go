// main.go
package main

import (
	"html/template"
	"log"
	"myapp/db"
	"myapp/handlers"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Fetch DATABASE_URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatalf("DATABASE_URL is not set. Ensure it's available in your environment variables.")
	}

	dbConn, err := db.NewPostgresDB(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer dbConn.Close()

	log.Println("Successfully connected to the database!")

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	templates, err := parseTemplates(
		"templates/*.html",
		"templates/users/*.html",
		"templates/countries/*.html",
		"templates/disease_types/*.html",
		"templates/diseases/*.html",
		"templates/discovers/*.html",
		"templates/specializes/*.html",
		"templates/patients/*.html",
		"templates/public_servants/*.html",
		"templates/doctors/*.html",
		"templates/patient_diseases/*.html", 
		"templates/records/*.html", 
	)
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	dashboardHandler := handlers.NewDashboardHandler(templates)
	userHandler := handlers.NewUserHandler(dbConn, templates)
	countryHandler := handlers.NewCountryHandler(dbConn, templates) 
	diseaseTypeHandler := handlers.NewDiseaseTypeHandler(dbConn, templates)
	diseaseHandler := handlers.NewDiseaseHandler(dbConn, templates)  
	discoverHandler := handlers.NewDiscoverHandler(dbConn, templates)
	specializeHandler := handlers.NewSpecializeHandler(dbConn, templates)
	patientHandler := handlers.NewPatientHandler(dbConn, templates)
	publicServantHandler := handlers.NewPublicServantHandler(dbConn, templates)
	doctorHandler := handlers.NewDoctorHandler(dbConn, templates)
	patientDiseaseHandler := handlers.NewPatientDiseaseHandler(dbConn, templates)
	recordHandler := handlers.NewRecordHandler(dbConn, templates)

	// Dashboard route
	http.HandleFunc("/", dashboardHandler.Dashboard)

	// Routes for CRUD
	http.HandleFunc("/users", userHandler.ListUsers)
	http.HandleFunc("/users/view", userHandler.ViewUser)
	http.HandleFunc("/users/create", userHandler.CreateUser)
	http.HandleFunc("/users/edit", userHandler.UpdateUser)
	http.HandleFunc("/users/delete", userHandler.DeleteUser)

	http.HandleFunc("/countries", countryHandler.ListCountries)
	http.HandleFunc("/countries/view", countryHandler.ViewCountry)
	http.HandleFunc("/countries/create", countryHandler.CreateCountry)
	http.HandleFunc("/countries/edit", countryHandler.UpdateCountry)
	http.HandleFunc("/countries/delete", countryHandler.DeleteCountry)

	http.HandleFunc("/disease_types", diseaseTypeHandler.ListDiseaseTypes)
	http.HandleFunc("/disease_types/view", diseaseTypeHandler.ViewDiseaseType)
	http.HandleFunc("/disease_types/create", diseaseTypeHandler.CreateDiseaseType)
	http.HandleFunc("/disease_types/edit", diseaseTypeHandler.UpdateDiseaseType)
	http.HandleFunc("/disease_types/delete", diseaseTypeHandler.DeleteDiseaseType)

	http.HandleFunc("/diseases", diseaseHandler.ListDiseases)
	http.HandleFunc("/diseases/view", diseaseHandler.ViewDisease)
	http.HandleFunc("/diseases/create", diseaseHandler.CreateDisease)
	http.HandleFunc("/diseases/edit", diseaseHandler.UpdateDisease)
	http.HandleFunc("/diseases/delete", diseaseHandler.DeleteDisease)

	http.HandleFunc("/discovers", discoverHandler.ListDiscovers)
	http.HandleFunc("/discovers/view", discoverHandler.ViewDiscover)
	http.HandleFunc("/discovers/create", discoverHandler.CreateDiscover)
	http.HandleFunc("/discovers/edit", discoverHandler.UpdateDiscover)
	http.HandleFunc("/discovers/delete", discoverHandler.DeleteDiscover)

	http.HandleFunc("/specializes", specializeHandler.ListSpecializes)
	http.HandleFunc("/specializes/view", specializeHandler.ViewSpecialize)
	http.HandleFunc("/specializes/create", specializeHandler.CreateSpecialize)
	http.HandleFunc("/specializes/edit", specializeHandler.UpdateSpecialize)
	http.HandleFunc("/specializes/delete", specializeHandler.DeleteSpecialize)

	http.HandleFunc("/patients", patientHandler.ListPatients)
	http.HandleFunc("/patients/view", patientHandler.ViewPatient)
	http.HandleFunc("/patients/create", patientHandler.CreatePatient)
	http.HandleFunc("/patients/edit", patientHandler.UpdatePatient)
	http.HandleFunc("/patients/delete", patientHandler.DeletePatient)

	http.HandleFunc("/public_servants", publicServantHandler.ListPublicServants)
	http.HandleFunc("/public_servants/view", publicServantHandler.ViewPublicServant)
	http.HandleFunc("/public_servants/create", publicServantHandler.CreatePublicServant)
	http.HandleFunc("/public_servants/edit", publicServantHandler.UpdatePublicServant)
	http.HandleFunc("/public_servants/delete", publicServantHandler.DeletePublicServant)

	http.HandleFunc("/doctors", doctorHandler.ListDoctors)
	http.HandleFunc("/doctors/view", doctorHandler.ViewDoctor)
	http.HandleFunc("/doctors/create", doctorHandler.CreateDoctor)
	http.HandleFunc("/doctors/edit", doctorHandler.UpdateDoctor)
	http.HandleFunc("/doctors/delete", doctorHandler.DeleteDoctor)

	http.HandleFunc("/patient_diseases", patientDiseaseHandler.ListPatientDiseases)
	http.HandleFunc("/patient_diseases/view", patientDiseaseHandler.ViewPatientDisease)
	http.HandleFunc("/patient_diseases/create", patientDiseaseHandler.CreatePatientDisease)
	http.HandleFunc("/patient_diseases/edit", patientDiseaseHandler.UpdatePatientDisease)
	http.HandleFunc("/patient_diseases/delete", patientDiseaseHandler.DeletePatientDisease)

	http.HandleFunc("/records", recordHandler.ListRecords)
	http.HandleFunc("/records/view", recordHandler.ViewRecord)
	http.HandleFunc("/records/create", recordHandler.CreateRecord)
	http.HandleFunc("/records/edit", recordHandler.UpdateRecord)
	http.HandleFunc("/records/delete", recordHandler.DeleteRecord)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}

	log.Printf("Server starting on port %s", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func parseTemplates(patterns ...string) (map[string]*template.Template, error) {
	tmplMap := make(map[string]*template.Template)
	layoutFiles, err := filepath.Glob("templates/*.html")
	if err != nil {
		return nil, err
	}

	for _, pattern := range patterns {
		files, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			tmplFiles := append(layoutFiles, file)
			tmpl, err := template.ParseFiles(tmplFiles...)
			if err != nil {
				return nil, err
			}

			// Create key based on relative path
			relPath, err := filepath.Rel("templates", file)
			if err != nil {
				return nil, err
			}
			// Remove .html extension
			key := relPath[:len(relPath)-len(".html")]

			tmplMap[key] = tmpl
		}
	}

	return tmplMap, nil
}
