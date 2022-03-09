package spannertest

import (
	"fmt"
	"io"
	"sort"
	"cloud.google.com/go/spanner/spansql"
	"log"
)

type rowIter interface {
	Cols() []colInfo
	Next() (row, error)
}
type aggSentinel struct {
	spansql.Expr
	Type     spansql.Type
	AggIndex int
}
type nullIter struct {
	done bool
}

func (ni *nullIter) gologoo__Cols_ba55c87de92a16706542bb88fa5fc839() []colInfo {
	return nil
}
func (ni *nullIter) gologoo__Next_ba55c87de92a16706542bb88fa5fc839() (row, error) {
	if ni.done {
		return nil, io.EOF
	}
	ni.done = true
	return nil, nil
}

type tableIter struct {
	t        *table
	rowIndex int
	alias    spansql.ID
}

func (ti *tableIter) gologoo__Cols_ba55c87de92a16706542bb88fa5fc839() []colInfo {
	cis := make([]colInfo, len(ti.t.cols))
	for _, ci := range ti.t.cols {
		if ti.alias != "" {
			ci.Alias = spansql.PathExp{ti.alias, ci.Name}
		}
		cis[ti.t.origIndex[ci.Name]] = ci
	}
	return cis
}
func (ti *tableIter) gologoo__Next_ba55c87de92a16706542bb88fa5fc839() (row, error) {
	if ti.rowIndex >= len(ti.t.rows) {
		return nil, io.EOF
	}
	r := ti.t.rows[ti.rowIndex]
	ti.rowIndex++
	res := make(row, len(r))
	for i, ci := range ti.t.cols {
		res[ti.t.origIndex[ci.Name]] = r[i]
	}
	return res, nil
}

type rawIter struct {
	cols []colInfo
	rows []row
}

func (raw *rawIter) gologoo__Cols_ba55c87de92a16706542bb88fa5fc839() []colInfo {
	return raw.cols
}
func (raw *rawIter) gologoo__Next_ba55c87de92a16706542bb88fa5fc839() (row, error) {
	if len(raw.rows) == 0 {
		return nil, io.EOF
	}
	res := raw.rows[0]
	raw.rows = raw.rows[1:]
	return res, nil
}
func (raw *rawIter) gologoo__add_ba55c87de92a16706542bb88fa5fc839(src row, colIndexes []int) {
	raw.rows = append(raw.rows, src.copyData(colIndexes))
}
func (raw *rawIter) gologoo__clone_ba55c87de92a16706542bb88fa5fc839() *rawIter {
	return &rawIter{cols: raw.cols, rows: raw.rows}
}
func gologoo__toRawIter_ba55c87de92a16706542bb88fa5fc839(ri rowIter) (*rawIter, error) {
	if raw, ok := ri.(*rawIter); ok {
		return raw, nil
	}
	raw := &rawIter{cols: ri.Cols()}
	for {
		row, err := ri.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		raw.rows = append(raw.rows, row.copyAllData())
	}
	return raw, nil
}

type whereIter struct {
	ri    rowIter
	ec    evalContext
	where spansql.BoolExpr
}

func (wi whereIter) gologoo__Cols_ba55c87de92a16706542bb88fa5fc839() []colInfo {
	return wi.ri.Cols()
}
func (wi whereIter) gologoo__Next_ba55c87de92a16706542bb88fa5fc839() (row, error) {
	for {
		row, err := wi.ri.Next()
		if err != nil {
			return nil, err
		}
		wi.ec.row = row
		b, err := wi.ec.evalBoolExpr(wi.where)
		if err != nil {
			return nil, err
		}
		if b != nil && *b {
			return row, nil
		}
	}
}

type selIter struct {
	ri       rowIter
	ec       evalContext
	cis      []colInfo
	list     []spansql.Expr
	distinct bool
	seen     []row
}

