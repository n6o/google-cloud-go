package spanner

import (
	"bytes"
	"fmt"
	"math/big"
	"time"
	"cloud.google.com/go/civil"
	proto3 "github.com/golang/protobuf/ptypes/struct"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc/codes"
	"log"
)

type Key []interface {
}

func gologoo__errInvdKeyPartType_497610871a9228c39fa0d4df4bea64a8(part interface {
}) error {
	return spannerErrorf(codes.InvalidArgument, "key part has unsupported type %T", part)
}
func gologoo__keyPartValue_497610871a9228c39fa0d4df4bea64a8(part interface {
}) (pb *proto3.Value, err error) {
	switch v := part.(type) {
	case int:
		pb, _, err = encodeValue(int64(v))
	case int8:
		pb, _, err = encodeValue(int64(v))
	case int16:
		pb, _, err = encodeValue(int64(v))
	case int32:
		pb, _, err = encodeValue(int64(v))
	case uint8:
		pb, _, err = encodeValue(int64(v))
	case uint16:
		pb, _, err = encodeValue(int64(v))
	case uint32:
		pb, _, err = encodeValue(int64(v))
	case float32:
		pb, _, err = encodeValue(float64(v))
	case int64, float64, NullInt64, NullFloat64, bool, NullBool, []byte, string, NullString, time.Time, civil.Date, NullTime, NullDate, big.Rat, NullNumeric:
		pb, _, err = encodeValue(v)
	case Encoder:
		part, err = v.EncodeSpanner()
		if err != nil {
			return nil, err
		}
		pb, err = keyPartValue(part)
	default:
		return nil, errInvdKeyPartType(v)
	}
	return pb, err
}
func (key Key) gologoo__proto_497610871a9228c39fa0d4df4bea64a8() (*proto3.ListValue, error) {
	lv := &proto3.ListValue{}
	lv.Values = make([]*proto3.Value, 0, len(key))
	for _, part := range key {
		v, err := keyPartValue(part)
		if err != nil {
			return nil, err
		}
		lv.Values = append(lv.Values, v)
	}
	return lv, nil
}
func (key Key) gologoo__keySetProto_497610871a9228c39fa0d4df4bea64a8() (*sppb.KeySet, error) {
	kp, err := key.proto()
	if err != nil {
		return nil, err
	}
	return &sppb.KeySet{Keys: []*proto3.ListValue{kp}}, nil
}
func (key Key) gologoo__String_497610871a9228c39fa0d4df4bea64a8() string {
	b := &bytes.Buffer{}
	fmt.Fprint(b, "(")
	for i, part := range []interface {
	}(key) {
		if i != 0 {
			fmt.Fprint(b, ",")
		}
		key.elemString(b, part)
	}
	fmt.Fprint(b, ")")
	return b.String()
}
func (key Key) gologoo__elemString_497610871a9228c39fa0d4df4bea64a8(b *bytes.Buffer, part interface {
}) {
	switch v := part.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, float32, float64, bool:
		fmt.Fprintf(b, "%v", v)
	case string:
		fmt.Fprintf(b, "%q", v)
	case []byte:
		if v != nil {
			fmt.Fprintf(b, "%q", v)
		} else {
			fmt.Fprint(b, nullString)
		}
	case NullInt64, NullFloat64, NullBool, NullNumeric:
		fmt.Fprintf(b, "%s", v)
	case NullString, NullDate, NullTime:
		if v.(NullableValue).IsNull() {
			fmt.Fprintf(b, "%s", nullString)
		} else {
			fmt.Fprintf(b, "%q", v)
		}
	case civil.Date:
		fmt.Fprintf(b, "%q", v)
	case time.Time:
		fmt.Fprintf(b, "%q", v.Format(time.RFC3339Nano))
	case big.Rat:
		fmt.Fprintf(b, "%v", NumericString(&v))
	case Encoder:
		var err error
		part, err = v.EncodeSpanner()
		if err != nil {
			fmt.Fprintf(b, "error")
		} else {
			key.elemString(b, part)
		}
	default:
		fmt.Fprintf(b, "%v", v)
	}
}
func (key Key) gologoo__AsPrefix_497610871a9228c39fa0d4df4bea64a8() KeyRange {
	return KeyRange{Start: key, End: key, Kind: ClosedClosed}
}

