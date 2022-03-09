package database

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	longrunningpb "google.golang.org/genproto/googleapis/longrunning"
	databasepb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

var retryer = gax.OnCodes([]codes.Code{codes.DeadlineExceeded, codes.Unavailable}, gax.Backoff{Initial: time.Millisecond, Max: time.Millisecond, Multiplier: 1.0})

func (c *DatabaseAdminClient) gologoo__CreateDatabaseWithRetry_af67a53e7410e2053de873d6ffae273f(ctx context.Context, req *databasepb.CreateDatabaseRequest, opts ...gax.CallOption) (*CreateDatabaseOperation, error) {
	for {
		db, createErr := c.CreateDatabase(ctx, req, opts...)
		if createErr == nil {
			return db, nil
		}
		delay, shouldRetry := retryer.Retry(createErr)
		if !shouldRetry {
			return nil, createErr
		}
		if err := gax.Sleep(ctx, delay); err != nil {
			return nil, err
		}
		dbName := extractDBName(req.CreateStatement)
		iter := c.ListDatabaseOperations(ctx, &databasepb.ListDatabaseOperationsRequest{Parent: req.Parent, Filter: fmt.Sprintf("(metadata.@type:type.googleapis.com/google.spanner.admin.database.v1.CreateDatabaseMetadata) AND (name:%s/databases/%s/operations/)", req.Parent, dbName)}, opts...)
		var mostRecentOp *longrunningpb.Operation
		for {
			op, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return nil, err
			}
			if !op.Done {
				return c.CreateDatabaseOperation(op.Name), nil
			}
			if op.GetError() == nil {
				mostRecentOp = op
			}
		}
		if mostRecentOp == nil {
			continue
		}
		_, getErr := c.GetDatabase(ctx, &databasepb.GetDatabaseRequest{Name: fmt.Sprintf("%s/databases/%s", req.Parent, dbName)})
		if getErr == nil {
			return c.CreateDatabaseOperation(mostRecentOp.Name), nil
		}
		if status.Code(getErr) == codes.NotFound {
			continue
		}
		return nil, getErr
	}
}

var dbNameRegEx = regexp.MustCompile("\\s*CREATE\\s+DATABASE\\s+(.+)\\s*")

func gologoo__extractDBName_af67a53e7410e2053de873d6ffae273f(createStatement string) string {
	if dbNameRegEx.MatchString(createStatement) {
		namePossiblyWithQuotes := strings.TrimRightFunc(dbNameRegEx.FindStringSubmatch(createStatement)[1], unicode.IsSpace)
		if len(namePossiblyWithQuotes) > 0 && namePossiblyWithQuotes[0] == '`' {
			if len(namePossiblyWithQuotes) > 5 && namePossiblyWithQuotes[1] == '`' && namePossiblyWithQuotes[2] == '`' {
				return string(namePossiblyWithQuotes[3 : len(namePossiblyWithQuotes)-3])
			}
			return string(namePossiblyWithQuotes[1 : len(namePossiblyWithQuotes)-1])
		}
		return string(namePossiblyWithQuotes)
	}
	return ""
}
func (c *DatabaseAdminClient) CreateDatabaseWithRetry(ctx context.Context, req *databasepb.CreateDatabaseRequest, opts ...gax.CallOption) (*CreateDatabaseOperation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__CreateDatabaseWithRetry_af67a53e7410e2053de873d6ffae273f")
	log.Printf("Input : %v %v %v\n", ctx, req, opts)
	r0, r1 := c.gologoo__CreateDatabaseWithRetry_af67a53e7410e2053de873d6ffae273f(ctx, req, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
func extractDBName(createStatement string) string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__extractDBName_af67a53e7410e2053de873d6ffae273f")
	log.Printf("Input : %v\n", createStatement)
	r0 := gologoo__extractDBName_af67a53e7410e2053de873d6ffae273f(createStatement)
	log.Printf("Output: %v\n", r0)
	return r0
}
