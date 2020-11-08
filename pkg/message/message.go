package message

import (
	"encoding/binary"
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
	Cursor       uint64
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
		topicStartPosition := byte(headerBytes + headerBytesWithTopic)
		topicEndPosition := topicStartPosition + raw[2]
		topic := raw[topicStartPosition:topicEndPosition]
		payload := raw[topicEndPosition:dataLength]
		return Message{
			Command:      command,
			Topic:        topic,
			TopicStart:   int(topicStartPosition),
			Payload:      payload,
			PayloadStart: int(topicEndPosition),
			Raw:          raw}
	case CommandSubscribe:
		topicLength := binary.LittleEndian.Uint16(raw[2:4])
		cursor := binary.LittleEndian.Uint64(raw[4:12])
		topicEndPosition := 12 + topicLength
		check := uint16(len(raw))
		if check != topicEndPosition {
			fmt.Printf("check:failed:%d:%d\n", check, topicEndPosition)
			fmt.Println(raw)
		}
		topic := raw[12:topicEndPosition]
		return Message{Command: command, Topic: topic, Raw: raw, Cursor: cursor}
	case CommandUnsubscribe:
		topicLength := binary.LittleEndian.Uint16(raw[2:4])
		topicEndPosition := 12 + topicLength
		topic := raw[12:topicEndPosition]
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
const headerBytesWithTopic = 2 + 8 // topic length + cursor

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
		topicLength := len(message.Topic)
		payloadLength := len(message.Payload)
		restOfMessageLength := headerBytesWithTopic + topicLength + payloadLength
		raw := make([]byte, headerBytes+headerBytesWithTopic+topicLength+payloadLength)
		raw[0] = CommandInform
		raw[1] = byte(restOfMessageLength)
		topicBytes := topicBytes(uint16(topicLength))
		copy(raw[2:], topicBytes)
		copy(raw[4:], cursorBytes(uint64(message.Cursor)))
		copy(raw[12:12+topicLength], message.Topic)
		copy(raw[12+topicLength:], message.Payload)
		message.Raw = raw
		return
	default:
		raw := make([]byte, headerBytes)
		raw[0] = message.Command
		raw[1] = 0
		message.Raw = raw
	}
}

func CommandAsString(command byte) string {
	switch command {
	case CommandConnect:
		return "CommandConnect"
	case CommandConnack:
		return "CommandConnack"
	case CommandPublish:
		return "CommandPublish"
	case CommandSubscribe:
		return "CommandSubscribe"
	case CommandSuback:
		return "CommandSuback"
	case CommandUnsubscribe:
		return "CommandUnsubscribe"
	case CommandUnsuback:
		return "CommandUnsuback"
	case CommandInform:
		return "CommandInform"
	default:
		return "Unmapped"
	}
}
