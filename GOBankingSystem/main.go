package main

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
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
func login(c *fiber.Ctx) error {
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
			claims := jwt.MapClaims{
				"name":          user.Name,
				"email":         user.Email,
				"phoneNumber":   user.Email,
				"address":       user.Address,
				"accountNumber": user.AccountNUM,
				"password":      user.Password,
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			t, err := token.SignedString([]byte("secret"))
			if err != nil {
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			return c.JSON(fiber.Map{
				"token": t,
			})
		}
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "BIG ERROR",
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

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome, " + name)
}

func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Post("/createuser", handleCreateUser)

	app.Post("/login", login)

	app.Static("/", "./public")

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))

	app.Get("/restricted", restricted)

	app.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(listaUsuarios)
	})
	app.Listen(":3000")
}
