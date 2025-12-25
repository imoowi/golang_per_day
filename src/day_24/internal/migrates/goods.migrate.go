package migrates

import (
	"golang_per_day_24/internal/models"

	"gorm.io/gorm"
)

func init() {
	regMigrate(func(d *gorm.DB) {
		d.AutoMigrate(&models.Goods{})
	})
}
