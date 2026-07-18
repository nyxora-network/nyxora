<div align="center">
  <img src="https://img.shields.io/badge/go-%3E%3D1.24-00ADD8?style=flat&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/license-MIT-blue?style=flat" alt="License">
  <img src="https://img.shields.io/badge/status-active-success?style=flat" alt="Status">
  <img src="https://img.shields.io/badge/PRs-welcome-brightgreen?style=flat" alt="PRs Welcome">
  <br>
  <img src="https://img.shields.io/badge/transports-12-ff69b4?style=flat" alt="12 Transports">
  <img src="https://img.shields.io/badge/platform-linux%20%7C%20macOS-lightgrey?style=flat" alt="Platform">
</div>

<br>

<div align="center">
  <a href="README.md"><img src="https://img.shields.io/badge/English-00ADD8?style=for-the-badge&logo=github&logoColor=white" alt="English"></a>
  <a href="README.fa.md"><img src="https://img.shields.io/badge/فارسی-DC143C?style=for-the-badge&logo=iran&logoColor=white" alt="فارسی"></a>
  <a href="README.ru.md"><img src="https://img.shields.io/badge/Русский-0052CC?style=for-the-badge&logo=russia&logoColor=white" alt="Русский"></a>
  <a href="README.zh.md"><img src="https://img.shields.io/badge/中文-FF4500?style=for-the-badge&logo=china&logoColor=white" alt="中文"></a>
  <a href="README.hi.md"><img src="https://img.shields.io/badge/हिन्दी-FF9933?style=for-the-badge&logo=india&logoColor=white" alt="हिन्दी"></a>
  <a href="README.es.md"><img src="https://img.shields.io/badge/Español-FF6B35?style=for-the-badge&logo=spain&logoColor=white" alt="Español"></a>
  <a href="README.ar.md"><img src="https://img.shields.io/badge/العربية-006233?style=for-the-badge&logo=saudi-arabia&logoColor=white" alt="العربية"></a>
</div>

<br>

<h1>NYXORA</h1>
  <h3>自适应隧道编排器</h3>
  <p>
    <b>自愈多传输 VPN/隧道管理器</b><br>
    在<i>一台</i>服务器上安装。连接到<i>任意</i>远程服务器。<br>
    无需远程代理。自动部署。自动故障转移。交互式 TUI。
  </p>
  <br>
  <p>
    <a href="#-功能特性">功能特性</a> •
    <a href="#-快速开始">快速开始</a> •
    <a href="#-一行安装">安装</a> •
    <a href="#-使用指南">使用指南</a> •
    <a href="#-架构">架构</a> •
    <a href="#-开发">开发</a>
  </p>
</div>

<br>

---

## ✨ 功能特性

<table>
<tr>
<td width="50%">

**🧠 自愈编排**
- 12 种隧道传输：WireGuard、OpenVPN、SSH、QUIC、FRP、Rathole、IPsec、Shadowsocks、Hysteria、Backhaul、TCP、WebSocket
- 自动故障转移 — 检测降级隧道，瞬间切换
- 5 种多路径调度模式（加权、最低延迟、最低丢包、均衡、全部激活）
- 实时评分引擎（延迟 + 丢包率 + 权重）

</td>
<td width="50%">

**🚀 零配置远程端**
- 远程服务器无需安装代理或软件
- 只需 SSH 访问（密码或密钥）
- 自动检测操作系统（Ubuntu、Debian、CentOS）
- 自动在远程服务器上安装隧道程序

</td>
</tr>
<tr>
<td width="50%">

**🖥️ 丰富的终端界面**
- 基于 Bubble Tea 的交互式 TUI，支持键盘导航
- 3 套专业色彩主题（Catppuccin Mocha、Tokyo Night、Catppuccin Latte）
- 实时仪表板，显示实时统计信息
- 动画渐变进度条
- ASCII 艺术标志启动画面
- 隧道拓扑视图
- 逐步连接向导

</td>
<td width="50%">

**🔐 企业级安全**
- 内核级 WireGuard VPN
- 支持 IPsec/strongSwan
- Shadowsocks 加密代理
- Hysteria 2（改进版 QUIC，抗审查）
- 自动生成密钥（密码、PSK、令牌）

