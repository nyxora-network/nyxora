package transport

import "fmt"

type TunnelMeta struct {
	Name        string
	Type        string
	Category    string
	Description string
	Port        int
	Protocol    string
	Binary      string
	Deps        []string
	Score       float64
	Weight      int
}

var TunnelRegistry = []TunnelMeta{
	{
		Name:        "wireguard",
		Type:        "wireguard",
		Category:    CatVPN,
		Description: "Fastest kernel-level VPN, minimal overhead",
		Port:        51820,
		Protocol:    "udp",
		Binary:      "wg",
		Deps:        []string{"wg", "wg-quick"},
		Score:       95,
		Weight:      30,
	},
	{
		Name:        "openvpn",
		Type:        "openvpn",
		Category:    CatVPN,
		Description: "Mature SSL/TLS VPN, maximum compatibility",
		Port:        1194,
		Protocol:    "udp",
		Binary:      "openvpn",
		Deps:        []string{"openvpn"},
		Score:       75,
		Weight:      10,
	},
	{
		Name:        "ssh",
		Type:        "ssh",
		Category:    CatTunnel,
		Description: "Built-in on every server, zero install",
		Port:        22,
		Protocol:    "tcp",
		Binary:      "ssh",
		Deps:        []string{"ssh"},
		Score:       60,
		Weight:      5,
	},
	{
		Name:        "quic",
		Type:        "quic",
		Category:    CatTunnel,
		Description: "UDP-based, handles packet loss well",
		Port:        9923,
		Protocol:    "udp",
		Binary:      "",
		Deps:        []string{},
		Score:       80,
		Weight:      15,
	},
	{
		Name:        "frp",
		Type:        "frp",
		Category:    CatRelay,
		Description: "Fast reverse proxy, NAT traversal king",
		Port:        7000,
		Protocol:    "tcp",
		Binary:      "frpc",
		Deps:        []string{"frpc"},
		Score:       70,
		Weight:      10,
	},
	{
		Name:        "rathole",
		Type:        "rathole",
		Category:    CatRelay,
		Description: "Rust-based reverse proxy, blistering speed",
		Port:        2333,
		Protocol:    "tcp",
		Binary:      "rathole",
		Deps:        []string{"rathole"},
		Score:       85,
		Weight:      12,
	},
	{
		Name:        "ipsec",
		Type:        "ipsec",
		Category:    CatVPN,
		Description: "Enterprise-grade IPsec VPN (strongSwan)",
		Port:        500,
		Protocol:    "udp",
		Binary:      "ipsec",
		Deps:        []string{"strongswan"},
		Score:       70,
		Weight:      5,
	},
	{
		Name:        "shadowsocks",
		Type:        "shadowsocks",
		Category:    CatProxy,
		Description: "Lightweight secure proxy, protocol obfuscation",
		Port:        8388,
		Protocol:    "tcp",
		Binary:      "ss-server",
		Deps:        []string{"shadowsocks-libev"},
		Score:       55,
		Weight:      3,
	},
	{
		Name:        "hysteria",
		Type:        "hysteria",
		Category:    CatTunnel,
		Description: "Modified QUIC, brute-force throughput",
		Port:        8444,
		Protocol:    "udp",
		Binary:      "hysteria",
		Deps:        []string{"hysteria"},
		Score:       90,
		Weight:      12,
	},
	{
		Name:        "backhaul",
		Type:        "backhaul",
		Category:    CatRelay,
		Description: "Lightning-fast reverse tunnel, NAT traversal with TCP/UDP/WS/WSS",
		Port:        3080,
		Protocol:    "tcp",
		Binary:      "backhaul",
		Deps:        []string{"backhaul"},
		Score:       82,
		Weight:      10,
	},
	{
		Name:        "tcp",
		Type:        "tcp",
		Category:    CatTunnel,
		Description: "Raw TCP direct connection, lowest overhead",
		Port:        9924,
		Protocol:    "tcp",
		Binary:      "",
		Deps:        []string{},
		Score:       50,
		Weight:      3,
	},
}

func LookupTunnel(name string) (*TunnelMeta, error) {
	for _, t := range TunnelRegistry {
		if t.Name == name {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("tunnel %s not found in registry", name)
}

func ListTunnels() []TunnelMeta {
	return TunnelRegistry
}

func CategoryList(category string) []TunnelMeta {
	var result []TunnelMeta
	for _, t := range TunnelRegistry {
		if t.Category == category {
			result = append(result, t)
		}
	}
	return result
}

func InstallScript(name string) string {
	switch name {
	case "wireguard":
		return `apt-get install -y wireguard wireguard-tools`
	case "openvpn":
		return `apt-get install -y openvpn`
	case "ssh":
		return "" // pre-installed
	case "quic":
		return "" // built-in
	case "frp":
		return `TAG=$(curl -sL https://api.github.com/repos/fatedier/frp/releases/latest | grep tag_name | cut -d'"' -f4) && VER=${TAG#v} && curl -sL "https://github.com/fatedier/frp/releases/download/${TAG}/frp_${VER}_linux_amd64.tar.gz" -o /tmp/frp.tar.gz && tar -xzf /tmp/frp.tar.gz -C /tmp/ && cp /tmp/frp_*/frpc /usr/local/bin/frpc && cp /tmp/frp_*/frps /usr/local/bin/frps && chmod +x /usr/local/bin/frpc /usr/local/bin/frps`
	case "rathole":
		return `apt-get install -y unzip >/dev/null 2>&1; TAG=$(curl -sL https://api.github.com/repos/rathole-org/rathole/releases/latest | grep tag_name | cut -d'"' -f4) && curl -sL "https://github.com/rathole-org/rathole/releases/download/${TAG}/rathole-x86_64-unknown-linux-gnu.zip" -o /tmp/rathole.zip && unzip -o /tmp/rathole.zip -d /usr/local/bin/ && chmod +x /usr/local/bin/rathole`
	case "ipsec":
		return `apt-get install -y strongswan`
	case "shadowsocks":
		return `apt-get install -y shadowsocks-libev`
	case "hysteria":
		return `bash -c "$(curl -sL https://get.hy2.sh)"`
	case "backhaul":
		return `curl -sL https://github.com/Musixal/Backhaul/releases/latest/download/backhaul_linux_amd64.tar.gz -o /tmp/backhaul.tar.gz && tar -xzf /tmp/backhaul.tar.gz -C /usr/local/bin/ backhaul && chmod +x /usr/local/bin/backhaul`
	default:
		return ""
	}
}
