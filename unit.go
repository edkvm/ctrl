package ctrl

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
)

type Action struct {
	Name        string
	ExecId      string
	execName    string
	handlerPath string
	configPath  string
	paramsPath  string
	sockPath    string
	ctrlCh      chan struct{}
}

func NewAction(name string) *Action {
	execId := genULID()
	actionPath := buildActionPath(name)
	return &Action{
		Name:        name,
		ExecId:      execId,
		execName:    "node",
		handlerPath: fmt.Sprintf("%s/index.js", actionPath),
		configPath:  fmt.Sprintf("%s/config.json", actionPath),
		paramsPath:  fmt.Sprintf("%s/params.json", actionPath),
		sockPath:    fmt.Sprintf("%s/tmp/%s_%s.sock", actionPath, name, execId),
	}
}

func (fr *Action) IsExists() bool {
	if _, err := os.Stat(fr.handlerPath); os.IsNotExist(err) {
		return false
	}

	return true
}

func (fr *Action) Execute(args []string) string {

	input := fr.parseArgs(args)

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
			log.Printf("[%s] stdout: %v\n", fr.Name, outScanner.Text())
		}
	}()

	go func() {
		for errScanner.Scan() {
			log.Printf("[%s] stderr: %v\n", fr.Name, errScanner.Text())
		}
	}()

	inputCh := make(chan []byte, 0)
	outCh := make(chan []byte, 0)

	fr.openSock(inputCh, outCh)

	inputCh <- []byte(input)

	// Wait for result
	result := <-outCh

	if err := cmd.Wait(); err != nil {
		log.Print(err)
	}

	return string(result)
}

func (fr *Action) openSock(inputCh <-chan []byte, outCh chan []byte) {

	addr, err := net.ResolveUnixAddr("unix", fr.sockPath)
	if err != nil {
		log.Println("failed to resolve: %v", err)
		os.Exit(1)
	}

	sock, err := net.ListenUnix("unix", addr)
	if err != nil {
		log.Printf("failed to open listener: %v\n", err)
	}

	go func() {
		conn, err := sock.AcceptUnix()
		if err != nil {
			log.Printf("error start accept on conn: %v\n", err)
			return
		}
		defer conn.Close()
		if err != nil {
			log.Printf("error in accepting new connection: %v\n", err)
			return
		}

		// Wait for connection
		buf := make([]byte, 256)
		n, _, _ := conn.ReadFromUnix(buf)
		op := string(buf[:n])
		if op != "op|start" {
			log.Println(fmt.Sprintf("wrong handshake from client: %v\n", op))
			return
		}

		// TODO: Add system context
		select {
		case op := <-inputCh:
			conn.Write(op)
		}


		// Wait for function output
		result := make([]byte, 256)
		for {
			n, _, err = conn.ReadFromUnix(buf)
			if err != nil {
				log.Println(err)
				break
			}

			header := string(buf[:2])
			if header == "op" {
				op = string(buf[3:8])
				if op == "close" {
					break
				}
			}
			copy(result, buf)
		}

		outCh <- result

		sock.Close()
	}()

}

func (fr *Action) parseArgs(args []string) string {
	conf := readFile(fr.configPath)
	data := readFile(fr.paramsPath)

	items := make([]interface{}, len(args))
	for i, _ := range args {
		items[i] = args[i]
	}
 	params := fmt.Sprintf(string(data), items...)

 	parsed := fmt.Sprintf(`{ "$": %s, "params": %s }`, string(conf), params)

	return parsed
}
