package spanner

import (
	"fmt"
	"reflect"
	proto3 "github.com/golang/protobuf/ptypes/struct"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc/codes"
	"log"
)

type Row struct {
	fields []*sppb.StructType_Field
	vals   []*proto3.Value
}

func (r *Row) gologoo__String_b293e741bc0fccb5d706cdc85160c281() string {
	return fmt.Sprintf("{fields: %s, values: %s}", r.fields, r.vals)
}
func gologoo__errNamesValuesMismatch_b293e741bc0fccb5d706cdc85160c281(columnNames []string, columnValues []interface {
}) error {
	return spannerErrorf(codes.FailedPrecondition, "different number of names(%v) and values(%v)", len(columnNames), len(columnValues))
}
func gologoo__NewRow_b293e741bc0fccb5d706cdc85160c281(columnNames []string, columnValues []interface {
}) (*Row, error) {
	if len(columnValues) != len(columnNames) {
		return nil, errNamesValuesMismatch(columnNames, columnValues)
	}
	r := Row{fields: make([]*sppb.StructType_Field, len(columnValues)), vals: make([]*proto3.Value, len(columnValues))}
	for i := range columnValues {
		val, typ, err := encodeValue(columnValues[i])
		if err != nil {
			return nil, err
		}
		r.fields[i] = &sppb.StructType_Field{Name: columnNames[i], Type: typ}
		r.vals[i] = val
	}
	return &r, nil
}
func (r *Row) gologoo__Size_b293e741bc0fccb5d706cdc85160c281() int {
	return len(r.fields)
}
func (r *Row) gologoo__ColumnName_b293e741bc0fccb5d706cdc85160c281(i int) string {
	if i < 0 || i >= len(r.fields) {
		return ""
	}
	return r.fields[i].Name
}
func (r *Row) gologoo__ColumnIndex_b293e741bc0fccb5d706cdc85160c281(name string) (int, error) {
	found := false
	var index int
	if len(r.vals) != len(r.fields) {
		return 0, errFieldsMismatchVals(r)
	}
	for i, f := range r.fields {
		if f == nil {
			return 0, errNilColType(i)
		}
		if name == f.Name {
			if found {
				return 0, errDupColName(name)
			}
			found = true
			index = i
		}
	}
	if !found {
		return 0, errColNotFound(name)
	}
	return index, nil
}
func (r *Row) gologoo__ColumnNames_b293e741bc0fccb5d706cdc85160c281() []string {
	var n []string
	for _, c := range r.fields {
		n = append(n, c.Name)
	}
	return n
}
func gologoo__errColIdxOutOfRange_b293e741bc0fccb5d706cdc85160c281(i int, r *Row) error {
	return spannerErrorf(codes.OutOfRange, "column index %d out of range [0,%d)", i, len(r.vals))
}
func gologoo__errDecodeColumn_b293e741bc0fccb5d706cdc85160c281(i int, err error) error {
	if err == nil {
		return nil
	}
	var se *Error
	if !errorAs(err, &se) {
		return spannerErrorf(codes.InvalidArgument, "failed to decode column %v, error = <%v>", i, err)
	}
	se.decorate(fmt.Sprintf("failed to decode column %v", i))
	return se
}
func gologoo__errFieldsMismatchVals_b293e741bc0fccb5d706cdc85160c281(r *Row) error {
	return spannerErrorf(codes.FailedPrecondition, "row has different number of fields(%v) and values(%v)", len(r.fields), len(r.vals))
}
func gologoo__errNilColType_b293e741bc0fccb5d706cdc85160c281(i int) error {
	return spannerErrorf(codes.FailedPrecondition, "column(%v)'s type is nil", i)
}
func (r *Row) gologoo__Column_b293e741bc0fccb5d706cdc85160c281(i int, ptr interface {
}) error {
	if len(r.vals) != len(r.fields) {
		return errFieldsMismatchVals(r)
	}
	if i < 0 || i >= len(r.fields) {
		return errColIdxOutOfRange(i, r)
	}
	if r.fields[i] == nil {
		return errNilColType(i)
	}
	if err := decodeValue(r.vals[i], r.fields[i].Type, ptr); err != nil {
		return errDecodeColumn(i, err)
	}
	return nil
}
func gologoo__errDupColName_b293e741bc0fccb5d706cdc85160c281(n string) error {
	return spannerErrorf(codes.FailedPrecondition, "ambiguous column name %q", n)
}
func gologoo__errColNotFound_b293e741bc0fccb5d706cdc85160c281(n string) error {
	return spannerErrorf(codes.NotFound, "column %q not found", n)
}
func (r *Row) gologoo__ColumnByName_b293e741bc0fccb5d706cdc85160c281(name string, ptr interface {
}) error {
	index, err := r.ColumnIndex(name)
	if err != nil {
		return err
	}
	return r.Column(index, ptr)
}
func gologoo__errNumOfColValue_b293e741bc0fccb5d706cdc85160c281(n int, r *Row) error {
	return spannerErrorf(codes.InvalidArgument, "Columns(): number of arguments (%d) does not match row size (%d)", n, len(r.vals))
}
func (r *Row) gologoo__Columns_b293e741bc0fccb5d706cdc85160c281(ptrs ...interface {
}) error {
	if len(ptrs) != len(r.vals) {
		return errNumOfColValue(len(ptrs), r)
	}
	if len(r.vals) != len(r.fields) {
		return errFieldsMismatchVals(r)
	}
	for i, p := range ptrs {
		if p == nil {
			continue
		}
		if err := r.Column(i, p); err != nil {
			return err
		}
	}
	return nil
}
func gologoo__errToStructArgType_b293e741bc0fccb5d706cdc85160c281(p interface {
}) error {
	return spannerErrorf(codes.InvalidArgument, "ToStruct(): type %T is not a valid pointer to Go struct", p)
}
func (r *Row) gologoo__ToStruct_b293e741bc0fccb5d706cdc85160c281(p interface {
}) error {
	if t := reflect.TypeOf(p); t == nil || t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errToStructArgType(p)
	}
	if len(r.vals) != len(r.fields) {
		return errFieldsMismatchVals(r)
	}
	return decodeStruct(&sppb.StructType{Fields: r.fields}, &proto3.ListValue{Values: r.vals}, p, false)
}
func (r *Row) gologoo__ToStructLenient_b293e741bc0fccb5d706cdc85160c281(p interface {
}) error {
	if t := reflect.TypeOf(p); t == nil || t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errToStructArgType(p)
	}
	if len(r.vals) != len(r.fields) {
		return errFieldsMismatchVals(r)
	}
	return decodeStruct(&sppb.StructType{Fields: r.fields}, &proto3.ListValue{Values: r.vals}, p, true)
}
func (r *Row) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : (none)\n")
	r0 := r.gologoo__String_b293e741bc0fccb5d706cdc85160c281()
	log.Printf("Output: %v\n", r0)
	return r0
}
func errNamesValuesMismatch(columnNames []string, columnValues []interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errNamesValuesMismatch_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : %v %v\n", columnNames, columnValues)
	r0 := gologoo__errNamesValuesMismatch_b293e741bc0fccb5d706cdc85160c281(columnNames, columnValues)
	log.Printf("Output: %v\n", r0)
	return r0
}
func NewRow(columnNames []string, columnValues []interface {
}) (*Row, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NewRow_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : %v %v\n", columnNames, columnValues)
	r0, r1 := gologoo__NewRow_b293e741bc0fccb5d706cdc85160c281(columnNames, columnValues)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (r *Row) Size() int {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Size_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : (none)\n")
	r0 := r.gologoo__Size_b293e741bc0fccb5d706cdc85160c281()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (r *Row) ColumnName(i int) string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ColumnName_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : %v\n", i)
	r0 := r.gologoo__ColumnName_b293e741bc0fccb5d706cdc85160c281(i)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (r *Row) ColumnIndex(name string) (int, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ColumnIndex_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : %v\n", name)
	r0, r1 := r.gologoo__ColumnIndex_b293e741bc0fccb5d706cdc85160c281(name)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (r *Row) ColumnNames() []string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ColumnNames_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : (none)\n")
	r0 := r.gologoo__ColumnNames_b293e741bc0fccb5d706cdc85160c281()
	log.Printf("Output: %v\n", r0)
	return r0
}
func errColIdxOutOfRange(i int, r *Row) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errColIdxOutOfRange_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : %v %v\n", i, r)
	r0 := gologoo__errColIdxOutOfRange_b293e741bc0fccb5d706cdc85160c281(i, r)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errDecodeColumn(i int, err error) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errDecodeColumn_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : %v %v\n", i, err)
	r0 := gologoo__errDecodeColumn_b293e741bc0fccb5d706cdc85160c281(i, err)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errFieldsMismatchVals(r *Row) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errFieldsMismatchVals_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : %v\n", r)
	r0 := gologoo__errFieldsMismatchVals_b293e741bc0fccb5d706cdc85160c281(r)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errNilColType(i int) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errNilColType_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : %v\n", i)
	r0 := gologoo__errNilColType_b293e741bc0fccb5d706cdc85160c281(i)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (r *Row) Column(i int, ptr interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Column_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : %v %v\n", i, ptr)
	r0 := r.gologoo__Column_b293e741bc0fccb5d706cdc85160c281(i, ptr)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errDupColName(n string) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errDupColName_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : %v\n", n)
	r0 := gologoo__errDupColName_b293e741bc0fccb5d706cdc85160c281(n)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errColNotFound(n string) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errColNotFound_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : %v\n", n)
	r0 := gologoo__errColNotFound_b293e741bc0fccb5d706cdc85160c281(n)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (r *Row) ColumnByName(name string, ptr interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ColumnByName_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : %v %v\n", name, ptr)
	r0 := r.gologoo__ColumnByName_b293e741bc0fccb5d706cdc85160c281(name, ptr)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errNumOfColValue(n int, r *Row) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errNumOfColValue_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : %v %v\n", n, r)
	r0 := gologoo__errNumOfColValue_b293e741bc0fccb5d706cdc85160c281(n, r)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (r *Row) Columns(ptrs ...interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Columns_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : %v\n", ptrs)
	r0 := r.gologoo__Columns_b293e741bc0fccb5d706cdc85160c281(ptrs...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errToStructArgType(p interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errToStructArgType_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : %v\n", p)
	r0 := gologoo__errToStructArgType_b293e741bc0fccb5d706cdc85160c281(p)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (r *Row) ToStruct(p interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ToStruct_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : %v\n", p)
	r0 := r.gologoo__ToStruct_b293e741bc0fccb5d706cdc85160c281(p)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (r *Row) ToStructLenient(p interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ToStructLenient_b293e741bc0fccb5d706cdc85160c281")
	log.Printf("Input : %v\n", p)
	r0 := r.gologoo__ToStructLenient_b293e741bc0fccb5d706cdc85160c281(p)
	log.Printf("Output: %v\n", r0)
	return r0
}
