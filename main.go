package main

import (
	"flag"
	"fmt"

	"ivy-importer/filemanager"
	"ivy-importer/importer"
	"ivy-importer/loader"
	"ivy-importer/mapper"
)

var importFilePath string

const defaultImportFilePath = ""

func init() {
	flag.StringVar(
		&importFilePath,
		"filepath",
		defaultImportFilePath,
		"the filepath where the file is",
	)
}

func main() {
	flag.Parse()
	run()

	fmt.Println("Successfully saved!")
}

func run() {
	imp, err := importer.New().ImportAndValidate(importFilePath)
	if err != nil {
	}

	loader, err := loader.New()
	if err != nil {
		panic(err)
	}

	mapper := mapper.New(imp, *loader)
	transactions, err := mapper.Exec()
	if err != nil {
		panic(err)
	}

	savedFilepath, err := loader.SaveAndReturnFilepath(transactions)
	if err != nil {
		panic(err)
	}

	err = filemanager.WriteOnZipFile(savedFilepath)
	if err != nil {
		panic(err)
	}
}