func (si *selIter) gologoo__Cols_ba55c87de92a16706542bb88fa5fc839() []colInfo {
	return si.cis
}
func (si *selIter) gologoo__Next_ba55c87de92a16706542bb88fa5fc839() (row, error) {
	for {
		r, err := si.next()
		if err != nil {
			return nil, err
		}
		if si.distinct && !si.keep(r) {
			continue
		}
		return r, nil
	}
}
func (si *selIter) gologoo__next_ba55c87de92a16706542bb88fa5fc839() (row, error) {
	r, err := si.ri.Next()
	if err != nil {
		return nil, err
	}
	si.ec.row = r
	var out row
	for _, e := range si.list {
		if e == spansql.Star {
			out = append(out, r...)
		} else {
			v, err := si.ec.evalExpr(e)
			if err != nil {
				return nil, err
			}
			out = append(out, v)
		}
	}
	return out, nil
}
func (si *selIter) gologoo__keep_ba55c87de92a16706542bb88fa5fc839(r row) bool {
	for _, prev := range si.seen {
		if rowEqual(prev, r) {
			return false
		}
	}
	si.seen = append(si.seen, r)
	return true
}

type offsetIter struct {
	ri   rowIter
	skip int64
}

func (oi *offsetIter) gologoo__Cols_ba55c87de92a16706542bb88fa5fc839() []colInfo {
	return oi.ri.Cols()
}
func (oi *offsetIter) gologoo__Next_ba55c87de92a16706542bb88fa5fc839() (row, error) {
	for oi.skip > 0 {
		_, err := oi.ri.Next()
		if err != nil {
			return nil, err
		}
		oi.skip--
	}
	row, err := oi.ri.Next()
	if err != nil {
		return nil, err
	}
	return row, nil
}

type limitIter struct {
	ri  rowIter
	rem int64
}

func (li *limitIter) gologoo__Cols_ba55c87de92a16706542bb88fa5fc839() []colInfo {
	return li.ri.Cols()
}
func (li *limitIter) gologoo__Next_ba55c87de92a16706542bb88fa5fc839() (row, error) {
	if li.rem <= 0 {
		return nil, io.EOF
	}
	row, err := li.ri.Next()
	if err != nil {
		return nil, err
	}
	li.rem--
	return row, nil
}

type queryParam struct {
	Value interface {
	}
	Type spansql.Type
}
type queryParams map[string]queryParam
type queryContext struct {
	params     queryParams
	tables     []*table
	tableIndex map[spansql.ID]*table
	locks      int
}

