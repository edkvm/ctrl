package ctrl

import "fmt"

const (
	servicePathDefult = "/usr/local/var/ctrl"
	functionsPath = "actions"
	templatePath = "assembler/templates"
)

var stacksList = map[string]StackConfig{
	"node8.9": {
		name:               "node8.9",
		tmplPath:           fmt.Sprintf("%s/%s/%s", servicePathDefult, templatePath, "node8.9/index.js.tmpl" ),
		entryPointFile:     "action",
		entryPointFunction: "main",
	},
}

func buildActionPath(name string) string {
	return fmt.Sprintf("%s/%s/%s", servicePathDefult, functionsPath, name)
}

//func buildActionHandlerPath() string{
//	return fmt.Sprintf("%s/index.js", actionPath)
//}
//
//func buildActionSockPath() string {
//	return fmt.Sprintf("%s/tmp/%s_%s.sock", actionPath, name, execId)
//}