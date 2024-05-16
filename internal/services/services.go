package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sshaparenko/restApiOnGo/internal/database"
	"github.com/sshaparenko/restApiOnGo/internal/models"
)

func GetAllItems() []models.Item {
	//create a variable to store items data
	var items []models.Item = []models.Item{}
	//get all data from database order by created_at
	database.DB.Order("created_at desc").Find(&items)
	//return all items from database
	return items
}

func GetItemById(id string) (models.Item, error) {
	// create a variable to store item data
	var item models.Item
	// get item data from the database by ID
	result := database.DB.First(&item, "id = ?", id)
	// if the item data is not found, return an error
	if result.RowsAffected == 0 {
		return models.Item{}, errors.New("item not found")
	}
	// return the item data from the database
	return item, nil
}

func CreateItem(itemRequest models.ItemRequest) models.Item {
	// create a new item
	// this item will be inserted to the database
	var newItem models.Item = models.Item{
		ID:        uuid.New().String(),
		Name:      itemRequest.Name,
		Price:     itemRequest.Price,
		Quantity:  itemRequest.Quantity,
		CreatedAt: time.Now(),
	}
	// insert the new item data into the database
	database.DB.Create(&newItem)

	// return the recently inserted item
	return newItem
}

func UpdateItem(itemRequest models.ItemRequest, id string) (models.Item, error) {
	item, err := GetItemById(id)

	if err != nil {
		return models.Item{}, err
	}

	item.Name = itemRequest.Name
	item.Price = itemRequest.Price
	item.Quantity = itemRequest.Quantity
	item.UpdatedAt = time.Now()

	database.DB.Save(&item)

	return item, nil
}

func DeleteItem(id string) bool {
	item, err := GetItemById(id)

	if err != nil {
		return false
	}

	database.DB.Delete(&item)

	return true
}
