package api

import (
	"github.com/EKKN/gotestdev/packages/config"
	"github.com/EKKN/gotestdev/packages/helpers"

	"github.com/gin-gonic/gin"
)

// GET: Laporan komisi per bulan berdasarkan omzet
func (s *APIServer) GetLaporanKomisi(c *gin.Context) (map[string]interface{}, map[string]interface{}) {
	var laporan []struct {
		MarketingID uint
		Marketing   string
		Bulan       string
		TotalOmzet  float64
	}

	// Query untuk mengambil omzet per marketing per bulan
	config.DB.Raw(`
		SELECT 
			m.id AS marketing_id,
			m.name AS marketing,
			DATE_FORMAT(p.date, '%Y-%m') AS bulan,
			SUM(p.grand_total) AS total_omzet
		FROM penjualans p
		JOIN marketings m ON p.marketing_id = m.id
		GROUP BY m.id, DATE_FORMAT(p.date, '%Y-%m')
		ORDER BY DATE_FORMAT(p.date, '%Y-%m') ASC, m.name ASC
	`).Scan(&laporan)

	// Proses data untuk menambahkan komisi
	var hasil []map[string]interface{}
	for _, l := range laporan {
		komisiPersen, komisiNominal := helpers.HitungKomisi(l.TotalOmzet)

		hasil = append(hasil, map[string]interface{}{
			"marketing":      l.Marketing,
			"bulan":          l.Bulan,
			"omzet":          l.TotalOmzet,
			"komisi_persen":  komisiPersen,
			"komisi_nominal": komisiNominal,
		})
	}

	return gin.H{"data": hasil}, nil
}
