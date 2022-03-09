package database

import "log"

func gologoo__DatabaseAdminInstancePath_35aa76d9e5bbca86106418c9d890a57a(project, instance string) string {
	return "" + "projects/" + project + "/instances/" + instance + ""
}
func gologoo__DatabaseAdminDatabasePath_35aa76d9e5bbca86106418c9d890a57a(project, instance, database string) string {
	return "" + "projects/" + project + "/instances/" + instance + "/databases/" + database + ""
}
func DatabaseAdminInstancePath(project, instance string) string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DatabaseAdminInstancePath_35aa76d9e5bbca86106418c9d890a57a")
	log.Printf("Input : %v %v\n", project, instance)
	r0 := gologoo__DatabaseAdminInstancePath_35aa76d9e5bbca86106418c9d890a57a(project, instance)
	log.Printf("Output: %v\n", r0)
	return r0
}
func DatabaseAdminDatabasePath(project, instance, database string) string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DatabaseAdminDatabasePath_35aa76d9e5bbca86106418c9d890a57a")
	log.Printf("Input : %v %v %v\n", project, instance, database)
	r0 := gologoo__DatabaseAdminDatabasePath_35aa76d9e5bbca86106418c9d890a57a(project, instance, database)
	log.Printf("Output: %v\n", r0)
	return r0
}
