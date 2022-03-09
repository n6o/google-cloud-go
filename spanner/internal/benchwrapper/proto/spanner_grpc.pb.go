package spanner_bench

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"log"
)

const _ = grpc.SupportPackageIsVersion7

type SpannerBenchWrapperClient interface {
	Read(ctx context.Context, in *ReadQuery, opts ...grpc.CallOption) (*EmptyResponse, error)
	Insert(ctx context.Context, in *InsertQuery, opts ...grpc.CallOption) (*EmptyResponse, error)
	Update(ctx context.Context, in *UpdateQuery, opts ...grpc.CallOption) (*EmptyResponse, error)
}
type spannerBenchWrapperClient struct {
	cc grpc.ClientConnInterface
}

func gologoo__NewSpannerBenchWrapperClient_3f924f5bee8266ebf2a5414294aea9e2(cc grpc.ClientConnInterface) SpannerBenchWrapperClient {
	return &spannerBenchWrapperClient{cc}
}
func (c *spannerBenchWrapperClient) gologoo__Read_3f924f5bee8266ebf2a5414294aea9e2(ctx context.Context, in *ReadQuery, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/spanner_bench.SpannerBenchWrapper/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *spannerBenchWrapperClient) gologoo__Insert_3f924f5bee8266ebf2a5414294aea9e2(ctx context.Context, in *InsertQuery, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/spanner_bench.SpannerBenchWrapper/Insert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *spannerBenchWrapperClient) gologoo__Update_3f924f5bee8266ebf2a5414294aea9e2(ctx context.Context, in *UpdateQuery, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/spanner_bench.SpannerBenchWrapper/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type SpannerBenchWrapperServer interface {
	Read(context.Context, *ReadQuery) (*EmptyResponse, error)
	Insert(context.Context, *InsertQuery) (*EmptyResponse, error)
	Update(context.Context, *UpdateQuery) (*EmptyResponse, error)
	mustEmbedUnimplementedSpannerBenchWrapperServer()
}
type UnimplementedSpannerBenchWrapperServer struct {
}

func (UnimplementedSpannerBenchWrapperServer) gologoo__Read_3f924f5bee8266ebf2a5414294aea9e2(context.Context, *ReadQuery) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedSpannerBenchWrapperServer) gologoo__Insert_3f924f5bee8266ebf2a5414294aea9e2(context.Context, *InsertQuery) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Insert not implemented")
}
func (UnimplementedSpannerBenchWrapperServer) gologoo__Update_3f924f5bee8266ebf2a5414294aea9e2(context.Context, *UpdateQuery) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedSpannerBenchWrapperServer) gologoo__mustEmbedUnimplementedSpannerBenchWrapperServer_3f924f5bee8266ebf2a5414294aea9e2() {
}

type UnsafeSpannerBenchWrapperServer interface {
	mustEmbedUnimplementedSpannerBenchWrapperServer()
}

