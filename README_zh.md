# cliface

轻量级命令行工具 GUI 封装器。通过 TOML 配置文件定义表单，生成图形界面来构建和执行 shell 命令。

## 特性

- 基于 TOML 配置驱动的 UI 生成
- 多种输入类型：字符串、数字、布尔、选择框
- 文件和目录选择器
- 多值字段（支持增删按钮）
- 字段验证（必填、正则、范围）
- 条件字段显示/隐藏
- 多种执行模式：可见窗口、弹窗输出、实时流式输出
- 跨平台支持（macOS、Windows、Linux）

[English](README.md)

## 仓库地址

https://github.com/listeng/cliface

## 安装

```bash
go build -o cliface .
```

## 使用

```bash
cliface -c examples/curl.toml
```

不指定 `-c` 时，默认读取可执行文件同目录下的 `config.toml`。

## 配置示例

```toml
title = "我的工具"
width = 500
height = 400

[[apps]]
[apps.command]
path = "/usr/bin/ffmpeg"
name = "视频转换"
args = ["-y"]
mode = "hidden"      # hidden | visible
output = "dialog"    # dialog | realtime
debug = true

[[apps.items]]
text = "将视频转换为 MP4 格式"  # 纯文本标签

[[apps.items]]
name = "i"
short = true
type = "string"
label = "输入文件"
picker = "file"
separator = " "

[[apps.items]]
name = "crf"
type = "number"
label = "质量"
default = 23

[[apps.items]]
name = "preset"
type = "choice"
label = "速度"
choices = ["ultrafast", "fast", "medium", "slow"]
default = "medium"
```

## 配置说明

### 全局配置

| 字段 | 说明 | 默认值 |
|------|------|--------|
| title | 窗口标题 | 命令名或 "cliface" |
| width | 窗口宽度 | 400 |
| height | 窗口高度 | 300 |

### Command 配置

| 字段 | 说明 |
|------|------|
| path | 可执行文件路径 |
| name | 显示名称（多 app 时作为 tab 标题） |
| args | 固定参数 |
| mode | `hidden` 隐藏执行 / `visible` 可见窗口 |
| output | `dialog` 完成后弹窗 / `realtime` 实时窗口 / `realtime-console` 终端输出 |
| debug | 显示"查看命令"按钮 |
| run_text / run_color | 运行按钮文字和颜色 (high/danger/warning/success/low) |
| debug_text / debug_color | 调试按钮文字和颜色 |
| env | 环境变量，键值对形式 |

### Item 配置

| 字段 | 说明 |
|------|------|
| text | 纯文本标签（设置后忽略其他字段） |
| name | 参数名 |
| type | `string` / `number` / `bool` / `choice` |
| short | true 时使用单横线 `-name` |
| positional | true 时为位置参数（无前缀） |
| label | 显示标签 |
| description | 字段说明 |
| default | 默认值 |
| choices | choice 类型的选项列表 |
| picker | `file` 或 `directory` 选择器 |
| picker_text | 自定义选择器按钮文字 |
| separator | 参数分隔符，`" "` 为空格，`"none"` 为无分隔符，默认 `=` |
| multi | 允许多值输入（带增删按钮） |
| required | 必填字段，运行前验证 |
| validate | 正则表达式验证 |
| min / max | 数字范围验证 |
| condition | 条件显示/启用（如 `field=value` 或 `field!=value`） |

## License

MIT
