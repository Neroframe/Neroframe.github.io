package handlers

import (
    "html/template"
    "net/http"
)

type DashboardHandler struct {
    Templates map[string]*template.Template
}

func NewDashboardHandler(templates map[string]*template.Template) *DashboardHandler {
    return &DashboardHandler{
        Templates: templates,
    }
}

func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
    tmpl, ok := h.Templates["dashboard"]
    if !ok {
        http.Error(w, "Template not found", http.StatusInternalServerError)
        return
    }

    data := struct {
        Title string
    }{
        Title: "Dashboard",
    }

    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
    }
}
