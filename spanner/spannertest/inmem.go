package spannertest

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"cloud.google.com/go/civil"
	"cloud.google.com/go/spanner/spansql"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	anypb "github.com/golang/protobuf/ptypes/any"
	emptypb "github.com/golang/protobuf/ptypes/empty"
	structpb "github.com/golang/protobuf/ptypes/struct"
	timestamppb "github.com/golang/protobuf/ptypes/timestamp"
	lropb "google.golang.org/genproto/googleapis/longrunning"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	spannerpb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	Addr string
	l    net.Listener
	srv  *grpc.Server
	s    *server
}
type server struct {
	logf     Logger
	db       database
	start    time.Time
	mu       sync.Mutex
	sessions map[string]*session
	lros     map[string]*lro
	adminpb.DatabaseAdminServer
	spannerpb.SpannerServer
	lropb.OperationsServer
}
type session struct {
	name         string
	creation     time.Time
	ctx          context.Context
	cancel       func()
	mu           sync.Mutex
	lastUse      time.Time
	transactions map[string]*transaction
}

func (s *session) gologoo__Proto_5d9035705d48cffd75ca99cea5c61bc1() *spannerpb.Session {
	s.mu.Lock()
	defer s.mu.Unlock()
	m := &spannerpb.Session{Name: s.name, CreateTime: timestampProto(s.creation), ApproximateLastUseTime: timestampProto(s.lastUse)}
	return m
}
func gologoo__timestampProto_5d9035705d48cffd75ca99cea5c61bc1(t time.Time) *timestamppb.Timestamp {
	if t.IsZero() {
		return nil
	}
	ts, err := ptypes.TimestampProto(t)
	if err != nil {
		return nil
	}
	return ts
}

type lro struct {
	mu    sync.Mutex
	state *lropb.Operation
	waitc chan struct {
	}
	waitatom int32
}

func gologoo__newLRO_5d9035705d48cffd75ca99cea5c61bc1(initState *lropb.Operation) *lro {
	return &lro{state: initState, waitc: make(chan struct {
	})}
}
func (l *lro) gologoo__noWait_5d9035705d48cffd75ca99cea5c61bc1() {
	if atomic.CompareAndSwapInt32(&l.waitatom, 0, 1) {
		close(l.waitc)
	}
}
func (l *lro) gologoo__State_5d9035705d48cffd75ca99cea5c61bc1() *lropb.Operation {
	l.mu.Lock()
	defer l.mu.Unlock()
	return proto.Clone(l.state).(*lropb.Operation)
}

type Logger func(format string, args ...interface {
})

