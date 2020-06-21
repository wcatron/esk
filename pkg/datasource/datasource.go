package datasource

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/wcatron/esk/pkg/websocket"
)

type FileDataSource struct {
	GenericHandler *websocket.DataSource
}

func NewDataSource() *FileDataSource {
	return &FileDataSource{
		GenericHandler: &websocket.DataSource{
			Write: make(chan websocket.Message),
			Read:  make(chan *websocket.SubscriptionNotification),
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

func Write(message websocket.Message) (cursor int64, err error) {
	fmt.Printf("ds:write:%s\n", message.Payload)
	topic := string(message.Topic)
	createDirectoryForTopic(topic)
	f, err := os.OpenFile(pathForTopic(topic), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Print(err)
		return 0, err
	}
	info, err := f.Stat()
	if err != nil {
		fmt.Print(err)
		return 0, err
	}
	_, err = f.Write(append(message.Raw[message.PayloadStart:], '\n'))
	return info.Size(), err
}

func Read(topic string, cursor int64, client *websocket.Client) {
	f, err := os.OpenFile(pathForTopic(topic), os.O_RDONLY, 0600)
	if err != nil {
		print(err)
		return
	}
	log.Printf("ds:read:%d\n", cursor)
	f.Seek(cursor, 0) // 0 here means offset from start of file
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		log.Printf("scanning:%d:%s\n", cursor, scanner.Text())
		client.Inform(topic, cursor, scanner.Bytes())
		cursor += int64(len(scanner.Bytes())) + 1 // 1 for the line split
	}
}

func (ds *FileDataSource) Listen() {
	for {
		select {
		case message := <-ds.GenericHandler.Write:
			cursor, err := Write(message)
			if err != nil {
				fmt.Print(err)
			}
			fmt.Printf("ds:wrote:%d:%s\n", cursor, message.Payload)
			break
		case notification := <-ds.GenericHandler.Read:
			Read(notification.Topic, 0, notification.Client)
			break
		}
	}
}
