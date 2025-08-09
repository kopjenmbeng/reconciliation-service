package usecase

import (
	"fmt"
	"math"
	"time"

	"github.com/kopjenmbeng/reconciliation-service/domains/reconciliation"
	"github.com/kopjenmbeng/reconciliation-service/models"
)

type UseCase struct {
	Repository reconciliation.Repository
}

func NewReconciliationUseCase(repository reconciliation.Repository) reconciliation.UseCase {
	u := UseCase{
		Repository: repository,
	}
	return &u
}

func (u *UseCase) Reconcile(systemTrxCsvPath string, bankStmntsCsvPath []string, start, end time.Time) models.ReconciliationReport {

	report := models.ReconciliationReport{
		TotalUnmatchedSystem: make([]models.Transaction, 0),
		TotalUnmatchedBank:   make([]models.BankStatement, 0),
	}

	var systemTrx []models.Transaction
	var bankStmnts []models.BankStatement

	// Create channels for system transactions and bank statements
	systemTrxChan := make(chan []models.Transaction)
	bankStmntsChan := make(chan []models.BankStatement)
	errChan := make(chan error, 2) // Buffer size 2 to handle potential errors from both goroutines

	// Goroutine to read system transactions to consider if the file is large
	// and we want to avoid blocking the main thread.
	go func() {
		datas, err := u.Repository.ReadTransactionsFromCSV(systemTrxCsvPath, start, end)
		if err != nil {
			errChan <- err
			close(systemTrxChan)
			return
		}
		systemTrxChan <- datas
		close(systemTrxChan)
	}()

	// Goroutine to read bank statements to consider if the file is large
	// and we want to avoid blocking the main thread.
	go func() {
		var allBankStmnts []models.BankStatement
		for _, path := range bankStmntsCsvPath {
			datas, err := u.Repository.ReadBankStatementsFromCSV(path, start, end)
			if err != nil {
				errChan <- err
				close(bankStmntsChan)
				return
			}
			allBankStmnts = append(allBankStmnts, datas...)
		}
		bankStmntsChan <- allBankStmnts
		close(bankStmntsChan)
	}()

	// Collect results from channels
	select {
	case systemTrx = <-systemTrxChan: // Collect system transactions
	case err := <-errChan:
		fmt.Printf("Error occurred: %v\n", err)
		return report
	}


	select {
	case bankStmnts = <-bankStmntsChan: // Collect bank statements
	case err := <-errChan:
		fmt.Printf("Error occurred: %v\n", err)
		return report
	}

	// counting processed system & bank statement
	report.TotalProcessedSystemTrx = len(systemTrx)
	report.TotalProcessedBankStmnt = len(bankStmnts)

	// Create a map for quick lookup of bank statements.
	// The key is a unique identifier based on date and absolute amount.
	bankStatementMap := make(map[string][]models.BankStatement)
	for _, bankStmnt := range bankStmnts {
		// Normalize amount to be positive for comparison
		normalizedAmount := math.Abs(bankStmnt.Amount)
		key := fmt.Sprintf("%s_%.2f", bankStmnt.Date.Format("2006-01-02"), normalizedAmount)
		bankStatementMap[key] = append(bankStatementMap[key], bankStmnt)
	}

	// Iterate through system transactions and attempt to match.
	matchedBankStatements := make(map[string]bool) // Track matched bank statements to avoid duplicates.
	for _, sysTrx := range systemTrx {
		// Normalize the system transaction amount based on its type.
		normalizedSysAmount := sysTrx.Amount
		if sysTrx.Type == models.DEBIT {
			normalizedSysAmount = math.Abs(sysTrx.Amount)
		}

		key := fmt.Sprintf("%s_%.2f", sysTrx.TransactionTime.Format("2006-01-02"), normalizedSysAmount)

		if matches, ok := bankStatementMap[key]; ok {
			// A match is found. Mark it as matched and handle potential multiple matches.
			// assuming the first available match is the correct one.
			for _, match := range matches {
				// check if this specific bank statement has already been matched.
				matchKey := fmt.Sprintf("%s_%s", match.UniqueIdentifier, key)
				if !matchedBankStatements[matchKey] {
					report.TotalMatched++
					matchedBankStatements[matchKey] = true

					// Check for discrepancies in the amount.
					if math.Abs(sysTrx.Amount) != math.Abs(match.Amount) {
						diff := math.Abs(sysTrx.Amount - math.Abs(match.Amount))
						report.TotalDiscrepancySum += diff
						fmt.Printf("Discrepancy found: System TrxID %s (%.2f) vs Bank ID %s (%.2f). Difference: %.2f\n",
							sysTrx.TrxID, sysTrx.Amount, match.UniqueIdentifier, match.Amount, diff)
					}
					// Break after finding one match to avoid duplicate matching for the same system transaction.
					goto nextSystemTrx
				}
			}
			// If we get here, no matches were available (all were already used).
			report.TotalUnmatchedSystem = append(report.TotalUnmatchedSystem, sysTrx)
		} else {
			// No match found in the bank statement map.
			report.TotalUnmatchedSystem = append(report.TotalUnmatchedSystem, sysTrx)
		}
	nextSystemTrx:
	}

	// Identify unmatched bank statements.
	for _, bs := range bankStmnts {
		// A bank statement is unmatched if its unique key wasn't marked as matched.
		normalizedAmount := math.Abs(bs.Amount)
		key := fmt.Sprintf("%s_%.2f", bs.Date.Format("2006-01-02"), normalizedAmount)
		matchKey := fmt.Sprintf("%s_%s", bs.UniqueIdentifier, key)
		if !matchedBankStatements[matchKey] {
			report.TotalUnmatchedBank = append(report.TotalUnmatchedBank, bs)
		}
	}

	return report

}
 