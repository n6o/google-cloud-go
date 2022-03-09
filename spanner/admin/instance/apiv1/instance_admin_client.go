package instance

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
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
	"log"
)

var newInstanceAdminClientHook clientHook

type InstanceAdminCallOptions struct {
	ListInstanceConfigs []gax.CallOption
	GetInstanceConfig   []gax.CallOption
	ListInstances       []gax.CallOption
	GetInstance         []gax.CallOption
	CreateInstance      []gax.CallOption
	UpdateInstance      []gax.CallOption
	DeleteInstance      []gax.CallOption
	SetIamPolicy        []gax.CallOption
	GetIamPolicy        []gax.CallOption
	TestIamPermissions  []gax.CallOption
}

func gologoo__defaultInstanceAdminGRPCClientOptions_9d780f59578540b8e87027995d65246b() []option.ClientOption {
	return []option.ClientOption{internaloption.WithDefaultEndpoint("spanner.googleapis.com:443"), internaloption.WithDefaultMTLSEndpoint("spanner.mtls.googleapis.com:443"), internaloption.WithDefaultAudience("https://spanner.googleapis.com/"), internaloption.WithDefaultScopes(DefaultAuthScopes()...), internaloption.EnableJwtWithScope(), option.WithGRPCDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)))}
}
func gologoo__defaultInstanceAdminCallOptions_9d780f59578540b8e87027995d65246b() *InstanceAdminCallOptions {
	return &InstanceAdminCallOptions{ListInstanceConfigs: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, GetInstanceConfig: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, ListInstances: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, GetInstance: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, CreateInstance: []gax.CallOption{}, UpdateInstance: []gax.CallOption{}, DeleteInstance: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, SetIamPolicy: []gax.CallOption{}, GetIamPolicy: []gax.CallOption{gax.WithRetry(func() gax.Retryer {
		return gax.OnCodes([]codes.Code{codes.Unavailable, codes.DeadlineExceeded}, gax.Backoff{Initial: 1000 * time.Millisecond, Max: 32000 * time.Millisecond, Multiplier: 1.30})
	})}, TestIamPermissions: []gax.CallOption{}}
}

type internalInstanceAdminClient interface {
	Close() error
	setGoogleClientInfo(...string)
	Connection() *grpc.ClientConn
	ListInstanceConfigs(context.Context, *instancepb.ListInstanceConfigsRequest, ...gax.CallOption) *InstanceConfigIterator
	GetInstanceConfig(context.Context, *instancepb.GetInstanceConfigRequest, ...gax.CallOption) (*instancepb.InstanceConfig, error)
	ListInstances(context.Context, *instancepb.ListInstancesRequest, ...gax.CallOption) *InstanceIterator
	GetInstance(context.Context, *instancepb.GetInstanceRequest, ...gax.CallOption) (*instancepb.Instance, error)
	CreateInstance(context.Context, *instancepb.CreateInstanceRequest, ...gax.CallOption) (*CreateInstanceOperation, error)
	CreateInstanceOperation(name string) *CreateInstanceOperation
	UpdateInstance(context.Context, *instancepb.UpdateInstanceRequest, ...gax.CallOption) (*UpdateInstanceOperation, error)
	UpdateInstanceOperation(name string) *UpdateInstanceOperation
	DeleteInstance(context.Context, *instancepb.DeleteInstanceRequest, ...gax.CallOption) error
	SetIamPolicy(context.Context, *iampb.SetIamPolicyRequest, ...gax.CallOption) (*iampb.Policy, error)
	GetIamPolicy(context.Context, *iampb.GetIamPolicyRequest, ...gax.CallOption) (*iampb.Policy, error)
	TestIamPermissions(context.Context, *iampb.TestIamPermissionsRequest, ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error)
}
type InstanceAdminClient struct {
	internalClient internalInstanceAdminClient
	CallOptions    *InstanceAdminCallOptions
	LROClient      *lroauto.OperationsClient
}

