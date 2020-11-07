package config

import (
	"github.com/wcatron/esk/pkg/websocket"
)

func NewConfig(datasource *websocket.DataSource) *websocket.Config {
	return &websocket.Config{}
}