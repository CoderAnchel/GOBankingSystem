package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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
	UserNum      string
	AccountNum   string
	balance      float64
	transferList []Transfer
	assets       []Asset
}

type Validation struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Transfer struct {
	ID              string
	TransactionType string
	Date            string
	From            string
	To              string
	Coin            string  `json:"coin" validate:"required"`
	Quantity        float64 `json:"quantity" validate:"required"`
	Cost            float32
}

type AssetsPrice struct {
	AAPL   float64
	GOOGL  float64
	TSLA   float64
	AMZN   float64
	MSFT   float64
	NFLX   float64
	FB     float64
	BTC    float64
	ETH    float64
	XRP    float64
	GOLD   float64
	SILVER float64
}

type Asset struct {
	Symbol         string
	Quantity_Owned float64
	Profit         string
	buyH           []Transfer
}

type AssetsTrans struct {
	AssetSymbol string
	Amount      float32
	Price       float32
}

var symbolsList = []string{"AAPL", "GOOGL", "TSLA", "AMZN", "MSFT", "NFLX", "FB", "BTC", "ETH", "GOLD", "SILVER"}

var listaUsuarios []User
var listaAccouts []Account
var globalTransferList []Transfer
var validate = validator.New()

func values(c *fiber.Ctx) error {
	prices := AssetsPrice{}

	res, err := http.Get("https://faas-lon1-917a94a7.doserverless.co/api/v1/web/fn-e0f31110-7521-4cb9-86a2-645f66eefb63/default/market-prices-simulator")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ERR:MSG": "Error fetching data :/",
		})
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return c.JSON("Error fetching data :/")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ERR:MSG": "Error fetching data :/",
		})
	}

	err = json.Unmarshal(body, &prices)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ERR:MSG": "ERROR UNMARSHALING",
		})
	}

	return c.JSON(prices)
}

func buyAsset(c *fiber.Ctx) error {
	asset := Asset{}
	transaction := Transfer{}
	assetsBuy := AssetsTrans{}
	finded := false

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

	//creating the assetTransaction object
	if err := c.BodyParser(&assetsBuy); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"err": "error wile creating assets object",
		})
	}

	//importing all of the actual stocks proces
	prices := make(map[string]float32)

	res, err := http.Get("https://faas-lon1-917a94a7.doserverless.co/api/v1/web/fn-e0f31110-7521-4cb9-86a2-645f66eefb63/default/market-prices-simulator")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ERR:MSG": "Error fetching data :/",
		})
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return c.JSON("Error fetching data :/")
	}

	if err := json.NewDecoder(res.Body).Decode(&prices); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ERR:MSG": "Error decoding data :/",
		})
	}

	if price, exists := prices[assetsBuy.AssetSymbol]; exists {
		transaction.ID = uuid.NewString()
		transaction.Quantity = float64(assetsBuy.Amount)
		assetsBuy.Price = price
		transaction.Coin = assetsBuy.AssetSymbol
		transaction.Cost = price * assetsBuy.Amount
		finded = true
	}

	if !finded {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "AN ERROR OCURRED symbol, account or balance are wrong ❌",
		})
	}

	for i := range listaAccouts {
		account := &listaAccouts[i]
		if userNum == account.UserNum {
			if account.balance-float64(transaction.Cost) >= 0 {
				transaction.From = account.AccountNum
				transaction.Date = time.Now().String()
				transaction.To = "N/A"
				transaction.TransactionType = "ASSETS_BUY"
				account.balance -= float64(transaction.Cost)
				account.transferList = append(account.transferList, transaction)
				asset.Symbol = transaction.Coin
				asset.Quantity_Owned += transaction.Quantity
				asset.buyH = append(asset.buyH, transaction)
				account.assets = append(account.assets, asset)
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"status":              "ASSETS BUY COMPLETED  ✅",
					"id":                  transaction.ID,
					"asset SYMBOL":        transaction.Coin,
					"assets number":       transaction.Quantity,
					"transactionType":     transaction.TransactionType,
					"transactionDate":     transaction.Date,
					"sourceAccountNumber": transaction.From,
					"targetAccountNumber": transaction.To,
					"asset price":         transaction.Cost,
				})
			} else {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "AN ERROR OCURRED symbol, account or balance are wrong ❌",
				})
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":              "ASSETS BUY COMPLETED  ✅",
		"id":                  transaction.ID,
		"asset SYMBOL":        transaction.Coin,
		"assets number":       transaction.Quantity,
		"transactionType":     transaction.TransactionType,
		"transactionDate":     transaction.Date,
		"sourceAccountNumber": transaction.From,
		"targetAccountNumber": transaction.To,
		"asset price":         transaction.Cost,
	})
}

