package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	toml := `
title = "Test App"
width = 600
height = 500

[[apps]]
[apps.command]
path = "echo"
name = "Echo"
args = ["hello"]
mode = "hidden"
output = "dialog"
debug = true
run_text = "Execute"
run_color = "high"
debug_text = "Preview"
debug_color = "low"

[apps.command.env]
FOO = "bar"

[[apps.items]]
name = "msg"
type = "string"
label = "Message"
default = "world"
`
	path := writeTempFile(t, toml)
	cfg := loadConfig(path)

	if cfg.Title != "Test App" {
		t.Errorf("Title = %q, want %q", cfg.Title, "Test App")
	}
	if cfg.Width != 600 {
		t.Errorf("Width = %v, want 600", cfg.Width)
	}
	if cfg.Height != 500 {
		t.Errorf("Height = %v, want 500", cfg.Height)
	}
	if len(cfg.Apps) != 1 {
		t.Fatalf("len(Apps) = %d, want 1", len(cfg.Apps))
	}

	app := cfg.Apps[0]
	if app.Command.Path != "echo" {
		t.Errorf("Command.Path = %q, want %q", app.Command.Path, "echo")
	}
	if app.Command.Name != "Echo" {
		t.Errorf("Command.Name = %q, want %q", app.Command.Name, "Echo")
	}
	if len(app.Command.Args) != 1 || app.Command.Args[0] != "hello" {
		t.Errorf("Command.Args = %v, want [hello]", app.Command.Args)
	}
	if app.Command.Mode != "hidden" {
		t.Errorf("Command.Mode = %q, want %q", app.Command.Mode, "hidden")
	}
	if app.Command.Output != "dialog" {
		t.Errorf("Command.Output = %q, want %q", app.Command.Output, "dialog")
	}
	if !app.Command.Debug {
		t.Error("Command.Debug = false, want true")
	}
	if app.Command.RunText != "Execute" {
		t.Errorf("Command.RunText = %q, want %q", app.Command.RunText, "Execute")
	}
	if app.Command.RunColor != "high" {
		t.Errorf("Command.RunColor = %q, want %q", app.Command.RunColor, "high")
	}
	if app.Command.DebugText != "Preview" {
		t.Errorf("Command.DebugText = %q, want %q", app.Command.DebugText, "Preview")
	}
	if app.Command.DebugColor != "low" {
		t.Errorf("Command.DebugColor = %q, want %q", app.Command.DebugColor, "low")
	}
	if app.Command.Env["FOO"] != "bar" {
		t.Errorf("Command.Env[FOO] = %q, want %q", app.Command.Env["FOO"], "bar")
	}

	if len(app.Items) != 1 {
		t.Fatalf("len(Items) = %d, want 1", len(app.Items))
	}
	item := app.Items[0]
	if item.Name != "msg" {
		t.Errorf("Item.Name = %q, want %q", item.Name, "msg")
	}
	if item.Type != "string" {
		t.Errorf("Item.Type = %q, want %q", item.Type, "string")
	}
	if item.Label != "Message" {
		t.Errorf("Item.Label = %q, want %q", item.Label, "Message")
	}
	if item.Default != "world" {
		t.Errorf("Item.Default = %v, want %q", item.Default, "world")
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	toml := `
[[apps]]
[apps.command]
path = "echo"
name = "Test"
`
	path := writeTempFile(t, toml)
	cfg := loadConfig(path)

	if cfg.Width != 400 {
		t.Errorf("default Width = %v, want 400", cfg.Width)
	}
	if cfg.Height != 300 {
		t.Errorf("default Height = %v, want 300", cfg.Height)
	}
	if cfg.Apps[0].Command.Mode != "hidden" {
		t.Errorf("default Mode = %q, want %q", cfg.Apps[0].Command.Mode, "hidden")
	}
	if cfg.Apps[0].Command.Output != "dialog" {
		t.Errorf("default Output = %q, want %q", cfg.Apps[0].Command.Output, "dialog")
	}
}

