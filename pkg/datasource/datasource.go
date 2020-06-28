package datasource

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/wcatron/esk/pkg/message"
	"github.com/wcatron/esk/pkg/websocket"
)

type FileDataSource struct {
	GenericHandler *websocket.DataSource
}

func NewDataSource() *FileDataSource {
	return &FileDataSource{
		GenericHandler: &websocket.DataSource{
			Write: make(chan websocket.WriteNotification),
			Read:  make(chan websocket.SubscriptionNotification),
		},
	}
}

func pathForTopic(topic string) string {
	return filepath.Join("data", topic+".bin")
}

func createDirectoryForTopic(topic string) {
	directory := filepath.Dir(pathForTopic(topic))
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		fmt.Printf("ds:createDirectoryForTopic:%s\n", directory)
		os.Mkdir(directory, 0700)
	}
}

func Write(msg message.Message) (cursor uint64, err error) {
	fmt.Printf("ds:write:%s\n", msg.Payload)
	topic := string(msg.Topic)
	createDirectoryForTopic(topic)
	path := pathForTopic(topic)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		fmt.Print(err)
		return 0, err
	}
	info, err := f.Stat()
	if err != nil {
		fmt.Print(err)
		return 0, err
	}
	data := append(msg.Payload, '\n')
	_, err = f.Write(data)
	return uint64(info.Size()), err
}

func Read(topic string, cursor uint64, client *websocket.Client) {
	f, err := os.OpenFile(pathForTopic(topic), os.O_RDONLY, 0600)
	if err != nil {
		print(err)
		return
	}
	log.Printf("ds:read:%d\n", cursor)
	f.Seek(int64(cursor), 0) // 0 here means offset from start of file
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if len(scanner.Bytes()) > 0 {
			fmt.Printf("scanning:%d:%s\n", cursor, scanner.Text())
			client.Inform(topic, cursor, scanner.Bytes())
			cursor += uint64(len(scanner.Bytes())) + 1 // 1 for the line split
		} else {
			fmt.Printf("scanning:%d:empty\n", cursor)
		}
	}
}

func (ds *FileDataSource) Listen() {
	for {
		select {
		case notification := <-ds.GenericHandler.Write:
			cursor, err := Write(notification.Message)
			if err != nil {
				fmt.Print(err)
			}
			notification.Message.Cursor = cursor
			notification.Pool.BroadcastWritten <- notification.Message
			break
		case notification := <-ds.GenericHandler.Read:
			Read(notification.Topic, notification.Cursor, notification.Client)
			break
		}
	}
}
