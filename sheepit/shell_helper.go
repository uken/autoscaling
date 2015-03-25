package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

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

func CommandProxy(bin string, args ...string) error {
	sigs := make(chan os.Signal, 1)
	var err error

	signal.Notify(sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGTTIN,
		syscall.SIGTTOU,
		syscall.SIGQUIT)

	cmd := prepareCommand(bin, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	cmd.Start()

	go func() {
		for {
			sig, ok := <-sigs
			if !ok {
				break
			}
			SLog.Println("Forwarding signal", sig, "to", cmd.Process.Pid)
			cmd.Process.Signal(sig)
		}
	}()

	err = cmd.Wait()
	return err
}

func prepareCommand(bin string, args ...string) *exec.Cmd {
	cmd := exec.Command(bin, args...)
	if os.Getenv("DEBUG") != "" {
		fmt.Println("running >>", bin, args)
	}
	return cmd
}

func currentGitRevision() string {
	rev, err := CommandOutput("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return ""
	}

	return rev
}
