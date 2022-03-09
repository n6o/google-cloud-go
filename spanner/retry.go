package spanner

import (
	"context"
	"strings"
	"time"
	"cloud.google.com/go/internal/trace"
	"github.com/golang/protobuf/ptypes"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

const retryInfoKey = "google.rpc.retryinfo-bin"

var DefaultRetryBackoff = gax.Backoff{Initial: 20 * time.Millisecond, Max: 32 * time.Second, Multiplier: 1.3}

type spannerRetryer struct {
	gax.Retryer
}

func gologoo__onCodes_3245d613cc8abcaff6eb34cbe3156069(bo gax.Backoff, cc ...codes.Code) gax.Retryer {
	return &spannerRetryer{Retryer: gax.OnCodes(cc, bo)}
}
func (r *spannerRetryer) gologoo__Retry_3245d613cc8abcaff6eb34cbe3156069(err error) (time.Duration, bool) {
	if status.Code(err) == codes.Internal && !strings.Contains(err.Error(), "stream terminated by RST_STREAM") && !strings.Contains(err.Error(), "HTTP/2 error code: INTERNAL_ERROR") && !strings.Contains(err.Error(), "Connection closed with unknown cause") && !strings.Contains(err.Error(), "Received unexpected EOS on DATA frame from server") {
		return 0, false
	}
	delay, shouldRetry := r.Retryer.Retry(err)
	if !shouldRetry {
		return 0, false
	}
	if serverDelay, hasServerDelay := ExtractRetryDelay(err); hasServerDelay {
		delay = serverDelay
	}
	return delay, true
}
func gologoo__runWithRetryOnAbortedOrSessionNotFound_3245d613cc8abcaff6eb34cbe3156069(ctx context.Context, f func(context.Context) error) error {
	retryer := onCodes(DefaultRetryBackoff, codes.Aborted)
	funcWithRetry := func(ctx context.Context) error {
		for {
			err := f(ctx)
			if err == nil {
				return nil
			}
			var retryErr error
			var se *Error
			if errorAs(err, &se) {
				retryErr = se
			} else {
				_, ok := status.FromError(err)
				if !ok {
					return err
				}
				retryErr = err
			}
			if isSessionNotFoundError(retryErr) {
				trace.TracePrintf(ctx, nil, "Retrying after Session not found")
				continue
			}
			delay, shouldRetry := retryer.Retry(retryErr)
			if !shouldRetry {
				return err
			}
			trace.TracePrintf(ctx, nil, "Backing off after ABORTED for %s, then retrying", delay)
			if err := gax.Sleep(ctx, delay); err != nil {
				return err
			}
		}
	}
	return funcWithRetry(ctx)
}
func gologoo__ExtractRetryDelay_3245d613cc8abcaff6eb34cbe3156069(err error) (time.Duration, bool) {
	var se *Error
	var s *status.Status
	if errorAs(err, &se) {
		s = status.Convert(se.Unwrap())
	} else {
		s = status.Convert(err)
	}
	if s == nil {
		return 0, false
	}
	for _, detail := range s.Details() {
		if retryInfo, ok := detail.(*errdetails.RetryInfo); ok {
			delay, err := ptypes.Duration(retryInfo.RetryDelay)
			if err != nil {
				return 0, false
			}
			return delay, true
		}
	}
	return 0, false
}
func onCodes(bo gax.Backoff, cc ...codes.Code) gax.Retryer {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__onCodes_3245d613cc8abcaff6eb34cbe3156069")
	log.Printf("Input : %v %v\n", bo, cc)
	r0 := gologoo__onCodes_3245d613cc8abcaff6eb34cbe3156069(bo, cc...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (r *spannerRetryer) Retry(err error) (time.Duration, bool) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Retry_3245d613cc8abcaff6eb34cbe3156069")
	log.Printf("Input : %v\n", err)
	r0, r1 := r.gologoo__Retry_3245d613cc8abcaff6eb34cbe3156069(err)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func runWithRetryOnAbortedOrSessionNotFound(ctx context.Context, f func(context.Context) error) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__runWithRetryOnAbortedOrSessionNotFound_3245d613cc8abcaff6eb34cbe3156069")
	log.Printf("Input : %v %v\n", ctx, f)
	r0 := gologoo__runWithRetryOnAbortedOrSessionNotFound_3245d613cc8abcaff6eb34cbe3156069(ctx, f)
	log.Printf("Output: %v\n", r0)
	return r0
}
func ExtractRetryDelay(err error) (time.Duration, bool) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ExtractRetryDelay_3245d613cc8abcaff6eb34cbe3156069")
	log.Printf("Input : %v\n", err)
	r0, r1 := gologoo__ExtractRetryDelay_3245d613cc8abcaff6eb34cbe3156069(err)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
