package fortune

// Person 是一次算命输入的基本信息。
type Person struct {
	Name   string `json:"name"`   // 称呼（可为空）
	Gender int    `json:"gender"` // 1=男 0=女
	Year   int    `json:"year"`
	Month  int    `json:"month"`
	Day    int    `json:"day"`
	Hour   int    `json:"hour"` // 24 小时制
	Minute int    `json:"minute"`
	IsLunar bool  `json:"is_lunar"` // 输入是否为农历
}

// Pillar 是天干地支构成的一柱。
type Pillar struct {
	Name       string   `json:"name"`        // 年柱/月柱/日柱/时柱
	GanZhi     string   `json:"gan_zhi"`     // 如 甲子
	Gan        string   `json:"gan"`         // 天干
	Zhi        string   `json:"zhi"`         // 地支
	GanWuXing  string   `json:"gan_wuxing"`  // 天干五行
	ZhiWuXing  string   `json:"zhi_wuxing"`  // 地支五行
	HiddenGan  []string `json:"hidden_gan"`  // 藏干
	ShiShen    string   `json:"shi_shen"`    // 天干十神（对日主）
	NaYin      string   `json:"na_yin"`      // 纳音
}

// DaYun 一步大运。
type DaYun struct {
	StartAge  int    `json:"start_age"`
	EndAge    int    `json:"end_age"`
	StartYear int    `json:"start_year"`
	EndYear   int    `json:"end_year"`
	GanZhi    string `json:"gan_zhi"`
	Gan       string `json:"gan"`
	Zhi       string `json:"zhi"`
}

// LiuNian 一个流年。
type LiuNian struct {
	Year   int    `json:"year"`
	Age    int    `json:"age"`
	GanZhi string `json:"gan_zhi"`
	Gan    string `json:"gan"`
	Zhi    string `json:"zhi"`
}

// Bazi 是排盘结果。
type Bazi struct {
	SolarText     string         `json:"solar_text"`   // 公历
	LunarText     string         `json:"lunar_text"`   // 农历
	GenderText    string         `json:"gender_text"`  // 男/女
	Animal        string         `json:"animal"`       // 生肖
	DayMaster     string         `json:"day_master"`   // 日主天干
	DayMasterWuXing string       `json:"day_master_wuxing"`
	BodyStrength  string         `json:"body_strength"`   // 身强/身弱/中
	XiYong        []string       `json:"xi_yong"`         // 喜用神五行
	JiShen        []string       `json:"ji_shen"`         // 忌神五行
	WuXingCount   map[string]int `json:"wuxing_count"`    // 五行统计
	Pillars       []Pillar       `json:"pillars"`         // 四柱
	DaYun         []DaYun        `json:"da_yun"`          // 大运
	LiuNian       []LiuNian      `json:"liu_nian"`        // 近流年
}

// Score 是某个维度的打分。
type Score struct {
	Value int    `json:"value"` // 0-100
	Grade string `json:"grade"` // 上上/上/中/中下/下
}

// Section 是一段解读内容。
type Section struct {
	Key   string `json:"key"`   // 稳定标识，如 personality / career / health
	Title string `json:"title"`
	Icon  string `json:"icon"`  // 前端图标标识
	Body  string `json:"body"`  // 正文
	Score *Score `json:"score,omitempty"`
}

// MatchResult 是两人（合婚/招聘）匹配结果。
type MatchResult struct {
	Score       Score    `json:"score"`
	Aspects     []Aspect `json:"aspects"` // 各维度
	Conclusion  string   `json:"conclusion"`
	Suggestion  string   `json:"suggestion"`
}

// Aspect 是匹配的一个维度。
type Aspect struct {
	Key     string `json:"key"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Good    bool   `json:"good"` // 该维度是否吉
	Weight  int    `json:"-"`
}

// Result 是一次算命生成并存储的完整结构化结果（存入 Reading.Payload）。
type Result struct {
	Scenario string       `json:"scenario"`
	Title    string       `json:"title"`
	Bazi     *Bazi        `json:"bazi,omitempty"`
	BaziB    *Bazi        `json:"bazi_b,omitempty"` // 合婚/招聘的第二人（可空）
	Match    *MatchResult `json:"match,omitempty"`
	Sections []Section    `json:"sections"`
}