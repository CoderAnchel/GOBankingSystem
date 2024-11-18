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
	UserNum             string
	AccountNum          string
	balance             float64
	transferList        []Transfer
	privateTransferList []privateTransfer
	Assets              map[string]Asset
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
	Price           float32
}

// for selling assets
type privateTransfer struct {
	ID              string
	TransactionType string
	Date            string
	From            string
	To              string
	Coin            string  `json:"coin" validate:"required"`
	Quantity        float64 `json:"quantity" validate:"required"`
	Cost            float32
	Price           float32
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
	Profit         float32
	buyH           []Transfer
}

type AssetsTrans struct {
	AssetSymbol string
	Amount      float32
	Price       float32
}

type assetCheck struct {
	AssetSymbol string
}

var symbolsList = []string{"AAPL", "GOOGL", "TSLA", "AMZN", "MSFT", "NFLX", "FB", "BTC", "ETH", "GOLD", "SILVER"}

var listaUsuarios []User
var listaAccouts []Account
var globalTransferList []Transfer
var validate = validator.New()

// in progress...
func sellAsset(c *fiber.Ctx) error {
	assetCheck := assetCheck{}
	if err := c.BodyParser(&assetCheck); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ERR:MSG": "ERROR WHILE PARSING DATA ❌ ",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

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

	for i := range listaAccouts {
		account := &listaAccouts[i]
		if account.UserNum == userNum {
			now := time.Now()
			var actualPrice float32
			if price, exists := prices[assetCheck.AssetSymbol]; exists {
				actualPrice = price
			}

			if asset, exists := account.Assets[assetCheck.AssetSymbol]; exists {
				asset.Profit = 0
				for i := range asset.buyH {
					costIfNow := actualPrice * float32(asset.buyH[i].Quantity)
					if costIfNow < asset.buyH[i].Cost {
						result := asset.buyH[i].Cost - costIfNow
						asset.Profit += result
					} else if costIfNow > asset.buyH[i].Cost {
						result := costIfNow - asset.buyH[i].Cost
						asset.Profit += result
					}
				}

				response := fiber.Map{
					"time":       now,
					"buyH":       asset.buyH,
					"depostited": asset.Profit,
				}

				// Send the response
				err := c.JSON(response)
				if err != nil {
					return err
				}

				return nil
			}
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERR:MSG": "ERROR WHILE PROCESSING PETITION ❌",
	})
}

// change from the history to the actual liquidation list
func checkAssetTrans(c *fiber.Ctx) error {
	assetCheck := assetCheck{}
	if err := c.BodyParser(&assetCheck); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ERR:MSG": "ERROR WHILE PARSING DATA ❌ ",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

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

	for i := range listaAccouts {
		account := listaAccouts[i]
		if account.UserNum == userNum {
			now := time.Now()
			var actualPrice float32
			if price, exists := prices[assetCheck.AssetSymbol]; exists {
				actualPrice = price
			}

			if asset, exists := account.Assets[assetCheck.AssetSymbol]; exists {
				for i := range asset.buyH {
					costIfNow := actualPrice * float32(asset.buyH[i].Quantity)
					if costIfNow < asset.buyH[i].Cost {
						result := asset.buyH[i].Cost - costIfNow
						asset.Profit += result
					} else if costIfNow > asset.buyH[i].Cost {
						result := costIfNow - asset.buyH[i].Cost
						asset.Profit += result
					}
				}

				response := fiber.Map{
					"time":   now,
					"buyH":   asset.buyH,
					"profit": asset.Profit,
				}

				// Send the response
				err := c.JSON(response)
				if err != nil {
					return err
				}

				return nil
			}
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERR:MSG": "ERROR WHILE SEARCHING DATA ❌ ",
	})
}

