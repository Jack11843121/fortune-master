// Package ai 封装大模型对话（OpenAI 兼容接口，DeepSeek/通义/Moonshot 通用），
// 未配置 API 时走本地规则兜底，保证离线也能用。
package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Message 一条对话消息。
type Message struct {
	Role    string `json:"role"` // system / user / assistant
	Content string `json:"content"`
}

// Config 大模型配置。
type Config struct {
	Enabled bool
	BaseURL string // 如 https://api.deepseek.com（自动补 /chat/completions）
	APIKey  string
	Model   string
}

// Client 大模型客户端。
type Client struct {
	cfg  Config
	http *http.Client
}

// New 创建客户端。
func New(c Config) *Client {
	return &Client{cfg: c, http: &http.Client{Timeout: 60 * time.Second}}
}

// Ready 是否已配置可用的远程模型。
func (c *Client) Ready() bool {
	return c.cfg.Enabled && c.cfg.APIKey != "" && c.cfg.BaseURL != ""
}

// Generate 调用远程大模型生成回复。
func (c *Client) Generate(msgs []Message) (string, error) {
	if !c.Ready() {
		return "", errors.New("未配置大模型")
	}
	endpoint := c.cfg.BaseURL
	if !strings.HasSuffix(endpoint, "/chat/completions") {
		endpoint = strings.TrimSuffix(endpoint, "/") + "/chat/completions"
	}
	body, err := json.Marshal(map[string]any{
		"model":       c.cfg.Model,
		"messages":    msgs,
		"temperature": 0.7,
		"stream":      false,
	})
	if err != nil {
		return "", err
	}
	req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("调用大模型失败: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("大模型返回 %d: %s", resp.StatusCode, truncate(string(data), 200))
	}
	var out struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return "", fmt.Errorf("解析大模型响应失败: %w", err)
	}
	if len(out.Choices) == 0 {
		return "", errors.New("大模型未返回内容")
	}
	return strings.TrimSpace(out.Choices[0].Message.Content), nil
}

// Local 离线兜底：由 handler 传入命主八字摘要与问题，生成一段吉祥、正面的命理回应。
// 不做远程调用，保证未配模型时功能可演示。
func Local(baziSummary, question string) string {
	q := strings.TrimSpace(question)
	if q == "" {
		q = "请大师为我指点一二"
	}
	var sb strings.Builder
	sb.WriteString("这位缘主，贫道为您细看。")
	if baziSummary != "" {
		sb.WriteString("您的命局：")
		sb.WriteString(baziSummary)
		sb.WriteString("。")
	}
	sb.WriteString("您所问「")
	sb.WriteString(q)
	sb.WriteString("」，贫道观之：命理云「一命二运三风水，四积阴德五读书」。")
	sb.WriteString("命为根基、运为时机，凡事知命而不认命——顺势而为、勤勉精进、广结善缘，自能趋吉避凶、时来运转。")
	sb.WriteString("您可先填出生信息排盘（个人命理），贫道再依您的八字逐项详批，更为精准。")
	return sb.String()
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}