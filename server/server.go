package server

import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"os"
	"os/signal"
	"syscall"
	"context"
	"time"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"
)

var upgrader = websocket.Upgrader{} // use default options

func webSocketHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		log.Printf("Message Rec from %v : %v", c.RemoteAddr().String(), string(message))
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func RunServer() {
	var addr = fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("port"))

	var stop = make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	routes := chi.NewRouter()
	routes.HandleFunc(viper.GetString("path"), webSocketHandler)

	server := http.Server{
		Addr:    addr,
		Handler: routes,
	}

	go server.ListenAndServe()
	log.Println(fmt.Sprintf("Server has been started on %s", addr))

	<-stop

	log.Println("Server is going down...")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	server.Shutdown(ctx)
}
