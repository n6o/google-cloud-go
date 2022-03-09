package spanner

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
	"cloud.google.com/go/internal/trace"
	vkit "cloud.google.com/go/spanner/apiv1"
	"github.com/golang/protobuf/proto"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
)

type transactionID []byte
type txReadEnv interface {
	acquire(ctx context.Context) (*sessionHandle, *sppb.TransactionSelector, error)
	setTimestamp(time.Time)
	release(error)
}
type txReadOnly struct {
	txReadEnv
	sequenceNumber     int64
	replaceSessionFunc func(ctx context.Context) error
	sp                 *sessionPool
	sh                 *sessionHandle
	qo                 QueryOptions
	txOpts             TransactionOptions
	ct                 *commonTags
}
type TransactionOptions struct {
	CommitOptions  CommitOptions
	TransactionTag string
	CommitPriority sppb.RequestOptions_Priority
}

func (to *TransactionOptions) gologoo__requestPriority_930eb0d2556f62a47e9151ba4c271fab() sppb.RequestOptions_Priority {
	return to.CommitPriority
}
func (to *TransactionOptions) gologoo__requestTag_930eb0d2556f62a47e9151ba4c271fab() string {
	return ""
}
func gologoo__errSessionClosed_930eb0d2556f62a47e9151ba4c271fab(sh *sessionHandle) error {
	return spannerErrorf(codes.FailedPrecondition, "session is already recycled / destroyed: session_id = %q, rpc_client = %v", sh.getID(), sh.getClient())
}
func (t *txReadOnly) gologoo__Read_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, table string, keys KeySet, columns []string) *RowIterator {
	return t.ReadWithOptions(ctx, table, keys, columns, nil)
}
func (t *txReadOnly) gologoo__ReadUsingIndex_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, table, index string, keys KeySet, columns []string) (ri *RowIterator) {
	return t.ReadWithOptions(ctx, table, keys, columns, &ReadOptions{Index: index})
}

type ReadOptions struct {
	Index      string
	Limit      int
	Priority   sppb.RequestOptions_Priority
	RequestTag string
}

func (t *txReadOnly) gologoo__ReadWithOptions_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, table string, keys KeySet, columns []string, opts *ReadOptions) (ri *RowIterator) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.Read")
	defer func() {
		trace.EndSpan(ctx, ri.err)
	}()
	var (
		sh  *sessionHandle
		ts  *sppb.TransactionSelector
		err error
	)
	kset, err := keys.keySetProto()
	if err != nil {
		return &RowIterator{err: err}
	}
	if sh, ts, err = t.acquire(ctx); err != nil {
		return &RowIterator{err: err}
	}
	client := sh.getClient()
	if client == nil {
		return &RowIterator{err: errSessionClosed(sh)}
	}
	index := ""
	limit := 0
	prio := sppb.RequestOptions_PRIORITY_UNSPECIFIED
	requestTag := ""
	if opts != nil {
		index = opts.Index
		if opts.Limit > 0 {
			limit = opts.Limit
		}
		prio = opts.Priority
		requestTag = opts.RequestTag
	}
	return streamWithReplaceSessionFunc(contextWithOutgoingMetadata(ctx, sh.getMetadata()), sh.session.logger, func(ctx context.Context, resumeToken []byte) (streamingReceiver, error) {
		client, err := client.StreamingRead(ctx, &sppb.ReadRequest{Session: t.sh.getID(), Transaction: ts, Table: table, Index: index, Columns: columns, KeySet: kset, ResumeToken: resumeToken, Limit: int64(limit), RequestOptions: createRequestOptions(prio, requestTag, t.txOpts.TransactionTag)})
		if err != nil {
			return client, err
		}
		md, err := client.Header()
		if getGFELatencyMetricsFlag() && md != nil && t.ct != nil {
			if err := createContextAndCaptureGFELatencyMetrics(ctx, t.ct, md, "ReadWithOptions"); err != nil {
				trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", err)
			}
		}
		return client, err
	}, t.replaceSessionFunc, t.setTimestamp, t.release)
}
func gologoo__errRowNotFound_930eb0d2556f62a47e9151ba4c271fab(table string, key Key) error {
	return spannerErrorf(codes.NotFound, "row not found(Table: %v, PrimaryKey: %v)", table, key)
}
func gologoo__errRowNotFoundByIndex_930eb0d2556f62a47e9151ba4c271fab(table string, key Key, index string) error {
	return spannerErrorf(codes.NotFound, "row not found(Table: %v, IndexKey: %v, Index: %v)", table, key, index)
}
func gologoo__errMultipleRowsFound_930eb0d2556f62a47e9151ba4c271fab(table string, key Key, index string) error {
	return spannerErrorf(codes.FailedPrecondition, "more than one row found by index(Table: %v, IndexKey: %v, Index: %v)", table, key, index)
}
func (t *txReadOnly) gologoo__ReadRow_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, table string, key Key, columns []string) (*Row, error) {
	return t.ReadRowWithOptions(ctx, table, key, columns, nil)
}
func (t *txReadOnly) gologoo__ReadRowWithOptions_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, table string, key Key, columns []string, opts *ReadOptions) (*Row, error) {
	iter := t.ReadWithOptions(ctx, table, key, columns, opts)
	defer iter.Stop()
	row, err := iter.Next()
	switch err {
	case iterator.Done:
		return nil, errRowNotFound(table, key)
	case nil:
		return row, nil
	default:
		return nil, err
	}
}
func (t *txReadOnly) gologoo__ReadRowUsingIndex_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, table string, index string, key Key, columns []string) (*Row, error) {
	iter := t.ReadUsingIndex(ctx, table, index, key, columns)
	defer iter.Stop()
	row, err := iter.Next()
	switch err {
	case iterator.Done:
		return nil, errRowNotFoundByIndex(table, key, index)
	case nil:
		_, err := iter.Next()
		switch err {
		case iterator.Done:
			return row, nil
		case nil:
			return nil, errMultipleRowsFound(table, key, index)
		default:
			return nil, err
		}
	default:
		return nil, err
	}
}

