package main

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

// gitSource represents a git local repository
// with specified name of remote upstream and branch
type gitSource struct {
	Dir    string
	Name   string
	Branch string
}

func (src gitSource) String() string {
	return src.Name + "/" + src.Branch
}

func (src gitSource) Context(stdout, stderr io.Writer) *gitContext {
	return &gitContext{src, stdout, stderr}
}

// gitContext represents a context to run git command
// includes a gitSource and the output writers
type gitContext struct {
	Src    gitSource
	Stdout io.Writer
	Stderr io.Writer
}

func (c *gitContext) Command(gitcmd string, v ...string) error {
	cmdSlice := append([]string{gitcmd}, v...)
	fmt.Fprintf(c.Stdout, "git %s\n", strings.Join(cmdSlice, " "))
	cmd := exec.Command("git", cmdSlice...)
	cmd.Dir = c.Src.Dir
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	return cmd.Run()
}

func (c *gitContext) HardPull() error {
	if err := c.Command("fetch", c.Src.Name, c.Src.Branch); err != nil {
		return err
	}
	if err := c.Command("reset", "--hard", c.Src.String()); err != nil {
		return err
	}
	if err := c.Command("checkout"); err != nil {
		return err
	}
	return io.EOF
}

// gitRootPath obtains root path of a git repository
// with rev-parse command
func gitRootPath(dir string) (rootPath string, err error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	if dir != "" {
		cmd.Dir = dir
	}
	out, err := cmd.CombinedOutput()
	rootPath = strings.Trim(string(out), "\n")
	return
}

func mustGitRootPath(dir string) (rootPath string) {
	rootPath, err := gitRootPath(dir)
	if err != nil {
		log.Fatalf("error: %s", rootPath)
	}
	return
}
