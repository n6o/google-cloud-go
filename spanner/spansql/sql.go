package spansql

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"log"
)

func gologoo__buildSQL_6948c200d52dc131145058c7d4919511(x interface {
	addSQL(*strings.Builder)
}) string {
	var sb strings.Builder
	x.addSQL(&sb)
	return sb.String()
}
func (ct CreateTable) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	str := "CREATE TABLE " + ct.Name.SQL() + " (\n"
	for _, c := range ct.Columns {
		str += "  " + c.SQL() + ",\n"
	}
	for _, tc := range ct.Constraints {
		str += "  " + tc.SQL() + ",\n"
	}
	str += ") PRIMARY KEY("
	for i, c := range ct.PrimaryKey {
		if i > 0 {
			str += ", "
		}
		str += c.SQL()
	}
	str += ")"
	if il := ct.Interleave; il != nil {
		str += ",\n  INTERLEAVE IN PARENT " + il.Parent.SQL() + " ON DELETE " + il.OnDelete.SQL()
	}
	if rdp := ct.RowDeletionPolicy; rdp != nil {
		str += ",\n  " + rdp.SQL()
	}
	return str
}
func (ci CreateIndex) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	str := "CREATE"
	if ci.Unique {
		str += " UNIQUE"
	}
	if ci.NullFiltered {
		str += " NULL_FILTERED"
	}
	str += " INDEX " + ci.Name.SQL() + " ON " + ci.Table.SQL() + "("
	for i, c := range ci.Columns {
		if i > 0 {
			str += ", "
		}
		str += c.SQL()
	}
	str += ")"
	if len(ci.Storing) > 0 {
		str += " STORING (" + idList(ci.Storing, ", ") + ")"
	}
	if ci.Interleave != "" {
		str += ", INTERLEAVE IN " + ci.Interleave.SQL()
	}
	return str
}
func (cv CreateView) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	str := "CREATE"
	if cv.OrReplace {
		str += " OR REPLACE"
	}
	str += " VIEW " + cv.Name.SQL() + " SQL SECURITY INVOKER AS " + cv.Query.SQL()
	return str
}
func (dt DropTable) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "DROP TABLE " + dt.Name.SQL()
}
func (di DropIndex) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "DROP INDEX " + di.Name.SQL()
}
func (dv DropView) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "DROP VIEW " + dv.Name.SQL()
}
func (at AlterTable) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "ALTER TABLE " + at.Name.SQL() + " " + at.Alteration.SQL()
}
func (ac AddColumn) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "ADD COLUMN " + ac.Def.SQL()
}
func (dc DropColumn) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "DROP COLUMN " + dc.Name.SQL()
}
func (ac AddConstraint) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "ADD " + ac.Constraint.SQL()
}
func (dc DropConstraint) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "DROP CONSTRAINT " + dc.Name.SQL()
}
func (sod SetOnDelete) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "SET ON DELETE " + sod.Action.SQL()
}
func (od OnDelete) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	switch od {
	case NoActionOnDelete:
		return "NO ACTION"
	case CascadeOnDelete:
		return "CASCADE"
	}
	panic("unknown OnDelete")
}
func (ac AlterColumn) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "ALTER COLUMN " + ac.Name.SQL() + " " + ac.Alteration.SQL()
}
func (ardp AddRowDeletionPolicy) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "ADD " + ardp.RowDeletionPolicy.SQL()
}
func (rrdp ReplaceRowDeletionPolicy) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "REPLACE " + rrdp.RowDeletionPolicy.SQL()
}
func (drdp DropRowDeletionPolicy) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "DROP ROW DELETION POLICY"
}
func (sct SetColumnType) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	str := sct.Type.SQL()
	if sct.NotNull {
		str += " NOT NULL"
	}
	return str
}
func (sco SetColumnOptions) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "SET " + sco.Options.SQL()
}
func (co ColumnOptions) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	str := "OPTIONS ("
	if co.AllowCommitTimestamp != nil {
		if *co.AllowCommitTimestamp {
			str += "allow_commit_timestamp = true"
		} else {
			str += "allow_commit_timestamp = null"
		}
	}
	str += ")"
	return str
}
func (ad AlterDatabase) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "ALTER DATABASE " + ad.Name.SQL() + " " + ad.Alteration.SQL()
}
func (sdo SetDatabaseOptions) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "SET " + sdo.Options.SQL()
}
func (do DatabaseOptions) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	str := "OPTIONS ("
	hasOpt := false
	if do.OptimizerVersion != nil {
		hasOpt = true
		if *do.OptimizerVersion == 0 {
			str += "optimizer_version=null"
		} else {
			str += fmt.Sprintf("optimizer_version=%v", *do.OptimizerVersion)
		}
	}
	if do.VersionRetentionPeriod != nil {
		if hasOpt {
			str += ", "
		}
		hasOpt = true
		if *do.VersionRetentionPeriod == "" {
			str += "version_retention_period=null"
		} else {
			str += fmt.Sprintf("version_retention_period='%s'", *do.VersionRetentionPeriod)
		}
	}
	if do.EnableKeyVisualizer != nil {
		if hasOpt {
			str += ", "
		}
		hasOpt = true
		if *do.EnableKeyVisualizer {
			str += "enable_key_visualizer=true"
		} else {
			str += "enable_key_visualizer=null"
		}
	}
	str += ")"
	return str
}
func (d *Delete) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "DELETE FROM " + d.Table.SQL() + " WHERE " + d.Where.SQL()
}
func (u *Update) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	str := "UPDATE " + u.Table.SQL() + " SET "
	for i, item := range u.Items {
		if i > 0 {
			str += ", "
		}
		str += item.Column.SQL() + " = "
		if item.Value != nil {
			str += item.Value.SQL()
		} else {
			str += "DEFAULT"
		}
	}
	str += " WHERE " + u.Where.SQL()
	return str
}
func (cd ColumnDef) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	str := cd.Name.SQL() + " " + cd.Type.SQL()
	if cd.NotNull {
		str += " NOT NULL"
	}
	if cd.Generated != nil {
		str += " AS (" + cd.Generated.SQL() + ") STORED"
	}
	if cd.Options != (ColumnOptions{}) {
		str += " " + cd.Options.SQL()
	}
	return str
}
func (tc TableConstraint) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	var str string
	if tc.Name != "" {
		str += "CONSTRAINT " + tc.Name.SQL() + " "
	}
	str += tc.Constraint.SQL()
	return str
}
func (rdp RowDeletionPolicy) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "ROW DELETION POLICY ( OLDER_THAN ( " + rdp.Column.SQL() + ", INTERVAL " + strconv.FormatInt(rdp.NumDays, 10) + " DAY ))"
}
func (fk ForeignKey) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	str := "FOREIGN KEY (" + idList(fk.Columns, ", ")
	str += ") REFERENCES " + fk.RefTable.SQL() + " ("
	str += idList(fk.RefColumns, ", ") + ")"
	return str
}
func (c Check) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return "CHECK (" + c.Expr.SQL() + ")"
}
func (t Type) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	str := t.Base.SQL()
	if t.Len > 0 && (t.Base == String || t.Base == Bytes) {
		str += "("
		if t.Len == MaxLen {
			str += "MAX"
		} else {
			str += strconv.FormatInt(t.Len, 10)
		}
		str += ")"
	}
	if t.Array {
		str = "ARRAY<" + str + ">"
	}
	return str
}
func (tb TypeBase) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	switch tb {
	case Bool:
		return "BOOL"
	case Int64:
		return "INT64"
	case Float64:
		return "FLOAT64"
	case Numeric:
		return "NUMERIC"
	case String:
		return "STRING"
	case Bytes:
		return "BYTES"
	case Date:
		return "DATE"
	case Timestamp:
		return "TIMESTAMP"
	case JSON:
		return "JSON"
	}
	panic("unknown TypeBase")
}
func (kp KeyPart) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	str := kp.Column.SQL()
	if kp.Desc {
		str += " DESC"
	}
	return str
}
func (q Query) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(q)
}
func (q Query) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	q.Select.addSQL(sb)
	if len(q.Order) > 0 {
		sb.WriteString(" ORDER BY ")
		for i, o := range q.Order {
			if i > 0 {
				sb.WriteString(", ")
			}
			o.addSQL(sb)
		}
	}
	if q.Limit != nil {
		sb.WriteString(" LIMIT ")
		sb.WriteString(q.Limit.SQL())
		if q.Offset != nil {
			sb.WriteString(" OFFSET ")
			sb.WriteString(q.Offset.SQL())
		}
	}
}
func (sel Select) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(sel)
}
func (sel Select) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	sb.WriteString("SELECT ")
	if sel.Distinct {
		sb.WriteString("DISTINCT ")
	}
	for i, e := range sel.List {
		if i > 0 {
			sb.WriteString(", ")
		}
		e.addSQL(sb)
		if len(sel.ListAliases) > 0 {
			alias := sel.ListAliases[i]
			if alias != "" {
				sb.WriteString(" AS ")
				sb.WriteString(alias.SQL())
			}
		}
	}
	if len(sel.From) > 0 {
		sb.WriteString(" FROM ")
		for i, f := range sel.From {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(f.SQL())
		}
	}
	if sel.Where != nil {
		sb.WriteString(" WHERE ")
		sel.Where.addSQL(sb)
	}
	if len(sel.GroupBy) > 0 {
		sb.WriteString(" GROUP BY ")
		addExprList(sb, sel.GroupBy, ", ")
	}
}
func (sft SelectFromTable) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	str := sft.Table.SQL()
	if len(sft.Hints) > 0 {
		str += "@{"
		kvs := make([]string, len(sft.Hints))
		i := 0
		for k, v := range sft.Hints {
			kvs[i] = fmt.Sprintf("%s=%s", k, v)
			i++
		}
		sort.Strings(kvs)
		str += strings.Join(kvs, ",")
		str += "}"
	}
	if sft.Alias != "" {
		str += " AS " + sft.Alias.SQL()
	}
	return str
}
func (sfj SelectFromJoin) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	str := sfj.LHS.SQL() + " " + joinTypes[sfj.Type] + " JOIN "
	str += sfj.RHS.SQL()
	if sfj.On != nil {
		str += " ON " + sfj.On.SQL()
	} else if len(sfj.Using) > 0 {
		str += " USING (" + idList(sfj.Using, ", ") + ")"
	}
	return str
}

