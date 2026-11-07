package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"fortune/internal/ai"
	"fortune/internal/fortune"
	"fortune/internal/store"
	"fortune/internal/webutil"
)

const aiSystemPrompt = "你是「算命大师」APP 的资深命理师，用积极、吉祥、通俗易懂的中文为用户解答八字命理、婚恋、事业、财运、起名、风水等问题。立场理性正面，不渲染宿命论、不制造焦虑，不给出医疗/法律/证券投资等专业判断；结尾可给一句吉祥寄语。若用户问与自己命盘有关的问题，请结合其提供的八字信息作答。"

type aiChatReq struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	ReadingID string `json:"reading_id"`
	Message   string `json:"message"`
}

// AIChat —— AI 算命师对话（会员可用，按次扣额度）。
func (h *Handler) AIChat(w http.ResponseWriter, r *http.Request) {
	if !h.aiLimiter.Allow(clientIP(r)) {
		webutil.Error(w, http.StatusTooManyRequests, "问得太频繁啦，稍后再问")
		return
	}
	var req aiChatReq
	if err := webutil.ReadJSON(r, &req); err != nil {
		webutil.Error(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	msg := strings.TrimSpace(req.Message)
	if msg == "" || len([]rune(msg)) > 500 {
		webutil.Error(w, http.StatusBadRequest, "请输入有效问题（500 字以内）")
		return
	}
	u := h.currentUser(r)
	if u == nil {
		webutil.Error(w, http.StatusUnauthorized, "请先登录后再使用 AI 算命师")
		return
	}
	if u.Status == "banned" {
		webutil.Error(w, http.StatusForbidden, "该账号已被限制")
		return
	}
	if u.AICredits <= 0 {
		webutil.Error(w, http.StatusPaymentRequired, "AI 次数不足：请开通会员（月付赠 1000 次）或购买充值包")
		return
	}
	uid := u.ID
	session := req.SessionID
	if session == "" {
		session = uuid.NewString()
	}

	// 保存用户消息
	h.db.Create(&store.AiMessage{UserID: uid, SessionID: session, Role: "user", Content: msg})

	// 组织上下文：系统提示 + 命盘摘要 + 最近历史
	msgs := []ai.Message{{Role: "system", Content: aiSystemPrompt}}
	if sum := h.baziSummary(req.ReadingID); sum != "" {
		msgs = append(msgs, ai.Message{Role: "system", Content: "命主八字：" + sum})
	}
	for _, m := range h.recentMessages(session, 10) {
		msgs = append(msgs, ai.Message{Role: m.Role, Content: m.Content})
	}
	msgs = append(msgs, ai.Message{Role: "user", Content: msg})

	// 生成回答：远程大模型优先，失败/未配置则本地兜底
	reply := ""
	if h.ai.Ready() {
		if out, err := h.ai.Generate(msgs); err == nil && strings.TrimSpace(out) != "" {
			reply = strings.TrimSpace(out)
		}
	}
	if reply == "" {
		reply = ai.Local(h.baziSummary(req.ReadingID), msg)
	}

	h.db.Create(&store.AiMessage{UserID: uid, SessionID: session, Role: "assistant", Content: reply, Credits: 1})
	h.db.Model(&store.User{}).Where("id = ?", uid).Update("ai_credits", gorm.Expr("ai_credits - 1"))
	h.db.Create(&store.MembershipLog{UserID: uid, Kind: "consume", Credits: -1, Note: "AI 算命师对话"})

	webutil.OK(w, map[string]any{
		"session_id": session,
		"reply":      reply,
		"ai_credits": u.AICredits - 1,
	})
}

// AIHistory 返回某会话的对话历史（或用户全部会话的最近消息）。
func (h *Handler) AIHistory(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		webutil.Error(w, http.StatusUnauthorized, "请先登录")
		return
	}
	session := r.URL.Query().Get("session_id")
	q := h.db.Model(&store.AiMessage{}).Where("user_id = ?", u.ID)
	if session != "" {
		q = q.Where("session_id = ?", session)
	}
	var list []store.AiMessage
	q.Order("id asc").Limit(200).Find(&list)
	webutil.OK(w, map[string]any{"messages": list})
}

// recentMessages 取会话内最近的若干条对话（正序）。
func (h *Handler) recentMessages(session string, n int) []store.AiMessage {
	var list []store.AiMessage
	h.db.Where("session_id = ?", session).Order("id desc").Limit(n).Find(&list)
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list
}

// baziSummary 组装命盘摘要（供 AI 上下文 / 本地兜底使用）。
func (h *Handler) baziSummary(readingID string) string {
	if readingID == "" {
		return ""
	}
	var rd store.Reading
	if h.db.First(&rd, "id = ?", readingID).Error != nil {
		return ""
	}
	var res fortune.Result
	if err := json.Unmarshal([]byte(rd.Payload), &res); err != nil || res.Bazi == nil {
		return ""
	}
	b := res.Bazi
	gz := make([]string, 0, 4)
	for _, p := range b.Pillars {
		gz = append(gz, p.GanZhi)
	}
	return fmt.Sprintf("日主%s（五行%s），四柱%s，属%s，身%s，喜用神%s",
		b.DayMaster, b.DayMasterWuXing, strings.Join(gz, " "), b.Animal, b.BodyStrength, strings.Join(b.XiYong, "、"))
}