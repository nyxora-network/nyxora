<div align="center">
  <img src="https://img.shields.io/badge/go-%3E%3D1.24-00ADD8?style=flat&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/license-MIT-blue?style=flat" alt="مجوز">
  <img src="https://img.shields.io/badge/status-active-success?style=flat" alt="وضعیت">
  <img src="https://img.shields.io/badge/PRs-welcome-brightgreen?style=flat" alt="PRs Welcome">
  <img src="https://img.shields.io/github/stars/nyxora-network/nyxora?style=flat&logo=github" alt="Stars">
  <br>
  <a href="https://github.com/nyxora-network/nyxora/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/nyxora-network/nyxora/ci.yml?branch=main&label=CI&logo=github" alt="CI"></a>
  <img src="https://img.shields.io/badge/transports-12-ff69b4?style=flat" alt="۱۲ ترنسپورت">
  <img src="https://img.shields.io/badge/platform-linux%20%7C%20macOS-lightgrey?style=flat" alt="پلتفرم">
  <img src="https://img.shields.io/github/v/release/nyxora-network/nyxora?style=flat" alt="Release">
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
  <h3>تست تونل‌ها رو یکی یکی متوقف کن — از NYXORA استفاده کن</h3>
  <p>
    <b>مدیریت تونل و VPN خود-ترمیم با چندین ترنسپورت</b><br>
    روی <i>یک</i> سرور نصب کن. به <i>هر</i> سرور راه دوری وصل شو.<br>
    بدون نیاز به عامل (Agent). نصب خودکار. Failover خودکار. TUI تعاملی.
  </p>
  <br>
  <p>
    <a href="#-ویژگی‌ها">ویژگی‌ها</a> •
    <a href="#-شروع-سریع">شروع سریع</a> •
    <a href="#-نصب-تک-خطی">نصب</a> •
    <a href="#-راهنما">راهنما</a> •
    <a href="#-معماری">معماری</a> •
    <a href="#-توسعه">توسعه</a>
  </p>
</div>

<br>

---

## ✨ ویژگی‌ها

<table>
<tr>
<td width="50%">

**🧠 ارکستریتور خود-ترمیم**
- ۱۲ ترنسپورت تونل: WireGuard, OpenVPN, SSH, QUIC, FRP, Rathole, IPsec, Shadowsocks, Hysteria, Backhaul, TCP, WebSocket
- Failover خودکار — تشخیص تونل‌های ضعیف و سوئیچ آنی
- ۵ حالت زمان‌بندی چندمسیره (weighted, lowest-latency, lowest-loss, even, all-active)
- موتور امتیازدهی بلادرنگ (تأخیر + افت بسته + وزن)

</td>
<td width="50%">

**🚀 راه دور بدون تنظیمات**
- بدون نیاز به عامل (Agent) یا نرم‌افزار روی سرور مقصد
- فقط دسترسی SSH (رمز یا کلید)
- تشخیص خودکار OS (Ubuntu, Debian, CentOS)
- نصب خودکار باینری‌های تونل روی سرور راه دور

</td>
</tr>
<tr>
<td width="50%">

**🖥️ رابط کاربری ترمینال غنی**
- TUI تعاملی با Bubble Tea و ناوبری صفحه‌کلید
- ۳ تم رنگی حرفه‌ای (Catppuccin Mocha, Tokyo Night, Catppuccin Latte)
- داشبورد زنده با آمار بلادرنگ
- نوار پیشرفت گرادیان متحرک
- لوگوی ASCII در صفحه شروع
- نمای توپولوژی تونل
- ویزارد اتصال گام‌به‌گام

</td>
<td width="50%">

**🔐 امنیت در سطح سازمانی**
- VPN WireGuard در سطح کرنل
- پشتیبانی IPsec/strongSwan
- پروکسی رمزنگاری‌شده Shadowsocks
- Hysteria 2 (QUIC اصلاح‌شده با ضدسانسور)
- پوشش TLS با تولید خودکار گواهی Self-Signed
- تولید خودکار رمزهای عبور (Passwords, PSKs, Tokens)

