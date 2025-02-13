# 🛠 go test dev - Golang Project


**go test dev** adalah proyek backend berbasis Golang menggunakan Gin Framework dan GORM untuk manajemen database. Proyek ini mendukung CRUD Marketing & Penjualan, Manajemen Pembayaran & Kredit, serta Komisi Marketing Otomatis.
Untuk frontend, proyek ini menggunakan Next.js.

### 📌 Fitur Utama
```
✅ CRUD Marketing & Penjualan
✅ Manajemen Pembayaran & Kredit
✅ Komisi Marketing otomatis berdasarkan transaksi
✅ RESTful API menggunakan Gin Framework
✅ Database Management dengan GORM (ORM untuk Golang)
✅ CORS Handling untuk integrasi dengan frontend
✅ Docker Support (Opsional)
```

***🚀 Cara Menjalankan Proyek***

### 1️⃣ Install Golang
Pastikan Golang telah terinstal di sistem Anda.
Jika belum, silakan install dari situs resmi Golang.

Install Golang dari [situs resmi](https://go.dev/doc/install)

### 2️⃣ Clone Repository
```bash
git clone https://github.com/EKKN/gotestdev.git
cd gotestdev
```


### 3️⃣ Jalankan Go mod tidy
```bash
go mod tidy
```


### 4️⃣ Buat Database
Buat database MySQL dengan nama testdev.
Gunakan perintah SQL berikut:
```sql
CREATE DATABASE testdev;
```


### 5️⃣ Konfigurasi .env
Buat file .env di root proyek dan tambahkan konfigurasi berikut:
```env
APP_PORT=:5010                # Port untuk REST API
TRUSTED_PROXY=localhost       # Jalankan di localhost
ALLOW_ORIGIN=http://localhost:3000  # Allow CORS untuk frontend

DB_USER=root                  # Username MySQL
DB_PASSWORD=root              # Password MySQL
DB_HOST=localhost             # Host database
DB_PORT=3306                  # Port MySQL
DB_NAME=testdev               # Nama database
```


### 6️⃣ Jalankan Backend
Jalankan perintah berikut untuk memulai backend:
```sh
go run .
```
✅ Backend akan berjalan di http://localhost:5010.