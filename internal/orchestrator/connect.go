package orchestrator

import (
	"fmt"
	"log"
	"strings"

	"github.com/nyxora/nyxora/internal/config"
	"github.com/nyxora/nyxora/internal/remote"
	"github.com/nyxora/nyxora/internal/transport"
)

func (o *Orchestrator) ConnectToRemote(addr string, port int, user, password string) error {
	o.phase = PhaseConnecting
	log.Printf("[orchestrator] connecting to %s@%s:%d", user, addr, port)

	o.remoteHost = remote.NewHost(addr, port, user, password)

	o.addStep("Pinging remote server", "RUNNING", "")
	lat, loss := o.remoteHost.Ping(4)
	if loss > 80 {
		o.addStep("Pinging remote server", "FAILED", fmt.Sprintf("loss: %.0f%%", loss))
		o.phase = PhaseFailed
		return fmt.Errorf("remote unreachable: %.0f%% loss", loss)
	}
	o.addStep("Pinging remote server", "OK", fmt.Sprintf("%.0fms, %.0f%% loss", lat, loss))

	o.addStep("SSH authentication", "RUNNING", "")
	msg, ok := o.remoteHost.CheckConnectivity()
	if !ok {
		o.addStep("SSH authentication", "FAILED", msg)
		o.phase = PhaseFailed
		return fmt.Errorf("ssh: %s", msg)
	}
	o.addStep("SSH authentication", "OK", msg)

	o.phase = PhaseSetup
	o.addStep("Detecting OS", "RUNNING", "")
	if err := o.remoteHost.DetectOS(); err != nil {
		o.addStep("Detecting OS", "FAILED", err.Error())
		o.phase = PhaseFailed
		return err
	}
	o.addStep("Detecting OS", "OK", fmt.Sprintf("%s | %s", o.remoteHost.OSInfo(), o.remoteHost.Arch()))

	// Check remote resources against mode
	minRAM := o.getMinRAMForMode()
	if minRAM > 0 {
		o.addStep("Checking remote resources", "RUNNING", "")
		ok, msg := o.remoteHost.CheckResourcesForMode(string(o.cfg.Mode), minRAM)
		if !ok {
			o.addStep("Checking remote resources", "FAILED", msg)
			o.phase = PhaseFailed
			return fmt.Errorf("resource check failed: %s", msg)
		}
		o.addStep("Checking remote resources", "OK", msg)
	}

	o.addStep("Installing dependencies", "RUNNING", "")
	var failedDeps []string
	for _, meta := range transport.ListTunnels() {
		if meta.Binary == "" {
			continue
		}
		if o.remoteHost.CheckTool(meta.Binary) {
			continue
		}
		script := transport.InstallScript(meta.Name)
		if script == "" {
			continue
		}
		log.Printf("[orchestrator] installing %s on remote...", meta.Name)
		_, err := o.remoteHost.SSHCommand(script)
		if err != nil {
			failedDeps = append(failedDeps, meta.Name)
			log.Printf("[orchestrator] %s install failed: %v", meta.Name, err)
		}
	}
	if len(failedDeps) > 0 {
		o.addStep("Installing dependencies", "WARN", fmt.Sprintf("%d failed: %s", len(failedDeps), strings.Join(failedDeps, ", ")))
	} else {
		o.addStep("Installing dependencies", "OK", "all tunnel dependencies ready")
	}

	o.phase = PhaseTunnel
	o.addStep("Generating WireGuard keys", "RUNNING", "")
	localPriv, localPub := generateLocalWGKey()
	o.addStep("Generating WireGuard keys", "OK", fmt.Sprintf("pub: %s...", localPub[:16]))

	o.addStep("Setting up remote WG endpoint", "RUNNING", "")
	remotePub, remoteIface, err := remote.SetupWireGuardRemote(o.remoteHost, localPub, 51820)
	o.remoteIface = remoteIface
	if err != nil {
		o.addStep("Setting up remote WG endpoint", "FAILED", err.Error())
		o.phase = PhaseFailed
		return err
	}
	o.addStep("Setting up remote WG endpoint", "OK", fmt.Sprintf("pub: %s... | iface: %s", remotePub[:16], o.remoteIface))

	o.addStep("Setting up local WG endpoint", "RUNNING", "")
	remoteIP, err := remote.GetRemotePublicIP(o.remoteHost)
	if err != nil {
		remoteIP = addr
	}

	localWG, ok := o.transportM.Get("wireguard")
	if !ok {
		localWG = transport.NewWireGuard()
		o.transportM.Register(localWG)
	}
	wgPort := 51820
	subnet := wgPort % 256
	localWG.Init(map[string]string{
		"private_key": localPriv,
		"remote_pub":  remotePub,
		"interface":   "nyxora0",
		"local_addr":  fmt.Sprintf("10.100.%d.2/24", subnet),
	})
	if err := localWG.Connect(remoteIP); err != nil {
		o.addStep("Setting up local WG endpoint", "FAILED", err.Error())
		o.phase = PhaseFailed
		return err
	}
	o.addStep("Setting up local WG endpoint", "OK", "interface nyxora0 ready")

	o.connected = true
	o.phase = PhaseMultipath

	o.initTransportSecrets(user, password, port, addr)

	o.addStep("Provisioning remote tunnel services", "RUNNING", "")
	if err := o.provisionRemoteTransports(); err != nil {
		o.addStep("Provisioning remote tunnel services", "WARN", fmt.Sprintf("some provisioning failed: %v", err))
	} else {
		o.addStep("Provisioning remote tunnel services", "OK", "remote services ready")
	}

	if o.cfg.AllTunnelsActive {
		o.addStep("Multipath mode", "OK", "all tunnels active simultaneously")
		o.transportM.SetAllActive(true)
		o.transportM.ConnectAll(remoteIP)
	} else {
		o.addStep("Smart mode", "OK", "best tunnel selected automatically")
	}

	o.addStep("Tunnel established", "OK",
		fmt.Sprintf("%s <-> %s (%s)", o.localNodeID[:8], o.remoteHost.Hostname(), remoteIP))

	log.Printf("[orchestrator] tunnel active: %s <-> %s (%s)",
		o.localNodeID[:8], o.remoteHost.Hostname(), remoteIP)

	o.routeEngine.SetCurrent("wireguard")
	go o.startMonitoring(remoteIP)

	return nil
}