</td>
</tr>
</table>

---

## 📦 一行安装

```bash
curl -L github.com/nyxorammd-lgtm/nyxora/releases/download/v0.2.0/nyxora_linux_amd64 -o /usr/local/bin/nyxora && chmod +x /usr/local/bin/nyxora
```

或使用 `wget`：

```bash
wget -q https://github.com/nyxorammd-lgtm/nyxora/releases/download/v0.2.0/nyxora_linux_amd64 -O /usr/local/bin/nyxora && chmod +x /usr/local/bin/nyxora
```

<details>
<summary><b>📋 手动安装（从源码）</b></summary>

```bash
# 安装依赖
sudo apt install golang-go git ssh sshpass wireguard curl
# 或：brew install go (macOS)

# 克隆仓库
git clone https://github.com/nyxorammd-lgtm/nyxora.git
cd nyxora

# 编译
make build

# 安装
sudo make install

# 验证
nyxora version
```
</details>

---

## 🚀 快速开始

```bash
# 1. 配置并检查依赖
nyxora install

# 2. 连接到远程服务器
nyxora connect 192.168.1.100 --user root --password your_password

# 3. 启动交互式 TUI
nyxora tui

# 4. 实时监控仪表板
nyxora dashboard
```

### 连接选项

```bash
nyxora connect <host> [options]

选项：
  --user, -u <name>       SSH 用户名（默认：root）
  --port, -p <port>       SSH 端口（默认：22）
  --password <pass>       SSH 密码
  --mode <mode>           服务器模式：full, lite, minimal
  --transports <list>     传输列表（覆盖 mode）
  --ports <pairs>         端口覆盖：wg=51820,ss=8388,...
```

#### 服务器模式

| 模式 | 传输 | 所需内存 |
|------|------|----------|
| `full` | 全部 11 种隧道 | 2GB+ |
| `lite` | 轻量选择 | 512MB–2GB |
| `minimal` | SSH + Shadowsocks 仅 | < 512MB |

---

## 🎮 交互式 TUI

