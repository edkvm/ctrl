package packing

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/edkvm/ctrl/fs"
	"log"
	"text/template"
)

type StackConfig struct {
	name               string
	runnerPath         string
	tmplPath           string
	fileNames          []string
	entryPointFile     string
	entryPointFunction string
}



type stackFile struct {
	path   string
	name   string
	output string
}

var stacksList = map[string]StackConfig{
	"node10": {
		name:               "node10",
		runnerPath:         fmt.Sprintf("/%s/%s", "node10", "runner.js"),
		tmplPath:           fmt.Sprintf("/%s/%s", "node10", "index.js.tmpl"),
		fileNames:          []string{"__func__.js", "config.json", "params.json"},
		entryPointFile:     "__func__.js",
		entryPointFunction: "action",
	},
	"go": {
		name: "go",
	},
}


func (sc *StackConfig) getTemplate() (*template.Template, error) {

	buf, err := fs.ReadStaticFile(sc.tmplPath)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("index").Parse(string(buf))
	if err != nil {
		return nil, fmt.Errorf("template for stack %s was not found, %s", sc.tmplPath, err)
	}

	return tmpl, nil
}

func (sc *StackConfig) writeEntryPoint(path string) error {

	// Write Runner
	data, err := fs.ReadStaticFile(sc.runnerPath)
	if err != nil {
		log.Fatal(err)
	}

	err = fs.WriteFile(fmt.Sprintf("%s/runner.js", path), data)
	if err != nil {
		return err
	}

	// Write Template
	buf := bytes.Buffer{}

	w := bufio.NewWriter(&buf)

	tmpl, _ := sc.getTemplate()

	tmpl.Execute(w, struct {
		HandlerPath string
		HandleName  string
	}{
		HandlerPath: sc.entryPointFile,
		HandleName:  sc.entryPointFunction,
	})

	w.Flush()

	err = fs.WriteFile(fmt.Sprintf("%s/index.js", path), buf.Bytes())
	if err != nil {
		return err
	}

	return nil

}


