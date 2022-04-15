package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/pflag"

	"github.com/alllomancer/k8s-slackbot/pkg/kubernetes"
	"github.com/alllomancer/k8s-slackbot/pkg/slack"
	"github.com/alllomancer/k8s-slackbot/server/options"
)

// SlackBotServer defines struct of slack bot server
type SlackBotServer struct {
	healthzPort    int
	configMap      string
	botToken       string
	kubeConfigFile string
	debugEnable    bool
	slackBot       slack.SlackBot
}

// NewSlackBotServerDefault returns new slack bot server
func NewSlackBotServerDefault(config *options.SlackBotServerConfig) *SlackBotServer {
	s := SlackBotServer{
		healthzPort:    config.HealthzPort,
		configMap:      config.ConfigMap,
		botToken:       config.BotToken,
		kubeConfigFile: config.KubeConfigFile,
		debugEnable:    config.DebugEnable,
	}
	_, err := kubernetes.RunGetNodes(s.kubeConfigFile)
	if err != nil {
		log.Fatal("Connect to kubernetes failed: %s\n", err)
	} else {
		log.Printf("Connect to kubernetes successful\n")
	}
	s.slackBot = slack.NewSlackBot(s.botToken)
	return &s
}

// Run starts server and health check
func (server *SlackBotServer) Run() {
	pflag.VisitAll(func(flag *pflag.Flag) {
		log.Printf("FLAG: --%s=%q", flag.Name, flag.Value)
	})
	server.setupHealthzHandlers()
	log.Printf("Setting up Healthz Handler(/readiness, /cache) on port :%d", server.healthzPort)

	server.start()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", server.healthzPort), nil))
}

// setupHealthzHandlers sets up a readiness and liveness endpoint for kube2sky.
func (server *SlackBotServer) setupHealthzHandlers() {
	http.HandleFunc("/readiness", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "ok\n")
	})
	http.HandleFunc("/cache", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "ok\n")
	})
}

// start starts server
func (server *SlackBotServer) start() {
	go server.slackBot.RunSlackRTMServer(server.kubeConfigFile)
}
