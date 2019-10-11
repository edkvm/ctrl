package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	ctrlFS "github.com/edkvm/ctrl/fs"
)


var ErrMissingStats = errors.New("no stats for action")

type RegID string

type LocalePath string

type ActionRepo interface {
	FindAll() []*Action
}

type Action struct {
	ID        RegID
	Name      string
	Path      LocalePath
	CreatedAt time.Time
	UpdatedAt time.Time
	configPath string
	paramsPath string
}

func NewAction(name string) *Action {
	actionPath := ctrlFS.BuildActionPath(name)
	return &Action{
		Name:     name,
		configPath:  fmt.Sprintf("%s/config.json", actionPath),
		paramsPath:  fmt.Sprintf("%s/params.json", actionPath),
	}
}

type Stat struct {
	ID          string
	ActionID 	RegID
	Duration    float32
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

func (fr *Action) BuildEnv() []string {
	configDef := ctrlFS.ReadFile(fr.configPath)

	var config map[string]interface{}
	json.Unmarshal(configDef, &config)

	env := make([]string, 0)
	for k, v := range config {
		env = append(env, fmt.Sprintf("%v=%v", k, v))
	}

	return env
}

func (fr *Action) EncodePayload(params map[string]interface{}) []byte {
	invReq := make(map[string]interface{}, 0)

	encParams, _ := json.Marshal(params)
	invReq["payload"] = encParams

	buf, _ := json.Marshal(invReq)

	return buf
}