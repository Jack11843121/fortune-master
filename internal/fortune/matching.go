package fortune

import (
	"fmt"
	"strings"
)

// scoreMarriage 合婚匹配（深度版）：多维度、逐项详析。
func scoreMarriage(a, b Bazi) MatchResult {
	total := 50
	aspects := make([]Aspect, 0, 12)

	// 1. 生肖合冲
	aYz, bYz := a.Pillars[0].Zhi, b.Pillars[0].Zhi
	aAn, bAn := zhiAnimal[aYz], zhiAnimal[bYz]
	switch {
	case zhiLiuHe[aYz] == bYz:
		total += 20
		aspects = append(aspects, Aspect{Key: "animal", Title: "生肖六合", Good: true,
			Body: fmt.Sprintf("二人生肖%s与%s乃「六合」，是上佳的缘分配对。六合者如磁石相吸，性情互补、心有灵犀，相处轻松自在，遇事能同心协力、彼此成就。", aAn, bAn)})
	case sameSanHe(aYz, bYz):
		total += 16
		aspects = append(aspects, Aspect{Key: "animal", Title: "生肖三合", Good: true,
			Body: fmt.Sprintf("二人生肖%s与%s属「三合」，志趣相投、缘定三生。三合者互相扶持、合作默契，越处越顺，是能携手走远的好搭配。", aAn, bAn)})
	case zhiLiuChong[aYz] == bYz:
		total -= 20
		aspects = append(aspects, Aspect{Key: "animal", Title: "生肖相冲", Good: false,
			Body: fmt.Sprintf("二人生肖%s与%s相冲，性格节奏差异明显，一个急一个缓、一动一静，相处中易起摩擦。若能多一份迁就与沟通，冲撞亦可化为互补的张力，反成别样缘分。", aAn, bAn)})
	case zhiXiangHai[aYz] == bYz:
		total -= 16
		aspects = append(aspects, Aspect{Key: "animal", Title: "生肖相害", Good: false,
			Body: fmt.Sprintf("二人生肖%s与%s相害，相处中易因小事生误会、互相猜忌。宜坦诚相待、多交心，把话说开，则嫌隙自消。", aAn, bAn)})
	default:
		total += 6
		aspects = append(aspects, Aspect{Key: "animal", Title: "生肖", Good: true,
			Body: fmt.Sprintf("二人生肖%s与%s无冲无害，相处自然平和，属于不温不火、细水长流的类型。", aAn, bAn)})
	}

	// 2. 纳音对照（年命）
	aN, bN := naYinWuXing(a.Pillars[0].NaYin), naYinWuXing(b.Pillars[0].NaYin)
	if aN != "" && bN != "" {
		switch {
		case aN == bN:
			total += 8
			aspects = append(aspects, Aspect{Key: "nayin", Title: "纳音同气", Good: true,
				Body: fmt.Sprintf("二人年命纳音同为「%s」，先天心性相投、价值观接近，是莫逆之交般的缘分。", aN)})
		case wuXingSheng[aN] == bN:
			total += 8
			aspects = append(aspects, Aspect{Key: "nayin", Title: "纳音相生", Good: true,
				Body: fmt.Sprintf("年命纳音「%s」生「%s」，一方先天滋养另一方，如春风化雨，是旺对方的缘分。", aN, bN)})
		case wuXingSheng[bN] == aN:
			total += 8
			aspects = append(aspects, Aspect{Key: "nayin", Title: "纳音相生", Good: true,
				Body: fmt.Sprintf("年命纳音「%s」生「%s」，彼此滋养、各有所得，缘分深厚。", bN, aN)})
		default:
			total -= 5
			aspects = append(aspects, Aspect{Key: "nayin", Title: "纳音相克", Good: false,
				Body: fmt.Sprintf("年命纳音「%s」与「%s」相克，先天脾性略有碰撞，需后天的理解与磨合。", aN, bN)})
		}
	}

	// 3. 日柱干支关系
	aDay, bDay := a.Pillars[2], b.Pillars[2]
	he := 0
	body := ""
	if ganHe[aDay.Gan] == bDay.Gan {
		he += 12
		body += fmt.Sprintf("日干%s与%s天干相合，心意相通、情投意合；", aDay.Gan, bDay.Gan)
	}
	if zhiLiuHe[aDay.Zhi] == bDay.Zhi {
		he += 16
		body += fmt.Sprintf("日支%s与%s六合，婚姻宫相合，婚后感情稳固、琴瑟和鸣；", aDay.Zhi, bDay.Zhi)
	}
	if zhiLiuChong[aDay.Zhi] == bDay.Zhi {
		he -= 16
		body += fmt.Sprintf("日支%s与%s相冲，日常琐事易起争执；", aDay.Zhi, bDay.Zhi)
	}
	if body == "" {
		he += 5
		body = "日柱无特别合冲，关系平稳，重在用心经营、细水长流。"
	}
	total += he
	aspects = append(aspects, Aspect{Key: "day", Title: "日柱·婚姻宫", Good: he >= 0, Body: body})

	// 4. 五行互补
	dwA, dwB := a.DayMasterWuXing, b.DayMasterWuXing
	switch {
	case wuXingSheng[dwA] == dwB:
		total += 12
		aspects = append(aspects, Aspect{Key: "wuxing", Title: "五行相生", Good: true,
			Body: fmt.Sprintf("日主「%s」生「%s」，一方的情感与付出恰是另一方所需的滋养，天然互补、互为成全。", dwA, dwB)})
	case wuXingSheng[dwB] == dwA:
		total += 12
		aspects = append(aspects, Aspect{Key: "wuxing", Title: "五行相生", Good: true,
			Body: fmt.Sprintf("日主「%s」生「%s」，彼此扶持、各取所需，情感上相互成就。", dwB, dwA)})
	case wuXingKe[dwA] == dwB || wuXingKe[dwB] == dwA:
		total -= 12
		aspects = append(aspects, Aspect{Key: "wuxing", Title: "五行相克", Good: false,
			Body: fmt.Sprintf("日主「%s」与「%s」相克，个性一刚一柔、彼此牵制，是「对方强则我弱」的较量型缘分，需多退让方能长久。", dwA, dwB)})
	default:
		total += 5
		aspects = append(aspects, Aspect{Key: "wuxing", Title: "五行相当", Good: true,
			Body: fmt.Sprintf("日主五行同为「%s」，性情相近、易成知己，是难得的同频之缘。", dwA)})
	}

	// 5. 喜用神互补
	xiBonus := 0
	if inWuXingList(a.XiYong, dwB) {
		xiBonus += 12
	}
	if inWuXingList(b.XiYong, dwA) {
		xiBonus += 12
	}
	if inWuXingList(a.JiShen, dwB) {
		xiBonus -= 8
	}
	if inWuXingList(b.JiShen, dwA) {
		xiBonus -= 8
	}
	total += xiBonus
	gb := "双方喜用神契合度一般，运势上相互助益有限，需靠二人共同经营。"
	if xiBonus >= 16 {
		gb = "双方互为喜用，是「旺夫旺妻」的上等格局，在一起彼此运势都更顺、财与福双收。"
	} else if xiBonus > 0 {
		gb = "一方能为另一方补喜用神，相处中常有「和则两利」的帮扶之效。"
	}
	aspects = append(aspects, Aspect{Key: "xiyong", Title: "喜用互补", Good: xiBonus >= 0, Body: gb})

	// 6. 阴阳调和
	yy := fmt.Sprintf("日主一「%s」一「%s」，阴阳调和，利于家庭和睦、刚柔并济。", a.DayMaster, b.DayMaster)
	if ganYinYang[a.DayMaster] == ganYinYang[b.DayMaster] {
		yy = fmt.Sprintf("日主「%s」「%s」阴阳相同，性情趋近，需留意步调与分工。", a.DayMaster, b.DayMaster)
		total -= 5
	} else {
		total += 5
	}
	aspects = append(aspects, Aspect{Key: "yinyang", Title: "阴阳调和", Good: ganYinYang[a.DayMaster] != ganYinYang[b.DayMaster], Body: yy})

	// 7. 十神互看（彼此在对方命中的角色）
	if ss := shiShen(a.DayMaster, b.DayMaster); ss != "" {
		switch ss {
		case "正官", "七杀":
			total += 6
			aspects = append(aspects, Aspect{Key: "shishen", Title: "命中角色", Good: true,
				Body: fmt.Sprintf("对方日主正是您的「%s」，是能「镇得住、管得住」您的人，也是您心中敬重依靠的归宿。", ss)})
		case "正财", "偏财":
			total += 6
			aspects = append(aspects, Aspect{Key: "shishen", Title: "命中角色", Good: true,
				Body: fmt.Sprintf("对方日主正是您的「%s」，是您甘愿为之付出、也最想守护的人，吸引力极强。", ss)})
		default:
			aspects = append(aspects, Aspect{Key: "shishen", Title: "命中角色", Good: true,
				Body: fmt.Sprintf("对方日主在您命中为「%s」，是亦师亦友、相伴成长的缘分。", ss)})
		}
	}

	// 8. 月令季节契合
	amZ, bmZ := a.Pillars[1].Zhi, b.Pillars[1].Zhi
	if amZ == bmZ || zhiLiuHe[amZ] == bmZ {
		total += 5
		aspects = append(aspects, Aspect{Key: "season", Title: "性情节奏", Good: true,
			Body: fmt.Sprintf("二人月支%s与%s同气或相合，成长节奏、生活习性相近，容易过到一块儿去。", amZ, bmZ)})
	} else {
		aspects = append(aspects, Aspect{Key: "season", Title: "性情节奏", Good: true,
			Body: fmt.Sprintf("二人月支一%s一%s，一快一慢、相得益彰，可在节奏上互相迁就调和。", amZ, bmZ)})
	}

	total = clampScore(total)
	return MatchResult{
		Score:      Score{Value: total, Grade: gradeByScore(total)},
		Aspects:    aspects,
		Conclusion: marriageConclusion(total, aspects),
		Suggestion: marriageSuggestion(total, aspects),
	}
}

