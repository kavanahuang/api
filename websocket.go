package api

import (
	"context"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"github.com/kavanahuang/logs"
	"github.com/kavanahuang/system"
	"time"
)

type websocketClient struct {
	conn   *neffos.NSConn
	client *neffos.Client
}

var Websocket = new(websocketClient)

func (ws *websocketClient) New(endpoint string, namespace string, timeout time.Duration) *websocketClient {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cancel()

	dialer := websocket.DefaultGobwasDialer
	client, err := websocket.Dial(ctx, dialer, endpoint, Event.WebsocketEvent(namespace))
	if err != nil {
		logs.Fatal(err)
	}
	ws.client = client
	// defer client.Close()

	c, err := client.Connect(ctx, namespace)
	if err != nil {
		logs.Fatal(err)
	}

	ws.conn = c
	return ws
}

func (ws *websocketClient) Send(text []byte) bool {
	return ws.conn.Emit("chat", text)
}

func (ws *websocketClient) Close() {
	if err := ws.conn.Disconnect(nil); err != nil {
		logs.Error("Reply from server: error: ", err)
	}
}

func (ws *websocketClient) Terminal() {
	defer ws.client.Close()
	ws.conn.Emit("chat", []byte("Hello from Go client side!"))
	system.Terminal.Call(ws.Close, ws.Send)
}
