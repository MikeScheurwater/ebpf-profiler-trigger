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

func main() {
	const filePath = "config.json"
	var lastValue bool
	var processId = -1
	var pollInterval = 2 * time.Second

	for {
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading file: %v", err)
			time.Sleep(pollInterval)
			continue
		}

		var config EbpfConfig
		err = json.Unmarshal(data, &config)
		if err != nil {
			log.Printf("Error parsing json: %v", err)
			time.Sleep(pollInterval)
			continue
		}

		if config.Enabled != lastValue {
			log.Printf("Profiler changed to %v", config.Enabled)
			if config.Enabled {
				processId = executeProfiler(config.Command, config.Args...)
			} else {

				if processId > 0 {
					stopProfiler(processId)
					processId = -1
				}
			}
			lastValue = config.Enabled
		}

		pollInterval = config.PollInterval * time.Second
		time.Sleep(pollInterval)
	}
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
