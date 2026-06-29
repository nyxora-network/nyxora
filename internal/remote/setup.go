package remote

import (
	"fmt"
	"log"
	"strings"
)

type SetupStep struct {
	Name   string
	Status string
	Detail string
	Done   bool
}

type SetupResult struct {
	Hostname  string
	OS        string
	Arch      string
	LatencyMs float64
	Steps     []SetupStep
	Success   bool
}

func SetupHost(host *Host) (*SetupResult, error) {
	result := &SetupResult{Success: true}

	// Step 1: connectivity check
	step1 := SetupStep{Name: "Checking connectivity"}
	msg, ok := host.CheckConnectivity()
	if !ok {
		step1.Status = "FAILED"
		step1.Detail = msg
		result.Steps = append(result.Steps, step1)
		result.Success = false
		return result, fmt.Errorf("connectivity: %s", msg)
	}
	step1.Status = "OK"
	step1.Detail = msg
	step1.Done = true
	result.Steps = append(result.Steps, step1)

	// Step 2: detect OS
	step2 := SetupStep{Name: "Detecting operating system"}
	if err := host.DetectOS(); err != nil {
		step2.Status = "FAILED"
		step2.Detail = err.Error()
		result.Success = false
		result.Steps = append(result.Steps, step2)
		return result, err
	}
	step2.Status = "OK"
	step2.Detail = fmt.Sprintf("%s | %s", host.OSInfo(), host.Arch())
	step2.Done = true
	result.Steps = append(result.Steps, step2)
	result.Hostname = host.Hostname()
	result.OS = host.OSInfo()
	result.Arch = host.Arch()

	// Step 3: ping test
	step3 := SetupStep{Name: "Measuring network quality"}
	lat, loss := host.Ping(4)
	step3.Status = "OK"
	step3.Detail = fmt.Sprintf("latency: %.0fms | loss: %.0f%%", lat, loss)
	step3.Done = true
	result.Steps = append(result.Steps, step3)
	result.LatencyMs = lat

	// Step 4: check and install dependencies
	step4 := SetupStep{Name: "Checking dependencies"}
	var deps []string
	for _, dep := range []string{"wg", "ssh", "curl", "ncat"} {
		if host.CheckTool(dep) {
			deps = append(deps, dep+": OK")
		} else {
			depName := dep
			if dep == "wg" {
				depName = "wireguard"
			}
			if dep == "ncat" {
				depName = "ncat"
			}
			log.Printf("[setup] installing %s on remote...", depName)
			if err := host.InstallTool(depName); err != nil {
				deps = append(deps, dep+": FAILED")
				log.Printf("[setup] install %s failed: %v", depName, err)
			} else {
				deps = append(deps, dep+": installed")
			}
		}
	}
	step4.Status = "OK"
	step4.Detail = strings.Join(deps, ", ")
	step4.Done = true
	result.Steps = append(result.Steps, step4)

	log.Printf("[setup] host %s ready | %s | %s | %.0fms",
		host.Hostname(), host.OSInfo(), host.Arch(), result.LatencyMs)

	return result, nil
}
