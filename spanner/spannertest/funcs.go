package spannertest

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"cloud.google.com/go/civil"
	"cloud.google.com/go/spanner/spansql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type function struct {
	Eval func(values []interface {
	}, types []spansql.Type) (interface {
	}, spansql.Type, error)
}

func gologoo__firstErr_3fabff857a70ddbd18407d4ff739347f(errors []error) error {
	for _, err := range errors {
		if err != nil {
			return err
		}
	}
	return nil
}

var functions = map[string]function{"STARTS_WITH": {Eval: func(values []interface {
}, types []spansql.Type) (interface {
}, spansql.Type, error) {
	if len(values) != 2 {
		return nil, spansql.Type{}, status.Error(codes.InvalidArgument, "No matching signature for function STARTS_WITH for the given argument types")
	}
	for _, v := range values {
		if _, ok := v.(string); !ok {
			return nil, spansql.Type{}, status.Error(codes.InvalidArgument, "No matching signature for function STARTS_WITH for the given argument types")
		}
	}
	s := values[0].(string)
	prefix := values[1].(string)
	return strings.HasPrefix(s, prefix), spansql.Type{Base: spansql.Bool}, nil
}}, "LOWER": {Eval: func(values []interface {
}, types []spansql.Type) (interface {
}, spansql.Type, error) {
	if len(values) != 1 {
		return nil, spansql.Type{}, status.Error(codes.InvalidArgument, "No matching signature for function LOWER for the given argument types")
	}
	if values[0] == nil {
		return nil, spansql.Type{Base: spansql.String}, nil
	}
	if _, ok := values[0].(string); !ok {
		return nil, spansql.Type{}, status.Error(codes.InvalidArgument, "No matching signature for function LOWER for the given argument types")
	}
	return strings.ToLower(values[0].(string)), spansql.Type{Base: spansql.String}, nil
}}, "CAST": {Eval: func(values []interface {
}, types []spansql.Type) (interface {
}, spansql.Type, error) {
	return cast(values, types, false)
}}, "SAFE_CAST": {Eval: func(values []interface {
}, types []spansql.Type) (interface {
}, spansql.Type, error) {
	return cast(values, types, true)
}}, "JSON_VALUE": {Eval: func(values []interface {
}, types []spansql.Type) (interface {
}, spansql.Type, error) {
	if len(values) != 2 {
		return nil, spansql.Type{}, status.Error(codes.InvalidArgument, "No matching signature for function JSON_VALUE for the given argument types")
	}
	if values[0] == nil || values[1] == nil {
		return nil, spansql.Type{Base: spansql.String}, nil
	}
	_, okArg1 := values[0].(string)
	_, okArg2 := values[1].(string)
	if !(okArg1 && okArg2) {
		return nil, spansql.Type{}, status.Error(codes.InvalidArgument, "No matching signature for function JSON_VALUE for the given argument types")
	}
	return "", spansql.Type{Base: spansql.String}, nil
}}, "EXTRACT": {Eval: func(values []interface {
}, types []spansql.Type) (interface {
}, spansql.Type, error) {
	date, okArg1 := values[0].(civil.Date)
	part, okArg2 := values[0].(int64)
	if !(okArg1 || okArg2) {
		return nil, spansql.Type{}, status.Error(codes.InvalidArgument, "No matching signature for function EXTRACT for the given argument types")
	}
	if okArg1 {
		return date, spansql.Type{Base: spansql.Date}, nil
	}
	return part, spansql.Type{Base: spansql.Int64}, nil
}}, "TIMESTAMP": {Eval: func(values []interface {
}, types []spansql.Type) (interface {
}, spansql.Type, error) {
	t, okArg1 := values[0].(string)
	if !(okArg1) {
		return nil, spansql.Type{}, status.Error(codes.InvalidArgument, "No matching signature for function TIMESTAMP for the given argument types")
	}
	timestamp, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return nil, spansql.Type{}, status.Error(codes.InvalidArgument, "No matching signature for function TIMESTAMP for the given argument types")
	}
	return timestamp, spansql.Type{Base: spansql.Timestamp}, nil
}}, "FARM_FINGERPRINT": {Eval: func(values []interface {
}, types []spansql.Type) (interface {
}, spansql.Type, error) {
	if len(values) != 1 {
		return nil, spansql.Type{}, status.Error(codes.InvalidArgument, "No matching signature for function FARM_FINGERPRINT for the given argument types")
	}
	if values[0] == nil {
		return int64(1), spansql.Type{Base: spansql.Int64}, nil
	}
	if _, ok := values[0].(string); !ok {
		return nil, spansql.Type{}, status.Error(codes.InvalidArgument, "No matching signature for function FARM_FINGERPRINT for the given argument types")
	}
	return int64(1), spansql.Type{Base: spansql.Int64}, nil
}}, "MOD": {Eval: func(values []interface {
}, types []spansql.Type) (interface {
}, spansql.Type, error) {
	if len(values) != 2 {
		return nil, spansql.Type{}, status.Error(codes.InvalidArgument, "No matching signature for function MOD for the given argument types")
	}
	x, okArg1 := values[0].(int64)
	y, okArg2 := values[1].(int64)
	if !(okArg1 && okArg2) {
		return nil, spansql.Type{}, status.Error(codes.InvalidArgument, "No matching signature for function MOD for the given argument types")
	}
	if y == 0 {
		return nil, spansql.Type{}, status.Error(codes.OutOfRange, "Division by zero")
	}
	return x % y, spansql.Type{Base: spansql.Int64}, nil
}}}

