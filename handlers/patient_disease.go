package handlers

import (
	"database/sql"
	"html/template"
	"myapp/models"
	"net/http"
)

type PatientDiseaseHandler struct {
	DB        *sql.DB
	Templates map[string]*template.Template
}

func NewPatientDiseaseHandler(db *sql.DB, templates map[string]*template.Template) *PatientDiseaseHandler {
	return &PatientDiseaseHandler{
		DB:        db,
		Templates: templates,
	}
}

func (h *PatientDiseaseHandler) ListPatientDiseases(w http.ResponseWriter, r *http.Request) {
	patientDiseases, err := models.GetAllPatientDiseases(h.DB)
	if err != nil {
		http.Error(w, "Error fetching patient diseases: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, ok := h.Templates["patient_diseases/list"]
	if !ok {
		http.Error(w, "Template not found: patient_diseases/list", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title           string
		PatientDiseases []models.PatientDisease
	}{
		Title:           "Patient Diseases",
		PatientDiseases: patientDiseases,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *PatientDiseaseHandler) ViewPatientDisease(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	diseaseCode := r.URL.Query().Get("disease_code")
	if email == "" || diseaseCode == "" {
		http.Error(w, "Missing patient email or disease code", http.StatusBadRequest)
		return
	}

	patientDisease, err := models.GetPatientDisease(h.DB, email, diseaseCode)
	if err != nil {
		http.Error(w, "Error fetching patient disease: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if patientDisease == nil {
		http.NotFound(w, r)
		return
	}

	tmpl, ok := h.Templates["patient_diseases/view"]
	if !ok {
		http.Error(w, "Template not found: patient_diseases/view", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title          string
		PatientDisease *models.PatientDisease
	}{
		Title:          "View Patient Disease",
		PatientDisease: patientDisease,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *PatientDiseaseHandler) CreatePatientDisease(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		patients, err := models.GetAllPatients(h.DB)
		if err != nil {
			http.Error(w, "Error fetching patients: "+err.Error(), http.StatusInternalServerError)
			return
		}

		diseases, err := models.GetAllDiseases(h.DB)
		if err != nil {
			http.Error(w, "Error fetching diseases: "+err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, ok := h.Templates["patient_diseases/form"]
		if !ok {
			http.Error(w, "Template not found: patient_diseases/form", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title          string
			PatientDisease *models.PatientDisease
			Patients       []models.Patient
			Diseases       []models.Disease
		}{
			Title:          "Create Patient Disease",
			PatientDisease: &models.PatientDisease{},
			Patients:       patients,
			Diseases:       diseases,
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		diseaseCode := r.FormValue("disease_code")

		if email == "" || diseaseCode == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		patientDisease := &models.PatientDisease{
			Email:       email,
			DiseaseCode: diseaseCode,
		}

		if err := models.CreatePatientDisease(h.DB, patientDisease); err != nil {
			http.Error(w, "Error creating patient disease: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/patient_diseases", http.StatusSeeOther)
	}
}

func (h *PatientDiseaseHandler) UpdatePatientDisease(w http.ResponseWriter, r *http.Request) {
	oldEmail := r.URL.Query().Get("email")
	oldDiseaseCode := r.URL.Query().Get("disease_code")
	if oldEmail == "" || oldDiseaseCode == "" {
		http.Error(w, "Missing patient email or disease code", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		patientDisease, err := models.GetPatientDisease(h.DB, oldEmail, oldDiseaseCode)
		if err != nil {
			http.Error(w, "Error fetching patient disease: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if patientDisease == nil {
			http.NotFound(w, r)
			return
		}

		tmpl, ok := h.Templates["patient_diseases/form"]
		if !ok {
			http.Error(w, "Template not found: patient_diseases/form", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title          string
			PatientDisease *models.PatientDisease
			Diseases       []models.Disease
		}{
			Title:          "Edit Patient Disease",
			PatientDisease: patientDisease,
			Diseases:       []models.Disease{}, // Empty slice since patients cannot change
		}

		// Fetch all diseases for the dropdown
		diseases, err := models.GetAllDiseases(h.DB)
		if err != nil {
			http.Error(w, "Error fetching diseases: "+err.Error(), http.StatusInternalServerError)
			return
		}
		data.Diseases = diseases

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
			return
		}

		newDiseaseCode := r.FormValue("disease_code")
		if newDiseaseCode == "" {
			http.Error(w, "Disease code is required", http.StatusBadRequest)
			return
		}

		if err := models.UpdatePatientDisease(h.DB, oldEmail, oldDiseaseCode, newDiseaseCode); err != nil {
			http.Error(w, "Error updating patient disease: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/patient_diseases", http.StatusSeeOther)
	}
}

func (h *PatientDiseaseHandler) DeletePatientDisease(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	diseaseCode := r.URL.Query().Get("disease_code")
	if email == "" || diseaseCode == "" {
		http.Error(w, "Missing patient email or disease code", http.StatusBadRequest)
		return
	}

	if err := models.DeletePatientDisease(h.DB, email, diseaseCode); err != nil {
		http.Error(w, "Error deleting patient disease: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/patient_diseases", http.StatusSeeOther)
}
