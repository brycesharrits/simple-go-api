package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"errors"
)

type book struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Author string `json:"author"`
	Quantity int `json:"quantity"`
	StartingQuantity int `json:"startingQuantity"`
}



var books = []book{
	{
		ID: "1",
		Name: "Book 1",
		Author: "Bryce",
		Quantity: 10,
		StartingQuantity: 10,
	},
	{
		ID: "2",
		Name: "Book 2",
		Author: "Jake",
		Quantity: 5,
		StartingQuantity: 5,
	},
	{
		ID: "3",
		Name: "Book 3",
		Author: "Stephen",
		Quantity: 2,
		StartingQuantity: 2,
	},

}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	
	if ok == false {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book is out of stock."})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	
	if ok == false {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity == book.StartingQuantity {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book is already at max quantity."})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("Book not found.")
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return-book", returnBook)
	router.Run("localhost:8080")
}