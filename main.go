package main

import (
	"log"

	"github.com/spf13/pflag"

	"github.com/alllomancer/k8s-slackbot/server"
	"github.com/alllomancer/k8s-slackbot/server/options"
)

var (
	// value overwritten during build. This can be used to resolve issues.
	version = "1.0"
	gitRepo = "https://github.com/alllomancer/k8s-slackbot"
)

func main() {
	config := options.NewSlackBotServerConfig()
	config.AddFlags(pflag.CommandLine)

	pflag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("Using build: %v - %v", gitRepo, version)

	s := server.NewSlackBotServerDefault(config)
	s.Run()
}
