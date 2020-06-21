package websocket

import (
	"fmt"
)

const (
	CommandConnect     = 1
	CommandConnack     = 2
	CommandPublish     = 3
	CommandSubscribe   = 8
	CommandSuback      = 9
	CommandUnsubscribe = 10
	CommandUnsuback    = 11
	CommandInform      = 103
)

type Message struct {
	Command      byte   `json:"command"`
	Topic        []byte `json:"topic"`
	TopicStart   int
	Payload      []byte `json:"payload"`
	PayloadStart int
	Cursor       int64
	Raw          []byte `json:"raw"`
}

func MessageFromRaw(raw []byte) Message {
	command := raw[0]
	dataLength := int(raw[1])
	fmt.Printf("client:command:%d\n", command)
	switch command {
	case CommandPublish:
		topicStartPosition := byte(header_bytes + header_bytes_publish)
		topicEndPosition := topicStartPosition + raw[2]
		topic := raw[3:topicEndPosition]
		payload := raw[topicEndPosition:dataLength]
		return Message{
			Command:      command,
			Topic:        topic,
			TopicStart:   int(topicStartPosition),
			Payload:      payload,
			PayloadStart: int(topicEndPosition),
			Raw:          raw}
	case CommandSubscribe:
		topicEndPosition := 3 + raw[2]
		topic := raw[3:topicEndPosition]
		return Message{Command: command, Topic: topic, Raw: raw}
	case CommandUnsubscribe:
		topicEndPosition := 3 + raw[2]
		topic := raw[3:topicEndPosition]
		return Message{Command: command, Topic: topic, Raw: raw}
	default:
		break
	}
	message := Message{Command: raw[0], Raw: raw}
	return message
}

// 0 - command
// 1 - size of rest of the message
const header_bytes = 2

// 0 - length of topic
const header_bytes_publish = 1

func MessageWriteRaw(message *Message) {
	switch message.Command {
	case CommandSuback:
		restOfMessageLength := 1 + len(message.Topic)
		raw := make([]byte, header_bytes+1)
		raw[0] = message.Command
		raw[1] = byte(restOfMessageLength)
		raw[2] = byte(len(message.Topic))
		message.Raw = append(raw, message.Topic...)
		return
	case CommandUnsuback:
		restOfMessageLength := 1 + len(message.Topic)
		raw := make([]byte, header_bytes+1)
		raw[0] = message.Command
		raw[1] = byte(restOfMessageLength)
		raw[2] = byte(len(message.Topic))
		message.Raw = append(raw, message.Topic...)
		return
	case CommandInform:
		restOfMessageLength := header_bytes_publish + len(message.Topic) + len(message.Payload)
		raw := make([]byte, header_bytes+header_bytes_publish)
		raw[0] = CommandPublish
		raw[1] = byte(restOfMessageLength)
		raw[2] = byte(len(message.Topic))
		message.Raw = append(append(raw, message.Topic...), message.Payload...)
		return
	default:
		raw := make([]byte, header_bytes)
		raw[0] = message.Command
		raw[1] = 0
		message.Raw = raw
	}
}
