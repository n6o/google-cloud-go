package spanner

import (
	"container/heap"
	"container/list"
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime/debug"
	"sync"
	"time"

	"cloud.google.com/go/internal/trace"
	vkit "cloud.google.com/go/spanner/apiv1"
	"cloud.google.com/go/spanner/internal"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	octrace "go.opencensus.io/trace"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

const healthCheckIntervalMins = 50

type sessionHandle struct {
	mu                   sync.Mutex
	session              *session
	checkoutTime         time.Time
	trackedSessionHandle *list.Element
	stack                []byte
}

func (sh *sessionHandle) gologoo__recycle_b42872e5776c89a8da4fd9993002e327() {
	sh.mu.Lock()
	if sh.session == nil {
		sh.mu.Unlock()
		return
	}
	p := sh.session.pool
	tracked := sh.trackedSessionHandle
	sh.session.recycle()
	sh.session = nil
	sh.trackedSessionHandle = nil
	sh.checkoutTime = time.Time{}
	sh.stack = nil
	sh.mu.Unlock()
	if tracked != nil {
		p.mu.Lock()
		p.trackedSessionHandles.Remove(tracked)
		p.mu.Unlock()
	}
}
func (sh *sessionHandle) gologoo__getID_b42872e5776c89a8da4fd9993002e327() string {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if sh.session == nil {
		return ""
	}
	return sh.session.getID()
}
func (sh *sessionHandle) gologoo__getClient_b42872e5776c89a8da4fd9993002e327() *vkit.Client {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if sh.session == nil {
		return nil
	}
	return sh.session.client
}
func (sh *sessionHandle) gologoo__getMetadata_b42872e5776c89a8da4fd9993002e327() metadata.MD {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if sh.session == nil {
		return nil
	}
	return sh.session.md
}
func (sh *sessionHandle) gologoo__getTransactionID_b42872e5776c89a8da4fd9993002e327() transactionID {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if sh.session == nil {
		return nil
	}
	return sh.session.tx
}
func (sh *sessionHandle) gologoo__destroy_b42872e5776c89a8da4fd9993002e327() {
	sh.mu.Lock()
	s := sh.session
	if s == nil {
		sh.mu.Unlock()
		return
	}
	tracked := sh.trackedSessionHandle
	sh.session = nil
	sh.trackedSessionHandle = nil
	sh.checkoutTime = time.Time{}
	sh.stack = nil
	sh.mu.Unlock()
	if tracked != nil {
		p := s.pool
		p.mu.Lock()
		p.trackedSessionHandles.Remove(tracked)
		p.mu.Unlock()
	}
	s.destroy(false)
}

type session struct {
	client         *vkit.Client
	id             string
	pool           *sessionPool
	createTime     time.Time
	logger         *log.Logger
	mu             sync.Mutex
	valid          bool
	hcIndex        int
	idleList       *list.Element
	nextCheck      time.Time
	checkingHealth bool
	md             metadata.MD
	tx             transactionID
	firstHCDone    bool
}

func (s *session) gologoo__isValid_b42872e5776c89a8da4fd9993002e327() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.valid
}
func (s *session) gologoo__isWritePrepared_b42872e5776c89a8da4fd9993002e327() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tx != nil
}
func (s *session) gologoo__String_b42872e5776c89a8da4fd9993002e327() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("<id=%v, hcIdx=%v, idleList=%p, valid=%v, create=%v, nextcheck=%v>", s.id, s.hcIndex, s.idleList, s.valid, s.createTime, s.nextCheck)
}
func (s *session) gologoo__ping_b42872e5776c89a8da4fd9993002e327() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, span := octrace.StartSpan(ctx, "cloud.google.com/go/spanner.ping", octrace.WithSampler(octrace.NeverSample()))
	defer span.End()
	_, err := s.client.ExecuteSql(contextWithOutgoingMetadata(ctx, s.md), &sppb.ExecuteSqlRequest{Session: s.getID(), Sql: "SELECT 1"})
	return err
}
func (s *session) gologoo__setHcIndex_b42872e5776c89a8da4fd9993002e327(i int) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	oi := s.hcIndex
	s.hcIndex = i
	return oi
}
func (s *session) gologoo__setIdleList_b42872e5776c89a8da4fd9993002e327(le *list.Element) *list.Element {
	s.mu.Lock()
	defer s.mu.Unlock()
	old := s.idleList
	s.idleList = le
	return old
}
func (s *session) gologoo__invalidate_b42872e5776c89a8da4fd9993002e327() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	ov := s.valid
	s.valid = false
	return ov
}
func (s *session) gologoo__setNextCheck_b42872e5776c89a8da4fd9993002e327(t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextCheck = t
}
func (s *session) gologoo__setTransactionID_b42872e5776c89a8da4fd9993002e327(tx transactionID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tx = tx
}
func (s *session) gologoo__getID_b42872e5776c89a8da4fd9993002e327() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.id
}
func (s *session) gologoo__getHcIndex_b42872e5776c89a8da4fd9993002e327() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.hcIndex
}
func (s *session) gologoo__getIdleList_b42872e5776c89a8da4fd9993002e327() *list.Element {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.idleList
}
func (s *session) gologoo__getNextCheck_b42872e5776c89a8da4fd9993002e327() time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.nextCheck
}
func (s *session) gologoo__recycle_b42872e5776c89a8da4fd9993002e327() {
	s.setTransactionID(nil)
	if !s.pool.recycle(s) {
		s.destroy(false)
	}
	s.pool.decNumInUse(context.Background())
}
func (s *session) gologoo__destroy_b42872e5776c89a8da4fd9993002e327(isExpire bool) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return s.destroyWithContext(ctx, isExpire)
}
func (s *session) gologoo__destroyWithContext_b42872e5776c89a8da4fd9993002e327(ctx context.Context, isExpire bool) bool {
	if !s.pool.remove(s, isExpire) {
		return false
	}
	s.pool.hc.unregister(s)
	s.delete(ctx)
	return true
}
func (s *session) gologoo__delete_b42872e5776c89a8da4fd9993002e327(ctx context.Context) {
	err := s.client.DeleteSession(contextWithOutgoingMetadata(ctx, s.md), &sppb.DeleteSessionRequest{Name: s.getID()})
	if err != nil && ErrCode(err) != codes.DeadlineExceeded {
		logf(s.logger, "Failed to delete session %v. Error: %v", s.getID(), err)
	}
}
func (s *session) gologoo__prepareForWrite_b42872e5776c89a8da4fd9993002e327(ctx context.Context) error {
	if s.isWritePrepared() {
		return nil
	}
	tx, err := beginTransaction(contextWithOutgoingMetadata(ctx, s.md), s.getID(), s.client)
	if isSessionNotFoundError(err) {
		s.pool.remove(s, false)
		s.pool.hc.unregister(s)
		return err
	}
	s.pool.mu.Lock()
	s.pool.disableBackgroundPrepareSessions = err != nil
	s.pool.mu.Unlock()
	if err != nil {
		return err
	}
	s.setTransactionID(tx)
	return nil
}

type SessionPoolConfig struct {
	MaxOpened                 uint64
	MinOpened                 uint64
	MaxIdle                   uint64
	MaxBurst                  uint64
	incStep                   uint64
	WriteSessions             float64
	HealthCheckWorkers        int
	HealthCheckInterval       time.Duration
	TrackSessionHandles       bool
	healthCheckSampleInterval time.Duration
	sessionLabels             map[string]string
}

