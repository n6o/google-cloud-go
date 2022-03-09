package spanner

import "log"

func (c *Client) gologoo__SetGoogleClientInfo_aa597e6dba14ba86c0418098d0194af2(keyval ...string) {
	c.setGoogleClientInfo(keyval...)
}
func (c *Client) SetGoogleClientInfo(keyval ...string) {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SetGoogleClientInfo_aa597e6dba14ba86c0418098d0194af2")
	log.Printf("Input : %v\n", keyval)
	c.gologoo__SetGoogleClientInfo_aa597e6dba14ba86c0418098d0194af2(keyval...)
	log.Printf("🚚 Output: %v\n", "(none)")
	return
}
