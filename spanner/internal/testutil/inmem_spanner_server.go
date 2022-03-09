package testutil

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"
	"github.com/golang/protobuf/ptypes"
	emptypb "github.com/golang/protobuf/ptypes/empty"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/genproto/googleapis/rpc/status"
	spannerpb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	gstatus "google.golang.org/grpc/status"
	"log"
)

var KvMeta = spannerpb.ResultSetMetadata{RowType: &spannerpb.StructType{Fields: []*spannerpb.StructType_Field{{Name: "Key", Type: &spannerpb.Type{Code: spannerpb.TypeCode_STRING}}, {Name: "Value", Type: &spannerpb.Type{Code: spannerpb.TypeCode_STRING}}}}}

type StatementResultType int

const (
	StatementResultError       StatementResultType = 0
	StatementResultResultSet   StatementResultType = 1
	StatementResultUpdateCount StatementResultType = 2
	MaxRowsPerPartialResultSet                     = 1
)
const (
	MethodBeginTransaction    string = "BEGIN_TRANSACTION"
	MethodCommitTransaction   string = "COMMIT_TRANSACTION"
	MethodBatchCreateSession  string = "BATCH_CREATE_SESSION"
	MethodCreateSession       string = "CREATE_SESSION"
	MethodDeleteSession       string = "DELETE_SESSION"
	MethodGetSession          string = "GET_SESSION"
	MethodExecuteSql          string = "EXECUTE_SQL"
	MethodExecuteStreamingSql string = "EXECUTE_STREAMING_SQL"
	MethodExecuteBatchDml     string = "EXECUTE_BATCH_DML"
	MethodStreamingRead       string = "EXECUTE_STREAMING_READ"
)

type StatementResult struct {
	Type         StatementResultType
	Err          error
	ResultSet    *spannerpb.ResultSet
	UpdateCount  int64
	ResumeTokens [][]byte
}
type PartialResultSetExecutionTime struct {
	ResumeToken   []byte
	ExecutionTime time.Duration
	Err           error
}

func (s *StatementResult) gologoo__ToPartialResultSets_705673bfbb22749846b5ab424c32ea33(resumeToken []byte) (result []*spannerpb.PartialResultSet, err error) {
	var startIndex uint64
	if len(resumeToken) > 0 {
		if startIndex, err = DecodeResumeToken(resumeToken); err != nil {
			return nil, err
		}
	}
	totalRows := uint64(len(s.ResultSet.Rows))
	if totalRows > 0 {
		for {
			rowCount := min(totalRows-startIndex, uint64(MaxRowsPerPartialResultSet))
			rows := s.ResultSet.Rows[startIndex : startIndex+rowCount]
			values := make([]*structpb.Value, len(rows)*len(s.ResultSet.Metadata.RowType.Fields))
			var idx int
			for _, row := range rows {
				for colIdx := range s.ResultSet.Metadata.RowType.Fields {
					values[idx] = row.Values[colIdx]
					idx++
				}
			}
			var rt []byte
			if len(s.ResumeTokens) == 0 {
				rt = EncodeResumeToken(startIndex + rowCount)
			} else {
				rt = s.ResumeTokens[startIndex]
			}
			result = append(result, &spannerpb.PartialResultSet{Metadata: s.ResultSet.Metadata, Values: values, ResumeToken: rt})
			startIndex += rowCount
			if startIndex == totalRows {
				break
			}
		}
	} else {
		result = append(result, &spannerpb.PartialResultSet{Metadata: s.ResultSet.Metadata})
	}
	return result, nil
}
func gologoo__min_705673bfbb22749846b5ab424c32ea33(x, y uint64) uint64 {
	if x > y {
		return y
	}
	return x
}
func (s *StatementResult) gologoo__updateCountToPartialResultSet_705673bfbb22749846b5ab424c32ea33(exact bool) *spannerpb.PartialResultSet {
	return &spannerpb.PartialResultSet{Stats: s.convertUpdateCountToResultSet(exact).Stats}
}
func (s *StatementResult) gologoo__convertUpdateCountToResultSet_705673bfbb22749846b5ab424c32ea33(exact bool) *spannerpb.ResultSet {
	if exact {
		return &spannerpb.ResultSet{Stats: &spannerpb.ResultSetStats{RowCount: &spannerpb.ResultSetStats_RowCountExact{RowCountExact: s.UpdateCount}}}
	}
	return &spannerpb.ResultSet{Stats: &spannerpb.ResultSetStats{RowCount: &spannerpb.ResultSetStats_RowCountLowerBound{RowCountLowerBound: s.UpdateCount}}}
}

type SimulatedExecutionTime struct {
	MinimumExecutionTime time.Duration
	RandomExecutionTime  time.Duration
	Errors               []error
	KeepError            bool
}
type InMemSpannerServer interface {
	spannerpb.SpannerServer
	Stop()
	Reset()
	SetError(err error)
	PutStatementResult(sql string, result *StatementResult) error
	PutPartitionResult(partitionToken []byte, result *StatementResult) error
	AddPartialResultSetError(sql string, err PartialResultSetExecutionTime)
	RemoveStatementResult(sql string)
	AbortTransaction(id []byte)
	PutExecutionTime(method string, executionTime SimulatedExecutionTime)
	Freeze()
	Unfreeze()
	TotalSessionsCreated() uint
	TotalSessionsDeleted() uint
	SetMaxSessionsReturnedByServerPerBatchRequest(sessionCount int32)
	SetMaxSessionsReturnedByServerInTotal(sessionCount int32)
	ReceivedRequests() chan interface {
	}
	DumpSessions() map[string]bool
	ClearPings()
	DumpPings() []string
}
type inMemSpannerServer struct {
	spannerpb.SpannerServer
	mu                                         sync.Mutex
	stopped                                    bool
	err                                        error
	sessionCounter                             uint64
	sessions                                   map[string]*spannerpb.Session
	sessionLastUseTime                         map[string]time.Time
	transactionCounters                        map[string]*uint64
	transactions                               map[string]*spannerpb.Transaction
	abortedTransactions                        map[string]bool
	partitionedDmlTransactions                 map[string]bool
	statementResults                           map[string]*StatementResult
	partitionResults                           map[string]*StatementResult
	executionTimes                             map[string]*SimulatedExecutionTime
	partialResultSetErrors                     map[string][]*PartialResultSetExecutionTime
	totalSessionsCreated                       uint
	totalSessionsDeleted                       uint
	maxSessionsReturnedByServerPerBatchRequest int32
	maxSessionsReturnedByServerInTotal         int32
	receivedRequests                           chan interface {
	}
	pings   []string
	freezed chan struct {
	}
}

