package migrates

import (
	"codee_jun/internal/models"

	"gorm.io/gorm"
)

func init() {
	regMigrate(func(d *gorm.DB) {
		d.AutoMigrate(&models.User{})
	})
}
