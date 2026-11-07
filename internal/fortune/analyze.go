package fortune

import (
	"fmt"
	"strings"
	"time"
)

// analyzePersonal 基于排盘结果生成个人命理的详细解读段。
func analyzePersonal(b Bazi) []Section {
	return []Section{
		overviewSection(b),
		pillarDetailSection(b),
		rizhuSection(b),
		personalitySection(b),
		animalSection(b),
		careerSection(b),
		wealthSection(b),
		loveSection(b),
		healthSection(b),
		studySection(b),
		shenshaSection(b),
		dayunSection(b),
		liunianSection(b),
		guiRenSection(b),
		lifeStageSection(b),
		luckSection(b),
	}
}

// ---------- 1. 命理总览 ----------

func pillarSummary(b Bazi) string {
	parts := make([]string, 0, 4)
	for _, p := range b.Pillars {
		s := p.Name + p.GanZhi
		if p.ShiShen != "" {
			s += "（" + p.ShiShen + "）"
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, "、")
}

func overviewSection(b Bazi) Section {
	dm := b.DayMaster
	nayin := ""
	if len(b.Pillars) > 0 {
		nayin = b.Pillars[0].NaYin
	}
	order := topWuXing(b.WuXingCount)
	strong, weak := order[0], order[len(order)-1]

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("您出生于%s（%s），日主为「%s」%s日，属%s%s，年命纳音为「%s」。",
		b.SolarText, b.LunarText, dm, b.DayMasterWuXing, b.Animal, b.GenderText, nayin))
	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("四柱排定：%s。年柱主祖上根基、月柱主父母家业、日柱主自身婚姻、时柱主子女晚景。",
		pillarSummary(b)))
	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("命局五行以「%s」偏旺、「%s」偏弱，日主%s；喜用神为%s，忌神为%s。",
		strong, weak, b.BodyStrength, formatWuXingList(b.XiYong), formatWuXingList(b.JiShen)))
	sb.WriteString("\n\n" + overviewComment(b))
	return Section{Key: "overview", Title: "命理总览", Icon: "compass", Body: sb.String()}
}

func overviewComment(b Bazi) string {
	role := ""
	if len(b.Pillars) > 2 {
		role = "「" + b.Pillars[2].GanZhi + "」"
	}
	switch b.BodyStrength {
	case "强":
		return fmt.Sprintf("综合看，您日主有力、身强能任，多能扛事、有主见，宜主动开拓、承担重任。但刚则易折，遇事宜留三分余地，多用喜用神之色与方位调和（婚姻宫在%s，尤为关键），则进退更从容。", role)
	case "弱":
		return "综合看，您日主偏弱，性情多温和、善顺应，贵在借势而行、亲近贵人。宜多与喜用神所对应的人事物相处，养精蓄锐、稳中求进，反而后劲更足。"
	default:
		return "综合看，您日主中和，刚柔相济、进退有度，是难得的平和之命，随遇而安、顺势而为即可得福。"
	}
}

// ---------- 2. 四柱全解 ----------

func pillarDetailSection(b Bazi) Section {
	roles := map[string]string{
		"年柱": "代表祖上根基与早年运势",
		"月柱": "代表父母家业与青年运势",
		"日柱": "代表自身、配偶与中年运势",
		"时柱": "代表子女、晚景与福报",
	}
	var sb strings.Builder
	sb.WriteString("四柱逐柱详批（天干十神 + 地支藏干 + 五行生克）：\n")
	for _, p := range b.Pillars {
		role := roles[p.Name]
		if role == "" {
			role = "命局一柱"
		}
		sb.WriteString(fmt.Sprintf("\n【%s】%s（纳音%s）——%s", p.Name, p.GanZhi, p.NaYin, role))
		if p.ShiShen != "" {
			if m, ok := shiShenMeaning[p.ShiShen]; ok {
				sb.WriteString(fmt.Sprintf("\n  天干「%s」为「%s」：%s", p.Gan, p.ShiShen, m))
			}
		}
		if len(p.HiddenGan) > 0 {
			parts := make([]string, 0, len(p.HiddenGan))
			for _, hg := range p.HiddenGan {
				if ss := shiShen(b.DayMaster, hg); ss != "" {
					parts = append(parts, hg+"（"+ss+"）")
				} else {
					parts = append(parts, hg)
				}
			}
			sb.WriteString("\n  地支藏干：" + strings.Join(parts, "、") + "，暗藏此柱的内在禀赋与心性。")
		}
		if zw, ok := zhiWuXing[p.Zhi]; ok {
			sb.WriteString(fmt.Sprintf("\n  地支属%s，", zw))
			switch {
			case wuXingSheng[zw] == b.DayMasterWuXing:
				sb.WriteString("生扶日主，是该柱对您的暗中助力。")
			case wuXingSheng[b.DayMasterWuXing] == zw:
				sb.WriteString("为日主所生，是您愿意为之付出的领域。")
			case zw == b.DayMasterWuXing:
				sb.WriteString("与日主同类，此领域与您同气相求。")
			default:
				sb.WriteString("与日主相制衡，是历练与成长的所在。")
			}
		}
	}
	sb.WriteString("\n\n四柱如同一棵大树的根、干、花、果，各安其位，则一生盘根错节、枝繁叶茂。")
	return Section{Key: "pillars", Title: "四柱全解", Icon: "gem", Body: sb.String()}
}

