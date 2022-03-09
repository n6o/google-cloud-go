package spanner

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"cloud.google.com/go/spanner/internal"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"google.golang.org/grpc/metadata"
	"log"
)

const statsPrefix = "cloud.google.com/go/spanner/"

var (
	tagKeyClientID           = tag.MustNewKey("client_id")
	tagKeyDatabase           = tag.MustNewKey("database")
	tagKeyInstance           = tag.MustNewKey("instance_id")
	tagKeyLibVersion         = tag.MustNewKey("library_version")
	tagKeyType               = tag.MustNewKey("type")
	tagCommonKeys            = []tag.Key{tagKeyClientID, tagKeyDatabase, tagKeyInstance, tagKeyLibVersion}
	tagNumInUseSessions      = tag.Tag{Key: tagKeyType, Value: "num_in_use_sessions"}
	tagNumBeingPrepared      = tag.Tag{Key: tagKeyType, Value: "num_sessions_being_prepared"}
	tagNumReadSessions       = tag.Tag{Key: tagKeyType, Value: "num_read_sessions"}
	tagNumWriteSessions      = tag.Tag{Key: tagKeyType, Value: "num_write_prepared_sessions"}
	tagKeyMethod             = tag.MustNewKey("grpc_client_method")
	gfeLatencyMetricsEnabled = false
	statsMu                  = sync.RWMutex{}
)

func gologoo__recordStat_869caf7b3e11320f46266f4992edbd1f(ctx context.Context, m *stats.Int64Measure, n int64) {
	stats.Record(ctx, m.M(n))
}

var (
	OpenSessionCount            = stats.Int64(statsPrefix+"open_session_count", "Number of sessions currently opened", stats.UnitDimensionless)
	OpenSessionCountView        = &view.View{Measure: OpenSessionCount, Aggregation: view.LastValue(), TagKeys: tagCommonKeys}
	MaxAllowedSessionsCount     = stats.Int64(statsPrefix+"max_allowed_sessions", "The maximum number of sessions allowed. Configurable by the user.", stats.UnitDimensionless)
	MaxAllowedSessionsCountView = &view.View{Measure: MaxAllowedSessionsCount, Aggregation: view.LastValue(), TagKeys: tagCommonKeys}
	SessionsCount               = stats.Int64(statsPrefix+"num_sessions_in_pool", "The number of sessions currently in use.", stats.UnitDimensionless)
	SessionsCountView           = &view.View{Measure: SessionsCount, Aggregation: view.LastValue(), TagKeys: append(tagCommonKeys, tagKeyType)}
	MaxInUseSessionsCount       = stats.Int64(statsPrefix+"max_in_use_sessions", "The maximum number of sessions in use during the last 10 minute interval.", stats.UnitDimensionless)
	MaxInUseSessionsCountView   = &view.View{Measure: MaxInUseSessionsCount, Aggregation: view.LastValue(), TagKeys: tagCommonKeys}
	GetSessionTimeoutsCount     = stats.Int64(statsPrefix+"get_session_timeouts", "The number of get sessions timeouts due to pool exhaustion.", stats.UnitDimensionless)
	GetSessionTimeoutsCountView = &view.View{Measure: GetSessionTimeoutsCount, Aggregation: view.Count(), TagKeys: tagCommonKeys}
	AcquiredSessionsCount       = stats.Int64(statsPrefix+"num_acquired_sessions", "The number of sessions acquired from the session pool.", stats.UnitDimensionless)
	AcquiredSessionsCountView   = &view.View{Measure: AcquiredSessionsCount, Aggregation: view.Count(), TagKeys: tagCommonKeys}
	ReleasedSessionsCount       = stats.Int64(statsPrefix+"num_released_sessions", "The number of sessions released by the user and pool maintainer.", stats.UnitDimensionless)
	ReleasedSessionsCountView   = &view.View{Measure: ReleasedSessionsCount, Aggregation: view.Count(), TagKeys: tagCommonKeys}
	GFELatency                  = stats.Int64(statsPrefix+"gfe_latency", "Latency between Google's network receiving an RPC and reading back the first byte of the response", stats.UnitMilliseconds)
	GFELatencyView              = &view.View{Name: "cloud.google.com/go/spanner/gfe_latency", Measure: GFELatency, Description: "Latency between Google's network receives an RPC and reads back the first byte of the response", Aggregation: view.Distribution(0.0, 0.01, 0.05, 0.1, 0.3, 0.6, 0.8, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 8.0, 10.0, 13.0, 16.0, 20.0, 25.0, 30.0, 40.0, 50.0, 65.0, 80.0, 100.0, 130.0, 160.0, 200.0, 250.0, 300.0, 400.0, 500.0, 650.0, 800.0, 1000.0, 2000.0, 5000.0, 10000.0, 20000.0, 50000.0, 100000.0), TagKeys: append(tagCommonKeys, tagKeyMethod)}
	GFEHeaderMissingCount       = stats.Int64(statsPrefix+"gfe_header_missing_count", "Number of RPC responses received without the server-timing header, most likely means that the RPC never reached Google's network", stats.UnitDimensionless)
	GFEHeaderMissingCountView   = &view.View{Name: "cloud.google.com/go/spanner/gfe_header_missing_count", Measure: GFEHeaderMissingCount, Description: "Number of RPC responses received without the server-timing header, most likely means that the RPC never reached Google's network", Aggregation: view.Count(), TagKeys: append(tagCommonKeys, tagKeyMethod)}
)

