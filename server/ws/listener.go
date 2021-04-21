package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Room(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	socket, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		panic(err)
	}

	for {
		msgType, msg, err := socket.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Mensagem recebida: ", string(msg))

		err = socket.WriteMessage(msgType, msg)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
