package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/sshaparenko/restApiOnGo/pkg/database"
	"github.com/sshaparenko/restApiOnGo/pkg/domain"
	"github.com/sshaparenko/restApiOnGo/pkg/utils"
	"github.com/steinfletcher/apitest"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newApp() *fiber.App {
	var app *fiber.App = NewFiberApp()
	initTestDatasource()
	return app
}

func initTestDatasource() {
	var db_host string = os.Getenv("POSTGRES_TEST_HOST")
	var db_pass string = os.Getenv("POSTGRES_TEST_PASSWORD")
	var db_user string = os.Getenv("POSTGRES_TEST_USERNAME")
	var db_name string = "pgtdb"
	var db_port string = "5433"

	fmt.Print("HOST " + db_host)
	fmt.Print("PASS " + db_pass)
	fmt.Print("USER " + db_user)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", db_host, db_user, db_pass, db_name, db_port)

	var err error

	database.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to Database")

	if err := database.DB.AutoMigrate(&domain.User{}, &domain.Item{}); err != nil {
		msg := fmt.Sprintf("Migration error: %s", err.Error())
		panic(msg)
	}
}

func getItem() domain.Item {
	initTestDatasource()

	item, err := database.SeedItem()
	if err != nil {
		panic(err)
	}
	return item
}

func cleanup(res *http.Response, req *http.Request, apiTest *apitest.APITest) {
	if http.StatusOK == res.StatusCode || http.StatusCreated == res.StatusCode {
		database.CleanSeeders()
	}
}

func FiberToHandleFunc(app *fiber.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//test app with fiber
		resp, err := app.Test(r)
		if err != nil {
			panic(err)
		}

		//copy headers
		for k, vv := range resp.Header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(resp.StatusCode)

		if _, err := io.Copy(w, resp.Body); err != nil {
			panic(err)
		}
	}
}

func getJWTToken(t *testing.T) string {
	// connect to the test database
	initTestDatasource()

	// insert a sample data for user into the database
	// the inserted sample data is returned into the "user variable"
	user, err := database.SeedUser()
	if err != nil {
		panic(err)
	}
	// create a request for login
	var userRequest *domain.UserRequest = &domain.UserRequest{
		Email:    user.Email,
		Password: user.Password,
	}
	// get the response from the login request
	var resp *http.Response = apitest.New().
		HandlerFunc(FiberToHandleFunc(newApp())).
		Post("/api/v1/login").
		JSON(userRequest).
		Expect(t).
		Status(http.StatusOK).End().Response
	// create a variable called "response"
	// to store the response body from the login request
	var response *domain.Response[string] = &domain.Response[string]{}
	// decode the response body into the "response" variable
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		msg := fmt.Sprintf("Error has occured while decoding API response: %s", err.Error())
		panic(msg)
	}
	// get the JWT token
	var token string = response.Data
	// create a bearer token
	var JWT_TOKEN string = "Bearer " + token
	// return the bearer token with JWT
	return JWT_TOKEN
}

