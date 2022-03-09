package database

import (
	"context"
	"fmt"
	"regexp"
	"time"
	pbt "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/googleapis/gax-go/v2"
	databasepb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	"log"
)

var validDBPattern = regexp.MustCompile("^projects/(?P<project>[^/]+)/instances/(?P<instance>[^/]+)/databases/(?P<database>[^/]+)$")

func (c *DatabaseAdminClient) gologoo__StartBackupOperation_42f9f61dedc83155b7588516de785fd4(ctx context.Context, backupID string, databasePath string, expireTime time.Time, opts ...gax.CallOption) (*CreateBackupOperation, error) {
	m := validDBPattern.FindStringSubmatch(databasePath)
	if m == nil {
		return nil, fmt.Errorf("database name %q should conform to pattern %q", databasePath, validDBPattern)
	}
	ts := &pbt.Timestamp{Seconds: expireTime.Unix(), Nanos: int32(expireTime.Nanosecond())}
	req := &databasepb.CreateBackupRequest{Parent: fmt.Sprintf("projects/%s/instances/%s", m[1], m[2]), BackupId: backupID, Backup: &databasepb.Backup{Database: databasePath, ExpireTime: ts}}
	return c.CreateBackup(ctx, req, opts...)
}
func (c *DatabaseAdminClient) StartBackupOperation(ctx context.Context, backupID string, databasePath string, expireTime time.Time, opts ...gax.CallOption) (*CreateBackupOperation, error) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__StartBackupOperation_42f9f61dedc83155b7588516de785fd4")
	log.Printf("Input : %v %v %v %v %v\n", ctx, backupID, databasePath, expireTime, opts)
	r0, r1 := c.gologoo__StartBackupOperation_42f9f61dedc83155b7588516de785fd4(ctx, backupID, databasePath, expireTime, opts...)
	log.Printf("Output: %v %v\n", r0, r1)
	return r0, r1
}
