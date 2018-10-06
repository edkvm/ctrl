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
	"strings"
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
		stacksPath: fmt.Sprintf("%s/assembler/stacks", dir),
		stacksList: map[string]string{"node8.9": "node8.9/index.js"},
	})
}

func run(conf mainConfig) {

	filename := createFunction()
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

	funcName := "test"
	scoektPath := fmt.Sprintf("/tmp/%s_%s.sock", "test", genULID())

	result := runFunction(
		fmt.Sprintf("%s/%s", conf.basePath, funcName),
		conf.stackPath("node8.9"),
		scoektPath,
	)

	fmt.Printf("[INFO] %s", result)

}


func runFunction(funcPath string, stackPath string, ctrlSocketPath string) string {
	cmdParams := []string{
		stackPath,
		ctrlSocketPath,
		funcPath,
	}

	addr, err := net.ResolveUnixAddr("unix", ctrlSocketPath)
	if err != nil {
		fmt.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}
	l, err := net.ListenUnix("unix", addr)
	defer func() {
		l.Close()
		os.Remove(ctrlSocketPath)
	}()
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

	go func() {
		conn, err := l.AcceptUnix()
		defer conn.Close()
		if err != nil {
			log.Println(fmt.Sprintf("Error in accepting new connection: %v\n", err))
			return
		}
		buf := make([]byte, 64)
		conn.ReadFromUnix(buf)
		op := strings.TrimSpace(string(buf))
		if op != "__connected" {
			log.Println(fmt.Sprintf("Wrong handshake from client: %v\n", op))
			return
		}
		conn.Write([]byte(fmt.Sprintf("socket message Id (%s)", genULID())))
	}()

	if err := cmd.Start(); err != nil {
		log.Println(fmt.Sprintf("Failed to start cmd: %v\n", err))
	}

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

func createFunction() string {
	filename := fmt.Sprintf("%s", "test.js")
	return filename
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

