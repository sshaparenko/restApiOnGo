package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/sshaparenko/restApiOnGo/internal/database"
	"github.com/sshaparenko/restApiOnGo/internal/models"
	"github.com/sshaparenko/restApiOnGo/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func Signup(userInput models.UserRequest) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	var user models.User = models.User{
		ID: uuid.New().String(),
		Email: userInput.Email,
		Password: string(password),
	}

	result := database.DB.Create(&user)

	if result.RowsAffected == 0 {
		return "", errors.New("user already exists")
	}

	token, err := utils.GenerateNewAccessToken()

	if err != nil {
		return "", err
	}

	return token, nil
}

func Login(userInput models.UserRequest) (string, error) {
	//create a var called user 
	var user models.User
	//find user based on email
	result := database.DB.First(&user, "email = ?", userInput.Email)
	//if user is not found => reurn error
	if result.RowsAffected == 0 {
		return "", errors.New("user not found")
	} 
	//compare password with the password from database
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userInput.Password))
	//if password not match => return error
	if err != nil {
		return "", errors.New("invalid password")
	}
	//generate JWT
	token, err := utils.GenerateNewAccessToken()
	//return error if JWT failed to generate
	if err != nil {
		return "", err
	}
	//return JWT
	return token, nil
}