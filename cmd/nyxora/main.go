package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/nyxora/nyxora/internal/config"
	"github.com/nyxora/nyxora/internal/orchestrator"
)

var version = "0.2.0"

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]
	switch command {
	case "install":
		cmdInstall()
	case "connect":
		cmdConnect()
	case "disconnect":
		cmdDisconnect()
	case "status":
		cmdStatus()
	case "dashboard":
		cmdDashboard()
	case "daemon":
		cmdDaemon()
	case "tui":
		cmdTUI()
	case "update":
		cmdUpdate()
	case "server":
		cmdServer()
	case "version":
		fmt.Printf("nyxora v%s\n", version)
	case "help":
		printUsage()
	default:
		fmt.Printf("unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`NYXORA  -  Adaptive Tunnel Orchestrator

USAGE
  nyxora install                    Check dependencies & setup directories
  nyxora connect <host>             Connect to remote server
  nyxora disconnect                 Close all tunnels
  nyxora status                     Show connection info
  nyxora dashboard                  Live terminal dashboard
  nyxora tui                        Interactive terminal UI
  nyxora update                     Check for updates
  nyxora daemon                     Run as background service
  nyxora server                     Show server info & suggest mode
  nyxora version                    Show version
  nyxora help                       This page

TELEGRAM
  https://t.me/NyxoraCore

CONNECT OPTIONS
  --user, -u        SSH user (default: root)
  --port, -p        SSH port (default: 22)
  --password        SSH password
  --mode            Server mode: full, lite, minimal (auto-detect if omitted)
  --transports      Comma-separated transport list (overrides mode)
  --ports           Port overrides: wg=51820,ss=8388,...

MODES
  full      All 11 tunnels (2GB+ RAM)
  lite      Lightweight tunnels only (512MB-2GB RAM)
  minimal   SSH + Shadowsocks only (<512MB RAM)

EXAMPLES
  nyxora connect 1.2.3.4 --user root --password secret
  nyxora connect 1.2.3.4 --mode lite
  nyxora connect 1.2.3.4 --transports ssh,shadowsocks,quic
  nyxora connect 1.2.3.4 --ports wg=51820,ss=9000,hy=8444
  nyxora server
  nyxora dashboard
`)
}

func loadConfig() *config.Config {
	cfgPath := os.Getenv("NYXORA_CONFIG")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		cfg = &config.DefaultConfig
	}
	if v := os.Getenv("NYXORA_ALL_ACTIVE"); v == "true" || v == "1" {
		cfg.AllTunnelsActive = true
	}
	return cfg
}

func parseFlags(cfg *config.Config) {
	for i, arg := range os.Args {
		switch arg {
		case "--mode":
			if i+1 < len(os.Args) {
				mode := config.ServerMode(os.Args[i+1])
				switch mode {
				case config.ModeFull, config.ModeLite, config.ModeMinimal:
					cfg.Mode = mode
				default:
					log.Fatalf("invalid mode: %s (use full, lite, or minimal)", os.Args[i+1])
				}
			}
		case "--transports":
			if i+1 < len(os.Args) {
				cfg.EnabledTransports = strings.Split(os.Args[i+1], ",")
			}
		case "--ports":
			if i+1 < len(os.Args) {
				cfg.PortOverrides = make(map[string]int)
				for _, pair := range strings.Split(os.Args[i+1], ",") {
					parts := strings.SplitN(pair, "=", 2)
					if len(parts) == 2 {
						port, err := strconv.Atoi(parts[1])
						if err == nil {
							cfg.PortOverrides[parts[0]] = port
						}
					}
				}
			}
		}
	}
}

func cmdInstall() {
	fmt.Println()
	fmt.Println("  " + "\033[38;5;141m● NYXORA\033[0m  \033[90mInstallation\033[0m")
	fmt.Println("  " + "\033[90m" + strings.Repeat("━", 40) + "\033[0m")
	fmt.Println()

	for _, dir := range []string{"/etc/nyxora", "/etc/nyxora/tunnels", "/etc/nyxora/cache", "/var/log/nyxora"} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("create dir %s: %v", dir, err)
		}
		fmt.Printf("  \033[32m✓\033[0m %s\n", dir)
	}

	cfg := loadConfig()
	parseFlags(cfg)
	cfg.Save("/etc/nyxora/config.json")

	fmt.Println()
	fmt.Println("  \033[1mChecking dependencies:\033[0m")
	for _, dep := range []string{"ping", "wg", "ssh", "sshpass", "curl"} {
		path, err := exec.LookPath(dep)
		if err == nil {
			fmt.Printf("  \033[32m✓\033[0m %-12s %s\n", dep, path)
		} else {
			fmt.Printf("  \033[33m△\033[0m %-12s \033[90mnot found (install: apt install %s)\033[0m\n", dep, dep)
		}
	}

	fmt.Println()
	info := config.ServerInfo()
	mode := info["suggested_mode"].(config.ServerMode)
	transports := config.GetTransportsForMode(mode)
	fmt.Printf("  \033[1mSuggested mode:\033[0m  %s\n", mode)
	fmt.Printf("  \033[1mTransports:\033[0m      %s\n", strings.Join(transports, ", "))
	if ram, ok := info["total_ram_mb"].(uint64); ok {
		fmt.Printf("  \033[1mTotal RAM:\033[0m       %d MB\n", ram)
	}

	fmt.Println()
	fmt.Println("  \033[32m● NYXORA installed successfully\033[0m")
	fmt.Println()
	fmt.Println("  \033[90mNext step:\033[0m")
	fmt.Println("  \033[1m  nyxora connect <remote-ip> --mode <full|lite|minimal>\033[0m")
	fmt.Println()
}

func cmdConnect() {
	addr := ""
	user := "root"
	port := 22
	password := ""

	if len(os.Args) >= 3 {
		addr = os.Args[2]
	}

	for i, arg := range os.Args {
		switch arg {
		case "--user", "-u":
			if i+1 < len(os.Args) {
				user = os.Args[i+1]
			}
		case "--port", "-p":
			if i+1 < len(os.Args) {
				p, err := strconv.Atoi(os.Args[i+1])
				if err == nil {
					port = p
				}
			}
		case "--password", "--pass":
			if i+1 < len(os.Args) {
				password = os.Args[i+1]
			}
		}
	}

	if addr == "" {
		fmt.Print("  Enter remote server IP: ")
		fmt.Scanln(&addr)
	}
	if addr == "" {
		log.Fatalf("remote address is required")
	}

	cfg := loadConfig()
	parseFlags(cfg)

	// Validate config BEFORE asking for password
	if err := cfg.Validate(); err != nil {
		log.Fatalf("\033[31m✗ config error: %v\033[0m", err)
	}

	if password == "" {
		fmt.Print("  Enter SSH password: ")
		fmt.Scanln(&password)
	}
	if password == "" {
		log.Fatalf("password is required")
	}

	transports := cfg.GetEffectiveTransports()
	fmt.Println()
	fmt.Printf("  \033[36m● Connecting to %s@%s:%d ...\033[0m\n", user, addr, port)
	fmt.Printf("  \033[90m  Mode: %s | Transports: %s\033[0m\n", cfg.Mode, strings.Join(transports, ", "))
	if cfg.PortOverrides != nil {
		fmt.Printf("  \033[90m  Port overrides: %v\033[0m\n", cfg.PortOverrides)
	}
	fmt.Println()

	o := orchestrator.New(cfg)
	if err := o.Init(); err != nil {
		log.Fatalf("\033[31m✗ init: %v\033[0m", err)
	}

	o.OnStepUpdate(func(step orchestrator.StepStatus) {
		icon := "\033[90m○\033[0m"
		switch step.Status {
		case "OK":
			icon = "\033[32m✓\033[0m"
		case "FAILED":
			icon = "\033[31m✗\033[0m"
		case "RUNNING":
			icon = "\033[33m◉\033[0m"
		case "WARN":
			icon = "\033[38;5;214m△\033[0m"
		}
		detail := ""
		if step.Detail != "" && step.Done {
			detail = " \033[90m" + step.Detail + "\033[0m"
		}
		fmt.Printf("  %s %s%s\n", icon, step.Name, detail)
	})

	fmt.Println()
	if err := o.ConnectToRemote(addr, port, user, password); err != nil {
		fmt.Printf("\n  \033[31m✗ Connection failed: %v\033[0m\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("  \033[32m● Tunnel established successfully\033[0m")
	fmt.Println("  \033[90m  Run 'nyxora dashboard' for live monitoring\033[0m")
	fmt.Println()
}

func cmdDisconnect() {
	fmt.Println("  \033[33m◉ disconnect ...\033[0m")
	cfg := loadConfig()
	o := orchestrator.New(cfg)
	o.Init()
	o.Stop()
	fmt.Println("  \033[32m✓ disconnected\033[0m")
}

func cmdStatus() {
	cfg := loadConfig()
	o := orchestrator.New(cfg)
	o.Init()

	status := o.Status()
	fmt.Println()
	fmt.Println("  \033[38;5;141m● NYXORA STATUS\033[0m")
	fmt.Println("  \033[90m" + strings.Repeat("━", 35) + "\033[0m")
	fmt.Printf("  \033[1mStatus:\033[0m        ")
	if c, _ := status["connected"].(bool); c {
		fmt.Println("\033[32mconnected\033[0m")
	} else {
		fmt.Println("\033[33midle\033[0m")
	}
	fmt.Printf("  \033[1mActive Tunnel:\033[0m %s\n", status["active_transport"])
	if r, ok := status["remote"].(map[string]interface{}); ok {
		fmt.Printf("  \033[1mRemote Host:\033[0m   %s (%s)\n", r["hostname"], r["address"])
	}
	if p, ok := status["phase"].(string); ok {
		fmt.Printf("  \033[1mPhase:\033[0m         %s\n", p)
	}
	fmt.Println()
}

func cmdDashboard() {
	fmt.Print("\033[?25l\033[2J")
	defer fmt.Print("\033[?25h")

	cfg := loadConfig()
	o := orchestrator.New(cfg)
	if err := o.Init(); err != nil {
		log.Fatalf("init: %v", err)
	}
	if err := o.Start(); err != nil {
		log.Fatalf("dashboard: %v", err)
	}
}

func cmdDaemon() {
	cfg := loadConfig()
	log.Printf("starting nyxora daemon v%s", version)
	o := orchestrator.New(cfg)
	if err := o.Init(); err != nil {
		log.Fatalf("init: %v", err)
	}
	if err := o.Start(); err != nil {
		log.Fatalf("daemon: %v", err)
	}
}

func cmdTUI() {
	fmt.Println("\n  Launching NYXORA TUI...")
	// TUI is launched via bubbletea - placeholder for interactive mode
	// Use 'nyxora dashboard' for live monitoring
}

func cmdUpdate() {
	// Check for updates via GitHub API
	fmt.Println("\n  Checking for updates...")
	resp, err := http.Get("https://api.github.com/repos/nyxorammd-lgtm/nyxora/releases/latest")
	if err != nil {
		fmt.Printf("  Error checking updates: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}
	json.NewDecoder(resp.Body).Decode(&release)

	latest := strings.TrimPrefix(release.TagName, "v")
	if latest == "0.2.0" || latest == "" {
		fmt.Println("  \033[32m✓\033[0m You're up to date! (v0.2.0)")
	} else {
		fmt.Printf("  \033[33m△\033[0m New version available: v%s\n", latest)
		fmt.Printf("  Download: https://github.com/nyxorammd-lgtm/nyxora/releases/tag/v%s\n", latest)
	}
	fmt.Println()
}

func cmdServer() {
	info := config.ServerInfo()
	mode := info["suggested_mode"].(config.ServerMode)
	transports := config.GetTransportsForMode(mode)

	fmt.Println()
	fmt.Println("  \033[38;5;141m● SERVER INFO\033[0m")
	fmt.Println("  \033[90m" + strings.Repeat("━", 35) + "\033[0m")

	if ram, ok := info["total_ram_mb"].(uint64); ok {
		fmt.Printf("  \033[1mTotal RAM:\033[0m       %d MB\n", ram)
	}
	if avail, ok := info["available_ram_mb"].(uint64); ok {
		fmt.Printf("  \033[1mAvailable RAM:\033[0m   %d MB\n", avail)
	}
	fmt.Printf("  \033[1mCPU Cores:\033[0m       %d\n", info["cpu_count"])

	fmt.Println()
	fmt.Printf("  \033[1mSuggested mode:\033[0m  \033[32m%s\033[0m\n", mode)
	fmt.Printf("  \033[1mTransports:\033[0m      %s\n", strings.Join(transports, ", "))

	fmt.Println()
	fmt.Println("  \033[90mUse --mode to override:\033[0m")
	fmt.Println("  \033[90m  full    - all 11 tunnels (2GB+ RAM)\033[0m")
	fmt.Println("  \033[90m  lite    - lightweight only (512MB-2GB)\033[0m")
	fmt.Println("  \033[90m  minimal - SSH + SS only (<512MB)\033[0m")
	fmt.Println()
}