func transfer(c *fiber.Ctx) error {
	transfer := Transfer{}

	var a bool = false
	var b bool = false

	if err := c.BodyParser(&transfer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Error generating transacction",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)
	for i := range listaAccouts {
		user := &listaAccouts[i]

		if userNum == user.UserNum {

			if user.balance-transfer.Quantity < 0 {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"ERR:MSG": "Not enought founds ❌",
				})
			}
			transfer.From = user.AccountNum
			user.balance -= transfer.Quantity
			a = true
		}

		if user.AccountNum == transfer.To {
			transfer.To = user.AccountNum
			b = true
			user.balance += transfer.Quantity
		}
	}

	if a && b {

		transfer.ID = uuid.NewString()
		transfer.TransactionType = "CASH_TRANSFER"
		transfer.Date = time.Now().String()

		globalTransferList = append(globalTransferList, transfer)
		for i := range listaAccouts {
			user := &listaAccouts[i]
			if user.UserNum == userNum {
				user.transferList = append(user.transferList, transfer)
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":              "DONE ✅",
			"id":                  transfer.ID,
			"amount":              transfer.Quantity,
			"transactionType":     transfer.TransactionType,
			"transactionDate":     transfer.Date,
			"sourceAccountNumber": transfer.From,
			"targetAccountNumber": transfer.To,
		})
	} else {
		for i := range listaAccouts {
			user := &listaAccouts[i]
			if userNum == user.UserNum {
				user.balance += transfer.Quantity
			}
		}
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERROR": "Transfer error :/",
	})
}

func transacctionHistoryAdmin(c *fiber.Ctx) error {
	return c.JSON(globalTransferList)
}

func transacctionHistory(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

	for i := range listaAccouts {
		user := &listaAccouts[i]
		if user.UserNum == userNum {
			return c.Status(fiber.StatusOK).JSON(user.transferList)
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERR:MSG": "Error processing petition ❌",
	})
}

func deposit(c *fiber.Ctx) error {
	transfer := Transfer{}

	if err := c.BodyParser(&transfer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Error generating transacction",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)
	for i := range listaAccouts {
		user := &listaAccouts[i]
		if userNum == user.UserNum {
			transfer.From = user.AccountNum
			transfer.To = "N/A"
			transfer.Date = time.Now().String()
			transfer.ID = uuid.NewString()
			transfer.TransactionType = "CASH_DEPOSIT"

			user.balance += transfer.Quantity
			user.transferList = append(user.transferList, transfer)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":              "DEPOSIT DONE ✅",
				"id":                  transfer.ID,
				"amount":              transfer.Quantity,
				"transactionType":     transfer.TransactionType,
				"transactionDate":     transfer.Date,
				"sourceAccountNumber": transfer.From,
				"targetAccountNumber": transfer.To,
			})
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERROR": "Deposit error :/",
	})
}

func withdraw(c *fiber.Ctx) error {
	transfer := Transfer{}

	if err := c.BodyParser(&transfer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Error generating transacction list :/",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)
	for i := range listaAccouts {
		user := &listaAccouts[i]
		if userNum == user.UserNum {
			transfer.From = user.AccountNum
			transfer.To = "N/A"
			transfer.Date = time.Now().String()
			transfer.ID = uuid.NewString()
			transfer.TransactionType = "CASH_WITHDRAW"

			user.balance -= transfer.Quantity
			user.transferList = append(user.transferList, transfer)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":              "WITHDRAW DONE ✅",
				"id":                  transfer.ID,
				"amount":              transfer.Quantity,
				"transactionType":     transfer.TransactionType,
				"transactionDate":     transfer.Date,
				"sourceAccountNumber": transfer.From,
				"targetAccountNumber": transfer.To,
			})
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERROR": "Withdraw error :/",
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

	app.Get("/values", values)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))

	app.Get("/checkAccounts", getAccount)

	app.Get("/restricted", restricted)

	app.Post("/transferTest", transfer)

	app.Post("/deposit", deposit)
	app.Post("/withdraw", withdraw)
	app.Get("/transacctionHistory", transacctionHistory)
	app.Post("/buyassets", buyAsset)

	app.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(listaUsuarios)
	})
	app.Listen(":3000")
}