var DefaultSessionPoolConfig = SessionPoolConfig{MinOpened: 100, MaxOpened: numChannels * 100, MaxBurst: 10, incStep: 25, WriteSessions: 0.2, HealthCheckWorkers: 10, HealthCheckInterval: healthCheckIntervalMins * time.Minute}

func gologoo__errMinOpenedGTMaxOpened_b42872e5776c89a8da4fd9993002e327(maxOpened, minOpened uint64) error {
	return spannerErrorf(codes.InvalidArgument, "require SessionPoolConfig.MaxOpened >= SessionPoolConfig.MinOpened, got %d and %d", maxOpened, minOpened)
}
func gologoo__errWriteFractionOutOfRange_b42872e5776c89a8da4fd9993002e327(writeFraction float64) error {
	return spannerErrorf(codes.InvalidArgument, "require SessionPoolConfig.WriteSessions >= 0.0 && SessionPoolConfig.WriteSessions <= 1.0, got %.2f", writeFraction)
}
func gologoo__errHealthCheckWorkersNegative_b42872e5776c89a8da4fd9993002e327(workers int) error {
	return spannerErrorf(codes.InvalidArgument, "require SessionPoolConfig.HealthCheckWorkers >= 0, got %d", workers)
}
func gologoo__errHealthCheckIntervalNegative_b42872e5776c89a8da4fd9993002e327(interval time.Duration) error {
	return spannerErrorf(codes.InvalidArgument, "require SessionPoolConfig.HealthCheckInterval >= 0, got %v", interval)
}
func (spc *SessionPoolConfig) gologoo__validate_b42872e5776c89a8da4fd9993002e327() error {
	if spc.MinOpened > spc.MaxOpened && spc.MaxOpened > 0 {
		return errMinOpenedGTMaxOpened(spc.MaxOpened, spc.MinOpened)
	}
	if spc.WriteSessions < 0.0 || spc.WriteSessions > 1.0 {
		return errWriteFractionOutOfRange(spc.WriteSessions)
	}
	if spc.HealthCheckWorkers < 0 {
		return errHealthCheckWorkersNegative(spc.HealthCheckWorkers)
	}
	if spc.HealthCheckInterval < 0 {
		return errHealthCheckIntervalNegative(spc.HealthCheckInterval)
	}
	return nil
}

type sessionPool struct {
	mu                    sync.Mutex
	valid                 bool
	sc                    *sessionClient
	trackedSessionHandles list.List
	idleList              list.List
	idleWriteList         list.List
	mayGetSession         chan struct {
	}
	sessionCreationError             error
	numOpened                        uint64
	createReqs                       uint64
	prepareReqs                      uint64
	numReadWaiters                   uint64
	numWriteWaiters                  uint64
	disableBackgroundPrepareSessions bool
	SessionPoolConfig
	hc            *healthChecker
	rand          *rand.Rand
	numInUse      uint64
	maxNumInUse   uint64
	lastResetTime time.Time
	numReads      uint64
	numWrites     uint64
	mw            *maintenanceWindow
	tagMap        *tag.Map
}