// ---------- 3. 性格解析 ----------

func personalitySection(b Bazi) Section {
	dm := b.DayMaster
	dw := b.DayMasterWuXing
	var sb strings.Builder

	if p, ok := ganPersonality[dm]; ok {
		sb.WriteString("日主气质：" + p)
	}
	if p, ok := wuXingPersonality[dw]; ok {
		sb.WriteString("\n\n五行气质：" + p)
	}
	primary := mostCommon(allShiShen(b))
	if primary != "" {
		if m, ok := shiShenMeaning[primary]; ok {
			sb.WriteString(fmt.Sprintf("\n\n命局十神以「%s」较显：%s", primary, m))
		}
	}
	sb.WriteString("\n\n性格优势：" + strengthText(b))
	sb.WriteString("\n\n需留意：" + weaknessText(b))
	if ad, ok := animalDetail[b.Animal]; ok {
		sb.WriteString(fmt.Sprintf("\n\n生肖气质（属%s）：%s", b.Animal, ad))
	}
	return Section{Key: "personality", Title: "性格解析", Icon: "heart", Body: sb.String()}
}

func allShiShen(b Bazi) []string {
	var out []string
	for _, pl := range b.Pillars {
		if pl.ShiShen != "" {
			out = append(out, pl.ShiShen)
		}
		for _, hg := range pl.HiddenGan {
			if s := shiShen(b.DayMaster, hg); s != "" {
				out = append(out, s)
			}
		}
	}
	return out
}

func strengthText(b Bazi) string {
	switch b.BodyStrength {
	case "强":
		return "意志坚定、行动果断，抗压能力强，遇事有担当，是团队里的主心骨。"
	case "弱":
		return "心思细腻、善解人意，适应力与可塑性强，低调踏实，容易赢得他人信赖。"
	default:
		return "刚柔兼济、能进能退，情绪稳定、为人随和，是难得的中和之性。"
	}
}

func weaknessText(b Bazi) string {
	switch b.BodyStrength {
	case "强":
		return "偶尔过于坚持己见、略显强势，宜多听善言、留些弹性，避免因刚致挫。"
	case "弱":
		return "偶尔思虑过多、易受外界影响，宜增强自信、当断则断，借喜用神之力补足。"
	default:
		return "宜保持主见、避免优柔寡断，关键时刻果断一些更佳。"
	}
}

// ---------- 3. 事业财运 ----------

