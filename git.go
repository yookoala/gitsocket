package main

import (
	"os"
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
