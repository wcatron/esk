package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/wcatron/esk/pkg/message"
)

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

		pretty, _ := json.MarshalIndent(msg, "", "  ")
		fmt.Printf("client:%s:received:%s\n", c.ID, pretty)

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

func (c *Client) Send(message message.Message) {
	pretty, _ := json.MarshalIndent(message, "", "  ")
	fmt.Printf("client:%s:send:%s\n", c.ID, pretty)
	err := c.Conn.WriteMessage(websocket.BinaryMessage, message.Raw)
	if err != nil {
		fmt.Println(err)
		return
	}
}

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

func (c *Client) Suback(topic string) {
	msg := message.Message{
		Command: message.CommandSuback,
		Topic:   []byte(topic),
	}
	message.MessageWriteRaw(&msg)
	c.Send(msg)
}

func (c *Client) Unsuback(topic string) {
	msg := message.Message{
		Command: message.CommandUnsuback,
		Topic:   []byte(topic),
	}
	message.MessageWriteRaw(&msg)
	c.Send(msg)
}

func (c *Client) Connack() {
	msg := message.Message{
		Command:  message.CommandConnack,
		ClientID: []byte(c.ID),
	}
	message.MessageWriteRaw(&msg)
	c.Send(msg)
}
