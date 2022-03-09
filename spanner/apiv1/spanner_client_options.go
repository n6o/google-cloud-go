package spanner

import (
	"google.golang.org/api/option"
	"log"
)

func gologoo__DefaultClientOptions_96fb67d5e672df001004eb67ad5f1189() []option.ClientOption {
	return defaultGRPCClientOptions()
}
func DefaultClientOptions() []option.ClientOption {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DefaultClientOptions_96fb67d5e672df001004eb67ad5f1189")
	log.Printf("Input : (none)\n")
	r0 := gologoo__DefaultClientOptions_96fb67d5e672df001004eb67ad5f1189()
	log.Printf("Output: %v\n", r0)
	return r0
}