func TestSignup_Success(t *testing.T) {
	userData, err := utils.CreateFaker[domain.User]()

	if err != nil {
		panic(err)
	}

	var userRequest *domain.UserRequest = &domain.UserRequest{
		Email:    userData.Email,
		Password: userData.Password,
	}

	apitest.New().
		//run the cleanup() after the test
		Observe(cleanup).
		//add an application to be tested
		HandlerFunc(FiberToHandleFunc(newApp())).
		Post("/api/v1/signup").
		JSON(userRequest).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestSignup_ValidationFailed(t *testing.T) {
	var userRequest *domain.UserRequest = &domain.UserRequest{
		Email:    "",
		Password: "",
	}

	apitest.New().
		HandlerFunc(FiberToHandleFunc(newApp())).
		Post("/api/v1/signup").
		JSON(userRequest).
		Expect(t).
		Status(http.StatusBadRequest).
		End()
}

func TestLogin_Success(t *testing.T) {
	initTestDatasource()

	user, err := database.SeedUser()
	if err != nil {
		panic(err)
	}

	var userRequest *domain.UserRequest = &domain.UserRequest{
		Email:    user.Email,
		Password: user.Password,
	}

	apitest.New().
		Observe(cleanup).
		HandlerFunc(FiberToHandleFunc(newApp())).
		Post("/api/v1/login").
		JSON(userRequest).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestLogin_ValidationFailed(t *testing.T) {
	var userRequest *domain.UserRequest = &domain.UserRequest{
		Email:    "",
		Password: "",
	}

	apitest.New().
		HandlerFunc(FiberToHandleFunc(newApp())).
		Post("/api/v1/login").
		JSON(userRequest).
		Expect(t).
		Status(http.StatusBadRequest).
		End()
}

func TestLogin_Failed(t *testing.T) {
	var userRequest *domain.UserRequest = &domain.UserRequest{
		Email:    "notfound@gmail.com",
		Password: "123123",
	}

	apitest.New().
		HandlerFunc(FiberToHandleFunc(newApp())).
		Post("/api/v1/login").
		JSON(userRequest).
		Expect(t).
		Status(http.StatusInternalServerError).
		End()
}

func TestGetItems_Success(t *testing.T) {
	apitest.New().
		HandlerFunc(FiberToHandleFunc(newApp())).
		Get("/api/v1/items").
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestGetItem_Success(t *testing.T) {
	var item domain.Item = getItem()

	apitest.New().
		HandlerFunc(FiberToHandleFunc(newApp())).
		Get("/api/v1/items/" + item.ID).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestGetItem_NotFound(t *testing.T) {
	apitest.New().
		HandlerFunc(FiberToHandleFunc(newApp())).
		Get("/api/v1/items/0").
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestCreateItem_Success(t *testing.T) {
	// create a sample data for item
	itemData, err := utils.CreateFaker[domain.Item]()
	if err != nil {
		panic(err)
	}
	// create a request body to create a new item
	// the request body is filled with the value
	// from the sample data
	var itemRequest *domain.ItemRequest = &domain.ItemRequest{
		Name:     itemData.Name,
		Price:    itemData.Price,
		Quantity: itemData.Quantity,
	}

	// get the JWT token for authentication
	var token string = getJWTToken(t)
	// create a test
	apitest.New().
		Observe(cleanup).
		HandlerFunc(FiberToHandleFunc(newApp())).
		Post("/api/v1/items").
		Header("Authorization", token).
		JSON(itemRequest).
		Expect(t).
		Status(http.StatusCreated).
		End()
}

func TestCreateItem_ValidationFailed(t *testing.T) {
	var itemRequest *domain.ItemRequest = &domain.ItemRequest{
		Name:     "",
		Price:    0,
		Quantity: 0,
	}

	var token string = getJWTToken(t)

	apitest.New().
		HandlerFunc(FiberToHandleFunc(newApp())).
		Post("/api/v1/items").
		Header("Authorization", token).
		JSON(itemRequest).
		Expect(t).
		Status(http.StatusBadRequest).
		End()
}

func TestUpdateItem_Success(t *testing.T) {
	var item domain.Item = getItem()

	var itemRequest *domain.ItemRequest = &domain.ItemRequest{
		Name:     item.Name,
		Price:    item.Price,
		Quantity: item.Quantity,
	}

	var token string = getJWTToken(t)

	apitest.New().
		Observe(cleanup).
		HandlerFunc(FiberToHandleFunc(newApp())).
		Put("/api/v1/items/"+item.ID).
		Header("Authorization", token).
		JSON(itemRequest).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestUpdateItem_Failed(t *testing.T) {
	var itemRequest *domain.ItemRequest = &domain.ItemRequest{
		Name:     "changed",
		Price:    10,
		Quantity: 10,
	}

	var token string = getJWTToken(t)

	apitest.New().
		HandlerFunc(FiberToHandleFunc(newApp())).
		Put("/api/v1/items/0").
		Header("Authorization", token).
		JSON(itemRequest).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}

func TestDeleteItem_Success(t *testing.T) {
	var item domain.Item = getItem()

	var token string = getJWTToken(t)

	apitest.New().
		Observe(cleanup).
		HandlerFunc(FiberToHandleFunc(newApp())).
		Delete("/api/v1/items/"+item.ID).
		Header("Authorization", token).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestDeleteItem_Failed(t *testing.T) {
	var token string = getJWTToken(t)

	apitest.
		New().
		HandlerFunc(FiberToHandleFunc(newApp())).
		Delete("/api/v1/items/0").
		Header("Authorization", token).
		Expect(t).
		Status(http.StatusNotFound).
		End()
}
