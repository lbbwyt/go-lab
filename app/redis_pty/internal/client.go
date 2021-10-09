package internal

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/peterh/liner"
	log "github.com/sirupsen/logrus"
	"go-lab/pkg/redis/prase"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type RedisServerInfo struct {
	Host string
	Port string
	DB   int
}

type Client struct {
	conn       *websocket.Conn
	send       chan []byte
	executor   *prase.RedisExecutor
	line       *liner.State
	serverInfo RedisServerInfo
}

func (c *Client) readPump() {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	c.ReadEvalPrintLoop()

	for {
		_, message, err := c.conn.ReadMessage()
		//客户端断开连接会返回error
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		//处理客户端发送的消息,解析用户的redis 命令，并返回执行结果
		res := c.executor.Execute(string(message))
		c.send <- []byte(res)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {

		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

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

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		conn:     conn,
		send:     make(chan []byte, 4),
		executor: prase.NewRedisExecutor("127.0.0.1:6379", "", ""),
	}
	go client.writePump()
	go client.readPump()
}

func (c *Client) ReadEvalPrintLoop() {
	c.line = liner.NewLiner()
	defer c.line.Close()
	c.line.SetCtrlCAborts(true)

	c.setCompletionHandler()

	reg, _ := regexp.Compile(`'.*?'|".*?"|\S+`)
	prompt := ""
	for {
		addr := c.addr()
		if c.serverInfo.DB > 0 && c.serverInfo.DB < 16 {
			prompt = fmt.Sprintf("%s[%d]> ", addr, c.serverInfo.DB)
		} else {
			prompt = fmt.Sprintf("%s> ", addr)
		}

		cmd, err := c.line.Prompt(prompt)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}

		cmds := reg.FindAllString(cmd, -1)
		if len(cmds) == 0 {
			continue
		} else {

			cmd := strings.ToLower(cmds[0])
			if cmd == "help" || cmd == "?" {
				printHelp(cmds)
			} else if cmd == "quit" || cmd == "exit" {
				os.Exit(0)
			} else if cmd == "clear" {
				println("Please use Ctrl + L instead")
			} else {
				//cliSendCommand(cmds)
			}
		}
	}
}

func printHelp(cmds []string) {
	args := cmds[1:]
	if len(args) == 0 {
		printGenericHelp()
	} else if len(args) > 1 {
		fmt.Println()
	} else {
		cmd := strings.ToUpper(args[0])
		for i := 0; i < len(prase.HelpCommands); i++ {
			if prase.HelpCommands[i][0] == cmd {
				printCommandHelp(prase.HelpCommands[i])
			}
		}
	}
}

func printGenericHelp() {
	msg :=
		`redis-cli
Type:	"help <command>" for help on <command>
	`
	fmt.Println(msg)
}

func printCommandHelp(arr []string) {
	fmt.Println()
	fmt.Printf("\t%s %s \n", arr[0], arr[1])
	fmt.Printf("\tGroup: %s \n", arr[2])
	fmt.Println()
}

func (c *Client) addr() string {
	var addr string

	host := "127.0.0.1"
	port := "6379"

	if c.serverInfo.Host != "" {
		host = c.serverInfo.Host
	}

	if c.serverInfo.Port != "" {
		port = c.serverInfo.Port
	}

	addr = fmt.Sprintf("%s:%s", host, port)
	return addr
}

func (c *Client) setCompletionHandler() {
	c.line.SetCompleter(func(line string) (c []string) {
		for _, i := range prase.HelpCommands {
			if strings.HasPrefix(i[0], strings.ToUpper(line)) {
				c = append(c, i[0])
			}
		}
		return
	})
}
