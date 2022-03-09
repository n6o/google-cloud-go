package spansql

import "log"

func gologoo__FuzzParseQuery_f08ab65815d0c22a464672319fe6c624(data []byte) int {
	if _, err := ParseQuery(string(data)); err != nil {
		return 0
	}
	return 1
}
func FuzzParseQuery(data []byte) int {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__FuzzParseQuery_f08ab65815d0c22a464672319fe6c624")
	log.Printf("Input : %v\n", data)
	r0 := gologoo__FuzzParseQuery_f08ab65815d0c22a464672319fe6c624(data)
	log.Printf("Output: %v\n", r0)
	return r0
}
