package handlers

import (
    "database/sql"
    "html/template"
    "myapp/models"
    "net/http"
)

type DoctorHandler struct {
    DB        *sql.DB
    Templates map[string]*template.Template
}

func NewDoctorHandler(db *sql.DB, templates map[string]*template.Template) *DoctorHandler {
    return &DoctorHandler{
        DB:        db,
        Templates: templates,
    }
}

func (h *DoctorHandler) ListDoctors(w http.ResponseWriter, r *http.Request) {
    doctors, err := models.GetAllDoctors(h.DB)
    if err != nil {
        http.Error(w, "Error fetching doctors: "+err.Error(), http.StatusInternalServerError)
        return
    }

    tmpl, ok := h.Templates["doctors/list"]
    if !ok {
        http.Error(w, "Template not found: doctors/list", http.StatusInternalServerError)
        return
    }

    data := struct {
        Title   string
        Doctors []models.Doctor
    }{
        Title:   "Doctors",
        Doctors: doctors,
    }

    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
    }
}

func (h *DoctorHandler) ViewDoctor(w http.ResponseWriter, r *http.Request) {
    email := r.URL.Query().Get("email")
    if email == "" {
        http.Error(w, "Missing doctor email", http.StatusBadRequest)
        return
    }

    doctor, err := models.GetDoctor(h.DB, email)
    if err != nil {
        http.Error(w, "Error fetching doctor: "+err.Error(), http.StatusInternalServerError)
        return
    }
    if doctor == nil {
        http.NotFound(w, r)
        return
    }

    tmpl, ok := h.Templates["doctors/view"]
    if !ok {
        http.Error(w, "Template not found: doctors/view", http.StatusInternalServerError)
        return
    }

    data := struct {
        Title  string
        Doctor *models.Doctor
    }{
        Title:  "View Doctor",
        Doctor: doctor,
    }

    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
    }
}

func (h *DoctorHandler) CreateDoctor(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        tmpl, ok := h.Templates["doctors/form"]
        if !ok {
            http.Error(w, "Template not found: doctors/form", http.StatusInternalServerError)
            return
        }

        data := struct {
            Title  string
            Doctor *models.Doctor
        }{
            Title:  "Create Doctor",
            Doctor: &models.Doctor{},
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

        degree := r.FormValue("degree")
        if degree == "" {
            http.Error(w, "Degree is required", http.StatusBadRequest)
            return
        }

        doctor := &models.Doctor{
            Email:  email,
            Degree: degree,
        }

        if err := models.CreateDoctor(h.DB, doctor); err != nil {
            http.Error(w, "Error creating doctor: "+err.Error(), http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/doctors", http.StatusSeeOther)
    }
}

func (h *DoctorHandler) UpdateDoctor(w http.ResponseWriter, r *http.Request) {
    email := r.URL.Query().Get("email")
    if email == "" {
        http.Error(w, "Missing doctor email", http.StatusBadRequest)
        return
    }

    if r.Method == "GET" {
        doctor, err := models.GetDoctor(h.DB, email)
        if err != nil {
            http.Error(w, "Error fetching doctor: "+err.Error(), http.StatusInternalServerError)
            return
        }
        if doctor == nil {
            http.NotFound(w, r)
            return
        }

        tmpl, ok := h.Templates["doctors/form"]
        if !ok {
            http.Error(w, "Template not found: doctors/form", http.StatusInternalServerError)
            return
        }

        data := struct {
            Title  string
            Doctor *models.Doctor
        }{
            Title:  "Edit Doctor",
            Doctor: doctor,
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

        degree := r.FormValue("degree")
        if degree == "" {
            http.Error(w, "Degree is required", http.StatusBadRequest)
            return
        }

        doctor := &models.Doctor{
            Email:  email,
            Degree: degree,
        }

        if err := models.UpdateDoctor(h.DB, doctor); err != nil {
            http.Error(w, "Error updating doctor: "+err.Error(), http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/doctors", http.StatusSeeOther)
    }
}

func (h *DoctorHandler) DeleteDoctor(w http.ResponseWriter, r *http.Request) {
    email := r.URL.Query().Get("email")
    if email == "" {
        http.Error(w, "Missing doctor email", http.StatusBadRequest)
        return
    }

    if err := models.DeleteDoctor(h.DB, email); err != nil {
        http.Error(w, "Error deleting doctor: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/doctors", http.StatusSeeOther)
}
