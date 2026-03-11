package export

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type (
	CSVExporter interface {
		Generate(data []map[string]interface{}, filename string) (string, error)
	}

	csvExporter struct{}
)

func NewCSVExporter() CSVExporter {
	return &csvExporter{}
}

// Generate writes the provided slice of key/value maps into a CSV file under /tmp.
func (e *csvExporter) Generate(data []map[string]interface{}, filename string) (string, error) {
	if len(data) == 0 {
		return "", errors.New("no data to export")
	}

	if filename == "" {
		filename = fmt.Sprintf("export_%d.csv", time.Now().Unix())
	}
	if filepath.Ext(filename) == "" {
		filename += ".csv"
	}

	headers := extractSortedKeys(data[0])
	filePath := filepath.Join(os.TempDir(), filename)

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.Write(headers); err != nil {
		return "", err
	}

	for _, row := range data {
		record := make([]string, len(headers))
		for i, header := range headers {
			if val, ok := row[header]; ok && val != nil {
				record[i] = fmt.Sprintf("%v", val)
			} else {
				record[i] = ""
			}
		}
		if err := writer.Write(record); err != nil {
			return "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return filePath, nil
}

func extractSortedKeys(row map[string]interface{}) []string {
	keys := make([]string, 0, len(row))
	for k := range row {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
