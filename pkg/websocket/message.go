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
	ClientID     []byte `json:"clientId"`
	Raw          []byte `json:"raw"`
}

func MessageFromRaw(raw []byte) Message {
	command := raw[0]
	dataLength := int(raw[1])
	fmt.Printf("client:command:%d\n", command)
	switch command {
	case CommandConnect:
		return Message{
			Command:  command,
			ClientID: raw[2:dataLength],
			Raw:      raw,
		}
	case CommandPublish:
		topicStartPosition := byte(headerBytes + headerBytesPublish)
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
const headerBytes = 2

// 0 - length of topic
const headerBytesPublish = 1

func MessageWriteRaw(message *Message) {
	switch message.Command {
	case CommandConnack:
		restOfMessageLength := len(message.ClientID)
		raw := make([]byte, headerBytes)
		raw[0] = message.Command
		raw[1] = byte(restOfMessageLength)
		message.Raw = append(raw, message.ClientID...)
		return
	case CommandSuback:
		restOfMessageLength := 1 + len(message.Topic)
		raw := make([]byte, headerBytes+1)
		raw[0] = message.Command
		raw[1] = byte(restOfMessageLength)
		raw[2] = byte(len(message.Topic))
		message.Raw = append(raw, message.Topic...)
		return
	case CommandUnsuback:
		restOfMessageLength := 1 + len(message.Topic)
		raw := make([]byte, headerBytes+1)
		raw[0] = message.Command
		raw[1] = byte(restOfMessageLength)
		raw[2] = byte(len(message.Topic))
		message.Raw = append(raw, message.Topic...)
		return
	case CommandInform:
		restOfMessageLength := headerBytesPublish + len(message.Topic) + len(message.Payload)
		raw := make([]byte, headerBytes+headerBytesPublish)
		raw[0] = CommandPublish
		raw[1] = byte(restOfMessageLength)
		raw[2] = byte(len(message.Topic))
		message.Raw = append(append(raw, message.Topic...), message.Payload...)
		return
	default:
		raw := make([]byte, headerBytes)
		raw[0] = message.Command
		raw[1] = 0
		message.Raw = raw
	}
}
