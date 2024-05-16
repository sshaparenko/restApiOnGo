package database

import (
	"errors"
	"fmt"
	"os"

	"github.com/sshaparenko/restApiOnGo/internal/models"
	"github.com/sshaparenko/restApiOnGo/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatasource(dbName string, dbPort string, dbUser string, dbPassword string) {

	var databaseHost string = os.Getenv("POSTGRES_HOST")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", databaseHost, dbUser, dbPassword, dbName, dbPort)

	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to Database")

	DB.AutoMigrate(&models.User{}, &models.Item{})
}

func SeedItem() (models.Item, error) {
	item, err := utils.CreateFaker[models.Item]()
	if err != nil {
		return models.Item{}, nil
	}
	DB.Create(&item)
	fmt.Println("Item seeded to the database")

	return item, nil
}

func SeedUser() (models.User, error) {
	user, err := utils.CreateFaker[models.User]()

	if err != nil {
		return models.User{}, err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return models.User{}, nil
	}

	var inputUser models.User = models.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: string(password),
	}

	DB.Create(&inputUser)
	fmt.Println("user seeded to database")

	return user, nil
}

func CleanSeeders() {
	itemResult := DB.Exec("TRUNCATE items")
	userResult := DB.Exec("TRUNCATE users")

	var isFailed bool = itemResult.Error != nil || userResult.Error != nil

	if isFailed {
		panic(errors.New("error when cleaning up seeders"))
	}
	fmt.Println("Seeders are cleaned up successfully")
}
