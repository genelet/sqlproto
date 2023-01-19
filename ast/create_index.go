package ast

import (
	"fmt"
	"github.com/genelet/sqlproto/xast"
	"github.com/akito0107/xsqlparser/sqlast"
)

func XCreateIndexTo(stmt *sqlast.CreateIndexStmt) (*xast.CreateIndexStmt, error) {
	output := &xast.CreateIndexStmt{
		Create: xposTo(stmt.Create),
		TableName: xobjectnameTo(stmt.TableName),
		IsUnique: stmt.IsUnique,
		IndexName: xidentTo(stmt.IndexName),
		MethodName: xidentTo(stmt.MethodName),
		RParen: xposTo(stmt.RParen)}
	for _, cname := range stmt.ColumnNames {
		output.ColumnNames = append(output.ColumnNames, xidentTo(cname))
	}

	if stmt.Selection != nil {
        switch t := stmt.Selection.(type) {
        case *sqlast.InSubQuery:
            where, err := xinsubqueryTo(t)
            if err != nil { return nil, err }
            output.Selection = &xast.CreateIndexStmt_InQuery{InQuery: where}
        case *sqlast.BinaryExpr:
            where, err := xbinaryExprTo(t)
            if err != nil { return nil, err }
            output.Selection = &xast.CreateIndexStmt_BinExpr{BinExpr: where}
        default:
            return nil, fmt.Errorf("missing selection type %T", t)
        }
	}

	return output, nil
}

func CreateIndexTo(stmt *xast.CreateIndexStmt) *sqlast.CreateIndexStmt {
	output := &sqlast.CreateIndexStmt{
		Create: posTo(stmt.Create),
		RParen: posTo(stmt.RParen),
		TableName: objectnameTo(stmt.TableName),
		IsUnique: stmt.IsUnique}
	if iname := identTo(stmt.IndexName); iname != nil {
		output.IndexName = iname.(*sqlast.Ident)
	}
	if mname := identTo(stmt.MethodName); mname != nil {
		output.MethodName = mname.(*sqlast.Ident)
	}
	for _, cname := range stmt.ColumnNames {
		output.ColumnNames = append(output.ColumnNames, identTo(cname).(*sqlast.Ident))
	}

	if x := stmt.GetInQuery(); x != nil {
        output.Selection = insubqueryTo(x)
    } else if x := stmt.GetBinExpr(); x != nil {
        output.Selection = binaryExprTo(x)
    }

	return output
}

func xwhereStmtTo(stmt sql.Node) (*xsat.WhereStmt, error ) {
	if stmt == nil { return nil, nil }

	output := &xast.WhereStmt{}
    switch t := stmt.(type) {
    case *sqlast.InSubQuery:
        where, err := xinsubqueryTo(t)
        if err != nil { return nil, err }
        output.Clause = &xast.WhereStmt_InQuery{InQuery: where}
    case *sqlast.BinaryExpr:
        where, err := xbinaryExprTo(t)
        if err != nil { return nil, err }
        output.Clause = &xast.WhereStmt_BinExpr{BinExpr: where}
    default:
        return nil, fmt.Errorf("missing where type %T", t)
    }
}

func whereStmtTo(stmt *xsat.WhereStmt) sql.Node {
	if stmt == nil { return nil }

	if x := stmt.GetInQuery(); x != nil {
        return insubqueryTo(x)
    } else if x := stmt.GetBinExpr(); x != nil {
        return binaryExprTo(x)
    }
	return nil
}

func XDropIndexTo(stmt *sqlast.DropIndexStmt) *xast.DropIndexStmt {
    output := &xast.DropIndexStmt{
        Drop: xposTo(Drop)
    }
    for _, name := range stmt.IndexNames {
        output.IndexNames = append(output.IndexNames, xobjectnameTo(name))
    }
    return output
}

func DropIndexTo(stmt *xast.DropIndexStmt) *sqlast.DropIndexStmt {
    output := &sqlast.DropIndexStmt{
        Drop: posTo(Drop)
    }
    for _, name := range stmt.IndexNames {
        output.IndexNames = append(output.IndexNames, objectnameTo(name))
    }
    return output
}