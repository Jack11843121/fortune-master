package fortune

import (
	"fmt"
	"strings"
	"time"
)

// Wealth 财运与投资深度解读（含股市/基金倾向）。场景标识为 "wealth"。
func (e *Engine) Wealth(p Person) (*Result, error) {
	if err := validate(p); err != nil {
		return nil, err
	}
	b, err := BuildBazi(p)
	if err != nil {
		return nil, err
	}
	return &Result{
		Scenario: "wealth",
		Title:    fmt.Sprintf("%s · 财运与投资", nameOr(p.Name, "命主")),
		Bazi:     b,
		Sections: analyzeWealth(*b),
	}, nil
}

func analyzeWealth(b Bazi) []Section {
	return []Section{
		caiOverviewSection(b),
		caiStagesSection(b),
		investStyleSection(b),
		stockSectorSection(b),
		caiTimingSection(b),
		caiWarningSection(b),
		wealthLuckSection(b),
	}
}

// caiWuXing 财星对应的五行（我克者）。
func caiWuXing(b Bazi) string {
	return wuXingKe[b.DayMasterWuXing]
}

// caiTomb 五行财对应的墓库（财库）。
var caiTomb = map[string]string{"木": "未", "火": "戌", "土": "辰", "金": "丑", "水": "辰"}

// 1. 财星总览
func caiOverviewSection(b Bazi) Section {
	wx := caiWuXing(b)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("您是「%s」日主，财为「%s」（日主所克者）。", b.DayMaster, wx))
	sb.WriteString("\n\n")

	hasZheng := findShiShenStar(b, []string{"正财"}) != ""
	hasPian := findShiShenStar(b, []string{"偏财"}) != ""
	switch {
	case hasZheng && hasPian:
		sb.WriteString("财星：正财、偏财并见，财路较宽——正财主稳定正职、偏财主机遇之财，正偏兼得，宜两条腿走路、主业副业并举。")
	case hasZheng:
		sb.WriteString("财星：命带正财，财源以稳定正职与踏实积累为主，宜细水长流、稳中求富。")
	case hasPian:
		sb.WriteString("财星：命带偏财，财源多来自机遇、投资与副业，眼光活络、善抓风口，但须守好风控、见好就收。")
	default:
		sb.WriteString("财星：命局财星不显，财多靠专业本事与口碑渐积，宜以技立身、精进专长，切忌投机冒进。")
	}

	if zhi := caiTomb[wx]; zhi != "" {
		if where := zhiInPillars(b, zhi); where != "" {
			sb.WriteString(fmt.Sprintf("命中见「%s」（财库，现于%s），主善储蓄、财有所藏，聚财守财之力较佳。", zhi, where))
		}
	}

	sb.WriteString("\n\n担财：")
	switch b.BodyStrength {
	case "强":
		sb.WriteString("日主强，能担财任财，来财能守、可做大，属「身旺任财」之命，宜放大格局。")
	case "弱":
		sb.WriteString("日主偏弱，财多则难担，宜借势合伙、稳字当先，切莫独扛大额投资、忌高杠杆。")
	default:
		sb.WriteString("日主中和，担财得宜，稳健进取即可渐入佳境。")
	}

	if way := caiWay(b); way != "" {
		sb.WriteString("\n\n得财之道：" + way)
	}
	return Section{Key: "cai_overview", Title: "财星总览", Icon: "coins", Body: sb.String()}
}

func caiWay(b Bazi) string {
	switch {
	case b.ShiShenIncluded("食神") || b.ShiShenIncluded("伤官"):
		return "宜凭才艺、技术、创意「以技生财」，把一技之长做实做深，自然生财有道。"
	case b.ShiShenIncluded("正官") || b.ShiShenIncluded("七杀"):
		return "宜凭职位、平台与担当取财，正职稳进、步步高升，靠名望与权责得财。"
	case b.ShiShenIncluded("正印") || b.ShiShenIncluded("偏印"):
		return "宜凭学识、专业技能与稳定薪资生财，靠脑袋与本事吃饭。"
	case b.ShiShenIncluded("比肩") || b.ShiShenIncluded("劫财"):
		return "宜合伙合作、靠人脉与实干打拼，但须防分财、慎借贷。"
	default:
		return "宜稳健经营、开源节流，正职与副业兼顾。"
	}
}

