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
  <h3>अनुकूली सुरंग ऑर्केस्ट्रेटर</h3>
  <p>
    <b>स्व-उपचार मल्टी-ट्रांसपोर्ट VPN/सुरंग प्रबंधक</b><br>
    एक सर्वर पर स्थापित करें। किसी भी रिमोट सर्वर से कनेक्ट करें।<br>
    कोई एजेंट नहीं चाहिए। स्वचालित प्रावधान। स्वचालित फेलओवर। इंटरैक्टिव TUI।
  </p>
  <br>
  <p>
    <a href="#-विशेषताएँ">विशेषताएँ</a> •
    <a href="#-त्वरित-आरंभ">त्वरित आरंभ</a> •
    <a href="#-एक-पंक्ति-स्थापना">स्थापना</a> •
    <a href="#-उपयोग">उपयोग</a> •
    <a href="#-आर्किटेक्चर">आर्किटेक्चर</a> •
    <a href="#-विकास">विकास</a>
  </p>
</div>

<br>

---

## ✨ विशेषताएँ

<table>
<tr>
<td width="50%">

**🧠 स्व-उपचार ऑर्केस्ट्रेशन**
- 12 सुरंग ट्रांसपोर्ट: WireGuard, OpenVPN, SSH, QUIC, FRP, Rathole, IPsec, Shadowsocks, Hysteria, Backhaul, TCP, WebSocket
- स्वचालित फेलओवर — खराब सुरंगों का पता लगाकर तुरंत स्विच करें
- 5 मल्टीपाथ शेड्यूलिंग मोड (भारित, न्यूनतम-विलंबता, न्यूनतम-हानि, समान, सभी-सक्रिय)
- रीयल-टाइम स्कोरिंग इंजन (विलंबता + पैकेट हानि + भार)

</td>
<td width="50%">

**🚀 शून्य-कॉन्फ़िग रिमोट**
- रिमोट सर्वर पर किसी एजेंट या सॉफ़्टवेयर की आवश्यकता नहीं
- केवल SSH एक्सेस (पासवर्ड या कुंजी)
- स्वचालित OS पहचान (Ubuntu, Debian, CentOS)
- रिमोट पर स्वचालित रूप से सुरंग बाइनरी स्थापित करें

</td>
</tr>
<tr>
<td width="50%">

**🖥️ समृद्ध टर्मिनल UI**
- कीबोर्ड नेविगेशन के साथ इंटरैक्टिव Bubble Tea TUI
- 3 पेशेवर रंग थीम (Catppuccin Mocha, Tokyo Night, Catppuccin Latte)
- रीयल-टाइम आँकड़ों के साथ लाइव डैशबोर्ड
- एनिमेटेड ग्रेडिएंट प्रगति पट्टियाँ
- ASCII कला लोगो के साथ बूट स्प्लैश
- सुरंग स्थलाकृति दृश्य
- चरण-दर-चरण कनेक्ट विज़ार्ड

</td>
<td width="50%">

**🔐 एंटरप्राइज़-ग्रेड सुरक्षा**
- कर्नेल स्तर पर WireGuard VPN
- IPsec/strongSwan समर्थन
- Shadowsocks एन्क्रिप्टेड प्रॉक्सी
- Hysteria 2 (संशोधित QUIC सेंसरशिप-विरोधी के साथ)
- स्वचालित गुप्त उत्पादन (पासवर्ड, PSK, टोकन)

</td>
</tr>
</table>

---

## 📦 एक-पंक्ति स्थापना

```bash
curl -L github.com/nyxorammd-lgtm/nyxora/releases/download/v0.2.0/nyxora_linux_amd64 -o /usr/local/bin/nyxora && chmod +x /usr/local/bin/nyxora
```

या `wget` के साथ:

```bash
wget -q https://github.com/nyxorammd-lgtm/nyxora/releases/download/v0.2.0/nyxora_linux_amd64 -O /usr/local/bin/nyxora && chmod +x /usr/local/bin/nyxora
```

<details>
<summary><b>📋 मैन्युअल स्थापना (स्रोत से)</b></summary>

```bash
# आवश्यकताएँ
sudo apt install golang-go git ssh sshpass wireguard curl
# या: brew install go (macOS)

# क्लोन करें
git clone https://github.com/nyxorammd-lgtm/nyxora.git
cd nyxora

# बिल्ड करें
make build

# स्थापित करें
sudo make install

# सत्यापित करें
nyxora version
```
</details>

---

## 🚀 त्वरित आरंभ

```bash
# 1. कॉन्फ़िग सेटअप और निर्भरताएँ जाँचें
nyxora install

# 2. रिमोट सर्वर से कनेक्ट करें
nyxora connect 192.168.1.100 --user root --password your_password

# 3. इंटरैक्टिव TUI लॉन्च करें
nyxora tui

# 4. लाइव मॉनिटरिंग डैशबोर्ड
nyxora dashboard
```

### कनेक्ट विकल्प

