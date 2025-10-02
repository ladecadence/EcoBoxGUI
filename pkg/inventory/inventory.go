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

	err = i.db.AutoMigrate(&models.Container{})
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (i *Inventory) GetContainers() ([]models.Container, error) {
	var containers []models.Container
	result := i.db.Find(&containers)
	return containers, result.Error
}

func (i *Inventory) InsertContainer(t models.Container) error {
	result := i.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&t)
	return result.Error
}

func (i *Inventory) DeleteContainer(t models.Container) error {
	result := i.db.Delete(&t)
	return result.Error
}

func (i *Inventory) DeleteContainerByCode(code string) error {
	result := i.db.Where("code = ?", code).Delete(&models.Container{})
	return result.Error
}

func (i *Inventory) DeleteAll() error {
	result := i.db.Exec("DELETE FROM containers")
	return result.Error
}

func (i *Inventory) GetContainer(code string) (models.Container, error) {
	var container models.Container
	result := i.db.Where("code=?", code).First(&container)
	return container, result.Error
}
