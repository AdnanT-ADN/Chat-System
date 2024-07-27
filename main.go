package main

import (
	"encoding/json"
	"fmt"
	"log"

	"html/template"
	"net/http"

	// "sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/gorilla/websocket"
)

type Client struct {
	Username   string
	Connection *websocket.Conn
}

var server_active_users = make(map[string]*Client)

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Auxiliary Functions
func sendJsonResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// End Points

func userLogin(w http.ResponseWriter, r *http.Request) {
	pageTemplate := template.Must(template.ParseFiles("./pages/chatRoomLogin.html"))
	pageTemplate.Execute(w, nil)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	// parse form
	err := r.ParseForm()
	if err != nil {
		w.Write([]byte("Error parsing form"))
	}

	// get username
	username := r.PostFormValue("username")

	_, user_already_exists := server_active_users[username]
	if user_already_exists {
		w.Write([]byte(fmt.Sprintf("User %s already exists", username)))
	} else {
		log.Println("Adding new user")
		connection, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			w.Write([]byte("Error was not able to establish websocket connection"))
		}
		server_active_users[username] = &Client{
			Username:   username,
			Connection: connection,
		}
	}

}

func main() {

	// define chi router
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// define paths
	router.Get("/", userLogin)
	router.Post("/add_user", addUser)

	// start server
	http.ListenAndServe(":8080", router)

}
