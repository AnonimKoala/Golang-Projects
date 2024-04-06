package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

func hashWithMD5(input string) string {
	hasher := md5.New()
	_, err := hasher.Write([]byte(input))
	if err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func scanRankRows(target *[]person, rows *sql.Rows) {
	for rows.Next() {
		var p person
		var status sql.NullString
		err := rows.Scan(&p.Rank, &p.FirstName, &p.LastName, &p.Card, &p.MPoints, &p.GPoints, &p.Sum, &status)
		if err != nil {
			log.Fatal(err)
		}

		if status.String == "M" {
			p.IsL = false
		} else {
			p.IsL = true
		}

		*target = append(*target, p)
	}
}

func scanActivityRows(target *[]activity, rows *sql.Rows) {
	for rows.Next() {
		var dA dbActivity
		err := rows.Scan(
			&dA.person.firstName,
			&dA.person.lastName,
			&dA.time,
			&dA.date,
		)

		var a activity
		a.Person.FirstName = dA.person.firstName.String
		a.Person.LastName = dA.person.lastName.String
		a.Time = dA.time.String

		ts := dA.date.String
		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			log.Fatalln(err)
		}

		a.Date = t.Format("02-01-2006")

		if err != nil {
			log.Fatal(err)

		}
		*target = append(*target, a)
	}

}

func queryRank() ([]person, []person) {
	baseQuery := "SELECT ROW_NUMBER() OVER(ORDER BY suma DESC) AS RowNum, osoby.Imie, osoby.Nazwisko, osoby.karta, osoby.zbiorki, osoby.Punkty, osoby.suma, osoby.stopien FROM `osoby` WHERE stopien = %s ORDER BY suma DESC"
	rA, err := db.Query(fmt.Sprintf(baseQuery, "\"M\""))
	if err != nil {
		log.Fatal(err)

	}

	rL, err := db.Query(fmt.Sprintf(baseQuery, "\"L\""))
	if err != nil {
		log.Fatal(err)

	}

	var l, a []person
	scanRankRows(&l, rL)
	scanRankRows(&a, rA)

	defer rA.Close()
	defer rL.Close()

	return l, a
}

func queryActivities(mode string) []activity {
	baseQuery := "SELECT osoby.Imie, osoby.Nazwisko, odczyt.msza_godz, odczyt.data_odczytu FROM osoby, odczyt WHERE osoby.karta = odczyt.karta"
	var query string

	switch mode {
	case "last10":
		query = fmt.Sprintf("%s ORDER BY odczyt.data_odczytu DESC, odczyt.msza_godz DESC LIMIT 10", baseQuery)
	case "thisMonth":
		query = fmt.Sprintf("%s AND odczyt.data_odczytu > CURRENT_DATE - INTERVAL 1 MONTH ORDER BY odczyt.data_odczytu DESC, odczyt.msza_godz DESC", baseQuery)
	case "thisWeek":
		query = fmt.Sprintf("%s AND odczyt.data_odczytu > CURRENT_DATE - INTERVAL 1 WEEK ORDER BY odczyt.data_odczytu DESC, odczyt.msza_godz DESC", baseQuery)
	default:
		log.Fatalln("Incorrect queryActivities() mode")
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var a []activity
	scanActivityRows(&a, rows)

	return a
}

func queryUsers() map[string]user {
	q := "SELECT users.username, users.password, users.role, users.fn, users.ln FROM `users`"
	rows, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var dbUsers = map[string]user{}
	var un, pwd, r, fn, ln sql.NullString

	for rows.Next() {
		err := rows.Scan(&un, &pwd, &r, &fn, &ln)
		if err != nil {
			log.Fatalln(err)
		}

		dbUsers[un.String] = user{un.String, pwd.String, fn.String, ln.String, r.String}
	}

	return dbUsers
}

func queryNewUser(u user) {
	q := fmt.Sprintf("INSERT INTO users (users.username, users.password, users.role, users.fn, users.ln) VALUES ('%s','%s','%s','%s','%s')", u.UserName, u.Password, u.Role, u.First, u.Last)
	stmt, err := db.Prepare(q)
	if err != nil {
		log.Fatalln(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Fatalln(err)
	}
}
