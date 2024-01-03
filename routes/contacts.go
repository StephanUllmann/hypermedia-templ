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
		r.Delete("/{id}", controllers.DeleteContact)
		r.Delete("/bulk", controllers.DeleteContacts)

		// deletion popover
		r.Get("/confirm-delete-popover/{id}", controllers.ConfirmDeleteContactPopover)

		// validating email
		r.Get("/email", controllers.ValidateEmail)

		// Contact Count
		r.Get("/count", controllers.GetContactCount)

		// Archive
		r.Post("/archive", controllers.ArchiveContacts)
		r.Get("/archive", controllers.ArchiveCheck)
		r.Get("/archive/file", controllers.ArchiveDownload)

		// Restore DB
		r.Post("/restore", controllers.RestoreDB)
		

	return r
}