type QueryOptions struct {
	Mode       *sppb.ExecuteSqlRequest_QueryMode
	Options    *sppb.ExecuteSqlRequest_QueryOptions
	Priority   sppb.RequestOptions_Priority
	RequestTag string
}

func (qo QueryOptions) gologoo__merge_930eb0d2556f62a47e9151ba4c271fab(opts QueryOptions) QueryOptions {
	merged := QueryOptions{Mode: qo.Mode, Options: &sppb.ExecuteSqlRequest_QueryOptions{}, RequestTag: qo.RequestTag, Priority: qo.Priority}
	if opts.Mode != nil {
		merged.Mode = opts.Mode
	}
	if opts.RequestTag != "" {
		merged.RequestTag = opts.RequestTag
	}
	if opts.Priority != sppb.RequestOptions_PRIORITY_UNSPECIFIED {
		merged.Priority = opts.Priority
	}
	proto.Merge(merged.Options, qo.Options)
	proto.Merge(merged.Options, opts.Options)
	return merged
}
func gologoo__createRequestOptions_930eb0d2556f62a47e9151ba4c271fab(prio sppb.RequestOptions_Priority, requestTag, transactionTag string) (ro *sppb.RequestOptions) {
	ro = &sppb.RequestOptions{}
	if prio != sppb.RequestOptions_PRIORITY_UNSPECIFIED {
		ro.Priority = prio
	}
	if requestTag != "" {
		ro.RequestTag = requestTag
	}
	if transactionTag != "" {
		ro.TransactionTag = transactionTag
	}
	return ro
}
func (t *txReadOnly) gologoo__Query_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, statement Statement) *RowIterator {
	mode := sppb.ExecuteSqlRequest_NORMAL
	return t.query(ctx, statement, QueryOptions{Mode: &mode, Options: t.qo.Options})
}
func (t *txReadOnly) gologoo__QueryWithOptions_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, statement Statement, opts QueryOptions) *RowIterator {
	return t.query(ctx, statement, t.qo.merge(opts))
}
func (t *txReadOnly) gologoo__QueryWithStats_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, statement Statement) *RowIterator {
	mode := sppb.ExecuteSqlRequest_PROFILE
	return t.query(ctx, statement, QueryOptions{Mode: &mode, Options: t.qo.Options})
}
func (t *txReadOnly) gologoo__AnalyzeQuery_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, statement Statement) (*sppb.QueryPlan, error) {
	mode := sppb.ExecuteSqlRequest_PLAN
	iter := t.query(ctx, statement, QueryOptions{Mode: &mode, Options: t.qo.Options})
	defer iter.Stop()
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	if iter.QueryPlan == nil {
		return nil, spannerErrorf(codes.Internal, "query plan unavailable")
	}
	return iter.QueryPlan, nil
}
func (t *txReadOnly) gologoo__query_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, statement Statement, options QueryOptions) (ri *RowIterator) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.Query")
	defer func() {
		trace.EndSpan(ctx, ri.err)
	}()
	req, sh, err := t.prepareExecuteSQL(ctx, statement, options)
	if err != nil {
		return &RowIterator{err: err}
	}
	client := sh.getClient()
	return streamWithReplaceSessionFunc(contextWithOutgoingMetadata(ctx, sh.getMetadata()), sh.session.logger, func(ctx context.Context, resumeToken []byte) (streamingReceiver, error) {
		req.ResumeToken = resumeToken
		req.Session = t.sh.getID()
		client, err := client.ExecuteStreamingSql(ctx, req)
		if err != nil {
			return client, err
		}
		md, err := client.Header()
		if getGFELatencyMetricsFlag() && md != nil && t.ct != nil {
			if err := createContextAndCaptureGFELatencyMetrics(ctx, t.ct, md, "query"); err != nil {
				trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", err)
			}
		}
		return client, err
	}, t.replaceSessionFunc, t.setTimestamp, t.release)
}
func (t *txReadOnly) gologoo__prepareExecuteSQL_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, stmt Statement, options QueryOptions) (*sppb.ExecuteSqlRequest, *sessionHandle, error) {
	sh, ts, err := t.acquire(ctx)
	if err != nil {
		return nil, nil, err
	}
	sid := sh.getID()
	if sid == "" {
		return nil, nil, errSessionClosed(sh)
	}
	params, paramTypes, err := stmt.convertParams()
	if err != nil {
		return nil, nil, err
	}
	mode := sppb.ExecuteSqlRequest_NORMAL
	if options.Mode != nil {
		mode = *options.Mode
	}
	req := &sppb.ExecuteSqlRequest{Session: sid, Transaction: ts, Sql: stmt.SQL, QueryMode: mode, Seqno: atomic.AddInt64(&t.sequenceNumber, 1), Params: params, ParamTypes: paramTypes, QueryOptions: options.Options, RequestOptions: createRequestOptions(options.Priority, options.RequestTag, t.txOpts.TransactionTag)}
	return req, sh, nil
}

type txState int

const (
	txNew txState = iota
	txInit
	txActive
	txClosed
)

func gologoo__errRtsUnavailable_930eb0d2556f62a47e9151ba4c271fab() error {
	return spannerErrorf(codes.Internal, "read timestamp is unavailable")
}
func gologoo__errTxClosed_930eb0d2556f62a47e9151ba4c271fab() error {
	return spannerErrorf(codes.InvalidArgument, "cannot use a closed transaction")
}
func gologoo__errUnexpectedTxState_930eb0d2556f62a47e9151ba4c271fab(ts txState) error {
	return spannerErrorf(codes.FailedPrecondition, "unexpected transaction state: %v", ts)
}

type ReadOnlyTransaction struct {
	mu sync.Mutex
	txReadOnly
	singleUse       bool
	tx              transactionID
	txReadyOrClosed chan struct {
	}
	state txState
	rts   time.Time
	tb    TimestampBound
}

