package main

import (
	"fmt"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"time"
)

func attendance(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Week  []activity
		Month []activity
	}

	var d data
	d.Month = queryActivities("thisMonth")
	d.Week = queryActivities("thisWeek")

	err := tpl.ExecuteTemplate(w, "attendance.gohtml", d)
	if err != nil {
		log.Fatal(err)
	}
}

func rank(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	status := mux.Vars(req)["status"]

	type PersonWithIndex struct {
		Index  int
		Person person
	}
	type data struct {
		Top3  []PersonWithIndex
		Other []PersonWithIndex
		IsL   bool
	}

	var persons []person
	var d data

	if status == "l" {
		persons, _ = queryRank()
		d.IsL = true
	} else if status == "a" {
		_, persons = queryRank()
		d.IsL = false
	}

	if len(persons) >= 3 {
		d.Top3 = make([]PersonWithIndex, 3)
		for i := 0; i < 3; i++ {
			d.Top3[i] = PersonWithIndex{i + 1, persons[i]}
		}
		d.Other = make([]PersonWithIndex, len(persons)-3)
		for i := 3; i < len(persons); i++ {
			d.Other[i-3] = PersonWithIndex{i + 1, persons[i]}
		}
	} else {
		d.Top3 = make([]PersonWithIndex, len(persons))
		for i := 0; i < len(persons); i++ {
			d.Top3[i] = PersonWithIndex{i + 1, persons[i]}
		}
		d.Other = make([]PersonWithIndex, 0)
	}

	if len(d.Top3) >= 3 {
		d.Top3[0], d.Top3[1], d.Top3[2] = d.Top3[1], d.Top3[0], d.Top3[2]
	}

	err := tpl.ExecuteTemplate(w, "rank.gohtml", d)
	if err != nil {
		log.Fatal(err)
	}
}

func home(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

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

func index(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/home", http.StatusSeeOther)
	} else {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
	}
}

func api(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	if u.Role != "api" {
		http.Error(w, "You are not allowed to this action!", http.StatusForbidden)
		return
	}

	vars := mux.Vars(req)
	card, isCard := vars["card"]
	recordTime, isTime := vars["time"]
	date, isDate := vars["date"]

	if !isCard {
		updateDb()
	}
	if isCard && !isTime && !isDate {
		addRowToDb(card, "", "")
	}
	if isCard && isTime && !isDate {
		addRowToDb(card, recordTime, "")
	}
	if isCard && isTime && isDate {
		addRowToDb(card, recordTime, date)
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Fatal(err)
	}
}

func signup(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	var u user
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		un := req.FormValue("username")
		p := req.FormValue("password")
		f := req.FormValue("firstname")
		l := req.FormValue("lastname")
		r := req.FormValue("role")
		// username taken?
		if _, ok := dbUsers[un]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}
		// create session
		sID := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		c.MaxAge = sessionLength
		http.SetCookie(w, c)
		dbSessions[c.Value] = session{un, time.Now()}
		// store user in dbUsers
		hash := hashWithMD5(p)
		//bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		//if err != nil {
		//	http.Error(w, "Internal server error", http.StatusInternalServerError)
		//	return
		//}
		u = user{un, hash, f, l, r}
		dbUsers[un] = u
		queryNewUser(u)

		// redirect
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "signup.gohtml", u)
}

func login(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	var u user
	// process form submission
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		p := req.FormValue("password")
		// is there a username?
		u, ok := dbUsers[un]
		if !ok {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		//err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		//if err != nil {
		//	http.Error(w, "Username and/or password do not match", http.StatusForbidden)
		//	return
		//}

		// does the entered password match the stored password?
		if hashWithMD5(p) != u.Password {
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}

		// create session
		sID := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		c.MaxAge = sessionLength
		http.SetCookie(w, c)
		dbSessions[c.Value] = session{un, time.Now()}
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "login.gohtml", u)
}

func logout(w http.ResponseWriter, req *http.Request) {
	c, _ := req.Cookie("session")
	// delete the session
	delete(dbSessions, c.Value)
	// remove the cookie
	c = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	// clean up dbSessions
	if time.Now().Sub(dbSessionsCleaned) > (time.Second * 30) {
		go cleanSessions()
	}

	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func authorized(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !alreadyLoggedIn(w, r) {
			//http.Error(w, "not logged in", http.StatusUnauthorized)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return // don't call original handler
		}
		h.ServeHTTP(w, r)
	})
}

func meeting(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	if !alreadyLoggedIn(w, req) {
		fmt.Println("not logged in")
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	if u.Role != "admin" {
		fmt.Println(u.Role)
		http.Error(w, "You must be admin to enter the meetings", http.StatusForbidden)
		return
	}

	vars := mux.Vars(req)
	status := vars["status"]

	var d struct {
		Person []person
		IsL    bool
	}

	if status == "l" {
		d.Person, _ = queryRank()
		d.IsL = true
	} else if status == "a" {
		_, d.Person = queryRank()
		d.IsL = false
	} else {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	err := tpl.ExecuteTemplate(w, "meeting.gohtml", d)
	if err != nil {
		log.Fatal(err)
	}
}

func justification(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if u.Role != "admin" {
		http.Error(w, "You must be admin to enter the justifications", http.StatusForbidden)
		return
	}

	var p []person
	{
		var tmp1, tmp2 []person
		tmp1, tmp2 = queryRank()
		p = append(tmp1, tmp2...)
	}

	err := tpl.ExecuteTemplate(w, "justification.gohtml", p)
	if err != nil {
		log.Fatal(err)
	}
}

func addAttendance(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {

		return
	}

	var p []person
	{
		var tmp1, tmp2 []person
		tmp1, tmp2 = queryRank()
		p = append(tmp1, tmp2...)
	}
	err := tpl.ExecuteTemplate(w, "addAttendance.gohtml", p)
	if err != nil {
		log.Fatal(err)
	}
}
