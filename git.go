package main

import (
	"io"
	"log"
	"os/exec"
	"strings"
)

type gitSource struct {
	Dir    string
	Name   string
	Branch string
}

func (src gitSource) String() string {
	return src.Name + "/" + src.Branch
}

// gitActionsFor returns commandFn to run all the
// gitsocket relevant git command for the given source src
func gitActionsFor(src gitSource) commandFn {
	return func(stdout, stderr io.Writer) (err error) {
		if err := gitFetch(src, stdout, stderr); err != nil {
			return err
		}
		if err := gitCheckOut(src, stdout, stderr); err != nil {
			return err
		}
		return io.EOF
	}
}

func gitFetch(src gitSource, stdout, stderr io.Writer) error {
	cmd := exec.Command("git", "fetch", src.Name, src.Branch)
	cmd.Dir = src.Dir
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func gitCheckOut(src gitSource, stdout, stderr io.Writer) error {
	cmd := exec.Command("git", "checkout", src.String())
	cmd.Dir = src.Dir
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

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