func (qc *queryContext) gologoo__Lock_ba55c87de92a16706542bb88fa5fc839() {
	for _, t := range qc.tables {
		t.mu.Lock()
		qc.locks++
	}
}
func (qc *queryContext) gologoo__Unlock_ba55c87de92a16706542bb88fa5fc839() {
	for _, t := range qc.tables {
		t.mu.Unlock()
		qc.locks--
	}
}
func (d *database) gologoo__Query_ba55c87de92a16706542bb88fa5fc839(q spansql.Query, params queryParams) (ri rowIter, err error) {
	qc, err := d.queryContext(q, params)
	if err != nil {
		return nil, err
	}
	qc.Lock()
	if qc.locks > 0 {
		defer func() {
			if err == nil {
				ri, err = toRawIter(ri)
			}
			qc.Unlock()
		}()
	}
	var aux []spansql.Expr
	var desc []bool
	for _, o := range q.Order {
		aux = append(aux, o.Expr)
		desc = append(desc, o.Desc)
	}
	si, err := d.evalSelect(q.Select, qc)
	if err != nil {
		return nil, err
	}
	ri = si
	if len(q.Order) > 0 {
		rows, keys, err := evalSelectOrder(si, aux)
		if err != nil {
			return nil, err
		}
		sort.Sort(externalRowSorter{rows: rows, keys: keys, desc: desc})
		ri = &rawIter{cols: si.cis, rows: rows}
	}
	if q.Limit != nil {
		if q.Offset != nil {
			off, err := evalLiteralOrParam(q.Offset, params)
			if err != nil {
				return nil, err
			}
			ri = &offsetIter{ri: ri, skip: off}
		}
		lim, err := evalLiteralOrParam(q.Limit, params)
		if err != nil {
			return nil, err
		}
		ri = &limitIter{ri: ri, rem: lim}
	}
	return ri, nil
}
func (d *database) gologoo__queryContext_ba55c87de92a16706542bb88fa5fc839(q spansql.Query, params queryParams) (*queryContext, error) {
	qc := &queryContext{params: params}
	addTable := func(name spansql.ID) error {
		if _, ok := qc.tableIndex[name]; ok {
			return nil
		}
		t, err := d.table(name)
		if err != nil {
			return err
		}
		if qc.tableIndex == nil {
			qc.tableIndex = make(map[spansql.ID]*table)
		}
		qc.tableIndex[name] = t
		return nil
	}
	var findTables func(sf spansql.SelectFrom) error
	findTables = func(sf spansql.SelectFrom) error {
		switch sf := sf.(type) {
		default:
			return fmt.Errorf("can't prepare query context for SelectFrom of type %T", sf)
		case spansql.SelectFromTable:
			return addTable(sf.Table)
		case spansql.SelectFromJoin:
			if err := findTables(sf.LHS); err != nil {
				return err
			}
			return findTables(sf.RHS)
		case spansql.SelectFromUnnest:
			return nil
		}
	}
	for _, sf := range q.Select.From {
		if err := findTables(sf); err != nil {
			return nil, err
		}
	}
	var names []spansql.ID
	for name := range qc.tableIndex {
		names = append(names, name)
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})
	for _, name := range names {
		qc.tables = append(qc.tables, qc.tableIndex[name])
	}
	return qc, nil
}
func (d *database) gologoo__evalSelect_ba55c87de92a16706542bb88fa5fc839(sel spansql.Select, qc *queryContext) (si *selIter, evalErr error) {
	var ri rowIter = &nullIter{}
	ec := evalContext{params: qc.params}
	if len(sel.From) > 1 {
		return nil, fmt.Errorf("selecting with more than one FROM clause not yet supported")
	}
	if len(sel.From) == 1 {
		var err error
		ec, ri, err = d.evalSelectFrom(qc, ec, sel.From[0])
		if err != nil {
			return nil, err
		}
	}
	if sel.Where != nil {
		ri = whereIter{ri: ri, ec: ec, where: sel.Where}
	}
	ec.aliases = make(map[spansql.ID]spansql.Expr)
	for i, alias := range sel.ListAliases {
		ec.aliases[alias] = sel.List[i]
	}
	var rowGroups [][2]int
	if len(sel.GroupBy) > 0 {
		raw, err := toRawIter(ri)
		if err != nil {
			return nil, err
		}
		keys := make([][]interface {
		}, 0, len(raw.rows))
		for _, row := range raw.rows {
			ec.row = row
			key, err := ec.evalExprList(sel.GroupBy)
			if err != nil {
				return nil, err
			}
			keys = append(keys, key)
		}
		ers := externalRowSorter{rows: raw.rows, keys: keys}
		sort.Sort(ers)
		raw.rows = ers.rows
		start := 0
		for i := 1; i < len(keys); i++ {
			if compareValLists(keys[i-1], keys[i], nil) == 0 {
				continue
			}
			rowGroups = append(rowGroups, [2]int{start, i})
			start = i
		}
		if len(keys) > 0 {
			rowGroups = append(rowGroups, [2]int{start, len(keys)})
		}
		ec.aliases = nil
		ri = raw
	}
	var aggI []int
	for i, e := range sel.List {
		f, ok := e.(spansql.Func)
		if !ok || len(f.Args) != 1 {
			continue
		}
		_, ok = aggregateFuncs[f.Name]
		if !ok {
			continue
		}
		aggI = append(aggI, i)
	}
	if len(aggI) > 0 || len(sel.GroupBy) > 0 {
		raw, err := toRawIter(ri)
		if err != nil {
			return nil, err
		}
		if len(sel.GroupBy) == 0 {
			rowGroups = [][2]int{{0, len(raw.rows)}}
		}
		rawOut := &rawIter{cols: append([]colInfo(nil), raw.cols...)}
		aggType := make([]*spansql.Type, len(aggI))
		for _, rg := range rowGroups {
			var outRow row
			if rg[0] < len(raw.rows) {
				repRow := raw.rows[rg[0]]
				for i := range repRow {
					outRow = append(outRow, repRow.copyDataElem(i))
				}
			} else {
				for i := 0; i < len(rawOut.cols); i++ {
					outRow = append(outRow, nil)
				}
			}
			for j, aggI := range aggI {
				fexpr := sel.List[aggI].(spansql.Func)
				fn := aggregateFuncs[fexpr.Name]
				starArg := fexpr.Args[0] == spansql.Star
				if starArg && !fn.AcceptStar {
					return nil, fmt.Errorf("aggregate function %s does not accept * as an argument", fexpr.Name)
				}
				var argType spansql.Type
				if !starArg {
					ci, err := ec.colInfo(fexpr.Args[0])
					if err != nil {
						return nil, fmt.Errorf("evaluating aggregate function %s arg type: %v", fexpr.Name, err)
					}
					argType = ci.Type
				}
				var values []interface {
				}
				for i := rg[0]; i < rg[1]; i++ {
					ec.row = raw.rows[i]
					if starArg {
						values = append(values, 1)
					} else {
						x, err := ec.evalExpr(fexpr.Args[0])
						if err != nil {
							return nil, err
						}
						values = append(values, x)
					}
				}
				x, typ, err := fn.Eval(values, argType)
				if err != nil {
					return nil, err
				}
				aggType[j] = &typ
				outRow = append(outRow, x)
			}
			rawOut.rows = append(rawOut.rows, outRow)
		}
		for j, aggI := range aggI {
			fexpr := sel.List[aggI].(spansql.Func)
			if aggType[j] == nil {
				aggType[j] = &int64Type
			}
			rawOut.cols = append(rawOut.cols, colInfo{Name: spansql.ID(fexpr.SQL()), Type: *aggType[j], AggIndex: aggI + 1})
			sel.List[aggI] = aggSentinel{Type: *aggType[j], AggIndex: aggI + 1}
		}
		ri = rawOut
		ec.cols = rawOut.cols
	}
	var colInfos []colInfo
	for i, e := range sel.List {
		if e == spansql.Star {
			colInfos = append(colInfos, ec.cols...)
		} else {
			ci, err := ec.colInfo(e)
			if err != nil {
				return nil, err
			}
			if len(sel.ListAliases) > 0 {
				alias := sel.ListAliases[i]
				if alias != "" {
					ci.Name = alias
				}
			}
			colInfos = append(colInfos, ci)
		}
	}
	return &selIter{ri: ri, ec: ec, cis: colInfos, list: sel.List, distinct: sel.Distinct}, nil
}
func (d *database) gologoo__evalSelectFrom_ba55c87de92a16706542bb88fa5fc839(qc *queryContext, ec evalContext, sf spansql.SelectFrom) (evalContext, rowIter, error) {
	switch sf := sf.(type) {
	default:
		return ec, nil, fmt.Errorf("selecting with FROM clause of type %T not yet supported", sf)
	case spansql.SelectFromTable:
		t, ok := qc.tableIndex[sf.Table]
		if !ok {
			return ec, nil, fmt.Errorf("unknown table %q", sf.Table)
		}
		ti := &tableIter{t: t}
		if sf.Alias != "" {
			ti.alias = sf.Alias
		} else {
			ti.alias = sf.Table
		}
		ec.cols = ti.Cols()
		return ec, ti, nil
	case spansql.SelectFromJoin:
		lhsEC, lhs, err := d.evalSelectFrom(qc, ec, sf.LHS)
		if err != nil {
			return ec, nil, err
		}
		lhsRaw, err := toRawIter(lhs)
		if err != nil {
			return ec, nil, err
		}
		rhsEC, rhs, err := d.evalSelectFrom(qc, ec, sf.RHS)
		if err != nil {
			return ec, nil, err
		}
		rhsRaw, err := toRawIter(rhs)
		if err != nil {
			return ec, nil, err
		}
		ji, ec, err := newJoinIter(lhsRaw, rhsRaw, lhsEC, rhsEC, sf)
		if err != nil {
			return ec, nil, err
		}
		return ec, ji, nil
	case spansql.SelectFromUnnest:
		col, err := ec.colInfo(sf.Expr)
		if err != nil {
			return ec, nil, fmt.Errorf("evaluating type of UNNEST arg: %v", err)
		}
		if !col.Type.Array {
			return ec, nil, fmt.Errorf("type of UNNEST arg is non-array %s", col.Type.SQL())
		}
		col.Name = sf.Alias
		col.Type.Array = false
		e, err := ec.evalExpr(sf.Expr)
		if err != nil {
			return ec, nil, fmt.Errorf("evaluating UNNEST arg: %v", err)
		}
		arr, ok := e.([]interface {
		})
		if !ok {
			return ec, nil, fmt.Errorf("evaluating UNNEST arg gave %t, want array", e)
		}
		var rows []row
		for _, v := range arr {
			rows = append(rows, row{v})
		}
		ri := &rawIter{cols: []colInfo{col}, rows: rows}
		ec.cols = ri.cols
		return ec, ri, nil
	}
}
func gologoo__newJoinIter_ba55c87de92a16706542bb88fa5fc839(lhs, rhs *rawIter, lhsEC, rhsEC evalContext, sfj spansql.SelectFromJoin) (*joinIter, evalContext, error) {
	if sfj.On != nil && len(sfj.Using) > 0 {
		return nil, evalContext{}, fmt.Errorf("JOIN may not have both ON and USING clauses")
	}
	if sfj.On == nil && len(sfj.Using) == 0 && sfj.Type != spansql.CrossJoin {
		return nil, evalContext{}, fmt.Errorf("non-CROSS JOIN must have ON or USING clause")
	}
	ji := &joinIter{jt: sfj.Type, ec: lhsEC, primary: lhs, secondaryOrig: rhs, primaryOffset: 0, secondaryOffset: len(lhsEC.cols)}
	switch ji.jt {
	case spansql.LeftJoin:
		ji.nullPad = true
	case spansql.RightJoin:
		ji.nullPad = true
		ji.ec = rhsEC
		ji.primary, ji.secondaryOrig = rhs, lhs
		ji.primaryOffset, ji.secondaryOffset = len(rhsEC.cols), 0
	case spansql.FullJoin:
		ji.nullPad = true
		ji.used = make([]bool, 0, 10)
	}
	ji.ec.cols, ji.ec.row = nil, nil
	if len(sfj.Using) == 0 {
		ji.prepNonUsing(sfj.On, lhsEC, rhsEC)
	} else {
		if err := ji.prepUsing(sfj.Using, lhsEC, rhsEC, ji.jt == spansql.RightJoin); err != nil {
			return nil, evalContext{}, err
		}
	}
	return ji, ji.ec, nil
}
func (ji *joinIter) gologoo__prepNonUsing_ba55c87de92a16706542bb88fa5fc839(on spansql.BoolExpr, lhsEC, rhsEC evalContext) {
	ji.ec.cols = append(ji.ec.cols, lhsEC.cols...)
	ji.ec.cols = append(ji.ec.cols, rhsEC.cols...)
	ji.ec.row = make(row, len(ji.ec.cols))
	ji.cond = func(primary, secondary row) (bool, error) {
		copy(ji.ec.row[ji.primaryOffset:], primary)
		copy(ji.ec.row[ji.secondaryOffset:], secondary)
		if on == nil {
			return true, nil
		}
		b, err := ji.ec.evalBoolExpr(on)
		if err != nil {
			return false, err
		}
		return b != nil && *b, nil
	}
	ji.zero = func(primary, secondary row) {
		for i := range ji.ec.row {
			ji.ec.row[i] = nil
		}
		copy(ji.ec.row[ji.primaryOffset:], primary)
		copy(ji.ec.row[ji.secondaryOffset:], secondary)
	}
}
func (ji *joinIter) gologoo__prepUsing_ba55c87de92a16706542bb88fa5fc839(using []spansql.ID, lhsEC, rhsEC evalContext, flipped bool) error {
	var lhsUsing, rhsUsing []int
	var lhsNotUsing, rhsNotUsing []int
	lhsUsed, rhsUsed := make(map[int]bool), make(map[int]bool)
	for _, id := range using {
		lhsi, err := lhsEC.resolveColumnIndex(id)
		if err != nil {
			return err
		}
		lhsUsing = append(lhsUsing, lhsi)
		lhsUsed[lhsi] = true
		rhsi, err := rhsEC.resolveColumnIndex(id)
		if err != nil {
			return err
		}
		rhsUsing = append(rhsUsing, rhsi)
		rhsUsed[rhsi] = true
		ji.ec.cols = append(ji.ec.cols, lhsEC.cols[lhsi])
	}
	for i, col := range lhsEC.cols {
		if !lhsUsed[i] {
			ji.ec.cols = append(ji.ec.cols, col)
			lhsNotUsing = append(lhsNotUsing, i)
		}
	}
	for i, col := range rhsEC.cols {
		if !rhsUsed[i] {
			ji.ec.cols = append(ji.ec.cols, col)
			rhsNotUsing = append(rhsNotUsing, i)
		}
	}
	ji.ec.row = make(row, len(ji.ec.cols))
	primaryUsing, secondaryUsing := lhsUsing, rhsUsing
	if flipped {
		primaryUsing, secondaryUsing = secondaryUsing, primaryUsing
	}
	orNil := func(r row, i int) interface {
	} {
		if r == nil {
			return nil
		}
		return r[i]
	}
	populate := func(primary, secondary row) {
		j := 0
		if primary != nil {
			for _, pi := range primaryUsing {
				ji.ec.row[j] = primary[pi]
				j++
			}
		} else {
			for _, si := range secondaryUsing {
				ji.ec.row[j] = secondary[si]
				j++
			}
		}
		lhs, rhs := primary, secondary
		if flipped {
			rhs, lhs = lhs, rhs
		}
		for _, i := range lhsNotUsing {
			ji.ec.row[j] = orNil(lhs, i)
			j++
		}
		for _, i := range rhsNotUsing {
			ji.ec.row[j] = orNil(rhs, i)
			j++
		}
	}
	ji.cond = func(primary, secondary row) (bool, error) {
		for i, pi := range primaryUsing {
			si := secondaryUsing[i]
			if compareVals(primary[pi], secondary[si]) != 0 {
				return false, nil
			}
		}
		populate(primary, secondary)
		return true, nil
	}
	ji.zero = func(primary, secondary row) {
		populate(primary, secondary)
	}
	return nil
}