</td>
<td width="50%">

**📊 مانیتورینگ در سطح Production**
- اندپوینت Prometheus metrics (/metrics)
- رزولور DNS داخلی با کش
- Rate Limiting به ازای هر ترنسپورت (Token Bucket)
- Hot-Reload کانفیگ (File Watcher)

</td>
</tr>
</table>

</td>
</tr>
</table>

---

## 📦 نصب تک خطی

```bash
curl -L github.com/nyxora-network/nyxora/releases/download/v0.2.0/nyxora_linux_amd64 -o /usr/local/bin/nyxora && chmod +x /usr/local/bin/nyxora
```

یا با `wget`:

```bash
wget -q https://github.com/nyxora-network/nyxora/releases/download/v0.2.0/nyxora_linux_amd64 -O /usr/local/bin/nyxora && chmod +x /usr/local/bin/nyxora
```

<details>
<summary><b>📋 نصب دستی (از سورس)</b></summary>

```bash
# پیش‌نیازها
sudo apt install golang-go git ssh sshpass wireguard curl
# یا: brew install go (macOS)

# Clone
git clone https://github.com/nyxora-network/nyxora.git
cd nyxora

# Build
make build

# Install
sudo make install

# Verify
nyxora version
```
</details>

---

## 🚀 شروع سریع

```bash
# 1. تنظیم کانفیگ و بررسی پیش‌نیازها
nyxora install

# 2. اتصال به سرور راه دور
nyxora connect 192.168.1.100 --user root --password your_password

# 3. اجرای TUI تعاملی
nyxora tui

# 4. داشبورد مانیتورینگ زنده
nyxora dashboard
```

### گزینه‌های اتصال

```bash
nyxora connect <host> [options]

Options:
  --user, -u <name>       نام کاربری SSH (پیش‌فرض: root)
  --port, -p <port>       پورت SSH (پیش‌فرض: 22)
  --password <pass>       رمز SSH
  --mode <mode>           حالت سرور: full, lite, minimal
  --transports <list>     لیست ترنسپورت‌ها (جایگزین mode)
  --ports <pairs>         تغییر پورت: wg=51820,ss=8388,...
```

#### حالت‌های سرور

| حالت | ترنسپورت‌ها | رم مورد نیاز |
|------|-------------|--------------|
| `full` | هر ۱۱ تونل | 2GB+ |
| `lite` | انتخاب سبک‌وزن | 512MB–2GB |
| `minimal` | SSH + Shadowsocks | < 512MB |

---

## 🎮 TUI تعاملی

