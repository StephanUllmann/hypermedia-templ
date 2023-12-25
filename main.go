package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", http.RedirectHandler("/contact", http.StatusMovedPermanently).ServeHTTP)

	r.Get("/contact", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Working"))
	})

	http.ListenAndServe(":3000", r)
}