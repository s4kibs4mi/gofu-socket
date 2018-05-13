package server

import (
	"testing"
	"time"
	"net/url"
	"log"
	"github.com/gorilla/websocket"
	"os"
	"os/signal"
	"syscall"
)

func TestEchoServer(t *testing.T) {
	var stop = make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	var addr = "localhost:7009"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/echo"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Websocket failed to dial :", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Error in reply : ", err)
				return
			}
			log.Printf("Reply received : %s", string(message))
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			log.Println("Connection halted")
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("Message send error : ", err)
				return
			}
			log.Println("Message send successfully")
		case <-stop:
			log.Println("Connection has been interrupted")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Something went wrong : ", err)
				return
			}
			select {
			case <-done:
				log.Println("Connection closed")
			case <-time.After(time.Second):
			}
			return
		}
	}
}
