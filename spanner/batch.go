package spanner

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"time"

	"cloud.google.com/go/internal/trace"
	"github.com/golang/protobuf/proto"
	"github.com/googleapis/gax-go/v2"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type BatchReadOnlyTransaction struct {
	ReadOnlyTransaction
	ID BatchReadOnlyTransactionID
}
type BatchReadOnlyTransactionID struct {
	tid transactionID
	sid string
	rts time.Time
}
type Partition struct {
	pt   []byte
	qreq *sppb.ExecuteSqlRequest
	rreq *sppb.ReadRequest
}
type PartitionOptions struct {
	PartitionBytes int64
	MaxPartitions  int64
}

func (opt PartitionOptions) gologoo__toProto_103284c03bd25c2e68afa447f399070e() *sppb.PartitionOptions {
	return &sppb.PartitionOptions{PartitionSizeBytes: opt.PartitionBytes, MaxPartitions: opt.MaxPartitions}
}
func (t *BatchReadOnlyTransaction) gologoo__PartitionRead_103284c03bd25c2e68afa447f399070e(ctx context.Context, table string, keys KeySet, columns []string, opt PartitionOptions) ([]*Partition, error) {
	return t.PartitionReadUsingIndex(ctx, table, "", keys, columns, opt)
}
func (t *BatchReadOnlyTransaction) gologoo__PartitionReadWithOptions_103284c03bd25c2e68afa447f399070e(ctx context.Context, table string, keys KeySet, columns []string, opt PartitionOptions, readOptions ReadOptions) ([]*Partition, error) {
	return t.PartitionReadUsingIndexWithOptions(ctx, table, "", keys, columns, opt, readOptions)
}
func (t *BatchReadOnlyTransaction) gologoo__PartitionReadUsingIndex_103284c03bd25c2e68afa447f399070e(ctx context.Context, table, index string, keys KeySet, columns []string, opt PartitionOptions) ([]*Partition, error) {
	return t.PartitionReadUsingIndexWithOptions(ctx, table, index, keys, columns, opt, ReadOptions{})
}
func (t *BatchReadOnlyTransaction) gologoo__PartitionReadUsingIndexWithOptions_103284c03bd25c2e68afa447f399070e(ctx context.Context, table, index string, keys KeySet, columns []string, opt PartitionOptions, readOptions ReadOptions) ([]*Partition, error) {
	sh, ts, err := t.acquire(ctx)
	if err != nil {
		return nil, err
	}
	sid, client := sh.getID(), sh.getClient()
	var (
		kset       *sppb.KeySet
		resp       *sppb.PartitionResponse
		partitions []*Partition
	)
	kset, err = keys.keySetProto()
	if err != nil {
		return nil, err
	}
	var md metadata.MD
	resp, err = client.PartitionRead(contextWithOutgoingMetadata(ctx, sh.getMetadata()), &sppb.PartitionReadRequest{Session: sid, Transaction: ts, Table: table, Index: index, Columns: columns, KeySet: kset, PartitionOptions: opt.toProto()}, gax.WithGRPCOptions(grpc.Header(&md)))
	if getGFELatencyMetricsFlag() && md != nil && t.ct != nil {
		if err := createContextAndCaptureGFELatencyMetrics(ctx, t.ct, md, "PartitionReadUsingIndexWithOptions"); err != nil {
			trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", err)
		}
	}
	req := &sppb.ReadRequest{Session: sid, Transaction: ts, Table: table, Index: index, Columns: columns, KeySet: kset, RequestOptions: createRequestOptions(readOptions.Priority, readOptions.RequestTag, "")}
	for _, p := range resp.GetPartitions() {
		partitions = append(partitions, &Partition{pt: p.PartitionToken, rreq: req})
	}
	return partitions, err
}
func (t *BatchReadOnlyTransaction) gologoo__PartitionQuery_103284c03bd25c2e68afa447f399070e(ctx context.Context, statement Statement, opt PartitionOptions) ([]*Partition, error) {
	return t.partitionQuery(ctx, statement, opt, t.ReadOnlyTransaction.txReadOnly.qo)
}
func (t *BatchReadOnlyTransaction) gologoo__PartitionQueryWithOptions_103284c03bd25c2e68afa447f399070e(ctx context.Context, statement Statement, opt PartitionOptions, qOpts QueryOptions) ([]*Partition, error) {
	return t.partitionQuery(ctx, statement, opt, t.ReadOnlyTransaction.txReadOnly.qo.merge(qOpts))
}
func (t *BatchReadOnlyTransaction) gologoo__partitionQuery_103284c03bd25c2e68afa447f399070e(ctx context.Context, statement Statement, opt PartitionOptions, qOpts QueryOptions) ([]*Partition, error) {
	sh, ts, err := t.acquire(ctx)
	if err != nil {
		return nil, err
	}
	sid, client := sh.getID(), sh.getClient()
	params, paramTypes, err := statement.convertParams()
	if err != nil {
		return nil, err
	}
	var md metadata.MD
	req := &sppb.PartitionQueryRequest{Session: sid, Transaction: ts, Sql: statement.SQL, PartitionOptions: opt.toProto(), Params: params, ParamTypes: paramTypes}
	resp, err := client.PartitionQuery(contextWithOutgoingMetadata(ctx, sh.getMetadata()), req, gax.WithGRPCOptions(grpc.Header(&md)))
	if getGFELatencyMetricsFlag() && md != nil && t.ct != nil {
		if err := createContextAndCaptureGFELatencyMetrics(ctx, t.ct, md, "partitionQuery"); err != nil {
			trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", err)
		}
	}
	r := &sppb.ExecuteSqlRequest{Session: sid, Transaction: ts, Sql: statement.SQL, Params: params, ParamTypes: paramTypes, QueryOptions: qOpts.Options, RequestOptions: createRequestOptions(qOpts.Priority, qOpts.RequestTag, "")}
	var partitions []*Partition
	for _, p := range resp.GetPartitions() {
		partitions = append(partitions, &Partition{pt: p.PartitionToken, qreq: r})
	}
	return partitions, err
}
func (t *BatchReadOnlyTransaction) gologoo__release_103284c03bd25c2e68afa447f399070e(err error) {
}
func (t *BatchReadOnlyTransaction) gologoo__setTimestamp_103284c03bd25c2e68afa447f399070e(ts time.Time) {
}
func (t *BatchReadOnlyTransaction) gologoo__Close_103284c03bd25c2e68afa447f399070e() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.state = txClosed
}
func (t *BatchReadOnlyTransaction) gologoo__Cleanup_103284c03bd25c2e68afa447f399070e(ctx context.Context) {
	t.Close()
	t.mu.Lock()
	defer t.mu.Unlock()
	sh := t.sh
	if sh == nil {
		return
	}
	t.sh = nil
	sid, client := sh.getID(), sh.getClient()
	var md metadata.MD
	err := client.DeleteSession(contextWithOutgoingMetadata(ctx, sh.getMetadata()), &sppb.DeleteSessionRequest{Name: sid}, gax.WithGRPCOptions(grpc.Header(&md)))
	if getGFELatencyMetricsFlag() && md != nil && t.ct != nil {
		if err := createContextAndCaptureGFELatencyMetrics(ctx, t.ct, md, "Cleanup"); err != nil {
			trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", err)
		}
	}
	if err != nil {
		var logger *log.Logger
		if sh.session != nil {
			logger = sh.session.logger
		}
		logf(logger, "Failed to delete session %v. Error: %v", sid, err)
	}
}
func (t *BatchReadOnlyTransaction) gologoo__Execute_103284c03bd25c2e68afa447f399070e(ctx context.Context, p *Partition) *RowIterator {
	var (
		sh  *sessionHandle
		err error
		rpc func(ct context.Context, resumeToken []byte) (streamingReceiver, error)
	)
	if sh, _, err = t.acquire(ctx); err != nil {
		return &RowIterator{err: err}
	}
	client := sh.getClient()
	if client == nil {
		return &RowIterator{err: errSessionClosed(sh)}
	}
	if p.rreq != nil {
		rpc = func(ctx context.Context, resumeToken []byte) (streamingReceiver, error) {
			client, err := client.StreamingRead(ctx, &sppb.ReadRequest{Session: p.rreq.Session, Transaction: p.rreq.Transaction, Table: p.rreq.Table, Index: p.rreq.Index, Columns: p.rreq.Columns, KeySet: p.rreq.KeySet, PartitionToken: p.pt, RequestOptions: p.rreq.RequestOptions, ResumeToken: resumeToken})
			if err != nil {
				return client, err
			}
			md, err := client.Header()
			if getGFELatencyMetricsFlag() && md != nil && t.ct != nil {
				if err := createContextAndCaptureGFELatencyMetrics(ctx, t.ct, md, "Execute"); err != nil {
					trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", err)
				}
			}
			return client, err
		}
	} else {
		rpc = func(ctx context.Context, resumeToken []byte) (streamingReceiver, error) {
			client, err := client.ExecuteStreamingSql(ctx, &sppb.ExecuteSqlRequest{Session: p.qreq.Session, Transaction: p.qreq.Transaction, Sql: p.qreq.Sql, Params: p.qreq.Params, ParamTypes: p.qreq.ParamTypes, QueryOptions: p.qreq.QueryOptions, PartitionToken: p.pt, RequestOptions: p.qreq.RequestOptions, ResumeToken: resumeToken})
			if err != nil {
				return client, err
			}
			md, err := client.Header()
			if getGFELatencyMetricsFlag() && md != nil && t.ct != nil {
				if err := createContextAndCaptureGFELatencyMetrics(ctx, t.ct, md, "Execute"); err != nil {
					trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", err)
				}
			}
			return client, err
		}
	}
	return stream(contextWithOutgoingMetadata(ctx, sh.getMetadata()), sh.session.logger, rpc, t.setTimestamp, t.release)
}
func (tid BatchReadOnlyTransactionID) gologoo__MarshalBinary_103284c03bd25c2e68afa447f399070e() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(tid.tid); err != nil {
		return nil, err
	}
	if err := enc.Encode(tid.sid); err != nil {
		return nil, err
	}
	if err := enc.Encode(tid.rts); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (tid *BatchReadOnlyTransactionID) gologoo__UnmarshalBinary_103284c03bd25c2e68afa447f399070e(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&tid.tid); err != nil {
		return err
	}
	if err := dec.Decode(&tid.sid); err != nil {
		return err
	}
	return dec.Decode(&tid.rts)
}
func (p Partition) gologoo__MarshalBinary_103284c03bd25c2e68afa447f399070e() (data []byte, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(p.pt); err != nil {
		return nil, err
	}
	var isReadPartition bool
	var req proto.Message
	if p.rreq != nil {
		isReadPartition = true
		req = p.rreq
	} else {
		isReadPartition = false
		req = p.qreq
	}
	if err := enc.Encode(isReadPartition); err != nil {
		return nil, err
	}
	if data, err = proto.Marshal(req); err != nil {
		return nil, err
	}
	if err := enc.Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (p *Partition) gologoo__UnmarshalBinary_103284c03bd25c2e68afa447f399070e(data []byte) error {
	var (
		isReadPartition bool
		d               []byte
		err             error
	)
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&p.pt); err != nil {
		return err
	}
	if err := dec.Decode(&isReadPartition); err != nil {
		return err
	}
	if err := dec.Decode(&d); err != nil {
		return err
	}
	if isReadPartition {
		p.rreq = &sppb.ReadRequest{}
		err = proto.Unmarshal(d, p.rreq)
	} else {
		p.qreq = &sppb.ExecuteSqlRequest{}
		err = proto.Unmarshal(d, p.qreq)
	}
	return err
}
func (opt PartitionOptions) toProto() *sppb.PartitionOptions {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__toProto_103284c03bd25c2e68afa447f399070e")
	log.Printf("Input : (none)\n")
	r0 := opt.gologoo__toProto_103284c03bd25c2e68afa447f399070e()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (t *BatchReadOnlyTransaction) PartitionRead(ctx context.Context, table string, keys KeySet, columns []string, opt PartitionOptions) ([]*Partition, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PartitionRead_103284c03bd25c2e68afa447f399070e")
	log.Printf("Input : %v %v %v %v %v\n", ctx, table, keys, columns, opt)
	r0, r1 := t.gologoo__PartitionRead_103284c03bd25c2e68afa447f399070e(ctx, table, keys, columns, opt)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *BatchReadOnlyTransaction) PartitionReadWithOptions(ctx context.Context, table string, keys KeySet, columns []string, opt PartitionOptions, readOptions ReadOptions) ([]*Partition, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PartitionReadWithOptions_103284c03bd25c2e68afa447f399070e")
	log.Printf("Input : %v %v %v %v %v %v\n", ctx, table, keys, columns, opt, readOptions)
	r0, r1 := t.gologoo__PartitionReadWithOptions_103284c03bd25c2e68afa447f399070e(ctx, table, keys, columns, opt, readOptions)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *BatchReadOnlyTransaction) PartitionReadUsingIndex(ctx context.Context, table, index string, keys KeySet, columns []string, opt PartitionOptions) ([]*Partition, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PartitionReadUsingIndex_103284c03bd25c2e68afa447f399070e")
	log.Printf("Input : %v %v %v %v %v %v\n", ctx, table, index, keys, columns, opt)
	r0, r1 := t.gologoo__PartitionReadUsingIndex_103284c03bd25c2e68afa447f399070e(ctx, table, index, keys, columns, opt)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *BatchReadOnlyTransaction) PartitionReadUsingIndexWithOptions(ctx context.Context, table, index string, keys KeySet, columns []string, opt PartitionOptions, readOptions ReadOptions) ([]*Partition, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PartitionReadUsingIndexWithOptions_103284c03bd25c2e68afa447f399070e")
	log.Printf("Input : %v %v %v %v %v %v %v\n", ctx, table, index, keys, columns, opt, readOptions)
	r0, r1 := t.gologoo__PartitionReadUsingIndexWithOptions_103284c03bd25c2e68afa447f399070e(ctx, table, index, keys, columns, opt, readOptions)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *BatchReadOnlyTransaction) PartitionQuery(ctx context.Context, statement Statement, opt PartitionOptions) ([]*Partition, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PartitionQuery_103284c03bd25c2e68afa447f399070e")
	log.Printf("Input : %v %v %v\n", ctx, statement, opt)
	r0, r1 := t.gologoo__PartitionQuery_103284c03bd25c2e68afa447f399070e(ctx, statement, opt)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *BatchReadOnlyTransaction) PartitionQueryWithOptions(ctx context.Context, statement Statement, opt PartitionOptions, qOpts QueryOptions) ([]*Partition, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PartitionQueryWithOptions_103284c03bd25c2e68afa447f399070e")
	log.Printf("Input : %v %v %v %v\n", ctx, statement, opt, qOpts)
	r0, r1 := t.gologoo__PartitionQueryWithOptions_103284c03bd25c2e68afa447f399070e(ctx, statement, opt, qOpts)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *BatchReadOnlyTransaction) partitionQuery(ctx context.Context, statement Statement, opt PartitionOptions, qOpts QueryOptions) ([]*Partition, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__partitionQuery_103284c03bd25c2e68afa447f399070e")
	log.Printf("Input : %v %v %v %v\n", ctx, statement, opt, qOpts)
	r0, r1 := t.gologoo__partitionQuery_103284c03bd25c2e68afa447f399070e(ctx, statement, opt, qOpts)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (t *BatchReadOnlyTransaction) release(err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__release_103284c03bd25c2e68afa447f399070e")
	log.Printf("Input : %v\n", err)
	t.gologoo__release_103284c03bd25c2e68afa447f399070e(err)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (t *BatchReadOnlyTransaction) setTimestamp(ts time.Time) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__setTimestamp_103284c03bd25c2e68afa447f399070e")
	log.Printf("Input : %v\n", ts)
	t.gologoo__setTimestamp_103284c03bd25c2e68afa447f399070e(ts)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (t *BatchReadOnlyTransaction) Close() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Close_103284c03bd25c2e68afa447f399070e")
	log.Printf("Input : (none)\n")
	t.gologoo__Close_103284c03bd25c2e68afa447f399070e()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (t *BatchReadOnlyTransaction) Cleanup(ctx context.Context) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Cleanup_103284c03bd25c2e68afa447f399070e")
	log.Printf("Input : %v\n", ctx)
	t.gologoo__Cleanup_103284c03bd25c2e68afa447f399070e(ctx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (t *BatchReadOnlyTransaction) Execute(ctx context.Context, p *Partition) *RowIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Execute_103284c03bd25c2e68afa447f399070e")
	log.Printf("Input : %v %v\n", ctx, p)
	r0 := t.gologoo__Execute_103284c03bd25c2e68afa447f399070e(ctx, p)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (tid BatchReadOnlyTransactionID) MarshalBinary() (data []byte, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__MarshalBinary_103284c03bd25c2e68afa447f399070e")
	log.Printf("Input : (none)\n")
	data, err = tid.gologoo__MarshalBinary_103284c03bd25c2e68afa447f399070e()
	log.Printf("Output: %v %v\n", data, err)
	return
}
func (tid *BatchReadOnlyTransactionID) UnmarshalBinary(data []byte) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UnmarshalBinary_103284c03bd25c2e68afa447f399070e")
	log.Printf("Input : %v\n", data)
	r0 := tid.gologoo__UnmarshalBinary_103284c03bd25c2e68afa447f399070e(data)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p Partition) MarshalBinary() (data []byte, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__MarshalBinary_103284c03bd25c2e68afa447f399070e")
	log.Printf("Input : (none)\n")
	data, err = p.gologoo__MarshalBinary_103284c03bd25c2e68afa447f399070e()
	log.Printf("Output: %v %v\n", data, err)
	return
}
func (p *Partition) UnmarshalBinary(data []byte) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UnmarshalBinary_103284c03bd25c2e68afa447f399070e")
	log.Printf("Input : %v\n", data)
	r0 := p.gologoo__UnmarshalBinary_103284c03bd25c2e68afa447f399070e(data)
	log.Printf("Output: %v\n", r0)
	return r0
}
