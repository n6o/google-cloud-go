package spannertest

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"cloud.google.com/go/civil"
	"cloud.google.com/go/spanner/spansql"
	"log"
)

type database struct {
	mu      sync.Mutex
	lastTS  time.Time
	tables  map[spansql.ID]*table
	indexes map[spansql.ID]struct {
	}
	views map[spansql.ID]struct {
	}
	rwMu sync.Mutex
}
type table struct {
	mu        sync.Mutex
	cols      []colInfo
	colIndex  map[spansql.ID]int
	origIndex map[spansql.ID]int
	pkCols    int
	pkDesc    []bool
	rdw       *spansql.RowDeletionPolicy
	rows      []row
}
type colInfo struct {
	Name      spansql.ID
	Type      spansql.Type
	Generated spansql.Expr
	NotNull   bool
	AggIndex  int
	Alias     spansql.PathExp
}

var commitTimestampSentinel = &struct {
}{}

type transaction struct {
	readOnly        bool
	d               *database
	commitTimestamp time.Time
	unlock          func()
}

func (d *database) gologoo__NewReadOnlyTransaction_48b215e0cf779586a3b42f53798ab710() *transaction {
	return &transaction{readOnly: true}
}
func (d *database) gologoo__NewTransaction_48b215e0cf779586a3b42f53798ab710() *transaction {
	return &transaction{d: d}
}
func (tx *transaction) gologoo__Start_48b215e0cf779586a3b42f53798ab710() {
	tx.d.rwMu.Lock()
	tx.d.mu.Lock()
	const tsRes = 1 * time.Microsecond
	now := time.Now().UTC().Truncate(tsRes)
	if !now.After(tx.d.lastTS) {
		now = tx.d.lastTS.Add(tsRes)
	}
	tx.d.lastTS = now
	tx.d.mu.Unlock()
	tx.commitTimestamp = now
	tx.unlock = tx.d.rwMu.Unlock
}
func (tx *transaction) gologoo__checkMutable_48b215e0cf779586a3b42f53798ab710() error {
	if tx.readOnly {
		return status.Errorf(codes.InvalidArgument, "transaction is read-only")
	}
	return nil
}
func (tx *transaction) gologoo__Commit_48b215e0cf779586a3b42f53798ab710() (time.Time, error) {
	if tx.unlock != nil {
		tx.unlock()
	}
	return tx.commitTimestamp, nil
}
func (tx *transaction) gologoo__Rollback_48b215e0cf779586a3b42f53798ab710() {
	if tx.unlock != nil {
		tx.unlock()
	}
}

type row []interface {
}

