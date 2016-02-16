package main

import (
	"io"
	"os/exec"
	"strings"
)

type gitSource struct {
	Name   string
	Branch string
}

func (src gitSource) String() string {
	return src.Name + "/" + src.Branch
}

func gitFetch(src gitSource, stdout, stderr io.Writer) error {
	cmd := exec.Command("git", "fetch", src.Name, src.Branch)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func gitCheckOut(src gitSource, stdout, stderr io.Writer) error {
	cmd := exec.Command("git", "checkout", src.String())
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func gitRootPath() (rootPath string, err error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.CombinedOutput()
	rootPath = strings.Trim(string(out), "\n")
	return
}
