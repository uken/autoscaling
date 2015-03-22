package main

import (
	"errors"
	"strings"
)

var (
	ErrorImageNotFound = errors.New("Image Not Found")
)

func RunImage(image string, version string, args []string) error {
	cmdArgs := []string{
		"run",
		"-t",
		"-i",
		"--rm",
		"--net",
		"host",
	}

	dockerImage := strings.Join([]string{image, version}, ":")

	cmdArgs = append(cmdArgs, dockerImage, "run")
	cmdArgs = append(cmdArgs, args...)

	return CommandStream("docker", cmdArgs...)
}
