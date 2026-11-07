package fortune

import (
	"fmt"
	"sort"
	"time"

	"github.com/6tail/lunar-go/calendar"
)

// BuildBazi 依据出生信息排八字盘。基于 lunar-go（6tail）实现，节气/干支/藏干/纳音/大运均准确。
func BuildBazi(p Person) (*Bazi, error) {
	// 公历日期合法性校验，避免非法日期（如 2月30日）被历法库静默“进位”成错误八字
	if !p.IsLunar {
		t := time.Date(p.Year, time.Month(p.Month), p.Day, 0, 0, 0, 0, time.UTC)
		if t.Year() != p.Year || int(t.Month()) != p.Month || t.Day() != p.Day {
			return nil, fmt.Errorf("出生日期无效：%d年%d月%d日不存在", p.Year, p.Month, p.Day)
		}
	}

	var lunar *calendar.Lunar
	var solar *calendar.Solar

	if p.IsLunar {
		// lunarMonth 可为负表示闰月
		lunar = calendar.NewLunar(p.Year, p.Month, p.Day, p.Hour, p.Minute, 0)
		solar = lunar.GetSolar()
	} else {
		solar = calendar.NewSolar(p.Year, p.Month, p.Day, p.Hour, p.Minute, 0)
		lunar = solar.GetLunar()
	}

	ec := lunar.GetEightChar()

	dayMaster := ec.GetDayGan()
	pillars := []Pillar{
		buildPillar("年柱", ec.GetYear(), ec.GetYearGan(), ec.GetYearZhi(), ec.GetYearHideGan(), ec.GetYearNaYin(), ec.GetYearShiShenGan()),
		buildPillar("月柱", ec.GetMonth(), ec.GetMonthGan(), ec.GetMonthZhi(), ec.GetMonthHideGan(), ec.GetMonthNaYin(), ec.GetMonthShiShenGan()),
		buildPillar("日柱", ec.GetDay(), ec.GetDayGan(), ec.GetDayZhi(), ec.GetDayHideGan(), ec.GetDayNaYin(), ec.GetDayShiShenGan()),
		buildPillar("时柱", ec.GetTime(), ec.GetTimeGan(), ec.GetTimeZhi(), ec.GetTimeHideGan(), ec.GetTimeNaYin(), ec.GetTimeShiShenGan()),
	}

	count := countWuXing(pillars)
	monthZhi := ec.GetMonthZhi()
	gender := "女"
	if p.Gender == 1 {
		gender = "男"
	}

	xi, ji := xiYong(dayMaster, monthZhi, count)

	yun := ec.GetYun(p.Gender)

	b := &Bazi{
		SolarText:       fmt.Sprintf("%d年%d月%d日 %02d:%02d", solar.GetYear(), solar.GetMonth(), solar.GetDay(), solar.GetHour(), solar.GetMinute()),
		LunarText:       "农历" + lunar.String(),
		GenderText:      gender,
		Animal:          lunar.GetYearShengXiao(),
		DayMaster:       dayMaster,
		DayMasterWuXing: ganWuXing[dayMaster],
		BodyStrength:    bodyStrength(dayMaster, monthZhi, count),
		XiYong:          xi,
		JiShen:          ji,
		WuXingCount:     count,
		Pillars:         pillars,
		DaYun:           buildDaYun(yun),
		LiuNian:         buildLiuNian(yun),
	}
	return b, nil
}

func buildPillar(name, ganZhi, gan, zhi string, hidden []string, naYin, shiShen string) Pillar {
	return Pillar{
		Name:      name,
		GanZhi:    ganZhi,
		Gan:       gan,
		Zhi:       zhi,
		GanWuXing: ganWuXing[gan],
		ZhiWuXing: zhiWuXing[zhi],
		HiddenGan: hidden,
		ShiShen:   shiShen,
		NaYin:     naYin,
	}
}

func buildDaYun(yun *calendar.Yun) []DaYun {
	dayunList := yun.GetDaYun()
	out := make([]DaYun, 0, len(dayunList))
	for _, d := range dayunList {
		gz := d.GetGanZhi()
		if gz == "" {
			// 起运前的占位段无干支，跳过（避免显示空干支）
			continue
		}
		gan, zhi := splitGanZhi(gz)
		out = append(out, DaYun{
			StartAge:  d.GetStartAge(),
			EndAge:    d.GetEndAge(),
			StartYear: d.GetStartYear(),
			EndYear:   d.GetEndYear(),
			GanZhi:    gz,
			Gan:       gan,
			Zhi:       zhi,
		})
	}
	return out
}

// buildLiuNian 从各步大运中收集当前年前后若干年的流年（带准确虚岁）。
func buildLiuNian(yun *calendar.Yun) []LiuNian {
	now := time.Now().Year()
	seen := map[int]bool{}
	out := make([]LiuNian, 0, 12)
	for _, d := range yun.GetDaYun() {
		for _, ln := range d.GetLiuNian() {
			y := ln.GetYear()
			if y >= now-1 && y <= now+5 && !seen[y] {
				seen[y] = true
				gan, zhi := splitGanZhi(ln.GetGanZhi())
				out = append(out, LiuNian{Year: y, Age: ln.GetAge(), GanZhi: ln.GetGanZhi(), Gan: gan, Zhi: zhi})
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Year < out[j].Year })
	return out
}

func splitGanZhi(ganZhi string) (string, string) {
	r := []rune(ganZhi)
	if len(r) < 2 {
		return ganZhi, ""
	}
	return string(r[0]), string(r[1])
}