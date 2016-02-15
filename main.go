package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"

	"github.com/codegangsta/cli"
)

func init() {

}

type gitSource struct {
	Name   string
	Branch string
}

func (src gitSource) String() string {
	return src.Name + "/" + src.Branch
}

func gitFetch(src gitSource) error {
	cmd := exec.Command("git", "fetch", src.Name, src.Branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}

func gitCheckOut(src gitSource) error {
	cmd := exec.Command("git", "checkout", src.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	return cmd.Run()
}

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

func actionHook(c *cli.Context) {
	l, err := net.Listen("unix", c.String("socket"))
	if err != nil {
		panic(err)
	}

	// cleanly disconnect the socket
	go handleShutdown(l)

	// define git source to update from
	src := gitSource{c.String("remote"), c.String("branch")}

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn, func() error {
			if err := gitFetch(src); err != nil {
				return err
			}
			if err := gitCheckOut(src); err != nil {
				return err
			}
			return nil
		})
	}
}

func main() {

	app := cli.NewApp()
	app.Name = "githook"
	app.Usage = "helps update local git repository on triggers"
	app.Version = "0.1.0"
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
	}

	app.Run(os.Args)
}
