package routes

import (
	"github.com/StephanUllmann/hypermedia-templ/controllers"
	"github.com/go-chi/chi/v5"
)

func ContactsRouter() chi.Router {

		r := chi.NewRouter()
	
		r.Get("/", controllers.GetContacts)
		r.Get("/{id}", controllers.GetContact)
	return r
}