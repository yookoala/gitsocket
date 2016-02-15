package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"

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

func gitRootPath() (rootPath string, err error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.CombinedOutput()
	rootPath = strings.Trim(string(out), "\n")
	return
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

func actionSetup(c *cli.Context) {
	rootPath, err := gitRootPath()
	if err != nil {
		log.Fatalf("error: %s", rootPath)
	}
	filename := path.Join(rootPath, ".git/hooks/post-checkout")

	// if file not exists, create the file
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Printf("file not exists")
		f, err := os.Create(filename)
		if err != nil {
			log.Fatalf("error: %s", err)
			return
		}
		fmt.Fprintf(f, "#!/bin/sh\n")
		fmt.Fprintf(f, "#\n")
		fmt.Fprintf(f, "# An example hook script to prepare a packed repository for use over\n")
		fmt.Fprintf(f, "# dumb transports.\n")
		fmt.Fprintf(f, "#\n")
		fmt.Fprintf(f, "# To enable this hook, rename this file to \"post-checkout\".\n")
		fmt.Fprintf(f, "exec echo \"checkout completed.\"\n")
		f.Close()

		os.Chmod(filename, 0777)
	}

	cmd := exec.Command("vi", filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatalf("error: %s", err)
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
