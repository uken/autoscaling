package main

import (
	"errors"
	"fmt"
	"strings"
	"time"

	consulapi "github.com/armon/consul-api"
)

var (
	ErrorHostFailedToBoot = errors.New("Host failed to boot up")
	ErrorKVNotSet         = errors.New("KV not set")
	ErrorKVUpdate         = errors.New("Failed to update KV")
	HostCheckDelay        = 5 * time.Second
	MaxHostRetries        = 30
)

type ReleaseConfig struct {
	TargetImage   string
	Version       string
	Registry      string
	Consul        string
	KeyNamespace  string
	RollingDeploy bool
}

func (cfg ReleaseConfig) registryTargetImage() string {
	return fmt.Sprintf("%s/%s:%s", cfg.Registry, cfg.TargetImage, cfg.Version)
}

func ReleaseImage(cfg ReleaseConfig) error {
	var err error

	SLog.Println("Tagging", cfg.registryTargetImage())

	err = releaseTag(cfg.TargetImage, cfg.Version, cfg.registryTargetImage())

	if err != nil {
		return err
	}

	err = releaseTag(cfg.TargetImage, "latest", cfg.registryTargetImage())

	if err != nil {
		return err
	}

	SLog.Println("Uploading", cfg.registryTargetImage())
	err = releaseUpload(cfg.registryTargetImage())

	if err != nil {
		return err
	}

	if cfg.RollingDeploy {
		err = rollingDeploy(cfg)

		if err != nil {
			return err
		}
	}

	SLog.Println("Updating master version")
	err = updateKey(cfg.Consul, fmt.Sprintf("%s/current", cfg.KeyNamespace), cfg.registryTargetImage())

	return err
}

func rollingDeploy(cfg ReleaseConfig) error {

	SLog.Println("Gathering active hosts")
	hosts, err := collectHosts(cfg.Consul, cfg.KeyNamespace)

	if err != nil {
		return err
	}

	SLog.Println("Starting rolling deploy on", hosts)

	// these keys are already namespaced
	for _, hostKey := range hosts {
		hostParts := strings.Split(hostKey, "/")
		host := hostParts[len(hostParts)-1]
		SLog.Println("Notifying", host, cfg.registryTargetImage())

		err = updateKey(cfg.Consul, fmt.Sprintf("%s/deploy", hostKey), cfg.registryTargetImage())
		if err != nil {
			SLog.Println("Failed to notify", host)
			return ErrorKVUpdate
		}

		retries := 0
		SLog.Println("Waiting for", host, "to boot up")
		for {
			time.Sleep(HostCheckDelay)

			nodeVer, err := getKey(cfg.Consul, fmt.Sprintf("%s/current", hostKey))
			if err != nil {
				SLog.Println("Failed to acquire node version, will keep trying")
				continue
			}

			if nodeVer == cfg.registryTargetImage() {
				SLog.Println(host, "is live")
				break
			}

			retries++

			if retries > MaxHostRetries {
				SLog.Println("Host is not coming back up. Aborting.")
				deleteKey(cfg.Consul, fmt.Sprintf("%s/deploy", hostKey))
				return ErrorHostFailedToBoot
			}
		}
	}

	return nil
}

func releaseTag(image string, version string, tag string) error {
	// we FORCE tag, which is not nice but allows devs to re-spin builds
	cmdArgs := []string{
		"tag",
		"-f",
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

	return Command("docker", cmdArgs...)
}

func collectHosts(consul string, key string) ([]string, error) {
	ret := []string{}
	conf := consulapi.DefaultConfig()
	conf.Address = consul

	client, err := consulapi.NewClient(conf)
	if err != nil {
		return ret, err
	}
	kv := client.KV()

	nodeKey := fmt.Sprintf("%s/nodes/", key)

	keys, _, err := kv.Keys(nodeKey, "/", nil)

	for _, k := range keys {
		ret = append(ret, strings.TrimSuffix(k, "/"))
	}

	return ret, nil
}

func updateKey(consul string, key string, val string) error {
	conf := consulapi.DefaultConfig()
	conf.Address = consul

	client, err := consulapi.NewClient(conf)
	if err != nil {
		return err
	}
	kv := client.KV()

	pair := &consulapi.KVPair{Key: key, Value: []byte(val)}
	_, err = kv.Put(pair, nil)

	return err
}

func deleteKey(consul string, key string) error {
	conf := consulapi.DefaultConfig()
	conf.Address = consul

	client, err := consulapi.NewClient(conf)
	if err != nil {
		return err
	}
	kv := client.KV()

	_, err = kv.Delete(key, nil)
	return err
}

func getKey(consul, key string) (string, error) {
	conf := consulapi.DefaultConfig()
	conf.Address = consul

	client, err := consulapi.NewClient(conf)
	if err != nil {
		return "", err
	}
	kv := client.KV()

	pair, _, err := kv.Get(key, nil)

	if err != nil {
		return "", err
	}

	if pair == nil {
		return "", ErrorKVNotSet
	}

	return string(pair.Value), nil
}