NYXORA 拥有一个功能齐全的终端界面，基于 [Bubble Tea](https://github.com/charmbracelet/bubbletea) 和 [Lip Gloss](https://github.com/charmbracelet/lipgloss) 构建。

```
┌──────────────────────────────────────────────────────────┐
│  NYXORA v0.2.0                                          │
│  ────────────────────────────────────────────────────    │
│                                                          │
│  CPU: 0.5  ████░░░░░░░░░░░░░░░░                        │
│  RAM: 45%  ██████████░░░░░░░░░░                        │
│                                                          │
│  [1] C  连接到服务器                                     │
│  [2] D  仪表板                                           │
│  [3] I  服务器信息                                       │
│  [4] N  安装                                             │
│  [5] U  检查更新                                         │
│  [6] X  断开连接                                         │
│  [7] T  隧道拓扑                                         │
│  [8] H  帮助                                             │
│  [9] Q  退出                                             │
│                                                          │
│  ┌────────────────────────────────────────────────────┐  │
│  │  连接到远程服务器                                  │  │
│  └────────────────────────────────────────────────────┘  │
│  ↑↓ 导航  ↵ 选择  1/2/3 主题  s 状态  ? 帮助  q 退出   │
│  https://t.me/NyxoraCore                                 │
└──────────────────────────────────────────────────────────┘
```

### 快捷键

| 按键 | 操作 |
|------|------|
| `↑` / `↓` | 导航菜单 |
| `Enter` | 选择项目 |
| `Esc` | 返回 |
| `q` | 退出 / 返回菜单 |
| `1` | Catppuccin Mocha（深色） |
| `2` | Tokyo Night（深色） |
| `3` | Catppuccin Latte（浅色） |
| `s` | 切换状态栏 |
| `?` | 打开帮助屏幕 |
| `t` | 隧道拓扑视图 |

---

## 🏗️ 架构

```
┌─────────────────────────────────────────────────────────────────┐
│  nyxora（本地服务器）                                           │
│                                                                 │
│  ┌──────────────┐  ┌────────────────────────────────────────┐   │
│  │  编排器       │  │  传输管理器                            │   │
│  │  (Orchestrator)│  ┌─────┐ ┌─────┐ ┌─────┐ ┌──────────┐ │   │
│  │              │  │ │ WG  │ │ SSH │ │ OVPN│ │ Hysteria  │ │   │
│  │  Init        │  │ ├─────┤ ├─────┤ ├─────┤ ├──────────┤ │   │
│  │  Connect     │  │ │ FRP │ │QUIC │ │Rhole│ │ Backhaul  │ │   │
│  │  Monitor     │  │ ├─────┤ ├─────┤ ├─────┤ ├──────────┤ │   │
│  │  Failover    │  │ │ TCP │ │IPsec│ │ SS  │ │          │ │   │
│  └──────┬───────┘  │ └─────┘ └─────┘ └─────┘ └──────────┘ │   │
│         │          └────────────────────────────────────────┘   │
│         │                                                       │
│  ┌──────┴───────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │  多路径调度器  │  │  故障转移引擎  │  │  仪表板 / TUI       │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
│         │                                                       │
│         │ SSH + 远程部署                                         │
│         ▼                                                       │
│  ┌────────────────────────────────────────────────────────┐     │
│  │  远程服务器（无需代理）                                  │     │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────┐  │     │
│  │  │ frps     │  │ rathole  │  │ hysteria │  │ WG    │  │     │
│  │  │ :7000    │  │ :2333    │  │ :8444    │  │:51820 │  │     │
│  │  ├──────────┤  ├──────────┤  ├──────────┤  ├───────┤  │     │
│  │  │ ss-srv   │  │ backhaul │  │ openvpn  │  │ SSHd  │  │     │
│  │  │ :8388    │  │ :3080    │  │ :1194    │  │ :22   │  │     │
│  │  └──────────┘  └──────────┘  └──────────┘  └───────┘  │     │
│  └────────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────────┘
```

### 连接流程

```
nyxora connect 91.107.243.237 --user root --password ...

  1. PING          → 测量延迟和丢包
  2. SSH           → 验证远程服务器
  3. DETECT OS     → 检测 Ubuntu / Debian / CentOS
  4. INSTALL       → 在远程部署隧道程序
  5. WG KEY        → 本地生成 WireGuard 密钥对
  6. REMOTE WG     → SSH：配置 + wg-quick up + iptables
  7. LOCAL WG      → 使用远程公钥启动 wg-quick up nyxora0
  8. PROVISION     → 启动守护进程：frps, rathole, ss, hys, backhaul
  9. ALL-ACTIVE    → 同时测试和激活所有隧道
  10. MONITOR      → 每 10 秒：ping、评分、故障转移检查
```

---

## 📋 命令

| 命令 | 说明 |
|------|------|
| `nyxora install` | 配置并检查依赖 |
| `nyxora connect <host>` | 连接到远程服务器 |
| `nyxora disconnect` | 关闭所有隧道 |
| `nyxora status` | 显示连接状态 |
| `nyxora dashboard` | 实时终端仪表板 |
| `nyxora tui` | 交互式 Bubble Tea 菜单 |
| `nyxora update` | 检查更新 |
| `nyxora server` | 显示服务器信息和推荐模式 |
| `nyxora version` | 显示版本 |
| `nyxora daemon` | 作为后台服务运行 |
| `nyxora help` | 显示帮助 |

---

## 🔧 配置

### 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `NYXORA_SS_PASSWORD` | Shadowsocks 密码 | 自动生成 |
| `NYXORA_SS_METHOD` | Shadowsocks 加密方法 | `aes-256-gcm` |
| `NYXORA_RATHOLE_TOKEN` | Rathole 认证令牌 | 自动生成 |
| `NYXORA_HYSTERIA_AUTH` | Hysteria 认证密码 | 自动生成 |
| `NYXORA_BACKHAUL_TOKEN` | Backhaul 认证令牌 | 自动生成 |
| `NYXORA_IPSEC_PSK` | IPsec 预共享密钥 | 自动生成 |
| `NYXORA_ALL_ACTIVE` | 同时启用所有隧道 | `false` |

### 配置文件

配置存储在 `/etc/nyxora/config.json`（通过 `nyxora install` 自动生成）。

---

## 📦 传输列表

| # | 名称 | 端口 | 协议 | 类别 | 基础评分 | 权重 |
|---|------|------|------|------|----------|------|
| 1 | **wireguard** | 51820 | UDP | VPN | 95 | 30 |
| 2 | **openvpn** | 1194 | UDP | VPN | 75 | 10 |
| 3 | **ssh** | 22 | TCP | 隧道 | 60 | 5 |
| 4 | **quic** | 9923 | UDP | 隧道 | 80 | 15 |
| 5 | **frp** | 7000 | TCP | 中继 | 70 | 10 |
| 6 | **rathole** | 2333 | TCP | 中继 | 85 | 12 |
| 7 | **ipsec** | 500 | UDP | VPN | 70 | 5 |
| 8 | **shadowsocks** | 8388 | TCP | 代理 | 55 | 3 |
| 9 | **hysteria** | 8444 | UDP | 隧道 | 90 | 12 |
| 10| **backhaul** | 3080 | TCP | 中继 | 82 | 10 |
| 11| **tcp** | 9924 | TCP | 隧道 | 50 | 3 |

### 多路径调度模式

| 模式 | 说明 |
|------|------|
| `weighted` | 根据隧道权重分配流量 |
| `lowest-latency` | 通过最低延迟路径路由所有流量 |
| `lowest-loss` | 通过最低丢包路径路由所有流量 |
| `even` | 在所有活跃隧道间均衡分配 |
| `all` | 所有隧道同时激活 |

---

## 🧑‍💻 开发

### 前提条件

- Go 1.24+
- Linux 或 macOS
- `ssh`, `sshpass`, `wg`, `curl`, `ping`

### 开始开发

```bash
git clone https://github.com/nyxorammd-lgtm/nyxora.git
cd nyxora

# 编译
make build

# 运行测试
make test

# 代码检查
make vet

# 本地运行
./nyxora version
```

### 项目结构

```
nyxora/
├── cmd/
│   ├── nyxora/           # CLI 入口
│   └── quic-server/      # QUIC 回显服务器
├── internal/
│   ├── config/           # 配置、密钥、服务器信息
│   ├── dashboard/        # ANSI 终端仪表板
│   ├── failover/         # 自动故障转移引擎
│   ├── interactive/      # Bubble Tea TUI
│   ├── monitor/          # 基于 ping 的监控
│   ├── multipath/        # 多路径调度器
│   ├── orchestrator/     # 核心引擎
│   ├── packager/         # 打包工具
│   ├── remote/           # SSH 客户端
│   ├── routing/          # 评分引擎
│   └── transport/        # 12 种传输实现
├── tunnels/              # 各隧道的安装脚本
├── Makefile              # 构建、测试、安装
├── Dockerfile            # Docker 构建
└── install.sh            # 一键安装脚本
```

### Makefile 目标

| 目标 | 说明 |
|------|------|
| `make build` | 编译二进制文件 |
| `make test` | 运行所有测试 |
| `make vet` | 运行 go vet |
| `make run` | 编译并运行 |
| `make clean` | 删除二进制文件和缓存 |
| `make install` | 安装到 `/usr/local/bin` |
| `make daemon` | 设置 systemd 服务 |
| `make tunnels` | 打包隧道脚本 |

---

## 🤝 贡献

我们欢迎贡献！请阅读我们的[贡献指南](CONTRIBUTING.md)。

**贡献方式：**
- 通过 [GitHub Issues](https://github.com/nyxorammd-lgtm/nyxora/issues) 报告错误
- 建议新的传输类型
- 改进 TUI / 仪表板
- 添加更多操作系统的支持
- 编写测试和文档
- 针对开放 Issue 提交 PR

---

## 📄 许可证

本项目采用 **MIT 许可证** 发布。

---

<div align="center">
  <br>
  <p>
    <a href="https://t.me/NyxoraCore">Telegram 频道</a> •
    <a href="https://github.com/nyxorammd-lgtm/nyxora/issues">报告错误</a> •
    <a href="https://github.com/nyxorammd-lgtm/nyxora/issues">功能请求</a>
  </p>
  <p>
    <sub>用 ❤️ 和 Go 与 Bubble Tea 构建</sub>
  </p>
</div>