func gologoo__errTxInitTimeout_930eb0d2556f62a47e9151ba4c271fab() error {
	return spannerErrorf(codes.Canceled, "timeout/context canceled in waiting for transaction's initialization")
}
func (t *ReadOnlyTransaction) gologoo__getTimestampBound_930eb0d2556f62a47e9151ba4c271fab() TimestampBound {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.tb
}
func (t *ReadOnlyTransaction) gologoo__begin_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context) error {
	var (
		locked bool
		tx     transactionID
		rts    time.Time
		sh     *sessionHandle
		err    error
		res    *sppb.Transaction
	)
	defer func() {
		if !locked {
			t.mu.Lock()
			locked = true
		}
		if t.state != txClosed {
			close(t.txReadyOrClosed)
			t.txReadyOrClosed = make(chan struct {
			})
		}
		t.mu.Unlock()
		if err != nil && sh != nil {
			if isSessionNotFoundError(err) {
				sh.destroy()
			}
			sh.recycle()
		}
	}()
	for {
		sh, err = t.sp.take(ctx)
		if err != nil {
			return err
		}
		var md metadata.MD
		res, err = sh.getClient().BeginTransaction(contextWithOutgoingMetadata(ctx, sh.getMetadata()), &sppb.BeginTransactionRequest{Session: sh.getID(), Options: &sppb.TransactionOptions{Mode: &sppb.TransactionOptions_ReadOnly_{ReadOnly: buildTransactionOptionsReadOnly(t.getTimestampBound(), true)}}}, gax.WithGRPCOptions(grpc.Header(&md)))
		if getGFELatencyMetricsFlag() && md != nil && t.ct != nil {
			if err := createContextAndCaptureGFELatencyMetrics(ctx, t.ct, md, "begin_BeginTransaction"); err != nil {
				trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", err)
			}
		}
		if isSessionNotFoundError(err) {
			sh.destroy()
			continue
		} else if err == nil {
			tx = res.Id
			if res.ReadTimestamp != nil {
				rts = time.Unix(res.ReadTimestamp.Seconds, int64(res.ReadTimestamp.Nanos))
			}
		} else {
			err = ToSpannerError(err)
		}
		break
	}
	t.mu.Lock()
	locked = true
	if t.state == txClosed {
		return errSessionClosed(sh)
	}
	t.tx = nil
	if err == nil {
		t.tx = tx
		t.rts = rts
		t.sh = sh
		t.state = txActive
	}
	return err
}
func (t *ReadOnlyTransaction) gologoo__acquire_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context) (*sessionHandle, *sppb.TransactionSelector, error) {
	if err := checkNestedTxn(ctx); err != nil {
		return nil, nil, err
	}
	if t.singleUse {
		return t.acquireSingleUse(ctx)
	}
	return t.acquireMultiUse(ctx)
}
func (t *ReadOnlyTransaction) gologoo__acquireSingleUse_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context) (*sessionHandle, *sppb.TransactionSelector, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	switch t.state {
	case txClosed:
		return nil, nil, errTxClosed()
	case txNew:
		t.state = txClosed
		ts := &sppb.TransactionSelector{Selector: &sppb.TransactionSelector_SingleUse{SingleUse: &sppb.TransactionOptions{Mode: &sppb.TransactionOptions_ReadOnly_{ReadOnly: buildTransactionOptionsReadOnly(t.tb, true)}}}}
		sh, err := t.sp.take(ctx)
		if err != nil {
			return nil, nil, err
		}
		t.sh = sh
		return sh, ts, nil
	}
	us := t.state
	return nil, nil, errUnexpectedTxState(us)
}
func (t *ReadOnlyTransaction) gologoo__acquireMultiUse_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context) (*sessionHandle, *sppb.TransactionSelector, error) {
	for {
		t.mu.Lock()
		switch t.state {
		case txClosed:
			t.mu.Unlock()
			return nil, nil, errTxClosed()
		case txNew:
			t.state = txInit
			t.mu.Unlock()
			continue
		case txInit:
			if t.tx != nil {
				txReadyOrClosed := t.txReadyOrClosed
				t.mu.Unlock()
				select {
				case <-txReadyOrClosed:
					continue
				case <-ctx.Done():
					return nil, nil, errTxInitTimeout()
				}
			}
			t.tx = transactionID{}
			t.mu.Unlock()
			if err := t.begin(ctx); err != nil {
				return nil, nil, err
			}
			continue
		case txActive:
			sh := t.sh
			ts := &sppb.TransactionSelector{Selector: &sppb.TransactionSelector_Id{Id: t.tx}}
			t.mu.Unlock()
			return sh, ts, nil
		}
		state := t.state
		t.mu.Unlock()
		return nil, nil, errUnexpectedTxState(state)
	}
}
func (t *ReadOnlyTransaction) gologoo__setTimestamp_930eb0d2556f62a47e9151ba4c271fab(ts time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.rts.IsZero() {
		t.rts = ts
	}
}
func (t *ReadOnlyTransaction) gologoo__release_930eb0d2556f62a47e9151ba4c271fab(err error) {
	t.mu.Lock()
	sh := t.sh
	t.mu.Unlock()
	if sh != nil {
		if isSessionNotFoundError(err) {
			sh.destroy()
		}
		if t.singleUse {
			sh.recycle()
		}
	}
}
func (t *ReadOnlyTransaction) gologoo__Close_930eb0d2556f62a47e9151ba4c271fab() {
	if t.singleUse {
		return
	}
	t.mu.Lock()
	if t.state != txClosed {
		t.state = txClosed
		close(t.txReadyOrClosed)
	}
	sh := t.sh
	t.mu.Unlock()
	if sh == nil {
		return
	}
	if sh != nil {
		sh.recycle()
	}
}
func (t *ReadOnlyTransaction) gologoo__Timestamp_930eb0d2556f62a47e9151ba4c271fab() (time.Time, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.rts.IsZero() {
		return t.rts, errRtsUnavailable()
	}
	return t.rts, nil
}
func (t *ReadOnlyTransaction) gologoo__WithTimestampBound_930eb0d2556f62a47e9151ba4c271fab(tb TimestampBound) *ReadOnlyTransaction {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.state == txNew {
		t.tb = tb
	}
	return t
}

