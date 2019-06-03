package main

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



type sysEnv struct {
}



func funcPath(name string) string {
	return fmt.Sprintf("%s/%s/%s", servicePathDefult, functionsPath, name)
}



type actionDef struct {

}



func deployFuncLocal(srcPath, dstPath string, tmpl *template.Template) error {

	entryFileName := "handler"
	entryFuncName := "myhandle"


	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		os.MkdirAll(dstPath, os.ModePerm)
		os.Mkdir(fmt.Sprintf("%s/tmp", dstPath), os.ModePerm)
	}

	indexFd, err := os.Create(fmt.Sprintf("%s/index.js", dstPath))
	if err != nil {
		return err
	}
	defer indexFd.Close()

	w := bufio.NewWriter(indexFd)

	tmpl.Execute(w, struct {
		HandlerPath string
		HandleName  string
	}{
		HandlerPath: entryFileName,
		HandleName:  entryFuncName,
	})
	w.Flush()

	// Copy handler
	srcFd, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFd.Close()

	dstFd, err := os.Create(fmt.Sprintf("%s/%s.js", dstPath, entryFileName))
	if err != nil {
		return err
	}
	defer dstFd.Close()

	_, err = io.Copy(dstFd, srcFd)
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
