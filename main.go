package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type EbpfConfig struct {
	Enabled      bool          `json:"enabled"`
	Command      string        `json:"command"`
	Args         []string      `json:"args"`
	PollInterval time.Duration `json:"poll_interval"`
}

const (
	filePath            = "config.json"
	defaultPollInterval = 2
)

func main() {
	var (
		lastEnabled  bool
		processId    = -1
		pollInterval = defaultPollInterval * time.Second
	)

	for {
		config, err := loadConfig(filePath)
		if err != nil {
			log.Printf("Error loading config: %v", err)
			time.Sleep(pollInterval)
			continue
		}

		if config.Enabled != lastEnabled && config.Command != "" {
			log.Printf("Profiler changed to %v", config.Enabled)
			if config.Enabled {
				processId = executeProfiler(config.Command, config.Args...)
			} else {

				if processId > 0 {
					stopProfiler(processId)
					processId = -1
				}
			}
			lastEnabled = config.Enabled
		}

		if config.PollInterval > 0 {
			pollInterval = config.PollInterval * time.Second
		}

		time.Sleep(pollInterval)
	}
}

func loadConfig(path string) (*EbpfConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg EbpfConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.PollInterval == 0 {
		cfg.PollInterval = defaultPollInterval
	}

	return &cfg, nil
}

func executeProfiler(command string, args ...string) int {
	fmt.Printf("Executing: %v %v\n", command, args)
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Starting process %d\n", cmd.Process.Pid)

	return cmd.Process.Pid
}

func stopProfiler(processId int) {
	log.Printf("Killing process %d\n", processId)
	err := syscall.Kill(processId, syscall.SIGTERM)
	if err != nil {
		log.Printf("Error killing process %d: %v", processId, err)
	}
}
