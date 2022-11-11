# Twitch Chat Frequency
Twitch is a livestreaming service with categories ranging from 'Game Shows' to 'Video Games'. Users can send chat messages, and read messages from others. This is a simple program (written in less than 100 lines in 30 minutes) to see the top K unique chat messages during a given polling period.

This repository is by no means a 'best practices' guide, but a quick demonstration of simple multi-threading. The command listener, chat callback and notify are all handled in separate Goroutines.

## TODO
- [x] Use sync.RWMutex for concurrent access
- [ ] Prettify this README

## How to Run
Download the repository using `git clone` or using the download button, then run:
```
go mod tidy
go build && ./TwitchChat # for macOS
```