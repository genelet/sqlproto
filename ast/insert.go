package ast

import (
	"fmt"
	"github.com/genelet/sqlproto/xast"
	"github.com/akito0107/xsqlparser/sqlast"
)

// XInsertTo translates a xsqlparser insert statement into xast statement
//
// see https://github.com/akito0107/xsqlparser for xsqlparser
//
func XInsertTo(stmt *sqlast.InsertStmt) (*xast.InsertStmt, error) {
	output := &xast.InsertStmt{
                Insert: xposTo(stmt.Insert),
                TableName: xobjectnameTo(stmt.TableName)}

	for _, column := range stmt.Columns {
		v := xidentTo(column)
                output.Columns = append(output.Columns, v)
	}

	switch source := stmt.Source.(type) {
	case *sqlast.SubQuerySource:
		query, err := XQueryTo(source.SubQuery)
		if err != nil { return nil, err }
		output.InsertSource = &xast.InsertStmt_QuerySource{QuerySource: query}
        case *sqlast.ConstructorSource:
		constructor, err := xconstructorTo(source)
		if err != nil { return nil, err }
		output.InsertSource = &xast.InsertStmt_Constructor{Constructor: constructor}
	default:
		return nil, fmt.Errorf("unknown insert data %T", source)
	}

	for _, item := range stmt.UpdateAssignments {
		v, err := xassignmentTo(item)
		if err != nil { return nil, err }
		output.UpdateAssignments = append(output.UpdateAssignments, v)
	}

	return output, nil
}

// InsertTo translates a xast insert statement into xsqlparser statement
//
func InsertTo(stmt *xast.InsertStmt) *sqlast.InsertStmt {
	output := &sqlast.InsertStmt{
		TableName: compoundToObjectname(stmt.TableName)}
	if stmt.Insert != nil {
		output.Insert = posTo(stmt.Insert)
	}

	if v := stmt.GetQuerySource(); v != nil {
		output.Source = &sqlast.SubQuerySource{SubQuery: QueryTo(v)}
        } else if v := stmt.GetConstructor(); v != nil {
		output.Source = constructorTo(v)
        }

	for _, item := range stmt.Columns {
		output.Columns = append(output.Columns, identTo(item).(*sqlast.Ident))
	}
	for _, item := range stmt.UpdateAssignments {
		output.UpdateAssignments = append(output.UpdateAssignments, assignmentTo(item))
	}

        return output

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
