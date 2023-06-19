package mapper

import (
	"fmt"

	"ivy-importer/importer"
	"ivy-importer/ivy"
	"ivy-importer/loader"
)

type ImportedDataToIvy struct {
	it     importer.DataIterator
	loader loader.Loader

	ivyTransactions []ivy.Transaction
}

func New(it importer.DataIterator, loader loader.Loader) *ImportedDataToIvy {
	return &ImportedDataToIvy{it, loader, []ivy.Transaction{}}
}

func (mapper *ImportedDataToIvy) Exec() ([]ivy.Transaction, error) {
	for mapper.it.HasNext() {
		err := mapper.assign(mapper.it.Next())
		if err != nil {
			return nil, err
		}
	}
	return mapper.ivyTransactions, nil
}

func (mapper *ImportedDataToIvy) assign(data *importer.ImportedData) error {
	iv := &ivy.Transaction{}

	iv.ID = data.ID()
	iv.Title = data.Title

	if data.IsCredit {
		iv.Description = "Credit"
	}

	iv.DateUnixMicro = fmt.Sprintf("%d", data.DateTime().UnixMilli())

	err := mapper.assignTransactionsDetails(iv, data)
	if err != nil {
		return err
	}

	mapper.ivyTransactions = append(mapper.ivyTransactions, *iv)

	return nil
}

func (mapper ImportedDataToIvy) assignTransactionsDetails(
	iv *ivy.Transaction,
	data *importer.ImportedData,
) error {
	account, err := mapper.loader.Account(data.AccountName)
	if err != nil {
		return err
	}

	iv.AccountID = account

	if data.IsTransfer() {
		iv.ToAmount = data.Amount
		iv.Amount = *data.TransferAmount

		toAccount, err := mapper.loader.Account(*data.TransferAccountName)
		if err != nil {
			return err
		}

		iv.AccountID = toAccount
		iv.ToAccountID = &account
		iv.Type = ivy.TransferTransactionType
		return nil
	}

	category, err := mapper.loader.Category(data.CategoryName)
	if err != nil {
		return err
	}

	iv.CategoryID = &category

	amount := data.Amount * -1

	iv.ToAmount = amount
	iv.Amount = amount

	if data.IsExpense() {
		iv.Type = ivy.ExpenseTransactionType
	} else {
		iv.Type = ivy.IncomeTransactionType
	}

	return nil
}
