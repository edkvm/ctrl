package main

import (
	"bufio"
	"fmt"
	"github.com/oklog/ulid"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"
)

const (
	SysWorkPath = "/tmp/exos"
	FuncPath    = "/var/lib/exos/functions"
)

type sysEnv struct {
}

type mainConfig struct {
	basePath   string
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
		basePath:   dir,
		stacksPath: fmt.Sprintf("%s/assembler/templates", dir),
		stacksList: map[string]string{"node8.9": "node8.9/index.js.tmpl"},
	})
}

func run(conf mainConfig) {

	funcName := "echo"

	curFuncPath := fmt.Sprintf("%s/handler.js", conf.basePath)

	remotePath, err := deployFunc(funcName, curFuncPath, conf.stackPath("node8.9"))
	if err != nil {
		log.Println(err)
	}

	fr := newFuncRunner(funcName, remotePath)

	result := fr.execute(fmt.Sprintf("socket message Id (%s)", genULID()))

	log.Printf("[INFO] %v", result)

}

type actionDef struct {

}

type funcRunner struct {
	Name        string
	ExecId      string
	execName    string
	handlerPath string
	pipePath    string
	ctrlCh      chan struct{}
}

func newFuncRunner(name string, path string) *funcRunner {
	execId := genULID()

	tmpDirPath := fmt.Sprintf("%s/tmp", path)
	if _, err := os.Stat(tmpDirPath); os.IsNotExist(err) {
		os.Mkdir(tmpDirPath, os.ModePerm)
	}

	ctrPipePath := fmt.Sprintf("%s/tmp/%s_%s.sock", path, name, execId)

	return &funcRunner{
		Name:        name,
		ExecId:      execId,
		execName:    "node",
		handlerPath: fmt.Sprintf("%s/index.js", path),
		pipePath:    ctrPipePath,
	}
}

func (fr *funcRunner) bindPipes() {

}

func (fr *funcRunner) execute(input string) string {
	cmdParams := []string{
		fr.handlerPath,
		fr.pipePath,
	}

	// TODO: Add context to cmd
	cmd := exec.Command(fr.execName, cmdParams...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Failed to connect stdout: %v\n", err)
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Failed to connect stderr: %v\n", err)
	}
	defer stderr.Close()

	// TODO: Add Instrumentation
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start cmd: %v\n", err)
	}

	outScanner := bufio.NewScanner(stdout)
	errScanner := bufio.NewScanner(stderr)
	go func() {
		for outScanner.Scan() {
			log.Printf("stdout scan: %v\n", outScanner.Text())
		}
	}()

	go func() {
		for errScanner.Scan() {
			log.Printf("stderr scan: %v\n", errScanner.Text())
		}
	}()

	inputCh := make(chan []byte, 0)
	outCh := make(chan []byte, 0)

	fr.openPipe(inputCh, outCh)

	inputCh <- []byte(input)
	result := <-outCh


	if err := cmd.Wait(); err != nil {
		log.Print(err)
	}

	return string(result)
}
func (fr *funcRunner) openPipe(inputCh <-chan []byte, outCh chan []byte) {

	addr, err := net.ResolveUnixAddr("unix", fr.pipePath)
	if err != nil {
		log.Println("Failed to resolve: %v", err)
		os.Exit(1)
	}

	l, err := net.ListenUnix("unix", addr)
	if err != nil {
		log.Printf("Failed to open listener: %v\n", err)
	}

	go func() {
		conn, err := l.AcceptUnix()
		if err != nil {
			log.Printf("Error start accept on conn: %v\n", err)
			return
		}
		defer conn.Close()
		if err != nil {
			log.Printf("Error in accepting new connection: %v\n", err)
			return
		}
		buf := make([]byte, 16)
		n, _, _ := conn.ReadFromUnix(buf)
		op := string(buf[:n])
		if op != "connected" {
			log.Println(fmt.Sprintf("Wrong handshake from client: %v\n", op))
			return
		}

		// TODO: Add system context
		select {
			case op := <-inputCh:
				conn.Write(op)
		}

		n, _, _ = conn.ReadFromUnix(buf)
		outCh <- buf[:n]

		l.Close()
	}()

}

func deployFunc(name string, origPath string, templatePath string) (string, error) {

	dirPath := fmt.Sprintf("%s/%s", FuncPath, name)
	os.Mkdir(dirPath, os.ModePerm)

	// Install index.js
	tmpl, _ := template.ParseFiles(templatePath)

	indexFd, err := os.Create(fmt.Sprintf("%s/index.js", dirPath))
	if err != nil {
		return "", err
	}
	defer indexFd.Close()

	w := bufio.NewWriter(indexFd)

	tmpl.Execute(w, struct {
		HandlerPath string
		HandleName  string
	}{
		HandlerPath: "handler",
		HandleName:  "myhandle",
	})
	w.Flush()

	// Copy handler
	srcFd, err := os.Open(origPath)
	if err != nil {
		return "", err
	}
	defer srcFd.Close()

	dstFd, err := os.Create(fmt.Sprintf("%s/handler.js", dirPath))
	if err != nil {
		return "", err
	}
	defer dstFd.Close()

	_, err = io.Copy(dstFd, srcFd)
	if err != nil {
		// TODO: File didn'tmpl open, Report as (SystemError)
		return "", err
	}

	return dirPath, nil
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