type KeyRangeKind int

const (
	ClosedOpen KeyRangeKind = iota
	ClosedClosed
	OpenClosed
	OpenOpen
)

type KeyRange struct {
	Start, End Key
	Kind       KeyRangeKind
}

func (r KeyRange) gologoo__String_497610871a9228c39fa0d4df4bea64a8() string {
	var left, right string
	switch r.Kind {
	case ClosedClosed:
		left, right = "[", "]"
	case ClosedOpen:
		left, right = "[", ")"
	case OpenClosed:
		left, right = "(", "]"
	case OpenOpen:
		left, right = "(", ")"
	default:
		left, right = "?", "?"
	}
	return fmt.Sprintf("%s%s,%s%s", left, r.Start, r.End, right)
}
func (r KeyRange) gologoo__proto_497610871a9228c39fa0d4df4bea64a8() (*sppb.KeyRange, error) {
	var err error
	var start, end *proto3.ListValue
	pb := &sppb.KeyRange{}
	if start, err = r.Start.proto(); err != nil {
		return nil, err
	}
	if end, err = r.End.proto(); err != nil {
		return nil, err
	}
	if r.Kind == ClosedClosed || r.Kind == ClosedOpen {
		pb.StartKeyType = &sppb.KeyRange_StartClosed{StartClosed: start}
	} else {
		pb.StartKeyType = &sppb.KeyRange_StartOpen{StartOpen: start}
	}
	if r.Kind == ClosedClosed || r.Kind == OpenClosed {
		pb.EndKeyType = &sppb.KeyRange_EndClosed{EndClosed: end}
	} else {
		pb.EndKeyType = &sppb.KeyRange_EndOpen{EndOpen: end}
	}
	return pb, nil
}
func (r KeyRange) gologoo__keySetProto_497610871a9228c39fa0d4df4bea64a8() (*sppb.KeySet, error) {
	rp, err := r.proto()
	if err != nil {
		return nil, err
	}
	return &sppb.KeySet{Ranges: []*sppb.KeyRange{rp}}, nil
}

type KeySet interface {
	keySetProto() (*sppb.KeySet, error)
}

func gologoo__AllKeys_497610871a9228c39fa0d4df4bea64a8() KeySet {
	return all{}
}

type all struct {
}

func (all) gologoo__keySetProto_497610871a9228c39fa0d4df4bea64a8() (*sppb.KeySet, error) {
	return &sppb.KeySet{All: true}, nil
}
func gologoo__KeySets_497610871a9228c39fa0d4df4bea64a8(keySets ...KeySet) KeySet {
	u := make(union, len(keySets))
	copy(u, keySets)
	return u
}
func gologoo__KeySetFromKeys_497610871a9228c39fa0d4df4bea64a8(keys ...Key) KeySet {
	u := make(union, len(keys))
	for i, k := range keys {
		u[i] = k
	}
	return u
}

type union []KeySet

