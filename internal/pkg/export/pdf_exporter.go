package export

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type (
	PDFExporter interface {
		Generate(data []map[string]interface{}, filename, title string) (string, error)
	}

	pdfExporter struct{}
)

func NewPDFExporter() PDFExporter {
	return &pdfExporter{}
}

// Generate creates a simple tabular PDF for the provided dataset.
func (e *pdfExporter) Generate(data []map[string]interface{}, filename, title string) (string, error) {
	if len(data) == 0 {
		return "", errors.New("no data to export")
	}

	if filename == "" {
		filename = fmt.Sprintf("export_%d.pdf", time.Now().Unix())
	}
	if filepath.Ext(filename) == "" {
		filename += ".pdf"
	}

	headers := extractSortedKeys(data[0])

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 14)
	if title == "" {
		title = "Export Data"
	}
	pdf.Cell(40, 10, title)
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 10)
	for _, header := range headers {
		pdf.CellFormat(40, 8, header, "1", 0, "", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 9)
	for _, row := range data {
		for _, header := range headers {
			text := fmt.Sprintf("%v", row[header])
			pdf.CellFormat(40, 7, text, "1", 0, "", false, 0, "")
		}
		pdf.Ln(-1)
	}

	filePath := filepath.Join(os.TempDir(), filename)
	if err := pdf.OutputFileAndClose(filePath); err != nil {
		return "", err
	}

	return filePath, nil
}
