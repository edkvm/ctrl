// +build darwin

package invoke

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

func (ap ActionProvider) InvokeAction(name string, payload []byte, env []string) interface{} {

	pod := ap.newExecuter(name, "")

	execName := pod.handlerPath


	// TODO: Add context to cmd
	cmd := exec.Command(execName)

	cmd.Env = append(env, fmt.Sprintf("CTRL_INT_SOCKET=%v",pod.sockPath))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("failed to connect stdout: %v\n", err)
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("failed to connect stderr: %v\n", err)
	}
	defer stderr.Close()

	// TODO: Add Instrumentation
	if err := cmd.Start(); err != nil {
		log.Printf("failed to start cmd: %v\n", err)
	}

	outScanner := bufio.NewScanner(stdout)
	errScanner := bufio.NewScanner(stderr)
	go func() {
		for outScanner.Scan() {
			log.Printf("[%s] stdout: %v\n", name, outScanner.Text())
		}
	}()

	go func() {
		for errScanner.Scan() {
			log.Printf("[%s] stderr: %v\n", name, errScanner.Text())
		}
	}()

	log.Println("[ctrl]", "starting", "action", name)
	result, err := pod.executeRPC(pod.sockPath, payload)
	log.Println("[ctrl]", "finished", "action", name, "result", result)

	cmd.Process.Kill()
	if err := cmd.Wait(); err != nil {
		log.Print(err)
	}

	return result

}
