package table

import (
	"fmt"
)

var _ Value = StringValue{}

type StringValue struct {
	value string
}

func (v StringValue) String() string {
	return v.value
}

func (v StringValue) Compare(val interface{}, op CompareOperationType) (bool, error) {
	compareValue, valid := val.(string)
	if !valid {
		return false, fmt.Errorf("invalid value for string: '%v'", val)
	}

	switch op {
	case CompareOperationTypeLess, CompareOperationTypeMore:
		return false, fmt.Errorf("invalid operation for type string: %s", op)
	case CompareOperationTypeEqual:
		return v.value == compareValue, nil
	}

	return false, fmt.Errorf("unknown operation for type string: %s", op)
}

var _ Value = NumberValue{}

type NumberValue struct {
	value float64
}

func (v NumberValue) String() string {
	return fmt.Sprint(v.value)
}

func (v NumberValue) Compare(val interface{}, op CompareOperationType) (bool, error) {
	compareValue, valid := val.(float64)
	if !valid {
		return false, fmt.Errorf("invalid value for number: '%v'", val)
	}

	switch op {
	case CompareOperationTypeLess:
		return v.value < compareValue, nil
	case CompareOperationTypeMore:
		return v.value > compareValue, nil
	case CompareOperationTypeEqual:
		return v.value == compareValue, nil
	}

	return false, fmt.Errorf("unknown operation for type number: %s", op)
}