func gologoo__NewServer_5d9035705d48cffd75ca99cea5c61bc1(laddr string) (*Server, error) {
	l, err := net.Listen("tcp", laddr)
	if err != nil {
		return nil, err
	}
	s := &Server{Addr: l.Addr().String(), l: l, srv: grpc.NewServer(), s: &server{logf: func(format string, args ...interface {
	}) {
		log.Printf("spannertest.inmem: "+format, args...)
	}, start: time.Now(), sessions: make(map[string]*session), lros: make(map[string]*lro)}}
	adminpb.RegisterDatabaseAdminServer(s.srv, s.s)
	spannerpb.RegisterSpannerServer(s.srv, s.s)
	lropb.RegisterOperationsServer(s.srv, s.s)
	go s.srv.Serve(s.l)
	return s, nil
}
func (s *Server) gologoo__SetLogger_5d9035705d48cffd75ca99cea5c61bc1(l Logger) {
	s.s.logf = l
}
func (s *Server) gologoo__Close_5d9035705d48cffd75ca99cea5c61bc1() {
	s.srv.Stop()
	s.l.Close()
}
func gologoo__genRandomSession_5d9035705d48cffd75ca99cea5c61bc1() string {
	var b [4]byte
	rand.Read(b[:])
	return fmt.Sprintf("%x", b)
}
func gologoo__genRandomTransaction_5d9035705d48cffd75ca99cea5c61bc1() string {
	var b [6]byte
	rand.Read(b[:])
	return fmt.Sprintf("tx-%x", b)
}
func gologoo__genRandomOperation_5d9035705d48cffd75ca99cea5c61bc1() string {
	var b [3]byte
	rand.Read(b[:])
	return fmt.Sprintf("op-%x", b)
}
func (s *server) gologoo__GetOperation_5d9035705d48cffd75ca99cea5c61bc1(ctx context.Context, req *lropb.GetOperationRequest) (*lropb.Operation, error) {
	s.mu.Lock()
	lro, ok := s.lros[req.Name]
	s.mu.Unlock()
	if !ok {
		return nil, status.Errorf(codes.NotFound, "unknown LRO %q", req.Name)
	}
	lro.noWait()
	return lro.State(), nil
}
func (s *server) gologoo__GetDatabase_5d9035705d48cffd75ca99cea5c61bc1(ctx context.Context, req *adminpb.GetDatabaseRequest) (*adminpb.Database, error) {
	s.logf("GetDatabase(%q)", req.Name)
	return &adminpb.Database{Name: req.Name, State: adminpb.Database_READY, CreateTime: timestampProto(s.start)}, nil
}
func (s *Server) gologoo__UpdateDDL_5d9035705d48cffd75ca99cea5c61bc1(ddl *spansql.DDL) error {
	ctx := context.Background()
	for _, stmt := range ddl.List {
		if st := s.s.runOneDDL(ctx, stmt); st.Code() != codes.OK {
			return st.Err()
		}
	}
	return nil
}
func (s *server) gologoo__UpdateDatabaseDdl_5d9035705d48cffd75ca99cea5c61bc1(ctx context.Context, req *adminpb.UpdateDatabaseDdlRequest) (*lropb.Operation, error) {
	var stmts []spansql.DDLStmt
	for _, s := range req.Statements {
		stmt, err := spansql.ParseDDLStmt(s)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "bad DDL statement %q: %v", s, err)
		}
		stmts = append(stmts, stmt)
	}
	id := "projects/fake-proj/instances/fake-instance/databases/fake-db/operations/" + genRandomOperation()
	lro := newLRO(&lropb.Operation{Name: id})
	s.mu.Lock()
	s.lros[id] = lro
	s.mu.Unlock()
	go lro.Run(s, stmts)
	return lro.State(), nil
}
func (l *lro) gologoo__Run_5d9035705d48cffd75ca99cea5c61bc1(s *server, stmts []spansql.DDLStmt) {
	ctx := context.Background()
	for _, stmt := range stmts {
		select {
		case <-time.After(100 * time.Millisecond):
		case <-l.waitc:
		}
		if st := s.runOneDDL(ctx, stmt); st.Code() != codes.OK {
			l.mu.Lock()
			l.state.Done = true
			l.state.Result = &lropb.Operation_Error{st.Proto()}
			l.mu.Unlock()
			return
		}
	}
	l.mu.Lock()
	l.state.Done = true
	l.state.Result = &lropb.Operation_Response{&anypb.Any{}}
	l.mu.Unlock()
}
func (s *server) gologoo__runOneDDL_5d9035705d48cffd75ca99cea5c61bc1(ctx context.Context, stmt spansql.DDLStmt) *status.Status {
	return s.db.ApplyDDL(stmt)
}
func (s *server) gologoo__GetDatabaseDdl_5d9035705d48cffd75ca99cea5c61bc1(ctx context.Context, req *adminpb.GetDatabaseDdlRequest) (*adminpb.GetDatabaseDdlResponse, error) {
	s.logf("GetDatabaseDdl(%q)", req.Database)
	var resp adminpb.GetDatabaseDdlResponse
	for _, stmt := range s.db.GetDDL() {
		resp.Statements = append(resp.Statements, stmt.SQL())
	}
	return &resp, nil
}
func (s *server) gologoo__CreateSession_5d9035705d48cffd75ca99cea5c61bc1(ctx context.Context, req *spannerpb.CreateSessionRequest) (*spannerpb.Session, error) {
	return s.newSession(), nil
}
func (s *server) gologoo__newSession_5d9035705d48cffd75ca99cea5c61bc1() *spannerpb.Session {
	id := genRandomSession()
	now := time.Now()
	sess := &session{name: id, creation: now, lastUse: now, transactions: make(map[string]*transaction)}
	sess.ctx, sess.cancel = context.WithCancel(context.Background())
	s.mu.Lock()
	s.sessions[id] = sess
	s.mu.Unlock()
	return sess.Proto()
}
func (s *server) gologoo__BatchCreateSessions_5d9035705d48cffd75ca99cea5c61bc1(ctx context.Context, req *spannerpb.BatchCreateSessionsRequest) (*spannerpb.BatchCreateSessionsResponse, error) {
	var sessions []*spannerpb.Session
	for i := int32(0); i < req.GetSessionCount(); i++ {
		sessions = append(sessions, s.newSession())
	}
	return &spannerpb.BatchCreateSessionsResponse{Session: sessions}, nil
}
func (s *server) gologoo__GetSession_5d9035705d48cffd75ca99cea5c61bc1(ctx context.Context, req *spannerpb.GetSessionRequest) (*spannerpb.Session, error) {
	s.mu.Lock()
	sess, ok := s.sessions[req.Name]
	s.mu.Unlock()
	if !ok {
		return nil, status.Errorf(codes.NotFound, "unknown session %q", req.Name)
	}
	return sess.Proto(), nil
}
func (s *server) gologoo__DeleteSession_5d9035705d48cffd75ca99cea5c61bc1(ctx context.Context, req *spannerpb.DeleteSessionRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	sess, ok := s.sessions[req.Name]
	delete(s.sessions, req.Name)
	s.mu.Unlock()
	if !ok {
		return nil, status.Errorf(codes.NotFound, "unknown session %q", req.Name)
	}
	sess.cancel()
	return &emptypb.Empty{}, nil
}
func (s *server) gologoo__popTx_5d9035705d48cffd75ca99cea5c61bc1(sessionID, tid string) (tx *transaction, err error) {
	s.mu.Lock()
	sess, ok := s.sessions[sessionID]
	s.mu.Unlock()
	if !ok {
		return nil, status.Errorf(codes.NotFound, "unknown session %q", sessionID)
	}
	sess.mu.Lock()
	sess.lastUse = time.Now()
	tx, ok = sess.transactions[tid]
	if ok {
		delete(sess.transactions, tid)
	}
	sess.mu.Unlock()
	if !ok {
		return nil, status.Errorf(codes.NotFound, "unknown transaction ID %q", tid)
	}
	return tx, nil
}
func (s *server) gologoo__readTx_5d9035705d48cffd75ca99cea5c61bc1(ctx context.Context, session string, tsel *spannerpb.TransactionSelector) (tx *transaction, cleanup func(), err error) {
	s.mu.Lock()
	sess, ok := s.sessions[session]
	s.mu.Unlock()
	if !ok {
		return nil, nil, status.Errorf(codes.NotFound, "unknown session %q", session)
	}
	sess.mu.Lock()
	sess.lastUse = time.Now()
	sess.mu.Unlock()
	singleUse := func() (*transaction, func(), error) {
		tx := s.db.NewReadOnlyTransaction()
		return tx, tx.Rollback, nil
	}
	if tsel.GetSelector() == nil {
		return singleUse()
	}
	switch sel := tsel.Selector.(type) {
	default:
		return nil, nil, fmt.Errorf("TransactionSelector type %T not supported", sel)
	case *spannerpb.TransactionSelector_SingleUse:
		switch mode := sel.SingleUse.Mode.(type) {
		case *spannerpb.TransactionOptions_ReadOnly_:
			return singleUse()
		case *spannerpb.TransactionOptions_ReadWrite_:
			return singleUse()
		default:
			return nil, nil, fmt.Errorf("single use transaction in mode %T not supported", mode)
		}
	case *spannerpb.TransactionSelector_Id:
		sess.mu.Lock()
		tx, ok := sess.transactions[string(sel.Id)]
		sess.mu.Unlock()
		if !ok {
			return nil, nil, fmt.Errorf("no transaction with id %q", sel.Id)
		}
		return tx, func() {
		}, nil
	}
}
func (s *server) gologoo__ExecuteSql_5d9035705d48cffd75ca99cea5c61bc1(ctx context.Context, req *spannerpb.ExecuteSqlRequest) (*spannerpb.ResultSet, error) {
	if req.Transaction.GetSelector() == nil || req.Transaction.GetSingleUse().GetReadOnly() != nil {
		ri, err := s.executeQuery(req)
		if err != nil {
			return nil, err
		}
		return s.resultSet(ri)
	}
	obj, ok := req.Transaction.Selector.(*spannerpb.TransactionSelector_Id)
	if !ok {
		return nil, fmt.Errorf("unsupported transaction type %T", req.Transaction.Selector)
	}
	tid := string(obj.Id)
	_ = tid
	stmt, err := spansql.ParseDMLStmt(req.Sql)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad DML: %v", err)
	}
	params, err := parseQueryParams(req.GetParams(), req.ParamTypes)
	if err != nil {
		return nil, err
	}
	s.logf("Executing: %s", stmt.SQL())
	if len(params) > 0 {
		s.logf("        ▹ %v", params)
	}
	n, err := s.db.Execute(stmt, params)
	if err != nil {
		return nil, err
	}
	return &spannerpb.ResultSet{Stats: &spannerpb.ResultSetStats{RowCount: &spannerpb.ResultSetStats_RowCountExact{int64(n)}}}, nil
}
func (s *server) gologoo__ExecuteStreamingSql_5d9035705d48cffd75ca99cea5c61bc1(req *spannerpb.ExecuteSqlRequest, stream spannerpb.Spanner_ExecuteStreamingSqlServer) error {
	tx, cleanup, err := s.readTx(stream.Context(), req.Session, req.Transaction)
	if err != nil {
		return err
	}
	defer cleanup()
	ri, err := s.executeQuery(req)
	if err != nil {
		return err
	}
	return s.readStream(stream.Context(), tx, stream.Send, ri)
}
func (s *server) gologoo__executeQuery_5d9035705d48cffd75ca99cea5c61bc1(req *spannerpb.ExecuteSqlRequest) (ri rowIter, err error) {
	q, err := spansql.ParseQuery(req.Sql)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "bad query: %v", err)
	}
	params, err := parseQueryParams(req.GetParams(), req.ParamTypes)
	if err != nil {
		return nil, err
	}
	s.logf("Querying: %s", q.SQL())
	if len(params) > 0 {
		s.logf("        ▹ %v", params)
	}
	return s.db.Query(q, params)
}
func (s *server) gologoo__StreamingRead_5d9035705d48cffd75ca99cea5c61bc1(req *spannerpb.ReadRequest, stream spannerpb.Spanner_StreamingReadServer) error {
	tx, cleanup, err := s.readTx(stream.Context(), req.Session, req.Transaction)
	if err != nil {
		return err
	}
	defer cleanup()
	if req.Index != "" {
		s.logf("Warning: index reads (%q) not supported", req.Index)
	}
	if len(req.ResumeToken) > 0 {
		return fmt.Errorf("read resumption not supported")
	}
	if len(req.PartitionToken) > 0 {
		return fmt.Errorf("partition restrictions not supported")
	}
	var ri rowIter
	if req.KeySet.All {
		s.logf("Reading all from %s (cols: %v)", req.Table, req.Columns)
		ri, err = s.db.ReadAll(spansql.ID(req.Table), idList(req.Columns), req.Limit)
	} else {
		s.logf("Reading rows from %d keys and %d ranges from %s (cols: %v)", len(req.KeySet.Keys), len(req.KeySet.Ranges), req.Table, req.Columns)
		ri, err = s.db.Read(spansql.ID(req.Table), idList(req.Columns), req.KeySet.Keys, makeKeyRangeList(req.KeySet.Ranges), req.Limit)
	}
	if err != nil {
		return err
	}
	return s.readStream(stream.Context(), tx, stream.Send, ri)
}
func (s *server) gologoo__resultSet_5d9035705d48cffd75ca99cea5c61bc1(ri rowIter) (*spannerpb.ResultSet, error) {
	rsm, err := s.buildResultSetMetadata(ri)
	if err != nil {
		return nil, err
	}
	rs := &spannerpb.ResultSet{Metadata: rsm}
	for {
		row, err := ri.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		values := make([]*structpb.Value, len(row))
		for i, x := range row {
			v, err := spannerValueFromValue(x)
			if err != nil {
				return nil, err
			}
			values[i] = v
		}
		rs.Rows = append(rs.Rows, &structpb.ListValue{Values: values})
	}
	return rs, nil
}
func (s *server) gologoo__readStream_5d9035705d48cffd75ca99cea5c61bc1(ctx context.Context, tx *transaction, send func(*spannerpb.PartialResultSet) error, ri rowIter) error {
	rsm, err := s.buildResultSetMetadata(ri)
	if err != nil {
		return err
	}
	for {
		row, err := ri.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		values := make([]*structpb.Value, len(row))
		for i, x := range row {
			v, err := spannerValueFromValue(x)
			if err != nil {
				return err
			}
			values[i] = v
		}
		prs := &spannerpb.PartialResultSet{Metadata: rsm, Values: values}
		if err := send(prs); err != nil {
			return err
		}
		rsm = nil
	}
	return nil
}
func (s *server) gologoo__buildResultSetMetadata_5d9035705d48cffd75ca99cea5c61bc1(ri rowIter) (*spannerpb.ResultSetMetadata, error) {
	rsm := &spannerpb.ResultSetMetadata{RowType: &spannerpb.StructType{}}
	for _, ci := range ri.Cols() {
		st, err := spannerTypeFromType(ci.Type)
		if err != nil {
			return nil, err
		}
		rsm.RowType.Fields = append(rsm.RowType.Fields, &spannerpb.StructType_Field{Name: string(ci.Name), Type: st})
	}
	return rsm, nil
}
func (s *server) gologoo__BeginTransaction_5d9035705d48cffd75ca99cea5c61bc1(ctx context.Context, req *spannerpb.BeginTransactionRequest) (*spannerpb.Transaction, error) {
	s.mu.Lock()
	sess, ok := s.sessions[req.Session]
	s.mu.Unlock()
	if !ok {
		return nil, status.Errorf(codes.NotFound, "unknown session %q", req.Session)
	}
	id := genRandomTransaction()
	tx := s.db.NewTransaction()
	sess.mu.Lock()
	sess.lastUse = time.Now()
	sess.transactions[id] = tx
	sess.mu.Unlock()
	tr := &spannerpb.Transaction{Id: []byte(id)}
	if req.GetOptions().GetReadOnly().GetReturnReadTimestamp() {
		tr.ReadTimestamp = timestampProto(s.db.LastCommitTimestamp())
	}
	return tr, nil
}
func (s *server) gologoo__Commit_5d9035705d48cffd75ca99cea5c61bc1(ctx context.Context, req *spannerpb.CommitRequest) (resp *spannerpb.CommitResponse, err error) {
	obj, ok := req.Transaction.(*spannerpb.CommitRequest_TransactionId)
	if !ok {
		return nil, fmt.Errorf("unsupported transaction type %T", req.Transaction)
	}
	tid := string(obj.TransactionId)
	tx, err := s.popTx(req.Session, tid)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	tx.Start()
	for _, m := range req.Mutations {
		switch op := m.Operation.(type) {
		default:
			return nil, fmt.Errorf("unsupported mutation operation type %T", op)
		case *spannerpb.Mutation_Insert:
			ins := op.Insert
			err := s.db.Insert(tx, spansql.ID(ins.Table), idList(ins.Columns), ins.Values)
			if err != nil {
				return nil, err
			}
		case *spannerpb.Mutation_Update:
			up := op.Update
			err := s.db.Update(tx, spansql.ID(up.Table), idList(up.Columns), up.Values)
			if err != nil {
				return nil, err
			}
		case *spannerpb.Mutation_InsertOrUpdate:
			iou := op.InsertOrUpdate
			err := s.db.InsertOrUpdate(tx, spansql.ID(iou.Table), idList(iou.Columns), iou.Values)
			if err != nil {
				return nil, err
			}
		case *spannerpb.Mutation_Delete_:
			del := op.Delete
			ks := del.KeySet
			err := s.db.Delete(tx, spansql.ID(del.Table), ks.Keys, makeKeyRangeList(ks.Ranges), ks.All)
			if err != nil {
				return nil, err
			}
		}
	}
	ts, err := tx.Commit()
	if err != nil {
		return nil, err
	}
	return &spannerpb.CommitResponse{CommitTimestamp: timestampProto(ts)}, nil
}
func (s *server) gologoo__Rollback_5d9035705d48cffd75ca99cea5c61bc1(ctx context.Context, req *spannerpb.RollbackRequest) (*emptypb.Empty, error) {
	s.logf("Rollback(%v)", req)
	tx, err := s.popTx(req.Session, string(req.TransactionId))
	if err != nil {
		return nil, err
	}
	tx.Rollback()
	return &emptypb.Empty{}, nil
}
func gologoo__parseQueryParams_5d9035705d48cffd75ca99cea5c61bc1(p *structpb.Struct, types map[string]*spannerpb.Type) (queryParams, error) {
	params := make(queryParams)
	for k, v := range p.GetFields() {
		p, err := parseQueryParam(v, types[k])
		if err != nil {
			return nil, err
		}
		params[k] = p
	}
	return params, nil
}
func gologoo__parseQueryParam_5d9035705d48cffd75ca99cea5c61bc1(v *structpb.Value, typ *spannerpb.Type) (queryParam, error) {
	rawv := v
	switch v := v.Kind.(type) {
	default:
		return queryParam{}, fmt.Errorf("unsupported well-known type value kind %T", v)
	case *structpb.Value_NullValue:
		return queryParam{Value: nil}, nil
	case *structpb.Value_BoolValue:
		return queryParam{Value: v.BoolValue, Type: boolType}, nil
	case *structpb.Value_NumberValue:
		return queryParam{Value: v.NumberValue, Type: float64Type}, nil
	case *structpb.Value_StringValue:
		t, err := typeFromSpannerType(typ)
		if err != nil {
			return queryParam{}, err
		}
		val, err := valForType(rawv, t)
		if err != nil {
			return queryParam{}, err
		}
		return queryParam{Value: val, Type: t}, nil
	case *structpb.Value_ListValue:
		var list []interface {
		}
		for _, elem := range v.ListValue.Values {
			p, err := parseQueryParam(elem, typ)
			if err != nil {
				return queryParam{}, err
			}
			list = append(list, p.Value)
		}
		t, err := typeFromSpannerType(typ)
		if err != nil {
			return queryParam{}, err
		}
		return queryParam{Value: list, Type: t}, nil
	}
}
func gologoo__typeFromSpannerType_5d9035705d48cffd75ca99cea5c61bc1(st *spannerpb.Type) (spansql.Type, error) {
	switch st.Code {
	default:
		return spansql.Type{}, fmt.Errorf("unhandled spanner type code %v", st.Code)
	case spannerpb.TypeCode_BOOL:
		return spansql.Type{Base: spansql.Bool}, nil
	case spannerpb.TypeCode_INT64:
		return spansql.Type{Base: spansql.Int64}, nil
	case spannerpb.TypeCode_FLOAT64:
		return spansql.Type{Base: spansql.Float64}, nil
	case spannerpb.TypeCode_TIMESTAMP:
		return spansql.Type{Base: spansql.Timestamp}, nil
	case spannerpb.TypeCode_DATE:
		return spansql.Type{Base: spansql.Date}, nil
	case spannerpb.TypeCode_STRING:
		return spansql.Type{Base: spansql.String}, nil
	case spannerpb.TypeCode_BYTES:
		return spansql.Type{Base: spansql.Bytes}, nil
	case spannerpb.TypeCode_ARRAY:
		typ, err := typeFromSpannerType(st.ArrayElementType)
		if err != nil {
			return spansql.Type{}, err
		}
		typ.Array = true
		return typ, nil
	}
}
func gologoo__spannerTypeFromType_5d9035705d48cffd75ca99cea5c61bc1(typ spansql.Type) (*spannerpb.Type, error) {
	var code spannerpb.TypeCode
	switch typ.Base {
	default:
		return nil, fmt.Errorf("unhandled base type %d", typ.Base)
	case spansql.Bool:
		code = spannerpb.TypeCode_BOOL
	case spansql.Int64:
		code = spannerpb.TypeCode_INT64
	case spansql.Float64:
		code = spannerpb.TypeCode_FLOAT64
	case spansql.String:
		code = spannerpb.TypeCode_STRING
	case spansql.Bytes:
		code = spannerpb.TypeCode_BYTES
	case spansql.Date:
		code = spannerpb.TypeCode_DATE
	case spansql.Timestamp:
		code = spannerpb.TypeCode_TIMESTAMP
	}
	st := &spannerpb.Type{Code: code}
	if typ.Array {
		st = &spannerpb.Type{Code: spannerpb.TypeCode_ARRAY, ArrayElementType: st}
	}
	return st, nil
}
func gologoo__spannerValueFromValue_5d9035705d48cffd75ca99cea5c61bc1(x interface {
}) (*structpb.Value, error) {
	switch x := x.(type) {
	default:
		return nil, fmt.Errorf("unhandled database value type %T", x)
	case bool:
		return &structpb.Value{Kind: &structpb.Value_BoolValue{x}}, nil
	case int64:
		s := strconv.FormatInt(x, 10)
		return &structpb.Value{Kind: &structpb.Value_StringValue{s}}, nil
	case float64:
		return &structpb.Value{Kind: &structpb.Value_NumberValue{x}}, nil
	case string:
		return &structpb.Value{Kind: &structpb.Value_StringValue{x}}, nil
	case []byte:
		return &structpb.Value{Kind: &structpb.Value_StringValue{base64.StdEncoding.EncodeToString(x)}}, nil
	case civil.Date:
		return &structpb.Value{Kind: &structpb.Value_StringValue{x.String()}}, nil
	case time.Time:
		s := x.Format("2006-01-02T15:04:05.999999999Z")
		return &structpb.Value{Kind: &structpb.Value_StringValue{s}}, nil
	case nil:
		return &structpb.Value{Kind: &structpb.Value_NullValue{}}, nil
	case []interface {
	}:
		var vs []*structpb.Value
		for _, elem := range x {
			v, err := spannerValueFromValue(elem)
			if err != nil {
				return nil, err
			}
			vs = append(vs, v)
		}
		return &structpb.Value{Kind: &structpb.Value_ListValue{&structpb.ListValue{Values: vs}}}, nil
	}
}
func gologoo__makeKeyRangeList_5d9035705d48cffd75ca99cea5c61bc1(ranges []*spannerpb.KeyRange) keyRangeList {
	var krl keyRangeList
	for _, r := range ranges {
		krl = append(krl, makeKeyRange(r))
	}
	return krl
}
func gologoo__makeKeyRange_5d9035705d48cffd75ca99cea5c61bc1(r *spannerpb.KeyRange) *keyRange {
	var kr keyRange
	switch s := r.StartKeyType.(type) {
	case *spannerpb.KeyRange_StartClosed:
		kr.start = s.StartClosed
		kr.startClosed = true
	case *spannerpb.KeyRange_StartOpen:
		kr.start = s.StartOpen
	}
	switch e := r.EndKeyType.(type) {
	case *spannerpb.KeyRange_EndClosed:
		kr.end = e.EndClosed
		kr.endClosed = true
	case *spannerpb.KeyRange_EndOpen:
		kr.end = e.EndOpen
	}
	return &kr
}
func gologoo__idList_5d9035705d48cffd75ca99cea5c61bc1(ss []string) (ids []spansql.ID) {
	for _, s := range ss {
		ids = append(ids, spansql.ID(s))
	}
	return
}
func (s *session) Proto() *spannerpb.Session {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Proto_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__Proto_5d9035705d48cffd75ca99cea5c61bc1()
	log.Printf("Output: %v\n", r0)
	return r0
}
func timestampProto(t time.Time) *timestamppb.Timestamp {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__timestampProto_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v\n", t)
	r0 := gologoo__timestampProto_5d9035705d48cffd75ca99cea5c61bc1(t)
	log.Printf("Output: %v\n", r0)
	return r0
}
func newLRO(initState *lropb.Operation) *lro {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__newLRO_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v\n", initState)
	r0 := gologoo__newLRO_5d9035705d48cffd75ca99cea5c61bc1(initState)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (l *lro) noWait() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__noWait_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : (none)\n")
	l.gologoo__noWait_5d9035705d48cffd75ca99cea5c61bc1()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (l *lro) State() *lropb.Operation {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__State_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : (none)\n")
	r0 := l.gologoo__State_5d9035705d48cffd75ca99cea5c61bc1()
	log.Printf("Output: %v\n", r0)
	return r0
}
func NewServer(laddr string) (*Server, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NewServer_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v\n", laddr)
	r0, r1 := gologoo__NewServer_5d9035705d48cffd75ca99cea5c61bc1(laddr)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *Server) SetLogger(l Logger) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SetLogger_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v\n", l)
	s.gologoo__SetLogger_5d9035705d48cffd75ca99cea5c61bc1(l)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *Server) Close() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Close_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : (none)\n")
	s.gologoo__Close_5d9035705d48cffd75ca99cea5c61bc1()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func genRandomSession() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__genRandomSession_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : (none)\n")
	r0 := gologoo__genRandomSession_5d9035705d48cffd75ca99cea5c61bc1()
	log.Printf("Output: %v\n", r0)
	return r0
}
func genRandomTransaction() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__genRandomTransaction_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : (none)\n")
	r0 := gologoo__genRandomTransaction_5d9035705d48cffd75ca99cea5c61bc1()
	log.Printf("Output: %v\n", r0)
	return r0
}
func genRandomOperation() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__genRandomOperation_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : (none)\n")
	r0 := gologoo__genRandomOperation_5d9035705d48cffd75ca99cea5c61bc1()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *server) GetOperation(ctx context.Context, req *lropb.GetOperationRequest) (*lropb.Operation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetOperation_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__GetOperation_5d9035705d48cffd75ca99cea5c61bc1(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *server) GetDatabase(ctx context.Context, req *adminpb.GetDatabaseRequest) (*adminpb.Database, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetDatabase_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__GetDatabase_5d9035705d48cffd75ca99cea5c61bc1(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *Server) UpdateDDL(ddl *spansql.DDL) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UpdateDDL_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v\n", ddl)
	r0 := s.gologoo__UpdateDDL_5d9035705d48cffd75ca99cea5c61bc1(ddl)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *server) UpdateDatabaseDdl(ctx context.Context, req *adminpb.UpdateDatabaseDdlRequest) (*lropb.Operation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__UpdateDatabaseDdl_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__UpdateDatabaseDdl_5d9035705d48cffd75ca99cea5c61bc1(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (l *lro) Run(s *server, stmts []spansql.DDLStmt) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Run_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", s, stmts)
	l.gologoo__Run_5d9035705d48cffd75ca99cea5c61bc1(s, stmts)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *server) runOneDDL(ctx context.Context, stmt spansql.DDLStmt) *status.Status {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__runOneDDL_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", ctx, stmt)
	r0 := s.gologoo__runOneDDL_5d9035705d48cffd75ca99cea5c61bc1(ctx, stmt)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *server) GetDatabaseDdl(ctx context.Context, req *adminpb.GetDatabaseDdlRequest) (*adminpb.GetDatabaseDdlResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetDatabaseDdl_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__GetDatabaseDdl_5d9035705d48cffd75ca99cea5c61bc1(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *server) CreateSession(ctx context.Context, req *spannerpb.CreateSessionRequest) (*spannerpb.Session, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateSession_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__CreateSession_5d9035705d48cffd75ca99cea5c61bc1(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *server) newSession() *spannerpb.Session {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__newSession_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__newSession_5d9035705d48cffd75ca99cea5c61bc1()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *server) BatchCreateSessions(ctx context.Context, req *spannerpb.BatchCreateSessionsRequest) (*spannerpb.BatchCreateSessionsResponse, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__BatchCreateSessions_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__BatchCreateSessions_5d9035705d48cffd75ca99cea5c61bc1(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *server) GetSession(ctx context.Context, req *spannerpb.GetSessionRequest) (*spannerpb.Session, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__GetSession_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__GetSession_5d9035705d48cffd75ca99cea5c61bc1(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *server) DeleteSession(ctx context.Context, req *spannerpb.DeleteSessionRequest) (*emptypb.Empty, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DeleteSession_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__DeleteSession_5d9035705d48cffd75ca99cea5c61bc1(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *server) popTx(sessionID, tid string) (tx *transaction, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__popTx_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", sessionID, tid)
	tx, err = s.gologoo__popTx_5d9035705d48cffd75ca99cea5c61bc1(sessionID, tid)
	log.Printf("Output: %v %v\n", tx, err)
	return
}
func (s *server) readTx(ctx context.Context, session string, tsel *spannerpb.TransactionSelector) (tx *transaction, cleanup func(), err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__readTx_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v %v\n", ctx, session, tsel)
	tx, cleanup, err = s.gologoo__readTx_5d9035705d48cffd75ca99cea5c61bc1(ctx, session, tsel)
	log.Printf("Output: %v %v %v\n", tx, cleanup, err)
	return
}
func (s *server) ExecuteSql(ctx context.Context, req *spannerpb.ExecuteSqlRequest) (*spannerpb.ResultSet, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ExecuteSql_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__ExecuteSql_5d9035705d48cffd75ca99cea5c61bc1(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *server) ExecuteStreamingSql(req *spannerpb.ExecuteSqlRequest, stream spannerpb.Spanner_ExecuteStreamingSqlServer) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ExecuteStreamingSql_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", req, stream)
	r0 := s.gologoo__ExecuteStreamingSql_5d9035705d48cffd75ca99cea5c61bc1(req, stream)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *server) executeQuery(req *spannerpb.ExecuteSqlRequest) (ri rowIter, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__executeQuery_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v\n", req)
	ri, err = s.gologoo__executeQuery_5d9035705d48cffd75ca99cea5c61bc1(req)
	log.Printf("Output: %v %v\n", ri, err)
	return
}
func (s *server) StreamingRead(req *spannerpb.ReadRequest, stream spannerpb.Spanner_StreamingReadServer) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__StreamingRead_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", req, stream)
	r0 := s.gologoo__StreamingRead_5d9035705d48cffd75ca99cea5c61bc1(req, stream)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *server) resultSet(ri rowIter) (*spannerpb.ResultSet, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__resultSet_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v\n", ri)
	r0, r1 := s.gologoo__resultSet_5d9035705d48cffd75ca99cea5c61bc1(ri)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *server) readStream(ctx context.Context, tx *transaction, send func(*spannerpb.PartialResultSet) error, ri rowIter) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__readStream_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v %v %v\n", ctx, tx, send, ri)
	r0 := s.gologoo__readStream_5d9035705d48cffd75ca99cea5c61bc1(ctx, tx, send, ri)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *server) buildResultSetMetadata(ri rowIter) (*spannerpb.ResultSetMetadata, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__buildResultSetMetadata_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v\n", ri)
	r0, r1 := s.gologoo__buildResultSetMetadata_5d9035705d48cffd75ca99cea5c61bc1(ri)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *server) BeginTransaction(ctx context.Context, req *spannerpb.BeginTransactionRequest) (*spannerpb.Transaction, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__BeginTransaction_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__BeginTransaction_5d9035705d48cffd75ca99cea5c61bc1(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (s *server) Commit(ctx context.Context, req *spannerpb.CommitRequest) (resp *spannerpb.CommitResponse, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Commit_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", ctx, req)
	resp, err = s.gologoo__Commit_5d9035705d48cffd75ca99cea5c61bc1(ctx, req)
	log.Printf("Output: %v %v\n", resp, err)
	return
}
func (s *server) Rollback(ctx context.Context, req *spannerpb.RollbackRequest) (*emptypb.Empty, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Rollback_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", ctx, req)
	r0, r1 := s.gologoo__Rollback_5d9035705d48cffd75ca99cea5c61bc1(ctx, req)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func parseQueryParams(p *structpb.Struct, types map[string]*spannerpb.Type) (queryParams, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseQueryParams_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", p, types)
	r0, r1 := gologoo__parseQueryParams_5d9035705d48cffd75ca99cea5c61bc1(p, types)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func parseQueryParam(v *structpb.Value, typ *spannerpb.Type) (queryParam, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseQueryParam_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v %v\n", v, typ)
	r0, r1 := gologoo__parseQueryParam_5d9035705d48cffd75ca99cea5c61bc1(v, typ)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func typeFromSpannerType(st *spannerpb.Type) (spansql.Type, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__typeFromSpannerType_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v\n", st)
	r0, r1 := gologoo__typeFromSpannerType_5d9035705d48cffd75ca99cea5c61bc1(st)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func spannerTypeFromType(typ spansql.Type) (*spannerpb.Type, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__spannerTypeFromType_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v\n", typ)
	r0, r1 := gologoo__spannerTypeFromType_5d9035705d48cffd75ca99cea5c61bc1(typ)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func spannerValueFromValue(x interface {
}) (*structpb.Value, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__spannerValueFromValue_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v\n", x)
	r0, r1 := gologoo__spannerValueFromValue_5d9035705d48cffd75ca99cea5c61bc1(x)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func makeKeyRangeList(ranges []*spannerpb.KeyRange) keyRangeList {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__makeKeyRangeList_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v\n", ranges)
	r0 := gologoo__makeKeyRangeList_5d9035705d48cffd75ca99cea5c61bc1(ranges)
	log.Printf("Output: %v\n", r0)
	return r0
}
func makeKeyRange(r *spannerpb.KeyRange) *keyRange {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__makeKeyRange_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v\n", r)
	r0 := gologoo__makeKeyRange_5d9035705d48cffd75ca99cea5c61bc1(r)
	log.Printf("Output: %v\n", r0)
	return r0
}
func idList(ss []string) (ids []spansql.ID) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__idList_5d9035705d48cffd75ca99cea5c61bc1")
	log.Printf("Input : %v\n", ss)
	ids = gologoo__idList_5d9035705d48cffd75ca99cea5c61bc1(ss)
	log.Printf("Output: %v\n", ids)
	return
}
