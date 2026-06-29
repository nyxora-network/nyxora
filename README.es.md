<div align="center">
  <img src="https://img.shields.io/badge/go-%3E%3D1.25-00ADD8?style=flat&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/license-MIT-blue?style=flat" alt="Licencia">
  <img src="https://img.shields.io/badge/status-active-success?style=flat" alt="Estado">
  <img src="https://img.shields.io/badge/PRs-welcome-brightgreen?style=flat" alt="PRs Bienvenidos">
  <br>
  <img src="https://img.shields.io/badge/transports-11-ff69b4?style=flat" alt="11 Transportes">
  <img src="https://img.shields.io/badge/platform-linux%20%7C%20macOS-lightgrey?style=flat" alt="Plataforma">
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
  <h3>Orquestador Adaptativo de Túneles</h3>
  <p>
    <b>Gestor de VPN/túneles multitransporte con autocuración</b><br>
    Instala en <i>un</i> servidor. Conéctate a <i>cualquier</i> servidor remoto.<br>
    Sin agente remoto. Aprovisionamiento automático. Failover automático. TUI interactiva.
  </p>
  <br>
  <p>
    <a href="#-características">Características</a> •
    <a href="#-inicio-rápido">Inicio Rápido</a> •
    <a href="#-instalación-en-una-línea">Instalación</a> •
    <a href="#-uso">Uso</a> •
    <a href="#-arquitectura">Arquitectura</a> •
    <a href="#-desarrollo">Desarrollo</a>
  </p>
</div>

<br>

---

## ✨ Características

<table>
<tr>
<td width="50%">

**🧠 Orquestación Autocurativa**
- 11 transportes de túnel: WireGuard, OpenVPN, SSH, QUIC, FRP, Rathole, IPsec, Shadowsocks, Hysteria, Backhaul, TCP
- Failover automático — detecta túneles degradados, cambia al instante
- 5 modos de programación multirruta (ponderado, menor latencia, menor pérdida, equitativo, todo-activo)
- Motor de puntuación en tiempo real (latencia + pérdida de paquetes + peso)

</td>
<td width="50%">

**🚀 Remoto Sin Configuración**
- No se necesita agente ni software en el servidor remoto
- Solo acceso SSH (contraseña o clave)
- Detecta automáticamente el SO (Ubuntu, Debian, CentOS)
- Instala automáticamente los binarios de túnel en el remoto

</td>
</tr>
<tr>
<td width="50%">

**🖥️ Interfaz de Terminal Enriquecida**
- TUI interactivo con Bubble Tea y navegación por teclado
- 3 temas de color profesionales (Catppuccin Mocha, Tokyo Night, Catppuccin Latte)
- Panel en vivo con estadísticas en tiempo real
- Barras de progreso animadas con degradado
- Pantalla de inicio con logotipo ASCII
- Vista de topología de túneles
- Asistente de conexión paso a paso

</td>
<td width="50%">

**🔐 Seguridad de Nivel Empresarial**
- VPN WireGuard a nivel de kernel
- Soporte IPsec/strongSwan
- Proxy cifrado Shadowsocks
- Hysteria 2 (QUIC modificado con anticensura)
- Generación automática de secretos (contraseñas, PSK, tokens)

</td>
</tr>
</table>

---

## 📦 Instalación en Una Línea

```bash
curl -fsSL https://raw.githubusercontent.com/nyxorammd-lgtm/nyxora/main/install.sh | sudo bash
```

O con `wget`:

```bash
wget -qO- https://raw.githubusercontent.com/nyxorammd-lgtm/nyxora/main/install.sh | sudo bash
```

<details>
<summary><b>📋 Instalación Manual (desde código fuente)</b></summary>

```bash
# Prerrequisitos
sudo apt install golang-go git ssh sshpass wireguard curl
# o: brew install go (macOS)

# Clonar
git clone https://github.com/nyxorammd-lgtm/nyxora.git
cd nyxora

# Compilar
make build

# Instalar
sudo make install

# Verificar
nyxora version
```
</details>

---

## 🚀 Inicio Rápido

```bash
# 1. Configurar y verificar dependencias
nyxora install

# 2. Conectar a un servidor remoto
nyxora connect 192.168.1.100 --user root --password tu_contraseña

# 3. Iniciar TUI interactiva
nyxora tui

# 4. Panel de monitoreo en vivo
nyxora dashboard
```

