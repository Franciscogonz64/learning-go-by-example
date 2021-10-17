package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token        string
	CleanContent string
	ImageName    string
)

const KuteGoAPIURL = "http://localhost:8080"

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

type Gopher struct {
	Name string `json: "name"`
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	CleanContent = CleanData(m.Content)
	ImageName = CleanContent + ".png"

	checkForResponse := func() {

		// Call the KuteGo API and retrieve Gopher images
		response, err := http.Get(KuteGoAPIURL + "/gopher/" + CleanContent)
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			_, err = s.ChannelFileSend(m.ChannelID, ImageName, response.Body)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Sprintf("Error: Can't get %s Gopher! :-(", CleanContent)
		}
	}

	switch m.Content {
	case "!5th-element", "!LICENSE", "!arrow-gopher", "!back-to-the-future-v2", "!baywatch", "!big-bang-theory", "!bike-gopher",
		"!blues-gophers", "!buffy-the-gopher-slayer", "!chandleur-gopher", "!cherry-gopher", "!devnation-france-gopher",
		"!dr-who", "!fire-gopher", "!firefly-gopher", "!fort-boyard", "!friends", "!gandalf-colored", "!gandalf", "!gladiator-gopher",
		"!gopher-dead", "!gopher-open", "!gopher-speaker", "!gopher", "!halloween-spider", "!happy-gopher", "!harry-gopher", "!idea-gopher",
		"!indiana-jones", "!jedi-gopher", "!jurassic-park", "!love-gopher", "!luigi-gopher", "!mac-gopher", "!marshal", "!men-in-black-v2",
		"!mojito-gopher", "!paris-gopher", "!sandcastle-gopher", "!saved-by-the-bell", "!star-wars", "!stargate", "!tadx-gopher", "!urgences",
		"!vampire-xmas", "!wired-gopher", "!x-files", "!yoda-gopher":

		checkForResponse()

	}

	if m.Content == "!random" {

		//Call the KuteGo API and retrieve a random Gopher
		response, err := http.Get(KuteGoAPIURL + "/gopher/random/")
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			_, err = s.ChannelFileSend(m.ChannelID, "random-gopher.png", response.Body)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error: Can't get random Gopher! :-(")
		}
	}

	if m.Content == "!gophers" {

		//Call the KuteGo API and display the list of available Gophers
		response, err := http.Get(KuteGoAPIURL + "/gophers/")
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			// Transform our response to a []byte
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
			}

			// Put only needed informations of the JSON document in our array of Gopher
			var data []Gopher
			err = json.Unmarshal(body, &data)
			if err != nil {
				fmt.Println(err)
			}

			// Create a string with all of the Gopher's name and a blank line as separator
			var gophers strings.Builder
			for _, gopher := range data {
				gophers.WriteString(gopher.Name + "\n")
			}

			// Send a text message with the list of Gophers
			_, err = s.ChannelMessageSend(m.ChannelID, gophers.String())
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Error: Can't get list of Gophers! :-(")
		}
	}
}

func CleanData(resp string) string {

	cleanString := strings.TrimPrefix(resp, "!")
	fmt.Printf("%s -> %s\n", resp, cleanString)

	return cleanString
}
