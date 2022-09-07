package api

import (
	"github.com/kataras/iris/v12/websocket"
	"github.com/kavanahuang/logs"
)

type events struct{}

var Event = new(events)

func (e *events) WebsocketEvent(namespace string) websocket.Namespaces {
	return websocket.Namespaces{
		namespace: websocket.Events{
			websocket.OnNamespaceConnected: func(c *websocket.NSConn, msg websocket.Message) error {
				logs.Info("Connected to namespace: ", msg.Namespace)
				return nil
			},
			websocket.OnNamespaceDisconnect: func(c *websocket.NSConn, msg websocket.Message) error {
				logs.Info("Disconnected to namespace: ", msg.Namespace)
				return nil
			},
			"chat": func(c *websocket.NSConn, msg websocket.Message) error {
				logs.Info("Chat msg: ", string(msg.Body))
				return nil
			},
		},
	}
}
