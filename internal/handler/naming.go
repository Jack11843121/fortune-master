package handler

import (
	"net/http"

	"fortune/internal/fortune"
	"fortune/internal/webutil"
)

type namingReq struct {
	Surname string         `json:"surname"`
	Gender  int            `json:"gender"`
	NameLen int            `json:"name_len"`
	Birth   fortune.Person `json:"birth"`
}

// Naming 宝宝起名：按八字喜用神生成起名建议。
func (h *Handler) Naming(w http.ResponseWriter, r *http.Request) {
	var req namingReq
	if err := webutil.ReadJSON(r, &req); err != nil {
		webutil.Error(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	if req.Surname == "" {
		webutil.Error(w, http.StatusBadRequest, "请填写姓氏")
		return
	}
	if req.Birth.Year < 1900 || req.Birth.Year > 2100 || req.Birth.Month < 1 || req.Birth.Month > 12 || req.Birth.Day < 1 || req.Birth.Day > 31 {
		webutil.Error(w, http.StatusBadRequest, "请输入正确的出生日期")
		return
	}

	b, err := fortune.BuildBazi(req.Birth)
	if err != nil {
		webutil.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	suggestions := fortune.SuggestNames(req.Surname, req.Gender, b, req.NameLen, 40)
	webutil.OK(w, map[string]any{
		"surname":     req.Surname,
		"day_master":  b.DayMaster,
		"xi_yong":     b.XiYong,
		"ji_shen":     b.JiShen,
		"suggestions": suggestions,
	})
}