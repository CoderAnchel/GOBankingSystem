package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"log"
	"net/smtp"
)

type CardRelated struct {
	CardNumber string
	CardHolder string
	CardType   string
	CardExpiry string
}

func sendEmail(c *fiber.Ctx) error {
	auth := smtp.PlainAuth(
		"",
		"anchelxlaid@gmail.com",
		"qiqt ajjj pwsl ybhq",
		"smtp.gmail.com",
	)

	msg := "Subject: Hello, World!\n\nThis is the email body."

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"anchelxlaid@gmail.com",
		[]string{"anchelxlaid@gmail.com"},
		[]byte(msg),
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Email error 🙅",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Email sent",
	})
}

func sendNewFriendRequestEmail(name string) {
	auth := smtp.PlainAuth(
		"",
		"anchelxlaid@gmail.com",
		"qiqt ajjj pwsl ybhq",
		"smtp.gmail.com",
	)

	msg := fmt.Sprintf("Subject: NEW FRIEND REQUEST ⭐️\n\nYou have a new friend request from: %s", name)

	smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"anchelxlaid@gmail.com",
		[]string{"anchelxlaid@gmail.com"},
		[]byte(msg),
	)
}

func sendEmail2() {
	htmlContent, err := ioutil.ReadFile("./public/NewCredit.html")
	if err != nil {
		log.Fatalf("Error al leer el archivo HTML: %v", err)
	}

	emailBody := string(htmlContent)

	subject := "Credit Card Created"
	from := "anchelxlaid@gmail.com"
	to := "anchelxlaid@gmail.com"
	password := "qiqt ajjj pwsl ybhq"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		emailBody)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		log.Fatalf("Error al enviar el correo: %v", err)
	}

	log.Println("Correo enviado exitosamente")
}
