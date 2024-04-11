package main

import "database/sql"

type person struct {
	FirstName string
	LastName  string
	Card      string
	MPoints   int
	GPoints   int
	Sum       int
	Rank      int
	IsL       bool
}

type activity struct {
	Person person
	Time   string
	Date   string
}

type dbPerson struct {
	firstName sql.NullString
	lastName  sql.NullString
}

type dbActivity struct {
	person dbPerson
	time   sql.NullString
	date   sql.NullString
}