// 2. 财运阶段
func caiStagesSection(b Bazi) Section {
	labels := []string{"早年（少年）", "青年（25-35）", "中年（36-50）", "晚年（50 后）"}
	var sb strings.Builder
	sb.WriteString("财运随四柱各管一段，每段侧重不同：\n")
	for i, p := range b.Pillars {
		if i >= len(labels) {
			break
		}
		good := ganWuXing[p.Gan] == caiWuXing(b) || zhiWuXing[p.Zhi] == caiWuXing(b) ||
			inWuXingList(b.XiYong, ganWuXing[p.Gan]) || inWuXingList(b.XiYong, zhiWuXing[p.Zhi])
		v := "此阶段财气平稳，宜稳打基础、以积累为主"
		if good {
			v = "此阶段财气较顺，宜勤勉经营、放大格局"
		}
		sb.WriteString(fmt.Sprintf("\n· %s（%s%s）：%s。", labels[i], p.Name, p.GanZhi, v))
	}
	sb.WriteString("\n\n财为养命之源，早积晚享、顺时而为，自能衣食丰足。")
	return Section{Key: "cai_stages", Title: "阶段财路", Icon: "trend", Body: sb.String()}
}

// 3. 投资性格
func investStyleSection(b Bazi) Section {
	var sb strings.Builder
	sb.WriteString("投资性格：")
	switch {
	case findShiShenStar(b, []string{"偏财"}) != "":
		sb.WriteString("命带偏财，眼光活络、善抓机会，属进取型，宜「进取 + 稳健」搭配，把好风控。")
	case findShiShenStar(b, []string{"正财"}) != "":
		sb.WriteString("命带正财，理性务实、细水长流，属稳健型，宜定投、储蓄与低位布局，忌追高。")
	default:
		sb.WriteString("财星不显，宜以本事立身，投资以「稳」为先，先保值再增值。")
	}

	sb.WriteString("\n\n风险承受：")
	switch b.BodyStrength {
	case "强":
		sb.WriteString("日主强、能扛波动，可适度参与中高风险品种，但须设好止损、分散配置。")
	case "弱":
		sb.WriteString("日主偏弱，宜以稳健品种为主，留足现金流，忌重仓单押、追涨杀跌。")
	default:
		sb.WriteString("日主中和，宜「核心 + 卫星」配置，稳中求进。")
	}

	sb.WriteString("\n\n适合的理财方式：")
	switch {
	case b.ShiShenIncluded("食神") || b.ShiShenIncluded("伤官"):
		sb.WriteString("偏创意与专业，宜选熟悉的行业主题基金、成长型标的，做透一个领域。")
	case b.ShiShenIncluded("正官") || b.ShiShenIncluded("正印"):
		sb.WriteString("偏稳健，宜存款、国债、宽基指数定投、货币基金等低波动品种。")
	default:
		sb.WriteString("宜「固收打底 + 权益浮动」的组合，定投为主、波段为辅。")
	}
	sb.WriteString("\n\n提醒：市场有风险，投资须量力而行、理性决策，以上为性格参考，非操作指令。")
	return Section{Key: "invest", Title: "投资性格", Icon: "coins", Body: sb.String()}
}

// 4. 股市板块倾向
func stockSectorSection(b Bazi) Section {
	var sb strings.Builder
	sb.WriteString("依喜用神五行，与您气场相合、可多关注的板块方向：\n")
	seen := map[string]bool{}
	for _, w := range b.XiYong {
		if seen[w] {
			continue
		}
		seen[w] = true
		if c, ok := stockSector[w]; ok {
			sb.WriteString(fmt.Sprintf("\n· 喜%s（%s）：%s", w, wuXingAttr[w].Color, c))
		}
	}

	sb.WriteString("\n\n操作风格建议：")
	switch {
	case b.BodyStrength == "强" || b.ShiShenIncluded("七杀") || b.ShiShenIncluded("劫财"):
		sb.WriteString("性子偏急、敢作敢为，可少量尝试波段，但务必严守纪律、即买即设止损。")
	case b.BodyStrength == "弱" || b.ShiShenIncluded("正印") || b.ShiShenIncluded("偏印"):
		sb.WriteString("宜中长线持有与基金定投，忌频繁交易、忌听小道消息。")
	default:
		sb.WriteString("宜「定投打底 + 少部分机动」，保持仓位弹性，进退有度。")
	}
	sb.WriteString("\n\n相合板块可作长期关注赛道，买卖时机仍须结合基本面与自身情况，切忌盲目跟风。")
	return Section{Key: "stock", Title: "股市·板块倾向", Icon: "trend", Body: sb.String()}
}

