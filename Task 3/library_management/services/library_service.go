package services

import (
	"fmt"
	"library_management/models"
)

func GetLibrary() *Library {
	return &Library{
		Books:        make(map[int]models.Book),
		nextBookId:   0,
		Members:      make(map[int]models.Member),
		nextMemberId: 0,
	}
}

type Library struct {
	Books        map[int]models.Book
	nextBookId   int
	Members      map[int]models.Member
	nextMemberId int
}

type LibraryManager interface {
	AddMember(member models.Member)
	RemoveMember(memberId int)
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAllMembers() []models.Member
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

func (library *Library) AddMember(member models.Member) {
	member.Id = library.nextMemberId
	library.Members[library.nextMemberId] = member
	library.nextMemberId++
}

func (library *Library) RemoveMember(memberId int) {
	member := library.Members[memberId]
	for _, borrowBook := range member.BorrowedBooks {
		book, e := library.Books[borrowBook.Id]
		if e {
			book.Status = models.Available
			library.Books[book.Id] = book
		}
	}
	delete(library.Members, memberId)
}

func (library *Library) AddBook(book models.Book) {
	book.Id = library.nextBookId
	book.Status = models.Available
	library.Books[library.nextBookId] = book
	library.nextBookId++
}

func (library *Library) RemoveBook(bookId int) {
	book, exists := library.Books[bookId]
	if exists && book.Status == models.Borrowed {
		fmt.Printf("Book with id %v is already borrowed\n", bookId)
	}
	delete(library.Books, bookId)
}

func (library *Library) BorrowBook(bookId int, memberId int) error {
	member, memberExists := library.Members[memberId]
	if !memberExists {
		return fmt.Errorf("member with id %v doesn't exist", memberId)
	}
	book, bookExists := library.Books[bookId]
	if !bookExists {
		return fmt.Errorf("book with id %v doesn't exist", bookId)
	}
	if book.Status == models.Borrowed {
		return fmt.Errorf("book with id %v is already Borrowed", bookId)
	}
	book.Status = models.Borrowed
	library.Books[book.Id] = book
	member.BorrowedBooks = append(member.BorrowedBooks, book)
	library.Members[memberId] = member
	return nil
}

func (library *Library) ReturnBook(bookId int, memberId int) error {
	member, memberExists := library.Members[memberId]
	if !memberExists {
		return fmt.Errorf("member with id %v doesn't exist", memberId)
	}
	book, bookExists := library.Books[bookId]
	if !bookExists {
		return fmt.Errorf("book with id %v doesn't exist", bookId)
	}
	if book.Status == models.Available {
		return fmt.Errorf("book with id %v isn't Borrowed", bookId)
	}
	for i, borrwoedBook := range member.BorrowedBooks {
		if borrwoedBook.Id == bookId {
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]...)
			library.Members[memberId] = member
			book.Status = models.Available
			library.Books[bookId] = book
			return nil
		}
	}
	return fmt.Errorf("book with id %v isn't Borrowed by Member with id %v", bookId, memberId)
}

func (library *Library) ListAllMembers() []models.Member {
	var members []models.Member
	for _, m := range library.Members {
		members = append(members, m)
	}
	return members
}

func (library *Library) ListAvailableBooks() []models.Book {
	availableBooks := []models.Book{}
	for _, book := range library.Books {
		if book.Status == models.Available {
			availableBooks = append(availableBooks, book)
		}
	}
	return availableBooks
}

func (library *Library) ListBorrowedBooks(memberId int) []models.Book {
	member, memberExists := library.Members[memberId]
	if !memberExists {
		return []models.Book{}
	}
	return member.BorrowedBooks
}
