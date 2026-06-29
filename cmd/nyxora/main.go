package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/nyxora/nyxora/internal/config"
	"github.com/nyxora/nyxora/internal/orchestrator"
)

var version = "0.1.0"

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
	fmt.Println(`NYXORA  -  Adaptive Tunnel Orchestrator

USAGE
  nyxora install                    Check dependencies & setup directories
  nyxora connect <host>             Connect to remote server (interactive)
  nyxora disconnect                 Close all tunnels
  nyxora status                     Show connection info
  nyxora dashboard                  Live terminal dashboard
  nyxora daemon                     Run as background service
  nyxora version                    Show version
  nyxora help                       This page

EXAMPLES
  nyxora connect 213.32.69.147
  nyxora connect 213.32.69.147 --user root --port 22
  nyxora dashboard

WHAT IS NYXORA?
  Install on ONE server. It connects to your remote server via SSH,
  installs what's needed, sets up the fastest tunnel, and keeps it
  healthy — automatically. No agent needed on the other side.
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
	fmt.Println("  \033[32m● NYXORA installed successfully\033[0m")
	fmt.Println()
	fmt.Println("  \033[90mNext step:\033[0m")
	fmt.Println("  \033[1m  nyxora connect <remote-ip>\033[0m")
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

	if password == "" {
		fmt.Print("  Enter SSH password: ")
		fmt.Scanln(&password)
	}
	if password == "" {
		log.Fatalf("password is required")
	}

	fmt.Println()
	fmt.Printf("  \033[36m● Connecting to %s@%s:%d ...\033[0m\n", user, addr, port)
	fmt.Println()

	cfg := loadConfig()
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
