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

			var ss []Pair

			mutex.RLock()
			isPolling = false

			for k, v := range frequencies {
				ss = append(ss, Pair{k, v})
			}

			mutex.RUnlock()

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
