package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"sync"
	"time"
)

type WebSocketClient struct {
	url           string
	origin        string
	conn          *websocket.Conn
	mutex         sync.Mutex
	messageChan   chan []byte
	errorChan     chan error
	doneChan      chan struct{}
	isConnected   bool
	autoReconnect bool
	reconnectWait time.Duration
}

// 客户端配置选项
type Config struct {
	Origin        string        // 来源地址，可选
	AutoReconnect bool          // 是否自动重连
	ReconnectWait time.Duration // 重连等待时间
	BufferSize    int           // 消息通道缓冲区大小
}

func NewWebSocketClient(url string, config *Config) *WebSocketClient {
	// 设置默认配置
	if config == nil {
		config = &Config{
			AutoReconnect: true,
			ReconnectWait: 5 * time.Second,
			BufferSize:    100,
			Origin:        "http://localhost:8080",
		}
	}

	return &WebSocketClient{
		url:           url,
		origin:        config.Origin,
		autoReconnect: config.AutoReconnect,
		reconnectWait: config.ReconnectWait,
		messageChan:   make(chan []byte, config.BufferSize),
		errorChan:     make(chan error, 1),
		doneChan:      make(chan struct{}),
	}
}

// 建立连接
func (c *WebSocketClient) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isConnected {
		return errors.New("already connected")
	}

	// 建立 WebSocket 连接
	ws, err := websocket.Dial(c.url, "", c.origin)
	if err != nil {
		return fmt.Errorf("websocket dial error: %w", err)
	}

	c.conn = ws
	c.isConnected = true
	go c.readPump()
	return nil
}

// 主消息读取循环
func (c *WebSocketClient) readPump() {
	defer func() {
		c.mutex.Lock()
		c.isConnected = false
		c.mutex.Unlock()
	}()

	for {
		var msg []byte
		err := websocket.Message.Receive(c.conn, &msg)
		if err != nil {
			c.errorChan <- fmt.Errorf("read error: %w", err)
			return
		}

		select {
		case c.messageChan <- msg:
		case <-c.doneChan:
			return
		}
	}
}

// 发送消息
func (c *WebSocketClient) Send(data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.isConnected {
		return errors.New("not connected")
	}

	return websocket.Message.Send(c.conn, data)
}

// 接收消息（阻塞模式）
func (c *WebSocketClient) Receive() <-chan []byte {
	return c.messageChan
}

// 错误通道
func (c *WebSocketClient) Errors() <-chan error {
	return c.errorChan
}

// 自动重连管理
func (c *WebSocketClient) AutoReconnect(ctx context.Context) {
	for {
		select {
		case err := <-c.errorChan:
			log.Printf("WebSocket error: %v", err)

			if c.autoReconnect {
				log.Printf("Reconnecting in %v...", c.reconnectWait)
				time.Sleep(c.reconnectWait)

				if reconnectErr := c.Connect(); reconnectErr != nil {
					c.errorChan <- fmt.Errorf("reconnect failed: %w", reconnectErr)
				} else {
					log.Println("Reconnected successfully")
				}
			}
		case <-c.doneChan:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (c *WebSocketClient) SendJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.Send(data)
}

// 安全关闭连接
func (c *WebSocketClient) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.isConnected {
		return nil
	}

	close(c.doneChan)
	err := c.conn.Close()
	c.isConnected = false
	return err
}
