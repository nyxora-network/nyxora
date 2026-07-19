<div align="center">
  <img src="https://img.shields.io/badge/go-%3E%3D1.24-00ADD8?style=flat&logo=go" alt="إصدار Go">
  <img src="https://img.shields.io/badge/license-MIT-blue?style=flat" alt="الترخيص">
  <img src="https://img.shields.io/badge/status-active-success?style=flat" alt="الحالة">
  <img src="https://img.shields.io/badge/PRs-welcome-brightgreen?style=flat" alt="PRs مرحب بها">
  <br>
  <img src="https://img.shields.io/badge/transports-12-ff69b4?style=flat" alt="12 ناقلاً">
  <img src="https://img.shields.io/badge/platform-linux%20%7C%20macOS-lightgrey?style=flat" alt="المنصة">
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

<div align="center" dir="rtl">
<h1>NYXORA</h1>
  <h3>منسق الأنفاق التكيفي</h3>
  <p>
    <b>مدير VPN/نفق متعدد الناقلات ذاتي الإصلاح</b><br>
    ثبته على <i>خادم</i> واحد. اتصل بأي <i>خادم</i> بعيد.<br>
    لا حاجة لعامل عن بعد. توفير تلقائي. تجاوز الفشل تلقائي. واجهة مستخدم تفاعلية.
  </p>
  <br>
  <p>
    <a href="#-الميزات">الميزات</a> •
    <a href="#-بداية-سريعة">بداية سريعة</a> •
    <a href="#-تثبيت-بسطر-واحد">تثبيت</a> •
    <a href="#-الاستخدام">الاستخدام</a> •
    <a href="#-الهندسة">الهندسة</a> •
    <a href="#-التطوير">التطوير</a>
  </p>
</div>

<br>

---

<div dir="rtl">

## ✨ الميزات

<table>
<tr>
<td width="50%">

**🧠 تنسيق ذاتي الإصلاح**
- 12 ناقلاً للأنفاق: WireGuard, OpenVPN, SSH, QUIC, FRP, Rathole, IPsec, Shadowsocks, Hysteria, Backhaul, TCP, WebSocket
- تجاوز الفشل تلقائي — يكتشف الأنفاق المتدهورة ويبدل فوراً
- 5 أوضاع جدولة متعددة المسارات (موزون، أقل زمن وصول، أقل فقدان، متساوٍ، الكل نشط)
- محرك تسجيل نقاط في الوقت الفعلي (زمن الوصول + فقدان الحزم + الوزن)

</td>
<td width="50%">

**🚀 عن بعد بدون إعدادات**
- لا حاجة لعامل أو برنامج على الخادم البعيد
- فقط الوصول عبر SSH (كلمة مرور أو مفتاح)
- كشف نظام التشغيل تلقائياً (Ubuntu, Debian, CentOS)
- تثبيت برامج الأنفاق على البعيد تلقائياً

</td>
</tr>
<tr>
<td width="50%">

**🖥️ واجهة طرفية غنية**
- واجهة TUI تفاعلية مع Bubble Tea وتنقل بلوحة المفاتيح
- 3 سمات ألوان احترافية (Catppuccin Mocha, Tokyo Night, Catppuccin Latte)
- لوحة تحكم حية بإحصائيات الوقت الفعلي
- أشرطة تقدم متحركة متدرجة
- شاشة بدء بشعار ASCII
- عرض طوبولوجيا الأنفاق
- معالج اتصال خطوة بخطوة

</td>
<td width="50%">

**🔐 أمان على مستوى المؤسسات**
- WireGuard VPN على مستوى النواة
- دعم IPsec/strongSwan
- وكيل Shadowsocks مشفر
- Hysteria 2 (QUIC معدل بمكافحة الرقابة)
- غلاف TLS مع توليد شهادات ذاتية التوقيع
- توليد تلقائي للأسرار (كلمات المرور، PSK، الرموز)

</td>
<td width="50%">

**📊 مراقبة على مستوى الإنتاج**
- نقطة نهاية مقاييس Prometheus (/metrics)
- محلل DNS داخلي مع التخزين المؤقت
- تقييد المعدل لكل ناقل (Token Bucket)
- إعادة تحميل التكوين (File Watcher)

