package fortune

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// 宝宝起名：依据八字喜用神，从吉祥常用字库中组合姓名，附拼音、五行、寓意与评分。

// 常见姓氏（含复姓）拼音
var surnamePinyin = map[string]string{
	"李": "lǐ", "王": "wáng", "张": "zhāng", "刘": "liú", "陈": "chén", "杨": "yáng", "黄": "huáng",
	"赵": "zhào", "吴": "wú", "周": "zhōu", "徐": "xú", "孙": "sūn", "马": "mǎ", "朱": "zhū",
	"胡": "hú", "郭": "guō", "何": "hé", "高": "gāo", "林": "lín", "罗": "luó", "郑": "zhèng",
	"梁": "liáng", "谢": "xiè", "宋": "sòng", "唐": "táng", "许": "xǔ", "韩": "hán", "冯": "féng",
	"邓": "dèng", "曹": "cáo", "彭": "péng", "曾": "zēng", "苏": "sū", "卢": "lú", "蒋": "jiǎng",
	"蔡": "cài", "贾": "jiǎ", "丁": "dīng", "魏": "wèi", "薛": "xuē", "叶": "yè", "阎": "yán",
	"余": "yú", "潘": "pān", "杜": "dù", "戴": "dài", "夏": "xià", "钟": "zhōng", "汪": "wāng",
	"田": "tián", "任": "rén", "姜": "jiāng", "范": "fàn", "方": "fāng", "石": "shí", "姚": "yáo",
	"谭": "tán", "廖": "liào", "邹": "zōu", "熊": "xióng", "金": "jīn", "陆": "lù", "郝": "hǎo",
	"欧阳": "ōuyáng", "司马": "sīmǎ", "上官": "shàngguān", "诸葛": "zhūgě", "令狐": "línghú", "东方": "dōngfāng", "皇甫": "huángfǔ", "慕容": "mùróng",
}

type nameChar struct {
	char    string
	pinyin  string
	wuxing  string
	gender  string // 男 / 女 / 通
	meaning string
}

