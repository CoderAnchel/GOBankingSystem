package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// modelo creacion user
type User struct {
	UserNUM            string
	Name               string `json:"name" validate:"required"`
	Password           string `json:"password" validate:"required"`
	Email              string `json:"email" validate:"required,email"`
	Address            string `json:"address" validate:"required"`
	PhoneNumber        string `json:"phoneNumber" validate:"required"`
	InitialAccountName string `json:"initialAccountName" validate:"required"`
	Requests           []string
	Friends            []string
	Color              string
}

type Account struct {
	UserNum             string
	AccountNum          string
	balance             float64
	transferList        []Transfer
	CardList            []Card
	privateTransferList []privateTransfer
	Assets              map[string]Asset
	Name                string
}

type Card struct {
	AccountNum     string
	Name           string `json:"name" validate:"required"`
	Number         string
	CensoredNumber string
	ExpDate        string
	CVV            int
	Balance        float64
	Color          string `json:"color" validate:"required"`
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
	Description     string
	Sender          string
	CardColor       string
	CardName        string
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

type Request struct {
	From string
	To   string
	Date string
}

type PFP struct {
	Color string
}

type Name struct {
	Name string
}

type Phone struct {
	Number string
}

type Search struct {
	UserNUM string
}

type CardOperation struct {
	CardNum string
	Amount  float64
}

var symbolsList = []string{"AAPL", "GOOGL", "TSLA", "AMZN", "MSFT", "NFLX", "FB", "BTC", "ETH", "GOLD", "SILVER"}

var listaUsuarios []User
var listaAccouts []Account
var globalTransferList []Transfer
var validate = validator.New()

var colors = []string{"Purple", "LightBlue", "LightYellow", "LightRed", "LightBrown"}

func cardToAccount(c *fiber.Ctx) error {
	cardOperation := CardOperation{}
	transfer := Transfer{}

	if err := c.BodyParser(&cardOperation); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Error depositing to card",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

	for i := range listaAccouts {
		account := &listaAccouts[i]
		if account.UserNum == userNum {
			for j := range account.CardList {
				card := &account.CardList[j]
				if card.Number == cardOperation.CardNum {

					transfer.ID = uuid.NewString()
					transfer.TransactionType = "CARD_TO_ACCOUNT"
					transfer.Sender = "You"
					transfer.From = account.AccountNum
					transfer.To = card.Number
					transfer.Quantity = cardOperation.Amount
					transfer.Coin = "EUR"
					transfer.Cost = 0
					transfer.Price = 0
					transfer.Date = time.Now().String()
					transfer.CardColor = card.Color
					transfer.CardName = card.Name
					account.transferList = append(account.transferList, transfer)
					card.Balance -= cardOperation.Amount
					account.balance += cardOperation.Amount
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						"status": "Card updated ✅",
						"card":   card,
					})
				}
			}
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"Error": "Error depositing to card",
	})
}

func withdrawFromCard(c *fiber.Ctx) error {
	cardOperation := CardOperation{}
	transfer := Transfer{}

	if err := c.BodyParser(&cardOperation); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Error depositing to card",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

	for i := range listaAccouts {
		account := &listaAccouts[i]
		if account.UserNum == userNum {
			for j := range account.CardList {
				card := &account.CardList[j]
				if card.Number == cardOperation.CardNum {

					transfer.ID = uuid.NewString()
					transfer.TransactionType = "CARD_WITHDRAW"
					transfer.Sender = "You"
					transfer.From = account.AccountNum
					transfer.To = card.Number
					transfer.Quantity = cardOperation.Amount
					transfer.Coin = "EUR"
					transfer.Cost = 0
					transfer.Price = 0
					transfer.Date = time.Now().String()
					transfer.CardColor = card.Color
					transfer.CardColor = card.Color
					transfer.CardName = card.Name
					account.transferList = append(account.transferList, transfer)
					card.Balance -= cardOperation.Amount
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						"status": "Card updated ✅",
						"card":   card,
					})
				}
			}
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"Error": "Error depositing to card",
	})
}

func depositToCard(c *fiber.Ctx) error {
	cardOperation := CardOperation{}
	transfer := Transfer{}

	if err := c.BodyParser(&cardOperation); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Error depositing to card",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

	for i := range listaAccouts {
		account := &listaAccouts[i]
		if account.UserNum == userNum {
			for j := range account.CardList {
				card := &account.CardList[j]
				if card.Number == cardOperation.CardNum {

					if account.balance-cardOperation.Amount < 0 {
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"Error": "Not enought founds ❌",
						})
					}

					account.balance -= cardOperation.Amount
					card.Balance += cardOperation.Amount
					transfer.ID = uuid.NewString()
					transfer.TransactionType = "CARD_DEPOSIT"
					transfer.Sender = "You"
					transfer.From = account.AccountNum
					transfer.To = card.Number
					transfer.Quantity = cardOperation.Amount
					transfer.Coin = "EUR"
					transfer.Cost = 0
					transfer.Price = 0
					transfer.Date = time.Now().String()
					transfer.CardColor = card.Color
					transfer.CardName = card.Name
					account.transferList = append(account.transferList, transfer)
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						"status": "Card updated ✅",
						"card":   card,
					})
				}
			}
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"Error": "Error depositing to card",
	})
}

