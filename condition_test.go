package main

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestParseCondition(t *testing.T) {
	tests := []struct {
		input     string
		wantField string
		wantOp    string
		wantValue string
	}{
		{"field=value", "field", "=", "value"},
		{"field!=value", "field", "!=", "value"},
		{"flag=true", "flag", "=", "true"},
		{"flag!=", "flag", "!=", ""},
		{"field", "field", "!=", ""},
	}
	for _, tt := range tests {
		field, op, value := parseCondition(tt.input)
		if field != tt.wantField || op != tt.wantOp || value != tt.wantValue {
			t.Errorf("parseCondition(%q) = (%q, %q, %q), want (%q, %q, %q)",
				tt.input, field, op, value, tt.wantField, tt.wantOp, tt.wantValue)
		}
	}
}

func TestCheckConditionEqual(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items: []Item{
			{Name: "mode", Type: "choice", Choices: []string{"simple", "advanced"}},
			{Name: "extra", Type: "string", Condition: "mode=advanced"},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// mode=simple，条件不满足
	ui.widgets["mode"].(*widget.Select).SetSelected("simple")
	if ui.checkCondition(&app.Items[1]) {
		t.Error("checkCondition() = true, want false (mode=simple)")
	}

	// mode=advanced，条件满足
	ui.widgets["mode"].(*widget.Select).SetSelected("advanced")
	if !ui.checkCondition(&app.Items[1]) {
		t.Error("checkCondition() = false, want true (mode=advanced)")
	}
}

func TestCheckConditionNotEqual(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items: []Item{
			{Name: "skip", Type: "bool"},
			{Name: "value", Type: "string", Condition: "skip!=true"},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// skip=false，条件满足
	if !ui.checkCondition(&app.Items[1]) {
		t.Error("checkCondition() = false, want true (skip=false)")
	}

	// skip=true，条件不满足
	ui.widgets["skip"].(*widget.Check).SetChecked(true)
	if ui.checkCondition(&app.Items[1]) {
		t.Error("checkCondition() = true, want false (skip=true)")
	}
}

func TestCheckConditionNoCondition(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "field", Type: "string"}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// 无条件，始终返回 true
	if !ui.checkCondition(&app.Items[0]) {
		t.Error("checkCondition() = false, want true (no condition)")
	}
}

func TestCheckConditionMissingField(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "field", Type: "string", Condition: "nonexistent=value"}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// 依赖字段不存在，返回 true
	if !ui.checkCondition(&app.Items[0]) {
		t.Error("checkCondition() = false, want true (missing field)")
	}
}