func gologoo__EnableStatViews_869caf7b3e11320f46266f4992edbd1f() error {
	return view.Register(OpenSessionCountView, MaxAllowedSessionsCountView, SessionsCountView, MaxInUseSessionsCountView, GetSessionTimeoutsCountView, AcquiredSessionsCountView, ReleasedSessionsCountView)
}
func gologoo__EnableGfeLatencyView_869caf7b3e11320f46266f4992edbd1f() error {
	setGFELatencyMetricsFlag(true)
	return view.Register(GFELatencyView)
}
func gologoo__EnableGfeHeaderMissingCountView_869caf7b3e11320f46266f4992edbd1f() error {
	setGFELatencyMetricsFlag(true)
	return view.Register(GFEHeaderMissingCountView)
}
func gologoo__EnableGfeLatencyAndHeaderMissingCountViews_869caf7b3e11320f46266f4992edbd1f() error {
	setGFELatencyMetricsFlag(true)
	return view.Register(GFELatencyView, GFEHeaderMissingCountView)
}
func gologoo__getGFELatencyMetricsFlag_869caf7b3e11320f46266f4992edbd1f() bool {
	statsMu.RLock()
	defer statsMu.RUnlock()
	return gfeLatencyMetricsEnabled
}
func gologoo__setGFELatencyMetricsFlag_869caf7b3e11320f46266f4992edbd1f(enable bool) {
	statsMu.Lock()
	gfeLatencyMetricsEnabled = enable
	statsMu.Unlock()
}
func gologoo__DisableGfeLatencyAndHeaderMissingCountViews_869caf7b3e11320f46266f4992edbd1f() {
	setGFELatencyMetricsFlag(false)
	view.Unregister(GFELatencyView, GFEHeaderMissingCountView)
}
func gologoo__captureGFELatencyStats_869caf7b3e11320f46266f4992edbd1f(ctx context.Context, md metadata.MD, keyMethod string) error {
	if len(md.Get("server-timing")) == 0 {
		recordStat(ctx, GFEHeaderMissingCount, 1)
		return nil
	}
	serverTiming := md.Get("server-timing")[0]
	gfeLatency, err := strconv.Atoi(strings.TrimPrefix(serverTiming, "gfet4t7; dur="))
	if !strings.HasPrefix(serverTiming, "gfet4t7; dur=") || err != nil {
		return err
	}
	ctx = tag.NewContext(ctx, tag.FromContext(ctx))
	ctx, err = tag.New(ctx, tag.Insert(tagKeyMethod, keyMethod))
	if err != nil {
		return err
	}
	recordStat(ctx, GFELatency, int64(gfeLatency))
	return nil
}
func gologoo__createContextAndCaptureGFELatencyMetrics_869caf7b3e11320f46266f4992edbd1f(ctx context.Context, ct *commonTags, md metadata.MD, keyMethod string) error {
	var ctxGFE, err = tag.New(ctx, tag.Upsert(tagKeyClientID, ct.clientID), tag.Upsert(tagKeyDatabase, ct.database), tag.Upsert(tagKeyInstance, ct.instance), tag.Upsert(tagKeyLibVersion, ct.libVersion))
	if err != nil {
		return err
	}
	return captureGFELatencyStats(ctxGFE, md, keyMethod)
}
func gologoo__getCommonTags_869caf7b3e11320f46266f4992edbd1f(sc *sessionClient) *commonTags {
	_, instance, database, err := parseDatabaseName(sc.database)
	if err != nil {
		return nil
	}
	return &commonTags{clientID: sc.id, database: database, instance: instance, libVersion: internal.Version}
}

