package main

import (
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type AppUI struct {
	app     *App
	widgets map[string]fyne.CanvasObject
	window  fyne.Window
}

func BuildUI(cfg *Config, w fyne.Window) fyne.CanvasObject {
	if len(cfg.Apps) == 1 {
		ui := NewAppUI(&cfg.Apps[0], w)
		return ui.Build()
	}
	tabs := container.NewAppTabs()
	for i := range cfg.Apps {
		ui := NewAppUI(&cfg.Apps[i], w)
		tabs.Append(container.NewTabItem(cfg.Apps[i].Command.Name, ui.Build()))
	}
	return tabs
}

func NewAppUI(app *App, w fyne.Window) *AppUI {
	return &AppUI{
		app:     app,
		widgets: make(map[string]fyne.CanvasObject),
		window:  w,
	}
}

func (u *AppUI) Build() fyne.CanvasObject {
	// 计算最大label宽度
	var maxWidth float32
	for _, item := range u.app.Items {
		if item.IsLabel() {
			continue
		}
		text := item.Label
		if text == "" {
			text = item.Name
		}
		lbl := widget.NewLabel(text)
		if w := lbl.MinSize().Width; w > maxWidth {
			maxWidth = w
		}
	}

	form := container.New(&noSpaceVBox{})
	for i := range u.app.Items {
		item := &u.app.Items[i]
		if item.IsLabel() {
			lbl := widget.NewLabel(item.Text)
			lbl.Wrapping = fyne.TextWrapWord
			form.Add(container.NewPadded(lbl))
			continue
		}
		w := u.createWidget(item)
		u.widgets[item.Name] = w
		lbl := widget.NewLabel(item.Label)
		if item.Label == "" {
			lbl.SetText(item.Name)
		}
		lbl.Alignment = fyne.TextAlignTrailing
		labelBox := container.NewHBox(layout.NewSpacer(), lbl)
		labelCol := container.NewGridWrap(fyne.NewSize(maxWidth+10, 0), labelBox)
		row := container.NewBorder(nil, nil, labelCol, nil, w)
		form.Add(container.NewPadded(row))
		if item.Description != "" {
			hint := widget.NewLabel(item.Description)
			hint.Wrapping = fyne.TextWrapWord
			hintRow := container.NewBorder(nil, nil, container.NewGridWrap(fyne.NewSize(maxWidth+10, 0)), nil, hint)
			form.Add(hintRow)
		}
	}

	runText := u.app.Command.RunText
	if runText == "" {
		runText = "Run"
	}
	runBtn := widget.NewButton(runText, func() { u.Execute() })
	runBtn.Importance = parseImportance(u.app.Command.RunColor)

	var buttons fyne.CanvasObject
	if u.app.Command.Debug {
		debugText := u.app.Command.DebugText
		if debugText == "" {
			debugText = "Show Command"
		}
		debugBtn := widget.NewButton(debugText, func() { u.showCommand() })
		debugBtn.Importance = parseImportance(u.app.Command.DebugColor)
		buttons = container.NewBorder(nil, nil, nil, debugBtn, runBtn)
	} else {
		buttons = runBtn
	}

	form.Add(buttons)

	// 设置条件监听
	u.setupConditions()

	return form
}

func (u *AppUI) createWidget(item *Item) fyne.CanvasObject {
	if item.Multi {
		return u.createMultiWidget(item)
	}
	switch item.Type {
	case "string":
		entry := widget.NewEntry()
		if item.Default != nil {
			entry.SetText(fmt.Sprintf("%v", item.Default))
		}
		if item.Picker == "file" || item.Picker == "directory" {
			btnText := item.PickerText
			if btnText == "" {
				btnText = "..."
			}
			picker := item.Picker
			btn := widget.NewButton(btnText, func() {
				if picker == "file" {
					dialog.ShowFileOpen(func(f fyne.URIReadCloser, err error) {
						if f != nil {
							entry.SetText(f.URI().Path())
							f.Close()
						}
					}, u.window)
				} else {
					dialog.ShowFolderOpen(func(f fyne.ListableURI, err error) {
						if f != nil {
							entry.SetText(f.Path())
						}
					}, u.window)
				}
			})
			return container.NewBorder(nil, nil, nil, btn, entry)
		}
		return entry
	case "number":
		entry := widget.NewEntry()
		if item.Default != nil {
			entry.SetText(fmt.Sprintf("%v", item.Default))
		}
		return entry
	case "bool":
		check := widget.NewCheck("", nil)
		if item.Default != nil {
			if v, ok := item.Default.(bool); ok {
				check.SetChecked(v)
			}
		}
		return check
	case "choice":
		sel := widget.NewSelect(item.Choices, nil)
		if item.Default != nil {
			sel.SetSelected(fmt.Sprintf("%v", item.Default))
		}
		return sel
	default:
		return widget.NewEntry()
	}
}

func (u *AppUI) createMultiWidget(item *Item) fyne.CanvasObject {
	mw := &multiWidget{entries: []*widget.Entry{}}
	mw.vbox = container.NewVBox()

	mw.addEntry = func() {
		entry := widget.NewEntry()
		mw.entries = append(mw.entries, entry)
		removeBtn := widget.NewButton("-", nil)
		row := container.NewBorder(nil, nil, nil, removeBtn, entry)
		removeBtn.OnTapped = func() {
			if len(mw.entries) <= 1 {
				return
			}
			for i, e := range mw.entries {
				if e == entry {
					mw.entries = append(mw.entries[:i], mw.entries[i+1:]...)
					break
				}
			}
			mw.vbox.Remove(row)
		}
		mw.vbox.Add(row)
	}

	mw.addEntry()
	addBtn := widget.NewButton("+", func() { mw.addEntry() })
	mw.vbox.Add(addBtn)

	return mw
}

type multiWidget struct {
	widget.BaseWidget
	vbox     *fyne.Container
	entries  []*widget.Entry
	addEntry func()
}

func (m *multiWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(m.vbox)
}

func (m *multiWidget) MinSize() fyne.Size {
	return m.vbox.MinSize()
}

func (m *multiWidget) Resize(size fyne.Size) {
	m.BaseWidget.Resize(size)
	m.vbox.Resize(size)
}

func (m *multiWidget) Values() []string {
	var vals []string
	for _, e := range m.entries {
		if e.Text != "" {
			vals = append(vals, e.Text)
		}
	}
	return vals
}

func (u *AppUI) BuildArgs() []string {
	args := append([]string{}, u.app.Command.Args...)
	for _, item := range u.app.Items {
		if item.IsLabel() {
			continue
		}
		w := u.widgets[item.Name]
		if item.Multi {
			if mw, ok := w.(*multiWidget); ok {
				for _, val := range mw.Values() {
					prefix := "--"
					if item.Short {
						prefix = "-"
					}
					if item.Separator == " " {
						args = append(args, prefix+item.Name, val)
					} else if item.Separator == "none" {
						args = append(args, prefix+item.Name+val)
					} else if item.Separator == "" {
						args = append(args, prefix+item.Name+"="+val)
					} else {
						args = append(args, prefix+item.Name+item.Separator+val)
					}
				}
			}
			continue
		}
		val := u.getWidgetValue(&item, w)
		if val == "" {
			continue
		}
		if item.Positional {
			args = append(args, val)
			continue
		}
		prefix := "--"
		if item.Short {
			prefix = "-"
		}
		if item.Type == "bool" {
			if val == "true" {
				args = append(args, prefix+item.Name)
			}
		} else {
			if item.Separator == " " {
				args = append(args, prefix+item.Name, val)
			} else if item.Separator == "none" {
				args = append(args, prefix+item.Name+val)
			} else if item.Separator == "" {
				args = append(args, prefix+item.Name+"="+val)
			} else {
				args = append(args, prefix+item.Name+item.Separator+val)
			}
		}
	}
	return args
}

func (u *AppUI) getWidgetValue(item *Item, w fyne.CanvasObject) string {
	var val string
	switch item.Type {
	case "string", "number":
		if entry, ok := w.(*widget.Entry); ok {
			val = entry.Text
		} else if c, ok := w.(*fyne.Container); ok {
			for _, obj := range c.Objects {
				if entry, ok := obj.(*widget.Entry); ok {
					val = entry.Text
					break
				}
			}
		}
	case "bool":
		if w.(*widget.Check).Checked {
			return "true"
		}
		return ""
	case "choice":
		return w.(*widget.Select).Selected
	}
	// trim quotes
	if len(val) >= 2 {
		if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
			val = val[1 : len(val)-1]
		}
	}
	return val
}

