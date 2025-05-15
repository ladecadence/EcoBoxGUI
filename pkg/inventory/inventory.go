package inventory

import (
	"github.com/ladecadence/EcoBoxGUI/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Inventory struct {
	databaseFile string
	db           *gorm.DB
}

func (i *Inventory) Connect() error {
	database, err := gorm.Open(sqlite.Open(i.databaseFile), &gorm.Config{})
	if err != nil {
		return err
	}

	i.db = database

	err = i.db.AutoMigrate(&models.Tupper{})
	if err != nil {
		return err
	}

	return nil
}

func (i *Inventory) GetTuppers() ([]models.Tupper, error) {
	var tuppers []models.Tupper
	result := i.db.Find(&tuppers)
	return tuppers, result.Error
}

func (i *Inventory) InsertTupper(t models.Tupper) error {
	result := i.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&t)
	return result.Error
}

func (i *Inventory) DeleteTupper(t models.Tupper) error {
	result := i.db.Delete(&t)
	return result.Error
}
