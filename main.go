package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/codegangsta/cli"
)

func handleShutdown(l net.Listener) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for {
		select {
		case <-c:
			l.Close()
			os.Exit(0)
		}
	}
}

func handleConnection(conn net.Conn, fn func() error) {
	log.Printf("server: handleConnection")
	for {
		bufbytes := make([]byte, 1024)
		nr, err := conn.Read(bufbytes)

		// handle error
		if err == io.EOF {
			log.Printf("server: client connect closed")
			return
		} else if err != nil {
			log.Printf("server read error: %#v", err.Error())
			return
		}

		data := bufbytes[0:nr]
		fmt.Fprintf(conn, "echo: ")
		conn.Write(data)
		log.Printf("server got: %s", data)

		if err := fn(); err != nil {
			log.Printf("callback error: %s", err.Error())
		}
	}
}

func main() {

	app := cli.NewApp()
	app.Name = "githook"
	app.Usage = "helps update local git repository on triggers"
	app.Version = "0.2.0"
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
					Name:  "socket, s",
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
