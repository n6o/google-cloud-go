package spanner

import (
	"context"
	"fmt"
	"github.com/googleapis/gax-go/v2/apierror"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type Error struct {
	Code                  codes.Code
	err                   error
	Desc                  string
	additionalInformation string
}
type TransactionOutcomeUnknownError struct {
	err error
}

const transactionOutcomeUnknownMsg = "transaction outcome unknown"

func (*TransactionOutcomeUnknownError) gologoo__Error_c64688aea675161795c04636e24c08f8() string {
	return transactionOutcomeUnknownMsg
}
func (e *TransactionOutcomeUnknownError) gologoo__Unwrap_c64688aea675161795c04636e24c08f8() error {
	return e.err
}
func (e *Error) gologoo__Error_c64688aea675161795c04636e24c08f8() string {
	if e == nil {
		return fmt.Sprintf("spanner: OK")
	}
	code := ErrCode(e)
	if e.additionalInformation == "" {
		return fmt.Sprintf("spanner: code = %q, desc = %q", code, e.Desc)
	}
	return fmt.Sprintf("spanner: code = %q, desc = %q, additional information = %s", code, e.Desc, e.additionalInformation)
}
func (e *Error) gologoo__Unwrap_c64688aea675161795c04636e24c08f8() error {
	return e.err
}
func (e *Error) gologoo__GRPCStatus_c64688aea675161795c04636e24c08f8() *status.Status {
	err := unwrap(e)
	for {
		if err == nil {
			return status.New(e.Code, e.Desc)
		}
		code := status.Code(err)
		if code != codes.Unknown {
			return status.New(code, e.Desc)
		}
		err = unwrap(err)
	}
}
func (e *Error) gologoo__decorate_c64688aea675161795c04636e24c08f8(info string) {
	e.Desc = fmt.Sprintf("%v, %v", info, e.Desc)
}
func gologoo__spannerErrorf_c64688aea675161795c04636e24c08f8(code codes.Code, format string, args ...interface {
}) error {
	msg := fmt.Sprintf(format, args...)
	wrapped, _ := apierror.FromError(status.Error(code, msg))
	return &Error{Code: code, err: wrapped, Desc: msg}
}
func gologoo__ToSpannerError_c64688aea675161795c04636e24c08f8(err error) error {
	return toSpannerErrorWithCommitInfo(err, false)
}
func gologoo__toAPIError_c64688aea675161795c04636e24c08f8(err error) error {
	if apiError, ok := apierror.FromError(err); ok {
		return apiError
	}
	return err
}
func gologoo__toSpannerErrorWithCommitInfo_c64688aea675161795c04636e24c08f8(err error, errorDuringCommit bool) error {
	if err == nil {
		return nil
	}
	var se *Error
	if errorAs(err, &se) {
		return se
	}
	switch {
	case err == context.DeadlineExceeded || err == context.Canceled:
		desc := err.Error()
		wrapped := status.FromContextError(err).Err()
		if errorDuringCommit {
			desc = fmt.Sprintf("%s, %s", desc, transactionOutcomeUnknownMsg)
			wrapped = &TransactionOutcomeUnknownError{err: wrapped}
		}
		return &Error{status.FromContextError(err).Code(), toAPIError(wrapped), desc, ""}
	case status.Code(err) == codes.Unknown:
		return &Error{codes.Unknown, toAPIError(err), err.Error(), ""}
	default:
		statusErr := status.Convert(err)
		code, desc := statusErr.Code(), statusErr.Message()
		wrapped := err
		if errorDuringCommit && (code == codes.DeadlineExceeded || code == codes.Canceled) {
			desc = fmt.Sprintf("%s, %s", desc, transactionOutcomeUnknownMsg)
			wrapped = &TransactionOutcomeUnknownError{err: wrapped}
		}
		return &Error{code, toAPIError(wrapped), desc, ""}
	}
}
func gologoo__ErrCode_c64688aea675161795c04636e24c08f8(err error) codes.Code {
	s, ok := status.FromError(err)
	if !ok {
		return codes.Unknown
	}
	return s.Code()
}
func gologoo__ErrDesc_c64688aea675161795c04636e24c08f8(err error) string {
	var se *Error
	if !errorAs(err, &se) {
		return err.Error()
	}
	return se.Desc
}
func gologoo__extractResourceType_c64688aea675161795c04636e24c08f8(err error) (string, bool) {
	var s *status.Status
	var se *Error
	if errorAs(err, &se) {
		s = status.Convert(se.Unwrap())
	} else {
		s = status.Convert(err)
	}
	if s == nil {
		return "", false
	}
	for _, detail := range s.Details() {
		if resourceInfo, ok := detail.(*errdetails.ResourceInfo); ok {
			return resourceInfo.ResourceType, true
		}
	}
	return "", false
}
func (recv *TransactionOutcomeUnknownError) Error() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Error_c64688aea675161795c04636e24c08f8")
	log.Printf("Input : (none)\n")
	r0 := recv.gologoo__Error_c64688aea675161795c04636e24c08f8()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (e *TransactionOutcomeUnknownError) Unwrap() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Unwrap_c64688aea675161795c04636e24c08f8")
	log.Printf("Input : (none)\n")
	r0 := e.gologoo__Unwrap_c64688aea675161795c04636e24c08f8()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (e *Error) Error() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Error_c64688aea675161795c04636e24c08f8")
	log.Printf("Input : (none)\n")
	r0 := e.gologoo__Error_c64688aea675161795c04636e24c08f8()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (e *Error) Unwrap() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Unwrap_c64688aea675161795c04636e24c08f8")
	log.Printf("Input : (none)\n")
	r0 := e.gologoo__Unwrap_c64688aea675161795c04636e24c08f8()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (e *Error) GRPCStatus() *status.Status {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GRPCStatus_c64688aea675161795c04636e24c08f8")
	log.Printf("Input : (none)\n")
	r0 := e.gologoo__GRPCStatus_c64688aea675161795c04636e24c08f8()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (e *Error) decorate(info string) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decorate_c64688aea675161795c04636e24c08f8")
	log.Printf("Input : %v\n", info)
	e.gologoo__decorate_c64688aea675161795c04636e24c08f8(info)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func spannerErrorf(code codes.Code, format string, args ...interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__spannerErrorf_c64688aea675161795c04636e24c08f8")
	log.Printf("Input : %v %v %v\n", code, format, args)
	r0 := gologoo__spannerErrorf_c64688aea675161795c04636e24c08f8(code, format, args...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func ToSpannerError(err error) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ToSpannerError_c64688aea675161795c04636e24c08f8")
	log.Printf("Input : %v\n", err)
	r0 := gologoo__ToSpannerError_c64688aea675161795c04636e24c08f8(err)
	log.Printf("Output: %v\n", r0)
	return r0
}
func toAPIError(err error) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__toAPIError_c64688aea675161795c04636e24c08f8")
	log.Printf("Input : %v\n", err)
	r0 := gologoo__toAPIError_c64688aea675161795c04636e24c08f8(err)
	log.Printf("Output: %v\n", r0)
	return r0
}
func toSpannerErrorWithCommitInfo(err error, errorDuringCommit bool) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__toSpannerErrorWithCommitInfo_c64688aea675161795c04636e24c08f8")
	log.Printf("Input : %v %v\n", err, errorDuringCommit)
	r0 := gologoo__toSpannerErrorWithCommitInfo_c64688aea675161795c04636e24c08f8(err, errorDuringCommit)
	log.Printf("Output: %v\n", r0)
	return r0
}
func ErrCode(err error) codes.Code {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ErrCode_c64688aea675161795c04636e24c08f8")
	log.Printf("Input : %v\n", err)
	r0 := gologoo__ErrCode_c64688aea675161795c04636e24c08f8(err)
	log.Printf("Output: %v\n", r0)
	return r0
}
func ErrDesc(err error) string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ErrDesc_c64688aea675161795c04636e24c08f8")
	log.Printf("Input : %v\n", err)
	r0 := gologoo__ErrDesc_c64688aea675161795c04636e24c08f8(err)
	log.Printf("Output: %v\n", r0)
	return r0
}
func extractResourceType(err error) (string, bool) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__extractResourceType_c64688aea675161795c04636e24c08f8")
	log.Printf("Input : %v\n", err)
	r0, r1 := gologoo__extractResourceType_c64688aea675161795c04636e24c08f8(err)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
