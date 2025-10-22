package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func main() {
	const filePath = "enable_profiler.txt"
	var lastValue string
	var processId = -1

	for {
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading file: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		value := strings.TrimSpace(string(data))
		if value != lastValue {
			if value == "true" {
				log.Println("Profiler changed to true!")
				processId = executeProfiler()
			} else {
				log.Printf("Profiler changed to %q", value)
				if processId > 0 {
					stopProfiler(processId)
					processId = -1
				}
			}
			lastValue = value
		}

		time.Sleep(2 * time.Second)
	}
}

func executeProfiler() int {
	cmd := exec.Command("ping", "8.8.8.8")
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
