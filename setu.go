package setu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

var instance *setu
var logger = utils.GetModuleLogger("internal.logging")

type setu struct {
}

func init() {
	instance = &setu{}
	bot.RegisterModule(instance)
}

func (s *setu) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "aimerneige.test.setu",
		Instance: instance,
	}
}

// Init 初始化过程
// 在此处可以进行 Module 的初始化配置
// 如配置读取
func (s *setu) Init() {
}

// PostInit 第二次初始化
// 再次过程中可以进行跨 Module 的动作
// 如通用数据库等等
func (s *setu) PostInit() {
}

// Serve 注册服务函数部分
func (s *setu) Serve(b *bot.Bot) {
	b.OnGroupMessage(func(c *client.QQClient, msg *message.GroupMessage) {
		if msg.ToString() != "setu" {
			return
		}
		imgData, err := getSetuImg()
		if err != nil {
			logger.WithError(err).Error("Unable to get img from Lolicon API.")
		}
		imgMsgElement, err := c.UploadGroupImage(msg.GroupCode, imgData)
		if err != nil {
			logger.WithError(err).Error("Unable to Upload img.")
		}
		imgMsg := message.NewSendingMessage().Append(imgMsgElement)
		c.SendGroupMessage(msg.GroupCode, imgMsg)
	})
}

// Start 此函数会新开携程进行调用
// ```go
// 		go exampleModule.Start()
// ```
// 可以利用此部分进行后台操作
// 如 http 服务器等等
func (s *setu) Start(b *bot.Bot) {
}

// Stop 结束部分
// 一般调用此函数时，程序接收到 os.Interrupt 信号
// 即将退出
// 在此处应该释放相应的资源或者对状态进行保存
func (s *setu) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	// 别忘了解锁
	defer wg.Done()
}

func getSetuImg() (io.ReadSeeker, error) {
	apiURL := "https://api.lolicon.app/setu/v2"
	type loliconResponse struct {
		Error string `json:"error"`
		Data  []struct {
			Pid        int      `json:"pid"`
			P          int      `json:"p"`
			UID        int      `json:"uid"`
			Title      string   `json:"title"`
			Author     string   `json:"author"`
			R18        bool     `json:"r18"`
			Width      int      `json:"width"`
			Height     int      `json:"height"`
			Tags       []string `json:"tags"`
			Ext        string   `json:"ext"`
			UploadDate int64    `json:"uploadDate"`
			Urls       struct {
				Original string `json:"original"`
			} `json:"urls"`
		} `json:"data"`
	}
	var apiResponse loliconResponse
	apiResp, err := getRequest(apiURL)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(apiResp, &apiResponse); err != nil {
		return nil, err
	}
	if apiResponse.Error != "" {
		return nil, fmt.Errorf(apiResponse.Error)
	}
	imgURL := apiResponse.Data[0].Urls.Original
	imgBytes, err := getRequest(imgURL)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(imgBytes), nil
}

func getRequest(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