// 5. 发财时机
func caiTimingSection(b Bazi) Section {
	var sb strings.Builder
	cai := caiWuXing(b)
	now := time.Now().Year()

	var dYun []string
	for _, d := range b.DaYun {
		if d.StartYear > now {
			if ganWuXing[d.Gan] == cai || zhiWuXing[d.Zhi] == cai || inWuXingList(b.XiYong, ganWuXing[d.Gan]) {
				dYun = append(dYun, fmt.Sprintf("%d-%d 年（%s）", d.StartYear, d.EndYear, d.GanZhi))
			}
		}
	}
	if len(dYun) > 0 {
		sb.WriteString("接下来财气较旺的大运：\n")
		n := min(len(dYun), 3)
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("\n· %s，宜提前布局、放大收益。", dYun[i]))
		}
	} else {
		sb.WriteString("近期大运财气平稳，宜稳中蓄力、等待喜用之运再放大。")
	}

	var liuNian []string
	for _, ln := range b.LiuNian {
		if ln.Year >= now && ln.Year <= now+4 {
			if ganWuXing[ln.Gan] == cai || zhiWuXing[ln.Zhi] == cai || inWuXingList(b.XiYong, ganWuXing[ln.Gan]) {
				liuNian = append(liuNian, fmt.Sprintf("%d 年（%s）", ln.Year, ln.GanZhi))
			}
		}
	}
	if len(liuNian) > 0 {
		sb.WriteString("\n\n近期财气较旺的流年：")
		sb.WriteString(strings.Join(liuNian, "、") + "，可多加经营、放大格局。")
	}
	sb.WriteString("\n\n时机如四季，顺喜用之年进取、逢忌神之年守成，自能踏准节奏。")
	return Section{Key: "cai_timing", Title: "发财时机", Icon: "calendar", Body: sb.String()}
}

// 6. 破财警示
func caiWarningSection(b Bazi) Section {
	var sb strings.Builder
	var warns []string
	if b.ShiShenIncluded("劫财") {
		warns = append(warns, "命带劫财，防合伙分财、借贷难还，重大财务务必亲兄弟明算账")
	}
	if b.ShiShenIncluded("比肩") && b.BodyStrength == "强" {
		warns = append(warns, "比肩较旺，财易被他人分走，宜守财有道、减少无谓应酬与大手大脚")
	}
	if b.BodyStrength == "弱" {
		warns = append(warns, "身弱财重时勿贪大，控杠杆、留足现金流，避免大额冒进")
	}
	for _, ln := range b.LiuNian {
		if ln.Year == time.Now().Year() && (inWuXingList(b.JiShen, ganWuXing[ln.Gan]) || inWuXingList(b.JiShen, zhiWuXing[ln.Zhi])) {
			warns = append(warns, fmt.Sprintf("今年（%s）逢忌神流年，理财宜守不宜攻", ln.GanZhi))
		}
	}

	if len(warns) == 0 {
		sb.WriteString("命局财路尚稳，主要风险在外部环境。守住「不贪、不赌、不盲目加杠杆」三原则即可。")
	} else {
		sb.WriteString("破财警示，需留意：\n")
		for _, w := range warns {
			sb.WriteString("\n· " + w)
		}
	}
	sb.WriteString("\n\n财不入急门，稳字当头、敬畏风险，方能细水长流。")
	return Section{Key: "cai_warning", Title: "破财警示", Icon: "shield", Body: sb.String()}
}

// 7. 求财开运
func wealthLuckSection(b Bazi) Section {
	var sb strings.Builder
	sb.WriteString("求财开运要点：\n")
	sb.WriteString(fmt.Sprintf("· 求财方位：%s（多向此方发力，事半功倍）\n", strings.Join(xiYongDirs(b.XiYong), "、")))
	sb.WriteString(fmt.Sprintf("· 幸运颜色：%s（办公、穿戴可多用）\n", strings.Join(xiYongColors(b.XiYong), "、")))
	sb.WriteString(fmt.Sprintf("· 幸运数字：%s\n", strings.Join(luckyNumbers(b.XiYong), "、")))

	seen := map[string]bool{}
	for _, w := range b.XiYong {
		if lk, ok := wuXingLucky[w]; ok && !seen[w] {
			seen[w] = true
			sb.WriteString(fmt.Sprintf("· 喜%s开运物：%s；%s\n", w, lk.Gem, lk.Item))
		}
	}
	sb.WriteString("\n财富之道，取之有道：开源、节流、增值、惜福，四者并举，方能财禄绵长。")
	return Section{Key: "wealth_luck", Title: "求财开运", Icon: "star", Body: sb.String()}
}

// stockSector 五行 -> 可关注的投资/板块方向（正面表述）。
var stockSector = map[string]string{
	"金": "金融、银行、证券、贵金属、高端制造、军工",
	"木": "医药、教育、农业、消费、新能源、成长型",
	"水": "物流、航运、传媒、食品饮料、平台消费",
	"火": "科技、电子、半导体、互联网、能源电力",
	"土": "地产基建、资源矿产、消费零售、公用事业",
}