package ctrl

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/oklog/ulid"

	_ "github.com/edkvm/ctrl/statik"
)


type Pack struct {
	stack StackConfig
	name  string
	files  map[string][]byte}

func BuildPack(stackName, wd string) (*Pack, error) {
	// TODO: Add more error handeling
	dirs := strings.Split(wd, "/")
	if len(dirs) < 2 {
		// TODO: return error, name is not absolute
	}

	// Action name is the folder name
	funcName := dirs[len(dirs) - 1]

	log.Println("building function:", funcName)

	pk := &Pack{
		stack: stacksList[stackName],
		name: funcName,
		files: make(map[string][]byte, 3),
	}

	pk.build(wd)

	return pk, nil
}

func (pk *Pack) build(wd string) error {
	// Read action
	for i := 0; i < len(pk.stack.fileNames); i++ {
		fileName := pk.stack.fileNames[i]
		pk.files[fileName] = readFile(fmt.Sprintf("%s/%s", wd, fileName))
	}

	return nil
}


func (pk *Pack) Deploy() error {


	actionPath := buildActionPath(pk.name)

	// Create tmp folder
	if _, err := os.Stat(actionPath); os.IsNotExist(err) {
		err := os.MkdirAll(fmt.Sprintf("%s/tmp", actionPath), os.ModePerm)
		if err != nil {
			return err
		}
	}


	err := pk.stack.writeEntryPoint(actionPath)
	if err != nil {
		return nil
	}

	// Write Action files
	for name, _ := range pk.files {
		// TODO move to function
		file := pk.files[name]
		srcReader := bytes.NewReader(file)

		dstFd, err := os.Create(fmt.Sprintf("%s/%s", actionPath, name))
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


func deployLocal() {

}

func updateFunction() {

}

func updateFunctionResources() {

}

func createFunctionSchedual() {

}

func pauseFunction() {

}

func deleteFunction() {

}

func genULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id, err := ulid.New(ulid.Timestamp(t), entropy)
	if err != nil {

	}

	return fmt.Sprintf("%s", id)
}
