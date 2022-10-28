package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gempir/go-twitch-irc/v3"
)

var (
	isPolling   = false
	frequencies = map[string]int{}
	resultLimit = 20
)

type Pair struct {
	Key   string
	Value int
}

func main() {
	client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if isPolling {
			if val, ok := frequencies[message.Message]; ok {
				frequencies[message.Message] = val + 1
			} else {
				frequencies[message.Message] = 1
			}
		}

		//fmt.Println(message.Message)
	})

	client.Join("xQc")

	go func() {
		commandHandler(client)
	}()

	err := client.Connect()
	if err != nil {
		panic(err)
	}
}

func commandHandler(client *twitch.Client) {
	var temp string
	for {
		_, err := fmt.Scanln(&temp)
		if err != nil {
			panic(err)
		}

		command := strings.ToLower(temp)

		switch command {
		case "exit":
			client.Disconnect()
		case "poll":
			frequencies = make(map[string]int) // Clear
			isPolling = true
			fmt.Println("Beginning polling...")
		case "res":
			fmt.Println("Calculating results...")
			isPolling = false

			var ss []Pair

			for k, v := range frequencies {
				ss = append(ss, Pair{k, v})
			}

			sort.Slice(ss, func(i, j int) bool {
				return ss[i].Value > ss[j].Value
			})

			var count = 0

			fmt.Println("Out of", len(frequencies), "unique messages, top", resultLimit, "results shown")
			for _, kv := range ss {
				fmt.Printf("[%d] %s\n", kv.Value, kv.Key)
				count += 1
				if count == resultLimit {
					break
				}
			}
		default:
			fmt.Println("Unknown command.")
		}
	}
}
