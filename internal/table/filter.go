package table

import (
	"context"
)

type CompareOperationType string

const (
	CompareOperationTypeEqual CompareOperationType = "="
	CompareOperationTypeLess  CompareOperationType = "<"
	CompareOperationTypeMore  CompareOperationType = ">"
)

type CompareFilter struct {
	FieldName string
	Op        CompareOperationType
	Val       interface{}
}

type LogicalOperationType string

const (
	LogicalOperationTypeAnd   LogicalOperationType = "and"
	LogicalOperationTypeOr    LogicalOperationType = "or"
	LogicalOperationTypeDummy LogicalOperationType = "dummy"
)

type LogicalFilter struct {
	Op    LogicalOperationType
	Left  *LogicalFilter
	Right *LogicalFilter
}

type LogicalFilterV2 interface {
	Filtrate(ctx context.Context, t Table) ([]int, error)
}

type AndFilter struct {
	Left  LogicalFilterV2
	Right LogicalFilterV2
}

func (f AndFilter) Filtrate(ctx context.Context, t Table) ([]int, error) {
	ret, err := f.Left.Filtrate(ctx, t)
	if err != nil {
		return nil, err
	}
	if len(ret) == 0 {
		return nil, nil
	}

	subTable := t.GetSubTable(ret)
	return f.Right.Filtrate(ctx, subTable)
}

type OrFilter struct {
	Left  LogicalFilterV2
	Right LogicalFilterV2
}

func (f OrFilter) Filtrate(ctx context.Context, t Table) ([]int, error) {
	res1, err := f.Left.Filtrate(ctx, t)
	if err != nil {
		return nil, err
	}

	res2, err := f.Right.Filtrate(ctx, t)
	if err != nil {
		return nil, err
	}

	if len(res1) == 0 {
		return res2, nil
	}
	if len(res2) == 0 {
		return res1, nil
	}

	return mergeIndexes(res1, res2), nil
}

type DummyFilter struct {
	filter CompareFilter
}

func (f DummyFilter) Filtrate(ctx context.Context, t Table) ([]int, error) {
	fieldIndex, err := t.GetFieldIndex(f.filter.FieldName)
	if err != nil {
		return nil, err
	}

	column := t.Columns[fieldIndex]

	var ret []int
	for i, val := range column {
		accept, err := val.Compare(f.filter.Val, f.filter.Op)
		if err != nil {
			return nil, err
		}
		if accept {
			ret = append(ret, i)
		}
	}
	return ret, nil
}

// mergeIndexes объединяет два слайса, сохраняя их порядок и исключая дубликаты
func mergeIndexes(first []int, second []int) []int {
	var ret []int

	if len(first) > len(second) {
		ret = make([]int, 0, len(first))
	} else {
		ret = make([]int, 0, len(second))
	}

	i, j := 0, 0
	lastValue := -1

	for {
		// Проверяем, не выходим ли мы за пределы слайсов
		if i == len(first) || j == len(second) {
			break
		}
		// Игнорируем дубликаты
		if first[i] == lastValue {
			i++
			continue
		}
		if second[j] == lastValue {
			j++
			continue
		}
		if first[i] < second[j] {
			ret = append(ret, first[i])
			lastValue = first[i]
			i++
			continue
		}
		ret = append(ret, second[j])
		lastValue = second[j]
		j++
	}

	// копируем оставшиеся элементы из первого слайса при их наличии
	if i < len(first) {
		for ; i < len(first); i++ {
			if lastValue != first[i] {
				ret = append(ret, first[i])
			}
		}
	}

	// копируем оставшиеся элементы из второго слайса при их наличии
	if j < len(second) {
		for ; j < len(second); j++ {
			if lastValue != second[j] {
				ret = append(ret, second[j])
			}
		}
	}

	return ret
}
