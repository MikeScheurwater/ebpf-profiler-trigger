package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
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
		cancelFn     context.CancelFunc
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
				ctx, cancel := context.WithCancel(context.Background())
				cancelFn = cancel
				go executeProfiler(ctx, config.Command, config.Args...)
			} else if cancelFn != nil {
				cancelFn()
				cancelFn = nil
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

func executeProfiler(ctx context.Context, command string, args ...string) {
	log.Printf("Executing: %v %v\n", command, args)
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Printf("Failed to start command: %v", err)
		return
	}

	log.Printf("Started process %d", cmd.Process.Pid)
	err := cmd.Wait()
	if err != nil {
		// Check if command was canceled by context
		if errors.Is(ctx.Err(), context.Canceled) {
			log.Printf("Process %d stopped (context canceled)", cmd.Process.Pid)
		} else {
			log.Printf("Process %d exited with error: %v", cmd.Process.Pid, err)
		}
	} else {
		log.Printf("Process %d finished successfully", cmd.Process.Pid)
	}
}