// 起名常用吉祥字库（按五行归类，寓意正面）
var nameChars = []nameChar{
	// 金
	{"鑫", "xīn", "金", "通", "金多兴旺，财源广进"}, {"铭", "míng", "金", "通", "铭记于心，才华出众"},
	{"钰", "yù", "金", "通", "珍宝美玉，珍贵不凡"}, {"锋", "fēng", "金", "男", "锐意进取，锋芒毕露"},
	{"锐", "ruì", "金", "男", "敏锐精明，锐不可当"}, {"锦", "jǐn", "金", "通", "前程似锦，华美绚烂"},
	{"钧", "jūn", "金", "男", "千钧之重，气度不凡"}, {"铮", "zhēng", "金", "男", "铮铮铁骨，正直刚强"},
	{"钦", "qīn", "金", "男", "令人钦佩，品德高尚"}, {"琛", "chēn", "金", "通", "如珠如宝，才华内蕴"},
	{"玥", "yuè", "金", "女", "掌上明珠，美好珍贵"}, {"珊", "shān", "金", "女", "珊瑚之美，温润贵重"},
	{"瑞", "ruì", "金", "通", "祥瑞之兆，福气满盈"}, {"铃", "líng", "金", "女", "铃声清脆，灵动美好"},
	{"皓", "hào", "金", "通", "洁白明亮，光明磊落"}, {"晨", "chén", "金", "通", "晨光清朗，朝气蓬勃"},
	{"靖", "jìng", "金", "男", "安定从容，稳重大方"}, {"舒", "shū", "金", "通", "舒展自在，从容优雅"},
	// 木
	{"林", "lín", "木", "通", "林木葱郁，生机盎然"}, {"森", "sēn", "木", "男", "森林繁茂，气度宽宏"},
	{"楠", "nán", "木", "通", "楠木良才，坚实贵重"}, {"桐", "tóng", "木", "通", "梧桐引凤，高洁不凡"},
	{"梓", "zǐ", "木", "通", "桑梓之情，才德兼备"}, {"彬", "bīn", "木", "男", "文质彬彬，温文尔雅"},
	{"杰", "jié", "木", "通", "人中豪杰，出类拔萃"}, {"栋", "dòng", "木", "男", "栋梁之才，可堪大任"},
	{"楚", "chǔ", "木", "通", "清楚脱俗，楚楚不凡"}, {"茉", "mò", "木", "女", "茉莉清香，纯洁美好"},
	{"莉", "lì", "木", "女", "茉莉芬芳，温婉可人"}, {"萱", "xuān", "木", "女", "萱草忘忧，快乐安康"},
	{"芳", "fāng", "木", "女", "芬芳馥郁，德艺双馨"}, {"茗", "míng", "木", "通", "香茗清雅，气韵悠长"},
	{"蔚", "wèi", "木", "通", "蔚然成风，文采斐然"}, {"菡", "hàn", "木", "女", "菡萏初开，清新脱俗"},
	{"柏", "bǎi", "木", "男", "松柏长青，坚贞不屈"}, {"晴", "qíng", "木", "女", "晴空万里，明朗欢快"},
	// 水
	{"涵", "hán", "水", "通", "涵养深厚，包容大度"}, {"浩", "hào", "水", "男", "浩然正气，胸襟宽广"},
	{"洋", "yáng", "水", "男", "洋洋大观，豁达开朗"}, {"涛", "tāo", "水", "男", "波涛澎湃，气势磅礴"},
	{"泽", "zé", "水", "通", "泽被四方，恩泽绵长"}, {"沐", "mù", "水", "通", "如沐春风，温润和煦"},
	{"汐", "xī", "水", "女", "潮汐之美，灵动温柔"}, {"沁", "qìn", "水", "女", "沁人心脾，清新怡人"},
	{"清", "qīng", "水", "通", "清正廉洁，清澈纯粹"}, {"洁", "jié", "水", "女", "纯洁无瑕，品性高洁"},
	{"淑", "shū", "水", "女", "贤淑温婉，德容兼备"}, {"源", "yuán", "水", "通", "饮水思源，源远流长"},
	{"泓", "hóng", "水", "通", "水深而广，气度沉稳"}, {"淼", "miǎo", "水", "女", "水势浩淼，心胸开阔"},
	{"澜", "lán", "水", "通", "波澜壮阔，气魄宏大"}, {"潇", "xiāo", "水", "男", "潇洒自如，风度翩翩"},
	{"淇", "qí", "水", "女", "淇水汤汤，温婉明净"}, {"雨", "yǔ", "水", "通", "润物无声，恩泽绵绵"},
	// 火
	{"烨", "yè", "火", "通", "光辉灿烂，才华闪耀"}, {"煜", "yù", "火", "通", "煜煜生辉，光耀门楣"},
	{"炜", "wěi", "火", "男", "光明显赫，气宇轩昂"}, {"炫", "xuàn", "火", "男", "光彩炫目，才华横溢"},
	{"烁", "shuò", "火", "通", "闪烁夺目，熠熠生辉"}, {"煌", "huáng", "火", "男", "辉煌灿烂，前程似锦"},
	{"灿", "càn", "火", "通", "灿若星辰，光明美好"}, {"灵", "líng", "火", "通", "灵秀聪慧，一点即通"},
	{"曜", "yào", "火", "男", "日月同辉，光耀乾坤"}, {"煊", "xuān", "火", "男", "煊赫有声，声名远播"},
	{"焕", "huàn", "火", "通", "焕然一新，精神焕发"}, {"曦", "xī", "火", "通", "晨曦初照，希望蓬勃"},
	{"昕", "xīn", "火", "女", "旭日初升，明亮美好"}, {"明", "míng", "火", "通", "光明磊落，智慧通达"},
	{"昀", "yún", "火", "通", "阳光和煦，温和明亮"}, {"晴", "qíng", "火", "女", "风和日丽，明朗欢快"},
	// 土
	{"坤", "kūn", "土", "通", "厚德载物，包容大度"}, {"培", "péi", "土", "男", "培育英才，根基厚实"},
	{"城", "chéng", "土", "男", "固若金城，稳重可靠"}, {"垚", "yáo", "土", "男", "土高成山，志存高远"},
	{"辰", "chén", "土", "通", "星辰大海，胸怀宽广"}, {"均", "jūn", "土", "通", "均衡平和，公正无私"},
	{"圣", "shèng", "土", "通", "圣明睿智，德行高尚"}, {"堂", "táng", "土", "男", "堂堂正正，光明磊落"},
	{"墨", "mò", "土", "通", "翰墨书香，才华沉淀"}, {"境", "jìng", "土", "通", "心境高远，从容自若"},
	{"垣", "yuán", "土", "男", "根基稳固，踏实可靠"}, {"恩", "ēn", "土", "通", "知恩感恩，宅心仁厚"},
	{"岳", "yuè", "土", "男", "高山仰止，稳重如山"}, {"岩", "yán", "土", "男", "坚若磐石，意志坚定"},
	{"勋", "xūn", "土", "通", "功勋卓著，有所作为"}, {"安", "ān", "土", "通", "平安喜乐，安稳从容"},
}

