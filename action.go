package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"regexp"

	"github.com/codegangsta/cli"
)

func handleShutdown(l net.Listener, pidfile string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for {
		select {
		case <-c:
			l.Close()
			if pidfile != "" {
				os.Remove(pidfile)
			}
			os.Exit(0)
		}
	}
}

type commandFn func(stdout, stderr io.Writer) error

func handleConnection(conn net.Conn, fn commandFn) {
	log.Printf("server: handleConnection")
	defer conn.Close()

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

		// TODO: allow overriding Stdout with log output
		w := io.MultiWriter(conn, os.Stdout)

		if err := fn(w, w); err == io.EOF {
			log.Printf("server: connection terminated")
			return
		} else if err != nil {
			log.Printf("callback error: %s", err.Error())
			return
		}
	}
}

// address returns networkk and address that fits
// the use of either net.Dial or net.Listen
func address(listen string) (network, address string) {
	reIP := regexp.MustCompile("^(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})\\:(\\d{2,5}$)")
	rePort := regexp.MustCompile("^(\\d+)$")
	switch {
	case reIP.MatchString(listen):
		network = "tcp"
		address = listen
	case rePort.MatchString(listen):
		network = "tcp"
		address = ":" + listen
	default:
		network = "unix"
		address = listen
	}
	return
}

func actionHook(c *cli.Context) {
	l, err := net.Listen(address(c.String("listen")))
	if err != nil {
		panic(err)
	}

	pidfile := c.String("pidfile")

	// cleanly disconnect the socket
	go handleShutdown(l, pidfile)

	if pidfile != "" {
		// get current pid and write to file
		pid := fmt.Sprintf("%d", os.Getpid())
		ioutil.WriteFile(pidfile, []byte(pid), 0600)
	}

	// define git source to update from
	src := gitSource{c.String("remote"), c.String("branch")}

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn, func(stdout, stderr io.Writer) error {
			if err := gitFetch(src, stdout, stderr); err != nil {
				return err
			}
			if err := gitCheckOut(src, stdout, stderr); err != nil {
				return err
			}
			return io.EOF
		})
	}
}

func actionClient(c *cli.Context) {
	conn, err := net.Dial(address(c.String("conn")))
	if err != nil {
		log.Fatalf("connection error: %s", err.Error())
		return
	}

	conn.Write([]byte("hello\n"))

	bufbytes := make([]byte, 1024)
	for {
		nr, err := conn.Read(bufbytes)

		// handle error
		if err == io.EOF {
			log.Printf("client: server connect closed")
			return
		} else if err != nil {
			log.Printf("client read error: %#v", err.Error())
			return
		}

		data := bufbytes[0:nr]
		fmt.Printf("%s", data)
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
