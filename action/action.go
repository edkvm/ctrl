package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	ctrlFS "github.com/edkvm/ctrl/fs"
)

var ErrMissingStats = errors.New("no stats for action")

type ActionRepo interface {
	FindAll() []*Action
}

type Codename string

type Code struct {
	Name        string
	Description string
	Author      string
	Keywords    string
	Version     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Action struct {
	CodeID      Codename
	Name        string
	createdAt   time.Time
	updatedAt   time.Time
	configPath  string
	paramsPath  string
}

type ActionParams map[string]interface{}

func NewAction(name string) *Action {
	actionPath := ctrlFS.BuildActionPath(name)
	return &Action{
		Name:       name,
		configPath: fmt.Sprintf("%s/config.json", actionPath),
		paramsPath: fmt.Sprintf("%s/params.json", actionPath),
	}
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
