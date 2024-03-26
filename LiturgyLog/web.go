package main

import (
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	type data struct {
		ACount      int
		LCount      int
		CMonthCount int
		Top5L       []person
		Top5A       []person
		Last10      []activity
	}

	var d data
	d.Top5L, d.Top5A = queryRank()
	d.ACount = len(d.Top5A)
	d.LCount = len(d.Top5L)
	d.CMonthCount = len(queryActivities("thisMonth"))
	d.Last10 = queryActivities("last10")

	if len(d.Top5A) >= 5 {
		d.Top5A = d.Top5A[:5]
	} else {
		missing := 5 - len(d.Top5A)
		d.Top5A = append(d.Top5A, make([]person, missing)...)
	}

	if len(d.Top5L) >= 5 {
		d.Top5L = d.Top5L[:5]
	} else {
		missing := 5 - len(d.Top5L)
		d.Top5L = append(d.Top5L, make([]person, missing)...)
	}

	err := tpl.ExecuteTemplate(w, "home.gohtml", d)
	if err != nil {
		log.Fatal(err)
	}
}
