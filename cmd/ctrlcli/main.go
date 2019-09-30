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
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "stack, s",
					Value: "node10",
				},
			},
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

				pk, err := ctrl.BuildPack("node10", wdPath)
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
			Name: "ls",
			Usage: "list all the available actions",
			Action: func(c *cli.Context) error {
				ar := ctrl.NewActionRepo()

				fmt.Println(ar.List())
				return nil
			},
		},
		{
			Name: "run",
			Usage: "runs the specified function",
			Action: func(c *cli.Context) error{

				actionName := c.Args().First()


				args := c.Args()[1:]
				ar := ctrl.NewActionRepo()

				if !ar.ActionExists(actionName) {
					log.Fatal(fmt.Sprintf("action does not exists: %s", actionName))
				}

				fr := ctrl.NewAction(actionName)

				log.Println("running action:", actionName)
				log.Println("%v", []byte(args.First()))
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