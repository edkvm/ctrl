package ctrl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/oklog/ulid"

	"github.com/rakyll/statik/fs"
	_ "github.com/edkvm/ctrl/statik"
)


type Pack struct {
	stack StackConfig
	name  string
	files  map[string][]byte}

type pack struct {

}

func BuildPack(stackName, wd string) (*Pack, error) {
	// TODO: Add more error handeling
	dirs := strings.Split(wd, "/")
	if len(dirs) < 2 {
		// TODO: return error, name is not absolute
	}

	// Action name is the folder name
	funcName := dirs[len(dirs) - 1]

	log.Println("building function:", funcName)

	pk := &Pack{
		stack: stacksList[stackName],
		name: funcName,
		files: make(map[string][]byte, 3),
	}

	pk.build(wd)

	return pk, nil
}

func (pk *Pack) build(wd string) error {
	// Read action
	for i := 0; i < len(pk.stack.fileNames); i++ {
		fileName := pk.stack.fileNames[i]
		pk.files[fileName] = readFile(fmt.Sprintf("%s/%s", wd, fileName))
	}

	return nil
}



func getTemplate(path string) (*template.Template, error) {
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

	tmpl, err := template.New("index").Parse(string(buf))
	if err != nil {
		return nil, fmt.Errorf("template for stack %s was not found, %s", path, err)
	}

	return tmpl, nil
}

func writeWrapper(path string, stack StackConfig) error {

	indexFd, err := os.Create(fmt.Sprintf("%s/index.js", path))
	if err != nil {
		return err
	}
	defer indexFd.Close()

	w := bufio.NewWriter(indexFd)


	tmpl, _ := getTemplate(stack.tmplPath)

	tmpl.Execute(w, struct {
		HandlerPath string
		HandleName  string
	}{
		HandlerPath: stack.entryPointFile,
		HandleName:  stack.entryPointFunction,
	})
	w.Flush()

	return nil

}

func (pk *Pack) Deploy() error {


	actionPath := buildActionPath(pk.name)

	if _, err := os.Stat(actionPath); os.IsNotExist(err) {
		err := os.MkdirAll(fmt.Sprintf("%s/tmp", actionPath), os.ModePerm)
		if err != nil {
			return err
		}
	}

	if err := writeWrapper(actionPath, pk.stack); err != nil {
		return nil
	}

	for name, _ := range pk.files {
		file := pk.files[name]
		src := bytes.NewReader(file)

		dstFd, err := os.Create(fmt.Sprintf("%s/%s", actionPath, name))
		if err != nil {
			return err
		}
		defer dstFd.Close()

		_, err = io.Copy(dstFd, src)
		if err != nil {
			// TODO: File didn't open, Report as (SystemError)
			return err
		}

	}

	return nil
}

func deployLocal() {

}

func updateFunction() {

}

func updateFunctionResources() {

}

func createFunctionSchedual() {

}

func pauseFunction() {

}

func deleteFunction() {

}

func genULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id, err := ulid.New(ulid.Timestamp(t), entropy)
	if err != nil {

	}

	return fmt.Sprintf("%s", id)
}
