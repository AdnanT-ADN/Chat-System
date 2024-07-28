package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/gorilla/sessions"
)

// server variables
var session_store = sessions.NewCookieStore([]byte(uuid.NewString())) // TODO Implement better way to generate the session key
var server_users = make(map[string]bool)
var mu sync.Mutex

const (
	LoginPage = "./pages/chatRoomLogin.html"
	HomePage  = "./pages/home.html"
)

// auxiliary methods
func navigateToPage(page_path string, w http.ResponseWriter, data any) {
	page := template.Must(template.ParseFiles(page_path))
	page.Execute(w, data)
}

func setSessionData(w http.ResponseWriter, r *http.Request, key string, value any) {
	// create session
	session, _ := session_store.Get(r, "user-session")

	// set username
	session.Values[key] = value

	// save session data
	session.Save(r, w)
}

func getSessionData(r *http.Request, key string) any {
	// retrieve session
	session, _ := session_store.Get(r, "user-session")

	if value, ok := session.Values[key]; ok {
		return value
	} else {
		log.Printf("Failed to obtain value for specified key: %s\n", key)
	}
	return nil
}

// end points
func attemptUserLogin(w http.ResponseWriter, r *http.Request) {
	// parse form
	err := r.ParseForm()
	if err != nil {
		navigateToPage(LoginPage, w, struct {
			Message string
		}{
			Message: "Error when attempting to parse form",
		})
	}

	// check if user already created
	mu.Lock()
	defer mu.Unlock()

	var username string = r.PostFormValue("username")
	_, user_exists := server_users[username]
	// if created return error
	if user_exists {
		navigateToPage(LoginPage, w, struct {
			Message string
		}{
			Message: fmt.Sprintf("Sorry, the username '%s' is already taken", username),
		})
	} else {
		// Create session and navigate to new page
		setSessionData(w, r, "username", username)
		server_users[username] = true
		homePage(w, r)
	}
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// navigate to login page
		navigateToPage(LoginPage, w, nil)
	case http.MethodPost:
		// check credentials
		attemptUserLogin(w, r)
	default:
		http.Error(w, "Error an unsupported method was used to reach the endpoint.", http.StatusMethodNotAllowed)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	user := getSessionData(r, "username")

	if username, ok := user.(string); ok {
		navigateToPage(HomePage, w, struct {
			Username string
		}{
			Username: username,
		})
	} else {
		navigateToPage(LoginPage, w, nil)
	}
}

func main() {
	// routers
	router := chi.NewRouter()

	loginRouter := chi.NewRouter()
	loginRouter.Get("/", loginPage)
	loginRouter.Post("/", loginPage)

	// end points
	router.HandleFunc("/", homePage)
	router.Mount("/login", loginRouter)

	// start server
	http.ListenAndServe(":8080", router)
}