type joinIter struct {
	jt                             spansql.JoinType
	ec                             evalContext
	primary, secondaryOrig         *rawIter
	primaryOffset, secondaryOffset int
	nullPad                        bool
	primaryRow                     row
	secondary                      *rawIter
	secondaryRead                  int
	any                            bool
	cond                           func(primary, secondary row) (bool, error)
	zero                           func(primary, secondary row)
	used                           []bool
	zeroUnused                     bool
	unusedIndex                    int
}

func (ji *joinIter) gologoo__Cols_ba55c87de92a16706542bb88fa5fc839() []colInfo {
	return ji.ec.cols
}
func (ji *joinIter) gologoo__nextPrimary_ba55c87de92a16706542bb88fa5fc839() error {
	var err error
	ji.primaryRow, err = ji.primary.Next()
	if err != nil {
		return err
	}
	ji.secondary = ji.secondaryOrig.clone()
	ji.secondaryRead = 0
	ji.any = false
	return nil
}
func (ji *joinIter) gologoo__Next_ba55c87de92a16706542bb88fa5fc839() (row, error) {
	if ji.primaryRow == nil && !ji.zeroUnused {
		err := ji.nextPrimary()
		if err == io.EOF && ji.used != nil {
			ji.zeroUnused = true
			ji.secondary = nil
			goto scanJiUsed
		}
		if err != nil {
			return nil, err
		}
	}
scanJiUsed:
	if ji.zeroUnused {
		if ji.secondary == nil {
			ji.secondary = ji.secondaryOrig.clone()
			ji.secondaryRead = 0
		}
		for ji.unusedIndex < len(ji.used) && ji.used[ji.unusedIndex] {
			ji.unusedIndex++
		}
		if ji.unusedIndex >= len(ji.used) || ji.secondaryRead >= len(ji.used) {
			return nil, io.EOF
		}
		var secondaryRow row
		for ji.secondaryRead <= ji.unusedIndex {
			var err error
			secondaryRow, err = ji.secondary.Next()
			if err != nil {
				return nil, err
			}
			ji.secondaryRead++
		}
		ji.zero(nil, secondaryRow)
		return ji.ec.row, nil
	}
	for {
		secondaryRow, err := ji.secondary.Next()
		if err == io.EOF {
			if !ji.any && ji.nullPad {
				ji.zero(ji.primaryRow, nil)
				ji.primaryRow = nil
				return ji.ec.row, nil
			}
			err := ji.nextPrimary()
			if err == io.EOF && ji.used != nil {
				ji.zeroUnused = true
				ji.secondary = nil
				goto scanJiUsed
			}
			if err != nil {
				return nil, err
			}
			continue
		}
		if err != nil {
			return nil, err
		}
		ji.secondaryRead++
		if ji.used != nil {
			for len(ji.used) < ji.secondaryRead {
				ji.used = append(ji.used, false)
			}
		}
		match, err := ji.cond(ji.primaryRow, secondaryRow)
		if err != nil {
			return nil, err
		}
		if !match {
			continue
		}
		ji.any = true
		if ji.used != nil {
			ji.used[ji.secondaryRead-1] = true
		}
		return ji.ec.row, nil
	}
}
func gologoo__evalSelectOrder_ba55c87de92a16706542bb88fa5fc839(si *selIter, aux []spansql.Expr) (rows []row, keys [][]interface {
}, err error) {
	for {
		r, err := si.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, nil, err
		}
		key, err := si.ec.evalExprList(aux)
		if err != nil {
			return nil, nil, err
		}
		rows = append(rows, r.copyAllData())
		keys = append(keys, key)
	}
	return
}

