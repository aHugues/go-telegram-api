# Go Telegram API

This project is a Go API providing methods to interact with [Telegram](https://core.telegram.org/bots/api).
It provides methods to send messages, and listen for incoming messages from specific channels (i.e.
also specific users)

## Example

The most simple way to use this library is to create a bot that will send a message to a given user.

### Prerequise

It's supposed that you already have [created a bot](https://core.telegram.org/bots/features#botfather)
and you have a channel into which you want to send messages.

As such you have
- a `token`: a string, for instance `110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw`
- a `channelID`: an integer, for instance `814576567`

You need those two values to use your bot

### The program

A simple program that can use the bot can be as such

```go
package main

import (
	"context"

	"github.com/ahugues/go-telegram-api/bot"
)

func main() {
	bot := bot.New("110201543:AAHdqTcvCH1vGWJxfSeofSAs0K5PALDsaw")
	bot.SendMessage(context.Background(), 814576567, "Hello world!")
}
```

This will send a simple "Hello World!" message from your bot to the given channel.
