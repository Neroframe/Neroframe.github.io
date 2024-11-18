package handlers

import (
	"database/sql"
	"html/template"
	"myapp/models"
	"net/http"
	"strconv"
)

type DiseaseTypeHandler struct {
	DB        *sql.DB
	Templates map[string]*template.Template
}

func NewDiseaseTypeHandler(db *sql.DB, templates map[string]*template.Template) *DiseaseTypeHandler {
	return &DiseaseTypeHandler{
		DB:        db,
		Templates: templates,
	}
}

func (h *DiseaseTypeHandler) ListDiseaseTypes(w http.ResponseWriter, r *http.Request) {
	diseaseTypes, err := models.GetAllDiseaseTypes(h.DB)
	if err != nil {
		http.Error(w, "Error fetching disease types: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, ok := h.Templates["disease_types/list"]
	if !ok {
		http.Error(w, "Template not found: disease_types/list", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title        string
		DiseaseTypes []models.DiseaseType
	}{
		Title:        "Disease Types",
		DiseaseTypes: diseaseTypes,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *DiseaseTypeHandler) ViewDiseaseType(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing disease type ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid disease type ID", http.StatusBadRequest)
		return
	}

	diseaseType, err := models.GetDiseaseType(h.DB, id)
	if err != nil {
		http.Error(w, "Error fetching disease type: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if diseaseType == nil {
		http.NotFound(w, r)
		return
	}

	tmpl, ok := h.Templates["disease_types/view"]
	if !ok {
		http.Error(w, "Template not found: disease_types/view", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title       string
		DiseaseType *models.DiseaseType
	}{
		Title:       "View Disease Type",
		DiseaseType: diseaseType,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *DiseaseTypeHandler) CreateDiseaseType(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, ok := h.Templates["disease_types/form"]
		if !ok {
			http.Error(w, "Template not found: disease_types/form", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title       string
			DiseaseType *models.DiseaseType
		}{
			Title:       "Create Disease Type",
			DiseaseType: &models.DiseaseType{},
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

		description := r.FormValue("description")
		if description == "" {
			http.Error(w, "Description is required", http.StatusBadRequest)
			return
		}

		diseaseType := &models.DiseaseType{
			ID:          0, // ID will be auto-incremented by the database
			Description: description,
		}

		if err := models.CreateDiseaseType(h.DB, diseaseType); err != nil {
			http.Error(w, "Error creating disease type: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/disease_types", http.StatusSeeOther)
	}
}

func (h *DiseaseTypeHandler) UpdateDiseaseType(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing disease type ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid disease type ID", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		diseaseType, err := models.GetDiseaseType(h.DB, id)
		if err != nil {
			http.Error(w, "Error fetching disease type: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if diseaseType == nil {
			http.NotFound(w, r)
			return
		}

		tmpl, ok := h.Templates["disease_types/form"]
		if !ok {
			http.Error(w, "Template not found: disease_types/form", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title       string
			DiseaseType *models.DiseaseType
		}{
			Title:       "Edit Disease Type",
			DiseaseType: diseaseType,
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

		description := r.FormValue("description")
		if description == "" {
			http.Error(w, "Description is required", http.StatusBadRequest)
			return
		}

		diseaseType := &models.DiseaseType{
			ID:          id,
			Description: description,
		}

		if err := models.UpdateDiseaseType(h.DB, diseaseType); err != nil {
			http.Error(w, "Error updating disease type: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/disease_types", http.StatusSeeOther)
	}
}

func (h *DiseaseTypeHandler) DeleteDiseaseType(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing disease type ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid disease type ID", http.StatusBadRequest)
		return
	}

	if err := models.DeleteDiseaseType(h.DB, id); err != nil {
		http.Error(w, "Error deleting disease type: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/disease_types", http.StatusSeeOther)
}
