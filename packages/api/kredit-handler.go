package api

import (
	"github.com/EKKN/gotestdev/packages/config"
	"github.com/EKKN/gotestdev/packages/models"
	"github.com/gin-gonic/gin"
)

func (s *APIServer) GetAllKredit(c *gin.Context) (map[string]interface{}, map[string]interface{}) {
	var kreditList []models.Kredit
	if err := config.DB.Order("months ASC").Find(&kreditList).Error; err != nil {
		return nil, gin.H{"message": "Failed to fetch kredit data", "error": err.Error()}
	}
	return gin.H{"data": kreditList}, nil
}
