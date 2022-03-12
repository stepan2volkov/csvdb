package loader

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/stepan2volkov/csvdb/internal/app/table"
	"github.com/stepan2volkov/csvdb/internal/app/table/value"
)

func LoadFromCSV(csvPath string, configPath string) (table.Table, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return table.Table{}, err
	}
	tableConfig, err := loadConfig(file)
	if err != nil {
		return table.Table{}, err
	}

	fields, err := tableConfig.getFields()
	if err != nil {
		return table.Table{}, err
	}

	t, err := load(tableConfig.Name, csvPath, tableConfig.getSep(), tableConfig.LazyQuotes, fields)
	if err != nil {
		return table.Table{}, err
	}

	return t, nil
}

func load(tableName string, path string, sep rune, lazyQuotes bool, fields []table.Field) (table.Table, error) {
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

		col := table.Column{Field: field}

		for rowIndex := 0; rowIndex < len(records); rowIndex++ {
			switch field.Type {
			case table.FieldTypeNumber:
				val, err := value.NewNumberValue(records[rowIndex][columnIndex])
				if err != nil {
					return table.Table{}, fmt.Errorf("error when parsing column %s, line %d: %w", field.Name, rowIndex+2, err)
				}
				col.Values = append(col.Values, val)
			case table.FieldTypeString:
				val := value.NewStringValue(records[rowIndex][columnIndex])
				col.Values = append(col.Values, val)
			default:
				return table.Table{}, fmt.Errorf("unknown field type for %s", field.Name)
			}
		}
		cols = append(cols, col)
	}

	return table.NewTable(tableName, cols), nil
}
