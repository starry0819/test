package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type UserInfo struct {
	id    string	// user_number + "_" + user_role + "_" + token
	isWeb bool		// true if device contains web_
}


type SendItem struct {
	MessageType string `json:"message_type"`
	Device      string `json:"device"`
	Token       string `json:"token"`
	UserRole    string `json:"user_role"`
	UserNumber  string `json:"user_number"`
	Param       string `json:"param"`
}

type WsConnect struct {
	wsConn *websocket.Conn
	inChan chan []byte
	outChan chan []byte
	closeChan chan byte
	isClose bool
	mutex sync.Mutex
}

func (conn *WsConnect) Read() {
	defer func() {
		log.Println("wsConn closed in read")
		if e := conn.wsConn.Close(); e != nil {
			// 注意这个方法内部直接调用了os.exit(1)
			log.Fatal("wsConn Close occurred error")
		}
	}()
	for {
		_, message, err := conn.wsConn.ReadMessage()
		// 如果有错误信息，注销该连接然后关闭
		if err != nil {
			goto ERR
		}
		sendItem := SendItem{}
		if err := json.Unmarshal(message, sendItem); err != nil {
			log.Println("发送消息Json解析错误")
			goto ERR
		}

	}
	ERR:

}