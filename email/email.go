package email

import (
	"fmt"
	"net/smtp"

	database "github.com/jmsp-test/stori-code-challenge/database"
	transactions "github.com/jmsp-test/stori-code-challenge/transactions"
)

type User struct {
	Name         string
	Email        string
	Transactions []database.Transaction
}

func BuildEmail(user User) (string, string, error) {

	var subject string
	var body string
	var strgNumTx string

	totalBalance, avrgCredit, avrgDebit, err := transactions.GetTotalBalance(user.Transactions)
	if err != nil {
		return "", "", err
	}
	_, bodyMonthTx, err := transactions.GetMonthlyTx(user.Transactions)
	if err != nil {
		return "", "", err
	}

	for _, tx := range bodyMonthTx {
		//Create string for number of transactions for the email body
		strgNumTx = strgNumTx + tx + "\n"
	}

	//Construct subject and body of the email
	subject = "Subject: Summary Information of your Account Balance\n"
	body = "Dear " + user.Name + ",\n\n" +
		"Thanks again for your preference. Please find your account balance in the mail below.\n\n" +
		"Total Balance is " + totalBalance + "\n" +
		strgNumTx +
		"Average debit amount: " + avrgDebit + "\n" +
		"Average credit amount: " + avrgCredit + "\n\n" +
		"Best Regards,\nStori Card\n"

	fmt.Println(subject)
	fmt.Println(body)

	return subject, body, nil
}

func SendEmail(subject string, body string) error {

	// Sender data.
	from := "jmsosa.tests@gmail.com"
	password := "zdjeymjaozynvdlp"

	// Receiver email address.
	to := []string{
		"jmsosa.tests@gmail.com",
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	address := smtpHost + ":" + smtpPort

	// Message.
	message := []byte(subject + body)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		return err
	}
	fmt.Println("Email Sent Successfully!")

	return nil
}