func (u *AppUI) buildCommandLine() string {
	args := u.BuildArgs()
	quote := "'"
	if runtime.GOOS == "windows" {
		quote = "\""
	}
	var quoted []string
	for _, arg := range args {
		if strings.ContainsAny(arg, " \t'\"") {
			arg = quote + arg + quote
		}
		quoted = append(quoted, arg)
	}
	return u.app.Command.Path + " " + strings.Join(quoted, " ")
}

func (u *AppUI) showCommand() {
	cmdLine := u.buildCommandLine()
	entry := widget.NewEntry()
	entry.SetText(cmdLine)
	d := dialog.NewCustomConfirm("Command", "Copy", "Close", entry, func(copy bool) {
		if copy {
			u.window.Clipboard().SetContent(cmdLine)
		}
	}, u.window)
	size := u.window.Canvas().Size()
	d.Resize(fyne.NewSize(size.Width*2/3, size.Height*2/3))
	d.Show()
}

func parseImportance(color string) widget.Importance {
	switch color {
	case "high":
		return widget.HighImportance
	case "danger":
		return widget.DangerImportance
	case "warning":
		return widget.WarningImportance
	case "success":
		return widget.SuccessImportance
	case "low":
		return widget.LowImportance
	default:
		return widget.MediumImportance
	}
}

// 无间距垂直布局
type noSpaceVBox struct{}

func (n *noSpaceVBox) MinSize(objects []fyne.CanvasObject) fyne.Size {
	var w, h float32
	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		s := o.MinSize()
		if s.Width > w {
			w = s.Width
		}
		h += s.Height
	}
	return fyne.NewSize(w, h)
}

