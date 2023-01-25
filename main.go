package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ahugues/go-telegram-api/bot"
	"github.com/ahugues/go-telegram-api/notifier"
	"github.com/ahugues/go-telegram-api/structs"
)

// func main() {
// 	platypus := bot.New(token)

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	updates, err := platypus.GetUpdates(ctx, 0)
// 	if err != nil {
// 		panic(err)
// 	}
// 	for _, update := range updates {
// 		fmt.Printf("update: %+v", update)
// 	}

// 	if err := platypus.SendMessage(ctx, channel, "plop"); err != nil {
// 		panic(err)
// 	}
// }

func main() {
	token := "redacted"
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	n := notifier.New(token)
	_, subChan1 := n.Subscribe([]structs.UpdateType{structs.UpdateMessage})
	_, subChan2 := n.Subscribe([]structs.UpdateType{structs.UpdateEditedMessage, structs.UpdateChannelPost})

	platypus := bot.New(token)

	fmt.Println("Testing bot")
	getMeCtx, getMeCancel := context.WithTimeout(context.Background(), 2*time.Second)
	if usr, err := platypus.GetMe(getMeCtx); err != nil {
		panic(err)
	} else {
		fmt.Printf("Successful login for bot %s [%s %s]\n", usr.Username, usr.FirstName, usr.LastName)
	}
	getMeCancel()

	go n.Run(ctx)

	for {
		select {
		case <-c:
			fmt.Println("Stopping listener")
			cancel()
			return
		case u := <-subChan1:
			fmt.Printf("Got update from listener 1: %s\n", u.Message.Text)
		case u := <-subChan2:
			fmt.Printf("Got update from listener 2: %+v\n", u)
		case err := <-n.ErrChan():
			fmt.Printf("Got an error: %s\n", err.Error())
		}
	}

}