func careerSection(b Bazi) Section {
	var sb strings.Builder
	cai := findShiShenStar(b, []string{"正财", "偏财"})
	if cai != "" {
		sb.WriteString(fmt.Sprintf("财星分析：命局见「%s」，主有财缘、善积累，", cai))
		if cai == "正财" {
			sb.WriteString("宜以稳定职业与踏实经营为正财之本，精打细算、细水长流。")
		} else {
			sb.WriteString("偏财灵动，善于捕捉机会，可适度经营投资，但需守好风险、见好就收。")
		}
	} else {
		sb.WriteString("财星分析：命局财星不显，宜以专业本事立身、靠技术或手艺吃饭，忌投机冒进，积蓄自会渐丰。")
	}

	sb.WriteString("\n\n")
	guan := findShiShenStar(b, []string{"正官", "七杀"})
	if guan != "" {
		if guan == "正官" {
			sb.WriteString("官星分析：命带正官，责任心强、守规矩、重名誉，适合公职、管理、稳定平台，易得上级赏识。")
		} else {
			sb.WriteString("官杀分析：命带七杀，魄力足、敢决断、争强好胜，适合军警、管理、创业等需决断力的领域。")
		}
	} else {
		sb.WriteString("官星分析：官星偏弱，更适合自由职业、技术专长或自主创业，靠本事与口碑立足。")
	}

	sb.WriteString("\n\n适合行业（按喜用神对应）：\n")
	for _, w := range b.XiYong {
		if c, ok := wuXingCareer[w]; ok {
			sb.WriteString(fmt.Sprintf("· %s（%s）：%s\n", w, wuXingAttr[w].Color, c))
		}
	}
	sb.WriteString(fmt.Sprintf("\n助力方位与颜色：喜用方位为%s；日常着衣、办公可用%s。",
		strings.Join(xiYongDirs(b.XiYong), "、"), strings.Join(xiYongColors(b.XiYong), "、")))
	return Section{Key: "career", Title: "事业财运", Icon: "briefcase", Body: sb.String()}
}

// ---------- 4. 情感姻缘 ----------

func loveSection(b Bazi) Section {
	gender := b.GenderText
	starNames := []string{"正财", "偏财"}
	starLabel := "妻星"
	if gender == "女" {
		starNames = []string{"正官", "七杀"}
		starLabel = "夫星"
	}

	var sb strings.Builder
	star := findShiShenStar(b, starNames)
	if star != "" {
		where := ""
		for _, pl := range b.Pillars {
			if pl.ShiShen == star {
				where = pl.Name
				break
			}
			for _, hg := range pl.HiddenGan {
				if shiShen(b.DayMaster, hg) == star {
					where = pl.Name + "藏干"
				}
			}
		}
		sb.WriteString(fmt.Sprintf("配偶星（%s）现于%s，主命中姻缘有定、易得良配。", starLabel+"「"+star+"」", where))
	} else {
		sb.WriteString(fmt.Sprintf("命局中%s不显，姻缘多需机缘与主动寻觅，心态从容、广结善缘则佳偶自来。", starLabel))
	}

	if len(b.Pillars) >= 3 {
		sb.WriteString("\n\n婚姻宫（日支「" + b.Pillars[2].Zhi + "」）：" + zhiRelToDay(b))
	}
	if tao := b.taohua(); tao != "" {
		sb.WriteString(fmt.Sprintf("\n\n命带桃花（%s），异性缘与魅力俱佳；然桃花多者，须以专一为重、远离烂桃花。", tao))
	}
	sb.WriteString("\n\n姻缘建议：感情贵在理解与经营，细水长流；相处中多换位思考，适时给对方空间，则琴瑟和鸣。")
	return Section{Key: "love", Title: "情感姻缘", Icon: "users", Body: sb.String()}
}

func zhiRelToDay(b Bazi) string {
	if len(b.Pillars) < 3 {
		return ""
	}
	dz := b.Pillars[2].Zhi
	dw := b.DayMasterWuXing
	zw := zhiWuXing[dz]
	switch {
	case wuXingSheng[zw] == dw:
		return fmt.Sprintf("日支属%s而生日主，配偶或伴侣多体贴助益，是旺您的缘分。", zw)
	case wuXingSheng[dw] == zw:
		return fmt.Sprintf("日主生其日支%s，您在感情中多主动付出、懂得照顾对方。", zw)
	case zw == dw:
		return fmt.Sprintf("日支与日主同为%s，两人性情相近、相处默契。", zw)
	case wuXingKe[zw] == dw:
		return fmt.Sprintf("日支%s克日主，感情中需多退让包容、避免较劲。", zw)
	default:
		return fmt.Sprintf("日主克日支%s，您在感情里较有主导力，宜多顾及对方感受。", zw)
	}
}

// ---------- 5. 健康养生 ----------

