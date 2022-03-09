package spanner

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"cloud.google.com/go/internal/trace"
	vkit "cloud.google.com/go/spanner/apiv1"
	"cloud.google.com/go/spanner/internal"
	"github.com/googleapis/gax-go/v2"
	"go.opencensus.io/tag"
	"google.golang.org/api/option"
	gtransport "google.golang.org/api/transport/grpc"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var cidGen = newClientIDGenerator()

type clientIDGenerator struct {
	mu  sync.Mutex
	ids map[string]int
}

func gologoo__newClientIDGenerator_829b5d6d60227b496d673506b210fee2() *clientIDGenerator {
	return &clientIDGenerator{ids: make(map[string]int)}
}
func (cg *clientIDGenerator) gologoo__nextID_829b5d6d60227b496d673506b210fee2(database string) string {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	var id int
	if val, ok := cg.ids[database]; ok {
		id = val + 1
	} else {
		id = 1
	}
	cg.ids[database] = id
	return fmt.Sprintf("client-%d", id)
}

type sessionConsumer interface {
	sessionReady(s *session)
	sessionCreationFailed(err error, numSessions int32)
}
type sessionClient struct {
	mu            sync.Mutex
	closed        bool
	connPool      gtransport.ConnPool
	database      string
	id            string
	sessionLabels map[string]string
	md            metadata.MD
	batchTimeout  time.Duration
	logger        *log.Logger
	callOptions   *vkit.CallOptions
}

