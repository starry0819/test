package main

import (
	"encoding/json"
	"fmt"
	"github.com/cihub/seelog"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"net/http"
)

type Client struct {
	// 用户id标识
	id string
	// 连接的websocket
	socket *websocket.Conn
	// 发送的消息
	send chan []byte
}

type ClientManger struct {
	// 客户端map
	clients map[*Client]bool
	// 广播的消息
	broadcast chan []byte
	// 新增的客户端
	register chan *Client
	// 注销的客户端
	unregister chan *Client
}

type Message struct {
	MessageType string `json:"message_type"`
	Device      string `json:"device"`
	Token       string `json:"token"`
	UserRole    string `json:"user_role"`
	UserNumber  string `json:"user_number"`
	Param       string `json:"param"`
}

type ResponseMessage struct {
	MessageType string `json:"message_type"`
	Response    string `json:"response"`
}

var (
	// 客户端管理者
	manager = ClientManger{
		broadcast:  make(chan []byte),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
)

func (manager *ClientManger) start() {
	for {
		select {
		// 有新连接接入
		case conn := <-manager.register:
			manager.clients[conn] = true
			// 将连接成功的消息json格式化
		// 有连接断开
		case conn := <-manager.unregister:
			// 判断连接状态, 如果是true, 则关闭send, 删除连接client的值
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
				//jsonMessage, _ := json.Marshal(&ResponseMessage{Response: "/A socket has disconnected."})
				//manager.send(jsonMessage, conn)
			}
		// 广播
		case message := <-manager.broadcast:
			// 遍历已经连接的客户端，把消息发送给他们
			for conn := range manager.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

func (manager *ClientManger) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}

//定义客户端结构体的read方法
func (c *Client) read() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		//读取消息
		messageType, message, err := c.socket.ReadMessage()
		fmt.Println("messageType: " , messageType)
		//如果有错误信息，就注销这个连接然后关闭
		fmt.Println("before message in read")
		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			seelog.Error(err.Error())
			break
		}
		fmt.Println("after message in read")
		//如果没有错误信息就把信息放入broadcast
		jsonMessage, _ := json.Marshal(&Message{UserNumber: c.id, Param: string(message)})
		manager.broadcast <- jsonMessage
	}
}

func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		//从send里读消息
		case message, ok := <-c.send:
			//如果没有消息
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			//有消息就写入，发送给web端
			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func main() {
	fmt.Println("Starting application...")
	//开一个goroutine执行开始程序
	go manager.start()
	//注册默认路由为 /ws ，并使用wsHandler这个方法
	http.HandleFunc("/ws", wsHandler)
	//监听本地的8011端口
	http.ListenAndServe(":8011", nil)
}

var conns = map[string]*websocket.Conn{}


func wsHandler(res http.ResponseWriter, req *http.Request) {
	// 支持跨域
	var upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	//将http协议升级成websocket协议
	conn, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return
	}

	uuidStr := uuid.Must(uuid.NewV4()).String()
	//每一次连接都会新开一个client，client.id通过uuid生成保证每次都是不同的
	client := &Client{id: uuidStr, socket: conn, send: make(chan []byte)}
	conns[uuidStr] = conn
	//注册一个新的链接
	manager.register <- client

	//启动协程收web端传过来的消息
	go client.read()
	//启动协程把消息返回给web端
	//go client.write()
}