</td>
</tr>
</table>

</td>
</tr>
</table>

---

## 📦 تثبيت بسطر واحد

```bash
curl -L github.com/nyxora-network/nyxora/releases/download/v0.2.0/nyxora_linux_amd64 -o /usr/local/bin/nyxora && chmod +x /usr/local/bin/nyxora
```

أو باستخدام `wget`:

```bash
wget -q https://github.com/nyxora-network/nyxora/releases/download/v0.2.0/nyxora_linux_amd64 -O /usr/local/bin/nyxora && chmod +x /usr/local/bin/nyxora
```

<details>
<summary><b>📋 تثبيت يدوي (من المصدر)</b></summary>

```bash
# المتطلبات
sudo apt install golang-go git ssh sshpass wireguard curl
# أو: brew install go (macOS)

# استنساخ
git clone https://github.com/nyxora-network/nyxora.git
cd nyxora

# بناء
make build

# تثبيت
sudo make install

# تحقق
nyxora version
```
</details>

---

## 🚀 بداية سريعة

```bash
# 1. إعداد التكوين والتحقق من التبعيات
nyxora install

# 2. الاتصال بخادم بعيد
nyxora connect 192.168.1.100 --user root --password كلمة_السر

# 3. تشغيل TUI التفاعلية
nyxora tui

# 4. لوحة تحكم المراقبة الحية
nyxora dashboard
```

### خيارات الاتصال

```bash
nyxora connect <host> [options]

الخيارات:
  --user, -u <name>       اسم مستخدم SSH (افتراضي: root)
  --port, -p <port>       منفذ SSH (افتراضي: 22)
  --password <pass>       كلمة مرور SSH
  --mode <mode>           وضع الخادم: full, lite, minimal
  --transports <list>     قائمة الناقلات مفصولة بفواصل
  --ports <pairs>         تجاوز المنافذ: wg=51820,ss=8388,...
```

#### أوضاع الخادم

| الوضع | الناقلات | RAM المطلوبة |
|-------|----------|--------------|
| `full` | جميع الأنفاق الـ11 | 2GB+ |
| `lite` | اختيار خفيف | 512MB–2GB |
| `minimal` | SSH + Shadowsocks فقط | < 512MB |

---

## 🎮 TUI تفاعلية

