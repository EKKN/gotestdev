package api

import (
	"log"
	"time"

	"github.com/EKKN/gotestdev/packages/config"
	"github.com/EKKN/gotestdev/packages/helpers"
	"github.com/EKKN/gotestdev/packages/models"
	"github.com/gin-gonic/gin"
)

func (s *APIServer) GetAllPenjualan(c *gin.Context) (map[string]interface{}, map[string]interface{}) {
	var penjualan []models.Penjualan

	// Query data penjualan dengan join ke tabel marketing
	if err := config.DB.Preload("Marketing").Order("date DESC").Find(&penjualan).Error; err != nil {
		return nil, gin.H{"message": "Failed to fetch penjualan data", "error": err.Error()}
	}

	return gin.H{"data": penjualan}, nil
}

func (s *APIServer) CreatePenjualan(c *gin.Context) (map[string]interface{}, map[string]interface{}) {
	var input models.PenjualanInput

	// Bind JSON ke struct sementara
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("Error binding JSON:", err)
		return nil, gin.H{"message": "Invalid request data", "error": err.Error()}
	}

	// Validasi input menggunakan validator/v10
	errors := helpers.ValidateStruct(input)
	if len(errors) > 0 {
		return nil, gin.H{"message": "Validation errors", "errors": errors}
	}

	// Validasi data Marketing
	var marketing models.Marketing
	if err := config.DB.First(&marketing, input.MarketingID).Error; err != nil {
		return nil, gin.H{"message": "Marketing not found", "error": err.Error()}
	}

	// Periksa apakah TransactionNumber sudah ada di tabel Penjualan
	var existingPenjualan models.Penjualan
	if err := config.DB.Where("transaction_number = ?", input.TransactionNumber).First(&existingPenjualan).Error; err == nil {
		return nil, gin.H{
			"errors": []gin.H{
				{"transaction_number": "Transaction number already exists"},
			},
			"message": "Duplicate entry",
		}

	}

	// Parsing date dari string ke time.Time
	const layout = "2006-01-02"
	parsedDate, err := time.Parse(layout, input.Date)
	if err != nil {
		log.Println("Error parsing date:", err)
		return nil, gin.H{"message": "Invalid date format", "error": err.Error()}
	}

	// Konversi data ke model Penjualan asli
	penjualan := models.Penjualan{
		TransactionNumber: input.TransactionNumber,
		MarketingID:       input.MarketingID,
		Date:              parsedDate,
		CargoFee:          input.CargoFee,
		TotalBalance:      input.TotalBalance,
		GrandTotal:        input.GrandTotal,
	}

	// Simpan ke database jika valid
	if err := config.DB.Create(&penjualan).Error; err != nil {
		return nil, gin.H{"message": "Failed to create transaction", "error": err.Error()}
	}

	return gin.H{"message": "Transaction created successfully", "data": penjualan}, nil
}
