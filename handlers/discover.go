package handlers

import (
    "database/sql"
    "html/template"
    "myapp/models"
    "net/http"
    "time"
)

type DiscoverHandler struct {
    DB        *sql.DB
    Templates map[string]*template.Template
}

func NewDiscoverHandler(db *sql.DB, templates map[string]*template.Template) *DiscoverHandler {
    return &DiscoverHandler{
        DB:        db,
        Templates: templates,
    }
}

func (h *DiscoverHandler) ListDiscovers(w http.ResponseWriter, r *http.Request) {
    discovers, err := models.GetAllDiscovers(h.DB)
    if err != nil {
        http.Error(w, "Error fetching discoveries: "+err.Error(), http.StatusInternalServerError)
        return
    }

    tmpl, ok := h.Templates["discovers/list"]
    if !ok {
        http.Error(w, "Template not found: discovers/list", http.StatusInternalServerError)
        return
    }

    data := struct {
        Title     string
        Discovers []models.Discover
    }{
        Title:     "Discoveries",
        Discovers: discovers,
    }

    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
    }
}

func (h *DiscoverHandler) ViewDiscover(w http.ResponseWriter, r *http.Request) {
    cname := r.URL.Query().Get("cname")
    diseaseCode := r.URL.Query().Get("disease_code")
    if cname == "" || diseaseCode == "" {
        http.Error(w, "Missing country name or disease code", http.StatusBadRequest)
        return
    }

    discover, err := models.GetDiscover(h.DB, cname, diseaseCode)
    if err != nil {
        http.Error(w, "Error fetching discovery: "+err.Error(), http.StatusInternalServerError)
        return
    }
    if discover == nil {
        http.NotFound(w, r)
        return
    }

    tmpl, ok := h.Templates["discovers/view"]
    if !ok {
        http.Error(w, "Template not found: discovers/view", http.StatusInternalServerError)
        return
    }

    data := struct {
        Title    string
        Discover *models.Discover
    }{
        Title:    "View Discovery",
        Discover: discover,
    }

    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
    }
}

func (h *DiscoverHandler) CreateDiscover(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
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

        tmpl, ok := h.Templates["discovers/form"]
        if !ok {
            http.Error(w, "Template not found: discovers/form", http.StatusInternalServerError)
            return
        }

        data := struct {
            Title     string
            Discover  *models.Discover
            Countries []models.Country
            Diseases  []models.Disease
        }{
            Title:     "Create Discovery",
            Discover:  &models.Discover{},
            Countries: countries,
            Diseases:  diseases,
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

        cname := r.FormValue("cname")
        diseaseCode := r.FormValue("disease_code")
        firstEncDateStr := r.FormValue("first_enc_date")

        if cname == "" || diseaseCode == "" || firstEncDateStr == "" {
            http.Error(w, "All fields are required", http.StatusBadRequest)
            return
        }

        firstEncDate, err := time.Parse("2006-01-02", firstEncDateStr)
        if err != nil {
            http.Error(w, "Invalid date format. Use YYYY-MM-DD.", http.StatusBadRequest)
            return
        }

        discover := &models.Discover{
            CName:        cname,
            DiseaseCode:  diseaseCode,
            FirstEncDate: firstEncDate,
        }

        if err := models.CreateDiscover(h.DB, discover); err != nil {
            http.Error(w, "Error creating discovery: "+err.Error(), http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/discovers", http.StatusSeeOther)
    }
}

func (h *DiscoverHandler) UpdateDiscover(w http.ResponseWriter, r *http.Request) {
    cname := r.URL.Query().Get("cname")
    diseaseCode := r.URL.Query().Get("disease_code")
    if cname == "" || diseaseCode == "" {
        http.Error(w, "Missing country name or disease code", http.StatusBadRequest)
        return
    }

    if r.Method == "GET" {
        discover, err := models.GetDiscover(h.DB, cname, diseaseCode)
        if err != nil {
            http.Error(w, "Error fetching discovery: "+err.Error(), http.StatusInternalServerError)
            return
        }
        if discover == nil {
            http.NotFound(w, r)
            return
        }

        tmpl, ok := h.Templates["discovers/form"]
        if !ok {
            http.Error(w, "Template not found: discovers/form", http.StatusInternalServerError)
            return
        }

        data := struct {
            Title    string
            Discover *models.Discover
        }{
            Title:    "Edit Discovery",
            Discover: discover,
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

        firstEncDateStr := r.FormValue("first_enc_date")
        if firstEncDateStr == "" {
            http.Error(w, "First Encounter Date is required", http.StatusBadRequest)
            return
        }

        firstEncDate, err := time.Parse("2006-01-02", firstEncDateStr)
        if err != nil {
            http.Error(w, "Invalid date format. Use YYYY-MM-DD.", http.StatusBadRequest)
            return
        }

        discover := &models.Discover{
            CName:        cname,
            DiseaseCode:  diseaseCode,
            FirstEncDate: firstEncDate,
        }

        if err := models.UpdateDiscover(h.DB, discover); err != nil {
            http.Error(w, "Error updating discovery: "+err.Error(), http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/discovers", http.StatusSeeOther)
    }
}

func (h *DiscoverHandler) DeleteDiscover(w http.ResponseWriter, r *http.Request) {
    cname := r.URL.Query().Get("cname")
    diseaseCode := r.URL.Query().Get("disease_code")
    if cname == "" || diseaseCode == "" {
        http.Error(w, "Missing country name or disease code", http.StatusBadRequest)
        return
    }

    if err := models.DeleteDiscover(h.DB, cname, diseaseCode); err != nil {
        http.Error(w, "Error deleting discovery: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/discovers", http.StatusSeeOther)
}
