package main

type Config struct {
	Title  string  `toml:"title"`
	Width  float32 `toml:"width"`
	Height float32 `toml:"height"`
	Apps   []App   `toml:"apps"`
}

type App struct {
	Command Command `toml:"command"`
	Items   []Item  `toml:"items"`
}

type Command struct {
	Path       string            `toml:"path"`
	Name       string            `toml:"name"`
	Args       []string          `toml:"args"`
	Mode       string            `toml:"mode"`
	Output     string            `toml:"output"`
	Debug      bool              `toml:"debug"`
	RunText    string            `toml:"run_text"`
	RunColor   string            `toml:"run_color"`
	DebugText  string            `toml:"debug_text"`
	DebugColor string            `toml:"debug_color"`
	Env        map[string]string `toml:"env"`
}

type Item struct {
	// label 类型
	Text string `toml:"text"`
	// argument 类型
	Name        string   `toml:"name"`
	Short       bool     `toml:"short"`
	Positional  bool     `toml:"positional"`
	Type        string   `toml:"type"`
	Label       string   `toml:"label"`
	Description string   `toml:"description"`
	Default     any      `toml:"default"`
	Choices     []string `toml:"choices"`
	Placeholder string   `toml:"placeholder"`
	Picker      string   `toml:"picker"`
	PickerText  string   `toml:"picker_text"`
	Separator   string   `toml:"separator"`
	Multi       bool     `toml:"multi"`
	// 验证
	Required  bool   `toml:"required"`
	Validate  string `toml:"validate"`
	Min       any    `toml:"min"`
	Max       any    `toml:"max"`
	Condition string `toml:"condition"`
}

func (i *Item) IsLabel() bool {
	return i.Text != "" && i.Name == ""
}
