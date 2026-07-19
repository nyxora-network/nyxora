# 贡献 NYXORA

首先，感谢您考虑为 NYXORA 做出贡献！我们欢迎所有人的贡献。

## 行为准则

本项目遵循[行为准则](CODE_OF_CONDUCT.md)。参与即表示您同意遵守此准则。

## 如何贡献？

### 报告 Bug

提交 Bug 报告前：
- 检查 [issues](https://github.com/nyxora-network/nyxora/issues) 是否已报告
- 收集信息：操作系统版本、Go 版本、重现步骤、错误输出

**提交 Bug 报告**：打开 [新 issue](https://github.com/nyxora-network/nyxora/issues/new?template=bug_report.md)。

### 建议新功能

打开 [功能请求](https://github.com/nyxora-network/nyxora/issues/new?template=feature_request.md)，描述：
- 您要解决的问题
- 您设想的解决方案
- 考虑过的替代方案

### 添加新传输

1. 创建 `internal/transport/<name>.go`，实现 `Transport` 接口
2. 在 `internal/transport/registry.go` 中注册
3. 创建 `tunnels/<name>/`，包含安装脚本和清单
4. 在 `internal/transport/scoring.go` 中添加评分权重
5. 编写测试并运行 `make test`

### 改进 TUI

交互式 TUI 位于 `internal/interactive/`，使用：
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI 框架
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — 样式
- [Bubbles](https://github.com/charmbracelet/bubbles) — 组件（textinput、spinner、progress）

主题颜色使用 `theme.go` 中的 Catppuccin TrueColor 十六进制值。

## Pull Request 流程

1. Fork 仓库并从 `main` 创建您的分支
2. 运行 `make test` 和 `make vet` — 两者都必须通过
3. 为新功能添加测试
4. 需要时更新文档
5. 确保您的代码遵循现有约定
6. 提交 PR 并附上清晰的描述

### 提交信息风格

使用 conventional commit 格式：
- `feat:` — 新功能
- `fix:` — Bug 修复
- `refactor:` — 不含修复/功能的代码变更
- `docs:` — 仅文档
- `test:` — 测试添加/修复
- `style:` — 格式化、样式变更
- `chore:` — 维护、依赖

## 开发环境设置

```bash
# Fork 并 clone
git clone https://github.com/YOUR_USERNAME/nyxora.git
cd nyxora

# 添加 upstream remote
git remote add upstream https://github.com/nyxora-network/nyxora.git

# 创建功能分支
git checkout -b feat/your-feature

# 进行修改，然后：
make test
make vet
make build

# 提交并推送
git commit -m "feat: add your feature"
git push origin feat/your-feature
```

## 有问题？

打开 [讨论](https://github.com/nyxora-network/nyxora/discussions) 或加入我们的 [Telegram](https://t.me/NyxoraCore)。
