package main

import (
	"errors"
	"net"
	"net/http"
	"strings"
	"text/template"
)

const envTemplate = `# Generated via sheepit
{{range .}}export {{.}}
{{end}}
`

var (
	ErrorBuildFailed = errors.New("Failed to build")
)

type BuildConfig struct {
	Env          []string
	SSHKey       string
	HTTPEndpoint string
	BuildFile    string
	TargetImage  string
	Version      string
}

func BuildImage(cfg BuildConfig) error {
	var err error

	srv, err := startHTTPServer(cfg)

	if err != nil {
		return err
	}

	defer func() {
		srv.Close()
	}()

	err = buildStreamOutput(cfg)

	if err != nil {
		return err
	}

	SLog.Println("Saved", cfg.TargetImage, "version", cfg.Version)
	return err
}

func buildStreamOutput(cfg BuildConfig) error {
	cmdArgs := []string{
		"build",
		"--pull=true", // make sure we have the latest upstream image (ex: ruby_20)
		"--rm=false",  // make sure we don't trash the docker cache between builds
		"-f",
		cfg.BuildFile,
		"-t",
		strings.Join([]string{cfg.TargetImage, cfg.Version}, ":"),
		".",
	}

	return CommandStream("docker", cmdArgs...)
}

func startHTTPServer(cfg BuildConfig) (net.Listener, error) {
	tcpListener, err := net.Listen("tcp", cfg.HTTPEndpoint)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/ssh_key", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, cfg.SSHKey)
	})

	mux.HandleFunc("/env", func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.New("env").Parse(envTemplate)
		err = tmpl.Execute(w, cfg.Env)
		if err != nil {
			http.Error(w, "Template failed", http.StatusBadGateway)
		}
	})

	server := &http.Server{
		Handler: mux,
	}

	go server.Serve(tcpListener)
	return tcpListener, nil
}
