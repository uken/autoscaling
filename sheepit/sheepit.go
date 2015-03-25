package main

import (
	"errors"
	"os"
	"strings"

	kingpin "gopkg.in/alecthomas/kingpin.v1"
)

type EnvKV struct {
	KV map[string]string
}

func (e *EnvKV) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	if len(parts) != 2 {
		return errors.New("Invalid environment variable")
	}

	e.KV[parts[0]] = parts[1]
	return nil
}

func (e *EnvKV) String() string {
	return ""
}

func BuildEnv(s kingpin.Settings) *EnvKV {
	target := &EnvKV{
		KV: make(map[string]string),
	}
	s.SetValue((*EnvKV)(target))
	return target
}

func main() {
	var err error

	curDir, _ := os.Getwd()

	var (
		app = kingpin.New("sheepit", "Deploy helper")

		build             = app.Command("build", "Prepare a docker image for release.")
		buildAppDir       = build.Flag("dir", "App source directory").Default(curDir).ExistingDir()
		buildTargetImage  = build.Flag("target", "Target Image").Required().String()
		buildVersion      = build.Flag("version", "Build version").Default(currentGitRevision()).String()
		buildHttpEndpoint = build.Flag("http", "HTTP Endpoint").Default("127.0.0.1:9090").String()
		buildSSHKey       = build.Flag("sshkey", "SSH Key").Default("/dev/null").ExistingFile()
		buildFile         = build.Flag("buildfile", "Dockerfile").Short('f').Default("Dockerfile").ExistingFile()
		buildEnv          = BuildEnv(build.Flag("env", "Set image environment variable").Short('e').PlaceHolder("ENV=VALUE"))

		run        = app.Command("run", "Run commands on a specific build.")
		runImage   = run.Flag("image", "Pre-built image").Short('i').Required().String()
		runVersion = run.Flag("version", "Build version").Default(currentGitRevision()).String()
		runCmd     = run.Arg("cmd", "Command").Default("/bin/bash").Strings()

		release               = app.Command("release", "Upload and release build.")
		releaseImage          = release.Flag("image", "Pre-built image").Short('i').Required().String()
		releaseVersion        = release.Flag("version", "Build version").Default(currentGitRevision()).String()
		releaseDockerRegistry = release.Flag("registry", "Docker Registry").Required().String()
		releaseConsul         = release.Flag("consul", "Consul Address").Default("127.0.0.1:8500").String()
		releaseKey            = release.Flag("namespace", "Consul Key Namespace").Short('n').Required().String()
		releaseRolling        = release.Flag("rolling", "Rolling deploy").Bool()

		setup          = app.Command("setup", "Setup consul template")
		setupConfig    = setup.Flag("config", "Configuration file (yml)").Short('f').Required().String()
		setupTemplate  = setup.Flag("template", "Consul template destination").Required().String()
		setupConsul    = setup.Flag("consul", "Consul Address").Default("127.0.0.1:8500").String()
		setupNamespace = setup.Flag("namespace", "Consul Key Namespace").Short('n').Required().String()
		setupPort      = setup.Flag("port", "setup Port").Short('p').Int()
		setupName      = setup.Flag("name", "Docker Container Name").Required().String()
		setupWorker    = setup.Arg("worker", "Procfile entry").Required().String()

		service       = app.Command("service", "Setup and start local docker container")
		serviceConfig = service.Flag("config", "Configuration file (yml)").Short('f').Required().ExistingFile()
		serviceConsul = setup.Flag("consul", "Consul Address").Default("127.0.0.1:8500").String()
	)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case build.FullCommand():
		err := os.Chdir(*buildAppDir)
		if err != nil {
			break
		}

		env := []string{}
		for k, v := range buildEnv.KV {
			env = append(env, k+"="+v)
		}

		cfg := BuildConfig{
			Env:          env,
			Version:      *buildVersion,
			BuildFile:    *buildFile,
			HTTPEndpoint: *buildHttpEndpoint,
			TargetImage:  *buildTargetImage,
			SSHKey:       *buildSSHKey,
		}
		err = BuildImage(cfg)
	case run.FullCommand():
		err = RunImage(*runImage, *runVersion, *runCmd)
	case release.FullCommand():
		cfg := ReleaseConfig{
			TargetImage:   *releaseImage,
			Version:       *releaseVersion,
			Registry:      *releaseDockerRegistry,
			Consul:        *releaseConsul,
			KeyNamespace:  *releaseKey,
			RollingDeploy: *releaseRolling,
		}
		err = ReleaseImage(cfg)
	case setup.FullCommand():
		cfg := SetupConfig{
			Config:   *setupConfig,
			Consul:   *setupConsul,
			Template: *setupTemplate,
			ServiceDescription: ServiceDescription{
				Name:      *setupName,
				Port:      *setupPort,
				Namespace: *setupNamespace,
				Service:   *setupWorker,
			},
		}
		err = SetupService(cfg)
	case service.FullCommand():
		cfg := ServiceConfig{
			Config: *serviceConfig,
			Consul: *serviceConsul,
		}
		err = RunService(cfg)
	default:
		app.Usage(os.Stdout)
	}

	if err != nil {
		SLog.Println(err)
		os.Exit(1)
	}
}
