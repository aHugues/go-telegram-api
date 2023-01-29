// package notifier provides a way to get events notifications
//
// Create a new notifier using `n := notifier.New(xxx)` using the login token for your bot.
// Then the notifier can start using `go n.Run()`
//
// Example
//
// Using the subscribe command for a list of event type returns a channel that will receive notifications
// of this type when received by the bot.
package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/ahugues/go-telegram-api/baseclt"
	"github.com/ahugues/go-telegram-api/servererror"
	"github.com/ahugues/go-telegram-api/structs"
	"github.com/google/uuid"
)

func Example() {
	n := New("test")
	id, rcvChan := n.Subscribe([]structs.UpdateType{structs.UpdateMessage})

	ctx, cancel := context.WithCancel(context.Background())
	go n.Run(ctx)

	cancel()
	<-rcvChan
	n.Unsubscribe(id)

}

type EventNotifier interface {
	Run(context.Context)
	Subscribe(eventType []structs.UpdateType) (uuid.UUID, <-chan structs.Update)
	Unsubscribe(uuid.UUID) error
	ErrChan() <-chan error
}

type updateAnswer struct {
	OK     bool             `json:"ok"`
	Result []structs.Update `json:"result"`
}

type subscriber struct {
	uuid    uuid.UUID
	events  []structs.UpdateType
	pubChan chan structs.Update
}

func (s *subscriber) subscribed(evenType structs.UpdateType) bool {
	for _, et := range s.events {
		if et == evenType {
			return true
		}
	}
	return false
}

type ConcreteUpdateNotifer struct {
	token         string
	lastOffset    int64
	httpClt       *http.Client
	ctx           context.Context
	mux           sync.Mutex
	cancel        context.CancelFunc
	stopChan      chan struct{}
	subscribers   map[uuid.UUID]*subscriber
	pollFrequency time.Duration
	errChan       chan error
}

func New(token string) *ConcreteUpdateNotifer {
	return &ConcreteUpdateNotifer{
		token:         token,
		httpClt:       http.DefaultClient,
		ctx:           context.TODO(),
		stopChan:      make(chan struct{}),
		subscribers:   make(map[uuid.UUID]*subscriber),
		pollFrequency: 1 * time.Second,
		errChan:       make(chan error),
	}
}

func (n *ConcreteUpdateNotifer) getUpdates() (updates []structs.Update, err error) {
	url := fmt.Sprintf("%s/bot%s/getUpdates?offset=%d", baseclt.BaseTelegramAPIURL, n.token, n.lastOffset+1)

	req, err := http.NewRequestWithContext(n.ctx, http.MethodGet, url, nil)
	if err != nil {
		return updates, fmt.Errorf("error building request: %w", err)
	}

	resp, err := n.httpClt.Do(req)
	if err != nil {
		return updates, fmt.Errorf("error sending request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return updates, servererror.FromResponse(resp)
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer([]byte{})
	maxRead := baseclt.HTTPMaxRead
	if resp.ContentLength != -1 {
		if resp.ContentLength > maxRead {
			return updates, fmt.Errorf("response too big (%d)", resp.ContentLength)
		}
		maxRead = resp.ContentLength
	}

	if _, err := io.CopyN(buf, resp.Body, maxRead); err != nil {
		return updates, fmt.Errorf("error reading body: %w", err)
	}

	var rawAnswer structs.Updates
	if err := json.Unmarshal(buf.Bytes(), &rawAnswer); err != nil {
		return updates, fmt.Errorf("error parsing body: %w", err)
	}
	if !rawAnswer.OK {
		return updates, servererror.FromResponse(resp)
	}
	return rawAnswer.Result, nil
}

func (n *ConcreteUpdateNotifer) handleUpdate(update structs.Update) {
	// TODO: maybe migrate this to a second indexed map updateType[uuid] to make search quicker
	// fmt.Printf("Handling update %+v\n", update)
	for _, sub := range n.subscribers {
		if sub.subscribed(update.Type()) {
			sub.pubChan <- update
		}
	}
	fmt.Printf("Updating offset from %d to %d\n", n.lastOffset, update.ID)
	n.lastOffset = update.ID
}

func (n *ConcreteUpdateNotifer) handleUpdates(updates []structs.Update) {
	n.mux.Lock()
	defer n.mux.Unlock()
	for _, u := range updates {
		n.handleUpdate(u)
	}
}

func (n *ConcreteUpdateNotifer) mainLoop() {
	defer n.cancel()
	for {
		select {
		case <-time.After(n.pollFrequency):
			// fmt.Println("Getting updates")
			updates, err := n.getUpdates()
			if err != nil {
				n.errChan <- fmt.Errorf("error getting updates: %w", err)
				continue
			}
			// fmt.Printf("Got updates %v\n", updates)
			n.handleUpdates(updates)
		case <-n.stopChan:
			return
		}
	}
}

func (n *ConcreteUpdateNotifer) Run(ctx context.Context) {
	n.ctx, n.cancel = context.WithCancel(ctx)
	go n.mainLoop()

	<-ctx.Done()
	n.stopChan <- struct{}{}
}

func (n *ConcreteUpdateNotifer) Subscribe(eventType []structs.UpdateType) (uuid.UUID, <-chan structs.Update) {
	n.mux.Lock()
	defer n.mux.Unlock()
	sub := subscriber{
		uuid:    uuid.New(),
		events:  eventType,
		pubChan: make(chan structs.Update),
	}
	n.subscribers[sub.uuid] = &sub
	return sub.uuid, sub.pubChan
}

func (n *ConcreteUpdateNotifer) Unsubscribe(id uuid.UUID) {
	n.mux.Lock()
	defer n.mux.Unlock()
	delete(n.subscribers, id)
}

func (n *ConcreteUpdateNotifer) ErrChan() <-chan error {
	return n.errChan
}
