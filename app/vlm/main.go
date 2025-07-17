package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-lab/app/vlm/conf"
	"go-lab/pkg/client/websocket"
	"go-lab/pkg/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	cfgPath string
)

func init() {
	flag.StringVar(&cfgPath, "c", "./config.yaml", "")
	flag.Parse()
}

type SendData struct {
	DeviceId    string `json:"device_id"`
	DataType    string `json:"data_type"`
	ImageBase64 string `json:"image_base64"`
	Tmp         int64  `json:"tmp"`
}

type RevData struct {
	DeviceId string    `json:"device_id"`
	Content  string    `json:"content"`
	Labels   VlmLabels `json:"labels"`
	Tmp      int64     `json:"tmp"`
}

type VlmLabels struct {
	DataType string       `json:"data_type"`
	Data     []LabelValue `json:"data"`
}

type LabelValue struct {
	Type  string `json:"type"`
	Value string `json:"labels"`
}

func main() {

	log.Info("[main] starting project")
	testDir()
	// 解析配置文件
	err := conf.Init(cfgPath)
	if err != nil {
		log.WithFields(log.Fields{"cfg_path": cfgPath}).WithError(err).Error("[main] config init error")
		return
	}
	log.Infof("[main] image path : %s", conf.GConfig.ImagePath)

	var (
		imgInChan  = make(chan *utils.File, 1024*10)
		imgOutChan = make(chan *RevData, 1024)
		mp         = make(map[int64]interface{})
	)
	getImage(imgInChan)

	go doImg2Label(imgInChan, imgOutChan)
	go dealVLmLabel(imgOutChan, mp)

	for {
		if conf.GConfig.Out >= conf.GConfig.In {
			log.Infof("[main] Done : %d, ！！！", conf.GConfig.Out)

		}
		time.Sleep(5 * time.Second)
	}

}

func getImage(imgInChan chan *utils.File) {
	// 读取顶层目录下的文件和文件夹
	entries, err := os.ReadDir(conf.GConfig.ImagePath)
	if err != nil {
		panic(err)
	}
	for _, v := range entries {
		if v.IsDir() {
			peopleImgPath := filepath.Join(conf.GConfig.ImagePath, v.Name())
			err := ReadImgDir(peopleImgPath, imgInChan, ".jpg")
			if err != nil {
				log.WithFields(log.Fields{"getImage": getImage}).WithError(err).Error("[main] read dir error")
				return
			}
		}
	}

	inputLength := len(imgInChan)
	conf.GConfig.In = atomic.AddInt32(&conf.GConfig.In, int32(inputLength))

	fmt.Println(fmt.Printf("At：%s , 待处理图片共 【%d】张图片", time.Now().Format("2006-01-02 15:04:05"), conf.GConfig.In))

}

func doImg2Label(imgInChan chan *utils.File, imgOutChan chan *RevData) {

	for i := 0; i < conf.GConfig.ClientNum; i++ {
		client := NewC(conf.GConfig.WsAddr, imgInChan, imgOutChan)
		go client.Run()
	}

}

/*
*
标签处理
*/
func dealVLmLabel(labels chan *RevData, mp map[int64]interface{}) {

	for {
		select {
		case l := <-labels:

			atomic.AddInt32(&conf.GConfig.Out, 1)

			if value, okk := mp[l.Tmp]; okk {
				if arr, ok := value.([]*RevData); ok {
					arr = append(arr, l)
				}
			} else {
				arr := make([]*RevData, 5)
				arr = append(arr, l)
				mp[l.Tmp] = arr

			}

		}
	}

}

type C struct {
	Ws         *websocket.WebSocketClient
	InputChan  chan *utils.File
	OutputChan chan *RevData
}

func NewC(url string, in chan *utils.File, out chan *RevData) *C {

	ws := newWsClient(url)
	return &C{
		Ws:         ws,
		InputChan:  in,
		OutputChan: out,
	}
}

func (c *C) Run() {

	sinChan := make(chan int, 1)

	sinChan <- 1

	// 消息发送
	go func() {
		for {

			select {
			case f := <-c.InputChan:

				<-sinChan
				d := &SendData{
					DeviceId:    f.DeviceId,
					DataType:    "image",
					ImageBase64: f.Base64Str,
					Tmp:         f.At,
				}
				err := c.Ws.SendJSON(d)
				if err != nil {
					log.Errorf("发送异常， Err： %v", err)
				}
			default:
			}
		}
	}()

	// 接收消息
	for {
		select {
		case msg := <-c.Ws.Receive():
			rv := &RevData{}
			err := json.Unmarshal(msg, rv)
			if err != nil {
				log.Errorf("Unmarshal error， Err： %v", err)
				continue
			}
			c.OutputChan <- rv
			sinChan <- 1

			log.Printf("Received: %s At %s", rv.DeviceId, time.Now().Format("2006-01-02 15:04:05"))
		case err := <-c.Ws.Errors():
			log.Printf("Received Error: %v", err)
		}
	}
}

func newWsClient(url string) *websocket.WebSocketClient {
	//"ws://172.19.202.20:18765/ws"

	client := websocket.NewWebSocketClient(url, &websocket.Config{
		AutoReconnect: true,
		ReconnectWait: 3 * time.Second,
		BufferSize:    50,
		Origin:        "http://localhost:8080",
	})

	if err := client.Connect(); err != nil {
		log.Fatal("Initial connect failed:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动自动重连管理器
	go client.AutoReconnect(ctx)
	return client
}

func testDir() {
	// 打印当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取工作目录失败: %v", err)
	}
	fmt.Printf("⚠️ 程序工作目录: %s\n", wd)

	// 测试路径解析
	testPath := "./config.yaml"
	absPath, _ := filepath.Abs(testPath)
	fmt.Printf("⚠️ 试图加载的绝对路径: %s\n", absPath)
}

/*
*
仅遍历一层, cam03_20250709153110.jpg
*/
func ReadImgDir(dirPath string, dataChan chan *utils.File, fileType string) error {
	// 读取顶层目录下的文件和文件夹
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}
	// 遍历目录项
	for _, entry := range entries {

		if strings.Contains(entry.Name(), fileType) {
			info := strings.Split(entry.Name(), ".")
			devAndTmp := strings.Split(info[0], "_")
			fullPath := filepath.Join(dirPath, entry.Name())
			// 转换（10 进制, 64 位）
			num, err := strconv.ParseInt(devAndTmp[1], 10, 64)
			if err != nil {
				return err
			}

			if num%conf.GConfig.Interval != 0 {
				continue
			}

			f := &utils.File{
				FileName: entry.Name(),
				DeviceId: devAndTmp[0],
				At:       num,
			}
			if content, err := os.ReadFile(fullPath); err == nil {
				f.Content = content
				f.Base64Str = base64.StdEncoding.EncodeToString(content)
				dataChan <- f
			}
		}
	}
	return nil
}
