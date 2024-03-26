package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func scanRankRows(target *[]person, rows *sql.Rows) {
	for rows.Next() {
		var p person
		err := rows.Scan(&p.FirstName, &p.LastName, &p.Card, &p.MPoints, &p.GPoints, &p.Sum)
		if err != nil {
			log.Fatal(err)
		}

		*target = append(*target, p)
	}
}

func scanActivityRows(target *[]activity, rows *sql.Rows) {
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
		*target = append(*target, a)
	}

}

func queryRank() ([]person, []person) {
	baseQuery := "SELECT osoby.Imie, osoby.Nazwisko, osoby.karta, osoby.zbiorki, osoby.Punkty, osoby.suma FROM `osoby` WHERE stopien = %s ORDER BY suma DESC"
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
