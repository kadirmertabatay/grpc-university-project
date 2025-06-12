# gRPC Üniversite Kütüphane Sistemi

Bu proje, Açık Kaynak Kodlu Yazılımlar dersi ödevi kapsamında gRPC kullanılarak geliştirilmiş basit bir üniversite kütüphane yönetim sistemidir.

Proje, Go (Golang) programlama dili kullanılarak geliştirilmiştir ve API tanımlaması için Protocol Buffers (Protobuf) kullanılmıştır.

## 🎯 Özellikler

Sistem, üç ana servis üzerine kuruludur:

* **`BookService`**: Kitap ekleme, listeleme, güncelleme, silme ve tekil görüntüleme.
* **`StudentService`**: Öğrenci ekleme, listeleme, güncelleme, silme ve tekil görüntüleme.
* **`LoanService`**: Kitap ödünç alma, iade etme, listeleme ve tekil görüntüleme.

## 🛠️ Kullanılan Teknolojiler

* **Programlama Dili:** Go (Golang)
* **RPC Framework:** gRPC
* **API Tanımlama:** Protocol Buffers (proto3)

## 🚀 Kurulum ve Çalıştırma

### Gereksinimler

* Go (v1.18 veya üstü)
* Protobuf Derleyicisi (`protoc`)
* `grpcurl` (Test için)

### Adımlar

1. **Projeyi bilgisayarınıza klonlayın:**

   ```sh
   git clone [GitHub Repo Linkiniz]
   cd [Proje Klasör Adınız]
   ```
2. **Gerekli Go araçlarını ve bağımlılıkları yükleyin:**

   ```sh
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   go mod tidy
   ```
3. **`.proto` dosyasından Go arayüz kodlarını (stub) üretin:**

   ```sh
   protoc --go_out=. --go-grpc_out=. university.proto
   ```
4. **Sunucuyu çalıştırın (bir terminal penceresinde):**

   ```sh
   go run server/main.go
   ```

   Sunucu `localhost:50051` portunda çalışmaya başlayacaktır. Bu terminali açık bırakın.
5. **İstemciyi çalıştırın (yeni bir terminal penceresinde):**

   ```sh
   go run client/main.go
   ```

   İstemci, sunucuya bağlanarak bir dizi örnek işlemi gerçekleştirecek ve sonuçları ekrana yazdıracaktır.

## 🧪 `grpcurl` ile Test Etme

Sunucu çalışır durumdayken, servisler `grpcurl` aracıyla doğrudan test edilebilir. Detaylı komutlar ve çıktılar için `grpcurl-tests.md` dosyasına bakınız.

**Örnek: Tüm kitapları listeleme**

```sh
grpcurl -plaintext -d '{}' localhost:50051 university.BookService.ListBooks
```


Bu proje, Açık Kaynak Kodlu Yazılımlar dersi ödevi kapsamında gRPC kullanılarak geliştirilmiş basit bir üniversite kütüphane yönetim sistemidir.

Proje, Go (Golang) programlama dili kullanılarak geliştirilmiştir ve Protocol Buffers (Protobuf) ile API tanımlaması yapılmıştır.

## Servisler ve Metotlar

Sistem, üç ana servis üzerine kuruludur:

* **BookService:** Kitap ekleme, listeleme, güncelleme, silme ve tekil görüntüleme.
* **StudentService:** Öğrenci ekleme, listeleme, güncelleme, silme ve tekil görüntüleme.
* **LoanService:** Kitap ödünç alma, iade etme, listeleme ve tekil görüntüleme.

## Teknoloji Stack'i

* **Programlama Dili:** Go
* **RPC Framework:** gRPC
* **API Tanımlama:** Protocol Buffers (Protobuf v3)

## Kurulum ve Çalıştırma

### Gereksinimler

* Go (v1.18+)
* Protobuf Derleyicisi (`protoc`)
* `grpcurl` (Test için)

### Adımlar

1. **Projeyi klonlayın:**

   ```sh
   git clone https://[GitHub_Repo_Linkiniz]
   cd grpc-university-project
   ```
2. **Gerekli Go araçlarını ve bağımlılıkları kurun:**

   ```sh
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   go mod tidy
   ```
3. **`.proto` dosyasından Go stub'larını üretin:**

   ```sh
   protoc --go_out=. --go-grpc_out=. university.proto
   ```
4. **Sunucuyu çalıştırın (bir terminalde):**

   ```sh
   go run server/main.go
   ```

   Sunucu `localhost:50051` portunda çalışmaya başlayacaktır.
5. **İstemciyi çalıştırın (yeni bir terminalde):**

   ```sh
   go run client/main.go
   ```

## `grpcurl` ile Test Etme

Sunucu çalışır durumdayken, servisleri `grpcurl` ile test edebilirsiniz. Detaylı komutlar ve çıktılar için `grpcurl-tests.md` dosyasına bakınız.

Örnek: Tüm kitapları listeleme

```sh
grpcurl -plaintext -d '{}' localhost:50051 university.BookService.ListBooks
```
