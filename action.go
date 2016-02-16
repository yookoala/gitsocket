package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path"

	"github.com/codegangsta/cli"
)

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
