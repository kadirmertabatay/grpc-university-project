// server/main.go (POINTER HATALARI İÇİN DÜZELTİLMİŞ SÜRÜM)

package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "grpc-university-project/university/protos"
)

// --- Sunucu Yapısı ---
type universityServer struct {
	pb.UnimplementedBookServiceServer
	pb.UnimplementedStudentServiceServer
	pb.UnimplementedLoanServiceServer

	mu       sync.Mutex
	books    map[string]*pb.Book
	students map[string]*pb.Student
	loans    map[string]*pb.Loan
}

func newServer() *universityServer {
	s := &universityServer{
		books:    make(map[string]*pb.Book),
		students: make(map[string]*pb.Student),
		loans:    make(map[string]*pb.Loan),
	}
	addInitialData(s)
	return s
}

// ... (BookService ve StudentService metotları öncekiyle aynı, değişiklik yok) ...
// =================================================================
// --- BookService Metotları ---
// =================================================================

func (s *universityServer) AddBook(ctx context.Context, req *pb.AddBookRequest) (*pb.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	book := req.GetBook()
	if book.GetId() == "" {
		book.Id = uuid.New().String()
	}
	s.books[book.Id] = book
	log.Printf("Added book: %s", book.Title)
	return book, nil
}

func (s *universityServer) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	book, exists := s.books[req.GetId()]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "book with ID %s not found", req.GetId())
	}
	return book, nil
}

func (s *universityServer) ListBooks(ctx context.Context, req *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var bookList []*pb.Book
	for _, book := range s.books {
		bookList = append(bookList, book)
	}
	return &pb.ListBooksResponse{Books: bookList}, nil
}

func (s *universityServer) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.Book, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	bookToUpdate := req.GetBook()
	_, exists := s.books[bookToUpdate.GetId()]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "cannot update, book with ID %s not found", bookToUpdate.GetId())
	}
	s.books[bookToUpdate.GetId()] = bookToUpdate
	log.Printf("Updated book: %s", bookToUpdate.Title)
	return bookToUpdate, nil
}

func (s *universityServer) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*pb.DeleteBookResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	idToDelete := req.GetId()
	if _, exists := s.books[idToDelete]; !exists {
		return nil, status.Errorf(codes.NotFound, "cannot delete, book with ID %s not found", idToDelete)
	}
	delete(s.books, idToDelete)
	log.Printf("Deleted book with ID: %s", idToDelete)
	return &pb.DeleteBookResponse{Message: "Book deleted successfully"}, nil
}

// =================================================================
// --- StudentService Metotları ---
// =================================================================

func (s *universityServer) AddStudent(ctx context.Context, req *pb.AddStudentRequest) (*pb.Student, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	student := req.GetStudent()
	if student.GetId() == "" {
		student.Id = uuid.New().String()
	}
	s.students[student.Id] = student
	log.Printf("Added student: %s", student.Name)
	return student, nil
}

func (s *universityServer) GetStudent(ctx context.Context, req *pb.GetStudentRequest) (*pb.Student, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	student, exists := s.students[req.GetId()]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "student with ID %s not found", req.GetId())
	}
	return student, nil
}

func (s *universityServer) ListStudents(ctx context.Context, req *pb.ListStudentsRequest) (*pb.ListStudentsResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var studentList []*pb.Student
	for _, student := range s.students {
		studentList = append(studentList, student)
	}
	return &pb.ListStudentsResponse{Students: studentList}, nil
}

func (s *universityServer) UpdateStudent(ctx context.Context, req *pb.UpdateStudentRequest) (*pb.Student, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	studentToUpdate := req.GetStudent()
	if _, exists := s.students[studentToUpdate.GetId()]; !exists {
		return nil, status.Errorf(codes.NotFound, "cannot update, student with ID %s not found", studentToUpdate.GetId())
	}
	s.students[studentToUpdate.GetId()] = studentToUpdate
	log.Printf("Updated student: %s", studentToUpdate.Name)
	return studentToUpdate, nil
}

func (s *universityServer) DeleteStudent(ctx context.Context, req *pb.DeleteStudentRequest) (*pb.DeleteStudentResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	idToDelete := req.GetId()
	if _, exists := s.students[idToDelete]; !exists {
		return nil, status.Errorf(codes.NotFound, "cannot delete, student with ID %s not found", idToDelete)
	}
	delete(s.students, idToDelete)
	log.Printf("Deleted student with ID: %s", idToDelete)
	return &pb.DeleteStudentResponse{Message: "Student deleted successfully"}, nil
}

