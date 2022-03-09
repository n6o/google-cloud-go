package spansql

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
	"cloud.google.com/go/civil"
	"log"
)

type CreateTable struct {
	Name              ID
	Columns           []ColumnDef
	Constraints       []TableConstraint
	PrimaryKey        []KeyPart
	Interleave        *Interleave
	RowDeletionPolicy *RowDeletionPolicy
	Position          Position
}

func (ct *CreateTable) gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32() string {
	return fmt.Sprintf("%#v", ct)
}
func (*CreateTable) gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (ct *CreateTable) gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32() Position {
	return ct.Position
}
func (ct *CreateTable) gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32() {
	for i := range ct.Columns {
		ct.Columns[i].clearOffset()
	}
	for i := range ct.Constraints {
		ct.Constraints[i].clearOffset()
	}
	ct.Position.Offset = 0
}

type TableConstraint struct {
	Name       ID
	Constraint Constraint
	Position   Position
}

func (tc TableConstraint) gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32() Position {
	return tc.Position
}
func (tc *TableConstraint) gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32() {
	switch c := tc.Constraint.(type) {
	case ForeignKey:
		c.clearOffset()
		tc.Constraint = c
	case Check:
		c.clearOffset()
		tc.Constraint = c
	}
	tc.Position.Offset = 0
}

type Constraint interface {
	isConstraint()
	SQL() string
	Node
}
type Interleave struct {
	Parent   ID
	OnDelete OnDelete
}
type RowDeletionPolicy struct {
	Column  ID
	NumDays int64
}
type CreateIndex struct {
	Name         ID
	Table        ID
	Columns      []KeyPart
	Unique       bool
	NullFiltered bool
	Storing      []ID
	Interleave   ID
	Position     Position
}

func (ci *CreateIndex) gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32() string {
	return fmt.Sprintf("%#v", ci)
}
func (*CreateIndex) gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (ci *CreateIndex) gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32() Position {
	return ci.Position
}
func (ci *CreateIndex) gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32() {
	ci.Position.Offset = 0
}

type CreateView struct {
	Name      ID
	OrReplace bool
	Query     Query
	Position  Position
}

func (cv *CreateView) gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32() string {
	return fmt.Sprintf("%#v", cv)
}
func (*CreateView) gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (cv *CreateView) gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32() Position {
	return cv.Position
}
func (cv *CreateView) gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32() {
	cv.Position.Offset = 0
}

type DropTable struct {
	Name     ID
	Position Position
}

func (dt *DropTable) gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32() string {
	return fmt.Sprintf("%#v", dt)
}
func (*DropTable) gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (dt *DropTable) gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32() Position {
	return dt.Position
}
func (dt *DropTable) gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32() {
	dt.Position.Offset = 0
}

type DropIndex struct {
	Name     ID
	Position Position
}

func (di *DropIndex) gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32() string {
	return fmt.Sprintf("%#v", di)
}
func (*DropIndex) gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (di *DropIndex) gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32() Position {
	return di.Position
}
func (di *DropIndex) gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32() {
	di.Position.Offset = 0
}

type DropView struct {
	Name     ID
	Position Position
}

func (dv *DropView) gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32() string {
	return fmt.Sprintf("%#v", dv)
}
func (*DropView) gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (dv *DropView) gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32() Position {
	return dv.Position
}
func (dv *DropView) gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32() {
	dv.Position.Offset = 0
}

type AlterTable struct {
	Name       ID
	Alteration TableAlteration
	Position   Position
}

func (at *AlterTable) gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32() string {
	return fmt.Sprintf("%#v", at)
}
func (*AlterTable) gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (at *AlterTable) gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32() Position {
	return at.Position
}
func (at *AlterTable) gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32() {
	switch alt := at.Alteration.(type) {
	case AddColumn:
		alt.Def.clearOffset()
		at.Alteration = alt
	case AddConstraint:
		alt.Constraint.clearOffset()
		at.Alteration = alt
	}
	at.Position.Offset = 0
}

type TableAlteration interface {
	isTableAlteration()
	SQL() string
}

