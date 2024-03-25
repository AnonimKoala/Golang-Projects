package main

import (
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "tmp.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}
}
