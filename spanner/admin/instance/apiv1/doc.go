package instance

import (
	"context"
	"os"
	"runtime"
	"strconv"
	"strings"
	"unicode"
	"google.golang.org/api/option"
	"google.golang.org/grpc/metadata"
	"log"
)

type clientHookParams struct {
}
type clientHook func(context.Context, clientHookParams) ([]option.ClientOption, error)

var versionClient string

func gologoo__getVersionClient_5cdfc515d613c010eb27c6f02da4ccff() string {
	if versionClient == "" {
		return "UNKNOWN"
	}
	return versionClient
}
func gologoo__insertMetadata_5cdfc515d613c010eb27c6f02da4ccff(ctx context.Context, mds ...metadata.MD) context.Context {
	out, _ := metadata.FromOutgoingContext(ctx)
	out = out.Copy()
	for _, md := range mds {
		for k, v := range md {
			out[k] = append(out[k], v...)
		}
	}
	return metadata.NewOutgoingContext(ctx, out)
}
func gologoo__checkDisableDeadlines_5cdfc515d613c010eb27c6f02da4ccff() (bool, error) {
	raw, ok := os.LookupEnv("GOOGLE_API_GO_EXPERIMENTAL_DISABLE_DEFAULT_DEADLINE")
	if !ok {
		return false, nil
	}
	b, err := strconv.ParseBool(raw)
	return b, err
}
func gologoo__DefaultAuthScopes_5cdfc515d613c010eb27c6f02da4ccff() []string {
	return []string{"https://www.googleapis.com/auth/cloud-platform", "https://www.googleapis.com/auth/spanner.admin"}
}
func gologoo__versionGo_5cdfc515d613c010eb27c6f02da4ccff() string {
	const develPrefix = "devel +"
	s := runtime.Version()
	if strings.HasPrefix(s, develPrefix) {
		s = s[len(develPrefix):]
		if p := strings.IndexFunc(s, unicode.IsSpace); p >= 0 {
			s = s[:p]
		}
		return s
	}
	notSemverRune := func(r rune) bool {
		return !strings.ContainsRune("0123456789.", r)
	}
	if strings.HasPrefix(s, "go1") {
		s = s[2:]
		var prerelease string
		if p := strings.IndexFunc(s, notSemverRune); p >= 0 {
			s, prerelease = s[:p], s[p:]
		}
		if strings.HasSuffix(s, ".") {
			s += "0"
		} else if strings.Count(s, ".") < 2 {
			s += ".0"
		}
		if prerelease != "" {
			s += "-" + prerelease
		}
		return s
	}
	return "UNKNOWN"
}
func getVersionClient() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__getVersionClient_5cdfc515d613c010eb27c6f02da4ccff")
	log.Printf("Input : (none)\n")
	r0 := gologoo__getVersionClient_5cdfc515d613c010eb27c6f02da4ccff()
	log.Printf("Output: %v\n", r0)
	return r0
}
func insertMetadata(ctx context.Context, mds ...metadata.MD) context.Context {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__insertMetadata_5cdfc515d613c010eb27c6f02da4ccff")
	log.Printf("Input : %v %v\n", ctx, mds)
	r0 := gologoo__insertMetadata_5cdfc515d613c010eb27c6f02da4ccff(ctx, mds...)
	log.Printf("Output: %v\n", r0)
	return r0
}
func checkDisableDeadlines() (bool, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__checkDisableDeadlines_5cdfc515d613c010eb27c6f02da4ccff")
	log.Printf("Input : (none)\n")
	r0, r1 := gologoo__checkDisableDeadlines_5cdfc515d613c010eb27c6f02da4ccff()
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func DefaultAuthScopes() []string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DefaultAuthScopes_5cdfc515d613c010eb27c6f02da4ccff")
	log.Printf("Input : (none)\n")
	r0 := gologoo__DefaultAuthScopes_5cdfc515d613c010eb27c6f02da4ccff()
	log.Printf("Output: %v\n", r0)
	return r0
}
func versionGo() string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__versionGo_5cdfc515d613c010eb27c6f02da4ccff")
	log.Printf("Input : (none)\n")
	r0 := gologoo__versionGo_5cdfc515d613c010eb27c6f02da4ccff()
	log.Printf("Output: %v\n", r0)
	return r0
}
