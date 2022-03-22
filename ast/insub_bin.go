package ast

import (
	"fmt"
	"github.com/genelet/sqlproto/xast"
	"github.com/akito0107/xsqlparser/sqlast"
)

func xinsubqueryTo(sq *sqlast.InSubQuery) (*xast.InSubQuery, error) {
	query, err := XQueryTo(sq.SubQuery)
	if err != nil { return nil, err }

	output := &xast.InSubQuery{
		SubQuery: query,
		Negated: sq.Negated,
		RParen: xposTo(sq.RParen)}

	if sq.Expr == nil { return output, nil }

	switch t := sq.Expr.(type) {
	case *sqlast.Ident:
		output.Expr = xidentsTo(t)
	case *sqlast.CompoundIdent:
		output.Expr = xcompoundTo(t)
	default:
		return nil, fmt.Errorf("expr is %#v", sq.Expr)
	}

	return output, nil
}

func insubqueryTo(sq *xast.InSubQuery) *sqlast.InSubQuery {
	query := QueryTo(sq.SubQuery)

	return &sqlast.InSubQuery{
		Expr: compoundTo(sq.Expr),
		SubQuery: query,
		Negated: sq.Negated,
		RParen: posTo(sq.RParen)}
}

func xsubqueryTo(sq *sqlast.SubQuery) (*xast.SubQuery, error) {
	query, err := XQueryTo(sq.Query)
	if err != nil { return nil, err }

	return &xast.SubQuery{
		Query: query,
		LParen: xposTo(sq.LParen),
		RParen: xposTo(sq.RParen)}, nil
}

func subqueryTo(sq *xast.SubQuery) *sqlast.SubQuery {
	query := QueryTo(sq.Query)

	return &sqlast.SubQuery{
		Query: query,
		LParen: posTo(sq.LParen),
		RParen: posTo(sq.RParen)}
}

func xbinaryexprTo(binary *sqlast.BinaryExpr) (*xast.BinaryExpr, error) {
	if binary == nil { return nil, nil }

	item := &xast.BinaryExpr{Op: xoperatorTo(binary.Op)}

	switch left := binary.Left.(type) {
	case *sqlast.Ident:
		item.LeftOneOf = &xast.BinaryExpr_LeftIdents{LeftIdents:xidentsTo(left)}
	case *sqlast.CompoundIdent:
		item.LeftOneOf = &xast.BinaryExpr_LeftIdents{LeftIdents:xcompoundTo(left)}
	case *sqlast.BinaryExpr:
		middle, err := xbinaryexprTo(left)
		if err != nil { return nil, err }
		item.LeftOneOf = &xast.BinaryExpr_LeftBinary{LeftBinary:middle}
	default:
		return nil, fmt.Errorf("left type %#v", left)
	}

	switch right := binary.Right.(type) {
	case *sqlast.Ident:
		item.RightOneOf = &xast.BinaryExpr_RightIdents{RightIdents:xidentsTo(right)}
	case *sqlast.CompoundIdent:
		item.RightOneOf = &xast.BinaryExpr_RightIdents{RightIdents:xcompoundTo(right)}
	case *sqlast.BinaryExpr:
		middle, err := xbinaryexprTo(right)
		if err != nil { return nil, err }
		item.RightOneOf = &xast.BinaryExpr_RightBinary{RightBinary:middle}
	case *sqlast.SubQuery:
		insub, err := xsubqueryTo(right)
		if err != nil { return nil, err }
		item.RightOneOf = &xast.BinaryExpr_QueryValue{QueryValue: insub}
	case *sqlast.InSubQuery:
		insub, err := xinsubqueryTo(right)
		if err != nil { return nil, err }
		item.RightOneOf = &xast.BinaryExpr_InQueryValue{InQueryValue: insub}
	case *sqlast.LongValue:
		item.RightOneOf = &xast.BinaryExpr_LongValue{LongValue:xlongTo(right)}
	case *sqlast.SingleQuotedString:
		item.RightOneOf = &xast.BinaryExpr_SingleQuotedString{SingleQuotedString:xstringTo(right)}
	case *sqlast.DoubleValue:
		item.RightOneOf = &xast.BinaryExpr_DoubleValue{DoubleValue:xdoubleTo(right)}
	default:
		return nil, fmt.Errorf("right type %#v", right)
	}

	return item, nil
}

func binaryexprTo(binary *xast.BinaryExpr) *sqlast.BinaryExpr {
	if binary == nil { return nil }

	item := &sqlast.BinaryExpr{Op: operatorTo(binary.Op)}

	if v := binary.GetLeftIdents(); v != nil {
		item.Left = compoundTo(v)
	} else if v := binary.GetLeftBinary(); v != nil {
		item.Left = binaryexprTo(v)
	}

	if v := binary.GetRightIdents(); v != nil {
		item.Right = compoundTo(v)
	} else if v := binary.GetSingleQuotedString(); v != nil {
		item.Right = stringTo(v)
	} else if v := binary.GetDoubleValue(); v != nil {
		item.Right = doubleTo(v)
	} else if v := binary.GetLongValue(); v != nil {
		item.Right = longTo(v)
	} else if v := binary.GetQueryValue(); v != nil {
		item.Right = subqueryTo(v)
	} else if v := binary.GetInQueryValue(); v != nil {
		item.Right = insubqueryTo(v)
	} else if v := binary.GetRightBinary(); v != nil {
		item.Right = binaryexprTo(v)
	}

	return item
}
