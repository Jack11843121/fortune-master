package fortune

import "fmt"

// Engine 负责各算命场景的编排（排盘 + 解读 + 匹配）。
type Engine struct{}

func NewEngine() *Engine { return &Engine{} }

// Personal 个人本命命理。
func (e *Engine) Personal(p Person) (*Result, error) {
	if err := validate(p); err != nil {
		return nil, err
	}
	b, err := BuildBazi(p)
	if err != nil {
		return nil, err
	}
	return &Result{
		Scenario: "personal",
		Title:    fmt.Sprintf("%s · 个人命理", nameOr(p.Name, "命主")),
		Bazi:     b,
		Sections: analyzePersonal(*b),
	}, nil
}

// Marriage 两人合婚匹配。
func (e *Engine) Marriage(a, b Person) (*Result, error) {
	if err := validate(a); err != nil {
		return nil, fmt.Errorf("甲方信息有误：%w", err)
	}
	if err := validate(b); err != nil {
		return nil, fmt.Errorf("乙方信息有误：%w", err)
	}
	ba, err := BuildBazi(a)
	if err != nil {
		return nil, err
	}
	bb, err := BuildBazi(b)
	if err != nil {
		return nil, err
	}
	m := scoreMarriage(*ba, *bb)
	return &Result{
		Scenario: "marriage",
		Title:    fmt.Sprintf("%s 与 %s · 合婚匹配", nameOr(a.Name, "甲方"), nameOr(b.Name, "乙方")),
		Bazi:     ba,
		BaziB:    bb,
		Match:    &m,
		Sections: matchIntroSections(ba, bb),
	}, nil
}

// Hiring 招聘八字评估（可选提供老板/团队八字做匹配）。
func (e *Engine) Hiring(candidate Person, boss *Person, role string) (*Result, error) {
	if err := validate(candidate); err != nil {
		return nil, fmt.Errorf("候选人信息有误：%w", err)
	}
	if role == "" {
		role = "综合评估"
	}
	bc, err := BuildBazi(candidate)
	if err != nil {
		return nil, err
	}
	var bb *Bazi
	if boss != nil && boss.Year > 0 {
		if err := validate(*boss); err == nil {
			if b, err := BuildBazi(*boss); err == nil {
				bb = b
			}
		}
	}
	m := scoreHiring(*bc, bb, role)
	return &Result{
		Scenario: "hiring",
		Title:    fmt.Sprintf("%s · 招聘八字评估（%s）", nameOr(candidate.Name, "候选人"), role),
		Bazi:     bc,
		BaziB:    bb,
		Match:    &m,
		Sections: matchIntroSections(bc, bb),
	}, nil
}

// matchIntroSections 为合婚/招聘生成可免费看的前置概览段。
func matchIntroSections(a, b *Bazi) []Section {
	if a == nil {
		return nil
	}
	size := "单人"
	if b != nil {
		size = "双人"
	}
	return []Section{{
		Key:   "intro",
		Title: size + "命盘速览",
		Icon:  "compass",
		Body:  fmt.Sprintf("此为您二人（或人选）命盘速览：先看四柱与五行分布，详解维度、综合结论与化解建议见下方解锁内容。"),
	}}
}

func nameOr(name, def string) string {
	if name != "" {
		return name
	}
	return def
}

// validate 基本校验输入范围，避免传给历法库引发异常。
func validate(p Person) error {
	if p.Year < 1900 || p.Year > 2100 {
		return fmt.Errorf("出生年份需在 1900-2100 之间")
	}
	if p.Month < 1 || p.Month > 12 {
		return fmt.Errorf("月份需在 1-12 之间")
	}
	if p.Day < 1 || p.Day > 31 {
		return fmt.Errorf("日期需在 1-31 之间")
	}
	if p.Hour < 0 || p.Hour > 23 {
		return fmt.Errorf("时辰需在 0-23 之间")
	}
	if p.Minute < 0 || p.Minute > 59 {
		return fmt.Errorf("分钟需在 0-59 之间")
	}
	if p.Gender != 0 && p.Gender != 1 {
		return fmt.Errorf("性别取值需为 1（男）或 0（女）")
	}
	return nil
}