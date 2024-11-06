package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}

var books = []book{
	{ID: "1", Title: "The Go Programming Language", Author: "Alan A. A. Donovan", Price: 35.99},
	{ID: "2", Title: "Introducing Go", Author: "Caleb Doxsey", Price: 25.99},
	{ID: "3", Title: "Go in Action", Author: "William Kennedy", Price: 45.99},
	{ID: "4", Title: "Go Programming Blueprints", Author: "Mat Ryer", Price: 40.99},
	{ID: "5", Title: "Concurrency in Go", Author: "Katherine Cox-Buday", Price: 50.99},
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBookByID)
	router.POST("/books", postBook)

	_ = router.Run("localhost:8080")

}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func postBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func getBookByID(c *gin.Context) {
	id := c.Param("id")

	for _, b := range books {
		if b.ID == id {
			c.IndentedJSON(http.StatusOK, b)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
}