func gologoo__newSessionClient_829b5d6d60227b496d673506b210fee2(connPool gtransport.ConnPool, database string, sessionLabels map[string]string, md metadata.MD, logger *log.Logger, callOptions *vkit.CallOptions) *sessionClient {
	return &sessionClient{connPool: connPool, database: database, id: cidGen.nextID(database), sessionLabels: sessionLabels, md: md, batchTimeout: time.Minute, logger: logger, callOptions: callOptions}
}
func (sc *sessionClient) gologoo__close_829b5d6d60227b496d673506b210fee2() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	sc.closed = true
	return sc.connPool.Close()
}
func (sc *sessionClient) gologoo__createSession_829b5d6d60227b496d673506b210fee2(ctx context.Context) (*session, error) {
	sc.mu.Lock()
	if sc.closed {
		sc.mu.Unlock()
		return nil, spannerErrorf(codes.FailedPrecondition, "SessionClient is closed")
	}
	sc.mu.Unlock()
	client, err := sc.nextClient()
	if err != nil {
		return nil, err
	}
	ctx = contextWithOutgoingMetadata(ctx, sc.md)
	var md metadata.MD
	sid, err := client.CreateSession(ctx, &sppb.CreateSessionRequest{Database: sc.database, Session: &sppb.Session{Labels: sc.sessionLabels}}, gax.WithGRPCOptions(grpc.Header(&md)))
	if getGFELatencyMetricsFlag() && md != nil {
		_, instance, database, err := parseDatabaseName(sc.database)
		if err != nil {
			return nil, ToSpannerError(err)
		}
		ctxGFE, err := tag.New(ctx, tag.Upsert(tagKeyClientID, sc.id), tag.Upsert(tagKeyDatabase, database), tag.Upsert(tagKeyInstance, instance), tag.Upsert(tagKeyLibVersion, internal.Version))
		if err != nil {
			trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", ToSpannerError(err))
		}
		err = captureGFELatencyStats(ctxGFE, md, "createSession")
		if err != nil {
			trace.TracePrintf(ctx, nil, "Error in recording GFE Latency. Try disabling and rerunning. Error: %v", ToSpannerError(err))
		}
	}
	if err != nil {
		return nil, ToSpannerError(err)
	}
	return &session{valid: true, client: client, id: sid.Name, createTime: time.Now(), md: sc.md, logger: sc.logger}, nil
}
func (sc *sessionClient) gologoo__batchCreateSessions_829b5d6d60227b496d673506b210fee2(createSessionCount int32, distributeOverChannels bool, consumer sessionConsumer) error {
	var sessionCountPerChannel int32
	var remainder int32
	if distributeOverChannels {
		sessionCountPerChannel = createSessionCount / int32(sc.connPool.Num())
		remainder = createSessionCount % int32(sc.connPool.Num())
	} else {
		sessionCountPerChannel = createSessionCount
	}
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if sc.closed {
		return spannerErrorf(codes.FailedPrecondition, "SessionClient is closed")
	}
	var numBeingCreated int32
	for i := 0; i < sc.connPool.Num() && numBeingCreated < createSessionCount; i++ {
		client, err := sc.nextClient()
		if err != nil {
			return err
		}
		createCountForChannel := sessionCountPerChannel
		if i == 0 {
			createCountForChannel += remainder
		}
		if createCountForChannel > 0 {
			go sc.executeBatchCreateSessions(client, createCountForChannel, sc.sessionLabels, sc.md, consumer)
			numBeingCreated += createCountForChannel
		}
	}
	return nil
}
func (sc *sessionClient) gologoo__executeBatchCreateSessions_829b5d6d60227b496d673506b210fee2(client *vkit.Client, createCount int32, labels map[string]string, md metadata.MD, consumer sessionConsumer) {
	ctx, cancel := context.WithTimeout(context.Background(), sc.batchTimeout)
	defer cancel()
	ctx = contextWithOutgoingMetadata(ctx, sc.md)
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.BatchCreateSessions")
	defer func() {
		trace.EndSpan(ctx, nil)
	}()
	trace.TracePrintf(ctx, nil, "Creating a batch of %d sessions", createCount)
	remainingCreateCount := createCount
	for {
		sc.mu.Lock()
		closed := sc.closed
		sc.mu.Unlock()
		if closed {
			err := spannerErrorf(codes.Canceled, "Session client closed")
			trace.TracePrintf(ctx, nil, "Session client closed while creating a batch of %d sessions: %v", createCount, err)
			consumer.sessionCreationFailed(err, remainingCreateCount)
			break
		}
		if ctx.Err() != nil {
			trace.TracePrintf(ctx, nil, "Context error while creating a batch of %d sessions: %v", createCount, ctx.Err())
			consumer.sessionCreationFailed(ToSpannerError(ctx.Err()), remainingCreateCount)
			break
		}
		var mdForGFELatency metadata.MD
		response, err := client.BatchCreateSessions(ctx, &sppb.BatchCreateSessionsRequest{SessionCount: remainingCreateCount, Database: sc.database, SessionTemplate: &sppb.Session{Labels: labels}}, gax.WithGRPCOptions(grpc.Header(&mdForGFELatency)))
		if getGFELatencyMetricsFlag() && mdForGFELatency != nil {
			_, instance, database, err := parseDatabaseName(sc.database)
			if err != nil {
				trace.TracePrintf(ctx, nil, "Error getting instance and database name: %v", err)
			}
			ctxGFE, err := tag.New(ctx, tag.Upsert(tagKeyClientID, sc.id), tag.Upsert(tagKeyDatabase, database), tag.Upsert(tagKeyInstance, instance), tag.Upsert(tagKeyLibVersion, internal.Version))
			if err != nil {
				trace.TracePrintf(ctx, nil, "Error in adding tags in BatchCreateSessions for GFE Latency: %v", err)
			}
			err = captureGFELatencyStats(ctxGFE, mdForGFELatency, "executeBatchCreateSessions")
			if err != nil {
				trace.TracePrintf(ctx, nil, "Error in Capturing GFE Latency and Header Missing count. Try disabling and rerunning. Error: %v", err)
			}
		}
		if err != nil {
			trace.TracePrintf(ctx, nil, "Error creating a batch of %d sessions: %v", remainingCreateCount, err)
			consumer.sessionCreationFailed(ToSpannerError(err), remainingCreateCount)
			break
		}
		actuallyCreated := int32(len(response.Session))
		trace.TracePrintf(ctx, nil, "Received a batch of %d sessions", actuallyCreated)
		for _, s := range response.Session {
			consumer.sessionReady(&session{valid: true, client: client, id: s.Name, createTime: time.Now(), md: md, logger: sc.logger})
		}
		if actuallyCreated < remainingCreateCount {
			remainingCreateCount -= actuallyCreated
		} else {
			trace.TracePrintf(ctx, nil, "Finished creating %d sessions", createCount)
			break
		}
	}
}
func (sc *sessionClient) gologoo__sessionWithID_829b5d6d60227b496d673506b210fee2(id string) (*session, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	client, err := sc.nextClient()
	if err != nil {
		return nil, err
	}
	return &session{valid: true, client: client, id: id, createTime: time.Now(), md: sc.md, logger: sc.logger}, nil
}
func (sc *sessionClient) gologoo__nextClient_829b5d6d60227b496d673506b210fee2() (*vkit.Client, error) {
	client, err := vkit.NewClient(context.Background(), option.WithGRPCConn(sc.connPool.Conn()))
	if err != nil {
		return nil, err
	}
	client.SetGoogleClientInfo("gccl", internal.Version)
	if sc.callOptions != nil {
		client.CallOptions = mergeCallOptions(client.CallOptions, sc.callOptions)
	}
	return client, nil
}
func gologoo__mergeCallOptions_829b5d6d60227b496d673506b210fee2(a *vkit.CallOptions, b *vkit.CallOptions) *vkit.CallOptions {
	res := &vkit.CallOptions{}
	resVal := reflect.ValueOf(res).Elem()
	aVal := reflect.ValueOf(a).Elem()
	bVal := reflect.ValueOf(b).Elem()
	t := aVal.Type()
	for i := 0; i < aVal.NumField(); i++ {
		fieldName := t.Field(i).Name
		aFieldVal := aVal.Field(i).Interface().([]gax.CallOption)
		bFieldVal := bVal.Field(i).Interface().([]gax.CallOption)
		merged := append(aFieldVal, bFieldVal...)
		resVal.FieldByName(fieldName).Set(reflect.ValueOf(merged))
	}
	return res
}
func newClientIDGenerator() *clientIDGenerator {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__newClientIDGenerator_829b5d6d60227b496d673506b210fee2")
	log.Printf("Input : (none)\n")
	r0 := gologoo__newClientIDGenerator_829b5d6d60227b496d673506b210fee2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (cg *clientIDGenerator) nextID(database string) string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__nextID_829b5d6d60227b496d673506b210fee2")
	log.Printf("Input : %v\n", database)
	r0 := cg.gologoo__nextID_829b5d6d60227b496d673506b210fee2(database)
	log.Printf("Output: %v\n", r0)
	return r0
}
func newSessionClient(connPool gtransport.ConnPool, database string, sessionLabels map[string]string, md metadata.MD, logger *log.Logger, callOptions *vkit.CallOptions) *sessionClient {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__newSessionClient_829b5d6d60227b496d673506b210fee2")
	log.Printf("Input : %v %v %v %v %v %v\n", connPool, database, sessionLabels, md, logger, callOptions)
	r0 := gologoo__newSessionClient_829b5d6d60227b496d673506b210fee2(connPool, database, sessionLabels, md, logger, callOptions)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (sc *sessionClient) close() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__close_829b5d6d60227b496d673506b210fee2")
	log.Printf("Input : (none)\n")
	r0 := sc.gologoo__close_829b5d6d60227b496d673506b210fee2()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (sc *sessionClient) createSession(ctx context.Context) (*session, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__createSession_829b5d6d60227b496d673506b210fee2")
	log.Printf("Input : %v\n", ctx)
	r0, r1 := sc.gologoo__createSession_829b5d6d60227b496d673506b210fee2(ctx)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (sc *sessionClient) batchCreateSessions(createSessionCount int32, distributeOverChannels bool, consumer sessionConsumer) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__batchCreateSessions_829b5d6d60227b496d673506b210fee2")
	log.Printf("Input : %v %v %v\n", createSessionCount, distributeOverChannels, consumer)
	r0 := sc.gologoo__batchCreateSessions_829b5d6d60227b496d673506b210fee2(createSessionCount, distributeOverChannels, consumer)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (sc *sessionClient) executeBatchCreateSessions(client *vkit.Client, createCount int32, labels map[string]string, md metadata.MD, consumer sessionConsumer) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__executeBatchCreateSessions_829b5d6d60227b496d673506b210fee2")
	log.Printf("Input : %v %v %v %v %v\n", client, createCount, labels, md, consumer)
	sc.gologoo__executeBatchCreateSessions_829b5d6d60227b496d673506b210fee2(client, createCount, labels, md, consumer)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (sc *sessionClient) sessionWithID(id string) (*session, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__sessionWithID_829b5d6d60227b496d673506b210fee2")
	log.Printf("Input : %v\n", id)
	r0, r1 := sc.gologoo__sessionWithID_829b5d6d60227b496d673506b210fee2(id)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (sc *sessionClient) nextClient() (*vkit.Client, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__nextClient_829b5d6d60227b496d673506b210fee2")
	log.Printf("Input : (none)\n")
	r0, r1 := sc.gologoo__nextClient_829b5d6d60227b496d673506b210fee2()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func mergeCallOptions(a *vkit.CallOptions, b *vkit.CallOptions) *vkit.CallOptions {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__mergeCallOptions_829b5d6d60227b496d673506b210fee2")
	log.Printf("Input : %v %v\n", a, b)
	r0 := gologoo__mergeCallOptions_829b5d6d60227b496d673506b210fee2(a, b)
	log.Printf("Output: %v\n", r0)
	return r0
}
