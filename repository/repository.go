package repository

// The repository layer is responsible for connecting directly to the database to retrieve and/or modify records.

import (
	"errors"
	"log"

	"github.com/nazerkhh/shop/models"
	"gorm.io/gorm"
)

// ItemRepository is an interface for the struct itemRepository
type ItemRepository interface {
	ListItems() ([]models.Item, error)
	FindItem(id int) (*models.Item, bool, error)
	CreateItem(input *models.Item) (*models.Item, error)
	UpdateItem(id int, input *models.Item) (*models.Item, error)
	DeleteItem(id int) error
}

type itemRepository struct {
	db *gorm.DB
}

func connectDB(dialector gorm.Dialector, config *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, config)
	return db, err
}

// NewItemRepository is a constructor for ItemRepository
func NewItemRepository(dialector gorm.Dialector, config *gorm.Config) *ItemRepository {
	db, err := connectDB(dialector, config)
	if err != nil {
		log.Fatalf("Failed to connect to the database due to error: %s", err)
		return nil
	}

	var r ItemRepository = &itemRepository{db: db}
	return &r
}

func (r *itemRepository) isItemComplete(item *models.Item) bool {
	return item.Name != nil && item.Price != nil && item.Sold != nil // i.e. all non-GORM fields are not nil
}

func (r *itemRepository) ListItems() ([]models.Item, error) {
	var items []models.Item
	err := r.db.Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *itemRepository) FindItem(id int) (*models.Item, bool, error) {
	var item models.Item
	err := r.db.First(&item, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	if !r.isItemComplete(&item) {
		return nil, false, errors.New("broken item found")
	}

	return &item, true, nil
}

func (r *itemRepository) CreateItem(input *models.Item) (*models.Item, error) {
	f := false
	item := models.Item{
		Name:  input.Name,
		Price: input.Price,
		Sold:  &f,
	}
	if err := r.db.Save(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *itemRepository) UpdateItem(id int, input *models.Item) (*models.Item, error) {
	item := models.Item{
		Name:  input.Name,
		Price: input.Price,
		Sold:  input.Sold,
	}

	result := r.db.Model(&item).Where("`items`.`id` = ?", id).Updates(item)
	if err := result.Error; err != nil {
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("no item found to update")
	}
	return &item, nil
}

func (r *itemRepository) DeleteItem(id int) error {
	var item models.Item
	result := r.db.Where("id = ? ", id).Delete(&item)
	if err := result.Error; err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		return errors.New("no item found to delete")
	}
	return nil
}
