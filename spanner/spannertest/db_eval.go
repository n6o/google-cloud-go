package spannertest

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"cloud.google.com/go/civil"
	"cloud.google.com/go/spanner/spansql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type evalContext struct {
	cols    []colInfo
	row     row
	aliases map[spansql.ID]spansql.Expr
	params  queryParams
}
type coercedValue struct {
	spansql.Expr
	val interface {
	}
	orig spansql.Expr
}

func (cv coercedValue) gologoo__SQL_71e383cb26f1785bb4841431523136cf() string {
	return cv.orig.SQL()
}
func (ec evalContext) gologoo__evalExprList_71e383cb26f1785bb4841431523136cf(list []spansql.Expr) ([]interface {
}, error) {
	var out []interface {
	}
	for _, e := range list {
		x, err := ec.evalExpr(e)
		if err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, nil
}
func (ec evalContext) gologoo__evalBoolExpr_71e383cb26f1785bb4841431523136cf(be spansql.BoolExpr) (*bool, error) {
	switch be := be.(type) {
	default:
		return nil, fmt.Errorf("unhandled BoolExpr %T", be)
	case spansql.BoolLiteral:
		b := bool(be)
		return &b, nil
	case spansql.ID, spansql.Param, spansql.Paren, spansql.Func, spansql.InOp:
		e, err := ec.evalExpr(be)
		if err != nil {
			return nil, err
		}
		if e == nil {
			return nil, nil
		}
		b, ok := e.(bool)
		if !ok {
			return nil, fmt.Errorf("got %T, want bool", e)
		}
		return &b, nil
	case spansql.LogicalOp:
		var lhs, rhs *bool
		var err error
		if be.LHS != nil {
			lhs, err = ec.evalBoolExpr(be.LHS)
			if err != nil {
				return nil, err
			}
		}
		rhs, err = ec.evalBoolExpr(be.RHS)
		if err != nil {
			return nil, err
		}
		switch be.Op {
		case spansql.And:
			if lhs != nil {
				if *lhs {
					return rhs, nil
				}
				return lhs, nil
			}
			if rhs != nil && !*rhs {
				return rhs, nil
			}
			return nil, nil
		case spansql.Or:
			if lhs != nil {
				if *lhs {
					return lhs, nil
				}
				return rhs, nil
			}
			if rhs != nil && *rhs {
				return rhs, nil
			}
			return nil, nil
		case spansql.Not:
			if rhs == nil {
				return nil, nil
			}
			b := !*rhs
			return &b, nil
		default:
			return nil, fmt.Errorf("unhandled LogicalOp %d", be.Op)
		}
	case spansql.ComparisonOp:
		be, err := ec.coerceComparisonOpArgs(be)
		if err != nil {
			return nil, err
		}
		lhs, err := ec.evalExpr(be.LHS)
		if err != nil {
			return nil, err
		}
		rhs, err := ec.evalExpr(be.RHS)
		if err != nil {
			return nil, err
		}
		if lhs == nil || rhs == nil {
			return nil, nil
		}
		var b bool
		switch be.Op {
		default:
			return nil, fmt.Errorf("TODO: ComparisonOp %d", be.Op)
		case spansql.Lt:
			b = compareVals(lhs, rhs) < 0
		case spansql.Le:
			b = compareVals(lhs, rhs) <= 0
		case spansql.Gt:
			b = compareVals(lhs, rhs) > 0
		case spansql.Ge:
			b = compareVals(lhs, rhs) >= 0
		case spansql.Eq:
			b = compareVals(lhs, rhs) == 0
		case spansql.Ne:
			b = compareVals(lhs, rhs) != 0
		case spansql.Like, spansql.NotLike:
			left, ok := lhs.(string)
			if !ok {
				return nil, fmt.Errorf("LHS of LIKE is %T, not string", lhs)
			}
			right, ok := rhs.(string)
			if !ok {
				return nil, fmt.Errorf("RHS of LIKE is %T, not string", rhs)
			}
			b = evalLike(left, right)
			if be.Op == spansql.NotLike {
				b = !b
			}
		case spansql.Between, spansql.NotBetween:
			rhs2, err := ec.evalExpr(be.RHS2)
			if err != nil {
				return nil, err
			}
			b = compareVals(rhs, lhs) <= 0 && compareVals(lhs, rhs2) <= 0
			if be.Op == spansql.NotBetween {
				b = !b
			}
		}
		return &b, nil
	case spansql.IsOp:
		lhs, err := ec.evalExpr(be.LHS)
		if err != nil {
			return nil, err
		}
		var b bool
		switch rhs := be.RHS.(type) {
		default:
			return nil, fmt.Errorf("unhandled IsOp %T", rhs)
		case spansql.BoolLiteral:
			if lhs == nil {
				lhs = !bool(rhs)
			}
			lhsBool, ok := lhs.(bool)
			if !ok {
				return nil, fmt.Errorf("non-bool value %T on LHS for %s", lhs, be.SQL())
			}
			b = (lhsBool == bool(rhs))
		case spansql.NullLiteral:
			b = (lhs == nil)
		}
		if be.Neg {
			b = !b
		}
		return &b, nil
	}
}
func (ec evalContext) gologoo__evalArithOp_71e383cb26f1785bb4841431523136cf(e spansql.ArithOp) (interface {
}, error) {
	switch e.Op {
	case spansql.Neg:
		rhs, err := ec.evalExpr(e.RHS)
		if err != nil {
			return nil, err
		}
		if rhs == nil {
			return nil, nil
		}
		switch rhs := rhs.(type) {
		case float64:
			return -rhs, nil
		case int64:
			return -rhs, nil
		}
		return nil, fmt.Errorf("RHS of %s evaluates to %T, want FLOAT64 or INT64", e.SQL(), rhs)
	case spansql.BitNot:
		rhs, err := ec.evalExpr(e.RHS)
		if err != nil {
			return nil, err
		}
		if rhs == nil {
			return nil, nil
		}
		switch rhs := rhs.(type) {
		case int64:
			return ^rhs, nil
		case []byte:
			b := append([]byte(nil), rhs...)
			for i := range b {
				b[i] = ^b[i]
			}
			return b, nil
		}
		return nil, fmt.Errorf("RHS of %s evaluates to %T, want INT64 or BYTES", e.SQL(), rhs)
	case spansql.Div:
		lhs, err := ec.evalFloat64(e.LHS)
		if err != nil {
			return nil, err
		}
		rhs, err := ec.evalFloat64(e.RHS)
		if err != nil {
			return nil, err
		}
		if rhs == 0 {
			return nil, fmt.Errorf("divide by zero")
		}
		return lhs / rhs, nil
	case spansql.Add, spansql.Sub, spansql.Mul:
		lhs, err := ec.evalExpr(e.LHS)
		if err != nil {
			return nil, err
		}
		if lhs == nil {
			return nil, nil
		}
		rhs, err := ec.evalExpr(e.RHS)
		if err != nil {
			return nil, err
		}
		if rhs == nil {
			return nil, nil
		}
		i1, ok1 := lhs.(int64)
		i2, ok2 := rhs.(int64)
		if ok1 && ok2 {
			switch e.Op {
			case spansql.Add:
				return i1 + i2, nil
			case spansql.Sub:
				return i1 - i2, nil
			case spansql.Mul:
				return i1 * i2, nil
			}
		}
		f1, err := asFloat64(e.LHS, lhs)
		if err != nil {
			return nil, err
		}
		f2, err := asFloat64(e.RHS, rhs)
		if err != nil {
			return nil, err
		}
		switch e.Op {
		case spansql.Add:
			return f1 + f2, nil
		case spansql.Sub:
			return f1 - f2, nil
		case spansql.Mul:
			return f1 * f2, nil
		}
	case spansql.BitAnd, spansql.BitXor, spansql.BitOr:
		lhs, err := ec.evalExpr(e.LHS)
		if err != nil {
			return nil, err
		}
		if lhs == nil {
			return nil, nil
		}
		rhs, err := ec.evalExpr(e.RHS)
		if err != nil {
			return nil, err
		}
		if rhs == nil {
			return nil, nil
		}
		i1, ok1 := lhs.(int64)
		i2, ok2 := rhs.(int64)
		if ok1 && ok2 {
			switch e.Op {
			case spansql.BitAnd:
				return i1 & i2, nil
			case spansql.BitXor:
				return i1 ^ i2, nil
			case spansql.BitOr:
				return i1 | i2, nil
			}
		}
		b1, ok1 := lhs.([]byte)
		b2, ok2 := rhs.([]byte)
		if !ok1 || !ok2 {
			return nil, fmt.Errorf("arguments of %s evaluate to (%T, %T), want (INT64, INT64) or (BYTES, BYTES)", e.SQL(), lhs, rhs)
		}
		if len(b1) != len(b2) {
			return nil, fmt.Errorf("arguments of %s evaluate to BYTES of unequal lengths (%d vs %d)", e.SQL(), len(b1), len(b2))
		}
		var f func(x, y byte) byte
		switch e.Op {
		case spansql.BitAnd:
			f = func(x, y byte) byte {
				return x & y
			}
		case spansql.BitXor:
			f = func(x, y byte) byte {
				return x ^ y
			}
		case spansql.BitOr:
			f = func(x, y byte) byte {
				return x | y
			}
		}
		b := make([]byte, len(b1))
		for i := range b1 {
			b[i] = f(b1[i], b2[i])
		}
		return b, nil
	}
	return nil, fmt.Errorf("TODO: evalArithOp(%s %v)", e.SQL(), e.Op)
}
func (ec evalContext) gologoo__evalFunc_71e383cb26f1785bb4841431523136cf(e spansql.Func) (interface {
}, spansql.Type, error) {
	var err error
	if f, ok := functions[e.Name]; ok {
		args := make([]interface {
		}, len(e.Args))
		types := make([]spansql.Type, len(e.Args))
		for i, arg := range e.Args {
			if args[i], err = ec.evalExpr(arg); err != nil {
				return nil, spansql.Type{}, err
			}
			if te, ok := arg.(spansql.TypedExpr); ok {
				types[i] = te.Type
			}
		}
		return f.Eval(args, types)
	}
	return nil, spansql.Type{}, status.Errorf(codes.Unimplemented, "function %q is not implemented", e.Name)
}
func (ec evalContext) gologoo__evalFloat64_71e383cb26f1785bb4841431523136cf(e spansql.Expr) (float64, error) {
	v, err := ec.evalExpr(e)
	if err != nil {
		return 0, err
	}
	return asFloat64(e, v)
}
func gologoo__asFloat64_71e383cb26f1785bb4841431523136cf(e spansql.Expr, v interface {
}) (float64, error) {
	switch v := v.(type) {
	default:
		return 0, fmt.Errorf("expression %s evaluates to %T, want FLOAT64 or INT64", e.SQL(), v)
	case float64:
		return v, nil
	case int64:
		return float64(v), nil
	}
}
func (ec evalContext) gologoo__evalExpr_71e383cb26f1785bb4841431523136cf(e spansql.Expr) (interface {
}, error) {
	evalBool := func(be spansql.BoolExpr) (interface {
	}, error) {
		b, err := ec.evalBoolExpr(be)
		if err != nil {
			return nil, err
		}
		if b == nil {
			return nil, nil
		}
		return *b, nil
	}
	switch e := e.(type) {
	default:
		return nil, fmt.Errorf("TODO: evalExpr(%s %T)", e.SQL(), e)
	case coercedValue:
		return e.val, nil
	case spansql.PathExp:
		return ec.evalPathExp(e)
	case spansql.ID:
		return ec.evalID(e)
	case spansql.Param:
		qp, ok := ec.params[string(e)]
		if !ok {
			return 0, fmt.Errorf("unbound param %s", e.SQL())
		}
		return qp.Value, nil
	case spansql.IntegerLiteral:
		return int64(e), nil
	case spansql.FloatLiteral:
		return float64(e), nil
	case spansql.StringLiteral:
		return string(e), nil
	case spansql.BytesLiteral:
		return []byte(e), nil
	case spansql.NullLiteral:
		return nil, nil
	case spansql.BoolLiteral:
		return bool(e), nil
	case spansql.Paren:
		return ec.evalExpr(e.Expr)
	case spansql.TypedExpr:
		return ec.evalTypedExpr(e)
	case spansql.ExtractExpr:
		return ec.evalExtractExpr(e)
	case spansql.AtTimeZoneExpr:
		return ec.evalAtTimeZoneExpr(e)
	case spansql.Func:
		v, _, err := ec.evalFunc(e)
		if err != nil {
			return nil, err
		}
		return v, nil
	case spansql.Array:
		var arr []interface {
		}
		for _, elt := range e {
			v, err := ec.evalExpr(elt)
			if err != nil {
				return nil, err
			}
			arr = append(arr, v)
		}
		return arr, nil
	case spansql.ArithOp:
		return ec.evalArithOp(e)
	case spansql.LogicalOp:
		return evalBool(e)
	case spansql.ComparisonOp:
		return evalBool(e)
	case spansql.InOp:
		if len(e.RHS) == 0 {
			return e.Neg, nil
		}
		lhs, err := ec.evalExpr(e.LHS)
		if err != nil {
			return false, err
		}
		if lhs == nil {
			return nil, nil
		}
		var b bool
		for _, rhse := range e.RHS {
			rhs, err := ec.evalExpr(rhse)
			if err != nil {
				return false, err
			}
			if !e.Unnest {
				if lhs == rhs {
					b = true
				}
			} else {
				if rhs == nil {
					return e.Neg, nil
				}
				arr, ok := rhs.([]interface {
				})
				if !ok {
					return nil, fmt.Errorf("UNNEST argument evaluated as %T, want array", rhs)
				}
				for _, rhs := range arr {
					if compareVals(lhs, rhs) == 0 {
						b = true
					}
				}
			}
		}
		if e.Neg {
			b = !b
		}
		return b, nil
	case spansql.IsOp:
		return evalBool(e)
	case aggSentinel:
		ci := -1
		for i, col := range ec.cols {
			if col.AggIndex == e.AggIndex {
				ci = i
				break
			}
		}
		if ci < 0 {
			return 0, fmt.Errorf("internal error: did not find aggregate column %d", e.AggIndex)
		}
		return ec.row[ci], nil
	}
}
func (ec evalContext) gologoo__resolveColumnIndex_71e383cb26f1785bb4841431523136cf(e spansql.Expr) (int, error) {
	switch e := e.(type) {
	case spansql.ID:
		for i, col := range ec.cols {
			if col.Name == e {
				return i, nil
			}
		}
	case spansql.PathExp:
		for i, col := range ec.cols {
			if pathExpEqual(e, col.Alias) {
				return i, nil
			}
		}
	}
	return 0, fmt.Errorf("couldn't resolve [%s] as a table column", e.SQL())
}
func (ec evalContext) gologoo__evalPathExp_71e383cb26f1785bb4841431523136cf(pe spansql.PathExp) (interface {
}, error) {
	if i, err := ec.resolveColumnIndex(pe); err == nil {
		return ec.row.copyDataElem(i), nil
	}
	return nil, fmt.Errorf("couldn't resolve path expression %s", pe.SQL())
}
func (ec evalContext) gologoo__evalID_71e383cb26f1785bb4841431523136cf(id spansql.ID) (interface {
}, error) {
	if i, err := ec.resolveColumnIndex(id); err == nil {
		return ec.row.copyDataElem(i), nil
	}
	if e, ok := ec.aliases[id]; ok {
		innerEC := ec
		innerEC.aliases = make(map[spansql.ID]spansql.Expr)
		for alias, e := range ec.aliases {
			if alias != id {
				innerEC.aliases[alias] = e
			}
		}
		return innerEC.evalExpr(e)
	}
	return nil, fmt.Errorf("couldn't resolve identifier %s", id)
}
func (ec evalContext) gologoo__coerceComparisonOpArgs_71e383cb26f1785bb4841431523136cf(co spansql.ComparisonOp) (spansql.ComparisonOp, error) {
	if co.RHS2 != nil {
		return co, nil
	}
	var err error
	if slit, ok := co.LHS.(spansql.StringLiteral); ok {
		co.LHS, err = ec.coerceString(co.RHS, slit)
		return co, err
	}
	if slit, ok := co.RHS.(spansql.StringLiteral); ok {
		co.RHS, err = ec.coerceString(co.LHS, slit)
		return co, err
	}
	return co, nil
}
func (ec evalContext) gologoo__coerceString_71e383cb26f1785bb4841431523136cf(target spansql.Expr, slit spansql.StringLiteral) (spansql.Expr, error) {
	ci, err := ec.colInfo(target)
	if err != nil {
		return nil, err
	}
	if ci.Type.Array {
		return nil, fmt.Errorf("unable to coerce string literal %q to match array type", slit)
	}
	switch ci.Type.Base {
	case spansql.String:
		return slit, nil
	case spansql.Date:
		d, err := parseAsDate(string(slit))
		if err != nil {
			return nil, fmt.Errorf("coercing string literal %q to DATE: %v", slit, err)
		}
		return coercedValue{val: d, orig: slit}, nil
	case spansql.Timestamp:
		t, err := parseAsTimestamp(string(slit))
		if err != nil {
			return nil, fmt.Errorf("coercing string literal %q to TIMESTAMP: %v", slit, err)
		}
		return coercedValue{val: t, orig: slit}, nil
	}
	return nil, fmt.Errorf("unable to coerce string literal %q to match %v", slit, ci.Type)
}
func (ec evalContext) gologoo__evalTypedExpr_71e383cb26f1785bb4841431523136cf(expr spansql.TypedExpr) (result interface {
}, err error) {
	val, err := ec.evalExpr(expr.Expr)
	if err != nil {
		return nil, err
	}
	return convert(val, expr.Type)
}
func (ec evalContext) gologoo__evalExtractExpr_71e383cb26f1785bb4841431523136cf(expr spansql.ExtractExpr) (result interface {
}, err error) {
	val, err := ec.evalExpr(expr.Expr)
	if err != nil {
		return nil, err
	}
	switch expr.Part {
	case "DATE":
		switch v := val.(type) {
		case time.Time:
			return civil.DateOf(v), nil
		case civil.Date:
			return v, nil
		}
	case "DAY":
		switch v := val.(type) {
		case time.Time:
			return int64(v.Day()), nil
		case civil.Date:
			return int64(v.Day), nil
		}
	case "MONTH":
		switch v := val.(type) {
		case time.Time:
			return int64(v.Month()), nil
		case civil.Date:
			return int64(v.Month), nil
		}
	case "YEAR":
		switch v := val.(type) {
		case time.Time:
			return int64(v.Year()), nil
		case civil.Date:
			return int64(v.Year), nil
		}
	}
	return nil, fmt.Errorf("Extract with part %v not supported", expr.Part)
}
func (ec evalContext) gologoo__evalAtTimeZoneExpr_71e383cb26f1785bb4841431523136cf(expr spansql.AtTimeZoneExpr) (result interface {
}, err error) {
	val, err := ec.evalExpr(expr.Expr)
	if err != nil {
		return nil, err
	}
	switch v := val.(type) {
	case time.Time:
		loc, err := time.LoadLocation(expr.Zone)
		if err != nil {
			return nil, fmt.Errorf("AtTimeZone with %T not supported", v)
		}
		return v.In(loc), nil
	default:
		return nil, fmt.Errorf("AtTimeZone with %T not supported", val)
	}
}
func gologoo__evalLiteralOrParam_71e383cb26f1785bb4841431523136cf(lop spansql.LiteralOrParam, params queryParams) (int64, error) {
	switch v := lop.(type) {
	case spansql.IntegerLiteral:
		return int64(v), nil
	case spansql.Param:
		return paramAsInteger(v, params)
	default:
		return 0, fmt.Errorf("LiteralOrParam with %T not supported", v)
	}
}
func gologoo__paramAsInteger_71e383cb26f1785bb4841431523136cf(p spansql.Param, params queryParams) (int64, error) {
	qp, ok := params[string(p)]
	if !ok {
		return 0, fmt.Errorf("unbound param %s", p.SQL())
	}
	switch v := qp.Value.(type) {
	default:
		return 0, fmt.Errorf("can't interpret parameter %s (%s) value of type %T as integer", p.SQL(), qp.Type.SQL(), v)
	case int64:
		return v, nil
	case string:
		x, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("bad int64 string %q: %v", v, err)
		}
		return x, nil
	}
}
func gologoo__compareValLists_71e383cb26f1785bb4841431523136cf(a, b []interface {
}, desc []bool) int {
	for i := range a {
		cmp := compareVals(a[i], b[i])
		if cmp == 0 {
			continue
		}
		if desc != nil && desc[i] {
			cmp = -cmp
		}
		return cmp
	}
	return 0
}
func gologoo__compareVals_71e383cb26f1785bb4841431523136cf(x, y interface {
}) int {
	if x == nil && y == nil {
		return 0
	} else if x == nil {
		return -1
	} else if y == nil {
		return 1
	}
	switch x := x.(type) {
	default:
		panic(fmt.Sprintf("unhandled comparison on %T", x))
	case bool:
		y := y.(bool)
		if !x && y {
			return -1
		} else if x && !y {
			return 1
		}
		return 0
	case int64:
		if s, ok := y.(string); ok {
			var err error
			y, err = strconv.ParseInt(s, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("bad int64 string %q: %v", s, err))
			}
		}
		if f, ok := y.(float64); ok {
			return compareVals(x, f)
		}
		y := y.(int64)
		if x < y {
			return -1
		} else if x > y {
			return 1
		}
		return 0
	case float64:
		if i, ok := y.(int64); ok {
			y = float64(i)
		}
		y := y.(float64)
		if x < y {
			return -1
		} else if x > y {
			return 1
		}
		return 0
	case string:
		return strings.Compare(x, y.(string))
	case civil.Date:
		y := y.(civil.Date)
		if x.Before(y) {
			return -1
		} else if x.After(y) {
			return 1
		}
		return 0
	case time.Time:
		y := y.(time.Time)
		if x.Before(y) {
			return -1
		} else if x.After(y) {
			return 1
		}
		return 0
	case []byte:
		return bytes.Compare(x, y.([]byte))
	}
}

