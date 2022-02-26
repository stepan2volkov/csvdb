package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/stepan2volkov/csvdb/internal/config"
	"github.com/stepan2volkov/csvdb/internal/parser"
	"github.com/stepan2volkov/csvdb/internal/table"
	"github.com/stepan2volkov/csvdb/internal/tableloader"
)

func NewTableHolder() *TableHolder {
	return &TableHolder{
		tables: make(map[string]table.Table),
	}
}

type TableHolder struct {
	tables map[string]table.Table
}

func (h *TableHolder) Load(csvPath string, configPath string) error {
	tc, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}

	if _, found := h.tables[tc.Name]; found {
		return fmt.Errorf("table %s already exists", tc.Name)
	}

	fields, err := tc.GetFields()
	if err != nil {
		return err
	}

	t, err := tableloader.Load(csvPath, tc.GetSep(), tc.LazyQuotes, fields)
	if err != nil {
		return err
	}

	h.tables[tc.Name] = t

	return nil
}

func (h *TableHolder) Delete(tableName string) error {
	if _, found := h.tables[tableName]; !found {
		return fmt.Errorf("table '%s' doesn't exist", tableName)
	}
	delete(h.tables, tableName)
	return nil
}

func (h *TableHolder) GetList() string {
	if len(h.tables) == 0 {
		return "\ttables haven't loaded yet"
	}
	tables := make([]string, 0, len(h.tables))

	for k := range h.tables {
		tables = append(tables, k)
	}

	return strings.Join(tables, "\n")
}

func (h *TableHolder) Exists(tableName string) bool {
	_, found := h.tables[tableName]
	return found
}

func (h *TableHolder) GetTable(tableName string) (table.Table, error) {
	t, found := h.tables[tableName]
	if !found {
		return table.Table{}, fmt.Errorf("table '%s' doesn't exist", tableName)
	}

	return t, nil
}

func (h *TableHolder) Query(tableName string, query string) (string, error) {
	t, err := h.GetTable(tableName)
	if err != nil {
		return "", err
	}
	p := parser.NewParser()
	tokens, err := p.Parse(strings.NewReader(query))
	if err != nil {
		return "", err
	}

	filter, err := parser.MakeFilters(tokens)
	if err != nil {
		return "", err
	}
	indexes, err := filter.Filtrate(context.Background(), t)
	if err != nil {
		return "", err
	}
	t2 := t.GetSubTable(indexes)

	return t2.String(), nil
}

func main() {
	var tableName string
	th := NewTableHolder()
	// th.Load("./grades.csv", "grades.yaml")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("%s~# ", tableName)

		openedBuffer := scanner.Scan()
		if !openedBuffer {
			return
		}
		in := scanner.Text()
		if in == "" {
			continue
		}
		switch {
		case in == `\q`:
			return
		case strings.HasPrefix(in, `\load`):
			in = strings.TrimSpace(strings.TrimPrefix(in, `\load`))
			args := strings.Split(in, " ")
			if len(args) != 2 {
				fmt.Println("\tThe right syntax: \n\t\\load <csv-path> <config-path>")
				continue
			}
			if err := th.Load(args[0], args[1]); err != nil {
				fmt.Printf("ERROR: %v\n", err)
			}
		case in == `\list`:
			fmt.Println(th.GetList())
		case in == `\help`:
		case strings.HasPrefix(in, `\drop`):
			newTableName := strings.TrimPrefix(in, `\drop`)
			newTableName = strings.TrimSpace(newTableName)
			if newTableName == "" {
				fmt.Println("\tThe right syntax: \n\t\\drop <tablename>")
				continue
			}
			if tableName == newTableName {
				tableName = ""
			}
			if err := th.Delete(newTableName); err != nil {
				fmt.Printf("\t%v\n", err)
			}
		case strings.HasPrefix(in, `\use`):
			newTableName := strings.TrimPrefix(in, `\use`)
			newTableName = strings.TrimSpace(newTableName)
			if newTableName == "" {
				fmt.Println("\tThe right syntax: \n\t\\use <tablename>")
				continue
			}
			if th.Exists(newTableName) {
				tableName = newTableName
			} else {
				fmt.Printf("\ttable '%s' doesn't exist\n", newTableName)
			}
		default:
			if tableName == "" {
				fmt.Println("\tyou should select table:\n\t\\use <tablename>")
			}
			res, err := th.Query(tableName, in)
			if err != nil {
				fmt.Printf("ERROR: %v\n", err)
				continue
			}
			fmt.Println(res)
		}
	}
}
