package main

import (
	_ "embed"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

// How often the checkin with channel indicating bot still up.
const checkinInterval = time.Duration(24 * time.Hour)

// How often to execute the job.
const execInterval = time.Duration(3 * time.Hour)

var channels []string

var (
	//go:embed checkin.txt
	checkinMessagesBody string

	checkinMessages = strings.Split(checkinMessagesBody, "\n")
)

func main() {
	rand.Seed(time.Now().Unix())

	apiKey := os.Getenv("SLACK_API_KEY")
	if len(apiKey) == 0 {
		log.Fatal("SLACK_API_KEY must be specified")
	}

	chStr := os.Getenv("SLACK_CHANNELS")
	if len(chStr) == 0 {
		log.Fatal("At least one Slack channel required via SLACK_CHANNELS")
	}

	if len(chStr) == 0 {
		log.Fatal("at least one Slack channel required via SLACK_CHANNELS")
	}
	channels = strings.Split(chStr, ",")
	log.Printf("Failures broadcast to channels: %v\n", channels)

	if len(os.Args) < 2 {
		log.Fatal("Path to executable must be provided as argument")
	} else if len(os.Args) > 2 {
		log.Println("WARN: executable args not currently supported")
	}
	name := os.Args[1]

	api := slack.New(apiKey)
	resp, err := api.AuthTest()
	if err != nil {
		log.Fatal("Unable to authenticate against Slack API")
	}
	log.Printf("Authenticated as %q\n", resp.User)

	// Setup a check in with Slack so you notice if the bot dissapears.
	// TODO: replace this with a Cloud deadman?
	checkin(api)
	go func() {
		for range time.NewTicker(checkinInterval).C {
			checkin(api)
		}
	}()

	execJob(api, name)
	go func() {
		for range time.NewTicker(execInterval).C {
			execJob(api, name)
		}
	}()

	select {}
}

var lastResult = true

func execJob(api *slack.Client, job string) {
	log.Printf("Running job: %q\n", job)
	start := time.Now()
	cmd := exec.Command(job)
	stdoutStderr, err := cmd.CombinedOutput()

	if err != nil {
		api.UploadFile(slack.FileUploadParameters{Content: string(stdoutStderr), Filetype: "text", Title: "Run Output", Channels: channels, InitialComment: ":apple: IOx build failed on M1. Please see attached output."})
	} else if !lastResult {
		jobGreen(api)
	}

	log.Printf("Job finished in : %v. Succeeded: %v \n", time.Since(start), err == nil)
	lastResult = (err == nil)
}

func checkin(api *slack.Client) {
	log.Printf("Bot checking in with channels: %v\n", channels)
	msg := checkinMessages[rand.Intn(len(checkinMessages))]
	sendMessage(api, fmt.Sprintf(":green_apple: %s", msg))
}

func jobGreen(api *slack.Client) {
	log.Println("Job status now good")
	sendMessage(api, ":white_tick: Last job completed successfully.")
}

func sendMessage(api *slack.Client, msg string) {
	for _, channel := range channels {
		if _, _, err := api.PostMessage(channel, slack.MsgOptionText(msg, false), slack.MsgOptionAsUser(true)); err != nil {
			log.Printf("ERROR: Message %q to channel %q failed: %v\n", msg, channel, err)
		}
	}
}
