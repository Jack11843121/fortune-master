package fortune

import (
	"fmt"
	"strings"
)

// 神煞：以日干、年支、日支为依据推定的吉星，为命局添彩、点明天赋与贵人方向。
// 均为正面表述，贴近普通用户（含老人）的阅读习惯。

// 三合局分组（用于驿马/华盖/将星等以三合局为基准的神煞）。
var sanHeGroup = map[string]string{
	"申": "申子辰", "子": "申子辰", "辰": "申子辰",
	"寅": "寅午戌", "午": "寅午戌", "戌": "寅午戌",
	"巳": "巳酉丑", "酉": "巳酉丑", "丑": "巳酉丑",
	"亥": "亥卯未", "卯": "亥卯未", "未": "亥卯未",
}

var yimaBranch = map[string]string{"申子辰": "寅", "寅午戌": "申", "巳酉丑": "亥", "亥卯未": "巳"}
var huagaiBranch = map[string]string{"申子辰": "辰", "寅午戌": "戌", "巳酉丑": "丑", "亥卯未": "未"}
var jiangxingBranch = map[string]string{"申子辰": "子", "寅午戌": "午", "巳酉丑": "酉", "亥卯未": "卯"}

// 文昌贵人（日干）：主聪慧好学、文思敏捷。
var wenChang = map[string]string{
	"甲": "巳", "乙": "午", "丙": "申", "丁": "酉", "戊": "申",
	"己": "酉", "庚": "亥", "辛": "子", "壬": "寅", "癸": "卯",
}

// 羊刃（日干）：主刚强果断、执行力强，锋芒可英武亦可伤身，宜善用。
var yangRen = map[string]string{
	"甲": "卯", "丙": "午", "戊": "午", "庚": "酉", "壬": "子",
}

// 禄神（日干，临官）：主福禄自足、衣食无忧。
var luShen = map[string]string{
	"甲": "寅", "乙": "卯", "丙": "巳", "丁": "午", "戊": "巳",
	"己": "午", "庚": "申", "辛": "酉", "壬": "亥", "癸": "子",
}

// 金舆（日干）：主得贵人提携、出行得助、见重而贵。
var jinYu = map[string]string{
	"甲": "辰", "乙": "巳", "丙": "未", "丁": "申", "戊": "未",
	"己": "申", "庚": "戌", "辛": "亥", "壬": "丑", "癸": "寅",
}

// 红鸾（年支）：主婚恋桃花、人缘喜事。
var hongLuan = map[string]string{
	"子": "卯", "丑": "寅", "寅": "丑", "卯": "子", "辰": "亥", "巳": "戌",
	"午": "酉", "未": "申", "申": "未", "酉": "午", "戌": "巳", "亥": "辰",
}

// 神煞的通俗解读文案。
var shenshaMeaning = map[string]string{
	"文昌贵人": "聪慧好学、才思敏捷，读书明理、文采过人，一生不乏学习与晋升之机。",
	"驿马": "好动进取、志在四方，宜向外发展、走动求财，出差远行往往有机遇。",
	"华盖": "独享清雅、悟性偏高，宜研习文艺、哲思、技艺等，于热闹中自有一方天地。",
	"羊刃": "刚毅果敢、敢打敢拼，胆识与执行力俱佳，宜善用锋芒、以柔济刚则大吉。",
	"禄神": "福禄自足、衣食无忧，勤奋踏实即有所得，稳中自见富贵。",
	"将星": "有统御之才、领导气场，善带团队、能堪重任，是天然的领路人。",
	"金舆": "得贵人提携、出行得助，见重于人，常有「车马之贵」，顺遂多福。",
	"红鸾": "人缘桃花、喜事临门，社交迷人、易结良缘，宜珍惜善缘、守诚专一。",
}

// shenshaSection 神煞格局：列出命带的神煞并逐项解读。
func shenshaSection(b Bazi) Section {
	if len(b.Pillars) == 0 {
		return Section{Key: "shensha", Title: "神煞格局", Icon: "sparkles", Body: "暂无数据。"}
	}
	dayGan := b.DayMaster
	yearZhi := b.Pillars[0].Zhi

	found := []struct{ name, where string }{}

	// 文昌（查日支是否落于文昌位）
	if z, ok := wenChang[dayGan]; ok {
		if where := zhiInPillars(b, z); where != "" {
			found = append(found, struct{ name, where string }{"文昌贵人", where})
		}
	}
	// 驿马 / 华盖 / 将星（由年支定三合局，查落点）
	if grp, ok := sanHeGroup[yearZhi]; ok {
		if z := yimaBranch[grp]; z != "" {
			if where := zhiInPillars(b, z); where != "" {
				found = append(found, struct{ name, where string }{"驿马", where})
			}
		}
		if z := huagaiBranch[grp]; z != "" {
			if where := zhiInPillars(b, z); where != "" {
				found = append(found, struct{ name, where string }{"华盖", where})
			}
		}
		if z := jiangxingBranch[grp]; z != "" {
			if where := zhiInPillars(b, z); where != "" {
				found = append(found, struct{ name, where string }{"将星", where})
			}
		}
	}
	// 羊刃、禄神、金舆（查日干对应地支）
	for _, fn := range []struct {
		name string
		m    map[string]string
	}{
		{"羊刃", yangRen}, {"禄神", luShen}, {"金舆", jinYu},
	} {
		if z, ok := fn.m[dayGan]; ok {
			if where := zhiInPillars(b, z); where != "" {
				found = append(found, struct{ name, where string }{fn.name, where})
			}
		}
	}
	// 红鸾（查年支对应）
	if z, ok := hongLuan[yearZhi]; ok {
		if where := zhiInPillars(b, z); where != "" {
			found = append(found, struct{ name, where string }{"红鸾", where})
		}
	}

	var sb strings.Builder
	if len(found) == 0 {
		sb.WriteString("命局神煞不显，属平和清贵之局，不靠奇门而贵在自身踏实进取。\n\n")
		sb.WriteString("无奇门者自有机缘：脚踏实地、与人为善、广结善缘，福报自会水到渠成。")
	} else {
		sb.WriteString(fmt.Sprintf("命带以下神煞，是您命中的亮点与天赋所在：\n"))
		for _, f := range found {
			if m, ok := shenshaMeaning[f.name]; ok {
				sb.WriteString(fmt.Sprintf("\n· 【%s】（现于%s）：%s", f.name, f.where, m))
			} else {
				sb.WriteString(fmt.Sprintf("\n· 【%s】（现于%s）", f.name, f.where))
			}
		}
		sb.WriteString("\n\n神煞乃点睛之笔，好星须善加运用、以德养之，方能长久生辉。")
	}
	return Section{Key: "shensha", Title: "神煞格局", Icon: "sparkles", Body: sb.String()}
}

// zhiInPillars 判断某地支是否落于四柱，返回所在柱名（如「月柱」），未找到返回空串。
func zhiInPillars(b Bazi, z string) string {
	for _, p := range b.Pillars {
		if p.Zhi == z {
			return p.Name
		}
	}
	return ""
}