package handlers

import (
    "database/sql"
    "html/template"
    "myapp/models"
    "net/http"
    "strconv"
)

type DiseaseHandler struct {
    DB        *sql.DB
    Templates map[string]*template.Template
}

func NewDiseaseHandler(db *sql.DB, templates map[string]*template.Template) *DiseaseHandler {
    return &DiseaseHandler{
        DB:        db,
        Templates: templates,
    }
}

func (h *DiseaseHandler) ListDiseases(w http.ResponseWriter, r *http.Request) {
    diseases, err := models.GetAllDiseases(h.DB)
    if err != nil {
        http.Error(w, "Error fetching diseases: "+err.Error(), http.StatusInternalServerError)
        return
    }

    tmpl, ok := h.Templates["diseases/list"]
    if !ok {
        http.Error(w, "Template not found: diseases/list", http.StatusInternalServerError)
        return
    }

    data := struct {
        Title    string
        Diseases []models.Disease
    }{
        Title:    "Diseases",
        Diseases: diseases,
    }

    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
    }
}

func (h *DiseaseHandler) ViewDisease(w http.ResponseWriter, r *http.Request) {
    diseaseCode := r.URL.Query().Get("disease_code")
    if diseaseCode == "" {
        http.Error(w, "Missing disease code", http.StatusBadRequest)
        return
    }

    disease, err := models.GetDisease(h.DB, diseaseCode)
    if err != nil {
        http.Error(w, "Error fetching disease: "+err.Error(), http.StatusInternalServerError)
        return
    }
    if disease == nil {
        http.NotFound(w, r)
        return
    }

    tmpl, ok := h.Templates["diseases/view"]
    if !ok {
        http.Error(w, "Template not found: diseases/view", http.StatusInternalServerError)
        return
    }

    data := struct {
        Title   string
        Disease *models.Disease
    }{
        Title:   "View Disease",
        Disease: disease,
    }

    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
    }
}

func (h *DiseaseHandler) CreateDisease(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        diseaseTypes, err := models.GetAllDiseaseTypes(h.DB)
        if err != nil {
            http.Error(w, "Error fetching disease types: "+err.Error(), http.StatusInternalServerError)
            return
        }

        tmpl, ok := h.Templates["diseases/form"]
        if !ok {
            http.Error(w, "Template not found: diseases/form", http.StatusInternalServerError)
            return
        }

        data := struct {
            Title        string
            Disease      *models.Disease
            DiseaseTypes []models.DiseaseType
        }{
            Title:        "Create Disease",
            Disease:      &models.Disease{},
            DiseaseTypes: diseaseTypes,
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

        diseaseCode := r.FormValue("disease_code")
        pathogen := r.FormValue("pathogen")
        description := r.FormValue("description")
        idStr := r.FormValue("id")

        if diseaseCode == "" || pathogen == "" || description == "" || idStr == "" {
            http.Error(w, "All fields are required", http.StatusBadRequest)
            return
        }

        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, "Invalid disease type ID", http.StatusBadRequest)
            return
        }

        disease := &models.Disease{
            DiseaseCode: diseaseCode,
            Pathogen:    pathogen,
            Description: description,
            ID:          id,
        }

        if err := models.CreateDisease(h.DB, disease); err != nil {
            http.Error(w, "Error creating disease: "+err.Error(), http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/diseases", http.StatusSeeOther)
    }
}

func (h *DiseaseHandler) UpdateDisease(w http.ResponseWriter, r *http.Request) {
    diseaseCode := r.URL.Query().Get("disease_code")
    if diseaseCode == "" {
        http.Error(w, "Missing disease code", http.StatusBadRequest)
        return
    }

    if r.Method == "GET" {
        disease, err := models.GetDisease(h.DB, diseaseCode)
        if err != nil {
            http.Error(w, "Error fetching disease: "+err.Error(), http.StatusInternalServerError)
            return
        }
        if disease == nil {
            http.NotFound(w, r)
            return
        }

        diseaseTypes, err := models.GetAllDiseaseTypes(h.DB)
        if err != nil {
            http.Error(w, "Error fetching disease types: "+err.Error(), http.StatusInternalServerError)
            return
        }

        tmpl, ok := h.Templates["diseases/form"]
        if !ok {
            http.Error(w, "Template not found: diseases/form", http.StatusInternalServerError)
            return
        }

        data := struct {
            Title        string
            Disease      *models.Disease
            DiseaseTypes []models.DiseaseType
        }{
            Title:        "Edit Disease",
            Disease:      disease,
            DiseaseTypes: diseaseTypes,
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

        pathogen := r.FormValue("pathogen")
        description := r.FormValue("description")
        idStr := r.FormValue("id")

        if pathogen == "" || description == "" || idStr == "" {
            http.Error(w, "All fields are required", http.StatusBadRequest)
            return
        }

        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, "Invalid disease type ID", http.StatusBadRequest)
            return
        }

        disease := &models.Disease{
            DiseaseCode: diseaseCode,
            Pathogen:    pathogen,
            Description: description,
            ID:          id,
        }

        if err := models.UpdateDisease(h.DB, disease); err != nil {
            http.Error(w, "Error updating disease: "+err.Error(), http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/diseases", http.StatusSeeOther)
    }
}

func (h *DiseaseHandler) DeleteDisease(w http.ResponseWriter, r *http.Request) {
    diseaseCode := r.URL.Query().Get("disease_code")
    if diseaseCode == "" {
        http.Error(w, "Missing disease code", http.StatusBadRequest)
        return
    }

    if err := models.DeleteDisease(h.DB, diseaseCode); err != nil {
        http.Error(w, "Error deleting disease: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/diseases", http.StatusSeeOther)
}