func (o *Orchestrator) getMinRAMForMode() uint64 {
	t := config.DefaultThresholds
	if o.cfg.Thresholds != nil {
		t = *o.cfg.Thresholds
	}
	switch o.cfg.Mode {
	case config.ModeMinimal:
		return 0
	case config.ModeLite:
		return t.MinimalMaxMB
	default:
		return t.LiteMaxMB
	}
}

func (o *Orchestrator) initTransportSecrets(user, password string, port int, addr string) {
	if t, ok := o.transportM.Get("ssh"); ok {
		t.Init(map[string]string{"password": password, "user": user, "port": fmt.Sprintf("%d", port)})
	}
	if t, ok := o.transportM.Get("shadowsocks"); ok {
		ssPort := o.cfg.GetPort("shadowsocks", 8388)
		t.Init(map[string]string{"password": o.cfg.Secrets.SSPassword, "method": o.cfg.Secrets.SSMethod, "port": fmt.Sprintf("%d", ssPort)})
	}
	if t, ok := o.transportM.Get("rathole"); ok {
		ratholePort := o.cfg.GetPort("rathole", 2333)
		t.Init(map[string]string{"token": o.cfg.Secrets.RatholeToken, "port": fmt.Sprintf("%d", ratholePort)})
	}
	if t, ok := o.transportM.Get("hysteria"); ok {
		hyPort := o.cfg.GetPort("hysteria", 8444)
		t.Init(map[string]string{"auth": o.cfg.Secrets.HysteriaAuth, "port": fmt.Sprintf("%d", hyPort)})
	}
	if t, ok := o.transportM.Get("backhaul"); ok {
		bhPort := o.cfg.GetPort("backhaul", 3080)
		t.Init(map[string]string{"token": o.cfg.Secrets.BackhaulToken, "port": fmt.Sprintf("%d", bhPort), "backhaul_transport": "tcp"})
	}
	if t, ok := o.transportM.Get("ipsec"); ok {
		t.Init(map[string]string{"local_ip": getLocalIP(), "remote_ip": addr, "psk": o.cfg.Secrets.IPsecPSK})
	}
}

func (o *Orchestrator) provisionRemoteTransports() error {
	if o.remoteHost == nil {
		return fmt.Errorf("no remote host")
	}
	var errs []string

	frpsPort := o.cfg.GetPort("frp", 7000)
	if err := remote.ProvisionFRPServer(o.remoteHost, frpsPort); err != nil {
		errs = append(errs, "frp:"+err.Error())
	}
	ratholePort := o.cfg.GetPort("rathole", 2333)
	if err := remote.ProvisionRatholeServer(o.remoteHost, ratholePort, o.cfg.Secrets.RatholeToken); err != nil {
		errs = append(errs, "rathole:"+err.Error())
	}
	ssPort := o.cfg.GetPort("shadowsocks", 8388)
	if err := remote.ProvisionShadowSOCKSServer(o.remoteHost, ssPort, o.cfg.Secrets.SSPassword, o.cfg.Secrets.SSMethod); err != nil {
		errs = append(errs, "shadowsocks:"+err.Error())
	}
	hyPort := o.cfg.GetPort("hysteria", 8444)
	if err := remote.ProvisionHysteriaServer(o.remoteHost, hyPort, o.cfg.Secrets.HysteriaAuth); err != nil {
		errs = append(errs, "hysteria:"+err.Error())
	}
	bhPort := o.cfg.GetPort("backhaul", 3080)
	if err := remote.ProvisionBackhaulServer(o.remoteHost, bhPort, o.cfg.Secrets.BackhaulToken, "tcp"); err != nil {
		errs = append(errs, "backhaul:"+err.Error())
	}
	ovpnPort := o.cfg.GetPort("openvpn", 1194)
	if err := remote.ProvisionOpenVPNServer(o.remoteHost, ovpnPort); err != nil {
		errs = append(errs, "openvpn:"+err.Error())
	}
	if err := remote.ProvisionIPsecServer(o.remoteHost, getLocalIP(), o.cfg.Secrets.IPsecPSK); err != nil {
		errs = append(errs, "ipsec:"+err.Error())
	}

	if len(errs) > 0 {
		return fmt.Errorf("provision errors: %s", strings.Join(errs, "; "))
	}
	return nil
}
