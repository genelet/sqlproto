package ast

import (
	"fmt"
	"strings"
	"github.com/genelet/sqlproto/xast"

	"github.com/akito0107/xsqlparser/sqlast"
	"github.com/akito0107/xsqlparser/sqltoken"
)

func xposTo(pos sqltoken.Pos) *xast.Pos {
	return &xast.Pos{
		Line:int32(pos.Line),
		Col:int32(pos.Col)}
}

func posTo(pos *xast.Pos) sqltoken.Pos {
	return sqltoken.Pos{
		Line:int(pos.Line),
		Col:int(pos.Col)}
}

func xposplusTo(pos sqltoken.Pos) *xast.Pos {
	return &xast.Pos{
		Line:int32(pos.Line),
		Col:int32(pos.Col)+1}
}

func xidentTo(ident *sqlast.Ident) *xast.Ident {
	if ident == nil { return nil }

	return &xast.Ident{
		Value: ident.Value,
		From: xposTo(ident.From),
		To: xposTo(ident.To)}
}

func xwildcardTo(card *sqlast.Wildcard) *xast.Ident {
	if card == nil { return nil }

	return &xast.Ident{
		Value: "*",
		From: xposTo(card.Wildcard),
		To: xposplusTo(card.Wildcard)}
}

func identTo(ident *xast.Ident) sqlast.Node {
	if ident == nil { return nil }

	if ident.Value== "*" {
		return &sqlast.Wildcard{Wildcard: posTo(ident.From)}
	}
	return &sqlast.Ident{
		Value: ident.Value,
		From: posTo(ident.From),
		To: posTo(ident.To)}
}

func xidentsTo(ident *sqlast.Ident) *xast.CompoundIdent {
	if ident == nil { return nil }

	return &xast.CompoundIdent{Idents:[]*xast.Ident{xidentTo(ident)}}
}

func xwildcardsTo(card *sqlast.Wildcard) *xast.CompoundIdent {
	if card == nil { return nil }

	return &xast.CompoundIdent{Idents:[]*xast.Ident{xwildcardTo(card)}}
}

func xcompoundTo(idents *sqlast.CompoundIdent) *xast.CompoundIdent {
	if idents == nil { return nil }

	var xs []*xast.Ident
	for _, item := range idents.Idents {
		xs = append(xs, xidentTo(item))
	}
	return &xast.CompoundIdent{Idents:xs}
}

func xwildcarditemTo(t *sqlast.QualifiedWildcardSelectItem) *xast.CompoundIdent {
	if t == nil { return nil }

	prefix := xobjectnameTo(t.Prefix)
	comp := &xast.CompoundIdent{}
	comp.Idents = append(prefix.Idents, &xast.Ident{
		Value: "*",
		From: xposTo(t.Pos()),
		To: xposplusTo(t.Pos())})
	return comp
}

func compoundTo(idents *xast.CompoundIdent) sqlast.Node {
	if idents == nil { return nil }

	if len(idents.Idents) == 1 {
		return identTo(idents.Idents[0])
	}
	var xs []*sqlast.Ident
	for _, item := range idents.Idents {
		switch t := identTo(item).(type) {
		case *sqlast.Wildcard:
			return &sqlast.QualifiedWildcardSelectItem{
				Prefix: &sqlast.ObjectName{Idents:xs}}
		case *sqlast.Ident:
			xs = append(xs, t)
		default:
		}
	}
	return &sqlast.CompoundIdent{Idents:xs}
}

func xobjectnameTo(idents *sqlast.ObjectName) *xast.ObjectName {
	if idents == nil { return nil }

	var xs []*xast.Ident
	for _, item := range idents.Idents {
		xs = append(xs, xidentTo(item))
	}
	return &xast.ObjectName{Idents:xs}
}

func objectnameTo(idents *xast.ObjectName) *sqlast.ObjectName {
	if idents == nil { return nil }

	var xs []*sqlast.Ident
	for _, item := range idents.Idents {
		xs = append(xs, identTo(item).(*sqlast.Ident))
	}
	return &sqlast.ObjectName{Idents:xs}
}

