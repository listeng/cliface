package main

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestValidateRequired(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items: []Item{
			{Name: "required_field", Type: "string", Label: "Required Field", Required: true},
			{Name: "optional_field", Type: "string"},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// 必填字段为空
	err := ui.validateRequired()
	if err == nil {
		t.Error("validateRequired() = nil, want error")
	}
	if err.Error() != "Required Field is required" {
		t.Errorf("error = %q, want %q", err.Error(), "Required Field is required")
	}

	// 填写必填字段
	setEntryText(ui.widgets["required_field"], "value")
	err = ui.validateRequired()
	if err != nil {
		t.Errorf("validateRequired() = %v, want nil", err)
	}
}

func TestValidateRequiredUsesNameIfNoLabel(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "myfield", Type: "string", Required: true}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	err := ui.validateRequired()
	if err == nil || err.Error() != "myfield is required" {
		t.Errorf("error = %v, want 'myfield is required'", err)
	}
}

func TestValidateRegex(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "email", Type: "string", Label: "Email", Validate: `^[a-z]+@[a-z]+\.[a-z]+$`}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// 无效格式
	setEntryText(ui.widgets["email"], "invalid")
	err := ui.validateAll()
	if err == nil {
		t.Error("validateAll() = nil, want error")
	}
	if err.Error() != "Email format invalid" {
		t.Errorf("error = %q, want %q", err.Error(), "Email format invalid")
	}

	// 有效格式
	setEntryText(ui.widgets["email"], "test@example.com")
	err = ui.validateAll()
	if err != nil {
		t.Errorf("validateAll() = %v, want nil", err)
	}
}

func TestValidateRegexEmpty(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "field", Type: "string", Validate: `^[a-z]+$`}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// 空值不验证正则
	err := ui.validateAll()
	if err != nil {
		t.Errorf("validateAll() = %v, want nil for empty value", err)
	}
}

func TestValidateNumberRange(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "port", Type: "number", Label: "Port", Min: int64(1), Max: int64(65535)}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// 小于最小值
	setEntryText(ui.widgets["port"], "0")
	err := ui.validateAll()
	if err == nil || err.Error() != "Port must be >= 1" {
		t.Errorf("error = %v, want 'Port must be >= 1'", err)
	}

	// 大于最大值
	setEntryText(ui.widgets["port"], "70000")
	err = ui.validateAll()
	if err == nil || err.Error() != "Port must be <= 65535" {
		t.Errorf("error = %v, want 'Port must be <= 65535'", err)
	}

	// 有效范围
	setEntryText(ui.widgets["port"], "8080")
	err = ui.validateAll()
	if err != nil {
		t.Errorf("validateAll() = %v, want nil", err)
	}
}

func TestValidateNumberInvalid(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "num", Type: "number", Label: "Number", Min: int64(0)}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	setEntryText(ui.widgets["num"], "abc")
	err := ui.validateAll()
	if err == nil || err.Error() != "Number must be a number" {
		t.Errorf("error = %v, want 'Number must be a number'", err)
	}
}

func TestValidateMinOnly(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "count", Type: "number", Label: "Count", Min: int64(1)}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	setEntryText(ui.widgets["count"], "0")
	err := ui.validateAll()
	if err == nil {
		t.Error("validateAll() = nil, want error")
	}

	setEntryText(ui.widgets["count"], "100")
	err = ui.validateAll()
	if err != nil {
		t.Errorf("validateAll() = %v, want nil", err)
	}
}

func TestValidateMaxOnly(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "percent", Type: "number", Label: "Percent", Max: int64(100)}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	setEntryText(ui.widgets["percent"], "150")
	err := ui.validateAll()
	if err == nil {
		t.Error("validateAll() = nil, want error")
	}

	setEntryText(ui.widgets["percent"], "50")
	err = ui.validateAll()
	if err != nil {
		t.Errorf("validateAll() = %v, want nil", err)
	}
}

func TestValidateWithConditionDisabled(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items: []Item{
			{Name: "enable", Type: "bool"},
			{Name: "value", Type: "string", Required: true, Condition: "enable=true"},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// enable=false，条件不满足，required 不生效
	err := ui.validateRequired()
	if err != nil {
		t.Errorf("validateRequired() = %v, want nil (condition not met)", err)
	}

	// enable=true，条件满足，required 生效
	ui.widgets["enable"].(*widget.Check).SetChecked(true)
	err = ui.validateRequired()
	if err == nil {
		t.Error("validateRequired() = nil, want error (condition met)")
	}
}

func TestToFloat(t *testing.T) {
	tests := []struct {
		input any
		want  float64
		ok    bool
	}{
		{int(10), 10.0, true},
		{int64(20), 20.0, true},
		{float64(30.5), 30.5, true},
		{"string", 0, false},
		{nil, 0, false},
	}
	for _, tt := range tests {
		got, ok := toFloat(tt.input)
		if ok != tt.ok || (ok && got != tt.want) {
			t.Errorf("toFloat(%v) = (%v, %v), want (%v, %v)", tt.input, got, ok, tt.want, tt.ok)
		}
	}
}
