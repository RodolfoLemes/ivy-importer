package importer

import (
	"time"

	"ivy-importer/filemanager"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
)

type DataIterator interface {
	HasNext() bool
	Next() *ImportedData
}

type Importer interface {
	ImportAndValidate(fileName string) (DataIterator, error)
}

func New() Importer {
	return &DefaultImporter{}
}

type DefaultImporter struct {
	data DataIterator

	src []*ImportedData
}

func (i *DefaultImporter) ImportAndValidate(fileName string) (DataIterator, error) {
	if err := i.importFromJson(fileName); err != nil {
		return nil, err
	}

	if err := i.validate(); err != nil {
		return nil, err
	}

	return &ImportedDataIterator{
		ImportedDatas: i.src,
	}, nil
}

func (i DefaultImporter) importFile(fileName string) error {
	mtype, err := mimetype.DetectFile(fileName)
	if err != nil {
		return newImporterError(err)
	}

	if mtype.Is("text/csv") {
		return i.importFromCsv(fileName)
	}

	if mtype.Is("application/json") {
		return i.importFromJson(fileName)
	}

	return newImporterError(invalidMimetype)
}

func (i *DefaultImporter) importFromJson(fileName string) error {
	i.src = []*ImportedData{}

	err := filemanager.ReadFromJsonFile(fileName, &i.src, false)
	if err != nil {
		return newImporterError(err)
	}

	i.data = &ImportedDataIterator{
		ImportedDatas: i.src,
	}

	return nil
}

func (i *DefaultImporter) importFromCsv(fileName string) error {
	i.src = []*ImportedData{}

	err := filemanager.ReadFromCsvFile(fileName, &i.src)
	if err != nil {
		return newImporterError(err)
	}

	i.data = &ImportedDataIterator{
		ImportedDatas: i.src,
	}

	return nil
}

func (i *DefaultImporter) validate() error {
	for i.data.HasNext() {
		if err := i.data.Next().validate(); err != nil {
			return newImporterError(err)
		}
	}

	return nil
}

var dateLayout string = "02/01/2006"

type ImportedData struct {
	id                  string    `json:"-"`
	dateTime            time.Time `json:"-"`
	Date                string
	Title               string
	Amount              float64
	IsCredit            bool
	AccountName         string
	CategoryName        string
	TransferAccountName *string
	TransferAmount      *float64
}

func (i ImportedData) ID() string {
	return i.id
}

func (i ImportedData) DateTime() time.Time {
	return i.dateTime
}

func (i ImportedData) IsTransfer() bool {
	return i.TransferAccountName != nil && i.TransferAmount != nil
}

func (i ImportedData) IsExpense() bool {
	return i.Amount < 0
}

func (i *ImportedData) validate() error {
	if err := i.parseDate(); err != nil {
		return err
	}

	i.assignID()

	return nil
}

func (i *ImportedData) assignID() {
	i.id = uuid.New().String()
}

func (i *ImportedData) parseDate() error {
	dateTime, err := time.Parse(dateLayout, i.Date)
	if err != nil {
		return newImporterError(err)
	}

	i.dateTime = dateTime.Add(8 * time.Hour)
	return nil
}

type ImportedDataIterator struct {
	cursor        int
	ImportedDatas []*ImportedData
}

func (it *ImportedDataIterator) HasNext() bool {
	return it.cursor < len(it.ImportedDatas)
}

func (it *ImportedDataIterator) Next() *ImportedData {
	if it.HasNext() {
		data := it.ImportedDatas[it.cursor]
		it.cursor++
		return data
	}

	return nil
}
