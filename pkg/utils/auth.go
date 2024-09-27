package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type TokenMetadata struct {
	Expires int64
}

func ExtractTokenMetata(c *fiber.Ctx) (*TokenMetadata, error) {
	token, err := verifyToken(c)

	if err != nil {
		return nil, fmt.Errorf("verifying token in utils.ExtractTokenMetata: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		expires := int64(claims["exp"].(float64))

		return &TokenMetadata{
			Expires: expires,
		}, nil
	}

	return nil, errors.New("claim data is not found or invalid")
}

// CheckToken returns tocken check result
func CheckToken(c *fiber.Ctx) (bool, error) {
	now := time.Now().Unix()
	claims, err := ExtractTokenMetata(c)

	if err != nil {
		return false, fmt.Errorf("in utils.CheckToken: %w", err)
	}

	expires := claims.Expires

	if now > expires {
		return false, fmt.Errorf("expired token: in utils.CheckToken: %w", err)
	}
	return true, nil
}

func extractToken(c *fiber.Ctx) string {
	//get bearer token from Authoriztion header
	bearToken := c.Get("Authorization")
	//get the JWT token from the bearer
	onlyToken := strings.Split(bearToken, " ")
	if len(onlyToken) == 2 {
		return onlyToken[1]
	}
	//return empty if bearer is empty
	return ""
}

func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	//get token from bearer
	tokenString := extractToken(c)
	//verify the token with JWT secret key
	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	//if verification is failed, rreturn error
	if err != nil {
		return nil, fmt.Errorf("parsing jwt in utils.verifyToken: %w", err)
	}
	//return valid token
	return token, nil
}

// jwtKeyFunc return the JWT secret key
// used to varify the token
func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}

func GenerateNewAccessToken() (string, error) {
	// get JWT secret from .env
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return "", errors.New("jwt secret key is not set")
	}
	// get JWT token expire time from .env
	jwtExpirationTime := os.Getenv("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT")
	if jwtExpirationTime == "" {
		return "", errors.New("jwt expiration time is not set")
	}
	minutesCount, err := strconv.Atoi(jwtExpirationTime)
	if err != nil {
		return "", fmt.Errorf("\"%v\" is ivalid token expiration time", jwtExpirationTime)
	}

	//create JWT claims object
	claims := jwt.MapClaims{}
	//add exp claim to claims
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix()
	//create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//sign token with secret
	t, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", fmt.Errorf("signing token is utils.GenerateNewAccessToken: %w", err)
	}

	return t, nil
}
