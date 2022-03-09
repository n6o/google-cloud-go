package spanner

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"cloud.google.com/go/internal/trace"
	vkit "cloud.google.com/go/spanner/apiv1"
	"cloud.google.com/go/spanner/internal"
	"google.golang.org/api/option"
	"google.golang.org/api/option/internaloption"
	gtransport "google.golang.org/api/transport/grpc"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

const (
	resourcePrefixHeader = "google-cloud-resource-prefix"
	numChannels          = 4
)
const (
	Scope      = "https://www.googleapis.com/auth/spanner.data"
	AdminScope = "https://www.googleapis.com/auth/spanner.admin"
)

var validDBPattern = regexp.MustCompile("^projects/(?P<project>[^/]+)/instances/(?P<instance>[^/]+)/databases/(?P<database>[^/]+)$")

func gologoo__validDatabaseName_86210e11dc0e68d491b41da30cf09e12(db string) error {
	if matched := validDBPattern.MatchString(db); !matched {
		return fmt.Errorf("database name %q should conform to pattern %q", db, validDBPattern.String())
	}
	return nil
}
func gologoo__parseDatabaseName_86210e11dc0e68d491b41da30cf09e12(db string) (project, instance, database string, err error) {
	matches := validDBPattern.FindStringSubmatch(db)
	if len(matches) == 0 {
		return "", "", "", fmt.Errorf("Failed to parse database name from %q according to pattern %q", db, validDBPattern.String())
	}
	return matches[1], matches[2], matches[3], nil
}

type Client struct {
	sc           *sessionClient
	idleSessions *sessionPool
	logger       *log.Logger
	qo           QueryOptions
	ct           *commonTags
}

func (c *Client) gologoo__DatabaseName_86210e11dc0e68d491b41da30cf09e12() string {
	return c.sc.database
}

type ClientConfig struct {
	NumChannels int
	SessionPoolConfig
	SessionLabels map[string]string
	QueryOptions  QueryOptions
	CallOptions   *vkit.CallOptions
	logger        *log.Logger
}

