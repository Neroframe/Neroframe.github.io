package handlers

import (
    "database/sql"
    "html/template"
    "myapp/models"
    "net/http"
)

type PublicServantHandler struct {
    DB        *sql.DB
    Templates map[string]*template.Template
}

func NewPublicServantHandler(db *sql.DB, templates map[string]*template.Template) *PublicServantHandler {
    return &PublicServantHandler{
        DB:        db,
        Templates: templates,
    }
}

func (h *PublicServantHandler) ListPublicServants(w http.ResponseWriter, r *http.Request) {
    publicServants, err := models.GetAllPublicServants(h.DB)
    if err != nil {
        http.Error(w, "Error fetching public servants: "+err.Error(), http.StatusInternalServerError)
        return
    }

    tmpl, ok := h.Templates["public_servants/list"]
    if !ok {
        http.Error(w, "Template not found: public_servants/list", http.StatusInternalServerError)
        return
    }

    data := struct {
        Title          string
        PublicServants []models.PublicServant
    }{
        Title:          "Public Servants",
        PublicServants: publicServants,
    }

    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
    }
}

func (h *PublicServantHandler) ViewPublicServant(w http.ResponseWriter, r *http.Request) {
    email := r.URL.Query().Get("email")
    if email == "" {
        http.Error(w, "Missing public servant email", http.StatusBadRequest)
        return
    }

    publicServant, err := models.GetPublicServant(h.DB, email)
    if err != nil {
        http.Error(w, "Error fetching public servant: "+err.Error(), http.StatusInternalServerError)
        return
    }
    if publicServant == nil {
        http.NotFound(w, r)
        return
    }

    tmpl, ok := h.Templates["public_servants/view"]
    if !ok {
        http.Error(w, "Template not found: public_servants/view", http.StatusInternalServerError)
        return
    }

    data := struct {
        Title         string
        PublicServant *models.PublicServant
    }{
        Title:         "View Public Servant",
        PublicServant: publicServant,
    }

    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
    }
}

func (h *PublicServantHandler) CreatePublicServant(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        tmpl, ok := h.Templates["public_servants/form"]
        if !ok {
            http.Error(w, "Template not found: public_servants/form", http.StatusInternalServerError)
            return
        }

        data := struct {
            Title         string
            PublicServant *models.PublicServant
        }{
            Title:         "Create Public Servant",
            PublicServant: &models.PublicServant{},
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

        publicServant := &models.PublicServant{
            Email:    email,
            Department: r.FormValue("department"),
        }

        if err := models.CreatePublicServant(h.DB, publicServant); err != nil {
            http.Error(w, "Error creating public servant: "+err.Error(), http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/public_servants", http.StatusSeeOther)
    }
}

func (h *PublicServantHandler) UpdatePublicServant(w http.ResponseWriter, r *http.Request) {
    email := r.URL.Query().Get("email")
    if email == "" {
        http.Error(w, "Missing public servant email", http.StatusBadRequest)
        return
    }

    if r.Method == "GET" {
        publicServant, err := models.GetPublicServant(h.DB, email)
        if err != nil {
            http.Error(w, "Error fetching public servant: "+err.Error(), http.StatusInternalServerError)
            return
        }
        if publicServant == nil {
            http.NotFound(w, r)
            return
        }

        tmpl, ok := h.Templates["public_servants/form"]
        if !ok {
            http.Error(w, "Template not found: public_servants/form", http.StatusInternalServerError)
            return
        }

        data := struct {
            Title         string
            PublicServant *models.PublicServant
        }{
            Title:         "Edit Public Servant",
            PublicServant: publicServant,
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

        publicServant := &models.PublicServant{
            Email:    email,
            Department: r.FormValue("department"),
        }

        if err := models.UpdatePublicServant(h.DB, publicServant); err != nil {
            http.Error(w, "Error updating public servant: "+err.Error(), http.StatusInternalServerError)
            return
        }

        http.Redirect(w, r, "/public_servants", http.StatusSeeOther)
    }
}

func (h *PublicServantHandler) DeletePublicServant(w http.ResponseWriter, r *http.Request) {
    email := r.URL.Query().Get("email")
    if email == "" {
        http.Error(w, "Missing public servant email", http.StatusBadRequest)
        return
    }

    if err := models.DeletePublicServant(h.DB, email); err != nil {
        http.Error(w, "Error deleting public servant: "+err.Error(), http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/public_servants", http.StatusSeeOther)
}
