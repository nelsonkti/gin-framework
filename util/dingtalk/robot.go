package dingtalk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Robot 是一个封装钉钉机器人的结构体
type Robot struct {
	webhookURL string
	secret     string
	httpClient *http.Client
}

// NewRobot 创建一个新的 Robot 实例
func NewRobot(webhookURL, secret string) *Robot {
	return &Robot{
		webhookURL: webhookURL,
		secret:     secret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type MessageOptions func(message *Message)

// WithAtAll 设置需要@所有人
func WithAtAll() MessageOptions {
	return func(message *Message) {
		message.At["isAtAll"] = true
	}
}

// WithAtMobiles 设置需要@的手机号
func WithAtMobiles(mobiles []string) MessageOptions {
	return func(message *Message) {
		message.At["atMobiles"] = mobiles
	}
}

// WithAtUserIds 设置需要@的用户ID
func WithAtUserIds(userIds []string) MessageOptions {
	return func(message *Message) {
		message.At["atUserIds"] = userIds
	}
}

// Message 是发送到钉钉的消息结构
type Message struct {
	MsgType  string                 `json:"msgtype"`
	Text     map[string]string      `json:"text,omitempty"`
	Link     map[string]interface{} `json:"link,omitempty"`
	Markdown map[string]string      `json:"markdown,omitempty"`
	At       map[string]interface{} `json:"at,omitempty"`
}

// SendText 发送文本消息到钉钉群
func (bot *Robot) SendText(content string, opts ...MessageOptions) error {
	message := &Message{
		MsgType: "text",
		Text: map[string]string{
			"content": content,
		},
		At: make(map[string]interface{}),
	}
	for _, opt := range opts {
		opt(message)
	}
	return bot.sendMessage(message)
}

// SendLink 发送链接消息到钉钉群
func (bot *Robot) SendLink(title, text, messageURL, picURL string, opts ...MessageOptions) error {
	message := &Message{
		MsgType: "link",
		Link: map[string]interface{}{
			"title":      title,
			"text":       text,
			"messageUrl": messageURL,
			"picUrl":     picURL,
		},
		At: make(map[string]interface{}),
	}
	for _, opt := range opts {
		opt(message)
	}
	return bot.sendMessage(message)
}

func (bot *Robot) SendMarkdown(title, text string, opts ...MessageOptions) error {
	message := &Message{
		MsgType: "markdown",
		Markdown: map[string]string{
			"title": title,
			"text":  text,
		},
		At: make(map[string]interface{}),
	}
	for _, opt := range opts {
		opt(message)
	}
	return bot.sendMessage(message)
}

// sign 生成签名
func (bot *Robot) sign() (string, string) {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	stringToSign := fmt.Sprintf("%s\n%s", timestamp, bot.secret)
	h := hmac.New(sha256.New, []byte(bot.secret))
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return timestamp, url.QueryEscape(signature)
}

// sendMessage 发送消息到钉钉
func (bot *Robot) sendMessage(message *Message) error {
	timestamp, sign := bot.sign()
	webhookURL := fmt.Sprintf("%s&timestamp=%s&sign=%s", bot.webhookURL, timestamp, sign)

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := bot.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 response: %v", resp.Status)
	}

	return nil
}
