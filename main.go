package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/gempir/go-twitch-irc/v3"
)

var (
	mutex       = sync.RWMutex{}
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
		mutex.Lock()
		defer mutex.Unlock()

		if isPolling {
			trimmed := strings.ToLower(strings.TrimSpace(message.Message))
			frequencies[trimmed] += 1
		}
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
			fmt.Println("Beginning polling...")

			mutex.Lock()
			frequencies = make(map[string]int) // Clear
			isPolling = true
			mutex.Unlock()
		case "res":
			fmt.Println("Calculating results...")

			var pairs []Pair

			mutex.RLock()
			isPolling = false

			for k, v := range frequencies {
				pairs = append(pairs, Pair{k, v})
			}

			mutex.RUnlock()

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
		default:
			fmt.Println("Unknown command.")
		}
	}
}
