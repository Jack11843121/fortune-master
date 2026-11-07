package fortune

import (
	"fmt"
	"sort"
	"strings"
)

// 本文件包含五行、天干、地支、十神等基础数术表与推演逻辑，
// 用于解读与匹配；具体排盘（干支/藏干/纳音/大运）交由 calendar.go 基于 lunar-go 实现。

// 天干五行
var ganWuXing = map[string]string{
	"甲": "木", "乙": "木",
	"丙": "火", "丁": "火",
	"戊": "土", "己": "土",
	"庚": "金", "辛": "金",
	"壬": "水", "癸": "水",
}

// 天干阴阳（true = 阳）
var ganYinYang = map[string]bool{
	"甲": true, "丙": true, "戊": true, "庚": true, "壬": true,
	"乙": false, "丁": false, "己": false, "辛": false, "癸": false,
}

// 地支五行
var zhiWuXing = map[string]string{
	"子": "水", "丑": "土", "寅": "木", "卯": "木",
	"辰": "土", "巳": "火", "午": "火", "未": "土",
	"申": "金", "酉": "金", "戌": "土", "亥": "水",
}

// 地支藏干（本气在前，中余气在后）
var zhiHiddenGan = map[string][]string{
	"子": {"癸"},
	"丑": {"己", "癸", "辛"},
	"寅": {"甲", "丙", "戊"},
	"卯": {"乙"},
	"辰": {"戊", "乙", "癸"},
	"巳": {"丙", "庚", "戊"},
	"午": {"丁", "己"},
	"未": {"己", "丁", "乙"},
	"申": {"庚", "壬", "戊"},
	"酉": {"辛"},
	"戌": {"戊", "辛", "丁"},
	"亥": {"壬", "甲"},
}

// 生肖
var zhiAnimal = map[string]string{
	"子": "鼠", "丑": "牛", "寅": "虎", "卯": "兔",
	"辰": "龙", "巳": "蛇", "午": "马", "未": "羊",
	"申": "猴", "酉": "鸡", "戌": "狗", "亥": "猪",
}

// 五行相生：我 -> 生出的五行
var wuXingSheng = map[string]string{"木": "火", "火": "土", "土": "金", "金": "水", "水": "木"}

// 五行相克：我 -> 克制的五行
var wuXingKe = map[string]string{"木": "土", "土": "水", "水": "火", "火": "金", "金": "木"}

// 五行顺序（用于展示/统计）
var wuXingOrder = []string{"金", "木", "水", "火", "土"}

// 五行对应的方向 / 颜色 / 脏腑（用于趋吉避凶与健康解读）
var wuXingAttr = map[string]struct {
	Direction string
	Color     string
	Organ     string
}{
	"金": {"西方", "白色/金色", "肺、呼吸道、皮肤"},
	"木": {"东方", "绿色/青色", "肝、胆、筋"},
	"水": {"北方", "黑色/蓝色", "肾、泌尿、耳"},
	"火": {"南方", "红色/紫色", "心、小肠、血脉"},
	"土": {"中央", "黄色/棕色", "脾胃、消化"},
}

// ShiShen 计算某天干相对“日主”的十神。
func shiShen(dayGan, otherGan string) string {
	if dayGan == "" || otherGan == "" {
		return ""
	}
	dw := ganWuXing[dayGan]
	ow := ganWuXing[otherGan]
	sameYY := ganYinYang[dayGan] == ganYinYang[otherGan]

	if ow == dw {
		if sameYY {
			return "比肩"
		}
		return "劫财"
	}
	// 我生者：日主五行 -> 生出 other
	if wuXingSheng[dw] == ow {
		if sameYY {
			return "食神"
		}
		return "伤官"
	}
	// 我克者
	if wuXingKe[dw] == ow {
		if sameYY {
			return "偏财"
		}
		return "正财"
	}
	// 克我者：other 克日主
	if wuXingKe[ow] == dw {
		if sameYY {
			return "七杀"
		}
		return "正官"
	}
	// 生我者
	if wuXingSheng[ow] == dw {
		if sameYY {
			return "偏印"
		}
		return "正印"
	}
	return ""
}

// countWuXing 统计五行出现次数（天干权重 1，地支权重 1，藏干权重 0.5）。
func countWuXing(pillars []Pillar, weights ...bool) map[string]int {
	includeHidden := true
	if len(weights) > 0 {
		includeHidden = weights[0]
	}
	res := map[string]int{"金": 0, "木": 0, "水": 0, "火": 0, "土": 0}
	for _, p := range pillars {
		if w, ok := ganWuXing[p.Gan]; ok {
			res[w] += 2
		}
		if w, ok := zhiWuXing[p.Zhi]; ok {
			res[w] += 2
		}
		if includeHidden {
			for _, hg := range p.HiddenGan {
				if w, ok := ganWuXing[hg]; ok {
					res[w]++
				}
			}
		}
	}
	return res
}

