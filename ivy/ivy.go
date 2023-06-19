package ivy

type TransactionType string

const (
	ExpenseTransactionType  TransactionType = "EXPENSE"
	TransferTransactionType TransactionType = "TRANSFER"
	IncomeTransactionType   TransactionType = "INCOME"
)

type Transaction struct {
	ID            string          `json:"id"`
	Title         string          `json:"title"`
	CategoryID    *string         `json:"categoryId,omitempty"`
	AccountID     string          `json:"accountId"`
	Description   string          `json:"description"`
	DateUnixMicro string          `json:"dateTime"`
	Type          TransactionType `json:"type"`
	Amount        float64         `json:"amount"`
	ToAmount      float64         `json:"toAmount"`
	ToAccountID   *string         `json:"toAccountId,omitempty"`
	IsDeleted     bool            `json:"isDeleted"`
	IsSynced      bool            `json:"isSynced"`
}
