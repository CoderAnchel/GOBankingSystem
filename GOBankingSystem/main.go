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
	UserNUM     string
	Name        string `json:"name" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Address     string `json:"address" validate:"required"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
}

type Account struct {
	UserNum    string
	AccountNum string
	balance    float64
}

type Validation struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Transfer struct {
	From     string
	To       string
	Coin     string  `json:"coin" validate:"required"`
	Quantity float64 `json:"quantity" validate:"required"`
}

var listaUsuarios []User
var listaAccouts []Account
var transferList []Transfer
var validate = validator.New()

func transferTest(c *fiber.Ctx) error {
	transfer := Transfer{}

	if err := c.BodyParser(&transfer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Error generating transacction",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)
	for _, user := range listaAccouts {
		if userNum == user.UserNum {

			transfer.To = user.UserNum
			user.balance += transfer.Quantity

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"accountNumber": user.AccountNum,
				"balance":       user.balance,
			})
		}
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERROR": "Transfer error :/",
	})
}

func getAccount(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

	for _, user := range listaAccouts {
		if userNum == user.UserNum {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"accountNumber": user.AccountNum,
				"balance":       user.balance,
			})
		}
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"error": "no accounts asociated",
	})
}

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
				"name":        user.Name,
				"email":       user.Email,
				"phoneNumber": user.Email,
				"address":     user.Address,
				"userNumber":  user.UserNUM,
				"password":    user.Password,
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
		"error": "Usuario o Contraseña Incorrecta",
	})
}

// handle create user
func handleCreateUser(c *fiber.Ctx) error {
	user := User{}
	account := Account{}
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

	account.AccountNum = uuid.NewString()
	user.UserNUM = uuid.NewString()
	account.UserNum = user.UserNUM
	account.balance = 0
	listaUsuarios = append(listaUsuarios, user)
	listaAccouts = append(listaAccouts, account)
	return c.Status(fiber.StatusOK).JSON(user)
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	email := claims["email"].(string)
	return c.SendString("Welcome, " + name + " " + email)
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

	app.Get("/checkAccounts", getAccount)

	app.Get("/restricted", restricted)

	app.Post("/transferTest", transferTest)

	app.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(listaUsuarios)
	})
	app.Listen(":3000")
}
