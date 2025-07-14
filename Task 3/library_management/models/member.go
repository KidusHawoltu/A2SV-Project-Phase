package models

import "fmt"

type Member struct {
	Id            int
	Name          string
	BorrowedBooks []Book
}

func (m Member) String() string {
	return fmt.Sprintf("Id: %v, Name: %s, Borrowed Books: [%v]", m.Id, m.Name, m.BorrowedBooks)
}
