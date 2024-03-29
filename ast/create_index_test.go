package ast

import (
	"bytes"
	"strings"
	"testing"

//	"github.com/k0kubun/pp/v3"

	"github.com/akito0107/xsqlparser"
	"github.com/akito0107/xsqlparser/sqlast"
	"github.com/akito0107/xsqlparser/dialect"
)

func TestCreateIndex(t *testing.T) {
	strs := []string{
	`CREATE UNIQUE INDEX customers_idx ON customers USING gist (name) WHERE name = 'test';`,
	`CREATE UNIQUE INDEX customers_idx ON customers (name, email);`,
	`CREATE UNIQUE INDEX customers_idx ON customers USING gist (name);`}

	for i, str := range strs {
		if i == 100 { continue }
		parser, err := xsqlparser.NewParser(bytes.NewBufferString(str), &dialect.GenericSQLDialect{}, xsqlparser.ParseComment())
		if err != nil { t.Fatal(err) }

		file, err := parser.ParseFile()
		if err != nil { t.Fatal(err) }

		createIndexStmt := file.Stmts[0].(*sqlast.CreateIndexStmt)
//pp.Println(createIndexStmt)
		createIndex, err := XCreateIndexTo(createIndexStmt)
		if err != nil { t.Fatal(err) }

		reverse := CreateIndexTo(createIndex)
//pp.Println(reverse)
		if strings.ToLower(createIndexStmt.ToSQLString()) != strings.ToLower(reverse.ToSQLString()) {
			t.Errorf("%d=>%s", i, createIndexStmt.ToSQLString())
			t.Errorf("%d=>%s", i, reverse.ToSQLString())
		}
	}
}