func gologoo__newSessionPool_b42872e5776c89a8da4fd9993002e327(sc *sessionClient, config SessionPoolConfig) (*sessionPool, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}
	pool := &sessionPool{sc: sc, valid: true, mayGetSession: make(chan struct {
	}), SessionPoolConfig: config, mw: newMaintenanceWindow(config.MaxOpened), rand: rand.New(rand.NewSource(time.Now().UnixNano()))}
	if config.HealthCheckWorkers == 0 {
		config.HealthCheckWorkers = 10
	}
	if config.HealthCheckInterval == 0 {
		config.HealthCheckInterval = healthCheckIntervalMins * time.Minute
	}
	if config.healthCheckSampleInterval == 0 {
		config.healthCheckSampleInterval = time.Minute
	}
	_, instance, database, err := parseDatabaseName(sc.database)
	if err != nil {
		return nil, err
	}
	ctx, err := tag.New(context.Background(), tag.Upsert(tagKeyClientID, sc.id), tag.Upsert(tagKeyDatabase, database), tag.Upsert(tagKeyInstance, instance), tag.Upsert(tagKeyLibVersion, internal.Version))
	if err != nil {
		logf(pool.sc.logger, "Failed to create tag map, error: %v", err)
	}
	pool.tagMap = tag.FromContext(ctx)
	pool.hc = newHealthChecker(config.HealthCheckInterval, config.HealthCheckWorkers, config.healthCheckSampleInterval, pool)
	if config.MinOpened > 0 {
		numSessions := minUint64(config.MinOpened, math.MaxInt32)
		if err := pool.initPool(numSessions); err != nil {
			return nil, err
		}
	}
	pool.recordStat(context.Background(), MaxAllowedSessionsCount, int64(config.MaxOpened))
	close(pool.hc.ready)
	return pool, nil
}
func (p *sessionPool) gologoo__recordStat_b42872e5776c89a8da4fd9993002e327(ctx context.Context, m *stats.Int64Measure, n int64, tags ...tag.Tag) {
	ctx = tag.NewContext(ctx, p.tagMap)
	mutators := make([]tag.Mutator, len(tags))
	for i, t := range tags {
		mutators[i] = tag.Upsert(t.Key, t.Value)
	}
	ctx, err := tag.New(ctx, mutators...)
	if err != nil {
		logf(p.sc.logger, "Failed to tag metrics, error: %v", err)
	}
	recordStat(ctx, m, n)
}
func (p *sessionPool) gologoo__initPool_b42872e5776c89a8da4fd9993002e327(numSessions uint64) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.growPoolLocked(numSessions, true)
}
func (p *sessionPool) gologoo__growPoolLocked_b42872e5776c89a8da4fd9993002e327(numSessions uint64, distributeOverChannels bool) error {
	numSessions = minUint64(numSessions, math.MaxInt32)
	p.numOpened += uint64(numSessions)
	p.recordStat(context.Background(), OpenSessionCount, int64(p.numOpened))
	p.createReqs += uint64(numSessions)
	return p.sc.batchCreateSessions(int32(numSessions), distributeOverChannels, p)
}
func (p *sessionPool) gologoo__sessionReady_b42872e5776c89a8da4fd9993002e327(s *session) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.sessionCreationError = nil
	s.pool = p
	p.hc.register(s)
	p.createReqs--
	if p.idleList.Len() > 0 {
		pos := rand.Intn(p.idleList.Len())
		before := p.idleList.Front()
		for i := 0; i < pos; i++ {
			before = before.Next()
		}
		s.setIdleList(p.idleList.InsertBefore(s, before))
	} else {
		s.setIdleList(p.idleList.PushBack(s))
	}
	p.incNumReadsLocked(context.Background())
	close(p.mayGetSession)
	p.mayGetSession = make(chan struct {
	})
}
func (p *sessionPool) gologoo__sessionCreationFailed_b42872e5776c89a8da4fd9993002e327(err error, numSessions int32) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.createReqs -= uint64(numSessions)
	p.numOpened -= uint64(numSessions)
	p.recordStat(context.Background(), OpenSessionCount, int64(p.numOpened))
	p.sessionCreationError = err
	close(p.mayGetSession)
	p.mayGetSession = make(chan struct {
	})
}
func (p *sessionPool) gologoo__isValid_b42872e5776c89a8da4fd9993002e327() bool {
	if p == nil {
		return false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.valid
}
func (p *sessionPool) gologoo__close_b42872e5776c89a8da4fd9993002e327(ctx context.Context) {
	if p == nil {
		return
	}
	p.mu.Lock()
	if !p.valid {
		p.mu.Unlock()
		return
	}
	p.valid = false
	p.mu.Unlock()
	p.hc.close()
	p.hc.mu.Lock()
	allSessions := make([]*session, len(p.hc.queue.sessions))
	copy(allSessions, p.hc.queue.sessions)
	p.hc.mu.Unlock()
	wg := sync.WaitGroup{}
	for _, s := range allSessions {
		wg.Add(1)
		go deleteSession(ctx, s, &wg)
	}
	wg.Wait()
}
func gologoo__deleteSession_b42872e5776c89a8da4fd9993002e327(ctx context.Context, s *session, wg *sync.WaitGroup) {
	defer wg.Done()
	s.destroyWithContext(ctx, false)
}

var errInvalidSessionPool = spannerErrorf(codes.InvalidArgument, "invalid session pool")
var errGetSessionTimeout = spannerErrorf(codes.Canceled, "timeout / context canceled during getting session")

func (p *sessionPool) gologoo__newSessionHandle_b42872e5776c89a8da4fd9993002e327(s *session) (sh *sessionHandle) {
	sh = &sessionHandle{session: s, checkoutTime: time.Now()}
	if p.TrackSessionHandles {
		p.mu.Lock()
		sh.trackedSessionHandle = p.trackedSessionHandles.PushBack(sh)
		p.mu.Unlock()
		sh.stack = debug.Stack()
	}
	return sh
}
func (p *sessionPool) gologoo__errGetSessionTimeout_b42872e5776c89a8da4fd9993002e327(ctx context.Context) error {
	var code codes.Code
	if ctx.Err() == context.DeadlineExceeded {
		code = codes.DeadlineExceeded
	} else {
		code = codes.Canceled
	}
	if p.TrackSessionHandles {
		return p.errGetSessionTimeoutWithTrackedSessionHandles(code)
	}
	return p.errGetBasicSessionTimeout(code)
}
func (p *sessionPool) gologoo__errGetBasicSessionTimeout_b42872e5776c89a8da4fd9993002e327(code codes.Code) error {
	return spannerErrorf(code, "timeout / context canceled during getting session.\n"+"Enable SessionPoolConfig.TrackSessionHandles if you suspect a session leak to get more information about the checked out sessions.")
}
func (p *sessionPool) gologoo__errGetSessionTimeoutWithTrackedSessionHandles_b42872e5776c89a8da4fd9993002e327(code codes.Code) error {
	err := spannerErrorf(code, "timeout / context canceled during getting session.")
	err.(*Error).additionalInformation = p.getTrackedSessionHandleStacksLocked()
	return err
}
func (p *sessionPool) gologoo__getTrackedSessionHandleStacksLocked_b42872e5776c89a8da4fd9993002e327() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	stackTraces := ""
	i := 1
	element := p.trackedSessionHandles.Front()
	for element != nil {
		sh := element.Value.(*sessionHandle)
		sh.mu.Lock()
		if sh.stack != nil {
			stackTraces = fmt.Sprintf("%s\n\nSession %d checked out of pool at %s by goroutine:\n%s", stackTraces, i, sh.checkoutTime.Format(time.RFC3339), sh.stack)
		}
		sh.mu.Unlock()
		element = element.Next()
		i++
	}
	return stackTraces
}
func (p *sessionPool) gologoo__shouldPrepareWriteLocked_b42872e5776c89a8da4fd9993002e327() bool {
	return !p.disableBackgroundPrepareSessions && float64(p.numOpened)*p.WriteSessions > float64(p.idleWriteList.Len()+int(p.prepareReqs))
}
func (p *sessionPool) gologoo__createSession_b42872e5776c89a8da4fd9993002e327(ctx context.Context) (*session, error) {
	trace.TracePrintf(ctx, nil, "Creating a new session")
	doneCreate := func(done bool) {
		p.mu.Lock()
		if !done {
			p.numOpened--
			p.recordStat(ctx, OpenSessionCount, int64(p.numOpened))
		}
		p.createReqs--
		close(p.mayGetSession)
		p.mayGetSession = make(chan struct {
		})
		p.mu.Unlock()
	}
	s, err := p.sc.createSession(ctx)
	if err != nil {
		doneCreate(false)
		return nil, err
	}
	s.pool = p
	p.hc.register(s)
	doneCreate(true)
	return s, nil
}
func (p *sessionPool) gologoo__isHealthy_b42872e5776c89a8da4fd9993002e327(s *session) bool {
	if s.getNextCheck().Add(2 * p.hc.getInterval()).Before(time.Now()) {
		if err := s.ping(); isSessionNotFoundError(err) {
			s.destroy(false)
			return false
		}
		p.hc.scheduledHC(s)
	}
	return true
}
func (p *sessionPool) gologoo__take_b42872e5776c89a8da4fd9993002e327(ctx context.Context) (*sessionHandle, error) {
	trace.TracePrintf(ctx, nil, "Acquiring a read-only session")
	for {
		var s *session
		p.mu.Lock()
		if !p.valid {
			p.mu.Unlock()
			return nil, errInvalidSessionPool
		}
		if p.idleList.Len() > 0 {
			s = p.idleList.Remove(p.idleList.Front()).(*session)
			trace.TracePrintf(ctx, map[string]interface {
			}{"sessionID": s.getID()}, "Acquired read-only session")
			p.decNumReadsLocked(ctx)
		} else if p.idleWriteList.Len() > 0 {
			s = p.idleWriteList.Remove(p.idleWriteList.Front()).(*session)
			trace.TracePrintf(ctx, map[string]interface {
			}{"sessionID": s.getID()}, "Acquired read-write session")
			p.decNumWritesLocked(ctx)
		}
		if s != nil {
			s.setIdleList(nil)
			numCheckedOut := p.currSessionsCheckedOutLocked()
			p.mu.Unlock()
			p.mw.updateMaxSessionsCheckedOutDuringWindow(numCheckedOut)
			if !p.isHealthy(s) {
				continue
			}
			p.incNumInUse(ctx)
			return p.newSessionHandle(s), nil
		}
		if p.numReadWaiters+p.numWriteWaiters >= p.createReqs {
			numSessions := minUint64(p.MaxOpened-p.numOpened, p.incStep)
			if err := p.growPoolLocked(numSessions, false); err != nil {
				p.mu.Unlock()
				return nil, err
			}
		}
		p.numReadWaiters++
		mayGetSession := p.mayGetSession
		p.mu.Unlock()
		trace.TracePrintf(ctx, nil, "Waiting for read-only session to become available")
		select {
		case <-ctx.Done():
			trace.TracePrintf(ctx, nil, "Context done waiting for session")
			p.recordStat(ctx, GetSessionTimeoutsCount, 1)
			p.mu.Lock()
			p.numReadWaiters--
			p.mu.Unlock()
			return nil, p.errGetSessionTimeout(ctx)
		case <-mayGetSession:
			p.mu.Lock()
			p.numReadWaiters--
			if p.sessionCreationError != nil {
				trace.TracePrintf(ctx, nil, "Error creating session: %v", p.sessionCreationError)
				err := p.sessionCreationError
				p.mu.Unlock()
				return nil, err
			}
			p.mu.Unlock()
		}
	}
}
func (p *sessionPool) gologoo__takeWriteSession_b42872e5776c89a8da4fd9993002e327(ctx context.Context) (*sessionHandle, error) {
	trace.TracePrintf(ctx, nil, "Acquiring a read-write session")
	for {
		var (
			s   *session
			err error
		)
		p.mu.Lock()
		if !p.valid {
			p.mu.Unlock()
			return nil, errInvalidSessionPool
		}
		if p.idleWriteList.Len() > 0 {
			s = p.idleWriteList.Remove(p.idleWriteList.Front()).(*session)
			trace.TracePrintf(ctx, map[string]interface {
			}{"sessionID": s.getID()}, "Acquired read-write session")
			p.decNumWritesLocked(ctx)
		} else if p.idleList.Len() > 0 {
			s = p.idleList.Remove(p.idleList.Front()).(*session)
			trace.TracePrintf(ctx, map[string]interface {
			}{"sessionID": s.getID()}, "Acquired read-only session")
			p.decNumReadsLocked(ctx)
		}
		if s != nil {
			s.setIdleList(nil)
			numCheckedOut := p.currSessionsCheckedOutLocked()
			p.mu.Unlock()
			p.mw.updateMaxSessionsCheckedOutDuringWindow(numCheckedOut)
			if !p.isHealthy(s) {
				continue
			}
		} else {
			if p.numReadWaiters+p.numWriteWaiters >= p.createReqs {
				numSessions := minUint64(p.MaxOpened-p.numOpened, p.incStep)
				if err := p.growPoolLocked(numSessions, false); err != nil {
					p.mu.Unlock()
					return nil, err
				}
			}
			p.numWriteWaiters++
			mayGetSession := p.mayGetSession
			p.mu.Unlock()
			trace.TracePrintf(ctx, nil, "Waiting for read-write session to become available")
			select {
			case <-ctx.Done():
				trace.TracePrintf(ctx, nil, "Context done waiting for session")
				p.recordStat(ctx, GetSessionTimeoutsCount, 1)
				p.mu.Lock()
				p.numWriteWaiters--
				p.mu.Unlock()
				return nil, p.errGetSessionTimeout(ctx)
			case <-mayGetSession:
				p.mu.Lock()
				p.numWriteWaiters--
				if p.sessionCreationError != nil {
					err := p.sessionCreationError
					p.mu.Unlock()
					return nil, err
				}
				p.mu.Unlock()
			}
			continue
		}
		if !s.isWritePrepared() {
			p.incNumBeingPrepared(ctx)
			defer p.decNumBeingPrepared(ctx)
			if err = s.prepareForWrite(ctx); err != nil {
				if isSessionNotFoundError(err) {
					s.destroy(false)
					trace.TracePrintf(ctx, map[string]interface {
					}{"sessionID": s.getID()}, "Session not found for write")
					return nil, ToSpannerError(err)
				}
				s.recycle()
				trace.TracePrintf(ctx, map[string]interface {
				}{"sessionID": s.getID()}, "Error preparing session for write")
				return nil, ToSpannerError(err)
			}
		}
		p.incNumInUse(ctx)
		return p.newSessionHandle(s), nil
	}
}
func (p *sessionPool) gologoo__recycle_b42872e5776c89a8da4fd9993002e327(s *session) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !s.isValid() || !p.valid {
		return false
	}
	ctx := context.Background()
	if s.isWritePrepared() {
		s.setIdleList(p.idleWriteList.PushFront(s))
		p.incNumWritesLocked(ctx)
	} else {
		s.setIdleList(p.idleList.PushFront(s))
		p.incNumReadsLocked(ctx)
	}
	close(p.mayGetSession)
	p.mayGetSession = make(chan struct {
	})
	return true
}
func (p *sessionPool) gologoo__remove_b42872e5776c89a8da4fd9993002e327(s *session, isExpire bool) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if isExpire && (p.numOpened <= p.MinOpened || s.getIdleList() == nil) {
		return false
	}
	ol := s.setIdleList(nil)
	ctx := context.Background()
	if ol != nil {
		p.idleList.Remove(ol)
		p.idleWriteList.Remove(ol)
		if s.isWritePrepared() {
			p.decNumWritesLocked(ctx)
		} else {
			p.decNumReadsLocked(ctx)
		}
	}
	if s.invalidate() {
		p.numOpened--
		p.recordStat(ctx, OpenSessionCount, int64(p.numOpened))
		close(p.mayGetSession)
		p.mayGetSession = make(chan struct {
		})
		return true
	}
	return false
}
func (p *sessionPool) gologoo__currSessionsCheckedOutLocked_b42872e5776c89a8da4fd9993002e327() uint64 {
	return p.numOpened - uint64(p.idleList.Len()) - uint64(p.idleWriteList.Len())
}
func (p *sessionPool) gologoo__incNumInUse_b42872e5776c89a8da4fd9993002e327(ctx context.Context) {
	p.mu.Lock()
	p.incNumInUseLocked(ctx)
	p.mu.Unlock()
}
func (p *sessionPool) gologoo__incNumInUseLocked_b42872e5776c89a8da4fd9993002e327(ctx context.Context) {
	p.numInUse++
	p.recordStat(ctx, SessionsCount, int64(p.numInUse), tagNumInUseSessions)
	p.recordStat(ctx, AcquiredSessionsCount, 1)
	if p.numInUse > p.maxNumInUse {
		p.maxNumInUse = p.numInUse
		p.recordStat(ctx, MaxInUseSessionsCount, int64(p.maxNumInUse))
	}
}
func (p *sessionPool) gologoo__decNumInUse_b42872e5776c89a8da4fd9993002e327(ctx context.Context) {
	p.mu.Lock()
	p.decNumInUseLocked(ctx)
	p.mu.Unlock()
}
func (p *sessionPool) gologoo__decNumInUseLocked_b42872e5776c89a8da4fd9993002e327(ctx context.Context) {
	p.numInUse--
	p.recordStat(ctx, SessionsCount, int64(p.numInUse), tagNumInUseSessions)
	p.recordStat(ctx, ReleasedSessionsCount, 1)
}
func (p *sessionPool) gologoo__incNumReadsLocked_b42872e5776c89a8da4fd9993002e327(ctx context.Context) {
	p.numReads++
	p.recordStat(ctx, SessionsCount, int64(p.numReads), tagNumReadSessions)
}
func (p *sessionPool) gologoo__decNumReadsLocked_b42872e5776c89a8da4fd9993002e327(ctx context.Context) {
	p.numReads--
	p.recordStat(ctx, SessionsCount, int64(p.numReads), tagNumReadSessions)
}
func (p *sessionPool) gologoo__incNumWritesLocked_b42872e5776c89a8da4fd9993002e327(ctx context.Context) {
	p.numWrites++
	p.recordStat(ctx, SessionsCount, int64(p.numWrites), tagNumWriteSessions)
}
func (p *sessionPool) gologoo__decNumWritesLocked_b42872e5776c89a8da4fd9993002e327(ctx context.Context) {
	p.numWrites--
	p.recordStat(ctx, SessionsCount, int64(p.numWrites), tagNumWriteSessions)
}
func (p *sessionPool) gologoo__incNumBeingPrepared_b42872e5776c89a8da4fd9993002e327(ctx context.Context) {
	p.mu.Lock()
	p.incNumBeingPreparedLocked(ctx)
	p.mu.Unlock()
}
func (p *sessionPool) gologoo__incNumBeingPreparedLocked_b42872e5776c89a8da4fd9993002e327(ctx context.Context) {
	p.prepareReqs++
	p.recordStat(ctx, SessionsCount, int64(p.prepareReqs), tagNumBeingPrepared)
}
func (p *sessionPool) gologoo__decNumBeingPrepared_b42872e5776c89a8da4fd9993002e327(ctx context.Context) {
	p.mu.Lock()
	p.decNumBeingPreparedLocked(ctx)
	p.mu.Unlock()
}
func (p *sessionPool) gologoo__decNumBeingPreparedLocked_b42872e5776c89a8da4fd9993002e327(ctx context.Context) {
	p.prepareReqs--
	p.recordStat(ctx, SessionsCount, int64(p.prepareReqs), tagNumBeingPrepared)
}

