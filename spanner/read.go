package spanner

import (
	"bytes"
	"context"
	"io"
	"log"
	"sync/atomic"
	"time"

	"cloud.google.com/go/internal/protostruct"
	"cloud.google.com/go/internal/trace"
	"github.com/golang/protobuf/proto"
	proto3 "github.com/golang/protobuf/ptypes/struct"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc/codes"
)

type streamingReceiver interface {
	Recv() (*sppb.PartialResultSet, error)
}

func gologoo__errEarlyReadEnd_7fd4d53bdad5ef404a290eb8dfb05849() error {
	return spannerErrorf(codes.FailedPrecondition, "read completed with active stream")
}
func gologoo__stream_7fd4d53bdad5ef404a290eb8dfb05849(ctx context.Context, logger *log.Logger, rpc func(ct context.Context, resumeToken []byte) (streamingReceiver, error), setTimestamp func(time.Time), release func(error)) *RowIterator {
	return streamWithReplaceSessionFunc(ctx, logger, rpc, nil, setTimestamp, release)
}
func gologoo__streamWithReplaceSessionFunc_7fd4d53bdad5ef404a290eb8dfb05849(ctx context.Context, logger *log.Logger, rpc func(ct context.Context, resumeToken []byte) (streamingReceiver, error), replaceSession func(ctx context.Context) error, setTimestamp func(time.Time), release func(error)) *RowIterator {
	ctx, cancel := context.WithCancel(ctx)
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.RowIterator")
	return &RowIterator{streamd: newResumableStreamDecoder(ctx, logger, rpc, replaceSession), rowd: &partialResultSetDecoder{}, setTimestamp: setTimestamp, release: release, cancel: cancel}
}

type RowIterator struct {
	QueryPlan  *sppb.QueryPlan
	QueryStats map[string]interface {
	}
	RowCount     int64
	Metadata     *sppb.ResultSetMetadata
	streamd      *resumableStreamDecoder
	rowd         *partialResultSetDecoder
	setTimestamp func(time.Time)
	release      func(error)
	cancel       func()
	err          error
	rows         []*Row
	sawStats     bool
}

func (r *RowIterator) gologoo__Next_7fd4d53bdad5ef404a290eb8dfb05849() (*Row, error) {
	if r.err != nil {
		return nil, r.err
	}
	for len(r.rows) == 0 && r.streamd.next() {
		prs := r.streamd.get()
		if prs.Stats != nil {
			r.sawStats = true
			r.QueryPlan = prs.Stats.QueryPlan
			r.QueryStats = protostruct.DecodeToMap(prs.Stats.QueryStats)
			if prs.Stats.RowCount != nil {
				rc, err := extractRowCount(prs.Stats)
				if err != nil {
					return nil, err
				}
				r.RowCount = rc
			}
		}
		var metadata *sppb.ResultSetMetadata
		r.rows, metadata, r.err = r.rowd.add(prs)
		if metadata != nil {
			r.Metadata = metadata
		}
		if r.err != nil {
			return nil, r.err
		}
		if !r.rowd.ts.IsZero() && r.setTimestamp != nil {
			r.setTimestamp(r.rowd.ts)
			r.setTimestamp = nil
		}
	}
	if len(r.rows) > 0 {
		row := r.rows[0]
		r.rows = r.rows[1:]
		return row, nil
	}
	if err := r.streamd.lastErr(); err != nil {
		r.err = ToSpannerError(err)
	} else if !r.rowd.done() {
		r.err = errEarlyReadEnd()
	} else {
		r.err = iterator.Done
	}
	return nil, r.err
}
func gologoo__extractRowCount_7fd4d53bdad5ef404a290eb8dfb05849(stats *sppb.ResultSetStats) (int64, error) {
	if stats.RowCount == nil {
		return 0, spannerErrorf(codes.Internal, "missing RowCount")
	}
	switch rc := stats.RowCount.(type) {
	case *sppb.ResultSetStats_RowCountExact:
		return rc.RowCountExact, nil
	case *sppb.ResultSetStats_RowCountLowerBound:
		return rc.RowCountLowerBound, nil
	default:
		return 0, spannerErrorf(codes.Internal, "unknown RowCount type %T", stats.RowCount)
	}
}
func (r *RowIterator) gologoo__Do_7fd4d53bdad5ef404a290eb8dfb05849(f func(r *Row) error) error {
	defer r.Stop()
	for {
		row, err := r.Next()
		switch err {
		case iterator.Done:
			return nil
		case nil:
			if err = f(row); err != nil {
				return err
			}
		default:
			return err
		}
	}
}
func (r *RowIterator) gologoo__Stop_7fd4d53bdad5ef404a290eb8dfb05849() {
	if r.streamd != nil {
		if r.err != nil && r.err != iterator.Done {
			defer trace.EndSpan(r.streamd.ctx, r.err)
		} else {
			defer trace.EndSpan(r.streamd.ctx, nil)
		}
	}
	if r.cancel != nil {
		r.cancel()
	}
	if r.release != nil {
		r.release(r.err)
		if r.err == nil {
			r.err = spannerErrorf(codes.FailedPrecondition, "Next called after Stop")
		}
		r.release = nil
	}
}

