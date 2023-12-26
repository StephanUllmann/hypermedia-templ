package main

import (
	"net/http"

	"github.com/StephanUllmann/hypermedia-templ/db"
	"github.com/StephanUllmann/hypermedia-templ/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)



func main() {

	db.Init()
	

	r := chi.NewRouter()

	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", http.RedirectHandler("/contacts", http.StatusMovedPermanently).ServeHTTP)
	r.Get("/contact", http.RedirectHandler("/contacts", http.StatusMovedPermanently).ServeHTTP)


	r.Mount("/contacts", routes.ContactsRouter())

	http.ListenAndServe(":3000", r)
}