package routes

import (
	"github.com/StephanUllmann/hypermedia-templ/controllers"
	"github.com/go-chi/chi/v5"
)

func ContactsRouter() chi.Router {

		r := chi.NewRouter()
	
		// getting contacts
		r.Get("/", controllers.GetContacts)
		r.Get("/{id}", controllers.GetContact)

		// creating contacts
		r.Get("/new", controllers.GetNewContact)
		r.Post("/new", controllers.PostNewContact)

		// editing contacts
		r.Get("/{id}/edit", controllers.GetEditContact)
		r.Post("/{id}/edit", controllers.PostEditContact)

		// deleting contacts
		r.Post("/{id}/delete", controllers.DeleteContact)

	return r
}