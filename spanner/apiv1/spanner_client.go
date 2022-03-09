package spanner

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"time"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	gtransport "google.golang.org/api/transport/grpc"
	spannerpb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"log"
)

var newClientHook clientHook

type CallOptions struct {
	CreateSession       []gax.CallOption
	BatchCreateSessions []gax.CallOption
	GetSession          []gax.CallOption
	ListSessions        []gax.CallOption
	DeleteSession       []gax.CallOption
	ExecuteSql          []gax.CallOption
	ExecuteStreamingSql []gax.CallOption
	ExecuteBatchDml     []gax.CallOption
	Read                []gax.CallOption
	StreamingRead       []gax.CallOption
	BeginTransaction    []gax.CallOption
	Commit              []gax.CallOption
	Rollback            []gax.CallOption
	PartitionQuery      []gax.CallOption
	PartitionRead       []gax.CallOption
}

func gologoo__defaultGRPCClientOptions_e4a94186c1a82597974e13e1f8bd3bbe() []option.ClientOption {
	return []option.ClientOption{internaloption.WithDefaultEndpoint("spanner.googleapis.com:443"), internaloption.WithDefaultMTLSEndpoint("spanner.mtls.googleapis.com:443"), internaloption.WithDefaultAudience("https://spanner.googleapis.com/"), internaloption.WithDefaultScopes(DefaultAuthScopes()...), internaloption.EnableJwtWithScope(), option.WithGRPCDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)))}
}
func gologoo__defaultCallOptions_e4a94186c1a82597974e13e1f8bd3bbe() *CallOptions {
	return &CallOptions{CreateSession: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable}, gax.Backoff{Initial: 250 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, BatchCreateSessions: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable}, gax.Backoff{Initial: 250 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, GetSession: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable}, gax.Backoff{Initial: 250 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, ListSessions: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable}, gax.Backoff{Initial: 250 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, DeleteSession: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable}, gax.Backoff{Initial: 250 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, ExecuteSql: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable}, gax.Backoff{Initial: 250 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, ExecuteStreamingSql: []gax.CallOption{}, ExecuteBatchDml: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable}, gax.Backoff{Initial: 250 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, Read: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable}, gax.Backoff{Initial: 250 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, StreamingRead: []gax.CallOption{}, BeginTransaction: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable}, gax.Backoff{Initial: 250 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, Commit: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable}, gax.Backoff{Initial: 250 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, Rollback: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable}, gax.Backoff{Initial: 250 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, PartitionQuery: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable}, gax.Backoff{Initial: 250 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, PartitionRead: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable}, gax.Backoff{Initial: 250 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}}
}

type internalClient interface {
	Close() error
	setGoogleClientInfo(...string)
	Connection() *grpc.ClientConn
	CreateSession(context.Context, *spannerpb.CreateSessionRequest, ...gax.CallOption) (*spannerpb.Session, error)
	BatchCreateSessions(context.Context, *spannerpb.BatchCreateSessionsRequest, ...gax.CallOption) (*spannerpb.BatchCreateSessionsResponse, error)
	GetSession(context.Context, *spannerpb.GetSessionRequest, ...gax.CallOption) (*spannerpb.Session, error)
	ListSessions(context.Context, *spannerpb.ListSessionsRequest, ...gax.CallOption) *SessionIterator
	DeleteSession(context.Context, *spannerpb.DeleteSessionRequest, ...gax.CallOption) error
	ExecuteSql(context.Context, *spannerpb.ExecuteSqlRequest, ...gax.CallOption) (*spannerpb.ResultSet, error)
	ExecuteStreamingSql(context.Context, *spannerpb.ExecuteSqlRequest, ...gax.CallOption) (spannerpb.Spanner_ExecuteStreamingSqlClient, error)
	ExecuteBatchDml(context.Context, *spannerpb.ExecuteBatchDmlRequest, ...gax.CallOption) (*spannerpb.ExecuteBatchDmlResponse, error)
	Read(context.Context, *spannerpb.ReadRequest, ...gax.CallOption) (*spannerpb.ResultSet, error)
	StreamingRead(context.Context, *spannerpb.ReadRequest, ...gax.CallOption) (spannerpb.Spanner_StreamingReadClient, error)
	BeginTransaction(context.Context, *spannerpb.BeginTransactionRequest, ...gax.CallOption) (*spannerpb.Transaction, error)
	Commit(context.Context, *spannerpb.CommitRequest, ...gax.CallOption) (*spannerpb.CommitResponse, error)
	Rollback(context.Context, *spannerpb.RollbackRequest, ...gax.CallOption) error
	PartitionQuery(context.Context, *spannerpb.PartitionQueryRequest, ...gax.CallOption) (*spannerpb.PartitionResponse, error)
	PartitionRead(context.Context, *spannerpb.PartitionReadRequest, ...gax.CallOption) (*spannerpb.PartitionResponse, error)
}
type Client struct {
	internalClient internalClient
	CallOptions    *CallOptions
}

