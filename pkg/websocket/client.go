package websocket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/wcatron/esk/pkg/message"
)

// Client structure
type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

func isMessageTypeNonstandard(messageType int) bool {
	return messageType != websocket.TextMessage && messageType != websocket.BinaryMessage
}

func handleError(c *Client, err error) {
	fmt.Print("client:error\n")
	log.Println(err)
}

func handleNonstandardMessage(c *Client, messageType int) {
	if messageType == websocket.PingMessage {
		c.Conn.WriteMessage(websocket.PongMessage, make([]byte, 0))
		return
	}
	fmt.Print("client:error:Only binary or text messages allowed.")
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			handleError(c, err)
			return
		}
		if isMessageTypeNonstandard(messageType) {
			handleNonstandardMessage(c, messageType)
			return
		}
		msg := message.MessageFromRaw(p)

		fmt.Printf("client:%s:received:%s\n", c.ID, message.CommandAsString(msg.Command))

		switch msg.Command {
		case message.CommandConnect:
			// TODO: Use ClientID sent by client.
			fmt.Printf("client:connect:%s\n", msg.ClientID)
			c.Connack()
			break
		case message.CommandPublish:
			c.Pool.Broadcast <- msg
			break
		case message.CommandSubscribe:
			c.Pool.Subscribe <- SubscriptionNotification{Topic: string(msg.Topic), Cursor: msg.Cursor, Client: c}
			break
		case message.CommandUnsubscribe:
			c.Pool.Unsubscribe <- SubscriptionNotification{Topic: string(msg.Topic), Client: c}
			break
		default:
			fmt.Printf("client:error:unknown_command:%d\n", msg.Command)
			break
		}
	}
}

// Send message to client
func (c *Client) Send(msg message.Message) {
	fmt.Printf("client:%s:send:%s\n", c.ID, message.CommandAsString(msg.Command))
	err := c.Conn.WriteMessage(websocket.BinaryMessage, msg.Raw)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Inform client of some payload
func (c *Client) Inform(topic string, cursor uint64, payload []byte) {
	msg := message.Message{
		Command: message.CommandInform,
		Topic:   []byte(topic),
		Payload: payload,
		Cursor:  cursor,
	}
	message.MessageWriteRaw(&msg)
	c.Send(msg)
}

// Suback : Acknowledge subscription to topic
func (c *Client) Suback(topic string) {
	msg := message.Message{
		Command: message.CommandSuback,
		Topic:   []byte(topic),
	}
	message.MessageWriteRaw(&msg)
	c.Send(msg)
}

// Unsuback : Acknowledge unsubscribe from topic
func (c *Client) Unsuback(topic string) {
	msg := message.Message{
		Command: message.CommandUnsuback,
		Topic:   []byte(topic),
	}
	message.MessageWriteRaw(&msg)
	c.Send(msg)
}

// Connack : Acknowledge connection of client
func (c *Client) Connack() {
	msg := message.Message{
		Command:  message.CommandConnack,
		ClientID: []byte(c.ID),
	}
	message.MessageWriteRaw(&msg)
	c.Send(msg)
}
