package api

import (
	"time"

	"github.com/EKKN/gotestdev/packages/config"
	"github.com/EKKN/gotestdev/packages/helpers"
	"github.com/EKKN/gotestdev/packages/models"
	"github.com/gin-gonic/gin"
)

func (s *APIServer) GetAllPenjualanNotPembayaran(c *gin.Context) (map[string]interface{}, map[string]interface{}) {
	var penjualans []models.Penjualan

	// Ambil data Penjualan yang belum ada di tabel Pembayaran
	if err := config.DB.
		Preload("Marketing").
		Joins("LEFT JOIN pembayarans ON penjualans.id = pembayarans.penjualan_id").
		Joins("JOIN marketings ON marketings.id = penjualans.marketing_id").
		Where("pembayarans.id IS NULL").
		Order("penjualans.date DESC").
		Find(&penjualans).Error; err != nil {
		return nil, gin.H{"message": "Failed to fetch data", "error": err.Error()}
	}

	return gin.H{"data": penjualans}, nil
}

func (s *APIServer) GetAllPembayaran(c *gin.Context) (map[string]interface{}, map[string]interface{}) {
	var pembayaran []models.Pembayaran

	// Query data pembayaran dengan join ke tabel Penjualan
	if err := config.DB.
		Preload("Penjualan.Marketing").
		Joins("JOIN penjualans ON penjualans.id = pembayarans.penjualan_id").
		Joins("JOIN marketings ON marketings.id = penjualans.marketing_id").
		Order("penjualans.date DESC").
		Find(&pembayaran).Error; err != nil {
		return nil, gin.H{"message": "Failed to fetch pembayaran data", "error": err.Error()}
	}

	return gin.H{"data": pembayaran}, nil
}

func (s *APIServer) CreatePembayaran(c *gin.Context) (map[string]interface{}, map[string]interface{}) {
	type RequestBody struct {
		PenjualanID uint `json:"penjualan_id" validate:"required"`
		KreditID    uint `json:"kredit_id" validate:"required"`
	}

	var input RequestBody

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		return nil, gin.H{"message": "Invalid request data", "error": err.Error()}
	}

	// Validasi input menggunakan validator/v10
	validationErrors := helpers.ValidateStruct(input)
	if len(validationErrors) > 0 {
		return nil, gin.H{"message": "Validation errors", "errors": validationErrors}
	}

	// Validasi data Penjualan
	var penjualan models.Penjualan
	if err := config.DB.First(&penjualan, input.PenjualanID).Error; err != nil {
		return nil, gin.H{"message": "Penjualan not found", "error": err.Error()}
	}

	// Validasi data Kredit
	var kredit models.Kredit
	if err := config.DB.First(&kredit, input.KreditID).Error; err != nil {
		return nil, gin.H{"message": "Kredit not found", "error": err.Error()}
	}

	// Periksa apakah PenjualanID sudah ada di tabel Pembayaran
	var existingPembayaran models.Pembayaran
	if err := config.DB.Where("penjualan_id = ?", input.PenjualanID).First(&existingPembayaran).Error; err == nil {
		return nil, gin.H{
			"errors": []gin.H{
				{"name": "Pembayaran for this Penjualan already exists"},
			},
			"message": "Duplicate entry",
		}

	}

	// Buat entri Pembayaran
	pembayaran := models.Pembayaran{
		PenjualanID:  input.PenjualanID,
		CreditMonths: kredit.Months,
		InterestRate: kredit.InterestRate,
		Installment:  helpers.CalculateInstallment(penjualan.GrandTotal, kredit.InterestRate, kredit.Months),
	}

	// Simpan entri Pembayaran ke database
	if err := config.DB.Create(&pembayaran).Error; err != nil {
		return nil, gin.H{"message": "Failed to create payment", "error": err.Error()}
	}

	// Buat detail pembayaran untuk setiap bulan
	for month := uint(1); month <= pembayaran.CreditMonths; month++ {
		currentTime := time.Now().AddDate(0, int(month), 0)
		detail := models.PembayaranDetail{
			PembayaranID: pembayaran.ID,
			Year:         uint(currentTime.Year()),
			Month:        uint(currentTime.Month()),
			Amount:       pembayaran.Installment,
			Status:       "Pending",
		}
		if err := config.DB.Create(&detail).Error; err != nil {
			return nil, gin.H{"message": "Failed to create payment detail", "error": err.Error()}
		}
	}

	return gin.H{"message": "Payment created successfully", "data": pembayaran}, nil
}
