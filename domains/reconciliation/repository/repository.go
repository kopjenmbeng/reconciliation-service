package repository

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/kopjenmbeng/reconciliation-service/domains/reconciliation"
	"github.com/kopjenmbeng/reconciliation-service/models"
)

type Repository struct {
}

func NewRepository() reconciliation.Repository {
	return &Repository{}
}

func (r *Repository) ReadTransactionsFromCSV(filePath string, start, end time.Time) ([]models.Transaction, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var transactions []models.Transaction

	// Skip header row
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read record: %w", err)
		}

		// Ensure the record has the correct number of fields
		if len(record) < 4 {
			log.Printf("Skipping malformed transaction record: %v\n", record)
			continue
		}

		amount, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Printf("Skipping record due to amount parsing error: %v\n", err)
			continue
		}

		transactionTime, err := time.Parse("2006-01-02", record[3])
		if err != nil {
			log.Printf("Skipping record due to time parsing error: %v\n", err)
			continue
		}

		// Filter by date range
		if !transactionTime.Before(start) && !transactionTime.After(end) {
			transactions = append(transactions, models.Transaction{
				TrxID:           record[0],
				Amount:          amount,
				Type:            models.TransactionType(record[2]),
				TransactionTime: transactionTime,
			})
		}
	}
	return transactions, nil
}

func (r *Repository) ReadBankStatementsFromCSV(filePath string, start, end time.Time) ([]models.BankStatement, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var statements []models.BankStatement

	// Skip header row
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read record: %w", err)
		}

		// Ensure the record has the correct number of fields
		if len(record) < 3 {
			log.Printf("Skipping malformed bank statement record: %v\n", record)
			continue
		}

		amount, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Printf("Skipping record due to amount parsing error: %v\n", err)
			continue
		}

		date, err := time.Parse("2006-01-02", record[2])
		if err != nil {
			log.Printf("Skipping record due to date parsing error: %v\n", err)
			continue
		}

		// Filter by date range
		if !date.Before(start) && !date.After(end) {
			statements = append(statements, models.BankStatement{
				UniqueIdentifier: record[0],
				Amount:           amount,
				Date:             date,
			})
		}
	}
	return statements, nil
}
