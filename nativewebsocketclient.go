package api

import (
	"fmt"
	"golang.org/x/net/websocket"
)

type nativeWebSocketClient struct {
	ws *websocket.Conn
}

var NativeWebSocketClient = new(nativeWebSocketClient)

// New web socket client connected.
func (w *nativeWebSocketClient) New(url string, origin string, protocol string) *nativeWebSocketClient {
	var err error
	w.ws, err = websocket.Dial(url, protocol, origin)
	if err != nil {
		fmt.Println(err)
	}

	return w
}

// Send string message
func (w *nativeWebSocketClient) Send(msg string) (err error) {
	defer func() { _ = w.ws.Close() }()

	err = w.SendMessage(msg)
	if err != nil {
		fmt.Println(err)
	}

	return err
}

// Send any message
func (w *nativeWebSocketClient) SendAny(a any) (err error) {
	code := websocket.Codec{}
	err = code.Send(w.ws, a)
	if err != nil {
		fmt.Println(err)
	}

	return
}

// Send a message.
func (w *nativeWebSocketClient) SendMessage(message string) error {
	buffer := []byte(message)
	return w.toBytes(buffer)
}

// Send broadcast a message to server.
func (w *nativeWebSocketClient) toBytes(message []byte) error {
	_, err := w.ws.Write(message)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
