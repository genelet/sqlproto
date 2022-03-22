package ast

import (
	"fmt"
	"github.com/genelet/sqlproto/xast"
	"github.com/akito0107/xsqlparser/sqlast"
)

func xbodywhereTo(body *sqlast.SQLSelect, query *xast.QueryStmt_SQLSelect) error {
	if body.WhereClause == nil { return nil }

	switch t := body.WhereClause.(type) {
	case *sqlast.InSubQuery:
		where, err := xinsubqueryTo(t)
		if err != nil { return err }
		query.WhereClause = &xast.QueryStmt_SQLSelect_InQuery{InQuery: where}
	case *sqlast.BinaryExpr:
		where, err := xbinaryexprTo(t)
		if err != nil { return err }
		query.WhereClause = &xast.QueryStmt_SQLSelect_BinExpr{BinExpr: where}
	default:
		return fmt.Errorf("'where' type %#v", t)
	}
	return nil
}

func bodywhereTo(body *xast.QueryStmt_SQLSelect, query *sqlast.SQLSelect) {
	if body == nil { return }
	if v := body.GetInQuery(); v != nil {
		query.WhereClause = insubqueryTo(v)
	} else if v := body.GetBinExpr(); v != nil {
		query.WhereClause = binaryexprTo(v)
	}
}

func xupdatewhereTo(stmt *sqlast.UpdateStmt, update *xast.UpdateStmt) error {
	if stmt.Selection == nil { return nil }

	switch t := stmt.Selection.(type) {
	case *sqlast.InSubQuery:
		where, err := xinsubqueryTo(t)
		if err != nil { return err }
		update.Selection = &xast.UpdateStmt_InQuery{InQuery: where}
	case *sqlast.BinaryExpr:
		where, err := xbinaryexprTo(t)
		if err != nil { return err }
		update.Selection = &xast.UpdateStmt_BinExpr{BinExpr: where}
	default:
		return fmt.Errorf("'where' type %#v", t)
	}
	return nil
}

func updatewhereTo(body *xast.UpdateStmt, query *sqlast.UpdateStmt) {
	if body == nil { return }
	if v := body.GetInQuery(); v != nil {
		query.Selection = insubqueryTo(v)
	} else if v := body.GetBinExpr(); v != nil {
		query.Selection = binaryexprTo(v)
	}
}

func xdeletewhereTo(stmt *sqlast.DeleteStmt, update *xast.DeleteStmt) error {
	if stmt.Selection == nil { return nil }

	switch t := stmt.Selection.(type) {
	case *sqlast.InSubQuery:
		where, err := xinsubqueryTo(t)
		if err != nil { return err }
		update.Selection = &xast.DeleteStmt_InQuery{InQuery: where}
	case *sqlast.BinaryExpr:
		where, err := xbinaryexprTo(t)
		if err != nil { return err }
		update.Selection = &xast.DeleteStmt_BinExpr{BinExpr: where}
	default:
		return fmt.Errorf("'where' type %#v", t)
	}
	return nil
}

func deletewhereTo(body *xast.DeleteStmt, query *sqlast.DeleteStmt) {
	if body == nil { return }
	if v := body.GetInQuery(); v != nil {
		query.Selection = insubqueryTo(v)
	} else if v := body.GetBinExpr(); v != nil {
		query.Selection = binaryexprTo(v)
	}
}
