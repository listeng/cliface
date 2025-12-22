package main

import (
	"reflect"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestBuildArgsBasic(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd", Args: []string{"sub"}},
		Items: []Item{
			{Name: "flag", Type: "bool"},
			{Name: "opt", Type: "string"},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// 设置值
	ui.widgets["flag"].(*widget.Check).SetChecked(true)
	setEntryText(ui.widgets["opt"], "value")

	args := ui.BuildArgs()
	want := []string{"sub", "--flag", "--opt=value"}
	if !reflect.DeepEqual(args, want) {
		t.Errorf("BuildArgs() = %v, want %v", args, want)
	}
}

func TestBuildArgsShortFlag(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items: []Item{
			{Name: "v", Type: "bool", Short: true},
			{Name: "o", Type: "string", Short: true},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	ui.widgets["v"].(*widget.Check).SetChecked(true)
	setEntryText(ui.widgets["o"], "out.txt")

	args := ui.BuildArgs()
	want := []string{"-v", "-o=out.txt"}
	if !reflect.DeepEqual(args, want) {
		t.Errorf("BuildArgs() = %v, want %v", args, want)
	}
}

func TestBuildArgsPositional(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd", Args: []string{"clone"}},
		Items: []Item{
			{Name: "url", Type: "string", Positional: true},
			{Name: "dir", Type: "string", Positional: true},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	setEntryText(ui.widgets["url"], "https://example.com/repo.git")
	setEntryText(ui.widgets["dir"], "/tmp/repo")

	args := ui.BuildArgs()
	want := []string{"clone", "https://example.com/repo.git", "/tmp/repo"}
	if !reflect.DeepEqual(args, want) {
		t.Errorf("BuildArgs() = %v, want %v", args, want)
	}
}

func TestBuildArgsSeparators(t *testing.T) {
	tests := []struct {
		name      string
		separator string
		want      string
	}{
		{"default (=)", "", "--opt=val"},
		{"space", " ", "--opt"},
		{"none", "none", "--optval"},
		{"custom", ":", "--opt:val"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				Command: Command{Path: "cmd"},
				Items:   []Item{{Name: "opt", Type: "string", Separator: tt.separator}},
			}
			w := test.NewWindow(nil)
			ui := NewAppUI(app, w)
			ui.Build()
			setEntryText(ui.widgets["opt"], "val")

			args := ui.BuildArgs()
			if tt.separator == " " {
				want := []string{"--opt", "val"}
				if !reflect.DeepEqual(args, want) {
					t.Errorf("BuildArgs() = %v, want %v", args, want)
				}
			} else {
				if len(args) != 1 || args[0] != tt.want {
					t.Errorf("BuildArgs() = %v, want [%s]", args, tt.want)
				}
			}
		})
	}
}

func TestBuildArgsChoice(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "format", Type: "choice", Choices: []string{"json", "xml", "csv"}}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	ui.widgets["format"].(*widget.Select).SetSelected("json")

	args := ui.BuildArgs()
	want := []string{"--format=json"}
	if !reflect.DeepEqual(args, want) {
		t.Errorf("BuildArgs() = %v, want %v", args, want)
	}
}

func TestBuildArgsMulti(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "H", Type: "string", Short: true, Multi: true, Separator: " "}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	mw := ui.widgets["H"].(*multiWidget)
	mw.entries[0].SetText("Content-Type: application/json")
	mw.addEntry()
	mw.entries[1].SetText("Authorization: Bearer token")

	args := ui.BuildArgs()
	want := []string{"-H", "Content-Type: application/json", "-H", "Authorization: Bearer token"}
	if !reflect.DeepEqual(args, want) {
		t.Errorf("BuildArgs() = %v, want %v", args, want)
	}
}

