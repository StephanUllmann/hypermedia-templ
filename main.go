package main

import (
	"log"
	"net/http"
	"os"

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

	r.Get("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))).ServeHTTP)

	r.Get("/", http.RedirectHandler("/contacts", http.StatusMovedPermanently).ServeHTTP)
	r.Get("/contact", http.RedirectHandler("/contacts", http.StatusMovedPermanently).ServeHTTP)


	r.Mount("/contacts", routes.ContactsRouter())

	port := ":" + os.Getenv("PORT")
 
	http.ListenAndServe(":3000", r)
	log.Println("Server started on port 3000")
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal(err)
	}
}