func compoundToObjectname(idents *xast.CompoundIdent) *sqlast.ObjectName {
	if idents == nil { return nil }

	var xs []*sqlast.Ident
	for _, item := range idents.Idents {
		xs = append(xs, identTo(item).(*sqlast.Ident))
	}
	return &sqlast.ObjectName{Idents:xs}
}

func xoperatorTo(op *sqlast.Operator) *xast.Operator {
	if op == nil { return nil }

	return &xast.Operator{
		Type: xast.OperatorType(op.Type),
		From: xposTo(op.From),
		To: xposTo(op.To)}
}

func operatorTo(op *xast.Operator) *sqlast.Operator {
	if op == nil { return nil }

	return &sqlast.Operator{
		Type: sqlast.OperatorType(op.Type),
		From: posTo(op.From),
		To: posTo(op.To)}
}

func xjoinTypeTo(t *sqlast.JoinType) *xast.JoinType {
	if t == nil { return nil }

	return &xast.JoinType{
		Condition: xast.JoinTypeCondition(t.Condition),
		From: xposTo(t.From),
		To: xposTo(t.To)}
}

func joinTypeTo(t *xast.JoinType) *sqlast.JoinType {
	if t == nil { return nil }

	return &sqlast.JoinType{
		Condition: sqlast.JoinTypeCondition(t.Condition),
		From: posTo(t.From),
		To: posTo(t.To)}
}

func xstringTo(t *sqlast.SingleQuotedString) *xast.SingleQuotedString {
    if t == nil { return nil }

    return &xast.SingleQuotedString{
        Value: t.String,
        From: xposTo(t.From),
        To: xposTo(t.To)}
}

func stringTo(t *xast.SingleQuotedString) *sqlast.SingleQuotedString {
    if t == nil { return nil }

    return &sqlast.SingleQuotedString{
        String: t.Value,
        From: posTo(t.From),
        To: posTo(t.To)}
}

func xdoubleTo(t *sqlast.DoubleValue) *xast.DoubleValue {
    if t == nil { return nil }

    return &xast.DoubleValue{
        Value: t.Double,
        From: xposTo(t.From),
        To: xposTo(t.To)}
}

func doubleTo(t *xast.DoubleValue) *sqlast.DoubleValue {
    if t == nil { return nil }

    return &sqlast.DoubleValue{
        Double: t.Value,
        From: posTo(t.From),
        To: posTo(t.To)}
}

func xlongTo(t *sqlast.LongValue) *xast.LongValue {
    if t == nil { return nil }

    return &xast.LongValue{
        Value: t.Long,
        From: xposTo(t.From),
        To: xposTo(t.To)}
}

func longTo(t *xast.LongValue) *sqlast.LongValue {
    if t == nil { return nil }

    return &sqlast.LongValue{
        Long: t.Value,
        From: posTo(t.From),
        To: posTo(t.To)}
}

func xnullValueTo(t *sqlast.NullValue) *xast.NullValue {
    if t == nil { return nil }

    return &xast.NullValue{
        From: xposTo(t.From),
        To: xposTo(t.To)}
}

func nullValueTo(t *xast.NullValue) *sqlast.NullValue {
    if t == nil { return nil }

    return &sqlast.NullValue{
        From: posTo(t.From),
        To: posTo(t.To)}
}

func xintTo(t *sqlast.Int) *xast.Int {
    if t == nil { return nil }

    return &xast.Int{
        From: xposTo(t.From),
        To: xposTo(t.To),
        IsUnsigned: t.IsUnsigned,
		Unsigned: xposTo(t.Unsigned)}
}

func intTo(t *xast.Int) *sqlast.Int {
    if t == nil { return nil }

    return &sqlast.Int{
        From: posTo(t.From),
        To: posTo(t.To),
        IsUnsigned: t.IsUnsigned,
		Unsigned: posTo(t.Unsigned)}
}

func xsmallIntTo(t *sqlast.SmallInt) *xast.SmallInt {
    if t == nil { return nil }

    return &xast.SmallInt{
        From: xposTo(t.From),
        To: xposTo(t.To),
        IsUnsigned: t.IsUnsigned,
		Unsigned: xposTo(t.Unsigned)}
}

func smallIntTo(t *xast.SmallInt) *sqlast.SmallInt {
    if t == nil { return nil }

    return &sqlast.SmallInt{
        From: posTo(t.From),
        To: posTo(t.To),
        IsUnsigned: t.IsUnsigned,
		Unsigned: posTo(t.Unsigned)}
}