func (AddColumn) gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (DropColumn) gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (AddConstraint) gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (DropConstraint) gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (SetOnDelete) gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (AlterColumn) gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (AddRowDeletionPolicy) gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (ReplaceRowDeletionPolicy) gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (DropRowDeletionPolicy) gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type AddColumn struct {
	Def ColumnDef
}
type DropColumn struct {
	Name ID
}
type AddConstraint struct {
	Constraint TableConstraint
}
type DropConstraint struct {
	Name ID
}
type SetOnDelete struct {
	Action OnDelete
}
type AlterColumn struct {
	Name       ID
	Alteration ColumnAlteration
}
type AddRowDeletionPolicy struct {
	RowDeletionPolicy RowDeletionPolicy
}
type ReplaceRowDeletionPolicy struct {
	RowDeletionPolicy RowDeletionPolicy
}
type DropRowDeletionPolicy struct {
}
type ColumnAlteration interface {
	isColumnAlteration()
	SQL() string
}

func (SetColumnType) gologoo__isColumnAlteration_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (SetColumnOptions) gologoo__isColumnAlteration_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type SetColumnType struct {
	Type    Type
	NotNull bool
}
type SetColumnOptions struct {
	Options ColumnOptions
}
type OnDelete int

const (
	NoActionOnDelete OnDelete = iota
	CascadeOnDelete
)

type AlterDatabase struct {
	Name       ID
	Alteration DatabaseAlteration
	Position   Position
}

func (ad *AlterDatabase) gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32() string {
	return fmt.Sprintf("%#v", ad)
}
func (*AlterDatabase) gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (ad *AlterDatabase) gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32() Position {
	return ad.Position
}
func (ad *AlterDatabase) gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32() {
	ad.Position.Offset = 0
}

type DatabaseAlteration interface {
	isDatabaseAlteration()
	SQL() string
}
type SetDatabaseOptions struct {
	Options DatabaseOptions
}

func (SetDatabaseOptions) gologoo__isDatabaseAlteration_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type DatabaseOptions struct {
	OptimizerVersion       *int
	VersionRetentionPeriod *string
	EnableKeyVisualizer    *bool
}
type Delete struct {
	Table ID
	Where BoolExpr
}

func (d *Delete) gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32() string {
	return fmt.Sprintf("%#v", d)
}
func (*Delete) gologoo__isDMLStmt_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type Update struct {
	Table ID
	Items []UpdateItem
	Where BoolExpr
}

func (u *Update) gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32() string {
	return fmt.Sprintf("%#v", u)
}
func (*Update) gologoo__isDMLStmt_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type UpdateItem struct {
	Column ID
	Value  Expr
}
type ColumnDef struct {
	Name      ID
	Type      Type
	NotNull   bool
	Generated Expr
	Options   ColumnOptions
	Position  Position
}

func (cd ColumnDef) gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32() Position {
	return cd.Position
}
func (cd *ColumnDef) gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32() {
	cd.Position.Offset = 0
}

type ColumnOptions struct {
	AllowCommitTimestamp *bool
}
type ForeignKey struct {
	Columns    []ID
	RefTable   ID
	RefColumns []ID
	Position   Position
}

func (fk ForeignKey) gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32() Position {
	return fk.Position
}
func (fk *ForeignKey) gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32() {
	fk.Position.Offset = 0
}
func (ForeignKey) gologoo__isConstraint_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type Check struct {
	Expr     BoolExpr
	Position Position
}

func (c Check) gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32() Position {
	return c.Position
}
func (c *Check) gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32() {
	c.Position.Offset = 0
}
func (Check) gologoo__isConstraint_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type Type struct {
	Array bool
	Base  TypeBase
	Len   int64
}

const MaxLen = math.MaxInt64

type TypeBase int

const (
	Bool TypeBase = iota
	Int64
	Float64
	Numeric
	String
	Bytes
	Date
	Timestamp
	JSON
)

type KeyPart struct {
	Column ID
	Desc   bool
}
type Query struct {
	Select        Select
	Order         []Order
	Limit, Offset LiteralOrParam
}
type Select struct {
	Distinct     bool
	List         []Expr
	From         []SelectFrom
	Where        BoolExpr
	GroupBy      []Expr
	TableSamples []*TableSample
	ListAliases  []ID
}
type SelectFrom interface {
	isSelectFrom()
	SQL() string
}
type SelectFromTable struct {
	Table ID
	Alias ID
	Hints map[string]string
}

