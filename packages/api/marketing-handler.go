package api

import (
	"github.com/EKKN/gotestdev/packages/config"
	"github.com/EKKN/gotestdev/packages/helpers"
	"github.com/EKKN/gotestdev/packages/models"
	"github.com/gin-gonic/gin"
)

func (s *APIServer) GetAllMarketing(c *gin.Context) (map[string]interface{}, map[string]interface{}) {

	var marketing []models.Marketing
	config.DB.Find(&marketing)
	return gin.H{"data": marketing}, nil
}

func (s *APIServer) CreateMarketing(c *gin.Context) (map[string]interface{}, map[string]interface{}) {
	var input models.Marketing

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		return nil, gin.H{"message": "Invalid Request Data", "error": err.Error()}
	}

	// Validasi input menggunakan helper ValidateStruct
	validationErrors := helpers.ValidateStruct(input)
	if len(validationErrors) > 0 {
		return nil, gin.H{"message": "Validation errors", "errors": validationErrors}
	}

	// **Cek apakah nama sudah ada di database**
	var existingMarketing models.Marketing
	if err := config.DB.Where("name = ?", input.Name).First(&existingMarketing).Error; err == nil {
		return nil, gin.H{
			"errors": []gin.H{
				{"name": "Marketing name already exists"},
			},
			"message": "Duplicate entry",
		}

	}

	// Simpan ke database jika valid
	if err := config.DB.Create(&input).Error; err != nil {
		return nil, gin.H{"message": "Failed to create marketing", "error": err.Error()}
	}

	return gin.H{"data": input}, nil
}

func (s *APIServer) ToggleMarketingActive(c *gin.Context) (map[string]interface{}, map[string]interface{}) {
	id := c.Param("id")

	// 1. Cek apakah data ada
	var marketing models.Marketing
	if err := config.DB.First(&marketing, id).Error; err != nil {
		return nil, gin.H{"message": "Marketing data not found", "error": err.Error()}
	}

	// 2. Toggle status is_active
	newStatus := !marketing.IsActive // Jika true -> false, jika false -> true
	if err := config.DB.Model(&marketing).Update("is_active", newStatus).Error; err != nil {
		return nil, gin.H{"message": "Failed to update marketing status", "error": err.Error()}
	}

	// 3. Return data yang diperbarui
	return gin.H{"message": "Marketing status updated", "data": marketing}, nil
}
