package utils

import (
	"fmt"

	"github.com/bxcodec/faker/v3"
)

// CreateFaker creates faker data for provided type
func CreateFaker[T any]() (T, error) {
	//create variable ro store faker data
	var fakerData *T = new(T)

	err := faker.FakeData(fakerData)
	if err != nil {
		return *fakerData, fmt.Errorf("generating fake data in utils.CreateFaker: %w", err)
	}

	return *fakerData, nil
}
