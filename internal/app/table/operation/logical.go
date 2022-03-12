package operation

import (
	"context"
	"sort"

	"github.com/stepan2volkov/csvdb/internal/app/table"
)

var _ table.LogicalOperation = AndOperation{}
var _ table.LogicalOperation = OrOperation{}
var _ table.LogicalOperation = DummyValueOperation{}

type AndOperation struct {
	Left  table.LogicalOperation
	Right table.LogicalOperation
}

func (o AndOperation) Apply(ctx context.Context, t table.Table) ([]int, error) {
	res1, err := o.Left.Apply(ctx, t)
	if err != nil {
		return nil, err
	}
	if len(res1) == 0 {
		return nil, nil
	}
	res2, err := o.Right.Apply(ctx, t)
	if err != nil {
		return nil, err
	}
	if len(res2) == 0 {
		return nil, nil
	}
	firstMap := map[int]struct{}{}
	for _, i := range res1 {
		firstMap[i] = struct{}{}
	}

	var ret []int
	for _, i := range res2 {
		if _, found := firstMap[i]; found {
			ret = append(ret, i)
		}
	}
	sort.Ints(ret)

	return ret, nil
}

type OrOperation struct {
	Left  table.LogicalOperation
	Right table.LogicalOperation
}

func (o OrOperation) Apply(ctx context.Context, t table.Table) ([]int, error) {
	res1, err := o.Left.Apply(ctx, t)
	if err != nil {
		return nil, err
	}

	res2, err := o.Right.Apply(ctx, t)
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

type DummyValueOperation struct {
	CompareOperation table.CompareValueOperation
}

func (o DummyValueOperation) Apply(ctx context.Context, t table.Table) ([]int, error) {
	var ret []int

	if o.CompareOperation.Type == table.CompareOperationTypeDummy {
		ret = make([]int, 0, len(t.Columns[0].Values))
		for i := range t.Columns[0].Values {
			ret = append(ret, i)
		}
		return ret, nil
	}

	column, err := t.GetColumnByName(o.CompareOperation.ColumnName)
	if err != nil {
		return nil, err
	}

	for i, val := range column.Values {
		accept, err := val.Compare(o.CompareOperation.Val, o.CompareOperation.Type)
		if err != nil {
			return nil, err
		}
		if accept {
			ret = append(ret, i)
		}
	}
	return ret, nil
}