type hcHeap struct {
	sessions []*session
}

func (h hcHeap) gologoo__Len_b42872e5776c89a8da4fd9993002e327() int {
	return len(h.sessions)
}
func (h hcHeap) gologoo__Less_b42872e5776c89a8da4fd9993002e327(i, j int) bool {
	return h.sessions[i].getNextCheck().Before(h.sessions[j].getNextCheck())
}
func (h hcHeap) gologoo__Swap_b42872e5776c89a8da4fd9993002e327(i, j int) {
	h.sessions[i], h.sessions[j] = h.sessions[j], h.sessions[i]
	h.sessions[i].setHcIndex(i)
	h.sessions[j].setHcIndex(j)
}
func (h *hcHeap) gologoo__Push_b42872e5776c89a8da4fd9993002e327(s interface {
}) {
	ns := s.(*session)
	ns.setHcIndex(len(h.sessions))
	h.sessions = append(h.sessions, ns)
}
func (h *hcHeap) gologoo__Pop_b42872e5776c89a8da4fd9993002e327() interface {
} {
	old := h.sessions
	n := len(old)
	s := old[n-1]
	h.sessions = old[:n-1]
	s.setHcIndex(-1)
	return s
}

const maintenanceWindowSize = 10

type maintenanceWindow struct {
	mu                    sync.Mutex
	maxSessionsCheckedOut [maintenanceWindowSize]uint64
}

