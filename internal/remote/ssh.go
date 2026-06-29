package remote

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Host struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password,omitempty"`
	KeyPath  string `json:"key_path,omitempty"`

	mu       sync.Mutex
	hostname string
	osInfo   string
	arch     string
	latency  float64
	loss     float64
}

func NewHost(addr string, port int, user, password string) *Host {
	return &Host{
		Address:  addr,
		Port:     port,
		User:     user,
		Password: password,
	}
}

func (h *Host) SSHCommand(cmd string) (string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	sshCmd := exec.Command("sshpass",
		"-p", h.Password,
		"ssh",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-o", "ConnectTimeout=10",
		"-o", "LogLevel=QUIET",
		"-p", fmt.Sprintf("%d", h.Port),
		fmt.Sprintf("%s@%s", h.User, h.Address),
		cmd,
	)

	output, err := sshCmd.Output()
	if err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			stderr := strings.TrimSpace(string(exitErr.Stderr))
			if stderr != "" {
				return "", fmt.Errorf("%s", stderr)
			}
		}
		return "", fmt.Errorf("ssh: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func (h *Host) SCP(localPath, remotePath string) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	scpCmd := exec.Command("sshpass",
		"-p", h.Password,
		"scp",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-P", fmt.Sprintf("%d", h.Port),
		localPath,
		fmt.Sprintf("%s@%s:%s", h.User, h.Address, remotePath),
	)

	output, err := scpCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("scp: %s: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

func (h *Host) Ping(count int) (latency, loss float64) {
	var rtts []float64
	var lossCount float64

	for i := 0; i < count; i++ {
		start := time.Now()
		cmd := exec.Command("ping", "-c", "1", "-W", "2", h.Address)
		if err := cmd.Run(); err == nil {
			rtt := time.Since(start).Seconds() * 1000
			rtts = append(rtts, rtt)
		} else {
			lossCount++
		}
	}

	loss = (lossCount / float64(count)) * 100
	if len(rtts) == 0 {
		return 999, 100
	}

	var sum float64
	for _, r := range rtts {
		sum += r
	}
	latency = sum / float64(len(rtts))

	h.mu.Lock()
	h.latency = latency
	h.loss = loss
	h.mu.Unlock()

	return
}

func (h *Host) DetectOS() error {
	out, err := h.SSHCommand("cat /etc/os-release 2>/dev/null | head -5")
	if err != nil {
		out2, err2 := h.SSHCommand("uname -a")
		if err2 != nil {
			return fmt.Errorf("detect OS: %v / %v", err, err2)
		}
		h.osInfo = out2
		h.arch = "unknown"
		h.hostname = strings.Split(out2, " ")[1]
		return nil
	}

	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			h.osInfo = strings.Trim(strings.SplitN(line, "=", 2)[1], "\"")
		}
	}

	out, _ = h.SSHCommand("uname -m")
	h.arch = strings.TrimSpace(out)

	out, _ = h.SSHCommand("hostname")
	h.hostname = strings.TrimSpace(out)

	log.Printf("[remote] detected: %s | %s | %s", h.hostname, h.osInfo, h.arch)
	return nil
}

func (h *Host) CheckTool(tool string) bool {
	out, err := h.SSHCommand(fmt.Sprintf("which %s 2>/dev/null", tool))
	return err == nil && out != ""
}

func (h *Host) InstallTool(tool string) error {
	log.Printf("[remote] installing %s on %s...", tool, h.Address)

	pkgMap := map[string]string{
		"wireguard": "wireguard wireguard-tools",
		"wg":        "wireguard wireguard-tools",
		"ssh":       "openssh-client",
		"curl":      "curl",
		"wget":      "wget",
		"git":       "git",
		"python3":   "python3",
		"ncat":     "nmap-ncat",
	}

	pkg, ok := pkgMap[tool]
	if !ok {
		pkg = tool
	}

	cmds := []string{
		fmt.Sprintf("apt-get update -qq && apt-get install -y -qq %s 2>&1", pkg),
		fmt.Sprintf("yum install -y -q %s 2>&1", pkg),
		fmt.Sprintf("apk add %s 2>&1", pkg),
	}

	var lastErr error
	for _, cmd := range cmds {
		_, err := h.SSHCommand(cmd)
		if err == nil {
			return nil
		}
		lastErr = err
	}

	return fmt.Errorf("install %s: %v", tool, lastErr)
}

func (h *Host) WriteFile(path, content string, mode string) error {
	modeFlag := "644"
	if mode != "" {
		modeFlag = mode
	}

	tmpFile := fmt.Sprintf("/tmp/nyxora-tmp-%d", time.Now().UnixNano())
	writeCmd := fmt.Sprintf("cat > %s << 'NYXORAEOF'\n%s\nNYXORAEOF\nchmod %s %s",
		tmpFile, content, modeFlag, tmpFile)

	_, err := h.SSHCommand(writeCmd)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	mvCmd := fmt.Sprintf("mv %s %s", tmpFile, path)
	if _, err := h.SSHCommand(mvCmd); err != nil {
		return fmt.Errorf("move file: %w", err)
	}

	return nil
}

func (h *Host) ReadFile(path string) (string, error) {
	return h.SSHCommand(fmt.Sprintf("cat %s", path))
}

func (h *Host) RunDaemon(command string) (string, error) {
	return h.SSHCommand(fmt.Sprintf("nohup %s > /dev/null 2>&1 &", command))
}

func (h *Host) CheckPort(port int, proto string) bool {
	if proto == "" {
		proto = "tcp"
	}
	cmd := fmt.Sprintf("ss -%slnp | grep ':%d '", string(proto[0]), port)
	result, err := h.SSHCommand(cmd)
	return err == nil && result != ""
}

func (h *Host) CheckConnectivity() (string, bool) {
	dialer := net.Dialer{Timeout: 5 * time.Second}
	conn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", h.Address, h.Port))
	if err != nil {
		return fmt.Sprintf("cannot reach %s:%d - %v", h.Address, h.Port, err), false
	}
	conn.Close()

	_, err = h.SSHCommand("echo connected")
	if err != nil {
		return fmt.Sprintf("ssh auth failed: %v", err), false
	}

	return fmt.Sprintf("connected as %s@%s", h.User, h.Address), true
}

func (h *Host) Hostname() string {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.hostname
}

func (h *Host) OSInfo() string {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.osInfo
}

func (h *Host) Arch() string {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.arch
}

func (h *Host) Latency() float64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.latency
}

func (h *Host) Loss() float64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.loss
}

func SSHKeyExists() bool {
	home, _ := os.UserHomeDir()
	paths := []string{
		home + "/.ssh/id_rsa",
		home + "/.ssh/id_ed25519",
		home + "/.ssh/id_ecdsa",
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return true
		}
	}
	return false
}

func SCPFile(local, remote, password string) error {
	host := strings.Split(remote, ":")[0]
	remotePath := strings.Split(remote, ":")[1]
	user := "root"
	port := 22

	if strings.Contains(host, "@") {
		parts := strings.Split(host, "@")
		user = parts[0]
		host = parts[1]
	}
	if strings.Contains(host, ":") {
		parts := strings.Split(host, ":")
		host = parts[0]
		fmt.Sscanf(parts[1], "%d", &port)
	}

	cmd := exec.Command("sshpass",
		"-p", password,
		"scp",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-P", fmt.Sprintf("%d", port),
		local,
		fmt.Sprintf("%s@%s:%s", user, host, remotePath),
	)
	return cmd.Run()
}
