package parser

import (
	"fmt"

	"github.com/stepan2volkov/csvdb/internal/table"
)

func MakeFilters(tokens []Token) (table.LogicalFilter, error) {
	var ret []table.LogicalFilter
	tokenQueue := make([]Token, 0, len(tokens))

	for _, token := range tokens {
		switch token.Type() {
		case TokenTypeOpMore, TokenTypeOpLess, TokenTypeOpEqual:

			val, id := tokenQueue[len(tokenQueue)-1], tokenQueue[len(tokenQueue)-2]
			tokenQueue = tokenQueue[:len(tokenQueue)-2]
			if id.Type() != TokenTypeID {
				val, id = id, val
			}
			if id.Type() != TokenTypeID || (val.Type() != TokenTypeString && val.Type() != TokenTypeNumber) {
				return nil, fmt.Errorf("invalid where format")
			}
			var op table.CompareOperationType

			switch token.Type() {
			case TokenTypeOpMore:
				op = table.CompareOperationTypeMore
			case TokenTypeOpLess:
				op = table.CompareOperationTypeLess
			case TokenTypeOpEqual:
				op = table.CompareOperationTypeEqual
			}

			ret = append(ret, table.DummyFilter{
				Filter: table.CompareFilter{
					FieldName: id.value.(string),
					Op:        op,
					Val:       val.value,
				},
			})
		case TokenTypeOpAnd:
			if len(ret) < 2 {
				return nil, fmt.Errorf("invalid where format")
			}
			arg1, arg2 := ret[len(ret)-1], ret[len(ret)-2]
			ret = ret[:len(ret)-2]
			ret = append(ret, table.AndFilter{
				Left:  arg1,
				Right: arg2,
			})
		case TokenTypeOpOr:
			if len(ret) < 2 {
				return nil, fmt.Errorf("invalid where format")
			}
			arg1, arg2 := ret[len(ret)-1], ret[len(ret)-2]
			ret = ret[:len(ret)-2]
			ret = append(ret, table.OrFilter{
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
