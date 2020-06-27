package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
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
		message := MessageFromRaw(p)

		pretty, _ := json.MarshalIndent(message, "", "  ")
		fmt.Printf("client:%s:received:%s\n", c.ID, pretty)

		switch message.Command {
		case CommandConnect:
			// TODO: Use ClientID sent by client.
			fmt.Printf("client:connect:%s\n", message.ClientID)
			c.Connack()
			break
		case CommandPublish:
			c.Pool.Broadcast <- message
			break
		case CommandSubscribe:
			c.Pool.Subscribe <- SubscriptionNotification{Topic: string(message.Topic), Cursor: message.Cursor, Client: c}
			break
		case CommandUnsubscribe:
			c.Pool.Unsubscribe <- SubscriptionNotification{Topic: string(message.Topic), Client: c}
			break
		default:
			fmt.Printf("client:error:unknown_command:%d\n", message.Command)
			break
		}
	}
}

func (c *Client) Send(message Message) {
	pretty, _ := json.MarshalIndent(message, "", "  ")
	fmt.Printf("client:%s:send:%s\n", c.ID, pretty)
	err := c.Conn.WriteMessage(websocket.BinaryMessage, message.Raw)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (c *Client) Inform(topic string, cursor uint64, payload []byte) {
	message := Message{
		Command: CommandInform,
		Topic:   []byte(topic),
		Payload: payload,
		Cursor:  cursor,
	}
	MessageWriteRaw(&message)
	c.Send(message)
}

func (c *Client) Suback(topic string) {
	message := Message{
		Command: CommandSuback,
		Topic:   []byte(topic),
	}
	MessageWriteRaw(&message)
	c.Send(message)
}

func (c *Client) Unsuback(topic string) {
	message := Message{
		Command: CommandUnsuback,
		Topic:   []byte(topic),
	}
	MessageWriteRaw(&message)
	c.Send(message)
}

func (c *Client) Connack() {
	message := Message{
		Command:  CommandConnack,
		ClientID: []byte(c.ID),
	}
	MessageWriteRaw(&message)
	c.Send(message)
}
