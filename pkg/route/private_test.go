package route

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hrshadhin/fiber-go-boilerplate/app/controller"
	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	"github.com/hrshadhin/fiber-go-boilerplate/app/repository"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/config"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

var adminToken, token string
var userID int
var userIDs [2]int
var bookIDS [10]string

func TestPrivateUserRoutes(t *testing.T) {
	setUpUser()
	defer tearDownUser()

	app := fiber.New()
	PrivateRoutes(app)

	getURLTests := []struct {
		description   string
		route         string
		expectedError bool
		expectedCode  int
		tokenString   string
	}{
		{
			description:   "get user by ID",
			route:         "/api/v1/users/" + strconv.Itoa(userID),
			expectedError: false,
			expectedCode:  200,
			tokenString:   "Bearer " + adminToken,
		},
		{
			description:   "get user by invalid ID",
			route:         "/api/v1/users/123456",
			expectedError: false,
			expectedCode:  404,
			tokenString:   "Bearer " + adminToken,
		},
		{
			description:   "get all user",
			route:         "/api/v1/users",
			expectedError: false,
			expectedCode:  200,
			tokenString:   "Bearer " + adminToken,
		},
		{
			description:   "get all user",
			route:         "/api/v1/users",
			expectedError: false,
			expectedCode:  403,
			tokenString:   "Bearer " + token,
		},
		{
			description:   "get all user",
			route:         "/api/v1/users",
			expectedError: false,
			expectedCode:  400,
			tokenString:   "",
		},
	}

	for _, test := range getURLTests {
		req := httptest.NewRequest("GET", test.route, nil)
		req.Header.Set("accept", "application/json")
		req.Header.Set("Authorization", test.tokenString)

		resp, err := app.Test(req, -1)

		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses,
		// the next test case needs to be processed.
		if test.expectedError {
			continue
		}

		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}

	userString := `{"IsAdmin": false, "IsActive": true, "username": "testrun1", "email": "testrun1@gmail.com", "password": "test123456", "first_name": "T1", "last_name": "R1"}`
	userStringInvalid := `{"IsAdmin": false, "IsActive": true, "username": "testrun2", "email": "testrun1@gmail.com", "password": "test", "first_name": "T1", "last_name": "R1"}`
	updateUser := `{"first_name": "H.R. ", "last_name": " Updated"}`

	postURLTests := []struct {
		description   string
		route         string
		method        string
		tokenString   string
		body          io.Reader
		expectedError bool
		expectedCode  int
	}{
		{
			description:   "create user",
			route:         "/api/v1/users",
			method:        "POST",
			tokenString:   "Bearer " + adminToken,
			body:          strings.NewReader(userString),
			expectedError: false,
			expectedCode:  200,
		},
		{
			description:   "create user invalid payload",
			route:         "/api/v1/users",
			method:        "POST",
			tokenString:   "Bearer " + adminToken,
			body:          strings.NewReader(userStringInvalid),
			expectedError: false,
			expectedCode:  400,
		},
		{
			description:   "create user without right access",
			route:         "/api/v1/users",
			method:        "POST",
			tokenString:   "Bearer " + token,
			body:          strings.NewReader(userString),
			expectedError: false,
			expectedCode:  403,
		},
		{
			description:   "update user",
			route:         "/api/v1/users/" + strconv.Itoa(userID),
			method:        "PUT",
			tokenString:   "Bearer " + adminToken,
			body:          strings.NewReader(updateUser),
			expectedError: false,
			expectedCode:  200,
		},
		{
			description:   "delete a user",
			route:         "/api/v1/users/" + strconv.Itoa(userIDs[1]),
			method:        "DELETE",
			tokenString:   "Bearer " + adminToken,
			body:          nil,
			expectedError: false,
			expectedCode:  200,
		},
	}

	for _, test := range postURLTests {
		req := httptest.NewRequest(test.method, test.route, test.body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", test.tokenString)

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

func TestPrivateBookRoutes(t *testing.T) {
	setUpBook()
	defer tearDownBook()

	app := fiber.New()
	PrivateRoutes(app)

	meta := model.Meta{
		Picture:     RandomString(5) + ".png",
		Description: RandomWord(30),
		Rating:      RandomInt(0, 10),
	}
	book := &model.Book{
		Title:  RandomWord(3),
		Author: RandomWord(2),
		UserID: userIDs[RandomInt(0, 1)],
		Status: RandomInt(0, 1),
		Meta:   meta,
	}

	bookString, _ := json.Marshal(book)
	book.Meta.Rating = 15
	bookStringInvalid, _ := json.Marshal(book)

	book.Title = "Updated title"
	book.Author = "Updated Author"
	book.Meta.Rating = 4
	updateBook, _ := json.Marshal(book)

	postURLTests := []struct {
		description   string
		route         string
		method        string
		tokenString   string
		body          io.Reader
		expectedError bool
		expectedCode  int
	}{
		{
			description:   "create book",
			route:         "/api/v1/books",
			method:        "POST",
			tokenString:   "Bearer " + adminToken,
			body:          strings.NewReader(string(bookString)),
			expectedError: false,
			expectedCode:  200,
		},
		{
			description:   "create book invalid payload",
			route:         "/api/v1/books",
			method:        "POST",
			tokenString:   "Bearer " + adminToken,
			body:          strings.NewReader(string(bookStringInvalid)),
			expectedError: false,
			expectedCode:  400,
		},
		{
			description:   "create book without right access",
			route:         "/api/v1/books",
			method:        "POST",
			tokenString:   "wrong token",
			body:          strings.NewReader(string(bookString)),
			expectedError: false,
			expectedCode:  400,
		},
		{
			description:   "update book",
			route:         "/api/v1/books/" + bookIDS[RandomInt(0, 9)],
			method:        "PUT",
			tokenString:   "Bearer " + token,
			body:          strings.NewReader(string(updateBook)),
			expectedError: false,
			expectedCode:  200,
		},
		{
			description:   "delete a book",
			route:         "/api/v1/books/" + bookIDS[RandomInt(0, 9)],
			method:        "DELETE",
			tokenString:   "Bearer " + token,
			body:          nil,
			expectedError: false,
			expectedCode:  200,
		},
		{
			description:   "delete a unknown book",
			route:         "/api/v1/books/" + uuid.New().String(),
			method:        "DELETE",
			tokenString:   "Bearer " + token,
			body:          nil,
			expectedError: false,
			expectedCode:  404,
		},
	}

	for _, test := range postURLTests {
		req := httptest.NewRequest(test.method, test.route, test.body)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", test.tokenString)

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

func setUpUser() {
	config.LoadAllConfigs("../../.env.test")
	if err := database.ConnectDB(); err != nil {
		panic(err)
	}

	users := []model.CreateUser{
		{
			IsAdmin:   true,
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

	user, err := userRepo.GetByUsername("hrshadhin")
	if err != nil {
		panic(err)
	}
	userID = user.ID
	userIDs[0] = user.ID

	adminToken, err = controller.GenerateNewAccessToken(user.ID, user.IsAdmin)
	if err != nil {
		panic(err)
	}

	user, err = userRepo.GetByUsername("demo1")
	if err != nil {
		panic(err)
	}
	userIDs[1] = user.ID
	token, err = controller.GenerateNewAccessToken(user.ID, user.IsAdmin)
	if err != nil {
		panic(err)
	}

}

func setUpBook() {
	setUpUser()

	bookRepo := repository.NewBookRepo(database.GetDB())
	n := 0
	for n < 10 {
		bookID := uuid.New()
		bookIDS[n] = bookID.String()

		meta := model.Meta{
			Picture:     RandomString(5) + ".png",
			Description: RandomWord(30),
			Rating:      RandomInt(0, 10),
		}
		book := &model.Book{
			ID:     bookID,
			Title:  RandomWord(3),
			Author: RandomWord(2),
			UserID: userIDs[RandomInt(0, 1)],
			Status: RandomInt(0, 1),
			Meta:   meta,
		}
		if err := bookRepo.Create(book); err != nil {
			panic(err)
		}

		n++
	}
}

func tearDownUser() {
	db := database.GetDB()
	_, err := db.Exec(`TRUNCATE TABLE "user" CASCADE;`)
	if err != nil {
		panic(err)
	}

}

func tearDownBook() {
	db := database.GetDB()
	_, err := db.Exec(`TRUNCATE TABLE book`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`TRUNCATE TABLE "user" CASCADE;`)
	if err != nil {
		panic(err)
	}

}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func RandomString(n int) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz"

	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomWord(n int) string {
	words := make([]string, n)
	for i := 0; i < n; i++ {
		words[i] = RandomString(5)
	}

	return strings.Join(words, " ")
}
