package ast

import (
	"github.com/genelet/sqlproto/xast"
	"github.com/akito0107/xsqlparser/sqlast"
)

// XDeleteTo translates a xsqlparser delete statement into xast statement
//
// see https://github.com/akito0107/xsqlparser for xsqlparser
//
func XDeleteTo(stmt *sqlast.DeleteStmt) (*xast.DeleteStmt, error) {
	output := &xast.DeleteStmt{
		Delete: xposTo(stmt.Delete),
		TableName: xobjectnameTo(stmt.TableName)}

	err := xdeletewhereTo(stmt, output)
	return output, err
}

// DeleteTo translates a xast delete statement into xsqlparser statement
//
func DeleteTo(stmt *xast.DeleteStmt) *sqlast.DeleteStmt {
	output := &sqlast.DeleteStmt{
		TableName: compoundToObjectname(stmt.TableName)}
	if stmt.Delete != nil {
		output.Delete = posTo(stmt.Delete)
	}

	deletewhereTo(stmt, output)

	return output
}
