package packing

import (
	"bytes"
	"fmt"
	"github.com/edkvm/ctrl"
	"github.com/edkvm/ctrl/packing/stacks"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type StackConfig interface {
	Deploy(path string) error
	Build(wd string) (map[string][]byte, error)
}

var stacksList = map[string]StackConfig{
	"node10": stacks.NewNodev10(),
	"go": stacks.NewGoV1(),
}



type ActionPack struct {
	sl *ctrl.ServiceLoc
}

func NewActionPack(sl *ctrl.ServiceLoc) *ActionPack {
	return &ActionPack{
		sl: sl,
	}
}

func (ap *ActionPack) Create(name string) error {
	dir := fmt.Sprintf("%s/%s.git", ap.sl.GitPath(), name)

	if st, err := os.Stat(dir); err == nil && st.IsDir() {
		return fmt.Errorf("function already exists")
	}

	// Create Bare git repo
	args := []string{"init", "--bare", fmt.Sprintf("%s.git", name) }
	command := exec.Command("/usr/bin/git", args...)
	command.Dir = ap.sl.GitPath()

	err := command.Run()

	return err
}

func (ap *ActionPack) Pull() {

}

func (ap *ActionPack) Pack() {

}


type Pack struct {
	stack      StackConfig
	actionName string
	files      map[string][]byte
	sl ctrl.ServiceLoc
}

func BuildPack(stackName, wd string) (*Pack, error) {
	// TODO: Add more error handeling
	dirs := strings.Split(wd, "/")
	if len(dirs) < 2 {
		// TODO: return error, actionName is not absolute
	}

	// Action actionName is the folder actionName
	actionName := dirs[len(dirs) - 1]

	pk := &Pack{
		stack:      stacksList[stackName],
		actionName: actionName,
		files:      make(map[string][]byte, 3),
	}

	files, err := pk.stack.Build(wd)
	if err != nil {
		return nil, err
	}

	pk.files = files

	log.Println("built action:", actionName)
	return pk, nil
}

func (pk *Pack) Deploy() error {
	actionPath := pk.sl.ActionPath(pk.actionName)

	// Create tmp folder
	if _, err := os.Stat(actionPath); os.IsNotExist(err) {
		err := os.MkdirAll(fmt.Sprintf("%s/tmp", actionPath), os.ModePerm)
		if err != nil {
			return err
		}
	}


	err := pk.stack.Deploy(actionPath)
	if err != nil {
		return nil
	}

	// Write Action files
	for name, _ := range pk.files {
		// TODO move to function
		file := pk.files[name]
		srcReader := bytes.NewReader(file)

		dstFd, err := os.OpenFile(fmt.Sprintf("%s/%s", actionPath, name),os.O_RDWR|os.O_CREATE|os.O_TRUNC,os.ModePerm)
		if err != nil {
			return err
		}
		defer dstFd.Close()

		_, err = io.Copy(dstFd, srcReader)
		if err != nil {
			// TODO: File didn't open, Report as (SystemError)
			return err
		}

	}

	return nil
}
