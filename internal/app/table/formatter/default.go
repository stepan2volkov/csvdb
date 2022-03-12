package formatter

import (
	prettyTable "github.com/jedib0t/go-pretty/v6/table"

	"github.com/stepan2volkov/csvdb/internal/app/table"
)

var _ table.Formatter = &DefaultFormatter{}

type DefaultFormatter struct{}

func (f *DefaultFormatter) Format(t table.Table) string {
	writer := prettyTable.NewWriter()
	writer.SetStyle(prettyTable.StyleLight)

	header := prettyTable.Row{}
	for _, col := range t.Columns {
		header = append(header, col.Field.Name)
	}
	writer.AppendHeader(header)

	rowCount := len(t.Columns[0].Values)
	rows := make([]prettyTable.Row, 0, rowCount)

	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {
		row := prettyTable.Row{}

		for columnIndex := range t.Columns {
			row = append(row, t.Columns[columnIndex].Values[rowIndex])
		}
		rows = append(rows, row)
	}

	writer.AppendRows(rows)

	return writer.Render()
}
