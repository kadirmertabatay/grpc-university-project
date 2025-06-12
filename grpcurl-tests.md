# gRPC Servis Testleri (`grpcurl`)

Bu doküman, gRPC sunucusunun metotlarını test etmek için kullanılan `grpcurl` komutlarını ve sonuçlarını içerir.

## 1. Genel Testler

### Sunucudaki Servisleri Listeleme

Bu komut, sunucunun hangi servisleri desteklediğini kontrol eder.

**Komut:**

```sh
grpcurl -plaintext localhost:50051 list
```

**Örnek Çıktı:**

```
grpc.reflection.v1.ServerReflection
grpc.reflection.v1alpha.ServerReflection
university.BookService
university.LoanService
university.StudentService
```

## 2. BookService Testleri

### ListBooks (Tüm Kitapları Listele)

**Komut:**

```
grpcurl -plaintext -d '{}' localhost:50051 university.BookService.ListBooks
```

**Örnek Çıktı:**

```
{
  "books": [
    {
      "id": "1b5d5349-e392-4843-ae55-246577531a11",
      "title": "The Go Programming Language",
      "author": "Donovan & Kernighan",
      "isbn": "978-0134190440",
      "publisher": "AW",
      "pageCount": 380,
      "stock": 10
    },
    {
      "id": "40467862-2d1c-4843-a97f-acaa4ec295b0",
      "title": "Clean Architecture",
      "author": "Robert C. Martin",
      "isbn": "978-0134494166",
      "publisher": "Prentice Hall",
      "pageCount": 432,
      "stock": 5
    }
  ]
}
```

GetBook (ID ile Tek Kitap Getir)

**Komut:** **(ID'yi kendi çıktınızdaki bir ID ile değiştirin)**

```
grpcurl -plaintext -d '{"id": "1b5d5349-e392-4843-ae55-246577531a11"}' localhost:50051 university.BookService.GetBook
```

## 3. StudentService Testleri

### ListStudents (Tüm Öğrencileri Listele)

**Komut:**

```
grpcurl -plaintext -d '{}' localhost:50051 university.StudentService.ListStudents
```

**Örnek Çıktı:**

```
{
  "students": [
    {
      "id": "422a98fe-826e-40c6-99d1-c7a419d36f59",
      "name": "Ali Veli",
      "studentNumber": "123456",
      "email": "ali.veli@example.com",
      "isActive": true
    },
    {
      "id": "5bbb3f43-6090-4fb3-a671-fbedc9c1a705",
      "name": "Ayşe Yılmaz",
      "studentNumber": "654321",
      "email": "ayse.yilmaz@example.com",
      "isActive": true
    }
  ]
}
```

## 4. LoanService Testleri

### ListLoans (Tüm Ödünç İşlemlerini Listele)

**Komut:**

```
grpcurl -plaintext -d '{}' localhost:50051 university.LoanService.ListLoans
```

**Örnek Çıktı:**

```
{
  "loans": [
    {
      "id": "907ce512-a9b1-4f0f-921e-652280637462",
      "studentId": "422a98fe-826e-40c6-99d1-c7a419d36f59",
      "bookId": "1b5d5349-e392-4843-ae55-246577531a11",
      "loanDate": "2025-06-12",
      "status": "ONGOING"
    }
  ]
}
```

## 5. Hatalı Durum Senaryoları

### Bulunmayan Bir Kitabı Getirme (**NotFound** **Hatası)**

**Komut:**

```
grpcurl -plaintext -d '{"id": "yok-boyle-bir-id"}' localhost:50051 university.BookService.GetBook
```

**Çıktı:**

```
{
  "code": 5,
  "message": "book with ID yok-boyle-bir-id not found",
  "details": []
}
```
