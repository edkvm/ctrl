package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
)


const (
	servicePathDefult = "/usr/local/var/pico"
	functionsPath = "functions"
	tmpPath = "tmp"
	templatePath = "assembler/templates"
)

type mainConfig struct {
	wdPath     string
	stacksPath string
	stacksList map[string]string
}

func (mc mainConfig) stackPath(name string) string {
	val := mc.stacksList[name]
	return fmt.Sprintf("%s/%s", mc.stacksPath, val)
}

func main() {
	//currentPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	wd, _ := os.Getwd()
	app := cli.NewApp()

	app.Name = "pico"
	envConf := mainConfig{
		wdPath:     wd,
		stacksPath: fmt.Sprintf("%s/assembler/templates", servicePathDefult),
		stacksList: map[string]string{"node8.9": "node8.9/index.js.tmpl"},
	}

	app.Commands = []cli.Command{
		{
			Name: "deploy",
			Usage: "copy function to the runner, if the function does not exist it will be created",
			Action: func(c *cli.Context) error {

				stackNameDefault := "node8.9"
				dirs := strings.Split(envConf.wdPath, "/")

				if len(dirs) < 2 {
					// TODO: return error, path is not absolute
				}
				// TODO: check fo root


				funcName := dirs[len(dirs) - 1]
				log.Println("Func name", funcName)
				srcPath := fmt.Sprintf("%s/handler.js", envConf.wdPath)


				tmpl, err := template.ParseFiles(envConf.stackPath(stackNameDefault))
				if err != nil {
					return fmt.Errorf("template for stack %s was not found, %s", stackNameDefault, err)
				}
				err = deployFuncLocal(srcPath, funcPath(funcName), tmpl)
				if err != nil {
					log.Println(err)
				}

				return nil
			},

		},
		{
			Name: "init",
			Usage: "not usage yet",
			Action: func(c *cli.Context) error {
				log.Println("init not implemented")
				return nil
			},
		},
		{
			Name: "list",
		},
		{
			Name: "run",
			Usage: "runs the specified function",
			Action: func(c *cli.Context) error{

				funcName := c.Args().First()
				log.Printf("args: %s", funcName)
				fr := newFuncRunner(funcName, funcPath(funcName))

				if _, err := os.Stat(fr.handlerPath); os.IsNotExist(err) {
					log.Fatal(fmt.Sprintf("function %s does not exists", funcName))
				}

				result := fr.execute(fmt.Sprintf("socket message Id (%s)", genULID()))

				log.Printf("[INFO] %v", result)

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}