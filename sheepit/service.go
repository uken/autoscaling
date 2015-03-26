package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type ServiceDescription struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
	Port      int    `yaml:"port"`
	Service   string `yaml:"service"`
	Image     string `yaml:"image"`
}

func (sd ServiceDescription) NodeCurrent() string {
	return sd.Namespace + "/nodes/" + hostName() + "/current"
}

func (sd ServiceDescription) NodeDeploy() string {
	return sd.Namespace + "/nodes/" + hostName() + "/deploy"
}

func (sd ServiceDescription) AppCurrent() string {
	return sd.Namespace + "/current"
}

// we cannot use yaml to write the template directly as it would escape tpl vars
func (sd ServiceDescription) SaveTo(w io.Writer) error {
	bw := bufio.NewWriter(w)
	_, err := bw.WriteString(fmt.Sprintf("name: %s\nnamespace: %s\nport: %d\nservice: %s\nimage: %s\n",
		sd.Name,
		sd.Namespace,
		sd.Port,
		sd.Service,
		sd.Image))

	if err != nil {
		return err
	}

	return bw.Flush()
}

type ServiceConfig struct {
	Config string
	Consul string
	ServiceDescription
}

func RunService(cfg ServiceConfig) error {
	buf, err := ioutil.ReadFile(cfg.Config)
	err = yaml.Unmarshal(buf, &cfg.ServiceDescription)
	err = fetchImage(cfg.Image)
	if err != nil {
		return err
	}

	err = killPreviousDockerName(cfg.Name)
	if err != nil {
		return err
	}

	err = setupNewDockerName(cfg)

	return notifyService(cfg)
}

func setupNewDockerName(cfg ServiceConfig) error {
	cmdArgs := []string{
		"create",
		"--net=host",
		"--name=" + cfg.Name,
	}

	if cfg.Port > 0 {
		cmdArgs = append(cmdArgs, "-e", fmt.Sprintf("PORT=%d", cfg.Port))
	}

	cmdArgs = append(cmdArgs, cfg.Image, fmt.Sprintf("forego start %s", cfg.Service))

	return CommandProxy("docker", cmdArgs...)
}

func killPreviousDockerName(name string) error {
	inspectArgs := []string{
		"inspect",
		name,
	}

	err := Command("docker", inspectArgs...)
	// in this case, the container doesn't exist
	if err != nil {
		return nil
	}

	stopArgs := []string{
		"stop",
		name,
	}

	err = Command("docker", stopArgs...)
	if err != nil {
		return nil
	}

	rmArgs := []string{
		"rm",
		name,
	}

	return Command("docker", rmArgs...)
}

func fetchImage(img string) error {
	cmdArgs := []string{
		"pull",
		img,
	}

	return CommandStream("docker", cmdArgs...)
}

func notifyService(cfg ServiceConfig) error {
	return updateKey(cfg.Consul, cfg.NodeCurrent(), cfg.Image)
}

func hostName() string {
	h, _ := os.Hostname()
	return h
}