type partialResultQueue struct {
	q     []*sppb.PartialResultSet
	first int
	last  int
	n     int
}

func (q *partialResultQueue) gologoo__empty_7fd4d53bdad5ef404a290eb8dfb05849() bool {
	return q.n == 0
}
func gologoo__errEmptyQueue_7fd4d53bdad5ef404a290eb8dfb05849() error {
	return spannerErrorf(codes.OutOfRange, "empty partialResultQueue")
}
func (q *partialResultQueue) gologoo__peekLast_7fd4d53bdad5ef404a290eb8dfb05849() (*sppb.PartialResultSet, error) {
	if q.empty() {
		return nil, errEmptyQueue()
	}
	return q.q[(q.last+cap(q.q)-1)%cap(q.q)], nil
}
func (q *partialResultQueue) gologoo__push_7fd4d53bdad5ef404a290eb8dfb05849(r *sppb.PartialResultSet) {
	if q.q == nil {
		q.q = make([]*sppb.PartialResultSet, 8)
	}
	if q.n == cap(q.q) {
		buf := make([]*sppb.PartialResultSet, cap(q.q)*2)
		for i := 0; i < q.n; i++ {
			buf[i] = q.q[(q.first+i)%cap(q.q)]
		}
		q.q = buf
		q.first = 0
		q.last = q.n
	}
	q.q[q.last] = r
	q.last = (q.last + 1) % cap(q.q)
	q.n++
}
func (q *partialResultQueue) gologoo__pop_7fd4d53bdad5ef404a290eb8dfb05849() *sppb.PartialResultSet {
	if q.n == 0 {
		return nil
	}
	r := q.q[q.first]
	q.q[q.first] = nil
	q.first = (q.first + 1) % cap(q.q)
	q.n--
	return r
}
func (q *partialResultQueue) gologoo__clear_7fd4d53bdad5ef404a290eb8dfb05849() {
	*q = partialResultQueue{}
}
func (q *partialResultQueue) gologoo__dump_7fd4d53bdad5ef404a290eb8dfb05849() []*sppb.PartialResultSet {
	var dq []*sppb.PartialResultSet
	for i := q.first; len(dq) < q.n; i = (i + 1) % cap(q.q) {
		dq = append(dq, q.q[i])
	}
	return dq
}

type resumableStreamDecoderState int

const (
	unConnected resumableStreamDecoderState = iota
	queueingRetryable
	queueingUnretryable
	aborted
	finished
)

type resumableStreamDecoder struct {
	state                       resumableStreamDecoderState
	stateWitness                func(resumableStreamDecoderState)
	ctx                         context.Context
	rpc                         func(ctx context.Context, restartToken []byte) (streamingReceiver, error)
	replaceSessionFunc          func(ctx context.Context) error
	logger                      *log.Logger
	stream                      streamingReceiver
	q                           partialResultQueue
	bytesBetweenResumeTokens    int32
	maxBytesBetweenResumeTokens int32
	np                          *sppb.PartialResultSet
	resumeToken                 []byte
	err                         error
	backoff                     gax.Backoff
}

