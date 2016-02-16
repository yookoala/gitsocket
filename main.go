package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "githook"
	app.Usage = "helps update local git repository on triggers"
	app.Version = "0.3.0"
	app.Commands = []cli.Command{
		{
			Name: "server",
			Usage: "socket server. listen to unix socket and update " +
				"local git repository accordingly",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "remote, r",
					Value: "origin",
					Usage: "name of remote repository",
				},
				cli.StringFlag{
					Name:  "branch, b",
					Value: "master",
					Usage: "branch of remote repository",
				},
				cli.StringFlag{
					Name:  "listen, l",
					Value: "./githook.sock",
					Usage: "path to socket to listen for connection",
				},
			},
			Action: actionHook,
		},
		{
			Name: "setup",
			Usage: "help setting up the post-checkout hook in the current " +
				"repository folder. depends on vi",
			Action: actionSetup,
		},
	}

	app.Run(os.Args)
}
