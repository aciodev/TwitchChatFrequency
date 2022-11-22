package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/gempir/go-twitch-irc/v3"
)

var (
	twitchClient  = twitch.NewAnonymousClient()
	twitchChannel = ""
	mutex         = sync.RWMutex{}
	isPolling     = false
	frequencies   = map[string]int{}
	resultLimit   = 20
	regex, _      = regexp.Compile("\\s\\s+")
)

type Pair struct {
	Key   string
	Value int
}

func main() {
	twitchClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		mutex.Lock()
		defer mutex.Unlock()

		if isPolling {
			trimmedInnerSpaces := regex.ReplaceAllString(message.Message, " ")
			trimmedNewLinesAndSpaces := strings.TrimSpace(trimmedInnerSpaces)
			lowercase := strings.ToLower(trimmedNewLinesAndSpaces)
			frequencies[lowercase] += 1
		}
	})

	go func() {
		commandHandler()
	}()

	_ = twitchClient.Connect()
	fmt.Println("Thank you for using the app!")
}

func commandHandler() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		commandString, _ := reader.ReadString('\n')
		commandString = strings.Replace(commandString, "\n", "", -1)
		commandArgs := strings.Split(commandString, " ")

		switch strings.ToLower(commandArgs[0]) {
		case "exit":
			_ = twitchClient.Disconnect()
			fmt.Println("Stopping the application.")
			return
		case "join":
			if len(commandArgs) != 2 {
				fmt.Println("Usage: join <username>")
				continue
			}

			leaveExistingTwitchChannel()
			joinTwitchChannel(commandArgs[1])
		case "leave":
			leaveExistingTwitchChannel()
		case "poll":
			beginPolling()
		case "res":
			printResults()
		default:
			printHelp()
		}
	}
}

func printHelp() {
	fmt.Println("Unknown command. Try the following:")
	fmt.Println("join <channel> - to join a Twitch chat")
	fmt.Println("leave          - to leave a Twitch chat")
	fmt.Println("poll           - to begin collecting chat messages (resets collected data)")
	fmt.Println("res            - to print the top K frequent chat messages")
	fmt.Println("exit           - to close this app")
}

func validateInChat() bool {
	if len(twitchChannel) == 0 {
		fmt.Println("You must first join a Twitch Channel using: join <username>")
		return false
	}

	return true
}

func leaveExistingTwitchChannel() {
	if len(twitchChannel) != 0 {
		fmt.Println("Leaving existing Twitch channel: " + twitchChannel)
		twitchClient.Depart(twitchChannel)
	}
}

func joinTwitchChannel(channel string) {
	twitchChannel = channel
	fmt.Println("Joined Twitch channel: " + twitchChannel)
	twitchClient.Join(twitchChannel)
}

func beginPolling() {
	if !validateInChat() {
		return
	}

	fmt.Println("Monitoring " + twitchChannel + "'s chat...")
	mutex.Lock()
	frequencies = make(map[string]int) // Clear
	isPolling = true
	mutex.Unlock()
}

func printResults() {
	if !validateInChat() {
		return
	}

	fmt.Println("Calculating results for " + twitchChannel + "'s chat...")

	var pairs []Pair

	// -- read lock the frequencies --
	mutex.RLock()

	for k, v := range frequencies {
		pairs = append(pairs, Pair{k, v})
	}

	mutex.RUnlock()

	// -- pairs available, end concurrent access --

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})

	var count = 0

	fmt.Println("Out of", len(pairs), "unique messages, top", resultLimit, "results shown")
	for _, kv := range pairs {
		fmt.Printf("[%d] %s\n", kv.Value, kv.Key)
		count += 1
		if count == resultLimit {
			break
		}
	}
}
