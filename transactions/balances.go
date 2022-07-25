package transactions

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Transaction struct {
	ID   string
	Date string
	Txn  string
}

func ParseTxnsData(data [][]string) []Transaction {
	var txn Transaction
	var txnsList []Transaction

	//Loop the data in the .csv file and append it in a Transaction array
	for i, tx := range data {
		if i == 0 { //Avoid header
			continue
		}
		for j, field := range tx {
			if j == 0 {
				txn.ID = field
			} else if j == 1 {
				txn.Date = field
			} else if j == 2 {
				txn.Txn = field
			}
		}
		txnsList = append(txnsList, txn)
	}

	return txnsList
}

func GetTotalBalance(txns []Transaction) (float64, float64, float64, error) {
	var valueTx string
	totalBalance, avrgCredit, avrgDebit := float64(0), float64(0), float64(0)
	numCredits, numDebits := float64(0), float64(0)

	for _, tx := range txns {
		//Check that transaction values starts with + or -
		formatCredit := strings.Index(tx.Txn, "+")
		formatDebit := strings.Index(tx.Txn, "-")

		if formatCredit != 0 && formatDebit != 0 {
			return 0, 0, 0, fmt.Errorf("TX ID %s: expects '+' or '-' as first character in Transaction value", tx.ID)
		}

		//If +, calculate credit related variables
		if formatCredit == 0 {
			valueTx = strings.TrimPrefix(tx.Txn, "+")
			n, err := strconv.ParseFloat(valueTx, 64)
			if err != nil {
				return 0, 0, 0, err
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
				return 0, 0, 0, err
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

	return totalBalance, avrgCredit, avrgDebit, nil
}

func GetMonthlyTx(txns []Transaction) (map[string]int, error) {
	txMonth := make(map[string]int)
	monthValid := map[string]bool{
		"1": true, "2": true, "3": true, "4": true,
		"5": true, "6": true, "7": true, "8": true,
		"9": true, "10": true, "11": true, "12": true}

	for _, tx := range txns {
		//Check that date values have format mm/dd, m/dd, mm/d or m/d
		formatDate := strings.LastIndex(tx.Date, "/")

		if formatDate != 1 && formatDate != 2 {
			return nil, fmt.Errorf("TX ID %s: wrong date format '%s' in Date value: expects mm/dd, m/dd, mm/d or m/d", tx.ID, tx.Date)
		}

		//Validate that month value is between 1 - 12.
		month := strings.Split(tx.Date, "/")
		if !monthValid[month[0]] {
			return nil, fmt.Errorf("TX ID %s: wrong day value '%s' in Date value: expects a value between 1-12", tx.ID, month[0])
		}

		//Add +1 to corresponding month
		txMonth[month[0]] += 1
	}

	return txMonth, nil
}