// NameSuggestion 单个起名建议
type NameSuggestion struct {
	Name    string `json:"name"`
	Pinyin  string `json:"pinyin"`
	Wuxing  string `json:"wuxing"`
	Meaning string `json:"meaning"` // 拆字寓意
	Reason  string `json:"reason"`  // 为什么推荐（命理依据）
	Score   int    `json:"score"`
}

// SuggestNames 依八字喜用神生成起名建议；gender 1=男 0=女；nameLen 为名字字数（1 单字名 / 2 双字名）。
func SuggestNames(surname string, gender int, b *Bazi, nameLen, count int) []NameSuggestion {
	if count <= 0 {
		count = 40
	}
	if nameLen < 1 || nameLen > 2 {
		nameLen = 2
	}
	surname = strings.TrimSpace(surname)
	if surname == "" {
		surname = "王"
	}
	xi, ji := b.XiYong, b.JiShen
	dm, dw := b.DayMaster, b.DayMasterWuXing
	order := topWuXing(b.WuXingCount)
	weak := order[len(order)-1]

	genderOK := func(c nameChar) bool {
		if gender == 1 {
			return c.gender == "男" || c.gender == "通"
		}
		return c.gender == "女" || c.gender == "通"
	}

	primary := []nameChar{}
	secondary := []nameChar{}
	for _, c := range nameChars {
		if !genderOK(c) {
			continue
		}
		if inWuXingList(xi, c.wuxing) {
			primary = append(primary, c)
		} else if !inWuXingList(ji, c.wuxing) {
			secondary = append(secondary, c)
		}
	}
	if len(primary) == 0 {
		primary = secondary
	}
	if len(primary) == 0 {
		for _, c := range nameChars {
			if genderOK(c) {
				primary = append(primary, c)
			}
		}
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	pick := func(pool []nameChar) nameChar {
		return pool[rng.Intn(len(pool))]
	}
	seen := map[string]bool{}
	py := surname
	if p, ok := surnamePinyin[surname]; ok {
		py = p
	}

	var out []NameSuggestion
	for attempt := 0; len(out) < count && attempt < 600; attempt++ {
		c1 := pick(primary)

		// 单字名（姓 + 1 字）
		if nameLen == 1 {
			name := surname + c1.char
			if seen[name] {
				continue
			}
			seen[name] = true
			score := 60
			if inWuXingList(xi, c1.wuxing) {
				score += 22
			}
			if inWuXingList(ji, c1.wuxing) {
				score -= 12
			}
			reason := fmt.Sprintf("宝宝日主「%s」属%s，命局喜用%s（%s偏弱，宜补）。单字名取「%s」属%s，字简意深、朗朗上口，既补喜用所需，又勾勒清朗美好的意韵，助宝宝五行调衡、顺遂成长。",
				dm, dw, formatWuXingList(xi), weak, c1.char, c1.wuxing)
			out = append(out, NameSuggestion{
				Name:    name,
				Pinyin:  py + " " + c1.pinyin,
				Wuxing:  c1.wuxing,
				Meaning: "「" + c1.char + "」" + c1.meaning + "。",
				Reason:  reason,
				Score:   score,
			})
			continue
		}

		// 双字名（姓 + 2 字）
		secondPool := secondary
		if len(secondPool) == 0 {
			secondPool = primary
		}
		c2 := pick(secondPool)
		if c1.char == c2.char || c1.pinyin == c2.pinyin {
			continue
		}
		name := surname + c1.char + c2.char
		if seen[name] {
			continue
		}
		seen[name] = true
		score := 60
		if inWuXingList(xi, c1.wuxing) {
			score += 20
		}
		if inWuXingList(xi, c2.wuxing) {
			score += 15
		}
		if inWuXingList(ji, c1.wuxing) {
			score -= 10
		}
		if inWuXingList(ji, c2.wuxing) {
			score -= 5
		}
		reason := fmt.Sprintf("宝宝日主「%s」属%s，命局喜用%s（其中%s偏弱，宜补）。此名取「%s」「%s」二字，五行属%s、%s，恰补喜用，助五行调衡、趋吉避凶；且二字寓意美好、音形俱佳，故此名可扶助宝宝健康成长、顺遂如意。",
			dm, dw, formatWuXingList(xi), weak, c1.char, c2.char, c1.wuxing, c2.wuxing)
		out = append(out, NameSuggestion{
			Name:    name,
			Pinyin:  py + " " + c1.pinyin + " " + c2.pinyin,
			Wuxing:  c1.wuxing + "+" + c2.wuxing,
			Meaning: "「" + c1.char + "」" + c1.meaning + "；「" + c2.char + "」" + c2.meaning + "。",
			Reason:  reason,
			Score:   score,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}