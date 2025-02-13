package api

import (
	"github.com/EKKN/gotestdev/packages/config"
	"github.com/EKKN/gotestdev/packages/helpers"
	"github.com/EKKN/gotestdev/packages/models"
	"github.com/gin-gonic/gin"
)

func (s *APIServer) GetAllPembayaranDetailByPembayaranId(c *gin.Context) (map[string]interface{}, map[string]interface{}) {
	// Ambil `id` dari parameter URL
	pembayaranID := c.Param("pembayaranid")

	var pembayaranDetails []models.PembayaranDetail

	// Cari detail pembayaran berdasarkan ID pembayaran dengan preload yang benar
	if err := config.DB.
		Preload("Pembayaran").                     // Load data Pembayaran
		Preload("Pembayaran.Penjualan").           // Load data Penjualan dari Pembayaran
		Preload("Pembayaran.Penjualan.Marketing"). // Load data Marketing dari Penjualan
		Where("pembayaran_id = ?", pembayaranID).
		Order("year ASC, month ASC").
		Find(&pembayaranDetails).Error; err != nil {
		return nil, gin.H{"message": "Pembayaran not found", "error": err.Error()}
	}

	// Return data jika ditemukan
	return gin.H{"data": pembayaranDetails}, nil
}

func (s *APIServer) UpdatePembayaranDetailStatus(c *gin.Context) (map[string]interface{}, map[string]interface{}) {
	type RequestBody struct {
		Status string `json:"status" validate:"required"`
	}

	var input RequestBody
	pembayaranDetailID := c.Param("id")

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		return nil, gin.H{"message": "Invalid request data", "error": err.Error()}
	}

	// Validasi input menggunakan validator/v10
	validationErrors := helpers.ValidateStruct(input)
	if len(validationErrors) > 0 {
		return nil, gin.H{"message": "Validation errors", "errors": validationErrors}
	}

	// Temukan PembayaranDetail berdasarkan ID
	var pembayaranDetail models.PembayaranDetail
	if err := config.DB.First(&pembayaranDetail, pembayaranDetailID).Error; err != nil {
		return nil, gin.H{"message": "PembayaranDetail not found", "error": err.Error()}
	}

	// Update status PembayaranDetail
	pembayaranDetail.Status = input.Status
	if err := config.DB.Save(&pembayaranDetail).Error; err != nil {
		return nil, gin.H{"message": "Failed to update PembayaranDetail status", "error": err.Error()}
	}

	return gin.H{"message": "PembayaranDetail status updated successfully", "data": pembayaranDetail}, nil
}
