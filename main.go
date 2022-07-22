package main

import (
	"log"

	database "github.com/jmsp-test/stori-code-challenge/database"
	email "github.com/jmsp-test/stori-code-challenge/email"
)

func main() {

	data := database.ReadCSV("./database/txns.csv")
	txData := database.ParseTxnsData(data)

	user := email.User{
		ID:           0,
		Name:         "José Pérez",
		Email:        "jmsosa.tests@gmail.com",
		Transactions: txData,
	}

	sender, err := database.ReadSenderEmail("./database/sender_email.yaml")
	if err != nil {
		log.Fatal(err)
	}

	subject, body, err := email.BuildEmail(user)
	if err != nil {
		log.Fatal(err)
	}

	err = email.SendEmail(user, sender, subject, body)
	if err != nil {
		log.Fatal(err)
	}
}
