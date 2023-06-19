package loader

import (
	"fmt"
	"io/ioutil"
	"sort"
	"time"

	"ivy-importer/filemanager"
	"ivy-importer/ivy"
)

type loaderData struct {
	Accounts   []account  `json:"accounts"`
	Categories []category `json:"categories"`
}

type account struct {
	Color                int    `json:"color"`
	Currency             string `json:"currency"`
	Icon                 string `json:"icon"`
	ID                   string `json:"id"`
	HasIncludedInBalance bool   `json:"includeInBalance"`
	IsDeleted            bool   `json:"isDeleted"`
	Name                 string `json:"name"`
	IsSynced             bool   `json:"isSynced"`
	OrderNumber          int    `json:"orderNum"`
}

type category struct {
	Color       int    `json:"color"`
	Icon        string `json:"icon"`
	ID          string `json:"id"`
	IsDeleted   bool   `json:"isDeleted"`
	Name        string `json:"name"`
	IsSynced    bool   `json:"isSynced"`
	OrderNumber int    `json:"orderNum"`
}

type Loader struct {
	Categories map[string]string
	Accounts   map[string]string

	filepath string
}

func New() (*Loader, error) {
	s := &Loader{
		Categories: make(map[string]string),
		Accounts:   make(map[string]string),
	}

	return s, s.load()
}

func (s *Loader) load() error {
	err := s.findAndSetLastSavedFilepath()
	if err != nil {
		return err
	}

	var ld loaderData
	err = filemanager.ReadFromJsonFile(s.filepath, &ld, true)
	if err != nil {
		return err
	}

	for _, c := range ld.Categories {
		s.Categories[c.Name] = c.ID
	}

	for _, a := range ld.Accounts {
		s.Accounts[a.Name] = a.ID
	}

	return nil
}

func (l *Loader) findAndSetLastSavedFilepath() error {
	dirName := "saves"

	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].ModTime().Unix() > files[j].ModTime().Unix()
	})

	l.filepath = dirName + "/" + files[0].Name()

	return nil
}

func (s *Loader) Category(categoryName string) (string, error) {
	c, found := s.Categories[categoryName]
	if !found {
		return "", fmt.Errorf("category not found: %s", categoryName)
	}

	return c, nil
}

func (s *Loader) Account(accountName string) (string, error) {
	c, found := s.Accounts[accountName]
	if !found {
		return "", fmt.Errorf("account not found: %s", accountName)
	}

	return c, nil
}

func (l *Loader) SaveAndReturnFilepath(newTransactions []ivy.Transaction) (string, error) {
	var sd SaveData
	err := filemanager.ReadFromJsonFile(l.filepath, &sd, true)
	if err != nil {
		return "", err
	}

	l.sync(&sd)

	alreadySyncedTransactions := sd.Transactions
	transactions := append(newTransactions, alreadySyncedTransactions...)

	sd.Transactions = transactions

	filepath := l.generateNewSaveFilepath()

	return filepath, filemanager.WriteOnJsonFile(filepath, sd)
}

func (l *Loader) sync(sd *SaveData) {
	for i := range sd.Categories {
		sd.Categories[i].IsSynced = true
	}

	for i := range sd.Settings {
		sd.Settings[i].IsSynced = true
	}

	for i := range sd.Accounts {
		sd.Accounts[i].IsSynced = true
	}

	for i := range sd.Transactions {
		sd.Transactions[i].IsSynced = true
	}
}

func (l Loader) generateNewSaveFilepath() string {
	return fmt.Sprintf("saves/%d.json", time.Now().UnixMilli())
}

type SaveData struct {
	Accounts            []account  `json:"accounts"`
	Budgets             []any      `json:"budgets"`
	Categories          []category `json:"categories"`
	LoanRecords         []any      `json:"loanRecords"`
	Loans               []any      `json:"loans"`
	PlannedPaymentRules []any      `json:"plannedPaymentRules"`
	Settings            []struct {
		BufferAmount int    `json:"bufferAmount"`
		Currency     string `json:"currency"`
		ID           string `json:"id"`
		IsDeleted    bool   `json:"isDeleted"`
		IsSynced     bool   `json:"isSynced"`
		Name         string `json:"name"`
		Theme        string `json:"theme"`
	} `json:"settings"`
	SharedPrefs struct {
		TransfersAsIncExp  string `json:"transfers_as_inc_exp"`
		ShowNotifications  string `json:"show_notifications"`
		HideCurrentBalance string `json:"hide_current_balance"`
		LockApp            string `json:"lock_app"`
	} `json:"sharedPrefs"`
	Transactions []ivy.Transaction `json:"transactions"`
}
