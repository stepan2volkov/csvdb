package config

import (
	"fmt"
	"os"

	"github.com/stepan2volkov/csvdb/internal/table"
	"gopkg.in/yaml.v3"
)

type Field struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type TableConfig struct {
	Name       string  `yaml:"name"`
	Sep        string  `yaml:"sep" default:";"`
	LazyQuotes bool    `yaml:"lazyQuotes" default:"false"`
	Fields     []Field `yaml:"fields"`
}

func (c TableConfig) GetSep() rune {
	for _, r := range c.Sep {
		return r
	}
	return ';'
}

func (c TableConfig) GetFields() ([]table.Field, error) {
	ret := make([]table.Field, 0, len(c.Fields))

	for _, f := range c.Fields {
		var t table.FieldType
		if f.Type == "number" {
			t = table.FieldTypeNumber
		} else if f.Type == "string" {
			t = table.FieldTypeString
		} else {
			return nil, fmt.Errorf("unknown type '%s'", f.Type)
		}

		ret = append(ret, table.Field{
			Name: f.Name,
			Type: t,
		})
	}

	return ret, nil
}

func LoadConfig(path string) (TableConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return TableConfig{}, err
	}
	tc := TableConfig{}
	if err = yaml.NewDecoder(file).Decode(&tc); err != nil {
		return TableConfig{}, err
	}
	if tc.Name == "" {
		return TableConfig{}, fmt.Errorf("name cannot be empty")
	}
	if len(tc.Sep) != 1 {
		return TableConfig{}, fmt.Errorf("sep should be presented by only one character")
	}
	return tc, nil
}
