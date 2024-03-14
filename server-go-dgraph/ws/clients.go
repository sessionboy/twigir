package ws

import (
	"log"
	"net/http"
	"server/utils"
	"time"

	res "server/shares/response"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kataras/i18n"
)

const (
	// Time allowed to write a message to the peer.
	// 向对等方写消息所允许的时间。
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	// 允许时间从对等方读取下一个pong消息。
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	// 在此期间将ping发送给同级。 必须小于pongWait。
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	// 对等方允许的最大消息大小。
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn     *websocket.Conn
	send     chan []byte // 出站消息的缓冲通道
	appId    string      // `platform_userid`
	userid   string
	username string
	ip       string
	platform string
}

func Handler(c *gin.Context) {
	lang := c.GetString("lang")
	userid := c.GetString("user_id")
	username := c.GetString("user_username")
	agent := utils.Parseua(c)
	appId := agent.Platform + "_" + userid

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("ws connect error:", err)
		c.JSON(http.StatusInternalServerError, res.Err(i18n.Tr(lang, "fail_socket_conn")))
		return
	}
	client := &Client{
		conn:     conn,
		send:     make(chan []byte),
		appId:    appId,
		userid:   userid,
		username: username,
		ip:       agent.Ip,
		platform: agent.Platform,
	}
	Hub.register <- client

	// 通过在新的goroutine中进行所有工作，允许收集调用方引用的内存。
	go client.Write()
}

// 向客户端写数据，将hub消息推送到客户端
// 为每个连接启动一个运行 Write 的 goroutine。
// 这应用程序通过执行此goroutine的所有写入操作来确保连接最多有一个writer(写入器)
func (c *Client) Write() {
	// 写超时控制
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			// 当接收消息写入时，延长写超时时间。
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// 连接断开的情况
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
