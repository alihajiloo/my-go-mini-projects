package main

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {

		m := gomail.NewMessage()
		m.SetHeader("From", "info@alihajiloo.ir")
		m.SetHeader("To", "haji.hd99@gmail.com", "alihajiloo.ir@gmail.com")
		//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
		m.SetHeader("Subject", "Hello!")
		m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
		m.Attach("./inoti.pdf")

		d := gomail.NewDialer("mail.alihajiloo.ir", 587, "info@alihajiloo.ir", "09375750383aLi")

		// Send the email to Bob, Cora and Dan.
		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":2020") // listen and serve on 0.0.0.0:8080
}