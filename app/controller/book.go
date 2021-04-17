package controller

import (
	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hrshadhin/fiber-go-boilerplate/app/model"
	repo "github.com/hrshadhin/fiber-go-boilerplate/app/repository"
	"github.com/hrshadhin/fiber-go-boilerplate/pkg/validator"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/database"
)

// GetBooks func gets all exists book.
// @Description Get all exists book.
// @Summary get all exists book
// @Tags Book
// @Accept json
// @Produce json
// @Param page query integer false "Page no"
// @Param page_size query integer false "records per page"
// @Success 200 {array} model.Book
// @Failure 400 {object} ErrorResponse "Error"
// @Router /v1/books [get]
func GetBooks(c *fiber.Ctx) error {
	pageNo, pageSize := GetPagination(c)
	bookRepo := repo.NewBookRepo(database.GetDB())
	books, err := bookRepo.All(pageSize, uint(pageSize*(pageNo-1)))

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "books were not found",
		})
	}

	return c.JSON(fiber.Map{
		"page":      pageNo,
		"page_size": pageSize,
		"count":     len(books),
		"books":     books,
	})
}

// GetBook func gets a book.
// @Description a book.
// @Summary get a book
// @Tags Book
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} model.Book
// @Failure 400,404 {object} ErrorResponse "Error"
// @Router /v1/books/{id} [get]
func GetBook(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	bookRepo := repo.NewBookRepo(database.GetDB())
	book, err := bookRepo.Get(ID)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "book were not found",
		})
	}

	return c.JSON(fiber.Map{
		"book": book,
	})
}

// CreateBook func for creates a new book.
// @Description Create a new book.
// @Summary create a new book
// @Tags Book
// @Accept json
// @Produce json
// @Param createbook body model.Book true "Create new book"
// @Failure 400,401,500 {object} ErrorResponse status "Error"
// @Success 200 {object} model.Book status "Ok"
// @Security ApiKeyAuth
// @Router /v1/books [post]
func CreateBook(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID, ok := claims["user_id"]
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "can't extract user info from request",
		})
	}

	// Create new Book struct
	book := &model.Book{}
	if err := c.BodyParser(book); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	book.ID = uuid.New()
	book.UserID = int(userID.(float64))
	book.Status = 1 // Active

	// Create a new validator for a Book model.
	validate := validator.NewValidator()
	if err := validate.Struct(book); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	bookRepo := repo.NewBookRepo(database.GetDB())
	if err := bookRepo.Create(book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"book": book,
	})
}

// UpdateBook func update a book.
// @Description update book
// @Summary update a book
// @Tags Book
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param updatebook body model.Book true "Update a book"
// @Success 200 {object} model.Book
// @Failure 400,401,403,404,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /v1/books/{id} [put]
func UpdateBook(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	bookRepo := repo.NewBookRepo(database.GetDB())
	_, err = bookRepo.Get(ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "book were not found",
		})
	}

	book := &model.Book{}
	if err := c.BodyParser(book); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	book.ID = ID

	// Create a new validator for a Book model.
	validate := validator.NewValidator()
	if err := validate.Struct(book); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	if err := bookRepo.Update(ID, book); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	dbBook, err := bookRepo.Get(ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"book": dbBook,
	})
}

// DeleteBook func delete a book.
// @Description delete book
// @Summary delete a book
// @Tags Book
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} interface{} "Ok"
// @Failure 401,403,404,500 {object} ErrorResponse "Error"
// @Security ApiKeyAuth
// @Router /v1/books/{id} [delete]
func DeleteBook(c *fiber.Ctx) error {
	ID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	bookRepo := repo.NewBookRepo(database.GetDB())
	_, err = bookRepo.Get(ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"msg": "book were not found",
		})
	}

	err = bookRepo.Delete(ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{})
}
