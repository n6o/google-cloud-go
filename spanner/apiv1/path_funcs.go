package spanner

import "log"

func gologoo__DatabasePath_bf593097abed8566e2f7ad6c210cd992(project, instance, database string) string {
	return "" + "projects/" + project + "/instances/" + instance + "/databases/" + database + ""
}
func gologoo__SessionPath_bf593097abed8566e2f7ad6c210cd992(project, instance, database, session string) string {
	return "" + "projects/" + project + "/instances/" + instance + "/databases/" + database + "/sessions/" + session + ""
}
func DatabasePath(project, instance, database string) string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__DatabasePath_bf593097abed8566e2f7ad6c210cd992")
	log.Printf("Input : %v %v %v\n", project, instance, database)
	r0 := gologoo__DatabasePath_bf593097abed8566e2f7ad6c210cd992(project, instance, database)
	log.Printf("Output: %v\n", r0)
	return r0
}
func SessionPath(project, instance, database, session string) string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__SessionPath_bf593097abed8566e2f7ad6c210cd992")
	log.Printf("Input : %v %v %v %v\n", project, instance, database, session)
	r0 := gologoo__SessionPath_bf593097abed8566e2f7ad6c210cd992(project, instance, database, session)
	log.Printf("Output: %v\n", r0)
	return r0
}