var joinTypes = map[JoinType]string{InnerJoin: "INNER", CrossJoin: "CROSS", FullJoin: "FULL", LeftJoin: "LEFT", RightJoin: "RIGHT"}

func (sfu SelectFromUnnest) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	str := "UNNEST(" + sfu.Expr.SQL() + ")"
	if sfu.Alias != "" {
		str += " AS " + sfu.Alias.SQL()
	}
	return str
}
func (o Order) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(o)
}
func (o Order) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	o.Expr.addSQL(sb)
	if o.Desc {
		sb.WriteString(" DESC")
	}
}

var arithOps = map[ArithOperator]string{Mul: "*", Div: "/", Concat: "||", Add: "+", Sub: "-", BitShl: "<<", BitShr: ">>", BitAnd: "&", BitXor: "^", BitOr: "|"}

func (ao ArithOp) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(ao)
}
func (ao ArithOp) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	switch ao.Op {
	case Neg:
		sb.WriteString("-(")
		ao.RHS.addSQL(sb)
		sb.WriteString(")")
		return
	case Plus:
		sb.WriteString("+(")
		ao.RHS.addSQL(sb)
		sb.WriteString(")")
		return
	case BitNot:
		sb.WriteString("~(")
		ao.RHS.addSQL(sb)
		sb.WriteString(")")
		return
	}
	op, ok := arithOps[ao.Op]
	if !ok {
		panic("unknown ArithOp")
	}
	sb.WriteString("(")
	ao.LHS.addSQL(sb)
	sb.WriteString(")")
	sb.WriteString(op)
	sb.WriteString("(")
	ao.RHS.addSQL(sb)
	sb.WriteString(")")
}
func (lo LogicalOp) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(lo)
}
func (lo LogicalOp) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	switch lo.Op {
	default:
		panic("unknown LogicalOp")
	case And:
		lo.LHS.addSQL(sb)
		sb.WriteString(" AND ")
	case Or:
		lo.LHS.addSQL(sb)
		sb.WriteString(" OR ")
	case Not:
		sb.WriteString("NOT ")
	}
	lo.RHS.addSQL(sb)
}