type ReadWriteTransaction struct {
	txReadOnly
	tx    transactionID
	mu    sync.Mutex
	state txState
	wb    []*Mutation
}

func (t *ReadWriteTransaction) gologoo__BufferWrite_930eb0d2556f62a47e9151ba4c271fab(ms []*Mutation) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.state == txClosed {
		return errTxClosed()
	}
	if t.state != txActive {
		return errUnexpectedTxState(t.state)
	}
	t.wb = append(t.wb, ms...)
	return nil
}
func (t *ReadWriteTransaction) gologoo__Update_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, stmt Statement) (rowCount int64, err error) {
	mode := sppb.ExecuteSqlRequest_NORMAL
	return t.update(ctx, stmt, QueryOptions{Mode: &mode, Options: t.qo.Options})
}
func (t *ReadWriteTransaction) gologoo__UpdateWithOptions_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, stmt Statement, opts QueryOptions) (rowCount int64, err error) {
	return t.update(ctx, stmt, t.qo.merge(opts))
}
func (t *ReadWriteTransaction) gologoo__update_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, stmt Statement, opts QueryOptions) (rowCount int64, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.Update")
	defer func() {
		trace.EndSpan(ctx, err)
	}()
	req, sh, err := t.prepareExecuteSQL(ctx, stmt, opts)
	if err != nil {
		return 0, err
	}
	var md metadata.MD
	resultSet, err := sh.getClient().ExecuteSql(contextWithOutgoingMetadata(ctx, sh.getMetadata()), req, gax.WithGRPCOptions(grpc.Header(&md)))
	if getGFELatencyMetricsFlag() && md != nil && t.ct != nil {
		if err := createContextAndCaptureGFELatencyMetrics(ctx, t.ct, md, "update"); err != nil {
			trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", err)
		}
	}
	if err != nil {
		return 0, ToSpannerError(err)
	}
	if resultSet.Stats == nil {
		return 0, spannerErrorf(codes.InvalidArgument, "query passed to Update: %q", stmt.SQL)
	}
	return extractRowCount(resultSet.Stats)
}
func (t *ReadWriteTransaction) gologoo__BatchUpdate_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, stmts []Statement) (_ []int64, err error) {
	return t.BatchUpdateWithOptions(ctx, stmts, QueryOptions{})
}
func (t *ReadWriteTransaction) gologoo__BatchUpdateWithOptions_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, stmts []Statement, opts QueryOptions) (_ []int64, err error) {
	return t.batchUpdateWithOptions(ctx, stmts, t.qo.merge(opts))
}
func (t *ReadWriteTransaction) gologoo__batchUpdateWithOptions_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, stmts []Statement, opts QueryOptions) (_ []int64, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.BatchUpdate")
	defer func() {
		trace.EndSpan(ctx, err)
	}()
	sh, ts, err := t.acquire(ctx)
	if err != nil {
		return nil, err
	}
	sid := sh.getID()
	if sid == "" {
		return nil, errSessionClosed(sh)
	}
	var sppbStmts []*sppb.ExecuteBatchDmlRequest_Statement
	for _, st := range stmts {
		params, paramTypes, err := st.convertParams()
		if err != nil {
			return nil, err
		}
		sppbStmts = append(sppbStmts, &sppb.ExecuteBatchDmlRequest_Statement{Sql: st.SQL, Params: params, ParamTypes: paramTypes})
	}
	var md metadata.MD
	resp, err := sh.getClient().ExecuteBatchDml(contextWithOutgoingMetadata(ctx, sh.getMetadata()), &sppb.ExecuteBatchDmlRequest{Session: sh.getID(), Transaction: ts, Statements: sppbStmts, Seqno: atomic.AddInt64(&t.sequenceNumber, 1), RequestOptions: createRequestOptions(opts.Priority, opts.RequestTag, t.txOpts.TransactionTag)}, gax.WithGRPCOptions(grpc.Header(&md)))
	if getGFELatencyMetricsFlag() && md != nil && t.ct != nil {
		if err := createContextAndCaptureGFELatencyMetrics(ctx, t.ct, md, "batchUpdateWithOptions"); err != nil {
			trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", ToSpannerError(err))
		}
	}
	if err != nil {
		return nil, ToSpannerError(err)
	}
	var counts []int64
	for _, rs := range resp.ResultSets {
		count, err := extractRowCount(rs.Stats)
		if err != nil {
			return nil, err
		}
		counts = append(counts, count)
	}
	if resp.Status != nil && resp.Status.Code != 0 {
		return counts, spannerErrorf(codes.Code(uint32(resp.Status.Code)), resp.Status.Message)
	}
	return counts, nil
}
func (t *ReadWriteTransaction) gologoo__acquire_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context) (*sessionHandle, *sppb.TransactionSelector, error) {
	ts := &sppb.TransactionSelector{Selector: &sppb.TransactionSelector_Id{Id: t.tx}}
	t.mu.Lock()
	defer t.mu.Unlock()
	switch t.state {
	case txClosed:
		return nil, nil, errTxClosed()
	case txActive:
		return t.sh, ts, nil
	}
	return nil, nil, errUnexpectedTxState(t.state)
}
func (t *ReadWriteTransaction) gologoo__release_930eb0d2556f62a47e9151ba4c271fab(err error) {
	t.mu.Lock()
	sh := t.sh
	t.mu.Unlock()
	if sh != nil && isSessionNotFoundError(err) {
		sh.destroy()
	}
}
func gologoo__beginTransaction_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, sid string, client *vkit.Client) (transactionID, error) {
	res, err := client.BeginTransaction(ctx, &sppb.BeginTransactionRequest{Session: sid, Options: &sppb.TransactionOptions{Mode: &sppb.TransactionOptions_ReadWrite_{ReadWrite: &sppb.TransactionOptions_ReadWrite{}}}})
	if err != nil {
		return nil, err
	}
	if res.Id == nil {
		return nil, spannerErrorf(codes.Unknown, "BeginTransaction returned a transaction with a nil ID.")
	}
	return res.Id, nil
}
func (t *ReadWriteTransaction) gologoo__begin_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context) error {
	if t.tx != nil {
		t.state = txActive
		return nil
	}
	tx, err := beginTransaction(contextWithOutgoingMetadata(ctx, t.sh.getMetadata()), t.sh.getID(), t.sh.getClient())
	if err == nil {
		t.tx = tx
		t.state = txActive
		return nil
	}
	if isSessionNotFoundError(err) {
		t.sh.destroy()
	}
	return err
}

