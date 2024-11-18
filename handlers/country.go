package handlers

import (
	"database/sql"
	"html/template"
	"myapp/models"
	"net/http"
	"strconv"
)

type CountryHandler struct {
	DB        *sql.DB
	Templates map[string]*template.Template
}

// Constructor
func NewCountryHandler(db *sql.DB, templates map[string]*template.Template) *CountryHandler {
	return &CountryHandler{
		DB:        db,
		Templates: templates,
	}
}

func (h *CountryHandler) ListCountries(w http.ResponseWriter, r *http.Request) {
	countries, err := models.GetAllCountries(h.DB)
	if err != nil {
		http.Error(w, "Error fetching countries: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, ok := h.Templates["countries/list"]
	if !ok {
		http.Error(w, "Template not found: countries/list", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title     string
		Countries []models.Country
	}{
		Title:     "Countries",
		Countries: countries,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *CountryHandler) ViewCountry(w http.ResponseWriter, r *http.Request) {
	cname := r.URL.Query().Get("cname")
	if cname == "" {
		http.Error(w, "Missing country name", http.StatusBadRequest)
		return
	}

	country, err := models.GetCountry(h.DB, cname)
	if err != nil {
		http.Error(w, "Error fetching country: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if country == nil {
		http.NotFound(w, r)
		return
	}

	tmpl, ok := h.Templates["countries/view"]
	if !ok {
		http.Error(w, "Template not found: countries/view", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title   string
		Country *models.Country
	}{
		Title:   "View Country",
		Country: country,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *CountryHandler) CreateCountry(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, ok := h.Templates["countries/form"]
		if !ok {
			http.Error(w, "Template not found: countries/form", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title   string
			Country *models.Country
		}{
			Title:   "Create Country",
			Country: &models.Country{},
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

		populationStr := r.FormValue("population")
		population, err := strconv.ParseInt(populationStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid population", http.StatusBadRequest)
			return
		}

		country := &models.Country{
			CName:      r.FormValue("cname"),
			Population: population,
		}

		if err := models.CreateCountry(h.DB, country); err != nil {
			http.Error(w, "Error creating country: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/countries", http.StatusSeeOther)
	}
}

func (h *CountryHandler) UpdateCountry(w http.ResponseWriter, r *http.Request) {
	cname := r.URL.Query().Get("cname")
	if cname == "" {
		http.Error(w, "Missing country name", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		country, err := models.GetCountry(h.DB, cname)
		if err != nil {
			http.Error(w, "Error fetching country: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if country == nil {
			http.NotFound(w, r)
			return
		}

		tmpl, ok := h.Templates["countries/form"]
		if !ok {
			http.Error(w, "Template not found: countries/form", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title   string
			Country *models.Country
		}{
			Title:   "Edit Country",
			Country: country,
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

		populationStr := r.FormValue("population")
		population, err := strconv.ParseInt(populationStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid population", http.StatusBadRequest)
			return
		}

		country := &models.Country{
			CName:      cname,
			Population: population,
		}

		if err := models.UpdateCountry(h.DB, country); err != nil {
			http.Error(w, "Error updating country: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/countries", http.StatusSeeOther)
	}
}

func (h *CountryHandler) DeleteCountry(w http.ResponseWriter, r *http.Request) {
	cname := r.URL.Query().Get("cname")
	if cname == "" {
		http.Error(w, "Missing country name", http.StatusBadRequest)
		return
	}

	if err := models.DeleteCountry(h.DB, cname); err != nil {
		http.Error(w, "Error deleting country: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/countries", http.StatusSeeOther)
}