var compOps = map[ComparisonOperator]string{Lt: "<", Le: "<=", Gt: ">", Ge: ">=", Eq: "=", Ne: "!=", Like: "LIKE", NotLike: "NOT LIKE", Between: "BETWEEN", NotBetween: "NOT BETWEEN"}

func (co ComparisonOp) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(co)
}
func (co ComparisonOp) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	op, ok := compOps[co.Op]
	if !ok {
		panic("unknown ComparisonOp")
	}
	co.LHS.addSQL(sb)
	sb.WriteString(" ")
	sb.WriteString(op)
	sb.WriteString(" ")
	co.RHS.addSQL(sb)
	if co.Op == Between || co.Op == NotBetween {
		sb.WriteString(" AND ")
		co.RHS2.addSQL(sb)
	}
}
func (io InOp) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(io)
}
func (io InOp) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	io.LHS.addSQL(sb)
	if io.Neg {
		sb.WriteString(" NOT")
	}
	sb.WriteString(" IN ")
	if io.Unnest {
		sb.WriteString("UNNEST")
	}
	sb.WriteString("(")
	addExprList(sb, io.RHS, ", ")
	sb.WriteString(")")
}
func (io IsOp) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(io)
}
func (io IsOp) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	io.LHS.addSQL(sb)
	sb.WriteString(" IS ")
	if io.Neg {
		sb.WriteString("NOT ")
	}
	io.RHS.addSQL(sb)
}
func (f Func) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(f)
}
func (f Func) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	sb.WriteString(f.Name)
	sb.WriteString("(")
	addExprList(sb, f.Args, ", ")
	sb.WriteString(")")
}
func (te TypedExpr) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(te)
}
func (te TypedExpr) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	te.Expr.addSQL(sb)
	sb.WriteString(" AS ")
	sb.WriteString(te.Type.SQL())
}
func (ee ExtractExpr) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(ee)
}
func (ee ExtractExpr) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	sb.WriteString(ee.Part)
	sb.WriteString(" FROM ")
	ee.Expr.addSQL(sb)
}
func (aze AtTimeZoneExpr) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(aze)
}
func (aze AtTimeZoneExpr) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	aze.Expr.addSQL(sb)
	sb.WriteString(" AT TIME ZONE ")
	sb.WriteString(aze.Zone)
}
func gologoo__idList_6948c200d52dc131145058c7d4919511(l []ID, join string) string {
	var ss []string
	for _, s := range l {
		ss = append(ss, s.SQL())
	}
	return strings.Join(ss, join)
}
func gologoo__addExprList_6948c200d52dc131145058c7d4919511(sb *strings.Builder, l []Expr, join string) {
	for i, s := range l {
		if i > 0 {
			sb.WriteString(join)
		}
		s.addSQL(sb)
	}
}
func gologoo__addIDList_6948c200d52dc131145058c7d4919511(sb *strings.Builder, l []ID, join string) {
	for i, s := range l {
		if i > 0 {
			sb.WriteString(join)
		}
		s.addSQL(sb)
	}
}
func (pe PathExp) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(pe)
}
func (pe PathExp) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	addIDList(sb, []ID(pe), ".")
}
func (p Paren) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(p)
}
func (p Paren) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	sb.WriteString("(")
	p.Expr.addSQL(sb)
	sb.WriteString(")")
}
func (a Array) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(a)
}
func (a Array) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	sb.WriteString("[")
	addExprList(sb, []Expr(a), ", ")
	sb.WriteString("]")
}
func (id ID) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(id)
}
func (id ID) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	if IsKeyword(string(id)) {
		sb.WriteString("`")
		sb.WriteString(string(id))
		sb.WriteString("`")
		return
	}
	sb.WriteString(string(id))
}
func (p Param) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(p)
}
func (p Param) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	sb.WriteString("@")
	sb.WriteString(string(p))
}
func (b BoolLiteral) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(b)
}
func (b BoolLiteral) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	if b {
		sb.WriteString("TRUE")
	} else {
		sb.WriteString("FALSE")
	}
}
func (NullLiteral) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(NullLiteral(0))
}
func (NullLiteral) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	sb.WriteString("NULL")
}
func (StarExpr) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(StarExpr(0))
}
func (StarExpr) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	sb.WriteString("*")
}
func (il IntegerLiteral) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(il)
}
func (il IntegerLiteral) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	fmt.Fprintf(sb, "%d", il)
}
func (fl FloatLiteral) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(fl)
}
func (fl FloatLiteral) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	fmt.Fprintf(sb, "%g", fl)
}
func (sl StringLiteral) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(sl)
}
func (sl StringLiteral) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	fmt.Fprintf(sb, "%q", sl)
}
func (bl BytesLiteral) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(bl)
}
func (bl BytesLiteral) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	fmt.Fprintf(sb, "B%q", bl)
}
func (dl DateLiteral) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(dl)
}
func (dl DateLiteral) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	fmt.Fprintf(sb, "DATE '%04d-%02d-%02d'", dl.Year, dl.Month, dl.Day)
}
func (tl TimestampLiteral) gologoo__SQL_6948c200d52dc131145058c7d4919511() string {
	return buildSQL(tl)
}
func (tl TimestampLiteral) gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb *strings.Builder) {
	fmt.Fprintf(sb, "TIMESTAMP '%s'", time.Time(tl).Format("2006-01-02 15:04:05.000000 -07:00"))
}
func buildSQL(x interface {
	addSQL(*strings.Builder)
}) string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__buildSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", x)
	r0 := gologoo__buildSQL_6948c200d52dc131145058c7d4919511(x)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ct CreateTable) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := ct.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ci CreateIndex) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := ci.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (cv CreateView) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := cv.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (dt DropTable) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := dt.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (di DropIndex) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := di.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (dv DropView) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := dv.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (at AlterTable) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := at.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ac AddColumn) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := ac.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (dc DropColumn) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := dc.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ac AddConstraint) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := ac.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (dc DropConstraint) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := dc.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (sod SetOnDelete) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := sod.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (od OnDelete) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := od.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ac AlterColumn) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := ac.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ardp AddRowDeletionPolicy) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := ardp.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (rrdp ReplaceRowDeletionPolicy) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := rrdp.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (drdp DropRowDeletionPolicy) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := drdp.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (sct SetColumnType) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := sct.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (sco SetColumnOptions) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := sco.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (co ColumnOptions) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := co.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ad AlterDatabase) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := ad.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (sdo SetDatabaseOptions) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := sdo.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (do DatabaseOptions) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := do.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (d *Delete) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := d.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (u *Update) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := u.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (cd ColumnDef) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := cd.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (tc TableConstraint) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := tc.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (rdp RowDeletionPolicy) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := rdp.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (fk ForeignKey) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := fk.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c Check) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t Type) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := t.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (tb TypeBase) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := tb.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (kp KeyPart) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := kp.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (q Query) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := q.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (q Query) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	q.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (sel Select) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := sel.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (sel Select) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	sel.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (sft SelectFromTable) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := sft.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (sfj SelectFromJoin) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := sfj.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (sfu SelectFromUnnest) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := sfu.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (o Order) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := o.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (o Order) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	o.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (ao ArithOp) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := ao.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ao ArithOp) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	ao.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (lo LogicalOp) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := lo.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (lo LogicalOp) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	lo.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (co ComparisonOp) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := co.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (co ComparisonOp) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	co.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (io InOp) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := io.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (io InOp) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	io.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (io IsOp) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := io.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (io IsOp) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	io.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (f Func) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := f.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (f Func) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	f.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (te TypedExpr) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := te.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (te TypedExpr) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	te.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (ee ExtractExpr) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := ee.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ee ExtractExpr) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	ee.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (aze AtTimeZoneExpr) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := aze.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (aze AtTimeZoneExpr) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	aze.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func idList(l []ID, join string) string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__idList_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v %v\n", l, join)
	r0 := gologoo__idList_6948c200d52dc131145058c7d4919511(l, join)
	log.Printf("Output: %v\n", r0)
	return r0
}
func addExprList(sb *strings.Builder, l []Expr, join string) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addExprList_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v %v %v\n", sb, l, join)
	gologoo__addExprList_6948c200d52dc131145058c7d4919511(sb, l, join)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func addIDList(sb *strings.Builder, l []ID, join string) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addIDList_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v %v %v\n", sb, l, join)
	gologoo__addIDList_6948c200d52dc131145058c7d4919511(sb, l, join)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (pe PathExp) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := pe.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (pe PathExp) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	pe.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p Paren) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := p.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p Paren) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	p.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (a Array) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := a.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (a Array) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	a.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (id ID) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := id.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (id ID) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	id.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p Param) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := p.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p Param) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	p.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (b BoolLiteral) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := b.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (b BoolLiteral) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	b.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv NullLiteral) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := recv.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv NullLiteral) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	recv.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv StarExpr) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := recv.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv StarExpr) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	recv.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (il IntegerLiteral) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := il.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (il IntegerLiteral) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	il.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (fl FloatLiteral) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := fl.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (fl FloatLiteral) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	fl.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (sl StringLiteral) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := sl.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (sl StringLiteral) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	sl.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (bl BytesLiteral) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := bl.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (bl BytesLiteral) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	bl.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (dl DateLiteral) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := dl.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (dl DateLiteral) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	dl.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (tl TimestampLiteral) SQL() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : (none)\n")
	r0 := tl.gologoo__SQL_6948c200d52dc131145058c7d4919511()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (tl TimestampLiteral) addSQL(sb *strings.Builder) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__addSQL_6948c200d52dc131145058c7d4919511")
	log.Printf("Input : %v\n", sb)
	tl.gologoo__addSQL_6948c200d52dc131145058c7d4919511(sb)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