func (c *Client) gologoo__Close_e4a94186c1a82597974e13e1f8bd3bbe() error {
	return c.internalClient.Close()
}
func (c *Client) gologoo__setGoogleClientInfo_e4a94186c1a82597974e13e1f8bd3bbe(keyval ...string) {
	c.internalClient.setGoogleClientInfo(keyval...)
}
func (c *Client) gologoo__Connection_e4a94186c1a82597974e13e1f8bd3bbe() *grpc.ClientConn {
	return c.internalClient.Connection()
}
func (c *Client) gologoo__CreateSession_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.CreateSessionRequest, opts ...gax.CallOption) (*spannerpb.Session, error) {
	return c.internalClient.CreateSession(ctx, req, opts...)
}
func (c *Client) gologoo__BatchCreateSessions_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.BatchCreateSessionsRequest, opts ...gax.CallOption) (*spannerpb.BatchCreateSessionsResponse, error) {
	return c.internalClient.BatchCreateSessions(ctx, req, opts...)
}
func (c *Client) gologoo__GetSession_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.GetSessionRequest, opts ...gax.CallOption) (*spannerpb.Session, error) {
	return c.internalClient.GetSession(ctx, req, opts...)
}
func (c *Client) gologoo__ListSessions_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.ListSessionsRequest, opts ...gax.CallOption) *SessionIterator {
	return c.internalClient.ListSessions(ctx, req, opts...)
}
func (c *Client) gologoo__DeleteSession_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.DeleteSessionRequest, opts ...gax.CallOption) error {
	return c.internalClient.DeleteSession(ctx, req, opts...)
}
func (c *Client) gologoo__ExecuteSql_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.ExecuteSqlRequest, opts ...gax.CallOption) (*spannerpb.ResultSet, error) {
	return c.internalClient.ExecuteSql(ctx, req, opts...)
}
func (c *Client) gologoo__ExecuteStreamingSql_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.ExecuteSqlRequest, opts ...gax.CallOption) (spannerpb.Spanner_ExecuteStreamingSqlClient, error) {
	return c.internalClient.ExecuteStreamingSql(ctx, req, opts...)
}
func (c *Client) gologoo__ExecuteBatchDml_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.ExecuteBatchDmlRequest, opts ...gax.CallOption) (*spannerpb.ExecuteBatchDmlResponse, error) {
	return c.internalClient.ExecuteBatchDml(ctx, req, opts...)
}
func (c *Client) gologoo__Read_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.ReadRequest, opts ...gax.CallOption) (*spannerpb.ResultSet, error) {
	return c.internalClient.Read(ctx, req, opts...)
}
func (c *Client) gologoo__StreamingRead_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.ReadRequest, opts ...gax.CallOption) (spannerpb.Spanner_StreamingReadClient, error) {
	return c.internalClient.StreamingRead(ctx, req, opts...)
}
func (c *Client) gologoo__BeginTransaction_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.BeginTransactionRequest, opts ...gax.CallOption) (*spannerpb.Transaction, error) {
	return c.internalClient.BeginTransaction(ctx, req, opts...)
}
func (c *Client) gologoo__Commit_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.CommitRequest, opts ...gax.CallOption) (*spannerpb.CommitResponse, error) {
	return c.internalClient.Commit(ctx, req, opts...)
}
func (c *Client) gologoo__Rollback_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.RollbackRequest, opts ...gax.CallOption) error {
	return c.internalClient.Rollback(ctx, req, opts...)
}
func (c *Client) gologoo__PartitionQuery_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.PartitionQueryRequest, opts ...gax.CallOption) (*spannerpb.PartitionResponse, error) {
	return c.internalClient.PartitionQuery(ctx, req, opts...)
}
func (c *Client) gologoo__PartitionRead_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.PartitionReadRequest, opts ...gax.CallOption) (*spannerpb.PartitionResponse, error) {
	return c.internalClient.PartitionRead(ctx, req, opts...)
}