func TestBuildArgsEmptyValues(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items: []Item{
			{Name: "empty", Type: "string"},
			{Name: "filled", Type: "string"},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	setEntryText(ui.widgets["filled"], "value")
	// empty 保持空

	args := ui.BuildArgs()
	want := []string{"--filled=value"}
	if !reflect.DeepEqual(args, want) {
		t.Errorf("BuildArgs() = %v, want %v", args, want)
	}
}

func TestBuildArgsBoolFalse(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "verbose", Type: "bool"}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// bool 默认 false，不应该出现在参数中
	args := ui.BuildArgs()
	if len(args) != 0 {
		t.Errorf("BuildArgs() = %v, want []", args)
	}
}

func TestBuildArgsWithDefaults(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items: []Item{
			{Name: "str", Type: "string", Default: "default_val"},
			{Name: "num", Type: "number", Default: int64(10)},
			{Name: "flag", Type: "bool", Default: true},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	args := ui.BuildArgs()
	want := []string{"--str=default_val", "--num=10", "--flag"}
	if !reflect.DeepEqual(args, want) {
		t.Errorf("BuildArgs() = %v, want %v", args, want)
	}
}

func TestBuildArgsQuotedInput(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "msg", Type: "string"}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// 用户输入带引号的值，应该去掉引号
	setEntryText(ui.widgets["msg"], `"hello world"`)

	args := ui.BuildArgs()
	want := []string{"--msg=hello world"}
	if !reflect.DeepEqual(args, want) {
		t.Errorf("BuildArgs() = %v, want %v", args, want)
	}
}

