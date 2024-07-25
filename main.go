package main

import (
	"html/template"
	"log"
	"net/http"
)

func chatRoomLogin(w http.ResponseWriter, r *http.Request) {
	template := template.Must(template.ParseFiles("./chatRoomLogin.html"))
	template.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", chatRoomLogin)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
