package table

import "fmt"

type FieldType int

const (
	FieldTypeInt    FieldType = iota
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
	fields := make([]Field, 0, len(t.Fields))
	copy(fields, t.Fields)

	rows := make([]Column, 0, len(rowIndexes))
	for _, i := range rowIndexes {
		rows = append(rows, t.Columns[i])
	}

	return Table{Fields: fields, Columns: rows}
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