func gologoo__contextWithOutgoingMetadata_86210e11dc0e68d491b41da30cf09e12(ctx context.Context, md metadata.MD) context.Context {
	existing, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		md = metadata.Join(existing, md)
	}
	return metadata.NewOutgoingContext(ctx, md)
}
func gologoo__NewClient_86210e11dc0e68d491b41da30cf09e12(ctx context.Context, database string, opts ...option.ClientOption) (*Client, error) {
	return NewClientWithConfig(ctx, database, ClientConfig{SessionPoolConfig: DefaultSessionPoolConfig}, opts...)
}
func gologoo__NewClientWithConfig_86210e11dc0e68d491b41da30cf09e12(ctx context.Context, database string, config ClientConfig, opts ...option.ClientOption) (c *Client, err error) {
	if err := validDatabaseName(database); err != nil {
		return nil, err
	}
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.NewClient")
	defer func() {
		trace.EndSpan(ctx, err)
	}()
	if emulatorAddr := os.Getenv("SPANNER_EMULATOR_HOST"); emulatorAddr != "" {
		emulatorOpts := []option.ClientOption{option.WithEndpoint(emulatorAddr), option.WithGRPCDialOption(grpc.WithInsecure()), option.WithoutAuthentication(), internaloption.SkipDialSettingsValidation()}
		opts = append(emulatorOpts, opts...)
	}
	hasNumChannelsConfig := config.NumChannels > 0
	if config.NumChannels == 0 {
		config.NumChannels = numChannels
	}
	allOpts := allClientOpts(config.NumChannels, opts...)
	pool, err := gtransport.DialPool(ctx, allOpts...)
	if err != nil {
		return nil, err
	}
	if hasNumChannelsConfig && pool.Num() != config.NumChannels {
		pool.Close()
		return nil, spannerErrorf(codes.InvalidArgument, "Connection pool mismatch: NumChannels=%v, WithGRPCConnectionPool=%v. Only set one of these options, or set both to the same value.", config.NumChannels, pool.Num())
	}
	sessionLabels := make(map[string]string)
	for k, v := range config.SessionLabels {
		sessionLabels[k] = v
	}
	if config.MaxOpened == 0 {
		config.MaxOpened = uint64(pool.Num() * 100)
	}
	if config.MaxBurst == 0 {
		config.MaxBurst = DefaultSessionPoolConfig.MaxBurst
	}
	if config.incStep == 0 {
		config.incStep = DefaultSessionPoolConfig.incStep
	}
	sc := newSessionClient(pool, database, sessionLabels, metadata.Pairs(resourcePrefixHeader, database), config.logger, config.CallOptions)
	config.SessionPoolConfig.sessionLabels = sessionLabels
	sp, err := newSessionPool(sc, config.SessionPoolConfig)
	if err != nil {
		sc.close()
		return nil, err
	}
	c = &Client{sc: sc, idleSessions: sp, logger: config.logger, qo: getQueryOptions(config.QueryOptions), ct: getCommonTags(sc)}
	return c, nil
}
func gologoo__allClientOpts_86210e11dc0e68d491b41da30cf09e12(numChannels int, userOpts ...option.ClientOption) []option.ClientOption {
	generatedDefaultOpts := vkit.DefaultClientOptions()
	clientDefaultOpts := []option.ClientOption{option.WithGRPCConnectionPool(numChannels), option.WithUserAgent(fmt.Sprintf("spanner-go/v%s", internal.Version)), internaloption.EnableDirectPath(true)}
	allDefaultOpts := append(generatedDefaultOpts, clientDefaultOpts...)
	return append(allDefaultOpts, userOpts...)
}
func gologoo__getQueryOptions_86210e11dc0e68d491b41da30cf09e12(opts QueryOptions) QueryOptions {
	if opts.Options == nil {
		opts.Options = &sppb.ExecuteSqlRequest_QueryOptions{}
	}
	opv := os.Getenv("SPANNER_OPTIMIZER_VERSION")
	if opv != "" {
		opts.Options.OptimizerVersion = opv
	}
	opsp := os.Getenv("SPANNER_OPTIMIZER_STATISTICS_PACKAGE")
	if opsp != "" {
		opts.Options.OptimizerStatisticsPackage = opsp
	}
	return opts
}
func (c *Client) gologoo__Close_86210e11dc0e68d491b41da30cf09e12() {
	if c.idleSessions != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		c.idleSessions.close(ctx)
	}
	c.sc.close()
}
func (c *Client) gologoo__Single_86210e11dc0e68d491b41da30cf09e12() *ReadOnlyTransaction {
	t := &ReadOnlyTransaction{singleUse: true}
	t.txReadOnly.sp = c.idleSessions
	t.txReadOnly.txReadEnv = t
	t.txReadOnly.qo = c.qo
	t.txReadOnly.replaceSessionFunc = func(ctx context.Context) error {
		if t.sh == nil {
			return spannerErrorf(codes.InvalidArgument, "missing session handle on transaction")
		}
		t.sh.destroy()
		t.state = txNew
		sh, _, err := t.acquire(ctx)
		if err != nil {
			return err
		}
		t.sh = sh
		return nil
	}
	t.ct = c.ct
	return t
}
func (c *Client) gologoo__ReadOnlyTransaction_86210e11dc0e68d491b41da30cf09e12() *ReadOnlyTransaction {
	t := &ReadOnlyTransaction{singleUse: false, txReadyOrClosed: make(chan struct {
	})}
	t.txReadOnly.sp = c.idleSessions
	t.txReadOnly.txReadEnv = t
	t.txReadOnly.qo = c.qo
	t.ct = c.ct
	return t
}
func (c *Client) gologoo__BatchReadOnlyTransaction_86210e11dc0e68d491b41da30cf09e12(ctx context.Context, tb TimestampBound) (*BatchReadOnlyTransaction, error) {
	var (
		tx  transactionID
		rts time.Time
		s   *session
		sh  *sessionHandle
		err error
	)
	defer func() {
		if err != nil && sh != nil {
			s.delete(ctx)
		}
	}()
	s, err = c.sc.createSession(ctx)
	if err != nil {
		return nil, err
	}
	sh = &sessionHandle{session: s}
	res, err := sh.getClient().BeginTransaction(contextWithOutgoingMetadata(ctx, sh.getMetadata()), &sppb.BeginTransactionRequest{Session: sh.getID(), Options: &sppb.TransactionOptions{Mode: &sppb.TransactionOptions_ReadOnly_{ReadOnly: buildTransactionOptionsReadOnly(tb, true)}}})
	if err != nil {
		return nil, ToSpannerError(err)
	}
	tx = res.Id
	if res.ReadTimestamp != nil {
		rts = time.Unix(res.ReadTimestamp.Seconds, int64(res.ReadTimestamp.Nanos))
	}
	t := &BatchReadOnlyTransaction{ReadOnlyTransaction: ReadOnlyTransaction{tx: tx, txReadyOrClosed: make(chan struct {
	}), state: txActive, rts: rts}, ID: BatchReadOnlyTransactionID{tid: tx, sid: sh.getID(), rts: rts}}
	t.txReadOnly.sh = sh
	t.txReadOnly.txReadEnv = t
	t.txReadOnly.qo = c.qo
	t.ct = c.ct
	return t, nil
}
func (c *Client) gologoo__BatchReadOnlyTransactionFromID_86210e11dc0e68d491b41da30cf09e12(tid BatchReadOnlyTransactionID) *BatchReadOnlyTransaction {
	s, err := c.sc.sessionWithID(tid.sid)
	if err != nil {
		logf(c.logger, "unexpected error: %v\nThis is an indication of an internal error in the Spanner client library.", err)
		s = &session{}
	}
	sh := &sessionHandle{session: s}
	t := &BatchReadOnlyTransaction{ReadOnlyTransaction: ReadOnlyTransaction{tx: tid.tid, txReadyOrClosed: make(chan struct {
	}), state: txActive, rts: tid.rts}, ID: tid}
	t.txReadOnly.sh = sh
	t.txReadOnly.txReadEnv = t
	t.txReadOnly.qo = c.qo
	t.ct = c.ct
	return t
}

