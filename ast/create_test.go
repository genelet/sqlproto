package ast

import (
	"bytes"
	"strings"
	"testing"

	"github.com/k0kubun/pp"

	"github.com/akito0107/xsqlparser"
	"github.com/akito0107/xsqlparser/sqlast"
	"github.com/akito0107/xsqlparser/dialect"
)

func TestCreate(t *testing.T) {
	strs := []string{
	"CREATE TABLE persons (person_id int, CONSTRAINT production UNIQUE(test_column), PRIMARY KEY(person_id), CHECK(id > 100), FOREIGN KEY(test_id) REFERENCES other_table(col1, col2))",
	//"CREATE TABLE persons (person_id int PRIMARY KEY NOT NULL, last_name character varying(255) NOT NULL, test_id int NOT NULL REFERENCES test(id1, id2), email character varying(255) UNIQUE NOT NULL, age int NOT NULL CHECK(age > 0 AND age < 100), created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL, INDEX (age))",
	}

	for i, str := range strs {
		//if i != 11 { continue }
		parser, err := xsqlparser.NewParser(bytes.NewBufferString(str), &dialect.GenericSQLDialect{})
		if err != nil { t.Fatal(err) }

		var str1, str2 string
		istmt, err := parser.ParseStatement()
		if err != nil { t.Fatal(err) }
		switch stmt := istmt.(type) {
		case *sqlast.QueryStmt:
			str1 = stmt.ToSQLString()
			xquery, err := XQueryTo(stmt)
			if err != nil { t.Fatal(err) }
			reverse := QueryTo(xquery)
			str2 = reverse.ToSQLString()
		case *sqlast.UpdateStmt:
			str1 = stmt.ToSQLString()
			xupdate, err := XUpdateTo(stmt)
			if err != nil { t.Fatal(err) }
			reverse := UpdateTo(xupdate)
			str2 = reverse.ToSQLString()
		case *sqlast.InsertStmt:
			str1 = stmt.ToSQLString()
			xinsert, err := XInsertTo(stmt)
			if err != nil { t.Fatal(err) }
			reverse := InsertTo(xinsert)
			str2 = reverse.ToSQLString()
			str1 = stmt.ToSQLString()
		case *sqlast.DeleteStmt:
			str1 = stmt.ToSQLString()
			xdelete, err := XDeleteTo(stmt)
			if err != nil { t.Fatal(err) }
			reverse := DeleteTo(xdelete)
			str2 = reverse.ToSQLString()
			str1 = stmt.ToSQLString()
		default:
			pp.Println(stmt)
			panic(nil)
		}

		if strings.ToLower(str1) != strings.ToLower(str2) {
			t.Errorf("%d=>%s", i, str1)
			t.Errorf("%d=>%s", i, str2)
		}
	}
}
