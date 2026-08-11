package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DingTalkMessage 钉钉消息结构
type DingTalkMessage struct {
	MsgType  string                `json:"msgtype"`
	Text     *DingTalkTextContent  `json:"text,omitempty"`
	Markdown *DingTalkMarkdownContent `json:"markdown,omitempty"`
	At       *DingTalkAt           `json:"at,omitempty"`
}

// DingTalkTextContent 文本消息内容
type DingTalkTextContent struct {
	Content string `json:"content"`
}

// DingTalkMarkdownContent markdown消息内容
type DingTalkMarkdownContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// DingTalkAt @配置
type DingTalkAt struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}

// DingTalkResponse 钉钉响应
type DingTalkResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// SendDingTalkText 发送钉钉文本消息
func SendDingTalkText(webhookURL, content string, atMobiles []string, isAtAll bool) error {
	msg := &DingTalkMessage{
		MsgType: "text",
		Text: &DingTalkTextContent{
			Content: content,
		},
		At: &DingTalkAt{
			AtMobiles: atMobiles,
			IsAtAll:   isAtAll,
		},
	}
	return sendDingTalk(webhookURL, msg)
}

// SendDingTalkMarkdown 发送钉钉markdown消息
func SendDingTalkMarkdown(webhookURL, title, text string) error {
	msg := &DingTalkMessage{
		MsgType: "markdown",
		Markdown: &DingTalkMarkdownContent{
			Title: title,
			Text:  text,
		},
	}
	return sendDingTalk(webhookURL, msg)
}

func sendDingTalk(webhookURL string, msg *DingTalkMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化钉钉消息失败: %w", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("发送钉钉消息失败: %w", err)
	}
	defer resp.Body.Close()

	var dingResp DingTalkResponse
	if err := json.NewDecoder(resp.Body).Decode(&dingResp); err != nil {
		return fmt.Errorf("解析钉钉响应失败: %w", err)
	}

	if dingResp.ErrCode != 0 {
		return fmt.Errorf("钉钉返回错误: errcode=%d, errmsg=%s", dingResp.ErrCode, dingResp.ErrMsg)
	}

	return nil
}