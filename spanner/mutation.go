package spanner

import (
	"reflect"
	proto3 "github.com/golang/protobuf/ptypes/struct"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc/codes"
	"log"
)

type op int

const (
	opDelete op = iota
	opInsert
	opInsertOrUpdate
	opReplace
	opUpdate
)

type Mutation struct {
	op      op
	table   string
	keySet  KeySet
	columns []string
	values  []interface {
	}
}

func gologoo__mapToMutationParams_cdb1dee679c01e0eeca12fe18c1c7a89(in map[string]interface {
}) ([]string, []interface {
}) {
	cols := []string{}
	vals := []interface {
	}{}
	for k, v := range in {
		cols = append(cols, k)
		vals = append(vals, v)
	}
	return cols, vals
}
func gologoo__errNotStruct_cdb1dee679c01e0eeca12fe18c1c7a89(in interface {
}) error {
	return spannerErrorf(codes.InvalidArgument, "%T is not a go struct type", in)
}
func gologoo__structToMutationParams_cdb1dee679c01e0eeca12fe18c1c7a89(in interface {
}) ([]string, []interface {
}, error) {
	if in == nil {
		return nil, nil, errNotStruct(in)
	}
	v := reflect.ValueOf(in)
	t := v.Type()
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		if v.IsNil() {
			return nil, nil, nil
		}
		v = v.Elem()
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, nil, errNotStruct(in)
	}
	fields, err := fieldCache.Fields(t)
	if err != nil {
		return nil, nil, ToSpannerError(err)
	}
	var cols []string
	var vals []interface {
	}
	for _, f := range fields {
		cols = append(cols, f.Name)
		vals = append(vals, v.FieldByIndex(f.Index).Interface())
	}
	return cols, vals, nil
}
func gologoo__Insert_cdb1dee679c01e0eeca12fe18c1c7a89(table string, cols []string, vals []interface {
}) *Mutation {
	return &Mutation{op: opInsert, table: table, columns: cols, values: vals}
}
func gologoo__InsertMap_cdb1dee679c01e0eeca12fe18c1c7a89(table string, in map[string]interface {
}) *Mutation {
	cols, vals := mapToMutationParams(in)
	return Insert(table, cols, vals)
}
func gologoo__InsertStruct_cdb1dee679c01e0eeca12fe18c1c7a89(table string, in interface {
}) (*Mutation, error) {
	cols, vals, err := structToMutationParams(in)
	if err != nil {
		return nil, err
	}
	return Insert(table, cols, vals), nil
}
func gologoo__Update_cdb1dee679c01e0eeca12fe18c1c7a89(table string, cols []string, vals []interface {
}) *Mutation {
	return &Mutation{op: opUpdate, table: table, columns: cols, values: vals}
}
func gologoo__UpdateMap_cdb1dee679c01e0eeca12fe18c1c7a89(table string, in map[string]interface {
}) *Mutation {
	cols, vals := mapToMutationParams(in)
	return Update(table, cols, vals)
}
func gologoo__UpdateStruct_cdb1dee679c01e0eeca12fe18c1c7a89(table string, in interface {
}) (*Mutation, error) {
	cols, vals, err := structToMutationParams(in)
	if err != nil {
		return nil, err
	}
	return Update(table, cols, vals), nil
}
func gologoo__InsertOrUpdate_cdb1dee679c01e0eeca12fe18c1c7a89(table string, cols []string, vals []interface {
}) *Mutation {
	return &Mutation{op: opInsertOrUpdate, table: table, columns: cols, values: vals}
}
func gologoo__InsertOrUpdateMap_cdb1dee679c01e0eeca12fe18c1c7a89(table string, in map[string]interface {
}) *Mutation {
	cols, vals := mapToMutationParams(in)
	return InsertOrUpdate(table, cols, vals)
}
func gologoo__InsertOrUpdateStruct_cdb1dee679c01e0eeca12fe18c1c7a89(table string, in interface {
}) (*Mutation, error) {
	cols, vals, err := structToMutationParams(in)
	if err != nil {
		return nil, err
	}
	return InsertOrUpdate(table, cols, vals), nil
}
func gologoo__Replace_cdb1dee679c01e0eeca12fe18c1c7a89(table string, cols []string, vals []interface {
}) *Mutation {
	return &Mutation{op: opReplace, table: table, columns: cols, values: vals}
}
func gologoo__ReplaceMap_cdb1dee679c01e0eeca12fe18c1c7a89(table string, in map[string]interface {
}) *Mutation {
	cols, vals := mapToMutationParams(in)
	return Replace(table, cols, vals)
}
func gologoo__ReplaceStruct_cdb1dee679c01e0eeca12fe18c1c7a89(table string, in interface {
}) (*Mutation, error) {
	cols, vals, err := structToMutationParams(in)
	if err != nil {
		return nil, err
	}
	return Replace(table, cols, vals), nil
}
func gologoo__Delete_cdb1dee679c01e0eeca12fe18c1c7a89(table string, ks KeySet) *Mutation {
	return &Mutation{op: opDelete, table: table, keySet: ks}
}
func gologoo__prepareWrite_cdb1dee679c01e0eeca12fe18c1c7a89(table string, columns []string, vals []interface {
}) (*sppb.Mutation_Write, error) {
	v, err := encodeValueArray(vals)
	if err != nil {
		return nil, err
	}
	return &sppb.Mutation_Write{Table: table, Columns: columns, Values: []*proto3.ListValue{v}}, nil
}
func gologoo__errInvdMutationOp_cdb1dee679c01e0eeca12fe18c1c7a89(m Mutation) error {
	return spannerErrorf(codes.InvalidArgument, "Unknown op type: %d", m.op)
}
func (m Mutation) gologoo__proto_cdb1dee679c01e0eeca12fe18c1c7a89() (*sppb.Mutation, error) {
	var pb *sppb.Mutation
	switch m.op {
	case opDelete:
		var kp *sppb.KeySet
		if m.keySet != nil {
			var err error
			kp, err = m.keySet.keySetProto()
			if err != nil {
				return nil, err
			}
		}
		pb = &sppb.Mutation{Operation: &sppb.Mutation_Delete_{Delete: &sppb.Mutation_Delete{Table: m.table, KeySet: kp}}}
	case opInsert:
		w, err := prepareWrite(m.table, m.columns, m.values)
		if err != nil {
			return nil, err
		}
		pb = &sppb.Mutation{Operation: &sppb.Mutation_Insert{Insert: w}}
	case opInsertOrUpdate:
		w, err := prepareWrite(m.table, m.columns, m.values)
		if err != nil {
			return nil, err
		}
		pb = &sppb.Mutation{Operation: &sppb.Mutation_InsertOrUpdate{InsertOrUpdate: w}}
	case opReplace:
		w, err := prepareWrite(m.table, m.columns, m.values)
		if err != nil {
			return nil, err
		}
		pb = &sppb.Mutation{Operation: &sppb.Mutation_Replace{Replace: w}}
	case opUpdate:
		w, err := prepareWrite(m.table, m.columns, m.values)
		if err != nil {
			return nil, err
		}
		pb = &sppb.Mutation{Operation: &sppb.Mutation_Update{Update: w}}
	default:
		return nil, errInvdMutationOp(m)
	}
	return pb, nil
}
func gologoo__mutationsProto_cdb1dee679c01e0eeca12fe18c1c7a89(ms []*Mutation) ([]*sppb.Mutation, error) {
	l := make([]*sppb.Mutation, 0, len(ms))
	for _, m := range ms {
		pb, err := m.proto()
		if err != nil {
			return nil, err
		}
		l = append(l, pb)
	}
	return l, nil
}
func mapToMutationParams(in map[string]interface {
}) ([]string, []interface {
}) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__mapToMutationParams_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v\n", in)
	r0, r1 := gologoo__mapToMutationParams_cdb1dee679c01e0eeca12fe18c1c7a89(in)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func errNotStruct(in interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errNotStruct_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v\n", in)
	r0 := gologoo__errNotStruct_cdb1dee679c01e0eeca12fe18c1c7a89(in)
	log.Printf("Output: %v\n", r0)
	return r0
}
func structToMutationParams(in interface {
}) ([]string, []interface {
}, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__structToMutationParams_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v\n", in)
	r0, r1, r2 := gologoo__structToMutationParams_cdb1dee679c01e0eeca12fe18c1c7a89(in)
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
func Insert(table string, cols []string, vals []interface {
}) *Mutation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Insert_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v %v %v\n", table, cols, vals)
	r0 := gologoo__Insert_cdb1dee679c01e0eeca12fe18c1c7a89(table, cols, vals)
	log.Printf("Output: %v\n", r0)
	return r0
}
func InsertMap(table string, in map[string]interface {
}) *Mutation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__InsertMap_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v %v\n", table, in)
	r0 := gologoo__InsertMap_cdb1dee679c01e0eeca12fe18c1c7a89(table, in)
	log.Printf("Output: %v\n", r0)
	return r0
}
func InsertStruct(table string, in interface {
}) (*Mutation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__InsertStruct_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v %v\n", table, in)
	r0, r1 := gologoo__InsertStruct_cdb1dee679c01e0eeca12fe18c1c7a89(table, in)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func Update(table string, cols []string, vals []interface {
}) *Mutation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Update_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v %v %v\n", table, cols, vals)
	r0 := gologoo__Update_cdb1dee679c01e0eeca12fe18c1c7a89(table, cols, vals)
	log.Printf("Output: %v\n", r0)
	return r0
}
func UpdateMap(table string, in map[string]interface {
}) *Mutation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UpdateMap_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v %v\n", table, in)
	r0 := gologoo__UpdateMap_cdb1dee679c01e0eeca12fe18c1c7a89(table, in)
	log.Printf("Output: %v\n", r0)
	return r0
}
func UpdateStruct(table string, in interface {
}) (*Mutation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UpdateStruct_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v %v\n", table, in)
	r0, r1 := gologoo__UpdateStruct_cdb1dee679c01e0eeca12fe18c1c7a89(table, in)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func InsertOrUpdate(table string, cols []string, vals []interface {
}) *Mutation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__InsertOrUpdate_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v %v %v\n", table, cols, vals)
	r0 := gologoo__InsertOrUpdate_cdb1dee679c01e0eeca12fe18c1c7a89(table, cols, vals)
	log.Printf("Output: %v\n", r0)
	return r0
}
func InsertOrUpdateMap(table string, in map[string]interface {
}) *Mutation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__InsertOrUpdateMap_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v %v\n", table, in)
	r0 := gologoo__InsertOrUpdateMap_cdb1dee679c01e0eeca12fe18c1c7a89(table, in)
	log.Printf("Output: %v\n", r0)
	return r0
}
func InsertOrUpdateStruct(table string, in interface {
}) (*Mutation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__InsertOrUpdateStruct_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v %v\n", table, in)
	r0, r1 := gologoo__InsertOrUpdateStruct_cdb1dee679c01e0eeca12fe18c1c7a89(table, in)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func Replace(table string, cols []string, vals []interface {
}) *Mutation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Replace_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v %v %v\n", table, cols, vals)
	r0 := gologoo__Replace_cdb1dee679c01e0eeca12fe18c1c7a89(table, cols, vals)
	log.Printf("Output: %v\n", r0)
	return r0
}
func ReplaceMap(table string, in map[string]interface {
}) *Mutation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ReplaceMap_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v %v\n", table, in)
	r0 := gologoo__ReplaceMap_cdb1dee679c01e0eeca12fe18c1c7a89(table, in)
	log.Printf("Output: %v\n", r0)
	return r0
}
func ReplaceStruct(table string, in interface {
}) (*Mutation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ReplaceStruct_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v %v\n", table, in)
	r0, r1 := gologoo__ReplaceStruct_cdb1dee679c01e0eeca12fe18c1c7a89(table, in)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func Delete(table string, ks KeySet) *Mutation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Delete_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v %v\n", table, ks)
	r0 := gologoo__Delete_cdb1dee679c01e0eeca12fe18c1c7a89(table, ks)
	log.Printf("Output: %v\n", r0)
	return r0
}
func prepareWrite(table string, columns []string, vals []interface {
}) (*sppb.Mutation_Write, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__prepareWrite_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v %v %v\n", table, columns, vals)
	r0, r1 := gologoo__prepareWrite_cdb1dee679c01e0eeca12fe18c1c7a89(table, columns, vals)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func errInvdMutationOp(m Mutation) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errInvdMutationOp_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v\n", m)
	r0 := gologoo__errInvdMutationOp_cdb1dee679c01e0eeca12fe18c1c7a89(m)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (m Mutation) proto() (*sppb.Mutation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__proto_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : (none)\n")
	r0, r1 := m.gologoo__proto_cdb1dee679c01e0eeca12fe18c1c7a89()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func mutationsProto(ms []*Mutation) ([]*sppb.Mutation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__mutationsProto_cdb1dee679c01e0eeca12fe18c1c7a89")
	log.Printf("Input : %v\n", ms)
	r0, r1 := gologoo__mutationsProto_cdb1dee679c01e0eeca12fe18c1c7a89(ms)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