func gologoo__cast_3fabff857a70ddbd18407d4ff739347f(values []interface {
}, types []spansql.Type, safe bool) (interface {
}, spansql.Type, error) {
	name := "CAST"
	if safe {
		name = "SAFE_CAST"
	}
	if len(types) != 1 {
		return nil, spansql.Type{}, status.Errorf(codes.InvalidArgument, "No type information for function %s for the given arguments", name)
	}
	if len(values) != 1 {
		return nil, spansql.Type{}, status.Errorf(codes.InvalidArgument, "No matching signature for function %s for the given arguments", name)
	}
	if err, ok := values[0].(error); ok {
		if safe {
			return nil, types[0], nil
		}
		return nil, types[0], err
	}
	return values[0], types[0], nil
}
func gologoo__convert_3fabff857a70ddbd18407d4ff739347f(val interface {
}, tp spansql.Type) (interface {
}, error) {
	if tp.Array {
		return nil, status.Errorf(codes.Unimplemented, "conversion to ARRAY types is not implemented")
	}
	var res interface {
	}
	var convertErr, err error
	switch tp.Base {
	case spansql.Int64:
		res, convertErr, err = convertToInt64(val)
	case spansql.Float64:
		res, convertErr, err = convertToFloat64(val)
	case spansql.String:
		res, convertErr, err = convertToString(val)
	case spansql.Bool:
		res, convertErr, err = convertToBool(val)
	case spansql.Date:
		res, convertErr, err = convertToDate(val)
	case spansql.Timestamp:
		res, convertErr, err = convertToTimestamp(val)
	case spansql.Numeric:
	case spansql.JSON:
	}
	if err != nil {
		return nil, err
	}
	if convertErr != nil {
		res = convertErr
	}
	if res != nil {
		return res, nil
	}
	return nil, status.Errorf(codes.Unimplemented, "unsupported conversion for %v to %v", val, tp.Base.SQL())
}
func gologoo__convertToInt64_3fabff857a70ddbd18407d4ff739347f(val interface {
}) (res int64, convertErr error, err error) {
	switch v := val.(type) {
	case int64:
		return v, nil, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, status.Errorf(codes.InvalidArgument, "invalid value for INT64: %q", v), nil
		}
		return res, nil, nil
	}
	return 0, nil, status.Errorf(codes.Unimplemented, "unsupported conversion for %v to INT64", val)
}
func gologoo__convertToFloat64_3fabff857a70ddbd18407d4ff739347f(val interface {
}) (res float64, convertErr error, err error) {
	switch v := val.(type) {
	case int64:
		return float64(v), nil, nil
	case float64:
		return v, nil, nil
	case string:
		res, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, status.Errorf(codes.InvalidArgument, "invalid value for FLOAT64: %q", v), nil
		}
		return res, nil, nil
	}
	return 0, nil, status.Errorf(codes.Unimplemented, "unsupported conversion for %v to FLOAT64", val)
}
func gologoo__convertToString_3fabff857a70ddbd18407d4ff739347f(val interface {
}) (res string, convertErr error, err error) {
	switch v := val.(type) {
	case string:
		return v, nil, nil
	case bool, int64, float64:
		return fmt.Sprintf("%v", v), nil, nil
	case civil.Date:
		return v.String(), nil, nil
	case time.Time:
		return v.UTC().Format(time.RFC3339Nano), nil, nil
	}
	return "", nil, status.Errorf(codes.Unimplemented, "unsupported conversion for %v to STRING", val)
}
func gologoo__convertToBool_3fabff857a70ddbd18407d4ff739347f(val interface {
}) (res bool, convertErr error, err error) {
	switch v := val.(type) {
	case bool:
		return v, nil, nil
	case string:
		res, err := strconv.ParseBool(v)
		if err != nil {
			return false, status.Errorf(codes.InvalidArgument, "invalid value for BOOL: %q", v), nil
		}
		return res, nil, nil
	}
	return false, nil, status.Errorf(codes.Unimplemented, "unsupported conversion for %v to BOOL", val)
}
func gologoo__convertToDate_3fabff857a70ddbd18407d4ff739347f(val interface {
}) (res civil.Date, convertErr error, err error) {
	switch v := val.(type) {
	case civil.Date:
		return v, nil, nil
	case string:
		res, err := civil.ParseDate(v)
		if err != nil {
			return civil.Date{}, status.Errorf(codes.InvalidArgument, "invalid value for DATE: %q", v), nil
		}
		return res, nil, nil
	}
	return civil.Date{}, nil, status.Errorf(codes.Unimplemented, "unsupported conversion for %v to DATE", val)
}
func gologoo__convertToTimestamp_3fabff857a70ddbd18407d4ff739347f(val interface {
}) (res time.Time, convertErr error, err error) {
	switch v := val.(type) {
	case time.Time:
		return v, nil, nil
	case string:
		res, err := time.Parse(time.RFC3339Nano, v)
		if err != nil {
			return time.Time{}, status.Errorf(codes.InvalidArgument, "invalid value for TIMESTAMP: %q", v), nil
		}
		return res, nil, nil
	}
	return time.Time{}, nil, status.Errorf(codes.Unimplemented, "unsupported conversion for %v to TIMESTAMP", val)
}