type transactionInProgressKey struct {
}

func gologoo__checkNestedTxn_86210e11dc0e68d491b41da30cf09e12(ctx context.Context) error {
	if ctx.Value(transactionInProgressKey{}) != nil {
		return spannerErrorf(codes.FailedPrecondition, "Cloud Spanner does not support nested transactions")
	}
	return nil
}
func (c *Client) gologoo__ReadWriteTransaction_86210e11dc0e68d491b41da30cf09e12(ctx context.Context, f func(context.Context, *ReadWriteTransaction) error) (commitTimestamp time.Time, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.ReadWriteTransaction")
	defer func() {
		trace.EndSpan(ctx, err)
	}()
	resp, err := c.rwTransaction(ctx, f, TransactionOptions{})
	return resp.CommitTs, err
}
func (c *Client) gologoo__ReadWriteTransactionWithOptions_86210e11dc0e68d491b41da30cf09e12(ctx context.Context, f func(context.Context, *ReadWriteTransaction) error, options TransactionOptions) (resp CommitResponse, err error) {
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.ReadWriteTransactionWithOptions")
	defer func() {
		trace.EndSpan(ctx, err)
	}()
	resp, err = c.rwTransaction(ctx, f, options)
	return resp, err
}
func (c *Client) gologoo__rwTransaction_86210e11dc0e68d491b41da30cf09e12(ctx context.Context, f func(context.Context, *ReadWriteTransaction) error, options TransactionOptions) (resp CommitResponse, err error) {
	if err := checkNestedTxn(ctx); err != nil {
		return resp, err
	}
	var sh *sessionHandle
	defer func() {
		if sh != nil {
			sh.recycle()
		}
	}()
	err = runWithRetryOnAbortedOrSessionNotFound(ctx, func(ctx context.Context) error {
		var (
			err error
			t   *ReadWriteTransaction
		)
		if sh == nil || sh.getID() == "" || sh.getClient() == nil {
			sh, err = c.idleSessions.takeWriteSession(ctx)
			if err != nil {
				return err
			}
			t = &ReadWriteTransaction{tx: sh.getTransactionID()}
		} else {
			t = &ReadWriteTransaction{}
		}
		t.txReadOnly.sh = sh
		t.txReadOnly.txReadEnv = t
		t.txReadOnly.qo = c.qo
		t.txOpts = options
		t.ct = c.ct
		trace.TracePrintf(ctx, map[string]interface {
		}{"transactionID": string(sh.getTransactionID())}, "Starting transaction attempt")
		if err = t.begin(ctx); err != nil {
			return err
		}
		resp, err = t.runInTransaction(ctx, f)
		return err
	})
	return resp, err
}

type applyOption struct {
	atLeastOnce    bool
	transactionTag string
	priority       sppb.RequestOptions_Priority
}
type ApplyOption func(*applyOption)

