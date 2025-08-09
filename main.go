package main

import (
	"fmt"
	"time"

	reconciliationRepository "github.com/kopjenmbeng/reconciliation-service/domains/reconciliation/repository"
	reconciliationUseCase "github.com/kopjenmbeng/reconciliation-service/domains/reconciliation/usecase"
)

func main() {
	// Initialize the reconciliation service
	fmt.Println("Starting reconciliation service...")
	// init dependencyInjection
	reconciliationRepo := reconciliationRepository.NewRepository()
	reconciliationUCase := reconciliationUseCase.NewReconciliationUseCase(reconciliationRepo)

	bankStatementCsvFiles := []string{"./files/bank-statement.csv"}
	systemTransactionCsvFiles := "./files/system-transaction.csv"

	// init start and end date for reconciliation
	startDate, _ := time.Parse("2006-01-02", "2023-10-25")
	endDate, _ := time.Parse("2006-01-02", "2023-10-28")

	// Perform reconciliation
	report := reconciliationUCase.Reconcile(systemTransactionCsvFiles, bankStatementCsvFiles, startDate, endDate)

	// Print reconciliation report
	fmt.Println("\nReconciliation Summary:")
	fmt.Printf("Total transactions processed (System): %d\n", report.TotalProcessedSystemTrx)
	fmt.Printf("Total transactions processed (Bank): %d\n", report.TotalProcessedBankStmnt)
	fmt.Printf("Total matched transactions: %d\n", report.TotalMatched)
	fmt.Printf("Total unmatched transactions (System): %d\n", len(report.TotalUnmatchedSystem))
	fmt.Printf("Total unmatched transactions (Bank): %d\n", len(report.TotalUnmatchedBank))
	fmt.Printf("Total discrepancy sum: %.2f\n", report.TotalDiscrepancySum)

	fmt.Println("\nDetails of Unmatched Transactions:")
	fmt.Println("---------------------------------------")
	if len(report.TotalUnmatchedSystem) > 0 {
		fmt.Println("System Transactions missing in bank statement(s):")
		for _, trx := range report.TotalUnmatchedSystem {
			fmt.Printf("  - [TrxID: %s, Amount: %.2f, Date: %s]\n", trx.TrxID, trx.Amount, trx.TransactionTime.Format("2006-01-02"))
		}
	}

	if len(report.TotalUnmatchedBank) > 0 {
		fmt.Println("\nBank Statement Transactions missing in system:")
		for _, bs := range report.TotalUnmatchedBank {
			fmt.Printf("  - [ID: %s, Amount: %.2f, Date: %s]\n", bs.UniqueIdentifier, bs.Amount, bs.Date.Format("2006-01-02"))
		}
	}

}
