package main

import (
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
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

	err := tpl.ExecuteTemplate(w, "frekwencja.gohtml", d)
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

func bar(w http.ResponseWriter, req *http.Request) {
	u := getUser(w, req)
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	if u.Role != "007" {
		http.Error(w, "You must be 007 to enter the bar", http.StatusForbidden)
		return
	}

	tpl.ExecuteTemplate(w, "bar.gohtml", u)
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
		bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		u = user{un, bs, f, l, r}
		dbUsers[un] = u
		queryNewUser(u)

		// redirect
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	showSessions() // for demonstration purposes
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
		// does the entered password match the stored password?
		err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err != nil {
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
	showSessions() // for demonstration purposes
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