func (mw *maintenanceWindow) gologoo__maxSessionsCheckedOutDuringWindow_b42872e5776c89a8da4fd9993002e327() uint64 {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	var max uint64
	for _, cycleMax := range mw.maxSessionsCheckedOut {
		max = maxUint64(max, cycleMax)
	}
	return max
}
func (mw *maintenanceWindow) gologoo__updateMaxSessionsCheckedOutDuringWindow_b42872e5776c89a8da4fd9993002e327(currNumSessionsCheckedOut uint64) {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	mw.maxSessionsCheckedOut[0] = maxUint64(currNumSessionsCheckedOut, mw.maxSessionsCheckedOut[0])
}
func (mw *maintenanceWindow) gologoo__startNewCycle_b42872e5776c89a8da4fd9993002e327(currNumSessionsCheckedOut uint64) {
	mw.mu.Lock()
	defer mw.mu.Unlock()
	copy(mw.maxSessionsCheckedOut[1:], mw.maxSessionsCheckedOut[:9])
	mw.maxSessionsCheckedOut[0] = currNumSessionsCheckedOut
}
func gologoo__newMaintenanceWindow_b42872e5776c89a8da4fd9993002e327(maxOpened uint64) *maintenanceWindow {
	mw := &maintenanceWindow{}
	for i := 0; i < maintenanceWindowSize; i++ {
		mw.maxSessionsCheckedOut[i] = maxOpened
	}
	return mw
}

type healthChecker struct {
	mu             sync.Mutex
	queue          hcHeap
	interval       time.Duration
	workers        int
	waitWorkers    sync.WaitGroup
	pool           *sessionPool
	sampleInterval time.Duration
	ready          chan struct {
	}
	done chan struct {
	}
	once             sync.Once
	maintainerCancel func()
}