NYXORA دارای یک رابط کاربری ترمینال کامل است که با [Bubble Tea](https://github.com/charmbracelet/bubbletea) و [Lip Gloss](https://github.com/charmbracelet/lipgloss) ساخته شده.

```
┌──────────────────────────────────────────────────────────┐
│  NYXORA v0.2.0                                          │
│  ────────────────────────────────────────────────────    │
│                                                          │
│  CPU: 0.5  ████░░░░░░░░░░░░░░░░                        │
│  RAM: 45%  ██████████░░░░░░░░░░                        │
│                                                          │
│  [1] C  اتصال به سرور                                   │
│  [2] D  داشبورد                                         │
│  [3] I  اطلاعات سرور                                    │
│  [4] N  نصب                                             │
│  [5] U  بررسی بروزرسانی                                 │
│  [6] X  قطع اتصال                                       │
│  [7] T  توپولوژی تونل                                   │
│  [8] H  راهنما                                          │
│  [9] Q  خروج                                            │
│                                                          │
│  ┌────────────────────────────────────────────────────┐  │
│  │  اتصال به سرور راه دور                             │  │
│  └────────────────────────────────────────────────────┘  │
│  ↑↓ حرکت  ↵ انتخاب  1/2/3 تم  s وضعیت  ? راهنما  q خروج │
│  https://t.me/NyxoraCore                                 │
└──────────────────────────────────────────────────────────┘
```

### میانبرهای صفحه‌کلید

| کلید | عملکرد |
|------|--------|
| `↑` / `↓` | حرکت در منو |
| `Enter` | انتخاب |
| `Esc` | برگشت |
| `q` | خروج / برگشت به منو |
| `1` | Catppuccin Mocha (تیره) |
| `2` | Tokyo Night (تیره) |
| `3` | Catppuccin Latte (روشن) |
| `s` | نمایش/مخفی کردن نوار وضعیت |
| `?` | صفحه راهنما |
| `t` | نمای توپولوژی تونل |

---

## 🏗️ معماری

```
┌─────────────────────────────────────────────────────────────────┐
│  nyxora (سرور محلی)                                             │
│                                                                 │
│  ┌──────────────┐  ┌────────────────────────────────────────┐   │
│  │  Orchestrator │  │  Transport Manager                     │   │
│  │              │  │  ┌─────┐ ┌─────┐ ┌─────┐ ┌──────────┐ │   │
│  │  Init        │  │  │ WG  │ │ SSH │ │ OVPN│ │ Hysteria  │ │   │
│  │  Connect     │  │  ├─────┤ ├─────┤ ├─────┤ ├──────────┤ │   │
│  │  Monitor     │  │  │ FRP │ │QUIC │ │Rhole│ │ Backhaul  │ │   │
│  │  Failover    │  │  ├─────┤ ├─────┤ ├─────┤ ├──────────┤ │   │
│  │              │  │  │ TCP │ │IPsec│ │ SS  │ │          │ │   │
│  └──────┬───────┘  │  └─────┘ └─────┘ └─────┘ └──────────┘ │   │
│         │          └────────────────────────────────────────┘   │
│         │                                                       │
│  ┌──────┴───────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │  Multipath   │  │  Failover    │  │  Dashboard / TUI     │  │
│  │  Scheduler   │  │  Engine      │  │  (Bubble Tea)        │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
│         │                                                       │
│         │ SSH + provisioning                                    │
│         ▼                                                       │
│  ┌────────────────────────────────────────────────────────┐     │
│  │  سرور راه دور (بدون نیاز به Agent)                     │     │
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

### مراحل اتصال

```
nyxora connect 91.107.243.237 --user root --password ...

  1.  PING          → اندازه‌گیری تأخیر و افت بسته
  2.  SSH           → احراز هویت در سرور راه دور
  3.  DETECT OS     → تشخیص Ubuntu / Debian / CentOS
  4.  INSTALL       → نصب باینری‌های تونل روی سرور راه دور
  5.  WG KEY        → تولید کلید WireGuard به صورت محلی
  6.  REMOTE WG     → SSH: کانفیگ + wg-quick up + iptables
  7.  LOCAL WG      → wg-quick up nyxora0 با کلید عمومی راه دور
  8.  PROVISION     → راه‌اندازی: frps, rathole, ss, hys, backhaul
  9.  ALL-ACTIVE    → تست و فعالسازی همه تونل‌ها همزمان
  10. MONITOR       → هر ۱۰ ثانیه: پینگ، امتیاز، بررسی failover
```

---

## 📸 اسکرین‌شات‌های وب داشبورد

NYXORA شامل یک **وب داشبورد** برای مدیریت تونل است. از آدرس `http://YOUR_SERVER:8080/tunnel.html` دسترسی داشته باشید

### نقشه تونل
![نقشه تونل](screenshots/01-tunnel-map.png)
*نقشه بصری اتصال بین سرورها با خطوط تونل متحرک*

### اتصال جدید
![اتصال](screenshots/02-connect.png)
*فرم اتصال تونل با احراز هویت SSH و انتخاب ترنسل*

### تونل‌های فعال
![تونل‌ها](screenshots/03-tunnels.png)
*لیست تونل‌های فعال با وضعیت، تأخیر و کنترل‌ها*

### تشخیص خودکار بهترین ترنسل
![تشخیص خودکار](screenshots/04-auto-detect.png)
*پینگ خودکار 14 ترنسل برای انتخاب بهترین*

### سازنده کانفیگ Xray
![کانفیگ Xray](screenshots/05-xray-config.png)
*تولید کانفیگ Xray برای VLESS, VMess, Trojan, Shadowsocks*

### نصب پنل‌ها
![پنل‌ها](screenshots/06-panels.png)
*نصب یک خطی Marzban, 3X-UI, Xray Core*

### نمودارهای عملکرد
![نمودارها](screenshots/07-charts.png)
*نمودارهای زنده تأخیر، پهنای باند، افت بسته و امتیاز*

### امنیت و ممیزی
![امنیت](screenshots/08-security.png)
*وضعیت SSL، اثر انگشت سرور، تلاش‌های ناموفق، QoS، آپتایم*

### ترمینال
![ترمینال](screenshots/09-terminal.png)
*رابط ترمینال کامل با دستورات NYXORA*

**[مستندات کامل وب داشبورد →](web/README.md)**

---

## 📋 دستورات

| دستور | توضیح |
|-------|-------|
| `nyxora install` | تنظیم کانفیگ و بررسی پیش‌نیازها |
| `nyxora connect <host>` | اتصال به سرور راه دور |
| `nyxora disconnect` | قطع همه تونل‌ها |
| `nyxora status` | نمایش وضعیت اتصال |
| `nyxora dashboard` | داشبورد ترمینال زنده |
| `nyxora tui` | منوی تعاملی Bubble Tea |
| `nyxora update` | بررسی بروزرسانی |
| `nyxora server` | اطلاعات سرور و حالت پیشنهادی |
| `nyxora version` | نمایش نسخه |
| `nyxora daemon` | اجرا به عنوان سرویس پس‌زمینه |
| `nyxora help` | راهنما |

---

## 🔧 تنظیمات

### متغیرهای محیطی

| متغیر | توضیح | پیش‌فرض |
|-------|-------|---------|
| `NYXORA_SS_PASSWORD` | رمز Shadowsocks | تولید خودکار |
| `NYXORA_SS_METHOD` | الگوریتم Shadowsocks | `aes-256-gcm` |
| `NYXORA_RATHOLE_TOKEN` | توکن Rathole | تولید خودکار |
| `NYXORA_HYSTERIA_AUTH` | رمز Hysteria | تولید خودکار |
| `NYXORA_BACKHAUL_TOKEN` | توکن Backhaul | تولید خودکار |
| `NYXORA_IPSEC_PSK` | کلید IPsec | تولید خودکار |
| `NYXORA_ALL_ACTIVE` | فعالسازی همزمان همه تونل‌ها | `false` |

### فایل کانفیگ

تنظیمات در `/etc/nyxora/config.json` ذخیره می‌شود (تولید خودکار با `nyxora install`).

---

## 📦 ترنسپورت‌ها

| # | نام | پورت | پروتکل | دسته‌بندی | امتیاز پایه | وزن |
|---|------|------|--------|-----------|-------------|------|
| 1 | **wireguard** | 51820 | UDP | VPN | 95 | 30 |
| 2 | **openvpn** | 1194 | UDP | VPN | 75 | 10 |
| 3 | **ssh** | 22 | TCP | تونل | 60 | 5 |
| 4 | **quic** | 9923 | UDP | تونل | 80 | 15 |
| 5 | **frp** | 7000 | TCP | رله | 70 | 10 |
| 6 | **rathole** | 2333 | TCP | رله | 85 | 12 |
| 7 | **ipsec** | 500 | UDP | VPN | 70 | 5 |
| 8 | **shadowsocks** | 8388 | TCP | پروکسی | 55 | 3 |
| 9 | **hysteria** | 8444 | UDP | تونل | 90 | 12 |
| 10| **backhaul** | 3080 | TCP | رله | 82 | 10 |
| 11| **tcp** | 9924 | TCP | تونل | 50 | 3 |
| 12| **websocket** | 9925 | TCP | تونل | 70 | 8 |

### حالت‌های زمان‌بندی چندمسیره

| حالت | توضیح |
|------|-------|
| `weighted` | توزیع ترافیک بر اساس وزن تونل‌ها |
| `lowest-latency` | مسیریابی همه ترافیک از کم‌تأخیرترین مسیر |
| `lowest-loss` | مسیریابی همه ترافیک از کم‌افت‌ترین مسیر |
| `even` | توزیع برابر بین همه تونل‌های فعال |
| `all` | همه تونل‌ها همزمان فعال |

---

## 🧑‍💻 توسعه

### پیش‌نیازها

- Go 1.24+
- لینوکس یا macOS
- `ssh`, `sshpass`, `wg`, `curl`, `ping`

### راه‌اندازی

```bash
git clone https://github.com/nyxora-network/nyxora.git
cd nyxora

# Build
make build

# تست
make test

# Vet
make vet

# اجرای محلی
./nyxora version
```

### ساختار پروژه

```
nyxora/
├── cmd/
│   ├── nyxora/           # ورودی CLI
│   └── quic-server/      # سرور QUIC
├── internal/
│   ├── config/           # کانفیگ، رمزها، اطلاعات سرور
│   ├── dashboard/        # داشبورد ترمینال ANSI
│   ├── failover/         # موتور Failover خودکار
│   ├── interactive/      # TUI با Bubble Tea
│   ├── monitor/          # مانیتورینگ با پینگ
│   ├── multipath/        # زمان‌بند چندمسیره
│   ├── orchestrator/     # موتور اصلی
│   ├── packager/         # ابزار بسته‌بندی
│   ├── remote/           # کلاینت SSH
│   ├── routing/          # موتور امتیازدهی
│   └── transport/        # ۱۱ ترنسپورت
├── tunnels/              # اسکریپت‌های نصب تونل
├── Makefile              # Build, test, install
├── Dockerfile            # Docker build
└── install.sh            # نصب تک خطی
```

### اهداف Makefile

| هدف | توضیح |
|-----|-------|
| `make build` | ساخت باینری |
| `make test` | اجرای تست‌ها |
| `make vet` | اجرای go vet |
| `make run` | ساخت و اجرا |
| `make clean` | حذف باینری و کش |
| `make install` | نصب در `/usr/local/bin` |
| `make daemon` | راه‌اندازی سرویس systemd |
| `make tunnels` | بسته‌بندی اسکریپت‌های تونل |

---

## 🤝 مشارکت

> **ما عاشق مشارکت هستیم!** چه در حال رفع اشتباه تایپی باشید، گزارش باگ، یا افزودن ترنسپورت جدید — هر مشارکتی ارزشمند است. NYXORA پروژه‌ای جامعه‌محور است و ما می‌خواهیم شما بخشی از آن باشید.

**[راهنمای کامل مشارکت را بخوانید →](CONTRIBUTING.md)**

---

### 🎯 راه‌های مشارکت

| نوع | سختی | توضیح |
|-----|------|-------|
| 🐛 **گزارش باگ** | آسان | چیزی شکسته پیدا کردید؟ به ما بگویید! |
| 📝 **مستندات** | آسان | رفع اشتباهات تایپی، بهبود راهنماها، افزودن مثال |
| 🧪 **تست‌ها** | متوسط | افزودن پوشش تست برای کد موجود |
| 🎨 **بهبود TUI** | متوسط | بهبود تم‌ها، layout، انیمیشن‌ها |
| 🔧 **رفع باگ** | متوسط | رفع issuesهای گزارش شده |
| 🚀 **ویژگی‌های جدید** | سخت | افزودن دستورات، حالت‌ها، گزینه‌های جدید |
| 🌐 **ترنسپورت‌های جدید** | سخت | پیاده‌سازی پروتکل‌های تونل جدید

**اولین بار است که مشارکت می‌کنید؟** به issuesهایی با برچسب [`good first issue`](https://github.com/nyxora-network/nyxora/labels/good%20first%20issue) نگاه کنید — عالی برای شروع هستند!

---

### 🐛 گزارش باگ

**قبل از ارسال:**
1. [Issuesهای موجود](https://github.com/nyxora-network/nyxora/issues) را جستجو کنید — ممکن است باگ شما قبلاً گزارش شده باشد
2. سعی کنید در نصب تمیز بازتولید کنید (`nyxora install` جدید)
3. به آخرین نسخه آپدیت کنید (`nyxora update`) — ممکن است باگ قبلاً رفع شده باشد

**نحوه ارسال گزارش باگ عالی:**

1. روی [**گزارش باگ جدید**](https://github.com/nyxora-network/nyxora/issues/new?template=bug_report.md) کلیک کنید
2. **تمام بخش‌ها** را پر کنید (هرچه جزئیات بیشتر باشد، سریع‌تر رفع می‌شود):

| بخش | چه چیزی وارد کنید |
|-----|-------------------|
| **Describe the Bug** | چه اتفاقی افتاد؟ چه انتظاری داشتید؟ |
| **To Reproduce** | دستورات دقیق، مرحله به مرحله |
| **Expected Behavior** | چه اتفاقی باید می‌افتاد |
| **Terminal Output** | خروجی کامل خطا را کپی کنید (از ```code blocks``` استفاده کنید) |
| **Environment** | سیستم‌عامل، نسخه Go، نسخه NYXORA، سیستم‌عامل سرور ریموت |

**📋 مثال گزارش باگ:**

```markdown
**Describe the Bug**
دستور `nyxora connect` با خطای "connection refused" هنگام اتصال به سرور CentOS 8 شکست می‌خورد

**To Reproduce**
1. دستور `nyxora install` را اجرا کنید
2. دستور `nyxora connect 192.168.1.50 --user root --password mypassword` را اجرا کنید
3. خطا در مرحله 5 (INSTALL) ظاهر می‌شود

**Expected Behavior**
باینری‌های تونل باید با موفقیت روی CentOS 8 نصب شوند

**Terminal Output**
```bash
[STEP 1] PING: measuring latency...
[STEP 2] SSH: authenticating...
[STEP 3] DETECT OS: CentOS 8 detected
[STEP 4] INSTALL: deploying tunnel binaries...
Error: ssh: connect to host 192.168.1.50 port 22: connection refused
```

**Environment**
- سیستم‌عامل (لوکال): Ubuntu 22.04
- نسخه Go: 1.24.0
- نسخه NYXORA: 0.2.0
- سیستم‌عامل سرور ریموت: CentOS 8
- RAM: 512MB
```

> **💡 نکته:** اسکرین‌شات و ضبط ترمینال خیلی کمک می‌کند! می‌توانید از [asciinema](https://asciinema.org/) برای ضبط ترمینال استفاده کنید.

---

### 📋 قوانین Pull Request

#### 🚀 شروع سریع برای مشارکت‌کنندگان

```bash
# 1. fork و clone
git clone https://github.com/YOUR_USERNAME/nyxora.git
cd nyxora

# 2. ساخت شاخه feature
git checkout -b feat/your-feature-name

# 3. تغییرات و سپس تست
make test
make vet

# 4. commit و push
git commit -m "feat: add your feature"
git push origin feat/your-feature-name

# 5. باز کردن PR در GitHub
```

#### 📌 قبل از شروع

| ✅ انجام دهید | ❌ انجام ندهید |
|---------------|---------------|
| ابتدا issue باز کنید برای بحث | ارسال PRهای بزرگ بدون بحث قبلی |
| PRها را متمرکز نگه دارید (یک feature/fix) | ترکیب چند feature در یک PR |
| تست برای functionality جدید بنویسید | ارسال کد بدون تست |
| از استایل کد موجود پیروی کنید | بازنویسی همه چیز به روش خودتان |
| مستندات را به‌روز کنید | فراموش کردن به‌روزرسانی README |

#### 🏷️ قرارداد نام‌گذاری شاخه

| پیشوند | کی استفاده کنیم | مثال |
|--------|----------------|------|
| `feat/` | ویژگی جدید یا ترنسپورت | `feat/add-vless-transport` |
| `fix/` | رفع باگ | `fix/failover-race-condition` |
| `docs/` | فقط مستندات | `docs/update-api-reference` |
| `test/` | افزودن تست | `test/add-wireguard-unit-tests` |
| `refactor/` | بازسازی کد | `refactor/extract-ssh-client` |
| `hotfix/` | رفع اضطراری در production | `hotfix/crash-on-startup` |

#### ✅ الزامات PR (باید pass شوند)

PR شما **تا زمانی که همه اینها سبز نباشند، merge نمی‌شود**:

| چک | دستور | وضعیت |
|----|-------|-------|
| تست‌ها | `make test` | ✅ همه رد شوند |
| لینتینگ | `make vet` | ✅ بدون هشدار |
| استایل | [STYLE_GUIDE.md](STYLE_GUIDE.md) | ✅ یکپارچه |
| مستندات | اگر رفتار تغییر کرده به‌روز شده | ✅ کامل |

#### 📝 قالب توضیحات PR

این را در توضیحات PR خود کپی کنید:

```markdown
## 📋 توضیحات
<!-- به طور خلاصه توضیح دهید این PR چه کاری انجام می‌دهد و چرا -->

## 🔗 Issue مرتبط
<!-- issue مرتبط که این PR رفع می‌کند -->
Closes #123

## 📝 نوع تغییر
- [ ] 🐛 رفع باگ (تغییر non-breaking که یک issue را رفع می‌کند)
- [ ] ✨ ویژگی جدید (تغییر non-breaking که functionality اضافه می‌کند)
- [ ] 💥 تغییر breaking (رفع باگ یا ویژگی که باعث تغییر عملکرد موجود می‌شود)
- [ ] 📝 به‌روزرسانی مستندات
- [ ] 🧪 به‌روزرسانی تست
- [ ] 🔧 بازسازی کد (بدون تغییر عملکرد)

## 🧪 تست‌ها
<!-- تست‌هایی که اجرا کردید و نحوه بازتولید -->
- [ ] `make test` رد شده
- [ ] `make vet` رد شده
- [ ] تست دستی:
  - تست شده روی: [سیستم‌عامل]
  - تست شده با: [سیستم‌عامل سرور ریموت]
  - مراحل: [توضیح]

## 📸 اسکرین‌شات / لاگ‌ها
<!-- اگر تغییر UI دارید، اسکرین‌شات اضافه کنید. اگر رفع باگ، قبل/بعد لاگ -->

## ✅ چک‌لیست
- [ ] کد من از راهنمای استایل پروژه پیروی می‌کند
- [ ] تست‌هایی اضافه کرده‌ام که اثبات می‌کند fix/feature کار می‌کند
- [ ] تمام تست‌های جدید و موجود رد شده‌اند
- [ ] مستندات را بر این اساس به‌روزرسانی کرده‌ام
- [ ] ورودی به CHANGELOG.md اضافه کرده‌ام (در صورت نیاز)
- [ ] تغییرات من هیچ هشدار جدیدی تولید نمی‌کند
- [ ] تغییرات وابسته ادغام و منتشر شده‌اند

## 💬 نکات اضافی
<!-- هر زمینه دیگری درباره PR -->
```

#### 💬 قرارداد پیام commit

ما از [Conventional Commits](https://www.conventionalcommits.org/) استفاده می‌کنیم — این به ما کمک می‌کند changelog را به صورت خودکار تولید کنیم:

| پیشوند | کی | مثال |
|--------|---|------|
| `feat:` | ویژگی جدید | `feat: add VLESS transport support` |
| `fix:` | رفع باگ | `fix: resolve failover race condition` |
| `docs:` | مستندات | `docs: add Chinese translation` |
| `test:` | تست‌ها | `test: add WireGuard unit tests` |
| `refactor:` | بازسازی کد | `refactor: extract SSH client module` |
| `perf:` | عملکرد | `perf: optimize ping latency calculation` |
| `style:` | قالب‌بندی | `style: fix indentation in dashboard.go` |
| `chore:` | نگهداری | `chore: update Go dependencies` |
| `ci:` | CI/CD | `ci: add GitHub Actions workflow` |

**قالب:** `<type>: <description>`

**مثال‌ها:**
```
feat: add VLESS transport with TLS support
fix: resolve connection timeout on slow networks
docs: update API reference for v0.3.0
test: add integration tests for failover engine
```

#### 👀 فرآیند بررسی

1. **چک‌های خودکار** باید pass شوند (CI/CD)
2. **بررسی کد** توسط حداقل **۱ نگه‌دارنده**
3. **بدون conflict** با `main`
4. **Squash and merge** — تاریخچه تمیز commit

**معیارهای بررسی:**
- ✅ کد تمیز و از الگوهای موجود پیروی می‌کند
- ✅ تست‌ها functionality جدید را پوشش می‌دهند
- ✅ مستندات به‌روزرسانی شده
- ✅ بدون تغییرات breaking (یا به وضوح مستند شده)
- ✅ عملکرد کاهش نیافته

---

### 🏆 تقدیر

تمام مشارکت‌کنندگان در [بخش مشارکت‌کنندگان](#-مشارکت‌کنندگان) شناخته می‌شوند. همچنین ما:

- مشارکت‌کنندگان را در release notes ذکر می‌کنیم
- تقدیر ویژه برای مشارکت‌های مهم داریم
- مشارکت‌کنندگان جدید با خوش‌آمدگویی در PR comments استقبال می‌کنیم

---

### 💬 نیاز به کمک دارید؟

- 📖 [راهنمای مشارکت](CONTRIBUTING.md) را برای دستورالعمل‌های دقیق بخوانید
- 💬 به [کانال تلگرام](https://t.me/NyxoraCore) برای سؤالات بپیوندید
- 🐛 [Issuesهای موجود](https://github.com/nyxora-network/nyxora/issues) را برای الهام بررسی کنید
- 📧 برای نگرانی‌های خصوصی به maintainers ایمیل بزنید

> **به یاد داشته باشید:** سؤال احمقانه‌ای وجود ندارد! همه ما یک زمانی مبتدی بودیم. درخواست کمک را دریغ نکنید.

---

## 🌟 مشارکت‌کنندگان

<a href="https://github.com/nyxora-network/nyxora/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=nyxora-network/nyxora" alt="مشارکت‌کنندگان" />
</a>

---

## 📄 مجوز

این پروژه تحت مجوز **MIT** منتشر شده است.

---

<div align="center">
  <br>
  <p>
    <a href="https://t.me/NyxoraCore">کانال تلگرام</a> •
    <a href="https://github.com/nyxora-network/nyxora/issues">گزارش باگ</a> •
    <a href="https://github.com/nyxora-network/nyxora/issues">درخواست ویژگی</a>
  </p>
  <p>
    <sub>ساخته شده با ❤️ با استفاده از Go &amp; Bubble Tea</sub>
  </p>
</div>
