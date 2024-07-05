package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/sshaparenko/restApiOnGo/pkg/domain"
	"github.com/sshaparenko/restApiOnGo/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatasource(dbName string, dbPort string) {

	dsn := generateDsn(dbName, dbPort)

	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to Database")

	if err := DB.AutoMigrate(&domain.User{}, &domain.Item{}); err != nil {
		msg := fmt.Sprintf("Migration error: %s", err.Error())
		panic(msg)
	}
}

func generateDsn(dbName string, dbPort string) string {
	var databaseHost string = os.Getenv("POSTGRES_HOST")
	var db_pass string = string(readDBPass())
	var db_user string = string(readDBUser())

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", databaseHost, db_user, db_pass, dbName, dbPort)

	return dsn
}

func SeedItem() (domain.Item, error) {
	item, err := utils.CreateFaker[domain.Item]()
	if err != nil {
		return domain.Item{}, nil
	}
	DB.Create(&item)
	fmt.Println("Item seeded to the database")

	return item, nil
}

func SeedUser() (domain.User, error) {
	user, err := utils.CreateFaker[domain.User]()

	if err != nil {
		return domain.User{}, err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return domain.User{}, nil
	}

	var inputUser domain.User = domain.User{
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

func readDBPass() []byte {
	pass, err := os.ReadFile("/run/secrets/pg_pass")
	if err != nil {
		log.Fatalf("Failed to read secret: %v", err)
	}
	return pass
}

func readDBUser() []byte {
	username, err := os.ReadFile("/run/secrets/pg_user")
	if err != nil {
		log.Fatalf("Failed to read secret: %v", err)
	}
	return username
}