type externalRowSorter struct {
	rows []row
	keys [][]interface {
	}
	desc []bool
}

func (ers externalRowSorter) gologoo__Len_ba55c87de92a16706542bb88fa5fc839() int {
	return len(ers.rows)
}
func (ers externalRowSorter) gologoo__Less_ba55c87de92a16706542bb88fa5fc839(i, j int) bool {
	return compareValLists(ers.keys[i], ers.keys[j], ers.desc) < 0
}
func (ers externalRowSorter) gologoo__Swap_ba55c87de92a16706542bb88fa5fc839(i, j int) {
	ers.rows[i], ers.rows[j] = ers.rows[j], ers.rows[i]
	ers.keys[i], ers.keys[j] = ers.keys[j], ers.keys[i]
}
func (ni *nullIter) Cols() []colInfo {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Cols_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0 := ni.gologoo__Cols_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ni *nullIter) Next() (row, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Next_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0, r1 := ni.gologoo__Next_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (ti *tableIter) Cols() []colInfo {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Cols_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0 := ti.gologoo__Cols_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ti *tableIter) Next() (row, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Next_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0, r1 := ti.gologoo__Next_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (raw *rawIter) Cols() []colInfo {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Cols_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0 := raw.gologoo__Cols_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (raw *rawIter) Next() (row, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Next_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0, r1 := raw.gologoo__Next_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (raw *rawIter) add(src row, colIndexes []int) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__add_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : %v %v\n", src, colIndexes)
	raw.gologoo__add_ba55c87de92a16706542bb88fa5fc839(src, colIndexes)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (raw *rawIter) clone() *rawIter {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__clone_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0 := raw.gologoo__clone_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v\n", r0)
	return r0
}
func toRawIter(ri rowIter) (*rawIter, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__toRawIter_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : %v\n", ri)
	r0, r1 := gologoo__toRawIter_ba55c87de92a16706542bb88fa5fc839(ri)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (wi whereIter) Cols() []colInfo {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Cols_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0 := wi.gologoo__Cols_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (wi whereIter) Next() (row, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Next_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0, r1 := wi.gologoo__Next_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (si *selIter) Cols() []colInfo {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Cols_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0 := si.gologoo__Cols_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (si *selIter) Next() (row, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Next_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0, r1 := si.gologoo__Next_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (si *selIter) next() (row, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__next_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0, r1 := si.gologoo__next_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (si *selIter) keep(r row) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__keep_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : %v\n", r)
	r0 := si.gologoo__keep_ba55c87de92a16706542bb88fa5fc839(r)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (oi *offsetIter) Cols() []colInfo {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Cols_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0 := oi.gologoo__Cols_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (oi *offsetIter) Next() (row, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Next_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0, r1 := oi.gologoo__Next_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (li *limitIter) Cols() []colInfo {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Cols_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0 := li.gologoo__Cols_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (li *limitIter) Next() (row, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Next_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0, r1 := li.gologoo__Next_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (qc *queryContext) Lock() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Lock_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	qc.gologoo__Lock_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (qc *queryContext) Unlock() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Unlock_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	qc.gologoo__Unlock_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (d *database) Query(q spansql.Query, params queryParams) (ri rowIter, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Query_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : %v %v\n", q, params)
	ri, err = d.gologoo__Query_ba55c87de92a16706542bb88fa5fc839(q, params)
	log.Printf("Output: %v %v\n", ri, err)
	return
}
func (d *database) queryContext(q spansql.Query, params queryParams) (*queryContext, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__queryContext_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : %v %v\n", q, params)
	r0, r1 := d.gologoo__queryContext_ba55c87de92a16706542bb88fa5fc839(q, params)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (d *database) evalSelect(sel spansql.Select, qc *queryContext) (si *selIter, evalErr error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__evalSelect_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : %v %v\n", sel, qc)
	si, evalErr = d.gologoo__evalSelect_ba55c87de92a16706542bb88fa5fc839(sel, qc)
	log.Printf("Output: %v %v\n", si, evalErr)
	return
}
func (d *database) evalSelectFrom(qc *queryContext, ec evalContext, sf spansql.SelectFrom) (evalContext, rowIter, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__evalSelectFrom_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : %v %v %v\n", qc, ec, sf)
	r0, r1, r2 := d.gologoo__evalSelectFrom_ba55c87de92a16706542bb88fa5fc839(qc, ec, sf)
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
func newJoinIter(lhs, rhs *rawIter, lhsEC, rhsEC evalContext, sfj spansql.SelectFromJoin) (*joinIter, evalContext, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__newJoinIter_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : %v %v %v %v %v\n", lhs, rhs, lhsEC, rhsEC, sfj)
	r0, r1, r2 := gologoo__newJoinIter_ba55c87de92a16706542bb88fa5fc839(lhs, rhs, lhsEC, rhsEC, sfj)
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
func (ji *joinIter) prepNonUsing(on spansql.BoolExpr, lhsEC, rhsEC evalContext) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__prepNonUsing_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : %v %v %v\n", on, lhsEC, rhsEC)
	ji.gologoo__prepNonUsing_ba55c87de92a16706542bb88fa5fc839(on, lhsEC, rhsEC)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (ji *joinIter) prepUsing(using []spansql.ID, lhsEC, rhsEC evalContext, flipped bool) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__prepUsing_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : %v %v %v %v\n", using, lhsEC, rhsEC, flipped)
	r0 := ji.gologoo__prepUsing_ba55c87de92a16706542bb88fa5fc839(using, lhsEC, rhsEC, flipped)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ji *joinIter) Cols() []colInfo {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Cols_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0 := ji.gologoo__Cols_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ji *joinIter) nextPrimary() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__nextPrimary_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0 := ji.gologoo__nextPrimary_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ji *joinIter) Next() (row, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Next_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0, r1 := ji.gologoo__Next_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func evalSelectOrder(si *selIter, aux []spansql.Expr) (rows []row, keys [][]interface {
}, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__evalSelectOrder_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : %v %v\n", si, aux)
	rows, keys, err = gologoo__evalSelectOrder_ba55c87de92a16706542bb88fa5fc839(si, aux)
	log.Printf("Output: %v %v %v\n", rows, keys, err)
	return
}
func (ers externalRowSorter) Len() int {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Len_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : (none)\n")
	r0 := ers.gologoo__Len_ba55c87de92a16706542bb88fa5fc839()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ers externalRowSorter) Less(i, j int) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Less_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : %v %v\n", i, j)
	r0 := ers.gologoo__Less_ba55c87de92a16706542bb88fa5fc839(i, j)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ers externalRowSorter) Swap(i, j int) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Swap_ba55c87de92a16706542bb88fa5fc839")
	log.Printf("Input : %v %v\n", i, j)
	ers.gologoo__Swap_ba55c87de92a16706542bb88fa5fc839(i, j)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
