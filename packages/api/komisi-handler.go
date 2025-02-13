package api

import (
	"github.com/EKKN/gotestdev/packages/config"
	"github.com/EKKN/gotestdev/packages/models"
	"github.com/gin-gonic/gin"
)

func (s *APIServer) GetAllKomisi(c *gin.Context) (map[string]interface{}, map[string]interface{}) {

	var komisi []models.KomisiPersen
	config.DB.Order("min_omzet ASC").Find(&komisi) // Urutkan berdasarkan MinOmzet

	return gin.H{"data": komisi}, nil
}
