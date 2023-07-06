package csv

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

func WriteCSVFile(filename string, records [][]string) error {
	// Create new file
	file, err := os.Create(ParsePath(filename))
	if err != nil {
		return err
	}
	defer file.Close()

	// Write data to file
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		err := writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadFromCSV(filename string) ([][]string, error) {
	// Open CSV file for reading
	file, err := os.Open(ParsePath(filename))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read data from file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func ParsePath(filename string) string {
	filename = filepath.Clean(filename)
	if filepath.IsAbs(filename) {
		return filename
	}

	filename, _ = filepath.Abs(filename)
	return filename
}