func xbigIntTo(t *sqlast.BigInt) *xast.BigInt {
    if t == nil { return nil }

    return &xast.BigInt{
        From: xposTo(t.From),
        To: xposTo(t.To),
        IsUnsigned: t.IsUnsigned,
		Unsigned: xposTo(t.Unsigned)}
}

func bigIntTo(t *xast.BigInt) *sqlast.BigInt {
    if t == nil { return nil }

    return &sqlast.BigInt{
        From: posTo(t.From),
        To: posTo(t.To),
        IsUnsigned: t.IsUnsigned,
		Unsigned: posTo(t.Unsigned)}
}

func xdecimalTo(t *sqlast.Decimal) *xast.Decimal {
	if t == nil { return nil }

	return &xast.Decimal {
		Precision: uint32(*t.Precision),
		Scale:     uint32(*t.Scale),
		Numeric: xposTo(t.Numeric),
		RParen: xposTo(t.RParen),
		IsUnsigned: t.IsUnsigned,
		Unsigned: xposTo(t.Unsigned)}
}

func decimalTo(t *xast.Decimal) *sqlast.Decimal {
	if t == nil { return nil }

	x := uint(t.Precision)
	y := uint(t.Scale)
	return &sqlast.Decimal {
		Precision: &x,
		Scale:     &y,
		Numeric: posTo(t.Numeric),
		RParen: posTo(t.RParen),
		IsUnsigned: t.IsUnsigned,
		Unsigned: posTo(t.Unsigned)}
}

func xtimestampTo(t *sqlast.Timestamp) *xast.Timestamp {
    if t == nil { return nil }

    return &xast.Timestamp{
		WithTimeZone: t.WithTimeZone,
        Timestamp: xposTo(t.Timestamp),
        Zone: xposTo(t.Zone)}
}

func timestampTo(t *xast.Timestamp) *sqlast.Timestamp {
    if t == nil { return nil }

    return &sqlast.Timestamp{
		WithTimeZone: t.WithTimeZone,
        Timestamp: posTo(t.Timestamp),
        Zone: posTo(t.Zone)}
}

func xcustomTo(t *sqlast.Custom) *xast.Custom {
    if t == nil { return nil }

    return &xast.Custom{Ty: xobjectnameTo(t.Ty)}
}

func customTo(t *xast.Custom) *sqlast.Custom {
    if t == nil { return nil }

    return &sqlast.Custom{
		Ty: objectnameTo(t.Ty)}
}

func xuuidTo(t *sqlast.UUID) *xast.UUID {
    if t == nil { return nil }

    return &xast.UUID{
        From: xposTo(t.From),
        To: xposTo(t.To)}
}

func uuidTo(t *xast.UUID) *sqlast.UUID {
    if t == nil { return nil }

    return &sqlast.UUID{
        From: posTo(t.From),
        To: posTo(t.To)}
}

func xcharTypeTo(t *sqlast.CharType) *xast.CharType {
    if t == nil { return nil }

    return &xast.CharType{
        Size: uint32(*t.Size),
        From: xposTo(t.From),
        To: xposTo(t.To)}
}

func charTypeTo(t *xast.CharType) *sqlast.CharType {
    if t == nil { return nil }

	size := uint(t.Size)
    return &sqlast.CharType{
        Size: &size,
        From: posTo(t.From),
        To: posTo(t.To)}
}

func xvarcharTypeTo(t *sqlast.VarcharType) *xast.VarcharType {
    if t == nil { return nil }

    return &xast.VarcharType{
        Size: uint32(*t.Size),
        Character: xposTo(t.Character),
        Varying: xposTo(t.Varying),
        RParen: xposTo(t.RParen)}
}

func varcharTypeTo(t *xast.VarcharType) *sqlast.VarcharType {
    if t == nil { return nil }

	size := uint(t.Size)
    return &sqlast.VarcharType{
        Size: &size,
        Character: posTo(t.Character),
        Varying: posTo(t.Varying),
        RParen: posTo(t.RParen)}
}