type gRPCClient struct {
	connPool         gtransport.ConnPool
	disableDeadlines bool
	CallOptions      **CallOptions
	client           spannerpb.SpannerClient
	xGoogMetadata    metadata.MD
}

func gologoo__NewClient_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, opts ...option.ClientOption) (*Client, error) {
	clientOpts := defaultGRPCClientOptions()
	if newClientHook != nil {
		hookOpts, err := newClientHook(ctx, clientHookParams{})
		if err != nil {
			return nil, err
		}
		clientOpts = append(clientOpts, hookOpts...)
	}
	disableDeadlines, err := checkDisableDeadlines()
	if err != nil {
		return nil, err
	}
	connPool, err := gtransport.DialPool(ctx, append(clientOpts, opts...)...)
	if err != nil {
		return nil, err
	}
	client := Client{CallOptions: defaultCallOptions()}
	c := &gRPCClient{connPool: connPool, disableDeadlines: disableDeadlines, client: spannerpb.NewSpannerClient(connPool), CallOptions: &client.CallOptions}
	c.setGoogleClientInfo()
	client.internalClient = c
	return &client, nil
}
func (c *gRPCClient) gologoo__Connection_e4a94186c1a82597974e13e1f8bd3bbe() *grpc.ClientConn {
	return c.connPool.Conn()
}
func (c *gRPCClient) gologoo__setGoogleClientInfo_e4a94186c1a82597974e13e1f8bd3bbe(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", getVersionClient(), "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}
func (c *gRPCClient) gologoo__Close_e4a94186c1a82597974e13e1f8bd3bbe() error {
	return c.connPool.Close()
}
func (c *gRPCClient) gologoo__CreateSession_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.CreateSessionRequest, opts ...gax.CallOption) (*spannerpb.Session, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 30000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "database", url.QueryEscape(req.GetDatabase())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).CreateSession[0:len((*c.CallOptions).CreateSession):len((*c.CallOptions).CreateSession)], opts...)
	var resp *spannerpb.Session
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.CreateSession(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *gRPCClient) gologoo__BatchCreateSessions_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.BatchCreateSessionsRequest, opts ...gax.CallOption) (*spannerpb.BatchCreateSessionsResponse, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 60000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "database", url.QueryEscape(req.GetDatabase())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).BatchCreateSessions[0:len((*c.CallOptions).BatchCreateSessions):len((*c.CallOptions).BatchCreateSessions)], opts...)
	var resp *spannerpb.BatchCreateSessionsResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.BatchCreateSessions(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *gRPCClient) gologoo__GetSession_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.GetSessionRequest, opts ...gax.CallOption) (*spannerpb.Session, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 30000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).GetSession[0:len((*c.CallOptions).GetSession):len((*c.CallOptions).GetSession)], opts...)
	var resp *spannerpb.Session
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.GetSession(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *gRPCClient) gologoo__ListSessions_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.ListSessionsRequest, opts ...gax.CallOption) *SessionIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "database", url.QueryEscape(req.GetDatabase())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).ListSessions[0:len((*c.CallOptions).ListSessions):len((*c.CallOptions).ListSessions)], opts...)
	it := &SessionIterator{}
	req = proto.Clone(req).(*spannerpb.ListSessionsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*spannerpb.Session, string, error) {
		resp := &spannerpb.ListSessionsResponse{}
		if pageToken != "" {
			req.PageToken = pageToken
		}
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else if pageSize != 0 {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.client.ListSessions(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		it.Response = resp
		return resp.GetSessions(), resp.GetNextPageToken(), nil
	}
	fetch := func(pageSize int, pageToken string) (string, error) {
		items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)
		if err != nil {
			return "", err
		}
		it.items = append(it.items, items...)
		return nextPageToken, nil
	}
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)
	it.pageInfo.MaxSize = int(req.GetPageSize())
	it.pageInfo.Token = req.GetPageToken()
	return it
}
func (c *gRPCClient) gologoo__DeleteSession_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.DeleteSessionRequest, opts ...gax.CallOption) error {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 30000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).DeleteSession[0:len((*c.CallOptions).DeleteSession):len((*c.CallOptions).DeleteSession)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.client.DeleteSession(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}
func (c *gRPCClient) gologoo__ExecuteSql_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.ExecuteSqlRequest, opts ...gax.CallOption) (*spannerpb.ResultSet, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 30000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "session", url.QueryEscape(req.GetSession())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).ExecuteSql[0:len((*c.CallOptions).ExecuteSql):len((*c.CallOptions).ExecuteSql)], opts...)
	var resp *spannerpb.ResultSet
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.ExecuteSql(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *gRPCClient) gologoo__ExecuteStreamingSql_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.ExecuteSqlRequest, opts ...gax.CallOption) (spannerpb.Spanner_ExecuteStreamingSqlClient, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "session", url.QueryEscape(req.GetSession())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	var resp spannerpb.Spanner_ExecuteStreamingSqlClient
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.ExecuteStreamingSql(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *gRPCClient) gologoo__ExecuteBatchDml_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.ExecuteBatchDmlRequest, opts ...gax.CallOption) (*spannerpb.ExecuteBatchDmlResponse, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 30000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "session", url.QueryEscape(req.GetSession())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).ExecuteBatchDml[0:len((*c.CallOptions).ExecuteBatchDml):len((*c.CallOptions).ExecuteBatchDml)], opts...)
	var resp *spannerpb.ExecuteBatchDmlResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.ExecuteBatchDml(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *gRPCClient) gologoo__Read_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.ReadRequest, opts ...gax.CallOption) (*spannerpb.ResultSet, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 30000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "session", url.QueryEscape(req.GetSession())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).Read[0:len((*c.CallOptions).Read):len((*c.CallOptions).Read)], opts...)
	var resp *spannerpb.ResultSet
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.Read(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *gRPCClient) gologoo__StreamingRead_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.ReadRequest, opts ...gax.CallOption) (spannerpb.Spanner_StreamingReadClient, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "session", url.QueryEscape(req.GetSession())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	var resp spannerpb.Spanner_StreamingReadClient
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.StreamingRead(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *gRPCClient) gologoo__BeginTransaction_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.BeginTransactionRequest, opts ...gax.CallOption) (*spannerpb.Transaction, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 30000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "session", url.QueryEscape(req.GetSession())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).BeginTransaction[0:len((*c.CallOptions).BeginTransaction):len((*c.CallOptions).BeginTransaction)], opts...)
	var resp *spannerpb.Transaction
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.BeginTransaction(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *gRPCClient) gologoo__Commit_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.CommitRequest, opts ...gax.CallOption) (*spannerpb.CommitResponse, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 3600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "session", url.QueryEscape(req.GetSession())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).Commit[0:len((*c.CallOptions).Commit):len((*c.CallOptions).Commit)], opts...)
	var resp *spannerpb.CommitResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.Commit(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *gRPCClient) gologoo__Rollback_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.RollbackRequest, opts ...gax.CallOption) error {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 30000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "session", url.QueryEscape(req.GetSession())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).Rollback[0:len((*c.CallOptions).Rollback):len((*c.CallOptions).Rollback)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.client.Rollback(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}
func (c *gRPCClient) gologoo__PartitionQuery_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.PartitionQueryRequest, opts ...gax.CallOption) (*spannerpb.PartitionResponse, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 30000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "session", url.QueryEscape(req.GetSession())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).PartitionQuery[0:len((*c.CallOptions).PartitionQuery):len((*c.CallOptions).PartitionQuery)], opts...)
	var resp *spannerpb.PartitionResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.PartitionQuery(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *gRPCClient) gologoo__PartitionRead_e4a94186c1a82597974e13e1f8bd3bbe(ctx context.Context, req *spannerpb.PartitionReadRequest, opts ...gax.CallOption) (*spannerpb.PartitionResponse, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 30000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "session", url.QueryEscape(req.GetSession())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).PartitionRead[0:len((*c.CallOptions).PartitionRead):len((*c.CallOptions).PartitionRead)], opts...)
	var resp *spannerpb.PartitionResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.client.PartitionRead(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type SessionIterator struct {
	items    []*spannerpb.Session
	pageInfo *iterator.PageInfo
	nextFunc func() error
	Response interface {
	}
	InternalFetch func(pageSize int, pageToken string) (results []*spannerpb.Session, nextPageToken string, err error)
}

func (it *SessionIterator) gologoo__PageInfo_e4a94186c1a82597974e13e1f8bd3bbe() *iterator.PageInfo {
	return it.pageInfo
}
func (it *SessionIterator) gologoo__Next_e4a94186c1a82597974e13e1f8bd3bbe() (*spannerpb.Session, error) {
	var item *spannerpb.Session
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}
func (it *SessionIterator) gologoo__bufLen_e4a94186c1a82597974e13e1f8bd3bbe() int {
	return len(it.items)
}
func (it *SessionIterator) gologoo__takeBuf_e4a94186c1a82597974e13e1f8bd3bbe() interface {
} {
	b := it.items
	it.items = nil
	return b
}
func defaultGRPCClientOptions() []option.ClientOption {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__defaultGRPCClientOptions_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : (none)\n")
	r0 := gologoo__defaultGRPCClientOptions_e4a94186c1a82597974e13e1f8bd3bbe()
	log.Printf("Output: %v\n", r0)
	return r0
}
func defaultCallOptions() *CallOptions {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__defaultCallOptions_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : (none)\n")
	r0 := gologoo__defaultCallOptions_e4a94186c1a82597974e13e1f8bd3bbe()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *Client) Close() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Close_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__Close_e4a94186c1a82597974e13e1f8bd3bbe()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *Client) setGoogleClientInfo(keyval ...string) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__setGoogleClientInfo_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v\n", keyval)
	c.gologoo__setGoogleClientInfo_e4a94186c1a82597974e13e1f8bd3bbe(keyval...)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (c *Client) Connection() *grpc.ClientConn {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Connection_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__Connection_e4a94186c1a82597974e13e1f8bd3bbe()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *Client) CreateSession(ctx context.Context, req *spannerpb.CreateSessionRequest, opts ...gax.CallOption) (*spannerpb.Session, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateSession_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__CreateSession_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *Client) BatchCreateSessions(ctx context.Context, req *spannerpb.BatchCreateSessionsRequest, opts ...gax.CallOption) (*spannerpb.BatchCreateSessionsResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__BatchCreateSessions_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__BatchCreateSessions_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *Client) GetSession(ctx context.Context, req *spannerpb.GetSessionRequest, opts ...gax.CallOption) (*spannerpb.Session, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetSession_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__GetSession_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *Client) ListSessions(ctx context.Context, req *spannerpb.ListSessionsRequest, opts ...gax.CallOption) *SessionIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ListSessions_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__ListSessions_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *Client) DeleteSession(ctx context.Context, req *spannerpb.DeleteSessionRequest, opts ...gax.CallOption) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DeleteSession_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__DeleteSession_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *Client) ExecuteSql(ctx context.Context, req *spannerpb.ExecuteSqlRequest, opts ...gax.CallOption) (*spannerpb.ResultSet, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ExecuteSql_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__ExecuteSql_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *Client) ExecuteStreamingSql(ctx context.Context, req *spannerpb.ExecuteSqlRequest, opts ...gax.CallOption) (spannerpb.Spanner_ExecuteStreamingSqlClient, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ExecuteStreamingSql_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__ExecuteStreamingSql_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *Client) ExecuteBatchDml(ctx context.Context, req *spannerpb.ExecuteBatchDmlRequest, opts ...gax.CallOption) (*spannerpb.ExecuteBatchDmlResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ExecuteBatchDml_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__ExecuteBatchDml_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *Client) Read(ctx context.Context, req *spannerpb.ReadRequest, opts ...gax.CallOption) (*spannerpb.ResultSet, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Read_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__Read_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *Client) StreamingRead(ctx context.Context, req *spannerpb.ReadRequest, opts ...gax.CallOption) (spannerpb.Spanner_StreamingReadClient, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__StreamingRead_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__StreamingRead_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *Client) BeginTransaction(ctx context.Context, req *spannerpb.BeginTransactionRequest, opts ...gax.CallOption) (*spannerpb.Transaction, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__BeginTransaction_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__BeginTransaction_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *Client) Commit(ctx context.Context, req *spannerpb.CommitRequest, opts ...gax.CallOption) (*spannerpb.CommitResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Commit_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__Commit_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *Client) Rollback(ctx context.Context, req *spannerpb.RollbackRequest, opts ...gax.CallOption) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Rollback_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__Rollback_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *Client) PartitionQuery(ctx context.Context, req *spannerpb.PartitionQueryRequest, opts ...gax.CallOption) (*spannerpb.PartitionResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PartitionQuery_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__PartitionQuery_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *Client) PartitionRead(ctx context.Context, req *spannerpb.PartitionReadRequest, opts ...gax.CallOption) (*spannerpb.PartitionResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PartitionRead_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__PartitionRead_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func NewClient(ctx context.Context, opts ...option.ClientOption) (*Client, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NewClient_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v\n", ctx, opts)
	r0, r1 := gologoo__NewClient_e4a94186c1a82597974e13e1f8bd3bbe(ctx, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *gRPCClient) Connection() *grpc.ClientConn {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Connection_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__Connection_e4a94186c1a82597974e13e1f8bd3bbe()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *gRPCClient) setGoogleClientInfo(keyval ...string) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__setGoogleClientInfo_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v\n", keyval)
	c.gologoo__setGoogleClientInfo_e4a94186c1a82597974e13e1f8bd3bbe(keyval...)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (c *gRPCClient) Close() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Close_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__Close_e4a94186c1a82597974e13e1f8bd3bbe()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *gRPCClient) CreateSession(ctx context.Context, req *spannerpb.CreateSessionRequest, opts ...gax.CallOption) (*spannerpb.Session, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateSession_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__CreateSession_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *gRPCClient) BatchCreateSessions(ctx context.Context, req *spannerpb.BatchCreateSessionsRequest, opts ...gax.CallOption) (*spannerpb.BatchCreateSessionsResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__BatchCreateSessions_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__BatchCreateSessions_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *gRPCClient) GetSession(ctx context.Context, req *spannerpb.GetSessionRequest, opts ...gax.CallOption) (*spannerpb.Session, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetSession_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__GetSession_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *gRPCClient) ListSessions(ctx context.Context, req *spannerpb.ListSessionsRequest, opts ...gax.CallOption) *SessionIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ListSessions_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__ListSessions_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *gRPCClient) DeleteSession(ctx context.Context, req *spannerpb.DeleteSessionRequest, opts ...gax.CallOption) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DeleteSession_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__DeleteSession_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *gRPCClient) ExecuteSql(ctx context.Context, req *spannerpb.ExecuteSqlRequest, opts ...gax.CallOption) (*spannerpb.ResultSet, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ExecuteSql_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__ExecuteSql_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *gRPCClient) ExecuteStreamingSql(ctx context.Context, req *spannerpb.ExecuteSqlRequest, opts ...gax.CallOption) (spannerpb.Spanner_ExecuteStreamingSqlClient, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ExecuteStreamingSql_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__ExecuteStreamingSql_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *gRPCClient) ExecuteBatchDml(ctx context.Context, req *spannerpb.ExecuteBatchDmlRequest, opts ...gax.CallOption) (*spannerpb.ExecuteBatchDmlResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ExecuteBatchDml_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__ExecuteBatchDml_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *gRPCClient) Read(ctx context.Context, req *spannerpb.ReadRequest, opts ...gax.CallOption) (*spannerpb.ResultSet, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Read_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__Read_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *gRPCClient) StreamingRead(ctx context.Context, req *spannerpb.ReadRequest, opts ...gax.CallOption) (spannerpb.Spanner_StreamingReadClient, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__StreamingRead_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__StreamingRead_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *gRPCClient) BeginTransaction(ctx context.Context, req *spannerpb.BeginTransactionRequest, opts ...gax.CallOption) (*spannerpb.Transaction, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__BeginTransaction_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__BeginTransaction_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *gRPCClient) Commit(ctx context.Context, req *spannerpb.CommitRequest, opts ...gax.CallOption) (*spannerpb.CommitResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Commit_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__Commit_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *gRPCClient) Rollback(ctx context.Context, req *spannerpb.RollbackRequest, opts ...gax.CallOption) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Rollback_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__Rollback_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *gRPCClient) PartitionQuery(ctx context.Context, req *spannerpb.PartitionQueryRequest, opts ...gax.CallOption) (*spannerpb.PartitionResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PartitionQuery_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__PartitionQuery_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *gRPCClient) PartitionRead(ctx context.Context, req *spannerpb.PartitionReadRequest, opts ...gax.CallOption) (*spannerpb.PartitionResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PartitionRead_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__PartitionRead_e4a94186c1a82597974e13e1f8bd3bbe(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (it *SessionIterator) PageInfo() *iterator.PageInfo {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PageInfo_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__PageInfo_e4a94186c1a82597974e13e1f8bd3bbe()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (it *SessionIterator) Next() (*spannerpb.Session, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Next_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : (none)\n")
	r0, r1 := it.gologoo__Next_e4a94186c1a82597974e13e1f8bd3bbe()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (it *SessionIterator) bufLen() int {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__bufLen_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__bufLen_e4a94186c1a82597974e13e1f8bd3bbe()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (it *SessionIterator) takeBuf() interface {
} {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__takeBuf_e4a94186c1a82597974e13e1f8bd3bbe")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__takeBuf_e4a94186c1a82597974e13e1f8bd3bbe()
	log.Printf("Output: %v\n", r0)
	return r0
}
