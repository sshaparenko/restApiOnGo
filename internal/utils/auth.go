package utils

import (
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
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		expires := int64(claims["exp"].(float64))

		return &TokenMetadata{
			Expires: expires,
		}, nil
	}

	return nil, err
}
//returns tocken check result
func CheckToken(c *fiber.Ctx) (bool, error) {
	//get current time
	now := time.Now().Unix()
	//get the token claim data
	claims, err := ExtractTokenMetata(c)
	//if claim data is not found or invalid => return false
	if err != nil {
		return false, err
	}
	//get expiration time from the claim data
	expires := claims.Expires
	//if token is expired return false
	if now > expires {
		return false, err
	}
	//return true, which means that token is valid
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
		return nil, err
	}
	//return valid token
	return token, nil
}

//jwtKeyFunc return the JWT secret key
//used to varify the token
func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(GetValue("JWT_SECRET_KEY")), nil
}

func GenerateNewAccessToken() (string, error) {
	// get JWT secret from .env
	secret := GetValue("JWT_SECRET_KEY")

	// get JWT token expire time from .env
	minutesCount, _ := strconv.Atoi(GetValue("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT"))
	//create JWT claims object
	claims := jwt.MapClaims{}
	//add exp claim to claims
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix()
	//create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//sign token with secret
	t, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return t, nil
}