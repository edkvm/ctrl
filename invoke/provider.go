package invoke

import (
	"encoding/json"
	"fmt"
	"github.com/edkvm/ctrl"
	"github.com/edkvm/ctrl/action"

	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"time"

	ctrlID "github.com/edkvm/ctrl/pkg/id"
	ctrlFS "github.com/edkvm/ctrl/pkg/fs"
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

func (ap ActionProvider) EncodePayload(params map[string]interface{}) []byte {
	invReq := make(map[string]interface{}, 0)

	encParams, _ := json.Marshal(params)
	invReq["payload"] = encParams

	buf, _ := json.Marshal(invReq)

	return buf
}

type executor struct {
	ID string
	execName    string
	handlerPath string
	sockPath    string
}

func (ap ActionProvider) newExecuter(name string, stack string) *executor {
	id := ctrlID.GenULID()
	actionPath := ap.sl.ActionPath(name)
	return &executor{
		ID: id,
		// Stack related
		handlerPath: fmt.Sprintf("%s/action", actionPath),
		sockPath:    fmt.Sprintf("%s/tmp/%s_%s.sock", actionPath, name, id),
	}
}

func (ex *executor) executeRPC(fd string , payload []byte) (interface{}, error) {

	c, err := connectToRPC(fd)
	log.Println(ex)
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