func (SelectFromTable) gologoo__isSelectFrom_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type SelectFromJoin struct {
	Type     JoinType
	LHS, RHS SelectFrom
	On       BoolExpr
	Using    []ID
	Hints    map[string]string
}

func (SelectFromJoin) gologoo__isSelectFrom_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type JoinType int

const (
	InnerJoin JoinType = iota
	CrossJoin
	FullJoin
	LeftJoin
	RightJoin
)

type SelectFromUnnest struct {
	Expr  Expr
	Alias ID
}

func (SelectFromUnnest) gologoo__isSelectFrom_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type Order struct {
	Expr Expr
	Desc bool
}
type TableSample struct {
	Method   TableSampleMethod
	Size     Expr
	SizeType TableSampleSizeType
}
type TableSampleMethod int

const (
	Bernoulli TableSampleMethod = iota
	Reservoir
)

type TableSampleSizeType int

const (
	PercentTableSample TableSampleSizeType = iota
	RowsTableSample
)

type BoolExpr interface {
	isBoolExpr()
	Expr
}
type Expr interface {
	isExpr()
	SQL() string
	addSQL(*strings.Builder)
}
type LiteralOrParam interface {
	isLiteralOrParam()
	SQL() string
}
type ArithOp struct {
	Op       ArithOperator
	LHS, RHS Expr
}

func (ArithOp) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type ArithOperator int

const (
	Neg ArithOperator = iota
	Plus
	BitNot
	Mul
	Div
	Concat
	Add
	Sub
	BitShl
	BitShr
	BitAnd
	BitXor
	BitOr
)

type LogicalOp struct {
	Op       LogicalOperator
	LHS, RHS BoolExpr
}

func (LogicalOp) gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (LogicalOp) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type LogicalOperator int

const (
	And LogicalOperator = iota
	Or
	Not
)

type ComparisonOp struct {
	Op       ComparisonOperator
	LHS, RHS Expr
	RHS2     Expr
}

func (ComparisonOp) gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (ComparisonOp) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type ComparisonOperator int

const (
	Lt ComparisonOperator = iota
	Le
	Gt
	Ge
	Eq
	Ne
	Like
	NotLike
	Between
	NotBetween
)

type InOp struct {
	LHS    Expr
	Neg    bool
	RHS    []Expr
	Unnest bool
}

func (InOp) gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (InOp) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type IsOp struct {
	LHS Expr
	Neg bool
	RHS IsExpr
}

func (IsOp) gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (IsOp) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type IsExpr interface {
	isIsExpr()
	Expr
}
type PathExp []ID

func (PathExp) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type Func struct {
	Name string
	Args []Expr
}

func (Func) gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (Func) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type TypedExpr struct {
	Type Type
	Expr Expr
}

func (TypedExpr) gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (TypedExpr) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type ExtractExpr struct {
	Part string
	Type Type
	Expr Expr
}

func (ExtractExpr) gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (ExtractExpr) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type AtTimeZoneExpr struct {
	Expr Expr
	Type Type
	Zone string
}

func (AtTimeZoneExpr) gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (AtTimeZoneExpr) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type Paren struct {
	Expr Expr
}

func (Paren) gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (Paren) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type Array []Expr

func (Array) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type ID string

func (ID) gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (ID) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type Param string

func (Param) gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (Param) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (Param) gologoo__isLiteralOrParam_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type BoolLiteral bool

const (
	True  = BoolLiteral(true)
	False = BoolLiteral(false)
)

func (BoolLiteral) gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (BoolLiteral) gologoo__isIsExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (BoolLiteral) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type NullLiteral int

const Null = NullLiteral(0)

func (NullLiteral) gologoo__isIsExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (NullLiteral) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type IntegerLiteral int64

func (IntegerLiteral) gologoo__isLiteralOrParam_9a673e2e5e9cb6f58c17acc311e36f32() {
}
func (IntegerLiteral) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type FloatLiteral float64

func (FloatLiteral) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type StringLiteral string

func (StringLiteral) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type BytesLiteral string

