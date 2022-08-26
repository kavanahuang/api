package api

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kavanahuang/log"
	"os"
	"time"
)

type websocketClient struct{}

var Websocket = new(websocketClient)

func (w *websocketClient) New(endpoint string, namespace string, timeout time.Duration) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(timeout))
	defer cancel()

	dialer := websocket.DefaultGobwasDialer
	client, err := websocket.Dial(ctx, dialer, endpoint, Event.WebsocketEvent(namespace))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	c, err := client.Connect(ctx, namespace)
	if err != nil {
		panic(err)
	}

	c.Emit("chat", []byte("Hello from Go client side!"))

	fmt.Fprint(os.Stdout, ">> ")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			log.Logs.Error("Scanner error: ", scanner.Err())
			return
		}

		text := scanner.Bytes()

		if bytes.Equal(text, []byte("exit")) {
			if err := c.Disconnect(nil); err != nil {
				log.Logs.Error("Reply from server: error: ", err)
			}
			break
		}

		ok := c.Emit("chat", text)
		if !ok {
			break
		}

		fmt.Fprint(os.Stdout, ">> ")
	}
}
