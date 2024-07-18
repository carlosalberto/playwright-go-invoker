package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

var (
	logger = log.Default()
)

func main() {
	dir := os.Getenv("PLAYWRIGHT_DIR")
	if dir == "" {
		panic("Need to specify PLAYWRIGHT_DIR")
	}
	testdir := os.Getenv("TEST_DIR") // Optional, defaults to PLAYWRIGHT_DIR

	// Time period between full invocations.
	ticker := time.NewTicker(5 * time.Second)
	done := make(chan bool)

	go func() {
		InvokePlaywright(dir, testdir)

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				InvokePlaywright(dir, testdir)
			}
		}
	}()

	fmt.Printf("Press enter to stop reporting...")
	var input string
	fmt.Scanln(&input)
	fmt.Printf("Shutting down...\n")
	ticker.Stop()
	done <- true
}

func InvokePlaywright(dir string, testdir string) {
	args := []string{"playwright", "test", "--reporter=json"}
	if testdir != "" {
		args = append(args, testdir)
	}

	cmd := exec.Command("npx", args...)
	cmd.Dir = dir

	bytes, err := cmd.Output()
	if err != nil {
		if bytes == nil {
			logger.Println(fmt.Sprintf("Could not fetch output: %v", err))
			return
		}
	}

	res := &ReportResult{}
	err = json.Unmarshal(bytes, res)
	if err != nil {
		logger.Println(fmt.Sprintf("Could not unmarshal output: %v", err))
		return
	}

	logger.Println(fmt.Sprintf("\nTotal Errors: %d\n", res.Stats.Unexpected))
	for _, suite := range res.Suites {
		for _, spec := range suite.Specs {
			logger.Println(fmt.Sprintf("%s: %t\n", spec.Title, spec.Ok))
		}
	}
}
