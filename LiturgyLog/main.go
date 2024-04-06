package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"time"
)

type user struct {
	UserName string
	Password string
	First    string
	Last     string
	Role     string
}

type session struct {
	un           string
	lastActivity time.Time
}

var tpl *template.Template
var dbUsers = map[string]user{}       // map[ userID(string)] user
var dbSessions = map[string]session{} // session ID, session
var dbSessionsCleaned time.Time

const sessionLength int = 30

var db *sql.DB

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	dbSessionsCleaned = time.Now()

	var err error
	db, err = sql.Open("mysql", "root:@(localhost:3306)/4120471_czytnik?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	dbUsers = queryUsers()
}

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))

	r.HandleFunc("/home", home)
	r.HandleFunc("/attendance", attendance)
	r.HandleFunc("/addAttendance", addAttendance)
	r.HandleFunc("/justification", justification)

	r.HandleFunc("/rank/{status}", rank)
	r.HandleFunc("/meeting/{status}", meeting)

	r.HandleFunc("/", index)
	r.HandleFunc("/bar", bar)
	r.HandleFunc("/signup", signup)
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", authorized(logout))
	r.Handle("/favicon.ico", http.NotFoundHandler())

	http.ListenAndServe(":8080", r)
}
