package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/stepan2volkov/csvdb/internal/app/parser"
	"github.com/stepan2volkov/csvdb/internal/app/scanner"
	"github.com/stepan2volkov/csvdb/internal/app/table"
	"go.uber.org/zap"
)

func NewApp(logger *zap.Logger) *App {
	return &App{
		logger: logger,
		tables: make(map[string]table.Table),
	}
}

type App struct {
	tables map[string]table.Table
	logger *zap.Logger
}

func (a *App) LoadTable(t table.Table) error {
	if strings.TrimSpace(t.Name) == "" {
		return fmt.Errorf("table name should'n be empty")
	}
	if _, found := a.tables[t.Name]; found {
		return fmt.Errorf("table '%s' has already exist", t.Name)
	}
	a.tables[t.Name] = t
	return nil
}

func (a *App) TableList() []string {
	ret := make([]string, 0, len(a.tables))
	for tableName := range a.tables {
		ret = append(ret, tableName)
	}
	return ret
}

func (a *App) DropTable(tableName string) error {
	if _, found := a.tables[tableName]; !found {
		return fmt.Errorf("table '%s' doesn't exist", tableName)
	}
	delete(a.tables, tableName)
	return nil
}

func (a *App) Execute(ctx context.Context, query string) (table.Table, error) {
	stmtScanner := scanner.NewScanner()
	tokens, err := stmtScanner.Scan(strings.NewReader(query))
	if err != nil {
		return table.Table{}, err
	}

	stmt, err := parser.MakeSelectStmt(tokens)
	if err != nil {
		return table.Table{}, err
	}

	t, found := a.tables[stmt.Tablename]
	if !found {
		return table.Table{}, fmt.Errorf("table '%s' doesn't exist", stmt.Tablename)
	}

	indexes, err := stmt.Filter.Apply(ctx, t)
	if err != nil {
		return table.Table{}, err
	}

	ret := t.GetSubTableByIndexes(indexes)
	if stmt.AllField {
		return ret, nil
	}

	ret, err = t.GetSubTableByFields(stmt.Fields)
	if err != nil {
		return table.Table{}, err
	}

	return ret, nil
}