func gologoo__newHealthChecker_b42872e5776c89a8da4fd9993002e327(interval time.Duration, workers int, sampleInterval time.Duration, pool *sessionPool) *healthChecker {
	if workers <= 0 {
		workers = 1
	}
	hc := &healthChecker{interval: interval, workers: workers, pool: pool, sampleInterval: sampleInterval, ready: make(chan struct {
	}), done: make(chan struct {
	}), maintainerCancel: func() {
	}}
	hc.waitWorkers.Add(1)
	go hc.maintainer()
	for i := 1; i <= hc.workers; i++ {
		hc.waitWorkers.Add(1)
		go hc.worker(i)
	}
	return hc
}
func (hc *healthChecker) gologoo__close_b42872e5776c89a8da4fd9993002e327() {
	hc.mu.Lock()
	hc.maintainerCancel()
	hc.mu.Unlock()
	hc.once.Do(func() {
		close(hc.done)
	})
	hc.waitWorkers.Wait()
}
func (hc *healthChecker) gologoo__isClosing_b42872e5776c89a8da4fd9993002e327() bool {
	select {
	case <-hc.done:
		return true
	default:
		return false
	}
}
func (hc *healthChecker) gologoo__getInterval_b42872e5776c89a8da4fd9993002e327() time.Duration {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	return hc.interval
}
func (hc *healthChecker) gologoo__scheduledHCLocked_b42872e5776c89a8da4fd9993002e327(s *session) {
	var constPart, randPart float64
	if !s.firstHCDone {
		constPart = float64(hc.interval) * 0.2
		randPart = hc.pool.rand.Float64() * float64(hc.interval) * 0.9
		s.firstHCDone = true
	} else {
		constPart = float64(hc.interval) * 0.9
		randPart = hc.pool.rand.Float64() * float64(hc.interval) * 0.2
	}
	nsFromNow := int64(math.Ceil(constPart + randPart))
	s.setNextCheck(time.Now().Add(time.Duration(nsFromNow)))
	if hi := s.getHcIndex(); hi != -1 {
		heap.Fix(&hc.queue, hi)
	}
}
func (hc *healthChecker) gologoo__scheduledHC_b42872e5776c89a8da4fd9993002e327(s *session) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.scheduledHCLocked(s)
}
func (hc *healthChecker) gologoo__register_b42872e5776c89a8da4fd9993002e327(s *session) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.scheduledHCLocked(s)
	heap.Push(&hc.queue, s)
}
func (hc *healthChecker) gologoo__unregister_b42872e5776c89a8da4fd9993002e327(s *session) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	oi := s.setHcIndex(-1)
	if oi >= 0 {
		heap.Remove(&hc.queue, oi)
	}
}
func (hc *healthChecker) gologoo__markDone_b42872e5776c89a8da4fd9993002e327(s *session) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	s.checkingHealth = false
}
func (hc *healthChecker) gologoo__healthCheck_b42872e5776c89a8da4fd9993002e327(s *session) {
	defer hc.markDone(s)
	if !s.pool.isValid() {
		s.destroy(false)
		return
	}
	if err := s.ping(); isSessionNotFoundError(err) {
		s.destroy(false)
	}
}
func (hc *healthChecker) gologoo__worker_b42872e5776c89a8da4fd9993002e327(i int) {
	getNextForPing := func() *session {
		hc.pool.mu.Lock()
		defer hc.pool.mu.Unlock()
		hc.mu.Lock()
		defer hc.mu.Unlock()
		if hc.queue.Len() <= 0 {
			return nil
		}
		s := hc.queue.sessions[0]
		if s.getNextCheck().After(time.Now()) && hc.pool.valid {
			return nil
		}
		hc.scheduledHCLocked(s)
		if !s.checkingHealth {
			s.checkingHealth = true
			return s
		}
		return nil
	}
	getNextForTx := func() *session {
		hc.pool.mu.Lock()
		defer hc.pool.mu.Unlock()
		if hc.pool.shouldPrepareWriteLocked() {
			if hc.pool.idleList.Len() > 0 && hc.pool.valid {
				hc.mu.Lock()
				defer hc.mu.Unlock()
				if hc.pool.idleList.Front().Value.(*session).checkingHealth {
					return nil
				}
				session := hc.pool.idleList.Remove(hc.pool.idleList.Front()).(*session)
				ctx := context.Background()
				hc.pool.decNumReadsLocked(ctx)
				session.checkingHealth = true
				hc.pool.incNumBeingPreparedLocked(ctx)
				return session
			}
		}
		return nil
	}
	for {
		if hc.isClosing() {
			hc.waitWorkers.Done()
			return
		}
		ws := getNextForTx()
		if ws != nil {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			err := ws.prepareForWrite(ctx)
			if err != nil {
				serr := ToSpannerError(err).(*Error)
				if serr.Code != codes.PermissionDenied {
					logf(hc.pool.sc.logger, "Failed to prepare session, error: %v", serr)
				}
			}
			hc.pool.recycle(ws)
			hc.pool.mu.Lock()
			hc.pool.decNumBeingPreparedLocked(ctx)
			hc.pool.mu.Unlock()
			cancel()
			hc.markDone(ws)
		}
		rs := getNextForPing()
		if rs == nil {
			if ws == nil {
				pause := int64(100 * time.Millisecond)
				if pause > int64(hc.interval) {
					pause = int64(hc.interval)
				}
				select {
				case <-time.After(time.Duration(rand.Int63n(pause) + pause/2)):
				case <-hc.done:
				}
			}
			continue
		}
		hc.healthCheck(rs)
	}
}
func (hc *healthChecker) gologoo__maintainer_b42872e5776c89a8da4fd9993002e327() {
	<-hc.ready
	for iteration := uint64(0); ; iteration++ {
		if hc.isClosing() {
			hc.waitWorkers.Done()
			return
		}
		hc.pool.mu.Lock()
		currSessionsOpened := hc.pool.numOpened
		maxIdle := hc.pool.MaxIdle
		minOpened := hc.pool.MinOpened
		now := time.Now()
		if now.After(hc.pool.lastResetTime.Add(10 * time.Minute)) {
			hc.pool.maxNumInUse = hc.pool.numInUse
			hc.pool.recordStat(context.Background(), MaxInUseSessionsCount, int64(hc.pool.maxNumInUse))
			hc.pool.lastResetTime = now
		}
		hc.pool.mu.Unlock()
		maxSessionsInUseDuringWindow := hc.pool.mw.maxSessionsCheckedOutDuringWindow()
		hc.mu.Lock()
		ctx, cancel := context.WithTimeout(context.Background(), hc.sampleInterval)
		hc.maintainerCancel = cancel
		hc.mu.Unlock()
		if currSessionsOpened < minOpened {
			if err := hc.growPoolInBatch(ctx, minOpened); err != nil {
				logf(hc.pool.sc.logger, "failed to grow pool: %v", err)
			}
		} else if maxIdle+maxSessionsInUseDuringWindow < currSessionsOpened {
			hc.shrinkPool(ctx, maxIdle+maxSessionsInUseDuringWindow)
		}
		select {
		case <-ctx.Done():
		case <-hc.done:
			cancel()
		}
		hc.pool.mu.Lock()
		currSessionsInUse := hc.pool.currSessionsCheckedOutLocked()
		hc.pool.mu.Unlock()
		hc.pool.mw.startNewCycle(currSessionsInUse)
	}
}
func (hc *healthChecker) gologoo__growPoolInBatch_b42872e5776c89a8da4fd9993002e327(ctx context.Context, growToNumSessions uint64) error {
	hc.pool.mu.Lock()
	defer hc.pool.mu.Unlock()
	numSessions := growToNumSessions - hc.pool.numOpened
	return hc.pool.growPoolLocked(numSessions, false)
}
func (hc *healthChecker) gologoo__shrinkPool_b42872e5776c89a8da4fd9993002e327(ctx context.Context, shrinkToNumSessions uint64) {
	hc.pool.mu.Lock()
	maxSessionsToDelete := int(hc.pool.numOpened - shrinkToNumSessions)
	hc.pool.mu.Unlock()
	var deleted int
	var prevNumOpened uint64 = math.MaxUint64
	for {
		if ctx.Err() != nil {
			return
		}
		p := hc.pool
		p.mu.Lock()
		if p.numOpened >= prevNumOpened {
			p.mu.Unlock()
			break
		}
		prevNumOpened = p.numOpened
		if shrinkToNumSessions >= p.numOpened || deleted >= maxSessionsToDelete {
			p.mu.Unlock()
			break
		}
		var s *session
		if p.idleList.Len() > 0 {
			s = p.idleList.Front().Value.(*session)
		} else if p.idleWriteList.Len() > 0 {
			s = p.idleWriteList.Front().Value.(*session)
		}
		p.mu.Unlock()
		if s != nil {
			deleted++
			s.destroy(true)
		} else {
			break
		}
	}
}
func gologoo__maxUint64_b42872e5776c89a8da4fd9993002e327(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}
func gologoo__minUint64_b42872e5776c89a8da4fd9993002e327(a, b uint64) uint64 {
	if a > b {
		return b
	}
	return a
}

const sessionResourceType = "type.googleapis.com/google.spanner.v1.Session"

