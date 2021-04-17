package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hrshadhin/fiber-go-boilerplate/app/controller"
	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/app/repository"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/config"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPublicRoutes(t *testing.T) {
	setUpTPuR()
	defer tearDownTPuR()

	app := fiber.New()
	PublicRoutes(app)

	getURLTests := []struct {
		description   string
		route         string
		expectedError bool
		expectedCode  int
	}{
		{
			description:   "get book by ID",
			route:         "/api/v1/books/" + uuid.New().String(),
			expectedError: false,
			expectedCode:  404,
		},
		{
			description:   "get book by invalid ID (non UUID)",
			route:         "/api/v1/books/123456",
			expectedError: false,
			expectedCode:  400,
		},
		{
			description:   "get all books",
			route:         "/api/v1/books",
			expectedError: false,
			expectedCode:  200,
		},
	}

	for _, test := range getURLTests {
		req := httptest.NewRequest("GET", test.route, nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses,
		// the next test case needs to be processed.
		if test.expectedError {
			continue
		}

		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}

	postURLTests := []struct {
		description   string
		route         string
		expectedError bool
		expectedCode  int
		requestBody   io.Reader
	}{
		{
			description:   "get new access token",
			route:         "/api/v1/token/new",
			expectedError: false,
			expectedCode:  200,
			requestBody:   strings.NewReader(`{"username": "hrshadhin", "password": "demo#123"}`),
		},
		{
			description:   "get new access token",
			route:         "/api/v1/token/new",
			expectedError: false,
			expectedCode:  401,
			requestBody:   strings.NewReader(`{"username": "demo1", "password": "test23"}`),
		},
		{
			description:   "get new access token",
			route:         "/api/v1/token/new",
			expectedError: false,
			expectedCode:  404,
			requestBody:   strings.NewReader(`{"username": "test1", "password": "demo#123"}`),
		},
	}

	for _, test := range postURLTests {
		req := httptest.NewRequest("POST", test.route, test.requestBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses,
		// the next test case needs to be processed.
		if test.expectedError {
			continue
		}

		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}

}

func setUpTPuR() {
	config.LoadAllConfigs("../../.env.test")
	if err := database.ConnectDB(); err != nil {
		panic(err)
	}

	users := []model.CreateUser{
		{
			IsAdmin:   false,
			IsActive:  true,
			UserName:  "hrshadhin",
			Email:     "dev@hrshadhin.me",
			Password:  "demo#123",
			FirstName: "H.R.",
			LastName:  "Shadhin",
		},
		{
			IsAdmin:   false,
			IsActive:  true,
			UserName:  "demo1",
			Email:     "demo1@gmail.com",
			Password:  "demo#123",
			FirstName: "Mr.",
			LastName:  "Demo",
		},
	}
	userRepo := repository.NewUserRepo(database.GetDB())
	for _, user := range users {
		user.Password, _ = controller.GeneratePasswordHash([]byte(user.Password))

		if err := userRepo.Create(&user); err != nil {
			panic(err)
		}
	}
}

func tearDownTPuR() {
	db := database.GetDB()
	_, err := db.Exec(`TRUNCATE TABLE "user" CASCADE;`)
	if err != nil {
		panic(err)
	}

}