func healthSection(b Bazi) Section {
	var sb strings.Builder
	count := b.WuXingCount
	hasIssue := false
	for _, w := range wuXingOrder {
		if count[w] >= 6 {
			sb.WriteString(fmt.Sprintf("· 五行%s偏旺，%s\n", w, wuXingHealth[w]))
			hasIssue = true
		} else if count[w] <= 2 {
			sb.WriteString(fmt.Sprintf("· 五行%s偏弱，%s\n", w, wuXingHealth[w]))
			hasIssue = true
		}
	}
	if !hasIssue {
		sb.WriteString("· 五行较为均衡，体质总体平和，保持规律作息即可。\n")
	}
	sb.WriteString("\n养生建议：作息规律、饮食清淡、适度运动（散步、太极、八段锦皆宜），")
	sb.WriteString(fmt.Sprintf("并多亲近喜用神方位%s与颜色，有助阴阳调和、气血顺畅。", strings.Join(xiYongDirs(b.XiYong), "、")))
	return Section{Key: "health", Title: "健康养生", Icon: "shield", Body: sb.String()}
}

// ---------- 6. 大运走势 ----------

func dayunSection(b Bazi) Section {
	if len(b.DaYun) == 0 {
		return Section{Key: "dayun", Title: "大运走势", Icon: "trend", Body: "暂无大运数据。"}
	}
	var sb strings.Builder
	qy := b.DaYun[0].StartAge
	sb.WriteString(fmt.Sprintf("您约 %d 岁起运，此后每十年为一步大运，随大运流转而运势起伏。", qy))

	now := time.Now().Year()
	cur := b.DaYun[0]
	for _, d := range b.DaYun {
		if now >= d.StartYear && now <= d.EndYear {
			cur = d
		}
	}
	sb.WriteString(fmt.Sprintf("\n\n当前大运（%d-%d岁，约%d-%d年）为「%s」：%s",
		cur.StartAge, cur.EndAge, cur.StartYear, cur.EndYear, cur.GanZhi, daYunRelationRich(b, cur)))

	cnt := 0
	for _, d := range b.DaYun {
		if d.StartYear > now && cnt < 2 {
			sb.WriteString(fmt.Sprintf("\n· 下一步（%d-%d年）「%s」：%s", d.StartYear, d.EndYear, d.GanZhi, daYunRelationRich(b, d)))
			cnt++
		}
	}
	sb.WriteString("\n\n大运乃人生十年之阶段：行喜用神之运宜放大进取，逢忌神之运宜稳字当头，把握节奏、顺势而为，方能乘风破浪。")
	return Section{Key: "dayun", Title: "大运走势", Icon: "trend", Body: sb.String()}
}

func daYunRelationRich(b Bazi, d DaYun) string {
	var parts []string
	if inWuXingList(b.XiYong, ganWuXing[d.Gan]) {
		parts = append(parts, "天干「"+d.Gan+"」为喜用，事业多得助力")
	} else if inWuXingList(b.JiShen, ganWuXing[d.Gan]) {
		parts = append(parts, "天干「"+d.Gan+"」偏忌，宜稳扎稳打")
	}
	if inWuXingList(b.XiYong, zhiWuXing[d.Zhi]) {
		parts = append(parts, "地支「"+d.Zhi+"」为喜用，根基稳固")
	} else if inWuXingList(b.JiShen, zhiWuXing[d.Zhi]) {
		parts = append(parts, "地支「"+d.Zhi+"」偏忌，注意劳逸结合")
	}
	if len(parts) == 0 {
		return "干支平和，稳步前行即可。"
	}
	return "此步" + strings.Join(parts, "，") + "。"
}

// ---------- 7. 流年运势 ----------

func liunianSection(b Bazi) Section {
	if len(b.LiuNian) == 0 {
		return Section{Key: "liunian", Title: "流年运势", Icon: "calendar", Body: "暂无流年数据。"}
	}
	var sb strings.Builder
	nowYear := time.Now().Year()
	for _, ln := range b.LiuNian {
		if ln.Year == nowYear {
			sb.WriteString(fmt.Sprintf("今年 %d 年（%s，太岁属%s，虚岁%d岁）：%s",
				ln.Year, ln.GanZhi, zhiAnimal[ln.Zhi], ln.Age, liuNianRich(b, ln)))
		}
	}
	sb.WriteString("\n\n往后两年：")
	for _, ln := range b.LiuNian {
		if ln.Year == nowYear+1 || ln.Year == nowYear+2 {
			sb.WriteString(fmt.Sprintf("\n· %d 年（%s，太岁属%s）：%s", ln.Year, ln.GanZhi, zhiAnimal[ln.Zhi], liuNianRich(b, ln)))
		}
	}
	sb.WriteString("\n\n" + liuYueHint(b))
	sb.WriteString("\n\n流年如四季轮转，顺喜用之流年宜进取，逢冲克之流年宜守成；趋吉避凶，贵在顺势。")
	return Section{Key: "liunian", Title: "流年运势", Icon: "calendar", Body: sb.String()}
}