type aggregateFunc struct {
	AcceptStar bool
	Eval       func(values []interface {
	}, typ spansql.Type) (interface {
	}, spansql.Type, error)
}

var aggregateFuncs = map[string]aggregateFunc{"ANY_VALUE": {Eval: func(values []interface {
}, typ spansql.Type) (interface {
}, spansql.Type, error) {
	for _, v := range values {
		if v != nil {
			return v, typ, nil
		}
	}
	return nil, typ, nil
}}, "ARRAY_AGG": {Eval: func(values []interface {
}, typ spansql.Type) (interface {
}, spansql.Type, error) {
	if typ.Array {
		return nil, spansql.Type{}, fmt.Errorf("ARRAY_AGG unsupported on values of type %v", typ.SQL())
	}
	typ.Array = true
	if len(values) == 0 {
		return nil, typ, nil
	}
	return values, typ, nil
}}, "COUNT": {AcceptStar: true, Eval: func(values []interface {
}, typ spansql.Type) (interface {
}, spansql.Type, error) {
	var n int64
	for _, v := range values {
		if v != nil {
			n++
		}
	}
	return n, int64Type, nil
}}, "MAX": {Eval: func(values []interface {
}, typ spansql.Type) (interface {
}, spansql.Type, error) {
	return evalMinMax("MAX", false, values, typ)
}}, "MIN": {Eval: func(values []interface {
}, typ spansql.Type) (interface {
}, spansql.Type, error) {
	return evalMinMax("MIN", true, values, typ)
}}, "SUM": {Eval: func(values []interface {
}, typ spansql.Type) (interface {
}, spansql.Type, error) {
	if typ.Array || !(typ.Base == spansql.Int64 || typ.Base == spansql.Float64) {
		return nil, spansql.Type{}, fmt.Errorf("SUM only supports arguments of INT64 or FLOAT64 type, not %s", typ.SQL())
	}
	if typ.Base == spansql.Int64 {
		var seen bool
		var sum int64
		for _, v := range values {
			if v == nil {
				continue
			}
			seen = true
			sum += v.(int64)
		}
		if !seen {
			return nil, typ, nil
		}
		return sum, typ, nil
	}
	var seen bool
	var sum float64
	for _, v := range values {
		if v == nil {
			continue
		}
		seen = true
		sum += v.(float64)
	}
	if !seen {
		return nil, typ, nil
	}
	return sum, typ, nil
}}, "AVG": {Eval: func(values []interface {
}, typ spansql.Type) (interface {
}, spansql.Type, error) {
	if typ.Array || !(typ.Base == spansql.Int64 || typ.Base == spansql.Float64) {
		return nil, spansql.Type{}, fmt.Errorf("AVG only supports arguments of INT64 or FLOAT64 type, not %s", typ.SQL())
	}
	if typ.Base == spansql.Int64 {
		var sum int64
		var n float64
		for _, v := range values {
			if v == nil {
				continue
			}
			sum += v.(int64)
			n++
		}
		if n == 0 {
			return nil, typ, nil
		}
		return (float64(sum) / n), float64Type, nil
	}
	var sum float64
	var n float64
	for _, v := range values {
		if v == nil {
			continue
		}
		sum += v.(float64)
		n++
	}
	if n == 0 {
		return nil, typ, nil
	}
	return (sum / n), typ, nil
}}}

