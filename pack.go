package ctrl

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"text/template"
	"time"

	"github.com/oklog/ulid"
)





type StackConfig struct {
	name               string
	tmplPath           string
	entryPointFile     string
	entryPointFunction string
}




type Pack struct {
	stack StackConfig
	actionPath string
	actionTmpPath string
}

func NewPack(name string, stackName string) *Pack {
	actionPath := buildActionPath(name)
	return &Pack{
		stack: stacksList[stackName],
		actionPath: actionPath,
		actionTmpPath: fmt.Sprintf("%s/tmp", actionPath),
	}
}

func (pk *Pack) createActionWrapper() error {

	if _, err := os.Stat(pk.actionPath); os.IsNotExist(err) {
		os.MkdirAll(pk.actionPath, os.ModePerm)
		os.Mkdir(pk.actionTmpPath, os.ModePerm)
	}

	indexFd, err := os.Create(fmt.Sprintf("%s/index.js", pk.actionPath))
	if err != nil {
		return err
	}
	defer indexFd.Close()

	w := bufio.NewWriter(indexFd)

	tmpl, err := template.ParseFiles(pk.stack.tmplPath)
	if err != nil {
		return fmt.Errorf("template for stack %s was not found, %s", pk.stack.name, err)
	}

	tmpl.Execute(w, struct {
		HandlerPath string
		HandleName  string
	}{
		HandlerPath: pk.stack.entryPointFile,
		HandleName:  pk.stack.entryPointFunction,
	})
	w.Flush()

	return nil

}

func (pk *Pack) Deploy(srcData io.Reader) error {

	if err := pk.createActionWrapper(); err != nil {
		return nil
	}

	dstFd, err := os.Create(fmt.Sprintf("%s/%s.js", pk.actionPath, pk.stack.entryPointFile))
	if err != nil {
		return err
	}
	defer dstFd.Close()

	_, err = io.Copy(dstFd, srcData)
	if err != nil {
		// TODO: File didn'tmpl open, Report as (SystemError)
		return err
	}

	return nil
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
