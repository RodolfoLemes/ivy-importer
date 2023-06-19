package importer

import (
	"fmt"
)

var invalidMimetype = fmt.Errorf("invalid mimetype detect, only json and csv are accepted")

type ImporterError struct {
	err error
}

func newImporterError(err error) *ImporterError {
	return &ImporterError{err}
}

func (e *ImporterError) Error() string {
	return fmt.Sprintf("importer error: %s", e.err.Error())
}
