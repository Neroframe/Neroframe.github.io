package handlers

import (
	"database/sql"
	"html/template"
	"myapp/models"
	"net/http"
	"strconv"
)

type RecordHandler struct {
	DB        *sql.DB
	Templates map[string]*template.Template
}

func NewRecordHandler(db *sql.DB, templates map[string]*template.Template) *RecordHandler {
	return &RecordHandler{
		DB:        db,
		Templates: templates,
	}
}

func (h *RecordHandler) ListRecords(w http.ResponseWriter, r *http.Request) {
	records, err := models.GetAllRecords(h.DB)
	if err != nil {
		http.Error(w, "Error fetching records: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, ok := h.Templates["records/list"]
	if !ok {
		http.Error(w, "Template not found: records/list", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title   string
		Records []models.Record
	}{
		Title:   "Records",
		Records: records,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *RecordHandler) ViewRecord(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	cname := r.URL.Query().Get("cname")
	diseaseCode := r.URL.Query().Get("disease_code")
	if email == "" || cname == "" || diseaseCode == "" {
		http.Error(w, "Missing email, country name, or disease code", http.StatusBadRequest)
		return
	}

	record, err := models.GetRecord(h.DB, email, cname, diseaseCode)
	if err != nil {
		http.Error(w, "Error fetching record: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if record == nil {
		http.NotFound(w, r)
		return
	}

	tmpl, ok := h.Templates["records/view"]
	if !ok {
		http.Error(w, "Template not found: records/view", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title  string
		Record *models.Record
	}{
		Title:  "View Record",
		Record: record,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *RecordHandler) CreateRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		publicServants, err := models.GetAllPublicServants(h.DB)
		if err != nil {
			http.Error(w, "Error fetching public servants: "+err.Error(), http.StatusInternalServerError)
			return
		}

		countries, err := models.GetAllCountries(h.DB)
		if err != nil {
			http.Error(w, "Error fetching countries: "+err.Error(), http.StatusInternalServerError)
			return
		}

		diseases, err := models.GetAllDiseases(h.DB)
		if err != nil {
			http.Error(w, "Error fetching diseases: "+err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, ok := h.Templates["records/form"]
		if !ok {
			http.Error(w, "Template not found: records/form", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title          string
			Record         *models.Record
			PublicServants []models.PublicServant
			Countries      []models.Country
			Diseases       []models.Disease
		}{
			Title:          "Create Record",
			Record:         &models.Record{},
			PublicServants: publicServants,
			Countries:      countries,
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
		cname := r.FormValue("cname")
		diseaseCode := r.FormValue("disease_code")

		if email == "" || cname == "" || diseaseCode == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		record := &models.Record{
			Email:       email,
			CName:       cname,
			DiseaseCode: diseaseCode,
		}

		if err := models.CreateRecord(h.DB, record); err != nil {
			http.Error(w, "Error creating record: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/records", http.StatusSeeOther)
	}
}

func (h *RecordHandler) UpdateRecord(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	cname := r.URL.Query().Get("cname")
	diseaseCode := r.URL.Query().Get("disease_code")
	if email == "" || cname == "" || diseaseCode == "" {
		http.Error(w, "Missing email, country name, or disease code", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		record, err := models.GetRecord(h.DB, email, cname, diseaseCode)
		if err != nil {
			http.Error(w, "Error fetching record: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if record == nil {
			http.NotFound(w, r)
			return
		}

		tmpl, ok := h.Templates["records/form"]
		if !ok {
			http.Error(w, "Template not found: records/form", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title  string
			Record *models.Record
		}{
			Title:  "Edit Record",
			Record: record,
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

		totalDeathsStr := r.FormValue("total_deaths")
		totalPatientsStr := r.FormValue("total_patients")

		if totalDeathsStr == "" || totalPatientsStr == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		totalDeaths, err := strconv.Atoi(totalDeathsStr)
		if err != nil {
			http.Error(w, "Invalid total deaths value", http.StatusBadRequest)
			return
		}

		totalPatients, err := strconv.Atoi(totalPatientsStr)
		if err != nil {
			http.Error(w, "Invalid total patients value", http.StatusBadRequest)
			return
		}

		record := &models.Record{
			Email:         email,
			CName:         cname,
			DiseaseCode:   diseaseCode,
			TotalDeaths:   totalDeaths,
			TotalPatients: totalPatients,
		}

		if err := models.UpdateRecord(h.DB, record); err != nil {
			http.Error(w, "Error updating record: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/records", http.StatusSeeOther)
	}
}

func (h *RecordHandler) DeleteRecord(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	cname := r.URL.Query().Get("cname")
	diseaseCode := r.URL.Query().Get("disease_code")
	if email == "" || cname == "" || diseaseCode == "" {
		http.Error(w, "Missing email, country name, or disease code", http.StatusBadRequest)
		return
	}

	if err := models.DeleteRecord(h.DB, email, cname, diseaseCode); err != nil {
		http.Error(w, "Error deleting record: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/records", http.StatusSeeOther)
}
