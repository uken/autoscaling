package main

import (
	"fmt"

	"github.com/armon/consul-api"
)

type ReleaseConfig struct {
	TargetImage string
	Version     string
	Registry    string
	Consul      string
	Key         string
}

func ReleaseImage(cfg ReleaseConfig) error {
	SLog.Println("Tagging local image as latest")
	err := releaseTag(cfg.TargetImage, cfg.Version, "latest")
	if err != nil {
		return err
	}

	registryVersion := fmt.Sprintf("%s/%s:%s", cfg.Registry, cfg.TargetImage, cfg.Version)
	SLog.Println("Tagging and uploading", registryVersion)

	err = releaseTag(cfg.TargetImage, cfg.Version, registryVersion)

	if err != nil {
		return err
	}

	err = releaseUpload(registryVersion)

	if err != nil {
		return err
	}

	SLog.Println("Notifying fleet")

	err = notifyFleet(cfg.Consul, cfg.Key, registryVersion)

	return err
}

func releaseTag(image string, version string, tag string) error {
	cmdArgs := []string{
		"tag",
		fmt.Sprintf("%s:%s", image, version),
		tag,
	}

	return Command("docker", cmdArgs...)

}

func releaseUpload(registryUrl string) error {
	cmdArgs := []string{
		"push",
		registryUrl,
	}

	return CommandStream("docker", cmdArgs...)
}

func notifyFleet(consul string, key string, registryUrl string) error {
	conf := consulapi.DefaultConfig()
	conf.Address = consul

	client, err := consulapi.NewClient(conf)
	if err != nil {
		return err
	}
	kv := client.KV()

	pair := &consulapi.KVPair{Key: key, Value: []byte(registryUrl)}
	_, err = kv.Put(pair, nil)

	return err
}