func getCreditCards(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

	for i := range listaAccouts {
		account := &listaAccouts[i]
		if account.UserNum == userNum {
			return c.Status(fiber.StatusOK).JSON(account.CardList)
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERR:MSG": "Error processing petition ❌",
	})
}

func getAccountDetails(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

	for i := range listaAccouts {
		account := &listaAccouts[i]
		if account.UserNum == userNum {
			return c.Status(fiber.StatusOK).JSON(account)
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERR:MSG": "Error processing petition ❌",
	})
}

func createCard(c *fiber.Ctx) error {
	card := Card{}

	if err := c.BodyParser(&card); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": "Error creating card",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)
	for i := range listaAccouts {
		account := &listaAccouts[i]
		if account.UserNum == userNum {
			card.Number = uuid.NewString()
			card.CensoredNumber = "**** **** **** " + card.Number[len(card.Number)-4:]
			card.Balance = 0
			card.ExpDate = time.Now().AddDate(3, 0, 0).String()
			card.CVV = rand.Intn(900) + 100
			card.AccountNum = account.AccountNum
			account.CardList = append(account.CardList, card)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "Card created ✅",
				"card":   card,
			})
		}
	}

	sendEmail2()

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"Error": "Error creating card",
	})
}

func showFriendData(c *fiber.Ctx) error {
	search := Search{}

	if err := c.BodyParser(&search); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ERR:MSG": "ERROR WHILE PARSING DATA ❌ ",
		})
	}

	for i := range listaUsuarios {
		user := &listaUsuarios[i]
		if user.UserNUM == search.UserNUM {
			name := user.Name
			color := user.Color
			userNum := user.UserNUM
			for i := range listaAccouts {
				account := &listaAccouts[i]
				if account.UserNum == userNum {
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						"name":       name,
						"color":      color,
						"userNum":    userNum,
						"AccountNum": account.AccountNum,
					})
				}
			}
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERR:MSG": "ERROR WHILE PARSING DATA ❌ ",
	})
}

func showFriends(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

	for i := range listaUsuarios {
		account := &listaUsuarios[i]
		if account.UserNUM == userNum {
			return c.Status(fiber.StatusOK).JSON(account.Friends)
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERR:MSG": "ERROR WHILE PROCESSING PETITION ❌",
	})
}

func rejectFriendRequest(c *fiber.Ctx) error {
	search := Search{}

	if err := c.BodyParser(&search); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ERR:MSG": "ERROR WHILE PARSING DATA ❌ ",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

	for i := range listaUsuarios {
		account := &listaUsuarios[i]
		if account.UserNUM == userNum {
			for i := range account.Requests {
				if account.Requests[i] == search.UserNUM {
					account.Requests = append(account.Requests[:i], account.Requests[i+1:]...)
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						"status": "Request Rejected ✅❌",
						"from":   search.UserNUM,
						"to":     account.UserNUM,
					})
				}
			}
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERR:MSG": "ERROR WHILE PARSING DATA ❌ ",
	})
}

func acceptFriendRequest(c *fiber.Ctx) error {
	search := Search{}

	if err := c.BodyParser(&search); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ERR:MSG": "ERROR WHILE PARSING DATA ❌ ",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

	for i := range listaUsuarios {
		account := &listaUsuarios[i]
		if account.UserNUM == userNum {
			for i := range account.Requests {
				if account.Requests[i] == search.UserNUM {
					account.Friends = append(account.Friends, search.UserNUM)
					account.Requests = append(account.Requests[:i], account.Requests[i+1:]...)
				}
			}
		}
	}

	for i := range listaUsuarios {
		account := &listaUsuarios[i]
		if account.UserNUM == search.UserNUM {
			account.Friends = append(account.Friends, userNum)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "Request Accepted ✅",
				"from":   userNum,
				"to":     account.UserNUM,
			})
		}

	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERR:MSG": "ERROR WHILE PARSING DATA ❌ ",
	})
}

func basicUserData(c *fiber.Ctx) error {
	user := Search{}

	if err := c.BodyParser(&user); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ERR:MSG": "ERROR WHILE PARSING DATA ❌ ",
		})
	}

	for i := range listaUsuarios {
		account := &listaUsuarios[i]

		if account.UserNUM == user.UserNUM {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"name":  account.Name,
				"Color": account.Color,
			})
		}

	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERR:MSG": "ERROR WHILE PROCESSING PETITION ❌",
	})
}

