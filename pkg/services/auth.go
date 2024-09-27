package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sshaparenko/restApiOnGo/pkg/database"
	"github.com/sshaparenko/restApiOnGo/pkg/domain"
	"github.com/sshaparenko/restApiOnGo/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func Signup(userInput domain.UserRequest) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)

	if err != nil {
		return "", fmt.Errorf("generating password hash: %w", err)
	}

	var user domain.User = domain.User{
		ID:       uuid.New().String(),
		Email:    userInput.Email,
		Password: string(password),
	}

	result := database.DB.Create(&user)

	if result.RowsAffected == 0 {
		return "", errors.New("user already exists")
	}

	token, err := utils.GenerateNewAccessToken()

	if err != nil {
		return "", fmt.Errorf("generating access token in services.Signup: %w", err)
	}

	return token, nil
}

func Login(userInput domain.UserRequest) (string, error) {
	//create a var called user
	var user domain.User
	//find user based on email
	result := database.DB.First(&user, "email = ?", userInput.Email)
	//if user is not found => reurn error
	if result.Error != nil {
		return "", fmt.Errorf("find user by email in services.Login: %w", result.Error)
	}

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
		return "", fmt.Errorf("generating access token in services.Login: %w", err)
	}
	//return JWT
	return token, nil
}
