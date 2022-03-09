package spanner

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
	"cloud.google.com/go/civil"
	"cloud.google.com/go/internal/fields"
	"github.com/golang/protobuf/proto"
	proto3 "github.com/golang/protobuf/ptypes/struct"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc/codes"
	"log"
)

const (
	nullString                       = "<null>"
	commitTimestampPlaceholderString = "spanner.commit_timestamp()"
	NumericPrecisionDigits           = 38
	NumericScaleDigits               = 9
)

type LossOfPrecisionHandlingOption int

const (
	NumericRound LossOfPrecisionHandlingOption = iota
	NumericError
)

var LossOfPrecisionHandling LossOfPrecisionHandlingOption

func gologoo__NumericString_6a1b612599996733c2066401f14792cf(r *big.Rat) string {
	return r.FloatString(NumericScaleDigits)
}
func gologoo__validateNumeric_6a1b612599996733c2066401f14792cf(r *big.Rat) error {
	if r == nil {
		return nil
	}
	strRep := r.FloatString(NumericScaleDigits + 1)
	strRep = strings.TrimRight(strRep, "0")
	strRep = strings.TrimLeft(strRep, "-")
	s := strings.Split(strRep, ".")
	whole := s[0]
	scale := s[1]
	if len(scale) > NumericScaleDigits {
		return fmt.Errorf("max scale for a numeric is %d. The requested numeric has more", NumericScaleDigits)
	}
	if len(whole) > NumericPrecisionDigits-NumericScaleDigits {
		return fmt.Errorf("max precision for the whole component of a numeric is %d. The requested numeric has a whole component with precision %d", NumericPrecisionDigits-NumericScaleDigits, len(whole))
	}
	return nil
}

var (
	CommitTimestamp = commitTimestamp
	commitTimestamp = time.Unix(0, 0).In(time.FixedZone("CommitTimestamp placeholder", 0xDB))
	jsonNullBytes   = []byte("null")
)

type Encoder interface {
	EncodeSpanner() (interface {
	}, error)
}
type Decoder interface {
	DecodeSpanner(input interface {
	}) error
}
type NullableValue interface {
	IsNull() bool
}
type NullInt64 struct {
	Int64 int64
	Valid bool
}

func (n NullInt64) gologoo__IsNull_6a1b612599996733c2066401f14792cf() bool {
	return !n.Valid
}
func (n NullInt64) gologoo__String_6a1b612599996733c2066401f14792cf() string {
	if !n.Valid {
		return nullString
	}
	return fmt.Sprintf("%v", n.Int64)
}
func (n NullInt64) gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf() ([]byte, error) {
	if n.Valid {
		return []byte(fmt.Sprintf("%v", n.Int64)), nil
	}
	return jsonNullBytes, nil
}
func (n *NullInt64) gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("payload should not be nil")
	}
	if bytes.Equal(payload, jsonNullBytes) {
		n.Int64 = int64(0)
		n.Valid = false
		return nil
	}
	num, err := strconv.ParseInt(string(payload), 10, 64)
	if err != nil {
		return fmt.Errorf("payload cannot be converted to int64: got %v", string(payload))
	}
	n.Int64 = num
	n.Valid = true
	return nil
}
func (n NullInt64) gologoo__Value_6a1b612599996733c2066401f14792cf() (driver.Value, error) {
	if n.IsNull() {
		return nil, nil
	}
	return n.Int64, nil
}
func (n *NullInt64) gologoo__Scan_6a1b612599996733c2066401f14792cf(value interface {
}) error {
	if value == nil {
		n.Int64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	switch p := value.(type) {
	default:
		return spannerErrorf(codes.InvalidArgument, "invalid type for NullInt64: %v", p)
	case *int64:
		n.Int64 = *p
	case int64:
		n.Int64 = p
	case *NullInt64:
		n.Int64 = p.Int64
		n.Valid = p.Valid
	case NullInt64:
		n.Int64 = p.Int64
		n.Valid = p.Valid
	}
	return nil
}
func (n NullInt64) gologoo__GormDataType_6a1b612599996733c2066401f14792cf() string {
	return "INT64"
}

type NullString struct {
	StringVal string
	Valid     bool
}

func (n NullString) gologoo__IsNull_6a1b612599996733c2066401f14792cf() bool {
	return !n.Valid
}
func (n NullString) gologoo__String_6a1b612599996733c2066401f14792cf() string {
	if !n.Valid {
		return nullString
	}
	return n.StringVal
}
func (n NullString) gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf() ([]byte, error) {
	if n.Valid {
		return []byte(fmt.Sprintf("%q", n.StringVal)), nil
	}
	return jsonNullBytes, nil
}
func (n *NullString) gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("payload should not be nil")
	}
	if bytes.Equal(payload, jsonNullBytes) {
		n.StringVal = ""
		n.Valid = false
		return nil
	}
	var s *string
	if err := json.Unmarshal(payload, &s); err != nil {
		return err
	}
	if s != nil {
		n.StringVal = *s
		n.Valid = true
	} else {
		n.StringVal = ""
		n.Valid = false
	}
	return nil
}
func (n NullString) gologoo__Value_6a1b612599996733c2066401f14792cf() (driver.Value, error) {
	if n.IsNull() {
		return nil, nil
	}
	return n.StringVal, nil
}
func (n *NullString) gologoo__Scan_6a1b612599996733c2066401f14792cf(value interface {
}) error {
	if value == nil {
		n.StringVal, n.Valid = "", false
		return nil
	}
	n.Valid = true
	switch p := value.(type) {
	default:
		return spannerErrorf(codes.InvalidArgument, "invalid type for NullString: %v", p)
	case *string:
		n.StringVal = *p
	case string:
		n.StringVal = p
	case *NullString:
		n.StringVal = p.StringVal
		n.Valid = p.Valid
	case NullString:
		n.StringVal = p.StringVal
		n.Valid = p.Valid
	}
	return nil
}
func (n NullString) gologoo__GormDataType_6a1b612599996733c2066401f14792cf() string {
	return "STRING(MAX)"
}

type NullFloat64 struct {
	Float64 float64
	Valid   bool
}

func (n NullFloat64) gologoo__IsNull_6a1b612599996733c2066401f14792cf() bool {
	return !n.Valid
}
func (n NullFloat64) gologoo__String_6a1b612599996733c2066401f14792cf() string {
	if !n.Valid {
		return nullString
	}
	return fmt.Sprintf("%v", n.Float64)
}
func (n NullFloat64) gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf() ([]byte, error) {
	if n.Valid {
		return []byte(fmt.Sprintf("%v", n.Float64)), nil
	}
	return jsonNullBytes, nil
}
func (n *NullFloat64) gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("payload should not be nil")
	}
	if bytes.Equal(payload, jsonNullBytes) {
		n.Float64 = float64(0)
		n.Valid = false
		return nil
	}
	num, err := strconv.ParseFloat(string(payload), 64)
	if err != nil {
		return fmt.Errorf("payload cannot be converted to float64: got %v", string(payload))
	}
	n.Float64 = num
	n.Valid = true
	return nil
}
func (n NullFloat64) gologoo__Value_6a1b612599996733c2066401f14792cf() (driver.Value, error) {
	if n.IsNull() {
		return nil, nil
	}
	return n.Float64, nil
}
func (n *NullFloat64) gologoo__Scan_6a1b612599996733c2066401f14792cf(value interface {
}) error {
	if value == nil {
		n.Float64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	switch p := value.(type) {
	default:
		return spannerErrorf(codes.InvalidArgument, "invalid type for NullFloat64: %v", p)
	case *float64:
		n.Float64 = *p
	case float64:
		n.Float64 = p
	case *NullFloat64:
		n.Float64 = p.Float64
		n.Valid = p.Valid
	case NullFloat64:
		n.Float64 = p.Float64
		n.Valid = p.Valid
	}
	return nil
}
func (n NullFloat64) gologoo__GormDataType_6a1b612599996733c2066401f14792cf() string {
	return "FLOAT64"
}

type NullBool struct {
	Bool  bool
	Valid bool
}

func (n NullBool) gologoo__IsNull_6a1b612599996733c2066401f14792cf() bool {
	return !n.Valid
}
func (n NullBool) gologoo__String_6a1b612599996733c2066401f14792cf() string {
	if !n.Valid {
		return nullString
	}
	return fmt.Sprintf("%v", n.Bool)
}
func (n NullBool) gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf() ([]byte, error) {
	if n.Valid {
		return []byte(fmt.Sprintf("%v", n.Bool)), nil
	}
	return jsonNullBytes, nil
}
func (n *NullBool) gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("payload should not be nil")
	}
	if bytes.Equal(payload, jsonNullBytes) {
		n.Bool = false
		n.Valid = false
		return nil
	}
	b, err := strconv.ParseBool(string(payload))
	if err != nil {
		return fmt.Errorf("payload cannot be converted to bool: got %v", string(payload))
	}
	n.Bool = b
	n.Valid = true
	return nil
}
func (n NullBool) gologoo__Value_6a1b612599996733c2066401f14792cf() (driver.Value, error) {
	if n.IsNull() {
		return nil, nil
	}
	return n.Bool, nil
}
func (n *NullBool) gologoo__Scan_6a1b612599996733c2066401f14792cf(value interface {
}) error {
	if value == nil {
		n.Bool, n.Valid = false, false
		return nil
	}
	n.Valid = true
	switch p := value.(type) {
	default:
		return spannerErrorf(codes.InvalidArgument, "invalid type for NullBool: %v", p)
	case *bool:
		n.Bool = *p
	case bool:
		n.Bool = p
	case *NullBool:
		n.Bool = p.Bool
		n.Valid = p.Valid
	case NullBool:
		n.Bool = p.Bool
		n.Valid = p.Valid
	}
	return nil
}
func (n NullBool) gologoo__GormDataType_6a1b612599996733c2066401f14792cf() string {
	return "BOOL"
}

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (n NullTime) gologoo__IsNull_6a1b612599996733c2066401f14792cf() bool {
	return !n.Valid
}
func (n NullTime) gologoo__String_6a1b612599996733c2066401f14792cf() string {
	if !n.Valid {
		return nullString
	}
	return n.Time.Format(time.RFC3339Nano)
}
func (n NullTime) gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf() ([]byte, error) {
	if n.Valid {
		return []byte(fmt.Sprintf("%q", n.String())), nil
	}
	return jsonNullBytes, nil
}
func (n *NullTime) gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("payload should not be nil")
	}
	if bytes.Equal(payload, jsonNullBytes) {
		n.Time = time.Time{}
		n.Valid = false
		return nil
	}
	payload, err := trimDoubleQuotes(payload)
	if err != nil {
		return err
	}
	s := string(payload)
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return fmt.Errorf("payload cannot be converted to time.Time: got %v", string(payload))
	}
	n.Time = t
	n.Valid = true
	return nil
}
func (n NullTime) gologoo__Value_6a1b612599996733c2066401f14792cf() (driver.Value, error) {
	if n.IsNull() {
		return nil, nil
	}
	return n.Time, nil
}
func (n *NullTime) gologoo__Scan_6a1b612599996733c2066401f14792cf(value interface {
}) error {
	if value == nil {
		n.Time, n.Valid = time.Time{}, false
		return nil
	}
	n.Valid = true
	switch p := value.(type) {
	default:
		return spannerErrorf(codes.InvalidArgument, "invalid type for NullTime: %v", p)
	case *time.Time:
		n.Time = *p
	case time.Time:
		n.Time = p
	case *NullTime:
		n.Time = p.Time
		n.Valid = p.Valid
	case NullTime:
		n.Time = p.Time
		n.Valid = p.Valid
	}
	return nil
}
func (n NullTime) gologoo__GormDataType_6a1b612599996733c2066401f14792cf() string {
	return "TIMESTAMP"
}

type NullDate struct {
	Date  civil.Date
	Valid bool
}

