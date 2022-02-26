package table

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

type FieldType int

const (
	FieldTypeNumber FieldType = iota
	FieldTypeString FieldType = iota
)

type Field struct {
	Name string
	Type FieldType
}

type Value interface {
	String() string
	Compare(val interface{}, op CompareOperationType) (bool, error)
}

type Column []Value

func NewTable(fields []Field, cols []Column) Table {
	fieldIndexes := make(map[string]int)
	for i, field := range fields {
		fieldIndexes[field.Name] = i
	}

	return Table{
		fieldIndexes: fieldIndexes,
		Fields:       fields,
		Columns:      cols,
	}
}

type Table struct {
	fieldIndexes map[string]int
	Fields       []Field
	Columns      []Column
}

func (t Table) GetSubTable(rowIndexes []int) Table {
	fields := make([]Field, len(t.Fields))
	copy(fields, t.Fields)

	cols := make([]Column, 0, len(t.Columns))

	for columnIndex := range t.Fields {
		col := make(Column, 0, len(rowIndexes))
		for _, rowIndex := range rowIndexes {
			col = append(col, t.Columns[columnIndex][rowIndex])
		}
		cols = append(cols, col)
	}

	return NewTable(fields, cols)
}

func (t Table) GetField(fieldName string) (Field, error) {
	i, found := t.fieldIndexes[fieldName]
	if !found {
		return Field{}, fmt.Errorf("field '%s' not found", fieldName)
	}
	return t.Fields[i], nil
}

func (t Table) GetFieldIndex(fieldName string) (int, error) {
	i, found := t.fieldIndexes[fieldName]
	if !found {
		return 0, fmt.Errorf("field '%s' not found", fieldName)
	}
	return i, nil
}

func (t Table) String() string {
	writer := table.NewWriter()
	writer.SetStyle(table.StyleLight)

	header := table.Row{}
	for _, field := range t.Fields {
		header = append(header, field.Name)
	}
	writer.AppendHeader(header)

	rowCount := len(t.Columns[0])
	rows := make([]table.Row, 0, rowCount)

	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
		row := table.Row{}

		for columnIndex := range t.Fields {
			row = append(row, t.Columns[columnIndex][rowIndex])
		}
		rows = append(rows, row)
	}

	writer.AppendRows(rows)

	return writer.Render()
}