// scoreHiring 招聘八字评估（深度版）。
func scoreHiring(candidate Bazi, boss *Bazi, role string) MatchResult {
	aspects := make([]Aspect, 0, 10)
	total := 50
	dw := candidate.DayMasterWuXing

	// 1. 岗位五行契合
	if role == "" {
		role = "综合评估"
	}
	rw := roleWuXing(role)
	if inWuXingList(rw, dw) {
		total += 16
		aspects = append(aspects, Aspect{Key: "role", Title: "岗位五行契合", Good: true,
			Body: fmt.Sprintf("候选人日主「%s」与「%s」岗位所需的%s五行契合，天分与岗位对路，上手快、越做越得心应手。", dw, role, strings.Join(rw, "、"))})
	} else {
		total -= 6
		aspects = append(aspects, Aspect{Key: "role", Title: "岗位五行契合", Good: false,
			Body: fmt.Sprintf("候选人日主「%s」与「%s」岗位的%s五行略疏，需一段适应期，或可借喜用神方位与方式补足。", dw, role, strings.Join(rw, "、"))})
	}

	// 2. 性格画像
	var pbody strings.Builder
	if p, ok := ganPersonality[candidate.DayMaster]; ok {
		pbody.WriteString("日主：" + p)
	}
	if p, ok := wuXingPersonality[dw]; ok {
		pbody.WriteString("\n五行气质：" + p)
	}
	if primary := mostCommon(allShiShen(candidate)); primary != "" {
		if m, ok := shiShenMeaning[primary]; ok {
			pbody.WriteString(fmt.Sprintf("\n十神倾向（%s）：%s", primary, m))
		}
	}
	aspects = append(aspects, Aspect{Key: "personality", Title: "性格画像", Good: true, Body: pbody.String()})

	// 3. 抗压与稳定性
	stable := candidate.BodyStrength != "弱"
	if stable {
		total += 12
		aspects = append(aspects, Aspect{Key: "stable", Title: "抗压稳定", Good: true,
			Body: fmt.Sprintf("日主%s，承压与自我调节能力较好，能扛事、能担当，遇变故不易自乱阵脚。", candidate.BodyStrength)})
	} else {
		aspects = append(aspects, Aspect{Key: "stable", Title: "抗压稳定", Good: false,
			Body: "日主偏弱，宜安排在支持性充足的环境、配以良师引导，不宜长期高压独扛。"})
	}

	// 4. 执行力 / 开拓力
	if candidate.ShiShenIncluded("七杀") || candidate.ShiShenIncluded("劫财") || candidate.ShiShenIncluded("伤官") {
		total += 10
		aspects = append(aspects, Aspect{Key: "drive", Title: "冲劲开拓", Good: true,
			Body: "命带七杀/劫财/伤官，冲劲足、敢打敢拼、不服输，适合攻坚、开拓、竞赛类岗位。"})
	} else {
		total += 6
		aspects = append(aspects, Aspect{Key: "drive", Title: "稳进执行", Good: true,
			Body: "命局温和，执行稳健、服从与配合度高，适合流程化、团队协作、精细化类工作。"})
	}

	// 5. 团队亲和与人际
	if candidate.ShiShenIncluded("食神") || candidate.ShiShenIncluded("正印") || candidate.ShiShenIncluded("比肩") {
		total += 8
		aspects = append(aspects, Aspect{Key: "team", Title: "团队亲和", Good: true,
			Body: "命带食神/正印/比肩，性情和善、乐于合作，人缘与口碑俱佳，是团队里的「润滑油」。"})
	} else {
		aspects = append(aspects, Aspect{Key: "team", Title: "团队亲和", Good: true,
			Body: "命局偏独立、喜独自钻研，宜安排目标清晰的独立任务，或以项目制协作。"})
	}

	// 6. 财富观与忠诚
	if star := findShiShenStar(candidate, []string{"正财", "偏财"}); star != "" {
		total += 6
		aspects = append(aspects, Aspect{Key: "money", Title: "财富观念", Good: true,
			Body: fmt.Sprintf("命带「%s」，对收益敏感、目标感强，适当的激励与晋升通道最能调动其积极性。", star)})
	} else {
		aspects = append(aspects, Aspect{Key: "money", Title: "财富观念", Good: true,
			Body: "财星不显，更重成长与稳定，宜以平台、前景、氛围留人。"})
	}

	// 7. 与上级/团队匹配（可选）
	if boss != nil {
		m := scoreMarriage(candidate, *boss)
		total += (m.Score.Value - 50) / 2
		aspects = append(aspects, Aspect{Key: "boss", Title: "与上级契合", Good: m.Score.Value >= 60,
			Body: fmt.Sprintf("与上级/团队八字匹配度%s（%d分）：%s", m.Score.Grade, m.Score.Value, m.Conclusion)})
	}

	total = clampScore(total)
	return MatchResult{
		Score:      Score{Value: total, Grade: gradeByScore(total)},
		Aspects:    aspects,
		Conclusion: hiringConclusion(total, role, aspects),
		Suggestion: hiringSuggestion(total, candidate),
	}
}