type CommitResponse struct {
	CommitTs    time.Time
	CommitStats *sppb.CommitResponse_CommitStats
}
type CommitOptions struct {
	ReturnCommitStats bool
}

func (t *ReadWriteTransaction) gologoo__commit_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, options CommitOptions) (CommitResponse, error) {
	resp := CommitResponse{}
	t.mu.Lock()
	t.state = txClosed
	mPb, err := mutationsProto(t.wb)
	t.mu.Unlock()
	if err != nil {
		return resp, err
	}
	sid, client := t.sh.getID(), t.sh.getClient()
	if sid == "" || client == nil {
		return resp, errSessionClosed(t.sh)
	}
	res, e := client.Commit(contextWithOutgoingMetadata(ctx, t.sh.getMetadata()), &sppb.CommitRequest{Session: sid, Transaction: &sppb.CommitRequest_TransactionId{TransactionId: t.tx}, RequestOptions: createRequestOptions(t.txOpts.CommitPriority, "", t.txOpts.TransactionTag), Mutations: mPb, ReturnCommitStats: options.ReturnCommitStats})
	if e != nil {
		return resp, toSpannerErrorWithCommitInfo(e, true)
	}
	if tstamp := res.GetCommitTimestamp(); tstamp != nil {
		resp.CommitTs = time.Unix(tstamp.Seconds, int64(tstamp.Nanos))
	}
	if options.ReturnCommitStats {
		resp.CommitStats = res.CommitStats
	}
	if isSessionNotFoundError(err) {
		t.sh.destroy()
	}
	return resp, err
}
func (t *ReadWriteTransaction) gologoo__rollback_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context) {
	t.mu.Lock()
	t.state = txClosed
	t.mu.Unlock()
	sid, client := t.sh.getID(), t.sh.getClient()
	if sid == "" || client == nil {
		return
	}
	err := client.Rollback(contextWithOutgoingMetadata(ctx, t.sh.getMetadata()), &sppb.RollbackRequest{Session: sid, TransactionId: t.tx})
	if isSessionNotFoundError(err) {
		t.sh.destroy()
	}
}
func (t *ReadWriteTransaction) gologoo__runInTransaction_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, f func(context.Context, *ReadWriteTransaction) error) (CommitResponse, error) {
	var (
		resp            CommitResponse
		err             error
		errDuringCommit bool
	)
	if err = f(context.WithValue(ctx, transactionInProgressKey{}, 1), t); err == nil {
		resp, err = t.commit(ctx, t.txOpts.CommitOptions)
		errDuringCommit = err != nil
	}
	if err != nil {
		if isAbortedErr(err) {
			return resp, err
		}
		if isSessionNotFoundError(err) {
			t.sh.destroy()
			return resp, err
		}
		if !errDuringCommit {
			t.rollback(ctx)
		}
		return resp, err
	}
	return resp, nil
}

type ReadWriteStmtBasedTransaction struct {
	ReadWriteTransaction
	options TransactionOptions
}

func gologoo__NewReadWriteStmtBasedTransaction_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, c *Client) (*ReadWriteStmtBasedTransaction, error) {
	return NewReadWriteStmtBasedTransactionWithOptions(ctx, c, TransactionOptions{})
}
func gologoo__NewReadWriteStmtBasedTransactionWithOptions_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, c *Client, options TransactionOptions) (*ReadWriteStmtBasedTransaction, error) {
	var (
		sh  *sessionHandle
		err error
		t   *ReadWriteStmtBasedTransaction
	)
	sh, err = c.idleSessions.takeWriteSession(ctx)
	if err != nil {
		return nil, err
	}
	t = &ReadWriteStmtBasedTransaction{ReadWriteTransaction: ReadWriteTransaction{tx: sh.getTransactionID()}}
	t.txReadOnly.sh = sh
	t.txReadOnly.txReadEnv = t
	t.txReadOnly.qo = c.qo
	t.txOpts = options
	t.ct = c.ct
	if err = t.begin(ctx); err != nil {
		if sh != nil {
			sh.recycle()
		}
		return nil, err
	}
	return t, err
}
func (t *ReadWriteStmtBasedTransaction) gologoo__Commit_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context) (time.Time, error) {
	resp, err := t.CommitWithReturnResp(ctx)
	return resp.CommitTs, err
}
func (t *ReadWriteStmtBasedTransaction) gologoo__CommitWithReturnResp_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context) (CommitResponse, error) {
	resp, err := t.commit(ctx, t.txOpts.CommitOptions)
	if err != nil && status.Code(err) != codes.Aborted {
		t.rollback(ctx)
	}
	if t.sh != nil {
		t.sh.recycle()
	}
	return resp, err
}
func (t *ReadWriteStmtBasedTransaction) gologoo__Rollback_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context) {
	t.rollback(ctx)
	if t.sh != nil {
		t.sh.recycle()
	}
}

