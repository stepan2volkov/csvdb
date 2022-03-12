package value

import (
	"fmt"
	"strconv"

	"github.com/stepan2volkov/csvdb/internal/app/table"
)

var _ table.Value = NumberValue{}

func NewNumberValue(val string) (NumberValue, error) {
	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return NumberValue{}, err
	}
	return NumberValue{value: num}, nil
}

type NumberValue struct {
	value float64
}

func (v NumberValue) String() string {
	return fmt.Sprint(v.value)
}

func (v NumberValue) Compare(val interface{}, op table.CompareOperationType) (bool, error) {
	compareValue, valid := val.(float64)
	if !valid {
		return false, fmt.Errorf("invalid value for number: '%v'", val)
	}

	switch op {
	case table.CompareOperationTypeLess:
		return v.value < compareValue, nil
	case table.CompareOperationTypeMore:
		return v.value > compareValue, nil
	case table.CompareOperationTypeEqual:
		return v.value == compareValue, nil
	}

	return false, fmt.Errorf("unknown operation for type number: %s", op)
}