func gologoo__evalMinMax_3fabff857a70ddbd18407d4ff739347f(name string, isMin bool, values []interface {
}, typ spansql.Type) (interface {
}, spansql.Type, error) {
	if typ.Array {
		return nil, spansql.Type{}, fmt.Errorf("%s only supports non-array arguments, not %s", name, typ.SQL())
	}
	if len(values) == 0 {
		return nil, typ, nil
	}
	var minMax interface {
	}
	for _, v := range values {
		if v == nil {
			continue
		}
		if typ.Base == spansql.Float64 && math.IsNaN(v.(float64)) {
			return v, typ, nil
		}
		if minMax == nil || compareVals(v, minMax) < 0 == isMin {
			minMax = v
		}
	}
	return minMax, typ, nil
}
func firstErr(errors []error) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__firstErr_3fabff857a70ddbd18407d4ff739347f")
	log.Printf("Input : %v\n", errors)
	r0 := gologoo__firstErr_3fabff857a70ddbd18407d4ff739347f(errors)
	log.Printf("Output: %v\n", r0)
	return r0
}
func cast(values []interface {
}, types []spansql.Type, safe bool) (interface {
}, spansql.Type, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__cast_3fabff857a70ddbd18407d4ff739347f")
	log.Printf("Input : %v %v %v\n", values, types, safe)
	r0, r1, r2 := gologoo__cast_3fabff857a70ddbd18407d4ff739347f(values, types, safe)
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
func convert(val interface {
}, tp spansql.Type) (interface {
}, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__convert_3fabff857a70ddbd18407d4ff739347f")
	log.Printf("Input : %v %v\n", val, tp)
	r0, r1 := gologoo__convert_3fabff857a70ddbd18407d4ff739347f(val, tp)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func convertToInt64(val interface {
}) (res int64, convertErr error, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__convertToInt64_3fabff857a70ddbd18407d4ff739347f")
	log.Printf("Input : %v\n", val)
	res, convertErr, err = gologoo__convertToInt64_3fabff857a70ddbd18407d4ff739347f(val)
	log.Printf("Output: %v %v %v\n", res, convertErr, err)
	return
}
func convertToFloat64(val interface {
}) (res float64, convertErr error, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__convertToFloat64_3fabff857a70ddbd18407d4ff739347f")
	log.Printf("Input : %v\n", val)
	res, convertErr, err = gologoo__convertToFloat64_3fabff857a70ddbd18407d4ff739347f(val)
	log.Printf("Output: %v %v %v\n", res, convertErr, err)
	return
}
func convertToString(val interface {
}) (res string, convertErr error, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__convertToString_3fabff857a70ddbd18407d4ff739347f")
	log.Printf("Input : %v\n", val)
	res, convertErr, err = gologoo__convertToString_3fabff857a70ddbd18407d4ff739347f(val)
	log.Printf("Output: %v %v %v\n", res, convertErr, err)
	return
}
func convertToBool(val interface {
}) (res bool, convertErr error, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__convertToBool_3fabff857a70ddbd18407d4ff739347f")
	log.Printf("Input : %v\n", val)
	res, convertErr, err = gologoo__convertToBool_3fabff857a70ddbd18407d4ff739347f(val)
	log.Printf("Output: %v %v %v\n", res, convertErr, err)
	return
}
func convertToDate(val interface {
}) (res civil.Date, convertErr error, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__convertToDate_3fabff857a70ddbd18407d4ff739347f")
	log.Printf("Input : %v\n", val)
	res, convertErr, err = gologoo__convertToDate_3fabff857a70ddbd18407d4ff739347f(val)
	log.Printf("Output: %v %v %v\n", res, convertErr, err)
	return
}
func convertToTimestamp(val interface {
}) (res time.Time, convertErr error, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__convertToTimestamp_3fabff857a70ddbd18407d4ff739347f")
	log.Printf("Input : %v\n", val)
	res, convertErr, err = gologoo__convertToTimestamp_3fabff857a70ddbd18407d4ff739347f(val)
	log.Printf("Output: %v %v %v\n", res, convertErr, err)
	return
}
func evalMinMax(name string, isMin bool, values []interface {
}, typ spansql.Type) (interface {
}, spansql.Type, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__evalMinMax_3fabff857a70ddbd18407d4ff739347f")
	log.Printf("Input : %v %v %v %v\n", name, isMin, values, typ)
	r0, r1, r2 := gologoo__evalMinMax_3fabff857a70ddbd18407d4ff739347f(name, isMin, values, typ)
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
