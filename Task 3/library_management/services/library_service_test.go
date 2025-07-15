package services

import (
	"fmt"
	"library_management/models"
	"testing"
)

// TestAddMember verifies that members are added correctly with auto-incrementing IDs.
func TestAddMember(t *testing.T) {
	library := NewLibrary()

	t.Run("AddFirstMember", func(t *testing.T) {
		library.AddMember(models.Member{Name: "Alice"})

		if len(library.Members) != 1 {
			t.Fatalf("Expected 1 member, got %d", len(library.Members))
		}
		member, exists := library.Members[0]
		if !exists || member.Name != "Alice" || member.Id != 0 {
			t.Errorf("Member was not added correctly. Got: %+v", member)
		}
	})

	t.Run("AddSecondMember", func(t *testing.T) {
		library.AddMember(models.Member{Name: "Bob"})

		if len(library.Members) != 2 {
			t.Fatalf("Expected 2 members, got %d", len(library.Members))
		}
		member, exists := library.Members[1]
		if !exists || member.Name != "Bob" || member.Id != 1 {
			t.Errorf("Second member was not added correctly. Got: %+v", member)
		}
	})
}

// TestRemoveMember checks member removal, including the important side-effect
// of making their borrowed books available again.
func TestRemoveMember(t *testing.T) {
	t.Run("RemoveMemberWithBorrowedBooks", func(t *testing.T) {
		// Arrange
		library := NewLibrary()
		library.AddBook(models.Book{Title: "Test Book"})      // ID 0
		library.AddMember(models.Member{Name: "Test Member"}) // ID 0
		_ = library.BorrowBook(0, 0)

		// Pre-condition check
		if library.Books[0].Status != models.Borrowed {
			t.Fatal("Setup failed: Book should be borrowed before member removal")
		}

		// Act
		library.RemoveMember(0)

		// Assert
		if _, exists := library.Members[0]; exists {
			t.Error("Expected member to be removed, but they still exist")
		}
		if library.Books[0].Status != models.Available {
			t.Errorf("Expected removed member's book to become Available, but status is %s", library.Books[0].Status)
		}
	})

	t.Run("RemoveNonExistentMember", func(t *testing.T) {
		library := NewLibrary()
		library.AddMember(models.Member{Name: "Alice"})
		// Act & Assert: This should not panic or error
		library.RemoveMember(99)
		if len(library.Members) != 1 {
			t.Error("Removing a non-existent member should not affect the member list")
		}
	})
}

// TestAddBook is a straightforward test for adding books.
func TestAddBook(t *testing.T) {
	library := NewLibrary()
	book := models.Book{Title: "The Go Programming Language", Author: "Donovan & Kernighan"}
	library.AddBook(book)

	if len(library.Books) != 1 {
		t.Fatalf("Expected library to have 1 book, but it has %d", len(library.Books))
	}
	addedBook := library.Books[0]
	if addedBook.Title != book.Title || addedBook.Id != 0 || addedBook.Status != models.Available {
		t.Errorf("Book was not added correctly. Got: %+v", addedBook)
	}
}

// TestRemoveBook verifies that books can be removed, but not if they are borrowed.
func TestRemoveBook(t *testing.T) {
	t.Run("RemoveAvailableBook", func(t *testing.T) {
		library := NewLibrary()
		library.AddBook(models.Book{Title: "To Be Removed"})
		library.RemoveBook(0)
		if _, exists := library.Books[0]; exists {
			t.Error("Expected book to be removed, but it still exists")
		}
	})

	t.Run("CannotRemoveBorrowedBook", func(t *testing.T) {
		library := NewLibrary()
		library.AddBook(models.Book{Title: "Borrowed Book"})
		library.AddMember(models.Member{Name: "Test Member"})
		_ = library.BorrowBook(0, 0)

		library.RemoveBook(0)

		if _, exists := library.Books[0]; !exists {
			t.Error("A borrowed book was removed, but it should not have been")
		}
	})
}

// TestBorrowBook covers the success and failure cases of borrowing a book.
func TestBorrowBook(t *testing.T) {
	// Arrange
	library := NewLibrary()
	library.AddBook(models.Book{Title: "Test Book"}) // ID 0
	library.AddMember(models.Member{Name: "Alice"})  // ID 0
	library.AddMember(models.Member{Name: "Bob"})    // ID 1

	t.Run("SuccessfulBorrow", func(t *testing.T) {
		err := library.BorrowBook(0, 0)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v", err)
		}
		if library.Books[0].Status != models.Borrowed {
			t.Error("Book status did not change to Borrowed")
		}
		if len(library.Members[0].BorrowedBooks) != 1 {
			t.Error("Book was not added to member's borrowed list")
		}
	})

	t.Run("FailToBorrowAlreadyBorrowedBook", func(t *testing.T) {
		// Note: This runs after the successful test, so book 0 is already borrowed by Alice
		err := library.BorrowBook(0, 1) // Bob tries to borrow the same book
		if err == nil {
			t.Error("Expected an error when borrowing an already borrowed book, but got none")
		}
	})

	t.Run("FailWithNonExistentBook", func(t *testing.T) {
		err := library.BorrowBook(99, 0)
		if err == nil {
			t.Error("Expected an error for a non-existent book, but got none")
		}
	})

	t.Run("FailWithNonExistentMember", func(t *testing.T) {
		err := library.BorrowBook(0, 99)
		if err == nil {
			t.Error("Expected an error for a non-existent member, but got none")
		}
	})
}

