package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template
var db *sql.DB

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))

	var err error
	db, err = sql.Open("mysql", "root:secret@(localhost:3306)/db?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	r.HandleFunc("/", index)
	r.HandleFunc("/attendance", attendance)
	http.ListenAndServe(":8080", r)
}
