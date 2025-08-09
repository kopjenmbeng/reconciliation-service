package models

type ReconciliationReport struct {
	TotalProcessedSystemTrx int
	TotalProcessedBankStmnt int
	TotalMatched            int64
	TotalUnmatchedSystem    []Transaction
	TotalUnmatchedBank      []BankStatement
	TotalDiscrepancySum     float64
}