func gologoo__newResumableStreamDecoder_7fd4d53bdad5ef404a290eb8dfb05849(ctx context.Context, logger *log.Logger, rpc func(ct context.Context, restartToken []byte) (streamingReceiver, error), replaceSession func(ctx context.Context) error) *resumableStreamDecoder {
	return &resumableStreamDecoder{ctx: ctx, logger: logger, rpc: rpc, replaceSessionFunc: replaceSession, maxBytesBetweenResumeTokens: atomic.LoadInt32(&maxBytesBetweenResumeTokens), backoff: DefaultRetryBackoff}
}
func (d *resumableStreamDecoder) gologoo__changeState_7fd4d53bdad5ef404a290eb8dfb05849(target resumableStreamDecoderState) {
	if d.state == queueingRetryable && d.state != target {
		d.bytesBetweenResumeTokens = 0
	}
	d.state = target
	if d.stateWitness != nil {
		d.stateWitness(target)
	}
}
func (d *resumableStreamDecoder) gologoo__isNewResumeToken_7fd4d53bdad5ef404a290eb8dfb05849(rt []byte) bool {
	if rt == nil {
		return false
	}
	if bytes.Equal(rt, d.resumeToken) {
		return false
	}
	return true
}

var maxBytesBetweenResumeTokens = int32(128 * 1024 * 1024)

func (d *resumableStreamDecoder) gologoo__next_7fd4d53bdad5ef404a290eb8dfb05849() bool {
	retryer := onCodes(d.backoff, codes.Unavailable, codes.Internal)
	for {
		switch d.state {
		case unConnected:
			d.stream, d.err = d.rpc(d.ctx, d.resumeToken)
			if d.err == nil {
				d.changeState(queueingRetryable)
				continue
			}
			delay, shouldRetry := retryer.Retry(d.err)
			if !shouldRetry {
				d.changeState(aborted)
				continue
			}
			trace.TracePrintf(d.ctx, nil, "Backing off stream read for %s", delay)
			if err := gax.Sleep(d.ctx, delay); err == nil {
				d.changeState(unConnected)
			} else {
				d.err = err
				d.changeState(aborted)
			}
			continue
		case queueingRetryable:
			fallthrough
		case queueingUnretryable:
			last, err := d.q.peekLast()
			if err != nil {
				d.tryRecv(retryer)
				continue
			}
			if d.isNewResumeToken(last.ResumeToken) {
				d.np = d.q.pop()
				if d.q.empty() {
					d.bytesBetweenResumeTokens = 0
					d.resumeToken = d.np.ResumeToken
					d.changeState(queueingRetryable)
				}
				return true
			}
			if d.bytesBetweenResumeTokens >= d.maxBytesBetweenResumeTokens && d.state == queueingRetryable {
				d.changeState(queueingUnretryable)
				continue
			}
			if d.state == queueingUnretryable {
				d.np = d.q.pop()
				return true
			}
			d.tryRecv(retryer)
			continue
		case aborted:
			d.q.clear()
			return false
		case finished:
			if d.q.empty() {
				return false
			}
			d.np = d.q.pop()
			return true
		default:
			logf(d.logger, "Unexpected resumableStreamDecoder.state: %v", d.state)
			return false
		}
	}
}
func (d *resumableStreamDecoder) gologoo__tryRecv_7fd4d53bdad5ef404a290eb8dfb05849(retryer gax.Retryer) {
	var res *sppb.PartialResultSet
	res, d.err = d.stream.Recv()
	if d.err == nil {
		d.q.push(res)
		if d.state == queueingRetryable && !d.isNewResumeToken(res.ResumeToken) {
			d.bytesBetweenResumeTokens += int32(proto.Size(res))
		}
		d.changeState(d.state)
		return
	}
	if d.err == io.EOF {
		d.err = nil
		d.changeState(finished)
		return
	}
	if d.replaceSessionFunc != nil && isSessionNotFoundError(d.err) && d.resumeToken == nil {
		if err := d.replaceSessionFunc(d.ctx); err != nil {
			d.err = err
			d.changeState(aborted)
			return
		}
	} else {
		delay, shouldRetry := retryer.Retry(d.err)
		if !shouldRetry || d.state != queueingRetryable {
			d.changeState(aborted)
			return
		}
		if err := gax.Sleep(d.ctx, delay); err != nil {
			d.err = err
			d.changeState(aborted)
			return
		}
	}
	d.err = nil
	d.q.clear()
	d.stream = nil
	d.changeState(unConnected)
}
func (d *resumableStreamDecoder) gologoo__get_7fd4d53bdad5ef404a290eb8dfb05849() *sppb.PartialResultSet {
	return d.np
}
func (d *resumableStreamDecoder) gologoo__lastErr_7fd4d53bdad5ef404a290eb8dfb05849() error {
	return d.err
}

