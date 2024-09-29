package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sshaparenko/restApiOnGo/pkg/database"
	"github.com/sshaparenko/restApiOnGo/pkg/domain"
)

// GetAllItems returns array of all items
// of type [domain.Item] from database
func GetAllItems() (items []domain.Item, err error) {
	items = []domain.Item{}
	result := database.DB.Order("created_at desc").Find(&items)
	if result.Error != nil {
		return []domain.Item{}, fmt.Errorf("in services.GetAllItems: %w", result.Error)
	}
	return items, nil
}

// GetItemByID returns [domain.Item] from database
// by its ID
func GetItemByID(id string) (domain.Item, error) {
	var item domain.Item
	result := database.DB.First(&item, "id = ?", id)
	// if result.RowsAffected == 0 {
	// 	return domain.Item{}, errors.New("item not found")
	// }
	if result.Error != nil {
		return domain.Item{}, fmt.Errorf("in services.GetItemByID: %w", result.Error)
	}
	return item, nil
}

// CreateItem creates a new item in database
// with data specified in [domain.ItemRequest]
func CreateItem(itemRequest domain.ItemRequest) (domain.Item, error) {
	var newItem domain.Item = domain.Item{
		ID:        uuid.New().String(),
		Name:      itemRequest.Name,
		Price:     itemRequest.Price,
		Quantity:  itemRequest.Quantity,
		CreatedAt: time.Now(),
	}

	result := database.DB.Create(&newItem)
	if result.Error != nil {
		return domain.Item{}, fmt.Errorf("in services.CreateItem: %w", result.Error)
	}
	return newItem, nil
}

// UpdateItem updates data about item in database
func UpdateItem(itemRequest domain.ItemRequest, id string) (domain.Item, error) {
	item, err := GetItemByID(id)

	if err != nil {
		return domain.Item{}, fmt.Errorf("in services.UpdateItem: %w", err)
	}

	item.Name = itemRequest.Name
	item.Price = itemRequest.Price
	item.Quantity = itemRequest.Quantity
	item.UpdatedAt = time.Now()

	result := database.DB.Save(&item)
	if result.Error != nil {
		return domain.Item{}, fmt.Errorf("in services.UpdateItem: %w", result.Error)
	}

	return item, nil
}

// DeleteItem removes item from database
// based on its id
func DeleteItem(id string) (bool, error) {
	item, err := GetItemByID(id)

	if err != nil {
		return false, fmt.Errorf("in services.DeleteItem: %w", err)
	}

	result := database.DB.Delete(&item)
	if result.Error != nil {
		return false, fmt.Errorf("in services.DeleteItem: %w", err)
	}

	return true, nil
}
