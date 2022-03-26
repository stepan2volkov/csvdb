package loader

import (
	"fmt"
	"io"

	"github.com/stepan2volkov/csvdb/internal/app/table"

	"gopkg.in/yaml.v3"
)

type field struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type tableConfig struct {
	Name       string  `yaml:"name"`
	Sep        string  `yaml:"sep" default:";"`
	LazyQuotes bool    `yaml:"lazyQuotes" default:"false"`
	Fields     []field `yaml:"fields"`
}

func (c tableConfig) getSep() rune {
	for _, r := range c.Sep {
		return r
	}

	return ';'
}

func (c tableConfig) getFields() ([]table.Field, error) {
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

func loadConfig(file io.Reader) (tableConfig, error) {
	tc := tableConfig{}
	if err := yaml.NewDecoder(file).Decode(&tc); err != nil {
		return tableConfig{}, fmt.Errorf("error when decode config: %w", err)
	}

	if tc.Name == "" {
		return tableConfig{}, fmt.Errorf("name cannot be empty")
	}
	if len(tc.Sep) != 1 {
		return tableConfig{}, fmt.Errorf("sep should be presented by only one character")
	}

	return tc, nil
}