// TestReturnBook covers the success and failure cases of returning a book.
func TestReturnBook(t *testing.T) {
	// Arrange
	library := NewLibrary()
	library.AddBook(models.Book{Title: "Test Book"}) // ID 0
	library.AddMember(models.Member{Name: "Alice"})  // ID 0
	_ = library.BorrowBook(0, 0)                     // Alice borrows the book

	t.Run("SuccessfulReturn", func(t *testing.T) {
		err := library.ReturnBook(0, 0)
		if err != nil {
			t.Fatalf("Expected no error, but got: %v", err)
		}
		if library.Books[0].Status != models.Available {
			t.Error("Book status did not change to Available")
		}
		if len(library.Members[0].BorrowedBooks) != 0 {
			t.Error("Book was not removed from member's borrowed list")
		}
	})

	t.Run("FailToReturnUnborrowedBook", func(t *testing.T) {
		// Note: This runs after the successful return, so book 0 is now Available
		err := library.ReturnBook(0, 0)
		if err == nil {
			t.Error("Expected error when returning an already available book, but got none")
		}
	})

	t.Run("FailWhenMemberDidNotBorrowBook", func(t *testing.T) {
		// Arrange new scenario
		lib2 := NewLibrary()
		lib2.AddBook(models.Book{Title: "Book A"})
		lib2.AddMember(models.Member{Name: "Charlie"})
		lib2.AddMember(models.Member{Name: "Diana"})
		_ = lib2.BorrowBook(0, 0) // Charlie borrows Book A

		// Act: Diana tries to return Charlie's book
		err := lib2.ReturnBook(0, 1)
		expectedErr := fmt.Sprintf("book with id %v isn't Borrowed by Member with id %v", 0, 1)
		if err == nil || err.Error() != expectedErr {
			t.Errorf("Expected error '%s', but got: %v", expectedErr, err)
		}
	})
}

// TestListAllMembers checks if the listing function works correctly.
func TestListAllMembers(t *testing.T) {
	library := NewLibrary()
	if len(library.ListAllMembers()) != 0 {
		t.Error("Expected empty slice for a new library")
	}

	library.AddMember(models.Member{Name: "Alice"})
	library.AddMember(models.Member{Name: "Bob"})
	members := library.ListAllMembers()
	if len(members) != 2 {
		t.Errorf("Expected 2 members, got %d", len(members))
	}
}

// TestListAvailableBooks checks if it correctly filters out borrowed books.
func TestListAvailableBooks(t *testing.T) {
	library := NewLibrary()
	library.AddBook(models.Book{Title: "Available Book"})
	library.AddBook(models.Book{Title: "Borrowed Book"})
	library.AddMember(models.Member{Name: "Alice"})
	_ = library.BorrowBook(1, 0) // Borrow the second book

	available := library.ListAvailableBooks()
	if len(available) != 1 {
		t.Fatalf("Expected 1 available book, got %d", len(available))
	}
	if available[0].Id != 0 || available[0].Title != "Available Book" {
		t.Errorf("Incorrect book listed as available. Got: %+v", available[0])
	}
}

// TestListBorrowedBooks checks if it correctly lists books for a specific member.
func TestListBorrowedBooks(t *testing.T) {
	library := NewLibrary()
	library.AddBook(models.Book{Title: "Book A"})
	library.AddBook(models.Book{Title: "Book B"})
	library.AddMember(models.Member{Name: "Alice"})

	t.Run("NoBorrowedBooks", func(t *testing.T) {
		books := library.ListBorrowedBooks(0)
		if len(books) != 0 {
			t.Error("Expected empty slice for member with no borrowed books")
		}
	})

	_ = library.BorrowBook(1, 0) // Alice borrows Book B

	t.Run("OneBorrowedBook", func(t *testing.T) {
		books := library.ListBorrowedBooks(0)
		if len(books) != 1 {
			t.Fatalf("Expected 1 borrowed book, got %d", len(books))
		}
		if books[0].Id != 1 || books[0].Title != "Book B" {
			t.Errorf("Incorrect book listed as borrowed. Got: %+v", books[0])
		}
	})

	t.Run("NonExistentMember", func(t *testing.T) {
		books := library.ListBorrowedBooks(99)
		if len(books) != 0 {
			t.Error("Expected empty slice for a non-existent member")
		}
	})
}