func liuNianRich(b Bazi, ln LiuNian) string {
	var parts []string
	if inWuXingList(b.XiYong, ganWuXing[ln.Gan]) {
		parts = append(parts, "天干「"+ln.Gan+"」为喜用，多有助力")
	} else if inWuXingList(b.JiShen, ganWuXing[ln.Gan]) {
		parts = append(parts, "天干「"+ln.Gan+"」偏忌，宜谨慎行事")
	}
	if inWuXingList(b.XiYong, zhiWuXing[ln.Zhi]) {
		parts = append(parts, "地支「"+ln.Zhi+"」为喜用，根基顺遂")
	} else if inWuXingList(b.JiShen, zhiWuXing[ln.Zhi]) {
		parts = append(parts, "地支「"+ln.Zhi+"」偏忌，注意变数")
	}
	if len(parts) == 0 {
		return "流年总体平稳，宜按部就班、稳中有进。"
	}
	return "此年" + strings.Join(parts, "，") + "。"
}

// ---------- 财运分析 ----------

func wealthSection(b Bazi) Section {
	var sb strings.Builder
	if star := findShiShenStar(b, []string{"正财", "偏财"}); star != "" {
		sb.WriteString(fmt.Sprintf("财星：命带「%s」，", star))
		if star == "正财" {
			sb.WriteString("正财主稳定收入，宜以主业深耕、稳健理财为基，开源节流、细水长流。")
		} else {
			sb.WriteString("偏财主机遇之财，宜多线经营、把握风口，但需严守风险、见好就收。")
		}
	} else {
		sb.WriteString("财星：命局财星不显，财富多靠专业与口碑渐积，宜以技谋生、精进本事，切忌投机冒进。")
	}

	woKe := wuXingKe[b.DayMasterWuXing]
	sb.WriteString(fmt.Sprintf("\n\n求财五行：日主「%s」，财星属%s（%s），宜向%s发展、多经营%s相关领域，事半功倍。",
		b.DayMaster, woKe, wuXingAttr[woKe].Color, wuXingAttr[woKe].Direction, woKe))

	sb.WriteString("\n\n理财心法：宜「开源＋节流＋增值」并举，亲近喜用神之色、方位与伙伴；忌贪高利、忌冲动消费；行喜用大运宜放大格局，逢忌神流年则守成蓄力。")
	return Section{Key: "wealth", Title: "财运分析", Icon: "coins", Body: sb.String()}
}

// ---------- 生肖详解 ----------

func animalSection(b Bazi) Section {
	ani := b.Animal
	ad, _ := animalDetail[ani]
	if ad == "" {
		ad = "生肖为先天气韵，知命而用、扬长避短，福运自增。"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("您属%s：%s", ani, ad))

	if len(b.Pillars) > 0 {
		yz := b.Pillars[0].Zhi
		var friends []string
		if he := zhiLiuHe[yz]; he != "" {
			friends = append(friends, "六合·属"+zhiAnimal[he])
		}
		if grp, ok := sanHeGroup[yz]; ok {
			for _, z := range zhiSanHe[grp] {
				if z != yz {
					friends = append(friends, "三合·属"+zhiAnimal[z])
				}
			}
		}
		if len(friends) > 0 {
			sb.WriteString(fmt.Sprintf("与您最投缘的属相：%s，宜多亲近合作。", strings.Join(friends, "、")))
		}
	}
	sb.WriteString("生肖为先天气韵，知命而用、扬长避短，好运自会常伴。")
	return Section{Key: "animal", Title: "生肖详解", Icon: "paw", Body: sb.String()}
}

// ---------- 流月提示 ----------

