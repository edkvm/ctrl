package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
)


type funcRunner struct {
	Name        string
	ExecId      string
	execName    string
	handlerPath string
	sockPath    string
	ctrlCh      chan struct{}
}

func newFuncRunner(name string, path string) *funcRunner {
	execId := genULID()

	return &funcRunner{
		Name:        name,
		ExecId:      execId,
		execName:    "node",
		handlerPath: fmt.Sprintf("%s/index.js", path),
		sockPath:    fmt.Sprintf("%s/tmp/%s_%s.sock", path, name, execId),
	}
}

func (fr *funcRunner) bindPipes() {

}

func (fr *funcRunner) execute(input string) string {
	cmdParams := []string{
		fr.handlerPath,
		fr.sockPath,
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

	fr.openSock(inputCh, outCh)

	inputCh <- []byte(input)
	result := <-outCh


	if err := cmd.Wait(); err != nil {
		log.Print(err)
	}

	return string(result)
}

func (fr *funcRunner) openSock(inputCh <-chan []byte, outCh chan []byte) {

	addr, err := net.ResolveUnixAddr("unix", fr.sockPath)
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