func (r row) gologoo__copyDataElem_48b215e0cf779586a3b42f53798ab710(index int) interface {
} {
	v := r[index]
	if is, ok := v.([]interface {
	}); ok {
		v = append([]interface {
		}(nil), is...)
	}
	return v
}
func (r row) gologoo__copyAllData_48b215e0cf779586a3b42f53798ab710() row {
	dst := make(row, 0, len(r))
	for i := range r {
		dst = append(dst, r.copyDataElem(i))
	}
	return dst
}
func (r row) gologoo__copyData_48b215e0cf779586a3b42f53798ab710(indexes []int) row {
	if len(indexes) == 0 {
		return nil
	}
	dst := make(row, 0, len(indexes))
	for _, i := range indexes {
		dst = append(dst, r.copyDataElem(i))
	}
	return dst
}
func (d *database) gologoo__LastCommitTimestamp_48b215e0cf779586a3b42f53798ab710() time.Time {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.lastTS
}
func (d *database) gologoo__GetDDL_48b215e0cf779586a3b42f53798ab710() []spansql.DDLStmt {
	d.mu.Lock()
	defer d.mu.Unlock()
	var stmts []spansql.DDLStmt
	for name, t := range d.tables {
		ct := &spansql.CreateTable{Name: name}
		t.mu.Lock()
		for i, col := range t.cols {
			ct.Columns = append(ct.Columns, spansql.ColumnDef{Name: col.Name, Type: col.Type, NotNull: col.NotNull})
			if i < t.pkCols {
				ct.PrimaryKey = append(ct.PrimaryKey, spansql.KeyPart{Column: col.Name, Desc: t.pkDesc[i]})
			}
		}
		t.mu.Unlock()
		stmts = append(stmts, ct)
	}
	return stmts
}
func (d *database) gologoo__ApplyDDL_48b215e0cf779586a3b42f53798ab710(stmt spansql.DDLStmt) *status.Status {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.tables == nil {
		d.tables = make(map[spansql.ID]*table)
	}
	if d.indexes == nil {
		d.indexes = make(map[spansql.ID]struct {
		})
	}
	if d.views == nil {
		d.views = make(map[spansql.ID]struct {
		})
	}
	switch stmt := stmt.(type) {
	default:
		return status.Newf(codes.Unimplemented, "unhandled DDL statement type %T", stmt)
	case *spansql.CreateTable:
		if _, ok := d.tables[stmt.Name]; ok {
			return status.Newf(codes.AlreadyExists, "table %s already exists", stmt.Name)
		}
		if len(stmt.PrimaryKey) == 0 {
			return status.Newf(codes.InvalidArgument, "table %s has no primary key", stmt.Name)
		}
		orig := make(map[spansql.ID]int)
		for i, col := range stmt.Columns {
			orig[col.Name] = i
		}
		pk := make(map[spansql.ID]int)
		var pkDesc []bool
		for i, kp := range stmt.PrimaryKey {
			pk[kp.Column] = -1000 + i
			pkDesc = append(pkDesc, kp.Desc)
		}
		sort.SliceStable(stmt.Columns, func(i, j int) bool {
			a, b := pk[stmt.Columns[i].Name], pk[stmt.Columns[j].Name]
			return a < b
		})
		t := &table{colIndex: make(map[spansql.ID]int), origIndex: orig, pkCols: len(pk), pkDesc: pkDesc}
		for _, cd := range stmt.Columns {
			if st := t.addColumn(cd, true); st.Code() != codes.OK {
				return st
			}
		}
		for col := range pk {
			if _, ok := t.colIndex[col]; !ok {
				return status.Newf(codes.InvalidArgument, "primary key column %q not in table", col)
			}
		}
		t.rdw = stmt.RowDeletionPolicy
		d.tables[stmt.Name] = t
		return nil
	case *spansql.CreateIndex:
		if _, ok := d.indexes[stmt.Name]; ok {
			return status.Newf(codes.AlreadyExists, "index %s already exists", stmt.Name)
		}
		d.indexes[stmt.Name] = struct {
		}{}
		return nil
	case *spansql.CreateView:
		if !stmt.OrReplace {
			if _, ok := d.views[stmt.Name]; ok {
				return status.Newf(codes.AlreadyExists, "view %s already exists", stmt.Name)
			}
		}
		d.views[stmt.Name] = struct {
		}{}
		return nil
	case *spansql.DropTable:
		if _, ok := d.tables[stmt.Name]; !ok {
			return status.Newf(codes.NotFound, "no table named %s", stmt.Name)
		}
		delete(d.tables, stmt.Name)
		return nil
	case *spansql.DropIndex:
		if _, ok := d.indexes[stmt.Name]; !ok {
			return status.Newf(codes.NotFound, "no index named %s", stmt.Name)
		}
		delete(d.indexes, stmt.Name)
		return nil
	case *spansql.DropView:
		if _, ok := d.views[stmt.Name]; !ok {
			return status.Newf(codes.NotFound, "no view named %s", stmt.Name)
		}
		delete(d.views, stmt.Name)
		return nil
	case *spansql.AlterTable:
		t, ok := d.tables[stmt.Name]
		if !ok {
			return status.Newf(codes.NotFound, "no table named %s", stmt.Name)
		}
		switch alt := stmt.Alteration.(type) {
		default:
			return status.Newf(codes.Unimplemented, "unhandled DDL table alteration type %T", alt)
		case spansql.AddColumn:
			if st := t.addColumn(alt.Def, false); st.Code() != codes.OK {
				return st
			}
			return nil
		case spansql.DropColumn:
			if st := t.dropColumn(alt.Name); st.Code() != codes.OK {
				return st
			}
			return nil
		case spansql.AlterColumn:
			if st := t.alterColumn(alt); st.Code() != codes.OK {
				return st
			}
			return nil
		case spansql.AddRowDeletionPolicy:
			if st := t.addRowDeletionPolicy(alt); st.Code() != codes.OK {
				return st
			}
			return nil
		case spansql.ReplaceRowDeletionPolicy:
			if st := t.replaceRowDeletionPolicy(alt); st.Code() != codes.OK {
				return st
			}
			return nil
		case spansql.DropRowDeletionPolicy:
			if st := t.dropRowDeletionPolicy(alt); st.Code() != codes.OK {
				return st
			}
			return nil
		}
	}
}
func (d *database) gologoo__table_48b215e0cf779586a3b42f53798ab710(tbl spansql.ID) (*table, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	t, ok := d.tables[tbl]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "no table named %s", tbl)
	}
	return t, nil
}
func (d *database) gologoo__writeValues_48b215e0cf779586a3b42f53798ab710(tx *transaction, tbl spansql.ID, cols []spansql.ID, values []*structpb.ListValue, f func(t *table, colIndexes []int, r row) error) error {
	if err := tx.checkMutable(); err != nil {
		return err
	}
	t, err := d.table(tbl)
	if err != nil {
		return err
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	colIndexes, err := t.colIndexes(cols)
	if err != nil {
		return err
	}
	revIndex := make(map[int]int)
	for j, i := range colIndexes {
		revIndex[i] = j
	}
	for pki := 0; pki < t.pkCols; pki++ {
		_, ok := revIndex[pki]
		if !ok {
			return status.Errorf(codes.InvalidArgument, "primary key column %s not included in write", t.cols[pki].Name)
		}
	}
	for _, vs := range values {
		if len(vs.Values) != len(colIndexes) {
			return status.Errorf(codes.InvalidArgument, "row of %d values can't be written to %d columns", len(vs.Values), len(colIndexes))
		}
		r := make(row, len(t.cols))
		for j, v := range vs.Values {
			i := colIndexes[j]
			if t.cols[i].Generated != nil {
				return status.Error(codes.InvalidArgument, "values can't be written to a generated column")
			}
			x, err := valForType(v, t.cols[i].Type)
			if err != nil {
				return err
			}
			if x == commitTimestampSentinel {
				x = tx.commitTimestamp
			}
			if x == nil && t.cols[i].NotNull {
				return status.Errorf(codes.FailedPrecondition, "%s must not be NULL in table %s", t.cols[i].Name, tbl)
			}
			r[i] = x
		}
		if err := f(t, colIndexes, r); err != nil {
			return err
		}
		pk := r[:t.pkCols]
		rowNum, found := t.rowForPK(pk)
		if !found {
			return status.Error(codes.Internal, "row failed to be inserted")
		}
		row := t.rows[rowNum]
		ec := evalContext{cols: t.cols, row: row}
		for i, col := range t.cols {
			if col.Generated != nil {
				res, err := ec.evalExpr(col.Generated)
				if err != nil {
					return err
				}
				row[i] = res
			}
		}
	}
	return nil
}
func (d *database) gologoo__Insert_48b215e0cf779586a3b42f53798ab710(tx *transaction, tbl spansql.ID, cols []spansql.ID, values []*structpb.ListValue) error {
	return d.writeValues(tx, tbl, cols, values, func(t *table, colIndexes []int, r row) error {
		pk := r[:t.pkCols]
		rowNum, found := t.rowForPK(pk)
		if found {
			return status.Errorf(codes.AlreadyExists, "row already in table")
		}
		t.insertRow(rowNum, r)
		return nil
	})
}
func (d *database) gologoo__Update_48b215e0cf779586a3b42f53798ab710(tx *transaction, tbl spansql.ID, cols []spansql.ID, values []*structpb.ListValue) error {
	return d.writeValues(tx, tbl, cols, values, func(t *table, colIndexes []int, r row) error {
		if t.pkCols == 0 {
			return status.Errorf(codes.InvalidArgument, "cannot update table %s with no columns in primary key", tbl)
		}
		pk := r[:t.pkCols]
		rowNum, found := t.rowForPK(pk)
		if !found {
			return status.Errorf(codes.NotFound, "row not in table")
		}
		for _, i := range colIndexes {
			t.rows[rowNum][i] = r[i]
		}
		return nil
	})
}
func (d *database) gologoo__InsertOrUpdate_48b215e0cf779586a3b42f53798ab710(tx *transaction, tbl spansql.ID, cols []spansql.ID, values []*structpb.ListValue) error {
	return d.writeValues(tx, tbl, cols, values, func(t *table, colIndexes []int, r row) error {
		pk := r[:t.pkCols]
		rowNum, found := t.rowForPK(pk)
		if !found {
			t.insertRow(rowNum, r)
		} else {
			for _, i := range colIndexes {
				t.rows[rowNum][i] = r[i]
			}
		}
		return nil
	})
}
func (d *database) gologoo__Delete_48b215e0cf779586a3b42f53798ab710(tx *transaction, table spansql.ID, keys []*structpb.ListValue, keyRanges keyRangeList, all bool) error {
	if err := tx.checkMutable(); err != nil {
		return err
	}
	t, err := d.table(table)
	if err != nil {
		return err
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if all {
		t.rows = nil
		return nil
	}
	for _, key := range keys {
		pk, err := t.primaryKey(key.Values)
		if err != nil {
			return err
		}
		rowNum, found := t.rowForPK(pk)
		if found {
			copy(t.rows[rowNum:], t.rows[rowNum+1:])
			t.rows = t.rows[:len(t.rows)-1]
		}
	}
	for _, r := range keyRanges {
		r.startKey, err = t.primaryKeyPrefix(r.start.Values)
		if err != nil {
			return err
		}
		r.endKey, err = t.primaryKeyPrefix(r.end.Values)
		if err != nil {
			return err
		}
		startRow, endRow := t.findRange(r)
		if n := endRow - startRow; n > 0 {
			copy(t.rows[startRow:], t.rows[endRow:])
			t.rows = t.rows[:len(t.rows)-n]
		}
	}
	return nil
}
func (d *database) gologoo__readTable_48b215e0cf779586a3b42f53798ab710(table spansql.ID, cols []spansql.ID, f func(*table, *rawIter, []int) error) (*rawIter, error) {
	t, err := d.table(table)
	if err != nil {
		return nil, err
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	colIndexes, err := t.colIndexes(cols)
	if err != nil {
		return nil, err
	}
	ri := &rawIter{}
	for _, i := range colIndexes {
		ri.cols = append(ri.cols, t.cols[i])
	}
	return ri, f(t, ri, colIndexes)
}
func (d *database) gologoo__Read_48b215e0cf779586a3b42f53798ab710(tbl spansql.ID, cols []spansql.ID, keys []*structpb.ListValue, keyRanges keyRangeList, limit int64) (rowIter, error) {
	if len(keys) == 0 && len(keyRanges) == 0 {
		return nil, status.Error(codes.Unimplemented, "Cloud Spanner does not support reading no keys")
	}
	return d.readTable(tbl, cols, func(t *table, ri *rawIter, colIndexes []int) error {
		done := make(map[int]bool)
		for _, key := range keys {
			pk, err := t.primaryKey(key.Values)
			if err != nil {
				return err
			}
			rowNum, found := t.rowForPK(pk)
			if !found {
				continue
			}
			if done[rowNum] {
				continue
			}
			done[rowNum] = true
			ri.add(t.rows[rowNum], colIndexes)
			if limit > 0 && len(ri.rows) >= int(limit) {
				return nil
			}
		}
		for _, r := range keyRanges {
			var err error
			r.startKey, err = t.primaryKeyPrefix(r.start.Values)
			if err != nil {
				return err
			}
			r.endKey, err = t.primaryKeyPrefix(r.end.Values)
			if err != nil {
				return err
			}
			startRow, endRow := t.findRange(r)
			for rowNum := startRow; rowNum < endRow; rowNum++ {
				if done[rowNum] {
					continue
				}
				done[rowNum] = true
				ri.add(t.rows[rowNum], colIndexes)
				if limit > 0 && len(ri.rows) >= int(limit) {
					return nil
				}
			}
		}
		return nil
	})
}
func (d *database) gologoo__ReadAll_48b215e0cf779586a3b42f53798ab710(tbl spansql.ID, cols []spansql.ID, limit int64) (*rawIter, error) {
	return d.readTable(tbl, cols, func(t *table, ri *rawIter, colIndexes []int) error {
		for _, r := range t.rows {
			ri.add(r, colIndexes)
			if limit > 0 && len(ri.rows) >= int(limit) {
				break
			}
		}
		return nil
	})
}
func (t *table) gologoo__addColumn_48b215e0cf779586a3b42f53798ab710(cd spansql.ColumnDef, newTable bool) *status.Status {
	if !newTable && cd.NotNull {
		return status.Newf(codes.InvalidArgument, "new non-key columns cannot be NOT NULL")
	}
	if _, ok := t.colIndex[cd.Name]; ok {
		return status.Newf(codes.AlreadyExists, "column %s already exists", cd.Name)
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.rows) > 0 {
		if cd.NotNull {
			return status.Newf(codes.Unimplemented, "can't add NOT NULL columns to non-empty tables yet")
		}
		for i := range t.rows {
			if cd.Generated != nil {
				ec := evalContext{cols: t.cols, row: t.rows[i]}
				val, err := ec.evalExpr(cd.Generated)
				if err != nil {
					return status.Newf(codes.InvalidArgument, "could not backfill values for generated column: %v", err)
				}
				t.rows[i] = append(t.rows[i], val)
			} else {
				t.rows[i] = append(t.rows[i], nil)
			}
		}
	}
	t.cols = append(t.cols, colInfo{Name: cd.Name, Type: cd.Type, NotNull: cd.NotNull, Generated: cd.Generated})
	t.colIndex[cd.Name] = len(t.cols) - 1
	if !newTable {
		t.origIndex[cd.Name] = len(t.cols) - 1
	}
	return nil
}
func (t *table) gologoo__dropColumn_48b215e0cf779586a3b42f53798ab710(name spansql.ID) *status.Status {
	t.mu.Lock()
	defer t.mu.Unlock()
	ci, ok := t.colIndex[name]
	if !ok {
		return status.Newf(codes.InvalidArgument, "unknown column %q", name)
	}
	if ci < t.pkCols {
		return status.Newf(codes.InvalidArgument, "can't drop primary key column %q", name)
	}
	t.cols = append(t.cols[:ci], t.cols[ci+1:]...)
	delete(t.colIndex, name)
	for i, col := range t.cols {
		t.colIndex[col.Name] = i
	}
	pre := t.origIndex[name]
	delete(t.origIndex, name)
	for n, i := range t.origIndex {
		if i > pre {
			t.origIndex[n]--
		}
	}
	for i := range t.rows {
		t.rows[i] = append(t.rows[i][:ci], t.rows[i][ci+1:]...)
	}
	return nil
}
func (t *table) gologoo__alterColumn_48b215e0cf779586a3b42f53798ab710(alt spansql.AlterColumn) *status.Status {
	sct, ok := alt.Alteration.(spansql.SetColumnType)
	if !ok {
		return status.Newf(codes.InvalidArgument, "unsupported ALTER COLUMN %s", alt.SQL())
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	ci, ok := t.colIndex[alt.Name]
	if !ok {
		return status.Newf(codes.InvalidArgument, "unknown column %q", alt.Name)
	}
	oldT, newT := t.cols[ci].Type, sct.Type
	stringOrBytes := func(bt spansql.TypeBase) bool {
		return bt == spansql.String || bt == spansql.Bytes
	}
	if !t.cols[ci].NotNull && sct.NotNull {
		if ci < t.pkCols {
			return status.Newf(codes.InvalidArgument, "cannot set NOT NULL on primary key column %q", alt.Name)
		}
		if oldT.Array {
			return status.Newf(codes.InvalidArgument, "cannot set NOT NULL on array-typed column %q", alt.Name)
		}
		for _, row := range t.rows {
			if row[ci] == nil {
				return status.Newf(codes.InvalidArgument, "cannot set NOT NULL on column %q that contains NULL values", alt.Name)
			}
		}
	}
	var conv func(x interface {
	}) interface {
	}
	if stringOrBytes(oldT.Base) && stringOrBytes(newT.Base) && !oldT.Array && !newT.Array {
		if oldT.Base == spansql.Bytes && newT.Base == spansql.String {
			conv = func(x interface {
			}) interface {
			} {
				return string(x.([]byte))
			}
		} else if oldT.Base == spansql.String && newT.Base == spansql.Bytes {
			conv = func(x interface {
			}) interface {
			} {
				return []byte(x.(string))
			}
		}
	} else if oldT == newT {
	} else {
		return status.Newf(codes.InvalidArgument, "unsupported ALTER COLUMN %s", alt.SQL())
	}
	t.cols[ci].NotNull = sct.NotNull
	t.cols[ci].Type = newT
	if conv != nil {
		for _, row := range t.rows {
			if row[ci] != nil {
				row[ci] = conv(row[ci])
			}
		}
	}
	return nil
}
func (t *table) gologoo__addRowDeletionPolicy_48b215e0cf779586a3b42f53798ab710(ard spansql.AddRowDeletionPolicy) *status.Status {
	_, ok := t.colIndex[ard.RowDeletionPolicy.Column]
	if !ok {
		return status.Newf(codes.InvalidArgument, "unknown column %q", ard.RowDeletionPolicy.Column)
	}
	if t.rdw != nil {
		return status.New(codes.InvalidArgument, "table already has a row deletion policy")
	}
	t.rdw = &ard.RowDeletionPolicy
	return nil
}
func (t *table) gologoo__replaceRowDeletionPolicy_48b215e0cf779586a3b42f53798ab710(ard spansql.ReplaceRowDeletionPolicy) *status.Status {
	_, ok := t.colIndex[ard.RowDeletionPolicy.Column]
	if !ok {
		return status.Newf(codes.InvalidArgument, "unknown column %q", ard.RowDeletionPolicy.Column)
	}
	if t.rdw == nil {
		return status.New(codes.InvalidArgument, "table does not have a row deletion policy")
	}
	t.rdw = &ard.RowDeletionPolicy
	return nil
}
func (t *table) gologoo__dropRowDeletionPolicy_48b215e0cf779586a3b42f53798ab710(ard spansql.DropRowDeletionPolicy) *status.Status {
	if t.rdw == nil {
		return status.New(codes.InvalidArgument, "table does not have a row deletion policy")
	}
	t.rdw = nil
	return nil
}
func (t *table) gologoo__insertRow_48b215e0cf779586a3b42f53798ab710(rowNum int, r row) {
	t.rows = append(t.rows, nil)
	copy(t.rows[rowNum+1:], t.rows[rowNum:])
	t.rows[rowNum] = r
}
func (t *table) gologoo__findRange_48b215e0cf779586a3b42f53798ab710(r *keyRange) (int, int) {
	startRow := sort.Search(len(t.rows), func(i int) bool {
		return rowCmp(r.startKey, t.rows[i][:t.pkCols], t.pkDesc) <= 0
	})
	if startRow == len(t.rows) {
		return startRow, startRow
	}
	if !r.startClosed && rowCmp(r.startKey, t.rows[startRow][:t.pkCols], t.pkDesc) == 0 {
		startRow++
	}
	endRow := sort.Search(len(t.rows), func(i int) bool {
		return rowCmp(r.endKey, t.rows[i][:t.pkCols], t.pkDesc) < 0
	})
	if !r.endClosed && rowCmp(r.endKey, t.rows[endRow-1][:t.pkCols], t.pkDesc) == 0 {
		endRow--
	}
	return startRow, endRow
}
func (t *table) gologoo__colIndexes_48b215e0cf779586a3b42f53798ab710(cols []spansql.ID) ([]int, error) {
	var is []int
	for _, col := range cols {
		i, ok := t.colIndex[col]
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "column %s not in table", col)
		}
		is = append(is, i)
	}
	return is, nil
}
func (t *table) gologoo__primaryKey_48b215e0cf779586a3b42f53798ab710(values []*structpb.Value) ([]interface {
}, error) {
	if len(values) != t.pkCols {
		return nil, status.Errorf(codes.InvalidArgument, "primary key length mismatch: got %d values, table has %d", len(values), t.pkCols)
	}
	return t.primaryKeyPrefix(values)
}
func (t *table) gologoo__primaryKeyPrefix_48b215e0cf779586a3b42f53798ab710(values []*structpb.Value) ([]interface {
}, error) {
	if len(values) > t.pkCols {
		return nil, status.Errorf(codes.InvalidArgument, "primary key length too long: got %d values, table has %d", len(values), t.pkCols)
	}
	var pk []interface {
	}
	for i, value := range values {
		v, err := valForType(value, t.cols[i].Type)
		if err != nil {
			return nil, err
		}
		pk = append(pk, v)
	}
	return pk, nil
}
func (t *table) gologoo__rowForPK_48b215e0cf779586a3b42f53798ab710(pk []interface {
}) (row int, found bool) {
	if len(pk) != t.pkCols {
		panic(fmt.Sprintf("primary key length mismatch: got %d values, table has %d", len(pk), t.pkCols))
	}
	i := sort.Search(len(t.rows), func(i int) bool {
		return rowCmp(pk, t.rows[i][:t.pkCols], t.pkDesc) <= 0
	})
	if i == len(t.rows) {
		return i, false
	}
	return i, rowEqual(pk, t.rows[i][:t.pkCols])
}
func gologoo__rowCmp_48b215e0cf779586a3b42f53798ab710(a, b []interface {
}, desc []bool) int {
	for i := 0; i < len(a); i++ {
		if cmp := compareVals(a[i], b[i]); cmp != 0 {
			if desc[i] {
				cmp = -cmp
			}
			return cmp
		}
	}
	return 0
}
func gologoo__rowEqual_48b215e0cf779586a3b42f53798ab710(a, b []interface {
}) bool {
	for i := 0; i < len(a); i++ {
		if compareVals(a[i], b[i]) != 0 {
			return false
		}
	}
	return true
}
func gologoo__valForType_48b215e0cf779586a3b42f53798ab710(v *structpb.Value, t spansql.Type) (interface {
}, error) {
	if _, ok := v.Kind.(*structpb.Value_NullValue); ok {
		return nil, nil
	}
	if lv, ok := v.Kind.(*structpb.Value_ListValue); ok && t.Array {
		et := t
		et.Array = false
		arr := make([]interface {
		}, 0, len(lv.ListValue.Values))
		for _, v := range lv.ListValue.Values {
			x, err := valForType(v, et)
			if err != nil {
				return nil, err
			}
			arr = append(arr, x)
		}
		return arr, nil
	}
	switch t.Base {
	case spansql.Bool:
		bv, ok := v.Kind.(*structpb.Value_BoolValue)
		if ok {
			return bv.BoolValue, nil
		}
	case spansql.Int64:
		sv, ok := v.Kind.(*structpb.Value_StringValue)
		if ok {
			x, err := strconv.ParseInt(sv.StringValue, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("bad int64 string %q: %v", sv.StringValue, err)
			}
			return x, nil
		}
	case spansql.Float64:
		nv, ok := v.Kind.(*structpb.Value_NumberValue)
		if ok {
			return nv.NumberValue, nil
		}
	case spansql.String:
		sv, ok := v.Kind.(*structpb.Value_StringValue)
		if ok {
			return sv.StringValue, nil
		}
	case spansql.Bytes:
		sv, ok := v.Kind.(*structpb.Value_StringValue)
		if ok {
			return base64.StdEncoding.DecodeString(sv.StringValue)
		}
	case spansql.Date:
		sv, ok := v.Kind.(*structpb.Value_StringValue)
		if ok {
			s := sv.StringValue
			d, err := parseAsDate(s)
			if err != nil {
				return nil, fmt.Errorf("bad DATE string %q: %v", s, err)
			}
			return d, nil
		}
	case spansql.Timestamp:
		sv, ok := v.Kind.(*structpb.Value_StringValue)
		if ok {
			s := sv.StringValue
			if strings.ToLower(s) == "spanner.commit_timestamp()" {
				return commitTimestampSentinel, nil
			}
			t, err := parseAsTimestamp(s)
			if err != nil {
				return nil, fmt.Errorf("bad TIMESTAMP string %q: %v", s, err)
			}
			return t, nil
		}
	}
	return nil, fmt.Errorf("unsupported inserting value kind %T into column of type %s", v.Kind, t.SQL())
}

type keyRange struct {
	start, end             *structpb.ListValue
	startClosed, endClosed bool
	startKey, endKey       []interface {
	}
}

func (r *keyRange) gologoo__String_48b215e0cf779586a3b42f53798ab710() string {
	var sb bytes.Buffer
	if r.startClosed {
		sb.WriteString("[")
	} else {
		sb.WriteString("(")
	}
	fmt.Fprintf(&sb, "%v,%v", r.startKey, r.endKey)
	if r.endClosed {
		sb.WriteString("]")
	} else {
		sb.WriteString(")")
	}
	return sb.String()
}

type keyRangeList []*keyRange

func (d *database) gologoo__Execute_48b215e0cf779586a3b42f53798ab710(stmt spansql.DMLStmt, params queryParams) (int, error) {
	switch stmt := stmt.(type) {
	default:
		return 0, status.Errorf(codes.Unimplemented, "unhandled DML statement type %T", stmt)
	case *spansql.Delete:
		t, err := d.table(stmt.Table)
		if err != nil {
			return 0, err
		}
		t.mu.Lock()
		defer t.mu.Unlock()
		n := 0
		for i := 0; i < len(t.rows); {
			ec := evalContext{cols: t.cols, row: t.rows[i], params: params}
			b, err := ec.evalBoolExpr(stmt.Where)
			if err != nil {
				return 0, err
			}
			if b != nil && *b {
				copy(t.rows[i:], t.rows[i+1:])
				t.rows = t.rows[:len(t.rows)-1]
				n++
				continue
			}
			i++
		}
		return n, nil
	case *spansql.Update:
		t, err := d.table(stmt.Table)
		if err != nil {
			return 0, err
		}
		t.mu.Lock()
		defer t.mu.Unlock()
		ec := evalContext{cols: t.cols, params: params}
		var dstIndex []int
		var expr []spansql.Expr
		for _, ui := range stmt.Items {
			i, err := ec.resolveColumnIndex(ui.Column)
			if err != nil {
				return 0, err
			}
			if i < t.pkCols {
				return 0, status.Errorf(codes.InvalidArgument, "cannot update primary key %s", ui.Column)
			}
			dstIndex = append(dstIndex, i)
			expr = append(expr, ui.Value)
		}
		n := 0
		values := make(row, len(stmt.Items))
		for i := 0; i < len(t.rows); i++ {
			ec.row = t.rows[i]
			b, err := ec.evalBoolExpr(stmt.Where)
			if err != nil {
				return 0, err
			}
			if b != nil && *b {
				for j := range dstIndex {
					if expr[j] == nil {
						values[j] = nil
						continue
					}
					v, err := ec.evalExpr(expr[j])
					if err != nil {
						return 0, err
					}
					values[j] = v
				}
				for j, v := range values {
					t.rows[i][dstIndex[j]] = v
				}
				n++
			}
		}
		return n, nil
	}
}
func gologoo__parseAsDate_48b215e0cf779586a3b42f53798ab710(s string) (civil.Date, error) {
	return civil.ParseDate(s)
}
func gologoo__parseAsTimestamp_48b215e0cf779586a3b42f53798ab710(s string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05.999999999Z", s)
}
func (d *database) NewReadOnlyTransaction() *transaction {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NewReadOnlyTransaction_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : (none)\n")
	r0 := d.gologoo__NewReadOnlyTransaction_48b215e0cf779586a3b42f53798ab710()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (d *database) NewTransaction() *transaction {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NewTransaction_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : (none)\n")
	r0 := d.gologoo__NewTransaction_48b215e0cf779586a3b42f53798ab710()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (tx *transaction) Start() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Start_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : (none)\n")
	tx.gologoo__Start_48b215e0cf779586a3b42f53798ab710()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (tx *transaction) checkMutable() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__checkMutable_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : (none)\n")
	r0 := tx.gologoo__checkMutable_48b215e0cf779586a3b42f53798ab710()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (tx *transaction) Commit() (time.Time, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Commit_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : (none)\n")
	r0, r1 := tx.gologoo__Commit_48b215e0cf779586a3b42f53798ab710()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (tx *transaction) Rollback() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Rollback_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : (none)\n")
	tx.gologoo__Rollback_48b215e0cf779586a3b42f53798ab710()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (r row) copyDataElem(index int) interface {
} {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__copyDataElem_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v\n", index)
	r0 := r.gologoo__copyDataElem_48b215e0cf779586a3b42f53798ab710(index)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (r row) copyAllData() row {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__copyAllData_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : (none)\n")
	r0 := r.gologoo__copyAllData_48b215e0cf779586a3b42f53798ab710()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (r row) copyData(indexes []int) row {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__copyData_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v\n", indexes)
	r0 := r.gologoo__copyData_48b215e0cf779586a3b42f53798ab710(indexes)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (d *database) LastCommitTimestamp() time.Time {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__LastCommitTimestamp_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : (none)\n")
	r0 := d.gologoo__LastCommitTimestamp_48b215e0cf779586a3b42f53798ab710()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (d *database) GetDDL() []spansql.DDLStmt {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetDDL_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : (none)\n")
	r0 := d.gologoo__GetDDL_48b215e0cf779586a3b42f53798ab710()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (d *database) ApplyDDL(stmt spansql.DDLStmt) *status.Status {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ApplyDDL_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v\n", stmt)
	r0 := d.gologoo__ApplyDDL_48b215e0cf779586a3b42f53798ab710(stmt)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (d *database) table(tbl spansql.ID) (*table, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__table_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v\n", tbl)
	r0, r1 := d.gologoo__table_48b215e0cf779586a3b42f53798ab710(tbl)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (d *database) writeValues(tx *transaction, tbl spansql.ID, cols []spansql.ID, values []*structpb.ListValue, f func(t *table, colIndexes []int, r row) error) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__writeValues_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v %v %v %v %v\n", tx, tbl, cols, values, f)
	r0 := d.gologoo__writeValues_48b215e0cf779586a3b42f53798ab710(tx, tbl, cols, values, f)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (d *database) Insert(tx *transaction, tbl spansql.ID, cols []spansql.ID, values []*structpb.ListValue) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Insert_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v %v %v %v\n", tx, tbl, cols, values)
	r0 := d.gologoo__Insert_48b215e0cf779586a3b42f53798ab710(tx, tbl, cols, values)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (d *database) Update(tx *transaction, tbl spansql.ID, cols []spansql.ID, values []*structpb.ListValue) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Update_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v %v %v %v\n", tx, tbl, cols, values)
	r0 := d.gologoo__Update_48b215e0cf779586a3b42f53798ab710(tx, tbl, cols, values)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (d *database) InsertOrUpdate(tx *transaction, tbl spansql.ID, cols []spansql.ID, values []*structpb.ListValue) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__InsertOrUpdate_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v %v %v %v\n", tx, tbl, cols, values)
	r0 := d.gologoo__InsertOrUpdate_48b215e0cf779586a3b42f53798ab710(tx, tbl, cols, values)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (d *database) Delete(tx *transaction, table spansql.ID, keys []*structpb.ListValue, keyRanges keyRangeList, all bool) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Delete_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v %v %v %v %v\n", tx, table, keys, keyRanges, all)
	r0 := d.gologoo__Delete_48b215e0cf779586a3b42f53798ab710(tx, table, keys, keyRanges, all)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (d *database) readTable(table spansql.ID, cols []spansql.ID, f func(*table, *rawIter, []int) error) (*rawIter, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__readTable_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v %v %v\n", table, cols, f)
	r0, r1 := d.gologoo__readTable_48b215e0cf779586a3b42f53798ab710(table, cols, f)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (d *database) Read(tbl spansql.ID, cols []spansql.ID, keys []*structpb.ListValue, keyRanges keyRangeList, limit int64) (rowIter, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Read_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v %v %v %v %v\n", tbl, cols, keys, keyRanges, limit)
	r0, r1 := d.gologoo__Read_48b215e0cf779586a3b42f53798ab710(tbl, cols, keys, keyRanges, limit)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (d *database) ReadAll(tbl spansql.ID, cols []spansql.ID, limit int64) (*rawIter, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ReadAll_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v %v %v\n", tbl, cols, limit)
	r0, r1 := d.gologoo__ReadAll_48b215e0cf779586a3b42f53798ab710(tbl, cols, limit)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *table) addColumn(cd spansql.ColumnDef, newTable bool) *status.Status {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addColumn_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v %v\n", cd, newTable)
	r0 := t.gologoo__addColumn_48b215e0cf779586a3b42f53798ab710(cd, newTable)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *table) dropColumn(name spansql.ID) *status.Status {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__dropColumn_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v\n", name)
	r0 := t.gologoo__dropColumn_48b215e0cf779586a3b42f53798ab710(name)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *table) alterColumn(alt spansql.AlterColumn) *status.Status {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__alterColumn_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v\n", alt)
	r0 := t.gologoo__alterColumn_48b215e0cf779586a3b42f53798ab710(alt)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *table) addRowDeletionPolicy(ard spansql.AddRowDeletionPolicy) *status.Status {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addRowDeletionPolicy_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v\n", ard)
	r0 := t.gologoo__addRowDeletionPolicy_48b215e0cf779586a3b42f53798ab710(ard)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *table) replaceRowDeletionPolicy(ard spansql.ReplaceRowDeletionPolicy) *status.Status {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__replaceRowDeletionPolicy_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v\n", ard)
	r0 := t.gologoo__replaceRowDeletionPolicy_48b215e0cf779586a3b42f53798ab710(ard)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *table) dropRowDeletionPolicy(ard spansql.DropRowDeletionPolicy) *status.Status {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__dropRowDeletionPolicy_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v\n", ard)
	r0 := t.gologoo__dropRowDeletionPolicy_48b215e0cf779586a3b42f53798ab710(ard)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *table) insertRow(rowNum int, r row) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__insertRow_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v %v\n", rowNum, r)
	t.gologoo__insertRow_48b215e0cf779586a3b42f53798ab710(rowNum, r)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (t *table) findRange(r *keyRange) (int, int) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__findRange_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v\n", r)
	r0, r1 := t.gologoo__findRange_48b215e0cf779586a3b42f53798ab710(r)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *table) colIndexes(cols []spansql.ID) ([]int, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__colIndexes_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v\n", cols)
	r0, r1 := t.gologoo__colIndexes_48b215e0cf779586a3b42f53798ab710(cols)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *table) primaryKey(values []*structpb.Value) ([]interface {
}, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__primaryKey_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v\n", values)
	r0, r1 := t.gologoo__primaryKey_48b215e0cf779586a3b42f53798ab710(values)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *table) primaryKeyPrefix(values []*structpb.Value) ([]interface {
}, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__primaryKeyPrefix_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v\n", values)
	r0, r1 := t.gologoo__primaryKeyPrefix_48b215e0cf779586a3b42f53798ab710(values)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *table) rowForPK(pk []interface {
}) (row int, found bool) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__rowForPK_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v\n", pk)
	row, found = t.gologoo__rowForPK_48b215e0cf779586a3b42f53798ab710(pk)
	log.Printf("Output: %v %v\n", row, found)
	return
}
func rowCmp(a, b []interface {
}, desc []bool) int {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__rowCmp_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v %v %v\n", a, b, desc)
	r0 := gologoo__rowCmp_48b215e0cf779586a3b42f53798ab710(a, b, desc)
	log.Printf("Output: %v\n", r0)
	return r0
}
func rowEqual(a, b []interface {
}) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__rowEqual_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v %v\n", a, b)
	r0 := gologoo__rowEqual_48b215e0cf779586a3b42f53798ab710(a, b)
	log.Printf("Output: %v\n", r0)
	return r0
}
func valForType(v *structpb.Value, t spansql.Type) (interface {
}, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__valForType_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v %v\n", v, t)
	r0, r1 := gologoo__valForType_48b215e0cf779586a3b42f53798ab710(v, t)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (r *keyRange) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : (none)\n")
	r0 := r.gologoo__String_48b215e0cf779586a3b42f53798ab710()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (d *database) Execute(stmt spansql.DMLStmt, params queryParams) (int, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Execute_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v %v\n", stmt, params)
	r0, r1 := d.gologoo__Execute_48b215e0cf779586a3b42f53798ab710(stmt, params)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func parseAsDate(s string) (civil.Date, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseAsDate_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v\n", s)
	r0, r1 := gologoo__parseAsDate_48b215e0cf779586a3b42f53798ab710(s)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func parseAsTimestamp(s string) (time.Time, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseAsTimestamp_48b215e0cf779586a3b42f53798ab710")
	log.Printf("Input : %v\n", s)
	r0, r1 := gologoo__parseAsTimestamp_48b215e0cf779586a3b42f53798ab710(s)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
