package controllers

import (
	"fmt"
	"library_management/models"
	"library_management/services"
	"strconv"
)

var library *services.Library

func getInput(prompt string) string {
	fmt.Print(prompt)
	var str string
	_, err := fmt.Scanln(&str)
	for err != nil {
		fmt.Print("Invalid Input. Enter valid input: ")
		_, err = fmt.Scanln(&str)
	}
	return str
}

func getIntInput(prompt string) int {
	fmt.Print(prompt)
	var str string
	var val int
	_, err := fmt.Scanln(&str)
	if err == nil {
		val, err = strconv.Atoi(str)
	}
	for err != nil {
		fmt.Print("Invalid Input. Enter valid input: ")
		_, err = fmt.Scanln(&str)
		if err == nil {
			val, err = strconv.Atoi(str)
		}
	}
	return val
}

func printMenu() {
	fmt.Println("Menu: ")
	fmt.Println("Enter 0 to See the Menu")
	fmt.Println("Enter 1 to Add a Member to the Library")
	fmt.Println("Enter 2 to Remove a Member from the Library")
	fmt.Println("Enter 3 to Add a Book to the Library")
	fmt.Println("Enter 4 to Remove a Book from the Library")
	fmt.Println("Enter 5 to Borrow a Book from the Library")
	fmt.Println("Enter 6 to Return a Book to the Library")
	fmt.Println("Enter 7 to See All Members of the Library")
	fmt.Println("Enter 8 to See Available Books in the Library")
	fmt.Println("Enter 9 to See Borrowed Books from the Library")
	fmt.Println("Enter any other integer to exit the program")
}

func Operate() {
	library = services.GetLibrary()
	stop := false

	fmt.Println("Welcome to Console based Library Management")
	printMenu()
	for !stop {
		menu := getIntInput("Your Input: ")
		switch menu {
		case 0:
			printMenu()
		case 1:
			AddMember()
		case 2:
			RemoveMember()
		case 3:
			AddBook()
		case 4:
			RemoveBook()
		case 5:
			BorrowBook()
		case 6:
			ReturnBook()
		case 7:
			ListAllMembers()
		case 8:
			ListAvailableBooks()
		case 9:
			ListBorrowedBooks()
		default:
			stop = true
		}
	}

	fmt.Println("Good Bye")
}

func AddMember() {
	name := getInput("Enter the name of the new Member: ")
	library.AddMember(models.Member{
		Name: name,
	})
}

func RemoveMember() {
	memberId := getIntInput("Enter the id of the member you want to remove from the Library: ")
	library.RemoveMember(memberId)
}

func AddBook() {
	title, author := getInput("Enter the title of the Book: "), getInput("Enter the name of the Author: ")
	library.AddBook(models.Book{
		Title:  title,
		Author: author,
	})
}

func RemoveBook() {
	bookId := getIntInput("Enter the Id of the book you want to delete\n(Enter -1 to see list of available books): ")
	for bookId < 0 {
		ListAvailableBooks()
		bookId = getIntInput("Enter the Id of the book you want to delete\n(Enter -1 to see list of available books): ")
	}
	library.RemoveBook(bookId)
}

func BorrowBook() {
	memberId := getIntInput("Enter your Id: ")
	bookId := getIntInput("Enter the Id of the book you want to borrow\n(Enter -1 to see list of available books): ")
	for bookId < 0 {
		ListAvailableBooks()
		bookId = getIntInput("Enter the Id of the book you want to borrow\n(Enter -1 to see list of available books): ")
	}
	err := library.BorrowBook(bookId, memberId)
	if err == nil {
		fmt.Println("Successfully Borrowed the Book")
	} else {
		fmt.Println("Error borrowing book")
	}
}

func ReturnBook() {
	memberId := getIntInput("Enter your Id: ")
	bookId := getIntInput("Enter the Id of the book you want to return\n(Enter -1 to see list of your Borrowed Books): ")
	for bookId < 0 {
		listBorrowedBooks(memberId)
		bookId = getIntInput("Enter the Id of the book you want to return\n(Enter -1 to see list of your Borrowed Books): ")
	}
	err := library.ReturnBook(bookId, memberId)
	if err == nil {
		fmt.Println("Successfully Returned the Book")
	} else {
		fmt.Println("Error returning the book")
	}
}

func ListAllMembers() {
	members := library.ListAllMembers()
	if len(members) == 0 {
		fmt.Println("There are no members in this library")
		return
	}
	fmt.Printf("%-5v%-15s Borrowed Books\n", "Id", "Name")
	for _, m := range members {
		fmt.Printf("%-5v%-15s ", m.Id, m.Name)
		if len(m.BorrowedBooks) == 0 {
			fmt.Println("-")
		}
		for i, b := range m.BorrowedBooks {
			fmt.Print(b.Title)
			if i == len(m.BorrowedBooks)-1 {
				fmt.Print("\n")
			} else {
				fmt.Print(", ")
			}
		}
	}
}

func ListAvailableBooks() {
	books := library.ListAvailableBooks()
	if len(books) == 0 {
		fmt.Println("There are no available books")
		return
	}
	fmt.Printf("%-5v%-30s%-20s\n", "Id", "Title", "Author")
	for _, book := range books {
		fmt.Printf("%-5v%-30s%-20s\n", book.Id, book.Title, book.Author)
	}
}

func ListBorrowedBooks() {
	memberId := getIntInput("Enter your Id: ")
	listBorrowedBooks(memberId)
}

func listBorrowedBooks(memberId int) {
	books := library.ListBorrowedBooks(memberId)
	if len(books) == 0 {
		fmt.Println("This member hasn't borrowed any books")
		return
	}
	fmt.Printf("%-5v%-30s%-20s\n", "Id", "Title", "Author")
	for _, book := range books {
		fmt.Printf("%-5v%-30s%-20s\n", book.Id, book.Title, book.Author)
	}
}
