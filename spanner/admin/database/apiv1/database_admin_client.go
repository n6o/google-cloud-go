package database

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"time"
	"cloud.google.com/go/longrunning"
	lroauto "cloud.google.com/go/longrunning/autogen"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	gtransport "google.golang.org/api/transport/grpc"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
	databasepb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"log"
)

var newDatabaseAdminClientHook clientHook

type DatabaseAdminCallOptions struct {
	ListDatabases          []gax.CallOption
	CreateDatabase         []gax.CallOption
	GetDatabase            []gax.CallOption
	UpdateDatabaseDdl      []gax.CallOption
	DropDatabase           []gax.CallOption
	GetDatabaseDdl         []gax.CallOption
	SetIamPolicy           []gax.CallOption
	GetIamPolicy           []gax.CallOption
	TestIamPermissions     []gax.CallOption
	CreateBackup           []gax.CallOption
	GetBackup              []gax.CallOption
	UpdateBackup           []gax.CallOption
	DeleteBackup           []gax.CallOption
	ListBackups            []gax.CallOption
	RestoreDatabase        []gax.CallOption
	ListDatabaseOperations []gax.CallOption
	ListBackupOperations   []gax.CallOption
}

func gologoo__defaultDatabaseAdminGRPCClientOptions_224d23f03999a157bdd6d773efae62a9() []option.ClientOption {
	return []option.ClientOption{internaloption.WithDefaultEndpoint("spanner.googleapis.com:443"), internaloption.WithDefaultMTLSEndpoint("spanner.mtls.googleapis.com:443"), internaloption.WithDefaultAudience("https://spanner.googleapis.com/"), internaloption.WithDefaultScopes(DefaultAuthScopes()...), internaloption.EnableJwtWithScope(), option.WithGRPCDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)))}
}
func gologoo__defaultDatabaseAdminCallOptions_224d23f03999a157bdd6d773efae62a9() *DatabaseAdminCallOptions {
	return &DatabaseAdminCallOptions{ListDatabases: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, CreateDatabase: []gax.CallOption{}, GetDatabase: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, UpdateDatabaseDdl: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, DropDatabase: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, GetDatabaseDdl: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, SetIamPolicy: []gax.CallOption{}, GetIamPolicy: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, TestIamPermissions: []gax.CallOption{}, CreateBackup: []gax.CallOption{}, GetBackup: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, UpdateBackup: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, DeleteBackup: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, ListBackups: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, RestoreDatabase: []gax.CallOption{}, ListDatabaseOperations: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, ListBackupOperations: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}}
}

type internalDatabaseAdminClient interface {
	Close() error
	setGoogleClientInfo(...string)
	Connection() *grpc.ClientConn
	ListDatabases(context.Context, *databasepb.ListDatabasesRequest, ...gax.CallOption) *DatabaseIterator
	CreateDatabase(context.Context, *databasepb.CreateDatabaseRequest, ...gax.CallOption) (*CreateDatabaseOperation, error)
	CreateDatabaseOperation(name string) *CreateDatabaseOperation
	GetDatabase(context.Context, *databasepb.GetDatabaseRequest, ...gax.CallOption) (*databasepb.Database, error)
	UpdateDatabaseDdl(context.Context, *databasepb.UpdateDatabaseDdlRequest, ...gax.CallOption) (*UpdateDatabaseDdlOperation, error)
	UpdateDatabaseDdlOperation(name string) *UpdateDatabaseDdlOperation
	DropDatabase(context.Context, *databasepb.DropDatabaseRequest, ...gax.CallOption) error
	GetDatabaseDdl(context.Context, *databasepb.GetDatabaseDdlRequest, ...gax.CallOption) (*databasepb.GetDatabaseDdlResponse, error)
	SetIamPolicy(context.Context, *iampb.SetIamPolicyRequest, ...gax.CallOption) (*iampb.Policy, error)
	GetIamPolicy(context.Context, *iampb.GetIamPolicyRequest, ...gax.CallOption) (*iampb.Policy, error)
	TestIamPermissions(context.Context, *iampb.TestIamPermissionsRequest, ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error)
	CreateBackup(context.Context, *databasepb.CreateBackupRequest, ...gax.CallOption) (*CreateBackupOperation, error)
	CreateBackupOperation(name string) *CreateBackupOperation
	GetBackup(context.Context, *databasepb.GetBackupRequest, ...gax.CallOption) (*databasepb.Backup, error)
	UpdateBackup(context.Context, *databasepb.UpdateBackupRequest, ...gax.CallOption) (*databasepb.Backup, error)
	DeleteBackup(context.Context, *databasepb.DeleteBackupRequest, ...gax.CallOption) error
	ListBackups(context.Context, *databasepb.ListBackupsRequest, ...gax.CallOption) *BackupIterator
	RestoreDatabase(context.Context, *databasepb.RestoreDatabaseRequest, ...gax.CallOption) (*RestoreDatabaseOperation, error)
	RestoreDatabaseOperation(name string) *RestoreDatabaseOperation
	ListDatabaseOperations(context.Context, *databasepb.ListDatabaseOperationsRequest, ...gax.CallOption) *OperationIterator
	ListBackupOperations(context.Context, *databasepb.ListBackupOperationsRequest, ...gax.CallOption) *OperationIterator
}
type DatabaseAdminClient struct {
	internalClient internalDatabaseAdminClient
	CallOptions    *DatabaseAdminCallOptions
	LROClient      *lroauto.OperationsClient
}

