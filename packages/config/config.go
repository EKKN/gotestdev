package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/EKKN/gotestdev/packages/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using default environment variables")
	}

	// Ambil konfigurasi dari .env atau gunakan default jika tidak ada
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Format DSN menggunakan variabel dari .env
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Koneksi ke database
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}

	DB = database
	fmt.Println("Database MySQL terhubung!")

	// AutoMigrate semua model dari package models
	database.AutoMigrate(&models.Marketing{}, &models.Penjualan{}, &models.Pembayaran{}, &models.PembayaranDetail{}, &models.KomisiPersen{}, &models.Kredit{})

	// Jalankan seeding data
	seedData(database)
}

// Fungsi untuk menambahkan data default
func seedData(db *gorm.DB) {
	seedMarketing(db)
	seedKomisiPersen(db)
	seedPenjualan(db)
	seedKredit(db)
}

// ** 1️⃣ Seeder Marketing **
func seedMarketing(db *gorm.DB) {
	var count int64
	db.Model(&models.Marketing{}).Count(&count)

	if count == 0 { // Jika tabel kosong, tambahkan data default
		marketings := []models.Marketing{
			{ID: 1, Name: "Alfandy", IsActive: true},
			{ID: 2, Name: "Mery", IsActive: true},
			{ID: 3, Name: "Danang", IsActive: true},
		}

		db.Create(&marketings)
		fmt.Println("✅ Seed data Marketing berhasil ditambahkan!")
	} else {
		fmt.Println("⚠️ Seed data Marketing sudah ada, tidak perlu menambahkan ulang.")
	}
}

// ** 2️⃣ Seeder Komisi **
func seedKomisiPersen(db *gorm.DB) {
	var count int64
	db.Model(&models.KomisiPersen{}).Count(&count)

	if count == 0 { // Jika tabel kosong, tambahkan data default
		komisiData := []models.KomisiPersen{
			{MinOmzet: 0, MaxOmzet: 100000000, Persentase: 0},
			{MinOmzet: 100000001, MaxOmzet: 200000000, Persentase: 2.5},
			{MinOmzet: 200000001, MaxOmzet: 500000000, Persentase: 5},
			{MinOmzet: 500000001, MaxOmzet: 0, Persentase: 10}, // 0 berarti tidak terbatas
		}

		db.Create(&komisiData)
		fmt.Println("✅ Seed data KomisiPersen berhasil ditambahkan!")
	} else {
		fmt.Println("⚠️ Seed data KomisiPersen sudah ada, tidak perlu menambahkan ulang.")
	}
}

// ** 3️⃣ Seeder Kredit **
func seedKredit(db *gorm.DB) {
	var count int64
	db.Model(&models.Kredit{}).Count(&count)

	if count == 0 {
		kreditData := []models.Kredit{
			{Months: 6, InterestRate: 0.0},
			{Months: 12, InterestRate: 0.0},
			{Months: 24, InterestRate: 4.0},
			{Months: 36, InterestRate: 3.8},
			{Months: 48, InterestRate: 3.5},
			{Months: 60, InterestRate: 3.2},
		}

		db.Create(&kreditData)
		fmt.Println("✅ Seed data Kredit berhasil ditambahkan!")
	} else {
		fmt.Println("⚠️ Seed data Kredit sudah ada, tidak perlu menambahkan ulang.")
	}
}

// ** 4️⃣ Seeder Penjualan **
func seedPenjualan(db *gorm.DB) {
	var count int64
	db.Model(&models.Penjualan{}).Count(&count)

	if count == 0 {
		penjualanData := []models.Penjualan{
			{TransactionNumber: "TRX001", MarketingID: 1, Date: parseDate("2023-05-22"), CargoFee: 25000, TotalBalance: 3000000, GrandTotal: 3025000},
			{TransactionNumber: "TRX002", MarketingID: 3, Date: parseDate("2023-05-22"), CargoFee: 25000, TotalBalance: 320000, GrandTotal: 345000},
			{TransactionNumber: "TRX003", MarketingID: 1, Date: parseDate("2023-05-22"), CargoFee: 0, TotalBalance: 65000000, GrandTotal: 65000000},
			{TransactionNumber: "TRX004", MarketingID: 1, Date: parseDate("2023-05-23"), CargoFee: 10000, TotalBalance: 70000000, GrandTotal: 70010000},
			{TransactionNumber: "TRX005", MarketingID: 2, Date: parseDate("2023-05-23"), CargoFee: 10000, TotalBalance: 80000000, GrandTotal: 80010000},
			{TransactionNumber: "TRX006", MarketingID: 3, Date: parseDate("2023-05-23"), CargoFee: 12000, TotalBalance: 44000000, GrandTotal: 44012000},
			{TransactionNumber: "TRX007", MarketingID: 1, Date: parseDate("2023-06-01"), CargoFee: 0, TotalBalance: 75000000, GrandTotal: 75000000},
			{TransactionNumber: "TRX008", MarketingID: 2, Date: parseDate("2023-06-02"), CargoFee: 0, TotalBalance: 85000000, GrandTotal: 85000000},
			{TransactionNumber: "TRX009", MarketingID: 2, Date: parseDate("2023-06-01"), CargoFee: 0, TotalBalance: 175000000, GrandTotal: 175000000},
			{TransactionNumber: "TRX010", MarketingID: 3, Date: parseDate("2023-06-01"), CargoFee: 0, TotalBalance: 75000000, GrandTotal: 75000000},
			{TransactionNumber: "TRX011", MarketingID: 2, Date: parseDate("2023-06-01"), CargoFee: 0, TotalBalance: 750020000, GrandTotal: 750020000},
			{TransactionNumber: "TRX012", MarketingID: 3, Date: parseDate("2023-06-01"), CargoFee: 0, TotalBalance: 130000000, GrandTotal: 120000000},
		}

		db.Create(&penjualanData)
		fmt.Println("✅ Seed data Penjualan berhasil ditambahkan!")
	} else {
		fmt.Println("⚠️ Seed data Penjualan sudah ada, tidak perlu menambahkan ulang.")
	}
}

// ** Fungsi bantu: Parse tanggal **
func parseDate(dateStr string) time.Time {
	layout := "2006-01-02"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		fmt.Printf("❌ Error parsing date: %s\n", dateStr)
		return time.Now()
	}
	return t
}
