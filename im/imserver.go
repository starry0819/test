package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("Starting application...")
	//开一个goroutine执行开始程序
	//go manager.start()
	//注册默认路由为 /ws ，并使用wsHandler这个方法
	http.HandleFunc("/ws", wsHandler)
	//监听本地的8011端口
	if err := http.ListenAndServe(":8011", nil); err != nil {
		log.Fatal("ListenAndServe", err.Error())
	}
}

var (
	wsConnToUserInfo      = make(map[*websocket.Conn]*UserInfo)
	wsConnToHeartBeatTime = make(map[*websocket.Conn]int64)
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 支持跨域
	var upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// 将http协议升级为websocket
	wsConn, e := upgrader.Upgrade(w, r, nil);
	if e != nil {
		http.NotFound(w, r)
		return
	}

	if userInfo := wsConnToUserInfo[wsConn]; userInfo != nil {
		log.Fatal("on_open: Connection handle existing")
		return
	}
	wsConnToHeartBeatTime[wsConn] = time.Now().Unix()

	conn := WsConnect{
		wsConn:  wsConn,
		isClose: true,
	}

	go conn.Read()

}
