package database

import (
	"encoding/csv"
	"os"
)

type Transaction struct {
	ID   string
	Date string
	Txn  string
}

func ReadCSV(filepath string) ([][]string, error) {
	// Open file
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return data, nil
}
