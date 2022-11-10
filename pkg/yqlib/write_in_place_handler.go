package yqlib

import (
	"os"
)

type writeInPlaceHandler interface {
	CreateTempFile() (*os.File, error)
	FinishWriteInPlace(evaluatedSuccessfully bool) error
}

type writeInPlaceHandlerImpl struct {
	inputFilename string
	tempFile      *os.File
}

func NewWriteInPlaceHandler(inputFile string) writeInPlaceHandler {

	return &writeInPlaceHandlerImpl{inputFile, nil}
}

func (w *writeInPlaceHandlerImpl) CreateTempFile() (*os.File, error) {
	file, err := createTempFile()

	if err != nil {
		return nil, err
	}
	info, err := os.Stat(w.inputFilename)
	if err != nil {
		return nil, err
	}
	err = os.Chmod(file.Name(), info.Mode())

	if err != nil {
		return nil, err
	}
	w.tempFile = file
	return file, err
}

func (w *writeInPlaceHandlerImpl) FinishWriteInPlace(evaluatedSuccessfully bool) error {
	safelyCloseFile(w.tempFile)
	if evaluatedSuccessfully {
		return tryRenameFile(w.tempFile.Name(), w.inputFilename)
	}
	tryRemoveTempFile(w.tempFile.Name())

	return nil
}
