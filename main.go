// https://github.com/gin-gonic/gin

// Build Library API, store a bunch of books
// Check in book, Check out a book, add books, view books
// get book by ID
// routes for API

package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Quantity int `json:"quantity"`
}

var books = []book {
	{ID: "1", Title: "Book#1", Author: "Author#1", Quantity: 3},
	{ID: "2", Title: "Book#2", Author: "Author#2", Quantity: 10},
	{ID: "3", Title: "Book#3", Author: "Author#3", Quantity: 7},
}

// takes in a gin.contect, which is all the information about the request / allows to return a response.
func getBooks(c *gin.Context) {
	// this allows us to receive the data in JSON, sending response of StatusOk and we're sending data: books
	c.IndentedJSON(http.StatusOK, books)
}

func bookByID(c *gin.Context){
	id := c.Param("id")
	book, err := getBookByID(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id quaery parameter"})
		return
	}

	book, err := getBookByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available"})
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id quaery parameter"})
		return
	}

	book, err := getBookByID(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func getBookByID(id string) (*book, error) {
	for i, b := range books{
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

func createBook(c *gin.Context) {
	// create var newBook with type of book
	var newBook book 
	// bind JSON to the newBook / if error doing that then return
	if err := c.BindJSON(&newBook); err != nil {
		return 
	}
	// append books to newBook
	books = append(books, newBook)
	// return book we created with status of created 
	c.IndentedJSON(http.StatusCreated, newBook)
}

// set up router, handle different routes & end points
func main ()  {
	// create router / from gin / allows us to handle routes
	router := gin.Default()
	// route is /books which means when at localhost:8080, /books is going to call the 'getBooks' function
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookByID)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}