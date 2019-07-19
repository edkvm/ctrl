package ctrl

import (
	"fmt"
	"io/ioutil"
	"os"
)

type StackConfig struct {
	name               string
	tmplPath           string
	fileNames          []string
	entryPointFile     string
	entryPointFunction string
}

const (
	servicePathDefult = "/usr/local/var/ctrl"
	functionsPath     = "actions"
)

var stacksList = map[string]StackConfig{
	"node_v10": {
		name:               "node_v10",
		tmplPath:           fmt.Sprintf("/%s/%s", "node_v10", "index.js.tmpl"),
		fileNames:          []string{"__func__.js", "config.json", "params.json"},
		entryPointFile:     "__func__.js",
		entryPointFunction: "main",
	},
}

func buildActionPath(name string) string {
	return fmt.Sprintf("%s/%s/%s", servicePathDefult, functionsPath, name)
}

func readFile(filePath string) []byte {

	_, err := os.Stat(filePath)
	if err != nil {
		return nil
	}

	// Copy handler
	srcFd, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer srcFd.Close()

	data, err := ioutil.ReadAll(srcFd)
	if err != nil {
		return nil
	}

	return data
}