type commonTags struct {
	clientID   string
	database   string
	instance   string
	libVersion string
}

func recordStat(ctx context.Context, m *stats.Int64Measure, n int64) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__recordStat_869caf7b3e11320f46266f4992edbd1f")
	log.Printf("Input : %v %v %v\n", ctx, m, n)
	gologoo__recordStat_869caf7b3e11320f46266f4992edbd1f(ctx, m, n)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func EnableStatViews() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__EnableStatViews_869caf7b3e11320f46266f4992edbd1f")
	log.Printf("Input : (none)\n")
	r0 := gologoo__EnableStatViews_869caf7b3e11320f46266f4992edbd1f()
	log.Printf("Output: %v\n", r0)
	return r0
}
func EnableGfeLatencyView() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__EnableGfeLatencyView_869caf7b3e11320f46266f4992edbd1f")
	log.Printf("Input : (none)\n")
	r0 := gologoo__EnableGfeLatencyView_869caf7b3e11320f46266f4992edbd1f()
	log.Printf("Output: %v\n", r0)
	return r0
}
func EnableGfeHeaderMissingCountView() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__EnableGfeHeaderMissingCountView_869caf7b3e11320f46266f4992edbd1f")
	log.Printf("Input : (none)\n")
	r0 := gologoo__EnableGfeHeaderMissingCountView_869caf7b3e11320f46266f4992edbd1f()
	log.Printf("Output: %v\n", r0)
	return r0
}
func EnableGfeLatencyAndHeaderMissingCountViews() error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__EnableGfeLatencyAndHeaderMissingCountViews_869caf7b3e11320f46266f4992edbd1f")
	log.Printf("Input : (none)\n")
	r0 := gologoo__EnableGfeLatencyAndHeaderMissingCountViews_869caf7b3e11320f46266f4992edbd1f()
	log.Printf("Output: %v\n", r0)
	return r0
}
func getGFELatencyMetricsFlag() bool {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getGFELatencyMetricsFlag_869caf7b3e11320f46266f4992edbd1f")
	log.Printf("Input : (none)\n")
	r0 := gologoo__getGFELatencyMetricsFlag_869caf7b3e11320f46266f4992edbd1f()
	log.Printf("Output: %v\n", r0)
	return r0
}
func setGFELatencyMetricsFlag(enable bool) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__setGFELatencyMetricsFlag_869caf7b3e11320f46266f4992edbd1f")
	log.Printf("Input : %v\n", enable)
	gologoo__setGFELatencyMetricsFlag_869caf7b3e11320f46266f4992edbd1f(enable)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func DisableGfeLatencyAndHeaderMissingCountViews() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DisableGfeLatencyAndHeaderMissingCountViews_869caf7b3e11320f46266f4992edbd1f")
	log.Printf("Input : (none)\n")
	gologoo__DisableGfeLatencyAndHeaderMissingCountViews_869caf7b3e11320f46266f4992edbd1f()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
func captureGFELatencyStats(ctx context.Context, md metadata.MD, keyMethod string) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__captureGFELatencyStats_869caf7b3e11320f46266f4992edbd1f")
	log.Printf("Input : %v %v %v\n", ctx, md, keyMethod)
	r0 := gologoo__captureGFELatencyStats_869caf7b3e11320f46266f4992edbd1f(ctx, md, keyMethod)
	log.Printf("Output: %v\n", r0)
	return r0
}
func createContextAndCaptureGFELatencyMetrics(ctx context.Context, ct *commonTags, md metadata.MD, keyMethod string) error {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__createContextAndCaptureGFELatencyMetrics_869caf7b3e11320f46266f4992edbd1f")
	log.Printf("Input : %v %v %v %v\n", ctx, ct, md, keyMethod)
	r0 := gologoo__createContextAndCaptureGFELatencyMetrics_869caf7b3e11320f46266f4992edbd1f(ctx, ct, md, keyMethod)
	log.Printf("Output: %v\n", r0)
	return r0
}
func getCommonTags(sc *sessionClient) *commonTags {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getCommonTags_869caf7b3e11320f46266f4992edbd1f")
	log.Printf("Input : %v\n", sc)
	r0 := gologoo__getCommonTags_869caf7b3e11320f46266f4992edbd1f(sc)
	log.Printf("Output: %v\n", r0)
	return r0
}
