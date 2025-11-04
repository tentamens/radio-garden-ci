package helpers

import (
	"log"
	"os"
)

func LoadTextFile(filePath string) (string, error) {
	loadedTextFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("could not read map file: %v", err)
	}

	fileContent := string(loadedTextFile)

	// reading file seems to always add a new line this removes that
	removeEndingLine := fileContent[:len(fileContent)-1]

	return removeEndingLine, err
}