type writeOnlyTransaction struct {
	sp             *sessionPool
	transactionTag string
	commitPriority sppb.RequestOptions_Priority
}

func (t *writeOnlyTransaction) gologoo__applyAtLeastOnce_930eb0d2556f62a47e9151ba4c271fab(ctx context.Context, ms ...*Mutation) (time.Time, error) {
	var (
		ts time.Time
		sh *sessionHandle
	)
	mPb, err := mutationsProto(ms)
	if err != nil {
		return ts, err
	}
	for {
		if sh == nil || sh.getID() == "" || sh.getClient() == nil {
			sh, err = t.sp.take(ctx)
			if err != nil {
				return ts, err
			}
			defer sh.recycle()
		}
		res, err := sh.getClient().Commit(contextWithOutgoingMetadata(ctx, sh.getMetadata()), &sppb.CommitRequest{Session: sh.getID(), Transaction: &sppb.CommitRequest_SingleUseTransaction{SingleUseTransaction: &sppb.TransactionOptions{Mode: &sppb.TransactionOptions_ReadWrite_{ReadWrite: &sppb.TransactionOptions_ReadWrite{}}}}, Mutations: mPb, RequestOptions: createRequestOptions(t.commitPriority, "", t.transactionTag)})
		if err != nil && !isAbortedErr(err) {
			if isSessionNotFoundError(err) {
				sh.destroy()
			}
			return ts, toSpannerErrorWithCommitInfo(err, true)
		} else if err == nil {
			if tstamp := res.GetCommitTimestamp(); tstamp != nil {
				ts = time.Unix(tstamp.Seconds, int64(tstamp.Nanos))
			}
			break
		}
	}
	return ts, ToSpannerError(err)
}
func gologoo__isAbortedErr_930eb0d2556f62a47e9151ba4c271fab(err error) bool {
	if err == nil {
		return false
	}
	if ErrCode(err) == codes.Aborted {
		return true
	}
	return false
}
func (to *TransactionOptions) requestPriority() sppb.RequestOptions_Priority {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__requestPriority_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : (none)\n")
	r0 := to.gologoo__requestPriority_930eb0d2556f62a47e9151ba4c271fab()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (to *TransactionOptions) requestTag() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__requestTag_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : (none)\n")
	r0 := to.gologoo__requestTag_930eb0d2556f62a47e9151ba4c271fab()
	log.Printf("Output: %v\n", r0)
	return r0
}
func errSessionClosed(sh *sessionHandle) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errSessionClosed_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", sh)
	r0 := gologoo__errSessionClosed_930eb0d2556f62a47e9151ba4c271fab(sh)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *txReadOnly) Read(ctx context.Context, table string, keys KeySet, columns []string) *RowIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Read_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v %v\n", ctx, table, keys, columns)
	r0 := t.gologoo__Read_930eb0d2556f62a47e9151ba4c271fab(ctx, table, keys, columns)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *txReadOnly) ReadUsingIndex(ctx context.Context, table, index string, keys KeySet, columns []string) (ri *RowIterator) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ReadUsingIndex_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v %v %v\n", ctx, table, index, keys, columns)
	ri = t.gologoo__ReadUsingIndex_930eb0d2556f62a47e9151ba4c271fab(ctx, table, index, keys, columns)
	log.Printf("Output: %v\n", ri)
	return
}
func (t *txReadOnly) ReadWithOptions(ctx context.Context, table string, keys KeySet, columns []string, opts *ReadOptions) (ri *RowIterator) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ReadWithOptions_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v %v %v\n", ctx, table, keys, columns, opts)
	ri = t.gologoo__ReadWithOptions_930eb0d2556f62a47e9151ba4c271fab(ctx, table, keys, columns, opts)
	log.Printf("Output: %v\n", ri)
	return
}
func errRowNotFound(table string, key Key) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errRowNotFound_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v\n", table, key)
	r0 := gologoo__errRowNotFound_930eb0d2556f62a47e9151ba4c271fab(table, key)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errRowNotFoundByIndex(table string, key Key, index string) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errRowNotFoundByIndex_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v\n", table, key, index)
	r0 := gologoo__errRowNotFoundByIndex_930eb0d2556f62a47e9151ba4c271fab(table, key, index)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errMultipleRowsFound(table string, key Key, index string) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errMultipleRowsFound_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v\n", table, key, index)
	r0 := gologoo__errMultipleRowsFound_930eb0d2556f62a47e9151ba4c271fab(table, key, index)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *txReadOnly) ReadRow(ctx context.Context, table string, key Key, columns []string) (*Row, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ReadRow_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v %v\n", ctx, table, key, columns)
	r0, r1 := t.gologoo__ReadRow_930eb0d2556f62a47e9151ba4c271fab(ctx, table, key, columns)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *txReadOnly) ReadRowWithOptions(ctx context.Context, table string, key Key, columns []string, opts *ReadOptions) (*Row, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ReadRowWithOptions_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v %v %v\n", ctx, table, key, columns, opts)
	r0, r1 := t.gologoo__ReadRowWithOptions_930eb0d2556f62a47e9151ba4c271fab(ctx, table, key, columns, opts)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *txReadOnly) ReadRowUsingIndex(ctx context.Context, table string, index string, key Key, columns []string) (*Row, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ReadRowUsingIndex_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v %v %v\n", ctx, table, index, key, columns)
	r0, r1 := t.gologoo__ReadRowUsingIndex_930eb0d2556f62a47e9151ba4c271fab(ctx, table, index, key, columns)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (qo QueryOptions) merge(opts QueryOptions) QueryOptions {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__merge_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", opts)
	r0 := qo.gologoo__merge_930eb0d2556f62a47e9151ba4c271fab(opts)
	log.Printf("Output: %v\n", r0)
	return r0
}
func createRequestOptions(prio sppb.RequestOptions_Priority, requestTag, transactionTag string) (ro *sppb.RequestOptions) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__createRequestOptions_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v\n", prio, requestTag, transactionTag)
	ro = gologoo__createRequestOptions_930eb0d2556f62a47e9151ba4c271fab(prio, requestTag, transactionTag)
	log.Printf("Output: %v\n", ro)
	return
}
func (t *txReadOnly) Query(ctx context.Context, statement Statement) *RowIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Query_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v\n", ctx, statement)
	r0 := t.gologoo__Query_930eb0d2556f62a47e9151ba4c271fab(ctx, statement)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *txReadOnly) QueryWithOptions(ctx context.Context, statement Statement, opts QueryOptions) *RowIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__QueryWithOptions_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v\n", ctx, statement, opts)
	r0 := t.gologoo__QueryWithOptions_930eb0d2556f62a47e9151ba4c271fab(ctx, statement, opts)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *txReadOnly) QueryWithStats(ctx context.Context, statement Statement) *RowIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__QueryWithStats_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v\n", ctx, statement)
	r0 := t.gologoo__QueryWithStats_930eb0d2556f62a47e9151ba4c271fab(ctx, statement)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *txReadOnly) AnalyzeQuery(ctx context.Context, statement Statement) (*sppb.QueryPlan, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__AnalyzeQuery_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v\n", ctx, statement)
	r0, r1 := t.gologoo__AnalyzeQuery_930eb0d2556f62a47e9151ba4c271fab(ctx, statement)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *txReadOnly) query(ctx context.Context, statement Statement, options QueryOptions) (ri *RowIterator) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__query_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v\n", ctx, statement, options)
	ri = t.gologoo__query_930eb0d2556f62a47e9151ba4c271fab(ctx, statement, options)
	log.Printf("Output: %v\n", ri)
	return
}
func (t *txReadOnly) prepareExecuteSQL(ctx context.Context, stmt Statement, options QueryOptions) (*sppb.ExecuteSqlRequest, *sessionHandle, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__prepareExecuteSQL_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v\n", ctx, stmt, options)
	r0, r1, r2 := t.gologoo__prepareExecuteSQL_930eb0d2556f62a47e9151ba4c271fab(ctx, stmt, options)
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
func errRtsUnavailable() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errRtsUnavailable_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : (none)\n")
	r0 := gologoo__errRtsUnavailable_930eb0d2556f62a47e9151ba4c271fab()
	log.Printf("Output: %v\n", r0)
	return r0
}
func errTxClosed() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errTxClosed_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : (none)\n")
	r0 := gologoo__errTxClosed_930eb0d2556f62a47e9151ba4c271fab()
	log.Printf("Output: %v\n", r0)
	return r0
}
func errUnexpectedTxState(ts txState) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errUnexpectedTxState_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", ts)
	r0 := gologoo__errUnexpectedTxState_930eb0d2556f62a47e9151ba4c271fab(ts)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errTxInitTimeout() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errTxInitTimeout_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : (none)\n")
	r0 := gologoo__errTxInitTimeout_930eb0d2556f62a47e9151ba4c271fab()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *ReadOnlyTransaction) getTimestampBound() TimestampBound {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getTimestampBound_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : (none)\n")
	r0 := t.gologoo__getTimestampBound_930eb0d2556f62a47e9151ba4c271fab()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *ReadOnlyTransaction) begin(ctx context.Context) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__begin_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", ctx)
	r0 := t.gologoo__begin_930eb0d2556f62a47e9151ba4c271fab(ctx)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *ReadOnlyTransaction) acquire(ctx context.Context) (*sessionHandle, *sppb.TransactionSelector, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__acquire_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", ctx)
	r0, r1, r2 := t.gologoo__acquire_930eb0d2556f62a47e9151ba4c271fab(ctx)
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
func (t *ReadOnlyTransaction) acquireSingleUse(ctx context.Context) (*sessionHandle, *sppb.TransactionSelector, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__acquireSingleUse_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", ctx)
	r0, r1, r2 := t.gologoo__acquireSingleUse_930eb0d2556f62a47e9151ba4c271fab(ctx)
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
func (t *ReadOnlyTransaction) acquireMultiUse(ctx context.Context) (*sessionHandle, *sppb.TransactionSelector, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__acquireMultiUse_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", ctx)
	r0, r1, r2 := t.gologoo__acquireMultiUse_930eb0d2556f62a47e9151ba4c271fab(ctx)
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
func (t *ReadOnlyTransaction) setTimestamp(ts time.Time) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__setTimestamp_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", ts)
	t.gologoo__setTimestamp_930eb0d2556f62a47e9151ba4c271fab(ts)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (t *ReadOnlyTransaction) release(err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__release_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", err)
	t.gologoo__release_930eb0d2556f62a47e9151ba4c271fab(err)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (t *ReadOnlyTransaction) Close() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Close_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : (none)\n")
	t.gologoo__Close_930eb0d2556f62a47e9151ba4c271fab()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (t *ReadOnlyTransaction) Timestamp() (time.Time, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Timestamp_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : (none)\n")
	r0, r1 := t.gologoo__Timestamp_930eb0d2556f62a47e9151ba4c271fab()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *ReadOnlyTransaction) WithTimestampBound(tb TimestampBound) *ReadOnlyTransaction {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__WithTimestampBound_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", tb)
	r0 := t.gologoo__WithTimestampBound_930eb0d2556f62a47e9151ba4c271fab(tb)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *ReadWriteTransaction) BufferWrite(ms []*Mutation) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__BufferWrite_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", ms)
	r0 := t.gologoo__BufferWrite_930eb0d2556f62a47e9151ba4c271fab(ms)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *ReadWriteTransaction) Update(ctx context.Context, stmt Statement) (rowCount int64, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Update_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v\n", ctx, stmt)
	rowCount, err = t.gologoo__Update_930eb0d2556f62a47e9151ba4c271fab(ctx, stmt)
	log.Printf("Output: %v %v\n", rowCount, err)
	return
}
func (t *ReadWriteTransaction) UpdateWithOptions(ctx context.Context, stmt Statement, opts QueryOptions) (rowCount int64, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UpdateWithOptions_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v\n", ctx, stmt, opts)
	rowCount, err = t.gologoo__UpdateWithOptions_930eb0d2556f62a47e9151ba4c271fab(ctx, stmt, opts)
	log.Printf("Output: %v %v\n", rowCount, err)
	return
}
func (t *ReadWriteTransaction) update(ctx context.Context, stmt Statement, opts QueryOptions) (rowCount int64, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__update_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v\n", ctx, stmt, opts)
	rowCount, err = t.gologoo__update_930eb0d2556f62a47e9151ba4c271fab(ctx, stmt, opts)
	log.Printf("Output: %v %v\n", rowCount, err)
	return
}
func (t *ReadWriteTransaction) BatchUpdate(ctx context.Context, stmts []Statement) (_ []int64, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__BatchUpdate_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v\n", ctx, stmts)
	_, err = t.gologoo__BatchUpdate_930eb0d2556f62a47e9151ba4c271fab(ctx, stmts)
	log.Printf("Output: %v\n", err)
	return
}
func (t *ReadWriteTransaction) BatchUpdateWithOptions(ctx context.Context, stmts []Statement, opts QueryOptions) (_ []int64, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__BatchUpdateWithOptions_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v\n", ctx, stmts, opts)
	_, err = t.gologoo__BatchUpdateWithOptions_930eb0d2556f62a47e9151ba4c271fab(ctx, stmts, opts)
	log.Printf("Output: %v\n", err)
	return
}
func (t *ReadWriteTransaction) batchUpdateWithOptions(ctx context.Context, stmts []Statement, opts QueryOptions) (_ []int64, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__batchUpdateWithOptions_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v\n", ctx, stmts, opts)
	_, err = t.gologoo__batchUpdateWithOptions_930eb0d2556f62a47e9151ba4c271fab(ctx, stmts, opts)
	log.Printf("Output: %v\n", err)
	return
}
func (t *ReadWriteTransaction) acquire(ctx context.Context) (*sessionHandle, *sppb.TransactionSelector, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__acquire_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", ctx)
	r0, r1, r2 := t.gologoo__acquire_930eb0d2556f62a47e9151ba4c271fab(ctx)
	log.Printf("Output: %v %v %v\n", r0, r1, r2)
	return r0, r1, r2
}
func (t *ReadWriteTransaction) release(err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__release_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", err)
	t.gologoo__release_930eb0d2556f62a47e9151ba4c271fab(err)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func beginTransaction(ctx context.Context, sid string, client *vkit.Client) (transactionID, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__beginTransaction_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v\n", ctx, sid, client)
	r0, r1 := gologoo__beginTransaction_930eb0d2556f62a47e9151ba4c271fab(ctx, sid, client)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *ReadWriteTransaction) begin(ctx context.Context) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__begin_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", ctx)
	r0 := t.gologoo__begin_930eb0d2556f62a47e9151ba4c271fab(ctx)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *ReadWriteTransaction) commit(ctx context.Context, options CommitOptions) (CommitResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__commit_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v\n", ctx, options)
	r0, r1 := t.gologoo__commit_930eb0d2556f62a47e9151ba4c271fab(ctx, options)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *ReadWriteTransaction) rollback(ctx context.Context) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__rollback_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", ctx)
	t.gologoo__rollback_930eb0d2556f62a47e9151ba4c271fab(ctx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (t *ReadWriteTransaction) runInTransaction(ctx context.Context, f func(context.Context, *ReadWriteTransaction) error) (CommitResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__runInTransaction_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v\n", ctx, f)
	r0, r1 := t.gologoo__runInTransaction_930eb0d2556f62a47e9151ba4c271fab(ctx, f)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func NewReadWriteStmtBasedTransaction(ctx context.Context, c *Client) (*ReadWriteStmtBasedTransaction, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NewReadWriteStmtBasedTransaction_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v\n", ctx, c)
	r0, r1 := gologoo__NewReadWriteStmtBasedTransaction_930eb0d2556f62a47e9151ba4c271fab(ctx, c)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func NewReadWriteStmtBasedTransactionWithOptions(ctx context.Context, c *Client, options TransactionOptions) (*ReadWriteStmtBasedTransaction, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NewReadWriteStmtBasedTransactionWithOptions_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v %v\n", ctx, c, options)
	r0, r1 := gologoo__NewReadWriteStmtBasedTransactionWithOptions_930eb0d2556f62a47e9151ba4c271fab(ctx, c, options)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *ReadWriteStmtBasedTransaction) Commit(ctx context.Context) (time.Time, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Commit_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", ctx)
	r0, r1 := t.gologoo__Commit_930eb0d2556f62a47e9151ba4c271fab(ctx)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *ReadWriteStmtBasedTransaction) CommitWithReturnResp(ctx context.Context) (CommitResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CommitWithReturnResp_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", ctx)
	r0, r1 := t.gologoo__CommitWithReturnResp_930eb0d2556f62a47e9151ba4c271fab(ctx)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *ReadWriteStmtBasedTransaction) Rollback(ctx context.Context) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Rollback_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", ctx)
	t.gologoo__Rollback_930eb0d2556f62a47e9151ba4c271fab(ctx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (t *writeOnlyTransaction) applyAtLeastOnce(ctx context.Context, ms ...*Mutation) (time.Time, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__applyAtLeastOnce_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v %v\n", ctx, ms)
	r0, r1 := t.gologoo__applyAtLeastOnce_930eb0d2556f62a47e9151ba4c271fab(ctx, ms...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func isAbortedErr(err error) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isAbortedErr_930eb0d2556f62a47e9151ba4c271fab")
	log.Printf("Input : %v\n", err)
	r0 := gologoo__isAbortedErr_930eb0d2556f62a47e9151ba4c271fab(err)
	log.Printf("Output: %v\n", r0)
	return r0
}
