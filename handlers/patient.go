package handlers

import (
	"database/sql"
	"html/template"
	"myapp/models"
	"net/http"
)

type PatientHandler struct {
	DB        *sql.DB
	Templates map[string]*template.Template
}

func NewPatientHandler(db *sql.DB, templates map[string]*template.Template) *PatientHandler {
	return &PatientHandler{
		DB:        db,
		Templates: templates,
	}
}

func (h *PatientHandler) ListPatients(w http.ResponseWriter, r *http.Request) {
	patients, err := models.GetAllPatients(h.DB)
	if err != nil {
		http.Error(w, "Error fetching patients: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, ok := h.Templates["patients/list"]
	if !ok {
		http.Error(w, "Template not found: patients/list", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title    string
		Patients []models.Patient
	}{
		Title:    "Patients",
		Patients: patients,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *PatientHandler) ViewPatient(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Missing patient email", http.StatusBadRequest)
		return
	}

	patient, err := models.GetPatient(h.DB, email)
	if err != nil {
		http.Error(w, "Error fetching patient: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if patient == nil {
		http.NotFound(w, r)
		return
	}

	tmpl, ok := h.Templates["patients/view"]
	if !ok {
		http.Error(w, "Template not found: patients/view", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title   string
		Patient *models.Patient
	}{
		Title:   "View Patient",
		Patient: patient,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *PatientHandler) CreatePatient(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, ok := h.Templates["patients/form"]
		if !ok {
			http.Error(w, "Template not found: patients/form", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title   string
			Patient *models.Patient
		}{
			Title:   "Create Patient",
			Patient: &models.Patient{},
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
		if email == "" {
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}

		patient := &models.Patient{
			Email: email,
		}

		if err := models.CreatePatient(h.DB, patient); err != nil {
			http.Error(w, "Error creating patient: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/patients", http.StatusSeeOther)
	}
}

func (h *PatientHandler) UpdatePatient(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Missing patient email", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		patient, err := models.GetPatient(h.DB, email)
		if err != nil {
			http.Error(w, "Error fetching patient: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if patient == nil {
			http.NotFound(w, r)
			return
		}

		tmpl, ok := h.Templates["patients/form"]
		if !ok {
			http.Error(w, "Template not found: patients/form", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title   string
			Patient *models.Patient
		}{
			Title:   "Edit Patient",
			Patient: patient,
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

		updatedPatient := &models.Patient{
			Email: email, 
		}

		if err := models.UpdatePatient(h.DB, updatedPatient); err != nil {
			http.Error(w, "Error updating patient: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/patients", http.StatusSeeOther)
	}
}

func (h *PatientHandler) DeletePatient(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Missing patient email", http.StatusBadRequest)
		return
	}

	if err := models.DeletePatient(h.DB, email); err != nil {
		http.Error(w, "Error deleting patient: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/patients", http.StatusSeeOther)
}
