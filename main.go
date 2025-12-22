package main

import (
	"flag"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/BurntSushi/toml"
)

func main() {
	configPath := flag.String("c", "", "config file path")
	flag.StringVar(configPath, "config", "", "config file path")
	flag.Parse()

	if *configPath == "" {
		exe, _ := os.Executable()
		*configPath = filepath.Join(filepath.Dir(exe), "config.toml")
	}

	cfg := loadConfig(*configPath)

	a := app.New()
	title := cfg.Title
	if title == "" {
		if len(cfg.Apps) == 1 {
			title = cfg.Apps[0].Command.Name
		} else {
			title = "cliface"
		}
	}
	w := a.NewWindow(title)
	w.Resize(fyne.NewSize(cfg.Width, cfg.Height))
	w.SetContent(BuildUI(cfg, w))
	w.ShowAndRun()
}

func loadConfig(path string) *Config {
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		panic(err)
	}
	for i := range cfg.Apps {
		if cfg.Apps[i].Command.Mode == "" {
			cfg.Apps[i].Command.Mode = "hidden"
		}
		if cfg.Apps[i].Command.Output == "" {
			cfg.Apps[i].Command.Output = "dialog"
		}
	}
	if cfg.Width == 0 {
		cfg.Width = 400
	}
	if cfg.Height == 0 {
		cfg.Height = 300
	}
	return &cfg
}