// =================================================================
// --- LoanService Metotları ---
// =================================================================

func (s *universityServer) BorrowBook(ctx context.Context, req *pb.BorrowBookRequest) (*pb.Loan, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	student, studentExists := s.students[req.GetStudentId()]
	if !studentExists {
		return nil, status.Errorf(codes.NotFound, "student with ID %s not found", req.GetStudentId())
	}
	book, bookExists := s.books[req.GetBookId()]
	if !bookExists {
		return nil, status.Errorf(codes.NotFound, "book with ID %s not found", req.GetBookId())
	}

	if book.Stock <= 0 {
		return nil, status.Errorf(codes.FailedPrecondition, "book '%s' is out of stock", book.Title)
	}

	book.Stock--
	s.books[book.Id] = book

	loan := &pb.Loan{
		Id:         uuid.New().String(),
		StudentId:  student.Id,
		BookId:     book.Id,
		LoanDate:   time.Now().Format("2006-01-02"),
		ReturnDate: nil, // DÜZELTME: Boş (yok) bir pointer için 'nil' kullanılır.
		Status:     pb.Loan_ONGOING,
	}
	s.loans[loan.Id] = loan

	log.Printf("Student '%s' borrowed book '%s'", student.Name, book.Title)
	return loan, nil
}

func (s *universityServer) ReturnBook(ctx context.Context, req *pb.ReturnBookRequest) (*pb.Loan, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	loan, exists := s.loans[req.GetLoanId()]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "loan with ID %s not found", req.GetLoanId())
	}

	if loan.Status == pb.Loan_RETURNED {
		return nil, status.Errorf(codes.FailedPrecondition, "this loan has already been returned")
	}

	book, bookExists := s.books[loan.BookId]
	if bookExists {
		book.Stock++
		s.books[book.Id] = book
	}

	loan.Status = pb.Loan_RETURNED
	// DÜZELTME: Bir string'in adresini atamak için geçici bir değişkene ihtiyacımız var.
	returnDateStr := time.Now().Format("2006-01-02")
	loan.ReturnDate = &returnDateStr
	s.loans[loan.Id] = loan

	log.Printf("Loan with ID %s has been returned", loan.Id)
	return loan, nil
}

func (s *universityServer) GetLoan(ctx context.Context, req *pb.GetLoanRequest) (*pb.Loan, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	loan, exists := s.loans[req.GetId()]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "loan with ID %s not found", req.GetId())
	}
	return loan, nil
}

func (s *universityServer) ListLoans(ctx context.Context, req *pb.ListLoansRequest) (*pb.ListLoansResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var loanList []*pb.Loan
	for _, loan := range s.loans {
		loanList = append(loanList, loan)
	}
	return &pb.ListLoansResponse{Loans: loanList}, nil
}

// =================================================================
// --- Ana Fonksiyon ve Başlangıç Verisi ---
// =================================================================

func main() {
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := newServer()

	pb.RegisterBookServiceServer(grpcServer, server)
	pb.RegisterStudentServiceServer(grpcServer, server)
	pb.RegisterLoanServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	log.Printf("gRPC server listening on %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func addInitialData(s *universityServer) {
    // Önceki AddBook çağrılarındaki yazım hatası düzeltildi.
	book1Req := &pb.AddBookRequest{Book: &pb.Book{Title: "The Go Programming Language", Author: "Donovan & Kernighan", Isbn: "978-0134190440", Publisher: "AW", PageCount: 380, Stock: 10}}
	book1Resp, _ := s.AddBook(context.Background(), book1Req)

	book2Req := &pb.AddBookRequest{Book: &pb.Book{Title: "Clean Architecture", Author: "Robert C. Martin", Isbn: "978-0134494166", Publisher: "Prentice Hall", PageCount: 432, Stock: 5}}
	s.AddBook(context.Background(), book2Req)

	student1Req := &pb.AddStudentRequest{Student: &pb.Student{Name: "Ali Veli", StudentNumber: "123456", Email: "ali.veli@example.com", IsActive: true}}
	student1Resp, _ := s.AddStudent(context.Background(), student1Req)

	s.AddStudent(context.Background(), &pb.AddStudentRequest{Student: &pb.Student{Name: "Ayşe Yılmaz", StudentNumber: "654321", Email: "ayse.yilmaz@example.com", IsActive: true}})

	if student1Resp != nil && book1Resp != nil {
		s.BorrowBook(context.Background(), &pb.BorrowBookRequest{
			StudentId: student1Resp.Id,
			BookId:    book1Resp.Id,
		})
	}
}