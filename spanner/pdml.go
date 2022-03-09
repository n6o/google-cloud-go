package spanner

import (
	"context"
	"cloud.google.com/go/internal/trace"
	"github.com/googleapis/gax-go/v2"
	"go.opencensus.io/tag"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"log"
)

func (c *Client) gologoo__PartitionedUpdate_9deb973da575c1105b5b9bd7b09f7903(ctx context.Context, statement Statement) (count int64, err error) {
	return c.partitionedUpdate(ctx, statement, c.qo)
}
func (c *Client) gologoo__PartitionedUpdateWithOptions_9deb973da575c1105b5b9bd7b09f7903(ctx context.Context, statement Statement, opts QueryOptions) (count int64, err error) {
	return c.partitionedUpdate(ctx, statement, c.qo.merge(opts))
}
func (c *Client) gologoo__partitionedUpdate_9deb973da575c1105b5b9bd7b09f7903(ctx context.Context, statement Statement, options QueryOptions) (count int64, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.PartitionedUpdate")
	defer func() {
		trace.EndSpan(ctx, err)
	}()
	if err := checkNestedTxn(ctx); err != nil {
		return 0, err
	}
	sh, err := c.idleSessions.take(ctx)
	if err != nil {
		return 0, ToSpannerError(err)
	}
	if sh != nil {
		defer sh.recycle()
	}
	params, paramTypes, err := statement.convertParams()
	if err != nil {
		return 0, ToSpannerError(err)
	}
	req := &sppb.ExecuteSqlRequest{Session: sh.getID(), Sql: statement.SQL, Params: params, ParamTypes: paramTypes, QueryOptions: options.Options, RequestOptions: createRequestOptions(options.Priority, options.RequestTag, "")}
	retryer := onCodes(DefaultRetryBackoff, codes.Aborted, codes.Internal)
	executePdmlWithRetry := func(ctx context.Context) (int64, error) {
		for {
			count, err := executePdml(ctx, sh, req)
			if err == nil {
				return count, nil
			}
			delay, shouldRetry := retryer.Retry(err)
			if !shouldRetry {
				return 0, err
			}
			if err := gax.Sleep(ctx, delay); err != nil {
				return 0, err
			}
		}
	}
	return executePdmlWithRetry(ctx)
}
func gologoo__executePdml_9deb973da575c1105b5b9bd7b09f7903(ctx context.Context, sh *sessionHandle, req *sppb.ExecuteSqlRequest) (count int64, err error) {
	var md metadata.MD
	res, err := sh.getClient().BeginTransaction(contextWithOutgoingMetadata(ctx, sh.getMetadata()), &sppb.BeginTransactionRequest{Session: sh.getID(), Options: &sppb.TransactionOptions{Mode: &sppb.TransactionOptions_PartitionedDml_{PartitionedDml: &sppb.TransactionOptions_PartitionedDml{}}}})
	if err != nil {
		return 0, ToSpannerError(err)
	}
	req.Transaction = &sppb.TransactionSelector{Selector: &sppb.TransactionSelector_Id{Id: res.Id}}
	resultSet, err := sh.getClient().ExecuteSql(contextWithOutgoingMetadata(ctx, sh.getMetadata()), req, gax.WithGRPCOptions(grpc.Header(&md)))
	if getGFELatencyMetricsFlag() && md != nil && sh.session.pool != nil {
		err := captureGFELatencyStats(tag.NewContext(ctx, sh.session.pool.tagMap), md, "executePdml_ExecuteSql")
		if err != nil {
			trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", err)
		}
	}
	if err != nil {
		return 0, err
	}
	if resultSet.Stats == nil {
		return 0, spannerErrorf(codes.InvalidArgument, "query passed to Update: %q", req.Sql)
	}
	return extractRowCount(resultSet.Stats)
}
func (c *Client) PartitionedUpdate(ctx context.Context, statement Statement) (count int64, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PartitionedUpdate_9deb973da575c1105b5b9bd7b09f7903")
	log.Printf("Input : %v %v\n", ctx, statement)
	count, err = c.gologoo__PartitionedUpdate_9deb973da575c1105b5b9bd7b09f7903(ctx, statement)
	log.Printf("Output: %v %v\n", count, err)
	return
}
func (c *Client) PartitionedUpdateWithOptions(ctx context.Context, statement Statement, opts QueryOptions) (count int64, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PartitionedUpdateWithOptions_9deb973da575c1105b5b9bd7b09f7903")
	log.Printf("Input : %v %v %v\n", ctx, statement, opts)
	count, err = c.gologoo__PartitionedUpdateWithOptions_9deb973da575c1105b5b9bd7b09f7903(ctx, statement, opts)
	log.Printf("Output: %v %v\n", count, err)
	return
}
func (c *Client) partitionedUpdate(ctx context.Context, statement Statement, options QueryOptions) (count int64, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__partitionedUpdate_9deb973da575c1105b5b9bd7b09f7903")
	log.Printf("Input : %v %v %v\n", ctx, statement, options)
	count, err = c.gologoo__partitionedUpdate_9deb973da575c1105b5b9bd7b09f7903(ctx, statement, options)
	log.Printf("Output: %v %v\n", count, err)
	return
}
func executePdml(ctx context.Context, sh *sessionHandle, req *sppb.ExecuteSqlRequest) (count int64, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__executePdml_9deb973da575c1105b5b9bd7b09f7903")
	log.Printf("Input : %v %v %v\n", ctx, sh, req)
	count, err = gologoo__executePdml_9deb973da575c1105b5b9bd7b09f7903(ctx, sh, req)
	log.Printf("Output: %v %v\n", count, err)
	return
}
