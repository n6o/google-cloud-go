package testutil

import (
	"fmt"
	"net"
	"strconv"
	"testing"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/api/option"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	spannerpb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc"
	"log"
)

const SelectFooFromBar = "SELECT FOO FROM BAR"
const selectFooFromBarRowCount int64 = 2
const selectFooFromBarColCount int = 1

var selectFooFromBarResults = [...]int64{1, 2}

const SelectSingerIDAlbumIDAlbumTitleFromAlbums = "SELECT SingerId, AlbumId, AlbumTitle FROM Albums"
const SelectSingerIDAlbumIDAlbumTitleFromAlbumsRowCount int64 = 3
const SelectSingerIDAlbumIDAlbumTitleFromAlbumsColCount int = 3
const UpdateBarSetFoo = "UPDATE FOO SET BAR=1 WHERE BAZ=2"
const UpdateBarSetFooRowCount = 5

type MockedSpannerInMemTestServer struct {
	TestSpanner       InMemSpannerServer
	TestInstanceAdmin InMemInstanceAdminServer
	server            *grpc.Server
}

func gologoo__NewMockedSpannerInMemTestServer_1a39ae16c8b612982e14d70c47f32ccf(t *testing.T) (mockedServer *MockedSpannerInMemTestServer, opts []option.ClientOption, teardown func()) {
	return NewMockedSpannerInMemTestServerWithAddr(t, "localhost:0")
}
func gologoo__NewMockedSpannerInMemTestServerWithAddr_1a39ae16c8b612982e14d70c47f32ccf(t *testing.T, addr string) (mockedServer *MockedSpannerInMemTestServer, opts []option.ClientOption, teardown func()) {
	mockedServer = &MockedSpannerInMemTestServer{}
	opts = mockedServer.setupMockedServerWithAddr(t, addr)
	return mockedServer, opts, func() {
		mockedServer.TestSpanner.Stop()
		mockedServer.TestInstanceAdmin.Stop()
		mockedServer.server.Stop()
	}
}
func (s *MockedSpannerInMemTestServer) gologoo__setupMockedServerWithAddr_1a39ae16c8b612982e14d70c47f32ccf(t *testing.T, addr string) []option.ClientOption {
	s.TestSpanner = NewInMemSpannerServer()
	s.TestInstanceAdmin = NewInMemInstanceAdminServer()
	s.setupFooResults()
	s.setupSingersResults()
	s.server = grpc.NewServer()
	spannerpb.RegisterSpannerServer(s.server, s.TestSpanner)
	instancepb.RegisterInstanceAdminServer(s.server, s.TestInstanceAdmin)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	go s.server.Serve(lis)
	serverAddress := lis.Addr().String()
	opts := []option.ClientOption{option.WithEndpoint(serverAddress), option.WithGRPCDialOption(grpc.WithInsecure()), option.WithoutAuthentication()}
	return opts
}
func (s *MockedSpannerInMemTestServer) gologoo__setupFooResults_1a39ae16c8b612982e14d70c47f32ccf() {
	fields := make([]*spannerpb.StructType_Field, selectFooFromBarColCount)
	fields[0] = &spannerpb.StructType_Field{Name: "FOO", Type: &spannerpb.Type{Code: spannerpb.TypeCode_INT64}}
	rowType := &spannerpb.StructType{Fields: fields}
	metadata := &spannerpb.ResultSetMetadata{RowType: rowType}
	rows := make([]*structpb.ListValue, selectFooFromBarRowCount)
	for idx, value := range selectFooFromBarResults {
		rowValue := make([]*structpb.Value, selectFooFromBarColCount)
		rowValue[0] = &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: strconv.FormatInt(value, 10)}}
		rows[idx] = &structpb.ListValue{Values: rowValue}
	}
	resultSet := &spannerpb.ResultSet{Metadata: metadata, Rows: rows}
	result := &StatementResult{Type: StatementResultResultSet, ResultSet: resultSet}
	s.TestSpanner.PutStatementResult(SelectFooFromBar, result)
	s.TestSpanner.PutStatementResult(UpdateBarSetFoo, &StatementResult{Type: StatementResultUpdateCount, UpdateCount: UpdateBarSetFooRowCount})
}
func (s *MockedSpannerInMemTestServer) gologoo__setupSingersResults_1a39ae16c8b612982e14d70c47f32ccf() {
	metadata := createSingersMetadata()
	rows := make([]*structpb.ListValue, SelectSingerIDAlbumIDAlbumTitleFromAlbumsRowCount)
	var idx int64
	for idx = 0; idx < SelectSingerIDAlbumIDAlbumTitleFromAlbumsRowCount; idx++ {
		rows[idx] = createSingersRow(idx)
	}
	resultSet := &spannerpb.ResultSet{Metadata: metadata, Rows: rows}
	result := &StatementResult{Type: StatementResultResultSet, ResultSet: resultSet}
	s.TestSpanner.PutStatementResult(SelectSingerIDAlbumIDAlbumTitleFromAlbums, result)
}
func (s *MockedSpannerInMemTestServer) gologoo__CreateSingleRowSingersResult_1a39ae16c8b612982e14d70c47f32ccf(rowNum int64) *StatementResult {
	metadata := createSingersMetadata()
	var returnedRows int
	if rowNum < SelectSingerIDAlbumIDAlbumTitleFromAlbumsRowCount {
		returnedRows = 1
	} else {
		returnedRows = 0
	}
	rows := make([]*structpb.ListValue, returnedRows)
	if returnedRows > 0 {
		rows[0] = createSingersRow(rowNum)
	}
	resultSet := &spannerpb.ResultSet{Metadata: metadata, Rows: rows}
	return &StatementResult{Type: StatementResultResultSet, ResultSet: resultSet}
}
func gologoo__createSingersMetadata_1a39ae16c8b612982e14d70c47f32ccf() *spannerpb.ResultSetMetadata {
	fields := make([]*spannerpb.StructType_Field, SelectSingerIDAlbumIDAlbumTitleFromAlbumsColCount)
	fields[0] = &spannerpb.StructType_Field{Name: "SingerId", Type: &spannerpb.Type{Code: spannerpb.TypeCode_INT64}}
	fields[1] = &spannerpb.StructType_Field{Name: "AlbumId", Type: &spannerpb.Type{Code: spannerpb.TypeCode_INT64}}
	fields[2] = &spannerpb.StructType_Field{Name: "AlbumTitle", Type: &spannerpb.Type{Code: spannerpb.TypeCode_STRING}}
	rowType := &spannerpb.StructType{Fields: fields}
	return &spannerpb.ResultSetMetadata{RowType: rowType}
}
func gologoo__createSingersRow_1a39ae16c8b612982e14d70c47f32ccf(idx int64) *structpb.ListValue {
	rowValue := make([]*structpb.Value, SelectSingerIDAlbumIDAlbumTitleFromAlbumsColCount)
	rowValue[0] = &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: strconv.FormatInt(idx+1, 10)}}
	rowValue[1] = &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: strconv.FormatInt(idx*10+idx, 10)}}
	rowValue[2] = &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: fmt.Sprintf("Album title %d", idx)}}
	return &structpb.ListValue{Values: rowValue}
}
func NewMockedSpannerInMemTestServer(t *testing.T) (mockedServer *MockedSpannerInMemTestServer, opts []option.ClientOption, teardown func()) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NewMockedSpannerInMemTestServer_1a39ae16c8b612982e14d70c47f32ccf")
	log.Printf("Input : %v\n", t)
	mockedServer, opts, teardown = gologoo__NewMockedSpannerInMemTestServer_1a39ae16c8b612982e14d70c47f32ccf(t)
	log.Printf("Output: %v %v %v\n", mockedServer, opts, teardown)
	return
}
func NewMockedSpannerInMemTestServerWithAddr(t *testing.T, addr string) (mockedServer *MockedSpannerInMemTestServer, opts []option.ClientOption, teardown func()) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NewMockedSpannerInMemTestServerWithAddr_1a39ae16c8b612982e14d70c47f32ccf")
	log.Printf("Input : %v %v\n", t, addr)
	mockedServer, opts, teardown = gologoo__NewMockedSpannerInMemTestServerWithAddr_1a39ae16c8b612982e14d70c47f32ccf(t, addr)
	log.Printf("Output: %v %v %v\n", mockedServer, opts, teardown)
	return
}
func (s *MockedSpannerInMemTestServer) setupMockedServerWithAddr(t *testing.T, addr string) []option.ClientOption {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__setupMockedServerWithAddr_1a39ae16c8b612982e14d70c47f32ccf")
	log.Printf("Input : %v %v\n", t, addr)
	r0 := s.gologoo__setupMockedServerWithAddr_1a39ae16c8b612982e14d70c47f32ccf(t, addr)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *MockedSpannerInMemTestServer) setupFooResults() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__setupFooResults_1a39ae16c8b612982e14d70c47f32ccf")
	log.Printf("Input : (none)\n")
	s.gologoo__setupFooResults_1a39ae16c8b612982e14d70c47f32ccf()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *MockedSpannerInMemTestServer) setupSingersResults() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__setupSingersResults_1a39ae16c8b612982e14d70c47f32ccf")
	log.Printf("Input : (none)\n")
	s.gologoo__setupSingersResults_1a39ae16c8b612982e14d70c47f32ccf()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *MockedSpannerInMemTestServer) CreateSingleRowSingersResult(rowNum int64) *StatementResult {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateSingleRowSingersResult_1a39ae16c8b612982e14d70c47f32ccf")
	log.Printf("Input : %v\n", rowNum)
	r0 := s.gologoo__CreateSingleRowSingersResult_1a39ae16c8b612982e14d70c47f32ccf(rowNum)
	log.Printf("Output: %v\n", r0)
	return r0
}
func createSingersMetadata() *spannerpb.ResultSetMetadata {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__createSingersMetadata_1a39ae16c8b612982e14d70c47f32ccf")
	log.Printf("Input : (none)\n")
	r0 := gologoo__createSingersMetadata_1a39ae16c8b612982e14d70c47f32ccf()
	log.Printf("Output: %v\n", r0)
	return r0
}
func createSingersRow(idx int64) *structpb.ListValue {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__createSingersRow_1a39ae16c8b612982e14d70c47f32ccf")
	log.Printf("Input : %v\n", idx)
	r0 := gologoo__createSingersRow_1a39ae16c8b612982e14d70c47f32ccf(idx)
	log.Printf("Output: %v\n", r0)
	return r0
}
