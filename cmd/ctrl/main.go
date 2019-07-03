package main

import (
	"fmt"
	"github.com/edkvm/ctrl"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
)

type mainConfig struct {
	wdPath     string
}


func main() {
	//currentPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	wdPath, _ := os.Getwd()
	app := cli.NewApp()

	app.Name = "ctrl"


	app.Commands = []cli.Command{
		{
			Name: "deploy",
			Usage: "copy function to the runner, if the function does not exist it will be created",
			Action: func(c *cli.Context) error {

				dirs := strings.Split(wdPath, "/")

				if len(dirs) < 2 {
					// TODO: return error, path is not absolute
				}

				// Action name is the folder name
				funcName := dirs[len(dirs) - 1]
				log.Println("[ctrl] Deploying Action: ", funcName)

				// Read action
				srcPath := fmt.Sprintf("%s/action.js", wdPath)

				packAction := ctrl.NewPack(funcName, "node8.9")

				// Copy handler
				srcFd, err := os.Open(srcPath)
				if err != nil {
					return err
				}
				defer srcFd.Close()

				err = packAction.Deploy(srcFd)
				if err != nil {
					log.Println(err)
				}

				return nil
			},

		},
		{
			Name: "init",
			Usage: "no usage yet",
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
				c.String()
				funcName := c.Args().First()
				log.Println("[ctrl] Running Action: ", funcName)

				args := c.Args()[1:]

				fr := ctrl.NewAction(funcName)

				if !fr.IsExists() {
					log.Fatal(fmt.Sprintf("function %s does not exists", funcName))
				}

				// Parse params
				input := fmt.Sprintf(`{ "__action_config": {
 "apiKey": "ba9f19affad8428980ee3a66462295a9", "unitSys": "metric" },"data": { "city": "%s" } }`, args[0])

				log.Println("input: ", input)
				result := fr.Execute(input)

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