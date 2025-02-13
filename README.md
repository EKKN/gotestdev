# 🛠 go test dev - Golang Project

**go test dev** adalah proyek backend berbasis **Golang** yang menggunakan **Gin Framework** dan **GORM** untuk manajemen database.

## 📌 Fitur
- CRUD Marketing & Penjualan
- Manajemen Pembayaran & Kredit
- Komisi Marketing Otomatis
- Menggunakan **Next.js** sebagai frontend

## 🚀 Cara Menjalankan Proyek

### 1️⃣ Install golang
install golang
https://go.dev/doc/install

### 2️⃣ Clone Repository
git clone https://github.com/EKKN/gotestdev.git

cd gotestdev


### 3️⃣ Jalankan Go mod tidy
go mod tidy




### 4️⃣ CREATE Database
Buat database MySQL dengan nama testdev



### 5️⃣ File .env
Edit file .envexmaple di root dir menjadi .env untuk konfigurasi database:

APP_PORT=:5010 => ini adalah port running rest

TRUSTED_PROXY=localhost => untuk dijalankan di localhost

ALLOW_ORIGIN=http://localhost:3000  => allow origin

DB_USER=root

DB_PASSWORD=password

DB_HOST=localhost

DB_PORT=3306

DB_NAME=testdev




### 6️⃣ jalankan service
go run .