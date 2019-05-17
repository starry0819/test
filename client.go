package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
)

//定义连接的服务端的网址
var addr = flag.String("addr", "localhost:8011", "http service address")

func main() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	//var dialer *websocket.Dialer

	//通过Dialer连接websocket服务器
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//go timeWriter(conn)
	//打印接收到的消息或者错误

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}
		fmt.Printf("received: %s\n", message)
	}
}

//func timeWriter(conn *websocket.Conn) {
//	for {
//		time.Sleep(time.Second * 2)
//		conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
//	}
//}