func (n NullDate) gologoo__IsNull_6a1b612599996733c2066401f14792cf() bool {
	return !n.Valid
}
func (n NullDate) gologoo__String_6a1b612599996733c2066401f14792cf() string {
	if !n.Valid {
		return nullString
	}
	return n.Date.String()
}
func (n NullDate) gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf() ([]byte, error) {
	if n.Valid {
		return []byte(fmt.Sprintf("%q", n.String())), nil
	}
	return jsonNullBytes, nil
}
func (n *NullDate) gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("payload should not be nil")
	}
	if bytes.Equal(payload, jsonNullBytes) {
		n.Date = civil.Date{}
		n.Valid = false
		return nil
	}
	payload, err := trimDoubleQuotes(payload)
	if err != nil {
		return err
	}
	s := string(payload)
	t, err := civil.ParseDate(s)
	if err != nil {
		return fmt.Errorf("payload cannot be converted to civil.Date: got %v", string(payload))
	}
	n.Date = t
	n.Valid = true
	return nil
}
func (n NullDate) gologoo__Value_6a1b612599996733c2066401f14792cf() (driver.Value, error) {
	if n.IsNull() {
		return nil, nil
	}
	return n.Date, nil
}
func (n *NullDate) gologoo__Scan_6a1b612599996733c2066401f14792cf(value interface {
}) error {
	if value == nil {
		n.Date, n.Valid = civil.Date{}, false
		return nil
	}
	n.Valid = true
	switch p := value.(type) {
	default:
		return spannerErrorf(codes.InvalidArgument, "invalid type for NullDate: %v", p)
	case *civil.Date:
		n.Date = *p
	case civil.Date:
		n.Date = p
	case *NullDate:
		n.Date = p.Date
		n.Valid = p.Valid
	case NullDate:
		n.Date = p.Date
		n.Valid = p.Valid
	}
	return nil
}
func (n NullDate) gologoo__GormDataType_6a1b612599996733c2066401f14792cf() string {
	return "DATE"
}

type NullNumeric struct {
	Numeric big.Rat
	Valid   bool
}

func (n NullNumeric) gologoo__IsNull_6a1b612599996733c2066401f14792cf() bool {
	return !n.Valid
}
func (n NullNumeric) gologoo__String_6a1b612599996733c2066401f14792cf() string {
	if !n.Valid {
		return nullString
	}
	return fmt.Sprintf("%v", NumericString(&n.Numeric))
}
func (n NullNumeric) gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf() ([]byte, error) {
	if n.Valid {
		return []byte(fmt.Sprintf("%q", NumericString(&n.Numeric))), nil
	}
	return jsonNullBytes, nil
}
func (n *NullNumeric) gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("payload should not be nil")
	}
	if bytes.Equal(payload, jsonNullBytes) {
		n.Numeric = big.Rat{}
		n.Valid = false
		return nil
	}
	payload, err := trimDoubleQuotes(payload)
	if err != nil {
		return err
	}
	s := string(payload)
	val, ok := (&big.Rat{}).SetString(s)
	if !ok {
		return fmt.Errorf("payload cannot be converted to big.Rat: got %v", string(payload))
	}
	n.Numeric = *val
	n.Valid = true
	return nil
}
func (n NullNumeric) gologoo__Value_6a1b612599996733c2066401f14792cf() (driver.Value, error) {
	if n.IsNull() {
		return nil, nil
	}
	return n.Numeric, nil
}
func (n *NullNumeric) gologoo__Scan_6a1b612599996733c2066401f14792cf(value interface {
}) error {
	if value == nil {
		n.Numeric, n.Valid = big.Rat{}, false
		return nil
	}
	n.Valid = true
	switch p := value.(type) {
	default:
		return spannerErrorf(codes.InvalidArgument, "invalid type for NullNumeric: %v", p)
	case *big.Rat:
		n.Numeric = *p
	case big.Rat:
		n.Numeric = p
	case *NullNumeric:
		n.Numeric = p.Numeric
		n.Valid = p.Valid
	case NullNumeric:
		n.Numeric = p.Numeric
		n.Valid = p.Valid
	}
	return nil
}
func (n NullNumeric) gologoo__GormDataType_6a1b612599996733c2066401f14792cf() string {
	return "NUMERIC"
}

type NullJSON struct {
	Value interface {
	}
	Valid bool
}

func (n NullJSON) gologoo__IsNull_6a1b612599996733c2066401f14792cf() bool {
	return !n.Valid
}
func (n NullJSON) gologoo__String_6a1b612599996733c2066401f14792cf() string {
	if !n.Valid {
		return nullString
	}
	b, err := json.Marshal(n.Value)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return fmt.Sprintf("%v", string(b))
}
func (n NullJSON) gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Value)
	}
	return jsonNullBytes, nil
}
func (n *NullJSON) gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("payload should not be nil")
	}
	if bytes.Equal(payload, jsonNullBytes) {
		n.Valid = false
		return nil
	}
	var v interface {
	}
	err := json.Unmarshal(payload, &v)
	if err != nil {
		return fmt.Errorf("payload cannot be converted to a struct: got %v, err: %s", string(payload), err)
	}
	n.Value = v
	n.Valid = true
	return nil
}
func (n NullJSON) gologoo__GormDataType_6a1b612599996733c2066401f14792cf() string {
	return "JSON"
}

type PGNumeric struct {
	Numeric string
	Valid   bool
}

func (n PGNumeric) gologoo__IsNull_6a1b612599996733c2066401f14792cf() bool {
	return !n.Valid
}
func (n PGNumeric) gologoo__String_6a1b612599996733c2066401f14792cf() string {
	if !n.Valid {
		return nullString
	}
	return n.Numeric
}
func (n PGNumeric) gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf() ([]byte, error) {
	if n.Valid {
		return []byte(fmt.Sprintf("%q", n.Numeric)), nil
	}
	return jsonNullBytes, nil
}
func (n *PGNumeric) gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload []byte) error {
	if payload == nil {
		return fmt.Errorf("payload should not be nil")
	}
	if bytes.Equal(payload, jsonNullBytes) {
		n.Numeric = ""
		n.Valid = false
		return nil
	}
	payload, err := trimDoubleQuotes(payload)
	if err != nil {
		return err
	}
	n.Numeric = string(payload)
	n.Valid = true
	return nil
}

type NullRow struct {
	Row   Row
	Valid bool
}
type GenericColumnValue struct {
	Type  *sppb.Type
	Value *proto3.Value
}