func TestBuildArgsLabel(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items: []Item{
			{Text: "This is a label"},
			{Name: "opt", Type: "string"},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	setEntryText(ui.widgets["opt"], "value")

	args := ui.BuildArgs()
	want := []string{"--opt=value"}
	if !reflect.DeepEqual(args, want) {
		t.Errorf("BuildArgs() = %v, want %v", args, want)
	}
}

func setEntryText(w interface{}, text string) {
	switch v := w.(type) {
	case *widget.Entry:
		v.SetText(text)
	}
}

func TestBuildUI_SingleApp(t *testing.T) {
	cfg := &Config{
		Apps: []App{{Command: Command{Path: "cmd", Name: "Test"}}},
	}
	w := test.NewWindow(nil)
	content := BuildUI(cfg, w)

	// 单 app 不应该是 tabs
	if _, ok := content.(*container.AppTabs); ok {
		t.Error("single app should not create tabs")
	}
}

func TestBuildUI_MultiApp(t *testing.T) {
	cfg := &Config{
		Apps: []App{
			{Command: Command{Path: "cmd1", Name: "App1"}},
			{Command: Command{Path: "cmd2", Name: "App2"}},
		},
	}
	w := test.NewWindow(nil)
	content := BuildUI(cfg, w)

	tabs, ok := content.(*container.AppTabs)
	if !ok {
		t.Fatal("multi app should create tabs")
	}
	if len(tabs.Items) != 2 {
		t.Errorf("tabs count = %d, want 2", len(tabs.Items))
	}
	if tabs.Items[0].Text != "App1" {
		t.Errorf("tab[0].Text = %q, want %q", tabs.Items[0].Text, "App1")
	}
	if tabs.Items[1].Text != "App2" {
		t.Errorf("tab[1].Text = %q, want %q", tabs.Items[1].Text, "App2")
	}
}

func TestBuildCommandLine(t *testing.T) {
	app := &App{
		Command: Command{Path: "/usr/bin/echo", Args: []string{"hello"}},
		Items: []Item{
			{Name: "msg", Type: "string"},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	setEntryText(ui.widgets["msg"], "world")
	cmdLine := ui.buildCommandLine()

	want := "/usr/bin/echo hello --msg=world"
	if cmdLine != want {
		t.Errorf("buildCommandLine() = %q, want %q", cmdLine, want)
	}
}

func TestBuildCommandLineWithSpaces(t *testing.T) {
	app := &App{
		Command: Command{Path: "echo"},
		Items:   []Item{{Name: "msg", Type: "string"}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	setEntryText(ui.widgets["msg"], "hello world")
	cmdLine := ui.buildCommandLine()

	// 包含空格的参数应该被引号包裹
	if !strings.Contains(cmdLine, "'--msg=hello world'") && !strings.Contains(cmdLine, "\"--msg=hello world\"") {
		t.Errorf("buildCommandLine() = %q, should quote arg with spaces", cmdLine)
	}
}

func TestCreateWidgetPicker(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items: []Item{
			{Name: "file", Type: "string", Picker: "file", PickerText: "Browse"},
			{Name: "dir", Type: "string", Picker: "directory"},
		},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// file picker 应该是 container
	fileWidget := ui.widgets["file"]
	if _, ok := fileWidget.(*fyne.Container); !ok {
		t.Error("file picker should be a container")
	}

	// directory picker 应该是 container
	dirWidget := ui.widgets["dir"]
	if _, ok := dirWidget.(*fyne.Container); !ok {
		t.Error("directory picker should be a container")
	}
}

func TestCreateWidgetNumber(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "port", Type: "number", Default: int64(8080)}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	entry := ui.widgets["port"].(*widget.Entry)
	if entry.Text != "8080" {
		t.Errorf("number default = %q, want %q", entry.Text, "8080")
	}
}

func TestCreateWidgetUnknownType(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "field", Type: "unknown"}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// 未知类型应该创建 Entry
	if _, ok := ui.widgets["field"].(*widget.Entry); !ok {
		t.Error("unknown type should create Entry")
	}
}

func TestBuildWithDescription(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "field", Type: "string", Description: "Help text"}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	content := ui.Build()

	// 应该成功构建，不报错
	if content == nil {
		t.Error("Build() returned nil")
	}
}

func TestBuildWithDebugButton(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd", Debug: true, DebugText: "Preview", DebugColor: "warning"},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	content := ui.Build()

	if content == nil {
		t.Error("Build() returned nil")
	}
}

func TestBuildWithCustomRunButton(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd", RunText: "Execute", RunColor: "success"},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	content := ui.Build()

	if content == nil {
		t.Error("Build() returned nil")
	}
}

func TestMultiWidgetValues(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "tags", Type: "string", Multi: true}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	mw := ui.widgets["tags"].(*multiWidget)

	// 初始一个空 entry
	vals := mw.Values()
	if len(vals) != 0 {
		t.Errorf("initial Values() = %v, want []", vals)
	}

	// 添加值
	mw.entries[0].SetText("tag1")
	mw.addEntry()
	mw.entries[1].SetText("tag2")
	mw.addEntry()
	// 第三个保持空

	vals = mw.Values()
	want := []string{"tag1", "tag2"}
	if !reflect.DeepEqual(vals, want) {
		t.Errorf("Values() = %v, want %v", vals, want)
	}
}

func TestNoSpaceVBoxMinSize(t *testing.T) {
	layout := &noSpaceVBox{}
	label1 := widget.NewLabel("Hello")
	label2 := widget.NewLabel("World")

	size := layout.MinSize([]fyne.CanvasObject{label1, label2})
	if size.Height <= 0 {
		t.Error("MinSize height should be > 0")
	}
	if size.Width <= 0 {
		t.Error("MinSize width should be > 0")
	}
}

func TestGetWidgetValueFromContainer(t *testing.T) {
	app := &App{
		Command: Command{Path: "cmd"},
		Items:   []Item{{Name: "file", Type: "string", Picker: "file"}},
	}
	w := test.NewWindow(nil)
	ui := NewAppUI(app, w)
	ui.Build()

	// picker 类型是 container，需要从中提取 entry
	c := ui.widgets["file"].(*fyne.Container)
	for _, obj := range c.Objects {
		if entry, ok := obj.(*widget.Entry); ok {
			entry.SetText("/path/to/file")
			break
		}
	}

	args := ui.BuildArgs()
	want := []string{"--file=/path/to/file"}
	if !reflect.DeepEqual(args, want) {
		t.Errorf("BuildArgs() = %v, want %v", args, want)
	}
}