func updatePhone(c *fiber.Ctx) error {
	phone := Phone{}

	if err := c.BodyParser(&phone); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ERR:MSG": "ERROR WHILE PARSING DATA ❌ ",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

	for i := range listaUsuarios {
		account := &listaUsuarios[i]

		if account.UserNUM == userNum {
			account.PhoneNumber = phone.Number
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "Name UPDATED ✅",
				"Phone":  account.PhoneNumber,
			})
		}

	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERR:MSG": "ERROR WHILE PROCESSING PETITION ❌",
	})
}

func updateName(c *fiber.Ctx) error {
	name := Name{}

	if err := c.BodyParser(&name); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ERR:MSG": "ERROR WHILE PARSING DATA ❌ ",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

	for i := range listaUsuarios {
		account := &listaUsuarios[i]

		if account.UserNUM == userNum {
			account.Name = name.Name
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": "Name UPDATED ✅",
				"Name":   account.Name,
			})
		}

	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERR:MSG": "ERROR WHILE PROCESSING PETITION ❌",
	})
}

func updatePFP(c *fiber.Ctx) error {
	pfp := PFP{}

	if err := c.BodyParser(&pfp); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ERR:MSG": "ERROR WHILE PARSING DATA ❌ ",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

	for i := range listaUsuarios {
		account := &listaUsuarios[i]

		if account.UserNUM == userNum {
			for color := range colors {
				if pfp.Color == colors[color] {
					account.Color = pfp.Color
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						"status": "PFP UPDATED ✅",
						"color":  account.Color,
					})
				}
			}
		}

	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERR:MSG": "ERROR WHILE PROCESSING PETITION ❌",
	})
}

func sendFriendRequest(c *fiber.Ctx) error {
	request := Request{}
	from := false
	to := false
	name := ""

	if err := c.BodyParser(&request); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ERR:MSG": "ERROR WHILE PARSING DATA ❌ ",
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

	for i := range listaUsuarios {
		account := &listaUsuarios[i]

		if account.UserNUM == userNum {
			from = true
			request.From = account.UserNUM
			name = account.Name
		}

		if account.UserNUM == request.To {
			for i := range account.Requests {
				if account.Requests[i] == userNum {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"ERR:MSG": "ERROR WHILE PROCESSING PETITION ❌",
					})
				}
			}
			for i := range account.Friends {
				if account.Friends[i] == userNum {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"ERR:MSG": "ERROR WHILE PROCESSING PETITION ❌",
					})
				}
			}
			to = true
		}
	}

	if from && to {
		for i := range listaUsuarios {
			account := &listaUsuarios[i]

			if account.UserNUM == request.To {

				sendNewFriendRequestEmail(name)

				account.Requests = append(account.Requests, request.From)
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"status":              "Request sent ✅",
					"from":                request.From,
					"to":                  request.To,
					"To name":             account.Name,
					"All Friend Requests": account.Requests,
				})
			}

		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERR:MSG": "ERROR WHILE PROCESSING PETITION ❌",
	})
}

func getUserDetails(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userNum := claims["userNumber"].(string)

	for _, user := range listaUsuarios {
		if user.UserNUM == userNum {
			return c.Status(fiber.StatusOK).JSON(user)
		}
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"ERR:MSG": "ERROR WHILE PROCESSING PETITION ❌",
	})
}

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
		transaction.Sender = assetsBuy.AssetSymbol
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

			if user.AccountNum == transfer.To {
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
			transfer.Sender = "You"

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
			if user.balance-transfer.Quantity < 0 {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"ERR:MSG": "Not enought founds ❌",
				})
			}

			transfer.From = user.AccountNum
			transfer.To = "N/A"
			transfer.Date = time.Now().String()
			transfer.ID = uuid.NewString()
			transfer.TransactionType = "CASH_WITHDRAW"
			transfer.Sender = "You"

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
				"name":          user.Name,
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
	account.Name = user.InitialAccountName
	user.UserNUM = uuid.NewString()
	user.Color = "LightBlue"
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
	rand.Seed(time.Now().UnixNano())
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Use(logger.New())

	app.Post("/createuser", handleCreateUser)
	app.Post("/login", login)
	app.Static("/", "./public")
	app.Get("/values", values)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))

	app.Get("/user", getUserDetails)
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
	app.Get("showFriends", showFriends)
	app.Post("/showFriendData", showFriendData)
	app.Post("/showUser", basicUserData)
	app.Post("rejectFriendRequest", rejectFriendRequest)
	app.Post("/acceptFriendRequest", acceptFriendRequest)
	app.Post("/sendFriendRequest", sendFriendRequest)
	app.Post("/updatePFP", updatePFP)
	app.Post("/updateName", updateName)
	app.Post("/updatePhone", updatePhone)

	//credit card
	app.Post("/createCard", createCard)
	app.Get("/accountDetails", getAccountDetails)
	app.Get("/creditCards", getCreditCards)
	app.Post("depositCard", depositToCard)
	app.Post("withdrawFromCard", withdrawFromCard)
	app.Post("transferFromCard", cardToAccount)

	//EMAIL SERVER
	app.Get("/sendEmail", sendEmail)

	app.Listen(":3001")
}