func liuYueHint(b Bazi) string {
	monthZhi := []string{"寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥", "子", "丑"}
	names := []string{"正月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "冬月", "腊月"}
	var yi, ji []string
	for i, z := range monthZhi {
		if inWuXingList(b.XiYong, zhiWuXing[z]) {
			yi = append(yi, names[i])
		} else if inWuXingList(b.JiShen, zhiWuXing[z]) {
			ji = append(ji, names[i])
		}
	}
	if len(yi) == 0 && len(ji) == 0 {
		return "逐月运气较为平稳，宜按部就班、稳中有进。"
	}
	return fmt.Sprintf("逐月观之：宜进取的月份多在%s；宜稳守、少动的月份多为%s。", strings.Join(yi, "、"), strings.Join(ji, "、"))
}

// ---------- 开运秘笈 ----------

func luckSection(b Bazi) Section {
	var sb strings.Builder
	sb.WriteString("开运秘笈，助您趋吉避凶、顺风顺水：\n")
	sb.WriteString(fmt.Sprintf("· 幸运颜色：%s\n", strings.Join(xiYongColors(b.XiYong), "、")))
	sb.WriteString(fmt.Sprintf("· 有利方位：%s\n", strings.Join(xiYongDirs(b.XiYong), "、")))
	sb.WriteString(fmt.Sprintf("· 幸运数字：%s\n", strings.Join(luckyNumbers(b.XiYong), "、")))

	sb.WriteString("\n开运珠宝与吉祥物（按喜用神）：\n")
	seen := map[string]bool{}
	for _, w := range b.XiYong {
		if lk, ok := wuXingLucky[w]; ok && !seen[w] {
			seen[w] = true
			sb.WriteString(fmt.Sprintf("· 喜%s：宜%s；%s\n", w, lk.Gem, lk.Item))
		}
	}

	sb.WriteString("\n适合行业方向：\n")
	for _, w := range b.XiYong {
		if c, ok := wuXingCareer[w]; ok {
			sb.WriteString(fmt.Sprintf("· %s相关（%s）：%s\n", w, wuXingAttr[w].Color, c))
		}
	}
	sb.WriteString("\n日常宜保持规律作息、平和心态，多亲近喜用神之人事物，少涉忌神方位（" + strings.Join(xiYongDirs(b.JiShen), "、") + "）。顺势而为，事半功倍。")
	sb.WriteString("\n\n" + blessing)
	return Section{Key: "luck", Title: "开运秘笈", Icon: "star", Body: sb.String()}
}

// ---------- 9. 一生运势脉络 ----------

func lifeStageSection(b Bazi) Section {
	var sb strings.Builder
	sb.WriteString("人生运势随四柱各管一段，各阶段侧重点不同：\n")
	stages := []struct {
		role  string
		label string
	}{
		{"年柱", "少年·根基"},
		{"月柱", "青年·家业"},
		{"日柱", "中年·婚姻事业"},
		{"时柱", "晚年·子女福报"},
	}
	for i, s := range stages {
		if i >= len(b.Pillars) {
			break
		}
		p := b.Pillars[i]
		sb.WriteString(fmt.Sprintf("\n· %s（%s%s）：%s。", s.label, s.role, p.GanZhi, pillarLuck(b, p)))
	}
	sb.WriteString("\n\n知命而后安，年轻时打根基、壮年时主心骨、晚年时享清福，各安其位，则一生从容。")
	return Section{Key: "lifestage", Title: "一生运势脉络", Icon: "trend", Body: sb.String()}
}

func pillarLuck(b Bazi, p Pillar) string {
	gGood := inWuXingList(b.XiYong, ganWuXing[p.Gan]) || inWuXingList(b.XiYong, zhiWuXing[p.Zhi])
	gBad := inWuXingList(b.JiShen, ganWuXing[p.Gan]) || inWuXingList(b.JiShen, zhiWuXing[p.Zhi])
	switch {
	case gGood && !gBad:
		return "此阶段得喜用之气相助，较顺遂，宜放手发展"
	case gBad && !gGood:
		return "此阶段忌神偏显、考验较多，宜稳中求进、广结善缘"
	default:
		return "此阶段平稳中略有起伏，顺其自然即可"
	}
}

// ---------- 10. 学业智慧 ----------

func studySection(b Bazi) Section {
	var sb strings.Builder
	if b.ShiShenIncluded("正印") || b.ShiShenIncluded("偏印") {
		sb.WriteString("命带印星（文昌之源），天资聪颖、悟性高，善读书明理，学业与钻研多有所成，一生不乏学习贵人。")
	} else {
		sb.WriteString("命局印星不显，智慧多来自实践与阅历，宜边做边学、以勤补拙，后劲同样可观。")
	}
	sb.WriteString("\n\n适合深耕的方向（按喜用神对应）：\n")
	for _, w := range b.XiYong {
		if c, ok := wuXingCareer[w]; ok {
			sb.WriteString(fmt.Sprintf("· %s：%s\n", w, c))
		}
	}
	sb.WriteString("\n学问之道贵在坚持，找到兴趣与天赋的交汇点，自能厚积薄发、学有所成。")
	return Section{Key: "study", Title: "学业智慧", Icon: "book", Body: sb.String()}
}

// ---------- 11. 贵人属相 ----------

// 天乙贵人：日干对应的贵人地支（古诀）
var tianYiGuiRen = map[string][]string{
	"甲": {"丑", "未"}, "戊": {"丑", "未"}, "庚": {"丑", "未"},
	"乙": {"子", "申"}, "己": {"子", "申"},
	"丙": {"亥", "酉"}, "丁": {"亥", "酉"},
	"壬": {"卯", "巳"}, "癸": {"卯", "巳"},
	"辛": {"午", "寅"},
}

func guiRenSection(b Bazi) Section {
	var sb strings.Builder
	if zhiList, ok := tianYiGuiRen[b.DayMaster]; ok {
		animals := make([]string, len(zhiList))
		for i, z := range zhiList {
			animals[i] = zhiAnimal[z]
		}
		sb.WriteString(fmt.Sprintf("您的天乙贵人属相为「%s」。多与这些属相之人合作、相处，常有助益、逢凶化吉。", strings.Join(animals, "、")))
	} else {
		sb.WriteString("贵人贵在广结善缘，与人为善者处处有贵人。")
	}
	if len(b.Pillars) > 2 {
		if he := zhiLiuHe[b.Pillars[2].Zhi]; he != "" {
			sb.WriteString(fmt.Sprintf("此外，属%s之人与您有六合之缘，是天然的合作伙伴。", zhiAnimal[he]))
		}
	}
	sb.WriteString("\n\n命带贵人者，遇难有助；但贵人之力需以诚相待、常怀感恩，方能长久相伴。")
	return Section{Key: "guiren", Title: "贵人属相", Icon: "star", Body: sb.String()}
}

// ---------- 工具函数 ----------

func findShiShenStar(b Bazi, names []string) string {
	for _, pl := range b.Pillars {
		for _, n := range names {
			if pl.ShiShen == n {
				return n
			}
			for _, hg := range pl.HiddenGan {
				if shiShen(b.DayMaster, hg) == n {
					return n
				}
			}
		}
	}
	return ""
}

func xiYongColors(xi []string) []string {
	out := make([]string, 0, len(xi))
	for _, w := range xi {
		out = append(out, wuXingAttr[w].Color)
	}
	return dedupe(out)
}

func xiYongDirs(xi []string) []string {
	out := make([]string, 0, len(xi))
	for _, w := range xi {
		out = append(out, wuXingAttr[w].Direction)
	}
	return dedupe(out)
}

func luckyNumbers(xi []string) []string {
	out := make([]string, 0, len(xi))
	for _, w := range xi {
		out = append(out, wuXingNumber[w])
	}
	return out
}

// ---------- 供 matching.go 复用 ----------

func mostCommon(list []string) string {
	cnt := map[string]int{}
	for _, s := range list {
		cnt[s]++
	}
	best, bestN := "", -1
	for k, v := range cnt {
		if v > bestN {
			best, bestN = k, v
		}
	}
	return best
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func inWuXingList(list []string, w string) bool {
	return contains(list, w)
}

func (b Bazi) XiYongContains(names []string) bool {
	return b.ShiShenIncluded(names...)
}

func (b Bazi) ShiShenIncluded(names ...string) bool {
	for _, pl := range b.Pillars {
		if contains(names, pl.ShiShen) {
			return true
		}
		for _, hs := range pl.HiddenGan {
			if contains(names, shiShen(b.DayMaster, hs)) {
				return true
			}
		}
	}
	return false
}

func (b Bazi) taohua() string {
	taohuaSet := map[string]bool{"子": true, "午": true, "卯": true, "酉": true}
	var found []string
	for _, pl := range b.Pillars {
		if taohuaSet[pl.Zhi] {
			found = append(found, pl.Name+"·"+pl.Zhi)
		}
	}
	return strings.Join(found, "、")
}