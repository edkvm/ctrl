package invoke

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/edkvm/ctrl"
	"github.com/edkvm/ctrl/action"

	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"time"

	ctrlFS "github.com/edkvm/ctrl/pkg/fs"
	ctrlID "github.com/edkvm/ctrl/pkg/id"
)

type ActionProvider struct {
	sl *ctrl.ServiceLoc
}

func NewActionProvider(sl *ctrl.ServiceLoc) *ActionProvider {
	return &ActionProvider{
		sl,
	}
}

func (ap *ActionProvider) BuildAction(name string) (*action.Action, error) {
	actionPath := ap.sl.ActionPath(name)
	paramDef := ctrlFS.ReadFile(fmt.Sprintf("%s/params.json", actionPath))
	configDef := ctrlFS.ReadFile(fmt.Sprintf("%s/config.json", actionPath))

	return action.NewAction(paramDef, configDef), nil
}

func (ap *ActionProvider) List() []string {

	items, err := ioutil.ReadDir(ap.sl.ActionFolderPath())
	if err != nil {
		return nil
	}

	dirList := make([]string, 0)
	for _, v := range items {
		dirList = append(dirList, v.Name())
	}

	return dirList

}

func (ap *ActionProvider) ActionExists(name string) bool {
	actionPath := ap.sl.ActionPath(name)
	if _, err := os.Stat(actionPath); os.IsNotExist(err) {
		return false
	}

	return true
}



type executor struct {
	ID          string
	execName    string
	handlerPath string
	sockPath    string
	configPath  string
}

func (ap ActionProvider) InvokeAction(name string, params map[string]interface{}, env []string) interface{} {

	pod := ap.newExecuter(name, "node10")



	dat, err := ioutil.ReadFile(pod.configPath)
	if err != nil {
		log.Printf("msg=action config does not exist,path=%v", pod.configPath)
	}
	var actionConfig map[string]interface{}

	err = json.Unmarshal(dat, &actionConfig)
	if err != nil {
		log.Printf("failed to read action config")
	}

	payload, _ := pod.EncodePayload(params)
	if err != nil {
		log.Printf("msg=action config does not exist,path=%v", pod.configPath)
	}

	//execName := pod.handlerPath
	args := []string{
		pod.handlerPath,
		fmt.Sprintf("%v", string(payload)),
	}
	cmd := exec.Command("node", args...)

	for k, v := range actionConfig {
		env = append(env, fmt.Sprintf("%v=%v", k, v))
	}
	log.Println(env)
	cmd.Env = env

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

	log.Printf("cmd=%s", cmd.Args)
	// TODO: Add Instrumentation
	if err := cmd.Start(); err != nil {
		log.Printf("failed to start cmd: %v\n", err)
	}

	var tmpResult map[string]interface{}

	outScanner := bufio.NewScanner(stdout)
	errScanner := bufio.NewScanner(stderr)

	go func() {
		for outScanner.Scan() {
			out := outScanner.Text()
			log.Printf("[%s] stdout: %v\n", name, out)
			if strings.Contains(out, fmt.Sprintf("%s", pod.ID)) {
				log.Println("value=", out)
				err := json.Unmarshal([]byte(out), &tmpResult)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()

	go func() {
		for errScanner.Scan() {
			log.Printf("[%s] stderr: %v\n", name, errScanner.Text())
		}
	}()


	log.Println("[ctrl]", "starting", "action", name)



	//cmd.Process.Kill()
	if err := cmd.Wait(); err != nil {
		log.Print(err)
	}
	log.Println("[ctrl]", "finished", "action", name, "result", tmpResult[pod.ID])
	return tmpResult[pod.ID]

}

func (ap ActionProvider) newExecuter(name string, stack string) *executor {
	id := ctrlID.GenULID()
	actionPath := ap.sl.ActionPath(name)

	return &executor{
		ID: id,
		// Stack related
		handlerPath: fmt.Sprintf("%s/action", actionPath),
		configPath:    fmt.Sprintf("%s/action/config.json", actionPath),
	}
}


func (ex *executor) EncodePayload(params map[string]interface{}) ([]byte, error) {
	invReq := make(map[string]interface{}, 0)


	invReq["params"] = params
	invReq["ctx"] = map[string]interface{}{
		"id": ex.ID,
	}

	return json.Marshal(invReq)
}
func (ex *executor) executeRPC(fd string, payload []byte) (interface{}, error) {

	c, err := connectToRPC(fd)

	var raw []byte
	err = c.Call("Action.Invoke", payload, &raw)
	if err != nil {
		log.Fatalln(err)
	}

	var result struct {
		ID      string
		Payload []byte
	}
	err = json.Unmarshal(raw, &result)
	if err != nil {
		log.Fatalln(err)
	}

	var final interface{}
	err = json.Unmarshal(result.Payload, &final)
	if err != nil {
		log.Fatalln(err)
	}

	return final, nil
}

func connectToRPC(fd string) (*rpc.Client, error) {
	var c *rpc.Client
	var err error
	retries := 10
	n := 0
	for n < retries {
		c, err = rpc.DialHTTP("unix", fd)
		if err == nil {
			break
		}
		time.Sleep(5 * time.Millisecond)
		n = n + 1
	}

	return c, err

}
