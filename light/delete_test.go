package light

import (
	"bytes"
	"strings"
	"testing"

	 "github.com/genelet/sqlproto/ast"

//	"github.com/k0kubun/pp/v3"

	"github.com/akito0107/xsqlparser"
	"github.com/akito0107/xsqlparser/sqlast"
	"github.com/akito0107/xsqlparser/dialect"
)

func TestDelete(t *testing.T) {
	strs := []string{
	"DELETE FROM customers WHERE customer_id = 1"}

	for i, str := range strs {
		//if i != 17 { continue }
		parser, err := xsqlparser.NewParser(bytes.NewBufferString(str), &dialect.GenericSQLDialect{})
		if err != nil { t.Fatal(err) }

		istmt, err := parser.ParseStatement()
		if err != nil { t.Fatal(err) }
		stmt := istmt.(*sqlast.DeleteStmt)
//pp.Println(stmt)

		xdelete, err := ast.XDeleteTo(stmt)
		if err != nil { t.Fatal(err) }

		delete0 := DeleteTo(xdelete)
		reverse2 := XDeleteTo(delete0)
		reverse3 := ast.DeleteTo(reverse2)
//pp.Println(reverse)
		if strings.ToLower(stmt.ToSQLString()) != strings.ToLower(reverse3.ToSQLString()) {
			t.Errorf("%d=>%s", i, stmt.ToSQLString())
			t.Errorf("%d=>%s", i, reverse3.ToSQLString())
		}
	}
}
