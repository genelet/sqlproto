package ast

import (
	"fmt"
	"github.com/genelet/sqlproto/xast"
	"github.com/akito0107/xsqlparser/sqlast"
)

func xnativevalueTo(node sqlast.Node) (*xast.NativeValue, error) {
	item := &xast.NativeValue{}
	switch value := node.(type) {
	case *sqlast.LongValue:
		item.NativeOneOf = &xast.NativeValue_LongValue{LongValue:xlongTo(value)}
	case *sqlast.SingleQuotedString:
		item.NativeOneOf = &xast.NativeValue_SingleQuotedString{SingleQuotedString:xstringTo(value)}
	case *sqlast.DoubleValue:
		item.NativeOneOf = &xast.NativeValue_DoubleValue{DoubleValue:xdoubleTo(value)}
	default:
		return nil, fmt.Errorf("native value type %#v", value)
	}

	return item, nil
}

func xrowvalueexprTo(row *sqlast.RowValueExpr) (*xast.ConstructorSource_RowValueExpr, error) {
	if row == nil { return nil, nil }

	item := &xast.ConstructorSource_RowValueExpr{
		LParen: xposTo(row.LParen),
		RParen: xposTo(row.RParen)};

	for _, nativevalue := range row.Values { 
		value, err := xnativevalueTo(nativevalue)
		if err != nil { return nil, err }
		item.Values = append(item.Values, value)
	}

	return item, nil
}

func xconstructorTo(constructor *sqlast.ConstructorSource) (*xast.ConstructorSource, error) {
	if constructor == nil { return nil, nil }

	item := &xast.ConstructorSource {Values: xposTo(constructor.Values)}
	for _, row := range constructor.Rows {
		rowvalue, err := xrowvalueexprTo(row)	
		if err != nil { return nil, err }
		item.Rows = append(item.Rows, rowvalue)
	}

	return item, nil
}

func constructorTo(constructor *xast.ConstructorSource) *sqlast.ConstructorSource {
	if constructor == nil { return nil }

	item := &sqlast.ConstructorSource {Values: posTo(constructor.Values)}
	for _, row := range constructor.Rows {	
		item.Rows = append(item.Rows, rowvalueexprTo(row))
	}
	return item
}

func rowvalueexprTo(row *xast.ConstructorSource_RowValueExpr) *sqlast.RowValueExpr {
	if row == nil { return nil }

	item := &sqlast.RowValueExpr{
		LParen: posTo(row.LParen),
		RParen: posTo(row.RParen)};

	for _, nativevalue := range row.Values { 
		item.Values = append(item.Values, nativevalueTo(nativevalue))
	}

	return item
}

func nativevalueTo(node *xast.NativeValue) sqlast.Node {
	if v := node.GetSingleQuotedString(); v != nil {
		return stringTo(v)
	} else if v := node.GetDoubleValue(); v != nil {
		return doubleTo(v)
	} else if v := node.GetLongValue(); v != nil {
		return longTo(v)
	}

	return nil
}
