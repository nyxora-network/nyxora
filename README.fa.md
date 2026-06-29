<div align="center">
  <img src="https://img.shields.io/badge/go-%3E%3D1.25-00ADD8?style=flat&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/license-MIT-blue?style=flat" alt="مجوز">
  <img src="https://img.shields.io/badge/status-active-success?style=flat" alt="وضعیت">
  <img src="https://img.shields.io/badge/PRs-welcome-brightgreen?style=flat" alt="PRs Welcome">
  <br>
  <img src="https://img.shields.io/badge/transports-11-ff69b4?style=flat" alt="۱۱ ترنسپورت">
  <img src="https://img.shields.io/badge/platform-linux%20%7C%20macOS-lightgrey?style=flat" alt="پلتفرم">
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
- ۱۱ ترنسپورت تونل: WireGuard, OpenVPN, SSH, QUIC, FRP, Rathole, IPsec, Shadowsocks, Hysteria, Backhaul, TCP
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
- تولید خودکار رمزهای عبور (Passwords, PSKs, Tokens)

</td>
</tr>
</table>

---

## 📦 نصب تک خطی

```bash
curl -fsSL https://raw.githubusercontent.com/nyxorammd-lgtm/nyxora/main/install.sh | sudo bash
```

یا با `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/nyxorammd-lgtm/nyxora/main/install.sh | sudo bash
```

<details>
<summary><b>📋 نصب دستی (از سورس)</b></summary>

```bash
# پیش‌نیازها
sudo apt install golang-go git ssh sshpass wireguard curl
# یا: brew install go (macOS)

# Clone
git clone https://github.com/nyxorammd-lgtm/nyxora.git
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

- Go 1.25+
- لینوکس یا macOS
- `ssh`, `sshpass`, `wg`, `curl`, `ping`

### راه‌اندازی

```bash
git clone https://github.com/nyxorammd-lgtm/nyxora.git
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

مشارکت شما خوش‌آمد است! لطفاً [راهنمای مشارکت](CONTRIBUTING.md) را مطالعه کنید.

**راه‌های مشارکت:**
- گزارش باگ از طریق [GitHub Issues](https://github.com/nyxorammd-lgtm/nyxora/issues)
- پیشنهاد ترنسپورت جدید
- بهبود TUI / داشبورد
- افزودن پشتیبانی از سیستم‌عامل‌های بیشتر
- نوشتن تست و مستندات
- ارسال PR برای Issues باز

---

## 📄 مجوز

این پروژه تحت مجوز **MIT** منتشر شده است.

---

<div align="center">
  <br>
  <p>
    <a href="https://t.me/NyxoraCore">کانال تلگرام</a> •
    <a href="https://github.com/nyxorammd-lgtm/nyxora/issues">گزارش باگ</a> •
    <a href="https://github.com/nyxorammd-lgtm/nyxora/issues">درخواست ویژگی</a>
  </p>
  <p>
    <sub>ساخته شده با ❤️ با استفاده از Go &amp; Bubble Tea</sub>
  </p>
</div>
