package ast

import (
	"fmt"
	"github.com/genelet/sqlproto/xast"
	"github.com/akito0107/xsqlparser/sqlast"
)

// XUpdateTo translates a xsqlparser update statement into xast statement
//
// see https://github.com/akito0107/xsqlparser for xsqlparser
//
func XUpdateTo(stmt *sqlast.UpdateStmt) (*xast.UpdateStmt, error) {
	output := &xast.UpdateStmt{
		Update: xposTo(stmt.Update),
		TableName: xobjectnameTo(stmt.TableName)}

	for _, item := range stmt.Assignments {
		v, err := xassignmentTo(item)
		if err != nil { return nil, err }
		output.Assignments = append(output.Assignments, v)
	}

	err := xupdatewhereTo(stmt, output)
	if err != nil { return nil, err }

	return output, nil
}

// UpdateTo translates a xast update statement into xsqlparser statement
//
func UpdateTo(stmt *xast.UpdateStmt) *sqlast.UpdateStmt {
	output := &sqlast.UpdateStmt{
		TableName: compoundToObjectname(stmt.TableName)}
	if stmt.Update != nil {
		output.Update = posTo(stmt.Update)
	}

	for _, item := range stmt.Assignments {
		output.Assignments = append(output.Assignments, assignmentTo(item))
	}

	updatewhereTo(stmt, output)

	return output
}

func xassignmentTo(item *sqlast.Assignment) (*xast.Assignment, error) {
	output := &xast.Assignment{
		AssignmentID: xidentsTo(item.ID)}

	switch value := item.Value.(type) {
	case *sqlast.Ident:
		output.ValueOneOf = &xast.Assignment_RightIdents{RightIdents:xidentsTo(value)}
	case *sqlast.CompoundIdent:
		output.ValueOneOf = &xast.Assignment_RightIdents{RightIdents:xcompoundTo(value)}
	case *sqlast.SubQuery:
		insub, err := xsubqueryTo(value)
		if err != nil { return nil, err }
		output.ValueOneOf = &xast.Assignment_QueryValue{QueryValue: insub}
	case *sqlast.BinaryExpr:
		middle, err := xbinaryexprTo(value)
		if err != nil { return nil, err }
		output.ValueOneOf = &xast.Assignment_RightBinary{RightBinary:middle}
	case *sqlast.LongValue:
		output.ValueOneOf = &xast.Assignment_LongValue{LongValue:xlongTo(value)}
	case *sqlast.SingleQuotedString:
		output.ValueOneOf = &xast.Assignment_SingleQuotedString{SingleQuotedString:xstringTo(value)}
	case *sqlast.DoubleValue:
		output.ValueOneOf = &xast.Assignment_DoubleValue{DoubleValue:xdoubleTo(value)}
	default:
		return nil, fmt.Errorf("type %T not found", value)
        }

	return output, nil
}

func assignmentTo(item *xast.Assignment) *sqlast.Assignment {
	output := &sqlast.Assignment{
		ID: compoundTo(item.AssignmentID).(*sqlast.Ident)}

	if v := item.GetRightIdents(); v != nil {
		output.Value = compoundTo(v)
	} else if v := item.GetSingleQuotedString(); v != nil {
		output.Value = stringTo(v)
	} else if v := item.GetDoubleValue(); v != nil {
		output.Value = doubleTo(v)
	} else if v := item.GetLongValue(); v != nil {
		output.Value = longTo(v)
	} else if v := item.GetQueryValue(); v != nil {
		output.Value = subqueryTo(v)
	} else if v := item.GetRightBinary(); v != nil {
		output.Value = binaryexprTo(v)
	}

	return output
}
