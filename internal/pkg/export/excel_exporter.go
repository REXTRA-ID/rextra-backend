package export

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

type (
	ExcelExporter interface {
		Generate(data []map[string]interface{}, filename, sheetName string) (string, error)
		GenerateMultiSheet(sheets map[string][]map[string]interface{}, filename string) (string, error)
	}

	excelExporter struct{}
)

func NewExcelExporter() ExcelExporter {
	return &excelExporter{}
}

func (e *excelExporter) Generate(data []map[string]interface{}, filename, sheetName string) (string, error) {
	if len(data) == 0 {
		return "", errors.New("no data to export")
	}

	if filename == "" {
		filename = fmt.Sprintf("export_%d.xlsx", time.Now().Unix())
	}
	if filepath.Ext(filename) == "" {
		filename += ".xlsx"
	}
	if sheetName == "" {
		sheetName = "Sheet1"
	}

	f := excelize.NewFile()
	f.SetSheetName("Sheet1", sheetName)

	headers := extractSortedKeys(data[0])
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	for rowIdx, row := range data {
		for colIdx, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			f.SetCellValue(sheetName, cell, row[header])
		}
	}

	// Auto-adjust column width based on header length.
	for colIdx := range headers {
		colName, _ := excelize.ColumnNumberToName(colIdx + 1)
		f.SetColWidth(sheetName, colName, colName, 20)
	}

	filePath := filepath.Join(os.TempDir(), filename)
	if err := f.SaveAs(filePath); err != nil {
		return "", err
	}

	return filePath, nil
}

func (e *excelExporter) GenerateMultiSheet(sheets map[string][]map[string]interface{}, filename string) (string, error) {
	if len(sheets) == 0 {
		return "", errors.New("no sheet data to export")
	}

	if filename == "" {
		filename = fmt.Sprintf("export_%d.xlsx", time.Now().Unix())
	}
	if filepath.Ext(filename) == "" {
		filename += ".xlsx"
	}

	f := excelize.NewFile()
	defaultSheet := f.GetSheetName(0)

	isFirst := true
	for sheet, data := range sheets {
		if sheet == "" {
			sheet = "Sheet"
		}
		if isFirst {
			f.SetSheetName(defaultSheet, sheet)
			isFirst = false
		} else {
			f.NewSheet(sheet)
		}

		if len(data) == 0 {
			continue
		}

		headers := extractSortedKeys(data[0])
		for i, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheet, cell, header)
		}

		for rowIdx, row := range data {
			for colIdx, header := range headers {
				cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
				f.SetCellValue(sheet, cell, row[header])
			}
		}

		for colIdx := range headers {
			colName, _ := excelize.ColumnNumberToName(colIdx + 1)
			f.SetColWidth(sheet, colName, colName, 20)
		}
	}

	filePath := filepath.Join(os.TempDir(), filename)
	if err := f.SaveAs(filePath); err != nil {
		return "", err
	}

	return filePath, nil
}
