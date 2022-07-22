package transactions

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	database "github.com/jmsp-test/stori-code-challenge/database"
)

func GetTotalBalance(txns []database.Transaction) (string, string, string, error) {
	var valueTx string
	totalBalance, avrgCredit, avrgDebit := float64(0), float64(0), float64(0)
	numCredits, numDebits := float64(0), float64(0)

	for _, tx := range txns {
		//Check that transaction values starts with + or -
		formatCredit := strings.Index(tx.Txn, "+")
		formatDebit := strings.Index(tx.Txn, "-")

		if formatCredit != 0 && formatDebit != 0 {
			return "", "", "", fmt.Errorf("TX ID %s: expects '+' or '-' as first character in Transaction value", tx.ID)
		}

		//If +, calculate credit related variables
		if formatCredit == 0 {
			valueTx = strings.TrimPrefix(tx.Txn, "+")
			n, err := strconv.ParseFloat(valueTx, 64)
			if err != nil {
				return "", "", "", err
			}
			totalBalance += n
			avrgCredit += n
			numCredits += 1
		}

		//If -, calculate debit related variables
		if formatDebit == 0 {
			valueTx = strings.TrimPrefix(tx.Txn, "-")
			n, err := strconv.ParseFloat(valueTx, 64)
			if err != nil {
				return "", "", "", err
			}
			totalBalance -= n
			avrgDebit += n
			numDebits += 1
		}
	}

	//Calculate average values of credit and debit
	avrgCredit = avrgCredit / numCredits
	avrgDebit = avrgDebit / numDebits

	//Round values into closest roof or floor
	totalBalance = math.Round(totalBalance*100) / 100
	avrgCredit = math.Round(avrgCredit*100) / 100
	avrgDebit = math.Round(avrgDebit*100) / 100

	//Convert float values into strings
	strgTotalBalance := fmt.Sprintf("%.2f", totalBalance)
	strgAvrgCredit := fmt.Sprintf("%.2f", avrgCredit)
	strgAvrgDebit := fmt.Sprintf("-%.2f", avrgDebit)

	return strgTotalBalance, strgAvrgCredit, strgAvrgDebit, nil
}

func GetMonthlyTx(txns []database.Transaction) (map[string]int, []string, error) {
	txMonth := make(map[string]int)
	monthValid := map[string]bool{
		"1": true, "2": true, "3": true, "4": true,
		"5": true, "6": true, "7": true, "8": true,
		"9": true, "10": true, "11": true, "12": true}

	for _, tx := range txns {
		//Check that date values have format mm/dd, m/dd, mm/d or m/d
		formatDate := strings.LastIndex(tx.Date, "/")

		if formatDate != 1 && formatDate != 2 {
			return nil, nil, fmt.Errorf("TX ID %s: wrong date format '%s' in Date value: expects mm/dd, m/dd, mm/d or m/d", tx.ID, tx.Date)
		}

		//Validate that month value is between 1 - 12.
		month := strings.Split(tx.Date, "/")
		if !monthValid[month[0]] {
			return nil, nil, fmt.Errorf("TX ID %s: wrong day value '%s' in Date value: expects a value between 1-12", tx.ID, month[0])
		}

		//Add +1 to corresponding month
		txMonth[month[0]] += 1
	}

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

	return txMonth, strgTxMonth, nil
}
