package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

//globals variables
var client *redis.Client
var store = sessions.NewCookieStore([]byte("t0p-s3cr3t"))
var templates *template.Template

func main() {
	client = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	templates = template.Must(template.ParseGlob("templates/*.html"))
	r := mux.NewRouter()
	r.HandleFunc("/contact", contactHandler).Methods("GET")
	r.HandleFunc("/about", aboutHandler).Methods("GET")
	r.HandleFunc("/", indexGetHandler).Methods("GET")
	r.HandleFunc("/", indexPostHandler).Methods("POST")
	r.HandleFunc("/login", loginGetHandler).Methods("GET")
	r.HandleFunc("/login", loginPostHandler).Methods("POST")
	r.HandleFunc("/test", testGetHandler).Methods("GET")
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8001", nil))
}

//request index page handle
func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	comments, err := client.LRange(client.Context(), "comments", 0, 10).Result()
	if err != nil {
		return
	}
	templates.ExecuteTemplate(w, "index.html", comments)
}

//request index page POST handle
func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//get the comment in html tag comment
	comment := r.PostForm.Get("comment")
	//push the comment to the comments list
	client.LPush(client.Context(), "comments", comment)
	//redirect to / when the submit form
	http.Redirect(w, r, "/", 302)
}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	session, _ := store.Get(r, "session")
	session.Values["username"] = username
	session.Save(r, w)
}

func testGetHandler(w http.ResponseWriter, r *http.Request) {
	// grab the session
	session, _ := store.Get(r, "session")
	// grab the username from the session object
	untyped, ok := session.Values["username"]
	// Verify that the username was actually there
	if !ok {
		return
	}
	// quick type assertion because the object store the data as
	// empty interface
	username, ok := untyped.(string)
	if !ok {
		return
	}
	// if pass in all these tests, write a
	// byte array with the username
	w.Write([]byte(username))
}

//request contact page handle
func contactHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "contact.html", "This is the contact page!")
}

//request about page handle
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "about.html", "This is the about page!")
}
