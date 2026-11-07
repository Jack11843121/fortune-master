package webutil

import (
	"encoding/json"
	"net/http"
)

// JSON 以 JSON 形式写出响应。
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// OK 写出 200 响应。
func OK(w http.ResponseWriter, v any) { JSON(w, http.StatusOK, v) }

// Error 以统一格式写出错误。
func Error(w http.ResponseWriter, status int, msg string) {
	JSON(w, status, map[string]string{"error": msg})
}

// ReadJSON 解析请求体中的 JSON（限制大小，防止滥用）。
func ReadJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20))
	return dec.Decode(dst)
}