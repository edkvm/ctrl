package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var ErrMissingStats = errors.New("no stats for action")

type ActionRepo interface {
	FindAll() []*Action
}

type BlueprintID string

type Blueprint struct {
	Name        string
	Description string
	Author      string
	Keywords    string
	Version     string
	Stack       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Action struct {
	BlueprintID BlueprintID
	Name        string
	createdAt   time.Time
	updatedAt   time.Time
	configDef   []byte
	paramsDef   []byte
}

type Params map[string]interface{}

func NewAction(paramDef []byte, configDef []byte) *Action {
	return &Action{
		paramsDef: paramDef,
		configDef: configDef,
	}
}

func (a *Action) ParamsToJSON(args []string) map[string]interface{} {

	vals := make([]interface{}, len(args))
	for i, _ := range args {
		vals[i] = args[i]
	}

	var params map[string]interface{}
	json.Unmarshal(a.paramsDef, &params)

	idx := 0
	for k, _ := range params {
		params[k] = vals[idx]
		idx = idx + 1
	}

	return params
}

func (a *Action) BuildEnv() []string {
	var config map[string]interface{}
	json.Unmarshal(a.configDef, &config)

	env := make([]string, 0)
	for k, v := range config {
		env = append(env, fmt.Sprintf("%v=%v", k, v))
	}

	return env
}