func TestCheckConditionWithBool(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items: []Item{
			{Name: "enabled", Type: "bool"},
			{Name: "config", Type: "string", Condition: "enabled=true"},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// enabled=false
	if ui.checkCondition(&app.Items[1]) {
		t.Error("checkCondition() = true, want false (enabled=false)")
	}

	// enabled=true
	ui.widgets["enabled"].(*widget.Check).SetChecked(true)
	if !ui.checkCondition(&app.Items[1]) {
		t.Error("checkCondition() = false, want true (enabled=true)")
	}
}

func TestCheckConditionWithEntry(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items: []Item{
			{Name: "type", Type: "string"},
			{Name: "extra", Type: "string", Condition: "type=custom"},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// type 为空
	if ui.checkCondition(&app.Items[1]) {
		t.Error("checkCondition() = true, want false (type empty)")
	}

	// type=custom
	setEntryText(ui.widgets["type"], "custom")
	if !ui.checkCondition(&app.Items[1]) {
		t.Error("checkCondition() = false, want true (type=custom)")
	}

	// type=other
	setEntryText(ui.widgets["type"], "other")
	if ui.checkCondition(&app.Items[1]) {
		t.Error("checkCondition() = true, want false (type=other)")
	}
}

func TestParseImportance(t *testing.T) {
	tests := []struct {
		color string
		want  widget.Importance
	}{
		{"high", widget.HighImportance},
		{"danger", widget.DangerImportance},
		{"warning", widget.WarningImportance},
		{"success", widget.SuccessImportance},
		{"low", widget.LowImportance},
		{"", widget.MediumImportance},
		{"unknown", widget.MediumImportance},
	}
	for _, tt := range tests {
		got := parseImportance(tt.color)
		if got != tt.want {
			t.Errorf("parseImportance(%q) = %v, want %v", tt.color, got, tt.want)
		}
	}
}

func TestSetupConditionsWithSelect(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items: []Item{
			{Name: "mode", Type: "choice", Choices: []string{"simple", "advanced"}},
			{Name: "extra", Type: "string", Condition: "mode=advanced"},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// 初始状态 mode 为空，extra 应该被禁用
	extraWidget := ui.widgets["extra"].(*widget.Entry)
	if !extraWidget.Disabled() {
		t.Error("extra should be disabled when mode is empty")
	}

	// 设置 mode=advanced，extra 应该启用
	ui.widgets["mode"].(*widget.Select).SetSelected("advanced")
	if extraWidget.Disabled() {
		t.Error("extra should be enabled when mode=advanced")
	}

	// 设置 mode=simple，extra 应该禁用
	ui.widgets["mode"].(*widget.Select).SetSelected("simple")
	if !extraWidget.Disabled() {
		t.Error("extra should be disabled when mode=simple")
	}
}

func TestSetupConditionsWithCheck(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items: []Item{
			{Name: "enable", Type: "bool"},
			{Name: "config", Type: "string", Condition: "enable=true"},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	configWidget := ui.widgets["config"].(*widget.Entry)

	// 初始 enable=false，config 禁用
	if !configWidget.Disabled() {
		t.Error("config should be disabled when enable=false")
	}

	// enable=true，config 启用
	ui.widgets["enable"].(*widget.Check).SetChecked(true)
	if configWidget.Disabled() {
		t.Error("config should be enabled when enable=true")
	}
}

func TestSetupConditionsWithEntry(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items: []Item{
			{Name: "type", Type: "string"},
			{Name: "custom", Type: "string", Condition: "type=custom"},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	customWidget := ui.widgets["custom"].(*widget.Entry)
	typeEntry := ui.widgets["type"].(*widget.Entry)

	// 初始 type 为空
	if !customWidget.Disabled() {
		t.Error("custom should be disabled when type is empty")
	}

	// type=custom
	typeEntry.SetText("custom")
	if customWidget.Disabled() {
		t.Error("custom should be enabled when type=custom")
	}
}

func TestSetupConditionsWithPicker(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items: []Item{
			{Name: "source", Type: "string", Picker: "file"},
			{Name: "dest", Type: "string", Condition: "source!="},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	destWidget := ui.widgets["dest"].(*widget.Entry)

	// 初始 source 为空，dest 禁用
	if !destWidget.Disabled() {
		t.Error("dest should be disabled when source is empty")
	}

	// 设置 source
	c := ui.widgets["source"].(*fyne.Container)
	for _, obj := range c.Objects {
		if entry, ok := obj.(*widget.Entry); ok {
			entry.SetText("/path/to/file")
			break
		}
	}
	if destWidget.Disabled() {
		t.Error("dest should be enabled when source is not empty")
	}
}

func TestUpdateWidgetStateWithPicker(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items: []Item{
			{Name: "enable", Type: "bool"},
			{Name: "file", Type: "string", Picker: "file", Condition: "enable=true"},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// picker container 内的 entry 和 button 都应该被禁用/启用
	c := ui.widgets["file"].(*fyne.Container)
	var entry *widget.Entry
	var btn *widget.Button
	for _, obj := range c.Objects {
		if e, ok := obj.(*widget.Entry); ok {
			entry = e
		}
		if b, ok := obj.(*widget.Button); ok {
			btn = b
		}
	}

	// 初始禁用
	if !entry.Disabled() {
		t.Error("picker entry should be disabled")
	}
	if !btn.Disabled() {
		t.Error("picker button should be disabled")
	}

	// 启用
	ui.widgets["enable"].(*widget.Check).SetChecked(true)
	if entry.Disabled() {
		t.Error("picker entry should be enabled")
	}
	if btn.Disabled() {
		t.Error("picker button should be enabled")
	}
}
