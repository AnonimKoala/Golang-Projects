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

func rankL(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	type PersonWithIndex struct {
		index  int
		Person person
	}
	type data struct {
		top3  []PersonWithIndex
		other []PersonWithIndex
	}

	l, _ := queryRank()
	var d data
	if len(l) >= 3 {
		d.top3 = make([]PersonWithIndex, 3)
		for i := 0; i < 3; i++ {
			d.top3[i] = PersonWithIndex{i + 1, l[i]}
		}
		d.other = make([]PersonWithIndex, len(l)-3)
		for i := 3; i < len(l); i++ {
			d.other[i-3] = PersonWithIndex{i + 1, l[i]}
		}
	} else {
		d.top3 = make([]PersonWithIndex, len(l))
		for i := 0; i < len(l); i++ {
			d.top3[i] = PersonWithIndex{i + 1, l[i]}
		}
		d.other = make([]PersonWithIndex, 0)
	}

	err := tpl.ExecuteTemplate(w, "rankL.gohtml", d)
	if err != nil {
		log.Fatal(err)
	}
}

func rankA(w http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(w, req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	type PersonWithIndex struct {
		index  int
		person person
	}
	type data struct {
		top3  []PersonWithIndex
		other []PersonWithIndex
	}

	_, a := queryRank()
	var d data
	if len(a) >= 3 {
		d.top3 = make([]PersonWithIndex, 3)
		for i := 0; i < 3; i++ {
			d.top3[i] = PersonWithIndex{i + 1, a[i]}
		}
		d.other = make([]PersonWithIndex, len(a)-3)
		for i := 3; i < len(a); i++ {
			d.other[i-3] = PersonWithIndex{i + 1, a[i]}
		}
	} else {
		d.top3 = make([]PersonWithIndex, len(a))
		for i := 0; i < len(a); i++ {
			d.top3[i] = PersonWithIndex{i + 1, a[i]}
		}
		d.other = make([]PersonWithIndex, 0)
	}

	err := tpl.ExecuteTemplate(w, "rankL.gohtml", d)
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
