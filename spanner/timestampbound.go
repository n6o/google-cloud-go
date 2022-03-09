package spanner

import (
	"fmt"
	"time"
	pbd "github.com/golang/protobuf/ptypes/duration"
	pbt "github.com/golang/protobuf/ptypes/timestamp"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"log"
)

type timestampBoundType int

const (
	strong timestampBoundType = iota
	exactStaleness
	maxStaleness
	minReadTimestamp
	readTimestamp
)

type TimestampBound struct {
	mode timestampBoundType
	d    time.Duration
	t    time.Time
}

func gologoo__StrongRead_6bc0c64d39a8e012c46d055c205e1d9e() TimestampBound {
	return TimestampBound{mode: strong}
}
func gologoo__ExactStaleness_6bc0c64d39a8e012c46d055c205e1d9e(d time.Duration) TimestampBound {
	return TimestampBound{mode: exactStaleness, d: d}
}
func gologoo__MaxStaleness_6bc0c64d39a8e012c46d055c205e1d9e(d time.Duration) TimestampBound {
	return TimestampBound{mode: maxStaleness, d: d}
}
func gologoo__MinReadTimestamp_6bc0c64d39a8e012c46d055c205e1d9e(t time.Time) TimestampBound {
	return TimestampBound{mode: minReadTimestamp, t: t}
}
func gologoo__ReadTimestamp_6bc0c64d39a8e012c46d055c205e1d9e(t time.Time) TimestampBound {
	return TimestampBound{mode: readTimestamp, t: t}
}
func (tb TimestampBound) gologoo__String_6bc0c64d39a8e012c46d055c205e1d9e() string {
	switch tb.mode {
	case strong:
		return fmt.Sprintf("(strong)")
	case exactStaleness:
		return fmt.Sprintf("(exactStaleness: %s)", tb.d)
	case maxStaleness:
		return fmt.Sprintf("(maxStaleness: %s)", tb.d)
	case minReadTimestamp:
		return fmt.Sprintf("(minReadTimestamp: %s)", tb.t)
	case readTimestamp:
		return fmt.Sprintf("(readTimestamp: %s)", tb.t)
	default:
		return fmt.Sprintf("{mode=%v, d=%v, t=%v}", tb.mode, tb.d, tb.t)
	}
}
func gologoo__durationProto_6bc0c64d39a8e012c46d055c205e1d9e(d time.Duration) *pbd.Duration {
	n := d.Nanoseconds()
	return &pbd.Duration{Seconds: n / int64(time.Second), Nanos: int32(n % int64(time.Second))}
}
func gologoo__timestampProto_6bc0c64d39a8e012c46d055c205e1d9e(t time.Time) *pbt.Timestamp {
	return &pbt.Timestamp{Seconds: t.Unix(), Nanos: int32(t.Nanosecond())}
}
func gologoo__buildTransactionOptionsReadOnly_6bc0c64d39a8e012c46d055c205e1d9e(tb TimestampBound, returnReadTimestamp bool) *sppb.TransactionOptions_ReadOnly {
	pb := &sppb.TransactionOptions_ReadOnly{ReturnReadTimestamp: returnReadTimestamp}
	switch tb.mode {
	case strong:
		pb.TimestampBound = &sppb.TransactionOptions_ReadOnly_Strong{Strong: true}
	case exactStaleness:
		pb.TimestampBound = &sppb.TransactionOptions_ReadOnly_ExactStaleness{ExactStaleness: durationProto(tb.d)}
	case maxStaleness:
		pb.TimestampBound = &sppb.TransactionOptions_ReadOnly_MaxStaleness{MaxStaleness: durationProto(tb.d)}
	case minReadTimestamp:
		pb.TimestampBound = &sppb.TransactionOptions_ReadOnly_MinReadTimestamp{MinReadTimestamp: timestampProto(tb.t)}
	case readTimestamp:
		pb.TimestampBound = &sppb.TransactionOptions_ReadOnly_ReadTimestamp{ReadTimestamp: timestampProto(tb.t)}
	default:
		panic(fmt.Sprintf("buildTransactionOptionsReadOnly(%v,%v)", tb, returnReadTimestamp))
	}
	return pb
}
func StrongRead() TimestampBound {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__StrongRead_6bc0c64d39a8e012c46d055c205e1d9e")
	log.Printf("Input : (none)\n")
	r0 := gologoo__StrongRead_6bc0c64d39a8e012c46d055c205e1d9e()
	log.Printf("Output: %v\n", r0)
	return r0
}
func ExactStaleness(d time.Duration) TimestampBound {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ExactStaleness_6bc0c64d39a8e012c46d055c205e1d9e")
	log.Printf("Input : %v\n", d)
	r0 := gologoo__ExactStaleness_6bc0c64d39a8e012c46d055c205e1d9e(d)
	log.Printf("Output: %v\n", r0)
	return r0
}
func MaxStaleness(d time.Duration) TimestampBound {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__MaxStaleness_6bc0c64d39a8e012c46d055c205e1d9e")
	log.Printf("Input : %v\n", d)
	r0 := gologoo__MaxStaleness_6bc0c64d39a8e012c46d055c205e1d9e(d)
	log.Printf("Output: %v\n", r0)
	return r0
}
func MinReadTimestamp(t time.Time) TimestampBound {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__MinReadTimestamp_6bc0c64d39a8e012c46d055c205e1d9e")
	log.Printf("Input : %v\n", t)
	r0 := gologoo__MinReadTimestamp_6bc0c64d39a8e012c46d055c205e1d9e(t)
	log.Printf("Output: %v\n", r0)
	return r0
}
func ReadTimestamp(t time.Time) TimestampBound {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ReadTimestamp_6bc0c64d39a8e012c46d055c205e1d9e")
	log.Printf("Input : %v\n", t)
	r0 := gologoo__ReadTimestamp_6bc0c64d39a8e012c46d055c205e1d9e(t)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (tb TimestampBound) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_6bc0c64d39a8e012c46d055c205e1d9e")
	log.Printf("Input : (none)\n")
	r0 := tb.gologoo__String_6bc0c64d39a8e012c46d055c205e1d9e()
	log.Printf("Output: %v\n", r0)
	return r0
}
func durationProto(d time.Duration) *pbd.Duration {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__durationProto_6bc0c64d39a8e012c46d055c205e1d9e")
	log.Printf("Input : %v\n", d)
	r0 := gologoo__durationProto_6bc0c64d39a8e012c46d055c205e1d9e(d)
	log.Printf("Output: %v\n", r0)
	return r0
}
func timestampProto(t time.Time) *pbt.Timestamp {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__timestampProto_6bc0c64d39a8e012c46d055c205e1d9e")
	log.Printf("Input : %v\n", t)
	r0 := gologoo__timestampProto_6bc0c64d39a8e012c46d055c205e1d9e(t)
	log.Printf("Output: %v\n", r0)
	return r0
}
func buildTransactionOptionsReadOnly(tb TimestampBound, returnReadTimestamp bool) *sppb.TransactionOptions_ReadOnly {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__buildTransactionOptionsReadOnly_6bc0c64d39a8e012c46d055c205e1d9e")
	log.Printf("Input : %v %v\n", tb, returnReadTimestamp)
	r0 := gologoo__buildTransactionOptionsReadOnly_6bc0c64d39a8e012c46d055c205e1d9e(tb, returnReadTimestamp)
	log.Printf("Output: %v\n", r0)
	return r0
}
