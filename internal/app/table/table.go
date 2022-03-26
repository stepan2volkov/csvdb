package table

import (
	"context"
	"fmt"
)

type CompareOperationType string

const (
	CompareOperationTypeEqual CompareOperationType = "="
	CompareOperationTypeLess  CompareOperationType = "<"
	CompareOperationTypeMore  CompareOperationType = ">"
	CompareOperationTypeDummy CompareOperationType = "dummy"
)

type Formatter interface {
	Format(ctx context.Context, t Table) (string, error)
}

type CompareValueOperation struct {
	ColumnName string
	Type       CompareOperationType
	Val        interface{}
}

type FieldType int

type LogicalOperation interface {
	Apply(ctx context.Context, t Table) ([]int, error)
}

const (
	FieldTypeNumber FieldType = iota
	FieldTypeString FieldType = iota
)

type Field struct {
	Name string
	Type FieldType
}

type Value interface {
	fmt.Stringer
	Compare(val interface{}, op CompareOperationType) (bool, error)
}

type Column struct {
	Field  Field
	Values []Value
}

func NewTable(name string, cols []Column) Table {
	columnIndexes := make(map[string]int)
	for i, col := range cols {
		columnIndexes[col.Field.Name] = i
	}

	return Table{
		Name:          name,
		columnIndexes: columnIndexes,
		Columns:       cols,
	}
}

type Table struct {
	Name          string
	columnIndexes map[string]int
	Columns       []Column
}

func (t Table) GetColumnByName(name string) (Column, error) {
	i, found := t.columnIndexes[name]
	if !found {
		return Column{}, fmt.Errorf("column '%s' not found", name)
	}

	return t.Columns[i], nil
}

func (t Table) GetSubTableByIndexes(ctx context.Context, rowIndexes []int) (Table, error) {
	cols := make([]Column, len(t.Columns))

	for i, col := range t.Columns {
		cols[i].Field = col.Field
		cols[i].Values = make([]Value, 0, len(rowIndexes))
		for _, index := range rowIndexes {
			select {
			case <-ctx.Done():
				return Table{}, ctx.Err()
			default:
				cols[i].Values = append(cols[i].Values, col.Values[index])
			}

		}
	}

	return NewTable(t.Name, cols), nil
}

func (t Table) GetSubTableByFields(fields []string) (Table, error) {
	fieldMap := make(map[string]struct{})
	for _, f := range fields {
		fieldMap[f] = struct{}{}
	}

	cols := make([]Column, 0, len(fieldMap))

	for i, f := range t.Columns {
		if _, found := fieldMap[f.Field.Name]; found {
			cols = append(cols, t.Columns[i])
			delete(fieldMap, f.Field.Name)
		}
	}
	if len(fieldMap) != 0 {
		notFoundField := ""
		for k := range fieldMap {
			notFoundField = k

			break
		}

		return Table{}, fmt.Errorf("fields %s not found", notFoundField)
	}

	return NewTable(t.Name, cols), nil
}
