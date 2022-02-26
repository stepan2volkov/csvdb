package tableloader

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/stepan2volkov/csvdb/internal/table"
)

func Load(path string, sep rune, lazyQuotes bool, fields []table.Field) (table.Table, error) {
	file, err := os.Open(path)
	if err != nil {
		return table.Table{}, err
	}

	reader := csv.NewReader(file)
	reader.Comma = sep
	reader.LazyQuotes = lazyQuotes

	records, err := reader.ReadAll()
	if err != nil {
		return table.Table{}, err
	}
	if len(records) == 0 {
		return table.Table{}, fmt.Errorf("empty file")
	}

	header := records[0]
	records = records[1:]

	fieldMap := make(map[string]int)
	for i, fieldName := range header {
		fieldMap[fieldName] = i
	}

	cols := make([]table.Column, 0, len(fields))

	for _, field := range fields {
		columnIndex, found := fieldMap[field.Name]
		if !found {
			return table.Table{}, fmt.Errorf("column '%s' not found in file", field.Name)
		}
		col := make(table.Column, 0, len(records))
		for rowIndex := 0; rowIndex < len(records); rowIndex++ {
			switch field.Type {
			case table.FieldTypeNumber:
				val, err := table.NewNumberValue(records[rowIndex][columnIndex])
				if err != nil {
					return table.Table{}, fmt.Errorf("error when parsing column %s, line %d: %w", field.Name, rowIndex+2, err)
				}
				col = append(col, val)
			case table.FieldTypeString:
				col = append(col, table.NewStringValue(records[rowIndex][columnIndex]))
			default:
				return table.Table{}, fmt.Errorf("unknown field type for %s", field.Name)
			}
		}
		cols = append(cols, col)
	}

	return table.NewTable(fields, cols), nil
}
