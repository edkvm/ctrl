package ctrl

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/rakyll/statik/fs"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

type StackConfig struct {
	name               string
	runnerPath               string
	tmplPath           string
	fileNames          []string
	entryPointFile     string
	entryPointFunction string
}

const (
	servicePathDefult = "/usr/local/var/ctrl"
	actionsPath       = "actions"
	stacksPath        = "stacks"
)

type stackFile struct {
	path string
	name string
	output string
}

var stacksList = map[string]StackConfig{
	"node10": {
		name:               "node10",
		runnerPath:           fmt.Sprintf("/%s/%s", "node10", "runner.js"),
		tmplPath:           fmt.Sprintf("/%s/%s", "node10", "index.js.tmpl"),
		fileNames:          []string{"__func__.js", "config.json", "params.json"},
		entryPointFile:     "__func__.js",
		entryPointFunction: "action",
	},
	"go": {
		name: "go",
	},
}

func buildActionRepoPath() string {
	return fmt.Sprintf("%s/%s", servicePathDefult, actionsPath)
}

func buildActionPath(name string) string {
	return fmt.Sprintf("%s/%s/%s", servicePathDefult, actionsPath, name)
}

func writeFile(filepath string, data []byte) error {
	fd, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = fd.Write(data)

	return err
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

func (sc *StackConfig) getTemplate() (*template.Template, error) {

	buf, err := readStaticFile(sc.tmplPath)
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
	data, err := readStaticFile(sc.runnerPath)
	if err != nil {
		log.Fatal(err)
	}

	err = writeFile(fmt.Sprintf("%s/runner.js", path), data)
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

	err = writeFile(fmt.Sprintf("%s/index.js", path), buf.Bytes())
	if err != nil {
		return err
	}

	return nil

}

func readStaticFile(path string) ([]byte, error) {
	statickFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(path)
	fs, err := statickFS.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	buf, err := ioutil.ReadAll(fs)
	if err != nil {
		log.Fatal(err)
	}

	return buf, nil
}
