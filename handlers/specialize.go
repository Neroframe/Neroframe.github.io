package handlers

import (
    "database/sql"
    "html/template"
    "myapp/models"
    "net/http"
    "strconv"
)

type SpecializeHandler struct {
    DB        *sql.DB
    Templates map[string]*template.Template
}

func NewSpecializeHandler(db *sql.DB, templates map[string]*template.Template) *SpecializeHandler {
    return &SpecializeHandler{
        DB:        db,
        Templates: templates,
    }
}

func (h *SpecializeHandler) ListSpecializes(w http.ResponseWriter, r *http.Request) {
    specializes, err := models.GetAllSpecializes(h.DB)
    if err != nil {
        http.Error(w, "Error fetching specializations: "+err.Error(), http.StatusInternalServerError)
        return
    }

    tmpl, ok := h.Templates["specializes/list"]
    if !ok {
        http.Error(w, "Template not found: specializes/list", http.StatusInternalServerError)
        return
    }

    data := struct {
        Title       string
        Specializes []models.Specialize
    }{
        Title:       "Specializations",
        Specializes: specializes,
    }

    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
    }
}

func (h *SpecializeHandler) ViewSpecialize(w http.ResponseWriter, r *http.Request) {
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

    email := r.URL.Query().Get("email")
    if email == "" {
        http.Error(w, "Missing doctor email", http.StatusBadRequest)
        return
    }

    specialize, err := models.GetSpecialize(h.DB, id, email)
    if err != nil {
        http.Error(w, "Error fetching specialization: "+err.Error(), http.StatusInternalServerError)
        return
    }
    if specialize == nil {
        http.NotFound(w, r)
        return
    }

    tmpl, ok := h.Templates["specializes/view"]
    if !ok {
        http.Error(w, "Template not found: specializes/view", http.StatusInternalServerError)
        return
    }

    data := struct {
        Title      string
        Specialize *models.Specialize
    }{
        Title:      "View Specialization",
        Specialize: specialize,
    }

    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
    }
}

func (h *SpecializeHandler) CreateSpecialize(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        diseaseTypes, err := models.GetAllDiseaseTypes(h.DB)
        if err != nil {
            http.Error(w, "Error fetching disease types: "+err.Error(), http.StatusInternalServerError)
            return
        }

        doctors, err := models.GetAllDoctors(h.DB)
        if err != nil {
            http.Error(w, "Error fetching doctors: "+err.Error(), http.StatusInternalServerError)
            return
        }

        tmpl, ok := h.Templates["specializes/form"]
        if !ok {
            http.Error(w, "Template not found: specializes/form", http.StatusInternalServerError)
            return
        }

        data := struct {
            Title        string
            Specialize   *models.Specialize
            DiseaseTypes []models.DiseaseType
            Doctors      []models.Doctor
        }{
            Title:        "Create Specialization",
            Specialize:   &models.Specialize{},
            DiseaseTypes: diseaseTypes,
            Doctors:      doctors,
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

        idStr := r.FormValue("id")
        email := r.FormValue("email")
        if idStr == "" || email == "" {
            http.Error(w, "All fields are required", http.StatusBadRequest)
            return
        }

        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, "Invalid disease type ID", http.StatusBadRequest)
            return
        }

        specialize := &models.Specialize{
            ID:    id,
            Email: email,
        }

        if err := models.CreateSpecialize(h.DB, specialize); err != nil {
            http.Error(w, "Error creating specialization: "+err.Error(), http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/specializes", http.StatusSeeOther)
    }
}

func (h *SpecializeHandler) UpdateSpecialize(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    email := r.URL.Query().Get("email")
    if idStr == "" || email == "" {
        http.Error(w, "Missing disease type ID or doctor email", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid disease type ID", http.StatusBadRequest)
        return
    }

    if r.Method == "GET" {
        specialize, err := models.GetSpecialize(h.DB, id, email)
        if err != nil {
            http.Error(w, "Error fetching specialization: "+err.Error(), http.StatusInternalServerError)
            return
        }
        if specialize == nil {
            http.NotFound(w, r)
            return
        }

        diseaseTypes, err := models.GetAllDiseaseTypes(h.DB)
        if err != nil {
            http.Error(w, "Error fetching disease types: "+err.Error(), http.StatusInternalServerError)
            return
        }

        doctors, err := models.GetAllDoctors(h.DB)
        if err != nil {
            http.Error(w, "Error fetching doctors: "+err.Error(), http.StatusInternalServerError)
            return
        }

        tmpl, ok := h.Templates["specializes/form"]
        if !ok {
            http.Error(w, "Template not found: specializes/form", http.StatusInternalServerError)
            return
        }

        data := struct {
            Title        string
            Specialize   *models.Specialize
            DiseaseTypes []models.DiseaseType
            Doctors      []models.Doctor
        }{
            Title:        "Edit Specialization",
            Specialize:   specialize,
            DiseaseTypes: diseaseTypes,
            Doctors:      doctors,
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

        newIDStr := r.FormValue("id")
        newEmail := r.FormValue("email")
        if newIDStr == "" || newEmail == "" {
            http.Error(w, "All fields are required", http.StatusBadRequest)
            return
        }

        newID, err := strconv.Atoi(newIDStr)
        if err != nil {
            http.Error(w, "Invalid disease type ID", http.StatusBadRequest)
            return
        }

        // Delete the old record then insert since we cannot update primary keys directly
        if err := models.DeleteSpecialize(h.DB, id, email); err != nil {
            http.Error(w, "Error updating specialization: "+err.Error(), http.StatusInternalServerError)
            return
        }

        specialize := &models.Specialize{
            ID:    newID,
            Email: newEmail,
        }

        if err := models.CreateSpecialize(h.DB, specialize); err != nil {
            http.Error(w, "Error updating specialization: "+err.Error(), http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/specializes", http.StatusSeeOther)
    }
}

func (h *SpecializeHandler) DeleteSpecialize(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    email := r.URL.Query().Get("email")
    if idStr == "" || email == "" {
        http.Error(w, "Missing disease type ID or doctor email", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid disease type ID", http.StatusBadRequest)
        return
    }

    if err := models.DeleteSpecialize(h.DB, id, email); err != nil {
        http.Error(w, "Error deleting specialization: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/specializes", http.StatusSeeOther)
}