تمتلك NYXORA واجهة طرفية كاملة المواصفات مبنية باستخدام [Bubble Tea](https://github.com/charmbracelet/bubbletea) و [Lip Gloss](https://github.com/charmbracelet/lipgloss).

```
┌──────────────────────────────────────────────────────────┐
│  NYXORA v0.2.0                                          │
│  ────────────────────────────────────────────────────    │
│                                                          │
│  CPU: 0.5  ████░░░░░░░░░░░░░░░░                        │
│  RAM: 45%  ██████████░░░░░░░░░░                        │
│                                                          │
│  [1] C  الاتصال بالخادم                                 │
│  [2] D  لوحة التحكم                                     │
│  [3] I  معلومات الخادم                                  │
│  [4] N  تثبيت                                           │
│  [5] U  التحقق من التحديثات                             │
│  [6] X  قطع الاتصال                                     │
│  [7] T  طوبولوجيا النفق                                 │
│  [8] H  المساعدة                                        │
│  [9] Q  خروج                                            │
│                                                          │
│  ┌────────────────────────────────────────────────────┐  │
│  │  الاتصال بخادم بعيد                                │  │
│  └────────────────────────────────────────────────────┘  │
│  ↑↓ تنقل  ↵ اختر  1/2/3 سمة  s حالة  ? مساعدة  q خروج   │
│  https://t.me/NyxoraCore                                 │
└──────────────────────────────────────────────────────────┘
```

### اختصارات لوحة المفاتيح

| المفتاح | الإجراء |
|---------|---------|
| `↑` / `↓` | التنقل في القائمة |
| `Enter` | اختيار عنصر |
| `Esc` | العودة |
| `q` | خروج / العودة للقائمة |
| `1` | Catppuccin Mocha (داكن) |
| `2` | Tokyo Night (داكن) |
| `3` | Catppuccin Latte (فاتح) |
| `s` | تبديل شريط الحالة |
| `?` | فتح شاشة المساعدة |
| `t` | عرض طوبولوجيا النفق |

---

## 🏗️ الهندسة

```
┌─────────────────────────────────────────────────────────────────┐
│  nyxora (الخادم المحلي)                                         │
│                                                                 │
│  ┌──────────────┐  ┌────────────────────────────────────────┐   │
│  │  المنسق      │  │  مدير الناقلات                         │   │
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
│  │  المجدول     │  │  محرك تجاوز  │  │  لوحة التحكم / TUI  │  │
│  │  متعدد المسارات│  │  الفشل      │  │  (Bubble Tea)        │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
│         │                                                       │
│         │ SSH + توفير                                           │
│         ▼                                                       │
│  ┌────────────────────────────────────────────────────────┐     │
│  │  الخادم البعيد (لا حاجة لعامل)                         │     │
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

### تدفق الاتصال

```
nyxora connect 91.107.243.237 --user root --password ...

  1.  PING          → قياس زمن الوصول وفقدان الحزم
  2.  SSH           → المصادقة على الخادم البعيد
  3.  DETECT OS     → كشف النظام Ubuntu / Debian / CentOS
  4.  INSTALL       → تثبيت برامج الأنفاق على البعيد
  5.  WG KEY        → توليد زوج مفاتيح WireGuard محلياً
  6.  REMOTE WG     → SSH: إعداد + wg-quick up + iptables
  7.  LOCAL WG      → wg-quick up nyxora0 مع المفتاح العام البعيد
  8.  PROVISION     → تشغيل البرامج الخلفية: frps, rathole, ss, hys, backhaul
  9.  ALL-ACTIVE    → اختبار وتفعيل جميع الأنفاق في وقت واحد
  10. MONITOR       → كل 10 ثوان: ping، نقاط، فحص تجاوز الفشل
```

---

## 📸 لقطات شاشة لوحة التحكم

يتضمن NYXORA **لوحة تحكم ويب** لإدارة الأنفاق. الوصول عبر: `http://YOUR_SERVER:8080/tunnel.html`

### خريطة الأنفاق
![خريطة الأنفاق](screenshots/01-tunnel-map.png)
*خريطة بصرية للاتصالات بين الخوادم مع خطوط أنفاق متحركة*

### اتصال جديد
![اتصال](screenshots/02-connect.png)
*نموذج إنشاء اتصال الأنفاق مع مصادقة SSH*

### الأنفاق النشطة
![الأنفاق](screenshots/03-tunnels.png)
*قائمة اتصالات الأنفاق النشطة مع الحالة والتحكم*

### الكشف التلقائي عن أفضل ناقل
![الكشف التلقائي](screenshots/04-auto-detect.png)
*فحص بينغ تلقائي لجميع الناقلات الـ 14*

### منشئ إعدادات Xray
![إعدادات Xray](screenshots/05-xray-config.png)
*إنشاء إعدادات Xray لـ VLESS و VMess و Trojan و Shadowsocks*

### تثبيت اللوحات
![اللوحات](screenshots/06-panels.png)
*تثبيت بنقرة واحدة لـ Marzban و 3X-UI و Xray Core*

### رسوم بيانية للأداء
![الرسوم](screenshots/07-charts.png)
*رسوم بيانية مباشرة للتأخير وعرض النطاق وفقدان الحزم*

### الأمان والتدقيق
![الأمان](screenshots/08-security.png)
*حالة SSL وبصمات الخوادم ومحاولات الدخول الفاشلة*

### المحطة الطرفية
![المحطة](screenshots/09-terminal.png)
*واجهة محطة طرفية كاملة مع أوامر NYXORA*

**[التوثيق الكامل للوحة التحكم ←](web/README.md)**

---

## 📋 الأوامر

| الأمر | الوصف |
|-------|-------|
| `nyxora install` | إعداد التكوين والتحقق من التبعيات |
| `nyxora connect <host>` | الاتصال بخادم بعيد |
| `nyxora disconnect` | إغلاق جميع الأنفاق |
| `nyxora status` | عرض حالة الاتصال |
| `nyxora dashboard` | لوحة تحكم طرفية حية |
| `nyxora tui` | قائمة Bubble Tea التفاعلية |
| `nyxora update` | التحقق من التحديثات |
| `nyxora server` | عرض معلومات الخادم والوضع المقترح |
| `nyxora version` | عرض الإصدار |
| `nyxora daemon` | التشغيل كخدمة خلفية |
| `nyxora help` | عرض المساعدة |

---

## 🔧 الإعدادات

### متغيرات البيئة

| المتغير | الوصف | الافتراضي |
|---------|-------|-----------|
| `NYXORA_SS_PASSWORD` | كلمة مرور Shadowsocks | تلقائي |
| `NYXORA_SS_METHOD` | تشفير Shadowsocks | `aes-256-gcm` |
| `NYXORA_RATHOLE_TOKEN` | رمز Rathole | تلقائي |
| `NYXORA_HYSTERIA_AUTH` | كلمة مرور Hysteria | تلقائي |
| `NYXORA_BACKHAUL_TOKEN` | رمز Backhaul | تلقائي |
| `NYXORA_IPSEC_PSK` | مفتاح IPsec المشترك مسبقاً | تلقائي |
| `NYXORA_ALL_ACTIVE` | تفعيل جميع الأنفاق في وقت واحد | `false` |

---

## 📦 الناقلات

| # | الاسم | المنفذ | البروتوكول | الفئة | النقاط الأساسية | الوزن |
|---|-------|--------|-----------|-------|-----------------|-------|
| 1 | **wireguard** | 51820 | UDP | VPN | 95 | 30 |
| 2 | **openvpn** | 1194 | UDP | VPN | 75 | 10 |
| 3 | **ssh** | 22 | TCP | نفق | 60 | 5 |
| 4 | **quic** | 9923 | UDP | نفق | 80 | 15 |
| 5 | **frp** | 7000 | TCP | مرحل | 70 | 10 |
| 6 | **rathole** | 2333 | TCP | مرحل | 85 | 12 |
| 7 | **ipsec** | 500 | UDP | VPN | 70 | 5 |
| 8 | **shadowsocks** | 8388 | TCP | وكيل | 55 | 3 |
| 9 | **hysteria** | 8444 | UDP | نفق | 90 | 12 |
| 10| **backhaul** | 3080 | TCP | مرحل | 82 | 10 |
| 11| **tcp** | 9924 | TCP | نفق | 50 | 3 |
| 12| **websocket** | 9925 | TCP | نفق | 70 | 8 |

---

## 🧑‍💻 التطوير

### المتطلبات

- Go 1.24+
- Linux أو macOS
- `ssh`, `sshpass`, `wg`, `curl`, `ping`

### الإعداد

```bash
git clone https://github.com/nyxora-network/nyxora.git
cd nyxora

# بناء
make build

# اختبارات
make test

# فحص
make vet

# تشغيل محلي
./nyxora version
```

---

## 🤝 المساهمة

نرحب بالمساهمات! يرجى مراجعة [دليل المساهمة](CONTRIBUTING.md).

**طرق المساهمة:**
- الإبلاغ عن الأخطاء عبر [GitHub Issues](https://github.com/nyxora-network/nyxora/issues)
- اقتراح أنواع ناقلات جديدة
- تحسين TUI / لوحة التحكم
- إضافة دعم لأنظمة تشغيل إضافية
- كتابة الاختبارات والتوثيق
- إرسال PRs للمشكلات المفتوحة

---

## 🌟 المساهمون

<a href="https://github.com/nyxora-network/nyxora/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=nyxora-network/nyxora" alt="المساهمون" />
</a>

---

## 📄 الترخيص

هذا المشروع مرخص تحت **رخصة MIT**.

---

</div>

<div align="center">
  <br>
  <p>
    <a href="https://t.me/NyxoraCore">قناة تلغرام</a> •
    <a href="https://github.com/nyxora-network/nyxora/issues">الإبلاغ عن خطأ</a> •
    <a href="https://github.com/nyxora-network/nyxora/issues">طلب ميزة</a>
  </p>
  <p>
    <sub>بُني بـ ❤️ باستخدام Go &amp; Bubble Tea</sub>
  </p>
</div>
