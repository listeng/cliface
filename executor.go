package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (u *AppUI) Execute() {
	if err := u.validateRequired(); err != nil {
		dialog.ShowError(err, u.window)
		return
	}
	if err := u.validateAll(); err != nil {
		dialog.ShowError(err, u.window)
		return
	}

	args := u.BuildArgs()
	cmd := exec.Command(u.app.Command.Path, args...)

	// 设置环境变量
	if len(u.app.Command.Env) > 0 {
		cmd.Env = os.Environ()
		for k, v := range u.app.Command.Env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}

	if u.app.Command.Mode == "visible" {
		cmd.Start()
		return
	}

	switch u.app.Command.Output {
	case "realtime":
		u.executeRealtime(cmd)
	case "realtime-console":
		u.executeConsole(cmd)
	default:
		u.executeDialog(cmd)
	}
}

func (u *AppUI) executeDialog(cmd *exec.Cmd) {
	done := make(chan struct{})
	prog := dialog.NewCustomConfirm("执行中", "取消", "", widget.NewProgressBarInfinite(), func(cancel bool) {
		if cancel {
			cmd.Process.Kill()
		}
	}, u.window)
	prog.Show()
	go func() {
		output, err := cmd.CombinedOutput()
		close(done)
		prog.Hide()
		text := string(output)
		if err != nil {
			text += "\n\n错误: " + err.Error()
		}
		entry := widget.NewMultiLineEntry()
		entry.SetText(text)
		entry.Wrapping = fyne.TextWrapWord
		win := fyne.CurrentApp().NewWindow("Output")
		win.SetContent(container.NewScroll(entry))
		win.Resize(fyne.NewSize(500, 400))
		win.Show()
	}()
}

func (u *AppUI) executeRealtime(cmd *exec.Cmd) {
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	entry := widget.NewMultiLineEntry()
	entry.Wrapping = fyne.TextWrapWord

	cancelBtn := widget.NewButton("取消", nil)
	cancelBtn.Importance = widget.DangerImportance

	win := fyne.CurrentApp().NewWindow("Output")
	win.SetContent(container.NewBorder(cancelBtn, nil, nil, nil, container.NewScroll(entry)))
	win.Resize(fyne.NewSize(500, 400))
	win.Show()

	if err := cmd.Start(); err != nil {
		entry.SetText("Error: " + err.Error())
		return
	}

	cancelBtn.OnTapped = func() {
		cmd.Process.Kill()
		entry.SetText(entry.Text + "\n[已取消]")
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		var lines []string
		const maxLines = 500
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
			if len(lines) > maxLines {
				lines = lines[len(lines)-maxLines:]
			}
			entry.SetText(strings.Join(lines, "\n"))
		}
	}()
}

func (u *AppUI) executeConsole(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Printf(">>> %s %s\n", u.app.Command.Path, strings.Join(u.BuildArgs(), " "))
	cmd.Run()
	fmt.Println("<<< done")
}
