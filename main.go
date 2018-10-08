package main

import (
	"io/ioutil"
	"os/exec"
	"fmt"
	"bytes"
	"time"
	"github.com/oklog/ulid"
	"math/rand"
	"log"
	"path/filepath"
	"os"
	"net"
	"text/template"
	"bufio"
	"io"
)

const (
	SysWorkPath = "/tmp/exos"
	FuncPath = "/tmp/exos/functions"
)

type mainConfig struct {
	basePath string
	stacksPath string
	stacksList map[string]string
}

func (mc mainConfig) stackPath(name string) string {
	val := mc.stacksList[name]
	return fmt.Sprintf("%s/%s", mc.stacksPath, val)
}

func main() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	 dir = "/Users/klimslava/Projects/golang/src/github.com/edkvm/exos"
	run(mainConfig{
		basePath: dir,
		stacksPath: fmt.Sprintf("%s/assembler/templates", dir),
		stacksList: map[string]string{"node8.9": "node8.9/index.js.tmpl"},
	})
}

func run(conf mainConfig) {

	funcName := "echo"
	//createFunction(funcName, "", "js")

	curFuncPath := fmt.Sprintf("%s/handler.js", conf.basePath)

	deployFunc(funcName, curFuncPath, conf.stackPath("node8.9"))

	socketPath := fmt.Sprintf("/tmp/%s_%s.sock", funcName, genULID())

	result := runFunction(
		curFuncPath,
		conf.stackPath("node8.9"),
		socketPath,
	)

	log.Println("[INFO] %s", result)

}

func runFunction(funcPath string, templatePath string, ctrlPath string) string {
	cmdParams := []string{
		templatePath,
		ctrlPath,
		funcPath,
	}

	cmd := exec.Command("node", cmdParams...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(fmt.Sprintf("Failed to connect stdout: %v\n", err))
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(fmt.Sprintf("Failed to connect stderr: %v\n", err))
	}
	defer stderr.Close()

	// TODO: Add Instrumentation
	if err := cmd.Start(); err != nil {
		log.Println(fmt.Sprintf("Failed to start cmd: %v\n", err))
	}

	opCh := make(chan []byte, 0)
	openCtrl(ctrlPath, opCh)
	opCh <- []byte(fmt.Sprintf("socket message Id (%s)", genULID()))

	// Read from stdout
	buf := bytes.Buffer{}
	buf.ReadFrom(stdout)

	buf.ReadFrom(stderr)

	//
	if err := cmd.Wait(); err != nil {
		log.Print(err)
	}

	return buf.String()
}

func openCtrl(path string, opChan <-chan []byte) {

	addr, err := net.ResolveUnixAddr("unix", path)
	if err != nil {
		log.Println("Failed to resolve: %v", err)
		os.Exit(1)
	}

	l, err := net.ListenUnix("unix", addr)

	go func() {
		conn, err := l.AcceptUnix()
		if err != nil {
			log.Printf("Error start accept on conn: %v\n", err)
			return
		}
		defer conn.Close()
		if err != nil {
			log.Println(fmt.Sprintf("Error in accepting new connection: %v\n", err))
			return
		}
		buf := make([]byte, 16)
		n, _, _ := conn.ReadFromUnix(buf)
		op := string(buf[:n])
		if op != "connected" {
			log.Println(fmt.Sprintf("Wrong handshake from client: %v\n", op))
			return
		}

		select {
		case op := <- opChan:
			conn.Write(op)
		}

		l.Close()
		os.Remove(path)
	}()

}

func createFunction(name string, data string, ext string) string {
	filename := fmt.Sprintf("%s.%s", name, ext)
	exData := []byte(`
	exports.handler = (event, callback) => {
		console.log("with socket");
		console.log(event);
		console.log("after");
	}`)

	err := ioutil.WriteFile(filename, exData, 0640)
	if err != nil {
		// TODO: File didn't open, Report as (SystemError)
		log.Print(err)
	}

	return filename
}

func deployFunc(name string, path string, templatePath string) {


	funcDirPath := fmt.Sprintf("%s/%s", FuncPath, name)
 	os.Mkdir(funcDirPath, os.ModePerm)

 	// Install index.js
	tmpl, _ := template.ParseFiles(templatePath)

	indexFd, err := os.Create(fmt.Sprintf("%s/index.js", funcDirPath))
	if err != nil {
		log.Print(err)
		return
	}
	defer indexFd.Close()

	w := bufio.NewWriter(indexFd)

	tmpl.Execute(w, struct {
		HandlerPath string
	}{
		HandlerPath: "handler",
	})
	w.Flush()

	// Copy handler
	srcFd, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return
	}
	defer srcFd.Close()

	dstFd, err := os.Create(fmt.Sprintf("%s/handler.js", funcDirPath))
	if err != nil {
		log.Print(err)
		return
	}
	defer dstFd.Close()

	_, err = io.Copy(dstFd, srcFd)
	if err != nil {
		// TODO: File didn'tmpl open, Report as (SystemError)
		log.Print(err)
	}

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