func (c *InstanceAdminClient) gologoo__Close_9d780f59578540b8e87027995d65246b() error {
	return c.internalClient.Close()
}
func (c *InstanceAdminClient) gologoo__setGoogleClientInfo_9d780f59578540b8e87027995d65246b(keyval ...string) {
	c.internalClient.setGoogleClientInfo(keyval...)
}
func (c *InstanceAdminClient) gologoo__Connection_9d780f59578540b8e87027995d65246b() *grpc.ClientConn {
	return c.internalClient.Connection()
}
func (c *InstanceAdminClient) gologoo__ListInstanceConfigs_9d780f59578540b8e87027995d65246b(ctx context.Context, req *instancepb.ListInstanceConfigsRequest, opts ...gax.CallOption) *InstanceConfigIterator {
	return c.internalClient.ListInstanceConfigs(ctx, req, opts...)
}
func (c *InstanceAdminClient) gologoo__GetInstanceConfig_9d780f59578540b8e87027995d65246b(ctx context.Context, req *instancepb.GetInstanceConfigRequest, opts ...gax.CallOption) (*instancepb.InstanceConfig, error) {
	return c.internalClient.GetInstanceConfig(ctx, req, opts...)
}
func (c *InstanceAdminClient) gologoo__ListInstances_9d780f59578540b8e87027995d65246b(ctx context.Context, req *instancepb.ListInstancesRequest, opts ...gax.CallOption) *InstanceIterator {
	return c.internalClient.ListInstances(ctx, req, opts...)
}
func (c *InstanceAdminClient) gologoo__GetInstance_9d780f59578540b8e87027995d65246b(ctx context.Context, req *instancepb.GetInstanceRequest, opts ...gax.CallOption) (*instancepb.Instance, error) {
	return c.internalClient.GetInstance(ctx, req, opts...)
}
func (c *InstanceAdminClient) gologoo__CreateInstance_9d780f59578540b8e87027995d65246b(ctx context.Context, req *instancepb.CreateInstanceRequest, opts ...gax.CallOption) (*CreateInstanceOperation, error) {
	return c.internalClient.CreateInstance(ctx, req, opts...)
}
func (c *InstanceAdminClient) gologoo__CreateInstanceOperation_9d780f59578540b8e87027995d65246b(name string) *CreateInstanceOperation {
	return c.internalClient.CreateInstanceOperation(name)
}
func (c *InstanceAdminClient) gologoo__UpdateInstance_9d780f59578540b8e87027995d65246b(ctx context.Context, req *instancepb.UpdateInstanceRequest, opts ...gax.CallOption) (*UpdateInstanceOperation, error) {
	return c.internalClient.UpdateInstance(ctx, req, opts...)
}
func (c *InstanceAdminClient) gologoo__UpdateInstanceOperation_9d780f59578540b8e87027995d65246b(name string) *UpdateInstanceOperation {
	return c.internalClient.UpdateInstanceOperation(name)
}
func (c *InstanceAdminClient) gologoo__DeleteInstance_9d780f59578540b8e87027995d65246b(ctx context.Context, req *instancepb.DeleteInstanceRequest, opts ...gax.CallOption) error {
	return c.internalClient.DeleteInstance(ctx, req, opts...)
}
func (c *InstanceAdminClient) gologoo__SetIamPolicy_9d780f59578540b8e87027995d65246b(ctx context.Context, req *iampb.SetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	return c.internalClient.SetIamPolicy(ctx, req, opts...)
}
func (c *InstanceAdminClient) gologoo__GetIamPolicy_9d780f59578540b8e87027995d65246b(ctx context.Context, req *iampb.GetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	return c.internalClient.GetIamPolicy(ctx, req, opts...)
}
func (c *InstanceAdminClient) gologoo__TestIamPermissions_9d780f59578540b8e87027995d65246b(ctx context.Context, req *iampb.TestIamPermissionsRequest, opts ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error) {
	return c.internalClient.TestIamPermissions(ctx, req, opts...)
}

type instanceAdminGRPCClient struct {
	connPool            gtransport.ConnPool
	disableDeadlines    bool
	CallOptions         **InstanceAdminCallOptions
	instanceAdminClient instancepb.InstanceAdminClient
	LROClient           **lroauto.OperationsClient
	xGoogMetadata       metadata.MD
}