func (u union) gologoo__keySetProto_497610871a9228c39fa0d4df4bea64a8() (*sppb.KeySet, error) {
	upb := &sppb.KeySet{}
	for _, ks := range u {
		pb, err := ks.keySetProto()
		if err != nil {
			return nil, err
		}
		if pb.All {
			return pb, nil
		}
		upb.Keys = append(upb.Keys, pb.Keys...)
		upb.Ranges = append(upb.Ranges, pb.Ranges...)
	}
	return upb, nil
}
func errInvdKeyPartType(part interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errInvdKeyPartType_497610871a9228c39fa0d4df4bea64a8")
	log.Printf("Input : %v\n", part)
	r0 := gologoo__errInvdKeyPartType_497610871a9228c39fa0d4df4bea64a8(part)
	log.Printf("Output: %v\n", r0)
	return r0
}
func keyPartValue(part interface {
}) (pb *proto3.Value, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__keyPartValue_497610871a9228c39fa0d4df4bea64a8")
	log.Printf("Input : %v\n", part)
	pb, err = gologoo__keyPartValue_497610871a9228c39fa0d4df4bea64a8(part)
	log.Printf("Output: %v %v\n", pb, err)
	return
}
func (key Key) proto() (*proto3.ListValue, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__proto_497610871a9228c39fa0d4df4bea64a8")
	log.Printf("Input : (none)\n")
	r0, r1 := key.gologoo__proto_497610871a9228c39fa0d4df4bea64a8()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (key Key) keySetProto() (*sppb.KeySet, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__keySetProto_497610871a9228c39fa0d4df4bea64a8")
	log.Printf("Input : (none)\n")
	r0, r1 := key.gologoo__keySetProto_497610871a9228c39fa0d4df4bea64a8()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (key Key) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_497610871a9228c39fa0d4df4bea64a8")
	log.Printf("Input : (none)\n")
	r0 := key.gologoo__String_497610871a9228c39fa0d4df4bea64a8()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (key Key) elemString(b *bytes.Buffer, part interface {
}) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__elemString_497610871a9228c39fa0d4df4bea64a8")
	log.Printf("Input : %v %v\n", b, part)
	key.gologoo__elemString_497610871a9228c39fa0d4df4bea64a8(b, part)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (key Key) AsPrefix() KeyRange {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__AsPrefix_497610871a9228c39fa0d4df4bea64a8")
	log.Printf("Input : (none)\n")
	r0 := key.gologoo__AsPrefix_497610871a9228c39fa0d4df4bea64a8()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (r KeyRange) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_497610871a9228c39fa0d4df4bea64a8")
	log.Printf("Input : (none)\n")
	r0 := r.gologoo__String_497610871a9228c39fa0d4df4bea64a8()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (r KeyRange) proto() (*sppb.KeyRange, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__proto_497610871a9228c39fa0d4df4bea64a8")
	log.Printf("Input : (none)\n")
	r0, r1 := r.gologoo__proto_497610871a9228c39fa0d4df4bea64a8()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (r KeyRange) keySetProto() (*sppb.KeySet, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__keySetProto_497610871a9228c39fa0d4df4bea64a8")
	log.Printf("Input : (none)\n")
	r0, r1 := r.gologoo__keySetProto_497610871a9228c39fa0d4df4bea64a8()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func AllKeys() KeySet {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__AllKeys_497610871a9228c39fa0d4df4bea64a8")
	log.Printf("Input : (none)\n")
	r0 := gologoo__AllKeys_497610871a9228c39fa0d4df4bea64a8()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (recv all) keySetProto() (*sppb.KeySet, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__keySetProto_497610871a9228c39fa0d4df4bea64a8")
	log.Printf("Input : (none)\n")
	r0, r1 := recv.gologoo__keySetProto_497610871a9228c39fa0d4df4bea64a8()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func KeySets(keySets ...KeySet) KeySet {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__KeySets_497610871a9228c39fa0d4df4bea64a8")
	log.Printf("Input : %v\n", keySets)
	r0 := gologoo__KeySets_497610871a9228c39fa0d4df4bea64a8(keySets...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func KeySetFromKeys(keys ...Key) KeySet {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__KeySetFromKeys_497610871a9228c39fa0d4df4bea64a8")
	log.Printf("Input : %v\n", keys)
	r0 := gologoo__KeySetFromKeys_497610871a9228c39fa0d4df4bea64a8(keys...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (u union) keySetProto() (*sppb.KeySet, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__keySetProto_497610871a9228c39fa0d4df4bea64a8")
	log.Printf("Input : (none)\n")
	r0, r1 := u.gologoo__keySetProto_497610871a9228c39fa0d4df4bea64a8()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
