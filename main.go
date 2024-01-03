package main

import (
	"log"
	"net/http"
	"os"

	"github.com/StephanUllmann/hypermedia-templ/db"
	"github.com/StephanUllmann/hypermedia-templ/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	// "github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)



func main() {

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	db.Init()
	

	r := chi.NewRouter()

	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5, "text/html", "text/css", "text/javascript", "application/json", "text/plain"))

	r.Get("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))).ServeHTTP)

	r.Get("/", http.RedirectHandler("/contacts", http.StatusMovedPermanently).ServeHTTP)
	r.Get("/contact", http.RedirectHandler("/contacts", http.StatusMovedPermanently).ServeHTTP)


	r.Mount("/contacts", routes.ContactsRouter())

	port := ":" + os.Getenv("PORT")
	log.Printf("Server started on port %v\n", port)

	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal(err)
	}
}