package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware,logger)
	r.Get("/", func(w http.ResponseWirter, _ *http.Request){
		tmpl, _ := template.New("").PraseFile("templates/index.html")
		tmpl.ExcuteTemplate(w, "Base", nil)
	})
	http.ListenAndServe("localhost:3000", r)
}