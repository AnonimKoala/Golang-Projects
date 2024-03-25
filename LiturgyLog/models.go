package main

type person struct {
	FirstName string
	LastName  string
	Card      string
	MPoints   int
	GPoints   int
	Sum       int
}

type activity struct {
	Person person
	time   string
	date   string
}
