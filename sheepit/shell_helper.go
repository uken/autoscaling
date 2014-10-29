package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func shortBuildId(fullId string) string {
	return fullId[0:8]
}

func CommandOutput(bin string, args ...string) (string, error) {
	cmd := prepareCommand(bin, args...)
	rawOut, err := cmd.Output()
	if err != nil {
		return "", err
	}
	out := strings.TrimRight(string(rawOut), "\n")
	return out, err
}

func Command(bin string, args ...string) error {
	cmd := prepareCommand(bin, args...)
	return cmd.Run()
}

func CommandStream(bin string, args ...string) error {
	cmd := prepareCommand(bin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	return err
}

func prepareCommand(bin string, args ...string) *exec.Cmd {
	cmd := exec.Command(bin, args...)
	if os.Getenv("DEBUG") != "" {
		fmt.Println("running >>", bin, args)
	}
	return cmd
}
