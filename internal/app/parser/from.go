package parser

import (
	"fmt"

	"github.com/stepan2volkov/csvdb/internal/app/scanner"
	"github.com/stepan2volkov/csvdb/internal/app/table"
	"github.com/stepan2volkov/csvdb/internal/app/table/operation"
)

func makeWhere(tokens []scanner.Token) (table.LogicalOperation, error) {
	var ret []table.LogicalOperation
	tokenQueue := make([]scanner.Token, 0, len(tokens))

	for _, token := range tokens {
		switch token.Type() {
		case scanner.TokenTypeOpMore, scanner.TokenTypeOpLess, scanner.TokenTypeOpEqual:

			val, id := tokenQueue[len(tokenQueue)-1], tokenQueue[len(tokenQueue)-2]
			tokenQueue = tokenQueue[:len(tokenQueue)-2]
			if id.Type() != scanner.TokenTypeID {
				val, id = id, val
			}
			if id.Type() != scanner.TokenTypeID ||
				(val.Type() != scanner.TokenTypeString && val.Type() != scanner.TokenTypeNumber) {
				return nil, fmt.Errorf("invalid where format")
			}
			var op table.CompareOperationType

			switch token.Type() {
			case scanner.TokenTypeOpMore:
				op = table.CompareOperationTypeMore
			case scanner.TokenTypeOpLess:
				op = table.CompareOperationTypeLess
			case scanner.TokenTypeOpEqual:
				op = table.CompareOperationTypeEqual
			}

			ret = append(ret, operation.DummyValueOperation{
				CompareOperation: table.CompareValueOperation{
					ColumnName: id.Value().(string),
					Type:       op,
					Val:        val.Value(),
				},
			})
		case scanner.TokenTypeOpAnd:
			if len(ret) < 2 {
				return nil, fmt.Errorf("invalid where format")
			}
			arg1, arg2 := ret[len(ret)-1], ret[len(ret)-2]
			ret = ret[:len(ret)-2]
			ret = append(ret, operation.AndOperation{
				Left:  arg1,
				Right: arg2,
			})
		case scanner.TokenTypeOpOr:
			if len(ret) < 2 {
				return nil, fmt.Errorf("invalid where format")
			}
			arg1, arg2 := ret[len(ret)-1], ret[len(ret)-2]
			ret = ret[:len(ret)-2]
			ret = append(ret, operation.OrOperation{
				Left:  arg1,
				Right: arg2,
			})
		default:
			tokenQueue = append(tokenQueue, token)
		}
	}
	if len(ret) != 1 {
		return nil, fmt.Errorf("invalid where format")
	}
	return ret[0], nil
}
