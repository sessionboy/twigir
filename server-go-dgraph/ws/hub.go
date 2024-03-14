package ws

import (
	"encoding/json"
	"strings"
)

var (
	Hub *WsHub
)

// 一个用户可能多终端登录 appid:Client
type ClientMap = map[string]*Client

// ws管理中心
type WsHub struct {
	clients    map[string]ClientMap // 已注册客户端列表 userid:[]*Client
	broadcast  chan Notification    // 客户端的消息
	register   chan *Client         // 注册来自客户端的请求
	unregister chan *Client         // 取消客户端注册
}

// 注册客户端
func (h *WsHub) Register(c *Client) {
	if h.clients[c.userid] == nil {
		h.clients[c.userid] = make(ClientMap)
	}
	// 允许如果已存在则覆盖，以应对重连的情况(丢弃旧的连接)
	h.clients[c.userid][c.appId] = c
}

// 注销客户端
func (h *WsHub) Unregister(c *Client) {
	if h.clients[c.userid] != nil {
		delete(h.clients[c.userid], c.appId)
		_, ok := <-c.send
		if ok {
			close(c.send)
		}
	}
}

// 是否在线
func (h *WsHub) OnLine(userid string) bool {
	client := h.clients[userid]
	if client == nil {
		return false
	}
	if len(client) == 0 {
		return false
	}
	return true
}

// 给客户端发送消息
func (h *WsHub) Send(n Notification) {
	for i := 0; i < len(n.Recipient); i++ {
		userid := n.Recipient[i]
		if h.OnLine(userid) {
			clients := h.clients[userid]
			msg, _ := json.Marshal(n)
			for _, client := range clients {
				client.send <- msg
			}
			// if len(n.Appid) > 0 {
			// 	// 指定客户端单发
			// 	c := h.GetClientByAppId(n.Appid)
			// 	if c != nil {
			// 		c.send <- msg
			// 	}
			// }
		}
	}
}

// 根据userid查找用户的client集合
func (h *WsHub) GetClientsByUserId(userid string) ClientMap {
	return h.clients[userid]
}

// 根据appid查找client
func (h *WsHub) GetClientByAppId(id string) *Client {
	parts := strings.Split(id, "_")
	if len(parts[1]) == 0 {
		return nil
	}
	if h.clients[parts[1]] == nil {
		return nil
	}
	if h.clients[parts[1]][id] == nil {
		return nil
	}
	return h.clients[parts[1]][id]
}

func InitStart() {
	Hub = &WsHub{
		clients:    make(map[string]ClientMap),
		broadcast:  make(chan Notification),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go Hub.run()
}

func (h *WsHub) run() {
	for {
		select {
		// 注册新用户
		case client := <-h.register:
			h.Register(client)

		// 用户离开，取消注册
		case client := <-h.unregister:
			if h.clients[client.userid] != nil {
				h.Unregister(client)
			}

			// 广播消息
			// case message := <-h.broadcast:
			// 	msg, _ := json.Marshal(message.Msg)
			// 	if len(message.Appid) > 0 {
			// 		// 单发，指定客户端连接
			// 		c := h.GetClientByAppId(message.Appid)
			// 		if c != nil {
			// 			c.send <- msg
			// 		}
			// 	} else if len(message.Userid) > 0 {
			// 		// 单发，该user所有的客户端连接
			// 		userClients := h.GetClientsByUserId(message.Userid)
			// 		for _, client := range userClients {
			// 			client.send <- msg
			// 		}
			// 	} else {
			// 		// 群发，所有用户的所有连接
			// 		for _, client := range h.clients {
			// 			for _, c := range client {
			// 				c.send <- msg
			// 			}
			// 		}
			// 	}
		}
	}
}