// ---------- 文案 ----------

func naYinWuXing(n string) string {
	r := []rune(n)
	if len(r) == 0 {
		return ""
	}
	return string(r[len(r)-1])
}

func marriageConclusion(v int, aspects []Aspect) string {
	good, bad := []string{}, []string{}
	for _, a := range aspects {
		if a.Good {
			good = append(good, a.Title)
		} else {
			bad = append(bad, a.Title)
		}
	}

	var base string
	switch {
	case v >= 85:
		base = "综观二人命盘，属上等婚配，缘分深厚、五行互补，是天作之合。"
	case v >= 70:
		base = "综观二人命盘，属中上婚配，多数维度契合，是难得的良配。"
	case v >= 55:
		base = "综观二人命盘，属中平婚配，有和合之象亦有磨合之处。"
	case v >= 40:
		base = "综观二人命盘，属中下婚配，性格与节奏差异较大，需更多理解与迁就。"
	default:
		base = "综观二人命盘，匹配度偏低，冲克较多，宜审慎考虑。"
	}

	var sb strings.Builder
	sb.WriteString(base)
	if len(good) > 0 {
		sb.WriteString(fmt.Sprintf("其中以「%s」最为合拍，是这段关系最稳固的根基；", strings.Join(good, "、")))
	}
	if len(bad) > 0 {
		sb.WriteString(fmt.Sprintf("需留意「%s」等相冲之处，是相处中要重点经营的地方。", strings.Join(bad, "、")))
	}
	sb.WriteString("缘分三分天定、七分经营，只要方向相合、彼此珍惜，皆可修成正果。")
	return sb.String()
}

