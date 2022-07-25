package email

import (
	"fmt"
	"io/ioutil"
	"net/smtp"

	txns "github.com/jmsp-test/stori-code-challenge/transactions"
	"gopkg.in/yaml.v3"
)

type UserEmail struct {
	ID           string
	Name         string
	Email        string
	Transactions []txns.Transaction
}

type SenderEmail struct {
	ID       string `yaml:"id"`
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

func ParseUserEmailData(data [][]string) []UserEmail {
	var user UserEmail
	var usersList []UserEmail

	//Loop the data in the .csv file and append it in a UserEmail array
	for i, usr := range data {
		if i == 0 { //Avoid header
			continue
		}
		for j, field := range usr {
			if j == 0 {
				user.ID = field
			} else if j == 1 {
				user.Name = field
			} else if j == 2 {
				user.Email = field
			}
		}
		usersList = append(usersList, user)
	}

	return usersList
}

func BuildNumTxBody(txMonth map[string]int) []string {
	var strgTxMonth []string
	months := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}
	monthName := map[string]string{
		"1": "January", "2": "February", "3": "March", "4": "April",
		"5": "May", "6": "June", "7": "July", "8": "August",
		"9": "September", "10": "October", "11": "November", "12": "December"}

	for _, month := range months {
		//Check months with number of transactions and start building the body for email.
		numTx, exist := txMonth[month]
		if !exist {
			continue
		}
		sTx := fmt.Sprintf("%d", numTx)
		strgTxMonth = append(strgTxMonth, fmt.Sprint("Number of transactions in "+monthName[month]+": "+sTx))
	}

	return strgTxMonth
}

func ReadSenderEmail(filepath string) (SenderEmail, error) {
	configData, err := ioutil.ReadFile(filepath)
	if err != nil {
		return SenderEmail{}, err
	}

	configSender := SenderEmail{}
	err = yaml.Unmarshal(configData, &configSender)
	if err != nil {
		return SenderEmail{}, err
	}

	return configSender, nil
}

func BuildEmail(user UserEmail) (string, string, error) {

	var subject string
	var body string
	var strgNumTx string

	totalBalance, avrgCredit, avrgDebit, err := txns.GetTotalBalance(user.Transactions)
	if err != nil {
		return "", "", err
	}

	//Convert float values into strings
	strgTotalBalance := fmt.Sprintf("%.2f", totalBalance)
	strgAvrgCredit := fmt.Sprintf("%.2f", avrgCredit)
	strgAvrgDebit := fmt.Sprintf("-%.2f", avrgDebit)

	txMonth, err := txns.GetMonthlyTx(user.Transactions)
	if err != nil {
		return "", "", err
	}

	bodyMonthTx := BuildNumTxBody(txMonth)

	for _, tx := range bodyMonthTx {
		//Create string for number of transactions for the email body
		strgNumTx = strgNumTx + tx + "\n"
	}

	//Construct subject and body of the email
	subject = "Subject: Summary Information of your Account Balance\n"
	body = "Dear " + user.Name + ",\n\n" +
		"Thanks again for your preference. Please find your account balance in the mail below.\n\n" +
		"Total Balance is " + strgTotalBalance + "\n" +
		strgNumTx +
		"Average credit amount: " + strgAvrgCredit + "\n" +
		"Average debit amount: " + strgAvrgDebit + "\n\n" +
		"Best Regards,\nStori Card\n"

	// fmt.Println(subject)
	// fmt.Println(body)

	return subject, body, nil
}

func SendEmail(user UserEmail, sender SenderEmail, subject string, body string) error {

	// Sender data.
	from := sender.Email
	password := sender.Password

	// Receiver email address.
	to := []string{
		user.Email,
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