func gologoo__NewInMemSpannerServer_705673bfbb22749846b5ab424c32ea33() InMemSpannerServer {
	res := &inMemSpannerServer{}
	res.initDefaults()
	res.statementResults = make(map[string]*StatementResult)
	res.partitionResults = make(map[string]*StatementResult)
	res.executionTimes = make(map[string]*SimulatedExecutionTime)
	res.partialResultSetErrors = make(map[string][]*PartialResultSetExecutionTime)
	res.receivedRequests = make(chan interface {
	}, 1000000)
	res.Freeze()
	res.Unfreeze()
	return res
}
func (s *inMemSpannerServer) gologoo__Stop_705673bfbb22749846b5ab424c32ea33() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopped = true
	close(s.receivedRequests)
}
func (s *inMemSpannerServer) gologoo__Reset_705673bfbb22749846b5ab424c32ea33() {
	s.mu.Lock()
	defer s.mu.Unlock()
	close(s.receivedRequests)
	s.receivedRequests = make(chan interface {
	}, 1000000)
	s.initDefaults()
}
func (s *inMemSpannerServer) gologoo__SetError_705673bfbb22749846b5ab424c32ea33(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.err = err
}
func (s *inMemSpannerServer) gologoo__PutStatementResult_705673bfbb22749846b5ab424c32ea33(sql string, result *StatementResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.statementResults[sql] = result
	return nil
}
func (s *inMemSpannerServer) gologoo__RemoveStatementResult_705673bfbb22749846b5ab424c32ea33(sql string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.statementResults, sql)
}
func (s *inMemSpannerServer) gologoo__PutPartitionResult_705673bfbb22749846b5ab424c32ea33(partitionToken []byte, result *StatementResult) error {
	tokenString := string(partitionToken)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.partitionResults[tokenString] = result
	return nil
}
func (s *inMemSpannerServer) gologoo__AbortTransaction_705673bfbb22749846b5ab424c32ea33(id []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.abortedTransactions[string(id)] = true
}
func (s *inMemSpannerServer) gologoo__PutExecutionTime_705673bfbb22749846b5ab424c32ea33(method string, executionTime SimulatedExecutionTime) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.executionTimes[method] = &executionTime
}
func (s *inMemSpannerServer) gologoo__AddPartialResultSetError_705673bfbb22749846b5ab424c32ea33(sql string, partialResultSetError PartialResultSetExecutionTime) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.partialResultSetErrors[sql] = append(s.partialResultSetErrors[sql], &partialResultSetError)
}
func (s *inMemSpannerServer) gologoo__Freeze_705673bfbb22749846b5ab424c32ea33() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.freezed = make(chan struct {
	})
}
func (s *inMemSpannerServer) gologoo__Unfreeze_705673bfbb22749846b5ab424c32ea33() {
	s.mu.Lock()
	defer s.mu.Unlock()
	close(s.freezed)
}
func (s *inMemSpannerServer) gologoo__ready_705673bfbb22749846b5ab424c32ea33() {
	s.mu.Lock()
	freezed := s.freezed
	s.mu.Unlock()
	<-freezed
}
func (s *inMemSpannerServer) gologoo__TotalSessionsCreated_705673bfbb22749846b5ab424c32ea33() uint {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.totalSessionsCreated
}
func (s *inMemSpannerServer) gologoo__TotalSessionsDeleted_705673bfbb22749846b5ab424c32ea33() uint {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.totalSessionsDeleted
}
func (s *inMemSpannerServer) gologoo__SetMaxSessionsReturnedByServerPerBatchRequest_705673bfbb22749846b5ab424c32ea33(sessionCount int32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.maxSessionsReturnedByServerPerBatchRequest = sessionCount
}
func (s *inMemSpannerServer) gologoo__SetMaxSessionsReturnedByServerInTotal_705673bfbb22749846b5ab424c32ea33(sessionCount int32) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.maxSessionsReturnedByServerInTotal = sessionCount
}
func (s *inMemSpannerServer) gologoo__ReceivedRequests_705673bfbb22749846b5ab424c32ea33() chan interface {
} {
	return s.receivedRequests
}
func (s *inMemSpannerServer) gologoo__ClearPings_705673bfbb22749846b5ab424c32ea33() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pings = nil
}
func (s *inMemSpannerServer) gologoo__DumpPings_705673bfbb22749846b5ab424c32ea33() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.pings...)
}
func (s *inMemSpannerServer) gologoo__DumpSessions_705673bfbb22749846b5ab424c32ea33() map[string]bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	st := map[string]bool{}
	for s := range s.sessions {
		st[s] = true
	}
	return st
}
func (s *inMemSpannerServer) gologoo__initDefaults_705673bfbb22749846b5ab424c32ea33() {
	s.sessionCounter = 0
	s.maxSessionsReturnedByServerPerBatchRequest = 100
	s.sessions = make(map[string]*spannerpb.Session)
	s.sessionLastUseTime = make(map[string]time.Time)
	s.transactions = make(map[string]*spannerpb.Transaction)
	s.abortedTransactions = make(map[string]bool)
	s.partitionedDmlTransactions = make(map[string]bool)
	s.transactionCounters = make(map[string]*uint64)
}
func (s *inMemSpannerServer) gologoo__generateSessionNameLocked_705673bfbb22749846b5ab424c32ea33(database string) string {
	s.sessionCounter++
	return fmt.Sprintf("%s/sessions/%d", database, s.sessionCounter)
}
func (s *inMemSpannerServer) gologoo__findSession_705673bfbb22749846b5ab424c32ea33(name string) (*spannerpb.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session := s.sessions[name]
	if session == nil {
		return nil, newSessionNotFoundError(name)
	}
	return session, nil
}

