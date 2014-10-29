package main

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrorBuildFailed = errors.New("Failed to build")
)

type BuildConfig struct {
	Environment string
	BaseImage   string
	CacheDir    string
	BuildScript string
	TargetImage string
	Version     string
}

func BuildImage(cfg BuildConfig) error {
	var err error

	buildId, err := buildTargetId(cfg)

	if err != nil {
		return err
	}

	SLog.Println("Build log for container", shortBuildId(buildId))

	err = buildStreamOutput(buildId)

	if err != nil {
		return err
	}

	err = buildWaitTargetId(buildId)

	if err != nil {
		SLog.Println("Build failed")
		return err
	}

	SLog.Println("Build passed. Saving target image")

	targetImage := fmt.Sprintf("%s:%s", cfg.TargetImage, cfg.Version)
	targetId, err := buildCommit(buildId, targetImage)

	if err != nil {
		SLog.Println("Failed to commit target image")
		return err
	}

	SLog.Println("Saved", shortBuildId(targetId), "as", cfg.TargetImage, "version", cfg.Version)

	return err
}

func buildTargetId(cfg BuildConfig) (string, error) {
	appPath, _ := os.Getwd()

	cmdArgs := []string{
		"run",
		"-d",
		"--net",
		"host",
		"-e",
		fmt.Sprintf("DEPLOY_ENV=%s", cfg.Environment),
		"-v",
		fmt.Sprintf("%s:/build", appPath),
	}

	expandedCache := os.ExpandEnv(cfg.CacheDir)
	if expandedCache != "" {
		cmdArgs = append(cmdArgs, fmt.Sprintf("%s:/cache", expandedCache))
	}

	cmdArgs = append(cmdArgs, cfg.BaseImage, "build", fmt.Sprintf("/build/%s", cfg.BuildScript))

	return CommandOutput("docker", cmdArgs...)
}

func buildWaitTargetId(buildId string) error {
	cmdArgs := []string{
		"wait",
		buildId,
	}
	retCode, err := CommandOutput("docker", cmdArgs...)
	if err != nil {
		return err
	}

	if retCode != "0" {
		return ErrorBuildFailed
	}
	return nil
}

func buildStreamOutput(buildId string) error {
	cmdArgs := []string{
		"attach",
		buildId,
	}
	return CommandStream("docker", cmdArgs...)
}

func buildCommit(buildId string, target string) (string, error) {
	cmdArgs := []string{
		"commit",
		buildId,
		target,
	}
	return CommandOutput("docker", cmdArgs...)
}
