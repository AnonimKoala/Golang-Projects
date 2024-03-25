package main

import (
	"database/sql"
	"log"
)

func scanRankRows(target []person, rows *sql.Rows) {
	for rows.Next() {
		var p person
		err := rows.Scan(&p.FirstName, &p.LastName, &p.Card, &p.MPoints, &p.GPoints, &p.Sum)
		if err != nil {
			log.Fatal(err)

		}

		target = append(target, p)
	}
}

func scanActivityRows(target []activity, rows *sql.Rows) {
	for rows.Next() {
		var a activity

		err := rows.Scan(
			&a.Person.FirstName,
			&a.Person.LastName,
			&a.time,
			&a.date,
		)

		if err != nil {
			log.Fatal(err)

		}
		target = append(target, a)
	}

}

func queryRank() ([]person, []person) {
	q := "SELECT osoby.Imie, osoby.Nazwisko, osoby.karta, osoby.zbiorki, osoby.Punkty, osoby.suma FROM `osoby` WHERE stopien = '?' ORDER BY suma DESC"

	rA, err := db.Query(q, "M")
	if err != nil {
		log.Fatal(err)

	}

	rL, err := db.Query(q, "L")
	if err != nil {
		log.Fatal(err)

	}

	defer rA.Close()
	defer rL.Close()

	var l, a []person
	scanRankRows(l, rL)
	scanRankRows(a, rA)

	return l, a
}

func queryActivities(mode string) []activity {
	q := "SELECT osoby.Imie, osoby.Nazwisko, odczyt.msza_godz, odczyt.data_odczytu FROM osoby, odczyt WHERE osoby.karta = odczyt.karta ? ORDER BY odczyt.data_odczytu, odczyt.msza_godz ?"

	var err error
	var rows *sql.Rows

	defer rows.Close()

	switch mode {
	case "last10":
		rows, err = db.Query(q, "", "LIMIT 10")
	case "thisMonth":
		rows, err = db.Query(q, "\"and odczyt.data_odczytu > CURRENT_DATE - INTERVAL 1 MONTH\"", "")
	case "thisWeek":
		rows, err = db.Query(q, "\"and odczyt.data_odczytu > CURRENT_DATE - INTERVAL 1 WEEK\"", "")
	}
	if err != nil {
		log.Fatal(err)

	}

	var a []activity
	scanActivityRows(a, rows)

	return a
}
