package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"fortune/internal/store"
	"fortune/internal/webutil"
)

// sanitizeUserText 清理用户输入：去首尾空白、剔除控制字符、限制长度。
// 展示层由 Vue 的插值自动转义（本应用不使用 v-html），双重防 XSS。
func sanitizeUserText(s string, max int) string {
	s = strings.TrimSpace(s)
	s = strings.Map(func(r rune) rune {
		if r < 0x20 && r != '\n' && r != '\t' {
			return -1
		}
		return r
	}, s)
	r := []rune(s)
	if len(r) > max {
		return string(r[:max])
	}
	return s
}

// ListPosts 社区帖子列表（公开可见 published）。
func (h *Handler) ListPosts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("size"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}
	q := h.db.Model(&store.CommunityPost{}).Where("status = ?", "published")
	if cat := r.URL.Query().Get("category"); cat != "" {
		q = q.Where("category = ?", cat)
	}
	var total int64
	q.Count(&total)
	var posts []store.CommunityPost
	q.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&posts)
	webutil.OK(w, map[string]any{"posts": h.attachNick(posts), "total": total, "page": page, "size": size})
}

type postView struct {
	store.CommunityPost
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type commentView struct {
	store.CommunityComment
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// userMeta 取用户昵称与头像（用于社区展示，微信/支付宝登录后即真实头像）。
func (h *Handler) userMeta(id string) (nick, avatar string) {
	var u store.User
	if h.db.First(&u, "id = ?", id).Error == nil {
		return u.Nickname, u.Avatar
	}
	return "", ""
}

func (h *Handler) attachNick(posts []store.CommunityPost) []postView {
	out := make([]postView, len(posts))
	for i, p := range posts {
		n, a := h.userMeta(p.UserID)
		out[i] = postView{CommunityPost: p, Nickname: n, Avatar: a}
	}
	return out
}

func (h *Handler) attachComments(comments []store.CommunityComment) []commentView {
	out := make([]commentView, len(comments))
	for i, c := range comments {
		n, a := h.userMeta(c.UserID)
		out[i] = commentView{CommunityComment: c, Nickname: n, Avatar: a}
	}
	return out
}

// GetPost 帖子详情 + 评论。
func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	id := chiUint(r)
	var post store.CommunityPost
	if h.db.First(&post, "id = ?", id).Error != nil {
		webutil.Error(w, http.StatusNotFound, "帖子不存在")
		return
	}
	if post.Status != "published" {
		webutil.Error(w, http.StatusNotFound, "帖子不可见")
		return
	}
	var comments []store.CommunityComment
	h.db.Where("post_id = ? AND status = ?", id, "published").Order("id asc").Limit(200).Find(&comments)
	webutil.OK(w, map[string]any{"post": post, "comments": h.attachComments(comments)})
}

// CreatePost 发布帖子（需登录）。
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if !h.postLimiter.Allow(clientIP(r)) {
		webutil.Error(w, http.StatusTooManyRequests, "发帖过于频繁，请稍后再试")
		return
	}
	u := h.currentUser(r)
	if u == nil {
		webutil.Error(w, http.StatusUnauthorized, "请先登录后再发帖")
		return
	}
	uid := u.ID
	var req struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Category string `json:"category"`
	}
	if err := webutil.ReadJSON(r, &req); err != nil {
		webutil.Error(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	title := sanitizeUserText(req.Title, 120)
	content := sanitizeUserText(req.Content, 2000)
	if title == "" || content == "" {
		webutil.Error(w, http.StatusBadRequest, "标题与内容不能为空")
		return
	}
	cat := sanitizeUserText(req.Category, 32)
	if cat == "" {
		cat = "命理趣谈"
	}
	post := store.CommunityPost{UserID: uid, Title: title, Content: content, Category: cat, Status: "published"}
	h.db.Create(&post)
	h.incGrowth(uid, 5)
	webutil.OK(w, post)
}

// LikePost 点赞。
func (h *Handler) LikePost(w http.ResponseWriter, r *http.Request) {
	id := chiUint(r)
	h.db.Model(&store.CommunityPost{}).Where("id = ?", id).UpdateColumn("likes", gorm.Expr("likes + 1"))
	webutil.OK(w, map[string]any{"ok": true})
}

// CreateComment 发布评论（需登录）。
func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	if !h.postLimiter.Allow(clientIP(r)) {
		webutil.Error(w, http.StatusTooManyRequests, "评论过于频繁，请稍后再试")
		return
	}
	u := h.currentUser(r)
	if u == nil {
		webutil.Error(w, http.StatusUnauthorized, "请先登录后再评论")
		return
	}
	uid := u.ID
	id := chiUint(r)
	var post store.CommunityPost
	if h.db.First(&post, "id = ?", id).Error != nil {
		webutil.Error(w, http.StatusNotFound, "帖子不存在")
		return
	}
	var req struct {
		Content string `json:"content"`
	}
	if err := webutil.ReadJSON(r, &req); err != nil {
		webutil.Error(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	content := sanitizeUserText(req.Content, 500)
	if content == "" {
		webutil.Error(w, http.StatusBadRequest, "评论不能为空")
		return
	}
	c := store.CommunityComment{PostID: id, UserID: uid, Content: content, Status: "published"}
	h.db.Create(&c)
	h.incGrowth(uid, 2)
	webutil.OK(w, c)
}

// RecordShare 记录一次内容分享，累计成长值。
func (h *Handler) RecordShare(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		webutil.OK(w, map[string]any{"ok": false})
		return
	}
	uid := u.ID
	var req struct {
		ReadingID string `json:"reading_id"`
		Platform  string `json:"platform"`
	}
	_ = webutil.ReadJSON(r, &req)
	req.Platform = sanitizeUserText(req.Platform, 32)
	if req.Platform == "" {
		req.Platform = "wechat"
	}
	h.db.Create(&store.ShareRecord{UserID: uid, ReadingID: req.ReadingID, Platform: req.Platform})
	h.incGrowth(uid, 1)
	webutil.OK(w, map[string]any{"ok": true})
}

// incGrowth 累计成长值并重算等级。
func (h *Handler) incGrowth(uid string, delta int) {
	h.db.Model(&store.User{}).Where("id = ?", uid).Update("growth", gorm.Expr("growth + ?", delta))
	h.recalcLevel(uid)
}

// chiUint 解析路径中的 uint 参数（用于帖子/评论 ID）。
func chiUint(r *http.Request) uint {
	id, _ := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	return uint(id)
}