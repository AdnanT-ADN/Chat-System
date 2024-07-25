package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func chatRoomLogin(w http.ResponseWriter, r *http.Request) {
	template := template.Must(template.ParseFiles("./pages/chatRoomLogin.html"))
	template.Execute(w, nil)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", chatRoomLogin)
	log.Fatal(http.ListenAndServe(":8000", r))
}