type partialResultSetDecoder struct {
	row     Row
	tx      *sppb.Transaction
	chunked bool
	ts      time.Time
}

func (p *partialResultSetDecoder) gologoo__yield_7fd4d53bdad5ef404a290eb8dfb05849(chunked, last bool) *Row {
	if len(p.row.vals) == len(p.row.fields) && (!chunked || !last) {
		fresh := Row{fields: p.row.fields, vals: make([]*proto3.Value, len(p.row.vals))}
		copy(fresh.vals, p.row.vals)
		p.row.vals = p.row.vals[:0]
		return &fresh
	}
	return nil
}
func gologoo__errChunkedEmptyRow_7fd4d53bdad5ef404a290eb8dfb05849() error {
	return spannerErrorf(codes.FailedPrecondition, "got invalid chunked PartialResultSet with empty Row")
}
func (p *partialResultSetDecoder) gologoo__add_7fd4d53bdad5ef404a290eb8dfb05849(r *sppb.PartialResultSet) ([]*Row, *sppb.ResultSetMetadata, error) {
	var rows []*Row
	if r.Metadata != nil {
		if p.row.fields == nil {
			p.row.fields = r.Metadata.RowType.Fields
		}
		if p.tx == nil && r.Metadata.Transaction != nil {
			p.tx = r.Metadata.Transaction
			if p.tx.ReadTimestamp != nil {
				p.ts = time.Unix(p.tx.ReadTimestamp.Seconds, int64(p.tx.ReadTimestamp.Nanos))
			}
		}
	}
	if len(r.Values) == 0 {
		return nil, r.Metadata, nil
	}
	if p.chunked {
		p.chunked = false
		last := len(p.row.vals) - 1
		if last < 0 {
			return nil, nil, errChunkedEmptyRow()
		}
		var err error
		if p.row.vals[last], err = p.merge(p.row.vals[last], r.Values[0]); err != nil {
			return nil, r.Metadata, err
		}
		r.Values = r.Values[1:]
		if row := p.yield(r.ChunkedValue, len(r.Values) == 0); row != nil {
			rows = append(rows, row)
		}
	}
	for i, v := range r.Values {
		p.row.vals = append(p.row.vals, v)
		if row := p.yield(r.ChunkedValue, i == len(r.Values)-1); row != nil {
			rows = append(rows, row)
		}
	}
	if r.ChunkedValue {
		p.chunked = true
	}
	return rows, r.Metadata, nil
}
func (p *partialResultSetDecoder) gologoo__isMergeable_7fd4d53bdad5ef404a290eb8dfb05849(a *proto3.Value) bool {
	switch a.Kind.(type) {
	case *proto3.Value_StringValue:
		return true
	case *proto3.Value_ListValue:
		return true
	default:
		return false
	}
}
func gologoo__errIncompatibleMergeTypes_7fd4d53bdad5ef404a290eb8dfb05849(a, b *proto3.Value) error {
	return spannerErrorf(codes.FailedPrecondition, "incompatible type in chunked PartialResultSet. expected (%T), got (%T)", a.Kind, b.Kind)
}
func gologoo__errUnsupportedMergeType_7fd4d53bdad5ef404a290eb8dfb05849(a *proto3.Value) error {
	return spannerErrorf(codes.FailedPrecondition, "unsupported type merge (%T)", a.Kind)
}
func (p *partialResultSetDecoder) gologoo__merge_7fd4d53bdad5ef404a290eb8dfb05849(a, b *proto3.Value) (*proto3.Value, error) {
	var err error
	typeErr := errIncompatibleMergeTypes(a, b)
	switch t := a.Kind.(type) {
	case *proto3.Value_StringValue:
		s, ok := b.Kind.(*proto3.Value_StringValue)
		if !ok {
			return nil, typeErr
		}
		return &proto3.Value{Kind: &proto3.Value_StringValue{StringValue: t.StringValue + s.StringValue}}, nil
	case *proto3.Value_ListValue:
		l, ok := b.Kind.(*proto3.Value_ListValue)
		if !ok {
			return nil, typeErr
		}
		if l.ListValue == nil || len(l.ListValue.Values) <= 0 {
			return a, nil
		}
		if t.ListValue == nil || len(t.ListValue.Values) <= 0 {
			return b, nil
		}
		if la := len(t.ListValue.Values) - 1; p.isMergeable(t.ListValue.Values[la]) {
			t.ListValue.Values[la], err = p.merge(t.ListValue.Values[la], l.ListValue.Values[0])
			if err != nil {
				return nil, err
			}
			l.ListValue.Values = l.ListValue.Values[1:]
		}
		return &proto3.Value{Kind: &proto3.Value_ListValue{ListValue: &proto3.ListValue{Values: append(t.ListValue.Values, l.ListValue.Values...)}}}, nil
	default:
		return nil, errUnsupportedMergeType(a)
	}
}
func (p *partialResultSetDecoder) gologoo__done_7fd4d53bdad5ef404a290eb8dfb05849() bool {
	return len(p.row.vals) == 0 && !p.chunked
}
func errEarlyReadEnd() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errEarlyReadEnd_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : (none)\n")
	r0 := gologoo__errEarlyReadEnd_7fd4d53bdad5ef404a290eb8dfb05849()
	log.Printf("Output: %v\n", r0)
	return r0
}
func stream(ctx context.Context, logger *log.Logger, rpc func(ct context.Context, resumeToken []byte) (streamingReceiver, error), setTimestamp func(time.Time), release func(error)) *RowIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__stream_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : %v %v %v %v %v\n", ctx, logger, rpc, setTimestamp, release)
	r0 := gologoo__stream_7fd4d53bdad5ef404a290eb8dfb05849(ctx, logger, rpc, setTimestamp, release)
	log.Printf("Output: %v\n", r0)
	return r0
}
func streamWithReplaceSessionFunc(ctx context.Context, logger *log.Logger, rpc func(ct context.Context, resumeToken []byte) (streamingReceiver, error), replaceSession func(ctx context.Context) error, setTimestamp func(time.Time), release func(error)) *RowIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__streamWithReplaceSessionFunc_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : %v %v %v %v %v %v\n", ctx, logger, rpc, replaceSession, setTimestamp, release)
	r0 := gologoo__streamWithReplaceSessionFunc_7fd4d53bdad5ef404a290eb8dfb05849(ctx, logger, rpc, replaceSession, setTimestamp, release)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (r *RowIterator) Next() (*Row, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Next_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : (none)\n")
	r0, r1 := r.gologoo__Next_7fd4d53bdad5ef404a290eb8dfb05849()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func extractRowCount(stats *sppb.ResultSetStats) (int64, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__extractRowCount_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : %v\n", stats)
	r0, r1 := gologoo__extractRowCount_7fd4d53bdad5ef404a290eb8dfb05849(stats)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (r *RowIterator) Do(f func(r *Row) error) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Do_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : %v\n", f)
	r0 := r.gologoo__Do_7fd4d53bdad5ef404a290eb8dfb05849(f)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (r *RowIterator) Stop() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Stop_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : (none)\n")
	r.gologoo__Stop_7fd4d53bdad5ef404a290eb8dfb05849()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (q *partialResultQueue) empty() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__empty_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : (none)\n")
	r0 := q.gologoo__empty_7fd4d53bdad5ef404a290eb8dfb05849()
	log.Printf("Output: %v\n", r0)
	return r0
}
func errEmptyQueue() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errEmptyQueue_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : (none)\n")
	r0 := gologoo__errEmptyQueue_7fd4d53bdad5ef404a290eb8dfb05849()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (q *partialResultQueue) peekLast() (*sppb.PartialResultSet, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__peekLast_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : (none)\n")
	r0, r1 := q.gologoo__peekLast_7fd4d53bdad5ef404a290eb8dfb05849()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (q *partialResultQueue) push(r *sppb.PartialResultSet) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__push_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : %v\n", r)
	q.gologoo__push_7fd4d53bdad5ef404a290eb8dfb05849(r)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (q *partialResultQueue) pop() *sppb.PartialResultSet {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__pop_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : (none)\n")
	r0 := q.gologoo__pop_7fd4d53bdad5ef404a290eb8dfb05849()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (q *partialResultQueue) clear() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__clear_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : (none)\n")
	q.gologoo__clear_7fd4d53bdad5ef404a290eb8dfb05849()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (q *partialResultQueue) dump() []*sppb.PartialResultSet {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__dump_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : (none)\n")
	r0 := q.gologoo__dump_7fd4d53bdad5ef404a290eb8dfb05849()
	log.Printf("Output: %v\n", r0)
	return r0
}
func newResumableStreamDecoder(ctx context.Context, logger *log.Logger, rpc func(ct context.Context, restartToken []byte) (streamingReceiver, error), replaceSession func(ctx context.Context) error) *resumableStreamDecoder {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__newResumableStreamDecoder_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : %v %v %v %v\n", ctx, logger, rpc, replaceSession)
	r0 := gologoo__newResumableStreamDecoder_7fd4d53bdad5ef404a290eb8dfb05849(ctx, logger, rpc, replaceSession)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (d *resumableStreamDecoder) changeState(target resumableStreamDecoderState) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__changeState_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : %v\n", target)
	d.gologoo__changeState_7fd4d53bdad5ef404a290eb8dfb05849(target)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (d *resumableStreamDecoder) isNewResumeToken(rt []byte) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isNewResumeToken_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : %v\n", rt)
	r0 := d.gologoo__isNewResumeToken_7fd4d53bdad5ef404a290eb8dfb05849(rt)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (d *resumableStreamDecoder) next() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__next_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : (none)\n")
	r0 := d.gologoo__next_7fd4d53bdad5ef404a290eb8dfb05849()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (d *resumableStreamDecoder) tryRecv(retryer gax.Retryer) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__tryRecv_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : %v\n", retryer)
	d.gologoo__tryRecv_7fd4d53bdad5ef404a290eb8dfb05849(retryer)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (d *resumableStreamDecoder) get() *sppb.PartialResultSet {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__get_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : (none)\n")
	r0 := d.gologoo__get_7fd4d53bdad5ef404a290eb8dfb05849()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (d *resumableStreamDecoder) lastErr() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__lastErr_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : (none)\n")
	r0 := d.gologoo__lastErr_7fd4d53bdad5ef404a290eb8dfb05849()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *partialResultSetDecoder) yield(chunked, last bool) *Row {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__yield_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : %v %v\n", chunked, last)
	r0 := p.gologoo__yield_7fd4d53bdad5ef404a290eb8dfb05849(chunked, last)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errChunkedEmptyRow() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errChunkedEmptyRow_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : (none)\n")
	r0 := gologoo__errChunkedEmptyRow_7fd4d53bdad5ef404a290eb8dfb05849()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *partialResultSetDecoder) add(r *sppb.PartialResultSet) ([]*Row, *sppb.ResultSetMetadata, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__add_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : %v\n", r)
	r0, r1, r2 := p.gologoo__add_7fd4d53bdad5ef404a290eb8dfb05849(r)
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
func (p *partialResultSetDecoder) isMergeable(a *proto3.Value) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isMergeable_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : %v\n", a)
	r0 := p.gologoo__isMergeable_7fd4d53bdad5ef404a290eb8dfb05849(a)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errIncompatibleMergeTypes(a, b *proto3.Value) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errIncompatibleMergeTypes_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : %v %v\n", a, b)
	r0 := gologoo__errIncompatibleMergeTypes_7fd4d53bdad5ef404a290eb8dfb05849(a, b)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errUnsupportedMergeType(a *proto3.Value) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errUnsupportedMergeType_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : %v\n", a)
	r0 := gologoo__errUnsupportedMergeType_7fd4d53bdad5ef404a290eb8dfb05849(a)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *partialResultSetDecoder) merge(a, b *proto3.Value) (*proto3.Value, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__merge_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : %v %v\n", a, b)
	r0, r1 := p.gologoo__merge_7fd4d53bdad5ef404a290eb8dfb05849(a, b)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *partialResultSetDecoder) done() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__done_7fd4d53bdad5ef404a290eb8dfb05849")
	log.Printf("Input : (none)\n")
	r0 := p.gologoo__done_7fd4d53bdad5ef404a290eb8dfb05849()
	log.Printf("Output: %v\n", r0)
	return r0
}
