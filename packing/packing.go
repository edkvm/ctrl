package packing

import (
	"bytes"
	"fmt"
	"github.com/edkvm/ctrl/fs"
	"github.com/edkvm/ctrl/packing/stacks"
	"io"
	"log"
	"os"
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

type Pack struct {
	stack StackConfig
	name  string
	files map[string][]byte
}

func BuildPack(stackName, wd string) (*Pack, error) {
	// TODO: Add more error handeling
	dirs := strings.Split(wd, "/")
	if len(dirs) < 2 {
		// TODO: return error, name is not absolute
	}

	// Action name is the folder name
	funcName := dirs[len(dirs) - 1]

	pk := &Pack{
		stack: stacksList[stackName],
		name:  funcName,
		files: make(map[string][]byte, 3),
	}

	files, err := pk.stack.Build(wd)
	if err != nil {
		return nil, err
	}

	pk.files = files

	log.Println("built action:", funcName)
	return pk, nil
}

func (pk *Pack) Deploy() error {


	actionPath := fs.BuildActionPath(pk.name)

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
