// client/main.go

package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "grpc-university-project/university/protos"
)

const (
	serverAddr = "localhost:50051"
)

func main() {
	// Sunucuya güvenli olmayan (plaintext) bir bağlantı kuruyoruz.
	// Gerçek dünyada TLS sertifikaları kullanılmalıdır.
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Her bir servis için yeni bir istemci oluşturuyoruz.
	bookClient := pb.NewBookServiceClient(conn)
	studentClient := pb.NewStudentServiceClient(conn)
	loanClient := pb.NewLoanServiceClient(conn)

	// İstekler için bir context oluşturuyoruz. Timeout belirlemek iyi bir pratiktir.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// --- Örnek RPC Çağrıları ---

	log.Println("--- 1. Listing all books ---")
	listBooks(ctx, bookClient)

	log.Println("\n--- 2. Adding a new student ---")
	newStudent := addStudent(ctx, studentClient, "Zeynep Kaya", "987654", "zeynep.kaya@example.com")
	if newStudent == nil {
		log.Println("Could not add student, skipping further tests for this student.")
		return
	}
	log.Printf("Added student: %s (ID: %s)", newStudent.Name, newStudent.Id)

	log.Println("\n--- 3. Listing all students ---")
	listStudents(ctx, studentClient)

	log.Println("\n--- 4. Borrowing a book for the new student ---")
	// İlk önce tüm kitapları alıp bir tanesini seçelim.
	booksResponse, err := bookClient.ListBooks(ctx, &pb.ListBooksRequest{})
	if err != nil || len(booksResponse.GetBooks()) == 0 {
		log.Fatalf("Could not list books to borrow one: %v", err)
	}
	// Stokta olan bir kitabı ödünç alalım.
	var bookToBorrow *pb.Book
	for _, book := range booksResponse.GetBooks() {
		if book.Stock > 0 {
			bookToBorrow = book
			break
		}
	}

	if bookToBorrow == nil {
		log.Println("No books in stock to borrow.")
	} else {
		borrowBook(ctx, loanClient, newStudent.Id, bookToBorrow.Id)
	}

	log.Println("\n--- 5. Listing all loans ---")
	listLoans(ctx, loanClient)
}

// --- Yardımcı Fonksiyonlar ---

func listBooks(ctx context.Context, client pb.BookServiceClient) {
	r, err := client.ListBooks(ctx, &pb.ListBooksRequest{})
	if err != nil {
		log.Printf("could not list books: %v", err)
		return
	}
	log.Println("Books:")
	for _, book := range r.GetBooks() {
		log.Printf("- ID: %s, Title: %s, Stock: %d", book.Id, book.Title, book.Stock)
	}
}

func addStudent(ctx context.Context, client pb.StudentServiceClient, name, number, email string) *pb.Student {
	student := &pb.Student{
		Name:          name,
		StudentNumber: number,
		Email:         email,
		IsActive:      true,
	}
	r, err := client.AddStudent(ctx, &pb.AddStudentRequest{Student: student})
	if err != nil {
		log.Printf("could not add student: %v", err)
		return nil
	}
	return r
}

func listStudents(ctx context.Context, client pb.StudentServiceClient) {
	r, err := client.ListStudents(ctx, &pb.ListStudentsRequest{})
	if err != nil {
		log.Printf("could not list students: %v", err)
		return
	}
	log.Println("Students:")
	for _, student := range r.GetStudents() {
		log.Printf("- ID: %s, Name: %s, Email: %s", student.Id, student.Name, student.Email)
	}
}

func borrowBook(ctx context.Context, client pb.LoanServiceClient, studentID, bookID string) {
	log.Printf("Attempting to borrow book (ID: %s) for student (ID: %s)", bookID, studentID)
	r, err := client.BorrowBook(ctx, &pb.BorrowBookRequest{StudentId: studentID, BookId: bookID})
	if err != nil {
		log.Printf("could not borrow book: %v", err)
		return
	}
	log.Printf("Successfully borrowed. Loan ID: %s, Status: %s", r.Id, r.Status)
}

func listLoans(ctx context.Context, client pb.LoanServiceClient) {
	r, err := client.ListLoans(ctx, &pb.ListLoansRequest{})
	if err != nil {
		log.Printf("could not list loans: %v", err)
		return
	}
	log.Println("Loans:")
	for _, loan := range r.GetLoans() {
		log.Printf("- ID: %s, BookID: %s, StudentID: %s, Status: %s", loan.Id, loan.BookId, loan.StudentId, loan.Status)
	}
}