func gologoo__NewInstanceAdminClient_9d780f59578540b8e87027995d65246b(ctx context.Context, opts ...option.ClientOption) (*InstanceAdminClient, error) {
	clientOpts := defaultInstanceAdminGRPCClientOptions()
	if newInstanceAdminClientHook != nil {
		hookOpts, err := newInstanceAdminClientHook(ctx, clientHookParams{})
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
	client := InstanceAdminClient{CallOptions: defaultInstanceAdminCallOptions()}
	c := &instanceAdminGRPCClient{connPool: connPool, disableDeadlines: disableDeadlines, instanceAdminClient: instancepb.NewInstanceAdminClient(connPool), CallOptions: &client.CallOptions}
	c.setGoogleClientInfo()
	client.internalClient = c
	client.LROClient, err = lroauto.NewOperationsClient(ctx, gtransport.WithConnPool(connPool))
	if err != nil {
		return nil, err
	}
	c.LROClient = &client.LROClient
	return &client, nil
}
func (c *instanceAdminGRPCClient) gologoo__Connection_9d780f59578540b8e87027995d65246b() *grpc.ClientConn {
	return c.connPool.Conn()
}
func (c *instanceAdminGRPCClient) gologoo__setGoogleClientInfo_9d780f59578540b8e87027995d65246b(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", getVersionClient(), "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}
func (c *instanceAdminGRPCClient) gologoo__Close_9d780f59578540b8e87027995d65246b() error {
	return c.connPool.Close()
}
func (c *instanceAdminGRPCClient) gologoo__ListInstanceConfigs_9d780f59578540b8e87027995d65246b(ctx context.Context, req *instancepb.ListInstanceConfigsRequest, opts ...gax.CallOption) *InstanceConfigIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).ListInstanceConfigs[0:len((*c.CallOptions).ListInstanceConfigs):len((*c.CallOptions).ListInstanceConfigs)], opts...)
	it := &InstanceConfigIterator{}
	req = proto.Clone(req).(*instancepb.ListInstanceConfigsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*instancepb.InstanceConfig, string, error) {
		resp := &instancepb.ListInstanceConfigsResponse{}
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
			resp, err = c.instanceAdminClient.ListInstanceConfigs(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		it.Response = resp
		return resp.GetInstanceConfigs(), resp.GetNextPageToken(), nil
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
func (c *instanceAdminGRPCClient) gologoo__GetInstanceConfig_9d780f59578540b8e87027995d65246b(ctx context.Context, req *instancepb.GetInstanceConfigRequest, opts ...gax.CallOption) (*instancepb.InstanceConfig, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 3600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).GetInstanceConfig[0:len((*c.CallOptions).GetInstanceConfig):len((*c.CallOptions).GetInstanceConfig)], opts...)
	var resp *instancepb.InstanceConfig
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.instanceAdminClient.GetInstanceConfig(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *instanceAdminGRPCClient) gologoo__ListInstances_9d780f59578540b8e87027995d65246b(ctx context.Context, req *instancepb.ListInstancesRequest, opts ...gax.CallOption) *InstanceIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).ListInstances[0:len((*c.CallOptions).ListInstances):len((*c.CallOptions).ListInstances)], opts...)
	it := &InstanceIterator{}
	req = proto.Clone(req).(*instancepb.ListInstancesRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*instancepb.Instance, string, error) {
		resp := &instancepb.ListInstancesResponse{}
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
			resp, err = c.instanceAdminClient.ListInstances(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		it.Response = resp
		return resp.GetInstances(), resp.GetNextPageToken(), nil
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
func (c *instanceAdminGRPCClient) gologoo__GetInstance_9d780f59578540b8e87027995d65246b(ctx context.Context, req *instancepb.GetInstanceRequest, opts ...gax.CallOption) (*instancepb.Instance, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 3600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).GetInstance[0:len((*c.CallOptions).GetInstance):len((*c.CallOptions).GetInstance)], opts...)
	var resp *instancepb.Instance
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.instanceAdminClient.GetInstance(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *instanceAdminGRPCClient) gologoo__CreateInstance_9d780f59578540b8e87027995d65246b(ctx context.Context, req *instancepb.CreateInstanceRequest, opts ...gax.CallOption) (*CreateInstanceOperation, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 3600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", url.QueryEscape(req.GetParent())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).CreateInstance[0:len((*c.CallOptions).CreateInstance):len((*c.CallOptions).CreateInstance)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.instanceAdminClient.CreateInstance(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &CreateInstanceOperation{lro: longrunning.InternalNewOperation(*c.LROClient, resp)}, nil
}
func (c *instanceAdminGRPCClient) gologoo__UpdateInstance_9d780f59578540b8e87027995d65246b(ctx context.Context, req *instancepb.UpdateInstanceRequest, opts ...gax.CallOption) (*UpdateInstanceOperation, error) {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 3600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "instance.name", url.QueryEscape(req.GetInstance().GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).UpdateInstance[0:len((*c.CallOptions).UpdateInstance):len((*c.CallOptions).UpdateInstance)], opts...)
	var resp *longrunningpb.Operation
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.instanceAdminClient.UpdateInstance(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return &UpdateInstanceOperation{lro: longrunning.InternalNewOperation(*c.LROClient, resp)}, nil
}
func (c *instanceAdminGRPCClient) gologoo__DeleteInstance_9d780f59578540b8e87027995d65246b(ctx context.Context, req *instancepb.DeleteInstanceRequest, opts ...gax.CallOption) error {
	if _, ok := ctx.Deadline(); !ok && !c.disableDeadlines {
		cctx, cancel := context.WithTimeout(ctx, 3600000*time.Millisecond)
		defer cancel()
		ctx = cctx
	}
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", url.QueryEscape(req.GetName())))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append((*c.CallOptions).DeleteInstance[0:len((*c.CallOptions).DeleteInstance):len((*c.CallOptions).DeleteInstance)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.instanceAdminClient.DeleteInstance(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}
func (c *instanceAdminGRPCClient) gologoo__SetIamPolicy_9d780f59578540b8e87027995d65246b(ctx context.Context, req *iampb.SetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
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
		resp, err = c.instanceAdminClient.SetIamPolicy(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *instanceAdminGRPCClient) gologoo__GetIamPolicy_9d780f59578540b8e87027995d65246b(ctx context.Context, req *iampb.GetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
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
		resp, err = c.instanceAdminClient.GetIamPolicy(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func (c *instanceAdminGRPCClient) gologoo__TestIamPermissions_9d780f59578540b8e87027995d65246b(ctx context.Context, req *iampb.TestIamPermissionsRequest, opts ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error) {
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
		resp, err = c.instanceAdminClient.TestIamPermissions(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type CreateInstanceOperation struct {
	lro *longrunning.Operation
}

func (c *instanceAdminGRPCClient) gologoo__CreateInstanceOperation_9d780f59578540b8e87027995d65246b(name string) *CreateInstanceOperation {
	return &CreateInstanceOperation{lro: longrunning.InternalNewOperation(*c.LROClient, &longrunningpb.Operation{Name: name})}
}
func (op *CreateInstanceOperation) gologoo__Wait_9d780f59578540b8e87027995d65246b(ctx context.Context, opts ...gax.CallOption) (*instancepb.Instance, error) {
	var resp instancepb.Instance
	if err := op.lro.WaitWithInterval(ctx, &resp, time.Minute, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}
func (op *CreateInstanceOperation) gologoo__Poll_9d780f59578540b8e87027995d65246b(ctx context.Context, opts ...gax.CallOption) (*instancepb.Instance, error) {
	var resp instancepb.Instance
	if err := op.lro.Poll(ctx, &resp, opts...); err != nil {
		return nil, err
	}
	if !op.Done() {
		return nil, nil
	}
	return &resp, nil
}
func (op *CreateInstanceOperation) gologoo__Metadata_9d780f59578540b8e87027995d65246b() (*instancepb.CreateInstanceMetadata, error) {
	var meta instancepb.CreateInstanceMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}
func (op *CreateInstanceOperation) gologoo__Done_9d780f59578540b8e87027995d65246b() bool {
	return op.lro.Done()
}
func (op *CreateInstanceOperation) gologoo__Name_9d780f59578540b8e87027995d65246b() string {
	return op.lro.Name()
}

type UpdateInstanceOperation struct {
	lro *longrunning.Operation
}

func (c *instanceAdminGRPCClient) gologoo__UpdateInstanceOperation_9d780f59578540b8e87027995d65246b(name string) *UpdateInstanceOperation {
	return &UpdateInstanceOperation{lro: longrunning.InternalNewOperation(*c.LROClient, &longrunningpb.Operation{Name: name})}
}
func (op *UpdateInstanceOperation) gologoo__Wait_9d780f59578540b8e87027995d65246b(ctx context.Context, opts ...gax.CallOption) (*instancepb.Instance, error) {
	var resp instancepb.Instance
	if err := op.lro.WaitWithInterval(ctx, &resp, time.Minute, opts...); err != nil {
		return nil, err
	}
	return &resp, nil
}
func (op *UpdateInstanceOperation) gologoo__Poll_9d780f59578540b8e87027995d65246b(ctx context.Context, opts ...gax.CallOption) (*instancepb.Instance, error) {
	var resp instancepb.Instance
	if err := op.lro.Poll(ctx, &resp, opts...); err != nil {
		return nil, err
	}
	if !op.Done() {
		return nil, nil
	}
	return &resp, nil
}
func (op *UpdateInstanceOperation) gologoo__Metadata_9d780f59578540b8e87027995d65246b() (*instancepb.UpdateInstanceMetadata, error) {
	var meta instancepb.UpdateInstanceMetadata
	if err := op.lro.Metadata(&meta); err == longrunning.ErrNoMetadata {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &meta, nil
}
func (op *UpdateInstanceOperation) gologoo__Done_9d780f59578540b8e87027995d65246b() bool {
	return op.lro.Done()
}
func (op *UpdateInstanceOperation) gologoo__Name_9d780f59578540b8e87027995d65246b() string {
	return op.lro.Name()
}

type InstanceConfigIterator struct {
	items    []*instancepb.InstanceConfig
	pageInfo *iterator.PageInfo
	nextFunc func() error
	Response interface {
	}
	InternalFetch func(pageSize int, pageToken string) (results []*instancepb.InstanceConfig, nextPageToken string, err error)
}

func (it *InstanceConfigIterator) gologoo__PageInfo_9d780f59578540b8e87027995d65246b() *iterator.PageInfo {
	return it.pageInfo
}
func (it *InstanceConfigIterator) gologoo__Next_9d780f59578540b8e87027995d65246b() (*instancepb.InstanceConfig, error) {
	var item *instancepb.InstanceConfig
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}
func (it *InstanceConfigIterator) gologoo__bufLen_9d780f59578540b8e87027995d65246b() int {
	return len(it.items)
}
func (it *InstanceConfigIterator) gologoo__takeBuf_9d780f59578540b8e87027995d65246b() interface {
} {
	b := it.items
	it.items = nil
	return b
}

type InstanceIterator struct {
	items    []*instancepb.Instance
	pageInfo *iterator.PageInfo
	nextFunc func() error
	Response interface {
	}
	InternalFetch func(pageSize int, pageToken string) (results []*instancepb.Instance, nextPageToken string, err error)
}

func (it *InstanceIterator) gologoo__PageInfo_9d780f59578540b8e87027995d65246b() *iterator.PageInfo {
	return it.pageInfo
}
func (it *InstanceIterator) gologoo__Next_9d780f59578540b8e87027995d65246b() (*instancepb.Instance, error) {
	var item *instancepb.Instance
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}
func (it *InstanceIterator) gologoo__bufLen_9d780f59578540b8e87027995d65246b() int {
	return len(it.items)
}
func (it *InstanceIterator) gologoo__takeBuf_9d780f59578540b8e87027995d65246b() interface {
} {
	b := it.items
	it.items = nil
	return b
}
func defaultInstanceAdminGRPCClientOptions() []option.ClientOption {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__defaultInstanceAdminGRPCClientOptions_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0 := gologoo__defaultInstanceAdminGRPCClientOptions_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v\n", r0)
	return r0
}
func defaultInstanceAdminCallOptions() *InstanceAdminCallOptions {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__defaultInstanceAdminCallOptions_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0 := gologoo__defaultInstanceAdminCallOptions_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *InstanceAdminClient) Close() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Close_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__Close_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *InstanceAdminClient) setGoogleClientInfo(keyval ...string) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__setGoogleClientInfo_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v\n", keyval)
	c.gologoo__setGoogleClientInfo_9d780f59578540b8e87027995d65246b(keyval...)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (c *InstanceAdminClient) Connection() *grpc.ClientConn {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Connection_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__Connection_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *InstanceAdminClient) ListInstanceConfigs(ctx context.Context, req *instancepb.ListInstanceConfigsRequest, opts ...gax.CallOption) *InstanceConfigIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ListInstanceConfigs_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__ListInstanceConfigs_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *InstanceAdminClient) GetInstanceConfig(ctx context.Context, req *instancepb.GetInstanceConfigRequest, opts ...gax.CallOption) (*instancepb.InstanceConfig, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetInstanceConfig_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__GetInstanceConfig_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *InstanceAdminClient) ListInstances(ctx context.Context, req *instancepb.ListInstancesRequest, opts ...gax.CallOption) *InstanceIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ListInstances_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__ListInstances_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *InstanceAdminClient) GetInstance(ctx context.Context, req *instancepb.GetInstanceRequest, opts ...gax.CallOption) (*instancepb.Instance, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetInstance_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__GetInstance_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *InstanceAdminClient) CreateInstance(ctx context.Context, req *instancepb.CreateInstanceRequest, opts ...gax.CallOption) (*CreateInstanceOperation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateInstance_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__CreateInstance_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *InstanceAdminClient) CreateInstanceOperation(name string) *CreateInstanceOperation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateInstanceOperation_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v\n", name)
	r0 := c.gologoo__CreateInstanceOperation_9d780f59578540b8e87027995d65246b(name)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *InstanceAdminClient) UpdateInstance(ctx context.Context, req *instancepb.UpdateInstanceRequest, opts ...gax.CallOption) (*UpdateInstanceOperation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UpdateInstance_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__UpdateInstance_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *InstanceAdminClient) UpdateInstanceOperation(name string) *UpdateInstanceOperation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UpdateInstanceOperation_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v\n", name)
	r0 := c.gologoo__UpdateInstanceOperation_9d780f59578540b8e87027995d65246b(name)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *InstanceAdminClient) DeleteInstance(ctx context.Context, req *instancepb.DeleteInstanceRequest, opts ...gax.CallOption) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DeleteInstance_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__DeleteInstance_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *InstanceAdminClient) SetIamPolicy(ctx context.Context, req *iampb.SetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SetIamPolicy_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__SetIamPolicy_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *InstanceAdminClient) GetIamPolicy(ctx context.Context, req *iampb.GetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetIamPolicy_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__GetIamPolicy_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *InstanceAdminClient) TestIamPermissions(ctx context.Context, req *iampb.TestIamPermissionsRequest, opts ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__TestIamPermissions_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__TestIamPermissions_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func NewInstanceAdminClient(ctx context.Context, opts ...option.ClientOption) (*InstanceAdminClient, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NewInstanceAdminClient_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v\n", ctx, opts)
	r0, r1 := gologoo__NewInstanceAdminClient_9d780f59578540b8e87027995d65246b(ctx, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *instanceAdminGRPCClient) Connection() *grpc.ClientConn {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Connection_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__Connection_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *instanceAdminGRPCClient) setGoogleClientInfo(keyval ...string) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__setGoogleClientInfo_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v\n", keyval)
	c.gologoo__setGoogleClientInfo_9d780f59578540b8e87027995d65246b(keyval...)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (c *instanceAdminGRPCClient) Close() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Close_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__Close_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *instanceAdminGRPCClient) ListInstanceConfigs(ctx context.Context, req *instancepb.ListInstanceConfigsRequest, opts ...gax.CallOption) *InstanceConfigIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ListInstanceConfigs_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__ListInstanceConfigs_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *instanceAdminGRPCClient) GetInstanceConfig(ctx context.Context, req *instancepb.GetInstanceConfigRequest, opts ...gax.CallOption) (*instancepb.InstanceConfig, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetInstanceConfig_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__GetInstanceConfig_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *instanceAdminGRPCClient) ListInstances(ctx context.Context, req *instancepb.ListInstancesRequest, opts ...gax.CallOption) *InstanceIterator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ListInstances_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__ListInstances_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *instanceAdminGRPCClient) GetInstance(ctx context.Context, req *instancepb.GetInstanceRequest, opts ...gax.CallOption) (*instancepb.Instance, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetInstance_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__GetInstance_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *instanceAdminGRPCClient) CreateInstance(ctx context.Context, req *instancepb.CreateInstanceRequest, opts ...gax.CallOption) (*CreateInstanceOperation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateInstance_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__CreateInstance_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *instanceAdminGRPCClient) UpdateInstance(ctx context.Context, req *instancepb.UpdateInstanceRequest, opts ...gax.CallOption) (*UpdateInstanceOperation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UpdateInstance_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__UpdateInstance_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *instanceAdminGRPCClient) DeleteInstance(ctx context.Context, req *instancepb.DeleteInstanceRequest, opts ...gax.CallOption) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DeleteInstance_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0 := c.gologoo__DeleteInstance_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *instanceAdminGRPCClient) SetIamPolicy(ctx context.Context, req *iampb.SetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SetIamPolicy_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__SetIamPolicy_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *instanceAdminGRPCClient) GetIamPolicy(ctx context.Context, req *iampb.GetIamPolicyRequest, opts ...gax.CallOption) (*iampb.Policy, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetIamPolicy_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__GetIamPolicy_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *instanceAdminGRPCClient) TestIamPermissions(ctx context.Context, req *iampb.TestIamPermissionsRequest, opts ...gax.CallOption) (*iampb.TestIamPermissionsResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__TestIamPermissions_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__TestIamPermissions_9d780f59578540b8e87027995d65246b(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *instanceAdminGRPCClient) CreateInstanceOperation(name string) *CreateInstanceOperation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateInstanceOperation_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v\n", name)
	r0 := c.gologoo__CreateInstanceOperation_9d780f59578540b8e87027995d65246b(name)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (op *CreateInstanceOperation) Wait(ctx context.Context, opts ...gax.CallOption) (*instancepb.Instance, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Wait_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v\n", ctx, opts)
	r0, r1 := op.gologoo__Wait_9d780f59578540b8e87027995d65246b(ctx, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (op *CreateInstanceOperation) Poll(ctx context.Context, opts ...gax.CallOption) (*instancepb.Instance, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Poll_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v\n", ctx, opts)
	r0, r1 := op.gologoo__Poll_9d780f59578540b8e87027995d65246b(ctx, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (op *CreateInstanceOperation) Metadata() (*instancepb.CreateInstanceMetadata, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Metadata_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0, r1 := op.gologoo__Metadata_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (op *CreateInstanceOperation) Done() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Done_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0 := op.gologoo__Done_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (op *CreateInstanceOperation) Name() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Name_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0 := op.gologoo__Name_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *instanceAdminGRPCClient) UpdateInstanceOperation(name string) *UpdateInstanceOperation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UpdateInstanceOperation_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v\n", name)
	r0 := c.gologoo__UpdateInstanceOperation_9d780f59578540b8e87027995d65246b(name)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (op *UpdateInstanceOperation) Wait(ctx context.Context, opts ...gax.CallOption) (*instancepb.Instance, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Wait_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v\n", ctx, opts)
	r0, r1 := op.gologoo__Wait_9d780f59578540b8e87027995d65246b(ctx, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (op *UpdateInstanceOperation) Poll(ctx context.Context, opts ...gax.CallOption) (*instancepb.Instance, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Poll_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : %v %v\n", ctx, opts)
	r0, r1 := op.gologoo__Poll_9d780f59578540b8e87027995d65246b(ctx, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (op *UpdateInstanceOperation) Metadata() (*instancepb.UpdateInstanceMetadata, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Metadata_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0, r1 := op.gologoo__Metadata_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (op *UpdateInstanceOperation) Done() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Done_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0 := op.gologoo__Done_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (op *UpdateInstanceOperation) Name() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Name_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0 := op.gologoo__Name_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (it *InstanceConfigIterator) PageInfo() *iterator.PageInfo {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PageInfo_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__PageInfo_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (it *InstanceConfigIterator) Next() (*instancepb.InstanceConfig, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Next_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0, r1 := it.gologoo__Next_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (it *InstanceConfigIterator) bufLen() int {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__bufLen_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__bufLen_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (it *InstanceConfigIterator) takeBuf() interface {
} {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__takeBuf_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__takeBuf_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (it *InstanceIterator) PageInfo() *iterator.PageInfo {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PageInfo_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__PageInfo_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (it *InstanceIterator) Next() (*instancepb.Instance, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Next_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0, r1 := it.gologoo__Next_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (it *InstanceIterator) bufLen() int {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__bufLen_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__bufLen_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (it *InstanceIterator) takeBuf() interface {
} {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__takeBuf_9d780f59578540b8e87027995d65246b")
	log.Printf("Input : (none)\n")
	r0 := it.gologoo__takeBuf_9d780f59578540b8e87027995d65246b()
	log.Printf("Output: %v\n", r0)
	return r0
}