// change to liquidation list
func checkAssets(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

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
	for i := range listaAccouts {
		account := &listaAccouts[i]
		if account.UserNum == userNum {
			var actualPrice float32
			for symbol, asset := range account.Assets {
				asset.Profit = 0 // Reset profit before calculation
				if price, exists := prices[asset.Symbol]; exists {
					actualPrice = price
				}

				for j := range asset.buyH {
					costIfNow := actualPrice * float32(asset.buyH[j].Quantity)
					if costIfNow < asset.buyH[j].Cost {
						result := asset.buyH[j].Cost - costIfNow
						asset.Profit += result
					} else if costIfNow > asset.buyH[j].Cost {
						result := costIfNow - asset.buyH[j].Cost
						asset.Profit += result
					}
				}
				account.Assets[symbol] = asset
			}

			return c.Status(fiber.StatusOK).JSON(account.Assets)
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON("ERROR :/")
}

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

// create transactionPrivate to
func buyAsset(c *fiber.Ctx) error {
	//asset := Asset{}
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
		transaction.Price = price
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

				if asset, exists := account.Assets[assetsBuy.AssetSymbol]; exists {
					asset.Symbol = transaction.Coin
					asset.Quantity_Owned += transaction.Quantity
					asset.buyH = append(asset.buyH, transaction)
					account.Assets[assetsBuy.AssetSymbol] = asset
				}
				// asset.Symbol = transaction.Coin
				// asset.Quantity_Owned += transaction.Quantity
				// asset.buyH = append(asset.buyH, transaction)
				// account.assets = append(account.assets, asset)

				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"status":              "ASSETS BUY COMPLETED  ✅",
					"id":                  transaction.ID,
					"asset SYMBOL":        transaction.Coin,
					"assets quantity":     transaction.Quantity,
					"transactionType":     transaction.TransactionType,
					"transactionDate":     transaction.Date,
					"sourceAccountNumber": transaction.From,
					"targetAccountNumber": transaction.To,
					"assets cost":         transaction.Cost,
					"assets price":        transaction.Price,
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

	asset1 := Asset{
		Symbol:         "AAPL",
		Quantity_Owned: 0,
		Profit:         0,
		buyH:           []Transfer{},
	}

	asset2 := Asset{
		Symbol:         "GOOGL",
		Quantity_Owned: 0,
		Profit:         0,
		buyH:           []Transfer{},
	}

	asset3 := Asset{
		Symbol:         "TSLA",
		Quantity_Owned: 0,
		Profit:         0,
		buyH:           []Transfer{},
	}

	asset4 := Asset{
		Symbol:         "AMZN",
		Quantity_Owned: 0,
		Profit:         0,
		buyH:           []Transfer{},
	}

	asset5 := Asset{
		Symbol:         "MSFT",
		Quantity_Owned: 0,
		Profit:         0,
		buyH:           []Transfer{},
	}

	asset6 := Asset{
		Symbol:         "NFLX",
		Quantity_Owned: 0,
		Profit:         0,
		buyH:           []Transfer{},
	}

	asset7 := Asset{
		Symbol:         "FB",
		Quantity_Owned: 0,
		Profit:         0,
		buyH:           []Transfer{},
	}

	asset8 := Asset{
		Symbol:         "BTC",
		Quantity_Owned: 0,
		Profit:         0,
		buyH:           []Transfer{},
	}

	asset9 := Asset{
		Symbol:         "ETH",
		Quantity_Owned: 0,
		Profit:         0,
		buyH:           []Transfer{},
	}

	asset10 := Asset{
		Symbol:         "GOLD",
		Quantity_Owned: 0,
		Profit:         0,
		buyH:           []Transfer{},
	}

	asset11 := Asset{
		Symbol:         "SILVER",
		Quantity_Owned: 0,
		Profit:         0,
		buyH:           []Transfer{},
	}

	account.Assets = make(map[string]Asset)
	account.AccountNum = uuid.NewString()
	user.UserNUM = uuid.NewString()
	account.UserNum = user.UserNUM
	account.balance = 0
	account.Assets["AAPL"] = asset1
	account.Assets["GOOGL"] = asset2
	account.Assets["TSLA"] = asset3
	account.Assets["AMZN"] = asset4
	account.Assets["MSFT"] = asset5
	account.Assets["NFLX"] = asset6
	account.Assets["FB"] = asset7
	account.Assets["BTC"] = asset8
	account.Assets["ETH"] = asset9
	account.Assets["GOLD"] = asset10
	account.Assets["SILVER"] = asset11
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
	app.Get("/checkAssets", checkAssets)
	app.Post("/checkAssetTrans", checkAssetTrans)
	app.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(listaUsuarios)
	})

	app.Listen(":3000")
}