func (c *DatabaseAdminClient) gologoo__Close_224d23f03999a157bdd6d773efae62a9() error {
	return c.internalClient.Close()
}
func (c *DatabaseAdminClient) gologoo__setGoogleClientInfo_224d23f03999a157bdd6d773efae62a9(keyval ...string) {
	c.internalClient.setGoogleClientInfo(keyval...)
}
func (c *DatabaseAdminClient) gologoo__Connection_224d23f03999a157bdd6d773efae62a9() *grpc.ClientConn {
	return c.internalClient.Connection()
}
func (c *DatabaseAdminClient) gologoo__ListDatabases_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.ListDatabasesRequest, opts ...gax.CallOption) *DatabaseIterator {
	return c.internalClient.ListDatabases(ctx, req, opts...)
}
func (c *DatabaseAdminClient) gologoo__CreateDatabase_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.CreateDatabaseRequest, opts ...gax.CallOption) (*CreateDatabaseOperation, error) {
	return c.internalClient.CreateDatabase(ctx, req, opts...)
}
func (c *DatabaseAdminClient) gologoo__CreateDatabaseOperation_224d23f03999a157bdd6d773efae62a9(name string) *CreateDatabaseOperation {
	return c.internalClient.CreateDatabaseOperation(name)
}
func (c *DatabaseAdminClient) gologoo__GetDatabase_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.GetDatabaseRequest, opts ...gax.CallOption) (*databasepb.Database, error) {
	return c.internalClient.GetDatabase(ctx, req, opts...)
}
func (c *DatabaseAdminClient) gologoo__UpdateDatabaseDdl_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.UpdateDatabaseDdlRequest, opts ...gax.CallOption) (*UpdateDatabaseDdlOperation, error) {
	return c.internalClient.UpdateDatabaseDdl(ctx, req, opts...)
}
func (c *DatabaseAdminClient) gologoo__UpdateDatabaseDdlOperation_224d23f03999a157bdd6d773efae62a9(name string) *UpdateDatabaseDdlOperation {
	return c.internalClient.UpdateDatabaseDdlOperation(name)
}
func (c *DatabaseAdminClient) gologoo__DropDatabase_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.DropDatabaseRequest, opts ...gax.CallOption) error {
	return c.internalClient.DropDatabase(ctx, req, opts...)
}
func (c *DatabaseAdminClient) gologoo__GetDatabaseDdl_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.GetDatabaseDdlRequest, opts ...gax.CallOption) (*databasepb.GetDatabaseDdlResponse, error) {
	return c.internalClient.GetDatabaseDdl(ctx, req, opts...)
}
func (c *DatabaseAdminClient) gologoo__SetIamPolicy_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *iampb.SetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	return c.internalClient.SetIamPolicy(ctx, req, opts...)
}
func (c *DatabaseAdminClient) gologoo__GetIamPolicy_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *iampb.GetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	return c.internalClient.GetIamPolicy(ctx, req, opts...)
}
func (c *DatabaseAdminClient) gologoo__TestIamPermissions_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *iampb.TestIamPermissionsRequest, opts ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error) {
	return c.internalClient.TestIamPermissions(ctx, req, opts...)
}
func (c *DatabaseAdminClient) gologoo__CreateBackup_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.CreateBackupRequest, opts ...gax.CallOption) (*CreateBackupOperation, error) {
	return c.internalClient.CreateBackup(ctx, req, opts...)
}
func (c *DatabaseAdminClient) gologoo__CreateBackupOperation_224d23f03999a157bdd6d773efae62a9(name string) *CreateBackupOperation {
	return c.internalClient.CreateBackupOperation(name)
}
func (c *DatabaseAdminClient) gologoo__GetBackup_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.GetBackupRequest, opts ...gax.CallOption) (*databasepb.Backup, error) {
	return c.internalClient.GetBackup(ctx, req, opts...)
}
func (c *DatabaseAdminClient) gologoo__UpdateBackup_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.UpdateBackupRequest, opts ...gax.CallOption) (*databasepb.Backup, error) {
	return c.internalClient.UpdateBackup(ctx, req, opts...)
}
func (c *DatabaseAdminClient) gologoo__DeleteBackup_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.DeleteBackupRequest, opts ...gax.CallOption) error {
	return c.internalClient.DeleteBackup(ctx, req, opts...)
}
func (c *DatabaseAdminClient) gologoo__ListBackups_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.ListBackupsRequest, opts ...gax.CallOption) *BackupIterator {
	return c.internalClient.ListBackups(ctx, req, opts...)
}
func (c *DatabaseAdminClient) gologoo__RestoreDatabase_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.RestoreDatabaseRequest, opts ...gax.CallOption) (*RestoreDatabaseOperation, error) {
	return c.internalClient.RestoreDatabase(ctx, req, opts...)
}
func (c *DatabaseAdminClient) gologoo__RestoreDatabaseOperation_224d23f03999a157bdd6d773efae62a9(name string) *RestoreDatabaseOperation {
	return c.internalClient.RestoreDatabaseOperation(name)
}
func (c *DatabaseAdminClient) gologoo__ListDatabaseOperations_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.ListDatabaseOperationsRequest, opts ...gax.CallOption) *OperationIterator {
	return c.internalClient.ListDatabaseOperations(ctx, req, opts...)
}
func (c *DatabaseAdminClient) gologoo__ListBackupOperations_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.ListBackupOperationsRequest, opts ...gax.CallOption) *OperationIterator {
	return c.internalClient.ListBackupOperations(ctx, req, opts...)
}

type databaseAdminGRPCClient struct {
	connPool            gtransport.ConnPool
	disableDeadlines    bool
	CallOptions         **DatabaseAdminCallOptions
	databaseAdminClient databasepb.DatabaseAdminClient
	LROClient           **lroauto.OperationsClient
	xGoogMetadata       metadata.MD
}

