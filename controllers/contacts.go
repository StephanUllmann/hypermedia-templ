package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/StephanUllmann/hypermedia-templ/db"
	models "github.com/StephanUllmann/hypermedia-templ/model"
	"github.com/StephanUllmann/hypermedia-templ/templates"
	"github.com/go-chi/chi/v5"
)

func GetContacts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	fmt.Printf("Value of q: %s\n", q)
	var rows *sql.Rows
	var err error
	if len(q) == 0 {
		rows, err = db.DB.Query("SELECT * FROM contacts")
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
	} else {
		query := "%" + q + "%"
		rows, err = db.DB.Query("SELECT * FROM contacts WHERE first LIKE ? OR last LIKE ?", query, query)
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

	fmt.Printf("Data: %v\n", data)
	// w.Write([]byte("Working"))
	// templ.Handler(Layout("Contacts"))
	// templates.Layout("Contacts").Render(r.Context(), w)
	templates.Contacts(data).Render(r.Context(), w)
}

func GetContact(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	data := models.Contact{}
	row := db.DB.QueryRow("SELECT * FROM contacts WHERE id = ?", id)
	switch err := row.Scan(&data.ID, &data.First, &data.Last, &data.Phone, &data.Email); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
		case nil:
			templates.Contact(data).Render(r.Context(), w)
		default:
			panic(err)
		}
}