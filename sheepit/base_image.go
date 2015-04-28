package main

import (
	"fmt"
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

	SLog.Println("Downloading ", base)
	err = fetchImage(base)
	if err != nil {
		SLog.Println("Failed to download from public repo. This should be ok for private images")
	}

	SLog.Println("Tagging", cfg.BaseImage, "as", target)
	return releaseTag(cfg.BaseImage, cfg.BaseVersion, target)
}
