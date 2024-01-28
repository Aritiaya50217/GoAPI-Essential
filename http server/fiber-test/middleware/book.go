package main

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Book struct to hold book data
type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// Initialize in - memory data
var books []Book = []Book{
	{Id: 1, Title: "1984", Author: "George Orwell"},
	{Id: 2, Title: "The Great Gatsby", Author: "F. Scott Fitzgerald"},
}

// Handler function
func GetBooks(c *fiber.Ctx) error {
	// Retrive user data from the context
	user := c.Locals(userContextKey).(*UserData)

	//  use the user data
	fmt.Printf("user email : %v , Role : %s\n", user.Email, user.Role)
	return c.JSON(books)
}

func GetBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	for _, book := range books {
		if book.Id == id {
			return c.JSON(book)
		}
	}
	return c.SendStatus(fiber.StatusNotFound)
}

func CreateBook(c *fiber.Ctx) error {
	book := new(Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	book.Id = len(books) + 1
	books = append(books, *book)
	return c.JSON(book)
}

func UpdateBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	bookUpdate := new(Book)
	if err := c.BodyParser(bookUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for i, book := range books {
		if book.Id == id {
			book.Title = bookUpdate.Title
			book.Author = bookUpdate.Author
			books[i] = book
			return c.JSON(book)
		}
	}
	return c.SendStatus(fiber.StatusNotFound)
}

func DeleteBook(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	for i, book := range books {
		if book.Id == id {
			books = append(books[:i], books[i+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}
	return c.SendStatus(fiber.StatusNotFound)
}
