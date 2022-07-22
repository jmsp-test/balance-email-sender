package database

import (
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Transaction struct {
	ID   string
	Date string
	Txn  string
}

type Sender struct {
	ID       string `yaml:"id"`
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
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

func ReadSenderEmail(filepath string) (Sender, error) {
	configData, err := ioutil.ReadFile(filepath)
	if err != nil {
		return Sender{}, err
	}

	configSender := Sender{}
	err = yaml.Unmarshal(configData, &configSender)
	if err != nil {
		return Sender{}, err
	}

	return configSender, nil
}
