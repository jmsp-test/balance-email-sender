package main

import (
	"log"

	database "github.com/jmsp-test/stori-code-challenge/database"
	email "github.com/jmsp-test/stori-code-challenge/email"
	txns "github.com/jmsp-test/stori-code-challenge/transactions"
)

func main() {

	// Take all info from transactions csv file.
	data, err := database.ReadCSV("./database/txns.csv")
	if err != nil {
		log.Fatal(err)
	}
	txData := txns.ParseTxnsData(data)

	// Take user info from csv file and assign transactions data to it.
	data, err = database.ReadCSV("./database/users.csv")
	if err != nil {
		log.Fatal(err)
	}
	user := email.ParseUserEmailData(data)
	user[0].Transactions = txData

	// Obtain sender email from yaml configuration
	sender, err := email.ReadSenderEmail("./email/sender.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Create subject and body of email from user information
	subject, body, err := email.BuildEmail(user[0])
	if err != nil {
		log.Fatal(err)
	}

	// Send email with all previous data.
	err = email.SendEmail(user[0], sender, subject, body)
	if err != nil {
		log.Fatal(err)
	}
}