func xfunctionTo(s *sqlast.Function) (*xast.AggFunction, error) {
	name := s.Name.Idents[0]
	aggType := xast.AggType(xast.AggType_value[strings.ToUpper(name.Value)])
	output := &xast.AggFunction{
		TypeName: aggType,
		From: xposTo(name.From),
		To: xposTo(name.To)}
	for _, item := range s.Args {
		x, err := xargsNodeTo(item)
		if err != nil { return nil, err }
		output.RestArgs = append(output.RestArgs, x)
	}
	return output, nil
}

func functionTo(f *xast.AggFunction) *sqlast.Function {
    if f == nil { return nil }

	aggname := xast.AggType_name[int32(f.TypeName)]
	output := &sqlast.Function{
		Name: &sqlast.ObjectName{Idents:[]*sqlast.Ident{&sqlast.Ident{
			Value: aggname,
			From: posTo(f.From),
			To: posTo(f.To)}}}}
	for _, item := range f.RestArgs {
		output.Args = append(output.Args, argsNodeTo(item))
	}
	return output
}

func xsetoperatorTo(op sqlast.SQLSetOperator) (*xast.SetOperator, error) {
    xop := &xast.SetOperator{}
    switch t := op.(type) {
    case *sqlast.UnionOperator:
        xop.Type = xast.SetOperatorType_Union
        xop.From = xposTo(t.From)
        xop.To = xposTo(t.To)
    case *sqlast.IntersectOperator:
        xop.Type = xast.SetOperatorType_Intersect
        xop.From = xposTo(t.From)
        xop.To = xposTo(t.To)
    case *sqlast.ExceptOperator:
        xop.Type = xast.SetOperatorType_Except
        xop.From = xposTo(t.From)
        xop.To = xposTo(t.To)
    default:
        return nil, fmt.Errorf("unknow set operation %#v", op)
    }
	return xop, nil
}

func setoperatorTo(op *xast.SetOperator) sqlast.SQLSetOperator {
    switch op.Type {
    case xast.SetOperatorType_Union:
        return &sqlast.UnionOperator{
            From: posTo(op.From),
            To: posTo(op.To)}
    case xast.SetOperatorType_Intersect:
        return &sqlast.IntersectOperator{
            From: posTo(op.From),
            To: posTo(op.To)}
    case xast.SetOperatorType_Except:
        return &sqlast.ExceptOperator{
            From: posTo(op.From),
            To: posTo(op.To)}
    default:
    }
	return nil
}

func xorderbyTo(orderby *sqlast.OrderByExpr) (*xast.OrderByExpr, error) {
	if orderby == nil { return nil, nil }
	output := &xast.OrderByExpr{
		OrderingPos: xposTo(orderby.OrderingPos)}
	if orderby.ASC == nil {
		output.ASCBool = true
	} else {
		output.ASCBool = *orderby.ASC
	}

	switch t := orderby.Expr.(type) {
	case *sqlast.Ident:
		output.Expr = xidentsTo(t)
	case *sqlast.CompoundIdent:
		output.Expr = xcompoundTo(t)
	default:
		return nil, fmt.Errorf("order by is %#v", orderby.Expr)
	}

	return output, nil
}

func orderbyTo(orderby *xast.OrderByExpr) *sqlast.OrderByExpr {
	if orderby == nil { return nil }
	return &sqlast.OrderByExpr{
		OrderingPos: posTo(orderby.OrderingPos),
		ASC: &orderby.ASCBool,
		Expr: compoundTo(orderby.Expr)}
}

func xlimitTo(limit *sqlast.LimitExpr) *xast.LimitExpr {
	if limit == nil { return nil }
	return &xast.LimitExpr{
		AllBool: limit.All,
		AllPos: xposTo(limit.AllPos),
		Limit: xposTo(limit.Limit),
		LimitValue: xlongTo(limit.LimitValue),
		OffsetValue: xlongTo(limit.OffsetValue)}
}

func limitTo(limit *xast.LimitExpr) *sqlast.LimitExpr {
	if limit == nil { return nil }
	return &sqlast.LimitExpr{
		All: limit.AllBool,
		AllPos: posTo(limit.AllPos),
		Limit: posTo(limit.Limit),
		LimitValue: longTo(limit.LimitValue),
		OffsetValue: longTo(limit.OffsetValue)}
}
