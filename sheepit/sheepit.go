package main

import (
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v1"
)

func main() {
	var err error

	curDir, _ := os.Getwd()

	var (
		app = kingpin.New("sheepit", "Docker helper")

		build            = app.Command("build", "Prepare a docker image for release.")
		buildAppDir      = build.Flag("dir", "App source directory").Default(curDir).ExistingDir()
		buildCacheDir    = build.Flag("cache", "Cache Directory on Host").ExistingDir()
		buildDeployEnv   = build.Flag("env", "Environment to setup").Default("staging").String()
		buildBaseImage   = build.Flag("base", "Base Docker image").Default("uken/frontend_19:latest").String()
		buildTargetImage = build.Flag("target", "Target Image").Required().String()
		buildVersion     = build.Flag("version", "Build version").Default("0.0.1").String()
		buildScript      = build.Arg("script", "Build script").Required().String()

		run        = app.Command("run", "Run commands on a specific build.")
		runImage   = run.Flag("image", "Pre-built image").Required().String()
		runVersion = run.Flag("version", "Build version").Default("0.0.1").String()
		runCmd     = run.Arg("cmd", "Command").Default("/bin/bash").Strings()

		release               = app.Command("release", "Upload and release build.")
		releaseImage          = release.Flag("image", "Pre-built image").Required().String()
		releaseVersion        = release.Flag("version", "Build version").Default("0.0.1").String()
		releaseDockerRegistry = release.Flag("registry", "Docker Registry").Required().String()
		releaseConsul         = release.Flag("consul", "Consul Address").Required().String()
		releaseKey            = release.Flag("key", "Consul Key").Required().String()
	)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case build.FullCommand():

		cfg := BuildConfig{
			AppDir:      *buildAppDir,
			CacheDir:    *buildCacheDir,
			Environment: *buildDeployEnv,
			BaseImage:   *buildBaseImage,
			Version:     *buildVersion,
			TargetImage: *buildTargetImage,
			BuildScript: *buildScript,
		}
		err = BuildImage(cfg)
	case run.FullCommand():
		err = RunImage(*runImage, *runVersion, *runCmd)
	case release.FullCommand():
		cfg := ReleaseConfig{
			TargetImage: *releaseImage,
			Version:     *releaseVersion,
			Registry:    *releaseDockerRegistry,
			Consul:      *releaseConsul,
			Key:         *releaseKey,
		}
		err = ReleaseImage(cfg)
	default:
		app.Usage(os.Stdout)
	}

	if err != nil {
		SLog.Println(err)
		os.Exit(1)
	}
}
