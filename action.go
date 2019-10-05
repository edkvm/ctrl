package ctrl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/rpc"
	"os"
	"os/exec"
	"time"

	ctrlFS "github.com/edkvm/ctrl/fs"
	"github.com/oklog/ulid"
)

type ActionRepo struct {
	path string
}

func NewActionRepo() *ActionRepo {
	return &ActionRepo{
		ctrlFS.BuildActionRepoPath(),
	}
}

func (ar *ActionRepo) List() []string {

	items, err := ioutil.ReadDir(ar.path)
	if err != nil {
		return nil
	}

	dirList := make([]string, 0)
	for _, v := range items {
		dirList = append(dirList, v.Name())
	}

	return dirList

}

func (ar *ActionRepo) ActionExists(name string) bool {
	actionPath := ctrlFS.BuildActionPath(name)
	if _, err := os.Stat(actionPath); os.IsNotExist(err) {
		return false
	}

	return true
}

type Action struct {
	Name        string
	configPath  string
	paramsPath  string
	ctrlCh      chan struct{}
}

func NewAction(name string) *Action {
	actionPath := ctrlFS.BuildActionPath(name)
	return &Action{
		Name:     name,
		configPath:  fmt.Sprintf("%s/config.json", actionPath),
		paramsPath:  fmt.Sprintf("%s/params.json", actionPath),
	}
}

func genULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id, err := ulid.New(ulid.Timestamp(t), entropy)
	if err != nil {

	}

	return fmt.Sprintf("%s", id)
}

type executor struct {
	ID string
	execName    string
	handlerPath string
	sockPath    string
}

func newExecuter(name string, stack string) *executor {
	id := genULID()
	actionPath := ctrlFS.BuildActionPath(name)
	return &executor{
		ID: id,
		// Stack related
		handlerPath: fmt.Sprintf("%s/action", actionPath),
		sockPath:    fmt.Sprintf("%s/tmp/%s_%s.sock", actionPath, name, id),
	}
}

func (ar *ActionRepo) ExecuteAction(name string, payload []byte, env []string) interface{} {

	pod := newExecuter(name, "")

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

	log.Println("[ctrl]", "starting")
	result, err := pod.executeRPC(pod.sockPath, payload)


	cmd.Process.Kill()
	if err := cmd.Wait(); err != nil {
		log.Print(err)
	}

	return result

}

//
//import (
//	"fmt"
//	"log"
//	"net/rpc/jsonrpc"
//)
//
//const SockAddr = ":1234"//"eddy.sock"
//
//type Greeter struct {
//}
//
//func (g Greeter) Greet(name *string, reply *string) error {
//	*reply = "Hello, " + *name
//	return nil
//}
//
//func main() {
//
//	client, err := jsonrpc.Dial("tcp", SockAddr)
//	if err != nil {
//		log.Fatal("dialing:", err)
//	}
//	name := 4
//	log.Println("send")
//	var reply string
//	err = client.Call("add", &name, &reply)
//	log.Println("rcv")
//	if err != nil {
//		log.Fatal("arith error:", err)
//	}
//	fmt.Printf("Arith: %d=%d", name, reply)
//}
func (ex *executor) executeRPC(fd string , payload []byte) (interface{}, error) {

	c, err := connectToRPC(fd)

	var raw []byte
	err = c.Call("Action.Invoke", payload, &raw)
	if err != nil {
		log.Fatalln(err)
	}

	var result struct {
		ID string
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
		time.Sleep(5*time.Millisecond)
		n = n + 1
	}

	return c, nil

}

func (fr *Action) ParamsToJSON(args []string) map[string]interface{} {
	paramDef := ctrlFS.ReadFile(fr.paramsPath)

	vals := make([]interface{}, len(args))
	for i, _ := range args {
		vals[i] = args[i]
	}

	var params map[string]interface{}
	json.Unmarshal(paramDef, &params)

	idx := 0
	for k, _ := range params {
		params[k] = vals[idx]
		idx = idx + 1
	}

	return params
}

func (fr *Action) BuildPayload(params map[string]interface{}) ([]byte, []string) {
	configDef := ctrlFS.ReadFile(fr.configPath)

	var config map[string]interface{}
	json.Unmarshal(configDef, &config)

	invReq := make(map[string]interface{}, 0)

	encParams, _ := json.Marshal(params)
	invReq["payload"] = encParams
	invReq["ctx"] = ""


	buf, _ := json.Marshal(invReq)

	env := make([]string, 0)
	for k, v := range config {
		env = append(env, fmt.Sprintf("%v=%v", k, v))
	}

	return buf, env
}