func (v GenericColumnValue) gologoo__Decode_6a1b612599996733c2066401f14792cf(ptr interface {
}) error {
	return decodeValue(v.Value, v.Type, ptr)
}
func gologoo__newGenericColumnValue_6a1b612599996733c2066401f14792cf(v interface {
}) (*GenericColumnValue, error) {
	value, typ, err := encodeValue(v)
	if err != nil {
		return nil, err
	}
	return &GenericColumnValue{Value: value, Type: typ}, nil
}
func gologoo__errTypeMismatch_6a1b612599996733c2066401f14792cf(srcCode, elCode sppb.TypeCode, dst interface {
}) error {
	s := srcCode.String()
	if srcCode == sppb.TypeCode_ARRAY {
		s = fmt.Sprintf("%v[%v]", srcCode, elCode)
	}
	return spannerErrorf(codes.InvalidArgument, "type %T cannot be used for decoding %s", dst, s)
}
func gologoo__errNilSpannerType_6a1b612599996733c2066401f14792cf() error {
	return spannerErrorf(codes.FailedPrecondition, "unexpected nil Cloud Spanner data type in decoding")
}
func gologoo__errNilSrc_6a1b612599996733c2066401f14792cf() error {
	return spannerErrorf(codes.FailedPrecondition, "unexpected nil Cloud Spanner value in decoding")
}
func gologoo__errNilDst_6a1b612599996733c2066401f14792cf(dst interface {
}) error {
	return spannerErrorf(codes.InvalidArgument, "cannot decode into nil type %T", dst)
}
func gologoo__errNilArrElemType_6a1b612599996733c2066401f14792cf(t *sppb.Type) error {
	return spannerErrorf(codes.FailedPrecondition, "array type %v is with nil array element type", t)
}
func gologoo__errUnsupportedEmbeddedStructFields_6a1b612599996733c2066401f14792cf(fname string) error {
	return spannerErrorf(codes.InvalidArgument, "Embedded field: %s. Embedded and anonymous fields are not allowed "+"when converting Go structs to Cloud Spanner STRUCT values. To create a STRUCT value with an "+"unnamed field, use a `spanner:\"\"` field tag.", fname)
}
func gologoo__errDstNotForNull_6a1b612599996733c2066401f14792cf(dst interface {
}) error {
	return spannerErrorf(codes.InvalidArgument, "destination %T cannot support NULL SQL values", dst)
}
func gologoo__errBadEncoding_6a1b612599996733c2066401f14792cf(v *proto3.Value, err error) error {
	return spannerErrorf(codes.FailedPrecondition, "%v wasn't correctly encoded: <%v>", v, err)
}
func gologoo__parseNullTime_6a1b612599996733c2066401f14792cf(v *proto3.Value, p *NullTime, code sppb.TypeCode, isNull bool) error {
	if p == nil {
		return errNilDst(p)
	}
	if code != sppb.TypeCode_TIMESTAMP {
		return errTypeMismatch(code, sppb.TypeCode_TYPE_CODE_UNSPECIFIED, p)
	}
	if isNull {
		*p = NullTime{}
		return nil
	}
	x, err := getStringValue(v)
	if err != nil {
		return err
	}
	y, err := time.Parse(time.RFC3339Nano, x)
	if err != nil {
		return errBadEncoding(v, err)
	}
	p.Valid = true
	p.Time = y
	return nil
}
func gologoo__decodeValue_6a1b612599996733c2066401f14792cf(v *proto3.Value, t *sppb.Type, ptr interface {
}) error {
	if v == nil {
		return errNilSrc()
	}
	if t == nil {
		return errNilSpannerType()
	}
	code := t.Code
	typeAnnotation := t.TypeAnnotation
	acode := sppb.TypeCode_TYPE_CODE_UNSPECIFIED
	atypeAnnotation := sppb.TypeAnnotationCode_TYPE_ANNOTATION_CODE_UNSPECIFIED
	if code == sppb.TypeCode_ARRAY {
		if t.ArrayElementType == nil {
			return errNilArrElemType(t)
		}
		acode = t.ArrayElementType.Code
		atypeAnnotation = t.ArrayElementType.TypeAnnotation
	}
	_, isNull := v.Kind.(*proto3.Value_NullValue)
	switch p := ptr.(type) {
	case nil:
		return errNilDst(nil)
	case *string:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_STRING {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			return errDstNotForNull(ptr)
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		*p = x
	case *NullString, **string, *sql.NullString:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_STRING {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *NullString:
				*sp = NullString{}
			case **string:
				*sp = nil
			case *sql.NullString:
				*sp = sql.NullString{}
			}
			break
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *NullString:
			sp.Valid = true
			sp.StringVal = x
		case **string:
			*sp = &x
		case *sql.NullString:
			sp.Valid = true
			sp.String = x
		}
	case *[]NullString, *[]*string:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_STRING {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *[]NullString:
				*sp = nil
			case *[]*string:
				*sp = nil
			}
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *[]NullString:
			y, err := decodeNullStringArray(x)
			if err != nil {
				return err
			}
			*sp = y
		case *[]*string:
			y, err := decodeStringPointerArray(x)
			if err != nil {
				return err
			}
			*sp = y
		}
	case *[]string:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_STRING {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeStringArray(x)
		if err != nil {
			return err
		}
		*p = y
	case *[]byte:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_BYTES {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		y, err := base64.StdEncoding.DecodeString(x)
		if err != nil {
			return errBadEncoding(v, err)
		}
		*p = y
	case *[][]byte:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_BYTES {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeByteArray(x)
		if err != nil {
			return err
		}
		*p = y
	case *int64:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_INT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			return errDstNotForNull(ptr)
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		y, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return errBadEncoding(v, err)
		}
		*p = y
	case *NullInt64, **int64:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_INT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *NullInt64:
				*sp = NullInt64{}
			case **int64:
				*sp = nil
			}
			break
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		y, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return errBadEncoding(v, err)
		}
		switch sp := ptr.(type) {
		case *NullInt64:
			sp.Valid = true
			sp.Int64 = y
		case **int64:
			*sp = &y
		}
	case *[]NullInt64, *[]*int64:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_INT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *[]NullInt64:
				*sp = nil
			case *[]*int64:
				*sp = nil
			}
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *[]NullInt64:
			y, err := decodeNullInt64Array(x)
			if err != nil {
				return err
			}
			*sp = y
		case *[]*int64:
			y, err := decodeInt64PointerArray(x)
			if err != nil {
				return err
			}
			*sp = y
		}
	case *[]int64:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_INT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeInt64Array(x)
		if err != nil {
			return err
		}
		*p = y
	case *bool:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_BOOL {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			return errDstNotForNull(ptr)
		}
		x, err := getBoolValue(v)
		if err != nil {
			return err
		}
		*p = x
	case *NullBool, **bool:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_BOOL {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *NullBool:
				*sp = NullBool{}
			case **bool:
				*sp = nil
			}
			break
		}
		x, err := getBoolValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *NullBool:
			sp.Valid = true
			sp.Bool = x
		case **bool:
			*sp = &x
		}
	case *[]NullBool, *[]*bool:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_BOOL {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *[]NullBool:
				*sp = nil
			case *[]*bool:
				*sp = nil
			}
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *[]NullBool:
			y, err := decodeNullBoolArray(x)
			if err != nil {
				return err
			}
			*sp = y
		case *[]*bool:
			y, err := decodeBoolPointerArray(x)
			if err != nil {
				return err
			}
			*sp = y
		}
	case *[]bool:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_BOOL {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeBoolArray(x)
		if err != nil {
			return err
		}
		*p = y
	case *float64:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_FLOAT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			return errDstNotForNull(ptr)
		}
		x, err := getFloat64Value(v)
		if err != nil {
			return err
		}
		*p = x
	case *NullFloat64, **float64:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_FLOAT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *NullFloat64:
				*sp = NullFloat64{}
			case **float64:
				*sp = nil
			}
			break
		}
		x, err := getFloat64Value(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *NullFloat64:
			sp.Valid = true
			sp.Float64 = x
		case **float64:
			*sp = &x
		}
	case *[]NullFloat64, *[]*float64:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_FLOAT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *[]NullFloat64:
				*sp = nil
			case *[]*float64:
				*sp = nil
			}
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *[]NullFloat64:
			y, err := decodeNullFloat64Array(x)
			if err != nil {
				return err
			}
			*sp = y
		case *[]*float64:
			y, err := decodeFloat64PointerArray(x)
			if err != nil {
				return err
			}
			*sp = y
		}
	case *[]float64:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_FLOAT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeFloat64Array(x)
		if err != nil {
			return err
		}
		*p = y
	case *big.Rat:
		if code != sppb.TypeCode_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			return errDstNotForNull(ptr)
		}
		x := v.GetStringValue()
		y, ok := (&big.Rat{}).SetString(x)
		if !ok {
			return errUnexpectedNumericStr(x)
		}
		*p = *y
	case *NullJSON:
		if p == nil {
			return errNilDst(p)
		}
		if code == sppb.TypeCode_ARRAY {
			if acode != sppb.TypeCode_JSON {
				return errTypeMismatch(code, acode, ptr)
			}
			x, err := getListValue(v)
			if err != nil {
				return err
			}
			y, err := decodeNullJSONArrayToNullJSON(x)
			if err != nil {
				return err
			}
			*p = *y
		} else {
			if code != sppb.TypeCode_JSON {
				return errTypeMismatch(code, acode, ptr)
			}
			if isNull {
				*p = NullJSON{}
				break
			}
			x := v.GetStringValue()
			var y interface {
			}
			err := json.Unmarshal([]byte(x), &y)
			if err != nil {
				return err
			}
			*p = NullJSON{y, true}
		}
	case *[]NullJSON:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_JSON {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeNullJSONArray(x)
		if err != nil {
			return err
		}
		*p = y
	case *NullNumeric:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = NullNumeric{}
			break
		}
		x := v.GetStringValue()
		y, ok := (&big.Rat{}).SetString(x)
		if !ok {
			return errUnexpectedNumericStr(x)
		}
		*p = NullNumeric{*y, true}
	case **big.Rat:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x := v.GetStringValue()
		y, ok := (&big.Rat{}).SetString(x)
		if !ok {
			return errUnexpectedNumericStr(x)
		}
		*p = y
	case *[]NullNumeric, *[]*big.Rat:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *[]NullNumeric:
				*sp = nil
			case *[]*big.Rat:
				*sp = nil
			}
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *[]NullNumeric:
			y, err := decodeNullNumericArray(x)
			if err != nil {
				return err
			}
			*sp = y
		case *[]*big.Rat:
			y, err := decodeNumericPointerArray(x)
			if err != nil {
				return err
			}
			*sp = y
		}
	case *[]big.Rat:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeNumericArray(x)
		if err != nil {
			return err
		}
		*p = y
	case *PGNumeric:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_NUMERIC || typeAnnotation != sppb.TypeAnnotationCode_PG_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = PGNumeric{}
			break
		}
		*p = PGNumeric{v.GetStringValue(), true}
	case *[]PGNumeric:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_NUMERIC || atypeAnnotation != sppb.TypeAnnotationCode_PG_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodePGNumericArray(x)
		if err != nil {
			return err
		}
		*p = y
	case *time.Time:
		var nt NullTime
		if isNull {
			return errDstNotForNull(ptr)
		}
		err := parseNullTime(v, &nt, code, isNull)
		if err != nil {
			return err
		}
		*p = nt.Time
	case *NullTime:
		err := parseNullTime(v, p, code, isNull)
		if err != nil {
			return err
		}
	case **time.Time:
		var nt NullTime
		if isNull {
			*p = nil
			break
		}
		err := parseNullTime(v, &nt, code, isNull)
		if err != nil {
			return err
		}
		*p = &nt.Time
	case *[]NullTime, *[]*time.Time:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_TIMESTAMP {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *[]NullTime:
				*sp = nil
			case *[]*time.Time:
				*sp = nil
			}
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *[]NullTime:
			y, err := decodeNullTimeArray(x)
			if err != nil {
				return err
			}
			*sp = y
		case *[]*time.Time:
			y, err := decodeTimePointerArray(x)
			if err != nil {
				return err
			}
			*sp = y
		}
	case *[]time.Time:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_TIMESTAMP {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeTimeArray(x)
		if err != nil {
			return err
		}
		*p = y
	case *civil.Date:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_DATE {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			return errDstNotForNull(ptr)
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		y, err := civil.ParseDate(x)
		if err != nil {
			return errBadEncoding(v, err)
		}
		*p = y
	case *NullDate, **civil.Date:
		if p == nil {
			return errNilDst(p)
		}
		if code != sppb.TypeCode_DATE {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *NullDate:
				*sp = NullDate{}
			case **civil.Date:
				*sp = nil
			}
			break
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		y, err := civil.ParseDate(x)
		if err != nil {
			return errBadEncoding(v, err)
		}
		switch sp := ptr.(type) {
		case *NullDate:
			sp.Valid = true
			sp.Date = y
		case **civil.Date:
			*sp = &y
		}
	case *[]NullDate, *[]*civil.Date:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_DATE {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			switch sp := ptr.(type) {
			case *[]NullDate:
				*sp = nil
			case *[]*civil.Date:
				*sp = nil
			}
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		switch sp := ptr.(type) {
		case *[]NullDate:
			y, err := decodeNullDateArray(x)
			if err != nil {
				return err
			}
			*sp = y
		case *[]*civil.Date:
			y, err := decodeDatePointerArray(x)
			if err != nil {
				return err
			}
			*sp = y
		}
	case *[]civil.Date:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_DATE {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeDateArray(x)
		if err != nil {
			return err
		}
		*p = y
	case *[]NullRow:
		if p == nil {
			return errNilDst(p)
		}
		if acode != sppb.TypeCode_STRUCT {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			*p = nil
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeRowArray(t.ArrayElementType.StructType, x)
		if err != nil {
			return err
		}
		*p = y
	case *GenericColumnValue:
		*p = GenericColumnValue{Type: t, Value: v}
	default:
		if decodedVal, ok := ptr.(Decoder); ok {
			x, err := getGenericValue(t, v)
			if err != nil {
				return err
			}
			return decodedVal.DecodeSpanner(x)
		}
		decodableType := getDecodableSpannerType(ptr, true)
		if decodableType != spannerTypeUnknown {
			if isNull && !decodableType.supportsNull() {
				return errDstNotForNull(ptr)
			}
			return decodableType.decodeValueToCustomType(v, t, acode, atypeAnnotation, ptr)
		}
		if !(code == sppb.TypeCode_ARRAY && acode == sppb.TypeCode_STRUCT) {
			return errTypeMismatch(code, acode, ptr)
		}
		vp := reflect.ValueOf(p)
		if !vp.IsValid() {
			return errNilDst(p)
		}
		if !isPtrStructPtrSlice(vp.Type()) {
			return fmt.Errorf("the container is not a slice of struct pointers: %v", errTypeMismatch(code, acode, ptr))
		}
		if vp.IsNil() {
			return errNilDst(p)
		}
		if isNull {
			vp.Elem().Set(reflect.Zero(vp.Elem().Type()))
			break
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		if err = decodeStructArray(t.ArrayElementType.StructType, x, p); err != nil {
			return err
		}
	}
	return nil
}

type decodableSpannerType uint

const (
	spannerTypeUnknown decodableSpannerType = iota
	spannerTypeInvalid
	spannerTypeNonNullString
	spannerTypeByteArray
	spannerTypeNonNullInt64
	spannerTypeNonNullBool
	spannerTypeNonNullFloat64
	spannerTypeNonNullNumeric
	spannerTypeNonNullTime
	spannerTypeNonNullDate
	spannerTypeNullString
	spannerTypeNullInt64
	spannerTypeNullBool
	spannerTypeNullFloat64
	spannerTypeNullTime
	spannerTypeNullDate
	spannerTypeNullNumeric
	spannerTypeNullJSON
	spannerTypePGNumeric
	spannerTypeArrayOfNonNullString
	spannerTypeArrayOfByteArray
	spannerTypeArrayOfNonNullInt64
	spannerTypeArrayOfNonNullBool
	spannerTypeArrayOfNonNullFloat64
	spannerTypeArrayOfNonNullNumeric
	spannerTypeArrayOfNonNullTime
	spannerTypeArrayOfNonNullDate
	spannerTypeArrayOfNullString
	spannerTypeArrayOfNullInt64
	spannerTypeArrayOfNullBool
	spannerTypeArrayOfNullFloat64
	spannerTypeArrayOfNullNumeric
	spannerTypeArrayOfNullJSON
	spannerTypeArrayOfNullTime
	spannerTypeArrayOfNullDate
	spannerTypeArrayOfPGNumeric
)

func (d decodableSpannerType) gologoo__supportsNull_6a1b612599996733c2066401f14792cf() bool {
	switch d {
	case spannerTypeNonNullString, spannerTypeNonNullInt64, spannerTypeNonNullBool, spannerTypeNonNullFloat64, spannerTypeNonNullTime, spannerTypeNonNullDate, spannerTypeNonNullNumeric:
		return false
	default:
		return true
	}
}

var typeOfNonNullTime = reflect.TypeOf(time.Time{})
var typeOfNonNullDate = reflect.TypeOf(civil.Date{})
var typeOfNonNullNumeric = reflect.TypeOf(big.Rat{})
var typeOfNullString = reflect.TypeOf(NullString{})
var typeOfNullInt64 = reflect.TypeOf(NullInt64{})
var typeOfNullBool = reflect.TypeOf(NullBool{})
var typeOfNullFloat64 = reflect.TypeOf(NullFloat64{})
var typeOfNullTime = reflect.TypeOf(NullTime{})
var typeOfNullDate = reflect.TypeOf(NullDate{})
var typeOfNullNumeric = reflect.TypeOf(NullNumeric{})
var typeOfNullJSON = reflect.TypeOf(NullJSON{})
var typeOfPGNumeric = reflect.TypeOf(PGNumeric{})

func gologoo__getDecodableSpannerType_6a1b612599996733c2066401f14792cf(ptr interface {
}, isPtr bool) decodableSpannerType {
	var val reflect.Value
	var kind reflect.Kind
	if isPtr {
		val = reflect.Indirect(reflect.ValueOf(ptr))
	} else {
		val = reflect.ValueOf(ptr)
	}
	kind = val.Kind()
	if kind == reflect.Invalid {
		return spannerTypeInvalid
	}
	switch kind {
	case reflect.Invalid:
		return spannerTypeInvalid
	case reflect.String:
		return spannerTypeNonNullString
	case reflect.Int64:
		return spannerTypeNonNullInt64
	case reflect.Bool:
		return spannerTypeNonNullBool
	case reflect.Float64:
		return spannerTypeNonNullFloat64
	case reflect.Ptr:
		t := val.Type()
		if t.ConvertibleTo(typeOfNullNumeric) {
			return spannerTypeNullNumeric
		}
		if t.ConvertibleTo(typeOfNullJSON) {
			return spannerTypeNullJSON
		}
	case reflect.Struct:
		t := val.Type()
		if t.ConvertibleTo(typeOfNonNullNumeric) {
			return spannerTypeNonNullNumeric
		}
		if t.ConvertibleTo(typeOfNonNullTime) {
			return spannerTypeNonNullTime
		}
		if t.ConvertibleTo(typeOfNonNullDate) {
			return spannerTypeNonNullDate
		}
		if t.ConvertibleTo(typeOfNullString) {
			return spannerTypeNullString
		}
		if t.ConvertibleTo(typeOfNullInt64) {
			return spannerTypeNullInt64
		}
		if t.ConvertibleTo(typeOfNullBool) {
			return spannerTypeNullBool
		}
		if t.ConvertibleTo(typeOfNullFloat64) {
			return spannerTypeNullFloat64
		}
		if t.ConvertibleTo(typeOfNullTime) {
			return spannerTypeNullTime
		}
		if t.ConvertibleTo(typeOfNullDate) {
			return spannerTypeNullDate
		}
		if t.ConvertibleTo(typeOfNullNumeric) {
			return spannerTypeNullNumeric
		}
		if t.ConvertibleTo(typeOfNullJSON) {
			return spannerTypeNullJSON
		}
		if t.ConvertibleTo(typeOfPGNumeric) {
			return spannerTypePGNumeric
		}
	case reflect.Slice:
		kind := val.Type().Elem().Kind()
		switch kind {
		case reflect.Invalid:
			return spannerTypeUnknown
		case reflect.String:
			return spannerTypeArrayOfNonNullString
		case reflect.Uint8:
			return spannerTypeByteArray
		case reflect.Int64:
			return spannerTypeArrayOfNonNullInt64
		case reflect.Bool:
			return spannerTypeArrayOfNonNullBool
		case reflect.Float64:
			return spannerTypeArrayOfNonNullFloat64
		case reflect.Ptr:
			t := val.Type().Elem()
			if t.ConvertibleTo(typeOfNullNumeric) {
				return spannerTypeArrayOfNullNumeric
			}
		case reflect.Struct:
			t := val.Type().Elem()
			if t.ConvertibleTo(typeOfNonNullNumeric) {
				return spannerTypeArrayOfNonNullNumeric
			}
			if t.ConvertibleTo(typeOfNonNullTime) {
				return spannerTypeArrayOfNonNullTime
			}
			if t.ConvertibleTo(typeOfNonNullDate) {
				return spannerTypeArrayOfNonNullDate
			}
			if t.ConvertibleTo(typeOfNullString) {
				return spannerTypeArrayOfNullString
			}
			if t.ConvertibleTo(typeOfNullInt64) {
				return spannerTypeArrayOfNullInt64
			}
			if t.ConvertibleTo(typeOfNullBool) {
				return spannerTypeArrayOfNullBool
			}
			if t.ConvertibleTo(typeOfNullFloat64) {
				return spannerTypeArrayOfNullFloat64
			}
			if t.ConvertibleTo(typeOfNullTime) {
				return spannerTypeArrayOfNullTime
			}
			if t.ConvertibleTo(typeOfNullDate) {
				return spannerTypeArrayOfNullDate
			}
			if t.ConvertibleTo(typeOfNullNumeric) {
				return spannerTypeArrayOfNullNumeric
			}
			if t.ConvertibleTo(typeOfNullJSON) {
				return spannerTypeArrayOfNullJSON
			}
			if t.ConvertibleTo(typeOfPGNumeric) {
				return spannerTypeArrayOfPGNumeric
			}
		case reflect.Slice:
			kind := val.Type().Elem().Elem().Kind()
			switch kind {
			case reflect.Uint8:
				return spannerTypeArrayOfByteArray
			}
		}
	}
	return spannerTypeUnknown
}
func (dsc decodableSpannerType) gologoo__decodeValueToCustomType_6a1b612599996733c2066401f14792cf(v *proto3.Value, t *sppb.Type, acode sppb.TypeCode, atypeAnnotation sppb.TypeAnnotationCode, ptr interface {
}) error {
	code := t.Code
	typeAnnotation := t.TypeAnnotation
	_, isNull := v.Kind.(*proto3.Value_NullValue)
	if dsc == spannerTypeInvalid {
		return errNilDst(ptr)
	}
	if isNull && !dsc.supportsNull() {
		return errDstNotForNull(ptr)
	}
	var result interface {
	}
	switch dsc {
	case spannerTypeNonNullString, spannerTypeNullString:
		if code != sppb.TypeCode_STRING {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &NullString{}
			break
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		if dsc == spannerTypeNonNullString {
			result = &x
		} else {
			result = &NullString{x, !isNull}
		}
	case spannerTypeByteArray:
		if code != sppb.TypeCode_BYTES {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = []byte(nil)
			break
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		y, err := base64.StdEncoding.DecodeString(x)
		if err != nil {
			return errBadEncoding(v, err)
		}
		result = y
	case spannerTypeNonNullInt64, spannerTypeNullInt64:
		if code != sppb.TypeCode_INT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &NullInt64{}
			break
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		y, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return errBadEncoding(v, err)
		}
		if dsc == spannerTypeNonNullInt64 {
			result = &y
		} else {
			result = &NullInt64{y, !isNull}
		}
	case spannerTypeNonNullBool, spannerTypeNullBool:
		if code != sppb.TypeCode_BOOL {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &NullBool{}
			break
		}
		x, err := getBoolValue(v)
		if err != nil {
			return err
		}
		if dsc == spannerTypeNonNullBool {
			result = &x
		} else {
			result = &NullBool{x, !isNull}
		}
	case spannerTypeNonNullFloat64, spannerTypeNullFloat64:
		if code != sppb.TypeCode_FLOAT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &NullFloat64{}
			break
		}
		x, err := getFloat64Value(v)
		if err != nil {
			return err
		}
		if dsc == spannerTypeNonNullFloat64 {
			result = &x
		} else {
			result = &NullFloat64{x, !isNull}
		}
	case spannerTypeNonNullNumeric, spannerTypeNullNumeric:
		if code != sppb.TypeCode_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &NullNumeric{}
			break
		}
		x := v.GetStringValue()
		y, ok := (&big.Rat{}).SetString(x)
		if !ok {
			return errUnexpectedNumericStr(x)
		}
		if dsc == spannerTypeNonNullNumeric {
			result = y
		} else {
			result = &NullNumeric{*y, true}
		}
	case spannerTypePGNumeric:
		if code != sppb.TypeCode_NUMERIC || typeAnnotation != sppb.TypeAnnotationCode_PG_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &PGNumeric{}
			break
		}
		result = &PGNumeric{v.GetStringValue(), true}
	case spannerTypeNullJSON:
		if code != sppb.TypeCode_JSON {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &NullJSON{}
			break
		}
		x := v.GetStringValue()
		var y interface {
		}
		err := json.Unmarshal([]byte(x), &y)
		if err != nil {
			return err
		}
		result = &NullJSON{y, true}
	case spannerTypeNonNullTime, spannerTypeNullTime:
		var nt NullTime
		err := parseNullTime(v, &nt, code, isNull)
		if err != nil {
			return err
		}
		if dsc == spannerTypeNonNullTime {
			result = &nt.Time
		} else {
			result = &nt
		}
	case spannerTypeNonNullDate, spannerTypeNullDate:
		if code != sppb.TypeCode_DATE {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			result = &NullDate{}
			break
		}
		x, err := getStringValue(v)
		if err != nil {
			return err
		}
		y, err := civil.ParseDate(x)
		if err != nil {
			return errBadEncoding(v, err)
		}
		if dsc == spannerTypeNonNullDate {
			result = &y
		} else {
			result = &NullDate{y, !isNull}
		}
	case spannerTypeArrayOfNonNullString, spannerTypeArrayOfNullString:
		if acode != sppb.TypeCode_STRING {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, stringType(), "STRING")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfByteArray:
		if acode != sppb.TypeCode_BYTES {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, bytesType(), "BYTES")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfNonNullInt64, spannerTypeArrayOfNullInt64:
		if acode != sppb.TypeCode_INT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, intType(), "INT64")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfNonNullBool, spannerTypeArrayOfNullBool:
		if acode != sppb.TypeCode_BOOL {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, boolType(), "BOOL")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfNonNullFloat64, spannerTypeArrayOfNullFloat64:
		if acode != sppb.TypeCode_FLOAT64 {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, floatType(), "FLOAT64")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfNonNullNumeric, spannerTypeArrayOfNullNumeric:
		if acode != sppb.TypeCode_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, numericType(), "NUMERIC")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfPGNumeric:
		if acode != sppb.TypeCode_NUMERIC || atypeAnnotation != sppb.TypeAnnotationCode_PG_NUMERIC {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, pgNumericType(), "PGNUMERIC")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfNullJSON:
		if acode != sppb.TypeCode_JSON {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, jsonType(), "JSON")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfNonNullTime, spannerTypeArrayOfNullTime:
		if acode != sppb.TypeCode_TIMESTAMP {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, timeType(), "TIMESTAMP")
		if err != nil {
			return err
		}
		result = y
	case spannerTypeArrayOfNonNullDate, spannerTypeArrayOfNullDate:
		if acode != sppb.TypeCode_DATE {
			return errTypeMismatch(code, acode, ptr)
		}
		if isNull {
			ptr = nil
			return nil
		}
		x, err := getListValue(v)
		if err != nil {
			return err
		}
		y, err := decodeGenericArray(reflect.TypeOf(ptr).Elem(), x, dateType(), "DATE")
		if err != nil {
			return err
		}
		result = y
	default:
		return fmt.Errorf("unknown decodable type found: %v", dsc)
	}
	source := reflect.Indirect(reflect.ValueOf(result))
	destination := reflect.Indirect(reflect.ValueOf(ptr))
	destination.Set(source.Convert(destination.Type()))
	return nil
}
func gologoo__errSrcVal_6a1b612599996733c2066401f14792cf(v *proto3.Value, want string) error {
	return spannerErrorf(codes.FailedPrecondition, "cannot use %v(Kind: %T) as %s Value", v, v.GetKind(), want)
}
func gologoo__getStringValue_6a1b612599996733c2066401f14792cf(v *proto3.Value) (string, error) {
	if x, ok := v.GetKind().(*proto3.Value_StringValue); ok && x != nil {
		return x.StringValue, nil
	}
	return "", errSrcVal(v, "String")
}
func gologoo__getBoolValue_6a1b612599996733c2066401f14792cf(v *proto3.Value) (bool, error) {
	if x, ok := v.GetKind().(*proto3.Value_BoolValue); ok && x != nil {
		return x.BoolValue, nil
	}
	return false, errSrcVal(v, "Bool")
}
func gologoo__getListValue_6a1b612599996733c2066401f14792cf(v *proto3.Value) (*proto3.ListValue, error) {
	if x, ok := v.GetKind().(*proto3.Value_ListValue); ok && x != nil {
		return x.ListValue, nil
	}
	return nil, errSrcVal(v, "List")
}
func gologoo__getGenericValue_6a1b612599996733c2066401f14792cf(t *sppb.Type, v *proto3.Value) (interface {
}, error) {
	switch x := v.GetKind().(type) {
	case *proto3.Value_NumberValue:
		return x.NumberValue, nil
	case *proto3.Value_BoolValue:
		return x.BoolValue, nil
	case *proto3.Value_StringValue:
		return x.StringValue, nil
	case *proto3.Value_NullValue:
		return getTypedNil(t)
	default:
		return 0, errSrcVal(v, "Number, Bool, String")
	}
}
func gologoo__getTypedNil_6a1b612599996733c2066401f14792cf(t *sppb.Type) (interface {
}, error) {
	switch t.Code {
	case sppb.TypeCode_FLOAT64:
		var f *float64
		return f, nil
	case sppb.TypeCode_BOOL:
		var b *bool
		return b, nil
	default:
		var s *string
		return s, nil
	}
}
func gologoo__errUnexpectedNumericStr_6a1b612599996733c2066401f14792cf(s string) error {
	return spannerErrorf(codes.FailedPrecondition, "unexpected string value %q for numeric number", s)
}
func gologoo__errUnexpectedFloat64Str_6a1b612599996733c2066401f14792cf(s string) error {
	return spannerErrorf(codes.FailedPrecondition, "unexpected string value %q for float64 number", s)
}
func gologoo__getFloat64Value_6a1b612599996733c2066401f14792cf(v *proto3.Value) (float64, error) {
	switch x := v.GetKind().(type) {
	case *proto3.Value_NumberValue:
		if x == nil {
			break
		}
		return x.NumberValue, nil
	case *proto3.Value_StringValue:
		if x == nil {
			break
		}
		switch x.StringValue {
		case "NaN":
			return math.NaN(), nil
		case "Infinity":
			return math.Inf(1), nil
		case "-Infinity":
			return math.Inf(-1), nil
		default:
			return 0, errUnexpectedFloat64Str(x.StringValue)
		}
	}
	return 0, errSrcVal(v, "Number")
}
func gologoo__errNilListValue_6a1b612599996733c2066401f14792cf(sqlType string) error {
	return spannerErrorf(codes.FailedPrecondition, "unexpected nil ListValue in decoding %v array", sqlType)
}
func gologoo__errDecodeArrayElement_6a1b612599996733c2066401f14792cf(i int, v proto.Message, sqlType string, err error) error {
	var se *Error
	if !errorAs(err, &se) {
		return spannerErrorf(codes.Unknown, "cannot decode %v(array element %v) as %v, error = <%v>", v, i, sqlType, err)
	}
	se.decorate(fmt.Sprintf("cannot decode %v(array element %v) as %v", v, i, sqlType))
	return se
}
func gologoo__decodeGenericArray_6a1b612599996733c2066401f14792cf(tp reflect.Type, pb *proto3.ListValue, t *sppb.Type, sqlType string) (interface {
}, error) {
	if pb == nil {
		return nil, errNilListValue(sqlType)
	}
	a := reflect.MakeSlice(tp, len(pb.Values), len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, t, a.Index(i).Addr().Interface()); err != nil {
			return nil, errDecodeArrayElement(i, v, "STRING", err)
		}
	}
	return a.Interface(), nil
}
func gologoo__decodeNullStringArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]NullString, error) {
	if pb == nil {
		return nil, errNilListValue("STRING")
	}
	a := make([]NullString, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, stringType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "STRING", err)
		}
	}
	return a, nil
}
func gologoo__decodeStringPointerArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]*string, error) {
	if pb == nil {
		return nil, errNilListValue("STRING")
	}
	a := make([]*string, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, stringType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "STRING", err)
		}
	}
	return a, nil
}
func gologoo__decodeStringArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]string, error) {
	if pb == nil {
		return nil, errNilListValue("STRING")
	}
	a := make([]string, len(pb.Values))
	st := stringType()
	for i, v := range pb.Values {
		if err := decodeValue(v, st, &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "STRING", err)
		}
	}
	return a, nil
}
func gologoo__decodeNullInt64Array_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]NullInt64, error) {
	if pb == nil {
		return nil, errNilListValue("INT64")
	}
	a := make([]NullInt64, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, intType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "INT64", err)
		}
	}
	return a, nil
}
func gologoo__decodeInt64PointerArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]*int64, error) {
	if pb == nil {
		return nil, errNilListValue("INT64")
	}
	a := make([]*int64, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, intType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "INT64", err)
		}
	}
	return a, nil
}
func gologoo__decodeInt64Array_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]int64, error) {
	if pb == nil {
		return nil, errNilListValue("INT64")
	}
	a := make([]int64, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, intType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "INT64", err)
		}
	}
	return a, nil
}
func gologoo__decodeNullBoolArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]NullBool, error) {
	if pb == nil {
		return nil, errNilListValue("BOOL")
	}
	a := make([]NullBool, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, boolType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "BOOL", err)
		}
	}
	return a, nil
}
func gologoo__decodeBoolPointerArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]*bool, error) {
	if pb == nil {
		return nil, errNilListValue("BOOL")
	}
	a := make([]*bool, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, boolType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "BOOL", err)
		}
	}
	return a, nil
}
func gologoo__decodeBoolArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]bool, error) {
	if pb == nil {
		return nil, errNilListValue("BOOL")
	}
	a := make([]bool, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, boolType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "BOOL", err)
		}
	}
	return a, nil
}
func gologoo__decodeNullFloat64Array_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]NullFloat64, error) {
	if pb == nil {
		return nil, errNilListValue("FLOAT64")
	}
	a := make([]NullFloat64, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, floatType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "FLOAT64", err)
		}
	}
	return a, nil
}
func gologoo__decodeFloat64PointerArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]*float64, error) {
	if pb == nil {
		return nil, errNilListValue("FLOAT64")
	}
	a := make([]*float64, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, floatType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "FLOAT64", err)
		}
	}
	return a, nil
}
func gologoo__decodeFloat64Array_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]float64, error) {
	if pb == nil {
		return nil, errNilListValue("FLOAT64")
	}
	a := make([]float64, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, floatType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "FLOAT64", err)
		}
	}
	return a, nil
}
func gologoo__decodeNullNumericArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]NullNumeric, error) {
	if pb == nil {
		return nil, errNilListValue("NUMERIC")
	}
	a := make([]NullNumeric, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, numericType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "NUMERIC", err)
		}
	}
	return a, nil
}
func gologoo__decodeNullJSONArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]NullJSON, error) {
	if pb == nil {
		return nil, errNilListValue("JSON")
	}
	a := make([]NullJSON, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, jsonType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "JSON", err)
		}
	}
	return a, nil
}
func gologoo__decodeNullJSONArrayToNullJSON_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) (*NullJSON, error) {
	if pb == nil {
		return nil, errNilListValue("JSON")
	}
	strs := []string{}
	for _, v := range pb.Values {
		if _, ok := v.Kind.(*proto3.Value_NullValue); ok {
			strs = append(strs, "null")
		} else {
			strs = append(strs, v.GetStringValue())
		}
	}
	s := fmt.Sprintf("[%s]", strings.Join(strs, ","))
	var y interface {
	}
	err := json.Unmarshal([]byte(s), &y)
	if err != nil {
		return nil, err
	}
	return &NullJSON{y, true}, nil
}
func gologoo__decodeNumericPointerArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]*big.Rat, error) {
	if pb == nil {
		return nil, errNilListValue("NUMERIC")
	}
	a := make([]*big.Rat, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, numericType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "NUMERIC", err)
		}
	}
	return a, nil
}
func gologoo__decodeNumericArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]big.Rat, error) {
	if pb == nil {
		return nil, errNilListValue("NUMERIC")
	}
	a := make([]big.Rat, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, numericType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "NUMERIC", err)
		}
	}
	return a, nil
}
func gologoo__decodePGNumericArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]PGNumeric, error) {
	if pb == nil {
		return nil, errNilListValue("PGNUMERIC")
	}
	a := make([]PGNumeric, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, pgNumericType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "PGNUMERIC", err)
		}
	}
	return a, nil
}
func gologoo__decodeByteArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([][]byte, error) {
	if pb == nil {
		return nil, errNilListValue("BYTES")
	}
	a := make([][]byte, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, bytesType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "BYTES", err)
		}
	}
	return a, nil
}
func gologoo__decodeNullTimeArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]NullTime, error) {
	if pb == nil {
		return nil, errNilListValue("TIMESTAMP")
	}
	a := make([]NullTime, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, timeType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "TIMESTAMP", err)
		}
	}
	return a, nil
}
func gologoo__decodeTimePointerArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]*time.Time, error) {
	if pb == nil {
		return nil, errNilListValue("TIMESTAMP")
	}
	a := make([]*time.Time, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, timeType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "TIMESTAMP", err)
		}
	}
	return a, nil
}
func gologoo__decodeTimeArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]time.Time, error) {
	if pb == nil {
		return nil, errNilListValue("TIMESTAMP")
	}
	a := make([]time.Time, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, timeType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "TIMESTAMP", err)
		}
	}
	return a, nil
}
func gologoo__decodeNullDateArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]NullDate, error) {
	if pb == nil {
		return nil, errNilListValue("DATE")
	}
	a := make([]NullDate, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, dateType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "DATE", err)
		}
	}
	return a, nil
}
func gologoo__decodeDatePointerArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]*civil.Date, error) {
	if pb == nil {
		return nil, errNilListValue("DATE")
	}
	a := make([]*civil.Date, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, dateType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "DATE", err)
		}
	}
	return a, nil
}
func gologoo__decodeDateArray_6a1b612599996733c2066401f14792cf(pb *proto3.ListValue) ([]civil.Date, error) {
	if pb == nil {
		return nil, errNilListValue("DATE")
	}
	a := make([]civil.Date, len(pb.Values))
	for i, v := range pb.Values {
		if err := decodeValue(v, dateType(), &a[i]); err != nil {
			return nil, errDecodeArrayElement(i, v, "DATE", err)
		}
	}
	return a, nil
}
func gologoo__errNotStructElement_6a1b612599996733c2066401f14792cf(i int, v *proto3.Value) error {
	return errDecodeArrayElement(i, v, "STRUCT", spannerErrorf(codes.FailedPrecondition, "%v(type: %T) doesn't encode Cloud Spanner STRUCT", v, v))
}
func gologoo__decodeRowArray_6a1b612599996733c2066401f14792cf(ty *sppb.StructType, pb *proto3.ListValue) ([]NullRow, error) {
	if pb == nil {
		return nil, errNilListValue("STRUCT")
	}
	a := make([]NullRow, len(pb.Values))
	for i := range pb.Values {
		switch v := pb.Values[i].GetKind().(type) {
		case *proto3.Value_ListValue:
			a[i] = NullRow{Row: Row{fields: ty.Fields, vals: v.ListValue.Values}, Valid: true}
		case *proto3.Value_NullValue:
		default:
			return nil, errNotStructElement(i, pb.Values[i])
		}
	}
	return a, nil
}
func gologoo__errNilSpannerStructType_6a1b612599996733c2066401f14792cf() error {
	return spannerErrorf(codes.FailedPrecondition, "unexpected nil StructType in decoding Cloud Spanner STRUCT")
}
func gologoo__errDupGoField_6a1b612599996733c2066401f14792cf(s interface {
}, name string) error {
	return spannerErrorf(codes.InvalidArgument, "Go struct %+v(type %T) has duplicate fields for GO STRUCT field %s", s, s, name)
}
func gologoo__errUnnamedField_6a1b612599996733c2066401f14792cf(ty *sppb.StructType, i int) error {
	return spannerErrorf(codes.InvalidArgument, "unnamed field %v in Cloud Spanner STRUCT %+v", i, ty)
}
func gologoo__errNoOrDupGoField_6a1b612599996733c2066401f14792cf(s interface {
}, f string) error {
	return spannerErrorf(codes.InvalidArgument, "Go struct %+v(type %T) has no or duplicate fields for Cloud Spanner STRUCT field %v", s, s, f)
}
func gologoo__errDupSpannerField_6a1b612599996733c2066401f14792cf(f string, ty *sppb.StructType) error {
	return spannerErrorf(codes.InvalidArgument, "duplicated field name %q in Cloud Spanner STRUCT %+v", f, ty)
}
func gologoo__errDecodeStructField_6a1b612599996733c2066401f14792cf(ty *sppb.StructType, f string, err error) error {
	var se *Error
	if !errorAs(err, &se) {
		return spannerErrorf(codes.Unknown, "cannot decode field %v of Cloud Spanner STRUCT %+v, error = <%v>", f, ty, err)
	}
	se.decorate(fmt.Sprintf("cannot decode field %v of Cloud Spanner STRUCT %+v", f, ty))
	return se
}
func gologoo__decodeStruct_6a1b612599996733c2066401f14792cf(ty *sppb.StructType, pb *proto3.ListValue, ptr interface {
}, lenient bool) error {
	if reflect.ValueOf(ptr).IsNil() {
		return errNilDst(ptr)
	}
	if ty == nil {
		return errNilSpannerStructType()
	}
	t := reflect.TypeOf(ptr).Elem()
	v := reflect.ValueOf(ptr).Elem()
	fields, err := fieldCache.Fields(t)
	if err != nil {
		return ToSpannerError(err)
	}
	if lenient {
		fieldNames := getAllFieldNames(v)
		for _, f := range fieldNames {
			if fields.Match(f) == nil {
				return errDupGoField(ptr, f)
			}
		}
	}
	seen := map[string]bool{}
	for i, f := range ty.Fields {
		if f.Name == "" {
			return errUnnamedField(ty, i)
		}
		sf := fields.Match(f.Name)
		if sf == nil {
			if lenient {
				continue
			}
			return errNoOrDupGoField(ptr, f.Name)
		}
		if seen[f.Name] {
			return errDupSpannerField(f.Name, ty)
		}
		if err := decodeValue(pb.Values[i], f.Type, v.FieldByIndex(sf.Index).Addr().Interface()); err != nil {
			return errDecodeStructField(ty, f.Name, err)
		}
		seen[f.Name] = true
	}
	return nil
}
func gologoo__isPtrStructPtrSlice_6a1b612599996733c2066401f14792cf(t reflect.Type) bool {
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Slice {
		return false
	}
	if t = t.Elem(); t.Elem().Kind() != reflect.Ptr || t.Elem().Elem().Kind() != reflect.Struct {
		return false
	}
	return true
}
func gologoo__decodeStructArray_6a1b612599996733c2066401f14792cf(ty *sppb.StructType, pb *proto3.ListValue, ptr interface {
}) error {
	if pb == nil {
		return errNilListValue("STRUCT")
	}
	ts := reflect.TypeOf(ptr).Elem().Elem()
	v := reflect.ValueOf(ptr).Elem()
	v.Set(reflect.MakeSlice(v.Type(), 0, len(pb.Values)))
	for i, pv := range pb.Values {
		if _, isNull := pv.Kind.(*proto3.Value_NullValue); isNull {
			v.Set(reflect.Append(v, reflect.New(ts).Elem()))
			continue
		}
		s := reflect.New(ts.Elem())
		l, err := getListValue(pv)
		if err != nil {
			return errDecodeArrayElement(i, pv, "STRUCT", err)
		}
		if err = decodeStruct(ty, l, s.Interface(), false); err != nil {
			return errDecodeArrayElement(i, pv, "STRUCT", err)
		}
		v.Set(reflect.Append(v, s))
	}
	return nil
}
func gologoo__getAllFieldNames_6a1b612599996733c2066401f14792cf(v reflect.Value) []string {
	var names []string
	typeOfT := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fieldType := typeOfT.Field(i)
		exported := (fieldType.PkgPath == "")
		if !exported && !fieldType.Anonymous {
			continue
		}
		if f.Kind() == reflect.Struct {
			if fieldType.Anonymous {
				names = append(names, getAllFieldNames(reflect.ValueOf(f.Interface()))...)
			}
			continue
		}
		name, keep, _, _ := spannerTagParser(fieldType.Tag)
		if !keep {
			continue
		}
		if name == "" {
			name = fieldType.Name
		}
		names = append(names, name)
	}
	return names
}
func gologoo__errEncoderUnsupportedType_6a1b612599996733c2066401f14792cf(v interface {
}) error {
	return spannerErrorf(codes.InvalidArgument, "client doesn't support type %T", v)
}
func gologoo__encodeValue_6a1b612599996733c2066401f14792cf(v interface {
}) (*proto3.Value, *sppb.Type, error) {
	pb := &proto3.Value{Kind: &proto3.Value_NullValue{NullValue: proto3.NullValue_NULL_VALUE}}
	var pt *sppb.Type
	var err error
	switch v := v.(type) {
	case nil:
	case string:
		pb.Kind = stringKind(v)
		pt = stringType()
	case NullString:
		if v.Valid {
			return encodeValue(v.StringVal)
		}
		pt = stringType()
	case sql.NullString:
		if v.Valid {
			return encodeValue(v.String)
		}
		pt = stringType()
	case []string:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(stringType())
	case []NullString:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(stringType())
	case *string:
		if v != nil {
			return encodeValue(*v)
		}
		pt = stringType()
	case []*string:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(stringType())
	case []byte:
		if v != nil {
			pb.Kind = stringKind(base64.StdEncoding.EncodeToString(v))
		}
		pt = bytesType()
	case [][]byte:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(bytesType())
	case int:
		pb.Kind = stringKind(strconv.FormatInt(int64(v), 10))
		pt = intType()
	case []int:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(intType())
	case int64:
		pb.Kind = stringKind(strconv.FormatInt(v, 10))
		pt = intType()
	case []int64:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(intType())
	case NullInt64:
		if v.Valid {
			return encodeValue(v.Int64)
		}
		pt = intType()
	case []NullInt64:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(intType())
	case *int64:
		if v != nil {
			return encodeValue(*v)
		}
		pt = intType()
	case []*int64:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(intType())
	case bool:
		pb.Kind = &proto3.Value_BoolValue{BoolValue: v}
		pt = boolType()
	case []bool:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(boolType())
	case NullBool:
		if v.Valid {
			return encodeValue(v.Bool)
		}
		pt = boolType()
	case []NullBool:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(boolType())
	case *bool:
		if v != nil {
			return encodeValue(*v)
		}
		pt = boolType()
	case []*bool:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(boolType())
	case float64:
		pb.Kind = &proto3.Value_NumberValue{NumberValue: v}
		pt = floatType()
	case []float64:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(floatType())
	case NullFloat64:
		if v.Valid {
			return encodeValue(v.Float64)
		}
		pt = floatType()
	case []NullFloat64:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(floatType())
	case *float64:
		if v != nil {
			return encodeValue(*v)
		}
		pt = floatType()
	case []*float64:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(floatType())
	case big.Rat:
		switch LossOfPrecisionHandling {
		case NumericError:
			err = validateNumeric(&v)
			if err != nil {
				return nil, nil, err
			}
		case NumericRound:
		}
		pb.Kind = stringKind(NumericString(&v))
		pt = numericType()
	case []big.Rat:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(numericType())
	case NullNumeric:
		if v.Valid {
			return encodeValue(v.Numeric)
		}
		pt = numericType()
	case []NullNumeric:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(numericType())
	case PGNumeric:
		if v.Valid {
			pb.Kind = stringKind(v.Numeric)
		}
		return pb, pgNumericType(), nil
	case []PGNumeric:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(pgNumericType())
	case NullJSON:
		if v.Valid {
			b, err := json.Marshal(v.Value)
			if err != nil {
				return nil, nil, err
			}
			pb.Kind = stringKind(string(b))
		}
		return pb, jsonType(), nil
	case []NullJSON:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(jsonType())
	case *big.Rat:
		switch LossOfPrecisionHandling {
		case NumericError:
			err = validateNumeric(v)
			if err != nil {
				return nil, nil, err
			}
		case NumericRound:
		}
		if v != nil {
			pb.Kind = stringKind(NumericString(v))
		}
		pt = numericType()
	case []*big.Rat:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(numericType())
	case time.Time:
		if v == commitTimestamp {
			pb.Kind = stringKind(commitTimestampPlaceholderString)
		} else {
			pb.Kind = stringKind(v.UTC().Format(time.RFC3339Nano))
		}
		pt = timeType()
	case []time.Time:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(timeType())
	case NullTime:
		if v.Valid {
			return encodeValue(v.Time)
		}
		pt = timeType()
	case []NullTime:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(timeType())
	case *time.Time:
		if v != nil {
			return encodeValue(*v)
		}
		pt = timeType()
	case []*time.Time:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(timeType())
	case civil.Date:
		pb.Kind = stringKind(v.String())
		pt = dateType()
	case []civil.Date:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(dateType())
	case NullDate:
		if v.Valid {
			return encodeValue(v.Date)
		}
		pt = dateType()
	case []NullDate:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(dateType())
	case *civil.Date:
		if v != nil {
			return encodeValue(*v)
		}
		pt = dateType()
	case []*civil.Date:
		if v != nil {
			pb, err = encodeArray(len(v), func(i int) interface {
			} {
				return v[i]
			})
			if err != nil {
				return nil, nil, err
			}
		}
		pt = listType(dateType())
	case GenericColumnValue:
		pb = proto.Clone(v.Value).(*proto3.Value)
		pt = proto.Clone(v.Type).(*sppb.Type)
	case []GenericColumnValue:
		return nil, nil, errEncoderUnsupportedType(v)
	default:
		if encodedVal, ok := v.(Encoder); ok {
			nv, err := encodedVal.EncodeSpanner()
			if err != nil {
				return nil, nil, err
			}
			return encodeValue(nv)
		}
		decodableType := getDecodableSpannerType(v, false)
		if decodableType != spannerTypeUnknown && decodableType != spannerTypeInvalid {
			converted, err := convertCustomTypeValue(decodableType, v)
			if err != nil {
				return nil, nil, err
			}
			return encodeValue(converted)
		}
		if !isStructOrArrayOfStructValue(v) {
			return nil, nil, errEncoderUnsupportedType(v)
		}
		typ := reflect.TypeOf(v)
		if (typ.Kind() == reflect.Struct) || (typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct) {
			return encodeStruct(v)
		}
		if typ.Kind() == reflect.Slice {
			return encodeStructArray(v)
		}
	}
	return pb, pt, nil
}
func gologoo__convertCustomTypeValue_6a1b612599996733c2066401f14792cf(sourceType decodableSpannerType, v interface {
}) (interface {
}, error) {
	var destination reflect.Value
	switch sourceType {
	case spannerTypeInvalid:
		return nil, fmt.Errorf("cannot encode a value to type spannerTypeInvalid")
	case spannerTypeNonNullString:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf("")))
	case spannerTypeNullString:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(NullString{})))
	case spannerTypeByteArray:
		if reflect.ValueOf(v).IsNil() {
			return []byte(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]byte{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeNonNullInt64:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(int64(0))))
	case spannerTypeNullInt64:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(NullInt64{})))
	case spannerTypeNonNullBool:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(false)))
	case spannerTypeNullBool:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(NullBool{})))
	case spannerTypeNonNullFloat64:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(float64(0.0))))
	case spannerTypeNullFloat64:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(NullFloat64{})))
	case spannerTypeNonNullTime:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(time.Time{})))
	case spannerTypeNullTime:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(NullTime{})))
	case spannerTypeNonNullDate:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(civil.Date{})))
	case spannerTypeNullDate:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(NullDate{})))
	case spannerTypeNonNullNumeric:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(big.Rat{})))
	case spannerTypeNullNumeric:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(NullNumeric{})))
	case spannerTypeNullJSON:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(NullJSON{})))
	case spannerTypePGNumeric:
		destination = reflect.Indirect(reflect.New(reflect.TypeOf(PGNumeric{})))
	case spannerTypeArrayOfNonNullString:
		if reflect.ValueOf(v).IsNil() {
			return []string(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]string{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNullString:
		if reflect.ValueOf(v).IsNil() {
			return []NullString(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]NullString{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfByteArray:
		if reflect.ValueOf(v).IsNil() {
			return [][]byte(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([][]byte{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNonNullInt64:
		if reflect.ValueOf(v).IsNil() {
			return []int64(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]int64{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNullInt64:
		if reflect.ValueOf(v).IsNil() {
			return []NullInt64(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]NullInt64{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNonNullBool:
		if reflect.ValueOf(v).IsNil() {
			return []bool(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]bool{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNullBool:
		if reflect.ValueOf(v).IsNil() {
			return []NullBool(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]NullBool{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNonNullFloat64:
		if reflect.ValueOf(v).IsNil() {
			return []float64(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]float64{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNullFloat64:
		if reflect.ValueOf(v).IsNil() {
			return []NullFloat64(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]NullFloat64{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNonNullTime:
		if reflect.ValueOf(v).IsNil() {
			return []time.Time(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]time.Time{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNullTime:
		if reflect.ValueOf(v).IsNil() {
			return []NullTime(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]NullTime{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNonNullDate:
		if reflect.ValueOf(v).IsNil() {
			return []civil.Date(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]civil.Date{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNullDate:
		if reflect.ValueOf(v).IsNil() {
			return []NullDate(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]NullDate{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNonNullNumeric:
		if reflect.ValueOf(v).IsNil() {
			return []big.Rat(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]big.Rat{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNullNumeric:
		if reflect.ValueOf(v).IsNil() {
			return []NullNumeric(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]NullNumeric{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfNullJSON:
		if reflect.ValueOf(v).IsNil() {
			return []NullJSON(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]NullJSON{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	case spannerTypeArrayOfPGNumeric:
		if reflect.ValueOf(v).IsNil() {
			return []PGNumeric(nil), nil
		}
		destination = reflect.MakeSlice(reflect.TypeOf([]PGNumeric{}), reflect.ValueOf(v).Len(), reflect.ValueOf(v).Cap())
	default:
		return nil, fmt.Errorf("unknown decodable type found: %v", sourceType)
	}
	if destination.Kind() == reflect.Slice || destination.Kind() == reflect.Array {
		sourceSlice := reflect.ValueOf(v)
		for i := 0; i < destination.Len(); i++ {
			source := sourceSlice.Index(i)
			destination.Index(i).Set(source.Convert(destination.Type().Elem()))
		}
	} else {
		source := reflect.ValueOf(v)
		destination.Set(source.Convert(destination.Type()))
	}
	return destination.Interface(), nil
}
func gologoo__encodeStruct_6a1b612599996733c2066401f14792cf(v interface {
}) (*proto3.Value, *sppb.Type, error) {
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		typ = typ.Elem()
		if val.IsNil() {
			_, st, err := encodeStruct(reflect.Zero(typ).Interface())
			if err != nil {
				return nil, nil, err
			}
			return nullProto(), st, nil
		}
		val = val.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, nil, errEncoderUnsupportedType(v)
	}
	stf := make([]*sppb.StructType_Field, 0, typ.NumField())
	stv := make([]*proto3.Value, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		fval := val.Field(i)
		if sf.Anonymous {
			return nil, nil, errUnsupportedEmbeddedStructFields(sf.Name)
		}
		if !fval.CanInterface() {
			continue
		}
		fname, ok := sf.Tag.Lookup("spanner")
		if !ok {
			fname = sf.Name
		}
		eval, etype, err := encodeValue(fval.Interface())
		if err != nil {
			return nil, nil, err
		}
		stf = append(stf, mkField(fname, etype))
		stv = append(stv, eval)
	}
	return listProto(stv...), structType(stf...), nil
}
func gologoo__encodeStructArray_6a1b612599996733c2066401f14792cf(v interface {
}) (*proto3.Value, *sppb.Type, error) {
	etyp := reflect.TypeOf(v).Elem()
	sliceval := reflect.ValueOf(v)
	if etyp.Kind() == reflect.Ptr {
		etyp = etyp.Elem()
	}
	_, elemTyp, err := encodeStruct(reflect.Zero(etyp).Interface())
	if err != nil {
		return nil, nil, err
	}
	if sliceval.IsNil() {
		return nullProto(), listType(elemTyp), nil
	}
	values := make([]*proto3.Value, 0, sliceval.Len())
	for i := 0; i < sliceval.Len(); i++ {
		ev, _, err := encodeStruct(sliceval.Index(i).Interface())
		if err != nil {
			return nil, nil, err
		}
		values = append(values, ev)
	}
	return listProto(values...), listType(elemTyp), nil
}
func gologoo__isStructOrArrayOfStructValue_6a1b612599996733c2066401f14792cf(v interface {
}) bool {
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ.Kind() == reflect.Struct
}
func gologoo__isSupportedMutationType_6a1b612599996733c2066401f14792cf(v interface {
}) bool {
	switch v.(type) {
	case nil, string, *string, NullString, []string, []*string, []NullString, []byte, [][]byte, int, []int, int64, *int64, []int64, []*int64, NullInt64, []NullInt64, bool, *bool, []bool, []*bool, NullBool, []NullBool, float64, *float64, []float64, []*float64, NullFloat64, []NullFloat64, time.Time, *time.Time, []time.Time, []*time.Time, NullTime, []NullTime, civil.Date, *civil.Date, []civil.Date, []*civil.Date, NullDate, []NullDate, big.Rat, *big.Rat, []big.Rat, []*big.Rat, NullNumeric, []NullNumeric, GenericColumnValue:
		return true
	default:
		if _, ok := v.(Encoder); ok {
			return true
		}
		decodableType := getDecodableSpannerType(v, false)
		return decodableType != spannerTypeUnknown && decodableType != spannerTypeInvalid
	}
}
func gologoo__encodeValueArray_6a1b612599996733c2066401f14792cf(vs []interface {
}) (*proto3.ListValue, error) {
	lv := &proto3.ListValue{}
	lv.Values = make([]*proto3.Value, 0, len(vs))
	for _, v := range vs {
		if !isSupportedMutationType(v) {
			return nil, errEncoderUnsupportedType(v)
		}
		pb, _, err := encodeValue(v)
		if err != nil {
			return nil, err
		}
		lv.Values = append(lv.Values, pb)
	}
	return lv, nil
}
func gologoo__encodeArray_6a1b612599996733c2066401f14792cf(len int, at func(int) interface {
}) (*proto3.Value, error) {
	vs := make([]*proto3.Value, len)
	var err error
	for i := 0; i < len; i++ {
		vs[i], _, err = encodeValue(at(i))
		if err != nil {
			return nil, err
		}
	}
	return listProto(vs...), nil
}
func gologoo__spannerTagParser_6a1b612599996733c2066401f14792cf(t reflect.StructTag) (name string, keep bool, other interface {
}, err error) {
	if s := t.Get("spanner"); s != "" {
		if s == "-" {
			return "", false, nil, nil
		}
		return s, true, nil, nil
	}
	return "", true, nil, nil
}

var fieldCache = fields.NewCache(spannerTagParser, nil, nil)

func gologoo__trimDoubleQuotes_6a1b612599996733c2066401f14792cf(payload []byte) ([]byte, error) {
	if len(payload) <= 1 || payload[0] != '"' || payload[len(payload)-1] != '"' {
		return nil, fmt.Errorf("payload is too short or not wrapped with double quotes: got %q", string(payload))
	}
	return payload[1 : len(payload)-1], nil
}
func NumericString(r *big.Rat) string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NumericString_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", r)
	r0 := gologoo__NumericString_6a1b612599996733c2066401f14792cf(r)
	log.Printf("Output: %v\n", r0)
	return r0
}
func validateNumeric(r *big.Rat) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__validateNumeric_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", r)
	r0 := gologoo__validateNumeric_6a1b612599996733c2066401f14792cf(r)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullInt64) IsNull() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__IsNull_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__IsNull_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullInt64) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__String_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullInt64) MarshalJSON() ([]byte, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0, r1 := n.gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (n *NullInt64) UnmarshalJSON(payload []byte) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", payload)
	r0 := n.gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullInt64) Value() (driver.Value, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Value_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0, r1 := n.gologoo__Value_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (n *NullInt64) Scan(value interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Scan_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", value)
	r0 := n.gologoo__Scan_6a1b612599996733c2066401f14792cf(value)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullInt64) GormDataType() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GormDataType_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__GormDataType_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullString) IsNull() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__IsNull_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__IsNull_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullString) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__String_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullString) MarshalJSON() ([]byte, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0, r1 := n.gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (n *NullString) UnmarshalJSON(payload []byte) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", payload)
	r0 := n.gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullString) Value() (driver.Value, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Value_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0, r1 := n.gologoo__Value_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (n *NullString) Scan(value interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Scan_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", value)
	r0 := n.gologoo__Scan_6a1b612599996733c2066401f14792cf(value)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullString) GormDataType() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GormDataType_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__GormDataType_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullFloat64) IsNull() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__IsNull_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__IsNull_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullFloat64) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__String_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullFloat64) MarshalJSON() ([]byte, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0, r1 := n.gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (n *NullFloat64) UnmarshalJSON(payload []byte) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", payload)
	r0 := n.gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullFloat64) Value() (driver.Value, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Value_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0, r1 := n.gologoo__Value_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (n *NullFloat64) Scan(value interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Scan_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", value)
	r0 := n.gologoo__Scan_6a1b612599996733c2066401f14792cf(value)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullFloat64) GormDataType() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GormDataType_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__GormDataType_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullBool) IsNull() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__IsNull_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__IsNull_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullBool) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__String_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullBool) MarshalJSON() ([]byte, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0, r1 := n.gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (n *NullBool) UnmarshalJSON(payload []byte) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", payload)
	r0 := n.gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullBool) Value() (driver.Value, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Value_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0, r1 := n.gologoo__Value_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (n *NullBool) Scan(value interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Scan_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", value)
	r0 := n.gologoo__Scan_6a1b612599996733c2066401f14792cf(value)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullBool) GormDataType() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GormDataType_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__GormDataType_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullTime) IsNull() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__IsNull_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__IsNull_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullTime) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__String_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullTime) MarshalJSON() ([]byte, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0, r1 := n.gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (n *NullTime) UnmarshalJSON(payload []byte) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", payload)
	r0 := n.gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullTime) Value() (driver.Value, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Value_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0, r1 := n.gologoo__Value_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (n *NullTime) Scan(value interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Scan_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", value)
	r0 := n.gologoo__Scan_6a1b612599996733c2066401f14792cf(value)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullTime) GormDataType() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GormDataType_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__GormDataType_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullDate) IsNull() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__IsNull_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__IsNull_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullDate) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__String_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullDate) MarshalJSON() ([]byte, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0, r1 := n.gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (n *NullDate) UnmarshalJSON(payload []byte) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", payload)
	r0 := n.gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullDate) Value() (driver.Value, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Value_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0, r1 := n.gologoo__Value_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (n *NullDate) Scan(value interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Scan_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", value)
	r0 := n.gologoo__Scan_6a1b612599996733c2066401f14792cf(value)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullDate) GormDataType() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GormDataType_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__GormDataType_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullNumeric) IsNull() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__IsNull_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__IsNull_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullNumeric) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__String_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullNumeric) MarshalJSON() ([]byte, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0, r1 := n.gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (n *NullNumeric) UnmarshalJSON(payload []byte) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", payload)
	r0 := n.gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullNumeric) Value() (driver.Value, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Value_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0, r1 := n.gologoo__Value_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (n *NullNumeric) Scan(value interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Scan_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", value)
	r0 := n.gologoo__Scan_6a1b612599996733c2066401f14792cf(value)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullNumeric) GormDataType() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GormDataType_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__GormDataType_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullJSON) IsNull() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__IsNull_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__IsNull_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullJSON) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__String_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullJSON) MarshalJSON() ([]byte, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0, r1 := n.gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (n *NullJSON) UnmarshalJSON(payload []byte) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", payload)
	r0 := n.gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n NullJSON) GormDataType() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GormDataType_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__GormDataType_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n PGNumeric) IsNull() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__IsNull_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__IsNull_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n PGNumeric) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := n.gologoo__String_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (n PGNumeric) MarshalJSON() ([]byte, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0, r1 := n.gologoo__MarshalJSON_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (n *PGNumeric) UnmarshalJSON(payload []byte) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", payload)
	r0 := n.gologoo__UnmarshalJSON_6a1b612599996733c2066401f14792cf(payload)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (v GenericColumnValue) Decode(ptr interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Decode_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", ptr)
	r0 := v.gologoo__Decode_6a1b612599996733c2066401f14792cf(ptr)
	log.Printf("Output: %v\n", r0)
	return r0
}
func newGenericColumnValue(v interface {
}) (*GenericColumnValue, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__newGenericColumnValue_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", v)
	r0, r1 := gologoo__newGenericColumnValue_6a1b612599996733c2066401f14792cf(v)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func errTypeMismatch(srcCode, elCode sppb.TypeCode, dst interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errTypeMismatch_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v %v\n", srcCode, elCode, dst)
	r0 := gologoo__errTypeMismatch_6a1b612599996733c2066401f14792cf(srcCode, elCode, dst)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errNilSpannerType() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errNilSpannerType_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := gologoo__errNilSpannerType_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func errNilSrc() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errNilSrc_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := gologoo__errNilSrc_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func errNilDst(dst interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errNilDst_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", dst)
	r0 := gologoo__errNilDst_6a1b612599996733c2066401f14792cf(dst)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errNilArrElemType(t *sppb.Type) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errNilArrElemType_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", t)
	r0 := gologoo__errNilArrElemType_6a1b612599996733c2066401f14792cf(t)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errUnsupportedEmbeddedStructFields(fname string) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errUnsupportedEmbeddedStructFields_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", fname)
	r0 := gologoo__errUnsupportedEmbeddedStructFields_6a1b612599996733c2066401f14792cf(fname)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errDstNotForNull(dst interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errDstNotForNull_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", dst)
	r0 := gologoo__errDstNotForNull_6a1b612599996733c2066401f14792cf(dst)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errBadEncoding(v *proto3.Value, err error) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errBadEncoding_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v\n", v, err)
	r0 := gologoo__errBadEncoding_6a1b612599996733c2066401f14792cf(v, err)
	log.Printf("Output: %v\n", r0)
	return r0
}
func parseNullTime(v *proto3.Value, p *NullTime, code sppb.TypeCode, isNull bool) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseNullTime_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v %v %v\n", v, p, code, isNull)
	r0 := gologoo__parseNullTime_6a1b612599996733c2066401f14792cf(v, p, code, isNull)
	log.Printf("Output: %v\n", r0)
	return r0
}
func decodeValue(v *proto3.Value, t *sppb.Type, ptr interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeValue_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v %v\n", v, t, ptr)
	r0 := gologoo__decodeValue_6a1b612599996733c2066401f14792cf(v, t, ptr)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (d decodableSpannerType) supportsNull() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__supportsNull_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := d.gologoo__supportsNull_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func getDecodableSpannerType(ptr interface {
}, isPtr bool) decodableSpannerType {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getDecodableSpannerType_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v\n", ptr, isPtr)
	r0 := gologoo__getDecodableSpannerType_6a1b612599996733c2066401f14792cf(ptr, isPtr)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (dsc decodableSpannerType) decodeValueToCustomType(v *proto3.Value, t *sppb.Type, acode sppb.TypeCode, atypeAnnotation sppb.TypeAnnotationCode, ptr interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeValueToCustomType_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v %v %v %v\n", v, t, acode, atypeAnnotation, ptr)
	r0 := dsc.gologoo__decodeValueToCustomType_6a1b612599996733c2066401f14792cf(v, t, acode, atypeAnnotation, ptr)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errSrcVal(v *proto3.Value, want string) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errSrcVal_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v\n", v, want)
	r0 := gologoo__errSrcVal_6a1b612599996733c2066401f14792cf(v, want)
	log.Printf("Output: %v\n", r0)
	return r0
}
func getStringValue(v *proto3.Value) (string, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getStringValue_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", v)
	r0, r1 := gologoo__getStringValue_6a1b612599996733c2066401f14792cf(v)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func getBoolValue(v *proto3.Value) (bool, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getBoolValue_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", v)
	r0, r1 := gologoo__getBoolValue_6a1b612599996733c2066401f14792cf(v)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func getListValue(v *proto3.Value) (*proto3.ListValue, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getListValue_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", v)
	r0, r1 := gologoo__getListValue_6a1b612599996733c2066401f14792cf(v)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func getGenericValue(t *sppb.Type, v *proto3.Value) (interface {
}, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getGenericValue_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v\n", t, v)
	r0, r1 := gologoo__getGenericValue_6a1b612599996733c2066401f14792cf(t, v)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func getTypedNil(t *sppb.Type) (interface {
}, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getTypedNil_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", t)
	r0, r1 := gologoo__getTypedNil_6a1b612599996733c2066401f14792cf(t)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func errUnexpectedNumericStr(s string) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errUnexpectedNumericStr_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", s)
	r0 := gologoo__errUnexpectedNumericStr_6a1b612599996733c2066401f14792cf(s)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errUnexpectedFloat64Str(s string) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errUnexpectedFloat64Str_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", s)
	r0 := gologoo__errUnexpectedFloat64Str_6a1b612599996733c2066401f14792cf(s)
	log.Printf("Output: %v\n", r0)
	return r0
}
func getFloat64Value(v *proto3.Value) (float64, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getFloat64Value_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", v)
	r0, r1 := gologoo__getFloat64Value_6a1b612599996733c2066401f14792cf(v)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func errNilListValue(sqlType string) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errNilListValue_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", sqlType)
	r0 := gologoo__errNilListValue_6a1b612599996733c2066401f14792cf(sqlType)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errDecodeArrayElement(i int, v proto.Message, sqlType string, err error) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errDecodeArrayElement_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v %v %v\n", i, v, sqlType, err)
	r0 := gologoo__errDecodeArrayElement_6a1b612599996733c2066401f14792cf(i, v, sqlType, err)
	log.Printf("Output: %v\n", r0)
	return r0
}
func decodeGenericArray(tp reflect.Type, pb *proto3.ListValue, t *sppb.Type, sqlType string) (interface {
}, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeGenericArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v %v %v\n", tp, pb, t, sqlType)
	r0, r1 := gologoo__decodeGenericArray_6a1b612599996733c2066401f14792cf(tp, pb, t, sqlType)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeNullStringArray(pb *proto3.ListValue) ([]NullString, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeNullStringArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeNullStringArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeStringPointerArray(pb *proto3.ListValue) ([]*string, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeStringPointerArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeStringPointerArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeStringArray(pb *proto3.ListValue) ([]string, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeStringArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeStringArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeNullInt64Array(pb *proto3.ListValue) ([]NullInt64, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeNullInt64Array_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeNullInt64Array_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeInt64PointerArray(pb *proto3.ListValue) ([]*int64, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeInt64PointerArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeInt64PointerArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeInt64Array(pb *proto3.ListValue) ([]int64, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeInt64Array_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeInt64Array_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeNullBoolArray(pb *proto3.ListValue) ([]NullBool, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeNullBoolArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeNullBoolArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeBoolPointerArray(pb *proto3.ListValue) ([]*bool, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeBoolPointerArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeBoolPointerArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeBoolArray(pb *proto3.ListValue) ([]bool, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeBoolArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeBoolArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeNullFloat64Array(pb *proto3.ListValue) ([]NullFloat64, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeNullFloat64Array_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeNullFloat64Array_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeFloat64PointerArray(pb *proto3.ListValue) ([]*float64, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeFloat64PointerArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeFloat64PointerArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeFloat64Array(pb *proto3.ListValue) ([]float64, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeFloat64Array_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeFloat64Array_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeNullNumericArray(pb *proto3.ListValue) ([]NullNumeric, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeNullNumericArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeNullNumericArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeNullJSONArray(pb *proto3.ListValue) ([]NullJSON, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeNullJSONArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeNullJSONArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeNullJSONArrayToNullJSON(pb *proto3.ListValue) (*NullJSON, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeNullJSONArrayToNullJSON_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeNullJSONArrayToNullJSON_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeNumericPointerArray(pb *proto3.ListValue) ([]*big.Rat, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeNumericPointerArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeNumericPointerArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeNumericArray(pb *proto3.ListValue) ([]big.Rat, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeNumericArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeNumericArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodePGNumericArray(pb *proto3.ListValue) ([]PGNumeric, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodePGNumericArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodePGNumericArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeByteArray(pb *proto3.ListValue) ([][]byte, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeByteArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeByteArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeNullTimeArray(pb *proto3.ListValue) ([]NullTime, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeNullTimeArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeNullTimeArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeTimePointerArray(pb *proto3.ListValue) ([]*time.Time, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeTimePointerArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeTimePointerArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeTimeArray(pb *proto3.ListValue) ([]time.Time, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeTimeArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeTimeArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeNullDateArray(pb *proto3.ListValue) ([]NullDate, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeNullDateArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeNullDateArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeDatePointerArray(pb *proto3.ListValue) ([]*civil.Date, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeDatePointerArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeDatePointerArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func decodeDateArray(pb *proto3.ListValue) ([]civil.Date, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeDateArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", pb)
	r0, r1 := gologoo__decodeDateArray_6a1b612599996733c2066401f14792cf(pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func errNotStructElement(i int, v *proto3.Value) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errNotStructElement_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v\n", i, v)
	r0 := gologoo__errNotStructElement_6a1b612599996733c2066401f14792cf(i, v)
	log.Printf("Output: %v\n", r0)
	return r0
}
func decodeRowArray(ty *sppb.StructType, pb *proto3.ListValue) ([]NullRow, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeRowArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v\n", ty, pb)
	r0, r1 := gologoo__decodeRowArray_6a1b612599996733c2066401f14792cf(ty, pb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func errNilSpannerStructType() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errNilSpannerStructType_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : (none)\n")
	r0 := gologoo__errNilSpannerStructType_6a1b612599996733c2066401f14792cf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func errDupGoField(s interface {
}, name string) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errDupGoField_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v\n", s, name)
	r0 := gologoo__errDupGoField_6a1b612599996733c2066401f14792cf(s, name)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errUnnamedField(ty *sppb.StructType, i int) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errUnnamedField_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v\n", ty, i)
	r0 := gologoo__errUnnamedField_6a1b612599996733c2066401f14792cf(ty, i)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errNoOrDupGoField(s interface {
}, f string) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errNoOrDupGoField_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v\n", s, f)
	r0 := gologoo__errNoOrDupGoField_6a1b612599996733c2066401f14792cf(s, f)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errDupSpannerField(f string, ty *sppb.StructType) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errDupSpannerField_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v\n", f, ty)
	r0 := gologoo__errDupSpannerField_6a1b612599996733c2066401f14792cf(f, ty)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errDecodeStructField(ty *sppb.StructType, f string, err error) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errDecodeStructField_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v %v\n", ty, f, err)
	r0 := gologoo__errDecodeStructField_6a1b612599996733c2066401f14792cf(ty, f, err)
	log.Printf("Output: %v\n", r0)
	return r0
}
func decodeStruct(ty *sppb.StructType, pb *proto3.ListValue, ptr interface {
}, lenient bool) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeStruct_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v %v %v\n", ty, pb, ptr, lenient)
	r0 := gologoo__decodeStruct_6a1b612599996733c2066401f14792cf(ty, pb, ptr, lenient)
	log.Printf("Output: %v\n", r0)
	return r0
}
func isPtrStructPtrSlice(t reflect.Type) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isPtrStructPtrSlice_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", t)
	r0 := gologoo__isPtrStructPtrSlice_6a1b612599996733c2066401f14792cf(t)
	log.Printf("Output: %v\n", r0)
	return r0
}
func decodeStructArray(ty *sppb.StructType, pb *proto3.ListValue, ptr interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decodeStructArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v %v\n", ty, pb, ptr)
	r0 := gologoo__decodeStructArray_6a1b612599996733c2066401f14792cf(ty, pb, ptr)
	log.Printf("Output: %v\n", r0)
	return r0
}
func getAllFieldNames(v reflect.Value) []string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getAllFieldNames_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", v)
	r0 := gologoo__getAllFieldNames_6a1b612599996733c2066401f14792cf(v)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errEncoderUnsupportedType(v interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errEncoderUnsupportedType_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", v)
	r0 := gologoo__errEncoderUnsupportedType_6a1b612599996733c2066401f14792cf(v)
	log.Printf("Output: %v\n", r0)
	return r0
}
func encodeValue(v interface {
}) (*proto3.Value, *sppb.Type, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__encodeValue_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", v)
	r0, r1, r2 := gologoo__encodeValue_6a1b612599996733c2066401f14792cf(v)
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
func convertCustomTypeValue(sourceType decodableSpannerType, v interface {
}) (interface {
}, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__convertCustomTypeValue_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v\n", sourceType, v)
	r0, r1 := gologoo__convertCustomTypeValue_6a1b612599996733c2066401f14792cf(sourceType, v)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func encodeStruct(v interface {
}) (*proto3.Value, *sppb.Type, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__encodeStruct_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", v)
	r0, r1, r2 := gologoo__encodeStruct_6a1b612599996733c2066401f14792cf(v)
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
func encodeStructArray(v interface {
}) (*proto3.Value, *sppb.Type, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__encodeStructArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", v)
	r0, r1, r2 := gologoo__encodeStructArray_6a1b612599996733c2066401f14792cf(v)
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
func isStructOrArrayOfStructValue(v interface {
}) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isStructOrArrayOfStructValue_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", v)
	r0 := gologoo__isStructOrArrayOfStructValue_6a1b612599996733c2066401f14792cf(v)
	log.Printf("Output: %v\n", r0)
	return r0
}
func isSupportedMutationType(v interface {
}) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isSupportedMutationType_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", v)
	r0 := gologoo__isSupportedMutationType_6a1b612599996733c2066401f14792cf(v)
	log.Printf("Output: %v\n", r0)
	return r0
}
func encodeValueArray(vs []interface {
}) (*proto3.ListValue, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__encodeValueArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", vs)
	r0, r1 := gologoo__encodeValueArray_6a1b612599996733c2066401f14792cf(vs)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func encodeArray(len int, at func(int) interface {
}) (*proto3.Value, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__encodeArray_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v %v\n", len, at)
	r0, r1 := gologoo__encodeArray_6a1b612599996733c2066401f14792cf(len, at)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func spannerTagParser(t reflect.StructTag) (name string, keep bool, other interface {
}, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__spannerTagParser_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", t)
	name, keep, other, err = gologoo__spannerTagParser_6a1b612599996733c2066401f14792cf(t)
	log.Printf("Output: %v %v %v %v\n", name, keep, other, err)
	return
}
func trimDoubleQuotes(payload []byte) ([]byte, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__trimDoubleQuotes_6a1b612599996733c2066401f14792cf")
	log.Printf("Input : %v\n", payload)
	r0, r1 := gologoo__trimDoubleQuotes_6a1b612599996733c2066401f14792cf(payload)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
