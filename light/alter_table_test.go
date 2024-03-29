package light

import (
	"bytes"
	"strings"
	"testing"

	"github.com/genelet/sqlproto/ast"
//	"google.golang.org/protobuf/encoding/protojson"
//	"github.com/k0kubun/pp/v3"

	"github.com/akito0107/xsqlparser"
	"github.com/akito0107/xsqlparser/sqlast"
	"github.com/akito0107/xsqlparser/dialect"
)

func TestAlterTable(t *testing.T) {
	strs := []string{
	"ALTER TABLE customers ADD COLUMN email character varying(255)",
"ALTER TABLE products DROP COLUMN description CASCADE",
"ALTER TABLE products ADD FOREIGN KEY(test_id) REFERENCES other_table(col1, col2)",
"ALTER TABLE products ALTER COLUMN created_at SET DEFAULT current_timestamp",
"ALTER TABLE products ALTER COLUMN number TYPE numeric(255,10)",
"ALTER TABLE Property_Leans ADD CONSTRAINT constraint_Property_Leans_Property_id_Property_Property_id FOREIGN KEY (Property_id) REFERENCES Property (Property_id) ON UPDATE NO ACTION ON DELETE NO ACTION"}

	for i, str := range strs {
		if i != 5 { continue }
		parser, err := xsqlparser.NewParser(bytes.NewBufferString(str), &dialect.GenericSQLDialect{})
		if err != nil { t.Fatal(err) }

		istmt, err := parser.ParseStatement()
		if err != nil { t.Fatal(err) }
		stmt := istmt.(*sqlast.AlterTableStmt)
//pp.Println(stmt)

		xalterTable, err := ast.XAlterTableTo(stmt)
		if err != nil { t.Fatal(err) }

		alterTable := AlterTableTo(xalterTable)
//t.Errorf("%s", protojson.Format(alterTable))

		reverse2 := XAlterTableTo(alterTable)
		reverse3 := ast.AlterTableTo(reverse2)
//pp.Println(reverse)
		if strings.ToLower(stmt.ToSQLString()) != strings.ToLower(reverse3.ToSQLString()) {
			t.Errorf("%d=>%s", i, stmt.ToSQLString())
			t.Errorf("%d=>%s", i, reverse3.ToSQLString())
		}
	}
}
