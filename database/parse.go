package database

import (
	"encoding/csv"
	"log"
	"os"
)

type Transaction struct {
	ID   string
	Date string
	Txn  string
}

func ReadCSV(filepath string) [][]string {
	// Open file
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return data
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