func gologoo__RegisterSpannerBenchWrapperServer_3f924f5bee8266ebf2a5414294aea9e2(s grpc.ServiceRegistrar, srv SpannerBenchWrapperServer) {
	s.RegisterService(&SpannerBenchWrapper_ServiceDesc, srv)
}
func gologoo___SpannerBenchWrapper_Read_Handler_3f924f5bee8266ebf2a5414294aea9e2(srv interface {
}, ctx context.Context, dec func(interface {
}) error, interceptor grpc.UnaryServerInterceptor) (interface {
}, error) {
	in := new(ReadQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpannerBenchWrapperServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/spanner_bench.SpannerBenchWrapper/Read"}
	handler := func(ctx context.Context, req interface {
	}) (interface {
	}, error) {
		return srv.(SpannerBenchWrapperServer).Read(ctx, req.(*ReadQuery))
	}
	return interceptor(ctx, in, info, handler)
}
func gologoo___SpannerBenchWrapper_Insert_Handler_3f924f5bee8266ebf2a5414294aea9e2(srv interface {
}, ctx context.Context, dec func(interface {
}) error, interceptor grpc.UnaryServerInterceptor) (interface {
}, error) {
	in := new(InsertQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpannerBenchWrapperServer).Insert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/spanner_bench.SpannerBenchWrapper/Insert"}
	handler := func(ctx context.Context, req interface {
	}) (interface {
	}, error) {
		return srv.(SpannerBenchWrapperServer).Insert(ctx, req.(*InsertQuery))
	}
	return interceptor(ctx, in, info, handler)
}
func gologoo___SpannerBenchWrapper_Update_Handler_3f924f5bee8266ebf2a5414294aea9e2(srv interface {
}, ctx context.Context, dec func(interface {
}) error, interceptor grpc.UnaryServerInterceptor) (interface {
}, error) {
	in := new(UpdateQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpannerBenchWrapperServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/spanner_bench.SpannerBenchWrapper/Update"}
	handler := func(ctx context.Context, req interface {
	}) (interface {
	}, error) {
		return srv.(SpannerBenchWrapperServer).Update(ctx, req.(*UpdateQuery))
	}
	return interceptor(ctx, in, info, handler)
}

var SpannerBenchWrapper_ServiceDesc = grpc.ServiceDesc{ServiceName: "spanner_bench.SpannerBenchWrapper", HandlerType: (*SpannerBenchWrapperServer)(nil), Methods: []grpc.MethodDesc{{MethodName: "Read", Handler: _SpannerBenchWrapper_Read_Handler}, {MethodName: "Insert", Handler: _SpannerBenchWrapper_Insert_Handler}, {MethodName: "Update", Handler: _SpannerBenchWrapper_Update_Handler}}, Streams: []grpc.StreamDesc{}, Metadata: "spanner.proto"}

func NewSpannerBenchWrapperClient(cc grpc.ClientConnInterface) SpannerBenchWrapperClient {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NewSpannerBenchWrapperClient_3f924f5bee8266ebf2a5414294aea9e2")
	log.Printf("Input : %v\n", cc)
	r0 := gologoo__NewSpannerBenchWrapperClient_3f924f5bee8266ebf2a5414294aea9e2(cc)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *spannerBenchWrapperClient) Read(ctx context.Context, in *ReadQuery, opts ...grpc.CallOption) (*EmptyResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Read_3f924f5bee8266ebf2a5414294aea9e2")
	log.Printf("Input : %v %v %v\n", ctx, in, opts)
	r0, r1 := c.gologoo__Read_3f924f5bee8266ebf2a5414294aea9e2(ctx, in, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *spannerBenchWrapperClient) Insert(ctx context.Context, in *InsertQuery, opts ...grpc.CallOption) (*EmptyResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Insert_3f924f5bee8266ebf2a5414294aea9e2")
	log.Printf("Input : %v %v %v\n", ctx, in, opts)
	r0, r1 := c.gologoo__Insert_3f924f5bee8266ebf2a5414294aea9e2(ctx, in, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *spannerBenchWrapperClient) Update(ctx context.Context, in *UpdateQuery, opts ...grpc.CallOption) (*EmptyResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Update_3f924f5bee8266ebf2a5414294aea9e2")
	log.Printf("Input : %v %v %v\n", ctx, in, opts)
	r0, r1 := c.gologoo__Update_3f924f5bee8266ebf2a5414294aea9e2(ctx, in, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (recv UnimplementedSpannerBenchWrapperServer) Read(context.Context, *ReadQuery) (*EmptyResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Read_3f924f5bee8266ebf2a5414294aea9e2")
	log.Printf("Input : (none)\n")
	r0, r1 := recv.gologoo__Read_3f924f5bee8266ebf2a5414294aea9e2(nil, nil)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (recv UnimplementedSpannerBenchWrapperServer) Insert(context.Context, *InsertQuery) (*EmptyResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Insert_3f924f5bee8266ebf2a5414294aea9e2")
	log.Printf("Input : (none)\n")
	r0, r1 := recv.gologoo__Insert_3f924f5bee8266ebf2a5414294aea9e2(nil, nil)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (recv UnimplementedSpannerBenchWrapperServer) Update(context.Context, *UpdateQuery) (*EmptyResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Update_3f924f5bee8266ebf2a5414294aea9e2")
	log.Printf("Input : (none)\n")
	r0, r1 := recv.gologoo__Update_3f924f5bee8266ebf2a5414294aea9e2(nil, nil)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (recv UnimplementedSpannerBenchWrapperServer) mustEmbedUnimplementedSpannerBenchWrapperServer() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__mustEmbedUnimplementedSpannerBenchWrapperServer_3f924f5bee8266ebf2a5414294aea9e2")
	log.Printf("Input : (none)\n")
	recv.gologoo__mustEmbedUnimplementedSpannerBenchWrapperServer_3f924f5bee8266ebf2a5414294aea9e2()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func RegisterSpannerBenchWrapperServer(s grpc.ServiceRegistrar, srv SpannerBenchWrapperServer) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__RegisterSpannerBenchWrapperServer_3f924f5bee8266ebf2a5414294aea9e2")
	log.Printf("Input : %v %v\n", s, srv)
	gologoo__RegisterSpannerBenchWrapperServer_3f924f5bee8266ebf2a5414294aea9e2(s, srv)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func _SpannerBenchWrapper_Read_Handler(srv interface {
}, ctx context.Context, dec func(interface {
}) error, interceptor grpc.UnaryServerInterceptor) (interface {
}, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo___SpannerBenchWrapper_Read_Handler_3f924f5bee8266ebf2a5414294aea9e2")
	log.Printf("Input : %v %v %v %v\n", srv, ctx, dec, interceptor)
	r0, r1 := gologoo___SpannerBenchWrapper_Read_Handler_3f924f5bee8266ebf2a5414294aea9e2(srv, ctx, dec, interceptor)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func _SpannerBenchWrapper_Insert_Handler(srv interface {
}, ctx context.Context, dec func(interface {
}) error, interceptor grpc.UnaryServerInterceptor) (interface {
}, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo___SpannerBenchWrapper_Insert_Handler_3f924f5bee8266ebf2a5414294aea9e2")
	log.Printf("Input : %v %v %v %v\n", srv, ctx, dec, interceptor)
	r0, r1 := gologoo___SpannerBenchWrapper_Insert_Handler_3f924f5bee8266ebf2a5414294aea9e2(srv, ctx, dec, interceptor)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func _SpannerBenchWrapper_Update_Handler(srv interface {
}, ctx context.Context, dec func(interface {
}) error, interceptor grpc.UnaryServerInterceptor) (interface {
}, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo___SpannerBenchWrapper_Update_Handler_3f924f5bee8266ebf2a5414294aea9e2")
	log.Printf("Input : %v %v %v %v\n", srv, ctx, dec, interceptor)
	r0, r1 := gologoo___SpannerBenchWrapper_Update_Handler_3f924f5bee8266ebf2a5414294aea9e2(srv, ctx, dec, interceptor)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