// topWuXing 返回统计值最高的五行列表（从多到少）。
func topWuXing(count map[string]int) []string {
	type kv struct {
		k string
		v int
	}
	list := make([]kv, 0, len(count))
	for k, v := range count {
		list = append(list, kv{k, v})
	}
	sort.SliceStable(list, func(i, j int) bool {
		if list[i].v == list[j].v {
			return indexOf(wuXingOrder, list[i].k) < indexOf(wuXingOrder, list[j].k)
		}
		return list[i].v > list[j].v
	})
	out := make([]string, 0, len(list))
	for _, e := range list {
		out = append(out, e.k)
	}
	return out
}

func indexOf(list []string, s string) int {
	for i, v := range list {
		if v == s {
			return i
		}
	}
	return -1
}

// bodyStrength 粗略判断日主身强身弱，返回 "强" / "弱" / "中"。
func bodyStrength(dayGan string, monthZhi string, count map[string]int) string {
	dw := ganWuXing[dayGan]
	score := 0
	// 得令：月支五行生/同 日主
	if mw := zhiWuXing[monthZhi]; mw == dw {
		score += 40
	} else if wuXingSheng[mw] == dw {
		score += 20
	} else if wuXingKe[dw] == mw || wuXingKe[mw] == dw {
		score -= 20
	}
	// 得地/得势：其余计数
	score += count[dw] * 2
	for _, w := range wuXingOrder {
		if wuXingSheng[w] == dw { // 生我者
			score += count[w]
		}
	}
	if score >= 45 {
		return "强"
	}
	if score <= 25 {
		return "弱"
	}
	return "中"
}

// xiYong 根据身强弱推导喜用神五行。
func xiYong(dayGan, monthZhi string, count map[string]int) (xi []string, ji []string) {
	dw := ganWuXing[dayGan]
	strength := bodyStrength(dayGan, monthZhi, count)

	shengWo := wuXingInvert(wuXingSheng, dw) // 生我者（印）
	keWo := wuXingInvert(wuXingKe, dw)       // 克我者（官杀）
	woKe := wuXingKe[dw]                     // 我克者（财）
	woSheng := wuXingSheng[dw]               // 我生者（食伤）

	switch strength {
	case "弱", "中":
		xi = []string{shengWo, dw}
		ji = []string{keWo, woKe, woSheng}
	default: // 强
		xi = []string{woKe, woSheng, keWo}
		ji = []string{shengWo, dw}
	}
	xi = dedupe(xi)
	ji = dedupe(ji)
	return xi, ji
}

// wuXingInvert 反查五行生克：给定关系表和结果五行，返回“产生该结果的源五行”。
func wuXingInvert(rel map[string]string, target string) string {
	for k, v := range rel {
		if v == target {
			return k
		}
	}
	return ""
}

func dedupe(in []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(in))
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

// 天干五合 / 地支六合 / 六冲 / 相刑（用于合婚与招聘匹配）
var ganHe = map[string]string{
	"甲": "己", "己": "甲",
	"乙": "庚", "庚": "乙",
	"丙": "辛", "辛": "丙",
	"丁": "壬", "壬": "丁",
	"戊": "癸", "癸": "戊",
}

var zhiLiuHe = map[string]string{
	"子": "丑", "丑": "子",
	"寅": "亥", "亥": "寅",
	"卯": "戌", "戌": "卯",
	"辰": "酉", "酉": "辰",
	"巳": "申", "申": "巳",
	"午": "未", "未": "午",
}

var zhiSanHe = map[string][]string{
	"申子辰": {"申", "子", "辰"},
	"寅午戌": {"寅", "午", "戌"},
	"巳酉丑": {"巳", "酉", "丑"},
	"亥卯未": {"亥", "卯", "未"},
}

var zhiLiuChong = map[string]string{
	"子": "午", "午": "子",
	"丑": "未", "未": "丑",
	"寅": "申", "申": "寅",
	"卯": "酉", "酉": "卯",
	"辰": "戌", "戌": "辰",
	"巳": "亥", "亥": "巳",
}

var zhiXiangHai = map[string]string{
	"子": "未", "未": "子",
	"丑": "午", "午": "丑",
	"寅": "巳", "巳": "寅",
	"卯": "辰", "辰": "卯",
	"申": "亥", "亥": "申",
	"酉": "戌", "戌": "酉",
}

func formatWuXingList(wxs []string) string {
	parts := make([]string, len(wxs))
	for i, w := range wxs {
		parts[i] = fmt.Sprintf("%s（%s）", w, wuXingAttr[w].Color)
	}
	return strings.Join(parts, "、")
}