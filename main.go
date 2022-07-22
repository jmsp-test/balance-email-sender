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
		Name:         "José Sosa",
		Email:        "jmsosa30@gmail.com",
		Transactions: txData,
	}

	subject, body, err := email.BuildEmail(user)
	if err != nil {
		log.Fatal(err)
	}

	err = email.SendEmail(subject, body)
	if err != nil {
		log.Fatal(err)
	}
}
