package main

import (
	"fmt"
	"github.com/edkvm/ctrl"
	"log"
	"os"
	"path/filepath"


	"github.com/urfave/cli"
)


func main() {
	//currentPath, err := filepath.Abs(filepath.Dir(os.Args[0]))

	app := cli.NewApp()

	app.Name = "ctrl"

	app.Commands = []cli.Command{
		{
			Name: "deploy",
			Usage: "copy function to the runner, if the function does not exist it will be created",
			Action: func(c *cli.Context) error {
				wdPath, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}

				if c.NArg() > 0 {
					pathArg := c.Args().First()
					wdPath, err = filepath.Abs(pathArg)
					if err != nil {
						log.Fatal(err)
					}
				}

				log.Println("deploying function from: ", wdPath)

				pk, err := ctrl.BuildPack("node_v10", wdPath)
				if err != nil {
					return err
				}

				err = pk.Deploy()
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

				funcName := c.Args().First()
				log.Println("running action: ", funcName)

				args := c.Args()[1:]

				fr := ctrl.NewAction(funcName)

				if !fr.IsExists() {
					log.Fatal(fmt.Sprintf("function does not exists: %s", funcName))
				}

				// Parse params
				result := fr.Execute([]string(args))

				log.Println(result)

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}