func gologoo__isSessionNotFoundError_b42872e5776c89a8da4fd9993002e327(err error) bool {
	if err == nil {
		return false
	}
	if ErrCode(err) == codes.NotFound {
		if rt, ok := extractResourceType(err); ok {
			return rt == sessionResourceType
		}
	}
	return false
}
func (sh *sessionHandle) recycle() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__recycle_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	sh.gologoo__recycle_b42872e5776c89a8da4fd9993002e327()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (sh *sessionHandle) getID() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getID_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := sh.gologoo__getID_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (sh *sessionHandle) getClient() *vkit.Client {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getClient_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := sh.gologoo__getClient_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (sh *sessionHandle) getMetadata() metadata.MD {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getMetadata_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := sh.gologoo__getMetadata_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (sh *sessionHandle) getTransactionID() transactionID {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getTransactionID_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := sh.gologoo__getTransactionID_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (sh *sessionHandle) destroy() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__destroy_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	sh.gologoo__destroy_b42872e5776c89a8da4fd9993002e327()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *session) isValid() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isValid_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__isValid_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *session) isWritePrepared() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isWritePrepared_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__isWritePrepared_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *session) String() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__String_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__String_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *session) ping() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ping_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__ping_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *session) setHcIndex(i int) int {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__setHcIndex_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", i)
	r0 := s.gologoo__setHcIndex_b42872e5776c89a8da4fd9993002e327(i)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *session) setIdleList(le *list.Element) *list.Element {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__setIdleList_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", le)
	r0 := s.gologoo__setIdleList_b42872e5776c89a8da4fd9993002e327(le)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *session) invalidate() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__invalidate_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__invalidate_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *session) setNextCheck(t time.Time) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__setNextCheck_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", t)
	s.gologoo__setNextCheck_b42872e5776c89a8da4fd9993002e327(t)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *session) setTransactionID(tx transactionID) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__setTransactionID_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", tx)
	s.gologoo__setTransactionID_b42872e5776c89a8da4fd9993002e327(tx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *session) getID() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getID_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__getID_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *session) getHcIndex() int {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getHcIndex_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__getHcIndex_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *session) getIdleList() *list.Element {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getIdleList_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__getIdleList_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *session) getNextCheck() time.Time {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getNextCheck_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := s.gologoo__getNextCheck_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *session) recycle() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__recycle_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	s.gologoo__recycle_b42872e5776c89a8da4fd9993002e327()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *session) destroy(isExpire bool) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__destroy_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", isExpire)
	r0 := s.gologoo__destroy_b42872e5776c89a8da4fd9993002e327(isExpire)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *session) destroyWithContext(ctx context.Context, isExpire bool) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__destroyWithContext_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v %v\n", ctx, isExpire)
	r0 := s.gologoo__destroyWithContext_b42872e5776c89a8da4fd9993002e327(ctx, isExpire)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (s *session) delete(ctx context.Context) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__delete_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	s.gologoo__delete_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (s *session) prepareForWrite(ctx context.Context) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__prepareForWrite_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	r0 := s.gologoo__prepareForWrite_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errMinOpenedGTMaxOpened(maxOpened, minOpened uint64) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errMinOpenedGTMaxOpened_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v %v\n", maxOpened, minOpened)
	r0 := gologoo__errMinOpenedGTMaxOpened_b42872e5776c89a8da4fd9993002e327(maxOpened, minOpened)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errWriteFractionOutOfRange(writeFraction float64) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errWriteFractionOutOfRange_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", writeFraction)
	r0 := gologoo__errWriteFractionOutOfRange_b42872e5776c89a8da4fd9993002e327(writeFraction)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errHealthCheckWorkersNegative(workers int) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errHealthCheckWorkersNegative_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", workers)
	r0 := gologoo__errHealthCheckWorkersNegative_b42872e5776c89a8da4fd9993002e327(workers)
	log.Printf("Output: %v\n", r0)
	return r0
}
func errHealthCheckIntervalNegative(interval time.Duration) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errHealthCheckIntervalNegative_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", interval)
	r0 := gologoo__errHealthCheckIntervalNegative_b42872e5776c89a8da4fd9993002e327(interval)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (spc *SessionPoolConfig) validate() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__validate_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := spc.gologoo__validate_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func newSessionPool(sc *sessionClient, config SessionPoolConfig) (*sessionPool, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__newSessionPool_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v %v\n", sc, config)
	r0, r1 := gologoo__newSessionPool_b42872e5776c89a8da4fd9993002e327(sc, config)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *sessionPool) recordStat(ctx context.Context, m *stats.Int64Measure, n int64, tags ...tag.Tag) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__recordStat_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v %v %v %v\n", ctx, m, n, tags)
	p.gologoo__recordStat_b42872e5776c89a8da4fd9993002e327(ctx, m, n, tags...)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *sessionPool) initPool(numSessions uint64) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__initPool_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", numSessions)
	r0 := p.gologoo__initPool_b42872e5776c89a8da4fd9993002e327(numSessions)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *sessionPool) growPoolLocked(numSessions uint64, distributeOverChannels bool) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__growPoolLocked_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v %v\n", numSessions, distributeOverChannels)
	r0 := p.gologoo__growPoolLocked_b42872e5776c89a8da4fd9993002e327(numSessions, distributeOverChannels)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *sessionPool) sessionReady(s *session) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__sessionReady_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", s)
	p.gologoo__sessionReady_b42872e5776c89a8da4fd9993002e327(s)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *sessionPool) sessionCreationFailed(err error, numSessions int32) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__sessionCreationFailed_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v %v\n", err, numSessions)
	p.gologoo__sessionCreationFailed_b42872e5776c89a8da4fd9993002e327(err, numSessions)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *sessionPool) isValid() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isValid_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := p.gologoo__isValid_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *sessionPool) close(ctx context.Context) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__close_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	p.gologoo__close_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func deleteSession(ctx context.Context, s *session, wg *sync.WaitGroup) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__deleteSession_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v %v %v\n", ctx, s, wg)
	gologoo__deleteSession_b42872e5776c89a8da4fd9993002e327(ctx, s, wg)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *sessionPool) newSessionHandle(s *session) (sh *sessionHandle) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__newSessionHandle_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", s)
	sh = p.gologoo__newSessionHandle_b42872e5776c89a8da4fd9993002e327(s)
	log.Printf("Output: %v\n", sh)
	return
}
func (p *sessionPool) errGetSessionTimeout(ctx context.Context) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errGetSessionTimeout_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	r0 := p.gologoo__errGetSessionTimeout_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *sessionPool) errGetBasicSessionTimeout(code codes.Code) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errGetBasicSessionTimeout_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", code)
	r0 := p.gologoo__errGetBasicSessionTimeout_b42872e5776c89a8da4fd9993002e327(code)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *sessionPool) errGetSessionTimeoutWithTrackedSessionHandles(code codes.Code) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__errGetSessionTimeoutWithTrackedSessionHandles_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", code)
	r0 := p.gologoo__errGetSessionTimeoutWithTrackedSessionHandles_b42872e5776c89a8da4fd9993002e327(code)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *sessionPool) getTrackedSessionHandleStacksLocked() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getTrackedSessionHandleStacksLocked_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := p.gologoo__getTrackedSessionHandleStacksLocked_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *sessionPool) shouldPrepareWriteLocked() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__shouldPrepareWriteLocked_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := p.gologoo__shouldPrepareWriteLocked_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *sessionPool) createSession(ctx context.Context) (*session, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__createSession_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	r0, r1 := p.gologoo__createSession_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *sessionPool) isHealthy(s *session) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isHealthy_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", s)
	r0 := p.gologoo__isHealthy_b42872e5776c89a8da4fd9993002e327(s)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *sessionPool) take(ctx context.Context) (*sessionHandle, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__take_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	r0, r1 := p.gologoo__take_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *sessionPool) takeWriteSession(ctx context.Context) (*sessionHandle, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__takeWriteSession_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	r0, r1 := p.gologoo__takeWriteSession_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (p *sessionPool) recycle(s *session) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__recycle_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", s)
	r0 := p.gologoo__recycle_b42872e5776c89a8da4fd9993002e327(s)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *sessionPool) remove(s *session, isExpire bool) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__remove_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v %v\n", s, isExpire)
	r0 := p.gologoo__remove_b42872e5776c89a8da4fd9993002e327(s, isExpire)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *sessionPool) currSessionsCheckedOutLocked() uint64 {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__currSessionsCheckedOutLocked_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := p.gologoo__currSessionsCheckedOutLocked_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (p *sessionPool) incNumInUse(ctx context.Context) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__incNumInUse_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	p.gologoo__incNumInUse_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *sessionPool) incNumInUseLocked(ctx context.Context) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__incNumInUseLocked_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	p.gologoo__incNumInUseLocked_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *sessionPool) decNumInUse(ctx context.Context) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decNumInUse_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	p.gologoo__decNumInUse_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *sessionPool) decNumInUseLocked(ctx context.Context) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decNumInUseLocked_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	p.gologoo__decNumInUseLocked_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *sessionPool) incNumReadsLocked(ctx context.Context) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__incNumReadsLocked_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	p.gologoo__incNumReadsLocked_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *sessionPool) decNumReadsLocked(ctx context.Context) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decNumReadsLocked_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	p.gologoo__decNumReadsLocked_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *sessionPool) incNumWritesLocked(ctx context.Context) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__incNumWritesLocked_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	p.gologoo__incNumWritesLocked_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *sessionPool) decNumWritesLocked(ctx context.Context) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decNumWritesLocked_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	p.gologoo__decNumWritesLocked_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *sessionPool) incNumBeingPrepared(ctx context.Context) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__incNumBeingPrepared_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	p.gologoo__incNumBeingPrepared_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *sessionPool) incNumBeingPreparedLocked(ctx context.Context) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__incNumBeingPreparedLocked_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	p.gologoo__incNumBeingPreparedLocked_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *sessionPool) decNumBeingPrepared(ctx context.Context) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decNumBeingPrepared_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	p.gologoo__decNumBeingPrepared_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (p *sessionPool) decNumBeingPreparedLocked(ctx context.Context) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__decNumBeingPreparedLocked_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", ctx)
	p.gologoo__decNumBeingPreparedLocked_b42872e5776c89a8da4fd9993002e327(ctx)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (h hcHeap) Len() int {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Len_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := h.gologoo__Len_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (h hcHeap) Less(i, j int) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Less_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v %v\n", i, j)
	r0 := h.gologoo__Less_b42872e5776c89a8da4fd9993002e327(i, j)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (h hcHeap) Swap(i, j int) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Swap_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v %v\n", i, j)
	h.gologoo__Swap_b42872e5776c89a8da4fd9993002e327(i, j)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (h *hcHeap) Push(s interface {
}) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Push_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", s)
	h.gologoo__Push_b42872e5776c89a8da4fd9993002e327(s)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (h *hcHeap) Pop() interface {
} {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Pop_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := h.gologoo__Pop_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (mw *maintenanceWindow) maxSessionsCheckedOutDuringWindow() uint64 {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__maxSessionsCheckedOutDuringWindow_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := mw.gologoo__maxSessionsCheckedOutDuringWindow_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (mw *maintenanceWindow) updateMaxSessionsCheckedOutDuringWindow(currNumSessionsCheckedOut uint64) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__updateMaxSessionsCheckedOutDuringWindow_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", currNumSessionsCheckedOut)
	mw.gologoo__updateMaxSessionsCheckedOutDuringWindow_b42872e5776c89a8da4fd9993002e327(currNumSessionsCheckedOut)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (mw *maintenanceWindow) startNewCycle(currNumSessionsCheckedOut uint64) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__startNewCycle_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", currNumSessionsCheckedOut)
	mw.gologoo__startNewCycle_b42872e5776c89a8da4fd9993002e327(currNumSessionsCheckedOut)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func newMaintenanceWindow(maxOpened uint64) *maintenanceWindow {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__newMaintenanceWindow_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", maxOpened)
	r0 := gologoo__newMaintenanceWindow_b42872e5776c89a8da4fd9993002e327(maxOpened)
	log.Printf("Output: %v\n", r0)
	return r0
}
func newHealthChecker(interval time.Duration, workers int, sampleInterval time.Duration, pool *sessionPool) *healthChecker {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__newHealthChecker_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v %v %v %v\n", interval, workers, sampleInterval, pool)
	r0 := gologoo__newHealthChecker_b42872e5776c89a8da4fd9993002e327(interval, workers, sampleInterval, pool)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (hc *healthChecker) close() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__close_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	hc.gologoo__close_b42872e5776c89a8da4fd9993002e327()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (hc *healthChecker) isClosing() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isClosing_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := hc.gologoo__isClosing_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (hc *healthChecker) getInterval() time.Duration {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getInterval_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	r0 := hc.gologoo__getInterval_b42872e5776c89a8da4fd9993002e327()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (hc *healthChecker) scheduledHCLocked(s *session) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__scheduledHCLocked_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", s)
	hc.gologoo__scheduledHCLocked_b42872e5776c89a8da4fd9993002e327(s)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (hc *healthChecker) scheduledHC(s *session) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__scheduledHC_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", s)
	hc.gologoo__scheduledHC_b42872e5776c89a8da4fd9993002e327(s)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (hc *healthChecker) register(s *session) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__register_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", s)
	hc.gologoo__register_b42872e5776c89a8da4fd9993002e327(s)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (hc *healthChecker) unregister(s *session) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__unregister_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", s)
	hc.gologoo__unregister_b42872e5776c89a8da4fd9993002e327(s)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (hc *healthChecker) markDone(s *session) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__markDone_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", s)
	hc.gologoo__markDone_b42872e5776c89a8da4fd9993002e327(s)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (hc *healthChecker) healthCheck(s *session) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__healthCheck_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", s)
	hc.gologoo__healthCheck_b42872e5776c89a8da4fd9993002e327(s)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (hc *healthChecker) worker(i int) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__worker_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", i)
	hc.gologoo__worker_b42872e5776c89a8da4fd9993002e327(i)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (hc *healthChecker) maintainer() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__maintainer_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : (none)\n")
	hc.gologoo__maintainer_b42872e5776c89a8da4fd9993002e327()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (hc *healthChecker) growPoolInBatch(ctx context.Context, growToNumSessions uint64) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__growPoolInBatch_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v %v\n", ctx, growToNumSessions)
	r0 := hc.gologoo__growPoolInBatch_b42872e5776c89a8da4fd9993002e327(ctx, growToNumSessions)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (hc *healthChecker) shrinkPool(ctx context.Context, shrinkToNumSessions uint64) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__shrinkPool_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v %v\n", ctx, shrinkToNumSessions)
	hc.gologoo__shrinkPool_b42872e5776c89a8da4fd9993002e327(ctx, shrinkToNumSessions)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func maxUint64(a, b uint64) uint64 {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__maxUint64_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v %v\n", a, b)
	r0 := gologoo__maxUint64_b42872e5776c89a8da4fd9993002e327(a, b)
	log.Printf("Output: %v\n", r0)
	return r0
}
func minUint64(a, b uint64) uint64 {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__minUint64_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v %v\n", a, b)
	r0 := gologoo__minUint64_b42872e5776c89a8da4fd9993002e327(a, b)
	log.Printf("Output: %v\n", r0)
	return r0
}
func isSessionNotFoundError(err error) bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__isSessionNotFoundError_b42872e5776c89a8da4fd9993002e327")
	log.Printf("Input : %v\n", err)
	r0 := gologoo__isSessionNotFoundError_b42872e5776c89a8da4fd9993002e327(err)
	log.Printf("Output: %v\n", r0)
	return r0
}
