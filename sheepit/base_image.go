package main

import (
	"fmt"
	"strings"
)

type BaseConfig struct {
	BaseImage   string
	BaseVersion string
	TargetImage string
}

func BaseImage(cfg BaseConfig) error {
	var err error

	base := fmt.Sprintf("%s:%s", cfg.BaseImage, cfg.BaseVersion)
	target := fmt.Sprintf("%s:base", cfg.TargetImage)

	hasLocalImage := hasLocal(cfg.TargetImage, "base")
	if hasLocalImage {
		SLog.Println("Base image", base, "already present. docker untag it first")
		return nil
	}

	SLog.Println("Downloading ", base)
	err = fetchImage(base)
	if err != nil {
		SLog.Println("Failed to download from public repo. This should be ok for private images")
	}

	SLog.Println("Tagging", cfg.BaseImage, "as", target)
	return releaseTag(cfg.BaseImage, cfg.BaseVersion, target)
}

// naive check
func hasLocal(image string, version string) bool {
	cmdArgs := []string{
		"images",
		image,
	}

	out, err := CommandOutput("docker", cmdArgs...)
	if err != nil {
		return false
	}

	return strings.Contains(out, version)
}
