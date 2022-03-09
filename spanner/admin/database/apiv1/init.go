package database

import (
	"context"
	"os"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"log"
)

func gologoo__init_920555eb839f823d5b8e9ddea7e7be08() {
	newDatabaseAdminClientHook = func(ctx context.Context, p clientHookParams) ([]option.ClientOption, error) {
		if emulator := os.Getenv("SPANNER_EMULATOR_HOST"); emulator != "" {
			return []option.ClientOption{option.WithEndpoint(emulator), option.WithGRPCDialOption(grpc.WithInsecure()), option.WithoutAuthentication()}, nil
		}
		return nil, nil
	}
}
func init() {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__init_920555eb839f823d5b8e9ddea7e7be08")
	log.Printf("Input : (none)\n")
	gologoo__init_920555eb839f823d5b8e9ddea7e7be08()
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
