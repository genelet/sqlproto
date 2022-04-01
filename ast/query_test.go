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

func TestQuery(t *testing.T) {
	strs := []string{
	"UPDATE y.test_table SET cola = 'a' WHERE bb in (SELECT region FROM top_regions)",
	"UPDATE wine SET vendorcost=vendor.ship-1, price=19.9, stock = stock - ( SELECT SUM (quantity) FROM order WHERE date = CURRENT_DATE AND order.wine_name = wine.name) WHERE x=1 AND y=2 AND z=m",
	"INSERT INTO CUSTOMERS SELECT * FROM CUSTOMERS WHERE ID IN (SELECT ID FROM CUSTOMERS)",
	"INSERT INTO CUSTOMERS VALUES (1, 1.45, 'xxx')",
	"INSERT INTO CUSTOMERS (id, y, x) VALUES (1, 1.45, 'xxx')",
	"DELETE FROM CUSTOMERS WHERE ID = 1",
	"DELETE FROM CUSTOMERS WHERE ID IN (SELECT ID FROM CUSTOMER_IDS)",
	"SELECT a from test_table",
	"SELECT * from test_table",
	"SELECT test_table.* from test_table",
	"SELECT y as z from test_table where x > 6 order by x asc, y.z desc",
	"SELECT aa from test_table where bb = 1.414 and cc = 'john'",
	"SELECT COUNT(*) from test_table group by x.a having x.a > 6",
	"SELECT *, aa, bb as s, COUNT(*), count(*) t, count(cc), count(cc) as u from test_table group by x having x > 6",
	"SELECT aa from test_table where c=1 and bb in (SELECT region FROM top_regions)",
	"SELECT orders.product FROM orders LEFT JOIN accounts as acs ON orders.account_id = accounts.id INNER JOIN accounts_type ON accounts_type.type_id = accounts.type_id ",
	"SELECT orders.product as prod, SUM(orders.quantity) AS product_units, accounts.* FROM orders LEFT JOIN accounts ON orders.account_id = accounts.id INNER JOIN accounts_type ON accounts_type.type_id = accounts.type_id WHERE orders.region IN (SELECT region FROM top_regions) ORDER BY product_units ASC LIMIT 100",
	"WITH regional_sales AS (" +
	"SELECT region, SUM(amount) AS total_sales " +
	"FROM orders GROUP BY region) " +
	"SELECT product, SUM(quantity) AS product_units " +
	"FROM orders " +
	"WHERE region IN (SELECT region FROM top_regions) " +
	"GROUP BY region, product",
	"SELECT x FROM y UNION SELECT x FROM z",
	"SELECT name FROM stadium EXCEPT SELECT T2.name FROM concert AS T1 JOIN stadium AS T2 ON T1.stadium_id  =  T2.stadium_id WHERE T1.year  =  2014",
	"SELECT x FROM a UNION SELECT x FROM b EXCEPT select x FROM c"}

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
		}

		if strings.ToLower(str1) != strings.ToLower(str2) {
			t.Errorf("%d=>%s", i, str1)
			t.Errorf("%d=>%s", i, str2)
		}
	}
}
