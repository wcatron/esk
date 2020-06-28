package websocket

import (
	"fmt"

	"github.com/rs/xid"
	"github.com/wcatron/esk/pkg/message"
)

type SubscriptionNotification struct {
	Client *Client
	Topic  string
	Cursor uint64
}

type WriteNotification struct {
	Pool    *Pool
	Message message.Message
}

type DataSource struct {
	Write chan WriteNotification
	Read  chan SubscriptionNotification
}

type Pool struct {
	Register         chan *Client
	Unregister       chan *Client
	Clients          map[*Client]bool
	Subscriptions    map[string]map[*Client]bool
	Subscribe        chan SubscriptionNotification
	Unsubscribe      chan SubscriptionNotification
	Broadcast        chan message.Message
	BroadcastWritten chan message.Message
	DataSource       *DataSource
}

func NewPool(datasource *DataSource) *Pool {
	return &Pool{
		Register:         make(chan *Client),
		Unregister:       make(chan *Client),
		Clients:          make(map[*Client]bool),
		Subscriptions:    make(map[string]map[*Client]bool),
		Subscribe:        make(chan SubscriptionNotification),
		Unsubscribe:      make(chan SubscriptionNotification),
		Broadcast:        make(chan message.Message),
		BroadcastWritten: make(chan message.Message),
		DataSource:       datasource,
	}
}

func suback(c *Client) {
	msg := message.Message{Command: message.CommandSuback}
	message.MessageWriteRaw(&msg)
	c.Send(msg)
}

func (pool *Pool) NextClientId() string {
	id := xid.New().String()
	for client := range pool.Clients {
		if client.ID == id {
			return pool.NextClientId()
		}
	}
	fmt.Printf("pool:NextClientId:%s\n", id)
	return id
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("pool:register:", len(pool.Clients))
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("pool:unregister:", len(pool.Clients))
			for _, clients := range pool.Subscriptions {
				delete(clients, client)
			}
			break
		case subscription := <-pool.Subscribe:
			fmt.Printf("pool:subscribe:%+v\n", subscription)
			clients := pool.Subscriptions[subscription.Topic]
			if clients == nil {
				pool.Subscriptions[subscription.Topic] = make(map[*Client]bool)
			}
			pool.Subscriptions[subscription.Topic][subscription.Client] = true
			subscription.Client.Suback(subscription.Topic)
			pool.DataSource.Read <- subscription
			break
		case subscription := <-pool.Unsubscribe:
			fmt.Printf("pool:unsubscribe:%+v\n", subscription)
			clients := pool.Subscriptions[subscription.Topic]
			if clients == nil {
				return
			}
			delete(pool.Subscriptions[subscription.Topic], subscription.Client)
			if len(pool.Subscriptions[subscription.Topic]) == 0 {
				delete(pool.Subscriptions, subscription.Topic)
			}
			subscription.Client.Unsuback(subscription.Topic)
			break
		case msg := <-pool.Broadcast:
			fmt.Printf("pool:broadcast\n")
			pool.DataSource.Write <- WriteNotification{
				Message: msg,
				Pool:    pool,
			}
			break
		case msg := <-pool.BroadcastWritten:
			fmt.Printf("pool:broadcastWritten\n")
			msg.Command = message.CommandInform
			message.MessageWriteRaw(&msg)
			topic := string(msg.Topic)
			clients := pool.Subscriptions[topic]
			if clients == nil {
				fmt.Printf("pool:broadcast:no clients for topic:%s\n", topic)
				break
			}
			for client := range clients {
				client.Send(msg)
			}
			break
		}
	}
}
