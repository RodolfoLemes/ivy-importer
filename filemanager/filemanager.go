package filemanager

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gocarina/gocsv"
	uni "golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func ReadFromJsonFile[T any](filePath string, src T, isUTC16BE bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	var r io.Reader = file
	if isUTC16BE {
		win16be := uni.UTF16(uni.BigEndian, uni.IgnoreBOM)
		utf16bom := uni.BOMOverride(win16be.NewDecoder())
		r = transform.NewReader(file, utf16bom)
	}

	err = json.NewDecoder(r).Decode(src)

	return err
}

func WriteOnJsonFile[T any](filePath string, data T) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer f.Close()

	t := transform.NewWriter(f, uni.UTF16(uni.BigEndian, uni.UseBOM).NewEncoder())

	defer t.Close()

	err = json.NewEncoder(t).Encode(data)

	return err
}

func WriteOnZipFile(filePath string) error {
	fileName := filepath.Base(filePath)

	zipFileName := fmt.Sprintf("backup-%s.zip", fileName)
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		panic(err)
	}

	defer zipFile.Close()

	w := zip.NewWriter(zipFile)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	wZipFile, err := w.Create(fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(wZipFile, file)
	if err != nil {
		return err
	}

	w.Close()

	return nil
}

func ReadFromCsvFile[T any](filePath string, src T) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	err = gocsv.UnmarshalFile(file, src)

	return err
}
