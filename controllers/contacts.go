package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/StephanUllmann/hypermedia-templ/db"
	models "github.com/StephanUllmann/hypermedia-templ/model"
	"github.com/StephanUllmann/hypermedia-templ/templates"
	"github.com/StephanUllmann/hypermedia-templ/utils"
	"github.com/a-h/templ"
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

	

	trigger := r.Header.Get("Hx-Trigger")
	if trigger == "search" {
		templates.Rows(data, page).Render(r.Context(), w)
	} else {
		utils.A.Init()
		progressWidth := "width:0%;"
		archiveProgressWidth := templ.Attributes{"style": progressWidth }
		archiverStatus := utils.A.GetStatus()
		archiverProgress := utils.A.GetProgress()
		templates.Contacts(data, q, page, archiverStatus, archiverProgress, archiveProgressWidth).Render(r.Context(), w)
	}
}

// For Lazy Loading Route
func GetContactCount(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second) // simulate slow response (for demo purposes only!
	var count int
	row := db.DB.QueryRow("SELECT COUNT(*) FROM contacts;")
	switch err := row.Scan(&count); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		w.Write([]byte(fmt.Sprintf("(%v total Contacts)", count)))
	default:
		panic(err)
	}

	
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
	if r.Header.Get("Hx-Trigger") == "delete-btn" {
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	} else {
		w.Write([]byte(""))
	}
}

// The DeleteContacts function reads the selected contact IDs from the request body, constructs a SQL
// query to delete the corresponding contacts from the database, and redirects the user to the contacts
// page.
func DeleteContacts(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
    if err != nil {
				fmt.Println(err)
		}
		pairs := strings.Split(string(body), "&")

		var ids []string
		for _, pair := range pairs {
			parts := strings.Split(pair, "=")

			var paramIsInt bool
			_, err = strconv.Atoi(parts[0])
			if err == nil {
				paramIsInt = true
			}

			if parts[0] == "selected_contact_ids" || paramIsInt  {
				_, err = strconv.Atoi(parts[1]) 
				if err == nil {
					ids = append(ids, parts[1])
				}
			}
		}
		fmt.Printf("IDs: %v\n", ids)
		query := "(" + strings.Join(ids, ", ") + ")"

		_, err = db.DB.Exec("DELETE FROM contacts WHERE id IN " + query + ";")
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

func ArchiveContacts(w http.ResponseWriter, r *http.Request) {
	utils.A.Run()
	progressWidth := "width:0%;"
	archiveProgressWidth := templ.Attributes{"style": progressWidth }
	archiverStatus := utils.A.GetStatus()
	archiverProgress := utils.A.GetProgress()
	fmt.Printf("Archiver Status: %v\n", archiverStatus)
	templates.Archive_UI(archiverStatus, archiverProgress, archiveProgressWidth).Render(r.Context(), w)
}

func ArchiveCheck(w http.ResponseWriter, r *http.Request) {
	archiverProgress := utils.A.GetProgress()
	archiverStatus := utils.A.GetStatus()
	progressWidth := fmt.Sprintf("width:%v%%;", archiverProgress*100)
	archiveProgressWidth := templ.Attributes{"style": progressWidth }
	templates.Archive_UI(archiverStatus, archiverProgress, archiveProgressWidth).Render(r.Context(), w)
}

func ArchiveDownload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Disposition", "attachment; filename=contacts.json")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	utils.A.Init()
	http.ServeFile(w, r, "public/contacts.json")
}


func ConfirmDeleteContactPopover(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	row := db.DB.QueryRow("SELECT * FROM contacts WHERE id = ?;", id)
	data := models.Contact{}
	switch err := row.Scan(&data.ID, &data.First, &data.Last, &data.Phone, &data.Email); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
		case nil:
			templates.ConfirmDeletePopover(data).Render(r.Context(), w)
		default:
			panic(err)
	}
}

func RestoreDB(w http.ResponseWriter, r *http.Request) {
	
	contacts, err := os.ReadFile("public/contacts.json")
	if err != nil {
		fmt.Println(err)
	}
	var data []models.Contact
	err = json.Unmarshal(contacts, &data)
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.DB.Exec("DELETE FROM contacts;")
	if err != nil {
		fmt.Println(err)
	}

	for _, contact := range data {
		_, err = db.DB.Exec("INSERT INTO contacts (first, last, email, phone) VALUES (?, ?, ?, ?);", contact.First, contact.Last, contact.Email, contact.Phone)
		if err != nil {
			fmt.Println(err)
		}
	}
	

	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}