func gologoo__ApplyAtLeastOnce_86210e11dc0e68d491b41da30cf09e12() ApplyOption {
	return func(ao *applyOption) {
		ao.atLeastOnce = true
	}
}
func gologoo__TransactionTag_86210e11dc0e68d491b41da30cf09e12(tag string) ApplyOption {
	return func(ao *applyOption) {
		ao.transactionTag = tag
	}
}
func gologoo__Priority_86210e11dc0e68d491b41da30cf09e12(priority sppb.RequestOptions_Priority) ApplyOption {
	return func(ao *applyOption) {
		ao.priority = priority
	}
}
func (c *Client) gologoo__Apply_86210e11dc0e68d491b41da30cf09e12(ctx context.Context, ms []*Mutation, opts ...ApplyOption) (commitTimestamp time.Time, err error) {
	ao := &applyOption{}
	for _, opt := range opts {
		opt(ao)
	}
	ctx = trace.StartSpan(ctx, "cloud.google.com/go/spanner.Apply")
	defer func() {
		trace.EndSpan(ctx, err)
	}()
	if !ao.atLeastOnce {
		resp, err := c.ReadWriteTransactionWithOptions(ctx, func(ctx context.Context, t *ReadWriteTransaction) error {
			return t.BufferWrite(ms)
		}, TransactionOptions{CommitPriority: ao.priority, TransactionTag: ao.transactionTag})
		return resp.CommitTs, err
	}
	t := &writeOnlyTransaction{sp: c.idleSessions, commitPriority: ao.priority, transactionTag: ao.transactionTag}
	return t.applyAtLeastOnce(ctx, ms...)
}
func gologoo__logf_86210e11dc0e68d491b41da30cf09e12(logger *log.Logger, format string, v ...interface {
}) {
	if logger == nil {
		log.Printf(format, v...)
	} else {
		logger.Printf(format, v...)
	}
}
func validDatabaseName(db string) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__validDatabaseName_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : %v\n", db)
	r0 := gologoo__validDatabaseName_86210e11dc0e68d491b41da30cf09e12(db)
	log.Printf("Output: %v\n", r0)
	return r0
}
func parseDatabaseName(db string) (project, instance, database string, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__parseDatabaseName_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : %v\n", db)
	project, instance, database, err = gologoo__parseDatabaseName_86210e11dc0e68d491b41da30cf09e12(db)
	log.Printf("Output: %v %v %v %v\n", project, instance, database, err)
	return
}
func (c *Client) DatabaseName() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DatabaseName_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__DatabaseName_86210e11dc0e68d491b41da30cf09e12()
	log.Printf("Output: %v\n", r0)
	return r0
}
func contextWithOutgoingMetadata(ctx context.Context, md metadata.MD) context.Context {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__contextWithOutgoingMetadata_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : %v %v\n", ctx, md)
	r0 := gologoo__contextWithOutgoingMetadata_86210e11dc0e68d491b41da30cf09e12(ctx, md)
	log.Printf("Output: %v\n", r0)
	return r0
}
func NewClient(ctx context.Context, database string, opts ...option.ClientOption) (*Client, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NewClient_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : %v %v %v\n", ctx, database, opts)
	r0, r1 := gologoo__NewClient_86210e11dc0e68d491b41da30cf09e12(ctx, database, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func NewClientWithConfig(ctx context.Context, database string, config ClientConfig, opts ...option.ClientOption) (c *Client, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__NewClientWithConfig_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : %v %v %v %v\n", ctx, database, config, opts)
	c, err = gologoo__NewClientWithConfig_86210e11dc0e68d491b41da30cf09e12(ctx, database, config, opts...)
	log.Printf("Output: %v %v\n", c, err)
	return
}
func allClientOpts(numChannels int, userOpts ...option.ClientOption) []option.ClientOption {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__allClientOpts_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : %v %v\n", numChannels, userOpts)
	r0 := gologoo__allClientOpts_86210e11dc0e68d491b41da30cf09e12(numChannels, userOpts...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func getQueryOptions(opts QueryOptions) QueryOptions {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getQueryOptions_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : %v\n", opts)
	r0 := gologoo__getQueryOptions_86210e11dc0e68d491b41da30cf09e12(opts)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *Client) Close() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Close_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : (none)\n")
	c.gologoo__Close_86210e11dc0e68d491b41da30cf09e12()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func (c *Client) Single() *ReadOnlyTransaction {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Single_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__Single_86210e11dc0e68d491b41da30cf09e12()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *Client) ReadOnlyTransaction() *ReadOnlyTransaction {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ReadOnlyTransaction_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : (none)\n")
	r0 := c.gologoo__ReadOnlyTransaction_86210e11dc0e68d491b41da30cf09e12()
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *Client) BatchReadOnlyTransaction(ctx context.Context, tb TimestampBound) (*BatchReadOnlyTransaction, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__BatchReadOnlyTransaction_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : %v %v\n", ctx, tb)
	r0, r1 := c.gologoo__BatchReadOnlyTransaction_86210e11dc0e68d491b41da30cf09e12(ctx, tb)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func (c *Client) BatchReadOnlyTransactionFromID(tid BatchReadOnlyTransactionID) *BatchReadOnlyTransaction {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__BatchReadOnlyTransactionFromID_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : %v\n", tid)
	r0 := c.gologoo__BatchReadOnlyTransactionFromID_86210e11dc0e68d491b41da30cf09e12(tid)
	log.Printf("Output: %v\n", r0)
	return r0
}
func checkNestedTxn(ctx context.Context) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__checkNestedTxn_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : %v\n", ctx)
	r0 := gologoo__checkNestedTxn_86210e11dc0e68d491b41da30cf09e12(ctx)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *Client) ReadWriteTransaction(ctx context.Context, f func(context.Context, *ReadWriteTransaction) error) (commitTimestamp time.Time, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ReadWriteTransaction_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : %v %v\n", ctx, f)
	commitTimestamp, err = c.gologoo__ReadWriteTransaction_86210e11dc0e68d491b41da30cf09e12(ctx, f)
	log.Printf("Output: %v %v\n", commitTimestamp, err)
	return
}
func (c *Client) ReadWriteTransactionWithOptions(ctx context.Context, f func(context.Context, *ReadWriteTransaction) error, options TransactionOptions) (resp CommitResponse, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ReadWriteTransactionWithOptions_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : %v %v %v\n", ctx, f, options)
	resp, err = c.gologoo__ReadWriteTransactionWithOptions_86210e11dc0e68d491b41da30cf09e12(ctx, f, options)
	log.Printf("Output: %v %v\n", resp, err)
	return
}
func (c *Client) rwTransaction(ctx context.Context, f func(context.Context, *ReadWriteTransaction) error, options TransactionOptions) (resp CommitResponse, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__rwTransaction_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : %v %v %v\n", ctx, f, options)
	resp, err = c.gologoo__rwTransaction_86210e11dc0e68d491b41da30cf09e12(ctx, f, options)
	log.Printf("Output: %v %v\n", resp, err)
	return
}
func ApplyAtLeastOnce() ApplyOption {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__ApplyAtLeastOnce_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : (none)\n")
	r0 := gologoo__ApplyAtLeastOnce_86210e11dc0e68d491b41da30cf09e12()
	log.Printf("Output: %v\n", r0)
	return r0
}
func TransactionTag(tag string) ApplyOption {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__TransactionTag_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : %v\n", tag)
	r0 := gologoo__TransactionTag_86210e11dc0e68d491b41da30cf09e12(tag)
	log.Printf("Output: %v\n", r0)
	return r0
}
func Priority(priority sppb.RequestOptions_Priority) ApplyOption {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Priority_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : %v\n", priority)
	r0 := gologoo__Priority_86210e11dc0e68d491b41da30cf09e12(priority)
	log.Printf("Output: %v\n", r0)
	return r0
}
func (c *Client) Apply(ctx context.Context, ms []*Mutation, opts ...ApplyOption) (commitTimestamp time.Time, err error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__Apply_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : %v %v %v\n", ctx, ms, opts)
	commitTimestamp, err = c.gologoo__Apply_86210e11dc0e68d491b41da30cf09e12(ctx, ms, opts...)
	log.Printf("Output: %v %v\n", commitTimestamp, err)
	return
}
func logf(logger *log.Logger, format string, v ...interface {
}) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__logf_86210e11dc0e68d491b41da30cf09e12")
	log.Printf("Input : %v %v %v\n", logger, format, v)
	gologoo__logf_86210e11dc0e68d491b41da30cf09e12(logger, format, v...)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
