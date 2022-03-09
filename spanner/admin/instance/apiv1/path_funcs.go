package instance

import "log"

func gologoo__InstanceAdminProjectPath_4620876520013748dcd9b45e831880fd(project string) string {
	return "" + "projects/" + project + ""
}
func gologoo__InstanceAdminInstanceConfigPath_4620876520013748dcd9b45e831880fd(project, instanceConfig string) string {
	return "" + "projects/" + project + "/instanceConfigs/" + instanceConfig + ""
}
func gologoo__InstanceAdminInstancePath_4620876520013748dcd9b45e831880fd(project, instance string) string {
	return "" + "projects/" + project + "/instances/" + instance + ""
}
func InstanceAdminProjectPath(project string) string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__InstanceAdminProjectPath_4620876520013748dcd9b45e831880fd")
	log.Printf("Input : %v\n", project)
	r0 := gologoo__InstanceAdminProjectPath_4620876520013748dcd9b45e831880fd(project)
	log.Printf("Output: %v\n", r0)
	return r0
}
func InstanceAdminInstanceConfigPath(project, instanceConfig string) string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__InstanceAdminInstanceConfigPath_4620876520013748dcd9b45e831880fd")
	log.Printf("Input : %v %v\n", project, instanceConfig)
	r0 := gologoo__InstanceAdminInstanceConfigPath_4620876520013748dcd9b45e831880fd(project, instanceConfig)
	log.Printf("Output: %v\n", r0)
	return r0
}
func InstanceAdminInstancePath(project, instance string) string {
	log.SetFlags(19)
	log.Printf("📨 Call %s\n", "gologoo__InstanceAdminInstancePath_4620876520013748dcd9b45e831880fd")
	log.Printf("Input : %v %v\n", project, instance)
	r0 := gologoo__InstanceAdminInstancePath_4620876520013748dcd9b45e831880fd(project, instance)
	log.Printf("Output: %v\n", r0)
	return r0
}
