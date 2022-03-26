package value

import (
	"fmt"

	"github.com/stepan2volkov/csvdb/internal/app/table"
)

var _ table.Value = StringValue{}

func NewStringValue(val string) StringValue {
	return StringValue{value: val}
}

type StringValue struct {
	value string
}

func (v StringValue) String() string {
	return v.value
}

func (v StringValue) Compare(val interface{}, op table.CompareOperationType) (bool, error) {
	compareValue, valid := val.(string)
	if !valid {
		return false, fmt.Errorf("invalid value for string: '%v'", val)
	}

	if op == table.CompareOperationTypeEqual {
		return v.value == compareValue, nil
	}

	return false, fmt.Errorf("invalid operation for type string: %s", op)
}
