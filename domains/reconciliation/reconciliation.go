package reconciliation

import (
	"time"

	"github.com/kopjenmbeng/reconciliation-service/models"
)

type Repository interface {
	ReadTransactionsFromCSV(filePath string, start, end time.Time) ([]models.Transaction, error)
	ReadBankStatementsFromCSV(filePath string, start, end time.Time) ([]models.BankStatement, error)
}

type UseCase interface {
	Reconcile(systemTrxCsvPath string, bankStmntsCsvPath []string,start, end time.Time) models.ReconciliationReport
}
