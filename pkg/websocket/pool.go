package websocket

import (
	"fmt"
)

type SubscriptionNotification struct {
	Client *Client
	Topic  string
}

type DataSource struct {
	Write chan Message
	Read  chan *SubscriptionNotification
}

type Pool struct {
	Register      chan *Client
	Unregister    chan *Client
	Clients       map[*Client]bool
	Subscriptions map[string]map[*Client]bool
	Subscribe     chan SubscriptionNotification
	Unsubscribe   chan SubscriptionNotification
	Broadcast     chan Message
	DataSource    *DataSource
}

func NewPool(datasource *DataSource) *Pool {
	return &Pool{
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Clients:       make(map[*Client]bool),
		Subscriptions: make(map[string]map[*Client]bool),
		Subscribe:     make(chan SubscriptionNotification),
		Unsubscribe:   make(chan SubscriptionNotification),
		Broadcast:     make(chan Message),
		DataSource:    datasource,
	}
}

func suback(c *Client) {
	message := Message{Command: CommandSuback}
	MessageWriteRaw(&message)
	c.Publish(message)
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("pool:size:", len(pool.Clients))
			suback(client)
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("pool:size:", len(pool.Clients))
			for _, clients := range pool.Subscriptions {
				delete(clients, client)
			}
			break
		case subscription := <-pool.Subscribe:
			fmt.Printf("pool:subscribe:%+v\n", subscription)
			pool.DataSource.Read <- &subscription
			clients := pool.Subscriptions[subscription.Topic]
			if clients == nil {
				pool.Subscriptions[subscription.Topic] = make(map[*Client]bool)
			}
			pool.Subscriptions[subscription.Topic][subscription.Client] = true
			subscription.Client.Suback(subscription.Topic)
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
		case message := <-pool.Broadcast:
			fmt.Printf("pool:broadcast\n")
			pool.DataSource.Write <- message
			topic := string(message.Topic)
			clients := pool.Subscriptions[topic]
			if clients == nil {
				fmt.Printf("pool:broadcast:no clients for topic:%s\n", topic)
				break
			}
			for client := range clients {
				client.Publish(message)
			}
		}
	}
}
