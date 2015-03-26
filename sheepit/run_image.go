package main

import (
	"errors"
	"strings"
)

var (
	ErrorImageNotFound = errors.New("Image Not Found")
	RunDefaultCommand  = "/bin/bash"
)

func RunImage(image string, version string, args []string) error {
	cmdArgs := []string{
		"run",
		"--rm",
		"--net",
		"host",
	}

	// TODO figure out when to start interactive session
	if args[0] == RunDefaultCommand {
		cmdArgs = append(cmdArgs, "-t", "-i")
	}

	dockerImage := strings.Join([]string{image, version}, ":")
	cmdArgs = append(cmdArgs, dockerImage)
	cmdArgs = append(cmdArgs, args...)

	return CommandStream("docker", cmdArgs...)
}