func gologoo__NewDatabaseAdminClient_224d23f03999a157bdd6d773efae62a9(ctx context.Context, opts ...option.ClientOption) (*DatabaseAdminClient, error) {
	clientOpts := defaultDatabaseAdminGRPCClientOptions()
	if newDatabaseAdminClientHook != nil {
		hookOpts, err := newDatabaseAdminClientHook(ctx, clientHookParams{})
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
	client := DatabaseAdminClient{CallOptions: defaultDatabaseAdminCallOptions()}
	c := &databaseAdminGRPCClient{connPool: connPool, disableDeadlines: disableDeadlines, databaseAdminClient: databasepb.NewDatabaseAdminClient(connPool), CallOptions: &client.CallOptions}
	c.setGoogleClientInfo()
	client.internalClient = c
	client.LROClient, err = lroauto.NewOperationsClient(ctx, gtransport.WithConnPool(connPool))
	if err != nil {
		return nil, err
	}
	c.LROClient = &client.LROClient
	return &client, nil
}
func (c *databaseAdminGRPCClient) gologoo__Connection_224d23f03999a157bdd6d773efae62a9() *grpc.ClientConn {
	return c.connPool.Conn()
}
func (c *databaseAdminGRPCClient) gologoo__setGoogleClientInfo_224d23f03999a157bdd6d773efae62a9(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", getVersionClient(), "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}
func (c *databaseAdminGRPCClient) gologoo__Close_224d23f03999a157bdd6d773efae62a9() error {
	return c.connPool.Close()
}
func (c *databaseAdminGRPCClient) gologoo__ListDatabases_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.ListDatabasesRequest, opts ...gax.CallOption) *DatabaseIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).ListDatabases[0:len((*c.CallOptions).ListDatabases):len((*c.CallOptions).ListDatabases)], opts...)
	it := &DatabaseIterator{}
	req = proto.Clone(req).(*databasepb.ListDatabasesRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*databasepb.Database, string, error) {
		resp := &databasepb.ListDatabasesResponse{}
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
			resp, err = c.databaseAdminClient.ListDatabases(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		it.Response = resp
		return resp.GetDatabases(), resp.GetNextPageToken(), nil
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
func (c *databaseAdminGRPCClient) gologoo__CreateDatabase_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.CreateDatabaseRequest, opts ...gax.CallOption) (*CreateDatabaseOperation, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 3600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).CreateDatabase[0:len((*c.CallOptions).CreateDatabase):len((*c.CallOptions).CreateDatabase)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.databaseAdminClient.CreateDatabase(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &CreateDatabaseOperation{lro: longrunning.InternalNewOperation(*c.LROClient, resp)}, nil
}
func (c *databaseAdminGRPCClient) gologoo__GetDatabase_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.GetDatabaseRequest, opts ...gax.CallOption) (*databasepb.Database, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 3600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).GetDatabase[0:len((*c.CallOptions).GetDatabase):len((*c.CallOptions).GetDatabase)], opts...)
	var resp *databasepb.Database
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.databaseAdminClient.GetDatabase(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *databaseAdminGRPCClient) gologoo__UpdateDatabaseDdl_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.UpdateDatabaseDdlRequest, opts ...gax.CallOption) (*UpdateDatabaseDdlOperation, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 3600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "database", url.QueryEscape(req.GetDatabase())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).UpdateDatabaseDdl[0:len((*c.CallOptions).UpdateDatabaseDdl):len((*c.CallOptions).UpdateDatabaseDdl)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.databaseAdminClient.UpdateDatabaseDdl(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &UpdateDatabaseDdlOperation{lro: longrunning.InternalNewOperation(*c.LROClient, resp)}, nil
}
func (c *databaseAdminGRPCClient) gologoo__DropDatabase_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.DropDatabaseRequest, opts ...gax.CallOption) error {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 3600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "database", url.QueryEscape(req.GetDatabase())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).DropDatabase[0:len((*c.CallOptions).DropDatabase):len((*c.CallOptions).DropDatabase)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.databaseAdminClient.DropDatabase(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}
func (c *databaseAdminGRPCClient) gologoo__GetDatabaseDdl_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.GetDatabaseDdlRequest, opts ...gax.CallOption) (*databasepb.GetDatabaseDdlResponse, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 3600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "database", url.QueryEscape(req.GetDatabase())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).GetDatabaseDdl[0:len((*c.CallOptions).GetDatabaseDdl):len((*c.CallOptions).GetDatabaseDdl)], opts...)
	var resp *databasepb.GetDatabaseDdlResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.databaseAdminClient.GetDatabaseDdl(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *databaseAdminGRPCClient) gologoo__SetIamPolicy_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *iampb.SetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 30000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "resource", url.QueryEscape(req.GetResource())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).SetIamPolicy[0:len((*c.CallOptions).SetIamPolicy):len((*c.CallOptions).SetIamPolicy)], opts...)
	var resp *iampb.Policy
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.databaseAdminClient.SetIamPolicy(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *databaseAdminGRPCClient) gologoo__GetIamPolicy_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *iampb.GetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 30000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "resource", url.QueryEscape(req.GetResource())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).GetIamPolicy[0:len((*c.CallOptions).GetIamPolicy):len((*c.CallOptions).GetIamPolicy)], opts...)
	var resp *iampb.Policy
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.databaseAdminClient.GetIamPolicy(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *databaseAdminGRPCClient) gologoo__TestIamPermissions_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *iampb.TestIamPermissionsRequest, opts ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 30000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "resource", url.QueryEscape(req.GetResource())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).TestIamPermissions[0:len((*c.CallOptions).TestIamPermissions):len((*c.CallOptions).TestIamPermissions)], opts...)
	var resp *iampb.TestIamPermissionsResponse
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.databaseAdminClient.TestIamPermissions(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *databaseAdminGRPCClient) gologoo__CreateBackup_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.CreateBackupRequest, opts ...gax.CallOption) (*CreateBackupOperation, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 3600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).CreateBackup[0:len((*c.CallOptions).CreateBackup):len((*c.CallOptions).CreateBackup)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.databaseAdminClient.CreateBackup(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &CreateBackupOperation{lro: longrunning.InternalNewOperation(*c.LROClient, resp)}, nil
}
func (c *databaseAdminGRPCClient) gologoo__GetBackup_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.GetBackupRequest, opts ...gax.CallOption) (*databasepb.Backup, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 3600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).GetBackup[0:len((*c.CallOptions).GetBackup):len((*c.CallOptions).GetBackup)], opts...)
	var resp *databasepb.Backup
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.databaseAdminClient.GetBackup(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *databaseAdminGRPCClient) gologoo__UpdateBackup_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.UpdateBackupRequest, opts ...gax.CallOption) (*databasepb.Backup, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 3600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "backup.name", url.QueryEscape(req.GetBackup().GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).UpdateBackup[0:len((*c.CallOptions).UpdateBackup):len((*c.CallOptions).UpdateBackup)], opts...)
	var resp *databasepb.Backup
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.databaseAdminClient.UpdateBackup(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *databaseAdminGRPCClient) gologoo__DeleteBackup_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.DeleteBackupRequest, opts ...gax.CallOption) error {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 3600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).DeleteBackup[0:len((*c.CallOptions).DeleteBackup):len((*c.CallOptions).DeleteBackup)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.databaseAdminClient.DeleteBackup(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}
func (c *databaseAdminGRPCClient) gologoo__ListBackups_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.ListBackupsRequest, opts ...gax.CallOption) *BackupIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).ListBackups[0:len((*c.CallOptions).ListBackups):len((*c.CallOptions).ListBackups)], opts...)
	it := &BackupIterator{}
	req = proto.Clone(req).(*databasepb.ListBackupsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*databasepb.Backup, string, error) {
		resp := &databasepb.ListBackupsResponse{}
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
			resp, err = c.databaseAdminClient.ListBackups(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		it.Response = resp
		return resp.GetBackups(), resp.GetNextPageToken(), nil
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
func (c *databaseAdminGRPCClient) gologoo__RestoreDatabase_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.RestoreDatabaseRequest, opts ...gax.CallOption) (*RestoreDatabaseOperation, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 3600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).RestoreDatabase[0:len((*c.CallOptions).RestoreDatabase):len((*c.CallOptions).RestoreDatabase)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.databaseAdminClient.RestoreDatabase(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &RestoreDatabaseOperation{lro: longrunning.InternalNewOperation(*c.LROClient, resp)}, nil
}
func (c *databaseAdminGRPCClient) gologoo__ListDatabaseOperations_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.ListDatabaseOperationsRequest, opts ...gax.CallOption) *OperationIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).ListDatabaseOperations[0:len((*c.CallOptions).ListDatabaseOperations):len((*c.CallOptions).ListDatabaseOperations)], opts...)
	it := &OperationIterator{}
	req = proto.Clone(req).(*databasepb.ListDatabaseOperationsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*longrunningpb.Operation, string, error) {
		resp := &databasepb.ListDatabaseOperationsResponse{}
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
			resp, err = c.databaseAdminClient.ListDatabaseOperations(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		it.Response = resp
		return resp.GetOperations(), resp.GetNextPageToken(), nil
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
func (c *databaseAdminGRPCClient) gologoo__ListBackupOperations_224d23f03999a157bdd6d773efae62a9(ctx context.Context, req *databasepb.ListBackupOperationsRequest, opts ...gax.CallOption) *OperationIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).ListBackupOperations[0:len((*c.CallOptions).ListBackupOperations):len((*c.CallOptions).ListBackupOperations)], opts...)
	it := &OperationIterator{}
	req = proto.Clone(req).(*databasepb.ListBackupOperationsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*longrunningpb.Operation, string, error) {
		resp := &databasepb.ListBackupOperationsResponse{}
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
			resp, err = c.databaseAdminClient.ListBackupOperations(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		it.Response = resp
		return resp.GetOperations(), resp.GetNextPageToken(), nil
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

type CreateBackupOperation struct {
	lro *longrunning.Operation
}

func (c *databaseAdminGRPCClient) gologoo__CreateBackupOperation_224d23f03999a157bdd6d773efae62a9(name string) *CreateBackupOperation {
	return &CreateBackupOperation{lro: longrunning.InternalNewOperation(*c.LROClient, &longrunningpb.Operation{Name: name})}
}
func (op *CreateBackupOperation) gologoo__Wait_224d23f03999a157bdd6d773efae62a9(ctx context.Context, opts ...gax.CallOption) (*databasepb.Backup, error) {
	var resp databasepb.Backup
	if err := op.lro.WaitWithInterval(ctx, &resp, time.Minute, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}
func (op *CreateBackupOperation) gologoo__Poll_224d23f03999a157bdd6d773efae62a9(ctx context.Context, opts ...gax.CallOption) (*databasepb.Backup, error) {
	var resp databasepb.Backup
	if err := op.lro.Poll(ctx, &resp, opts...); err != nil {
		return nil, err
	}
	if !op.Done() {
		return nil, nil
	}
	return &resp, nil
}
func (op *CreateBackupOperation) gologoo__Metadata_224d23f03999a157bdd6d773efae62a9() (*databasepb.CreateBackupMetadata, error) {
	var meta databasepb.CreateBackupMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}
func (op *CreateBackupOperation) gologoo__Done_224d23f03999a157bdd6d773efae62a9() bool {
	return op.lro.Done()
}
func (op *CreateBackupOperation) gologoo__Name_224d23f03999a157bdd6d773efae62a9() string {
	return op.lro.Name()
}

type CreateDatabaseOperation struct {
	lro *longrunning.Operation
}

func (c *databaseAdminGRPCClient) gologoo__CreateDatabaseOperation_224d23f03999a157bdd6d773efae62a9(name string) *CreateDatabaseOperation {
	return &CreateDatabaseOperation{lro: longrunning.InternalNewOperation(*c.LROClient, &longrunningpb.Operation{Name: name})}
}
func (op *CreateDatabaseOperation) gologoo__Wait_224d23f03999a157bdd6d773efae62a9(ctx context.Context, opts ...gax.CallOption) (*databasepb.Database, error) {
	var resp databasepb.Database
	if err := op.lro.WaitWithInterval(ctx, &resp, time.Minute, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}
func (op *CreateDatabaseOperation) gologoo__Poll_224d23f03999a157bdd6d773efae62a9(ctx context.Context, opts ...gax.CallOption) (*databasepb.Database, error) {
	var resp databasepb.Database
	if err := op.lro.Poll(ctx, &resp, opts...); err != nil {
		return nil, err
	}
	if !op.Done() {
		return nil, nil
	}
	return &resp, nil
}
func (op *CreateDatabaseOperation) gologoo__Metadata_224d23f03999a157bdd6d773efae62a9() (*databasepb.CreateDatabaseMetadata, error) {
	var meta databasepb.CreateDatabaseMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}
func (op *CreateDatabaseOperation) gologoo__Done_224d23f03999a157bdd6d773efae62a9() bool {
	return op.lro.Done()
}
func (op *CreateDatabaseOperation) gologoo__Name_224d23f03999a157bdd6d773efae62a9() string {
	return op.lro.Name()
}

type RestoreDatabaseOperation struct {
	lro *longrunning.Operation
}

func (c *databaseAdminGRPCClient) gologoo__RestoreDatabaseOperation_224d23f03999a157bdd6d773efae62a9(name string) *RestoreDatabaseOperation {
	return &RestoreDatabaseOperation{lro: longrunning.InternalNewOperation(*c.LROClient, &longrunningpb.Operation{Name: name})}
}
func (op *RestoreDatabaseOperation) gologoo__Wait_224d23f03999a157bdd6d773efae62a9(ctx context.Context, opts ...gax.CallOption) (*databasepb.Database, error) {
	var resp databasepb.Database
	if err := op.lro.WaitWithInterval(ctx, &resp, time.Minute, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}
func (op *RestoreDatabaseOperation) gologoo__Poll_224d23f03999a157bdd6d773efae62a9(ctx context.Context, opts ...gax.CallOption) (*databasepb.Database, error) {
	var resp databasepb.Database
	if err := op.lro.Poll(ctx, &resp, opts...); err != nil {
		return nil, err
	}
	if !op.Done() {
		return nil, nil
	}
	return &resp, nil
}
func (op *RestoreDatabaseOperation) gologoo__Metadata_224d23f03999a157bdd6d773efae62a9() (*databasepb.RestoreDatabaseMetadata, error) {
	var meta databasepb.RestoreDatabaseMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}
func (op *RestoreDatabaseOperation) gologoo__Done_224d23f03999a157bdd6d773efae62a9() bool {
	return op.lro.Done()
}
func (op *RestoreDatabaseOperation) gologoo__Name_224d23f03999a157bdd6d773efae62a9() string {
	return op.lro.Name()
}

type UpdateDatabaseDdlOperation struct {
	lro *longrunning.Operation
}

func (c *databaseAdminGRPCClient) gologoo__UpdateDatabaseDdlOperation_224d23f03999a157bdd6d773efae62a9(name string) *UpdateDatabaseDdlOperation {
	return &UpdateDatabaseDdlOperation{lro: longrunning.InternalNewOperation(*c.LROClient, &longrunningpb.Operation{Name: name})}
}
func (op *UpdateDatabaseDdlOperation) gologoo__Wait_224d23f03999a157bdd6d773efae62a9(ctx context.Context, opts ...gax.CallOption) error {
	return op.lro.WaitWithInterval(ctx, nil, time.Minute, opts...)
}
func (op *UpdateDatabaseDdlOperation) gologoo__Poll_224d23f03999a157bdd6d773efae62a9(ctx context.Context, opts ...gax.CallOption) error {
	return op.lro.Poll(ctx, nil, opts...)
}
func (op *UpdateDatabaseDdlOperation) gologoo__Metadata_224d23f03999a157bdd6d773efae62a9() (*databasepb.UpdateDatabaseDdlMetadata, error) {
	var meta databasepb.UpdateDatabaseDdlMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}
func (op *UpdateDatabaseDdlOperation) gologoo__Done_224d23f03999a157bdd6d773efae62a9() bool {
	return op.lro.Done()
}
func (op *UpdateDatabaseDdlOperation) gologoo__Name_224d23f03999a157bdd6d773efae62a9() string {
	return op.lro.Name()
}

type BackupIterator struct {
	items    []*databasepb.Backup
	pageInfo *iterator.PageInfo
	nextFunc func() error
	Response interface {
	}
	InternalFetch func(pageSize int, pageToken string) (results []*databasepb.Backup, nextPageToken string, err error)
}

func (it *BackupIterator) gologoo__PageInfo_224d23f03999a157bdd6d773efae62a9() *iterator.PageInfo {
	return it.pageInfo
}
func (it *BackupIterator) gologoo__Next_224d23f03999a157bdd6d773efae62a9() (*databasepb.Backup, error) {
	var item *databasepb.Backup
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}
func (it *BackupIterator) gologoo__bufLen_224d23f03999a157bdd6d773efae62a9() int {
	return len(it.items)
}
func (it *BackupIterator) gologoo__takeBuf_224d23f03999a157bdd6d773efae62a9() interface {
} {
	b := it.items
	it.items = nil
	return b
}

type DatabaseIterator struct {
	items    []*databasepb.Database
	pageInfo *iterator.PageInfo
	nextFunc func() error
	Response interface {
	}
	InternalFetch func(pageSize int, pageToken string) (results []*databasepb.Database, nextPageToken string, err error)
}

func (it *DatabaseIterator) gologoo__PageInfo_224d23f03999a157bdd6d773efae62a9() *iterator.PageInfo {
	return it.pageInfo
}
func (it *DatabaseIterator) gologoo__Next_224d23f03999a157bdd6d773efae62a9() (*databasepb.Database, error) {
	var item *databasepb.Database
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}
func (it *DatabaseIterator) gologoo__bufLen_224d23f03999a157bdd6d773efae62a9() int {
	return len(it.items)
}
func (it *DatabaseIterator) gologoo__takeBuf_224d23f03999a157bdd6d773efae62a9() interface {
} {
	b := it.items
	it.items = nil
	return b
}

type OperationIterator struct {
	items    []*longrunningpb.Operation
	pageInfo *iterator.PageInfo
	nextFunc func() error
	Response interface {
	}
	InternalFetch func(pageSize int, pageToken string) (results []*longrunningpb.Operation, nextPageToken string, err error)
}

func (it *OperationIterator) gologoo__PageInfo_224d23f03999a157bdd6d773efae62a9() *iterator.PageInfo {
	return it.pageInfo
}
func (it *OperationIterator) gologoo__Next_224d23f03999a157bdd6d773efae62a9() (*longrunningpb.Operation, error) {
	var item *longrunningpb.Operation
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}
func (it *OperationIterator) gologoo__bufLen_224d23f03999a157bdd6d773efae62a9() int {
	return len(it.items)
}
func (it *OperationIterator) gologoo__takeBuf_224d23f03999a157bdd6d773efae62a9() interface {
} {
	b := it.items
	it.items = nil
	return b
}
func defaultDatabaseAdminGRPCClientOptions() []option.ClientOption {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__defaultDatabaseAdminGRPCClientOptions_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := gologoo__defaultDatabaseAdminGRPCClientOptions_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func defaultDatabaseAdminCallOptions() *DatabaseAdminCallOptions {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__defaultDatabaseAdminCallOptions_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := gologoo__defaultDatabaseAdminCallOptions_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *DatabaseAdminClient) Close() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Close_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__Close_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *DatabaseAdminClient) setGoogleClientInfo(keyval ...string) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__setGoogleClientInfo_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v\n", keyval)
	c.gologoo__setGoogleClientInfo_224d23f03999a157bdd6d773efae62a9(keyval...)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (c *DatabaseAdminClient) Connection() *grpc.ClientConn {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Connection_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__Connection_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *DatabaseAdminClient) ListDatabases(ctx context.Context, req *databasepb.ListDatabasesRequest, opts ...gax.CallOption) *DatabaseIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ListDatabases_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__ListDatabases_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *DatabaseAdminClient) CreateDatabase(ctx context.Context, req *databasepb.CreateDatabaseRequest, opts ...gax.CallOption) (*CreateDatabaseOperation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateDatabase_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__CreateDatabase_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *DatabaseAdminClient) CreateDatabaseOperation(name string) *CreateDatabaseOperation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateDatabaseOperation_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v\n", name)
	r0 := c.gologoo__CreateDatabaseOperation_224d23f03999a157bdd6d773efae62a9(name)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *DatabaseAdminClient) GetDatabase(ctx context.Context, req *databasepb.GetDatabaseRequest, opts ...gax.CallOption) (*databasepb.Database, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetDatabase_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__GetDatabase_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *DatabaseAdminClient) UpdateDatabaseDdl(ctx context.Context, req *databasepb.UpdateDatabaseDdlRequest, opts ...gax.CallOption) (*UpdateDatabaseDdlOperation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UpdateDatabaseDdl_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__UpdateDatabaseDdl_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *DatabaseAdminClient) UpdateDatabaseDdlOperation(name string) *UpdateDatabaseDdlOperation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UpdateDatabaseDdlOperation_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v\n", name)
	r0 := c.gologoo__UpdateDatabaseDdlOperation_224d23f03999a157bdd6d773efae62a9(name)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *DatabaseAdminClient) DropDatabase(ctx context.Context, req *databasepb.DropDatabaseRequest, opts ...gax.CallOption) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DropDatabase_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__DropDatabase_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *DatabaseAdminClient) GetDatabaseDdl(ctx context.Context, req *databasepb.GetDatabaseDdlRequest, opts ...gax.CallOption) (*databasepb.GetDatabaseDdlResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetDatabaseDdl_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__GetDatabaseDdl_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *DatabaseAdminClient) SetIamPolicy(ctx context.Context, req *iampb.SetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SetIamPolicy_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__SetIamPolicy_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *DatabaseAdminClient) GetIamPolicy(ctx context.Context, req *iampb.GetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetIamPolicy_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__GetIamPolicy_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *DatabaseAdminClient) TestIamPermissions(ctx context.Context, req *iampb.TestIamPermissionsRequest, opts ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__TestIamPermissions_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__TestIamPermissions_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *DatabaseAdminClient) CreateBackup(ctx context.Context, req *databasepb.CreateBackupRequest, opts ...gax.CallOption) (*CreateBackupOperation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateBackup_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__CreateBackup_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *DatabaseAdminClient) CreateBackupOperation(name string) *CreateBackupOperation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateBackupOperation_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v\n", name)
	r0 := c.gologoo__CreateBackupOperation_224d23f03999a157bdd6d773efae62a9(name)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *DatabaseAdminClient) GetBackup(ctx context.Context, req *databasepb.GetBackupRequest, opts ...gax.CallOption) (*databasepb.Backup, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetBackup_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__GetBackup_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *DatabaseAdminClient) UpdateBackup(ctx context.Context, req *databasepb.UpdateBackupRequest, opts ...gax.CallOption) (*databasepb.Backup, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UpdateBackup_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__UpdateBackup_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *DatabaseAdminClient) DeleteBackup(ctx context.Context, req *databasepb.DeleteBackupRequest, opts ...gax.CallOption) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DeleteBackup_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__DeleteBackup_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *DatabaseAdminClient) ListBackups(ctx context.Context, req *databasepb.ListBackupsRequest, opts ...gax.CallOption) *BackupIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ListBackups_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__ListBackups_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *DatabaseAdminClient) RestoreDatabase(ctx context.Context, req *databasepb.RestoreDatabaseRequest, opts ...gax.CallOption) (*RestoreDatabaseOperation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__RestoreDatabase_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__RestoreDatabase_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *DatabaseAdminClient) RestoreDatabaseOperation(name string) *RestoreDatabaseOperation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__RestoreDatabaseOperation_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v\n", name)
	r0 := c.gologoo__RestoreDatabaseOperation_224d23f03999a157bdd6d773efae62a9(name)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *DatabaseAdminClient) ListDatabaseOperations(ctx context.Context, req *databasepb.ListDatabaseOperationsRequest, opts ...gax.CallOption) *OperationIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ListDatabaseOperations_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__ListDatabaseOperations_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *DatabaseAdminClient) ListBackupOperations(ctx context.Context, req *databasepb.ListBackupOperationsRequest, opts ...gax.CallOption) *OperationIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ListBackupOperations_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__ListBackupOperations_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func NewDatabaseAdminClient(ctx context.Context, opts ...option.ClientOption) (*DatabaseAdminClient, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NewDatabaseAdminClient_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v\n", ctx, opts)
	r0, r1 := gologoo__NewDatabaseAdminClient_224d23f03999a157bdd6d773efae62a9(ctx, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *databaseAdminGRPCClient) Connection() *grpc.ClientConn {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Connection_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__Connection_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *databaseAdminGRPCClient) setGoogleClientInfo(keyval ...string) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__setGoogleClientInfo_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v\n", keyval)
	c.gologoo__setGoogleClientInfo_224d23f03999a157bdd6d773efae62a9(keyval...)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (c *databaseAdminGRPCClient) Close() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Close_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__Close_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *databaseAdminGRPCClient) ListDatabases(ctx context.Context, req *databasepb.ListDatabasesRequest, opts ...gax.CallOption) *DatabaseIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ListDatabases_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__ListDatabases_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *databaseAdminGRPCClient) CreateDatabase(ctx context.Context, req *databasepb.CreateDatabaseRequest, opts ...gax.CallOption) (*CreateDatabaseOperation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateDatabase_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__CreateDatabase_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *databaseAdminGRPCClient) GetDatabase(ctx context.Context, req *databasepb.GetDatabaseRequest, opts ...gax.CallOption) (*databasepb.Database, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetDatabase_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__GetDatabase_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *databaseAdminGRPCClient) UpdateDatabaseDdl(ctx context.Context, req *databasepb.UpdateDatabaseDdlRequest, opts ...gax.CallOption) (*UpdateDatabaseDdlOperation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UpdateDatabaseDdl_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__UpdateDatabaseDdl_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *databaseAdminGRPCClient) DropDatabase(ctx context.Context, req *databasepb.DropDatabaseRequest, opts ...gax.CallOption) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DropDatabase_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__DropDatabase_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *databaseAdminGRPCClient) GetDatabaseDdl(ctx context.Context, req *databasepb.GetDatabaseDdlRequest, opts ...gax.CallOption) (*databasepb.GetDatabaseDdlResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetDatabaseDdl_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__GetDatabaseDdl_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *databaseAdminGRPCClient) SetIamPolicy(ctx context.Context, req *iampb.SetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SetIamPolicy_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__SetIamPolicy_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *databaseAdminGRPCClient) GetIamPolicy(ctx context.Context, req *iampb.GetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetIamPolicy_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__GetIamPolicy_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *databaseAdminGRPCClient) TestIamPermissions(ctx context.Context, req *iampb.TestIamPermissionsRequest, opts ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__TestIamPermissions_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__TestIamPermissions_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *databaseAdminGRPCClient) CreateBackup(ctx context.Context, req *databasepb.CreateBackupRequest, opts ...gax.CallOption) (*CreateBackupOperation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateBackup_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__CreateBackup_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *databaseAdminGRPCClient) GetBackup(ctx context.Context, req *databasepb.GetBackupRequest, opts ...gax.CallOption) (*databasepb.Backup, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetBackup_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__GetBackup_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *databaseAdminGRPCClient) UpdateBackup(ctx context.Context, req *databasepb.UpdateBackupRequest, opts ...gax.CallOption) (*databasepb.Backup, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UpdateBackup_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__UpdateBackup_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *databaseAdminGRPCClient) DeleteBackup(ctx context.Context, req *databasepb.DeleteBackupRequest, opts ...gax.CallOption) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DeleteBackup_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__DeleteBackup_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *databaseAdminGRPCClient) ListBackups(ctx context.Context, req *databasepb.ListBackupsRequest, opts ...gax.CallOption) *BackupIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ListBackups_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__ListBackups_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *databaseAdminGRPCClient) RestoreDatabase(ctx context.Context, req *databasepb.RestoreDatabaseRequest, opts ...gax.CallOption) (*RestoreDatabaseOperation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__RestoreDatabase_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__RestoreDatabase_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *databaseAdminGRPCClient) ListDatabaseOperations(ctx context.Context, req *databasepb.ListDatabaseOperationsRequest, opts ...gax.CallOption) *OperationIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ListDatabaseOperations_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__ListDatabaseOperations_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *databaseAdminGRPCClient) ListBackupOperations(ctx context.Context, req *databasepb.ListBackupOperationsRequest, opts ...gax.CallOption) *OperationIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ListBackupOperations_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__ListBackupOperations_224d23f03999a157bdd6d773efae62a9(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *databaseAdminGRPCClient) CreateBackupOperation(name string) *CreateBackupOperation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateBackupOperation_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v\n", name)
	r0 := c.gologoo__CreateBackupOperation_224d23f03999a157bdd6d773efae62a9(name)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (op *CreateBackupOperation) Wait(ctx context.Context, opts ...gax.CallOption) (*databasepb.Backup, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Wait_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v\n", ctx, opts)
	r0, r1 := op.gologoo__Wait_224d23f03999a157bdd6d773efae62a9(ctx, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (op *CreateBackupOperation) Poll(ctx context.Context, opts ...gax.CallOption) (*databasepb.Backup, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Poll_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v\n", ctx, opts)
	r0, r1 := op.gologoo__Poll_224d23f03999a157bdd6d773efae62a9(ctx, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (op *CreateBackupOperation) Metadata() (*databasepb.CreateBackupMetadata, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Metadata_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0, r1 := op.gologoo__Metadata_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (op *CreateBackupOperation) Done() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Done_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := op.gologoo__Done_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (op *CreateBackupOperation) Name() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Name_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := op.gologoo__Name_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *databaseAdminGRPCClient) CreateDatabaseOperation(name string) *CreateDatabaseOperation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateDatabaseOperation_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v\n", name)
	r0 := c.gologoo__CreateDatabaseOperation_224d23f03999a157bdd6d773efae62a9(name)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (op *CreateDatabaseOperation) Wait(ctx context.Context, opts ...gax.CallOption) (*databasepb.Database, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Wait_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v\n", ctx, opts)
	r0, r1 := op.gologoo__Wait_224d23f03999a157bdd6d773efae62a9(ctx, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (op *CreateDatabaseOperation) Poll(ctx context.Context, opts ...gax.CallOption) (*databasepb.Database, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Poll_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v\n", ctx, opts)
	r0, r1 := op.gologoo__Poll_224d23f03999a157bdd6d773efae62a9(ctx, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (op *CreateDatabaseOperation) Metadata() (*databasepb.CreateDatabaseMetadata, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Metadata_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0, r1 := op.gologoo__Metadata_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (op *CreateDatabaseOperation) Done() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Done_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := op.gologoo__Done_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (op *CreateDatabaseOperation) Name() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Name_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := op.gologoo__Name_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *databaseAdminGRPCClient) RestoreDatabaseOperation(name string) *RestoreDatabaseOperation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__RestoreDatabaseOperation_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v\n", name)
	r0 := c.gologoo__RestoreDatabaseOperation_224d23f03999a157bdd6d773efae62a9(name)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (op *RestoreDatabaseOperation) Wait(ctx context.Context, opts ...gax.CallOption) (*databasepb.Database, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Wait_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v\n", ctx, opts)
	r0, r1 := op.gologoo__Wait_224d23f03999a157bdd6d773efae62a9(ctx, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (op *RestoreDatabaseOperation) Poll(ctx context.Context, opts ...gax.CallOption) (*databasepb.Database, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Poll_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v\n", ctx, opts)
	r0, r1 := op.gologoo__Poll_224d23f03999a157bdd6d773efae62a9(ctx, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (op *RestoreDatabaseOperation) Metadata() (*databasepb.RestoreDatabaseMetadata, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Metadata_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0, r1 := op.gologoo__Metadata_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (op *RestoreDatabaseOperation) Done() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Done_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := op.gologoo__Done_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (op *RestoreDatabaseOperation) Name() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Name_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := op.gologoo__Name_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *databaseAdminGRPCClient) UpdateDatabaseDdlOperation(name string) *UpdateDatabaseDdlOperation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UpdateDatabaseDdlOperation_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v\n", name)
	r0 := c.gologoo__UpdateDatabaseDdlOperation_224d23f03999a157bdd6d773efae62a9(name)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (op *UpdateDatabaseDdlOperation) Wait(ctx context.Context, opts ...gax.CallOption) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Wait_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v\n", ctx, opts)
	r0 := op.gologoo__Wait_224d23f03999a157bdd6d773efae62a9(ctx, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (op *UpdateDatabaseDdlOperation) Poll(ctx context.Context, opts ...gax.CallOption) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Poll_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : %v %v\n", ctx, opts)
	r0 := op.gologoo__Poll_224d23f03999a157bdd6d773efae62a9(ctx, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (op *UpdateDatabaseDdlOperation) Metadata() (*databasepb.UpdateDatabaseDdlMetadata, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Metadata_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0, r1 := op.gologoo__Metadata_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (op *UpdateDatabaseDdlOperation) Done() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Done_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := op.gologoo__Done_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (op *UpdateDatabaseDdlOperation) Name() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Name_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := op.gologoo__Name_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (it *BackupIterator) PageInfo() *iterator.PageInfo {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PageInfo_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__PageInfo_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (it *BackupIterator) Next() (*databasepb.Backup, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Next_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0, r1 := it.gologoo__Next_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (it *BackupIterator) bufLen() int {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__bufLen_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__bufLen_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (it *BackupIterator) takeBuf() interface {
} {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__takeBuf_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__takeBuf_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (it *DatabaseIterator) PageInfo() *iterator.PageInfo {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PageInfo_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__PageInfo_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (it *DatabaseIterator) Next() (*databasepb.Database, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Next_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0, r1 := it.gologoo__Next_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (it *DatabaseIterator) bufLen() int {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__bufLen_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__bufLen_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (it *DatabaseIterator) takeBuf() interface {
} {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__takeBuf_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__takeBuf_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (it *OperationIterator) PageInfo() *iterator.PageInfo {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PageInfo_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__PageInfo_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (it *OperationIterator) Next() (*longrunningpb.Operation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Next_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0, r1 := it.gologoo__Next_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (it *OperationIterator) bufLen() int {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__bufLen_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__bufLen_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (it *OperationIterator) takeBuf() interface {
} {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__takeBuf_224d23f03999a157bdd6d773efae62a9")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__takeBuf_224d23f03999a157bdd6d773efae62a9()
	log.Printf("Output: %v\n", r0)
	return r0
}