func (n *noSpaceVBox) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	var y float32
	for _, o := range objects {
		if !o.Visible() {
			continue
		}
		h := o.MinSize().Height
		o.Resize(fyne.NewSize(size.Width, h))
		o.Move(fyne.NewPos(0, y))
		y += h
	}
}

// 验证必填字段
func (u *AppUI) validateRequired() error {
	for _, item := range u.app.Items {
		if !item.Required || item.IsLabel() {
			continue
		}
		if !u.checkCondition(&item) {
			continue
		}
		val := u.getWidgetValue(&item, u.widgets[item.Name])
		if val == "" {
			label := item.Label
			if label == "" {
				label = item.Name
			}
			return fmt.Errorf("%s is required", label)
		}
	}
	return nil
}

// 验证所有字段
func (u *AppUI) validateAll() error {
	for _, item := range u.app.Items {
		if item.IsLabel() {
			continue
		}
		if !u.checkCondition(&item) {
			continue
		}
		val := u.getWidgetValue(&item, u.widgets[item.Name])
		if err := u.validateItem(&item, val); err != nil {
			return err
		}
	}
	return nil
}

// 验证单个字段
func (u *AppUI) validateItem(item *Item, val string) error {
	if val == "" {
		return nil
	}
	label := item.Label
	if label == "" {
		label = item.Name
	}

	// 正则验证
	if item.Validate != "" {
		re, err := regexp.Compile(item.Validate)
		if err != nil {
			return fmt.Errorf("invalid regex for %s", label)
		}
		if !re.MatchString(val) {
			return fmt.Errorf("%s format invalid", label)
		}
	}

	// 数字范围验证
	if item.Type == "number" && (item.Min != nil || item.Max != nil) {
		num, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return fmt.Errorf("%s must be a number", label)
		}
		if item.Min != nil {
			if min, ok := toFloat(item.Min); ok && num < min {
				return fmt.Errorf("%s must be >= %v", label, item.Min)
			}
		}
		if item.Max != nil {
			if max, ok := toFloat(item.Max); ok && num > max {
				return fmt.Errorf("%s must be <= %v", label, item.Max)
			}
		}
	}
	return nil
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case float64:
		return n, true
	}
	return 0, false
}

// 解析条件
func parseCondition(cond string) (field, op, value string) {
	if i := strings.Index(cond, "!="); i > 0 {
		return cond[:i], "!=", cond[i+2:]
	}
	if i := strings.Index(cond, "="); i > 0 {
		return cond[:i], "=", cond[i+1:]
	}
	return cond, "!=", ""
}

// 检查条件是否满足
func (u *AppUI) checkCondition(item *Item) bool {
	if item.Condition == "" {
		return true
	}
	field, op, expected := parseCondition(item.Condition)
	w := u.widgets[field]
	if w == nil {
		return true
	}
	var fieldItem *Item
	for i := range u.app.Items {
		if u.app.Items[i].Name == field {
			fieldItem = &u.app.Items[i]
			break
		}
	}
	if fieldItem == nil {
		return true
	}
	actual := u.getWidgetValue(fieldItem, w)
	if op == "!=" {
		return actual != expected
	}
	return actual == expected
}

// 设置条件监听
func (u *AppUI) setupConditions() {
	// 收集依赖关系: field -> []dependentItems
	deps := make(map[string][]*Item)
	for i := range u.app.Items {
		item := &u.app.Items[i]
		if item.Condition == "" {
			continue
		}
		field, _, _ := parseCondition(item.Condition)
		deps[field] = append(deps[field], item)
	}

	// 为每个被依赖的字段添加监听
	for field, items := range deps {
		w := u.widgets[field]
		if w == nil {
			continue
		}
		dependents := items
		updateFunc := func() {
			for _, item := range dependents {
				u.updateWidgetState(item)
			}
		}
		// 根据 widget 类型添加监听
		switch wt := w.(type) {
		case *widget.Entry:
			wt.OnChanged = func(s string) { updateFunc() }
		case *widget.Select:
			wt.OnChanged = func(s string) { updateFunc() }
		case *widget.Check:
			wt.OnChanged = func(b bool) { updateFunc() }
		case *fyne.Container:
			// picker 类型的 entry
			for _, obj := range wt.Objects {
				if entry, ok := obj.(*widget.Entry); ok {
					entry.OnChanged = func(s string) { updateFunc() }
					break
				}
			}
		}
		// 初始化状态
		updateFunc()
	}
}

// 更新 widget 启用/禁用状态
func (u *AppUI) updateWidgetState(item *Item) {
	w := u.widgets[item.Name]
	if w == nil {
		return
	}
	enabled := u.checkCondition(item)
	if dw, ok := w.(fyne.Disableable); ok {
		if enabled {
			dw.Enable()
		} else {
			dw.Disable()
		}
	}
	// 处理 container 类型 (picker)
	if c, ok := w.(*fyne.Container); ok {
		for _, obj := range c.Objects {
			if dw, ok := obj.(fyne.Disableable); ok {
				if enabled {
					dw.Enable()
				} else {
					dw.Disable()
				}
			}
		}
	}
}