```bash
nyxora connect <host> [options]

विकल्प:
  --user, -u <name>       SSH उपयोगकर्ता (डिफ़ॉल्ट: root)
  --port, -p <port>       SSH पोर्ट (डिफ़ॉल्ट: 22)
  --password <pass>       SSH पासवर्ड
  --mode <mode>           सर्वर मोड: full, lite, minimal
  --transports <list>     अल्पविराम से अलग ट्रांसपोर्ट सूची
  --ports <pairs>         पोर्ट ओवरराइड: wg=51820,ss=8388,...
```

#### सर्वर मोड

| मोड | ट्रांसपोर्ट | RAM आवश्यकता |
|-----|-------------|---------------|
| `full` | सभी 11 सुरंगें | 2GB+ |
| `lite` | हल्का चयन | 512MB–2GB |
| `minimal` | केवल SSH + Shadowsocks | < 512MB |

---

## 🎮 इंटरैक्टिव TUI

NYXORA में [Bubble Tea](https://github.com/charmbracelet/bubbletea) और [Lip Gloss](https://github.com/charmbracelet/lipgloss) के साथ निर्मित एक पूर्ण-विशेषताओं वाला टर्मिनल UI है।

```
┌──────────────────────────────────────────────────────────┐
│  NYXORA v0.2.0                                          │
│  ────────────────────────────────────────────────────    │
│                                                          │
│  CPU: 0.5  ████░░░░░░░░░░░░░░░░                        │
│  RAM: 45%  ██████████░░░░░░░░░░                        │
│                                                          │
│  [1] C  सर्वर से कनेक्ट करें                            │
│  [2] D  डैशबोर्ड                                        │
│  [3] I  सर्वर जानकारी                                   │
│  [4] N  स्थापित करें                                    │
│  [5] U  अपडेट जाँचें                                    │
│  [6] X  डिस्कनेक्ट करें                                 │
│  [7] T  सुरंग स्थलाकृति                                 │
│  [8] H  सहायता                                          │
│  [9] Q  बाहर निकलें                                     │
│                                                          │
│  ┌────────────────────────────────────────────────────┐  │
│  │  रिमोट सर्वर से कनेक्ट करें                       │  │
│  └────────────────────────────────────────────────────┘  │
│  ↑↓ नेविगेट  ↵ चुनें  1/2/3 थीम  s स्थिति  ? सहायता   │
│  https://t.me/NyxoraCore                                 │
└──────────────────────────────────────────────────────────┘
```

### कीबोर्ड शॉर्टकट

| कुंजी | कार्य |
|--------|-------|
| `↑` / `↓` | मेनू नेविगेट करें |
| `Enter` | आइटम चुनें |
| `Esc` | वापस जाएँ |
| `q` | बाहर निकलें / मेनू पर वापस जाएँ |
| `1` | Catppuccin Mocha (गहरा) |
| `2` | Tokyo Night (गहरा) |
| `3` | Catppuccin Latte (हल्का) |
| `s` | स्थिति पट्टी टॉगल करें |
| `?` | सहायता स्क्रीन खोलें |
| `t` | सुरंग स्थलाकृति दृश्य |

---

## 🏗️ आर्किटेक्चर

```
┌─────────────────────────────────────────────────────────────────┐
│  nyxora (स्थानीय सर्वर)                                        │
│                                                                 │
│  ┌──────────────┐  ┌────────────────────────────────────────┐   │
│  │  ऑर्केस्ट्रेटर │  │  ट्रांसपोर्ट मैनेजर                    │   │
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
│  │  मल्टीपाथ   │  │  फेलओवर     │  │  डैशबोर्ड / TUI     │  │
│  │  शेड्यूलर   │  │  इंजन       │  │  (Bubble Tea)        │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
│         │                                                       │
│         │ SSH + प्रावधान                                        │
│         ▼                                                       │
│  ┌────────────────────────────────────────────────────────┐     │
│  │  रिमोट सर्वर (कोई एजेंट नहीं चाहिए)                    │     │
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

### कनेक्शन प्रवाह

```
nyxora connect 91.107.243.237 --user root --password ...

  1.  PING          → विलंबता और पैकेट हानि मापें
  2.  SSH           → रिमोट सर्वर पर प्रमाणित करें
  3.  DETECT OS     → Ubuntu / Debian / CentOS पहचान
  4.  INSTALL       → रिमोट पर सुरंग बाइनरी स्थापित करें
  5.  WG KEY        → स्थानीय रूप से WireGuard कुंजी जोड़ी उत्पन्न करें
  6.  REMOTE WG     → SSH: कॉन्फ़िग + wg-quick up + iptables
  7.  LOCAL WG      → रिमोट सार्वजनिक कुंजी के साथ wg-quick up nyxora0
  8.  PROVISION     → डेमॉन शुरू करें: frps, rathole, ss, hys, backhaul
  9.  ALL-ACTIVE    → सभी सुरंगों का एक साथ परीक्षण और सक्रियण
  10. MONITOR       → हर 10 सेकंड: ping, स्कोर, फेलओवर जाँच
```

---

## 📋 कमांड

| कमांड | विवरण |
|--------|--------|
| `nyxora install` | कॉन्फ़िग सेटअप और निर्भरताएँ जाँचें |
| `nyxora connect <host>` | रिमोट सर्वर से कनेक्ट करें |
| `nyxora disconnect` | सभी सुरंगें बंद करें |
| `nyxora status` | कनेक्शन स्थिति दिखाएँ |
| `nyxora dashboard` | लाइव टर्मिनल डैशबोर्ड |
| `nyxora tui` | इंटरैक्टिव Bubble Tea मेनू |
| `nyxora update` | अपडेट जाँचें |
| `nyxora server` | सर्वर जानकारी और सुझाया मोड दिखाएँ |
| `nyxora version` | संस्करण दिखाएँ |
| `nyxora daemon` | बैकग्राउंड सेवा के रूप में चलाएँ |
| `nyxora help` | सहायता दिखाएँ |

---

## 🔧 कॉन्फ़िगरेशन

### पर्यावरण चर

| चर | विवरण | डिफ़ॉल्ट |
|-----|--------|-----------|
| `NYXORA_SS_PASSWORD` | Shadowsocks पासवर्ड | स्वतः-उत्पन्न |
| `NYXORA_SS_METHOD` | Shadowsocks सिफर | `aes-256-gcm` |
| `NYXORA_RATHOLE_TOKEN` | Rathole प्रमाणीकरण टोकन | स्वतः-उत्पन्न |
| `NYXORA_HYSTERIA_AUTH` | Hysteria प्रमाणीकरण पासवर्ड | स्वतः-उत्पन्न |
| `NYXORA_BACKHAUL_TOKEN` | Backhaul प्रमाणीकरण टोकन | स्वतः-उत्पन्न |
| `NYXORA_IPSEC_PSK` | IPsec पूर्व-साझा कुंजी | स्वतः-उत्पन्न |
| `NYXORA_ALL_ACTIVE` | सभी सुरंगों को एक साथ सक्षम करें | `false` |

---

## 📦 ट्रांसपोर्ट

| # | नाम | पोर्ट | प्रोटोकॉल | श्रेणी | आधार स्कोर | भार |
|---|-----|-------|-----------|---------|-------------|------|
| 1 | **wireguard** | 51820 | UDP | VPN | 95 | 30 |
| 2 | **openvpn** | 1194 | UDP | VPN | 75 | 10 |
| 3 | **ssh** | 22 | TCP | सुरंग | 60 | 5 |
| 4 | **quic** | 9923 | UDP | सुरंग | 80 | 15 |
| 5 | **frp** | 7000 | TCP | रिले | 70 | 10 |
| 6 | **rathole** | 2333 | TCP | रिले | 85 | 12 |
| 7 | **ipsec** | 500 | UDP | VPN | 70 | 5 |
| 8 | **shadowsocks** | 8388 | TCP | प्रॉक्सी | 55 | 3 |
| 9 | **hysteria** | 8444 | UDP | सुरंग | 90 | 12 |
| 10| **backhaul** | 3080 | TCP | रिले | 82 | 10 |
| 11| **tcp** | 9924 | TCP | सुरंग | 50 | 3 |

---

## 🧑‍💻 विकास

### आवश्यकताएँ

- Go 1.24+
- Linux या macOS
- `ssh`, `sshpass`, `wg`, `curl`, `ping`

### सेटअप

```bash
git clone https://github.com/nyxorammd-lgtm/nyxora.git
cd nyxora

# बिल्ड
make build

# टेस्ट
make test

# वेट
make vet

# स्थानीय रूप से चलाएँ
./nyxora version
```

---

## 🤝 योगदान

हम योगदान का स्वागत करते हैं! कृपया हमारा [योगदान गाइड](CONTRIBUTING.md) देखें।

**योगदान के तरीके:**
- [GitHub Issues](https://github.com/nyxorammd-lgtm/nyxora/issues) पर बग रिपोर्ट करें
- नए ट्रांसपोर्ट प्रकार सुझाएँ
- TUI / डैशबोर्ड में सुधार करें
- अधिक OS लक्ष्यों के लिए समर्थन जोड़ें
- परीक्षण और दस्तावेज़ लिखें
- खुले Issues के लिए PR सबमिट करें

---

## 📄 लाइसेंस

यह प्रोजेक्ट **MIT लाइसेंस** के तहत लाइसेंस प्राप्त है।

---

<div align="center">
  <br>
  <p>
    <a href="https://t.me/NyxoraCore">टेलीग्राम चैनल</a> •
    <a href="https://github.com/nyxorammd-lgtm/nyxora/issues">बग रिपोर्ट करें</a> •
    <a href="https://github.com/nyxorammd-lgtm/nyxora/issues">सुविधा अनुरोध</a>
  </p>
  <p>
    <sub>❤️ के साथ Go और Bubble Tea का उपयोग करके बनाया गया</sub>
  </p>
</div>
