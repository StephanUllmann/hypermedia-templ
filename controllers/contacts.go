package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/StephanUllmann/hypermedia-templ/db"
	models "github.com/StephanUllmann/hypermedia-templ/model"
	"github.com/StephanUllmann/hypermedia-templ/templates"
	"github.com/go-chi/chi/v5"
)

func GetContacts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	fmt.Printf("Value of q: %s\n", q)
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		fmt.Println(err)
		page = 1
	}
	offset := (page - 1) * 10
	// fmt.Printf("Value of offset: %v\n", offset)
	var rows *sql.Rows
	if len(q) == 0 {
		rows, err = db.DB.Query("SELECT * FROM contacts LIMIT 10 OFFSET ?;", offset)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
	} else {
		query := "%" + q + "%"
		rows, err = db.DB.Query("SELECT * FROM contacts WHERE first LIKE ? OR last LIKE ? OR email LIKE ? LIMIT 10 OFFSET ?;", query, query, query, offset)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
	}
	data := []models.Contact{}
	
	for rows.Next() {
		contact := models.Contact{}
		err = rows.Scan(&contact.ID, &contact.First, &contact.Last, &contact.Phone, &contact.Email)
		if err != nil {
			fmt.Println(err)
		}
		data = append(data, contact)
	}
	
	templates.Contacts(data, page).Render(r.Context(), w)
}

func GetContact(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	data := models.Contact{}
	row := db.DB.QueryRow("SELECT * FROM contacts WHERE id = ?;", id)
	switch err := row.Scan(&data.ID, &data.First, &data.Last, &data.Phone, &data.Email); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
		case nil:
			fmt.Printf("Data: %v\n", data)
			templates.Contact(data).Render(r.Context(), w)
		default:
			panic(err)
		}
}

func GetNewContact(w http.ResponseWriter, r *http.Request) {
	templates.NewContact(models.Contact{}).Render(r.Context(), w)
}

func PostNewContact(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	contact := models.Contact{
		First: r.FormValue("first"),
		Last: r.FormValue("last"),
		Email: r.FormValue("email"),
		Phone: r.FormValue("phone"),
}
	_, err := db.DB.Exec("INSERT INTO contacts (first, last, email, phone) VALUES (?, ?, ?, ?);", contact.First, contact.Last, contact.Email, contact.Phone)
	if err != nil {
		fmt.Println(err)
		templates.NewContact(models.Contact{}).Render(r.Context(), w)
	}
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func GetEditContact(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	data := models.Contact{}
	row := db.DB.QueryRow("SELECT * FROM contacts WHERE id = ?;", id)
	switch err := row.Scan(&data.ID, &data.First, &data.Last, &data.Phone, &data.Email); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
		case nil:
			templates.EditContact(data).Render(r.Context(), w)
		default:
			panic(err)
		}
}

func PostEditContact(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	r.ParseForm()
	contact := models.Contact{
		First: r.FormValue("first"),
		Last: r.FormValue("last"),
		Email: r.FormValue("email"),
		Phone: r.FormValue("phone"),
	}
	_, err := db.DB.Exec("UPDATE contacts SET first = ?, last = ?, email = ?, phone = ? WHERE id = ?;", contact.First, contact.Last, contact.Email, contact.Phone, id)
	if err != nil {
		fmt.Println(err)
		templates.EditContact(models.Contact{}).Render(r.Context(), w)
	}
	http.Redirect(w, r, fmt.Sprintf("/contacts/%v", id), http.StatusSeeOther)
}

func DeleteContact(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	_, err := db.DB.Exec("DELETE FROM contacts WHERE id = ?;", id)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func ValidateEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	fmt.Printf("Request: %v\n", email)
	addr, err := mail.ParseAddress(email)
	if err != nil {
		message := fmt.Sprintf("%v\n", err)
		templates.Message(message, "error").Render(r.Context(), w)
	} else {
		row := db.DB.QueryRow("SELECT email FROM contacts WHERE email LIKE ?;", addr.Address)
		switch err := row.Scan(&email); err {
			case sql.ErrNoRows:
				templates.Message("Email is available!", "success").Render(r.Context(), w)
			case nil:
				message := fmt.Sprintf("Email %v is already taken!", email)
				templates.Message(message, "error").Render(r.Context(), w)
			default:
				panic(err)
	}	
}
}