const sessionResourceType = "type.googleapis.com/google.spanner.v1.Session"

func gologoo__newSessionNotFoundError_705673bfbb22749846b5ab424c32ea33(name string) error {
	s := gstatus.Newf(codes.NotFound, "Session not found: Session with id %s not found", name)
	s, _ = s.WithDetails(&errdetails.ResourceInfo{ResourceType: sessionResourceType, ResourceName: name})
	return s.Err()
}
func (s *inMemSpannerServer) gologoo__updateSessionLastUseTime_705673bfbb22749846b5ab424c32ea33(session string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessionLastUseTime[session] = time.Now()
}
func gologoo__getCurrentTimestamp_705673bfbb22749846b5ab424c32ea33() *timestamp.Timestamp {
	t := time.Now()
	return &timestamp.Timestamp{Seconds: t.Unix(), Nanos: int32(t.Nanosecond())}
}
func (s *inMemSpannerServer) gologoo__getTransactionID_705673bfbb22749846b5ab424c32ea33(session *spannerpb.Session, txSelector *spannerpb.TransactionSelector) []byte {
	var res []byte
	if txSelector.GetBegin() != nil {
		res = s.beginTransaction(session, txSelector.GetBegin()).Id
	} else if txSelector.GetId() != nil {
		res = txSelector.GetId()
	}
	return res
}
func (s *inMemSpannerServer) gologoo__generateTransactionName_705673bfbb22749846b5ab424c32ea33(session string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	counter, ok := s.transactionCounters[session]
	if !ok {
		counter = new(uint64)
		s.transactionCounters[session] = counter
	}
	*counter++
	return fmt.Sprintf("%s/transactions/%d", session, *counter)
}
func (s *inMemSpannerServer) gologoo__beginTransaction_705673bfbb22749846b5ab424c32ea33(session *spannerpb.Session, options *spannerpb.TransactionOptions) *spannerpb.Transaction {
	id := s.generateTransactionName(session.Name)
	res := &spannerpb.Transaction{Id: []byte(id), ReadTimestamp: getCurrentTimestamp()}
	s.mu.Lock()
	s.transactions[id] = res
	s.partitionedDmlTransactions[id] = options.GetPartitionedDml() != nil
	s.mu.Unlock()
	return res
}
func (s *inMemSpannerServer) gologoo__getTransactionByID_705673bfbb22749846b5ab424c32ea33(id []byte) (*spannerpb.Transaction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	tx, ok := s.transactions[string(id)]
	if !ok {
		return nil, gstatus.Error(codes.NotFound, "Transaction not found")
	}
	aborted, ok := s.abortedTransactions[string(id)]
	if ok && aborted {
		return nil, newAbortedErrorWithMinimalRetryDelay()
	}
	return tx, nil
}
func gologoo__newAbortedErrorWithMinimalRetryDelay_705673bfbb22749846b5ab424c32ea33() error {
	st := gstatus.New(codes.Aborted, "Transaction has been aborted")
	retry := &errdetails.RetryInfo{RetryDelay: ptypes.DurationProto(time.Nanosecond)}
	st, _ = st.WithDetails(retry)
	return st.Err()
}
func (s *inMemSpannerServer) gologoo__removeTransaction_705673bfbb22749846b5ab424c32ea33(tx *spannerpb.Transaction) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.transactions, string(tx.Id))
	delete(s.partitionedDmlTransactions, string(tx.Id))
}
func (s *inMemSpannerServer) gologoo__getPartitionResult_705673bfbb22749846b5ab424c32ea33(partitionToken []byte) (*StatementResult, error) {
	tokenString := string(partitionToken)
	s.mu.Lock()
	defer s.mu.Unlock()
	result, ok := s.partitionResults[tokenString]
	if !ok {
		return nil, gstatus.Error(codes.Internal, fmt.Sprintf("No result found for partition token %v", tokenString))
	}
	return result, nil
}
func (s *inMemSpannerServer) gologoo__getStatementResult_705673bfbb22749846b5ab424c32ea33(sql string) (*StatementResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result, ok := s.statementResults[sql]
	if !ok {
		return nil, gstatus.Error(codes.Internal, fmt.Sprintf("No result found for statement %v", sql))
	}
	return result, nil
}
func (s *inMemSpannerServer) gologoo__simulateExecutionTime_705673bfbb22749846b5ab424c32ea33(method string, req interface {
}) error {
	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		return gstatus.Error(codes.Unavailable, "server has been stopped")
	}
	s.receivedRequests <- req
	s.mu.Unlock()
	s.ready()
	s.mu.Lock()
	if s.err != nil {
		err := s.err
		s.err = nil
		s.mu.Unlock()
		return err
	}
	executionTime, ok := s.executionTimes[method]
	s.mu.Unlock()
	if ok {
		var randTime int64
		if executionTime.RandomExecutionTime > 0 {
			randTime = rand.Int63n(int64(executionTime.RandomExecutionTime))
		}
		totalExecutionTime := time.Duration(int64(executionTime.MinimumExecutionTime) + randTime)
		<-time.After(totalExecutionTime)
		s.mu.Lock()
		if executionTime.Errors != nil && len(executionTime.Errors) > 0 {
			err := executionTime.Errors[0]
			if !executionTime.KeepError {
				executionTime.Errors = executionTime.Errors[1:]
			}
			s.mu.Unlock()
			return err
		}
		s.mu.Unlock()
	}
	return nil
}
func (s *inMemSpannerServer) gologoo__CreateSession_705673bfbb22749846b5ab424c32ea33(ctx context.Context, req *spannerpb.CreateSessionRequest) (*spannerpb.Session, error) {
	if err := s.simulateExecutionTime(MethodCreateSession, req); err != nil {
		return nil, err
	}
	if req.Database == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing database")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.maxSessionsReturnedByServerInTotal > int32(0) && int32(len(s.sessions)) == s.maxSessionsReturnedByServerInTotal {
		return nil, gstatus.Error(codes.ResourceExhausted, "No more sessions available")
	}
	sessionName := s.generateSessionNameLocked(req.Database)
	ts := getCurrentTimestamp()
	session := &spannerpb.Session{Name: sessionName, CreateTime: ts, ApproximateLastUseTime: ts}
	s.totalSessionsCreated++
	s.sessions[sessionName] = session
	return session, nil
}
func (s *inMemSpannerServer) gologoo__BatchCreateSessions_705673bfbb22749846b5ab424c32ea33(ctx context.Context, req *spannerpb.BatchCreateSessionsRequest) (*spannerpb.BatchCreateSessionsResponse, error) {
	if err := s.simulateExecutionTime(MethodBatchCreateSession, req); err != nil {
		return nil, err
	}
	if req.Database == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing database")
	}
	if req.SessionCount <= 0 {
		return nil, gstatus.Error(codes.InvalidArgument, "Session count must be >= 0")
	}
	sessionsToCreate := req.SessionCount
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.maxSessionsReturnedByServerInTotal > int32(0) && int32(len(s.sessions)) >= s.maxSessionsReturnedByServerInTotal {
		return nil, gstatus.Error(codes.ResourceExhausted, "No more sessions available")
	}
	if sessionsToCreate > s.maxSessionsReturnedByServerPerBatchRequest {
		sessionsToCreate = s.maxSessionsReturnedByServerPerBatchRequest
	}
	if s.maxSessionsReturnedByServerInTotal > int32(0) && (sessionsToCreate+int32(len(s.sessions))) > s.maxSessionsReturnedByServerInTotal {
		sessionsToCreate = s.maxSessionsReturnedByServerInTotal - int32(len(s.sessions))
	}
	sessions := make([]*spannerpb.Session, sessionsToCreate)
	for i := int32(0); i < sessionsToCreate; i++ {
		sessionName := s.generateSessionNameLocked(req.Database)
		ts := getCurrentTimestamp()
		sessions[i] = &spannerpb.Session{Name: sessionName, CreateTime: ts, ApproximateLastUseTime: ts}
		s.totalSessionsCreated++
		s.sessions[sessionName] = sessions[i]
	}
	header := metadata.New(map[string]string{"server-timing": "gfet4t7; dur=123"})
	if err := grpc.SendHeader(ctx, header); err != nil {
		return nil, gstatus.Errorf(codes.Internal, "unable to send 'server-timing' header")
	}
	return &spannerpb.BatchCreateSessionsResponse{Session: sessions}, nil
}
func (s *inMemSpannerServer) gologoo__GetSession_705673bfbb22749846b5ab424c32ea33(ctx context.Context, req *spannerpb.GetSessionRequest) (*spannerpb.Session, error) {
	if err := s.simulateExecutionTime(MethodGetSession, req); err != nil {
		return nil, err
	}
	if req.Name == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing session name")
	}
	session, err := s.findSession(req.Name)
	if err != nil {
		return nil, err
	}
	return session, nil
}
func (s *inMemSpannerServer) gologoo__ListSessions_705673bfbb22749846b5ab424c32ea33(ctx context.Context, req *spannerpb.ListSessionsRequest) (*spannerpb.ListSessionsResponse, error) {
	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		return nil, gstatus.Error(codes.Unavailable, "server has been stopped")
	}
	s.receivedRequests <- req
	s.mu.Unlock()
	if req.Database == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing database")
	}
	expectedSessionName := req.Database + "/sessions/"
	var sessions []*spannerpb.Session
	s.mu.Lock()
	for _, session := range s.sessions {
		if strings.Index(session.Name, expectedSessionName) == 0 {
			sessions = append(sessions, session)
		}
	}
	s.mu.Unlock()
	sort.Slice(sessions[:], func(i, j int) bool {
		return sessions[i].Name < sessions[j].Name
	})
	res := &spannerpb.ListSessionsResponse{Sessions: sessions}
	return res, nil
}
func (s *inMemSpannerServer) gologoo__DeleteSession_705673bfbb22749846b5ab424c32ea33(ctx context.Context, req *spannerpb.DeleteSessionRequest) (*emptypb.Empty, error) {
	if err := s.simulateExecutionTime(MethodDeleteSession, req); err != nil {
		return nil, err
	}
	if req.Name == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing session name")
	}
	if _, err := s.findSession(req.Name); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.totalSessionsDeleted++
	delete(s.sessions, req.Name)
	return &emptypb.Empty{}, nil
}
func (s *inMemSpannerServer) gologoo__ExecuteSql_705673bfbb22749846b5ab424c32ea33(ctx context.Context, req *spannerpb.ExecuteSqlRequest) (*spannerpb.ResultSet, error) {
	if err := s.simulateExecutionTime(MethodExecuteSql, req); err != nil {
		return nil, err
	}
	if req.Sql == "SELECT 1" {
		s.mu.Lock()
		s.pings = append(s.pings, req.Session)
		s.mu.Unlock()
	}
	if req.Session == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing session name")
	}
	session, err := s.findSession(req.Session)
	if err != nil {
		return nil, err
	}
	var id []byte
	s.updateSessionLastUseTime(session.Name)
	if id = s.getTransactionID(session, req.Transaction); id != nil {
		_, err = s.getTransactionByID(id)
		if err != nil {
			return nil, err
		}
	}
	var statementResult *StatementResult
	if req.PartitionToken != nil {
		statementResult, err = s.getPartitionResult(req.PartitionToken)
	} else {
		statementResult, err = s.getStatementResult(req.Sql)
	}
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	isPartitionedDml := s.partitionedDmlTransactions[string(id)]
	s.mu.Unlock()
	switch statementResult.Type {
	case StatementResultError:
		return nil, statementResult.Err
	case StatementResultResultSet:
		return statementResult.ResultSet, nil
	case StatementResultUpdateCount:
		return statementResult.convertUpdateCountToResultSet(!isPartitionedDml), nil
	}
	return nil, gstatus.Error(codes.Internal, "Unknown result type")
}
func (s *inMemSpannerServer) gologoo__ExecuteStreamingSql_705673bfbb22749846b5ab424c32ea33(req *spannerpb.ExecuteSqlRequest, stream spannerpb.Spanner_ExecuteStreamingSqlServer) error {
	if err := s.simulateExecutionTime(MethodExecuteStreamingSql, req); err != nil {
		return err
	}
	return s.executeStreamingSQL(req, stream)
}
func (s *inMemSpannerServer) gologoo__executeStreamingSQL_705673bfbb22749846b5ab424c32ea33(req *spannerpb.ExecuteSqlRequest, stream spannerpb.Spanner_ExecuteStreamingSqlServer) error {
	if req.Session == "" {
		return gstatus.Error(codes.InvalidArgument, "Missing session name")
	}
	session, err := s.findSession(req.Session)
	if err != nil {
		return err
	}
	s.updateSessionLastUseTime(session.Name)
	var id []byte
	if id = s.getTransactionID(session, req.Transaction); id != nil {
		_, err = s.getTransactionByID(id)
		if err != nil {
			return err
		}
	}
	var statementResult *StatementResult
	if req.PartitionToken != nil {
		statementResult, err = s.getPartitionResult(req.PartitionToken)
	} else {
		statementResult, err = s.getStatementResult(req.Sql)
	}
	if err != nil {
		return err
	}
	s.mu.Lock()
	isPartitionedDml := s.partitionedDmlTransactions[string(id)]
	s.mu.Unlock()
	switch statementResult.Type {
	case StatementResultError:
		return statementResult.Err
	case StatementResultResultSet:
		parts, err := statementResult.ToPartialResultSets(req.ResumeToken)
		if err != nil {
			return err
		}
		var nextPartialResultSetError *PartialResultSetExecutionTime
		s.mu.Lock()
		pErrors := s.partialResultSetErrors[req.Sql]
		if len(pErrors) > 0 {
			nextPartialResultSetError = pErrors[0]
			s.partialResultSetErrors[req.Sql] = pErrors[1:]
		}
		s.mu.Unlock()
		for _, part := range parts {
			if nextPartialResultSetError != nil && bytes.Equal(part.ResumeToken, nextPartialResultSetError.ResumeToken) {
				if nextPartialResultSetError.ExecutionTime > 0 {
					<-time.After(nextPartialResultSetError.ExecutionTime)
				}
				if nextPartialResultSetError.Err != nil {
					return nextPartialResultSetError.Err
				}
			}
			if err := stream.Send(part); err != nil {
				return err
			}
		}
		return nil
	case StatementResultUpdateCount:
		part := statementResult.updateCountToPartialResultSet(!isPartitionedDml)
		if err := stream.Send(part); err != nil {
			return err
		}
		return nil
	}
	return gstatus.Error(codes.Internal, "Unknown result type")
}
func (s *inMemSpannerServer) gologoo__ExecuteBatchDml_705673bfbb22749846b5ab424c32ea33(ctx context.Context, req *spannerpb.ExecuteBatchDmlRequest) (*spannerpb.ExecuteBatchDmlResponse, error) {
	if err := s.simulateExecutionTime(MethodExecuteBatchDml, req); err != nil {
		return nil, err
	}
	if req.Session == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing session name")
	}
	session, err := s.findSession(req.Session)
	if err != nil {
		return nil, err
	}
	s.updateSessionLastUseTime(session.Name)
	var id []byte
	if id = s.getTransactionID(session, req.Transaction); id != nil {
		_, err = s.getTransactionByID(id)
		if err != nil {
			return nil, err
		}
	}
	s.mu.Lock()
	isPartitionedDml := s.partitionedDmlTransactions[string(id)]
	s.mu.Unlock()
	resp := &spannerpb.ExecuteBatchDmlResponse{}
	resp.ResultSets = make([]*spannerpb.ResultSet, len(req.Statements))
	resp.Status = &status.Status{Code: int32(codes.OK)}
	for idx, batchStatement := range req.Statements {
		statementResult, err := s.getStatementResult(batchStatement.Sql)
		if err != nil {
			return nil, err
		}
		switch statementResult.Type {
		case StatementResultError:
			resp.Status = &status.Status{Code: int32(gstatus.Code(statementResult.Err)), Message: statementResult.Err.Error()}
			resp.ResultSets = resp.ResultSets[:idx]
			return resp, nil
		case StatementResultResultSet:
			return nil, gstatus.Error(codes.InvalidArgument, fmt.Sprintf("Not an update statement: %v", batchStatement.Sql))
		case StatementResultUpdateCount:
			resp.ResultSets[idx] = statementResult.convertUpdateCountToResultSet(!isPartitionedDml)
		}
	}
	return resp, nil
}
func (s *inMemSpannerServer) gologoo__Read_705673bfbb22749846b5ab424c32ea33(ctx context.Context, req *spannerpb.ReadRequest) (*spannerpb.ResultSet, error) {
	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		return nil, gstatus.Error(codes.Unavailable, "server has been stopped")
	}
	s.receivedRequests <- req
	s.mu.Unlock()
	header := metadata.New(map[string]string{"server-timing": "gfet4t7; dur=123"})
	if err := grpc.SendHeader(ctx, header); err != nil {
		return nil, gstatus.Errorf(codes.Internal, "unable to send 'server-timing' header")
	}
	return nil, gstatus.Error(codes.Unimplemented, "Method not yet implemented")
}
func (s *inMemSpannerServer) gologoo__StreamingRead_705673bfbb22749846b5ab424c32ea33(req *spannerpb.ReadRequest, stream spannerpb.Spanner_StreamingReadServer) error {
	if err := s.simulateExecutionTime(MethodStreamingRead, req); err != nil {
		return err
	}
	sqlReq := &spannerpb.ExecuteSqlRequest{Session: req.Session, Transaction: req.Transaction, PartitionToken: req.PartitionToken, ResumeToken: req.ResumeToken, Sql: fmt.Sprintf("SELECT %s FROM %s", strings.Join(req.Columns, ", "), req.Table)}
	return s.executeStreamingSQL(sqlReq, stream)
}
func (s *inMemSpannerServer) gologoo__BeginTransaction_705673bfbb22749846b5ab424c32ea33(ctx context.Context, req *spannerpb.BeginTransactionRequest) (*spannerpb.Transaction, error) {
	if err := s.simulateExecutionTime(MethodBeginTransaction, req); err != nil {
		return nil, err
	}
	if req.Session == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing session name")
	}
	session, err := s.findSession(req.Session)
	if err != nil {
		return nil, err
	}
	s.updateSessionLastUseTime(session.Name)
	tx := s.beginTransaction(session, req.Options)
	return tx, nil
}
func (s *inMemSpannerServer) gologoo__Commit_705673bfbb22749846b5ab424c32ea33(ctx context.Context, req *spannerpb.CommitRequest) (*spannerpb.CommitResponse, error) {
	if err := s.simulateExecutionTime(MethodCommitTransaction, req); err != nil {
		return nil, err
	}
	if req.Session == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing session name")
	}
	session, err := s.findSession(req.Session)
	if err != nil {
		return nil, err
	}
	s.updateSessionLastUseTime(session.Name)
	var tx *spannerpb.Transaction
	if req.GetSingleUseTransaction() != nil {
		tx = s.beginTransaction(session, req.GetSingleUseTransaction())
	} else if req.GetTransactionId() != nil {
		tx, err = s.getTransactionByID(req.GetTransactionId())
		if err != nil {
			return nil, err
		}
	} else {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing transaction in commit request")
	}
	s.removeTransaction(tx)
	resp := &spannerpb.CommitResponse{CommitTimestamp: getCurrentTimestamp()}
	if req.ReturnCommitStats {
		resp.CommitStats = &spannerpb.CommitResponse_CommitStats{MutationCount: int64(1)}
	}
	return resp, nil
}
func (s *inMemSpannerServer) gologoo__Rollback_705673bfbb22749846b5ab424c32ea33(ctx context.Context, req *spannerpb.RollbackRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		return nil, gstatus.Error(codes.Unavailable, "server has been stopped")
	}
	s.receivedRequests <- req
	s.mu.Unlock()
	if req.Session == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing session name")
	}
	session, err := s.findSession(req.Session)
	if err != nil {
		return nil, err
	}
	s.updateSessionLastUseTime(session.Name)
	tx, err := s.getTransactionByID(req.TransactionId)
	if err != nil {
		return nil, err
	}
	s.removeTransaction(tx)
	return &emptypb.Empty{}, nil
}
func (s *inMemSpannerServer) gologoo__PartitionQuery_705673bfbb22749846b5ab424c32ea33(ctx context.Context, req *spannerpb.PartitionQueryRequest) (*spannerpb.PartitionResponse, error) {
	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		return nil, gstatus.Error(codes.Unavailable, "server has been stopped")
	}
	s.receivedRequests <- req
	s.mu.Unlock()
	if req.Session == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing session name")
	}
	session, err := s.findSession(req.Session)
	if err != nil {
		return nil, err
	}
	var id []byte
	var tx *spannerpb.Transaction
	s.updateSessionLastUseTime(session.Name)
	if id = s.getTransactionID(session, req.Transaction); id != nil {
		tx, err = s.getTransactionByID(id)
		if err != nil {
			return nil, err
		}
	}
	var partitions []*spannerpb.Partition
	for i := int64(0); i < req.PartitionOptions.MaxPartitions; i++ {
		token := make([]byte, 10)
		_, err := rand.Read(token)
		if err != nil {
			return nil, gstatus.Error(codes.Internal, "failed to generate random partition token")
		}
		partitions = append(partitions, &spannerpb.Partition{PartitionToken: token})
	}
	return &spannerpb.PartitionResponse{Partitions: partitions, Transaction: tx}, nil
}
func (s *inMemSpannerServer) gologoo__PartitionRead_705673bfbb22749846b5ab424c32ea33(ctx context.Context, req *spannerpb.PartitionReadRequest) (*spannerpb.PartitionResponse, error) {
	return s.PartitionQuery(ctx, &spannerpb.PartitionQueryRequest{Session: req.Session, Transaction: req.Transaction, PartitionOptions: req.PartitionOptions, Sql: fmt.Sprintf("SELECT %s FROM %s", strings.Join(req.Columns, ", "), req.Table)})
}
func gologoo__EncodeResumeToken_705673bfbb22749846b5ab424c32ea33(t uint64) []byte {
	rt := make([]byte, 16)
	binary.PutUvarint(rt, t)
	return rt
}
func gologoo__DecodeResumeToken_705673bfbb22749846b5ab424c32ea33(t []byte) (uint64, error) {
	s, n := binary.Uvarint(t)
	if n <= 0 {
		return 0, fmt.Errorf("invalid resume token: %v", t)
	}
	return s, nil
}
func (s *StatementResult) ToPartialResultSets(resumeToken []byte) (result []*spannerpb.PartialResultSet, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ToPartialResultSets_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", resumeToken)
	result, err = s.gologoo__ToPartialResultSets_705673bfbb22749846b5ab424c32ea33(resumeToken)
	log.Printf("Output: %v %v\n", result, err)
	return
}
func min(x, y uint64) uint64 {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__min_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", x, y)
	r0 := gologoo__min_705673bfbb22749846b5ab424c32ea33(x, y)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *StatementResult) updateCountToPartialResultSet(exact bool) *spannerpb.PartialResultSet {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__updateCountToPartialResultSet_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", exact)
	r0 := s.gologoo__updateCountToPartialResultSet_705673bfbb22749846b5ab424c32ea33(exact)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *StatementResult) convertUpdateCountToResultSet(exact bool) *spannerpb.ResultSet {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__convertUpdateCountToResultSet_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", exact)
	r0 := s.gologoo__convertUpdateCountToResultSet_705673bfbb22749846b5ab424c32ea33(exact)
	log.Printf("Output: %v\n", r0)
	return r0
}
func NewInMemSpannerServer() InMemSpannerServer {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NewInMemSpannerServer_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : (none)\n")
	r0 := gologoo__NewInMemSpannerServer_705673bfbb22749846b5ab424c32ea33()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) Stop() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Stop_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : (none)\n")
	s.gologoo__Stop_705673bfbb22749846b5ab424c32ea33()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemSpannerServer) Reset() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Reset_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : (none)\n")
	s.gologoo__Reset_705673bfbb22749846b5ab424c32ea33()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemSpannerServer) SetError(err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SetError_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", err)
	s.gologoo__SetError_705673bfbb22749846b5ab424c32ea33(err)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemSpannerServer) PutStatementResult(sql string, result *StatementResult) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PutStatementResult_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", sql, result)
	r0 := s.gologoo__PutStatementResult_705673bfbb22749846b5ab424c32ea33(sql, result)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) RemoveStatementResult(sql string) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__RemoveStatementResult_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", sql)
	s.gologoo__RemoveStatementResult_705673bfbb22749846b5ab424c32ea33(sql)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemSpannerServer) PutPartitionResult(partitionToken []byte, result *StatementResult) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PutPartitionResult_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", partitionToken, result)
	r0 := s.gologoo__PutPartitionResult_705673bfbb22749846b5ab424c32ea33(partitionToken, result)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) AbortTransaction(id []byte) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__AbortTransaction_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", id)
	s.gologoo__AbortTransaction_705673bfbb22749846b5ab424c32ea33(id)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemSpannerServer) PutExecutionTime(method string, executionTime SimulatedExecutionTime) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PutExecutionTime_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", method, executionTime)
	s.gologoo__PutExecutionTime_705673bfbb22749846b5ab424c32ea33(method, executionTime)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemSpannerServer) AddPartialResultSetError(sql string, partialResultSetError PartialResultSetExecutionTime) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__AddPartialResultSetError_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", sql, partialResultSetError)
	s.gologoo__AddPartialResultSetError_705673bfbb22749846b5ab424c32ea33(sql, partialResultSetError)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemSpannerServer) Freeze() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Freeze_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : (none)\n")
	s.gologoo__Freeze_705673bfbb22749846b5ab424c32ea33()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemSpannerServer) Unfreeze() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Unfreeze_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : (none)\n")
	s.gologoo__Unfreeze_705673bfbb22749846b5ab424c32ea33()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemSpannerServer) ready() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ready_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : (none)\n")
	s.gologoo__ready_705673bfbb22749846b5ab424c32ea33()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemSpannerServer) TotalSessionsCreated() uint {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__TotalSessionsCreated_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__TotalSessionsCreated_705673bfbb22749846b5ab424c32ea33()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) TotalSessionsDeleted() uint {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__TotalSessionsDeleted_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__TotalSessionsDeleted_705673bfbb22749846b5ab424c32ea33()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) SetMaxSessionsReturnedByServerPerBatchRequest(sessionCount int32) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SetMaxSessionsReturnedByServerPerBatchRequest_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", sessionCount)
	s.gologoo__SetMaxSessionsReturnedByServerPerBatchRequest_705673bfbb22749846b5ab424c32ea33(sessionCount)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemSpannerServer) SetMaxSessionsReturnedByServerInTotal(sessionCount int32) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SetMaxSessionsReturnedByServerInTotal_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", sessionCount)
	s.gologoo__SetMaxSessionsReturnedByServerInTotal_705673bfbb22749846b5ab424c32ea33(sessionCount)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemSpannerServer) ReceivedRequests() chan interface {
} {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ReceivedRequests_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__ReceivedRequests_705673bfbb22749846b5ab424c32ea33()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) ClearPings() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ClearPings_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : (none)\n")
	s.gologoo__ClearPings_705673bfbb22749846b5ab424c32ea33()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemSpannerServer) DumpPings() []string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DumpPings_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__DumpPings_705673bfbb22749846b5ab424c32ea33()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) DumpSessions() map[string]bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DumpSessions_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__DumpSessions_705673bfbb22749846b5ab424c32ea33()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) initDefaults() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__initDefaults_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : (none)\n")
	s.gologoo__initDefaults_705673bfbb22749846b5ab424c32ea33()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemSpannerServer) generateSessionNameLocked(database string) string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__generateSessionNameLocked_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", database)
	r0 := s.gologoo__generateSessionNameLocked_705673bfbb22749846b5ab424c32ea33(database)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) findSession(name string) (*spannerpb.Session, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__findSession_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", name)
	r0, r1 := s.gologoo__findSession_705673bfbb22749846b5ab424c32ea33(name)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func newSessionNotFoundError(name string) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__newSessionNotFoundError_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", name)
	r0 := gologoo__newSessionNotFoundError_705673bfbb22749846b5ab424c32ea33(name)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) updateSessionLastUseTime(session string) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__updateSessionLastUseTime_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", session)
	s.gologoo__updateSessionLastUseTime_705673bfbb22749846b5ab424c32ea33(session)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func getCurrentTimestamp() *timestamp.Timestamp {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getCurrentTimestamp_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : (none)\n")
	r0 := gologoo__getCurrentTimestamp_705673bfbb22749846b5ab424c32ea33()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) getTransactionID(session *spannerpb.Session, txSelector *spannerpb.TransactionSelector) []byte {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getTransactionID_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", session, txSelector)
	r0 := s.gologoo__getTransactionID_705673bfbb22749846b5ab424c32ea33(session, txSelector)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) generateTransactionName(session string) string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__generateTransactionName_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", session)
	r0 := s.gologoo__generateTransactionName_705673bfbb22749846b5ab424c32ea33(session)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) beginTransaction(session *spannerpb.Session, options *spannerpb.TransactionOptions) *spannerpb.Transaction {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__beginTransaction_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", session, options)
	r0 := s.gologoo__beginTransaction_705673bfbb22749846b5ab424c32ea33(session, options)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) getTransactionByID(id []byte) (*spannerpb.Transaction, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getTransactionByID_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", id)
	r0, r1 := s.gologoo__getTransactionByID_705673bfbb22749846b5ab424c32ea33(id)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func newAbortedErrorWithMinimalRetryDelay() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__newAbortedErrorWithMinimalRetryDelay_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : (none)\n")
	r0 := gologoo__newAbortedErrorWithMinimalRetryDelay_705673bfbb22749846b5ab424c32ea33()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) removeTransaction(tx *spannerpb.Transaction) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__removeTransaction_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", tx)
	s.gologoo__removeTransaction_705673bfbb22749846b5ab424c32ea33(tx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *inMemSpannerServer) getPartitionResult(partitionToken []byte) (*StatementResult, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getPartitionResult_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", partitionToken)
	r0, r1 := s.gologoo__getPartitionResult_705673bfbb22749846b5ab424c32ea33(partitionToken)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *inMemSpannerServer) getStatementResult(sql string) (*StatementResult, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getStatementResult_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", sql)
	r0, r1 := s.gologoo__getStatementResult_705673bfbb22749846b5ab424c32ea33(sql)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *inMemSpannerServer) simulateExecutionTime(method string, req interface {
}) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__simulateExecutionTime_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", method, req)
	r0 := s.gologoo__simulateExecutionTime_705673bfbb22749846b5ab424c32ea33(method, req)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) CreateSession(ctx context.Context, req *spannerpb.CreateSessionRequest) (*spannerpb.Session, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateSession_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__CreateSession_705673bfbb22749846b5ab424c32ea33(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *inMemSpannerServer) BatchCreateSessions(ctx context.Context, req *spannerpb.BatchCreateSessionsRequest) (*spannerpb.BatchCreateSessionsResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__BatchCreateSessions_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__BatchCreateSessions_705673bfbb22749846b5ab424c32ea33(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *inMemSpannerServer) GetSession(ctx context.Context, req *spannerpb.GetSessionRequest) (*spannerpb.Session, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetSession_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__GetSession_705673bfbb22749846b5ab424c32ea33(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *inMemSpannerServer) ListSessions(ctx context.Context, req *spannerpb.ListSessionsRequest) (*spannerpb.ListSessionsResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ListSessions_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__ListSessions_705673bfbb22749846b5ab424c32ea33(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *inMemSpannerServer) DeleteSession(ctx context.Context, req *spannerpb.DeleteSessionRequest) (*emptypb.Empty, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DeleteSession_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__DeleteSession_705673bfbb22749846b5ab424c32ea33(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *inMemSpannerServer) ExecuteSql(ctx context.Context, req *spannerpb.ExecuteSqlRequest) (*spannerpb.ResultSet, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ExecuteSql_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__ExecuteSql_705673bfbb22749846b5ab424c32ea33(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *inMemSpannerServer) ExecuteStreamingSql(req *spannerpb.ExecuteSqlRequest, stream spannerpb.Spanner_ExecuteStreamingSqlServer) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ExecuteStreamingSql_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", req, stream)
	r0 := s.gologoo__ExecuteStreamingSql_705673bfbb22749846b5ab424c32ea33(req, stream)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) executeStreamingSQL(req *spannerpb.ExecuteSqlRequest, stream spannerpb.Spanner_ExecuteStreamingSqlServer) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__executeStreamingSQL_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", req, stream)
	r0 := s.gologoo__executeStreamingSQL_705673bfbb22749846b5ab424c32ea33(req, stream)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) ExecuteBatchDml(ctx context.Context, req *spannerpb.ExecuteBatchDmlRequest) (*spannerpb.ExecuteBatchDmlResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ExecuteBatchDml_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__ExecuteBatchDml_705673bfbb22749846b5ab424c32ea33(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *inMemSpannerServer) Read(ctx context.Context, req *spannerpb.ReadRequest) (*spannerpb.ResultSet, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Read_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__Read_705673bfbb22749846b5ab424c32ea33(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *inMemSpannerServer) StreamingRead(req *spannerpb.ReadRequest, stream spannerpb.Spanner_StreamingReadServer) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__StreamingRead_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", req, stream)
	r0 := s.gologoo__StreamingRead_705673bfbb22749846b5ab424c32ea33(req, stream)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *inMemSpannerServer) BeginTransaction(ctx context.Context, req *spannerpb.BeginTransactionRequest) (*spannerpb.Transaction, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__BeginTransaction_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__BeginTransaction_705673bfbb22749846b5ab424c32ea33(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *inMemSpannerServer) Commit(ctx context.Context, req *spannerpb.CommitRequest) (*spannerpb.CommitResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Commit_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__Commit_705673bfbb22749846b5ab424c32ea33(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *inMemSpannerServer) Rollback(ctx context.Context, req *spannerpb.RollbackRequest) (*emptypb.Empty, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Rollback_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__Rollback_705673bfbb22749846b5ab424c32ea33(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *inMemSpannerServer) PartitionQuery(ctx context.Context, req *spannerpb.PartitionQueryRequest) (*spannerpb.PartitionResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PartitionQuery_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__PartitionQuery_705673bfbb22749846b5ab424c32ea33(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *inMemSpannerServer) PartitionRead(ctx context.Context, req *spannerpb.PartitionReadRequest) (*spannerpb.PartitionResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__PartitionRead_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__PartitionRead_705673bfbb22749846b5ab424c32ea33(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func EncodeResumeToken(t uint64) []byte {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__EncodeResumeToken_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", t)
	r0 := gologoo__EncodeResumeToken_705673bfbb22749846b5ab424c32ea33(t)
	log.Printf("Output: %v\n", r0)
	return r0
}
func DecodeResumeToken(t []byte) (uint64, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DecodeResumeToken_705673bfbb22749846b5ab424c32ea33")
	log.Printf("Input : %v\n", t)
	r0, r1 := gologoo__DecodeResumeToken_705673bfbb22749846b5ab424c32ea33(t)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
