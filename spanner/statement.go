package spanner

import (
	"fmt"
	proto3 "github.com/golang/protobuf/ptypes/struct"
	structpb "github.com/golang/protobuf/ptypes/struct"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc/codes"
	"log"
)

type Statement struct {
	SQL    string
	Params map[string]interface {
	}
}

func gologoo__NewStatement_bad61f76b47d0204298c4f04080f8ec9(sql string) Statement {
	return Statement{SQL: sql, Params: map[string]interface {
	}{}}
}
func (s *Statement) gologoo__convertParams_bad61f76b47d0204298c4f04080f8ec9() (*structpb.Struct, map[string]*sppb.Type, error) {
	params := &proto3.Struct{Fields: map[string]*proto3.Value{}}
	paramTypes := map[string]*sppb.Type{}
	for k, v := range s.Params {
		val, t, err := encodeValue(v)
		if err != nil {
			return nil, nil, errBindParam(k, v, err)
		}
		params.Fields[k] = val
		if t != nil {
			paramTypes[k] = t
		}
	}
	return params, paramTypes, nil
}
func gologoo__errBindParam_bad61f76b47d0204298c4f04080f8ec9(k string, v interface {
}, err error) error {
	if err == nil {
		return nil
	}
	var se *Error
	if !errorAs(err, &se) {
		return spannerErrorf(codes.InvalidArgument, "failed to bind query parameter(name: %q, value: %v), error = <%v>", k, v, err)
	}
	se.decorate(fmt.Sprintf("failed to bind query parameter(name: %q, value: %v)", k, v))
	return se
}
func NewStatement(sql string) Statement {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NewStatement_bad61f76b47d0204298c4f04080f8ec9")
	log.Printf("Input : %v\n", sql)
	r0 := gologoo__NewStatement_bad61f76b47d0204298c4f04080f8ec9(sql)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *Statement) convertParams() (*structpb.Struct, map[string]*sppb.Type, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__convertParams_bad61f76b47d0204298c4f04080f8ec9")
	log.Printf("Input : (none)\n")
	r0, r1, r2 := s.gologoo__convertParams_bad61f76b47d0204298c4f04080f8ec9()
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
func errBindParam(k string, v interface {
}, err error) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errBindParam_bad61f76b47d0204298c4f04080f8ec9")
	log.Printf("Input : %v %v %v\n", k, v, err)
	r0 := gologoo__errBindParam_bad61f76b47d0204298c4f04080f8ec9(k, v, err)
	log.Printf("Output: %v\n", r0)
	return r0
}