var (
	boolType    = spansql.Type{Base: spansql.Bool}
	int64Type   = spansql.Type{Base: spansql.Int64}
	float64Type = spansql.Type{Base: spansql.Float64}
	stringType  = spansql.Type{Base: spansql.String}
)

func (ec evalContext) gologoo__colInfo_71e383cb26f1785bb4841431523136cf(e spansql.Expr) (colInfo, error) {
	switch e := e.(type) {
	case spansql.BoolLiteral:
		return colInfo{Type: boolType}, nil
	case spansql.IntegerLiteral:
		return colInfo{Type: int64Type}, nil
	case spansql.StringLiteral:
		return colInfo{Type: stringType}, nil
	case spansql.BytesLiteral:
		return colInfo{Type: spansql.Type{Base: spansql.Bytes}}, nil
	case spansql.ArithOp:
		t, err := ec.arithColType(e)
		if err != nil {
			return colInfo{}, err
		}
		return colInfo{Type: t}, nil
	case spansql.LogicalOp, spansql.ComparisonOp, spansql.IsOp:
		return colInfo{Type: spansql.Type{Base: spansql.Bool}}, nil
	case spansql.PathExp, spansql.ID:
		i, err := ec.resolveColumnIndex(e)
		if err == nil {
			return ec.cols[i], nil
		}
	case spansql.Param:
		qp, ok := ec.params[string(e)]
		if !ok {
			return colInfo{}, fmt.Errorf("unbound param %s", e.SQL())
		}
		return colInfo{Type: qp.Type}, nil
	case spansql.Paren:
		return ec.colInfo(e.Expr)
	case spansql.Func:
		_, t, err := ec.evalFunc(e)
		if err != nil {
			return colInfo{}, err
		}
		return colInfo{Type: t}, nil
	case spansql.Array:
		if len(e) == 0 {
			return colInfo{Type: spansql.Type{Base: spansql.Int64, Array: true}}, nil
		}
		ci, err := ec.colInfo(e[0])
		if err != nil {
			return colInfo{}, err
		}
		if ci.Type.Array {
			return colInfo{}, fmt.Errorf("can't nest array literals")
		}
		ci.Type.Array = true
		return ci, nil
	case spansql.NullLiteral:
		return colInfo{Type: int64Type}, nil
	case aggSentinel:
		return colInfo{Type: e.Type, AggIndex: e.AggIndex}, nil
	}
	return colInfo{}, fmt.Errorf("can't deduce column type from expression [%s] (type %T)", e.SQL(), e)
}
func (ec evalContext) gologoo__arithColType_71e383cb26f1785bb4841431523136cf(ao spansql.ArithOp) (spansql.Type, error) {
	var lhs, rhs spansql.Type
	var err error
	if ao.LHS != nil {
		ci, err := ec.colInfo(ao.LHS)
		if err != nil {
			return spansql.Type{}, err
		}
		lhs = ci.Type
	}
	ci, err := ec.colInfo(ao.RHS)
	if err != nil {
		return spansql.Type{}, err
	}
	rhs = ci.Type
	switch ao.Op {
	default:
		return spansql.Type{}, fmt.Errorf("can't deduce column type from ArithOp [%s]", ao.SQL())
	case spansql.Neg, spansql.BitNot:
		return rhs, nil
	case spansql.Add, spansql.Sub, spansql.Mul:
		if lhs == int64Type && rhs == int64Type {
			return int64Type, nil
		}
		return float64Type, nil
	case spansql.Div:
		return float64Type, nil
	case spansql.Concat:
		if !lhs.Array {
			return stringType, nil
		}
		return lhs, nil
	case spansql.BitShl, spansql.BitShr, spansql.BitAnd, spansql.BitXor, spansql.BitOr:
		return lhs, nil
	}
}
func gologoo__pathExpEqual_71e383cb26f1785bb4841431523136cf(a, b spansql.PathExp) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
func gologoo__evalLike_71e383cb26f1785bb4841431523136cf(str, pat string) bool {
	pat = regexp.QuoteMeta(pat)
	if !strings.HasPrefix(pat, "%") {
		pat = "^" + pat
	}
	if !strings.HasSuffix(pat, "%") {
		pat = pat + "$"
	}
	pat = strings.Replace(pat, "%", ".*", -1)
	pat = strings.Replace(pat, "_", ".", -1)
	match, err := regexp.MatchString(pat, str)
	if err != nil {
		panic(fmt.Sprintf("internal error: constructed bad regexp /%s/: %v", pat, err))
	}
	return match
}
func (cv coercedValue) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : (none)\n")
	r0 := cv.gologoo__SQL_71e383cb26f1785bb4841431523136cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ec evalContext) evalExprList(list []spansql.Expr) ([]interface {
}, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__evalExprList_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v\n", list)
	r0, r1 := ec.gologoo__evalExprList_71e383cb26f1785bb4841431523136cf(list)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (ec evalContext) evalBoolExpr(be spansql.BoolExpr) (*bool, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__evalBoolExpr_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v\n", be)
	r0, r1 := ec.gologoo__evalBoolExpr_71e383cb26f1785bb4841431523136cf(be)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (ec evalContext) evalArithOp(e spansql.ArithOp) (interface {
}, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__evalArithOp_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v\n", e)
	r0, r1 := ec.gologoo__evalArithOp_71e383cb26f1785bb4841431523136cf(e)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (ec evalContext) evalFunc(e spansql.Func) (interface {
}, spansql.Type, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__evalFunc_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v\n", e)
	r0, r1, r2 := ec.gologoo__evalFunc_71e383cb26f1785bb4841431523136cf(e)
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
func (ec evalContext) evalFloat64(e spansql.Expr) (float64, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__evalFloat64_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v\n", e)
	r0, r1 := ec.gologoo__evalFloat64_71e383cb26f1785bb4841431523136cf(e)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func asFloat64(e spansql.Expr, v interface {
}) (float64, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__asFloat64_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v %v\n", e, v)
	r0, r1 := gologoo__asFloat64_71e383cb26f1785bb4841431523136cf(e, v)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (ec evalContext) evalExpr(e spansql.Expr) (interface {
}, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__evalExpr_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v\n", e)
	r0, r1 := ec.gologoo__evalExpr_71e383cb26f1785bb4841431523136cf(e)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (ec evalContext) resolveColumnIndex(e spansql.Expr) (int, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__resolveColumnIndex_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v\n", e)
	r0, r1 := ec.gologoo__resolveColumnIndex_71e383cb26f1785bb4841431523136cf(e)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (ec evalContext) evalPathExp(pe spansql.PathExp) (interface {
}, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__evalPathExp_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v\n", pe)
	r0, r1 := ec.gologoo__evalPathExp_71e383cb26f1785bb4841431523136cf(pe)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (ec evalContext) evalID(id spansql.ID) (interface {
}, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__evalID_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v\n", id)
	r0, r1 := ec.gologoo__evalID_71e383cb26f1785bb4841431523136cf(id)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (ec evalContext) coerceComparisonOpArgs(co spansql.ComparisonOp) (spansql.ComparisonOp, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__coerceComparisonOpArgs_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v\n", co)
	r0, r1 := ec.gologoo__coerceComparisonOpArgs_71e383cb26f1785bb4841431523136cf(co)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (ec evalContext) coerceString(target spansql.Expr, slit spansql.StringLiteral) (spansql.Expr, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__coerceString_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v %v\n", target, slit)
	r0, r1 := ec.gologoo__coerceString_71e383cb26f1785bb4841431523136cf(target, slit)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (ec evalContext) evalTypedExpr(expr spansql.TypedExpr) (result interface {
}, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__evalTypedExpr_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v\n", expr)
	result, err = ec.gologoo__evalTypedExpr_71e383cb26f1785bb4841431523136cf(expr)
	log.Printf("Output: %v %v\n", result, err)
	return
}
func (ec evalContext) evalExtractExpr(expr spansql.ExtractExpr) (result interface {
}, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__evalExtractExpr_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v\n", expr)
	result, err = ec.gologoo__evalExtractExpr_71e383cb26f1785bb4841431523136cf(expr)
	log.Printf("Output: %v %v\n", result, err)
	return
}
func (ec evalContext) evalAtTimeZoneExpr(expr spansql.AtTimeZoneExpr) (result interface {
}, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__evalAtTimeZoneExpr_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v\n", expr)
	result, err = ec.gologoo__evalAtTimeZoneExpr_71e383cb26f1785bb4841431523136cf(expr)
	log.Printf("Output: %v %v\n", result, err)
	return
}
func evalLiteralOrParam(lop spansql.LiteralOrParam, params queryParams) (int64, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__evalLiteralOrParam_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v %v\n", lop, params)
	r0, r1 := gologoo__evalLiteralOrParam_71e383cb26f1785bb4841431523136cf(lop, params)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func paramAsInteger(p spansql.Param, params queryParams) (int64, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__paramAsInteger_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v %v\n", p, params)
	r0, r1 := gologoo__paramAsInteger_71e383cb26f1785bb4841431523136cf(p, params)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func compareValLists(a, b []interface {
}, desc []bool) int {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__compareValLists_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v %v %v\n", a, b, desc)
	r0 := gologoo__compareValLists_71e383cb26f1785bb4841431523136cf(a, b, desc)
	log.Printf("Output: %v\n", r0)
	return r0
}
func compareVals(x, y interface {
}) int {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__compareVals_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v %v\n", x, y)
	r0 := gologoo__compareVals_71e383cb26f1785bb4841431523136cf(x, y)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ec evalContext) colInfo(e spansql.Expr) (colInfo, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__colInfo_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v\n", e)
	r0, r1 := ec.gologoo__colInfo_71e383cb26f1785bb4841431523136cf(e)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (ec evalContext) arithColType(ao spansql.ArithOp) (spansql.Type, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__arithColType_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v\n", ao)
	r0, r1 := ec.gologoo__arithColType_71e383cb26f1785bb4841431523136cf(ao)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func pathExpEqual(a, b spansql.PathExp) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__pathExpEqual_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v %v\n", a, b)
	r0 := gologoo__pathExpEqual_71e383cb26f1785bb4841431523136cf(a, b)
	log.Printf("Output: %v\n", r0)
	return r0
}
func evalLike(str, pat string) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__evalLike_71e383cb26f1785bb4841431523136cf")
	log.Printf("Input : %v %v\n", str, pat)
	r0 := gologoo__evalLike_71e383cb26f1785bb4841431523136cf(str, pat)
	log.Printf("Output: %v\n", r0)
	return r0
}
