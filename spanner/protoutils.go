package spanner

import (
	"encoding/base64"
	"math/big"
	"strconv"
	"time"
	"cloud.google.com/go/civil"
	proto3 "github.com/golang/protobuf/ptypes/struct"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"log"
)

func gologoo__stringProto_39faf8fa5af488454e09be583d7c271d(s string) *proto3.Value {
	return &proto3.Value{Kind: stringKind(s)}
}
func gologoo__stringKind_39faf8fa5af488454e09be583d7c271d(s string) *proto3.Value_StringValue {
	return &proto3.Value_StringValue{StringValue: s}
}
func gologoo__stringType_39faf8fa5af488454e09be583d7c271d() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_STRING}
}
func gologoo__boolProto_39faf8fa5af488454e09be583d7c271d(b bool) *proto3.Value {
	return &proto3.Value{Kind: &proto3.Value_BoolValue{BoolValue: b}}
}
func gologoo__boolType_39faf8fa5af488454e09be583d7c271d() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_BOOL}
}
func gologoo__intProto_39faf8fa5af488454e09be583d7c271d(n int64) *proto3.Value {
	return &proto3.Value{Kind: &proto3.Value_StringValue{StringValue: strconv.FormatInt(n, 10)}}
}
func gologoo__intType_39faf8fa5af488454e09be583d7c271d() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_INT64}
}
func gologoo__floatProto_39faf8fa5af488454e09be583d7c271d(n float64) *proto3.Value {
	return &proto3.Value{Kind: &proto3.Value_NumberValue{NumberValue: n}}
}
func gologoo__floatType_39faf8fa5af488454e09be583d7c271d() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_FLOAT64}
}
func gologoo__numericProto_39faf8fa5af488454e09be583d7c271d(n *big.Rat) *proto3.Value {
	return &proto3.Value{Kind: &proto3.Value_StringValue{StringValue: NumericString(n)}}
}
func gologoo__numericType_39faf8fa5af488454e09be583d7c271d() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_NUMERIC}
}
func gologoo__pgNumericType_39faf8fa5af488454e09be583d7c271d() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_NUMERIC, TypeAnnotation: sppb.TypeAnnotationCode_PG_NUMERIC}
}
func gologoo__jsonType_39faf8fa5af488454e09be583d7c271d() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_JSON}
}
func gologoo__bytesProto_39faf8fa5af488454e09be583d7c271d(b []byte) *proto3.Value {
	return &proto3.Value{Kind: &proto3.Value_StringValue{StringValue: base64.StdEncoding.EncodeToString(b)}}
}
func gologoo__bytesType_39faf8fa5af488454e09be583d7c271d() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_BYTES}
}
func gologoo__timeProto_39faf8fa5af488454e09be583d7c271d(t time.Time) *proto3.Value {
	return stringProto(t.UTC().Format(time.RFC3339Nano))
}
func gologoo__timeType_39faf8fa5af488454e09be583d7c271d() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_TIMESTAMP}
}
func gologoo__dateProto_39faf8fa5af488454e09be583d7c271d(d civil.Date) *proto3.Value {
	return stringProto(d.String())
}
func gologoo__dateType_39faf8fa5af488454e09be583d7c271d() *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_DATE}
}
func gologoo__listProto_39faf8fa5af488454e09be583d7c271d(p ...*proto3.Value) *proto3.Value {
	return &proto3.Value{Kind: &proto3.Value_ListValue{ListValue: &proto3.ListValue{Values: p}}}
}
func gologoo__listValueProto_39faf8fa5af488454e09be583d7c271d(p ...*proto3.Value) *proto3.ListValue {
	return &proto3.ListValue{Values: p}
}
func gologoo__listType_39faf8fa5af488454e09be583d7c271d(t *sppb.Type) *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_ARRAY, ArrayElementType: t}
}
func gologoo__mkField_39faf8fa5af488454e09be583d7c271d(n string, t *sppb.Type) *sppb.StructType_Field {
	return &sppb.StructType_Field{Name: n, Type: t}
}
func gologoo__structType_39faf8fa5af488454e09be583d7c271d(fields ...*sppb.StructType_Field) *sppb.Type {
	return &sppb.Type{Code: sppb.TypeCode_STRUCT, StructType: &sppb.StructType{Fields: fields}}
}
func gologoo__nullProto_39faf8fa5af488454e09be583d7c271d() *proto3.Value {
	return &proto3.Value{Kind: &proto3.Value_NullValue{NullValue: proto3.NullValue_NULL_VALUE}}
}
func stringProto(s string) *proto3.Value {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__stringProto_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : %v\n", s)
	r0 := gologoo__stringProto_39faf8fa5af488454e09be583d7c271d(s)
	log.Printf("Output: %v\n", r0)
	return r0
}
func stringKind(s string) *proto3.Value_StringValue {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__stringKind_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : %v\n", s)
	r0 := gologoo__stringKind_39faf8fa5af488454e09be583d7c271d(s)
	log.Printf("Output: %v\n", r0)
	return r0
}
func stringType() *sppb.Type {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__stringType_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : (none)\n")
	r0 := gologoo__stringType_39faf8fa5af488454e09be583d7c271d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func boolProto(b bool) *proto3.Value {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__boolProto_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : %v\n", b)
	r0 := gologoo__boolProto_39faf8fa5af488454e09be583d7c271d(b)
	log.Printf("Output: %v\n", r0)
	return r0
}
func boolType() *sppb.Type {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__boolType_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : (none)\n")
	r0 := gologoo__boolType_39faf8fa5af488454e09be583d7c271d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func intProto(n int64) *proto3.Value {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__intProto_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : %v\n", n)
	r0 := gologoo__intProto_39faf8fa5af488454e09be583d7c271d(n)
	log.Printf("Output: %v\n", r0)
	return r0
}
func intType() *sppb.Type {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__intType_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : (none)\n")
	r0 := gologoo__intType_39faf8fa5af488454e09be583d7c271d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func floatProto(n float64) *proto3.Value {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__floatProto_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : %v\n", n)
	r0 := gologoo__floatProto_39faf8fa5af488454e09be583d7c271d(n)
	log.Printf("Output: %v\n", r0)
	return r0
}
func floatType() *sppb.Type {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__floatType_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : (none)\n")
	r0 := gologoo__floatType_39faf8fa5af488454e09be583d7c271d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func numericProto(n *big.Rat) *proto3.Value {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__numericProto_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : %v\n", n)
	r0 := gologoo__numericProto_39faf8fa5af488454e09be583d7c271d(n)
	log.Printf("Output: %v\n", r0)
	return r0
}
func numericType() *sppb.Type {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__numericType_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : (none)\n")
	r0 := gologoo__numericType_39faf8fa5af488454e09be583d7c271d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func pgNumericType() *sppb.Type {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__pgNumericType_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : (none)\n")
	r0 := gologoo__pgNumericType_39faf8fa5af488454e09be583d7c271d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func jsonType() *sppb.Type {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__jsonType_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : (none)\n")
	r0 := gologoo__jsonType_39faf8fa5af488454e09be583d7c271d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func bytesProto(b []byte) *proto3.Value {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__bytesProto_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : %v\n", b)
	r0 := gologoo__bytesProto_39faf8fa5af488454e09be583d7c271d(b)
	log.Printf("Output: %v\n", r0)
	return r0
}
func bytesType() *sppb.Type {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__bytesType_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : (none)\n")
	r0 := gologoo__bytesType_39faf8fa5af488454e09be583d7c271d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func timeProto(t time.Time) *proto3.Value {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__timeProto_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : %v\n", t)
	r0 := gologoo__timeProto_39faf8fa5af488454e09be583d7c271d(t)
	log.Printf("Output: %v\n", r0)
	return r0
}
func timeType() *sppb.Type {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__timeType_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : (none)\n")
	r0 := gologoo__timeType_39faf8fa5af488454e09be583d7c271d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func dateProto(d civil.Date) *proto3.Value {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__dateProto_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : %v\n", d)
	r0 := gologoo__dateProto_39faf8fa5af488454e09be583d7c271d(d)
	log.Printf("Output: %v\n", r0)
	return r0
}
func dateType() *sppb.Type {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__dateType_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : (none)\n")
	r0 := gologoo__dateType_39faf8fa5af488454e09be583d7c271d()
	log.Printf("Output: %v\n", r0)
	return r0
}
func listProto(p ...*proto3.Value) *proto3.Value {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__listProto_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : %v\n", p)
	r0 := gologoo__listProto_39faf8fa5af488454e09be583d7c271d(p...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func listValueProto(p ...*proto3.Value) *proto3.ListValue {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__listValueProto_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : %v\n", p)
	r0 := gologoo__listValueProto_39faf8fa5af488454e09be583d7c271d(p...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func listType(t *sppb.Type) *sppb.Type {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__listType_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : %v\n", t)
	r0 := gologoo__listType_39faf8fa5af488454e09be583d7c271d(t)
	log.Printf("Output: %v\n", r0)
	return r0
}
func mkField(n string, t *sppb.Type) *sppb.StructType_Field {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__mkField_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : %v %v\n", n, t)
	r0 := gologoo__mkField_39faf8fa5af488454e09be583d7c271d(n, t)
	log.Printf("Output: %v\n", r0)
	return r0
}
func structType(fields ...*sppb.StructType_Field) *sppb.Type {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__structType_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : %v\n", fields)
	r0 := gologoo__structType_39faf8fa5af488454e09be583d7c271d(fields...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func nullProto() *proto3.Value {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__nullProto_39faf8fa5af488454e09be583d7c271d")
	log.Printf("Input : (none)\n")
	r0 := gologoo__nullProto_39faf8fa5af488454e09be583d7c271d()
	log.Printf("Output: %v\n", r0)
	return r0
}
