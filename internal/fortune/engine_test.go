package fortune

import "testing"

func TestBuildBazi(t *testing.T) {
	b, err := BuildBazi(Person{Gender: 1, Year: 1990, Month: 6, Day: 15, Hour: 10, Minute: 30})
	if err != nil {
		t.Fatal(err)
	}
	if b.DayMaster == "" {
		t.Fatal("日主为空")
	}
	if len(b.Pillars) != 4 {
		t.Fatalf("四柱数量=%d，期望 4", len(b.Pillars))
	}
	if len(b.DaYun) == 0 {
		t.Fatal("大运为空")
	}
	for _, d := range b.DaYun {
		if d.GanZhi == "" {
			t.Fatal("大运存在空干支，应被过滤")
		}
	}
	if len(b.WuXingCount) != 5 {
		t.Fatalf("五行统计项数=%d，期望 5", len(b.WuXingCount))
	}
	t.Logf("日主=%s（%s） 身%s 生肖=%s 性别=%s", b.DayMaster, b.DayMasterWuXing, b.BodyStrength, b.Animal, b.GenderText)
	for _, p := range b.Pillars {
		t.Logf("%s %s 干%s(十神%s) 支%s 藏干=%v 纳音=%s", p.Name, p.GanZhi, p.Gan, p.ShiShen, p.Zhi, p.HiddenGan, p.NaYin)
	}
	t.Logf("喜用=%v 忌=%v 五行=%v", b.XiYong, b.JiShen, b.WuXingCount)
	t.Logf("流年=%+v", b.LiuNian)
}

func TestPersonalResult(t *testing.T) {
	e := NewEngine()
	res, err := e.Personal(Person{Gender: 0, Year: 1988, Month: 11, Day: 2, Hour: 8, Minute: 0})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Sections) < 5 {
		t.Fatalf("解读段数=%d，期望 >=5", len(res.Sections))
	}
	for _, s := range res.Sections {
		if s.Title == "" || s.Body == "" {
			t.Fatalf("段 %q 内容为空", s.Key)
		}
		t.Logf("[%s] %s -> %s", s.Key, s.Title, truncate(s.Body, 60))
	}
}

func TestMarriageResult(t *testing.T) {
	e := NewEngine()
	res, err := e.Marriage(
		Person{Gender: 1, Year: 1992, Month: 3, Day: 5, Hour: 12},
		Person{Gender: 0, Year: 1994, Month: 7, Day: 9, Hour: 14},
	)
	if err != nil {
		t.Fatal(err)
	}
	if res.Match == nil {
		t.Fatal("匹配结果为空")
	}
	t.Logf("合婚评分=%d(%s)", res.Match.Score.Value, res.Match.Score.Grade)
	for _, a := range res.Match.Aspects {
		t.Logf("  %s: %s", a.Title, truncate(a.Body, 50))
	}
	t.Logf("结论=%q", res.Match.Conclusion)
}

func TestHiringResult(t *testing.T) {
	e := NewEngine()
	res, err := e.Hiring(
		Person{Gender: 1, Year: 1990, Month: 1, Day: 1, Hour: 10},
		&Person{Gender: 1, Year: 1978, Month: 5, Day: 20, Hour: 9},
		"技术研发",
	)
	if err != nil {
		t.Fatal(err)
	}
	if res.Match == nil {
		t.Fatal("匹配结果为空")
	}
	t.Logf("招聘评分=%d(%s) 结论=%q", res.Match.Score.Value, res.Match.Score.Grade, res.Match.Conclusion)
}

func TestInvalidSolarDate(t *testing.T) {
	if _, err := BuildBazi(Person{Year: 2025, Month: 2, Day: 30}); err == nil {
		t.Fatal("预期 2月30日 应返回错误")
	}
	if _, err := BuildBazi(Person{Year: 2025, Month: 4, Day: 31}); err == nil {
		t.Fatal("预期 4月31日 应返回错误")
	}
	if _, err := BuildBazi(Person{Year: 2024, Month: 2, Day: 29, Gender: 1}); err != nil {
		t.Fatalf("闰年2月29日应合法，却返回错误: %v", err)
	}
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}