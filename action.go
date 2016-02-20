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
	"syscall"
	"text/template"
	"time"

	"github.com/codegangsta/cli"
	godaemon "github.com/yookoala/go-daemon"
)

func handleShutdown(l net.Listener, pidfile string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL)

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

func handleConnection(conn net.Conn, src gitSource, stdout, stderr io.Writer) {
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

		rw := io.MultiWriter(conn, stdout)
		ew := io.MultiWriter(conn, stderr)

		if err := src.Context(rw, ew).HardPull(); err == io.EOF {
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

func actionServer(c *cli.Context) {

	var stdout io.Writer = os.Stdout
	var stderr io.Writer = os.Stderr

	if output := c.String("output"); output != "" {
		var f *os.File
		var err error
		if f, err = os.Create(output); err != nil {
			log.Fatalf("error opening output logfile %#v: %s",
				output, err.Error())
			return
		}
		stdout = f
		stderr = f
		log.SetOutput(f)
	} else if c.Bool("daemon") {
		var f *os.File
		var err error
		if f, err = os.Create(os.DevNull); err != nil {
			log.Fatalf("error opening output logfile %#v: %s",
				output, err.Error())
			return
		}
		stdout = f
		stderr = f
		log.SetOutput(f)
	}

	// daemonized server
	if c.Bool("daemon") {
		context := new(godaemon.Context)
		if child, _ := context.Reborn(); child != nil {

			// set timeout time
			timeout := time.After(time.Second * 30)

			// test if the socket is ready
			ready := make(chan int)
			go func() {
				for {
					conn, err := net.Dial(address(c.String("listen")))
					if err == nil {
						conn.Close() // close the test connection
						break
					}
				}
				ready <- 0
			}()

			// wait until timeout or socket ready by child
			select {
			case <-timeout:
				log.Fatalf("timeout: socket not ready in %d seconds", 30)
			case <-ready:
				return
			}
		}
		defer context.Release()
		actionServerMain(c, stdout, stderr)
		return
	}

	// normal server output
	actionServerMain(c, stdout, stderr)
}

func actionServerMain(c *cli.Context, stdout, stderr io.Writer) {

	l, err := net.Listen(address(c.String("listen")))
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	pidfile := c.String("pidfile")
	if pidfile != "" {
		// get current pid and write to file
		pid := fmt.Sprintf("%d", os.Getpid())
		ioutil.WriteFile(pidfile, []byte(pid), 0600)
	}

	// cleanly disconnect the socket
	go handleShutdown(l, pidfile)

	// define git source to update from
	src := gitSource{mustGitRootPath(c.String("gitrepo")),
		c.String("remote"), c.String("branch")}

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn, src, stdout, stderr)
	}
}

func actionOnce(c *cli.Context) {

	var stdout io.Writer = os.Stdout
	var stderr io.Writer = os.Stderr
	if output := c.String("output"); output != "" {
		var f *os.File
		var err error
		if f, err = os.Create(output); err != nil {
			log.Fatalf("error opening output logfile %#v: %s",
				output, err.Error())
			return
		}
		stdout = f
		stderr = f
		log.SetOutput(f)
	}

	// define git source to update from
	src := gitSource{mustGitRootPath(c.String("gitrepo")),
		c.String("remote"), c.String("branch")}

	if err := src.Context(stdout, stderr).HardPull(); err != io.EOF {
		log.Fatalf("error: %s", err.Error())
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

func createHookScript(filename, command string) (err error) {

	// template for git hook script
	tpl := template.Must(template.New("gitsocket").Parse(`#!/bin/sh
#
# An example hook script to prepare a packed repository for use over
# dumb transports.
#
# To enable this hook, rename this file to "post-checkout".
{{ .Command }}
`))

	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("error: %s", err)
		return
	}
	err = tpl.Execute(f, map[string]interface{}{
		"Command": command,
	})
	f.Close()
	if err != nil {
		return
	}

	err = os.Chmod(filename, 0777)
	return
}

func actionSetup(c *cli.Context) {

	// define git source to update from
	rootPath := mustGitRootPath(c.String("gitrepo"))
	filename := path.Join(rootPath, ".git/hooks/post-checkout")

	if command := c.String("command"); command != "" {
		// if file not exists, create the file
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			createHookScript(filename, command)
			return
		} else if c.Bool("force") {
			createHookScript(filename, command)
			return
		}
		fmt.Println("post-checkout script already exists. If you want to " +
			"overwrite, please use the -f flag")
		os.Exit(1)
		return
	}

	// if file not exists, create the file
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		createHookScript(filename, "exec echo \"checkout completed.\"\n")
	}

	cmd := exec.Command("vi", filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatalf("error: %s", err)
	}
}