func marriageSuggestion(v int, aspects []Aspect) string {
	var sb strings.Builder
	if v >= 70 {
		sb.WriteString("建议择喜用神吉日办事，多借对方喜用神之色、方位与职业，彼此助旺；")
	} else {
		sb.WriteString("若仍愿携手，建议多沟通、明确分工，重大分歧时请教专业人士，择吉日化解冲克；")
	}
	sb.WriteString("日常相处中，男重担当、女重体谅，遇事把话说开、莫积怨。")
	sb.WriteString("感情如种花，常浇水施肥，方能常开常艳。")
	return sb.String()
}

func hiringConclusion(v int, role string, aspects []Aspect) string {
	switch {
	case v >= 80:
		return fmt.Sprintf("综合评估，此人与「%s」岗位高度契合，能力、性格与团队匹配俱佳，是重点考虑的人选；建议尽快锁定、委以重任。", role)
	case v >= 65:
		return fmt.Sprintf("综合评估，此人与「%s」岗位基本契合，可胜任岗位要求，入职后注意扬长避短、配以合适的导师与机制。", role)
	case v >= 50:
		return fmt.Sprintf("综合评估，此人与「%s」岗位契合一般，可作为备选，或调整至更匹配的岗位方向观察。", role)
	default:
		return "综合评估，此人与目标岗位契合度偏低，建议转岗、谨慎录用，或以试用期充分考察。"
	}
}

