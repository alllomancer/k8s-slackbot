package slack

import (
	"fmt"
	"github.com/alllomancer/k8s-slackbot/pkg/kubernetes"
	"github.com/oriser/regroup"
	"log"
	"os"
	"strings"
	"github.com/slack-go/slack"
)

// SlackBot defines a slack client
type SlackBot struct {
	Client *slack.Client
}

type service struct {
	Name    string
	Age     string
	Version string
}

var services []service

func NewSlackBot(token string) SlackBot {
	return SlackBot{
		Client: slack.New(
			token,
			slack.OptionDebug(true),
			slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
		),
	}
}

// RunSlackRTMServer runs rtm server
func (bot SlackBot) RunSlackRTMServer(kubeconfig string) {
	rtm := bot.Client.NewRTM()
	go rtm.ManageConnection()
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {

		case *slack.MessageEvent:
			input := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")
			if len(input) > 0 {
				command := input[0]
				args := input[1:]
				log.Printf("command: %v, args: %v", command, args)
				switch command {
				case "list":

					result, err := kubernetes.RunGet(kubeconfig)
					if err != nil {
						rtm.SendMessage(rtm.NewOutgoingMessage(err.Error(), ev.Msg.Channel))
					}
					var re = regroup.MustCompile(`(?m)[a-z].* (?P<name>[a-z].*\S) .*\d/\d.* (?P<age>\d.*[d])`)
					for _, line := range strings.Split(strings.TrimSuffix(result, "\n"), "\n") {

						matches, err := re.Groups(line)
						if err == nil {
							version, err := kubernetes.RunExec(kubeconfig, matches["name"])
							if err != nil {
								rtm.SendMessage(rtm.NewOutgoingMessage(err.Error(), ev.Msg.Channel))
							}
							services = append(services, service{Name: matches["name"], Age: matches["age"], Version: version})
						}

					}
					if err != nil {
						rtm.SendMessage(rtm.NewOutgoingMessage(err.Error(), ev.Msg.Channel))
					} else {
						rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("%+v", services), ev.Msg.Channel))
					}

				case "logs":
					if len(args) == 2 {

						output, err := kubernetes.RunLogs(kubeconfig, args[0], args[1])
						if err != nil {
							rtm.SendMessage(rtm.NewOutgoingMessage(err.Error(), ev.Msg.Channel))
						} else {
							rtm.SendMessage(rtm.NewOutgoingMessage(output, ev.Msg.Channel))
						}
					}
					rtm.SendMessage(rtm.NewOutgoingMessage("wrong number of args, command is `logs [pod] [tail limit]`", ev.Msg.Channel))
				default:
					rtm.SendMessage(rtm.NewOutgoingMessage("allowed commands are `list` and `logs [pod] [tail limit]`", ev.Msg.Channel))
				}

			}

		case *slack.InvalidAuthEvent:
			log.Printf("Invalid credentials")
			return

		case *slack.RTMError:
			log.Printf("Error: %s\n", ev.Error())
		default:
		}
	}
}