func TestItemIsLabel(t *testing.T) {
	tests := []struct {
		item Item
		want bool
	}{
		{Item{Text: "Hello", Name: ""}, true},
		{Item{Text: "", Name: "field"}, false},
		{Item{Text: "Hello", Name: "field"}, false},
		{Item{Text: "", Name: ""}, false},
	}
	for _, tt := range tests {
		got := tt.item.IsLabel()
		if got != tt.want {
			t.Errorf("Item{Text:%q, Name:%q}.IsLabel() = %v, want %v", tt.item.Text, tt.item.Name, got, tt.want)
		}
	}
}

func TestLoadConfigAllItemTypes(t *testing.T) {
	toml := `
[[apps]]
[apps.command]
path = "test"
name = "Test"

[[apps.items]]
name = "str"
type = "string"
default = "hello"

[[apps.items]]
name = "num"
type = "number"
default = 42

[[apps.items]]
name = "flag"
type = "bool"
default = true

[[apps.items]]
name = "opt"
type = "choice"
choices = ["a", "b", "c"]
default = "b"
`
	path := writeTempFile(t, toml)
	cfg := loadConfig(path)

	items := cfg.Apps[0].Items
	if items[0].Default != "hello" {
		t.Errorf("string default = %v, want %q", items[0].Default, "hello")
	}
	if items[1].Default != int64(42) {
		t.Errorf("number default = %v (%T), want 42", items[1].Default, items[1].Default)
	}
	if items[2].Default != true {
		t.Errorf("bool default = %v, want true", items[2].Default)
	}
	if items[3].Default != "b" {
		t.Errorf("choice default = %v, want %q", items[3].Default, "b")
	}
	if len(items[3].Choices) != 3 {
		t.Errorf("choices len = %d, want 3", len(items[3].Choices))
	}
}

func TestLoadConfigItemOptions(t *testing.T) {
	toml := `
[[apps]]
[apps.command]
path = "test"
name = "Test"

[[apps.items]]
name = "file"
type = "string"
picker = "file"
picker_text = "Browse"

[[apps.items]]
name = "dir"
type = "string"
picker = "directory"

[[apps.items]]
name = "multi"
type = "string"
multi = true

[[apps.items]]
name = "pos"
type = "string"
positional = true

[[apps.items]]
name = "s"
type = "bool"
short = true

[[apps.items]]
name = "sep"
type = "string"
separator = " "

[[apps.items]]
name = "req"
type = "string"
required = true

[[apps.items]]
name = "val"
type = "string"
validate = "^[a-z]+$"

[[apps.items]]
name = "range"
type = "number"
min = 0
max = 100

[[apps.items]]
name = "cond"
type = "string"
condition = "flag=true"
`
	path := writeTempFile(t, toml)
	cfg := loadConfig(path)

	items := cfg.Apps[0].Items
	if items[0].Picker != "file" {
		t.Errorf("picker = %q, want %q", items[0].Picker, "file")
	}
	if items[0].PickerText != "Browse" {
		t.Errorf("picker_text = %q, want %q", items[0].PickerText, "Browse")
	}
	if items[1].Picker != "directory" {
		t.Errorf("picker = %q, want %q", items[1].Picker, "directory")
	}
	if !items[2].Multi {
		t.Error("multi = false, want true")
	}
	if !items[3].Positional {
		t.Error("positional = false, want true")
	}
	if !items[4].Short {
		t.Error("short = false, want true")
	}
	if items[5].Separator != " " {
		t.Errorf("separator = %q, want %q", items[5].Separator, " ")
	}
	if !items[6].Required {
		t.Error("required = false, want true")
	}
	if items[7].Validate != "^[a-z]+$" {
		t.Errorf("validate = %q, want %q", items[7].Validate, "^[a-z]+$")
	}
	if items[8].Min != int64(0) {
		t.Errorf("min = %v, want 0", items[8].Min)
	}
	if items[8].Max != int64(100) {
		t.Errorf("max = %v, want 100", items[8].Max)
	}
	if items[9].Condition != "flag=true" {
		t.Errorf("condition = %q, want %q", items[9].Condition, "flag=true")
	}
}

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}
