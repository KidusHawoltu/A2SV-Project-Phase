package models

import "fmt"

type BookStatus string

const (
	Available BookStatus = "Available"
	Borrowed  BookStatus = "Borrowed"
)

type Book struct {
	Id     int
	Title  string
	Author string
	Status BookStatus
}

func (book *Book) String() string {
	return fmt.Sprintf("Id: %v, %q by %q (%s)", book.Id, book.Title, book.Author, book.Status)
}