func (BytesLiteral) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type DateLiteral civil.Date

func (DateLiteral) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type TimestampLiteral time.Time

func (TimestampLiteral) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type StarExpr int

const Star = StarExpr(0)

func (StarExpr) gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32() {
}

type DDL struct {
	List     []DDLStmt
	Filename string
	Comments []*Comment
}

func (d *DDL) gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32() {
	for _, stmt := range d.List {
		stmt.clearOffset()
	}
	for _, c := range d.Comments {
		c.clearOffset()
	}
}

type DDLStmt interface {
	isDDLStmt()
	clearOffset()
	SQL() string
	Node
}
type DMLStmt interface {
	isDMLStmt()
	SQL() string
}
type Comment struct {
	Marker     string
	Isolated   bool
	Start, End Position
	Text       []string
}

func (c *Comment) gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32() string {
	return fmt.Sprintf("%#v", c)
}
func (c *Comment) gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32() Position {
	return c.Start
}
func (c *Comment) gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32() {
	c.Start.Offset, c.End.Offset = 0, 0
}

type Node interface {
	Pos() Position
}
type Position struct {
	Line   int
	Offset int
}

func (pos Position) gologoo__IsValid_9a673e2e5e9cb6f58c17acc311e36f32() bool {
	return pos.Line > 0
}
func (pos Position) gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32() string {
	if pos.Line == 0 {
		return ":<invalid>"
	}
	return fmt.Sprintf(":%d", pos.Line)
}
func (ddl *DDL) gologoo__LeadingComment_9a673e2e5e9cb6f58c17acc311e36f32(n Node) *Comment {
	lineEnd := n.Pos().Line - 1
	ci := sort.Search(len(ddl.Comments), func(i int) bool {
		return ddl.Comments[i].End.Line >= lineEnd
	})
	if ci >= len(ddl.Comments) || ddl.Comments[ci].End.Line != lineEnd {
		return nil
	}
	if !ddl.Comments[ci].Isolated {
		return nil
	}
	return ddl.Comments[ci]
}
func (ddl *DDL) gologoo__InlineComment_9a673e2e5e9cb6f58c17acc311e36f32(n Node) *Comment {
	pos := n.Pos()
	ci := sort.Search(len(ddl.Comments), func(i int) bool {
		return ddl.Comments[i].Start.Line >= pos.Line
	})
	if ci >= len(ddl.Comments) {
		return nil
	}
	c := ddl.Comments[ci]
	if c.Start.Line != pos.Line {
		return nil
	}
	if c.Start.Line != c.End.Line || len(c.Text) != 1 {
		return nil
	}
	return c
}
func (ct *CreateTable) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := ct.gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *CreateTable) isDDLStmt() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (ct *CreateTable) Pos() Position {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := ct.gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ct *CreateTable) clearOffset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	ct.gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (tc TableConstraint) Pos() Position {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := tc.gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (tc *TableConstraint) clearOffset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	tc.gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (ci *CreateIndex) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := ci.gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *CreateIndex) isDDLStmt() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (ci *CreateIndex) Pos() Position {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := ci.gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ci *CreateIndex) clearOffset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	ci.gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (cv *CreateView) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := cv.gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *CreateView) isDDLStmt() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (cv *CreateView) Pos() Position {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := cv.gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (cv *CreateView) clearOffset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	cv.gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (dt *DropTable) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := dt.gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *DropTable) isDDLStmt() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (dt *DropTable) Pos() Position {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := dt.gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (dt *DropTable) clearOffset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	dt.gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (di *DropIndex) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := di.gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *DropIndex) isDDLStmt() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (di *DropIndex) Pos() Position {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := di.gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (di *DropIndex) clearOffset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	di.gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (dv *DropView) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := dv.gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *DropView) isDDLStmt() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (dv *DropView) Pos() Position {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := dv.gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (dv *DropView) clearOffset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	dv.gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (at *AlterTable) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := at.gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *AlterTable) isDDLStmt() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (at *AlterTable) Pos() Position {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := at.gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (at *AlterTable) clearOffset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	at.gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv AddColumn) isTableAlteration() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv DropColumn) isTableAlteration() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv AddConstraint) isTableAlteration() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv DropConstraint) isTableAlteration() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv SetOnDelete) isTableAlteration() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv AlterColumn) isTableAlteration() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv AddRowDeletionPolicy) isTableAlteration() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv ReplaceRowDeletionPolicy) isTableAlteration() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv DropRowDeletionPolicy) isTableAlteration() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isTableAlteration_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv SetColumnType) isColumnAlteration() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isColumnAlteration_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isColumnAlteration_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv SetColumnOptions) isColumnAlteration() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isColumnAlteration_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isColumnAlteration_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (ad *AlterDatabase) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := ad.gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *AlterDatabase) isDDLStmt() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isDDLStmt_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (ad *AlterDatabase) Pos() Position {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := ad.gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ad *AlterDatabase) clearOffset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	ad.gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv SetDatabaseOptions) isDatabaseAlteration() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isDatabaseAlteration_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isDatabaseAlteration_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (d *Delete) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := d.gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *Delete) isDMLStmt() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isDMLStmt_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isDMLStmt_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (u *Update) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := u.gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv *Update) isDMLStmt() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isDMLStmt_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isDMLStmt_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (cd ColumnDef) Pos() Position {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := cd.gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (cd *ColumnDef) clearOffset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	cd.gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (fk ForeignKey) Pos() Position {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := fk.gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (fk *ForeignKey) clearOffset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	fk.gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv ForeignKey) isConstraint() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isConstraint_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isConstraint_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (c Check) Pos() Position {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *Check) clearOffset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	c.gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv Check) isConstraint() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isConstraint_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isConstraint_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv SelectFromTable) isSelectFrom() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isSelectFrom_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isSelectFrom_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv SelectFromJoin) isSelectFrom() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isSelectFrom_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isSelectFrom_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv SelectFromUnnest) isSelectFrom() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isSelectFrom_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isSelectFrom_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv ArithOp) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv LogicalOp) isBoolExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv LogicalOp) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv ComparisonOp) isBoolExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv ComparisonOp) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv InOp) isBoolExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv InOp) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv IsOp) isBoolExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv IsOp) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv PathExp) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv Func) isBoolExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv Func) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv TypedExpr) isBoolExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv TypedExpr) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv ExtractExpr) isBoolExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv ExtractExpr) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv AtTimeZoneExpr) isBoolExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv AtTimeZoneExpr) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv Paren) isBoolExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv Paren) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv Array) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv ID) isBoolExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv ID) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv Param) isBoolExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv Param) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv Param) isLiteralOrParam() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isLiteralOrParam_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isLiteralOrParam_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv BoolLiteral) isBoolExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isBoolExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv BoolLiteral) isIsExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isIsExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isIsExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv BoolLiteral) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv NullLiteral) isIsExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isIsExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isIsExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv NullLiteral) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv IntegerLiteral) isLiteralOrParam() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isLiteralOrParam_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isLiteralOrParam_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv IntegerLiteral) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv FloatLiteral) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv StringLiteral) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv BytesLiteral) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv DateLiteral) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv TimestampLiteral) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (recv StarExpr) isExpr() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	recv.gologoo__isExpr_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (d *DDL) clearOffset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	d.gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (c *Comment) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *Comment) Pos() Position {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__Pos_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *Comment) clearOffset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	c.gologoo__clearOffset_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (pos Position) IsValid() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__IsValid_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := pos.gologoo__IsValid_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (pos Position) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : (none)\n")
	r0 := pos.gologoo__String_9a673e2e5e9cb6f58c17acc311e36f32()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ddl *DDL) LeadingComment(n Node) *Comment {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__LeadingComment_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : %v\n", n)
	r0 := ddl.gologoo__LeadingComment_9a673e2e5e9cb6f58c17acc311e36f32(n)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (ddl *DDL) InlineComment(n Node) *Comment {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__InlineComment_9a673e2e5e9cb6f58c17acc311e36f32")
	log.Printf("Input : %v\n", n)
	r0 := ddl.gologoo__InlineComment_9a673e2e5e9cb6f58c17acc311e36f32(n)
	log.Printf("Output: %v\n", r0)
	return r0
}