func hiringSuggestion(v int, c Bazi) string {
	var sb strings.Builder
	if v >= 65 {
		sb.WriteString("可根据候选人的喜用神安排办公方位、职责与激励方式，助其快速进入状态、产出绩效。")
	} else {
		sb.WriteString("建议以结构化面试结合试用期观察为主，八字分析仅作性格与潜质的参考维度之一。")
	}
	if c.BodyStrength == "弱" {
		sb.WriteString("其人宜在有支持、有引导的环境下成长，避免一上来就压担子。")
	} else {
		sb.WriteString("其人可适度授权、委以挑战，激发其担当与潜力。")
	}
	return sb.String()
}

func sameSanHe(a, b string) bool {
	for _, grp := range zhiSanHe {
		if contains(grp, a) && contains(grp, b) {
			return true
		}
	}
	return false
}

func clampScore(v int) int {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}

func roleWuXing(role string) []string {
	switch strings.TrimSpace(role) {
	case "技术研发":
		return []string{"金", "水"}
	case "销售市场":
		return []string{"火", "木"}
	case "管理领导":
		return []string{"土", "金"}
	case "财务审计":
		return []string{"金", "土"}
	case "创意设计":
		return []string{"木", "水", "火"}
	case "行政人事":
		return []string{"土", "金"}
	case "运营执行":
		return []string{"土", "火"}
	default:
		return []string{"土", "金", "水"}
	}
}

// Roles 前端可选的岗位类型
var Roles = []string{"技术研发", "销售市场", "管理领导", "财务审计", "创意设计", "行政人事", "运营执行"}