### Opciones de Conexión

```bash
nyxora connect <host> [options]

Opciones:
  --user, -u <name>       Usuario SSH (por defecto: root)
  --port, -p <port>       Puerto SSH (por defecto: 22)
  --password <pass>       Contraseña SSH
  --mode <mode>           Modo servidor: full, lite, minimal
  --transports <list>     Lista de transportes separada por comas
  --ports <pairs>         Anulación de puertos: wg=51820,ss=8388,...
```

#### Modos de Servidor

| Modo | Transportes | RAM Requerida |
|------|-------------|---------------|
| `full` | Todos los 11 túneles | 2GB+ |
| `lite` | Selección ligera | 512MB–2GB |
| `minimal` | Solo SSH + Shadowsocks | < 512MB |

---

## 🎮 TUI Interactiva

NYXORA cuenta con una interfaz de terminal completa construida con [Bubble Tea](https://github.com/charmbracelet/bubbletea) y [Lip Gloss](https://github.com/charmbracelet/lipgloss).

```
┌──────────────────────────────────────────────────────────┐
│  NYXORA v0.2.0                                          │
│  ────────────────────────────────────────────────────    │
│                                                          │
│  CPU: 0.5  ████░░░░░░░░░░░░░░░░                        │
│  RAM: 45%  ██████████░░░░░░░░░░                        │
│                                                          │
│  [1] C  Conectar al Servidor                             │
│  [2] D  Panel                                            │
│  [3] I  Información del Servidor                         │
│  [4] N  Instalar                                         │
│  [5] U  Buscar Actualizaciones                           │
│  [6] X  Desconectar                                      │
│  [7] T  Topología de Túneles                             │
│  [8] H  Ayuda                                            │
│  [9] Q  Salir                                            │
│                                                          │
│  ┌────────────────────────────────────────────────────┐  │
│  │  Conectar a un servidor remoto                     │  │
│  └────────────────────────────────────────────────────┘  │
│  ↑↓ navegar  ↵ seleccionar  1/2/3 tema  s estado  ? ayuda│
│  https://t.me/NyxoraCore                                 │
└──────────────────────────────────────────────────────────┘
```

### Atajos de Teclado

| Tecla | Acción |
|-------|--------|
| `↑` / `↓` | Navegar menú |
| `Enter` | Seleccionar |
| `Esc` | Volver |
| `q` | Salir / Volver al menú |
| `1` | Catppuccin Mocha (oscuro) |
| `2` | Tokyo Night (oscuro) |
| `3` | Catppuccin Latte (claro) |
| `s` | Alternar barra de estado |
| `?` | Abrir pantalla de ayuda |
| `t` | Vista de topología de túneles |

---

## 🏗️ Arquitectura

```
┌─────────────────────────────────────────────────────────────────┐
│  nyxora (servidor local)                                        │
│                                                                 │
│  ┌──────────────┐  ┌────────────────────────────────────────┐   │
│  │  Orquestador  │  │  Gestor de Transportes                │   │
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
│  │  Programador │  │  Motor de    │  │  Panel / TUI         │  │
│  │  Multirruta  │  │  Failover    │  │  (Bubble Tea)        │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
│         │                                                       │
│         │ SSH + aprovisionamiento                                │
│         ▼                                                       │
│  ┌────────────────────────────────────────────────────────┐     │
│  │  Servidor remoto (sin agente)                          │     │
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

### Flujo de Conexión

```
nyxora connect 91.107.243.237 --user root --password ...

  1.  PING          → Medir latencia y pérdida de paquetes
  2.  SSH           → Autenticar en el servidor remoto
  3.  DETECT OS     → Detectar Ubuntu / Debian / CentOS
  4.  INSTALL       → Instalar binarios de túnel en remoto
  5.  WG KEY        → Generar par de claves WireGuard localmente
  6.  REMOTE WG     → SSH: configurar + wg-quick up + iptables
  7.  LOCAL WG      → wg-quick up nyxora0 con clave pública remota
  8.  PROVISION     → Iniciar daemons: frps, rathole, ss, hys, backhaul
  9.  ALL-ACTIVE    → Probar y activar todos los túneles simultáneamente
  10. MONITOR       → Cada 10s: ping, puntuación, verificación failover
```

---

## 📋 Comandos

| Comando | Descripción |
|---------|-------------|
| `nyxora install` | Configurar y verificar dependencias |
| `nyxora connect <host>` | Conectar a servidor remoto |
| `nyxora disconnect` | Cerrar todos los túneles |
| `nyxora status` | Mostrar estado de conexión |
| `nyxora dashboard` | Panel de terminal en vivo |
| `nyxora tui` | Menú interactivo Bubble Tea |
| `nyxora update` | Buscar actualizaciones |
| `nyxora server` | Mostrar información del servidor |
| `nyxora version` | Mostrar versión |
| `nyxora daemon` | Ejecutar como servicio de fondo |
| `nyxora help` | Mostrar ayuda |

---

## 🔧 Configuración

### Variables de Entorno

| Variable | Descripción | Por Defecto |
|----------|-------------|-------------|
| `NYXORA_SS_PASSWORD` | Contraseña Shadowsocks | auto-generada |
| `NYXORA_SS_METHOD` | Cifrado Shadowsocks | `aes-256-gcm` |
| `NYXORA_RATHOLE_TOKEN` | Token de autenticación Rathole | auto-generado |
| `NYXORA_HYSTERIA_AUTH` | Contraseña Hysteria | auto-generada |
| `NYXORA_BACKHAUL_TOKEN` | Token Backhaul | auto-generado |
| `NYXORA_IPSEC_PSK` | Clave precompartida IPsec | auto-generada |
| `NYXORA_ALL_ACTIVE` | Activar todos los túneles simultáneamente | `false` |

---

## 📦 Transportes

| # | Nombre | Puerto | Protocolo | Categoría | Puntuación Base | Peso |
|---|--------|--------|-----------|-----------|-----------------|------|
| 1 | **wireguard** | 51820 | UDP | VPN | 95 | 30 |
| 2 | **openvpn** | 1194 | UDP | VPN | 75 | 10 |
| 3 | **ssh** | 22 | TCP | Túnel | 60 | 5 |
| 4 | **quic** | 9923 | UDP | Túnel | 80 | 15 |
| 5 | **frp** | 7000 | TCP | Relay | 70 | 10 |
| 6 | **rathole** | 2333 | TCP | Relay | 85 | 12 |
| 7 | **ipsec** | 500 | UDP | VPN | 70 | 5 |
| 8 | **shadowsocks** | 8388 | TCP | Proxy | 55 | 3 |
| 9 | **hysteria** | 8444 | UDP | Túnel | 90 | 12 |
| 10| **backhaul** | 3080 | TCP | Relay | 82 | 10 |
| 11| **tcp** | 9924 | TCP | Túnel | 50 | 3 |

---

## 🧑‍💻 Desarrollo

### Prerrequisitos

- Go 1.25+
- Linux o macOS
- `ssh`, `sshpass`, `wg`, `curl`, `ping`

### Configuración

```bash
git clone https://github.com/nyxorammd-lgtm/nyxora.git
cd nyxora

# Compilar
make build

# Pruebas
make test

# Vet
make vet

# Ejecutar localmente
./nyxora version
```

---

## 🤝 Contribuir

¡Aceptamos contribuciones! Consulta nuestra [Guía de Contribución](CONTRIBUTING.md).

**Formas de contribuir:**
- Reporta errores en [GitHub Issues](https://github.com/nyxorammd-lgtm/nyxora/issues)
- Sugiere nuevos tipos de transporte
- Mejora la TUI / el panel
- Agrega soporte para más sistemas operativos
- Escribe pruebas y documentación
- Envía PRs para issues abiertos

---

## 📄 Licencia

Este proyecto está licenciado bajo **MIT License**.

---

<div align="center">
  <br>
  <p>
    <a href="https://t.me/NyxoraCore">Canal de Telegram</a> •
    <a href="https://github.com/nyxorammd-lgtm/nyxora/issues">Reportar Error</a> •
    <a href="https://github.com/nyxorammd-lgtm/nyxora/issues">Solicitar Característica</a>
  </p>
  <p>
    <sub>Hecho con ❤️ usando Go &amp; Bubble Tea</sub>
  </p>
</div>
