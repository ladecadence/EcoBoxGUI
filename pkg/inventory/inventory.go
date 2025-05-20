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

func New(file string) (*Inventory, error) {
	i := Inventory{databaseFile: file}

	database, err := gorm.Open(sqlite.Open(i.databaseFile), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	i.db = database

	err = i.db.AutoMigrate(&models.Tupper{})
	if err != nil {
		return nil, err
	}

	return &i, nil
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

func (i *Inventory) DeleteTupperByNum(number int) error {
	result := i.db.Where("number = ?", number).Delete(&models.Tupper{})
	return result.Error
}

func (i *Inventory) GetTupper(id string) (models.Tupper, error) {
	var tupper models.Tupper
	result := i.db.Where("id=?", id).First(&tupper)
	return tupper, result.Error
}
