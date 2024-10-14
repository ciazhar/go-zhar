package service

import (
	"bytes"
	"github.com/tealeg/xlsx"
)

// GenerateXLSX creates an XLSX file and returns it as a byte slice.
func GenerateXLSX() ([]byte, error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return nil, err
	}

	row := sheet.AddRow()
	cell := row.AddCell()
	cell.Value = "Hello, world"

	var buffer bytes.Buffer
	err = file.Write(&buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
