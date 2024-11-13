package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
)

// modelo creacion user
type User struct {
	AccountNUM  string
	Name        string `json:"name" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Address     string `json:"address" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
}

type Validation struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

var listaUsuarios []User
var accouts []User
var validate = validator.New()

// "log-in"XD
func handleSearchUser(c *fiber.Ctx) error {
	validation := Validation{}
	if err := c.BodyParser(&validation); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error wile parsing body",
		})
	}

	if err := validate.Struct(&validation); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing needed fields at request body",
		})
	}

	for _, user := range listaUsuarios {
		fmt.Println("Comparing user: " + user.Email + " " + user.Password + "with: " + validation.Email + " " + validation.Password)
		if validation.Email == user.Email && validation.Password == user.Password {
			return c.Status(fiber.StatusOK).JSON(user)
		}
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "WRONG EMAIL OR PASSWORD",
	})
}

// handle create user
func handleCreateUser(c *fiber.Ctx) error {
	user := User{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error pasing the body request",
		})
	}

	if err := validate.Struct(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing needed fields at request body",
		})
	}

	for _, u := range listaUsuarios {
		if user.Email == u.Email {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Account created previously with the same email address",
			})
		}

		if user.PhoneNumber == u.PhoneNumber {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Account created previously with the same Phone Number",
			})
		}
	}

	fmt.Println(user)

	user.AccountNUM = uuid.NewString()
	listaUsuarios = append(listaUsuarios, user)
	return c.Status(fiber.StatusOK).JSON(user)
}

func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Static("/", "./public")

	app.Post("/find", handleSearchUser)

	app.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(listaUsuarios)
	})

	app.Post("/createuser", handleCreateUser)

	app.Listen(":3000")
}
