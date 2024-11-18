package handlers

import (
	"database/sql"
	"html/template"
	"myapp/models"
	"net/http"
	"strconv"
)

type UserHandler struct {
	DB        *sql.DB
	Templates map[string]*template.Template
}

func NewUserHandler(db *sql.DB, templates map[string]*template.Template) *UserHandler {
	return &UserHandler{
		DB:        db,
		Templates: templates,
	}
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetAllUsers(h.DB)
	if err != nil {
		http.Error(w, "Error fetching users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, ok := h.Templates["users/list"]
	if !ok {
		http.Error(w, "Template not found: users/list", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
		Users []models.User
	}{
		Title: "Users",
		Users: users,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *UserHandler) ViewUser(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Missing user email", http.StatusBadRequest)
		return
	}

	user, err := models.GetUser(h.DB, email)
	if err != nil {
		http.Error(w, "Error fetching user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.NotFound(w, r)
		return
	}

	tmpl, ok := h.Templates["users/view"]
	if !ok {
		http.Error(w, "Template not found: users/view", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
		User  *models.User
	}{
		Title: "View User",
		User:  user,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, ok := h.Templates["users/form"]
		if !ok {
			http.Error(w, "Template not found: users/form", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title string
			User  *models.User
		}{
			Title: "Create User",
			User:  &models.User{},
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

		user := &models.User{
			Email:   r.FormValue("email"),
			Name:    r.FormValue("name"),
			Surname: r.FormValue("surname"),
			CName:   r.FormValue("cname"),
		}

		if salaryStr := r.FormValue("salary"); salaryStr != "" {
			salary, err := strconv.ParseInt(salaryStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid salary", http.StatusBadRequest)
				return
			}
			user.Salary = sql.NullInt64{Int64: salary, Valid: true}
		} else {
			user.Salary = sql.NullInt64{Valid: false}
		}

		if phone := r.FormValue("phone"); phone != "" {
			user.Phone = sql.NullString{String: phone, Valid: true}
		} else {
			user.Phone = sql.NullString{Valid: false}
		}

		if err := models.CreateUser(h.DB, user); err != nil {
			http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/users", http.StatusSeeOther)
	}
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Missing user email", http.StatusBadRequest)
		return
	}

	if r.Method == "GET" {
		user, err := models.GetUser(h.DB, email)
		if err != nil {
			http.Error(w, "Error fetching user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.NotFound(w, r)
			return
		}

		tmpl, ok := h.Templates["users/form"]
		if !ok {
			http.Error(w, "Template not found: users/form", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title string
			User  *models.User
		}{
			Title: "Edit User",
			User:  user,
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

		user := &models.User{
			Email:   email,
			Name:    r.FormValue("name"),
			Surname: r.FormValue("surname"),
			CName:   r.FormValue("cname"),
		}

		if salaryStr := r.FormValue("salary"); salaryStr != "" {
			salary, err := strconv.ParseInt(salaryStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid salary", http.StatusBadRequest)
				return
			}
			user.Salary = sql.NullInt64{Int64: salary, Valid: true}
		} else {
			user.Salary = sql.NullInt64{Valid: false}
		}

		if phone := r.FormValue("phone"); phone != "" {
			user.Phone = sql.NullString{String: phone, Valid: true}
		} else {
			user.Phone = sql.NullString{Valid: false}
		}

		if err := models.UpdateUser(h.DB, user); err != nil {
			http.Error(w, "Error updating user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/users", http.StatusSeeOther)
	}
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Missing user email", http.StatusBadRequest)
		return
	}

	if err := models.DeleteUser(h.DB, email); err != nil {
		http.Error(w, "Error deleting user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
